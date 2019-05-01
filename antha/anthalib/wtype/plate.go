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
	"bytes"
	"encoding/csv"
	"fmt"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/antha/anthalib/wutil"
)

// LHPlate is an alias for Plate to preserve backwards compatibility.
// Plate is the principle type for Liquid handling Plates in Antha.
// The structure describing a microplate.
type LHPlate = Plate

// Plate is the principle type for Liquid handling Plates in Antha.
// The structure describing a microplate.
type Plate struct {
	ID          string
	Inst        string
	Loc         string             // location of plate
	PlateName   string             // user-definable plate name
	Type        string             // plate type
	Mnfr        string             // manufacturer
	WlsX        int                // wells along long axis
	WlsY        int                // wells along short axis
	Nwells      int                // total number of wells
	HWells      map[string]*LHWell // map of well IDs to well
	Rows        [][]*LHWell
	Cols        [][]*LHWell
	Welltype    *LHWell
	Wellcoords  map[string]*LHWell // map of coords in A1 format to wells
	WellXOffset float64            // distance (mm) between well centres in X direction
	WellYOffset float64            // distance (mm) between well centres in Y direction
	WellXStart  float64            // offset (mm) to first well in X direction
	WellYStart  float64            // offset (mm) to first well in Y direction
	WellZStart  float64            // offset (mm) to bottom of well in Z direction
	Bounds      BBox               // (relative) position of the plate (mm), set by parent
	parent      LHObject
}

var (
	CONSTRAINTMARKER = "constraint-"
)

func (plate Plate) OutputLayout() {
	fmt.Println(plate.GetLayout())
}

func (plate Plate) GetLayout() string {
	s := ""
	for x := 0; x < plate.WellsX(); x += 1 {
		for y := 0; y < plate.WellsY(); y += 1 {
			well := plate.Cols[x][y]
			if well.IsEmpty() {
				continue
			}
			s += fmt.Sprint("\t\t")
			var wc WellCoords
			wc.X = x
			wc.Y = y
			s += fmt.Sprint(wc.FormatA1(), " ")
			//for _, c := range well.WContents {
			//	fmt.Print(well.WContents.CN, " ")
			if well.WContents.IsInstance() {
				s += fmt.Sprint(well.WContents.CNID(), " ")
			} else {
				s += fmt.Sprint(well.WContents.CName, " ")
			}
			//}
			s += fmt.Sprintln(well.Contents().Volume())
			s += fmt.Sprintln()
			s += fmt.Sprintln()
		}
	}
	return s
}

func (lhp *Plate) GetID() string {
	return lhp.ID
}

// Name returns the name of the plate.
func (lhp Plate) Name() string {
	return lhp.PlateName
}

// Set name sets the name of the plate.
func (lhp *Plate) SetName(name string) {
	lhp.PlateName = strings.TrimSpace(name)
}

func (lhp Plate) String() string {
	return fmt.Sprintf(
		`LHPlate {
	ID          : %s, 
	Inst        : %s,
	Loc         : %s,
	PlateName   : %s,
	Type        : %s,
	Mnfr        : %s,
	WlsX        : %d,
	WlsY        : %d,
	Nwells      : %d,
	HWells      : %p,
	Rows        : %p,
	Cols        : %p,
	Welltype    : %s,
	Wellcoords  : %p,
	WellXOffset : %f,
	WellYOffset : %f,
	WellXStart  : %f,
	WellYStart  : %f,
	WellZStart  : %f,
	Size  : %f x %f x %f,
}`,
		lhp.ID,
		lhp.Inst,
		lhp.Loc,
		lhp.PlateName,
		lhp.Type,
		lhp.Mnfr,
		lhp.WlsX,
		lhp.WlsY,
		lhp.Nwells,
		lhp.HWells,
		lhp.Rows,
		lhp.Cols,
		lhp.Welltype.String(),
		lhp.Wellcoords,
		lhp.WellXOffset,
		lhp.WellYOffset,
		lhp.WellXStart,
		lhp.WellYStart,
		lhp.WellZStart,
		lhp.Bounds.GetSize().X,
		lhp.Bounds.GetSize().Y,
		lhp.Bounds.GetSize().Z,
	)
}

// deprecated
/*
func (lhp *LHPlate) FindComponentsMulti(cmps ComponentVector, ori, multi int, independent bool) (plateIDs, wellCoords [][]string, vols [][]wunit.Volume, err error) {

	for _, c := range cmps {
		if independent && c == nil {
			// HERE HERE HERE -->  INDEPENDENT MULTI NEEDS THIS
			err = fmt.Errorf("Cannot do non-contiguous asks")
			return
		}
	}

	err = fmt.Errorf("Not found")

	var it VectorPlateIterator

	if ori == LHVChannel {
		//it = NewColVectorIterator(lhp, multi)

		tpw := multi / lhp.WellsY()
		wpt := lhp.WellsY() / multi

		if tpw == 0 {
			tpw = 1
		}

		if wpt == 0 {
			wpt = 1
		}

		it = NewTickingColVectorIterator(lhp, multi, tpw, wpt)
	} else {
		it = NewRowVectorIterator(lhp, multi)
	}

	best := 0.0
	bestMatch := ComponentMatch{}
	/// MIS --> debug multichannel leads me here
	//          -- for some reason it's not picking up ONE of the transfers..
	//	       clearly an annoying edge here somewhere
	//		well G6 in the M$ protocol
	for wv := it.Curr(); it.Valid(); wv = it.Next() {
		// cmps needs duping here
		mycmps := lhp.GetContentVector(wv)

		fmt.Println("INVOKE")
		match, errr := matchComponents(cmps.Dup(), mycmps, independent)

		if errr != nil {
			err = errr
			return
		}

		// issue here: this only ever keeps one match
		// matchComponents needs to return multiple matches

		sc := scoreMatch(match, independent)

		if sc > best {
			bestMatch = match
			best = sc
		}
	}

	for _, m := range bestMatch.Matches {
		plateIDs = append(plateIDs, m.IDs)
		wellCoords = append(wellCoords, m.WCs)
		vols = append(vols, m.Vols)
	}

	fmt.Println("BEST AMTCH CHRHE: ")
	fmt.Println(plateIDs)
	fmt.Println(wellCoords)
	fmt.Println(vols)
	fmt.Println("---")

	if best <= 0.0 {
		err = fmt.Errorf("Not found")
	} else {
		err = nil
	}

	return
}

*/

// this gets ONE component... possibly from several wells
func (lhp *Plate) BetterGetComponent(cmp *Liquid, mpv wunit.Volume, legacyVolume bool) ([]WellCoords, []wunit.Volume, bool) {
	// we first try to find a single well that satisfies us
	// should do DP to improve on this mess
	ret := make([]WellCoords, 0, 1)
	vols := make([]wunit.Volume, 0, 1)
	it := NewAddressIterator(lhp, ColumnWise, TopToBottom, LeftToRight, false)

	volGot := wunit.NewVolume(0.0, "ul")
	volWant := cmp.Volume().Dup()

	// find any well with at least as much as we need
	// if exists, return, if not then fall through

	for wc := it.Curr(); it.Valid(); wc = it.Next() {
		w := lhp.Wellcoords[wc.FormatA1()]

		if w.IsEmpty() {
			continue
		}

		//if w.Contents().CName == cmp.CName {
		if w.Contains(cmp) {
			v := w.CurrentWorkingVolume()

			// check volume unless this is an instance and we are tolerating this
			if !cmp.IsInstance() || !legacyVolume {
				if v.LessThan(volWant) {
					continue
				}
			}

			volGot.Add(volWant)
			ret = append(ret, wc)
			vols = append(vols, volGot)

			volWant.Subtract(volGot)

			if volGot.GreaterThan(cmp.Volume()) || volGot.EqualTo(cmp.Volume()) {
				break
			}
		}

	}

	if volGot.LessThan(cmp.Volume()) {
		return lhp.GetComponent(cmp, mpv)
	}
	//fmt.Println("FOUND: ", cmp.CName, " AT: ", ret[0].FormatA1(), " WANT ", cmp.Volume().ToString(), " GOT ", volGot.ToString(), "  ", ret)

	return ret, vols, true
}

// convenience method

func (lhp *Plate) AddComponent(cmp *Liquid, overflow bool) (wc []WellCoords, err error) {
	ret := make([]WellCoords, 0, 1)

	v := wunit.NewVolume(cmp.Vol, cmp.Vunit)
	wv := lhp.Welltype.MaxVolume()

	if v.GreaterThan(wv) && !overflow {
		return ret, fmt.Errorf("Too much to put in a single well of this type")
	}

	it := NewAddressIterator(lhp, ColumnWise, TopToBottom, LeftToRight, false)

	vt := wunit.ZeroVolume()

	for wc := it.Curr(); it.Valid(); wc = it.Next() {
		wl := lhp.Wellcoords[wc.FormatA1()]

		if !wl.IsEmpty() {
			continue
		}

		c, err := cmp.Sample(wv)
		if err != nil {
			return ret, err
		}

		err = wl.AddComponent(c)
		if err != nil {
			//this shouldn't happen because the well was empty
			//but we should check for linting
			return ret, err
		}

		ret = append(ret, wc)
		vt.Add(c.Volume())
		if vt.EqualTo(v) {
			return ret, nil
		}
	}

	return ret, fmt.Errorf("Not enough empty wells")
}

// convenience method

func (lhp *Plate) GetComponent(cmp *Liquid, mpv wunit.Volume) ([]WellCoords, []wunit.Volume, bool) {
	ret := make([]WellCoords, 0, 1)
	vols := make([]wunit.Volume, 0, 1)
	it := NewAddressIterator(lhp, ColumnWise, TopToBottom, LeftToRight, false)

	volGot := wunit.NewVolume(0.0, "ul")
	volWant := cmp.Volume().Dup()

	for wc := it.Curr(); it.Valid(); wc = it.Next() {
		w := lhp.Wellcoords[wc.FormatA1()]

		if w.Contains(cmp) {
			v := w.CurrentWorkingVolume()
			if v.LessThan(mpv) {
				continue
			}
			volGot.Add(v)
			ret = append(ret, wc)

			if volWant.GreaterThan(v) {
				vols = append(vols, v)
			} else {
				vols = append(vols, volWant.Dup())
			}

			volWant.Subtract(v)

			if volGot.GreaterThan(cmp.Volume()) || volGot.EqualTo(cmp.Volume()) {
				break
			}
		}
	}

	//fmt.Println("FOUND: ", cmp.CName, " WANT ", cmp.Volume().ToString(), " GOT ", volGot.ToString(), "  ", ret)

	if volGot.LessThan(cmp.Volume()) {
		return ret, vols, false
	}

	return ret, vols, true
}

func (lhp *Plate) ValidateVolumes() error {
	var lastErr error
	var errCoords []string

	for coords, well := range lhp.Wellcoords {
		err := well.ValidateVolume()
		if err != nil {
			lastErr = err
			errCoords = append(errCoords, coords)
		}
	}

	if len(errCoords) == 1 {
		return lastErr
	} else if len(errCoords) > 1 {
		return LHError(LH_ERR_VOL, fmt.Sprintf("invalid volumes found in %d wells in plate %s at well coordinates %s. Well Capacity on plate type %s: %s",
			len(errCoords), lhp.GetName(), strings.Join(errCoords, ", "), lhp.Type, lhp.Welltype.MaxVolume().ToString()))
	}

	return nil

}

func (lhp *Plate) Wells() [][]*LHWell {
	return lhp.Rows
}
func (lhp *Plate) WellMap() map[string]*LHWell {
	return lhp.Wellcoords
}

const (
	BYROW    = true
	BYCOLUMN = false
)

func (lhp *Plate) AllWellPositions(byrow bool) (wellpositionarray []string) {
	wellpositionarray = make([]string, 0, lhp.WlsX*lhp.WlsY)

	if byrow {

		// range through well coordinates
		for j := 0; j < lhp.WlsY; j++ {
			for i := 0; i < lhp.WlsX; i++ {
				wellposition := wutil.NumToAlpha(j+1) + strconv.Itoa(i+1)
				wellpositionarray = append(wellpositionarray, wellposition)
			}
		}

	} else {

		// range through well coordinates
		for j := 0; j < lhp.WlsX; j++ {
			for i := 0; i < lhp.WlsY; i++ {
				wellposition := wutil.NumToAlpha(i+1) + strconv.Itoa(j+1)
				wellpositionarray = append(wellpositionarray, wellposition)
			}
		}

	}
	return
}

func (lhp *Plate) GetWellCoordsFromOrdering(ordinals []int, byrow bool) []WellCoords {
	wc := lhp.GetA1WellCoordsFromOrdering(ordinals, byrow)
	return WCArrayFromStrings(wc)
}

func (lhp *Plate) GetA1WellCoordsFromOrdering(ordinals []int, byrow bool) []string {
	wps := lhp.AllWellPositions(byrow)

	ret := make([]string, 0, len(wps))

	for _, v := range ordinals {
		if v < 0 {
			panic("No negative wells allowed")
		}
		if v > len(wps)-1 {
			fmt.Println("LEN WPS - 1", len(wps)-1, " V: ", v)
			panic("No wells out of bounds allowed")
		}
		ret = append(ret, wps[v])
	}

	return ret
}
func (lhp *Plate) GetOrderingFromWellCoords(wc []WellCoords, byrow bool) []int {
	wa1 := A1ArrayFromWellCoords(wc)
	return lhp.GetOrderingFromA1WellCoords(wa1, byrow)
}

func (lhp *Plate) GetOrderingFromA1WellCoords(wa1 []string, byrow bool) []int {
	wps := lhp.AllWellPositions(byrow)

	ret := make([]int, len(wa1))

	for i, v := range wa1 {
		ret[i] = FirstIndexInStrArray(v, wps)
	}

	return ret
}

// @implement named

func (lhp *Plate) GetName() string {
	if lhp == nil {
		return "<nil>"
	}
	return lhp.PlateName
}

// @implement Typed
func (lhp *Plate) GetType() string {
	if lhp == nil {
		return "<nil>"
	}
	return lhp.Type
}

func (self *Plate) GetClass() string {
	return "plate"
}

func (lhp *Plate) WellAt(wc WellCoords) (*LHWell, bool) {
	w, ok := lhp.Wellcoords[wc.FormatA1()]
	return w, ok
}

func (lhp *Plate) WellAtString(s string) (*LHWell, bool) {

	//parse well coords, guessing the format
	wc := MakeWellCoords(s)
	if wc.IsZero() {
		//couldn't parse format
		return nil, false
	}

	return lhp.WellAt(wc)
}

func (lhp *Plate) WellsX() int {
	return lhp.WlsX
}

func (lhp *Plate) WellsY() int {
	return lhp.WlsY
}

func (lhp *Plate) IsEmpty() bool {
	for _, w := range lhp.Wellcoords {
		if !w.IsEmpty() {
			return false
		}
	}
	return true
}

//Clean empty all the wells of the plate so that IsEmpty returns true
func (lhp *Plate) Clean() {

	for _, w := range lhp.Wellcoords {
		w.Clean()
	}
	lhp.Welltype.Clean()
}

func (lhp *Plate) NextEmptyWell(it AddressIterator) WellCoords {
	c := 0
	for wc := it.Curr(); it.Valid(); wc = it.Next() {
		if c == lhp.Nwells {
			// prevent iterators from ever making this loop infinitely
			break
		}

		if lhp.Cols[wc.X][wc.Y].IsEmpty() {
			return wc
		}
	}

	return ZeroWellCoords()
}

func NewLHPlate(platetype, mfr string, nrows, ncols int, size Coordinates3D, welltype *LHWell, wellXOffset, wellYOffset, wellXStart, wellYStart, wellZStart float64) *Plate {
	var lhp Plate
	lhp.Type = platetype
	//lhp.ID = "plate-" + GetUUID()
	lhp.ID = GetUUID()
	lhp.PlateName = fmt.Sprintf("%s_%s", platetype, lhp.ID[1:len(lhp.ID)-2])
	lhp.Mnfr = mfr
	lhp.WlsX = ncols
	lhp.WlsY = nrows
	lhp.Nwells = ncols * nrows
	welltype.Plate = &lhp
	lhp.Welltype = welltype
	lhp.WellXOffset = wellXOffset
	lhp.WellYOffset = wellYOffset
	lhp.WellXStart = wellXStart
	lhp.WellYStart = wellYStart
	lhp.WellZStart = wellZStart
	lhp.Bounds.SetSize(size)

	wellcoords := make(map[string]*LHWell, ncols*nrows)

	// make wells
	rowarr := make([][]*LHWell, nrows)
	colarr := make([][]*LHWell, ncols)
	arr := make([][]*LHWell, nrows)
	wellmap := make(map[string]*LHWell, ncols*nrows)

	//offsets to corner of wells
	xOff := welltype.GetSize().X * 0.5
	yOff := welltype.GetSize().Y * 0.5

	for i := 0; i < nrows; i++ {
		arr[i] = make([]*LHWell, ncols)
		rowarr[i] = make([]*LHWell, ncols)
		for j := 0; j < ncols; j++ {
			if colarr[j] == nil {
				colarr[j] = make([]*LHWell, nrows)
			}
			arr[i][j] = welltype.CDup()

			//crds := wutil.NumToAlpha(i+1) + ":" + strconv.Itoa(j+1)
			crds := WellCoords{j, i}
			wellcoords[crds.FormatA1()] = arr[i][j]
			colarr[j][i] = arr[i][j]
			rowarr[i][j] = arr[i][j]
			wellmap[arr[i][j].ID] = arr[i][j]
			arr[i][j].Plate = &lhp
			arr[i][j].Crds = crds
			arr[i][j].WContents.Loc = lhp.ID + ":" + crds.FormatA1()
			arr[i][j].SetOffset(Coordinates3D{ //nolint
				wellXStart + float64(j)*wellXOffset - xOff,
				wellYStart + float64(i)*wellYOffset - yOff,
				wellZStart,
			})
		}
	}

	lhp.Wellcoords = wellcoords
	lhp.HWells = wellmap
	lhp.Cols = colarr
	lhp.Rows = rowarr

	return &lhp
}

func (lhp *Plate) Dup() *Plate {
	return lhp.dup(false)
}

func (lhp *Plate) DupKeepIDs() *Plate {
	return lhp.dup(true)
}

func (lhp *Plate) dup(keep_ids bool) *Plate {
	// protect yourself fgs
	if lhp == nil {
		panic(fmt.Sprintln("Can't dup nonexistent plate"))
	}

	var wellType *LHWell
	if keep_ids {
		wellType = lhp.Welltype.DupKeepIDs()
	} else {
		wellType = lhp.Welltype.Dup()
	}

	ret := NewLHPlate(lhp.Type, lhp.Mnfr, lhp.WlsY, lhp.WlsX, lhp.GetSize(), wellType, lhp.WellXOffset, lhp.WellYOffset, lhp.WellXStart, lhp.WellYStart, lhp.WellZStart)
	if keep_ids {
		ret.ID = lhp.ID
	}

	ret.PlateName = lhp.PlateName

	ret.HWells = make(map[string]*LHWell, len(ret.HWells))

	var d *LHWell
	for i, row := range lhp.Rows {
		for j, well := range row {
			if keep_ids {
				d = well.DupKeepIDs()
			} else {
				d = well.Dup()
				d.WContents.Loc = ret.ID + ":" + d.Crds.FormatA1()
			}
			d.SetParent(ret) //nolint - ret is certainly an lhplate
			ret.Rows[i][j] = d
			ret.Cols[j][i] = d
			ret.Wellcoords[d.Crds.FormatA1()] = d
			ret.HWells[d.ID] = d
		}
	}

	return ret
}

func Initialize_Wells(plate *Plate) {
	wells := (*plate).HWells
	newwells := make(map[string]*LHWell, len(wells))
	wellcrds := (*plate).Wellcoords
	for _, well := range wells {
		well.ID = GetUUID()
		newwells[well.ID] = well
		wellcrds[well.Crds.FormatA1()] = well
	}
	(*plate).HWells = newwells
	(*plate).Wellcoords = wellcrds
}

func (p *Plate) RemoveComponent(well string, vol wunit.Volume) *Liquid {
	w := p.Wellcoords[well]

	if w == nil {
		fmt.Println("RemoveComponent (plate) ERROR: ", well, " ", vol.ToString(), " Can't find well")
		return nil
	}

	cmp, _ := w.RemoveVolume(vol)

	return cmp
}

func (p *Plate) DeclareAutoallocated() {
	for _, w := range p.Wellcoords {
		w.DeclareAutoallocated()
	}
}

// AllAutoallocated returns true if the plate only contains liquids which were autoallocated
// i.e. nothing user-defined
func (p *Plate) AllAutoallocated() bool {
	for _, w := range p.Wellcoords {
		if !w.IsAutoallocated() && !w.IsEmpty() {
			return false
		}
	}

	return true
}

// SubComponentsHeader is the header found to be used in the plate CSV files to specify Sub components of liquids.
const SubComponentsHeader = "SubComponents"

// ExportPlateCSV a exports an LHPlate and its contents as a csv file.
// The caller is required to set the well locations and volumes explicitely with this function.
func ExportPlateCSV(outputFileName string, plate *Plate, plateName string, wells []string, liquids []*Liquid, volumes []wunit.Volume) (file File, err error) {
	if len(wells) != len(liquids) || len(liquids) != len(volumes) {
		return File{}, fmt.Errorf("Found %d liquids, %d wells and %d volumes. Cannot ExportPlateCSV unless these are all equal.", len(liquids), len(wells), len(volumes))
	}

	records := make([][]string, 0)
	headerrecord := []string{plate.Type, plateName, "LiquidType", "Vol", "Vol Unit", "Conc", "Conc Unit"}
	var includeSubComponentsHeader bool

	for i, well := range wells {
		volfloat := volumes[i].RawValue()
		volstr := strconv.FormatFloat(volfloat, 'G', -1, 64)

		var componentName string
		var liquidType string
		var conc float64
		var concUnit string
		var subComponents []string
		if liquids[i] != nil {
			componentName = liquids[i].Name()
			liquidType = liquids[i].TypeName()
			conc = liquids[i].Conc
			concUnit = liquids[i].Cunit
			if liquids[i].Conc == 0 && liquids[i].Cunit == "" {
				concUnit = "mg/l"
			}

			if liquids[i].SubComponents.Components != nil {
				includeSubComponentsHeader = true
				for componentName := range liquids[i].SubComponents.Components {
					subComponents = append(subComponents, componentName)
				}
			}

			sort.Strings(subComponents)
		}

		record := []string{well, componentName, liquidType, volstr, volumes[i].Unit().PrefixedSymbol(), fmt.Sprint(conc), concUnit}

		for _, componentName := range subComponents {
			record = append(record, componentName+":", liquids[i].SubComponents.Components[componentName].ToString())
		}

		records = append(records, record)
	}

	if includeSubComponentsHeader {
		headerrecord = append(headerrecord, SubComponentsHeader)
	}
	records = append([][]string{headerrecord}, records...)

	return exportCSV(records, outputFileName)
}

// AutoExportPlateCSV exports an LHPlate and its contents as a csv file.
// This is not 100% safe to use in elements since, currently,
// at the time of running an element, the scheduler  will not have allocated positions
// for the components so, for example, accurate well information cannot currently be obtained with this function.
// If allocating wells manually use the ExportPlateCSV function and explicitely set the sample locations and volumes.
func AutoExportPlateCSV(outputFileName string, plate *Plate) (file File, err error) {

	var platename string = plate.PlateName
	var wells = make([]string, 0)
	var liquids = make([]*Liquid, 0)
	var volumes = make([]wunit.Volume, 0)
	allpositions := plate.AllWellPositions(false)

	var nilComponent *Liquid

	for _, position := range allpositions {
		well := plate.WellMap()[position]

		if !well.IsEmpty() {
			wells = append(wells, position)
			liquids = append(liquids, well.Contents())
			volumes = append(volumes, well.CurrentVolume())
		} else {
			wells = append(wells, position)
			liquids = append(liquids, nilComponent)
			volumes = append(volumes, wunit.NewVolume(0.0, "ul"))
		}
	}
	return ExportPlateCSV(outputFileName, plate, platename, wells, liquids, volumes)
}

// Export a 2D array of string data as a csv file
func exportCSV(records [][]string, filename string) (File, error) {
	var anthafile File
	var buf bytes.Buffer

	/// use the buffer to create a csv writer
	w := csv.NewWriter(&buf)

	// write all records to the buffer
	err := w.WriteAll(records)
	if err != nil {
		return anthafile, err
	}

	if err := w.Error(); err != nil {
		return anthafile, fmt.Errorf("error writing csv: %s", err.Error())
	}

	//This code shows how to create an antha File from this buffer which can be downloaded through the UI:

	anthafile.Name = filename

	err = anthafile.WriteAll(buf.Bytes())
	if err != nil {
		return anthafile, err
	}

	///// to write this to a file on the command line this is what we'd do (or something similar)

	// also create a file on os
	file, _ := os.Create(filename)
	defer file.Close() // nolint

	// this time we'll use the file to create the writer instead of a buffer (anything which fulfils the writer interface can be used here ... checkout golang io.Writer and io.Reader)
	fw := csv.NewWriter(file)

	// same as before ...
	err = fw.WriteAll(records)
	return anthafile, err
}

func makeConstraintKeyFor(platform string) string {
	if isConstraintKey(platform) {
		return platform
	}

	return CONSTRAINTMARKER + platform
}

func unMakeConstraintKey(s string) string {
	if !isConstraintKey(s) {
		return s
	}

	return strings.Replace(s, CONSTRAINTMARKER, "", -1)
}

func isConstraintKey(s string) bool {
	return strings.HasPrefix(s, CONSTRAINTMARKER)
}

func (p *Plate) SetConstrained(platform string, positions []string) {

	cstrKey := makeConstraintKeyFor(platform)
	p.Welltype.Extra[cstrKey] = positions
}

func (p *Plate) IsConstrainedOn(platform string) ([]string, bool) {
	cstrKey := makeConstraintKeyFor(platform)
	par, ok := p.Welltype.Extra[cstrKey]
	if !ok {
		return nil, false
	}

	switch par := par.(type) {

	case []string:
		return par, true

	case []interface{}:
		var pos []string
		for _, v := range par {
			pos = append(pos, v.(string))
		}
		return pos, true

	default:
		panic(fmt.Sprintf("unknown type %T", par))
	}
}

func (p *Plate) GetAllConstraints() map[string][]string {
	ret := make(map[string][]string)
	for k, v := range p.Welltype.Extra {
		if isConstraintKey(k) {
			var pos []string
			pos = append(pos, v.([]string)...)
			ret[unMakeConstraintKey(k)] = pos
		}
	}

	return ret
}

//##############################################
//@implement LHObject
//##############################################

func (self *Plate) GetPosition() Coordinates3D {
	if self.parent != nil {
		return self.parent.GetPosition().Add(self.Bounds.GetPosition())
	}
	return self.Bounds.GetPosition()
}

func (self *Plate) GetSize() Coordinates3D {
	return self.Bounds.GetSize()
}

func (self *Plate) GetWellOffset() Coordinates3D {
	return Coordinates3D{self.WellXStart, self.WellYStart, self.WellZStart}
}

func (self *Plate) GetWellCorner() Coordinates3D {
	size := self.Welltype.GetSize()
	return Coordinates3D{self.WellXStart - 0.5*size.X, self.WellYStart - 0.5*size.Y, self.WellZStart}
}

func (self *Plate) GetWellSize() Coordinates3D {
	return Coordinates3D{
		self.WellXOffset*float64(self.NCols()-1) + self.Welltype.GetSize().X,
		self.WellYOffset*float64(self.NRows()-1) + self.Welltype.GetSize().Y,
		self.Welltype.GetSize().Z,
	}
}

func (self *Plate) GetWellBounds() BBox {
	return BBox{
		self.Bounds.GetPosition().Add(self.GetWellCorner()),
		self.GetWellSize(),
	}
}

func (self *Plate) GetBoxIntersections(box BBox) []LHObject {
	//relative to me
	box.SetPosition(box.GetPosition().Subtract(OriginOf(self)))
	ret := []LHObject{}

	if self.GetWellBounds().IntersectsBox(box) {
		for _, row := range self.Rows {
			for _, well := range row {
				ret = append(ret, well.GetBoxIntersections(box)...)
			}
		}
	}

	if len(ret) == 0 && self.Bounds.IntersectsBox(box) {
		ret = append(ret, self)
	}
	//todo, scan through wells
	return ret
}

func (self *Plate) GetPointIntersections(point Coordinates3D) []LHObject {
	//relative
	point = point.Subtract(OriginOf(self))
	ret := []LHObject{}

	if self.GetWellBounds().IntersectsPoint(point) {
		for _, row := range self.Rows {
			for _, well := range row {
				ret = append(ret, well.GetPointIntersections(point)...)
			}
		}
	}

	if len(ret) == 0 && self.Bounds.IntersectsPoint(point) {
		ret = append(ret, self)
	}
	return ret
}

func (self *Plate) SetOffset(o Coordinates3D) error {
	self.Bounds.SetPosition(o)
	return nil
}

func (self *Plate) SetParent(p LHObject) error {
	self.parent = p
	return nil
}

//@implement LHObject
func (self *Plate) ClearParent() {
	self.parent = nil
}

func (self *Plate) GetParent() LHObject {
	return self.parent
}

//Duplicate copies an LHObject
func (self *Plate) Duplicate(keepIDs bool) LHObject {
	return self.dup(keepIDs)
}

//DimensionsString returns a string description of the position and size of the object and its children.
func (self *Plate) DimensionsString() string {
	wb := self.GetWellBounds()
	ret := make([]string, 0, 1+len(self.Wellcoords))
	ret = append(ret, fmt.Sprintf("Plate %s at %v+%v, with %dx%d wells bounded by %v",
		self.GetName(), self.GetPosition(), self.GetSize(), self.NCols(), self.NRows(), wb))

	for _, adaptorName := range self.ListAdaptorsWithTargets() {
		ret = append(ret, fmt.Sprintf("\tTargets for adaptor \"%s\":", adaptorName))
		for i, t := range self.GetTargets(adaptorName) {
			ret = append(ret, fmt.Sprintf("\t\t%d: %v", i, t))
		}
	}

	for _, wellrow := range self.Rows {
		for _, well := range wellrow {
			ret = append(ret, "\t"+well.DimensionsString())
		}
	}

	return strings.Join(ret, "\n")
}

//##############################################
//@implement Addressable
//##############################################

func (self *Plate) AddressExists(c WellCoords) bool {
	return c.X >= 0 &&
		c.Y >= 0 &&
		c.X < self.WlsX &&
		c.Y < self.WlsY
}

func (lhp *Plate) NCols() int {
	return lhp.WlsX
}

func (lhp *Plate) NRows() int {
	return lhp.WlsY
}

func (self *Plate) GetChildByAddress(c WellCoords) LHObject {
	if !self.AddressExists(c) {
		return nil
	}
	//LHWells aren't LHObjects yet
	return self.Cols[c.X][c.Y]
}

func (self *Plate) CoordsToWellCoords(r Coordinates3D) (WellCoords, Coordinates3D) {
	rel := r.Subtract(self.GetPosition())
	wellSize := self.Welltype.GetSize()
	wc := WellCoords{
		int(math.Floor((rel.X - self.WellXStart + 0.5*wellSize.X) / self.WellXOffset)),
		int(math.Floor((rel.Y - self.WellYStart + 0.5*wellSize.Y) / self.WellYOffset)),
	}
	if wc.X < 0 {
		wc.X = 0
	} else if wc.X >= self.WlsX {
		wc.X = self.WlsX - 1
	}
	if wc.Y < 0 {
		wc.Y = 0
	} else if wc.Y >= self.WlsY {
		wc.Y = self.WlsY - 1
	}

	r2, _ := self.WellCoordsToCoords(wc, TopReference)

	return wc, r.Subtract(r2)
}

func (self *Plate) WellCoordsToCoords(wc WellCoords, r WellReference) (Coordinates3D, bool) {
	if !self.AddressExists(wc) {
		return Coordinates3D{}, false
	}

	child := self.GetChildByAddress(wc).(*LHWell)

	var z float64
	//return the bottom as a lower bound for liquid reference
	if r == BottomReference {
		z = child.GetPosition().Z + child.Bottomh
	} else if r == TopReference {
		z = child.GetPosition().Z + child.GetSize().Z
	} else if r == LiquidReference {
		// we don't know exactly where the liquid level is, so return the middle
		z = child.GetPosition().Z + 0.5*(child.Bottomh+child.GetSize().Z)
	}

	center := child.GetPosition().Add(child.GetSize().Multiply(0.5))

	return Coordinates3D{center.X, center.Y, z}, true
}

func (p *Plate) ResetID(newID string) {
	for _, w := range p.Wellcoords {
		w.ResetPlateID(newID)
	}
	p.ID = newID
}

func (p *Plate) Height() float64 {
	return p.Bounds.GetSize().Z
}

func (p *Plate) IsUserAllocated() bool {
	// true if any wells are user allocated

	for _, w := range p.Wellcoords {
		if w.IsUserAllocated() {
			return true
		}
	}

	return false
}

// semantics are: put stuff from p2 into p unless
// the well in p is declared as user allocated
func (p *Plate) MergeWith(p2 *Plate) {
	// do nothing if these are not same type

	if p.Type != p2.Type {
		return
	}

	// transfer any non-User-Allocated wells in here

	it := NewAddressIterator(p, ColumnWise, TopToBottom, LeftToRight, false)

	for ; it.Valid(); it.Next() {
		wc := it.Curr()

		if !it.Valid() {
			break
		}

		w1 := p.Wellcoords[wc.FormatA1()]
		w2 := p2.Wellcoords[wc.FormatA1()]

		if !w1.IsUserAllocated() {
			w1.WContents = w2.WContents
		}
	}
}

func (p *Plate) MarkNonEmptyWellsUserAllocated() {
	for _, w := range p.Wellcoords {
		if !w.IsEmpty() {
			w.SetUserAllocated()
		}
	}
}

func (p *Plate) AllNonEmptyWells() []*LHWell {
	ret := make([]*LHWell, 0, p.Nwells)

	it := NewAddressIterator(p, ColumnWise, TopToBottom, LeftToRight, false)

	for wc := it.Curr(); it.Valid(); wc = it.Next() {
		w := p.Wellcoords[wc.FormatA1()]

		if !w.IsEmpty() {
			ret = append(ret, w)
		}
	}

	return ret
}

//AreaWellTargetsEnabled should well targets be set with this plate?
//aim is to deprecate IsSpecial
func (p *Plate) AreWellTargetsEnabled(adaptorChannels int, channelSpacing float64) bool {

	if p.NRows() != 1 {
		return false
	}

	if !BoxShape.Equals(p.Welltype.Shape().Type) {
		return false
	}

	spaceNeeded := float64(adaptorChannels-1)*channelSpacing + 2.0

	return p.Welltype.GetSize().Y >= spaceNeeded

}

func (p *Plate) IsSpecial() bool {
	if p == nil || p.Welltype.Extra == nil {
		return false
	}

	s, ok := p.Welltype.Extra["IMSPECIAL"]

	if !ok || !s.(bool) {
		return false
	}

	return true
}

func (p *Plate) DeclareSpecial() {
	if p != nil && p.Welltype.Extra != nil {
		p.Welltype.Extra["IMSPECIAL"] = true
	}
}

// AvailableContents accepts a slice of well coordinates, wv, and returns the Liquids that could be taken from each well in subsequent steps.
// if any well coordinate in wv does not exist, a nil liquid is returned in that position
func (lhp *Plate) AvailableContents(wv []WellCoords) ComponentVector {
	ret := make(ComponentVector, len(wv))

	for i, wc := range wv {
		if well := lhp.Wellcoords[wc.FormatA1()]; well != nil {
			l := well.Contents().Dup()
			// don't include the residual volume, since it can't be accessed
			l.SetVolume(well.CurrentWorkingVolume())
			ret[i] = l
		}
	}

	return ret
}

func (p *Plate) FindAndUpdateID(before string, after *Liquid) bool {
	for _, w := range p.Wellcoords {
		if w.UpdateContentID(before, after) {
			return true
		}
	}
	return false
}

// SetData implements Annotatable
func (p *Plate) SetData(key string, data []byte) error {
	if err := p.checkExtra(fmt.Sprintf("cannot add data %s", key)); err != nil {
		return err
	}

	// nb -- in future disallow already set keys as well?
	if err := p.CheckExtraKey(key); err != nil {
		return fmt.Errorf("invalid key %s: %s", key, err)
	}

	p.Welltype.Extra[key] = data

	return nil

}

// ClearData removes data with the given name
func (p *Plate) ClearData(k string) error {
	err := p.checkExtra(fmt.Sprintf("cannot clear data %s", k))

	if err != nil {
		return err
	}

	delete(p.Welltype.Extra, k)

	return nil
}

func (p *Plate) checkExtra(s string) error {
	if p == nil {
		return fmt.Errorf("nil plate: %s", s)
	}

	if p.Welltype == nil {
		return fmt.Errorf("corrupt plate - missing well type: %s", s)
	}

	if p.Welltype.Extra == nil {
		return fmt.Errorf("corrupt well type - %s", s)
	}

	return nil
}

func (p Plate) GetData(key string) ([]byte, error) {
	if err := p.checkExtra(fmt.Sprintf("cannot get key %s", key)); err != nil {
		return nil, err
	}

	if err := p.CheckExtraKey(key); err != nil {
		return nil, fmt.Errorf("invalid key %s: %s", key, err)
	}

	bs, ok := p.Welltype.Extra[key].([]byte)
	if !ok {
		return nil, fmt.Errorf("key %s not found", key)
	}

	return bs, nil
}

// CheckExtraKey checks if the key is a reserved name
func (p Plate) CheckExtraKey(k string) error {
	reserved := []string{"IMSPECIAL", "Pipetmax"}

	if wutil.StrInStrArray(k, reserved) {
		return fmt.Errorf("%s is a system key used by plates", k)
	}

	if p.Welltype == nil {
		return fmt.Errorf("No valid well")
	}

	return p.Welltype.CheckExtraKey(k)
}

// AllContents returns all the components on the plate
func (p *Plate) AllContents() []*Liquid {
	ret := make([]*Liquid, 0, len(p.Wellcoords))
	for _, c := range p.Cols {
		for _, w := range c {
			ret = append(ret, w.WContents)
		}
	}

	return ret
}

func (p *Plate) ColVol() wunit.Volume {
	if p == nil {
		return wunit.ZeroVolume()
	}

	v := p.Welltype.MaxVolume()

	v.MultiplyBy(float64(p.WlsY))

	return v
}

//GetTargetOffset get the offset for addressing a well with the named adaptor and channel
func (p *Plate) GetTargetOffset(adaptorName string, channel int) Coordinates3D {
	targets := p.Welltype.GetWellTargets(adaptorName)
	if channel < 0 || channel >= len(targets) {
		return Coordinates3D{}
	}
	return targets[channel]
}

//GetTargets return all the defined targets for the named adaptor
func (p *Plate) GetTargets(adaptorName string) []Coordinates3D {
	return p.Welltype.GetWellTargets(adaptorName)
}

//ListAdaptorsWithTargets get a list of the names of all the adaptors with targets set
func (p *Plate) ListAdaptorsWithTargets() []string {
	return p.Welltype.ListAdaptorsWithTargets()
}
