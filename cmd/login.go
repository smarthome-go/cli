package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/MikMuellerDev/homescript-cli/cmd/log"
	"github.com/howeyc/gopass"
)

// The login function prompts the user to enter his password and username (if not provided via flag)
func PromptLogin() {
	if Username == "" {
		log.Info("\x1b[1;33mAuthentication required\x1b[1;0m: Enter your username below.")
		fmt.Printf("Username: ")
		var username string
		_, err := fmt.Scanln(&username)
		if err != nil {
			log.Fatal("Failed to scan username: ", err.Error())
		}
		Username = username
	} else {
		log.Debug("Username already set")
	}
	if Password == "" {
		log.Info("Please authenticate for user: ", Username)
		fmt.Printf("Password: ")
		pass, err := gopass.GetPasswd()
		if err != nil {
			log.Fatal("Failed to scan password: ", err.Error())
		}
		Password = string(pass)
	} else {
		log.Debug("Password already set")
	}
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Tests if the provided credentials are valid
func Login() {
	body, err := json.Marshal(
		LoginRequest{
			Username: Username,
			Password: Password,
		})
	if err != nil {
		log.Fatal("Failed to prepare login request: ", err.Error())
	}
	res, err := http.Post(
		fmt.Sprintf("%s/api/login", SmarthomeURL),
		"application/json",
		bytes.NewBuffer([]byte(body)),
	)
	if err != nil {
		log.Fatal("Failed to login: ", err.Error())
	}
	if res.StatusCode != 204 {
		log.Fatal("Failed to login: response does not indicate success: ", res.Status)
	}
	if len(res.Cookies()) != 1 {
		log.Fatal("Failed to login: invalid cookies")
	}
	SessionCookies = res.Cookies()
	log.Debug("Login successful: you are now authenticated as: ", Username)
}
