package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/smarthome-go/sdk"
)

const Version = "2.8.0"

var (
	Verbose  bool
	Username string
	Password string
	Url      string

	Connection *sdk.Connection
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
			content, err := ioutil.ReadFile(args[0])
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
			exitCode := RunCode(string(content), hmsArgs, args[0])
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
			RunCode(strings.Join(args, "\n"), make(map[string]string, 0), "stdin")
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

	rootCmd.AddCommand(cmdRun)
	rootCmd.AddCommand(cmdInfo)
	rootCmd.AddCommand(cmdPipeIn)
	rootCmd.AddCommand(cmdListSwitches)
	rootCmd.AddCommand(cmdPowerSummary)

	// Subcommands
	rootCmd.AddCommand(createCmdConfig())
	rootCmd.AddCommand(createCmdWs())

	rootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "verbose output")
	rootCmd.PersistentFlags().StringVarP(&Username, "username", "u", "", "smarthome user used for connection")
	rootCmd.PersistentFlags().StringVarP(&Password, "password", "p", "", "the user's password used for connection")
	rootCmd.PersistentFlags().StringVarP(&Url, "ip", "i", "http://localhost", "URL used for connecting to Smarthome")

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
