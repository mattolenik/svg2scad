package ast

import (
	"fmt"
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

func (b *Bezier) ToSCAD() (string, error) {
	r := strings.NewReplacer("{varName}", b.Name, "{points}", b.Points.String())
	return r.Replace(`{varName} = {points}; debug_bezier({varName}, N=len({varName})-1);`), nil
}

type ClosePath struct{}

type Scaddable interface {
	ToSCAD() (string, error)
}

type Module struct {
	Name     string
	Contents []any
}
