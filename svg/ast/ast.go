package ast

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

type Coord [2]float64

func (c Coord) String() string {
	return fmt.Sprintf("[%v, %v]", c[0], c[1])
}

type Coords []Coord

func (c Coords) String() string {
	strs := make([]string, len(c))
	for i, s := range c {
		strs[i] = s.String()
	}
	return "[" + strings.Join(strs, ", ") + "]"
}

type MoveTo struct {
	Coord Coord
}

type LineTo struct {
	Coord Coord
}

type Bezier struct {
	Points Coords
	Name   string
}

func (b *Bezier) ToSCAD(cw *CodeWriter) error {
	_, err := cw.WriteLines(
		fmt.Sprintf("%s = %v;", b.Name, b.Points),
		fmt.Sprintf("debug_bezier(%s, N=len(%s)-1);", b.Name, b.Name),
	)
	if err != nil {
		return fmt.Errorf("failed to generate OpenSCAD code: %v", err)
	}
	return nil
}

type ClosePath struct{}

type Scaddable interface {
	ToSCAD(cw *CodeWriter) error
}

type Module struct {
	Name     string
	Contents []any
}

type CodeWriter struct {
	writer   *bufio.Writer
	depth    int
	tabWidth int
	tabStr   string
}

func NewCodeWriter(writer io.Writer) *CodeWriter {
	return &CodeWriter{writer: bufio.NewWriter(writer), depth: 0, tabWidth: 4}
}

func (cw *CodeWriter) Flush() error {
	return cw.writer.Flush()
}

func (cw *CodeWriter) Tab() {
	cw.depth++
	cw.tabStr = strings.Repeat(" ", cw.depth*cw.tabWidth)
}

func (cw *CodeWriter) Untab() {
	cw.depth--
	cw.tabStr = strings.Repeat(" ", cw.depth*cw.tabWidth)
}

func (cw *CodeWriter) WriteLines(code ...string) (int, error) {
	sum := 0
	for _, line := range code {
		n, err := cw.WriteStrings(cw.tabStr, line, "\n")
		if err != nil {
			return sum, err
		}
		sum += n
	}
	return sum, nil
}

func (cw *CodeWriter) WriteStrings(strs ...string) (int, error) {
	sum := 0
	for _, s := range strs {
		n, err := cw.writer.WriteString(s)
		if err != nil {
			return sum, err
		}
		sum += n
	}
	return sum, nil
}
