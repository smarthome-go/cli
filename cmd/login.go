package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/howeyc/gopass"

	"github.com/MikMuellerDev/homescript-cli/cmd/log"
)

// The login function prompts the user to enter his password and username (if not provided via flag)
func PromptLogin() {
	if Username == "" {
		log.Logn("\x1b[1;33mAuthentication required\x1b[1;0m: Enter username in order to continue.")
		fmt.Printf("Username: ")
		var username string
		_, err := fmt.Scanln(&username)
		if err != nil {
			log.Loge("Failed to scan username: ", err.Error())
		}
		Username = username
	} else {
		if Verbose {
			log.Logn("Username already set from args")
		}
	}
	if Password == "" {
		log.Logn("Please authenticate as user: ", Username)
		fmt.Printf("Password: ")
		pass, err := gopass.GetPasswd()
		if err != nil {
			log.Loge("Failed to scan password: ", err.Error())
		}
		Password = string(pass)
	} else {
		if Verbose {
			log.Logn("Password already set from env / args")
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
		log.Loge("Failed to prepare login request: ", err.Error())
		os.Exit(1)
	}
	res, err := http.Post(
		fmt.Sprintf("%s/api/login", SmarthomeURL),
		"application/json",
		bytes.NewBuffer([]byte(body)),
	)
	if err != nil {
		log.Loge("Login failed: ", err.Error())
		os.Exit(1)
	}
	switch res.StatusCode {
	case 204:
	case 400:
		log.Loge("Failed to login (\x1b[33m400\x1b[0m): invalid request body. Is your homescript client up-to-date?")
		os.Exit(2)
	case 401:
		log.Loge("Failed to login (\x1b[33m401\x1b[0m): invalid credentials")
		os.Exit(3)
	case 500:
		log.Loge("Failed to login (\x1b[31m500\x1b[0m): smarthome server returned an error")
		os.Exit(4)
	case 503:
		log.Loge("Failed to login (\x1b[31m503\x1b[0m): smarthome is currently unavailable. Is the datbase online?")
		os.Exit(5)
	default:
		log.Loge("Failed to login: received unknown status code from smarthome: ", res.Status)
		os.Exit(6)
	}

	if len(res.Cookies()) != 1 {
		log.Loge("Login failed: smarthome returned invalid login cookies")
		os.Exit(1)
	}
	SessionCookies = res.Cookies()
	if Verbose {
		log.Logn("Login successful: you are now authenticated as: ", Username)
	}
}

// Tests if the server configuration is valid and the server is reachable
func PingServer() {
	if Verbose {
		log.Logn(fmt.Sprintf("Pinging Smarthome server at '%s'", SmarthomeURL))
	}
	res, err := http.Get(fmt.Sprintf("%s/health", SmarthomeURL))
	if err != nil {
		log.Loge(fmt.Sprintf("Server ping \x1b[31mfailed\x1b[0m: failed to connect to \x1b[35m%s\x1b[0m. Check your server configuration", SmarthomeURL))
		os.Exit(1)
	}
	switch res.StatusCode {
	case 200:
	case 503:
		log.Loge("Smarthome may be in a degraded state (\x1b[33m503\x1b[0m): some services are currently unavailable")
		return
	default:
		log.Loge(fmt.Sprintf("Server ping failed (\x1b[31m%s\x1b[0m): server responded with unknown response code", res.Status))
		os.Exit(3)
	}
	if Verbose {
		log.Logn("Server ping successful: Smarthome configuration is valid")
	}
}
