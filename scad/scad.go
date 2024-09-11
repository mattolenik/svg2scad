package scad

import "slices"

var DefaultImports = []string{"include <BOSL2/std.scad>", "include <BOSL2/beziers.scad>"}

func init() {
	slices.Sort(DefaultImports)
}
