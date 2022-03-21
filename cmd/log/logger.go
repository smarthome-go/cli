package log

import (
	"os"

	"github.com/sirupsen/logrus"
)

var log *logrus.Logger

func InitLogger(level logrus.Level) {
	log = logrus.New()
	log.SetLevel(level)
}

func Trace(args ...string) {
	output := ""
	for _, arg := range args {
		output += arg
	}
	log.Trace(output)
}

func Debug(args ...string) {
	output := "\x1b[1;36m"
	for _, arg := range args {
		output += arg
	}
	output += "\x1b[1;0m"
	log.Debug(output)
}

func Info(args ...string) {
	var output string
	for _, arg := range args {
		output += arg
	}
	log.Info(output)
}

func Warn(args ...string) {
	output := "\x1b[1;33m"
	for _, arg := range args {
		output += arg
	}
	output += "\x1b[1;0m"
	log.Warn(output)
}

func Error(args ...string) {
	var output string
	for _, arg := range args {
		output += arg
	}
	log.Error(output)
}

func Fatal(args ...string) {
	output := "\x1b[1;31m"
	for _, arg := range args {
		output += arg
	}
	output += "\x1b[1;0m"
	log.Fatal(output)
	os.Exit(1)
}
