package cmd

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/chzyer/readline"

	"github.com/MikMuellerDev/homescript-cli/cmd/homescript"
	"github.com/MikMuellerDev/homescript-cli/cmd/log"
)

var (
	History   []string
	Switches  []Switch
	DebugInfo DebugInfoData
)

var completer *readline.PrefixCompleter

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
		readline.PcItem("#verbose"),
		readline.PcItem("#switches"),
		readline.PcItem("#debug"),
	)
}

func StartRepl() {
	if Verbose {
		log.Logn("Fetching switches from Smarthome")
		log.Logn("Fetching server info from Smarthome")
	}
	GetDebugInfo()
	initCompleter()
	log.Logn(fmt.Sprintf("Server: v%s:%s on \x1b[35m%s\x1b[0m", DebugInfo.ServerVersion, DebugInfo.GoVersion, SmarthomeURL), fmt.Sprintf("\nWelcome to Homescript interactive v%s. CLI commands and comments start with \x1b[90m#\x1b[0m", Version))
	cacheDir, err := os.UserCacheDir()
	var historyFile string
	if err != nil {
		log.Loge("Failed to setup default history, user has no default caching directory, using fallback at `/tmp`")
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
			log.InitLog(true)
			log.Logn("Set output mode to verbose")
			continue
		}
		if strings.ReplaceAll(line, " ", "") == "#switches" {
			listSwitches()
			continue
		}
		if strings.ReplaceAll(line, " ", "") == "#debug" {
			debugInfo()
			continue
		}

		startTime := time.Now()
		exitCode := homescript.Run(line, SmarthomeURL, SessionCookies)
		var display string
		if exitCode != 0 {
			display = fmt.Sprintf(" \x1b[31m[%d]\x1b[0m", exitCode)
		}
		l.SetPrompt(fmt.Sprintf("\x1b[32m%s\x1b[0m@\x1b[34mhomescript\x1b[0m%s[\x1b[90m%.2fs\x1b[0m]> ", Username, display, time.Since(startTime).Seconds()))
	}
}
