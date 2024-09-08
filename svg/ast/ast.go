package ast

type Coord struct {
	X float64
	Y float64
}

type MoveTo struct {
	Coord *Coord
}

type LineTo struct {
	Coord *Coord
}

type Bezier struct {
	Points []*Coord
}

type ClosePath struct{}
