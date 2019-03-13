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
	"math"
)

type Coordinates3D struct {
	X float64
	Y float64
	Z float64
}

func (c Coordinates3D) Equals(c2 Coordinates3D) bool {
	return c.X == c2.X && c.Y == c2.Y && c.Z == c2.Z
}

//String implements Stringer
func (self Coordinates3D) String() string {
	return fmt.Sprintf("%.1fx%.1fx%.1f mm", self.X, self.Y, self.Z)
}

func (self Coordinates3D) StringXY() string {
	return fmt.Sprintf("%.1fx%.1f mm", self.X, self.Y)
}

//Dim Value for dimension
func (a Coordinates3D) Dim(x int) float64 {
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
func (self Coordinates3D) Add(rhs Coordinates3D) Coordinates3D {
	return Coordinates3D{self.X + rhs.X,
		self.Y + rhs.Y,
		self.Z + rhs.Z}
}

//Subtract returns a new wtype.Coordinates
func (self Coordinates3D) Subtract(rhs Coordinates3D) Coordinates3D {
	return Coordinates3D{self.X - rhs.X,
		self.Y - rhs.Y,
		self.Z - rhs.Z}
}

//Multiply returns a new wtype.Coordinates
func (self Coordinates3D) Multiply(v float64) Coordinates3D {
	return Coordinates3D{self.X * v,
		self.Y * v,
		self.Z * v}
}

//Divide returns a new wtype.Coordinates divided by v. If v is zero, inf will be returned
func (self Coordinates3D) Divide(v float64) Coordinates3D {
	return Coordinates3D{self.X / v,
		self.Y / v,
		self.Z / v}
}

//Dot product
func (self Coordinates3D) Dot(rhs Coordinates3D) float64 {
	return self.X*rhs.X + self.Y + rhs.Y + self.Z + rhs.Z
}

//Abs L2-Norm
func (self Coordinates3D) Abs() float64 {
	return math.Sqrt(self.X*self.X + self.Y*self.Y + self.Z*self.Z)
}

//AbsXY L2-Norm in XY only
func (self Coordinates3D) AbsXY() float64 {
	return math.Sqrt(self.X*self.X + self.Y*self.Y)
}

//Unit return a Unit vector in the same direction as the coordinates
func (self Coordinates3D) Unit() Coordinates3D {
	return self.Divide(self.Abs())
}

//To2D return a two dimensional coordinate by dropping z dimension
func (self Coordinates3D) To2D() Coordinates2D {
	return Coordinates2D{
		X: self.X,
		Y: self.Y,
	}
}

type PointSet []Coordinates3D

func (ps PointSet) CentreTo(c Coordinates3D) PointSet {
	ret := make(PointSet, len(ps))

	for i, p := range ps {
		ret[i] = p.Subtract(c)
	}

	return ret
}

type Coordinates2D struct {
	X float64
	Y float64
}

//String string representation of the coordinate
func (self Coordinates2D) String() string {
	return fmt.Sprintf("%.1fx%.1f mm", self.X, self.Y)
}

//Equals return true if the two coordinates are equal
func (self Coordinates2D) Equals(rhs Coordinates2D) bool {
	return self == rhs
}

//Add return a new coordinate which is the sum of the two
func (self Coordinates2D) Add(rhs Coordinates2D) Coordinates2D {
	return Coordinates2D{
		X: self.X + rhs.X,
		Y: self.Y + rhs.Y,
	}
}

//Subtract return a new coordinate which is self minus rhs
func (self Coordinates2D) Subtract(rhs Coordinates2D) Coordinates2D {
	return Coordinates2D{
		X: self.X - rhs.X,
		Y: self.Y - rhs.Y,
	}
}

//Multiply return a new coordinate scaled by factor
func (self Coordinates2D) Multiply(factor float64) Coordinates2D {
	return Coordinates2D{
		X: self.X * factor,
		Y: self.Y * factor,
	}
}

//Divide return a new coordinate scaled by the reciprocal of factor
//if factor is zero, inf will be returned
func (self Coordinates2D) Divide(factor float64) Coordinates2D {
	return Coordinates2D{
		X: self.X / factor,
		Y: self.Y / factor,
	}
}

//Abs return the L2 norm of the coordinate
func (self Coordinates2D) Abs() float64 {
	return math.Sqrt(self.X*self.X + self.Y*self.Y)
}

// SBSFootprint the size of standard SBS format plates
var SBSFootprint = Coordinates2D{X: 127.76, Y: 85.48}

//a rectangle
type Rectangle struct {
	lowerLeft  Coordinates2D
	upperRight Coordinates2D
}

//NewRectangle create a new rectangle from any two opposing corners
func NewRectangle(firstCorner, secondCorner Coordinates2D) Rectangle {
	return Rectangle{
		lowerLeft: Coordinates2D{
			X: math.Min(firstCorner.X, secondCorner.X),
			Y: math.Min(firstCorner.Y, secondCorner.Y),
		},
		upperRight: Coordinates2D{
			X: math.Max(firstCorner.X, secondCorner.X),
			Y: math.Max(firstCorner.Y, secondCorner.Y),
		},
	}
}

//NewBoundingRectangle create a new rectangle which is the smallest rectangle to
//include all the given coordinates
func NewBoundingRectangle(coords []Coordinates2D) Rectangle {
	if len(coords) == 0 {
		return Rectangle{}
	}
	ret := Rectangle{
		lowerLeft:  coords[0],
		upperRight: coords[0],
	}

	for _, coord := range coords {
		ret.upperRight.X = math.Max(ret.upperRight.X, coord.X)
		ret.upperRight.Y = math.Max(ret.upperRight.Y, coord.Y)
		ret.lowerLeft.X = math.Min(ret.lowerLeft.X, coord.X)
		ret.lowerLeft.Y = math.Min(ret.lowerLeft.Y, coord.Y)
	}

	return ret
}

//Width the width of the rectangle
func (self Rectangle) Width() float64 {
	return self.upperRight.X - self.lowerLeft.X
}

//Height the height of the rectangle
func (self Rectangle) Height() float64 {
	return self.upperRight.Y - self.lowerLeft.Y
}

//Center the central point of the rectangle
func (self Rectangle) Center() Coordinates2D {
	return self.upperRight.Add(self.lowerLeft).Divide(2.0)
}

//Expand return a new Rectangle with the same center point but whose width and
//height is increased by the given positive amount
func (self Rectangle) Expand(amount float64) Rectangle {
	delta := Coordinates2D{X: amount / 2.0, Y: amount / 2.0}
	return Rectangle{
		lowerLeft:  self.lowerLeft.Subtract(delta),
		upperRight: self.upperRight.Add(delta),
	}
}

//Contains return true if the given coordinate is within the rectangle
func (self Rectangle) Contains(pos Coordinates2D) bool {
	return pos.X > self.lowerLeft.X && pos.X < self.upperRight.X &&
		pos.Y > self.lowerLeft.Y && pos.Y < self.upperRight.Y
}
