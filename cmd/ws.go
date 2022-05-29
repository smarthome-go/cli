package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/smarthome-go/cli/cmd/workspace"
	"github.com/spf13/cobra"
)

func createCmdWs() *cobra.Command {
	// Parent ws commands
	cmdWS := &cobra.Command{
		Use:   "ws",
		Short: "workspace",
		Long:  "Allows the user to create workspaces and develop homescript files",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			if err := cmd.Help(); err != nil {
				panic(err.Error())
			}
		},
	}
	cmdWSInit := &cobra.Command{
		Use:   "new [hms-id] [project-name]",
		Short: "new project",
		Long:  "Creates a new project and creates a new Homescript on the remote",
		Args:  cobra.RangeArgs(1, 2),
		PreRun: func(cmd *cobra.Command, args []string) {
			readConfigFile()
		},
		Run: func(cmd *cobra.Command, args []string) {
			InitConn()
			if len(args) == 2 {
				workspace.New(args[0], args[1], Connection)
			} else {
				workspace.New(args[0], "", Connection)
			}
		},
	}
	cmdWSPush := &cobra.Command{
		Use:   "push",
		Short: "push project state",
		Long:  "Reads local changes and pushes local state to the remote",
		Args:  cobra.NoArgs,
		PreRun: func(cmd *cobra.Command, args []string) {
			readConfigFile()
		},
		Run: func(cmd *cobra.Command, args []string) {
			InitConn()
			workspace.PushLocal(Connection)
		},
	}
	cmdWSL := &cobra.Command{
		Use:   "ls",
		Short: "List remote projects",
		Long:  "Displays a list of available remote projects to clone.",
		Args:  cobra.NoArgs,
		PreRun: func(cmd *cobra.Command, args []string) {
			readConfigFile()
		},
		Run: func(cmd *cobra.Command, args []string) {
			InitConn()
			workspace.ListAll(Connection)
		},
	}
	cmdWSPull := &cobra.Command{
		Use:   "pull",
		Short: "pull project state",
		Long:  "Fetches remote changes and writes them to the local project",
		Args:  cobra.NoArgs,
		PreRun: func(cmd *cobra.Command, args []string) {
			readConfigFile()
		},
		Run: func(cmd *cobra.Command, args []string) {
			InitConn()
			workspace.PullLocal(Connection)
		},
	}
	var runOnlyLocal = false
	cmdWsRun := &cobra.Command{
		Use:   "run",
		Short: "Run current project",
		Long:  "Executes the local state of the homescript file on the server",
		Args:  cobra.ArbitraryArgs,
		PreRun: func(cmd *cobra.Command, args []string) {
			readConfigFile()
		},
		Run: func(cmd *cobra.Command, args []string) {
			startTime := time.Now()
			// Read local workspace data
			content, config, err := workspace.ReadLocalData(Connection)
			if err != nil {
				fmt.Printf("Error: %s\n", err.Error())
				os.Exit(1)
			}
			// Prepare Homescript arguments
			hmsArgs := make(map[string]string, 0)
			if len(args) > 0 {
				hmsArgsTemp, err := processHmsArgs(args)
				if err != nil {
					fmt.Println(err.Error())
					os.Exit(1)
				}
				hmsArgs = hmsArgsTemp
			}
			// Initialize connection to the Smarthome server
			InitConn()
			// Run the Homescript code
			var exitCode int
			if runOnlyLocal {
				if Verbose {
					fmt.Printf("Executing `%s.hms` using local state", config.Id)
				}
				exitCode = RunCode(string(content), hmsArgs, fmt.Sprintf("%s.hms", config.Id))
			} else {
				if Verbose {
					fmt.Printf("Executing `%s.hms` using remote state", config.Id)
				}
				exitCode = RunById(config.Id, hmsArgs)
			}
			if exitCode != 0 {
				fmt.Printf("Homescript terminated with exit code: %d \x1b[90m[%.2fs]\x1b[1;0m\n", exitCode, time.Since(startTime).Seconds())
			} else {
				fmt.Printf("Homescript was executed successfully: %d \x1b[90m[%.2fs]\x1b[1;0m\n", exitCode, time.Since(startTime).Seconds())
			}
			os.Exit(exitCode)
		},
	}
	cmdWsRun.Flags().BoolVarP(&runOnlyLocal, "local", "l", false, "whether the file should be executed using the local state instead of the remote state")

	var lintOnRemote = false
	cmdWsLint := &cobra.Command{
		Use:   "lint",
		Short: "Lint local homescript project",
		Long:  "Executes the script as dry-run in order to lint for errors",
		Args:  cobra.ArbitraryArgs,
		PreRun: func(cmd *cobra.Command, args []string) {
			readConfigFile()
		},
		Run: func(cmd *cobra.Command, args []string) {
			// Prepare Homescript arguments
			hmsArgs := make(map[string]string, 0)
			if len(args) > 0 {
				hmsArgsTemp, err := processHmsArgs(args)
				if err != nil {
					fmt.Println(err.Error())
					os.Exit(1)
				}
				hmsArgs = hmsArgsTemp
			}
			// Read the local workspace data
			content, config, err := workspace.ReadLocalData(Connection)
			if err != nil {
				fmt.Printf("Error: %s\n", err.Error())
				os.Exit(1)
			}
			// Initialize connection to the Smarthome server
			InitConn()
			// Lint the Homescript using the data and arguments
			var exitCode int
			if lintOnRemote {
				if Verbose {
					fmt.Printf("Linting `%s.hms` using current remote state", config.Id)
				}
				exitCode = LintById(config.Id, hmsArgs)
			} else {
				if Verbose {
					fmt.Printf("Linting `%s.hms` using local state", config.Id)
				}
				exitCode = LintCode(string(content), hmsArgs, fmt.Sprintf("%s.hms", config.Id))
			}
			os.Exit(exitCode)
		},
	}
	cmdWsLint.Flags().BoolVarP(&lintOnRemote, "remote", "r", false, "whether the file should be linted using the remote state instead of the local state")

	var purge bool
	cmdWSRemove := &cobra.Command{
		Use:   "rm [hms-id]",
		Short: "removes a project",
		Long:  "Removes project locally and will purge it on the remote if the -P flag is set.",
		Args:  cobra.ExactArgs(1),
		PreRun: func(cmd *cobra.Command, args []string) {
			readConfigFile()
		},
		Run: func(cmd *cobra.Command, args []string) {

			InitConn()
			workspace.Delete(args[0], purge, Connection)
		},
	}
	var all bool
	cmdWSClone := &cobra.Command{
		Use:   "clone [hms-id]",
		Short: "clones a project",
		Long:  "Downloads a remote project into a folder named equally to the ID. Sets it up for local development.",
		Args:  cobra.RangeArgs(0, 1),
		PreRun: func(cmd *cobra.Command, args []string) {
			readConfigFile()
		},
		Run: func(cmd *cobra.Command, args []string) {
			InitConn()
			if !all {
				if len(args) == 0 {
					fmt.Sprintln("Error: accepts 1 arg, received 0")
					if err := cmd.Help(); err != nil {
						panic(err.Error())
					}
				}
				workspace.Clone(Connection, args[0])
				os.Exit(0)
			} else {
				workspace.CloneAll(Connection)
				os.Exit(0)
			}
		},
	}
	cmdWSClone.Flags().BoolVarP(&all, "all", "a", false, "")
	cmdWSRemove.Flags().BoolVarP(&purge, "purge", "P", false, "whether the project should be deleted on the remote")
	cmdWS.AddCommand(cmdWSInit)
	cmdWS.AddCommand(cmdWSPush)
	cmdWS.AddCommand(cmdWSPull)
	cmdWS.AddCommand(cmdWSL)
	cmdWS.AddCommand(cmdWsRun)
	cmdWS.AddCommand(cmdWsLint)
	cmdWS.AddCommand(cmdWSRemove)
	cmdWS.AddCommand(cmdWSClone)
	return cmdWS
}
