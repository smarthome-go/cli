package workspace

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/pelletier/go-toml"
	"github.com/smarthome-go/sdk"
)

type ConfigToml struct {
	Id                  string `toml:"id"`
	Name                string `toml:"name"`
	Description         string `toml:"description"`
	QuickActionsEnabled bool   `toml:"quickActions"`
	SchedulerEnabled    bool   `toml:"scheduler"`
}

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
			default:
				fmt.Printf("Failed to remove project: unknown error: %s\n", err.Error())
			}
			os.Exit(1)
		}
		fmt.Printf("Deleted project `%s` from remote.\n", id)
	}
}

// Used by `createProject`
func createProjectConfigFile(id string, name string, path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	if err := toml.NewEncoder(file).Encode(ConfigToml{
		Id:   id,
		Name: name,
	}); err != nil {
		return err
	}
	if err := file.Close(); err != nil {
		return err
	}
	return nil
}

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
	if err := c.ModifyHomescript(sdk.HomescriptRequest{
		Id:                  configToml.Id,
		Name:                configToml.Name,
		Description:         configToml.Description,
		QuickActionsEnabled: configToml.QuickActionsEnabled,
		SchedulerEnabled:    configToml.SchedulerEnabled,
		Code:                string(hmsContent),
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
}

func PullLocal(c *sdk.Connection) {
	if _, err := os.Stat("hms.toml"); err != nil {
		if os.IsNotExist(err) {
			fmt.Println("You can only pull remote state inside a hms-project.")
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
	remote, err := c.GetHomescript(configToml.Id)
	if err != nil {
		switch err {
		case sdk.ErrUnprocessableEntity:
			fmt.Printf("Could not pull remote state: either the project does not exist on the remote or you don't have the required permission to access it.")
		case sdk.ErrPermissionDenied:
			fmt.Printf("Failed to pull remote state: permission denied: please ensure that you have the correct access rights to pull hms-objects.\n")
		default:
			fmt.Printf("Could not pull remote state: server responded with unknown error: %s\n", err.Error())
		}
		os.Exit(1)
	}
	data, err := toml.Marshal(ConfigToml{
		Id:                  remote.Id,
		Name:                remote.Name,
		Description:         remote.Description,
		QuickActionsEnabled: remote.QuickActionsEnabled,
		SchedulerEnabled:    remote.SchedulerEnabled,
	})
	if err != nil {
		fmt.Printf("Could not pull remote state: failed to parse server response: %s\n", err.Error())
		os.Exit(1)
	}
	if err := ioutil.WriteFile("hms.toml", data, 0775); err != nil {
		fmt.Printf("Could not pull remote state: failed to update `hms.toml` config file: %s\n", err.Error())
		os.Exit(1)
	}
	if err := ioutil.WriteFile(fmt.Sprintf("%s.hms", configToml.Id), []byte(remote.Code), 0775); err != nil {
		fmt.Printf("Could not pull remote state: failed to update local `.hms` file: %s\n", err.Error())
		os.Exit(1)
	}
}
