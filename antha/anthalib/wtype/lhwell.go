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
	"time"

	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/eng"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/antha/anthalib/wutil"
	"github.com/antha-lang/antha/microArch/logger"
)

const (
	LHWBFLAT = iota
	LHWBU
	LHWBV
)

func BottomType(well *LHWell) (desc string) {

	if well.Bottom == 0 {
		desc = "flat bottomed"
	}
	if well.Bottom == 1 {
		desc = "U bottomed"
	}
	if well.Bottom == 2 {
		desc = "V bottomed"
	}
	return
}

// structure representing a well on a microplate - description of a destination
type LHWell struct {
	ID        string
	Inst      string
	Plateinst string
	Plateid   string
	Platetype string
	Crds      string
	MaxVol    float64
	Vunit     string
	WContents *LHComponent
	Rvol      float64
	WShape    *Shape
	Bottom    int
	Xdim      float64
	Ydim      float64
	Zdim      float64
	Bottomh   float64
	Dunit     string
	Extra     map[string]interface{}
	Plate     *LHPlate `gotopb:"-" json:"-"`
}

func (w LHWell) String() string {
	return fmt.Sprintf(
		`LHWELL{
ID        : %s,
Inst      : %s,
Plateinst : %s,
Plateid   : %s,
Platetype : %s,
Crds      : %s,
MaxVol    : %g,
Vunit     : %s,
WContents : %v,
Rvol      : %g,
WShape    : %v,
Bottom    : %d,
Xdim      : %g,
Ydim      : %g,
Zdim      : %g,
Bottomh   : %g,
Dunit     : %s,
Extra     : %v,
Plate     : %v,
}`,
		w.ID,
		w.Inst,
		w.Plateinst,
		w.Plateid,
		w.Platetype,
		w.Crds,
		w.MaxVol,
		w.Vunit,
		w.WContents,
		w.Rvol,
		w.WShape,
		w.Bottom,
		w.Xdim,
		w.Ydim,
		w.Zdim,
		w.Bottomh,
		w.Dunit,
		w.Extra,
		w.Plate,
	)
}

func (w *LHWell) Protected() bool {
	if w.Extra == nil {
		return false
	}

	p, ok := w.Extra["protected"]

	if !ok || !(p.(bool)) {
		return false
	}

	return true
}

func (w *LHWell) Protect() {
	if w.Extra == nil {
		w.Extra = make(map[string]interface{}, 3)
	}

	w.Extra["protected"] = true
}

func (w *LHWell) UnProtect() {
	if w.Extra == nil {
		w.Extra = make(map[string]interface{}, 3)
	}
	w.Extra["protected"] = false
}

func (w *LHWell) Contents() *LHComponent {
	// be careful
	if w == nil {
		logger.Debug("CONTENTS OF NIL WELL REQUESTED")
		// XXX XXX XXX see below ... returning nil here makes a lot more sense
		return NewLHComponent()
	}
	// this makes no sense - we should maintain the
	// contract that the contents of a well are fixed,
	// not make a new content
	// XXX XXX XXX
	// --> mark this for imminent cleanup
	//     best replacement: set well contents then return that
	if w.WContents == nil {
		return NewLHComponent()
	}

	return w.WContents
}

func (w *LHWell) Currvol() float64 {
	if w == nil {
		return 0.0
	}
	return w.Contents().Vol
}

func (w *LHWell) CurrVolume() wunit.Volume {
	if w == nil {
		return wunit.ZeroVolume()
	}
	return w.Contents().Volume()
}

func (w *LHWell) MaxVolume() wunit.Volume {
	if w == nil {
		return wunit.ZeroVolume()
	}
	return wunit.NewVolume(w.MaxVol, w.Vunit)
}
func (w *LHWell) Add(c *LHComponent) {
	if w == nil {
		return
	}
	//wasEmpty := w.Empty()
	mv := wunit.NewVolume(w.MaxVol, w.Vunit)
	cv := wunit.NewVolume(c.Vol, c.Vunit)
	wv := w.CurrentVolume()
	cv.Add(wv)
	if cv.GreaterThan(mv) {
		// could make this fatal but we don't track state well enough
		// for that to be worthwhile
		logger.Debug("WARNING: OVERFULL WELL AT ", w.Crds)
	}

	w.Contents().Mix(c)

	//if wasEmpty {
	// get rid of junk ID
	//	logger.Track(fmt.Sprintf("MIX REPLACED WELL CONTENTS ID WAS %s NOW %s", w.WContents.ID, c.ID))
	//w.WContents.ID = c.ID
	//}
}

func (w *LHWell) Remove(v wunit.Volume) *LHComponent {
	if w == nil {
		return nil
	}
	// if the volume is too high we complain

	if v.GreaterThan(w.CurrentVolume()) {
		//logger.Debug("You ask too much: ", w.Crds, " ", v.ToString(), " I only have: ", w.CurrentVolume().ToString(), " PLATEID: ", w.Plateid)
		return nil
	}

	ret := w.Contents().Dup()
	ret.Vol = v.ConvertToString(w.Vunit)

	w.Contents().Remove(v)
	return ret
}

func (w *LHWell) PlateLocation() PlateLocation {
	if w == nil {
		return ZeroPlateLocation()
	}
	return w.WContents.PlateLocation()
}

func (w *LHWell) WorkingVolume() wunit.Volume {
	if w == nil {
		return wunit.ZeroVolume()
	}
	v := wunit.NewVolume(w.Currvol(), w.Vunit)
	v2 := wunit.NewVolume(w.Rvol, w.Vunit)
	v.Subtract(v2)
	return v
}

func (w *LHWell) ResidualVolume() wunit.Volume {
	if w == nil {
		return wunit.ZeroVolume()
	}
	v := wunit.NewVolume(w.Rvol, w.Vunit)
	return v
}

func (w *LHWell) CurrentVolume() wunit.Volume {
	if w == nil {
		return wunit.ZeroVolume()
	}
	return w.Contents().Volume()
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
	return lhw.Platetype
}

func (lhw *LHWell) Shape() *Shape {
	if lhw == nil {
		return NewNilShape()
	}
	if lhw.WShape == nil {
		// return the non-shape
		return NewNilShape()
	}
	return lhw.WShape
}

// @implement Well
// @deprecate Well

func (w *LHWell) ContainerType() string {
	if w == nil {
		return ""
	}
	return w.Platetype
}

func (w *LHWell) Clear() {
	if w == nil {
		return
	}
	w.WContents = NewLHComponent()
	w.WContents.Loc = w.Plateid + ":" + w.Crds
}

func (w *LHWell) Empty() bool {
	// nil wells are empty
	if w == nil {
		return true
	}

	if w.Currvol() <= 0.000001 {
		return true
	} else {
		return false
	}
}

// copy of instance
func (lhw *LHWell) Dup() *LHWell {
	if lhw == nil {
		return nil
	}
	cp := NewLHWell(lhw.Platetype, lhw.Plateid, lhw.Crds, lhw.Vunit, lhw.MaxVol, lhw.Rvol, lhw.Shape().Dup(), lhw.Bottom, lhw.Xdim, lhw.Ydim, lhw.Zdim, lhw.Bottomh, lhw.Dunit)

	for k, v := range lhw.Extra {
		cp.Extra[k] = v
	}

	cp.WContents = lhw.Contents().Dup()

	return cp
}

// copy of type
func (lhw *LHWell) CDup() *LHWell {
	if lhw == nil {
		return nil
	}
	cp := NewLHWell(lhw.Platetype, lhw.Plateid, lhw.Crds, lhw.Vunit, lhw.MaxVol, lhw.Rvol, lhw.Shape().Dup(), lhw.Bottom, lhw.Xdim, lhw.Ydim, lhw.Zdim, lhw.Bottomh, lhw.Dunit)
	for k, v := range lhw.Extra {
		cp.Extra[k] = v
	}

	return cp
}
func (lhw *LHWell) DupKeepIDs() *LHWell {
	if lhw == nil {
		return nil
	}
	cp := NewLHWell(lhw.Platetype, lhw.Plateid, lhw.Crds, lhw.Vunit, lhw.MaxVol, lhw.Rvol, lhw.Shape().Dup(), lhw.Bottom, lhw.Xdim, lhw.Ydim, lhw.Zdim, lhw.Bottomh, lhw.Dunit)

	for k, v := range lhw.Extra {
		cp.Extra[k] = v
	}

	// Dup here doesn't change ID
	cp.WContents = lhw.Contents().Dup()

	cp.ID = lhw.ID

	return cp
}

func (lhw *LHWell) CalculateMaxCrossSectionArea() (ca wunit.Area, err error) {
	if lhw == nil {
		return
	}

	ca, err = lhw.Shape().MaxCrossSectionalArea()

	return
}

func (lhw *LHWell) AreaForVolume() wunit.Area {

	if lhw == nil {
		return wunit.ZeroArea()
	}
	ret := wunit.NewArea(0.0, "m^2")

	vf := lhw.GetAfVFunc()

	if vf == nil {
		ret, _ := lhw.CalculateMaxCrossSectionArea()
		return ret
	} else {
		vol := lhw.WContents.Volume()
		r := vf.F(vol.ConvertToString("ul"))
		ret = wunit.NewArea(r, "mm^2")
	}

	return ret
}

func (lhw *LHWell) HeightForVolume() wunit.Length {
	if lhw == nil {
		return wunit.ZeroLength()
	}
	ret := wunit.NewLength(0.0, "m")

	return ret
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
	} else {
		x, err := wutil.UnmarshalFunc([]byte(f.(string)))
		if err != nil {
			logger.Fatal(fmt.Sprintf("Can't unmarshal function, error: %s", err.Error))
		}
		return x
	}
	return nil
}

func (lhw *LHWell) CalculateMaxVolume() (vol wunit.Volume, err error) {
	if lhw == nil {
		return wunit.ZeroVolume(), fmt.Errorf("Nil well has no max volume")
	}

	if lhw.Bottom == 0 { // flat
		vol, err = lhw.Shape().Volume()
	} /*else if lhw.Bottom == 1 { // round
		vol, err = lhw.Shape().Volume()
		// + additional calculation
	} else if lhw.Bottom == 2 { // Pointed / v-shaped /pyramid
		vol, err = lhw.Shape().Volume()
		// + additional calculation
	}
	*/
	return
}

// make a new well structure
func NewLHWell(platetype, plateid, crds, vunit string, vol, rvol float64, shape *Shape, bott int, xdim, ydim, zdim, bottomh float64, dunit string) *LHWell {
	var well LHWell

	well.WContents = NewLHComponent()
	well.WContents.DeclareInstance()

	//well.ID = "well-" + GetUUID()
	well.ID = GetUUID()
	well.Platetype = platetype
	well.Plateid = plateid
	well.Crds = crds
	well.MaxVol = vol
	well.Rvol = rvol
	well.Vunit = vunit
	well.WShape = shape.Dup()
	well.Bottom = bott
	well.Xdim = xdim
	well.Ydim = ydim
	well.Zdim = zdim
	well.Bottomh = bottomh
	well.Dunit = dunit
	well.Extra = make(map[string]interface{})
	return &well
}

// this function is somewhat buggy... need to define its responsibilities better
func Get_Next_Well(plate *LHPlate, component *LHComponent, curwell *LHWell) (*LHWell, bool) {
	vol := component.Vol

	it := NewOneTimeColumnWiseIterator(plate)

	if curwell != nil {
		// quick check to see if we have room
		vol_left := get_vol_left(curwell)

		if vol_left >= vol && curwell.Contains(component) {
			// fine we can just return this one
			return curwell, true
		}

		startcoords := MakeWellCoords(curwell.Crds)
		it.SetStartTo(startcoords)
		it.Rewind()
		it.Next()
	}

	var new_well *LHWell

	for wc := it.Curr(); it.Valid(); wc = it.Next() {

		crds := wc.FormatA1()

		new_well = plate.Wellcoords[crds]

		if new_well.Empty() {
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
			continue
		}
		vol_left := get_vol_left(new_well)

		if vol < vol_left {
			break
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
	Currvol := well.Currvol
	rvol := well.Rvol
	vol := well.MaxVol
	return vol - (Currvol() + total_carry_vol + rvol)
}

func (well *LHWell) DeclareTemporary() {
	if well != nil {

		if well.Extra == nil {
			well.Extra = make(map[string]interface{})
		}

		well.Extra["temporary"] = true
	} else {
		logger.Debug("Warning: Attempt to access nil well in DeclareTemporary()")
	}
}

func (well *LHWell) DeclareNotTemporary() {
	if well != nil {
		if well.Extra == nil {
			well.Extra = make(map[string]interface{})
		}
		well.Extra["temporary"] = false
	} else {
		logger.Debug("Warning: Attempt to access nil well in DeclareTemporary()")
	}
}

func (well *LHWell) IsTemporary() bool {
	if well != nil {
		if well.Extra == nil {
			return false
		}

		// user allocated wells are never temporary

		if well.IsUserAllocated() {
			return false
		}

		t, ok := well.Extra["temporary"]

		if !ok || !t.(bool) {
			return false
		}
		return true
	} else {
		logger.Debug("Warning: Attempt to access nil well in IsTemporary()")
	}
	return false
}

func (well *LHWell) DeclareAutoallocated() {
	if well != nil {

		if well.Extra == nil {
			well.Extra = make(map[string]interface{})
		}

		well.Extra["autoallocated"] = true
	} else {
		logger.Debug("Warning: Attempt to access nil well in DeclareAutoallocated()")
	}
}

func (well *LHWell) DeclareNotAutoallocated() {
	if well != nil {
		if well.Extra == nil {
			well.Extra = make(map[string]interface{})
		}
		well.Extra["autoallocated"] = false
	} else {
		logger.Debug("Warning: Attempt to access nil well in DeclareNotAutoallocated()")
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
		logger.Debug("Warning: Attempt to access nil well in IsAutoallocated()")
	}
	return false
}

func (well *LHWell) Evaporate(time time.Duration, env Environment) VolumeCorrection {
	var ret VolumeCorrection

	// don't let this happen
	if well == nil {
		return ret
	}

	if well.Empty() {
		return ret
	}

	// we need to use the evaporation calculator
	// we should likely decorate wells since we have different capabilities
	// for different well types

	vol := eng.EvaporationVolume(env.Temperature, "water", env.Humidity, time.Seconds(), env.MeanAirFlowVelocity, well.AreaForVolume(), env.Pressure)

	r := well.Remove(vol)

	if r == nil {
		well.WContents.Vol = 0.0
	}

	ret.Type = "Evaporation"
	ret.Volume = vol.Dup()
	ret.Location = well.WContents.Loc

	return ret
}

func (w *LHWell) ResetPlateID(newID string) {
	if w == nil {
		return
	}
	ltx := strings.Split(w.WContents.Loc, ":")
	w.WContents.Loc = newID + ":" + ltx[1]
	w.Plateid = newID
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

func (w *LHWell) Contains(cmp *LHComponent) bool {
	// obviously empty wells don't contain anything
	if w.Empty() {
		return false
	}
	// request for a specific component
	if cmp.IsInstance() {
		if cmp.IsSample() {
			//  look for the ID of its parent (we don't allow sampling from samples yet)
			return cmp.ParentID == w.WContents.ID
		} else {
			// if this is just the whole component we check for *its* Id
			return cmp.ID == w.WContents.ID
		}
	} else {
		// sufficient to be of same types
		return cmp.IsSameKindAs(w.WContents)
	}
}

func (w *LHWell) UpdateContentID(IDBefore string, after *LHComponent) bool {
	if w.WContents.ID == IDBefore {
		previous := w.WContents
		after.AddParentComponent(previous)
		w.WContents = after
		return true
	}

	return false
}

func (w LHWell) CheckExtraKey(s string) error {
	reserved := []string{"protected", "afvfunc", "temporary", "autoallocated", "UserAllocated"}

	if wutil.StrInStrArray(s, reserved) {
		return fmt.Errorf("%s is a system key used by plates", s)
	}

	return nil
}
