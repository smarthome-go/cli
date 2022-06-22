package cmd

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/briandowns/spinner"
	"github.com/chzyer/readline"

	"github.com/smarthome-go/sdk"
)

var (
	History   []string
	Switches  []sdk.Switch
	completer *readline.PrefixCompleter
)

func filterInput(r rune) (rune, bool) {
	switch r {
	// block CtrlZ feature
	case readline.CharCtrlZ:
		return r, false
	}
	return r, true
}

func initCompleter() {
	switchCompletions := make([]readline.PrefixCompleterInterface, 0)
	for _, switchItem := range Switches {
		switchCompletions = append(switchCompletions,
			readline.PcItem(fmt.Sprintf("('%s', on)", switchItem.Id)),
		)
		switchCompletions = append(switchCompletions,
			readline.PcItem(fmt.Sprintf("('%s', off)", switchItem.Id)),
		)
	}
	completer = readline.NewPrefixCompleter(
		readline.PcItem("switch",
			switchCompletions...,
		),
		readline.PcItem("sleep",
			readline.PcItem("(1)"),
		),
		readline.PcItem("print",
			readline.PcItem("(debugInfo)"),
			readline.PcItem("(weather)"),
			readline.PcItem("(temperature)"),
			readline.PcItem("(user)"),
		),
		readline.PcItem("#exit"),
		readline.PcItem("#switches"),
		readline.PcItem("#power"),
		readline.PcItem("#debug"),
		readline.PcItem("#config"),
		readline.PcItem("#verbose"),
		readline.PcItem("#wipe"),
		readline.PcItem("#reload"),
	)
}

func StartRepl() {
	s := spinner.New(spinner.CharSets[11], 100*time.Millisecond)
	s.Suffix = " Preparing REPL"
	if Verbose {
		fmt.Println("Fetching switches from Smarthome")
		fmt.Println("Fetching debug info from Smarthome")
	}
	s.Start()
	// Fetch the user switches
	switches, err := Connection.GetPersonalSwitches()
	if err != nil {
		defer s.Stop()
		fmt.Println(err.Error())
	}
	Switches = switches

	initCompleter()
	s.Stop()
	fmt.Printf("Welcome to Homescript interactive v%s. CLI commands and comments start with \x1b[90m#\x1b[0m\n", Version)
	fmt.Printf("Server: v%s:%s on \x1b[35m%s\x1b[0m\n", Connection.SmarthomeVersion, Connection.SmarthomeGoVersion, Url)
	cacheDir, err := os.UserCacheDir()
	var historyFile string
	if err != nil {
		fmt.Println("Failed to setup default history, user has no default caching directory, using fallback at `/tmp`")
		historyFile = "/tmp/homescript.history"
	} else {
		historyFile = fmt.Sprintf("%s/homescript.history", cacheDir)
	}
	l, err := readline.NewEx(&readline.Config{
		Prompt:          fmt.Sprintf("\x1b[32m%s\x1b[0m@\x1b[34mhomescript\x1b[0m> ", Username),
		HistoryFile:     historyFile,
		AutoComplete:    completer,
		InterruptPrompt: "^C",
		EOFPrompt:       "exit",

		HistorySearchFold:   true,
		FuncFilterInputRune: filterInput,
	})
	if err != nil {
		panic(err)
	}
	defer l.Close()

	for {
		line, err := l.Readline()
		if err == readline.ErrInterrupt {
			if len(line) == 0 {
				break
			} else {
				continue
			}
		} else if err == io.EOF {
			break
		}
		if strings.ReplaceAll(line, " ", "") == "#exit" {
			os.Exit(0)
		}
		if strings.ReplaceAll(line, " ", "") == "#verbose" {
			Verbose = true
			fmt.Println("Set output mode to verbose")
			continue
		}
		if strings.ReplaceAll(line, " ", "") == "#switches" {
			listSwitches()
			continue
		}
		if strings.ReplaceAll(line, " ", "") == "#power" {
			powerStats()
			continue
		}
		if strings.ReplaceAll(line, " ", "") == "#debug" {
			printDebugInfo()
			continue
		}
		if strings.ReplaceAll(line, " ", "") == "#config" {
			printConfig()
			continue
		}
		if strings.ReplaceAll(line, " ", "") == "#wipe" {
			if Verbose {
				fmt.Println("History has been deleted.")
			}
			l.ResetHistory()
			continue
		}
		if strings.ReplaceAll(line, " ", "") == "#reload" {
			if Verbose {
				fmt.Printf("Reconnecting: (using %s@%s)\n", Connection.Username, Connection.SmarthomeURL.Hostname())
			}
			// Reconnect
			InitConn()

			if Verbose {
				fmt.Println("Updating available switches...")
			}
			// Fetch the user switches again
			switches, err := Connection.GetPersonalSwitches()
			if err != nil {
				fmt.Println(err.Error())
			}
			Switches = switches

			// Generate new autocompletions based on new switches
			initCompleter()
			l.Refresh()

			// Reinitialize readline
			l, err = readline.NewEx(&readline.Config{
				Prompt:          fmt.Sprintf("\x1b[32m%s\x1b[0m@\x1b[34mhomescript\x1b[0m> ", Username),
				HistoryFile:     historyFile,
				AutoComplete:    completer,
				InterruptPrompt: "^C",
				EOFPrompt:       "exit",

				HistorySearchFold:   true,
				FuncFilterInputRune: filterInput,
			})
			if err != nil {
				fmt.Println(err.Error())
			}
			fmt.Println("Session has been reloaded.")
			continue
		}

		if Verbose {
			fmt.Printf("Executing current line. (using %s@%s)\n", Connection.Username, Connection.SmarthomeURL.Hostname())
		}
		startTime := time.Now()
		exitCode := RunCode(line, make(map[string]string, 0), "repl")
		var display string
		if exitCode != 0 {
			display = fmt.Sprintf(" \x1b[31m[%d]\x1b[0m", exitCode)
		}
		l.SetPrompt(fmt.Sprintf("\x1b[32m%s\x1b[0m@\x1b[34mhomescript\x1b[0m%s[\x1b[90m%.2fs\x1b[0m]> ", Username, display, time.Since(startTime).Seconds()))
	}
}
