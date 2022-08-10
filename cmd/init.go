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
	PromptLogin(false)
	s.Start()
	if (Config.Credentials.Username != "" || Config.Credentials.Password != "") && Config.Connection.UseToken {
		fmt.Println("Warning: username and / or password not empty whilst using token authentication\n=> Ist this intended?")
	}
	if !strings.HasPrefix(Config.Connection.SmarthomeUrl, "https://") && !strings.HasPrefix(Config.Connection.SmarthomeUrl, "http://") {
		fmt.Println("Warning: no URL scheme specified: using insecure HTTP")
		Config.Connection.SmarthomeUrl = "http://" + Config.Connection.SmarthomeUrl
	}
	var conn *sdk.Connection
	var err error
	if Config.Connection.UseToken {
		conn, err = sdk.NewConnection(
			Config.Connection.SmarthomeUrl,
			sdk.AuthMethodCookieToken,
		)
	} else {
		conn, err = sdk.NewConnection(
			Config.Connection.SmarthomeUrl,
			sdk.AuthMethodCookiePassword,
		)
	}
	if err != nil {
		s.FinalMSG = fmt.Sprintf("Could not prepare connection via SDK for Smarthome-server (url: '%s'). Error: %s", Config.Connection.SmarthomeUrl, err.Error())
		s.Stop()
		os.Exit(99)
	}
	Connection = conn
	if Config.Connection.UseToken {
		if Verbose {
			fmt.Println("Note: Using token authentication")
		}
		err = Connection.TokenLogin(Config.Credentials.Token)
	} else {
		if Verbose {
			fmt.Println("Note: Using token password")
		}
		err = Connection.UserLogin(
			Config.Credentials.Username,
			Config.Credentials.Password,
		)
	}
	if err != nil {
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
			s.FinalMSG = fmt.Sprintf("Could not initialize SDK for Smarthome-server (url: '%s').\n  Error: %s\n=> You can revise your local configuration using \x1b[32m'%s config'\x1b[0m\n", Config.Connection.SmarthomeUrl, err.Error(), os.Args[0])
		}
		s.Stop()
		os.Exit(99)
	}
	// Get the username from the connection if token authentication is used
	if Config.Connection.UseToken {
		username, err := Connection.GetUsername()
		if err != nil {
			panic(fmt.Sprintf("Encountered impossible error: %s", err.Error()))
		}
		Config.Credentials.Username = username
		if Verbose {
			fmt.Println("Successfully fetched username after token authentication")
		}
	}
	if Verbose {
		s.FinalMSG = fmt.Sprintf("Successfully connected to '%s' on port %s\n", Connection.SmarthomeURL.Hostname(), Connection.SmarthomeURL.Port())
	}
	s.Stop()
}

// The login function prompts the user to enter their credentials, only used if credentials are not specified beforehand (using config or flags)
func PromptLogin(force bool) {
	if Config.Credentials.Username == "" || force || Config.Connection.UseToken {
		fmt.Printf("\x1b[1;33mAuthentication required\x1b[1;0m: Please enter credentials for `%s`\n", Config.Connection.SmarthomeUrl)
		if Config.Connection.UseToken {
			if Config.Credentials.Token == "" {
				fmt.Printf("Please enter authentication token: (Obtain it here `%s/profile`)\nToken: ", Config.Connection.SmarthomeUrl)
				token, err := gopass.GetPasswd()
				if err != nil {
					fmt.Println("Failed to scan token from STDIN: ", err.Error())
					os.Exit(99)
				}
				Config.Credentials.Token = string(token)
			} else if Verbose {
				fmt.Println("Token already set via flags or config file")
			}
			return
		}
		fmt.Printf("Username: ")
		var username string
		_, err := fmt.Scanln(&username)
		if err != nil {
			fmt.Println("Failed to scan username from STDIN: ", err.Error())
			os.Exit(99)
		}
		Config.Credentials.Username = username
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
	if adminPassword, adminPasswordOk := os.LookupEnv("SMARTHOME_ADMIN_PASSWORD"); adminPasswordOk && Config.Credentials.Username == "admin" && Config.Credentials.Password == "" {
		Config.Credentials.Password = adminPassword
		if Verbose {
			fmt.Printf("Omitting password-prompt: found possible password from \x1b[1;33mSMARTHOME_ADMIN_PASSWORD\x1b[1;0m\n")
		}
	}
	if Config.Credentials.Password == "" || force {
		fmt.Printf("Please enter password for user '%s' in order to continue.\n", Config.Credentials.Username)
		fmt.Printf("Password: ")
		pass, err := gopass.GetPasswd()
		if err != nil {
			fmt.Println("Failed to scan password from STDIN: ", err.Error())
			os.Exit(99)
		}
		Config.Credentials.Password = string(pass)
	} else {
		if Verbose {
			fmt.Println("Password already set from env, args, or config file")
		}
	}
}
