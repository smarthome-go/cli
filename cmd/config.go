package cmd

import (
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/rodaine/table"
	"gopkg.in/yaml.v3"
)

func readConfigFile() {
	configDir, err := os.UserConfigDir()
	if err != nil {
		fmt.Println("Failed to determine user config directory, not reading config file")
		return
	}
	configFilePath := fmt.Sprintf("%s/smarthome-cli.yaml", configDir)
	_, err = os.Stat(configFilePath)
	if os.IsNotExist(err) {
		if Verbose {
			fmt.Println("Configuration file does not exist, creating...")
		}
		marshaled, err := yaml.Marshal(Config)
		if err != nil {
			fmt.Println("Could not create config file: ", err.Error())
			return
		}
		if err := os.WriteFile(configFilePath, marshaled, 0600); err != nil {
			fmt.Println("Could not create config file: ", err.Error())
			return
		}
		if Verbose {
			fmt.Printf("Created new configuration at %s\n", configFilePath)
		}
		return
	}
	fileContent, err := os.ReadFile(configFilePath)
	if err != nil {
		fmt.Println("Failed to read Homescript config file")
		os.Exit(1)
	}
	if err := yaml.Unmarshal(fileContent, &Config); err != nil {
		fmt.Printf("Failed to parse config file at %s: invalid YAML format: %s\n", configFilePath, err.Error())
		os.Exit(1)
	}
	if Username == "" && Config["Username"] != "" {
		if Verbose {
			fmt.Println("Selected username from config file.")
		}
		Username = Config["Username"]
	}
	if Password == "" && Config["Password"] != "" {
		if Verbose {
			fmt.Println("Selected password from config file.")
		}
		Password = Config["Password"]
	}
	if Url == "http://localhost" && Config["SmarthomeURL"] != "http://localhost" {
		if Verbose {
			fmt.Println("Selected smarthome-url from config file.")
		}
		Url = Config["SmarthomeURL"]
	}
	if Config["LintOnPush"] != "no" && Config["LintOnPush"] != "yes" {
		fmt.Printf("Unexpected value in config file: `LintOnPush` holds invalid value: `%s`\n", Config["LintOnPush"])
		os.Exit(1)
	}
	if LintOnPush && Config["LintOnPush"] != "yes" {
		if Verbose {
			fmt.Println("Selected lint-on-push option from config file.")
		}
		LintOnPush = false
	}
}

func printConfig() {
	readConfigFile()
	headerFmt := color.New(color.FgGreen, color.Underline).SprintfFunc()
	columnFmt := color.New(color.FgYellow).SprintfFunc()
	tbl := table.New("Username", "Password", "SmarthomeUrl")
	tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)
	tbl.AddRow(Config["Username"], Config["Password"], Config["SmarthomeURL"])
	tbl.Print()
}

func writeConfig(username string, password string, smarthomeUrl string, lintOnPush bool) {
	fmt.Println("Updating REPL configuration...")
	readConfigFile()
	if username == "" {
		username = Username
	}
	if password == "" {
		password = Password
	}
	if smarthomeUrl == "" {
		smarthomeUrl = Url
	}
	if !strings.HasPrefix(smarthomeUrl, "https://") && !strings.HasPrefix(smarthomeUrl, "http://") {
		smarthomeUrl = "http://" + smarthomeUrl
	}
	if _, err := url.Parse(smarthomeUrl); err != nil {
		fmt.Println("Invalid URL specified: please provide a valid URL.")
		os.Exit(1)
	}
	lintOnPushString := "no"
	if lintOnPush {
		lintOnPushString = "yes"
	}
	data := map[string]string{
		"Username":     username,
		"Password":     password,
		"SmarthomeURL": smarthomeUrl,
		"LintOnPush":   lintOnPushString,
	}
	output, err := yaml.Marshal(&data)
	if err != nil {
		fmt.Println("Failed to update configuration: could not encode config file", err.Error())
		os.Exit(1)
	}
	configDir, err := os.UserConfigDir()
	if err != nil {
		fmt.Println("Failed to update configuration: could not determine user's config directory")
		os.Exit(1)
	}
	configFilePath := fmt.Sprintf("%s/smarthome-cli.yaml", configDir)
	_, err = os.Stat(configFilePath)
	if os.IsNotExist(err) {
		fmt.Println("Config file does not exist, creating...")
		if err := os.WriteFile(configFilePath, []byte(output), 0600); err != nil {
			fmt.Println("Failed to update configuration: could not create new configuration file: ", err.Error())
			os.Exit(1)
		}
		fmt.Println("...created and written")
		return
	}

	if err := os.WriteFile(configFilePath, output, 0600); err != nil {
		fmt.Println("Failed to update configuration: could not write to config file: ", err.Error())
		os.Exit(1)
	}
	fmt.Println("...updated")
}

func deleteConfigFile() {
	if Verbose {
		fmt.Println("Deleting configuration file...")
	}
	configDir, err := os.UserConfigDir()
	if err != nil {
		fmt.Println("Failed to delete configuration file: could not determine user's config directory")
		os.Exit(1)
	}
	configFilePath := fmt.Sprintf("%s/smarthome-cli.yaml", configDir)
	_, err = os.Stat(configFilePath)
	if os.IsNotExist(err) {
		fmt.Println("Did not delete configuration file: file is already deleted")
		os.Exit(0)
	}
	if err := os.Remove(configFilePath); err != nil {
		fmt.Println("Failed to delete configuration file: ", err.Error())
		os.Exit(1)
	}
	fmt.Printf("Successfully deleted configuration file from %s\n", configFilePath)
}
