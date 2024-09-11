package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"reflect"

	"github.com/k0kubun/pp/v3"
	"github.com/mattolenik/svg2scad/log"
	"github.com/mattolenik/svg2scad/scad"
	"github.com/mattolenik/svg2scad/svg"
	"github.com/mattolenik/svg2scad/svg/ast"
)

func main() {
	if err := mainE(); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
		os.Exit(1)
	}
}

func mainE() error {
	outDir := flag.String("out", "./curves", "output directory for .scad files")
	//watch := flag.Bool("watch", false, "watch for changes to the .svg files and refresh .scad files automatically")
	flag.Parse()

	svgFiles := flag.Args()

	if len(svgFiles) == 0 {
		return fmt.Errorf("please provide one or more .svg files to convert")
	}

	if err := os.MkdirAll(*outDir, 0755); err != nil {
		return fmt.Errorf("couldn't create output directory %q: %w", *outDir, err)
	}

	for _, file := range svgFiles {
		svg, err := svg.ReadSVGFromFile(file)
		if err != nil {
			log.Errorf("the SVG file %q could not be read: %w", file, err)
		}
		err = convert(svg, *outDir)
		if err != nil {
			log.Errorf("the SVG file %q could not be converted: %w", file, err)
		}
	}

	return nil
}

func convert(svg *svg.SVG, outDir string) error {
	ext := filepath.Ext(svg.Filename)
	scadFilename := svg.Filename[:len(svg.Filename)-len(ext)] + ".scad"
	outPath := filepath.Join(outDir, scadFilename)
	file, err := os.Create(outPath)
	if err != nil {
	}
	defer file.Close()

	cw := ast.NewCodeWriter(file)
	defer cw.Close()
	cw.WriteLines(scad.DefaultImports...)

	modules := []string{}
	found := func(m *ast.Module) {
		modules = append(modules, m.Name)
	}

	for _, path := range svg.Paths {
		tree, err := ast.Parse(file.Name(), []byte(path.D))
		if err != nil {
			return fmt.Errorf("failed to parse path from SVG %q: %w", path, err)
		}
		err = walk(cw, tree, found)
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

func walk(cw *ast.CodeWriter, node any, foundModule func(m *ast.Module)) (err error) {
	switch node := node.(type) {
	case []any:
		for _, n := range node {
			err := walk(cw, n, foundModule)
			if err != nil {
				return err
			}
		}
	case *ast.MoveTo:
		cw.WriteLines(
			fmt.Sprintf("translate(%v) {", node.Coord))
		cw.Tab()
		for _, child := range node.Children {
			if err := walk(cw, child, foundModule); err != nil {
				return err
			}
		}
		cw.Untab()
		cw.WriteLines("}")
	case *ast.Bezier:
		cw.WriteLines(
			fmt.Sprintf("%s = %v;", node.Name, node.Points),
			fmt.Sprintf("%s_curve = bezpath_curve(%s, N=len(%s)-1);", node.Name, node.Name, node.Name),
			// fmt.Sprintf("debug_bezier(%s, N=len(%s)-1);", node.Name, node.Name))
			fmt.Sprintf(`stroke(%s_curve, width=5, dots=false, dots_color="red");`, node.Name))
	case *ast.LineTo:
		// TODO: line command
	case *ast.ClosePath:
		// TODO: close path
	case *ast.Module:
		cw.WriteLines(fmt.Sprintf("module %s(anchor, spin, orient) {", node.Name))

		cw.Tab()
		if err := walk(cw, node.Children, foundModule); err != nil {
			return err
		}
		cw.Untab()

		cw.WriteLines("}")
		foundModule(node)
	default:
		return fmt.Errorf("unsupported command: %q", reflect.TypeOf(node))
	}
	return nil
}
