package scad

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/mattolenik/svg2scad/files"
	"github.com/mattolenik/svg2scad/log"
	"github.com/mattolenik/svg2scad/svg"
	"github.com/mattolenik/svg2scad/svg/ast"
)

type SCADWriter struct {
	SplineSteps   int
	PrintExamples bool
}

func id(p *int) int {
	*p++
	return *p
}

func (sw *SCADWriter) ConvertSVG(svg *svg.SVG, outDir, filename string) error {
	outPath := filepath.Join(outDir, filename)
	writer, err := os.Create(outPath)
	if err != nil {
		return fmt.Errorf("couldn't create output .scad file %q: %w", outPath, err)
	}
	defer writer.Close()
	err = sw.ConvertSVGToSCAD(svg, writer, outPath)
	if err != nil {
		return err
	}
	err = files.WriteFileWithDir(filepath.Join(outDir, LibSubdir, LibFilename), []byte(LibFileData))
	if err != nil {
		return fmt.Errorf("failed to write support file: %w", err)
	}
	return nil
}

func (sw *SCADWriter) ConvertSVGToSCAD(svg *svg.SVG, output io.Writer, outPath string) error {
	cw := ast.NewCodeWriter()
	cw.Lines(Imports...)
	cw.BlankLine()

	pathNames := []string{}

	ids := map[string]int{} // for tracking path IDs in the loop below
	uniq := 0

	for _, path := range svg.Paths {
		if path.ID == "" {
			// Give unnamed paths a default name
			path.ID = fmt.Sprintf("path_%d", id(&uniq))
		} else {
			// Make sure path.ID is fully unique. There shouldn't be multiple <path> elements in the
			// SVG with the same ID, but if there are, they will be renamed with a numerical suffix.
			count, exists := ids[path.ID]
			if !exists {
				ids[path.ID] = 1
			} else {
				ids[path.ID] = count + 1
				path.ID = fmt.Sprintf("%s_%d", path.ID, ids[path.ID])
			}
		}
		tree, err := ast.Parse(path.ID, []byte(path.D))
		if err != nil {
			return fmt.Errorf("failed to parse path from SVG %q: %w", path, err)
		}

		state := walkState{
			paths:  []string{},
			pathID: path.ID,
		}
		_, err = sw.walk(cw, tree, &state)
		if err != nil {
			return fmt.Errorf("failed to generate OpenSCAD code: %w", err)
		}
		pathNames = append(pathNames, state.paths...)
	}
	cw.BlankLine()
	for _, name := range pathNames {
		cw.BlankLine()
		cw.Linef("module %s(depth=0, anchor, spin, orient)", name)
		cw.OpenBrace().Linef(
			"p = %s([ 0, 0 ]);", name).Linef(
			"exts = %s(p);", EXTENTS).Lines(
			"width = exts[0][0] - exts[1][0];",
			"height = exts[0][1] - exts[1][1];",
			"two_d = depth == 0;",
			"size = two_d ? [ width, height ] : [ width, height, depth ];",
			"attachable(anchor, spin, orient, two_d = two_d, size = size)").
			OpenBrace().
			Lines(
				"translate(-[ width / 2 + exts[1][0], height / 2 + exts[1][1], depth / 2 ])",
				"if (!two_d) { linear_extrude(depth) polygon(p); } else { polygon(p); }",
				"children();",
			).
			CloseBrace()
		cw.CloseBrace()
	}
	log.Userf("curves: %s", strings.Join(pathNames, ", "))
	if sw.PrintExamples {
		log.Userf("\n  Usage, assuming your .scad file is in the current folder:\n")
		log.Userf("  include <%s>", outPath)
		log.Userf("  %s(100);  // get a 3D object, your path extruded by 100mm", pathNames[0])
		log.Userf("  %s();     // get a 2D path", pathNames[0])
		log.Userf("")
	}
	return cw.Write(output)
}

type walkState struct {
	paths  []string
	pathID string
	points []ast.Coord
}

func (ws *walkState) addPoint(p ast.Coord) {
	ws.points = append(ws.points, p)
}

func (ws *walkState) firstPoint() ast.Coord {
	return ws.points[0]
}

func (ws *walkState) lastPoint() ast.Coord {
	return ws.points[len(ws.points)-1]
}

func (sw *SCADWriter) walk(cw *ast.CodeWriter, node any, state *walkState) (val any, err error) {
	log.Debugf(reflect.TypeOf(node).String())
	switch node := node.(type) {

	case *ast.MoveTo:
		if node.Relative {
			node.Coord = node.Coord.Add(state.lastPoint())
		}
		cw.Linef("let(%s = %s + %v)", ast.Cursor, ast.Cursor, node.Coord)
		state.addPoint(node.Coord)
		return nil, nil

	case ast.CommandList:
		curveCoords := []ast.Coords{}
		for _, child := range node {
			r, err := sw.walk(cw, child, state)
			if err != nil {
				return nil, fmt.Errorf("failed building curve: %w", err)
			}
			if r == nil {
				continue
			}
			switch r := r.(type) {
			case ast.Coords:
				curveCoords = append(curveCoords, r)
				state.addPoint(r[len(r)-1])
			default:
				return nil, fmt.Errorf("type %v is not supported", reflect.TypeOf(r))
			}
		}
		cw.Linef("let(curve = [ %s, ", ast.Cursor)
		cw.Indent()
		colWidths := make([][2]int, 3)
		for i, ws := range colWidths {
			for c := range ws {
				for _, r := range curveCoords {
					if len(r[i][c]) > colWidths[i][c] {
						colWidths[i][c] = len(r[i][c])
					}
				}
			}
		}
		for _, coord := range curveCoords {
			cw.Lines(coord.Columnized(colWidths) + ",")
		}
		cw.Unindent()
		cw.Lines("],")

		cw.Linef("path = bezpath_curve(curve, splinesteps = %d))", sw.SplineSteps)

	case *ast.CubicBezier:
		if node.Relative {
			node.Points = node.Points.Add(state.lastPoint())
		}
		return node.Points, nil

	case *ast.LineTo:
		if node.Relative {
			node.Coord = node.Coord.Add(state.lastPoint())
		}
		// Convert to a curve, it's easier to create the geometry in OpenSCAD as all bezier
		return ast.Coords{node.Coord, node.Coord, node.Coord}, nil

	case *ast.ClosePath:
		c := state.firstPoint()
		return ast.Coords{c, c, c}, nil

	case *ast.Path:
		node.Name = state.pathID
		cw.Linef("function %s(%s) =", node.Name, ast.Cursor)
		cw.Indent()
		defer cw.Unindent()
		defer func() { state.paths = append(state.paths, node.Name) }()
		return sw.walk(cw, node.Children, state)

	case []any:
		results := []any{}
		for _, n := range node {
			r, err := sw.walk(cw, n, state)
			if err != nil {
				return nil, err
			}
			results = append(results, r)
		}
		return results, nil

	default:
		return nil, fmt.Errorf("unsupported command: %q", reflect.TypeOf(node))
	}
	cw.Tab().Lines("path;")
	return nil, nil
}
