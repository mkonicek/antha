// liquidhandling/lhtypes.Go: Part of the Antha language
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

package wtype

import (
	"fmt"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
)

//  TODO remove BBox once shape implements LHObject
type LHTip struct {
	ID       string
	Type     string
	Mnfr     string
	Dirty    bool
	MaxVol   wunit.Volume
	MinVol   wunit.Volume
	Shape    *Shape
	Bounds   BBox
	parent   LHObject `gotopb:"-"`
	contents *LHComponent
	Filtered bool
}

//@implement Named
func (self *LHTip) GetName() string {
	if self == nil {
		return "<nil>"
	}
	if addr, ok := self.parent.(Addressable); ok {
		pos := self.GetPosition().Add(self.GetSize().Multiply(0.5))
		wc, _ := addr.CoordsToWellCoords(pos)
		return fmt.Sprintf("%s@%s", wc.FormatA1(), NameOf(self.parent))
	}
	return fmt.Sprintf("%s_%s", self.Mnfr, self.Type)
}

func (self *LHTip) GetID() string {
	return self.ID
}

//@implement Typed
func (self *LHTip) GetType() string {
	if self == nil {
		return "<nil>"
	}
	return self.Type
}

//@implement Classy
func (self *LHTip) GetClass() string {
	return "tip"
}

//@implement LHObject
func (self *LHTip) GetPosition() Coordinates {
	return OriginOf(self).Add(self.Bounds.GetPosition())
}

//@implement LHObject
func (self *LHTip) GetSize() Coordinates {
	return self.Bounds.GetSize()
}

//@implement LHObject
func (self *LHTip) GetBoxIntersections(box BBox) []LHObject {
	box.SetPosition(box.GetPosition().Subtract(OriginOf(self)))
	if self.Bounds.IntersectsBox(box) {
		return []LHObject{self}
	}
	return nil
}

//@implement LHObject
func (self *LHTip) GetPointIntersections(point Coordinates) []LHObject {
	point = point.Subtract(point)
	//TODO more accurate intersection detection with Shape
	if self.Bounds.IntersectsPoint(point) {
		return []LHObject{self}
	}
	return nil
}

//@implement LHObject
func (self *LHTip) SetOffset(point Coordinates) error {
	self.Bounds.SetPosition(point)
	return nil
}

//@implement LHObject
func (self *LHTip) SetParent(o LHObject) error {
	//parent should be LHTipbox (should accept LHAdaptor, but it doesn't implement LHObject yet)
	if _, ok := o.(*LHTipbox); ok {
		self.parent = o
		return nil
	}
	return fmt.Errorf("Cannot set %s \"%s\" as parent of tip", ClassOf(o), NameOf(o))
}

//@implement LHObject
func (self *LHTip) ClearParent() {
	self.parent = nil
}

//@implement LHObject
func (self *LHTip) GetParent() LHObject {
	return self.parent
}

//Duplicate copies an LHObject
func (self *LHTip) Duplicate(keepIDs bool) LHObject {
	return self.dup(keepIDs)
}

func (tip *LHTip) GetParams() *LHChannelParameter {
	// be safe
	if tip.IsNil() {
		return nil
	}

	lhcp := LHChannelParameter{Name: tip.Type + "Params", Minvol: tip.MinVol, Maxvol: tip.MaxVol, Multi: 1, Independent: false, Orientation: LHVChannel}
	return &lhcp
}

func (tip *LHTip) IsNil() bool {
	if tip == nil || tip.Type == "" || tip.MaxVol.IsZero() || tip.MinVol.IsZero() {
		return true
	}
	return false
}

//Dup copy the tip generating a new ID
func (tip *LHTip) Dup() *LHTip {
	return tip.dup(false)
}

//Dup copy the tip keeping the previous ID
func (tip *LHTip) DupKeepID() *LHTip {
	return tip.dup(true)
}

func (tip *LHTip) dup(keepIDs bool) *LHTip {
	t := NewLHTip(tip.Mnfr, tip.Type, tip.MinVol.RawValue(), tip.MaxVol.RawValue(), tip.MinVol.Unit().PrefixedSymbol(), tip.Filtered, tip.Shape.Dup())
	t.Dirty = tip.Dirty
	t.contents = tip.Contents().Dup()
	t.Bounds = tip.Bounds

	if keepIDs {
		t.ID = tip.ID
	}

	return t
}

func NewLHTip(mfr, ttype string, minvol, maxvol float64, volunit string, filtered bool, shape *Shape) *LHTip {
	lht := LHTip{
		GetUUID(),
		ttype,
		mfr,
		false, //dirty
		wunit.NewVolume(maxvol, volunit),
		wunit.NewVolume(minvol, volunit),
		shape,
		BBox{Coordinates{}, Coordinates{
			shape.Height().ConvertToString("mm"), //not a mistake, Shape currently has height&width as
			shape.Width().ConvertToString("mm"),  // XY coordinates and Depth as Z
			shape.Depth().ConvertToString("mm"),
		}},
		nil,
		NewLHComponent(),
		filtered,
	}

	return &lht
}

func CopyTip(tt LHTip) *LHTip {
	return &tt
}

//DimensionsString returns a string description of the position and size of the object and its children.
func (self *LHTip) DimensionsString() string {
	if self == nil {
		return "no tip"
	}
	return fmt.Sprintf("Tip %s at %v+%v", self.GetName(), self.GetPosition(), self.GetSize())
}

//@implement LHContainer
func (self *LHTip) Contents() *LHComponent {
	//Only happens with dodgy tip initialization
	if self.contents == nil {
		self.contents = NewLHComponent()
	}
	return self.contents
}

//@implement LHContainer
func (self *LHTip) CurrentVolume() wunit.Volume {
	return self.contents.Volume()
}

//@implement LHContainer
func (self *LHTip) ResidualVolume() wunit.Volume {
	//currently not really supported
	return wunit.NewVolume(0, "ul")
}

//@implement LHContainer
func (self *LHTip) CurrentWorkingVolume() wunit.Volume {
	return self.contents.Volume()
}

//@implement LHContainer
func (self *LHTip) AddComponent(v *LHComponent) error {
	fv := self.CurrentVolume()
	fv.Add(v.Volume())

	self.contents.Mix(v)

	if fv.GreaterThan(self.MaxVol) {
		return fmt.Errorf("Tip %s overfull, contains %v and maximum is %v", self.GetName(), fv, self.MaxVol)
	}
	if fv.LessThan(self.MinVol) {
		return fmt.Errorf("Added less than minimum volume to %s, contains %v and minimum working volume is %v", self.GetName(), fv, self.MinVol)
	}
	return nil
}

//SetContents set the contents of the tip, returns an error if the tip is overfilled
func (self *LHTip) SetContents(v *LHComponent) error {
	if v.Volume().GreaterThan(self.MaxVol) {
		return fmt.Errorf("Tip %s overfull, contains %v and maximum is %v", self.GetName(), v.Volume(), self.MaxVol)
	}
	if v.Volume().LessThan(self.MinVol) {
		return fmt.Errorf("Added less than minimum volume to %s, contains %v and minimum working volume is %v", self.GetName(), v.Volume(), self.MinVol)
	}

	self.contents = v
	return nil
}

//@implement LHContainer
func (self *LHTip) RemoveVolume(v wunit.Volume) (*LHComponent, error) {
	if v.GreaterThan(self.CurrentWorkingVolume()) {
		return nil, fmt.Errorf("Requested removal of %v from tip %s which only has %v working volume", v, self.GetName(), self.CurrentWorkingVolume())
	}
	ret := self.contents.Dup()
	ret.Vol = v.ConvertToString("ul")
	self.contents.Remove(v)
	return ret, nil
}
