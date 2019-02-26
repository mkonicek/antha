package wunit

import (
	"github.com/pkg/errors"
	"strings"
)

// Axis represent a particular direction
type Axis int

const (
	XAxis Axis = iota
	YAxis
	ZAxis
)

// String return the name of the axis
func (a Axis) String() string {
	switch a {
	case XAxis:
		return "X"
	case YAxis:
		return "Y"
	case ZAxis:
		return "Z"
	}
	panic("unknown axis")
}

// AxisFromString return the relevant axis from the string which should be
// "X", "Y", or "Z" (or lowercase), otherwise returns an invalid axis and
// an error
func AxisFromString(a string) (Axis, error) {
	switch strings.ToUpper(a) {
	case "X":
		return XAxis, nil
	case "Y":
		return YAxis, nil
	case "Z":
		return ZAxis, nil
	}
	return Axis(-1), errors.Errorf("unknown axis %q", a)
}
