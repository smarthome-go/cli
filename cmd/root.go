package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/MikMuellerDev/smarthome_sdk"
)

const Version = "2.0.0-beta"

var (
	Verbose  bool
	Username string
	Password string
	Url      string

	Connection *smarthome_sdk.Connection
)

// Map for the config file
var Config = map[string]string{
	"Username":     "",
	"Password":     "",
	"SmarthomeURL": "",
}

var (
	rootCmd = &cobra.Command{
		Use:     "homescript",
		Short:   "Homescript language CLI",
		Version: Version,
		Long: "" +
			fmt.Sprintf("homescript-cli v%s : ", Version) +
			"A command line interface for the smarthome server using homescript\n" +
			"A working and set-up Smarthome server instance is required.\n" +
			"For more information and usage documentation visit:\n" +
			"\n" +
			"  \x1b[1;32mThe Homescript Programming Language:\x1b[1;0m\n" +
			"  - https://github.com/MikMuellerDev/homescript\n\n" +
			"  \x1b[1;33mThe CLI Interface For Homescript:\x1b[1;0m\n" +
			"  - https://github.com/MikMuellerDev/homescript-cli\n\n" +
			"  \x1b[1;34mThe Smarthome Server:\x1b[1;0m\n" +
			"  - https://github.com/MikMuellerDev/smarthome\n",
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
		Use:   "run [filename]",
		Short: "Run a homescript file",
		Long:  "Runs a homescript file and connects to the server",
		Args:  cobra.ExactArgs(1),
		PreRun: func(cmd *cobra.Command, args []string) {
			readConfigFile()
		},
		Run: func(cmd *cobra.Command, args []string) {
			startTime := time.Now()
			InitConn()
			content, err := ioutil.ReadFile(args[0])
			if err != nil {
				fmt.Printf("Could not execute Homescript file '%s' due to fs error: %s", args[1], err.Error())
				os.Exit(1)
			}
			exitCode := RunCode(string(content), args[0])
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
		Short: "Smarthome Server Debug Info",
		Long:  "Prints debugging information about the server",
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
			RunCode(strings.Join(args, "\n"), "stdin")
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
	cmdPowerSummary := &cobra.Command{
		Use:   "power",
		Short: "Power Summary",
		Long:  "A compact overview of estimated power usage and states",
		Args:  cobra.NoArgs,
		PreRun: func(cmd *cobra.Command, args []string) {
			readConfigFile()
		},
		Run: func(cmd *cobra.Command, args []string) {
			InitConn()
			powerStats()
		},
	}
	rootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "verbose output")
	rootCmd.PersistentFlags().StringVarP(&Username, "username", "u", "", "smarthome user used for connection")
	rootCmd.PersistentFlags().StringVarP(&Password, "password", "p", "", "the user's password used for connection")
	rootCmd.PersistentFlags().StringVarP(&Url, "ip", "i", "http://localhost", "URL used for connecting to Smarthome")

	rootCmd.AddCommand(cmdRun)
	rootCmd.AddCommand(cmdInfo)
	rootCmd.AddCommand(cmdPipeIn)
	rootCmd.AddCommand(cmdListSwitches)
	rootCmd.AddCommand(cmdPowerSummary)

	// Parent configuration commands
	cmdConfig := &cobra.Command{
		Use:   "config",
		Short: "REPL configuration",
		Long:  "Retrieve and update the REPL configuration. If no arguments are provided, the configuration is printed. The configuration can be updated with [Username, Password, SmarthomeURL]",
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
			writeConfig(setUsername, setPassword, setURL)
		},
	}
	cmdConfigSet.Flags().StringVarP(&setUsername, "new-username", "n", "", "username to be updated")
	cmdConfigSet.Flags().StringVarP(&setPassword, "new-password", "t", "", "password to be updated")
	cmdConfigSet.Flags().StringVarP(&setURL, "new-ip", "a", "", "url / ip to be updated")
	cmdConfig.AddCommand(cmdConfigSet)
	rootCmd.AddCommand(cmdConfig)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
