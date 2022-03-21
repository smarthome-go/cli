package homescript

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/MikMuellerDev/homescript-cli/cmd/log"
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
	log.Info(fmt.Sprintf("[Homescript] script: '%s' user: '%s': %s", self.ScriptName, self.Username, output))
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
		log.Error(fmt.Sprintf("[Homescript] ERROR: script: '%s' user: '%s': failed to set power: %s", self.ScriptName, self.Username, err.Error()))
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
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Error(fmt.Sprintf("[Homescript] ERROR: script: '%s' user: '%s': failed to set power: %s", self.ScriptName, self.Username, err.Error()))
	}
	if res.StatusCode > 299 {
		log.Error(fmt.Sprintf("[Homescript] ERROR: script: '%s' user: '%s': failed to set power: %s", self.ScriptName, self.Username, res.Status))
	}
	defer res.Body.Close()
	onOffText := "on"
	if !powerOn {
		onOffText = "off"
	}
	log.Debug(fmt.Sprintf("[Homescript] script: '%s' user: '%s': turning switch %s %s", self.ScriptName, self.Username, switchId, onOffText))
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
		log.Trace(title, description)
	case 1:
		log.Debug(title, description)
	case 2:
		log.Info(title, description)
	case 3:
		log.Warn(title, description)
	case 4:
		log.Error(title, description)
	case 5:
		log.Fatal(title, description)
	default:
		log.Error(fmt.Sprintf("[Homescript] ERROR: script: '%s' user: '%s': failed to log event: invalid level", self.ScriptName, self.Username))
	}
	return nil
}

func (self *Executor) GetUser() string {
	return self.Username
}

func (self *Executor) GetWeather() (string, error) {
	log.Error(fmt.Sprintf("[Homescript] ERROR: script: '%s' user: '%s': weather is not implemented yet", self.ScriptName, self.Username))
	return "rainy", nil
}

func (self *Executor) GetTemperature() (int, error) {
	log.Error(fmt.Sprintf("[Homescript] ERROR: script: '%s' user: '%s': temperature is not implemented yet", self.ScriptName, self.Username))
	return 42, nil
}

func (self *Executor) GetDate() (int, int, int, int, int, int) {
	now := time.Now()
	return now.Year(), int(now.Month()), now.Day(), now.Hour(), now.Minute(), now.Second()
}
