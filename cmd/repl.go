package cmd

import (
	"fmt"
	"log"
	"strings"

	"github.com/MikMuellerDev/homescript-cli/cmd/homescript"
	"github.com/c-bata/go-prompt"
)

var (
	History  []string
	Switches []Switch
)

func StartRepl() {
	if Verbose {
		log.Println("Fetching switches from Smarthome")
	}
	getPersonalSwitches()
	if Verbose {
		log.Println("Switches have been successfully fetched")
	}

	log.Printf("Type exit(0) or CTRL+D to exit.\nWelcome to Homescript.\n")
	p := prompt.New(
		executor,
		completer,
		prompt.OptionPrefix("homescript> "),
		prompt.OptionTitle("Homescript"),
		prompt.OptionHistory(History),
		prompt.OptionSuggestionBGColor(prompt.Black),
		prompt.OptionSelectedSuggestionBGColor(prompt.Blue),
		prompt.OptionCompletionOnDown(),
	)
	p.Run()
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
		suggestions = append(suggestions, prompt.Suggest{Text: "print(debugInfo)", Description: "Print debug information"})
		suggestions = append(suggestions, prompt.Suggest{Text: "print('')", Description: "Print something to the console"})
		suggestions = append(suggestions, prompt.Suggest{Text: "exit(0)", Description: "Exit the repl"})
	}

	// return prompt.FilterHasPrefix(suggestions, d.GetWordBeforeCursor(), true)
	return prompt.FilterContains(suggestions, d.GetWordBeforeCursor(), true)
}
