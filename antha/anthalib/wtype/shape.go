// wtype/shape.go: Part of the Antha language
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
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"math"

	"github.com/antha-lang/antha/antha/anthalib/wunit"
)

type shapeType struct {
	Name  string
	Round bool
}

func newShapeType(name string, isRound bool) *shapeType {
	ret := &shapeType{
		Name:  name,
		Round: isRound,
	}
	shapeTypeByName[name] = ret
	return ret
}

// String returns a string description of the shape type
func (sti *shapeType) String() string {
	if sti == nil {
		return "unknown"
	}
	return sti.Name
}

// IsRound returns whether or not the shape is rounded
func (sti *shapeType) IsRound() bool {
	if sti == nil {
		return false
	}
	return sti.Round
}

// Equals returns true if the shape types are the same
func (sti *shapeType) Equals(rhs *shapeType) bool {
	return sti == rhs
}

var shapeTypeByName = make(map[string]*shapeType)

// shapeTypeFromName get a shape type by string name, required for deserialisation
// returns nil if the name is unknown
func ShapeTypeFromName(name string) *shapeType {
	return shapeTypeByName[name]
}

var (
	UnknownShape   = (*shapeType)(nil)
	CylinderShape  = newShapeType("cylinder", true)
	CircleShape    = newShapeType("circle", true)
	RoundShape     = newShapeType("round", true)
	SphereShape    = newShapeType("sphere", true)
	SquareShape    = newShapeType("square", false)
	BoxShape       = newShapeType("box", false)
	RectangleShape = newShapeType("rectangle", false)
	TrapezoidShape = newShapeType("trapezoid", false)
)

type Shape struct {
	Type       *shapeType
	LengthUnit string
	H          float64
	W          float64
	D          float64
}

func (sh *Shape) Equals(sh2 *Shape) bool {
	return sh.Type.Equals(sh2.Type) && sh.LengthUnit == sh2.LengthUnit && sh.H == sh2.H && sh.W == sh2.W && sh.D == sh2.D
}

// let shape implement geometry

func (sh *Shape) Height() wunit.Length { // y?
	return wunit.NewLength(sh.H, sh.LengthUnit)
}
func (sh *Shape) Width() wunit.Length { // X?
	return wunit.NewLength(sh.W, sh.LengthUnit)
}
func (sh *Shape) Depth() wunit.Length { // Z?
	return wunit.NewLength(sh.D, sh.LengthUnit)
}

func (sh *Shape) Dup() *Shape {
	return &Shape{
		Type:       sh.Type,
		LengthUnit: sh.LengthUnit,
		H:          sh.H,
		W:          sh.W,
		D:          sh.D,
	}
}

func (sh *Shape) String() string {
	return fmt.Sprintf("%s [%fx%fx%f]", sh.Type, sh.H, sh.W, sh.D)
}

func (sh *Shape) MaxCrossSectionalArea() (wunit.Area, error) {

	// attempt to get H and W in mm
	// nb. "Width" and "Height" are X and Y. Z is "Depth"
	if height, err := sh.Height().InStringUnit("mm"); err != nil {
		return wunit.ZeroArea(), errors.WithMessage(err, "while converting height to mm")
	} else if width, err := sh.Width().InStringUnit("mm"); err != nil {
		return wunit.ZeroArea(), errors.WithMessage(err, "while converting width to mm")
	} else if sh.Type.IsRound() {
		radius := width.RawValue() / 2.0
		return wunit.NewArea(math.Pi*radius*radius, "mm^2"), nil
	} else {
		return wunit.NewArea(height.RawValue()*width.RawValue(), "mm^2"), nil
	}
}

func (sh *Shape) Volume() (volume wunit.Volume, err error) {

	// attempt to get H, W and D in mm
	// nb. "Width" and "Height" are X and Y. Z is "Depth"
	if height, err := sh.Height().InStringUnit("mm"); err != nil {
		return wunit.ZeroVolume(), errors.WithMessage(err, "while converting height to mm")
	} else if width, err := sh.Width().InStringUnit("mm"); err != nil {
		return wunit.ZeroVolume(), errors.WithMessage(err, "while converting width to mm")
	} else if depth, err := sh.Depth().InStringUnit("mm"); err != nil {
		return wunit.ZeroVolume(), errors.WithMessage(err, "while converting depth to mm")
	} else if sh.Type.IsRound() {
		// assume the top shape is an ellipse
		return wunit.NewVolume(math.Pi*height.RawValue()*width.RawValue()*depth.RawValue(), "mm^3"), nil
	} else {
		// assume the shape is a cuboid
		return wunit.NewVolume(height.RawValue()*width.RawValue()*depth.RawValue(), "mm^3"), nil
	}
}

func NewShape(shapetype *shapeType, lengthunit string, h, w, d float64) *Shape {
	return &Shape{
		Type:       shapetype,
		LengthUnit: lengthunit,
		H:          h,
		W:          w,
		D:          d,
	}
}

func (st *Shape) MarshalJSON() ([]byte, error) {
	// avoid calling in a loop
	type ShapeAlias Shape
	return json.Marshal(struct {
		ShapeAlias
		Type string `json:"type"` // serialize type as string
	}{
		ShapeAlias: ShapeAlias(*st),
		Type:       st.Type.String(),
	})
}

func (st *Shape) UnmarshalJSON(bs []byte) error {
	// avoid calling in a loop
	type ShapeAlias Shape
	var s struct {
		ShapeAlias
		Type string `json:"type"` // serialize type as string
	}

	if err := json.Unmarshal(bs, &s); err != nil {
		return errors.WithMessage(err, "unmarshalling shape type")
	}

	*st = Shape(s.ShapeAlias)

	// set the correct *shapeType to preserve pointer equality
	st.Type = ShapeTypeFromName(s.Type)
	return nil
}
