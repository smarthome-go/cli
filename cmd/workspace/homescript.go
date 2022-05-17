package workspace

import (
	"github.com/smarthome-go/homescript/homescript"
	hmsError "github.com/smarthome-go/homescript/homescript/error"
)

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

func convertError(errorItem hmsError.Error) HomescriptError {
	return HomescriptError{
		ErrorType: errorItem.TypeName,
		Location: Location{
			Filename: errorItem.Location.Filename,
			Line:     errorItem.Location.Line,
			Column:   errorItem.Location.Column,
			Index:    errorItem.Location.Index,
		},
		Message: errorItem.Message,
	}
}

func convertErrors(errorItems ...hmsError.Error) []HomescriptError {
	var outputErrors []HomescriptError
	for _, errorItem := range errorItems {
		outputErrors = append(outputErrors, convertError(errorItem))
	}
	return outputErrors
}

// Executes a given homescript as a given user, returns the output and a possible error slice
func Run(username string, scriptLabel string, scriptCode string) (string, int, []HomescriptError) {
	executor := &Executor{
		Username:   username,
		ScriptName: scriptLabel,
	}
	exitCode, runtimeErrors := homescript.Run(
		executor,
		scriptLabel,
		scriptCode,
	)
	if len(runtimeErrors) > 0 {
		return executor.Output, 1, convertErrors(runtimeErrors...)
	}
	return executor.Output, exitCode, make([]HomescriptError, 0)
}
