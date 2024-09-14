package scad

import (
	"fmt"
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
	OutDir      string
}

func NewSCADWriter(outDir string) *SCADWriter {
	return &SCADWriter{
		StrokeWidth: 5,
		OutDir:      outDir,
	}
}

func (sw *SCADWriter) ConvertSVG(svg *svg.SVG, filename string) error {
	if filename == "" {
		ext := filepath.Ext(svg.Filename)
		filename = svg.Filename[:len(svg.Filename)-len(ext)] + ".scad"
	} else {
		filename = std.EnsureSuffix(filename, ".scad")
	}
	outPath := filepath.Join(sw.OutDir, filename)
	file, err := os.Create(outPath)
	if err != nil {
	}
	defer file.Close()

	cw := ast.NewCodeWriter(file)
	defer cw.Close()
	cw.WriteLines(DefaultImports...)
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
		cw.WriteStrings("", fmt.Sprintf("%s();", module))
	}
	return nil
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
		cw.WriteLines(fmt.Sprintf("let(cursor = %v)", node.Coord))
		cw.OpenBrace()
		for _, child := range node.Children {
			if err := sw.walk(cw, child, foundModule, lastPoint, pathID); err != nil {
				return err
			}
		}
		cw.CloseBrace()
	case *ast.Bezier:
		coords := ast.Coords(append([]ast.Coord{*lastPoint}, node.Points...))
		cw.WriteLinef("%s = %v;", node.Name, coords)
		cw.WriteLinef("%s_curve = bezier_curve(%s);", node.Name, node.Name)
		//cw.WriteLinef("debug_bezier(%s, N=len(%s)-1);", node.Name, node.Name)
		cw.WriteLinef(`stroke(%s_curve, width = %d);`, node.Name, sw.StrokeWidth)
		*lastPoint = node.Points[len(node.Points)-1]
	case *ast.LineTo:
		// TODO: line command
	case *ast.ClosePath:
		// TODO: close path
	case *ast.Module:
		if pathID != "" {
			node.Name = pathID
		}
		cw.WriteLinef("module %s(anchor, spin, orient)", node.Name)

		cw.OpenBrace()
		if err := sw.walk(cw, node.Children, foundModule, lastPoint, pathID); err != nil {
			return err
		}
		cw.CloseBrace()

		foundModule(node)
	default:
		return fmt.Errorf("unsupported command: %q", reflect.TypeOf(node))
	}
	return nil
}
