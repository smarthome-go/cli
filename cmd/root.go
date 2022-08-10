package cmd

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/smarthome-go/cli/cmd/workspace"
	"github.com/smarthome-go/sdk"
)

const Version = "2.14.1"

// Cli override configuration
var (
	Verbose bool
	// Configuration from the config file
	Config Configuration
	// Override parameters from the CLI
	overrideConfig Configuration
	// Connection used for Smarthome
	Connection *sdk.Connection
)

var (
	rootCmd = &cobra.Command{
		Use:     "smarthome-cli",
		Short:   "Smarthome CLI",
		Version: Version,
		Long: "" +
			fmt.Sprintf("smarthome-cli v%s : ", Version) +
			"A command line interface for the smarthome server\n" +
			"A working and set-up Smarthome server instance is required.\n" +
			"For more information and usage documentation visit:\n" +
			"\n" +
			"  \x1b[1;32mThe Homescript Programming Language:\x1b[1;0m\n" +
			"  - https://github.com/smarthome-go/homescript\n\n" +
			"  \x1b[1;33mThe CLI Interface For Homescript:\x1b[1;0m\n" +
			"  - https://github.com/smarthome-go/cli\n\n" +
			"  \x1b[1;34mThe Smarthome Server:\x1b[1;0m\n" +
			"  - https://github.com/smarthome-go/smarthome\n",
		PreRun: func(cmd *cobra.Command, args []string) {
			readConfigFile()
		},
		Run: func(cmd *cobra.Command, args []string) {
			InitConn()
			StartRepl()
		},
	}
)

func Execute() {
	cmdRun := &cobra.Command{
		Use:   "run [filename] [key:value]",
		Short: "Run a Homescript file",
		Long:  "Runs a local Homescript file with arguments on the Smarthome server",
		Args:  cobra.MinimumNArgs(1),
		PreRun: func(cmd *cobra.Command, args []string) {
			readConfigFile()
		},
		Run: func(cmd *cobra.Command, args []string) {
			startTime := time.Now()
			// Read file
			content, err := os.ReadFile(args[0])
			if err != nil {
				fmt.Printf("Could not execute Homescript file '%s' due to fs error: %s\n", args[0], err.Error())
				os.Exit(1)
			}
			// Prepare Homescript arguments
			hmsArgs, err := processHmsArgs(args[1:])
			if err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}
			// Initialize Smarthome connection
			InitConn()
			// Execute code
			exitCode := workspace.RunCode(
				Connection,
				string(content),
				hmsArgs,
				args[0],
			)
			if exitCode != 0 {
				fmt.Printf("Homescript terminated with exit code: %d \x1b[90m[%.2fs]\x1b[1;0m\n", exitCode, time.Since(startTime).Seconds())
			} else {
				fmt.Printf("Homescript was executed successfully: %d \x1b[90m[%.2fs]\x1b[1;0m\n", exitCode, time.Since(startTime).Seconds())
			}
			os.Exit(exitCode)
		},
	}
	cmdInfo := &cobra.Command{
		Use:   "debug",
		Short: "Server Debug Info",
		Long:  "Prints debugging information about the Smarthome server",
		Args:  cobra.NoArgs,
		PreRun: func(cmd *cobra.Command, args []string) {
			readConfigFile()
		},
		Run: func(cmd *cobra.Command, args []string) {
			InitConn()
			printDebugInfo()
		},
	}
	cmdPipeIn := &cobra.Command{
		Use:   "pipe",
		Short: "Run Code via Stdin",
		Long:  "Run code via Stdin without interactive prompts and output. Ideal for bash-based scripting.",
		Args:  cobra.MinimumNArgs(1),
		PreRun: func(cmd *cobra.Command, args []string) {
			readConfigFile()
		},
		Run: func(cmd *cobra.Command, args []string) {
			InitConn()
			workspace.RunCode(Connection,
				strings.Join(args, "\n"),
				make(map[string]string, 0),
				"stdin",
			)
		},
	}
	cmdListSwitches := &cobra.Command{
		Use:   "switches",
		Short: "List switches",
		Long:  "List switches of the current user",
		Args:  cobra.NoArgs,
		PreRun: func(cmd *cobra.Command, args []string) {
			readConfigFile()
		},
		Run: func(cmd *cobra.Command, args []string) {
			InitConn()
			listSwitches()
		},
	}

	rootCmd.AddCommand(cmdRun)
	rootCmd.AddCommand(cmdInfo)
	rootCmd.AddCommand(cmdPipeIn)
	rootCmd.AddCommand(cmdListSwitches)

	// Subcommands
	rootCmd.AddCommand(createCmdConfig())
	rootCmd.AddCommand(createCmdWs())
	rootCmd.AddCommand(createCmdPower())

	rootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "Enables verbose output")
	rootCmd.PersistentFlags().StringVarP(&overrideConfig.Credentials.Username, "username", "u", "", "Smarthome-user used for the connection")
	rootCmd.PersistentFlags().StringVarP(&overrideConfig.Credentials.Password, "password", "p", "", "The user's password used for connection")
	rootCmd.PersistentFlags().StringVarP(&overrideConfig.Connection.SmarthomeUrl, "ip", "i", "", "URL of the target Smarthome instance")

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
