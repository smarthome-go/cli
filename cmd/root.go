package cmd

import (
	"fmt"
	"net/http"
	"os"

	"log"

	"github.com/MikMuellerDev/homescript-cli/cmd/homescript"
	"github.com/spf13/cobra"
)

var (
	Verbose  bool
	ShowInfo bool

	SmarthomeURL   string
	SessionCookies []*http.Cookie

	Username string
	Password string
)

var (
	rootCmd = &cobra.Command{
		Use:   "homescript",
		Short: "Homescript language CLI",
		Long: "" +
			"A CLI interface for testing the Homescript Programming language for Smarthome.\n" +
			"A working and set-up Smarthome server instance is required.\n" +
			"For more information and usage documentation visit:\n" +
			"\n" +
			"  \x1b[1;32mThe Homescript Programming Language:\x1b[1;0m\n" +
			"  - https://github.com/MikMuellerDev/homescript\n\n" +
			"  \x1b[1;33mThe CLI Interface For Homescript:\x1b[1;0m\n" +
			"  - https://github.com/MikMuellerDev/homescript-cli\n\n" +
			"  \x1b[1;34mThe Smarthome Server:\x1b[1;0m\n" +
			"  - https://github.com/MikMuellerDev/smarthome\n",
		Run: func(cmd *cobra.Command, args []string) {
			PromptLogin()
			Login()
			StartRepl()
			fmt.Print("\x1b[3J\033c")
		},
	}
)

func Execute() {
	cmdRun := &cobra.Command{
		Use:   "run [filename]",
		Short: "Run a homescript file",
		Long:  "Runs a homescript file and connects to the server",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			PromptLogin()
			Login()
			homescript.RunFile(args[0], SmarthomeURL, SessionCookies)
		},
	}
	cmdInfo := &cobra.Command{
		Use:   "info",
		Short: "Server Debug Info",
		Long:  "Prints debugging information about the server",
		Args:  cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			PromptLogin()
			Login()
			homescript.Run(
				`print(debugInfo)`,
				SmarthomeURL,
				SessionCookies,
			)
		},
	}

	rootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "verbose output")
	rootCmd.PersistentFlags().StringVarP(&Username, "username", "u", "", "smarthome user used for connection")
	rootCmd.PersistentFlags().StringVarP(&Password, "password", "p", "", "smarthome password used for connection")
	rootCmd.PersistentFlags().StringVarP(&SmarthomeURL, "smarthome-url", "s", "http://localhost", "Url used for connecting to smarthome")
	// Environment variables, same as the ones used in the docker image
	/*
		`SMARTHOME_ADMIN_PASSWORD`: Checks for the smarthome admin user
	*/
	if adminPassword, adminPasswordOk := os.LookupEnv("SMARTHOME_ADMIN_PASSWORD"); adminPasswordOk && Password == "" {
		Password = adminPassword
		if Verbose {
			log.Println("Found password from \x1b[1;33mSMARTHOME_ADMIN_PASSWORD\x1b[1;0m")
		}
	}
	rootCmd.AddCommand(cmdRun)
	rootCmd.AddCommand(cmdInfo)
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
