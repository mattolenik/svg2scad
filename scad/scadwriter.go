package scad

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"

	"github.com/k0kubun/pp/v3"
	"github.com/mattolenik/svg2scad/std"
	"github.com/mattolenik/svg2scad/std/strs"
	"github.com/mattolenik/svg2scad/svg"
	"github.com/mattolenik/svg2scad/svg/ast"
)

type SCADWriter struct {
	StrokeWidth int
	SplineSteps int
	Cursor      string
	nextID      int
}

func NewSCADWriter(outDir string) *SCADWriter {
	return &SCADWriter{
		StrokeWidth: 2,
		SplineSteps: 32,
		Cursor:      "cursor",
		nextID:      0,
	}
}
func (sw *SCADWriter) id() int {
	id := sw.nextID
	sw.nextID++
	return id
}

func (sw *SCADWriter) ConvertSVG(svg *svg.SVG, outDir, filename string) error {
	if filename == "" {
		ext := filepath.Ext(svg.Filename)
		filename = svg.Filename[:len(svg.Filename)-len(ext)] + ".scad"
	} else {
		filename = std.EnsureSuffix(filename, ".scad")
	}
	outPath := filepath.Join(outDir, filename)
	file, err := os.Create(outPath)
	if err != nil {
		return fmt.Errorf("couldn't create output .scad file %q: %w", outPath, err)
	}
	defer file.Close()
	return sw.ConvertSVGToSCAD(svg, file)
}

func (sw *SCADWriter) ConvertSVGToSCAD(svg *svg.SVG, output io.Writer) error {
	cw := ast.NewCodeWriter()
	cw.Lines(DefaultImports...)
	cw.BlankLine()
	cw.Lines(Functions...)
	cw.BlankLine()

	pathFunctions := []string{}
	appendFunction := func(m *ast.Path) {
		pathFunctions = append(pathFunctions, m.Name)
	}

	ids := map[string]int{}
	for _, path := range svg.Paths {
		if path.ID == "" {
			path.ID = fmt.Sprintf("path_%d", sw.id())
		} else {
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
		if _, err := sw.walk(cw, tree,
			&walkState{
				addFunction: appendFunction,
				lastPoint:   ast.Coord{0, 0}.String(),
				pathID:      path.ID,
			},
		); err != nil {
			return fmt.Errorf("failed to generate OpenSCAD code: %w", err)
		}
		pp.Println(tree)
	}
	cw.BlankLine()
	for _, fn := range pathFunctions {
		cw.BlankLine()
		cw.Linef("module %s(depth=0, anchor, spin, orient)", fn)
		cw.OpenBrace().Linef(
			"p = %s([ 0, 0 ]);", fn).Lines(
			"exts = extents(p);",
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
	for _, module := range pathFunctions {
		cw.Linef("%s(100);", module)
	}
	return cw.Write(output)
}

type walkState struct {
	lastPoint   string
	addFunction func(*ast.Path)
	pathID      string
}

func (sw *SCADWriter) walk(cw *ast.CodeWriter, node any, state *walkState) (val any, err error) {
	fmt.Println(reflect.TypeOf(node))
	switch node := node.(type) {

	case *ast.MoveTo:
		state.lastPoint = sw.Cursor
		cw.Linef("let(%s = %s + %v,", sw.Cursor, sw.Cursor, node.Coord)
		return nil, nil

	case ast.CommandList:
		coords := []any{sw.Cursor}
		curveVars := []string{}
		lines := []string{}
		for i, child := range node {
			r, err := sw.walk(cw, child, state)
			if err != nil {
				return nil, fmt.Errorf("failed building curve: %w", err)
			}
			if r == nil {
				continue
			}
			switch r := r.(type) {
			case ast.Coords:
				varName := fmt.Sprintf("c%03d", i)
				curveVars = append(curveVars, varName)
				cw.Linef("%s = %v,", varName, r)
			case string:
				lines = append(lines, r)
			case []string:
				lines = append(lines, r...)
			default:
				return nil, fmt.Errorf("type %v is not supported", reflect.TypeOf(r))
			}
		}
		for _, varName := range curveVars {
			idx := func(i int) string { return fmt.Sprintf("%s[%d]", varName, i) }
			coords = append(coords, idx(0), idx(1), idx(2))
		}
		cw.Linef("curve = %s)", strs.Bracketed(coords))
		cw.Linef("let(path = bezpath_curve(curve, splinesteps=%d))", sw.SplineSteps)
		cw.Lines(lines...)

	case *ast.CubicBezier:
		return node.Points, nil

	case *ast.LineTo:
		return fmt.Sprintf("let(path = concat(path, [ %v ]))", node.Coord), nil
	case *ast.ClosePath:
		// TODO: close path
		return nil, nil

	case *ast.Path:
		node.Name = state.pathID
		cw.Linef("function %s(%s) =", node.Name, sw.Cursor)
		cw.Tab()
		defer cw.Untab()
		defer state.addFunction(node)
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
	cw.Lines("    path;")
	return nil, nil
}
