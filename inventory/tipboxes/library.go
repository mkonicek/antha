// Copyright (C) 2017 The Antha authors. All rights reserved.
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

package tipboxes

import (
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/laboratory/effects/id"
)

const (
	XOffset = 14.28
	YOffset = 11.24
	SbsX    = 127.76
	SbsY    = 85.48
)

func getTipboxSize() wtype.Coordinates3D {
	return wtype.Coordinates3D{X: SbsX, Y: SbsY, Z: 60.13}
}

func makeTipboxes(idGen *id.IDGenerator) (tipboxes []*wtype.LHTipbox) {

	shp := wtype.NewShape(wtype.CylinderShape, "mm", 7.3, 7.3, 51.2)
	w := wtype.NewLHWell(idGen, "ul", 250.0, 10.0, shp, 0, 7.3, 7.3, 51.2, 0.0, "mm")
	w.Extra["InnerL"] = 5.6
	w.Extra["InnerW"] = 5.6
	tip := wtype.NewLHTip(idGen, "cybio", "CyBio250", 10.0, 250.0, "ul", false, shp, 0.0)
	tb := wtype.NewLHTipbox(idGen, 8, 12, getTipboxSize(), "CyBio", "CyBio250Tipbox", tip, w, 9.0, 9.0, XOffset, YOffset, 0.0)
	tipboxes = append(tipboxes, tb)

	w = wtype.NewLHWell(idGen, "ul", 50.0, 0.5, shp, 0, 7.3, 7.3, 51.2, 0.0, "mm")
	w.Extra["InnerL"] = 5.6
	w.Extra["InnerW"] = 5.6

	tip = wtype.NewLHTip(idGen, "cybio", "CyBio50", 0.5, 50.0, "ul", false, shp, 0.0)
	tb = wtype.NewLHTipbox(idGen, 8, 12, getTipboxSize(), "CyBio", "CyBio50Tipbox", tip, w, 9.0, 9.0, XOffset, YOffset, 0.0)
	tipboxes = append(tipboxes, tb)

	// these details are incorrect and need fixing
	w = wtype.NewLHWell(idGen, "ul", 1000.0, 50.0, shp, 0, 7.3, 7.3, 51.2, 0.0, "mm")
	w.Extra["InnerL"] = 5.6
	w.Extra["InnerW"] = 5.6
	tip = wtype.NewLHTip(idGen, "cybio", "CyBio1000", 100.0, 1000.0, "ul", false, shp, 0.0)
	tb = wtype.NewLHTipbox(idGen, 8, 12, getTipboxSize(), "CyBio", "CyBio1000Tipbox", tip, w, 9.0, 9.0, XOffset, YOffset, 0.0)
	tipboxes = append(tipboxes, tb)

	tipboxes = append(tipboxes, makeGilsonTipboxes(idGen)...)

	tipboxes = append(tipboxes, makeTecanTipBoxes(idGen)...)

	tipboxes = append(tipboxes, makeHamiltonTipboxes(idGen)...)

	return tipboxes
}

func makeGilsonTipboxes(idGen *id.IDGenerator) []*wtype.LHTipbox {
	var ret []*wtype.LHTipbox

	shp := wtype.NewShape(wtype.CylinderShape, "mm", 7.3, 7.3, 51.2)

	//Non-filter tips

	w := wtype.NewLHWell(idGen, "ul", 200.0, 20.0, shp, 0, 7.3, 7.3, 51.2, 0.0, "mm")
	w.Extra["InnerL"] = 5.6
	w.Extra["InnerW"] = 5.6
	tip := wtype.NewLHTip(idGen, "gilson", "Gilson200", 20.0, 200.0, "ul", false, shp, 44.7)
	tb := wtype.NewLHTipbox(idGen, 8, 12, getTipboxSize(), "Gilson", "D200 Tip Rack (PIPETMAX 8x200)", tip, w, 9.0, 9.0, XOffset, YOffset, 24.78)
	ret = append(ret, tb)

	// this is the low volume version of the high-volume tip.
	w = wtype.NewLHWell(idGen, "ul", 20.0, 1.0, shp, 0, 7.3, 7.3, 51.2, 0.0, "mm")
	w.Extra["InnerL"] = 5.6
	w.Extra["InnerW"] = 5.6

	// the difference between the effective lengths for the two versions is equal to the difference in coordinate system
	// for the two heads: 5.4mm
	tip = wtype.NewLHTip(idGen, "gilson", "LVGilson200", 1.0, 20.0, "ul", false, shp, 39.3)
	tb = wtype.NewLHTipbox(idGen, 8, 12, getTipboxSize(), "Gilson", "D200 Tip Rack (PIPETMAX 8x20)", tip, w, 9.0, 9.0, XOffset, YOffset, 24.78)

	ret = append(ret, tb)
	w = wtype.NewLHWell(idGen, "ul", 20.0, 1.0, shp, 0, 7.3, 7.3, 46.0, 0.0, "mm")
	w.Extra["InnerL"] = 5.5
	w.Extra["InnerW"] = 5.5
	effectiveHeight := 33.9
	tip = wtype.NewLHTip(idGen, "gilson", "Gilson20", 0.5, 20.0, "ul", false, shp, effectiveHeight)
	tb = wtype.NewLHTipbox(idGen, 8, 12, getTipboxSize(), "Gilson", "DL10 Tip Rack (PIPETMAX 8x20)", tip, w, 9.0, 9.0, XOffset, YOffset, 28.93)
	ret = append(ret, tb)

	//Filter tips

	//Tipeffectiveheight values below are consistent with the values supplied by gilson
	//however, physical testing showed that the offset below was required to avoid collision with the bottom of the well
	filterHeightOffset := 0.00

	w = wtype.NewLHWell(idGen, "ul", 200.0, 20.0, shp, 0, 7.3, 7.3, 51.2, 0.0, "mm")
	w.Extra["InnerL"] = 5.6
	w.Extra["InnerW"] = 5.6
	tip = wtype.NewLHTip(idGen, "gilson", "GilsonFilter200", 20.0, 200.0, "ul", true, shp, 44.7+filterHeightOffset)
	tb = wtype.NewLHTipbox(idGen, 8, 12, getTipboxSize(), "Gilson", "DF200 Tip Rack (PIPETMAX 8x200)", tip, w, 9.0, 9.0, XOffset, YOffset, 24.78)
	ret = append(ret, tb)

	// this is the low volume version of the high-volume tip.
	w = wtype.NewLHWell(idGen, "ul", 20.0, 1.0, shp, 0, 7.3, 7.3, 51.2, 0.0, "mm")
	w.Extra["InnerL"] = 5.6
	w.Extra["InnerW"] = 5.6
	tip = wtype.NewLHTip(idGen, "gilson", "LVGilsonFilter200", 1.0, 20.0, "ul", true, shp, 39.3+filterHeightOffset)
	tb = wtype.NewLHTipbox(idGen, 8, 12, getTipboxSize(), "Gilson", "DF200 Tip Rack (PIPETMAX 8x20)", tip, w, 9.0, 9.0, XOffset, YOffset, 24.78)
	ret = append(ret, tb)

	//DF30 tip has 20ul max volume to avoid attempts to pick it up with the high volume head which currently causes a crash
	//HJK (20/3/18)
	w = wtype.NewLHWell(idGen, "ul", 20.0, 2.0, shp, 0, 7.3, 7.3, 51.2, 0.0, "mm")
	w.Extra["InnerL"] = 5.6
	w.Extra["InnerW"] = 5.6
	tip = wtype.NewLHTip(idGen, "gilson", "GilsonFilter30", 2.0, 20.0, "ul", true, shp, 39.3+filterHeightOffset)
	tb = wtype.NewLHTipbox(idGen, 8, 12, getTipboxSize(), "Gilson", "DF30 Tip Rack (PIPETMAX 8x20)", tip, w, 9.0, 9.0, XOffset, YOffset, 24.78)
	ret = append(ret, tb)

	w = wtype.NewLHWell(idGen, "ul", 10.0, 0.5, shp, 0, 7.3, 7.3, 46.0, 0.0, "mm")
	w.Extra["InnerL"] = 5.5
	w.Extra["InnerW"] = 5.5
	tip = wtype.NewLHTip(idGen, "gilson", "GilsonFilter10", 0.5, 10.0, "ul", true, shp, 33.9+filterHeightOffset)
	tb = wtype.NewLHTipbox(idGen, 8, 12, getTipboxSize(), "Gilson", "DFL10 Tip Rack (PIPETMAX 8x20)", tip, w, 9.0, 9.0, XOffset, YOffset, 28.93)
	ret = append(ret, tb)

	return ret
}

func makeTecanTipBoxes(idGen *id.IDGenerator) []*wtype.LHTipbox {
	shp := wtype.NewShape(wtype.CylinderShape, "mm", 7.3, 7.3, 51.2)

	ret := make([]*wtype.LHTipbox, 0, 4)

	w := wtype.NewLHWell(idGen, "ul", 1000.0, 200.0, shp, 0, 7.3, 7.3, 46.0, 0.0, "mm")
	tip := wtype.NewLHTip(idGen, "Tecan", "Tecan1000", 200.0, 1000.0, "ul", false, shp, 0.0)
	tb := wtype.NewLHTipbox(idGen, 8, 12, getTipboxSize(), "Tecan", "DiTi 1000uL LiHa", tip, w, 9.0, 9.0, XOffset, YOffset, 28.93)
	ret = append(ret, tb)

	w = wtype.NewLHWell(idGen, "ul", 200.0, 15.0, shp, 0, 7.3, 7.3, 46.0, 0.0, "mm")
	tip = wtype.NewLHTip(idGen, "Tecan", "Tecan200", 15.0, 200.0, "ul", false, shp, 0.0)
	tb = wtype.NewLHTipbox(idGen, 8, 12, getTipboxSize(), "Tecan", "DiTi 200uL LiHa", tip, w, 9.0, 9.0, XOffset, YOffset, 28.93)
	ret = append(ret, tb)

	w = wtype.NewLHWell(idGen, "ul", 50.0, 3.0, shp, 0, 7.3, 7.3, 46.0, 0.0, "mm")
	tip = wtype.NewLHTip(idGen, "Tecan", "Tecan50", 3.0, 50.0, "ul", false, shp, 0.0)
	tb = wtype.NewLHTipbox(idGen, 8, 12, getTipboxSize(), "Tecan", "DiTi 50uL LiHa", tip, w, 9.0, 9.0, XOffset, YOffset, 28.93)
	ret = append(ret, tb)

	w = wtype.NewLHWell(idGen, "ul", 10.0, 1.0, shp, 0, 7.3, 7.3, 46.0, 0.0, "mm")
	tip = wtype.NewLHTip(idGen, "Tecan", "Tecan10", 1.0, 10.0, "ul", false, shp, 0.0)
	tb = wtype.NewLHTipbox(idGen, 8, 12, getTipboxSize(), "Tecan", "DiTi 10uL LiHa", tip, w, 9.0, 9.0, XOffset, YOffset, 28.93)
	ret = append(ret, tb)

	return ret
}

func makeHamiltonTipboxes(idGen *id.IDGenerator) []*wtype.LHTipbox {
	var ret []*wtype.LHTipbox

	//Filter tips

	// 300ul, see https://www.hamiltoncompany.com/automated-liquid-handling/disposable-tips/300-%CE%BCl-conductive-sterile-filter-tips
	shp := wtype.NewShape(wtype.CylinderShape, "mm", 7.3, 7.3, 59.9)
	w := wtype.NewLHWell(idGen, "ul", 300.0, 20.0, shp, 0, 7.3, 7.3, 59.9, 0.0, "mm")
	tip := wtype.NewLHTip(idGen, "Hamilton", "Hx300F", 20.0, 300.0, "ul", false, shp, 59.9)
	tb := wtype.NewLHTipbox(idGen, 8, 12, getTipboxSize(), "Hamilton", "Hx300F Tipbox", tip, w, 9.0, 9.0, XOffset, YOffset, 59.9)
	ret = append(ret, tb)

	// 50ul, see https://www.hamiltoncompany.com/automated-liquid-handling/disposable-tips/50-%CE%BCl-conductive-sterile-filter-tips
	shp = wtype.NewShape(wtype.CylinderShape, "mm", 7.3, 7.3, 50.4)
	w = wtype.NewLHWell(idGen, "ul", 50.0, 1.0, shp, 0, 7.3, 7.3, 50.4, 0.0, "mm")
	tip = wtype.NewLHTip(idGen, "Hamilton", "Hx50F", 1.0, 50.0, "ul", false, shp, 50.4)
	tb = wtype.NewLHTipbox(idGen, 8, 12, getTipboxSize(), "Hamilton", "Hx50F Tipbox", tip, w, 9.0, 9.0, XOffset, YOffset, 50.4)
	ret = append(ret, tb)

	return ret
}
