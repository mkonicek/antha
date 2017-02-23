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
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/antha/anthalib/wutil"
	"github.com/antha-lang/antha/microArch/logger"
)

// structure describing a microplate
type LHPlate struct {
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
	Height      float64            // plate height (unit below)
	Hunit       string             // unit of plate height
	Rows        [][]*LHWell
	Cols        [][]*LHWell
	Welltype    *LHWell
	Wellcoords  map[string]*LHWell // map of coords in A1 format to wells
	WellXOffset float64            // distance (mm) between well centres in X direction
	WellYOffset float64            // distance (mm) between well centres in Y direction
	WellXStart  float64            // offset (mm) to first well in X direction
	WellYStart  float64            // offset (mm) to first well in Y direction
	WellZStart  float64            // offset (mm) to bottom of well in Z direction
}

func (lhp LHPlate) Name() string {
	return lhp.PlateName
}

func (lhp LHPlate) String() string {
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
	Height      : %f,
	Hunit       : %s,
	Rows        : %p,
	Cols        : %p,
	Welltype    : %s,
	Wellcoords  : %p,
	WellXOffset : %f,
	WellYOffset : %f,
	WellXStart  : %f,
	WellYStart  : %f,
	WellZStart  : %f,
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
		lhp.Height,
		lhp.Hunit,
		lhp.Rows,
		lhp.Cols,
		lhp.Welltype.String(),
		lhp.Wellcoords,
		lhp.WellXOffset,
		lhp.WellYOffset,
		lhp.WellXStart,
		lhp.WellYStart,
		lhp.WellZStart,
	)
}

func (lhp *LHPlate) GetContentVector(wv []WellCoords) ComponentVector {
	ret := make([]*LHComponent, len(wv))

	for i, wc := range wv {
		ret[i] = lhp.Wellcoords[wc.FormatA1()].WContents.Dup()
		wv := lhp.Wellcoords[wc.FormatA1()].WorkingVolume()
		ret[i].Vol = wv.ConvertToString(ret[i].Vunit)
	}

	return ret
}

//plateIDs, wellCoords, vols, err = p.FindComponentsMulti(cmps, ori, multi, contiguous)

func (lhp *LHPlate) FindComponentsMulti(cmps ComponentVector, ori, multi int, contiguous bool) (plateIDs, wellCoords []string, vols []wunit.Volume, err error) {

	for _, c := range cmps {
		if contiguous && c == nil {
			err = fmt.Errorf("Cannot do non-contiguous asks")
			return
		}
	}

	err = fmt.Errorf("Not found")

	var it VectorPlateIterator

	if ori == LHVChannel {
		it = NewColVectorIterator(lhp, multi)
	} else {
		it = NewRowVectorIterator(lhp, multi)
	}

	for wv := it.Curr(); it.Valid(); wv = it.Next() {
		mycmps := lhp.GetContentVector(wv)
		if canGet(cmps, mycmps) {
			plateIDs = mycmps.GetPlateIds()
			wellCoords = mycmps.GetWellCoords()
			vols = cmps.GetVols()
			err = nil
		}
	}

	return
}

// this gets ONE component... possibly from several wells
func (lhp *LHPlate) BetterGetComponent(cmp *LHComponent, exact bool, mpv wunit.Volume) ([]WellCoords, []wunit.Volume, bool) {
	// we first try to find a single well that satisfies us
	// should do DP to improve on this mess
	ret := make([]WellCoords, 0, 1)
	vols := make([]wunit.Volume, 0, 1)
	it := NewOneTimeColumnWiseIterator(lhp)

	volGot := wunit.NewVolume(0.0, "ul")
	volWant := cmp.Volume().Dup()

	// find any well with at least as much as we need
	// if exists, return, if not then fall through

	for wc := it.Curr(); it.Valid(); wc = it.Next() {
		w := lhp.Wellcoords[wc.FormatA1()]

		if w.Empty() {
			continue
		}

		if w.Contents().CName == cmp.CName {
			if exact && w.Contents().ID != cmp.ID {
				continue
			}

			v := w.WorkingVolume()

			if v.LessThan(volWant) {
				continue
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
		return lhp.GetComponent(cmp, exact, mpv)
	}
	//fmt.Println("FOUND: ", cmp.CName, " AT: ", ret[0].FormatA1(), " WANT ", cmp.Volume().ToString(), " GOT ", volGot.ToString(), "  ", ret)

	return ret, vols, true
}

// convenience method

func (lhp *LHPlate) AddComponent(cmp *LHComponent, overflow bool) (wc []WellCoords, err error) {
	ret := make([]WellCoords, 0, 1)

	v := wunit.NewVolume(cmp.Vol, cmp.Vunit)
	wv := wunit.NewVolume(lhp.Welltype.MaxVol, lhp.Welltype.Vunit)

	if v.GreaterThan(wv) && !overflow {
		return ret, fmt.Errorf("Too much to put in a single well of this type")
	}

	it := NewOneTimeColumnWiseIterator(lhp)

	vt := wunit.ZeroVolume()

	for wc := it.Curr(); it.Valid(); wc = it.Next() {
		wl := lhp.Wellcoords[wc.FormatA1()]

		if !wl.Empty() {
			continue
		}

		c, e := cmp.Sample(wv)

		if e != nil {
			return ret, e
		}

		ret = append(ret, wc)
		wl.Add(c)
		vt.Add(c.Volume())
		if vt.EqualTo(v) {
			return ret, nil
		}
	}

	return ret, fmt.Errorf("Not enough empty wells")
}

// convenience method

func (lhp *LHPlate) GetComponent(cmp *LHComponent, exact bool, mpv wunit.Volume) ([]WellCoords, []wunit.Volume, bool) {
	ret := make([]WellCoords, 0, 1)
	vols := make([]wunit.Volume, 0, 1)
	it := NewOneTimeColumnWiseIterator(lhp)

	volGot := wunit.NewVolume(0.0, "ul")
	volWant := cmp.Volume().Dup()

	for wc := it.Curr(); it.Valid(); wc = it.Next() {
		w := lhp.Wellcoords[wc.FormatA1()]

		/*
			if !w.Empty() {
				logger.Debug(fmt.Sprint("WANT: ", cmp.CName, " :: ", wc.FormatA1(), " ", w.Contents().CName, " ", w.CurrVolume().ToString()))
			}
		*/
		if w.Contents().CName == cmp.CName {
			if exact && w.Contents().ID != cmp.ID {
				continue
			}

			v := w.WorkingVolume()
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

func (lhp *LHPlate) Wells() [][]*LHWell {
	return lhp.Rows
}
func (lhp *LHPlate) WellMap() map[string]*LHWell {
	return lhp.Wellcoords
}

const (
	BYROW    = true
	BYCOLUMN = false
)

func (lhp *LHPlate) AllWellPositions(byrow bool) (wellpositionarray []string) {

	wellpositionarray = make([]string, 0)

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

// @implement named

func (lhp *LHPlate) GetName() string {
	return lhp.PlateName
}

func (lhp *LHPlate) WellAt(wc WellCoords) *LHWell {
	return lhp.Wellcoords[wc.FormatA1()]
}

func (lhp *LHPlate) WellAtString(s string) (*LHWell, bool) {
	// improve later, start by assuming these are in FormatA1()
	w, ok := lhp.Wellcoords[s]

	return w, ok
}

func (lhp *LHPlate) WellsX() int {
	return lhp.WlsX
}

func (lhp *LHPlate) WellsY() int {
	return lhp.WlsY
}

func (lhp *LHPlate) NextEmptyWell(it PlateIterator) WellCoords {
	c := 0
	for wc := it.Curr(); it.Valid(); wc = it.Next() {
		if c == lhp.Nwells {
			// prevent iterators from ever making this loop infinitely
			break
		}

		if lhp.Cols[wc.X][wc.Y].Empty() {
			return wc
		}
	}

	return ZeroWellCoords()
}

func NewLHPlate(platetype, mfr string, nrows, ncols int, height float64, hunit string, welltype *LHWell, wellXOffset, wellYOffset, wellXStart, wellYStart, wellZStart float64) *LHPlate {
	var lhp LHPlate
	lhp.Type = platetype
	//lhp.ID = "plate-" + GetUUID()
	lhp.ID = GetUUID()
	lhp.Mnfr = mfr
	lhp.WlsX = ncols
	lhp.WlsY = nrows
	lhp.Nwells = ncols * nrows
	lhp.Height = height
	lhp.Hunit = hunit
	lhp.Welltype = welltype
	lhp.WellXOffset = wellXOffset
	lhp.WellYOffset = wellYOffset
	lhp.WellXStart = wellXStart
	lhp.WellYStart = wellYStart
	lhp.WellZStart = wellZStart

	wellcoords := make(map[string]*LHWell, ncols*nrows)

	// make wells
	rowarr := make([][]*LHWell, nrows)
	colarr := make([][]*LHWell, ncols)
	arr := make([][]*LHWell, nrows)
	wellmap := make(map[string]*LHWell, ncols*nrows)

	for i := 0; i < nrows; i++ {
		arr[i] = make([]*LHWell, ncols)
		rowarr[i] = make([]*LHWell, ncols)
		for j := 0; j < ncols; j++ {
			if colarr[j] == nil {
				colarr[j] = make([]*LHWell, nrows)
			}
			arr[i][j] = welltype.CDup()

			//crds := wutil.NumToAlpha(i+1) + ":" + strconv.Itoa(j+1)
			crds := WellCoords{j, i}.FormatA1()
			wellcoords[crds] = arr[i][j]
			arr[i][j].Crds = crds
			colarr[j][i] = arr[i][j]
			rowarr[i][j] = arr[i][j]
			wellmap[arr[i][j].ID] = arr[i][j]
			arr[i][j].Plate = &lhp
			arr[i][j].Plateinst = lhp.Inst
			arr[i][j].Plateid = lhp.ID
			arr[i][j].Platetype = lhp.Type
			arr[i][j].Crds = crds
			arr[i][j].WContents.Loc = lhp.ID + ":" + crds
		}
	}

	lhp.Wellcoords = wellcoords
	lhp.HWells = wellmap
	lhp.Cols = colarr
	lhp.Rows = rowarr

	return &lhp
}

func (lhp *LHPlate) Dup() *LHPlate {
	// protect yourself fgs
	if lhp == nil {
		logger.Fatal(fmt.Sprintln("Can't dup nonexistent plate"))
	}
	ret := NewLHPlate(lhp.Type, lhp.Mnfr, lhp.WlsY, lhp.WlsX, lhp.Height, lhp.Hunit, lhp.Welltype, lhp.WellXOffset, lhp.WellYOffset, lhp.WellXStart, lhp.WellYStart, lhp.WellZStart)

	ret.PlateName = lhp.PlateName

	ret.HWells = make(map[string]*LHWell, len(ret.HWells))

	for i, row := range lhp.Rows {
		for j, well := range row {
			d := well.Dup()
			ret.Rows[i][j] = d
			ret.Cols[j][i] = d
			ret.Wellcoords[d.Crds] = d
			ret.HWells[d.ID] = d
			d.WContents.Loc = ret.ID + ":" + d.Crds
			d.Plate = ret
			d.Plateinst = ret.Inst
			d.Plateid = ret.ID
		}
	}

	return ret
}
func (lhp *LHPlate) DupKeepIDs() *LHPlate {
	// protect yourself fgs
	if lhp == nil {
		logger.Fatal(fmt.Sprintln("Can't dup nonexistent plate"))
	}
	ret := NewLHPlate(lhp.Type, lhp.Mnfr, lhp.WlsY, lhp.WlsX, lhp.Height, lhp.Hunit, lhp.Welltype, lhp.WellXOffset, lhp.WellYOffset, lhp.WellXStart, lhp.WellYStart, lhp.WellZStart)
	ret.ID = lhp.ID

	ret.PlateName = lhp.PlateName

	ret.HWells = make(map[string]*LHWell, len(ret.HWells))

	for i, row := range lhp.Rows {
		for j, well := range row {
			d := well.Dup()
			d.ID = well.ID
			ret.Rows[i][j] = d
			ret.Cols[j][i] = d
			ret.Wellcoords[d.Crds] = d
			ret.HWells[d.ID] = d
			d.WContents.ID = well.WContents.ID
			d.WContents.Loc = ret.ID + ":" + d.Crds
			d.Plate = ret
			d.Plateinst = ret.Inst
			d.Plateid = ret.ID
		}
	}

	return ret
}

func (p *LHPlate) ProtectAllWells() {
	for _, v := range p.Wellcoords {
		v.Protect()
	}
}

func (p *LHPlate) UnProtectAllWells() {
	for _, v := range p.Wellcoords {
		v.UnProtect()
	}
}

func New_Plate(platetype *LHPlate) *LHPlate {
	new_plate := NewLHPlate(platetype.Type, platetype.Mnfr, platetype.WlsY, platetype.WlsX, platetype.Height, platetype.Hunit, platetype.Welltype, platetype.WellXOffset, platetype.WellYOffset, platetype.WellXStart, platetype.WellYStart, platetype.WellZStart)
	//	Initialize_Wells(new_plate)
	return new_plate
}

func Initialize_Wells(plate *LHPlate) {
	id := (*plate).ID
	wells := (*plate).HWells
	newwells := make(map[string]*LHWell, len(wells))
	wellcrds := (*plate).Wellcoords
	for _, well := range wells {
		well.ID = GetUUID()
		well.Plateid = id
		newwells[well.ID] = well
		wellcrds[well.Crds] = well
	}
	(*plate).HWells = newwells
	(*plate).Wellcoords = wellcrds
}

func (p *LHPlate) RemoveComponent(well string, vol wunit.Volume) *LHComponent {
	w := p.Wellcoords[well]

	if w == nil {
		logger.Debug(fmt.Sprint("RemoveComponent (plate) ERROR: ", well, " ", vol.ToString(), " Can't find well"))
		return nil
	}

	err := w.Remove(vol)

	return err
}

func (p *LHPlate) DeclareTemporary() {
	for _, w := range p.Wellcoords {
		w.DeclareTemporary()
	}
}

func (p *LHPlate) IsTemporary() bool {
	for _, w := range p.Wellcoords {
		if !w.IsTemporary() {
			return false
		}
	}

	return true
}

func (p *LHPlate) DeclareAutoallocated() {
	for _, w := range p.Wellcoords {
		w.DeclareAutoallocated()
	}
}

func (p *LHPlate) IsAutoallocated() bool {
	for _, w := range p.Wellcoords {
		if !w.IsAutoallocated() {
			return false
		}
	}

	return true
}

func ExportPlateCSV(outputpilename string, plate *LHPlate, platename string, wells []string, liquids []*LHComponent, Volumes []wunit.Volume) error {

	csvfile, err := os.Create(outputpilename)
	if err != nil {
		return err
	}

	defer csvfile.Close()

	records := make([][]string, 0)

	//record := make([]string, 0)

	headerrecord := []string{plate.Type, platename, "", "", "", "", ""}

	records = append(records, headerrecord)

	for i, well := range wells {

		volfloat := Volumes[i].RawValue()

		volstr := strconv.FormatFloat(volfloat, 'G', -1, 64)

		// if no conc unit and conc is zero use a default concentration unit
		if liquids[i].Conc == 0 && liquids[i].Cunit == "" {
			liquids[i].Cunit = "mg/l"
		}

		record := []string{well, liquids[i].CName, liquids[i].TypeName(), volstr, Volumes[i].Unit().PrefixedSymbol(), fmt.Sprint(liquids[i].Conc), liquids[i].Cunit}
		records = append(records, record)
	}

	csvwriter := csv.NewWriter(csvfile)

	for _, record := range records {

		err = csvwriter.Write(record)

		if err != nil {
			return err
		}
	}
	csvwriter.Flush()

	return err
}

func AutoExportPlateCSV(outputfilename string, plate *LHPlate) error {

	var platename string = plate.PlateName
	var wells = make([]string, 0)
	var liquids = make([]*LHComponent, 0)
	var volumes = make([]wunit.Volume, 0)
	var concs = make([]wunit.Concentration, 0)
	allpositions := plate.AllWellPositions(false)

	for _, position := range allpositions {
		well := plate.WellMap()[position]

		if !well.Empty() {
			wells = append(wells, position)
			liquids = append(liquids, well.Contents())
			volumes = append(volumes, well.CurrentVolume())
			if well.Contents().Cunit != "" {
				concs = append(concs, wunit.NewConcentration(well.Contents().Conc, well.Contents().Cunit))
			} else {
				concs = append(concs, wunit.NewConcentration(well.Contents().Conc, "ng/ul"))
			}
		}
	}

	csvfile, err := os.Create(outputfilename)
	if err != nil {
		return err
	}
	defer csvfile.Close()

	records := make([][]string, 0)

	//record := make([]string, 0)

	headerrecord := []string{plate.Type, platename, "LiquidType ", "Vol", "Vol Unit", "Conc", "Conc Unit"}

	records = append(records, headerrecord)

	for i, well := range wells {

		volfloat := volumes[i].RawValue()
		concfloat := concs[i].RawValue()

		volstr := strconv.FormatFloat(volfloat, 'G', -1, 64)
		concstr := strconv.FormatFloat(concfloat, 'G', -1, 64)
		/*
			fmt.Println("len(wells)", len(wells))
			fmt.Println("len(liquids)", len(liquids))
			fmt.Println("len(Volumes)", len(Volumes))
		*/

		record := []string{well, liquids[i].CName, liquids[i].TypeName(), volstr, volumes[i].Unit().PrefixedSymbol(), concstr, concs[i].Unit().PrefixedSymbol()}
		records = append(records, record)
	}

	csvwriter := csv.NewWriter(csvfile)

	for _, record := range records {

		err = csvwriter.Write(record)

		if err != nil {
			return err
		}
	}
	csvwriter.Flush()

	return err
}

func (p *LHPlate) SetConstrained(platform string, positions []string) {
	p.Welltype.Extra[platform] = positions
}

func (p *LHPlate) IsConstrainedOn(platform string) ([]string, bool) {
	var pos []string

	par, ok := p.Welltype.Extra[platform]

	if ok {
		pos = par.([]string)
		return pos, true
	}

	return pos, false

}

func (p *LHPlate) Evaporate(time time.Duration, env Environment) []VolumeCorrection {
	ret := make([]VolumeCorrection, 0, 10)
	if p == nil {
		return ret
	}
	for _, w := range p.Wellcoords {
		if !w.Empty() {
			vc := w.Evaporate(time, env)
			if vc.Type != "" {
				ret = append(ret, vc)
			}
		}
	}

	return ret
}

func (p *LHPlate) ResetID(newID string) {
	for _, w := range p.Wellcoords {
		w.ResetPlateID(newID)
	}
	p.ID = newID
}

func (p *LHPlate) IsUserAllocated() bool {
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
func (p *LHPlate) MergeWith(p2 *LHPlate) {
	// do nothing if these are not same type

	if p.Type != p2.Type {
		return
	}

	// transfer any non-User-Allocated wells in here

	it := NewOneTimeColumnWiseIterator(p)

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

func (p *LHPlate) MarkNonEmptyWellsUserAllocated() {
	for _, w := range p.Wellcoords {
		if !w.Empty() {
			w.SetUserAllocated()
		}
	}
}

func (p *LHPlate) AllNonEmptyWells() []*LHWell {
	ret := make([]*LHWell, 0, p.Nwells)

	it := NewOneTimeColumnWiseIterator(p)

	for wc := it.Curr(); it.Valid(); wc = it.Next() {
		w := p.Wellcoords[wc.FormatA1()]

		if !w.Empty() {
			ret = append(ret, w)
		}
	}

	return ret
}
