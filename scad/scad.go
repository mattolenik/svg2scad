package scad

var Imports = []string{
	"include <BOSL2/beziers.scad>",
	"include <BOSL2/std.scad>",
	"include <svgsupport.scad>"}

const SupportingFile = `
// Finds the bounding box of a set of coordinates
function extents(coords, largest = [ -1e9, -1e9 ], smallest = [ 1e9, 1e9 ]) =
    len(coords) == 0 ? [ largest, smallest ]
                     : extents(list_tail(coords),
                               [ max(largest[0],  coords[0][0]), max(largest[1],  coords[0][1]) ],
                               [ min(smallest[0], coords[0][0]), min(smallest[1], coords[0][1]) ]);

`
