package cmd

import (
	"fmt"
	"io/ioutil"
	"os"

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
	configFilePath := fmt.Sprintf("%s/homescript.yaml", configDir)
	_, err = os.Stat(configFilePath)
	if os.IsNotExist(err) {
		fmt.Println("Config file does not exist, creating...")
		if err := os.WriteFile(configFilePath, []byte("Username: user\nPassword: password\nSmarthomeURL: http://localhost"), 0600); err != nil {
			fmt.Println("Could not create config file: ", err.Error())
			return
		}
		fmt.Println("...created")
		return
	}
	fileContent, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		fmt.Println("Failed to read Homescript config file")
		os.Exit(1)
	}
	if err := yaml.Unmarshal(fileContent, &Config); err != nil {
		fmt.Println(fmt.Sprintf("Failed to parse config file at %s: invalid YAML format: %s", configFilePath, err.Error()))
		os.Exit(1)
	}
	if Username == "" {
		if Verbose {
			fmt.Println("Selected username from config file.")
		}
		Username = Config["Username"]
	}
	if Password == "" {
		if Verbose {
			fmt.Println("Selected password from config file.")
		}
		Password = Config["Password"]
	}
	if Url == "http://localhost" {
		if Verbose {
			fmt.Println("Selected smarthome-url from config file.")
		}
		Url = Config["SmarthomeURL"]
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

func writeConfig(username string, password string, smarthomeUrl string) {
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
	data := map[string]string{
		"Username":     username,
		"Password":     password,
		"SmarthomeURL": smarthomeUrl,
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
	configFilePath := fmt.Sprintf("%s/homescript.yaml", configDir)
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

	if err := ioutil.WriteFile(configFilePath, output, 0600); err != nil {
		fmt.Println("Failed to update configuration: could not write to config file: ", err.Error())
		os.Exit(1)
	}
	fmt.Println("...updated")
}

func deleteConfigFile() {
	configDir, err := os.UserConfigDir()
	if err != nil {
		fmt.Println("Failed to delete configuration file: could not determine user's config directory")
		os.Exit(1)
	}
	configFilePath := fmt.Sprintf("%s/homescript.yaml", configDir)
	_, err = os.Stat(configFilePath)
	if os.IsNotExist(err) {
		fmt.Println("Did not delete configuration file: file is already deleted")
		os.Exit(0)
	}
	if err := os.Remove(configFilePath); err != nil {
		fmt.Println("Failed to delete configuration file: ", err.Error())
		os.Exit(1)
	}
}
