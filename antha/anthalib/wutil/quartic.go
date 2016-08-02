package wutil

import (
	"math"
)

// function for quadratics
type Quartic struct {
	Quartic string
	A       float64
	B       float64
	C       float64
	D       float64
	E       float64
}

func (c Quartic) F(v float64) float64 {
	return c.A*math.Pow(v, 4.0) + c.B*math.Pow(v, 3.0) + c.C*math.Pow(v, 2.0) + c.D*v + c.E
}

func (q Quartic) Name() string {
	return "Quartic"
}
