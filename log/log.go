package log

import (
	"fmt"
	"os"
	"strings"
)

var logOut = os.Stderr

var Debug = false

func Errorf(format string, a ...any) {
	errString := strings.Trim(fmt.Errorf(format, a...).Error(), "\n")
	fmt.Fprintln(logOut, errString)
}

func Debugf(format string, a ...any) {
	if Debug {
		fmt.Fprintln(logOut, fmt.Sprintf(format, a...))
	}
}
