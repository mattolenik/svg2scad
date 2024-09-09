package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

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
	defer cw.Flush()
	_, err = cw.WriteLines(scad.DefaultImports...)
	if err != nil {
		return fmt.Errorf("failed to write OpenSCAD file: %w", err)
	}

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
	}
	for _, module := range modules {
		if _, err := cw.WriteStrings("", fmt.Sprintf("%s();", module)); err != nil {
			return fmt.Errorf("failed to generate OpenScad code: %w", err)
		}
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
	// case *ast.MoveTo:
	// 	if _, err := cw.WriteLines(
	// 		fmt.Sprintf("translate(%v) {", node.Coord),
	// 	); err != nil {
	// 		return err
	// 	}
	// 	cw.Tab()
	case *ast.Bezier:
		if _, err := cw.WriteLines(
			fmt.Sprintf("%s = %v;", node.Name, node.Points),
			fmt.Sprintf("debug_bezier(%s, N=len(%s)-1);", node.Name, node.Name),
		); err != nil {
			return err
		}
	case *ast.Module:
		if _, err := cw.WriteLines(fmt.Sprintf("module %s(anchor, spin, orient) {", node.Name)); err != nil {
			return err
		}

		cw.Tab()
		for _, n := range node.Children {
			if err := walk(cw, n, foundModule); err != nil {
				return err
			}
		}
		cw.Untab()

		if _, err = cw.WriteLines("}"); err != nil {
			return err
		}
		foundModule(node)
	}
	return nil
}
