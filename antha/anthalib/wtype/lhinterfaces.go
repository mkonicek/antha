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

import (
	"github.com/pkg/errors"

	"github.com/antha-lang/antha/antha/anthalib/wunit"
)

type Named interface {
	GetName() string
}

//NameOf
func NameOf(o interface{}) string {
	if on, ok := o.(Named); ok {
		return on.GetName()
	}
	return "<unnamed>"
}

type Identifiable interface {
	GetID() string
}

//IDOf
func IDOf(o interface{}) string {
	if idf, ok := o.(Identifiable); ok {
		return idf.GetID()
	}
	return "<No-ID>"
}

type Typed interface {
	GetType() string
}

//TypeOf
func TypeOf(o interface{}) string {
	if ot, ok := o.(Typed); ok {
		return ot.GetType()
	}
	return "<unknown>"
}

//the class of thing it is, mainly for more helpful errors
type Classy interface {
	GetClass() string
}

//ClassOf
func ClassOf(o interface{}) string {
	if ot, ok := o.(Classy); ok {
		return ot.GetClass()
	}
	return "<unknown>"
}

//LHObject Provides a unified interface to physical size and location
//of objects that exist within a liquid handler
type LHObject interface {
	//GetPosition get the absolute position of the object (mm)
	GetPosition() Coordinates3D
	//GetSize get the size of the object (mm)
	GetSize() Coordinates3D
	//GetBoxIntersections get a list of LHObjects (can be this object or children) which intersect with the BBox
	GetBoxIntersections(BBox) []LHObject
	//GetPointIntersections get a list of LHObjects (can be this object or children) which intersect with the given point
	GetPointIntersections(Coordinates3D) []LHObject
	//SetOffset set the offset of the object relative to its parent (global if parent is nil)
	SetOffset(Coordinates3D) error
	//SetParent Store the offset of the object
	SetParent(LHObject) error
	//ClearParent Clear the parent of the object
	ClearParent()
	//GetParent
	GetParent() LHObject
	//GetID
	GetID() string
	//Dup duplicate the object, optionally keeping the IDs identical
	Duplicate(bool) LHObject
	//DimensionsString returns a string description of the position and size of the object and its children.
	//useful for debugging
	DimensionsString() string
}

//GetObjectRoot get the highest parent
func GetObjectRoot(o LHObject) LHObject {
	start := o
	for o.GetParent() != nil {
		o = o.GetParent()
		if o == start {
			panic("Infinite loop of LHObjects")
		}
	}
	return o
}

//get the origin for the objects coordinate system
func OriginOf(o LHObject) Coordinates3D {
	if p := o.GetParent(); p != nil {
		return p.GetPosition()
	}
	return Coordinates3D{}
}

//LHParent An LHObject that can hold other LHObjects
type LHParent interface {
	//GetChild get the child in the specified slot, nil if none. bool is false if the slot doesn't exists
	GetChild(string) (LHObject, bool)
	//GetSlotNames get a list of the slots
	GetSlotNames() []string
	//SetChild put the object in the slot
	SetChild(string, LHObject) error
	//Accepts test if slot can accept a certain class
	Accepts(string, LHObject) bool
	//GetSlotSize
	GetSlotSize(string) Coordinates2D
}

// WellReference used for specifying position within a well particularly for
// signify the intended locations for tip movements. Offsets are taken relative to these
type WellReference int

const (
	UnknownReference WellReference = iota - 1 // -1
	BottomReference                           //0
	TopReference                              //1
	LiquidReference                           //2
)

// NewWellReference construct a WellReference from its descriptive name.
// returns UnkonwnReference and an error if the reference is not recognised
func NewWellReference(s string) (WellReference, error) {
	names := map[string]WellReference{
		"well_bottom":  BottomReference,
		"well_top":     TopReference,
		"liquid_level": LiquidReference,
	}

	if r, ok := names[s]; ok {
		return r, nil
	}
	return UnknownReference, errors.Errorf(`unknown well reference "%s"`, s)
}

// String returns a descriptive name for the reference which is interpreted by the front end
func (self WellReference) String() string {
	names := []string{"well_bottom", "well_top", "liquid_level"}
	if int(self) >= 0 && int(self) < len(names) {
		return names[self]
	}
	return "unknown_reference"
}

func (self WellReference) AsInt() int {
	return int(self)
}

//Addressable unifies the interface to objects which have
//sub-components that can be addressed by WellCoords (e.g. "A1")
//for example tip-boxes, plates, etc
type Addressable interface {
	//AddressExists Do the given coordinates exist in the object?
	AddressExists(WellCoords) bool
	NRows() int
	NCols() int
	//GetChildByAddress Returns the object at the given well coords
	//nil if empty or position doesn't exist
	GetChildByAddress(WellCoords) LHObject
	//CoordsToWellCoords Convert Real world coordinates
	//(relative to the object origin) to WellCoords.
	//The returned WellCoords should be the closest
	//addressable location to the coorinates, and shold only be
	//invalid if the object has no adressable locations (e.g. wells
	//or tips). The second return value gives the offset from the top
	//of the center of the well/tip to the given coordinate
	//(this leaves the caller to ascertain whether any mis-alignment
	//is acceptable)
	CoordsToWellCoords(Coordinates3D) (WellCoords, Coordinates3D)
	//WellCoordsToCoords Get the physical location of an addressable
	//position relative to the object origin.
	//WellCoords should be valid in the object, or the bool will
	//return false and Coordinates are undefined.
	//WellReference is the position within a well.
	//Requesting LiquidReference on a LHTipbox will return false
	WellCoordsToCoords(WellCoords, WellReference) (Coordinates3D, bool)
}

type Targetted interface {
	//GetTargetOffset Gets the well target location for the numbered channel of the named adaptor
	GetTargetOffset(string, int) Coordinates3D
	//GetTargets return all the defined targets for the named adaptor
	GetTargets(string) []Coordinates3D
}

//LHContainer a tip or a well or something that holds liquids
type LHContainer interface {
	Contents() *Liquid
	CurrentVolume() wunit.Volume
	ResidualVolume() wunit.Volume
	//WorkingVolume = CurrentVolume - ResidualVolume
	CurrentWorkingVolume() wunit.Volume
	//Add to the container
	AddComponent(*Liquid) error
	//Set the contents of the container
	SetContents(*Liquid) error
	//Remove from the container
	RemoveVolume(wunit.Volume) (*Liquid, error)
}
