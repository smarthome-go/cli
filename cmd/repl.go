package cmd

import (
	"fmt"
	"strings"

	"github.com/MikMuellerDev/homescript-cli/cmd/homescript"
	"github.com/MikMuellerDev/homescript-cli/cmd/log"
	"github.com/c-bata/go-prompt"
)

var (
	History  []string
	Switches []Switch
)

func StartRepl() {
	log.Debug("Fetching switches from Smarthome")
	getPersonalSwitches()
	log.Debug("Switches have been successfully fetched")

	fmt.Printf("Welcome to Homescript.\nType exit(0) or CTRL+D to exit.\n")
	p := prompt.New(
		executor,
		completer,
		prompt.OptionPrefix("homescript> "),
		prompt.OptionTitle("Homescript"),
		prompt.OptionHistory(History),
		prompt.OptionSuggestionBGColor(prompt.DarkGray),
		prompt.OptionSelectedSuggestionBGColor(prompt.Blue),
	)
	p.Run()
	fmt.Print("\x1b[3J\033c")
}

func executor(input string) {
	homescript.Run(Username, "repl:01", input, SmarthomeURL, SessionCookies)
}

func completer(d prompt.Document) []prompt.Suggest {
	var suggestions []prompt.Suggest
	if strings.Contains(d.CurrentLineBeforeCursor(), "switch(") {
		for _, switchItem := range Switches {
			if strings.Contains(d.CurrentLineBeforeCursor(), "switch(") && !strings.Contains(d.CurrentLineBeforeCursor(), ",") {
				suggestions = append(suggestions, prompt.Suggest{
					Text:        fmt.Sprintf("'%s', ", switchItem.Id),
					Description: fmt.Sprintf("% 4s | %s", switchItem.Id, switchItem.Name),
				})
			}
		}

		for _, switchItem := range Switches {
			if strings.Contains(d.CurrentLineBeforeCursor(), "switch(") && strings.Contains(d.CurrentLineBeforeCursor(), switchItem.Id) {
				suggestions = append(suggestions, prompt.Suggest{Text: "on)", Description: "ON keyword"})
				suggestions = append(suggestions, prompt.Suggest{Text: "off)", Description: "OFF keyword"})
			}
		}
	} else {
		suggestions = append(suggestions, prompt.Suggest{Text: "switch(", Description: "Turn on / off a switch"})
		suggestions = append(suggestions, prompt.Suggest{Text: "print('')", Description: "Print something to the console"})
		suggestions = append(suggestions, prompt.Suggest{Text: "exit(0)", Description: "Exit the repl"})
	}

	return prompt.FilterHasPrefix(suggestions, d.GetWordBeforeCursor(), true)
}
