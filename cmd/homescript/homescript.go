package homescript

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"log"

	"github.com/MikMuellerDev/homescript/homescript"
)

// Executes a given homescript as a given user, returns the output and a possible error
func Run(username string, scriptLabel string, scriptCode string, serverUrl string, cookies []*http.Cookie) {
	executor := &Executor{
		Username:       username,
		ScriptName:     scriptLabel,
		ServerUrl:      serverUrl,
		SessionCookies: cookies,
	}
	_, err := homescript.Run(
		executor, scriptLabel, scriptCode,
	)
	if err != nil && len(err) > 0 {
		// TODO: do proper error handling
		log.Println(fmt.Sprintf("Homescript '%s' has terminated:\n\x1b[1;0m%s", scriptLabel, err[0].Message))
	}
}

func RunFile(username string, filename string, serverUrl string, cookies []*http.Cookie) {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Println("Failed to read file: ", err.Error())
	}
	Run(username, filename, string(content), serverUrl, cookies)
}
