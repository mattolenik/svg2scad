package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
)

func main() {
	if err := mainE(); err != nil {
		log.Fatal(err)
	}
}

func mainE() error {
	inFile := flag.String("in", "", "input .svg file")
	outFile := flag.String("out", "", "output .scad file")
	flag.Parse()

	if *inFile == "" {
		fmt.Println("Please provide an input .svg file using the -in flag")
		return nil
	}

	if *outFile == "" {
		fmt.Println("Please provide an output .scad file using the -out flag")
		return nil
	}
	return nil
}

func readSVGFromFile(path string) (*SVG, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	svg, err := readSVG(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read SVG: %w", err)
	}

	return svg, nil
}

func readSVG(r io.Reader) (*SVG, error) {
	var svg SVG
	err := xml.NewDecoder(r).Decode(&svg)
	if err != nil {
		return nil, fmt.Errorf("failed to decode SVG: %w", err)
	}
	return &svg, nil
}

type SVG struct {
	XMLName xml.Name `xml:"svg"`
	Version string   `xml:"version,attr"`
	Width   string   `xml:"width,attr"`
	Height  string   `xml:"height,attr"`
	Path    []Path   `xml:"path"`
}

type Path struct {
	ID string `xml:"id,attr"`
	D  string `xml:"d,attr"`
}
