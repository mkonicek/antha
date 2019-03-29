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
	"github.com/pkg/errors"
)

// MarshalDeckObject marshals any object which can be placed on a deck into valid JSON
// in a way such that it can be unmarshalled by UnMarshalDeckObject
func MarshalDeckObject(object LHObject) ([]byte, error) {
	if class := ClassOf(object); class == "" {
		// shouldn't happen since all current LHObjects implement Classy
		return nil, errors.Errorf("cannot serialise object of type %T", object)
	} else {
		return json.Marshal(struct {
			Class  string
			Object LHObject
		}{
			Class:  ClassOf(object),
			Object: object,
		})
	}
}

// UnmarshalDeckObject unmarshal an on-deck object serialised with MarshalDeckObject
// retaining the correct underlying type
func UnmarshalDeckObject(data []byte) (LHObject, error) {
	obj := struct {
		Class  string
		Object *json.RawMessage
	}{}
	if err := json.Unmarshal(data, &obj); err != nil {
		return nil, err
	}

	switch obj.Class {
	case "plate":
		var p Plate
		return &p, json.Unmarshal(*obj.Object, &p)
	case "tipwaste":
		var tw LHTipwaste
		return &tw, json.Unmarshal(*obj.Object, &tw)
	case "tipbox":
		var tb LHTipbox
		return &tb, json.Unmarshal(*obj.Object, &tb)
	default:
		return nil, errors.Errorf("cannot unmarshal object with class %q", obj.Class)
	}
}

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

type sTip struct {
	ID              string
	Type            string
	Mnfr            string
	Dirty           bool
	MaxVol          wunit.Volume
	MinVol          wunit.Volume
	Shape           *Shape
	Bounds          BBox
	EffectiveHeight float64
	Contents        *Liquid
	Filtered        bool
}

func NewSTip(tip *LHTip) *sTip {
	return &sTip{
		ID:              tip.ID,
		Type:            tip.Type,
		Mnfr:            tip.Mnfr,
		Dirty:           tip.Dirty,
		MaxVol:          tip.MaxVol,
		MinVol:          tip.MinVol,
		Shape:           tip.Shape,
		Bounds:          tip.Bounds,
		EffectiveHeight: tip.EffectiveHeight,
		Contents:        tip.contents,
		Filtered:        tip.Filtered,
	}
}

func (s *sTip) Fill(t *LHTip) {
	t.ID = s.ID
	t.Type = s.Type
	t.Mnfr = s.Mnfr
	t.Dirty = s.Dirty
	t.MaxVol = s.MaxVol
	t.MinVol = s.MinVol
	t.Shape = s.Shape
	t.Bounds = s.Bounds
	t.EffectiveHeight = s.EffectiveHeight
	t.contents = s.Contents
	t.Filtered = s.Filtered
}

type sTipbox struct {
	ID         string
	Boxname    string
	Type       string
	Mnfr       string
	Nrows      int
	Ncols      int
	Height     float64
	Tiptype    *LHTip
	AsWell     *LHWell
	NTips      int
	Tips       [][]*LHTip
	TipXOffset float64
	TipYOffset float64
	TipXStart  float64
	TipYStart  float64
	TipZStart  float64
	Bounds     BBox
}

func newSTipbox(tb *LHTipbox) *sTipbox {
	return &sTipbox{
		ID:         tb.ID,
		Boxname:    tb.Boxname,
		Type:       tb.Type,
		Mnfr:       tb.Mnfr,
		Nrows:      tb.Nrows,
		Ncols:      tb.Ncols,
		Height:     tb.Height,
		Tiptype:    tb.Tiptype,
		AsWell:     tb.AsWell,
		NTips:      tb.NTips,
		Tips:       tb.Tips,
		TipXOffset: tb.TipXOffset,
		TipYOffset: tb.TipYOffset,
		TipXStart:  tb.TipXStart,
		TipYStart:  tb.TipYStart,
		TipZStart:  tb.TipZStart,
		Bounds:     tb.Bounds,
	}
}

func (stb *sTipbox) Fill(tb *LHTipbox) {
	tb.ID = stb.ID
	tb.Boxname = stb.Boxname
	tb.Type = stb.Type
	tb.Mnfr = stb.Mnfr
	tb.Nrows = stb.Nrows
	tb.Ncols = stb.Ncols
	tb.Height = stb.Height
	tb.Tiptype = stb.Tiptype
	tb.AsWell = stb.AsWell
	tb.NTips = stb.NTips
	tb.Tips = stb.Tips
	tb.TipXOffset = stb.TipXOffset
	tb.TipYOffset = stb.TipYOffset
	tb.TipXStart = stb.TipXStart
	tb.TipYStart = stb.TipYStart
	tb.TipZStart = stb.TipZStart
	tb.Bounds = stb.Bounds

	for _, row := range tb.Tips {
		for _, tip := range row {
			if err := tip.SetParent(tb); err != nil {
				//Tip must accept tipbox as parent, so this should never happen
				panic(err)
			}
		}
	}
}

type sTipwaste struct {
	Name       string
	ID         string
	Type       string
	Mnfr       string
	Capacity   int
	Contents   int
	Height     float64
	WellXStart float64
	WellYStart float64
	WellZStart float64
	AsWell     *LHWell
	Bounds     BBox
}

func newSTipwaste(tw *LHTipwaste) *sTipwaste {
	return &sTipwaste{
		Name:       tw.Name,
		ID:         tw.ID,
		Type:       tw.Type,
		Mnfr:       tw.Mnfr,
		Capacity:   tw.Capacity,
		Contents:   tw.Contents,
		Height:     tw.Height,
		WellXStart: tw.WellXStart,
		WellYStart: tw.WellYStart,
		WellZStart: tw.WellZStart,
		AsWell:     tw.AsWell,
		Bounds:     tw.Bounds,
	}
}

func (stw *sTipwaste) Fill(tw *LHTipwaste) {
	tw.Name = stw.Name
	tw.ID = stw.ID
	tw.Type = stw.Type
	tw.Mnfr = stw.Mnfr
	tw.Capacity = stw.Capacity
	tw.Contents = stw.Contents
	tw.Height = stw.Height
	tw.WellXStart = stw.WellXStart
	tw.WellYStart = stw.WellYStart
	tw.WellZStart = stw.WellZStart
	tw.AsWell = stw.AsWell
	tw.Bounds = stw.Bounds

	if err := tw.AsWell.SetParent(tw); err != nil {
		//well should accept any tipwaste as parent, so this should never happen
		panic(err)
	}
}

type sHeadAssemblyPosition struct {
	Offset    Coordinates3D
	HeadIndex int
}

func newSHeadAssemblyPosition(hap *LHHeadAssemblyPosition, heads map[*LHHead]int) *sHeadAssemblyPosition {
	if headIndex, ok := heads[hap.Head]; !ok {
		// caller error, should not happen
		panic(errors.New("head not in head map"))
	} else {
		return &sHeadAssemblyPosition{
			Offset:    hap.Offset,
			HeadIndex: headIndex,
		}
	}
}

func (shap *sHeadAssemblyPosition) Fill(hap *LHHeadAssemblyPosition, heads []*LHHead) {
	hap.Offset = shap.Offset
	hap.Head = heads[shap.HeadIndex]
}

// SerializableHeadAssembly an easily serialisable representation of LHHEadAssembly referring to heads
// by index in some array
type SerializableHeadAssembly struct {
	Positions      []*sHeadAssemblyPosition
	MotionLimits   *BBox
	VelocityLimits *VelocityRange
}

// NewSerializableHeadAssembly convert to an easily serialisable representation of a head assembly
// heads is a map of heads to list index
func NewSerializableHeadAssembly(ha *LHHeadAssembly, heads map[*LHHead]int) *SerializableHeadAssembly {
	positions := make([]*sHeadAssemblyPosition, 0, len(ha.Positions))
	for _, pos := range ha.Positions {
		positions = append(positions, newSHeadAssemblyPosition(pos, heads))
	}

	return &SerializableHeadAssembly{
		Positions:      positions,
		MotionLimits:   ha.MotionLimits,
		VelocityLimits: ha.VelocityLimits,
	}
}

func (sha *SerializableHeadAssembly) Fill(ha *LHHeadAssembly, heads []*LHHead) {
	positions := make([]*LHHeadAssemblyPosition, 0, len(sha.Positions))
	for _, spos := range sha.Positions {
		pos := LHHeadAssemblyPosition{}
		spos.Fill(&pos, heads)
		positions = append(positions, &pos)
	}

	ha.Positions = positions
	ha.MotionLimits = sha.MotionLimits
	ha.VelocityLimits = sha.VelocityLimits
}

type SerializableHead struct {
	Name         string
	Manufacturer string
	ID           string
	AdaptorIndex int
	Params       *LHChannelParameter
	TipLoading   TipLoadingBehaviour
}

func NewSerializableHead(head *LHHead, adaptors map[*LHAdaptor]int) *SerializableHead {
	if adaptorIndex, ok := adaptors[head.Adaptor]; !ok {
		// unknown adaptor loaded, caller error
		panic(errors.New("unknown adaptor found in head"))
	} else {
		return &SerializableHead{
			Name:         head.Name,
			Manufacturer: head.Manufacturer,
			ID:           head.ID,
			AdaptorIndex: adaptorIndex,
			Params:       head.Params,
			TipLoading:   head.TipLoading,
		}
	}
}

func (sh *SerializableHead) Fill(head *LHHead, adaptors []*LHAdaptor) {
	head.Name = sh.Name
	head.Manufacturer = sh.Manufacturer
	head.ID = sh.ID
	head.Adaptor = adaptors[sh.AdaptorIndex]
	head.Params = sh.Params
	head.TipLoading = sh.TipLoading
}
