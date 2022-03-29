package cmd

import (
	"fmt"
	"net/http"
	"os"

	"github.com/MikMuellerDev/homescript-cli/cmd/homescript"
	"github.com/MikMuellerDev/homescript-cli/cmd/log"
	"github.com/spf13/cobra"
)

const Version = "0.4.1-beta"

var (
	Verbose  bool
	ShowInfo bool
	Silent   bool

	SmarthomeURL   string
	SessionCookies []*http.Cookie

	Username string
	Password string
)

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
		Run: func(cmd *cobra.Command, args []string) {
			log.InitLog(Verbose)
			log.Silent = Silent
			homescript.Silent = Silent
			PingServer()
			PromptLogin()
			Login(true)
			StartRepl()
			// fmt.Print("\x1b[3J\033c")
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
			log.InitLog(Verbose)
			log.Silent = Silent
			homescript.Silent = Silent
			PingServer()
			PromptLogin()
			Login(true)
			homescript.RunFile(args[0], SmarthomeURL, SessionCookies)
		},
	}
	cmdInfo := &cobra.Command{
		Use:   "info",
		Short: "Smarthome Server Debug Info",
		Long:  "Prints debugging information about the server",
		Args:  cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			log.InitLog(Verbose)
			log.Silent = Silent
			homescript.Silent = Silent
			PingServer()
			PromptLogin()
			Login(true)
			homescript.Run(
				`print(debugInfo)`,
				SmarthomeURL,
				SessionCookies,
			)
		},
	}
	cmdPipeIn := &cobra.Command{
		Use:   "pipe",
		Short: "Run Code via Stdin",
		Long:  "Run code via Stdin without interactive prompts and output. Ideal for bash-based scripting.",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			log.InitLog(Verbose)
			log.Silent = Silent
			homescript.Silent = Silent
			PingServer()
			Login(false)
			homescript.Run(
				`print(debugInfo)`,
				SmarthomeURL,
				SessionCookies,
			)
		},
	}
	cmdListSwitches := &cobra.Command{
		Use:   "switches",
		Short: "List switches",
		Long:  "List switches of the current user",
		Args:  cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			log.InitLog(Verbose)
			log.Silent = Silent
			homescript.Silent = Silent
			PingServer()
			Login(true)
			getPersonalSwitches()
			listSwitches()
		},
	}
	rootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "verbose output")
	rootCmd.PersistentFlags().BoolVarP(&Silent, "silent", "s", false, "no output")
	rootCmd.PersistentFlags().StringVarP(&Username, "username", "u", "", "smarthome user used for connection")
	rootCmd.PersistentFlags().StringVarP(&Password, "password", "p", "", "smarthome password used for connection")
	rootCmd.PersistentFlags().StringVarP(&SmarthomeURL, "ip", "i", "http://localhost", "Url used for connecting to smarthome")
	log.InitLog(true)
	// Environment variables, same as the ones used in the docker image
	/*
		`SMARTHOME_ADMIN_PASSWORD`: Checks for the smarthome admin user
	*/
	if adminPassword, adminPasswordOk := os.LookupEnv("SMARTHOME_ADMIN_PASSWORD"); adminPasswordOk && Password == "" {
		Password = adminPassword
		if Verbose {
			log.Logn("Found password from \x1b[1;33mSMARTHOME_ADMIN_PASSWORD\x1b[1;0m")
		}
	}
	rootCmd.AddCommand(cmdRun)
	rootCmd.AddCommand(cmdInfo)
	rootCmd.AddCommand(cmdPipeIn)
	rootCmd.AddCommand(cmdListSwitches)
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
