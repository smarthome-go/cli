package workspace

import (
	"fmt"
	"strings"
	"time"

	"github.com/briandowns/spinner"

	"github.com/smarthome-go/sdk"
)

// Pretty-prints a Homescript error
func printError(err sdk.HomescriptError, program string) {
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

// Executes an arbitrary Homescript given its id
// Error handling is done internally and printed directly
func RunById(connection *sdk.Connection, id string, args map[string]string) int {
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
	output, err := connection.RunHomescriptById(id, args, time.Minute*2)
	ch <- struct{}{}
	if err != nil {
		if err == sdk.ErrPermissionDenied {
			username, err := connection.GetUsername()
			if err != nil {
				panic(fmt.Sprintf("Encountered impossible error: %s", err.Error()))
			}
			fmt.Printf("Permission denied: you \x1b[90m(%s)\x1b[0m do not have the permission \x1b[90m(homescript)\x1b[0m which is required to use Homescript.\n", username)
			return 403
		}
		fmt.Println(err.Error())
		return 99
	}
	if !output.Success || output.Exitcode != 0 {
		fmt.Printf("Error: Program terminated abnormally with exit-code %d\n", output.Exitcode)
		// Retrieve remote code in order to pretty-print the error
		remoteData, err := connection.GetHomescript(id)
		if err != nil {
			fmt.Printf("Could not download remote code for error display:\n%s\n", err.Error())
			return 255
		}
		for _, errorItem := range output.Errors {
			errorItem.Location.Filename = fmt.Sprintf("%s.hms", id)
			printError(errorItem, remoteData.Data.Code)
		}
		return output.Exitcode
	}
	if output.Output != "" {
		fmt.Printf("\x1b[90m%s\x1b[0m\n", output.Output)
	}
	return output.Exitcode
}

// Executes an arbitrary string of Homescript code
// Error handling is done internally and printed directly
func RunCode(connection *sdk.Connection, code string, args map[string]string, filename string) int {
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
	output, err := connection.RunHomescriptCode(code, args, time.Minute*2)
	ch <- struct{}{}
	if err != nil {
		username, err := connection.GetUsername()
		if err != nil {
			panic(fmt.Sprintf("Encountered impossible error: %s", err.Error()))
		}
		if err == sdk.ErrPermissionDenied {
			fmt.Printf("Permission denied: you \x1b[90m(%s)\x1b[0m do not have the permission \x1b[90m(homescript)\x1b[0m which is required to use Homescript.\n", username)
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

// Lints an arbitrary Homescript given its id
// Error handling is done internally and printed directly
func LintById(connection *sdk.Connection, id string, args map[string]string) int {
	output, err := connection.LintHomescriptById(id, args, time.Minute)
	if err != nil {
		username, err := connection.GetUsername()
		if err != nil {
			panic(fmt.Sprintf("Encountered impossible error: %s", err.Error()))
		}
		if err == sdk.ErrPermissionDenied {
			fmt.Printf("Permission denied: you \x1b[90m(%s)\x1b[0m do not have the permission \x1b[90m(homescript)\x1b[0m which is required to use Homescript.\n", username)
			return 403
		}
		fmt.Println(err.Error())
		return 99
	}
	if !output.Success || output.Exitcode != 0 {
		fmt.Printf("FAIL: linting discovered problems in '%s.hms':\n", id)
		// Retrieve remote code in order to pretty-print the error
		remoteData, err := connection.GetHomescript(id)
		if err != nil {
			fmt.Printf("Could not download remote code for error display:\n%s\n", err.Error())
			return 255
		}
		for _, errorItem := range output.Errors {
			errorItem.Location.Filename = fmt.Sprintf("%s.hms", id)
			printError(errorItem, remoteData.Data.Code)
		}
		return output.Exitcode
	}
	if output.Output != "" {
		fmt.Printf("\x1b[90m%s\x1b[0m\n", output.Output)
	}
	fmt.Printf("PASS: linting discovered no problems in '%s.hms'\n", id)
	return output.Exitcode
}

// Lints an arbitrary string of Homescript code
// Error handling is done internally and printed directly
func LintCode(connection *sdk.Connection, code string, args map[string]string, filename string) int {
	output, err := connection.LintHomescriptCode(code, args, time.Minute*2)
	if err != nil {
		username, err := connection.GetUsername()
		if err != nil {
			panic(fmt.Sprintf("Encountered impossible error: %s", err.Error()))
		}
		if err == sdk.ErrPermissionDenied {
			fmt.Printf("Permission denied: you \x1b[90m(%s)\x1b[0m do not have the permission \x1b[90m(homescript)\x1b[0m which is required to use Homescript.\n", username)
			return 403
		}
		fmt.Println(err.Error())
		return 99
	}
	if !output.Success || output.Exitcode != 0 {
		fmt.Printf("FAIL: linting discovered problems in '%s':\n", filename)
		for _, errorItem := range output.Errors {
			errorItem.Location.Filename = filename
			printError(errorItem, code)
		}
		return output.Exitcode
	}
	if output.Output != "" {
		fmt.Printf("\x1b[90m%s\x1b[0m\n", output.Output)
	}
	fmt.Printf("PASS: linting discovered no problems in '%s'\n", filename)
	return output.Exitcode
}
