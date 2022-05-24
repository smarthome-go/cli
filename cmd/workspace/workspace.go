package workspace

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/fatih/color"
	"github.com/pelletier/go-toml"
	"github.com/rodaine/table"
	"github.com/sergi/go-diff/diffmatchpatch"

	"github.com/smarthome-go/sdk"
)

type ConfigToml struct {
	Id                  string `toml:"id"`
	Name                string `toml:"name"`
	Description         string `toml:"description"`
	QuickActionsEnabled bool   `toml:"quickActions"`
	SchedulerEnabled    bool   `toml:"scheduler"`
	MDIcon              string `toml:"icon"`
}

// Creates a new project on the remote and locally
func New(id string, name string, c *sdk.Connection) {
	if err := createProjectFiles(id, name); err != nil {
		if os.IsExist(err) {
			fmt.Printf("Failed to initialize project root at `./%s`: specified directory already exists.\n", id)
		} else {
			fmt.Printf("Failed to initialize project root at `./%s`: %s\n", id, err.Error())
		}
		os.Exit(1)
	}
	if err := c.CreateHomescript(sdk.HomescriptRequest{
		Id:   id,
		Name: name,
	}); err != nil {
		switch err {
		case sdk.ErrUnprocessableEntity:
			fmt.Printf("Failed to create remote project: id (`%s`) already exists on remote.\n", id)
		case sdk.ErrPermissionDenied:
			fmt.Printf("Failed to create remote project: permission denied: please ensure that you have the correct access rights to create new hms-objects.\n")
		default:
			fmt.Printf("Failed to create project `%s`: could not create remote object: unknown error: %s\n", id, err.Error())
		}
		if err := removeProjectFiles(id); err != nil {
			fmt.Printf("Revert: project root at `./%s` could not be removed.\n", id)
		}
		os.Exit(1)
	}
	fmt.Printf("Successfully created new remote project: '%s' at './%s'.\n", id, id)
}

// Removes a local project
// If `purgeOrigin` is set to `true`, the project is also deleted on the remote
func Delete(id string, purgeOrigin bool, c *sdk.Connection) {
	if err := removeProjectFiles(id); err != nil {
		fmt.Printf("Failed to remove local project files: %s\n", err.Error())
		os.Exit(1)
	}
	if purgeOrigin {
		if err := c.DeleteHomescript(id); err != nil {
			switch err {
			case sdk.ErrUnprocessableEntity:
				fmt.Printf("Failed to remove project: local project (`%s`) does not exist on remote.\n", id)
			case sdk.ErrPermissionDenied:
				fmt.Printf("Failed to remove project: permission denied: please ensure that you have the correct access rights to remove hms-objects.\n")
			case sdk.ErrConflict:
				fmt.Printf("Failed to remove project: safety prevention: one or more automations depend on this homescript.\n")
			default:
				fmt.Printf("Failed to remove project: unknown error: %s\n", err.Error())
			}
			os.Exit(1)
		}
		fmt.Printf("Deleted project `%s` from remote.\n", id)
	}
}

// Used by `createProject`, creates the `hms.toml` config file
func createProjectConfigFile(id string, name string, path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	if name == "" {
		name = id
	}
	if err := toml.NewEncoder(file).Encode(ConfigToml{
		Id:     id,
		Name:   name,
		MDIcon: "code",
	}); err != nil {
		return err
	}
	if err := file.Close(); err != nil {
		return err
	}
	return nil
}

// Creates all needed project files
func createProjectFiles(id string, name string) error {
	if err := os.Mkdir(id, 0755); err != nil {
		return err
	}
	if err := ioutil.WriteFile(
		fmt.Sprintf("./%s/%s.hms", id, id),
		[]byte(fmt.Sprintf("# Write your code for `%s` below", id)),
		0775,
	); err != nil {
		return err
	}
	return createProjectConfigFile(id, name, fmt.Sprintf("./%s/hms.toml", id))
}

// Deletes the project from the local file system
// Used by `Delete`
func removeProjectFiles(id string) error {
	if _, err := os.Stat(id); err != nil {
		if os.IsNotExist(err) {
			fmt.Printf("Project does not exist locally, therefore skipping local removal.\n")
			return nil
		}
		return err
	}
	if err := os.RemoveAll(id); err != nil {
		return err
	}
	fmt.Printf("Removed project root at ./%s\n", id)
	return nil
}

// Reads the local project state and uploads it to the remote
func PushLocal(c *sdk.Connection) {
	if _, err := os.Stat("hms.toml"); err != nil {
		if os.IsNotExist(err) {
			fmt.Println("You can only push local state inside a hms-project.")
		} else {
			fmt.Println("Unknown error: ", err.Error())
		}
		os.Exit(1)
	}
	content, err := ioutil.ReadFile("./hms.toml")
	if err != nil {
		fmt.Printf("Could not push local state: failed to read `hms.toml`: %s\n", err.Error())
		os.Exit(1)
	}
	var configToml ConfigToml
	if err := toml.Unmarshal(content, &configToml); err != nil {
		fmt.Printf("Could not push local state: failed to parse `hms.toml`: %s\n", err.Error())
		os.Exit(1)
	}
	hmsContent, err := ioutil.ReadFile(fmt.Sprintf("./%s.hms", configToml.Id))
	if err != nil {
		fmt.Printf("Could not push local state: failed to read homescript file: %s\n", err.Error())
		os.Exit(1)
	}
	// Fetch current remote state for diff
	remoteBef, err := c.GetHomescript(configToml.Id)
	if err != nil {
		switch err {
		case sdk.ErrUnprocessableEntity:
			fmt.Printf("Could not pull remote state: either the project does not exist on the remote or you don't have the required permission to access it.\n")
		case sdk.ErrPermissionDenied:
			fmt.Printf("Failed to pull remote state: permission denied: please ensure that you have the correct access rights to pull hms-objects.\n")
		default:
			fmt.Printf("Could not pull remote state: server responded with unknown error: %s\n", err.Error())
		}
		os.Exit(1)
	}
	// Send modification request
	if err := c.ModifyHomescript(sdk.HomescriptRequest{
		Id:                  configToml.Id,
		Name:                configToml.Name,
		Description:         configToml.Description,
		QuickActionsEnabled: configToml.QuickActionsEnabled,
		SchedulerEnabled:    configToml.SchedulerEnabled,
		Code:                string(hmsContent),
		MDIcon:              configToml.MDIcon,
	}); err != nil {
		switch err {
		case sdk.ErrUnprocessableEntity:
			fmt.Printf("Failed to push local project state: invalid data provided: %s\n", err.Error())
		case sdk.ErrPermissionDenied:
			fmt.Printf("Failed to push local project: permission denied: please ensure that you have the correct access rights to push hms-objects.\n")
		default:
			fmt.Printf("Failed to push local project: unknown error: %s\n", err.Error())
		}
		os.Exit(1)
	}
	// Display diff after successful change
	dmp := diffmatchpatch.New()
	diffs := dmp.DiffMain(remoteBef.Data.Code, string(hmsContent), false)
	noChange := true
	for _, d := range diffs {
		if d.Type != diffmatchpatch.DiffEqual {
			noChange = false
		}
	}
	if !noChange {
		fmt.Printf("Diff:\n%s\n", dmp.DiffPrettyText(diffs))
	}
	// Display `hms.toml` changes via diff
	tomlBef := ConfigToml{
		Id:                  remoteBef.Data.Id,
		Name:                remoteBef.Data.Name,
		Description:         remoteBef.Data.Description,
		QuickActionsEnabled: remoteBef.Data.QuickActionsEnabled,
		SchedulerEnabled:    remoteBef.Data.SchedulerEnabled,
		MDIcon:              remoteBef.Data.MDIcon,
	}
	// Display general Diff info
	if tomlBef != configToml {
		fmt.Println("Changes to `hms.toml` synced to remote")
	} else if noChange {
		fmt.Println("Everything up-to-date.")
	}
}

// Reads project state from the server and patches the local files accordingly
func PullLocal(c *sdk.Connection) {
	if _, err := os.Stat("hms.toml"); err != nil {
		if os.IsNotExist(err) {
			fmt.Println("Error: Not a hms-project.")
		} else {
			fmt.Println("Unknown error: ", err.Error())
		}
		os.Exit(1)
	}
	content, err := ioutil.ReadFile("./hms.toml")
	if err != nil {
		fmt.Printf("Could not pull remote state: failed to read `hms.toml`: %s\n", err.Error())
		os.Exit(1)
	}
	var configToml ConfigToml
	if err := toml.Unmarshal(content, &configToml); err != nil {
		fmt.Printf("Could not pull remote state: failed to parse `hms.toml`: %s\n", err.Error())
		os.Exit(1)
	}
	// Read hms-file content before change
	hmsContent, err := ioutil.ReadFile(fmt.Sprintf("./%s.hms", configToml.Id))
	if err != nil {
		fmt.Printf("Could not pull remote state: failed to read homescript file: %s\n", err.Error())
		os.Exit(1)
	}
	remote, err := c.GetHomescript(configToml.Id)
	if err != nil {
		switch err {
		case sdk.ErrUnprocessableEntity:
			fmt.Printf("Could not pull remote state: either the project does not exist on the remote or you don't have the required permission to access it.\n")
		case sdk.ErrPermissionDenied:
			fmt.Printf("Failed to pull remote state: permission denied: please ensure that you have the correct access rights to pull hms-objects.\n")
		default:
			fmt.Printf("Could not pull remote state: server responded with unknown error: %s\n", err.Error())
		}
		os.Exit(1)
	}
	data, err := toml.Marshal(ConfigToml{
		Id:                  remote.Data.Id,
		Name:                remote.Data.Name,
		Description:         remote.Data.Description,
		QuickActionsEnabled: remote.Data.QuickActionsEnabled,
		SchedulerEnabled:    remote.Data.SchedulerEnabled,
		MDIcon:              configToml.MDIcon,
	})
	if err != nil {
		fmt.Printf("Could not pull remote state: failed to parse server response: %s\n", err.Error())
		os.Exit(1)
	}
	if err := ioutil.WriteFile("hms.toml", data, 0775); err != nil {
		fmt.Printf("Could not pull remote state: failed to update `hms.toml` config file: %s\n", err.Error())
		os.Exit(1)
	}
	if err := ioutil.WriteFile(fmt.Sprintf("%s.hms", configToml.Id), []byte(remote.Data.Code), 0775); err != nil {
		fmt.Printf("Could not pull remote state: failed to update local `.hms` file: %s\n", err.Error())
		os.Exit(1)
	}
	// Display diff after successful change
	dmp := diffmatchpatch.New()
	diffs := dmp.DiffMain(string(hmsContent), remote.Data.Code, false)
	noChange := true
	for _, d := range diffs {
		if d.Type != diffmatchpatch.DiffEqual {
			noChange = false
		}
	}
	if !noChange {
		fmt.Printf("Diff:\n%s\n", dmp.DiffPrettyText(diffs))
	}
	// Display `hms.toml` changes via diff
	tomlBef := ConfigToml{
		Id:                  remote.Data.Id,
		Name:                remote.Data.Name,
		Description:         remote.Data.Description,
		QuickActionsEnabled: remote.Data.QuickActionsEnabled,
		SchedulerEnabled:    remote.Data.SchedulerEnabled,
		MDIcon:              configToml.MDIcon,
	}
	// Display general Diff info
	if tomlBef != configToml {
		fmt.Println("Changes to `hms.toml` synced from remote.")
	} else if noChange {
		fmt.Println("Everything up-to-date.")
	}
}

// Reads the `hms.toml` file in the current project and returns a struct
func ReadLocalData(c *sdk.Connection) (string, ConfigToml) {
	content, err := ioutil.ReadFile("./hms.toml")
	if err != nil {
		fmt.Printf("Could not run local file: failed to read `hms.toml`: %s\n", err.Error())
		os.Exit(1)
	}
	var configToml ConfigToml
	if err := toml.Unmarshal(content, &configToml); err != nil {
		fmt.Printf("Could not run local file: failed to parse `hms.toml`: %s\n", err.Error())
		os.Exit(1)
	}
	hmsContent, err := ioutil.ReadFile(fmt.Sprintf("./%s.hms", configToml.Id))
	if err != nil {
		fmt.Printf("Could not run local file: failed to read homescript file: %s\n", err.Error())
		os.Exit(1)
	}
	return string(hmsContent), configToml
}

// Displays a list of cloneable Homescripts of the current user
func ListAll(c *sdk.Connection) {
	scripts, err := c.ListHomescript()
	if err != nil {
		fmt.Printf("Could not list all homescripts: failed to load data from server: %s\n", err.Error())
		os.Exit(1)
	}

	headerFmt := color.New(color.FgGreen, color.Underline).SprintfFunc()
	columnFmt := color.New(color.FgYellow).SprintfFunc()

	tbl := table.New("ID", "Name", "MDIcon", "QuickActions", "Scheduler")
	tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)

	// Fill the table
	for _, script := range scripts {
		quickActionsIndicator, schedulerIndicator := "no", "no"
		if script.Data.QuickActionsEnabled {
			quickActionsIndicator = "yes"
		}
		if script.Data.SchedulerEnabled {
			schedulerIndicator = "yes"
		}

		tbl.AddRow(script.Data.Id, script.Data.Name, script.Data.MDIcon, quickActionsIndicator, schedulerIndicator)
	}
	tbl.Print()
}

func Clone(c *sdk.Connection, id string) {
	fmt.Printf("Cloning into `%s`...\nResolving remote project...\n", id)
	remote, err := c.GetHomescript(id)
	if err != nil {
		switch err {
		case sdk.ErrPermissionDenied:
			fmt.Printf("Error: Could not clone `%s`: you lack permission to access Homescript.\n", id)
		case sdk.ErrReadResponseBody:
			fmt.Printf("Error: Could not clone `%s`: server returned invalid response: %s.\n", id, err.Error())
		case sdk.ErrUnprocessableEntity:
			fmt.Printf("Error: Could not read from remote: Project `%s` not found.\nPlease ensure that you have the correct access rights and the project exists.\n", id)
		default:
			fmt.Printf("Fatal: Failed to clone `%s`: unknown error: %s\n", id, err.Error())
		}
		os.Exit(1)
	}
	json, err := json.Marshal(remote)
	if err != nil {
		fmt.Printf("Cannot display size of project: %s\n", err.Error())
	}
	fmt.Printf("Downloaded remote project (size: %dB).\n", len(json))
	if err := os.Mkdir(id, 0755); err != nil {
		if os.IsExist(err) {
			fmt.Printf("Error: Could not clone into `./%s`.\nFailed to initialize project root: specified directory already exists.\n", id)
		} else {
			fmt.Printf("Error: Could not clone into `./%s`.\nFailed to initialize project root: %s\n", id, err.Error())
		}
		os.Exit(1)
	}
	data, err := toml.Marshal(ConfigToml{
		Id:                  id,
		Name:                remote.Data.Name,
		Description:         remote.Data.Description,
		QuickActionsEnabled: remote.Data.QuickActionsEnabled,
		SchedulerEnabled:    remote.Data.SchedulerEnabled,
		MDIcon:              remote.Data.MDIcon,
	})
	if err != nil {
		fmt.Printf("Could not pull remote state: failed to parse server response: %s\n", err.Error())
		os.Exit(1)
	}
	if err := ioutil.WriteFile(fmt.Sprintf("./%s/hms.toml", id), data, 0775); err != nil {
		fmt.Printf("Could not pull remote state: failed to update `hms.toml` config file: %s\n", err.Error())
		os.Exit(1)
	}
	if err := ioutil.WriteFile(fmt.Sprintf("./%s/%s.hms", id, id), []byte(remote.Data.Code), 0775); err != nil {
		fmt.Printf("Could not pull remote state: failed to update local `.hms` file: %s\n", err.Error())
		os.Exit(1)
	}
}

func CloneAll(c *sdk.Connection) {
	fmt.Printf("Cloning all available Homescripts from `%s`...\n\n", c.SmarthomeURL.Host)
	start := time.Now()
	scripts, err := c.ListHomescript()
	if err != nil {
		fmt.Printf("Could not clone all homescripts: failed to load list from server: %s\n", err.Error())
		os.Exit(1)
	}
	scriptCount := len(scripts)
	for _, script := range scripts {
		Clone(c, script.Data.Id)
		fmt.Println()
	}

	projectSIndicator := "s"
	if scriptCount == 1 {
		projectSIndicator = ""
	}

	if scriptCount > 0 {
		fmt.Printf("Finished: cloned %d project%s in %.2fs.\n", scriptCount, projectSIndicator, time.Since(start).Seconds())
	}
}
