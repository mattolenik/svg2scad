package scad

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"

	"github.com/k0kubun/pp/v3"
	"github.com/mattolenik/svg2scad/std"
	"github.com/mattolenik/svg2scad/svg"
	"github.com/mattolenik/svg2scad/svg/ast"
)

type SCADWriter struct {
	StrokeWidth int
	SplineSteps int
}

func NewSCADWriter(outDir string) *SCADWriter {
	return &SCADWriter{
		StrokeWidth: 2,
		SplineSteps: 32,
	}
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

	modules := []string{}
	found := func(m *ast.Module) {
		modules = append(modules, m.Name)
	}

	for _, path := range svg.Paths {
		tree, err := ast.Parse(path.ID, []byte(path.D))
		if err != nil {
			return fmt.Errorf("failed to parse path from SVG %q: %w", path, err)
		}
		err = sw.walk(cw, tree, found, &ast.Coord{0, 0}, path.ID)
		if err != nil {
			return fmt.Errorf("failed to generate OpenSCAD code: %w", err)
		}
		pp.Println(tree)
	}
	for _, module := range modules {
		cw.Linef("%s();", module)
	}
	return cw.Write(output)
}

func (sw *SCADWriter) walk(cw *ast.CodeWriter, node any, foundModule func(m *ast.Module), lastPoint *ast.Coord, pathID string) (err error) {
	switch node := node.(type) {
	case []any:
		for _, n := range node {
			err := sw.walk(cw, n, foundModule, lastPoint, pathID)
			if err != nil {
				return err
			}
		}
	case *ast.MoveTo:
		*lastPoint = node.Coord
		cw.Lines(fmt.Sprintf("let(cursor = %v)", node.Coord))
		if err := cw.BraceIndent(func() error {
			return sw.walk(cw, node.Children, foundModule, lastPoint, pathID)
		}); err != nil {
			return err
		}
	case *ast.Bezier:
		coords := ast.Coords(append([]ast.Coord{*lastPoint}, node.Points...))
		cw.Linef("%s = bezier_curve(%v, splinesteps = %d);", node.Name, coords, sw.SplineSteps)
		cw.Linef(`stroke(%s, width = %d);`, node.Name, sw.StrokeWidth)
		*lastPoint = node.Points.End()
	case *ast.LineTo:
		// TODO: line command
	case *ast.ClosePath:
		// TODO: close path
	case *ast.Module:
		if pathID != "" {
			node.Name = pathID
		}
		cw.Linef("module %s(anchor, spin, orient)", node.Name)

		if err := cw.BraceIndent(func() error {
			return sw.walk(cw, node.Children, foundModule, lastPoint, pathID)
		}); err != nil {
			return err
		}

		foundModule(node)
	default:
		return fmt.Errorf("unsupported command: %q", reflect.TypeOf(node))
	}
	return nil
}
