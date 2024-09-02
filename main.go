package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/mattolenik/svg2scad/svg"
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
	fmt.Println(*outDir)
	fmt.Println(svgFiles)

	if len(svgFiles) == 0 {
		return fmt.Errorf("please provide one or more .svg files to convert")
	}
	for _, file := range svgFiles {
		svg, err := svg.ReadSVGFromFile(file)
		if err != nil {
			return fmt.Errorf("the SVG file %q could not be read: %w", file, err)
		}
		fmt.Println(svg)
	}

	return nil
}
