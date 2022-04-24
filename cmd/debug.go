package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
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

type DebugInfoData struct {
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
func GetDebugInfo() error {
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
		log.Println(fmt.Sprintf("Failed to fetch debug info: %s", err.Error()))
		DebugInfo = DebugInfoData{
			ServerVersion: "_unknown",
			GoVersion:     "go_unknown",
		}
		return err
	}
	if res.StatusCode > 299 {
		log.Println(fmt.Sprintf("Failed to fetch debug info: non-200 status code %s", res.Status))
		DebugInfo = DebugInfoData{
			ServerVersion: "_unknown",
			GoVersion:     "go_unknown",
		}
		return err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println(fmt.Sprintf("Failed to fetch debug info: could not parse response: %s", res.Status))
		DebugInfo = DebugInfoData{
			ServerVersion: "_unknown",
			GoVersion:     "go_unknown",
		}
		return err
	}
	var parsedBody DebugInfoData
	if err := json.Unmarshal(body, &parsedBody); err != nil {
		log.Println("Failed to fetch debug info: ", err.Error())
		DebugInfo = DebugInfoData{
			ServerVersion: "_unknown",
			GoVersion:     "go_unknown",
		}
		return err
	}
	DebugInfo = parsedBody
	return nil
}
