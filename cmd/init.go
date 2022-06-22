package cmd

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/briandowns/spinner"
	"github.com/howeyc/gopass"
	"github.com/smarthome-go/sdk"
)

func InitConn() {
	s := spinner.New(spinner.CharSets[59], 150*time.Millisecond)
	s.Prefix = "Connecting to Smarthome "
	PromptLogin()
	s.Start()
	if !strings.HasPrefix(Url, "https://") && !strings.HasPrefix(Url, "http://") {
		fmt.Println("Warning: no URL scheme specified: using insecure HTTP")
		Url = "http://" + Url
	}
	conn, err := sdk.NewConnection(Url, sdk.AuthMethodCookie)
	if err != nil {
		s.FinalMSG = fmt.Sprintf("Could not prepare connection via SDK for Smarthome-server (url: '%s'). Error: %s", Url, err.Error())
		s.Stop()
		os.Exit(99)
	}
	Connection = conn
	if err := Connection.Connect(Username, Password); err != nil {
		if err == sdk.ErrInvalidVersion {
			// The Server is not compatible with the current client
			// TODO: write a list of supported versions in README
			s.FinalMSG = fmt.Sprintf("Could not initialize SDK for Smarthome-server. Incompatible Server version: This client (v%s) requires minimal server version %s but the server is running %s.\nYou can try upgrading your server or downgrading this client.\n", sdk.Version, sdk.MinSmarthomeVersion, Connection.SmarthomeVersion)
		}
		s.FinalMSG = fmt.Sprintf("Could not initialize SDK for Smarthome-server (url: '%s'). Error: %s\nYou can validate you local configuration parameters using \x1b[32m'homescript config'\x1b[0m\n", Url, err.Error())
		s.Stop()
		os.Exit(99)
	}
	if Verbose {
		s.FinalMSG = fmt.Sprintf("Successfully connected to '%s' on port %s\n", Connection.SmarthomeURL.Hostname(), Connection.SmarthomeURL.Port())
	}
	s.Stop()
}

// The login function prompts the user to enter their credentials, only used if credentials are not specified beforehand (using config or flags)
func PromptLogin() {
	if Username == "" {
		fmt.Println("\x1b[1;33mAuthentication required\x1b[1;0m: Please enter your username in order to continue.")
		fmt.Printf("Username: ")
		var username string
		_, err := fmt.Scanln(&username)
		if err != nil {
			fmt.Println("Failed to scan username from STDIN: ", err.Error())
			os.Exit(99)
		}
		Username = username
	} else {
		if Verbose {
			fmt.Println("Username already set via flags or configuration file, not prompting")
		}
	}
	/*
		`SMARTHOME_ADMIN_PASSWORD`: Checks for the smarthome admin user
	*/
	// Uses the admin-password environment variable if the user is `admin` and has a omitted password
	// Only used inside the Smarthome Docker-container on initial setup because the environment variable is only used on the first start of the container
	if adminPassword, adminPasswordOk := os.LookupEnv("SMARTHOME_ADMIN_PASSWORD"); adminPasswordOk && Username == "admin" && Password == "" {
		Password = adminPassword
		if Verbose {
			fmt.Printf("Omitting password-prompt: found possible password from \x1b[1;33mSMARTHOME_ADMIN_PASSWORD\x1b[1;0m\n")
		}
	}
	if Password == "" {
		fmt.Printf("Please enter password for user '%s' in order to continue.\n", Username)
		fmt.Printf("Password: ")
		pass, err := gopass.GetPasswd()
		if err != nil {
			fmt.Println("Failed to scan password from STDIN: ", err.Error())
			os.Exit(99)
		}
		Password = string(pass)
	} else {
		if Verbose {
			fmt.Println("Password already set from env, args, or config file")
		}
	}
}
