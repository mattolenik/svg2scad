package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/mattolenik/svg2scad/files"
	"github.com/mattolenik/svg2scad/log"
	"github.com/mattolenik/svg2scad/scad"
	"github.com/mattolenik/svg2scad/svg"
)

func main() {
	if err := mainE(os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
		os.Exit(1)
	}
}

func mainE(args []string) error {
	sw := scad.SCADWriter{}
	help := flag.Bool("help", false, "Show help screen")
	outDir := flag.String("out", "./svg-scad", "Output directory for .scad files")
	//watch := flag.Bool("watch", false, "watch for changes to the .svg files and refresh .scad files automatically")
	flag.IntVar(&sw.SplineSteps, "detail", 32, "Higher values create smoother curves, excessive values may cause issues")
	flag.BoolVar(&log.Debug, "debug", false, "Print debug/tracing info, for development use")
	flag.BoolVar(&log.Quiet, "quiet", false, "Quiet mode, don't print info messages, only errors")

	flag.CommandLine.Parse(args)

	svgFiles := flag.Args()

	if *help {
		flag.Usage()
		os.Exit(1)
	}

	if len(svgFiles) == 0 {
		flag.Usage()
		log.Errorf("please provide one or more .svg files to convert")
		os.Exit(1)
	}

	if err := files.CreateDirIfNotExists(*outDir); err != nil {
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
		log.Userf("\n")
	}

	return nil
}
