package cmd

import (
	"fmt"
	"strings"
	"time"

	"github.com/briandowns/spinner"

	"github.com/MikMuellerDev/smarthome_sdk"
)

// Pretty-prints a Homescript error
func printError(err smarthome_sdk.HomescriptError, program string) {
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

// Executes an arbitrary string of Homescript code
// Error handling is done internally and printed directly
func RunCode(code string, filename string) int {
	s := spinner.New([]string{"⠏", "⠛", "⠹", "⢸", "⣰", "⣤", "⣆", "⡇"}, 100*time.Millisecond)
	s.Prefix = "Executing Homescript "
	s.FinalMSG = ""
	start := time.Now()
	ch := make(chan struct{})
	go func(ch *chan struct{}) {
		for {
			if time.Since(start).Milliseconds() > 200 {
				s.Start()
			}
			select {
			case <-*ch:
				s.Stop()
				return
			default:
			}
		}
	}(&ch)
	output, err := Connection.RunHomescript(code, time.Minute*2)
	ch <- struct{}{}
	if err != nil {
		if err == smarthome_sdk.ErrPermissionDenied {
			fmt.Printf("Permission denied: you \x1b[90m(%s)\x1b[0m do not have the permission \x1b[90m(homescript)\x1b[0m which is required to use Homescript.\n", Connection.Username)
			return 403
		}
		fmt.Println(err.Error())
		return 99
	}
	if !output.Success || output.Exitcode != 0 {
		fmt.Printf("Error: Program terminated abnormally with exit-code %d\n", output.Exitcode)
		for _, errorItem := range output.Errors {
			errorItem.Location.Filename = filename
			printError(errorItem, code)
		}
		return output.Exitcode
	}
	if output.Output != "" {
		fmt.Printf("\x1b[90m%s\x1b[0m\n", output.Output)
	}
	return output.Exitcode
}
