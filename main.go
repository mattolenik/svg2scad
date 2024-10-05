package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/mattolenik/svg2scad/scad"
	"github.com/mattolenik/svg2scad/svg"
)

func main() {
	if err := mainE(); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
		os.Exit(1)
	}
}

func mainE() error {

	sw := scad.SCADWriter{}
	outDir := flag.String("out", "./curves", "Output directory for .scad files")
	//watch := flag.Bool("watch", false, "watch for changes to the .svg files and refresh .scad files automatically")
	flag.IntVar(&sw.SplineSteps, "detail", 32, "The number of spline steps, a higher value results in a smoother shape. An excessive value may cause problems.")
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
			return fmt.Errorf("the SVG file %q could not be read: %w", file, err)
		}
		err = sw.ConvertSVG(svg, *outDir, "")
		if err != nil {
			return fmt.Errorf("the SVG file %q could not be converted: %w", file, err)
		}
	}

	return nil
}
