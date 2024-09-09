package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/mattolenik/svg2scad/log"
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
	_, err = cw.WriteLines("include <BOSL2/std.scad>", "include <BOSL2/beziers.scad>", "")
	if err != nil {
		return fmt.Errorf("failed to write OpenSCAD file: %w", err)
	}

	for _, path := range svg.Paths {
		tree, err := ast.Parse(file.Name(), []byte(path.D))
		if err != nil {
			return fmt.Errorf("failed to parse path from SVG %q: %w", path, err)
		}
		err = walk(cw, tree)
		if err != nil {
			return fmt.Errorf("failed to generate OpenSCAD code: %w", err)
		}
	}
	return nil
}

func walk(cw *ast.CodeWriter, node any) (err error) {
	switch node := node.(type) {
	case []any:
		for _, n := range node {
			err := walk(cw, n)
			if err != nil {
				return err
			}
		}
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
		for _, n := range node.Contents {
			if err := walk(cw, n); err != nil {
				return err
			}
		}
		cw.Untab()

		if _, err = cw.WriteLines("}"); err != nil {
			return err
		}
	}
	return nil
}
