package homescript

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/MikMuellerDev/homescript-cli/cmd/log"
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

	err := homescript.Run(
		executor, scriptCode,
	)
	if err != nil && len(err) > 0 {
		log.Fatal(fmt.Sprintf("Homescript '%s' ran by user '%s' has terminated:\n%s", scriptLabel, username, err[0].Error()))
		fmt.Printf("Output:\n%s\n", executor.Output)
	}
	log.Info(fmt.Sprintf("Homescript '%s' ran by user '%s' was executed successfully", scriptLabel, username))
}

func RunFile(username string, filename string, serverUrl string, cookies []*http.Cookie) {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal("Failed to read file: ", err.Error())
	}
	Run(username, filename, string(content), serverUrl, cookies)
}
