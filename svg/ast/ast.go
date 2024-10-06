package ast

import (
	"bytes"
	"fmt"
	"io"
	"math"
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

func (c Coord) Add(coord Coord) Coord {
	return Coord{c[0] + coord[0], c[1] + coord[1]}
}

var Cursor = [2]float64{math.NaN(), math.NaN()} // Represents "cursor" that isn't a numeric value

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

func (c Coords) GoString() string {
	strs := make([]string, len(c))
	for i, s := range c {
		strs[i] = s.String()
	}
	return strings.Join(strs, ", ")
}

func (c Coords) Add(coord Coord) Coords {
	result := make(Coords, len(c))
	for i, cc := range c {
		result[i] = cc.Add(coord)
	}
	return result
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

type Color struct {
	R int
	G int
	B int
	A int
}

type Path struct {
	Name     string
	Children any
}

type CodeWriter struct {
	buf         bytes.Buffer
	depth       int
	tabWidth    int
	tabStr      string
	indentation string
}

func NewCodeWriter() *CodeWriter {
	tw := 4
	return &CodeWriter{
		tabWidth: tw,
		tabStr:   strings.Repeat(" ", tw),
	}
}

func (cw *CodeWriter) Write(writer io.Writer) error {
	if _, err := writer.Write(cw.buf.Bytes()); err != nil {
		return fmt.Errorf("failed to write code: %w", err)
	}
	return nil
}

func (cw *CodeWriter) Printf(format string, args ...any) error {
	_, err := cw.buf.WriteString(cw.indentation + fmt.Sprintf(format, args...))
	return err
}

func (cw *CodeWriter) Linef(format string, args ...any) *CodeWriter {
	cw.Lines(fmt.Sprintf(format, args...))
	return cw
}

func (cw *CodeWriter) Lines(code ...string) *CodeWriter {
	for _, line := range code {
		cw.buf.WriteString(cw.indentation + line + "\n")
	}
	return cw
}

func (cw *CodeWriter) BlankLine() *CodeWriter {
	cw.BlankLines(1)
	return cw
}

func (cw *CodeWriter) BlankLines(num int) *CodeWriter {
	cw.buf.WriteString(strings.Repeat("\n", num))
	return cw
}

func (cw *CodeWriter) Tab() *CodeWriter {
	cw.buf.WriteString(cw.tabStr)
	return cw
}

func (cw *CodeWriter) Indent() *CodeWriter {
	cw.depth++
	cw.indentation = strings.Repeat(cw.tabStr, cw.depth)
	return cw
}

func (cw *CodeWriter) Unindent() *CodeWriter {
	cw.depth--
	cw.indentation = strings.Repeat(cw.tabStr, cw.depth)
	return cw
}

func (cw *CodeWriter) OpenBrace() *CodeWriter {
	cw.Lines("{")
	cw.Indent()
	return cw
}

func (cw *CodeWriter) CloseBrace() *CodeWriter {
	cw.Unindent()
	cw.Lines("}")
	return cw
}
