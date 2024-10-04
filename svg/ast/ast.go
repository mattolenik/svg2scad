package ast

import (
	"bytes"
	"fmt"
	"io"
	"strings"
)

type Curve struct {
	Children []any
}

type CommandList []any

type Coord [2]float64

func (c Coord) String() string {
	return fmt.Sprintf("[ %.3f, %.3f ]", c[0], c[1])
}

type Coords []Coord

func (c Coords) End() Coord {
	return c[len(c)-1]
}

func (c Coords) String() string {
	strs := make([]string, len(c))
	for i, s := range c {
		strs[i] = s.String()
	}
	return "[ " + strings.Join(strs, ", ") + " ]"
}

type MoveTo struct {
	Coord    Coord
	Relative bool
}

type LineTo struct {
	Coord    Coord
	Relative bool
}

type CubicBezier struct {
	Points   Coords
	Relative bool
}

type QuadraticBezier struct {
	Points Coords
}

type ClosePath struct{}

type Path struct {
	Name     string
	Children any
}

type CodeWriter struct {
	buf      *bytes.Buffer
	depth    int
	tabWidth int
	tabStr   string
}

func NewCodeWriter() *CodeWriter {
	return &CodeWriter{buf: &bytes.Buffer{}, depth: 0, tabWidth: 4}
}

func (cw *CodeWriter) Write(writer io.Writer) error {
	if _, err := writer.Write(cw.buf.Bytes()); err != nil {
		return fmt.Errorf("failed to write code: %w", err)
	}
	return nil
}

func (cw *CodeWriter) BraceIndent(action func() error) error {
	cw.OpenBrace()
	err := action()
	if err != nil {
		return fmt.Errorf("error printing indented code: %w", err)
	}
	cw.CloseBrace()
	return nil
}

func (cw *CodeWriter) Linef(format string, args ...any) {
	cw.Lines(fmt.Sprintf(format, args...))
}

func (cw *CodeWriter) Lines(code ...string) {
	for _, line := range code {
		cw.buf.WriteString(cw.tabStr + line + "\n")
	}
}

func (cw *CodeWriter) BlankLine() {
	cw.BlankLines(1)
}

func (cw *CodeWriter) BlankLines(num int) {
	cw.buf.WriteString(strings.Repeat("\n", num))
}

func (cw *CodeWriter) Tab() {
	cw.depth++
	cw.tabStr = strings.Repeat(" ", cw.depth*cw.tabWidth)
}

func (cw *CodeWriter) Untab() {
	cw.depth--
	cw.tabStr = strings.Repeat(" ", cw.depth*cw.tabWidth)
}

func (cw *CodeWriter) OpenBrace() {
	cw.Lines("{")
	cw.Tab()
}

func (cw *CodeWriter) CloseBrace() {
	cw.Untab()
	cw.Lines("}")
}
