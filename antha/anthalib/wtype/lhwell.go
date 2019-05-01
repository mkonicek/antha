// liquidhandling/lhwell.Go: Part of the Antha language
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
	"strings"

	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/antha/anthalib/wutil"
)

type WellBottomType int

const (
	FlatWellBottom WellBottomType = iota
	UWellBottom
	VWellBottom
)

var WellBottomNames []string = []string{"flat", "U", "V"}

func (bt WellBottomType) String() string {
	return WellBottomNames[bt]
}

// structure representing a well on a microplate - description of a destination
type LHWell struct {
	ID        string
	Inst      string
	Crds      WellCoords
	MaxVol    float64 //Maximum total capacity of the well
	WContents *Liquid
	Rvol      float64 //Residual volume which can't be removed from the well
	WShape    *Shape
	Bottom    WellBottomType
	Bounds    BBox
	Bottomh   float64
	Extra     map[string]interface{}
	Plate     LHObject `gotopb:"-" json:"-"`
}

//@implement Named
func (self *LHWell) GetName() string {
	return fmt.Sprintf("%s@%s", self.Crds.FormatA1(), NameOf(self.Plate))
}

func (self *LHWell) GetID() string {
	return self.ID
}

//@implement Typed
func (self *LHWell) GetType() string {
	return fmt.Sprintf("well_in_%s", TypeOf(self.Plate))
}

//@implement Classy
func (self *LHWell) GetClass() string {
	return "well"
}

func (self *LHWell) GetWellCoords() WellCoords {
	if self == nil {
		return ZeroWellCoords()
	}
	return self.Crds
}

//@implement LHObject
func (self *LHWell) GetPosition() Coordinates3D {
	return OriginOf(self).Add(self.Bounds.GetPosition())
}

//@implement LHObject
func (self *LHWell) GetSize() Coordinates3D {
	return self.Bounds.GetSize()
}

func (self *LHWell) GetVolumeUnit() string {
	return "ul"
}

//@implement LHObject
func (self *LHWell) GetBoxIntersections(box BBox) []LHObject {
	//relative box
	box.SetPosition(box.GetPosition().Subtract(OriginOf(self)))
	if self.Bounds.IntersectsBox(box) {
		return []LHObject{self}
	}
	return nil
}

//@implement LHObject
func (self *LHWell) GetPointIntersections(point Coordinates3D) []LHObject {
	//relative point
	point = point.Subtract(OriginOf(self))
	//At some point this should be called self.shape for a more accurate intersection test
	//see branch shape-changes
	if self.Bounds.IntersectsPoint(point) {
		return []LHObject{self}
	}
	return nil
}

//@implement LHObject
func (self *LHWell) SetOffset(point Coordinates3D) error {
	self.Bounds.SetPosition(point)
	return nil
}

//@implement LHObject
func (self *LHWell) SetParent(p LHObject) error {
	//Seems unlikely, but I suppose wells that you can take from one plate and insert
	//into another could be feasible with some funky labware
	if plate, ok := p.(*Plate); ok {
		self.Plate = plate
		return nil
	}
	if tb, ok := p.(*LHTipwaste); ok {
		self.Plate = tb
		return nil
	}
	return fmt.Errorf("Cannot set well parent to %s \"%s\", only plates allowed", ClassOf(p), NameOf(p))
}

//@implement LHObject
func (self *LHWell) ClearParent() {
	self.Plate = nil
}

//@implement LHObject
func (self *LHWell) GetParent() LHObject {
	if self == nil {
		return nil
	}
	return self.Plate
}

//Duplicate copies an LHObject
func (self *LHWell) Duplicate(keepIDs bool) LHObject {
	return self.dup(keepIDs)
}

//DimensionsString returns a string description of the position and size of the object and its children.
func (self *LHWell) DimensionsString() string {
	if self == nil {
		return "nil well"
	}
	return fmt.Sprintf("Well %s at %v+%v", self.GetName(), self.GetPosition(), self.GetSize())
}

func (w LHWell) String() string {
	return fmt.Sprintf(
		`LHWELL{
ID        : %s,
Inst      : %s,
Crds      : %s,
MaxVol    : %g ul,
WContents : %v,
Rvol      : %g ul,
WShape    : %v,
Bottom    : %s,
size      : [%v x %v x %v]mm,
Bottomh   : %g,
Extra     : %v,
}`,
		w.ID,
		w.Inst,
		w.Crds.FormatA1(),
		w.MaxVol,
		w.WContents,
		w.Rvol,
		w.WShape,
		WellBottomNames[w.Bottom],
		w.GetSize().X,
		w.GetSize().Y,
		w.GetSize().Z,
		w.Bottomh,
		w.Extra,
	)
}

func (w *LHWell) Contents() *Liquid {
	if w == nil {
		return nil
	}

	if w.WContents == nil {
		w.WContents = NewLHComponent()
	}

	return w.WContents
}

func (w *LHWell) SetContents(newContents *Liquid) error {
	if w == nil {
		return nil
	}
	maxVol := w.MaxVolume()
	if newContents.Volume().GreaterThan(maxVol) {
		//HJK: Disabling overflow errors until CarryVolume issues are resolved
		//return LHError(LH_ERR_VOL,
		//	fmt.Sprintf("Cannot set %s as contents of well %s as maximum volume is %s", newContents.GetName(), w.GetName(), maxVol)s

		fmt.Printf("setting %s as contents of well %s even though maximum volume is %s\n", newContents.GetName(), w.GetName(), maxVol)
	}

	w.WContents = newContents
	w.updateComponentLocation()
	return nil
}

//CurrentVolume return the volume of the component currently in the well
func (w *LHWell) CurrentVolume() wunit.Volume {
	if w == nil {
		return wunit.ZeroVolume()
	}
	return w.Contents().Volume()
}

//CurrentWorkingVolume return the available working volume in the well - i.e. current volume minus residual volume
func (w *LHWell) CurrentWorkingVolume() wunit.Volume {
	if w == nil {
		return wunit.ZeroVolume()
	}
	v := w.CurrentVolume()
	v.Subtract(w.ResidualVolume())

	if v.LessThan(wunit.ZeroVolume()) {
		return wunit.ZeroVolume()
	}
	return v
}

func (w *LHWell) ResidualVolume() wunit.Volume {
	if w == nil {
		return wunit.ZeroVolume()
	}
	v := wunit.NewVolume(w.Rvol, "ul")
	return v
}

//MaxVolume get the maximum working volume of the well
func (w *LHWell) MaxVolume() wunit.Volume {
	if w == nil {
		return wunit.ZeroVolume()
	}
	return wunit.NewVolume(w.MaxVol, "ul")
}

//MaxWorkingVolume get the total maximum working volume in the well, i.e. the max volume minus residual volume
func (w *LHWell) MaxWorkingVolume() wunit.Volume {
	if w == nil {
		return wunit.ZeroVolume()
	}
	ret := w.MaxVolume()
	ret.Subtract(w.ResidualVolume())
	if ret.LessThan(wunit.ZeroVolume()) {
		return wunit.ZeroVolume()
	}
	return ret
}

//AddComponent add some liquid to the well
func (w *LHWell) AddComponent(c *Liquid) error {
	if w == nil {
		return nil
	}
	maxVol := w.MaxVolume()
	curVol := w.CurrentVolume()
	finalVol := c.Volume()
	finalVol.Add(curVol)

	if finalVol.GreaterThan(maxVol) {
		//HJK: Disabled overflow errors while CarryVolume issues are resolved
		//return fmt.Errorf("Cannot add %s to well \"%s\", well already contains %s and maximum volume is %s", c.GetName(), w.GetName(), curVol, maxVol)
		fmt.Printf("Adding %s to well \"%s\", even though well already contains %s and maximum volume is %s\n", c.Summarize(), w.GetName(), curVol, maxVol)
	}

	w.Contents().Mix(c)

	return nil
}

//RemoveVolume remove some liquid from the well
func (w *LHWell) RemoveVolume(v wunit.Volume) (*Liquid, error) {
	if w == nil {
		return nil, nil
	}

	// if the volume is too high we complain
	if v.GreaterThan(w.CurrentWorkingVolume()) {
		//HJK: Disabled underflow errors while CarryVolume issues are resolved
		//return nil, fmt.Errorf("requested %s from well \"%s\" which only contains %s working volume", v, w.GetName(), w.CurrentWorkingVolume())
		fmt.Printf("requested %s from well \"%s\" which only contains %s working volume and %s total volume\n",
			v, w.GetName(), w.CurrentWorkingVolume(), w.CurrentVolume())
	}

	ret := w.Contents().Dup()
	ret.Vol = v.ConvertToString("ul")

	w.Contents().Remove(v)
	return ret, nil
}

//RemoveCarry Remove the carry volume
func (w *LHWell) RemoveCarry(v wunit.Volume) {
	if w == nil {
		return
	}

	w.Contents().Remove(v)
}

//IsVolumeValid tests whether the volume in the well is within the allowable range
func (w *LHWell) IsVolumeValid() bool {
	if w == nil {
		return true
	}
	vol := w.CurrentVolume()

	return vol.LessThan(w.MaxVolume()) && !vol.LessThan(wunit.ZeroVolume())
}

//ValidateVolume validates that the volume in the well is within allowable range
func (w *LHWell) ValidateVolume() error {
	if w.IsVolumeValid() {
		return nil
	}

	return LHError(LH_ERR_VOL, fmt.Sprintf("well %s contains invalid volume %s, maximum volume is %s", w.GetName(), w.CurrentVolume(), w.MaxVolume()))
}

func (w *LHWell) PlateLocation() PlateLocation {
	if w == nil {
		return ZeroPlateLocation()
	}
	return w.WContents.PlateLocation()
}

//@implement Location

func (lhw *LHWell) Location_ID() string {
	if lhw == nil {
		return ""
	}
	return lhw.ID
}

func (lhw *LHWell) Location_Name() string {
	if lhw == nil {
		return ""
	}
	return NameOf(lhw.Plate)
}

func (lhw *LHWell) Shape() *Shape {
	if lhw == nil {
		return nil
	}
	if lhw.WShape == nil {
		// return the non-shape
		return nil
	}
	return lhw.WShape
}

// @implement Well
// @deprecate Well

func (w *LHWell) ContainerType() string {
	if w == nil {
		return ""
	}
	return TypeOf(w.Plate)
}

func (w *LHWell) Clear() {
	if w == nil {
		return
	}
	w.WContents = NewLHComponent()
	w.updateComponentLocation()
}

//IsEmpty returns true if the well contains nothing, though this does not mean that the working volume is greater than zero
func (w *LHWell) IsEmpty() bool {
	// nil wells are empty
	if w == nil {
		return true
	}

	tolerance := wunit.NewVolume(0.000001, "ul")

	return w.CurrentVolume().LessThan(tolerance)
}

//Clean resets the volume in the well so that it's empty
func (w *LHWell) Clean() {
	w.WContents.Clean()
	w.updateComponentLocation()
	newExtra := make(map[string]interface{})

	// some keys must be retained

	for k, v := range w.Extra {
		if isConstraintKey(k) || k == wellTargetKey || k == "IMSPECIAL" || k == "afvfunc" {
			newExtra[k] = v
		}
	}

	w.Extra = newExtra
}

func (w *LHWell) updateComponentLocation() {
	if w.WContents != nil {
		w.WContents.Loc = fmt.Sprintf("%s:%s", IDOf(w.Plate), w.Crds.FormatA1())
	}
}

// copy of instance
func (lhw *LHWell) Dup() *LHWell {
	return lhw.dup(false)
}

// copy of type
func (lhw *LHWell) CDup() *LHWell {
	if lhw == nil {
		return nil
	}
	cp := NewLHWell("ul", lhw.MaxVol, lhw.Rvol, lhw.Shape().Dup(), lhw.Bottom, lhw.GetSize().X, lhw.GetSize().Y, lhw.GetSize().Z, lhw.Bottomh, "mm")
	cp.Plate = lhw.Plate
	cp.Crds = lhw.Crds
	cp.WContents = lhw.Contents().Dup()

	for k, v := range lhw.Extra {
		cp.Extra[k] = v
	}

	return cp
}

func (lhw *LHWell) DupKeepIDs() *LHWell {
	return lhw.dup(true)
}

func (lhw *LHWell) dup(keep_ids bool) *LHWell {
	if lhw == nil {
		return nil
	}
	cp := NewLHWell("ul", lhw.MaxVol, lhw.Rvol, lhw.Shape().Dup(), lhw.Bottom, lhw.GetSize().X, lhw.GetSize().Y, lhw.GetSize().Z, lhw.Bottomh, "mm")
	cp.Plate = lhw.Plate
	cp.Crds = lhw.Crds
	cp.Bounds = lhw.Bounds

	if keep_ids {
		cp.ID = lhw.ID
	}

	// Dup here doesn't change ID
	cp.WContents = lhw.Contents().Dup()

	for k, v := range lhw.Extra {
		cp.Extra[k] = v
	}

	return cp
}

func (lhw *LHWell) CalculateMaxCrossSectionArea() (ca wunit.Area, err error) {
	if lhw == nil {
		return
	}
	return lhw.Shape().MaxCrossSectionalArea()
}

func (lhw *LHWell) SetAfVFunc(f string) {
	if lhw == nil {
		return
	}
	lhw.Extra["afvfunc"] = f
}

func (lhw *LHWell) GetAfVFunc() wutil.Func1Prm {
	if lhw == nil {
		return wutil.Quadratic{}
	}
	f, ok := lhw.Extra["afvfunc"]

	if !ok {
		return nil
	}
	x, err := wutil.UnmarshalFunc([]byte(f.(string)))
	if err != nil {
		panic(fmt.Sprintf("Can't unmarshal function, error: %s", err))
	}
	return x
}

func (lhw *LHWell) CalculateMaxVolume() (vol wunit.Volume, err error) {
	if lhw == nil {
		return wunit.ZeroVolume(), fmt.Errorf("Nil well has no max volume")
	}

	if lhw.Bottom == FlatWellBottom { // flat
		vol, err = lhw.Shape().Volume()
	} /*else if lhw.Bottom == UWellBottom { // round
		vol, err = lhw.Shape().Volume()
		// + additional calculation
	} else if lhw.Bottom == VWellBottom { // Pointed / v-shaped /pyramid
		vol, err = lhw.Shape().Volume()
		// + additional calculation
	}
	*/
	return
}

// make a new well structure
func NewLHWell(vunit string, vol, rvol float64, shape *Shape, bott WellBottomType, xdim, ydim, zdim, bottomh float64, dunit string) *LHWell {
	var well LHWell

	well.Plate = nil
	crds := ZeroWellCoords()
	well.WContents = NewLHComponent()
	well.updateComponentLocation()

	well.ID = GetUUID()
	well.Crds = crds
	well.MaxVol = wunit.NewVolume(vol, vunit).ConvertToString("ul")
	well.Rvol = wunit.NewVolume(rvol, vunit).ConvertToString("ul")
	well.WShape = shape.Dup()
	well.Bottom = bott
	well.Bounds = BBox{Coordinates3D{}, Coordinates3D{
		wunit.NewLength(xdim, dunit).ConvertToString("mm"),
		wunit.NewLength(ydim, dunit).ConvertToString("mm"),
		wunit.NewLength(zdim, dunit).ConvertToString("mm"),
	}}
	well.Bottomh = wunit.NewLength(bottomh, dunit).ConvertToString("mm")
	well.Extra = make(map[string]interface{})
	return &well
}

// this function is somewhat buggy... need to define its responsibilities better
func Get_Next_Well(plate *Plate, component *Liquid, curwell *LHWell) (*LHWell, bool) {
	vol := component.Vol

	it := NewAddressIterator(plate, ColumnWise, TopToBottom, LeftToRight, false)

	if curwell != nil {
		// quick check to see if we have room
		vol_left := get_vol_left(curwell)

		if vol_left >= vol && curwell.Contains(component) {
			// fine we can just return this one
			return curwell, true
		}

		startcoords := curwell.Crds
		it.MoveTo(startcoords)
	}

	var new_well *LHWell

	for wc := it.Curr(); it.Valid(); wc = it.Next() {

		new_well = plate.GetChildByAddress(wc).(*LHWell)

		if new_well.IsEmpty() {
			break
		}
		/*
			cnts := new_well.Contents()

			cont := cnts.Name()
			// oops... need to check if this is an instance or not
			if cont != component.Name() {
				continue
			}
		*/
		if !new_well.Contains(component) {
			new_well = nil
			continue
		}
		vol_left := get_vol_left(new_well)

		if vol < vol_left {
			break
		} else {
			new_well = nil
		}
	}

	if new_well == nil {
		return nil, false
	}

	return new_well, true
}

//XXX sloboda? This makes no sense now; need to revise
func get_vol_left(well *LHWell) float64 {
	//cnts := well.WContents
	// this is very odd... I can see how this works as a heuristic
	// but it doesn't make much sense to me
	carry_vol := 10.0 // microlitres
	//	total_carry_vol := float64(len(cnts)) * carry_vol
	total_carry_vol := carry_vol // yeah right
	Currvol := well.CurrentVolume().ConvertToString("ul")
	rvol := well.ResidualVolume().ConvertToString("ul")
	vol := well.MaxVolume().ConvertToString("ul")
	return vol - (Currvol + total_carry_vol + rvol)
}

func (well *LHWell) DeclareAutoallocated() {
	if well != nil {

		if well.Extra == nil {
			well.Extra = make(map[string]interface{})
		}

		well.Extra["autoallocated"] = true
	} else {
		fmt.Println("Warning: Attempt to access nil well in DeclareAutoallocated()")
	}
}

func (well *LHWell) DeclareNotAutoallocated() {
	if well != nil {
		if well.Extra == nil {
			well.Extra = make(map[string]interface{})
		}
		well.Extra["autoallocated"] = false
	} else {
		fmt.Println("Warning: Attempt to access nil well in DeclareNotAutoallocated()")
	}
}

func (well *LHWell) IsAutoallocated() bool {
	if well != nil {
		if well.Extra == nil {
			return false
		}

		t, ok := well.Extra["autoallocated"]

		if !ok || !t.(bool) {
			return false
		}
		return true
	} else {
		fmt.Println("Warning: Attempt to access nil well in IsAutoallocated()")
	}
	return false
}

func (w *LHWell) ResetPlateID(newID string) {
	if w == nil {
		return
	}
	ltx := strings.Split(w.WContents.Loc, ":")
	w.WContents.Loc = newID + ":" + ltx[1]
	//w.Plateid = newID
}

func (w *LHWell) XDim() float64 {
	return w.Bounds.GetSize().X
}

func (w *LHWell) YDim() float64 {
	return w.Bounds.GetSize().Y
}
func (w *LHWell) ZDim() float64 {
	return w.Bounds.GetSize().Z
}

func (w *LHWell) IsUserAllocated() bool {
	if w == nil {
		return false
	}
	if w.Extra == nil {
		return false
	}

	ua, ok := w.Extra["UserAllocated"].(bool)

	if !ok {
		return false
	}

	return ua
}

func (w *LHWell) SetUserAllocated() {
	if w == nil {
		return
	}
	if w.Extra == nil {
		w.Extra = make(map[string]interface{})
	}
	w.Extra["UserAllocated"] = true
}

func (w *LHWell) ClearUserAllocated() {
	if w == nil {
		return
	}
	if w.Extra == nil {
		w.Extra = make(map[string]interface{})
	}
	w.Extra["UserAllocated"] = false
}

func (w *LHWell) Contains(cmp *Liquid) bool {
	// obviously empty wells don't contain anything
	if w.IsEmpty() || cmp == nil {
		return false
	}
	// components are the keepers of this information
	return cmp.Matches(w.WContents)
}

func (w *LHWell) UpdateContentID(IDBefore string, after *Liquid) bool {
	if w.WContents.ID == IDBefore {
		/*
			previous := w.WContents
			after.AddParentComponent(previous)
			after.Loc = w.WContents.Loc

			fmt.Println("UPDATE BEFORE: ", w.WContents.CName, " ", w.WContents.Vol, " ", after.CName, " ", after.Vol)

			w.WContents = after
		*/

		w.WContents.AddParentComponent(w.WContents)
		w.WContents.ID = after.ID
		w.WContents.CName = after.CName
		return true
	}

	return false
}

// CheckExtraKey checks if the key is a reserved name
func (w LHWell) CheckExtraKey(s string) error {
	reserved := []string{"afvfunc", "temporary", "autoallocated", "UserAllocated"}

	if wutil.StrInStrArray(s, reserved) {
		return fmt.Errorf("%s is a system key used by plates", s)
	}

	return nil
}

const wellTargetKey = "well_targets"

func (w *LHWell) getTargetMap() map[string][]Coordinates3D {
	m, ok := w.Extra[wellTargetKey]
	if !ok {
		ret := make(map[string][]Coordinates3D)
		w.Extra[wellTargetKey] = ret
		return ret
	}
	ret, ok := m.(map[string][]Coordinates3D)
	//gets broken by serialise/deserialise
	if !ok {
		ret := make(map[string][]Coordinates3D)
		for name, coords := range m.(map[string]interface{}) {
			ret[name] = make([]Coordinates3D, 0)
			for _, crd := range coords.([]interface{}) {
				c := crd.(map[string]interface{})
				ret[name] = append(ret[name], Coordinates3D{X: c["X"].(float64),
					Y: c["Y"].(float64),
					Z: c["Z"].(float64)})
			}
		}
		return ret
	}
	return ret
}

//SetWellTargets sets the targets for each channel when accessing the well using
//adaptor. offsets is a list of coordinates which specify the offset in mm from
//the well center for each channel.
//len(offsets) specifies the maximum number of channels which can access the well
//simultaneously using adaptor
func (w *LHWell) SetWellTargets(adaptor string, offsets []Coordinates3D) {
	wtMap := w.getTargetMap()
	wtMap[adaptor] = offsets
}

//GetWellTargets return the well targets for the given adaptor
//if the adaptor has no defined targets, simply returns the well center
func (w *LHWell) GetWellTargets(adaptor string) []Coordinates3D {
	wtMap := w.getTargetMap()
	adaptorTargets, ok := wtMap[adaptor]
	if !ok {
		return []Coordinates3D{}
	}
	return adaptorTargets
}

//GetAdaptorsWithTargets gets the list of the names of adaptors which have
//targets defined for the well
func (w *LHWell) ListAdaptorsWithTargets() []string {
	wtMap := w.getTargetMap()
	ret := make([]string, 0, len(wtMap))
	for adaptorName := range wtMap {
		ret = append(ret, adaptorName)
	}
	return ret
}
