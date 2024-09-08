package log

import (
	"fmt"
	"os"
	"strings"
)

func Errorf(format string, a ...any) {
	errString := strings.Trim(fmt.Errorf(format, a...).Error(), "\n")
	fmt.Fprintln(os.Stderr, errString)
}
