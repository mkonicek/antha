// liquidhandling/lhinterfaces.go: Part of the Antha language
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
// contact license@antha-lang.Org or write to the Antha team c/o
// Synthace Ltd. The London Bioscience Innovation Centre
// 2 Royal College St, London NW1 0NH UK

// defines types for dealing with liquid handling requests
package wtype

import "math"

//BBox is a simple LHObject representing a bounding box,
//useful for checking if there's stuff in the way
type BBox struct {
	Position Coordinates
	Size     Coordinates
}

//Equals defines an equivalence relation over bounding boxes
//returns true iff position and size are both equal
func (bb BBox) Equals(bb2 BBox) bool {
	return bb.Position.Equals(bb2.Position) && bb.Size.Equals(bb2.Size)
}

func NewBBox(pos, Size Coordinates) *BBox {
	if Size.X < 0. {
		pos.X = pos.X + Size.X
		Size.X = -Size.X
	}
	if Size.Y < 0. {
		pos.Y = pos.Y + Size.Y
		Size.Y = -Size.Y
	}
	if Size.Z < 0. {
		pos.Z = pos.Z + Size.Z
		Size.Z = -Size.Z
	}
	r := BBox{pos, Size}
	return &r
}

func NewBBox6f(pos_x, pos_y, pos_z, Size_x, Size_y, Size_z float64) *BBox {
	return NewBBox(Coordinates{pos_x, pos_y, pos_z},
		Coordinates{Size_x, Size_y, Size_z})
}

func NewXBox4f(pos_y, pos_z, Size_y, Size_z float64) *BBox {
	return NewBBox(Coordinates{-math.MaxFloat64 / 2., pos_y, pos_z},
		Coordinates{math.MaxFloat64, Size_y, Size_z})
}

func NewYBox4f(pos_x, pos_z, Size_x, Size_z float64) *BBox {
	return NewBBox(Coordinates{pos_x, -math.MaxFloat64 / 2., pos_z},
		Coordinates{Size_x, math.MaxFloat64, Size_z})
}

func NewZBox4f(pos_x, pos_y, Size_x, Size_y float64) *BBox {
	return NewBBox(Coordinates{pos_x, pos_y, -math.MaxFloat64 / 2.},
		Coordinates{Size_x, Size_y, math.MaxFloat64})
}

func (self BBox) GetPosition() Coordinates {
	return self.Position
}
func (self BBox) ZMax() float64 {
	return self.Position.Z + self.Size.Z
}

func (self BBox) GetSize() Coordinates {
	return self.Size
}

func (self *BBox) SetPosition(c Coordinates) {
	self.Position = c
}

func (self *BBox) SetSize(c Coordinates) {
	self.Size = c
}

func (self BBox) Contains(rhs Coordinates) bool {
	return (rhs.X >= self.Position.X && rhs.X < self.Position.X+self.Size.X &&
		rhs.Y >= self.Position.Y && rhs.Y < self.Position.Y+self.Size.Y &&
		rhs.Z >= self.Position.Z && rhs.Z < self.Position.Z+self.Size.Z)
}

//IntersectsBox checks for bounding box intersection
func (self BBox) IntersectsBox(rhs BBox) bool {
	//test a single dimension.
	//(a,b) are the start and end of the first Position
	//(c,d) are the start and end of the second pos
	// assert(a > b  and  d > c)
	f := func(a, b, c, d float64) bool {
		return !(c >= b || d <= a)
	}

	s := self.Position.Add(self.Size)
	r := rhs.GetPosition().Add(rhs.GetSize())
	return (f(self.Position.X, s.X, rhs.GetPosition().X, r.X) &&
		f(self.Position.Y, s.Y, rhs.GetPosition().Y, r.Y) &&
		f(self.Position.Z, s.Z, rhs.GetPosition().Z, r.Z))
}

//IntersectsPoint
func (self BBox) IntersectsPoint(rhs Coordinates) bool {
	return (rhs.X >= self.Position.X && rhs.X < self.Position.X+self.Size.X &&
		rhs.Y >= self.Position.Y && rhs.Y < self.Position.Y+self.Size.Y &&
		rhs.Z >= self.Position.Z && rhs.Z < self.Position.Z+self.Size.Z)
}
