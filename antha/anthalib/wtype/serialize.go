// anthalib//liquidhandling/serialize.go: Part of the Antha language
// Copyright (C) 2015 The Antha authors. All rights reserved.
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
	"github.com/antha-lang/antha/antha/anthalib/wunit"
)

func (lhp *Plate) MarshalJSON() ([]byte, error) {
	slhp := lhp.ToSLHPLate()

	return json.Marshal(slhp)
}

func (lhp *Plate) UnmarshalJSON(b []byte) error {
	var slhp SLHPlate

	err := json.Unmarshal(b, &slhp)

	if err != nil {
		return err
	}

	slhp.FillPlate(lhp)

	return nil
}

// serializable, stripped-down version of the LHPlate
type SLHPlate struct {
	ID          string
	Inst        string
	Loc         string
	Name        string
	Type        string
	Mnfr        string
	WellsX      int
	WellsY      int
	Nwells      int
	Bounds      BBox
	Welltype    *LHWell
	Wellcoords  map[string]*LHWell
	WellXOffset float64 // distance (mm) between well centres in X direction
	WellYOffset float64 // distance (mm) between well centres in Y direction
	WellXStart  float64 // offset (mm) to first well in X direction
	WellYStart  float64 // offset (mm) to first well in Y direction
	WellZStart  float64 // offset (mm) to bottom of well in Z direction
}

func (p *Plate) ToSLHPLate() SLHPlate {
	return SLHPlate{
		ID:          p.ID,
		Inst:        p.Inst,
		Loc:         p.Loc,
		Name:        p.PlateName,
		Type:        p.Type,
		Mnfr:        p.Mnfr,
		WellsX:      p.WlsX,
		WellsY:      p.WlsY,
		Nwells:      p.Nwells,
		Bounds:      p.Bounds,
		Welltype:    p.Welltype,
		Wellcoords:  p.Wellcoords,
		WellXOffset: p.WellXOffset,
		WellYOffset: p.WellYOffset,
		WellXStart:  p.WellXStart,
		WellYStart:  p.WellYStart,
		WellZStart:  p.WellZStart,
	}
}

func (slhp SLHPlate) FillPlate(plate *Plate) {
	plate.ID = slhp.ID
	plate.Inst = slhp.Inst
	plate.Loc = slhp.Loc
	plate.PlateName = slhp.Name
	plate.Type = slhp.Type
	plate.Mnfr = slhp.Mnfr
	plate.WlsX = slhp.WellsX
	plate.WlsY = slhp.WellsY
	plate.Nwells = slhp.Nwells
	plate.Bounds = slhp.Bounds
	//	plate.Width = slhp.Width
	//	plate.Length = slhp.Length
	//	plate.Height = slhp.Height
	//  plate.Hunit = slhp.Hunit
	plate.Welltype = slhp.Welltype
	plate.Wellcoords = slhp.Wellcoords
	plate.WellXOffset = slhp.WellXOffset
	plate.WellYOffset = slhp.WellYOffset
	plate.WellXStart = slhp.WellXStart
	plate.WellYStart = slhp.WellYStart
	plate.WellZStart = slhp.WellZStart
	makeRows(plate)
	makeCols(plate)
	plate.HWells = make(map[string]*LHWell, len(plate.Wellcoords))
	for _, w := range plate.Wellcoords {
		plate.HWells[w.ID] = w
		w.Plate = plate
	}
	plate.Welltype.Plate = plate
}

func makeRows(p *Plate) {
	p.Rows = make([][]*LHWell, p.WlsY)
	for i := 0; i < p.WlsY; i++ {
		p.Rows[i] = make([]*LHWell, p.WlsX)
		for j := 0; j < p.WlsX; j++ {
			wc := WellCoords{X: j, Y: i}
			p.Rows[i][j] = p.Wellcoords[wc.FormatA1()]
		}
	}
}
func makeCols(p *Plate) {
	p.Cols = make([][]*LHWell, p.WlsX)
	for i := 0; i < p.WlsX; i++ {
		p.Cols[i] = make([]*LHWell, p.WlsY)
		for j := 0; j < p.WlsY; j++ {
			wc := WellCoords{X: i, Y: j}
			p.Cols[i][j] = p.Wellcoords[wc.FormatA1()]
		}
	}
}

// this is for keeping track of the well type

type LHWellType struct {
	Vol       float64
	Vunit     string
	Rvol      float64
	ShapeName string
	Bottom    WellBottomType
	Xdim      float64
	Ydim      float64
	Zdim      float64
	Bottomh   float64
	Dunit     string
}

func (w *LHWell) AddDimensions(lhwt *LHWellType) {
	w.MaxVol = wunit.NewVolume(lhwt.Vol, lhwt.Vunit).ConvertToString("ul")
	w.Rvol = wunit.NewVolume(lhwt.Rvol, lhwt.Vunit).ConvertToString("ul")
	w.WShape = NewShape(lhwt.ShapeName, lhwt.Dunit, lhwt.Xdim, lhwt.Ydim, lhwt.Zdim)
	w.Bottom = lhwt.Bottom
	w.Bounds.SetSize(Coordinates{
		wunit.NewLength(lhwt.Xdim, lhwt.Dunit).ConvertToString("mm"),
		wunit.NewLength(lhwt.Ydim, lhwt.Dunit).ConvertToString("mm"),
		wunit.NewLength(lhwt.Zdim, lhwt.Dunit).ConvertToString("mm"),
	})
	w.Bottomh = wunit.NewLength(lhwt.Bottomh, lhwt.Dunit).ConvertToString("mm")
}

func (plate *Plate) Welldimensions() *LHWellType {
	t := plate.Welltype
	lhwt := LHWellType{t.MaxVol, "ul", t.Rvol, t.WShape.ShapeName, t.Bottom, t.GetSize().X, t.GetSize().Y, t.GetSize().Z, t.Bottomh, "mm"}
	return &lhwt
}

type SLHWell struct {
	ID       string
	Inst     string
	Coords   WellCoords
	Contents *Liquid
}

func (slw SLHWell) FillWell(lw *LHWell) {
	lw.ID = slw.ID
	lw.Inst = slw.Inst
	lw.Crds = slw.Coords
	lw.WContents = slw.Contents
}

type FromFactory struct {
	String string
}

func (f *FromFactory) MarshalJSON() ([]byte, error) {
	v, e := json.Marshal(f.String)
	return v, e
}

func (f *FromFactory) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	f.String = s
	return nil
}
