package log

import (
	"fmt"
	"os"
	"strings"
)

var Debug = false
var Quiet = false

var logOut = os.Stderr
var userOut = os.Stdout

var logln = func(format string, a ...any) {
	if !Quiet {
		fmt.Fprintln(logOut, strings.TrimSpace(fmt.Sprintf(format, a...)))
	}
}

var userln = func(format string, a ...any) {
	if !Quiet {
		fmt.Fprintln(userOut, fmt.Sprintf(format, a...))
	}
}

func Errorf(format string, a ...any) {
	errString := "❌ " + strings.Trim(fmt.Errorf(format, a...).Error(), "\n")
	logln(errString)
}

func Debugf(format string, a ...any) {
	if Debug {
		logln(format, a...)
	}
}

func Infof(format string, a ...any) {
	logln("ℹ "+format, a...)
}

func Userf(format string, a ...any) {
	userln(format, a...)
}
