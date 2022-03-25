package cmd

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"time"

	"github.com/MikMuellerDev/homescript-cli/cmd/debug"
	"github.com/MikMuellerDev/homescript-cli/cmd/homescript"
	"github.com/chzyer/readline"
)

var (
	History   []string
	Switches  []Switch
	DebugInfo debug.DebugInfo
)

func usage(w io.Writer) {
	io.WriteString(w, "commands:\n")
	io.WriteString(w, completer.Tree("    "))
}

// Function constructor - constructs new function for listing given directory
func listFiles(path string) func(string) []string {
	return func(line string) []string {
		names := make([]string, 0)
		files, _ := ioutil.ReadDir(path)
		for _, f := range files {
			names = append(names, f.Name())
		}
		return names
	}
}

var completer = readline.NewPrefixCompleter(
	readline.PcItem("switch",
		readline.PcItem("('', off)"),
	),
	readline.PcItem("print",
		readline.PcItem("(debugInfo)"),
	),
	readline.PcItem("exit",
		readline.PcItem("(1)"),
	),
)

func filterInput(r rune) (rune, bool) {
	switch r {
	// block CtrlZ feature
	case readline.CharCtrlZ:
		return r, false
	}
	return r, true
}

func StartRepl() {
	if Verbose {
		log.Println("Fetching switches from Smarthome")
		log.Println("Fetching server info from Smarthome")
	}
	getPersonalSwitches()
	serverInfo, err := debug.GetDebugInfo(SmarthomeURL, SessionCookies)
	if err != nil {
		log.Println(err.Error())
	}
	DebugInfo = serverInfo
	if Verbose {
		log.Println("Switches have been successfully fetched")
	}
	fmt.Println(fmt.Sprintf("Server: v%s:%s on \x1b[35m%s\x1b[0m", DebugInfo.ServerVersion, DebugInfo.GoVersion, SmarthomeURL), "\nWelcome to Homescript interactive.")
	l, err := readline.NewEx(&readline.Config{
		Prompt:          fmt.Sprintf("\x1b[32m%s\x1b[0m@\x1b[34mhomescript\x1b[0m> ", Username),
		HistoryFile:     "/tmp/homescript_history",
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

	log.SetOutput(l.Stderr())
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

		startTime := time.Now()
		exitCode := homescript.Run(line, SmarthomeURL, SessionCookies)
		var display string
		if exitCode != 0 {
			display = fmt.Sprintf(" \x1b[31m[%d]\x1b[0m", exitCode)
		}
		l.SetPrompt(fmt.Sprintf("\x1b[32m%s\x1b[0m@\x1b[34mhomescript\x1b[0m%s[\x1b[90m%.2fs\x1b[0m]> ", Username, display, time.Since(startTime).Seconds()))
	}
}
