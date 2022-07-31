package cmd

import (
	"fmt"
	"strings"
)

func processHmsArgs(args []string) (map[string]string, error) {
	argsTemp := make(map[string]string, 0)
	for indexArg, arg := range args {
		switch strings.Count(arg, ":") {
		case 0:
			return nil, fmt.Errorf("Bad Homescript argument formatting at position %d: '%s' does not contain the ':' separator.\nThe separator is required in order to distinguish between key and value", indexArg, arg)
		case 1:
			argsTemp[strings.Split(arg, ":")[0]] = strings.Split(arg, ":")[1]
		default:
			return nil, fmt.Errorf("Bad Homescript argument formatting at position %d: '%s' contains more than one ':' separator.\nThe separator is required in order to distinguish between key and value", indexArg, arg)
		}
	}
	return argsTemp, nil
}
