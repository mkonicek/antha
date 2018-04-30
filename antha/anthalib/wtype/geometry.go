// wtype/geometry.go: Part of the Antha language
// Copyright (C) 2014 the Antha authors. All rights reserved.
//
// This program is free software; you can redistribute it and/or
// modify it under the terms of the GNU General Public License
// as published by the Free Software Foundation; either version 2
// of the License, or (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program; if not, write to the Free Software
// Foundation, Inc., 51 Franklin Street, Fifth Floor, Boston, MA  02110-1301, USA.
//
// For more information relating to the software or licensing issues please
// contact license@antha-lang.org or write to the Antha team c/o
// Synthace Ltd. The London Bioscience Innovation Centre
// 2 Royal College St, London NW1 0NH UK

package wtype

import (
	"fmt"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"math"
)

type Coordinates struct {
	X float64
	Y float64
	Z float64
}

func (c Coordinates) Equals(c2 Coordinates) bool {
	return c.X == c2.X && c.Y == c2.Y && c.Z == c2.Z
}

//String implements Stringer
func (self Coordinates) String() string {
	return fmt.Sprintf("%vx%vx%vmm", self.X, self.Y, self.Z)
}

//Dim Value for dimension
func (a Coordinates) Dim(x int) float64 {
	switch x {
	case 0:
		return a.X
	case 1:
		return a.Y
	case 2:
		return a.Z
	default:
		return 0.0
	}
}

//Add Addition returns a new wtype.Coordinates
func (self Coordinates) Add(rhs Coordinates) Coordinates {
	return Coordinates{self.X + rhs.X,
		self.Y + rhs.Y,
		self.Z + rhs.Z}
}

//Subtract returns a new wtype.Coordinates
func (self Coordinates) Subtract(rhs Coordinates) Coordinates {
	return Coordinates{self.X - rhs.X,
		self.Y - rhs.Y,
		self.Z - rhs.Z}
}

//Multiply returns a new wtype.Coordinates
func (self Coordinates) Multiply(v float64) Coordinates {
	return Coordinates{self.X * v,
		self.Y * v,
		self.Z * v}
}

//Divide returns a new wtype.Coordinates
func (self Coordinates) Divide(v float64) Coordinates {
	return Coordinates{self.X / v,
		self.Y / v,
		self.Z / v}
}

//Dot product
func (self Coordinates) Dot(rhs Coordinates) float64 {
	return self.X*rhs.X + self.Y + rhs.Y + self.Z + rhs.Z
}

//Abs L2-Norm
func (self Coordinates) Abs() float64 {
	return math.Sqrt(self.X*self.X + self.Y*self.Y + self.Z*self.Z)
}

//AbsXY L2-Norm in XY only
func (self Coordinates) AbsXY() float64 {
	return math.Sqrt(self.X*self.X + self.Y*self.Y)
}

//Unit return a Unit vector in the same direction as the coordinates
func (self Coordinates) Unit() Coordinates {
	return self.Divide(self.Abs())
}

//Geometry interface for 3D geometry
type Geometry interface {
	Height() wunit.Length
	Width() wunit.Length
	Depth() wunit.Length
}

type PointSet []Coordinates

func (ps PointSet) CentreTo(c Coordinates) PointSet {
	ret := make(PointSet, len(ps))

	for i, p := range ps {
		ret[i] = p.Subtract(c)
	}

	return ret
}
