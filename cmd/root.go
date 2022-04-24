package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/MikMuellerDev/homescript-cli/cmd/homescript"
	"github.com/MikMuellerDev/homescript-cli/cmd/log"
)

const Version = "0.5.3-beta"

var (
	Verbose  bool
	ShowInfo bool
	Silent   bool

	SmarthomeURL   string
	SessionCookies []*http.Cookie

	Username string
	Password string
)

type Config struct {
	Username     string `json:"username"`
	Password     string `json:"password"`
	SmarthomeUrl string `json:"smarthomeUrl"`
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
		Run: func(cmd *cobra.Command, args []string) {
			log.InitLog(Verbose)
			log.Silent = Silent
			homescript.Silent = Silent
			PingServer()
			PromptLogin()
			Login(true)
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
		Use:   "debug",
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
			debugInfo()
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
				strings.Join(args, " \n"),
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
	readConfigFile()
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}

func readConfigFile() {
	configDir, err := os.UserConfigDir()
	if err != nil {
		log.Loge("Failed to determine user config directory")
		return
	}
	configFilePath := fmt.Sprintf("%s/homescript.json", configDir)
	_, err = os.Stat(configFilePath)
	if os.IsNotExist(err) {
		log.Logn("Config file does not exists")
		return
	}
	fileContent, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		log.Loge("Failed to read homescript config file")
		return
	}
	decoder := json.NewDecoder(bytes.NewReader(fileContent))
	decoder.DisallowUnknownFields()
	var config Config
	if err := decoder.Decode(&config); err != nil {
		log.Loge("Failed to parse config file to struct: invalid json format: ", err.Error())
		return
	}
	if Username == "" {
		Username = config.Username
	}
	if Password == "" {
		Username = config.Password
	}
	if SmarthomeURL == "" {
		SmarthomeURL = config.SmarthomeUrl
	}
}
