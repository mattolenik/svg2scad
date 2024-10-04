package scad

import (
	"fmt"
	"slices"
	"strings"
)

var DefaultImports = []string{"include <BOSL2/std.scad>", "include <BOSL2/beziers.scad>"}

var Functions = []string{
	`
function extents(coords, largest = [ -1e9, -1e9 ], smallest = [ 1e9, 1e9 ]) =
    len(coords) == 0 ? [ largest, smallest ]
                     : extents(list_tail(coords),
                               [ max(largest[0], coords[0][0]), max(largest[1], coords[0][1]) ],
                               [ min(smallest[0], coords[0][0]), min(smallest[1], coords[0][1]) ]);
`}

func init() {
	slices.Sort(DefaultImports)
	for i, f := range Functions {
		Functions[i] = fmt.Sprintf("\n%s\n", strings.TrimSpace(f)) // normalize spacing
	}
}
