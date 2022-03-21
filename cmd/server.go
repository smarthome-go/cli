package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/MikMuellerDev/homescript-cli/cmd/log"
)

var (
	ServerInfo DebugInfo
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
		fmt.Sprintf("%s/api/power/list/personal", SmarthomeURL),
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

// Debug info is shown at start when using -v
type DBStatus struct {
	OpenConnections int `json:"openConnections"`
	InUse           int `json:""`
	Idle            int `json:""`
}

type PowerJob struct {
	Id         int64  `json:"id"`
	SwitchName string `json:"switchName"`
	Power      bool   `json:"power"`
}

type JobResult struct {
	Id    int64 `json:"id"`
	Error error `json:"error"`
}

type DebugInfo struct {
	ServerVersion          string      `json:"version"`
	DatabaseOnline         bool        `json:"databaseOnline"`
	DatabaseStats          DBStatus    `json:"databaseStats"`
	CpuCores               uint8       `json:"cpuCores"`
	Goroutines             uint16      `json:"goroutines"`
	GoVersion              string      `json:"goVersion"`
	MemoryUsage            uint16      `json:"memoryUsage"`
	PowerJobCount          uint16      `json:"powerJobCount"`
	PowerJobWithErrorCount uint16      `json:"lastPowerJobErrorCount"`
	PowerJobs              []PowerJob  `json:"powerJobs"`
	PowerJobResults        []JobResult `json:"powerJobResults"`
}

// Fetches debug information from the smarthome server
func getDebugInfo() {
	req, err := http.NewRequest(
		"GET",
		fmt.Sprintf("%s/api/debug", SmarthomeURL),
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
		log.Fatal(fmt.Sprintf("Failed to fetch debug info: %s", err.Error()))
	}
	if res.StatusCode > 299 {
		log.Fatal(fmt.Sprintf("Failed to debug info: non-200 status code %s", res.Status))
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(fmt.Sprintf("Failed to fetch debug info: could not parse response: %s", res.Status))
	}
	var parsedBody DebugInfo
	if err := json.Unmarshal(body, &parsedBody); err != nil {
		log.Fatal("Failed to fetch debug info: ", err.Error())
	}
	log.Debug("Successfully fetched debug info")
	ServerInfo = parsedBody
}

func DisplayServerInfo() {
	if Username != "admin" || !ShowInfo {
		return
	}
	getDebugInfo()

	var output string

	output += fmt.Sprintf("%s\n", strings.Repeat("\u2015", 45))
	output += fmt.Sprintf(" Smarthome Server Version: %s │ v%s\n", strings.Repeat(" ", 30-len("Smarthome Server Version: ")), ServerInfo.ServerVersion)
	var databaseOnlineString = "\x1b[1;31mNO\x1b[1;0m"
	if ServerInfo.DatabaseOnline {
		databaseOnlineString = "\x1b[1;32mYES\x1b[1;0m"
	}
	output += fmt.Sprintf(" Database Online: %s │ %- 10s\n", strings.Repeat(" ", 30-len("Database Online: ")), databaseOnlineString)
	output += fmt.Sprintf(" Compiled with: %s │ %- 10s\n", strings.Repeat(" ", 30-len("Compiled with: ")), ServerInfo.GoVersion)
	output += fmt.Sprintf(" CPU Cores: %s │ %d\n", strings.Repeat(" ", 30-len("CPU Cores: ")), ServerInfo.CpuCores)
	output += fmt.Sprintf(" Current Goroutines: %s │ %d\n", strings.Repeat(" ", 30-len("Current Goroutines: ")), ServerInfo.Goroutines)
	output += fmt.Sprintf(" Current Memory Usage: %s │ %d\n", strings.Repeat(" ", 30-len("Current Memory Usage: ")), ServerInfo.MemoryUsage)
	output += fmt.Sprintf(" Current Power Jobs: %s │ %d\n", strings.Repeat(" ", 30-len("Current Power Jobs: ")), ServerInfo.PowerJobCount)
	output += fmt.Sprintf(" Last Power Job Error Count: %s │ %d\n", strings.Repeat(" ", 30-len("Last Power Job Error Count: ")), ServerInfo.PowerJobWithErrorCount)
	output += fmt.Sprintf("%s", strings.Repeat("\u2015", 45))
	fmt.Println(output)
}
