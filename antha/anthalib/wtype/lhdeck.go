// liquidhandling/lhdeck.Go: Part of the Antha language
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
	"fmt"
	"math"
	"sort"
	"strings"
)

type deckSlot struct {
	contents LHObject
	position Coordinates3D
	size     Coordinates2D
	accepts  []string
}

func newDeckSlot(position Coordinates3D, size Coordinates2D) *deckSlot {
	r := deckSlot{nil, position, size, make([]string, 0)}
	return &r
}

func (self *deckSlot) Fits(size Coordinates3D) bool {
	//.1mm tolerance for potential numerical error
	return math.Abs(self.size.X-size.X) < 0.1 &&
		math.Abs(self.size.Y-size.Y) < 0.1
}

func (self *deckSlot) AcceptsClass(class string) bool {
	for _, t := range self.accepts {
		if t == class {
			return true
		}
	}
	return false
}

func (self *deckSlot) SetAccepts(class string) {
	self.accepts = append(self.accepts, class)
}

func (self *deckSlot) GetAccepted() []string {
	return self.accepts
}

func (self *deckSlot) IsBelow(point Coordinates3D) bool {
	return (point.X >= self.position.X && point.X <= self.position.X+self.size.X &&
		point.Y >= self.position.Y && point.Y <= self.position.Y+self.size.Y)
}

//Duplicate copy the deckSlot and contained objects, optionally keeping IDs unchanged
func (self *deckSlot) Duplicate(keepIDs bool) *deckSlot {
	accepts := make([]string, len(self.accepts))
	copy(accepts, self.accepts)
	return &deckSlot{
		contents: self.contents.Duplicate(keepIDs),
		position: self.position,
		size:     self.size,
		accepts:  accepts,
	}
}

//LHDeck Represents a robot deck
type LHDeck struct {
	name     string
	mfg      string
	decktype string
	id       string
	slots    map[string]*deckSlot
}

//NewLHDeck make a new deck
func NewLHDeck(name, mfg, decktype string) *LHDeck {
	r := LHDeck{name, mfg, decktype, GetUUID(), make(map[string]*deckSlot)}
	return &r
}

//@implements Named

func (self *LHDeck) GetName() string {
	return self.name
}

//@implements Typed

func (self *LHDeck) GetType() string {
	return self.decktype
}

func (self *LHDeck) GetClass() string {
	return "deck"
}

//DimensionsString returns a string description of the position and size of the object and its children.
func (self *LHDeck) DimensionsString() string {
	ret := []string{fmt.Sprintf("Deck \"%s\" with %d slots:", self.GetName(), len(self.slots))}

	slotNames := self.GetSlotNames()
	sort.Strings(slotNames)

	for _, slotName := range slotNames {
		slot := self.slots[slotName]
		ret = append(ret, fmt.Sprintf("\tslot \"%s\" at %v+%v", slotName, slot.position, slot.size))
		if slot.contents == nil {
			ret = append(ret, "\t\tempty")
			continue
		}
		ds := slot.contents.DimensionsString()
		ret = append(ret, "\t\t"+strings.Replace(ds, "\n", "\n\t\t", -1))
	}

	return strings.Join(ret, "\n")
}

func (self *LHDeck) GetManufacturer() string {
	return self.mfg
}

func (self *LHDeck) GetID() string {
	return self.id
}

//@implements LHObject

func (self *LHDeck) GetPosition() Coordinates3D {
	return Coordinates3D{}
}

//zero size
func (self *LHDeck) GetSize() Coordinates3D {
	return Coordinates3D{}
}

func (self *LHDeck) GetBoxIntersections(box BBox) []LHObject {
	ret := []LHObject{}
	for _, ds := range self.slots {
		if ds.contents != nil {
			ret = append(ret, ds.contents.GetBoxIntersections(box)...)
		}
	}
	return ret
}

func (self *LHDeck) GetPointIntersections(point Coordinates3D) []LHObject {
	ret := []LHObject{}
	for _, ds := range self.slots {
		if ds.contents != nil {
			ret = append(ret, ds.contents.GetPointIntersections(point)...)
		}
	}
	return ret
}

func (self *LHDeck) SetOffset(o Coordinates3D) error {
	return fmt.Errorf("Can't set offset for deck \"%s\"", self.GetName())
}

func (self *LHDeck) SetParent(o LHObject) error {
	return fmt.Errorf("Can't set deck \"%s\"'s parent, tried to set to %s \"%s\"",
		self.GetName(), ClassOf(o), NameOf(o))
}

//@implement LHObject
func (self *LHDeck) ClearParent() {}

func (self *LHDeck) GetParent() LHObject {
	return nil
}

//Duplicate copy the LHObject, optionally keeping IDs unchanged
func (self *LHDeck) Duplicate(keepIDs bool) LHObject {

	slots := make(map[string]*deckSlot)
	for name, sl := range self.slots {
		slots[name] = sl.Duplicate(keepIDs)
	}

	uuid := self.id
	if !keepIDs {
		uuid = GetUUID()
	}

	r := &LHDeck{self.name, self.mfg, self.decktype, uuid, slots}
	return r
}

//@implements LHParent
func (self *LHDeck) GetChild(name string) (LHObject, bool) {
	if ds, ok := self.slots[name]; ok {
		return ds.contents, true
	}
	return nil, false
}

func (self *LHDeck) GetSlotNames() []string {
	ret := make([]string, 0, len(self.slots))
	for key := range self.slots {
		ret = append(ret, key)
	}
	return ret
}

func (self *LHDeck) GetSlotPosition(name string) Coordinates3D {
	if self == nil {
		return Coordinates3D{}
	}
	if ds, ok := self.slots[name]; ok {
		return ds.position
	}
	return Coordinates3D{}
}

func (self *LHDeck) GetSlotContaining(obj LHObject) string {
	if self == nil {
		return ""
	}
	for n, sl := range self.slots {
		if IDOf(sl.contents) == IDOf(obj) {
			return n
		}
	}
	return ""
}

func (self *LHDeck) SetChild(name string, child LHObject) error {
	if ds, ok := self.slots[name]; !ok {
		return fmt.Errorf("Cannot put %s \"%s\" at unknown slot \"%s\"", ClassOf(child), NameOf(child), name)
	} else if !ds.Fits(child.GetSize()) {
		return fmt.Errorf("Footprint of %s \"%s\"[%vmm x %vmm] doesn't fit slot \"%s\"[%vmm x %vmm]",
			ClassOf(child), NameOf(child), child.GetSize().X, child.GetSize().Y,
			name, ds.size.X, ds.size.Y)
	} else if !ds.AcceptsClass(ClassOf(child)) {
		return fmt.Errorf("Slot \"%s\" can't accept %s \"%s\", only %s allowed",
			name, ClassOf(child), NameOf(child), strings.Join(ds.GetAccepted(), ","))
	} else if ds.contents != nil {
		return fmt.Errorf("Couldn't add %s \"%s\" to location \"%s\" which already contains %s \"%s\"",
			ClassOf(child), NameOf(child), name, ClassOf(ds.contents), NameOf(ds.contents))
	} else {
		ds.contents = child
		if err := child.SetParent(self); err != nil {
			return err
		}

		if err := child.SetOffset(ds.position); err != nil {
			return err
		}
	}
	return nil
}

func (self *LHDeck) Clear(name string) error {
	if ds, ok := self.slots[name]; !ok {
		return fmt.Errorf("Cannot clear from unknown slot \"%s\"", name)
	} else {
		ds.contents = nil
	}
	return nil
}

func (self *LHDeck) Accepts(name string, child LHObject) bool {
	if ds, ok := self.slots[name]; ok {
		return ds.Fits(child.GetSize()) && ds.AcceptsClass(ClassOf(child))
	}
	return false
}

func (self *LHDeck) GetSlotSize(name string) Coordinates2D {
	return self.slots[name].size
}

//LHDeck specific methods

func (self *LHDeck) AddSlot(name string, position Coordinates3D, size Coordinates2D) {
	self.slots[name] = newDeckSlot(position, size)
}

func (self *LHDeck) SetSlotAccepts(name string, class string) {
	if sl, ok := self.slots[name]; ok {
		sl.SetAccepts(class)
	}
}

func (self *LHDeck) GetSlotNamesBelow(point Coordinates3D) []string {
	ret := make([]string, 0)
	for name, slot := range self.slots {
		if slot.IsBelow(point) {
			ret = append(ret, name)
		}
	}
	return ret
}

//get all objects above and below the point
func (self *LHDeck) GetVChildren(point Coordinates3D) []LHObject {
	//get all children in the same vertical plane
	box := NewBBox6f(point.X, point.Y, -math.MaxFloat64/2, 0, 0, math.MaxFloat64)
	return self.GetBoxIntersections(*box)

}

//Return the nearest object below the point, nil if none.
//The base of the object is used as reference, so e.g. a point within a well
//will return the plate
func (self *LHDeck) GetChildBelow(point Coordinates3D) LHObject {
	candidates := self.GetVChildren(point)
	//find the closest that's below
	z_off_min := math.MaxFloat64
	z_off_i := -1
	for i, c := range candidates {
		if z_off := (point.Z - (c.GetPosition().Z + c.GetSize().Z)); (point.Z-c.GetPosition().Z) > 0 && z_off < z_off_min {
			z_off_min = z_off
			z_off_i = i
		}
	}

	//len(candidates) == 0
	if z_off_i < 0 {
		return nil
	}
	return candidates[z_off_i]
}
