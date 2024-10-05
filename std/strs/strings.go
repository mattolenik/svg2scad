package strs

import (
	"fmt"
	"strings"

	"github.com/mattolenik/svg2scad/std"
)

func Bracketed(items []any) string {
	return "[ " + strings.Join(std.Map(items, func(v any) string { return fmt.Sprintf("%v", v) }), ", ") + " ]"
}

func IsLower(s string) bool {
	return strings.ToLower(s) == s
}

func IsUpper(s string) bool {
	return strings.ToUpper(s) == s
}

func Capitalize(s string) string {
	return strings.ToUpper(s[0:1]) + s[1:]
}
