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

// serializable, stripped-down version of the LHPlate
type SLHPlate struct {
	ID             string
	Inst           string
	Loc            string
	Name           string
	Type           string
	Mnfr           string
	WellsX         int
	WellsY         int
	Nwells         int
	Width          float64
	Length         float64
	Height         float64
	Hunit          string
	Welltype       *LHWell
	Wellcoords     map[string]*LHWell
	Welldimensions *LHWellType
}

func (slhp SLHPlate) FillPlate(plate *LHPlate) {
	plate.ID = slhp.ID
	plate.Inst = slhp.Inst
	plate.Loc = slhp.Loc
	plate.PlateName = slhp.Name
	plate.Type = slhp.Type
	plate.Mnfr = slhp.Mnfr
	plate.WlsX = slhp.WellsX
	plate.WlsY = slhp.WellsY
	plate.Nwells = slhp.Nwells
	//	plate.Width = slhp.Width
	//	plate.Length = slhp.Length
	//	plate.Height = slhp.Height
	//  plate.Hunit = slhp.Hunit
	plate.Welltype = slhp.Welltype
	plate.Wellcoords = slhp.Wellcoords
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
	w.Dunit = lhwt.Dunit
}

func (plate *LHPlate) Welldimensions() *LHWellType {
	t := plate.Welltype
	lhwt := LHWellType{t.MaxVol, t.Vunit, t.Rvol, NameOf(t.WShape), t.Bottom, t.GetSize().X, t.GetSize().Y, t.GetSize().Z, t.Bottomh, t.Dunit}
	return &lhwt
}

type SLHWell struct {
	ID       string
	Inst     string
	Coords   WellCoords
	Contents *LHComponent
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
