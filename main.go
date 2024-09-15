package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/mattolenik/svg2scad/log"
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
	outDir := flag.String("out", "./curves", "output directory for .scad files")
	strokeWidth := flag.Int("strokeWidth", 2, "how wide curves will appear on screen, doesn't affect geometry")
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
		sw := scad.NewSCADWriter(*outDir)
		sw.StrokeWidth = *strokeWidth
		err = sw.ConvertSVG(svg, *outDir, "")
		if err != nil {
			log.Errorf("the SVG file %q could not be converted: %w", file, err)
		}
	}

	return nil
}
