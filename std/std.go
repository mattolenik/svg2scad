package std

import "strings"

func EnsureSuffix(str, suffix string) string {
	if strings.HasSuffix(str, suffix) {
		return str
	}
	return str + suffix
}
