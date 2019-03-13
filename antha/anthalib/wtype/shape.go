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
	"fmt"
	"strings"

	"github.com/antha-lang/antha/antha/anthalib/wunit"

	"math"
)

type ShapeTypeID string

const (
	CylinderShape  ShapeTypeID = "cylinder"
	CircleShape    ShapeTypeID = "circle"
	RoundShape     ShapeTypeID = "round"
	SphereShape    ShapeTypeID = "sphere"
	SquareShape    ShapeTypeID = "square"
	BoxShape       ShapeTypeID = "box"
	RectangleShape ShapeTypeID = "rectangle"
)

var isRound = map[ShapeTypeID]bool{
	CylinderShape:  true,
	CircleShape:    true,
	RoundShape:     true,
	SphereShape:    true,
	SquareShape:    false,
	BoxShape:       false,
	RectangleShape: false,
}

// IsRound returns true if wells with this shapetype are round, otherwise they are assumed to square or rectangular
func (s ShapeTypeID) IsRound() bool {
	if r, ok := isRound[s]; ok {
		return r
	}
	panic(fmt.Sprintf("ShapeType %s not in isRound", s))
}

type Shape struct {
	ShapeName  ShapeTypeID
	LengthUnit string
	H          float64
	W          float64
	D          float64
}

func (sh *Shape) Equals(sh2 *Shape) bool {
	return sh.ShapeName == sh2.ShapeName && sh.LengthUnit == sh2.LengthUnit && sh.H == sh2.H && sh.W == sh2.W && sh.D == sh2.D
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
	return &(Shape{sh.ShapeName, sh.LengthUnit, sh.H, sh.W, sh.D})
}

func (sh *Shape) String() string {
	return fmt.Sprintf("%s [%fx%fx%f]", sh.ShapeName, sh.H, sh.W, sh.D)
}

func (sh *Shape) MaxCrossSectionalArea() (area wunit.Area, err error) {

	shapename := strings.ToLower(string(sh.ShapeName))
	var areaunit string
	if sh.LengthUnit == "mm" {
		areaunit = "mm^2" //sh.LengthUnit + `^` + strconv.Itoa(2)
	} else {
		err = fmt.Errorf("sh.Lengthunit = %s", sh.LengthUnit)
		fmt.Println(err.Error())
	}
	var circular bool
	var boxlike bool

	if shapename == "circle" || shapename == "cylinder" || shapename == "round" || shapename == "sphere" {
		circular = true
	} else if shapename == "square" || shapename == "rectangle" || shapename == "box" {
		boxlike = true
	}

	if circular /*&& sh.Height() == sh.Width() */ {
		area = wunit.NewArea(math.Pi*(sh.W/2)*(sh.W/2), areaunit)
	} else if boxlike {
		area = wunit.NewArea(sh.H*sh.W, areaunit)
	} else {
		err = fmt.Errorf("No method to work out cross sectional area for shape \"%s\" yet Circular? %t", sh.ShapeName, circular)
	}
	return
}

func (sh *Shape) Volume() (volume wunit.Volume, err error) {

	shapename := strings.ToLower(string(sh.ShapeName))
	var volumeunit string
	if sh.LengthUnit == "mm" {
		volumeunit = "ul"
	} else {
		err = fmt.Errorf("can't handle conversion of %s to volume unit yet", sh.LengthUnit)
	}

	var cylinder bool
	var boxlike bool

	if shapename == "cylinder" {
		cylinder = true
	} else if shapename == "square" || shapename == "rectangle" || shapename == "box" {
		boxlike = true
	}

	if cylinder && sh.Height().EqualTo(sh.Width()) {
		volume = wunit.NewVolume(math.Pi*sh.H*sh.H*sh.D, volumeunit)
	} else if boxlike {
		volume = wunit.NewVolume(sh.H*sh.W*sh.D, volumeunit)
	} else {
		err = fmt.Errorf("No method to work out volume for shape %s yet. ", sh.ShapeName)
	}
	return
}

func NewShape(name ShapeTypeID, lengthunit string, h, w, d float64) *Shape {
	sh := Shape{name, lengthunit, h, w, d}
	return &sh
}

func NewNilShape() *Shape {
	sh := Shape{"", "", 0.0, 0.0, 0.0}
	return &sh
}

func (sh *Shape) IsZero() bool {
	return len(sh.ShapeName) == 0
}
