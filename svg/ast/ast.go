package ast

type Coord struct {
	X any
	Y any
}

type MoveTo struct {
	Coord *Coord
}

type LineTo struct {
	Coord any
}

type CubicCoords struct {
	C1, C2, C3 any
}

type BezierCurve struct {
	Coords any
}

type ClosePath struct{}
