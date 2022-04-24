package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/MikMuellerDev/homescript-cli/cmd/log"
)

// Used when starting a REPL session (for autocompletion)
type Switch struct {
	Id      string `json:"id"`
	Name    string `json:"name"`
	RoomId  string `json:"roomId"`
	PowerOn bool   `json:"powerOn"`
	Watts   uint   `json:"watts"`
}

// Fetches the available user switches from the smarthome server
func getPersonalSwitches() {
	req, err := http.NewRequest(
		"GET",
		fmt.Sprintf("%s/api/switch/list/personal", SmarthomeURL),
		nil,
	)
	if err != nil {
		log.Loge("Failed to fetch switches: could not create request: ", err.Error())
		os.Exit(1)
	}
	for _, cookie := range SessionCookies {
		req.AddCookie(
			&http.Cookie{
				Name:  cookie.Name,
				Value: cookie.Value,
			},
		)
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Loge(fmt.Sprintf("Failed to fetch switches: %s", err.Error()))
		os.Exit(1)
	}

	switch res.StatusCode {
	case 200:
	case 400:
		log.Loge("Failed to fetch switches (\x1b[33m400\x1b[0m): invalid request body. Is your homescript client up-to-date?")
		os.Exit(2)
	case 401:
		log.Loge("Failed to fetch switches (\x1b[33m401\x1b[0m): invalid credentials")
		os.Exit(3)
	case 500:
		log.Loge("Failed to fetch switches (\x1b[31m500\x1b[0m): smarthome server returned an error")
		os.Exit(4)
	case 503:
		log.Loge("Failed to fetch switches (\x1b[31m503\x1b[0m): smarthome is currently unavailable. Is the datbase online?")
		os.Exit(5)
	default:
		log.Loge("Failed to fetch switches: received unknown status code from smarthome: ", res.Status)
		os.Exit(6)
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Loge(fmt.Sprintf("Failed to fetch switches: could not parse response: %s", res.Status))
		os.Exit(1)
	}
	var parsedBody []Switch
	if err := json.Unmarshal(body, &parsedBody); err != nil {
		log.Loge("Failed to fetch switches: ", err.Error())
		os.Exit(1)
	}
	Switches = parsedBody
}
