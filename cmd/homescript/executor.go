package homescript

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"log"

	"github.com/MikMuellerDev/homescript-cli/cmd/debug"
	"github.com/MikMuellerDev/homescript/homescript/interpreter"
)

type Executor struct {
	ScriptName     string
	Username       string
	ServerUrl      string
	SessionCookies []*http.Cookie
	Output         string
}

type PowerRequest struct {
	Switch  string `json:"switch"`
	PowerOn bool   `json:"powerOn"`
}

func (self *Executor) Exit(code int) {
	os.Exit(code)
}

func (self *Executor) Print(args ...string) {
	var output string
	for _, arg := range args {
		output += arg
	}
	self.Output += output
	log.Println(output)
}

func (self *Executor) SwitchOn(switchId string) (bool, error) {
	return false, nil
}

func (self *Executor) Switch(switchId string, powerOn bool) error {
	body, err := json.Marshal(PowerRequest{
		Switch:  switchId,
		PowerOn: powerOn,
	})
	if err != nil {
		log.Println(fmt.Sprintf("[Homescript] ERROR: script: '%s' user: '%s': failed to set power: %s", self.ScriptName, self.Username, err.Error()))
		return err
	}
	req, err := http.NewRequest(
		"POST",
		fmt.Sprintf("%s/api/power/set", self.ServerUrl),
		bytes.NewBuffer([]byte(body)),
	)
	req.Header.Add("Content-Type", "application/json")
	for _, cookie := range self.SessionCookies {
		req.AddCookie(
			&http.Cookie{
				Name:  cookie.Name,
				Value: cookie.Value,
			},
		)
	}
	if err != nil {
		return err
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println(fmt.Sprintf("[Homescript] ERROR: script: '%s' user: '%s': failed to set power: %s", self.ScriptName, self.Username, err.Error()))
		return err
	}
	if res.StatusCode > 299 {
		log.Println(fmt.Sprintf("[Homescript] ERROR: script: '%s' user: '%s': failed to set power: %s", self.ScriptName, self.Username, res.Status))
		return err
	}
	defer res.Body.Close()
	return nil
}

func (self *Executor) Play(server string, mode string) error {
	return errors.New("The feature 'radiGo' is not yet implemented")
}

func (self *Executor) Notify(
	title string,
	description string,
	level interpreter.LogLevel,
) error {
	return nil
}

func (self *Executor) Log(
	title string,
	description string,
	level interpreter.LogLevel,
) error {
	switch level {
	case 0:
		log.Println(title, description)
	case 1:
		log.Println(title, description)
	case 2:
		log.Println(title, description)
	case 3:
		log.Println(title, description)
	case 4:
		log.Println(title, description)
	case 5:
		log.Println(title, description)
	default:
		log.Println(fmt.Sprintf("[Homescript] ERROR: script: '%s' user: '%s': failed to log event: invalid level", self.ScriptName, self.Username))
	}
	return nil
}

func (self *Executor) GetUser() string {
	return self.Username
}

func (self *Executor) GetWeather() (string, error) {
	log.Println(fmt.Sprintf("[Homescript] ERROR: script: '%s' user: '%s': weather is not implemented yet", self.ScriptName, self.Username))
	return "rainy", nil
}

func (self *Executor) GetTemperature() (int, error) {
	log.Println(fmt.Sprintf("[Homescript] ERROR: script: '%s' user: '%s': temperature is not implemented yet", self.ScriptName, self.Username))
	return 42, nil
}

func (self *Executor) GetDate() (int, int, int, int, int, int) {
	now := time.Now()
	return now.Year(), int(now.Month()), now.Day(), now.Hour(), now.Minute(), now.Second()
}

func (self *Executor) GetDebugInfo() (string, error) {
	debugInfo, err := debug.GetServerInfo(self.ServerUrl, self.SessionCookies)
	if err != nil {
		log.Println(fmt.Sprintf("[Homescript] ERROR: script: '%s' user: '%s': could not get debug info: %s", self.ScriptName, self.Username, err.Error()))
		return "", err
	}
	return "\n" + debugInfo, nil
}
