package log

import (
	"log"
	"os"
)

var normalLogger *log.Logger
var errLogger *log.Logger

func Logf(args ...string) {
	output := ""
	for _, arg := range args {
		output += arg
	}
	normalLogger.Printf(output)
}

func Logn(args ...string) {
	output := ""
	for _, arg := range args {
		output += arg
	}
	normalLogger.Println(output)
}

func Loge(args ...string) {
	output := ""
	for _, arg := range args {
		output += arg
	}
	errLogger.Println(output)
}

func InitLog(useTime bool) {
	if useTime {
		errLogger = log.New(os.Stderr, "", log.Default().Flags())
		normalLogger = log.New(os.Stdout, "", log.Default().Flags())
	} else {
		errLogger = log.New(os.Stderr, "", 0)
		normalLogger = log.New(os.Stdout, "", 0)
	}
}
