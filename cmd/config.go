package cmd

import (
	"fmt"
	"net/url"
	"os"
	"strings"
	"unicode/utf8"

	"github.com/fatih/color"
	"github.com/pelletier/go-toml"
	"github.com/rodaine/table"
)

// Is appended to the user's configuration directory path
const filePathPrefix = "smarthome-cli"

// Is appended to the user's configuration directory path after `fileName`
const fileName = "config.toml"

var filePath = fmt.Sprintf("%s/%s", filePathPrefix, fileName)

type Configuration struct {
	Connection  ConnectionConfig `toml:"connection"`  // Connection settings
	Credentials Credentials      `toml:"credentials"` // Credential store
	Homescript  HomescriptConfig `toml:"homescript"`  // Homescript settings
}

type ConnectionConfig struct {
	SmarthomeUrl string `toml:"smarthome_url"`        // Connection URL
	UseToken     bool   `toml:"token_authentication"` // If token or user + password authentication should be used
}

type Credentials struct {
	Token    string `toml:"token"`    // For token-based authentication
	Username string `toml:"username"` // For username + password authentication
	Password string `toml:"password"` // For username + password authentication
}

type HomescriptConfig struct {
	// Whether to lint Homescript projects before push
	LintOnPush bool
}

func readConfigFile() {
	configDir, err := os.UserConfigDir()
	if err != nil {
		fmt.Println("Failed to determine user config directory, not reading config file")
		os.Exit(1)
	}
	configFilePath := fmt.Sprintf("%s/%s", configDir, filePath)
	_, err = os.Stat(configFilePath)
	if os.IsNotExist(err) {
		if Verbose {
			fmt.Println("Configuration file does not exist, creating...")
		}
		// Set a default configuration
		Config = Configuration{
			Connection: ConnectionConfig{
				SmarthomeUrl: "http://localhost",
				UseToken:     false,
			},
			Credentials: Credentials{
				Token:    "",
				Username: "",
				Password: "",
			},
			Homescript: HomescriptConfig{
				LintOnPush: true,
			},
		}
		marshaled, err := toml.Marshal(Config)
		if err != nil {
			fmt.Println("Could not create config file: ", err.Error())
			os.Exit(1)
		}
		if err := os.MkdirAll(fmt.Sprintf("%s/%s", configDir, filePathPrefix), 0755); err != nil {
			fmt.Println("Could not create Smarthome-Cli configuration directory: ", err.Error())
			os.Exit(1)
		}
		if err := os.WriteFile(configFilePath, marshaled, 0600); err != nil {
			fmt.Println("Could not create config file: ", err.Error())
			os.Exit(1)
		}
		if Verbose {
			fmt.Printf("Created new configuration at %s\n", configFilePath)
		}
		return
	}
	fileContent, err := os.ReadFile(configFilePath)
	if err != nil {
		fmt.Println("Failed to read configuration file")
		os.Exit(1)
	}
	if err := toml.Unmarshal(fileContent, &Config); err != nil {
		fmt.Printf("Failed to parse configuration file at `%s`: invalid TOML format: %s\n", configFilePath, err.Error())
		os.Exit(1)
	}
	if overrideConfig.Credentials.Username != "" {
		if Verbose {
			fmt.Println("Selected username from flags instead of file.")
		}
		Config.Credentials.Username = overrideConfig.Credentials.Username
	}
	if overrideConfig.Credentials.Password != "" {
		if Verbose {
			fmt.Println("Selected password from flags instead of file.")
		}
		Config.Credentials.Password = overrideConfig.Credentials.Password
	}
	if overrideConfig.Connection.SmarthomeUrl != "" {
		if Verbose {
			fmt.Println("Selected Smarthome URL from flags instead of file.")
		}
		Config.Connection.SmarthomeUrl = overrideConfig.Connection.SmarthomeUrl
	}
	if !overrideConfig.Homescript.LintOnPush {
		if Verbose {
			fmt.Println("Selected lint-on-push from flags instead of file.")
		}
		Config.Homescript.LintOnPush = false
	}
}

func printConfig() {
	readConfigFile()
	configDir, err := os.UserConfigDir()
	if err != nil {
		fmt.Println("Failed to determine user config directory, not reading config file")
		os.Exit(1)
	}
	fmt.Printf("You configuration file is located at `%s/%s`, you can edit it for more settings\n", configDir, filePath)
	headerFmt := color.New(color.FgGreen, color.Underline).SprintfFunc()
	columnFmt := color.New(color.FgYellow).SprintfFunc()
	tbl := table.New("Option", "Value")
	tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)
	tbl.AddRow("Smarthome URL", Config.Connection.SmarthomeUrl)
	// Authentication method
	authMethodString := "username + password"
	if Config.Connection.UseToken {
		authMethodString = "authentication token"
	}
	tbl.AddRow("Authentication Mode", authMethodString)
	// Credential display
	if Config.Connection.UseToken {
		tbl.AddRow("Token", strings.Repeat("*", utf8.RuneCount([]byte(Config.Credentials.Token))))
	} else {
		tbl.AddRow("Username", Config.Credentials.Username)
		tbl.AddRow("Password", strings.Repeat("*", utf8.RuneCount([]byte(Config.Credentials.Password))))
	}
	lintOnPushStr := "yes"
	if !Config.Homescript.LintOnPush {
		lintOnPushStr = "no"
	}
	tbl.AddRow("Lint HMS on push", lintOnPushStr)
	tbl.Print()
}

func writeConfig(newConfig Configuration) {
	fmt.Println("Updating configuration...")
	readConfigFile()
	if !strings.HasPrefix(newConfig.Connection.SmarthomeUrl, "https://") && !strings.HasPrefix(newConfig.Connection.SmarthomeUrl, "http://") {
		newConfig.Connection.SmarthomeUrl = "http://" + newConfig.Connection.SmarthomeUrl
	}
	if _, err := url.Parse(newConfig.Connection.SmarthomeUrl); err != nil {
		fmt.Println("Invalid URL specified: please provide a valid URL.")
		os.Exit(1)
	}
	output, err := toml.Marshal(newConfig)
	if err != nil {
		fmt.Println("Failed to update configuration: could not encode config file", err.Error())
		os.Exit(1)
	}
	configDir, err := os.UserConfigDir()
	if err != nil {
		fmt.Println("Failed to update configuration: could not determine user's config directory")
		os.Exit(1)
	}
	configFilePath := fmt.Sprintf("%s/%s", configDir, filePath)
	fmt.Println("Writing configuration to...", configFilePath)
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
	configFilePath := fmt.Sprintf("%s/%s", configDir, filePath)
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
