package cmd

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/Masterminds/semver"
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
		if err == sdk.ErrUnsupportedVersion {
			// The Server is not compatible with the current client
			s.FinalMSG = fmt.Sprintf("Could not establish connection to unsupported server.\nThis client (v%s) requires minimal server version '%s' but is using '%s'.\n", sdk.Version, sdk.MinSmarthomeVersion, Connection.SmarthomeVersion)
			serverV, err := semver.NewVersion(Connection.SmarthomeVersion)
			if err != nil {
				// These errors usually can't happen due to SDK validation
				s.FinalMSG = fmt.Sprintf("Invalid SemVer version of server: %s\n", err.Error())
				s.Stop()
				os.Exit(99)
			}
			supportV, err2 := semver.NewVersion(sdk.MinSmarthomeVersion)
			if err2 != nil {
				// These errors usually can't happen due to SDK validation
				s.FinalMSG = fmt.Sprintf("Invalid SemVer version of SDK requirement: %s\n", err.Error())
				s.Stop()
				os.Exit(99)
			}
			if serverV.Major() > supportV.Major() {
				s.FinalMSG += fmt.Sprintf("The supported major version has been superseded.\n  Required: %10s [deprecated]\n  Server:   %10s [current]\n=> Try installing the current version of the CLI.\n", "v"+sdk.MinSmarthomeVersion, "v"+Connection.SmarthomeVersion)
			} else if serverV.LessThan(supportV) {
				s.FinalMSG += fmt.Sprintf("The server is outdated.\n  Required: %10s [current]\n  Server:   %10s [deprecated]\n=> Try installing the current version of the server.\n", "v"+sdk.MinSmarthomeVersion, "v"+Connection.SmarthomeVersion)
			}
		} else {
			s.FinalMSG = fmt.Sprintf("Could not initialize SDK for Smarthome-server (url: '%s').\n  Error: %s\n=> You can revise your local configuration using \x1b[32m'%s config'\x1b[0m\n", Url, err.Error(), os.Args[0])
		}
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
