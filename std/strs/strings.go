package strs

import (
	"fmt"
	"strings"

	"github.com/mattolenik/svg2scad/std"
)

func Bracketed(items []any) string {
	return "[ " + strings.Join(std.Map(items, func(v any) string { return fmt.Sprintf("%v", v) }), ", ") + " ]"
}
