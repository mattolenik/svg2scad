package ast

import (
	"bytes"
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
	Coord    Coord
	Children []any
}

type LineTo struct {
	Coord Coord
}

type Bezier struct {
	Points Coords
	Name   string
}

type ClosePath struct{}

type Module struct {
	Name     string
	Children any
}

type CodeWriter struct {
	buf      *bytes.Buffer
	writer   io.WriteCloser
	depth    int
	tabWidth int
	tabStr   string
}

func NewCodeWriter(writer io.WriteCloser) *CodeWriter {
	return &CodeWriter{buf: &bytes.Buffer{}, writer: writer, depth: 0, tabWidth: 4}
}

func (cw *CodeWriter) Close() error {
	if _, err := cw.writer.Write(cw.buf.Bytes()); err != nil {
		return fmt.Errorf("failed to write code: %w", err)
	}
	return cw.writer.Close()
}

func (cw *CodeWriter) Tab() {
	cw.depth++
	cw.tabStr = strings.Repeat(" ", cw.depth*cw.tabWidth)
}

func (cw *CodeWriter) Untab() {
	cw.depth--
	cw.tabStr = strings.Repeat(" ", cw.depth*cw.tabWidth)
}

func (cw *CodeWriter) Indent(action func()) {
	cw.Tab()
	action()
	cw.Untab()
}

func (cw *CodeWriter) WriteLines(code ...string) {
	for _, line := range code {
		cw.WriteStrings(cw.tabStr, line, "\n")
	}
}

func (cw *CodeWriter) WriteStrings(strs ...string) {
	for _, s := range strs {
		cw.buf.WriteString(s)
	}
}
