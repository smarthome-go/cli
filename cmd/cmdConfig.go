package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func createCmdConfig() *cobra.Command {
	// Parent configuration commands
	cmdConfig := &cobra.Command{
		Use:   "config",
		Short: "CLI configuration",
		Long:  "Retrieve and update the CLI configuration. If no arguments are provided, the configuration is printed. The configuration can be updated with [Username, Password, SmarthomeURL]",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			if err := cmd.Help(); err != nil {
				panic(err.Error())
			}
			printConfig()
		},
	}

	// View current configuration
	cmdConfigGet := &cobra.Command{
		Use:   "get",
		Short: "View configuration",
		Long:  "View the parameters which are currently stored in the configuration file.",
		Args:  cobra.NoArgs,
		PreRun: func(cmd *cobra.Command, args []string) {
			readConfigFile()
		},
		Run: func(cmd *cobra.Command, args []string) {
			printConfig()
		},
	}
	cmdConfig.AddCommand(cmdConfigGet)

	// Delete configuration
	cmdConfigRm := &cobra.Command{
		Use:   "rm",
		Short: "Remove configuration",
		Long:  "Deletes the configuration file from the filesystem",
		Run: func(cmd *cobra.Command, args []string) {
			deleteConfigFile()
		},
	}
	cmdConfig.AddCommand(cmdConfigRm)

	// Update configuration
	setUsername := ""
	setPassword := ""
	setURL := ""
	setLintOnPush := false

	cmdConfigSet := &cobra.Command{
		Use:   "set",
		Short: "Update configuration",
		Long:  "Write new configuration values to the configuration file",
		PreRun: func(cmd *cobra.Command, args []string) {
			readConfigFile()
		},
		Run: func(cmd *cobra.Command, args []string) {
			if setPassword == "" && setUsername == "" && setURL == "" {
				fmt.Println("Provided at least one of the flags below in order to update the configuration.")
				if err := cmd.Help(); err != nil {
					panic(err.Error())
				}
				return
			}
			writeConfig(setUsername, setPassword, setURL, setLintOnPush)
		},
	}
	cmdConfigSet.Flags().StringVarP(&setUsername, "new-username", "n", "", "New username to be updated")
	cmdConfigSet.Flags().StringVarP(&setPassword, "new-password", "t", "", "New password to be updated")
	cmdConfigSet.Flags().StringVarP(&setURL, "new-ip", "a", "", "New URL / ip to be updated")
	cmdConfigSet.Flags().BoolVarP(&setLintOnPush, "new-lint-on-push", "l", false, "Whether to lint on every HMS workspace push action or not")
	cmdConfig.AddCommand(cmdConfigSet)
	return cmdConfig
}
