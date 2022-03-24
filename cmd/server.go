package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// Used when starting a REPL session (for autocompletion)
type Switch struct {
	Id     string `json:"id"`
	Name   string `json:"name"`
	RoomId string `json:"roomId"`
}

// Fetches the available user switches from the smarthome server
func getPersonalSwitches() {
	req, err := http.NewRequest(
		"GET",
		fmt.Sprintf("%s/api/switch/list/personal", SmarthomeURL),
		nil,
	)
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
		log.Fatal(fmt.Sprintf("Failed to fetch switches: %s", err.Error()))
	}
	if res.StatusCode > 299 {
		log.Fatal(fmt.Sprintf("Failed to fetch switches: non-200 status code: %s", res.Status))
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(fmt.Sprintf("Failed to fetch switches: could not parse response: %s", res.Status))
	}
	var parsedBody []Switch
	if err := json.Unmarshal(body, &parsedBody); err != nil {
		log.Fatal("Failed to fetch switches: ", err.Error())
	}
	Switches = parsedBody
}
