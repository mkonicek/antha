// anthalib//factory/make_tip_library.go: Part of the Antha language
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

package factory

import (
	"fmt"

	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/microArch/logger"
)

func makeTipLibrary() map[string]*wtype.LHTipbox {
	tips := make(map[string]*wtype.LHTipbox)

	// create a well representation of the tip holder... sometimes needed
	// heh, should have kept LHTipholder!
	cybio50tip, cybio50tipbox := makeCyBio50Tips()
	tips[cybio50tip.Type] = cybio50tipbox
	tips[cybio50tipbox.Type] = cybio50tipbox

	cybio250tip, cybio250tipbox := makeCyBio250Tips()

	tips[cybio250tip.Type] = cybio250tipbox
	tips[cybio250tipbox.Type] = cybio250tipbox

	cybio1ktip, cybio1ktipbox := makeCyBio1000Tips()
	tips[cybio1ktip.Type] = cybio1ktipbox
	tips[cybio1ktipbox.Type] = cybio1ktipbox

	g200Tip, g200Tipbox := makeGilson200Tips()
	tips[g200Tip.Type] = g200Tipbox
	tips[g200Tipbox.Type] = g200Tipbox

	gilson20Tip, gilson20Tipbox := makeGilson20Tips()
	tips[gilson20Tip.Type] = gilson20Tipbox
	tips[gilson20Tipbox.Type] = gilson20Tipbox

	tecan200Tip, tecan200Tipbox := makeTecan200Tips()
	tips[tecan200Tip.Type] = tecan200Tipbox
	tips[tecan200Tipbox.Type] = tecan200Tipbox

	return tips
}

func makeGilson200Tips() (*wtype.LHTip, *wtype.LHTipbox) {
	shp := wtype.NewShape("cylinder", "mm", 7.3, 7.3, 51.2)
	w := wtype.NewLHWell("Gilson200Tipbox", "", "A1", "ul", 200.0, 10.0, shp, 0, 7.3, 7.3, 51.2, 0.0, "mm")
	w.Extra["InnerL"] = 5.6
	w.Extra["InnerW"] = 5.6
	w.Extra["Tipeffectiveheight"] = 44.7
	tip := wtype.NewLHTip("gilson", "Gilson200", 10.0, 200.0, "ul")
	tb := wtype.NewLHTipbox(8, 12, 60.13, "Gilson", "DF200 Tip Rack (PIPETMAX 8x200)", tip, w, 9.0, 9.0, 0.0, 0.0, 24.78)

	return tip, tb
}
func makeGilson20Tips() (*wtype.LHTip, *wtype.LHTipbox) {
	shp := wtype.NewShape("cylinder", "mm", 7.3, 7.3, 51.2)
	w := wtype.NewLHWell("Gilson20Tipbox", "", "A1", "ul", 20.0, 1.0, shp, 0, 7.3, 7.3, 46.0, 0.0, "mm")
	w.Extra["InnerL"] = 5.5
	w.Extra["InnerW"] = 5.5
	w.Extra["Tipeffectiveheight"] = 34.6
	tip := wtype.NewLHTip("gilson", "Gilson20", 0.5, 20.0, "ul")
	tb := wtype.NewLHTipbox(8, 12, 60.13, "Gilson", "DL10 Tip Rack (PIPETMAX 8x20)", tip, w, 9.0, 9.0, 0.0, 0.0, 28.93)

	return tip, tb
}
func makeCyBio50Tips() (*wtype.LHTip, *wtype.LHTipbox) {
	shp := wtype.NewShape("cylinder", "mm", 7.3, 7.3, 51.2)
	w := wtype.NewLHWell("CyBio50Tipbox", "", "A1", "ul", 50.0, 0.5, shp, 0, 7.3, 7.3, 51.2, 0.0, "mm")
	w.Extra["InnerL"] = 5.6
	w.Extra["InnerW"] = 5.6

	tip := wtype.NewLHTip("cybio", "CyBio50", 0.5, 50.0, "ul")
	tb := wtype.NewLHTipbox(8, 12, 60.13, "CyBio", "CyBio50Tipbox", tip, w, 9.0, 9.0, 0.0, 0.0, 0.0)

	return tip, tb
}
func makeCyBio250Tips() (*wtype.LHTip, *wtype.LHTipbox) {
	shp := wtype.NewShape("cylinder", "mm", 7.3, 7.3, 51.2)
	w := wtype.NewLHWell("CyBio250Tipbox", "", "A1", "ul", 250.0, 10.0, shp, 0, 7.3, 7.3, 51.2, 0.0, "mm")
	w.Extra["InnerL"] = 5.6
	w.Extra["InnerW"] = 5.6
	tip := wtype.NewLHTip("cybio", "CyBio250", 10.0, 250.0, "ul")
	tb := wtype.NewLHTipbox(8, 12, 60.13, "CyBio", "CyBio250Tipbox", tip, w, 9.0, 9.0, 0.0, 0.0, 0.0)

	return tip, tb
}
func makeCyBio1000Tips() (*wtype.LHTip, *wtype.LHTipbox) {
	shp := wtype.NewShape("cylinder", "mm", 7.3, 7.3, 51.2)
	w := wtype.NewLHWell("Cybio1000Tipbox", "", "A1", "ul", 1000.0, 50.0, shp, 0, 7.3, 7.3, 51.2, 0.0, "mm")
	w.Extra["InnerL"] = 5.6
	w.Extra["InnerW"] = 5.6
	tip := wtype.NewLHTip("cybio", "CyBio1000", 100.0, 1000.0, "ul")
	tb := wtype.NewLHTipbox(8, 12, 60.13, "CyBio", "CyBio1000Tipbox", tip, w, 9.0, 9.0, 0.0, 0.0, 0.0)

	return tip, tb
}

func makeTecan200Tips() (*wtype.LHTip, *wtype.LHTipbox) {
	shp := wtype.NewShape("cylinder", "mm", 7.3, 7.3, 51.2)
	w := wtype.NewLHWell("Tecan200Tipbox", "", "A1", "ul", 200.0, 3.0, shp, 0, 7.3, 7.3, 51.2, 0.0, "mm")
	w.Extra["InnerL"] = 5.6
	w.Extra["InnerW"] = 5.6
	w.Extra["Tipeffectiveheight"] = 44.7
	tip := wtype.NewLHTip("Tecan", "Tecan200", 1.0, 200.0, "ul")
	tb := wtype.NewLHTipbox(8, 12, 60.13, "Tecan", "Tecan200Tiprack", tip, w, 9.0, 9.0, 0.0, 0.0, 24.78)
	return tip, tb
}

func GetTipBoxByTip(tip *wtype.LHTip) *wtype.LHTipbox {
	return GetTipByType(tip.Type)
}

func GetTipboxByType(typ string) *wtype.LHTipbox {
	return GetTipByType(typ)
}

func GetTipByType(typ string) *wtype.LHTipbox {
	tips := makeTipLibrary()
	t := tips[typ]

	if t == nil {
		logger.Debug(fmt.Sprintln("NO TIP TYPE: ", typ))
		return nil
	}

	return t.Dup()
}

func GetTipList() []string {
	tips := makeTipLibrary()
	kz := make([]string, len(tips))
	x := 0
	for name, _ := range tips {
		kz[x] = name
		x += 1
	}
	return kz
}
