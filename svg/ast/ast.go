package ast

import (
	"fmt"
	"strings"
)

type Coord []float64

func (c Coord) String() string {
	return fmt.Sprintf("[%v, %v]", c[0], c[1])
}

type Coords []Coord

func (c Coords) String() string {
	strs := make([]string, len(c))
	for i, s := range c {
		strs[i] = s.String()
	}
	return "[" + strings.Join(strs, ",") + "]"
}

type MoveTo struct {
	Coord Coord
}

type LineTo struct {
	Coord Coord
}

type Bezier struct {
	Points Coords
}

type ClosePath struct{}
