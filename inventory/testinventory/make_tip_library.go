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

package testinventory

import (
	"github.com/antha-lang/antha/antha/anthalib/wtype"
)

func makeTipboxes() (tipboxes []*wtype.LHTipbox) {
	// create a well representation of the tip holder... sometimes needed
	// heh, should have kept LHTipholder!
	shp := wtype.NewShape("cylinder", "mm", 7.3, 7.3, 51.2)
	w := wtype.NewLHWell("CyBio250Tipbox", "", "A1", "ul", 250.0, 10.0, shp, 0, 7.3, 7.3, 51.2, 0.0, "mm")
	w.Extra["InnerL"] = 5.6
	w.Extra["InnerW"] = 5.6
	tip := wtype.NewLHTip("cybio", "CyBio250", 10.0, 250.0, "ul")
	tb := wtype.NewLHTipbox(8, 12, 60.13, "CyBio", "CyBio250Tipbox", tip, w, 9.0, 9.0, 0.0, 0.0, 0.0)
	tipboxes = append(tipboxes, tb)

	w = wtype.NewLHWell("CyBio50Tipbox", "", "A1", "ul", 50.0, 0.5, shp, 0, 7.3, 7.3, 51.2, 0.0, "mm")
	w.Extra["InnerL"] = 5.6
	w.Extra["InnerW"] = 5.6

	tip = wtype.NewLHTip("cybio", "CyBio50", 0.5, 50.0, "ul")
	tb = wtype.NewLHTipbox(8, 12, 60.13, "CyBio", "CyBio50Tipbox", tip, w, 9.0, 9.0, 0.0, 0.0, 0.0)
	tipboxes = append(tipboxes, tb)

	// these details are incorrect and need fixing
	w = wtype.NewLHWell("Cybio1000Tipbox", "", "A1", "ul", 1000.0, 50.0, shp, 0, 7.3, 7.3, 51.2, 0.0, "mm")
	w.Extra["InnerL"] = 5.6
	w.Extra["InnerW"] = 5.6
	tip = wtype.NewLHTip("cybio", "CyBio1000", 100.0, 1000.0, "ul")
	tb = wtype.NewLHTipbox(8, 12, 60.13, "CyBio", "CyBio1000Tipbox", tip, w, 9.0, 9.0, 0.0, 0.0, 0.0)
	tipboxes = append(tipboxes, tb)

	w = wtype.NewLHWell("Gilson200Tipbox", "", "A1", "ul", 200.0, 10.0, shp, 0, 7.3, 7.3, 51.2, 0.0, "mm")
	w.Extra["InnerL"] = 5.6
	w.Extra["InnerW"] = 5.6
	w.Extra["Tipeffectiveheight"] = 44.7
	tip = wtype.NewLHTip("gilson", "Gilson200", 10.0, 200.0, "ul")
	tb = wtype.NewLHTipbox(8, 12, 60.13, "Gilson", "DF200 Tip Rack (PIPETMAX 8x200)", tip, w, 9.0, 9.0, 0.0, 0.0, 24.78)
	tipboxes = append(tipboxes, tb)

	w = wtype.NewLHWell("Gilson20Tipbox", "", "A1", "ul", 20.0, 1.0, shp, 0, 7.3, 7.3, 46.0, 0.0, "mm")
	w.Extra["InnerL"] = 5.5
	w.Extra["InnerW"] = 5.5
	w.Extra["Tipeffectiveheight"] = 34.6
	tip = wtype.NewLHTip("gilson", "Gilson20", 0.5, 20.0, "ul")
	tb = wtype.NewLHTipbox(8, 12, 60.13, "Gilson", "DL10 Tip Rack (PIPETMAX 8x20)", tip, w, 9.0, 9.0, 0.0, 0.0, 28.93)
	tipboxes = append(tipboxes, tb)

	tipboxes = append(tipboxes, makeTecanTipBoxes()...)

	return tipboxes
}

func makeTecanTipBoxes() []*wtype.LHTipbox {
	shp := wtype.NewShape("cylinder", "mm", 7.3, 7.3, 51.2)

	ret := make([]*wtype.LHTipbox, 0, 4)

	w := wtype.NewLHWell("Tecan1000Tipbox", "", "A1", "ul", 1000.0, 200.0, shp, 0, 7.3, 7.3, 46.0, 0.0, "mm")
	tip := wtype.NewLHTip("Tecan", "Tecan1000", 200.0, 1000.0, "ul")
	tb := wtype.NewLHTipbox(8, 12, 60.13, "Tecan", "DiTi 1000uL LiHa", tip, w, 9.0, 9.0, 0.0, 0.0, 28.93)
	ret = append(ret, tb)

	w = wtype.NewLHWell("Tecan200Tipbox", "", "A1", "ul", 200.0, 15.0, shp, 0, 7.3, 7.3, 46.0, 0.0, "mm")
	tip = wtype.NewLHTip("Tecan", "Tecan200", 15.0, 200.0, "ul")
	tb = wtype.NewLHTipbox(8, 12, 60.13, "Tecan", "DiTi 200uL LiHa", tip, w, 9.0, 9.0, 0.0, 0.0, 28.93)
	ret = append(ret, tb)

	w = wtype.NewLHWell("Tecan50Tipbox", "", "A1", "ul", 50.0, 3.0, shp, 0, 7.3, 7.3, 46.0, 0.0, "mm")
	tip = wtype.NewLHTip("Tecan", "Tecan50", 3.0, 50.0, "ul")
	tb = wtype.NewLHTipbox(8, 12, 60.13, "Tecan", "DiTi 50uL LiHa", tip, w, 9.0, 9.0, 0.0, 0.0, 28.93)
	ret = append(ret, tb)

	w = wtype.NewLHWell("Tecan10Tipbox", "", "A1", "ul", 10.0, 1.0, shp, 0, 7.3, 7.3, 46.0, 0.0, "mm")
	tip = wtype.NewLHTip("Tecan", "Tecan10", 1.0, 10.0, "ul")
	tb = wtype.NewLHTipbox(8, 12, 60.13, "Tecan", "DiTi 10uL LiHa", tip, w, 9.0, 9.0, 0.0, 0.0, 28.93)
	ret = append(ret, tb)

	return ret
}
