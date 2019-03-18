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

import "github.com/antha-lang/antha/antha/anthalib/wtype"

func makeTipwastes() (tipwastes []*wtype.LHTipwaste) {
	tipwastes = append(tipwastes, makeGilsonTipWaste(), makeGilsonTipChute(), makeCyBioTipwaste(), makeManualTipwaste(), makeTecanTipwaste())
	return
}

func makeGilsonTipWaste() *wtype.LHTipwaste {
	shp := wtype.NewShape(wtype.BoxShape, "mm", 123.0, 80.0, 92.0)
	w := wtype.NewLHWell("ul", 800000.0, 800000.0, shp, 0, 123.0, 80.0, 92.0, 0.0, "mm")
	lht := wtype.NewLHTipwaste(6000, "Gilsontipwaste", "gilson", wtype.Coordinates3D{X: sbsX, Y: sbsY, Z: 92.0}, w, 49.5+xOffset, 31.5+yOffset, 0.0)
	return lht
}

//makeGilsonTipChute this is the chute for position 1 from direct measurements
func makeGilsonTipChute() *wtype.LHTipwaste {
	shp := wtype.NewShape(wtype.BoxShape, "mm", 50.0, 63.8, 82.98)
	w := wtype.NewLHWell("ul", 800000.0, 800000.0, shp, wtype.FlatWellBottom, 50.0, 63.8, 82.98, 0.0, "mm")
	lht := wtype.NewLHTipwaste(6000, "GilsonTipChute", "gilson", wtype.Coordinates3D{X: sbsX, Y: sbsY, Z: 50.0}, w, sbsX/2.0, sbsY/2.0, 0.0)
	return lht
}

// TODO figure out tip capacity
func makeCyBioTipwaste() *wtype.LHTipwaste {
	shp := wtype.NewShape(wtype.BoxShape, "mm", 90.5, 171.0, 90.0)
	w := wtype.NewLHWell("ul", 800000.0, 800000.0, shp, 0, 90.5, 171.0, 90.0, 0.0, "mm")
	lht := wtype.NewLHTipwaste(700, "CyBiotipwaste", "cybio", wtype.Coordinates3D{X: sbsX, Y: sbsY, Z: 90.5}, w, 85.5+xOffset, 45.0+yOffset, 0.0)
	return lht
}

// TODO figure out tip capacity
func makeManualTipwaste() *wtype.LHTipwaste {
	shp := wtype.NewShape(wtype.BoxShape, "mm", 90.5, 171.0, 90.0)
	w := wtype.NewLHWell("ul", 800000.0, 800000.0, shp, 0, 90.5, 171.0, 90.0, 0.0, "mm")
	lht := wtype.NewLHTipwaste(1000000, "Manualtipwaste", "ACMEBagsInc", wtype.Coordinates3D{X: sbsX, Y: sbsY, Z: 90.5}, w, 85.5+xOffset, 45.0+yOffset, 0.0)
	return lht
}

func makeTecanTipwaste() *wtype.LHTipwaste {
	shp := wtype.NewShape(wtype.BoxShape, "mm", 90.5, 171.0, 90.0)
	w := wtype.NewLHWell("ul", 800000.0, 800000.0, shp, 0, 90.5, 171.0, 90.0, 0.0, "mm")
	lht := wtype.NewLHTipwaste(2000, "Tecantipwaste", "Tecan", wtype.Coordinates3D{X: sbsX, Y: sbsY, Z: 90.5}, w, 85.5+xOffset, 45.0+yOffset, 0.0)
	return lht
}
