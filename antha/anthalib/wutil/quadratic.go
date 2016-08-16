package wutil

import (
	"math"
)

// function for quadratics
type Quadratic struct {
	Quadratic string // just a label
	A         float64
	B         float64
	C         float64
}

func (c Quadratic) F(v float64) float64 {
	return c.A*math.Pow(v, 2.0) + c.B*v + c.C
}

func (c Quadratic) Name() string {
	return "Quadratic"
}
