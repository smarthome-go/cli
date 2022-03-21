package debug

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/MikMuellerDev/homescript-cli/cmd/log"
)

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
func getDebugInfo(url string, cookies []*http.Cookie) (DebugInfo, error) {
	req, err := http.NewRequest(
		"GET",
		fmt.Sprintf("%s/api/debug", url),
		nil,
	)
	for _, cookie := range cookies {
		req.AddCookie(
			&http.Cookie{
				Name:  cookie.Name,
				Value: cookie.Value,
			},
		)
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Error(fmt.Sprintf("Failed to fetch debug info: %s", err.Error()))
		return DebugInfo{}, err
	}
	if res.StatusCode > 299 {
		log.Error(fmt.Sprintf("Failed to debug info: non-200 status code %s", res.Status))
		return DebugInfo{}, err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Error(fmt.Sprintf("Failed to fetch debug info: could not parse response: %s", res.Status))
		return DebugInfo{}, err
	}
	var parsedBody DebugInfo
	if err := json.Unmarshal(body, &parsedBody); err != nil {
		log.Error("Failed to fetch debug info: ", err.Error())
		return DebugInfo{}, err
	}
	log.Debug("Successfully fetched debug info")
	return parsedBody, nil
}

func GetServerInfo(url string, cookies []*http.Cookie) (string, error) {
	debugInfo, err := getDebugInfo(url, cookies)
	if err != nil {
		return "", err
	}
	var output string
	output += fmt.Sprintf("%s\n", strings.Repeat("\u2015", 45))
	output += fmt.Sprintf(" Smarthome Server Version: %s │ v%s\n", strings.Repeat(" ", 30-len("Smarthome Server Version: ")), debugInfo.ServerVersion)
	var databaseOnlineString = "\x1b[1;31mNO\x1b[1;0m"
	if debugInfo.DatabaseOnline {
		databaseOnlineString = "\x1b[1;32mYES\x1b[1;0m"
	}
	output += fmt.Sprintf(" Database Online: %s │ %- 10s\n", strings.Repeat(" ", 30-len("Database Online: ")), databaseOnlineString)
	output += fmt.Sprintf(" Compiled with: %s │ %- 10s\n", strings.Repeat(" ", 30-len("Compiled with: ")), debugInfo.GoVersion)
	output += fmt.Sprintf(" CPU Cores: %s │ %d\n", strings.Repeat(" ", 30-len("CPU Cores: ")), debugInfo.CpuCores)
	output += fmt.Sprintf(" Current Goroutines: %s │ %d\n", strings.Repeat(" ", 30-len("Current Goroutines: ")), debugInfo.Goroutines)
	output += fmt.Sprintf(" Current Memory Usage: %s │ %d\n", strings.Repeat(" ", 30-len("Current Memory Usage: ")), debugInfo.MemoryUsage)
	output += fmt.Sprintf(" Current Power Jobs: %s │ %d\n", strings.Repeat(" ", 30-len("Current Power Jobs: ")), debugInfo.PowerJobCount)
	output += fmt.Sprintf(" Last Power Job Error Count: %s │ %d\n", strings.Repeat(" ", 30-len("Last Power Job Error Count: ")), debugInfo.PowerJobWithErrorCount)
	output += fmt.Sprintf("%s", strings.Repeat("\u2015", 45))
	return output, nil
}
