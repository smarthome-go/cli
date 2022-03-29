package homescript

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"log"
)

type RunRequest struct {
	Code string `json:"code"`
}

type Location struct {
	Filename string `json:"filename"`
	Line     uint   `json:"line"`
	Column   uint   `json:"column"`
	Index    uint   `json:"index"`
}

type HomescriptError struct {
	ErrorType string   `json:"errorType"`
	Location  Location `json:"location"`
	Message   string   `json:"message"`
}

type HomescriptResponse struct {
	Success  bool              `json:"success"`
	Exitcode int               `json:"exitCode"`
	Message  string            `json:"message"`
	Output   string            `json:"output"`
	Errors   []HomescriptError `json:"error"`
}

type ANSICode string

var (
	ANSIRedFg   ANSICode = "\x1b[31m"
	ANSIClearFg ANSICode = "\x1b[0m"
)

func startSpinner(text string, ch *chan bool) {
	positions := []string{"⠏", "⠛", "⠹", "⢸", "⣰", "⣤", "⣆", "⡇"}
	fmt.Println()
	fmt.Printf("\x1b[1F")
	startTime := time.Now()
	for {
		for _, pos := range positions {
			if time.Since(startTime).Milliseconds() > 300 {
				fmt.Printf("%s %s [%.1fs]\x1b[1F", pos, text, time.Since(startTime).Seconds())
				time.Sleep(time.Millisecond * 50)
				fmt.Println()
			}
			select {
			case <-*ch:
				return
			default:
			}
		}
	}
}

func printError(err HomescriptError, program string) {
	lines := strings.Split(program, "\n")

	line1 := ""
	if err.Location.Line > 1 {
		line1 = fmt.Sprintf("\n \x1b[90m%- 3d | \x1b[0m%s", err.Location.Line-1, lines[err.Location.Line-2])
	}
	line2 := fmt.Sprintf(" \x1b[90m%- 3d | \x1b[0m%s", err.Location.Line, lines[err.Location.Line-1])
	line3 := ""
	if int(err.Location.Line) < len(lines) {
		line3 = fmt.Sprintf("\n \x1b[90m%- 3d | \x1b[0m%s", err.Location.Line+1, lines[err.Location.Line])
	}

	marker := fmt.Sprintf("%s\x1b[1;31m^\x1b[0m", strings.Repeat(" ", int(err.Location.Column+6)))

	fmt.Printf(
		"\x1b[1;36m%s\x1b[39m at %s:%d:%d\x1b[0m\n%s\n%s\n%s%s\n\n\x1b[1;31m%s\x1b[0m\n",
		err.ErrorType,
		err.Location.Filename,
		err.Location.Line,
		err.Location.Column,
		line1,
		line2,
		marker,
		line3,
		err.Message,
	)
}

func Run(scriptCode string, serverUrl string, cookies []*http.Cookie) int {
	ch := make(chan bool)
	go startSpinner("Executing homescript", &ch)
	url := fmt.Sprintf("%s/api/homescript/run/live", serverUrl)
	requestBody, err := json.Marshal(RunRequest{
		Code: scriptCode,
	})
	if err != nil {
		ch <- true
		log.Printf("%sError%s: Could not encode request to JSON: %s\n", ANSIRedFg, ANSIClearFg, err.Error())
		return 2
	}
	req, err := http.NewRequest(
		"POST",
		url,
		bytes.NewReader(requestBody),
	)
	if err != nil {
		ch <- true
		log.Printf("%sError%s: Failed to create request from parameters: %s\n", ANSIRedFg, ANSIClearFg, err.Error())
		return 10
	}
	req.Header.Set("Content-Type", "application/json")
	for _, cookie := range cookies {
		req.AddCookie(
			&http.Cookie{
				Name:  cookie.Name,
				Value: cookie.Value,
			},
		)
	}
	client := &http.Client{
		Timeout: 60 * time.Second,
	}
	res, err := client.Do(req)
	if err != nil {
		ch <- true
		log.Printf("%sError%s: Failed to send request to server: %s\n", ANSIRedFg, ANSIClearFg, err.Error())
		return 11
	}
	ch <- true
	defer res.Body.Close()
	switch res.StatusCode {
	case 401:
		log.Printf("Failed to run Homescript: unauthorized\n")
		return 401
	case 200:
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			log.Printf("%sError%s: Failed to display output: could not evaluate server's response: %s\n", ANSIRedFg, ANSIClearFg, err.Error())
			return 12
		}
		var parsedBody HomescriptResponse
		if err := json.Unmarshal(body, &parsedBody); err != nil {
			log.Printf("%sError%s: Failed to parse server's response: %s. body:'%s'\n", ANSIRedFg, ANSIClearFg, err.Error(), body)
			return 13
		}
		if parsedBody.Output != "" {
			log.Printf("Output: %s\n", parsedBody.Output)
		}
		return parsedBody.Exitcode
	case 500:
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			log.Printf("%sError%s: Failed to display error: could not evaluate server's response: %s\n", ANSIRedFg, ANSIClearFg, err.Error())
			return 12
		}
		var parsedBody HomescriptResponse
		if err := json.Unmarshal(body, &parsedBody); err != nil {
			log.Printf("%sError%s: Failed to parse server's response: %s. body:'%s'\n", ANSIRedFg, ANSIClearFg, err.Error(), body)
			return 13
		}
		log.Printf("Homescript error: terminated with exit code: %d\n", parsedBody.Exitcode)
		for _, errorItem := range parsedBody.Errors {
			printError(errorItem, scriptCode)
		}
		return parsedBody.Exitcode
	default:
		log.Printf("%sError%s: Unknown response code from server: %s\n", ANSIRedFg, ANSIClearFg, res.Status)
		return 14
	}
}

func RunFile(filename string, serverUrl string, cookies []*http.Cookie) {
	startTime := time.Now()
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Println("Failed to read file: ", err.Error())
		os.Exit(1)
	}
	exitCode := Run(string(content), serverUrl, cookies)
	fmt.Printf("Homescript finished with exit code: %d \x1b[90m[%ds]%s\n", exitCode, time.Since(startTime).Milliseconds(), ANSIClearFg)
}
