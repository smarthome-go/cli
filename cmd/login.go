package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"log"

	"github.com/howeyc/gopass"
)

// The login function prompts the user to enter his password and username (if not provided via flag)
func PromptLogin() {
	if Username == "" {
		log.Println("\x1b[1;33mAuthentication required\x1b[1;0m: Enter username in order to continue.")
		fmt.Printf("Username: ")
		var username string
		_, err := fmt.Scanln(&username)
		if err != nil {
			loge("Failed to scan username: ", err.Error())
		}
		Username = username
	} else {
		if Verbose {
			logn("Username already set from args")
		}
	}
	if Password == "" {
		logn("Please authenticate as user: ", Username)
		fmt.Printf("Password: ")
		pass, err := gopass.GetPasswd()
		if err != nil {
			loge("Failed to scan password: ", err.Error())
		}
		Password = string(pass)
	} else {
		if Verbose {
			logn("Password already set from env / args")
		}
	}
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Tests if the provided credentials are valid
func Login(doOutput bool) {
	body, err := json.Marshal(
		LoginRequest{
			Username: Username,
			Password: Password,
		})
	if err != nil {
		loge("Failed to prepare login request: ", err.Error())
		os.Exit(1)
	}
	res, err := http.Post(
		fmt.Sprintf("%s/api/login", SmarthomeURL),
		"application/json",
		bytes.NewBuffer([]byte(body)),
	)
	if err != nil {
		loge("Login failed: ", err.Error())
		os.Exit(1)
	}
	switch res.StatusCode {
	case 204:
	case 400:
		loge("Failed to login (\x1b[33m400\x1b[0m): invalid request body. Is your homescript client up-to-date?")
		os.Exit(2)
	case 401:
		loge("Failed to login (\x1b[33m401\x1b[0m): invalid credentials")
		os.Exit(3)
	case 500:
		loge("Failed to login (\x1b[31m500\x1b[0m): smarthome server returned an error")
		os.Exit(4)
	case 503:
		loge("Failed to login (\x1b[31m503\x1b[0m): smarthome is currently unavailable. Is the datbase online?")
		os.Exit(5)
	default:
		loge("Failed to login: received unknown status code from smarthome: ", res.Status)
		os.Exit(6)
	}

	if len(res.Cookies()) != 1 {
		loge("Login failed: smarthome returned invalid login cookies")
		os.Exit(1)
	}
	SessionCookies = res.Cookies()
	if Verbose {
		logn("Login successful: you are now authenticated as: ", Username)
	}
}
