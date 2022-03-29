package cmd

import (
	"log"
	"os"
)

var normalLogger *log.Logger
var errLogger *log.Logger

func logf(args ...string) {
	output := ""
	for _, arg := range args {
		output += arg
	}
	normalLogger.Printf(output)
}

func logn(args ...string) {
	output := ""
	for _, arg := range args {
		output += arg
	}
	normalLogger.Println(output)
}

func loge(args ...string) {
	output := ""
	for _, arg := range args {
		output += arg
	}
	errLogger.Println(output)
}

func initLog(useTime bool) {
	if useTime {
		errLogger = log.New(os.Stderr, "", log.Default().Flags())
		normalLogger = log.New(os.Stdout, "", log.Default().Flags())
	} else {
		errLogger = log.New(os.Stderr, "", 0)
		normalLogger = log.New(os.Stdout, "", 0)
	}
}
