// /anthalib/factory/make_tip_waste_library.go: Part of the Antha language
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

import "github.com/antha-lang/antha/antha/anthalib/wtype"

func makeTipwastes() map[string]*wtype.LHTipwaste {
	ret := make(map[string]*wtype.LHTipwaste, 1)

	ret["Gilsontipwaste"] = makeGilsonTipWaste()
	ret["CyBiotipwaste"] = makeCyBioTipwaste()
	ret["Manualtipwaste"] = makeManualTipwaste()
	ret["Tecantipwaste"] = makeTecanTipwaste()
	return ret
}

func makeGilsonTipWaste() *wtype.LHTipwaste {
	shp := wtype.NewShape("box", "mm", 123.0, 80.0, 92.0)
	w := wtype.NewLHWell("Gilsontipwaste", "", "A1", "ul", 800000.0, 800000.0, shp, 0, 123.0, 80.0, 92.0, 0.0, "mm")
	lht := wtype.NewLHTipwaste(600, "gilsontipwaste", "gilson", 92.0, w, 49.5, 31.5, 0.0)
	return lht
}

// TODO figure out tip capacity
func makeCyBioTipwaste() *wtype.LHTipwaste {
	shp := wtype.NewShape("box", "mm", 90.5, 171.0, 90.0)
	w := wtype.NewLHWell("CyBiotipwaste", "", "A1", "ul", 800000.0, 800000.0, shp, 0, 90.5, 171.0, 90.0, 0.0, "mm")
	lht := wtype.NewLHTipwaste(700, "CyBiotipwaste", "cybio", 90.5, w, 85.5, 45.0, 0.0)
	return lht
}
func makeTecanTipwaste() *wtype.LHTipwaste {
	shp := wtype.NewShape("box", "mm", 90.5, 171.0, 90.0)
	w := wtype.NewLHWell("Tecantipwaste", "", "A1", "ul", 800000.0, 800000.0, shp, 0, 90.5, 171.0, 90.0, 0.0, "mm")
	lht := wtype.NewLHTipwaste(700, "Tecantipwaste", "cybio", 90.5, w, 85.5, 45.0, 0.0)
	return lht
}

// TODO figure out tip capacity
func makeManualTipwaste() *wtype.LHTipwaste {
	shp := wtype.NewShape("box", "mm", 90.5, 171.0, 90.0)
	w := wtype.NewLHWell("manualtipwaste", "", "A1", "ul", 800000.0, 800000.0, shp, 0, 90.5, 171.0, 90.0, 0.0, "mm")
	lht := wtype.NewLHTipwaste(1000000, "TipWasteBag", "ACMEBagsInc", 90.5, w, 85.5, 45.0, 0.0)
	return lht
}

func GetTipwasteByType(typ string) *wtype.LHTipwaste {
	tipwastes := makeTipwastes()
	t := tipwastes[typ]
	return t.Dup()
}

func TipwasteList() []string {
	tipwastes := makeTipwastes()
	kz := make([]string, len(tipwastes))
	x := 0
	for name, _ := range tipwastes {
		kz[x] = name
		x += 1
	}
	return kz
}
