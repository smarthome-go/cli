package cmd

import (
	"fmt"
	"net/http"
	"os"

	"github.com/MikMuellerDev/homescript-cli/cmd/homescript"
	"github.com/MikMuellerDev/homescript-cli/cmd/log"
	"github.com/sirupsen/logrus"
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
			if Verbose {
				log.InitLogger(logrus.TraceLevel)
			} else {
				log.InitLogger(logrus.InfoLevel)
			}
			PromptLogin()
			Login()
			DisplayServerInfo()
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
		Run: func(cmd *cobra.Command, args []string) {
			if Verbose {
				log.InitLogger(logrus.TraceLevel)
			} else {
				log.InitLogger(logrus.InfoLevel)
			}
			PromptLogin()
			Login()
			DisplayServerInfo()
			homescript.RunFile(Username, args[0], SmarthomeURL, SessionCookies)
		},
	}
	rootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "verbose output")
	rootCmd.PersistentFlags().BoolVarP(&ShowInfo, "info", "i", false, "show server info")
	rootCmd.PersistentFlags().StringVarP(&Username, "username", "u", "", "smarthome user used for connection")
	rootCmd.PersistentFlags().StringVarP(&Password, "password", "p", "", "smarthome password used for connection")
	rootCmd.PersistentFlags().StringVarP(&SmarthomeURL, "smarthome-url", "s", "http://localhost:8082", "Url used for connecting to smarthome")
	// Environment variables, same as the ones used in the docker image
	/*
		`SMARTHOME_ADMIN_PASSWORD`: Checks for the smarthome admin user
	*/
	if adminPassword, adminPasswordOk := os.LookupEnv("SMARTHOME_ADMIN_PASSWORD"); adminPasswordOk && Password == "" {
		Password = adminPassword
		log.Debug("Found password from \x1b[1;33mSMARTHOME_ADMIN_PASSWORD\x1b[1;0m")
	}
	rootCmd.AddCommand(cmdRun)
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
