package cmd

import (
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

	// Update login credentials
	cmdConfigSet := &cobra.Command{
		Use:   "login",
		Short: "Save login credentials",
		Long:  "Login to a Smarthome server and save the configutaion",
		PreRun: func(cmd *cobra.Command, args []string) {
			readConfigFile()
		},
		Run: func(cmd *cobra.Command, args []string) {
			// Prompt the user's login credentials
			PromptLogin(true)
			// Try to connect
			InitConn()
			// On success, write the configuration to the file
			if Config.Connection.UseToken {
				// If token authentication is used, do not write a username or a password to the new file
				Config.Credentials.Username = ""
				Config.Credentials.Password = ""
			}
			writeConfig(Config)
		},
	}
	cmdConfig.AddCommand(cmdConfigSet)
	return cmdConfig
}
