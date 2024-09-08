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
		return fmt.Errorf("couldn't write OpenSCAD file %q: %w", file.Name(), err)
	}
	defer file.Close()
	for _, path := range svg.Paths {
		tree, err := ast.Parse(file.Name(), []byte(path.D))
		if err != nil {
			return fmt.Errorf("failed to parse path from SVG %q: %w", path, err)
		}
		err = walk(tree)
		if err != nil {
			return fmt.Errorf("failed to generate OpenSCAD code: %w", err)
		}
	}
	return nil
}

func walk(node any) error {
	switch node := node.(type) {
	case []any:
		for _, n := range node {
			err := walk(n)
			if err != nil {
				return err
			}
		}
	case ast.Scaddable:
		scad, err := node.ToSCAD()
		if err != nil {
			return fmt.Errorf("failed to generate OpenSCAD code: %w", err)
		}
		fmt.Println(scad)
	case *ast.Module:
		fmt.Printf("module %s(anchor, spin, orient) {\n", node.Name)
		for _, n := range node.Contents {
			err := walk(n)
			if err != nil {
				return err
			}
		}
		fmt.Printf("}\n")
	}
	return nil
}
