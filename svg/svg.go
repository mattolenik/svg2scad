package svg

import (
	"encoding/xml"
	"fmt"
	"io"
	"os"
)

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

func ReadSVGFromFile(path string) (*SVG, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	svg, err := ReadSVG(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read SVG: %w", err)
	}

	return svg, nil
}

func ReadSVG(r io.Reader) (*SVG, error) {
	var svg SVG
	err := xml.NewDecoder(r).Decode(&svg)
	if err != nil {
		return nil, fmt.Errorf("failed to decode SVG: %w", err)
	}
	return &svg, nil
}