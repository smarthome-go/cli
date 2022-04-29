package cmd

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/rodaine/table"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"github.com/MikMuellerDev/homescript-cli/cmd/homescript"
	"github.com/MikMuellerDev/homescript-cli/cmd/log"
)

const Version = "0.8.0-beta"

var (
	Verbose  bool
	ShowInfo bool
	Silent   bool

	SmarthomeURL   string
	SessionCookies []*http.Cookie

	Username string
	Password string
)

var Config = map[string]string{
	Username:     "",
	Password:     "",
	SmarthomeURL: "",
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
	cmdConfig := &cobra.Command{
		Use:   "config",
		Short: "REPL configuration",
		Long:  "Retrieve and update the REPL configuration. If no arguments are provided, the configuration is printed. The configuration can be updated with [Username, Password, SmarthomeURL]",
		Args:  cobra.RangeArgs(0, 3),
		Run: func(cmd *cobra.Command, args []string) {
			log.InitLog(Verbose)
			log.Silent = Silent
			homescript.Silent = Silent
			if len(args) == 0 {
				listConfig()
				return
			}
			if len(args) != 3 {
				log.Loge("Setting configuration requires exactly 3 arguments: [Username, Password, SmarthomeURL]")
				return
			}
			writeConfig(args[0], args[1], args[2])
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

	rootCmd.AddCommand(cmdRun)
	rootCmd.AddCommand(cmdInfo)
	rootCmd.AddCommand(cmdPipeIn)
	rootCmd.AddCommand(cmdListSwitches)
	rootCmd.AddCommand(cmdConfig)

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
	configFilePath := fmt.Sprintf("%s/homescript.yaml", configDir)
	_, err = os.Stat(configFilePath)
	if os.IsNotExist(err) {
		log.Logn("Config file does not exist, creating...")
		if err := os.WriteFile(configFilePath, []byte("Username: user\nPassword: password\nSmarthomeURL: http://localhost"), 0600); err != nil {
			log.Loge("Could not create config file: ", err.Error())
			return
		}
		log.Logn("...created")
		return
	}
	fileContent, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		log.Loge("Failed to read homescript config file")
		return
	}
	if err := yaml.Unmarshal(fileContent, &Config); err != nil {
		log.Loge(fmt.Sprintf("Failed to parse config file at %s: invalid YAML format: %s", configFilePath, err.Error()))
		return
	}
	if Username == "" {
		if Verbose {
			log.Logn("Selected username from config file.")
		}
		Username = Config["Username"]
	}
	if Password == "" {
		if Verbose {
			log.Logn("Selected password from config file.")
		}
		Password = Config["Password"]
	}
	if SmarthomeURL == "http://localhost" {
		if Verbose {
			log.Logn("Selected smarthome-url from config file.")
		}
		SmarthomeURL = Config["SmarthomeURL"]
	}
}

func listConfig() {
	readConfigFile()
	headerFmt := color.New(color.FgGreen, color.Underline).SprintfFunc()
	columnFmt := color.New(color.FgYellow).SprintfFunc()
	tbl := table.New("Username", "Password", "SmarthomeUrl")
	tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)
	tbl.AddRow(Config["Username"], Config["Password"], Config["SmarthomeURL"])
	tbl.Print()
}

func writeConfig(username string, password string, smarthomeUrl string) {
	log.Logn("Updating REPL configuration...")
	data := map[string]string{
		"Username":     username,
		"Password":     password,
		"SmarthomeURL": smarthomeUrl,
	}
	output, err := yaml.Marshal(&data)
	if err != nil {
		log.Loge("Could not encode config file", err.Error())
		return
	}
	configDir, err := os.UserConfigDir()
	if err != nil {
		log.Loge("Failed to determine user config directory")
		return
	}
	configFilePath := fmt.Sprintf("%s/homescript.yaml", configDir)
	_, err = os.Stat(configFilePath)
	if os.IsNotExist(err) {
		log.Logn("Config file does not exist, creating...")
		if err := os.WriteFile(configFilePath, []byte(output), 0600); err != nil {
			log.Loge("Could not create config file: ", err.Error())
			return
		}
		log.Logn("...created and written")
		return
	}

	if err := ioutil.WriteFile(configFilePath, output, 0600); err != nil {
		log.Loge("Could not write config file: ", err.Error())
		return
	}
	log.Logn("...updated")
}
