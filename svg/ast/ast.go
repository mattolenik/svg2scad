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

type CubicCoords struct {
	C1, C2, C3 *Coord
}

type BezierCurve struct {
	Coords *CubicCoords
}

type ClosePath struct{}
