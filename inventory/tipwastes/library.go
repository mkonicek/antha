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

package tipwastes

import (
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/inventory/tipboxes"
	"github.com/antha-lang/antha/laboratory/effects/id"
)

func makeTipwastes(idGen *id.IDGenerator) (tipwastes []*wtype.LHTipwaste) {
	tipwastes = append(tipwastes, makeGilsonTipWaste(idGen), makeGilsonTipChute(idGen), makeCyBioTipwaste(idGen), makeManualTipwaste(idGen), makeTecanTipwaste(idGen))
	return
}

func makeGilsonTipWaste(idGen *id.IDGenerator) *wtype.LHTipwaste {
	shp := wtype.NewShape("box", "mm", 123.0, 80.0, 92.0)
	w := wtype.NewLHWell(idGen, "ul", 800000.0, 800000.0, shp, 0, 123.0, 80.0, 92.0, 0.0, "mm")
	lht := wtype.NewLHTipwaste(idGen, 6000, "Gilsontipwaste", "gilson", wtype.Coordinates{X: tipboxes.SbsX, Y: tipboxes.SbsY, Z: 92.0}, w, 49.5+tipboxes.XOffset, 31.5+tipboxes.YOffset, 0.0)
	return lht
}

//makeGilsonTipChute this is the chute for position 1 from direct measurements
func makeGilsonTipChute(idGen *id.IDGenerator) *wtype.LHTipwaste {
	shp := wtype.NewShape("box", "mm", 50.0, 63.8, 82.98)
	w := wtype.NewLHWell(idGen, "ul", 800000.0, 800000.0, shp, wtype.FlatWellBottom, 50.0, 63.8, 82.98, 0.0, "mm")
	lht := wtype.NewLHTipwaste(idGen, 6000, "GilsonTipChute", "gilson", wtype.Coordinates{X: tipboxes.SbsX, Y: tipboxes.SbsY, Z: 50.0}, w, tipboxes.SbsX/2.0, tipboxes.SbsY/2.0, 0.0)
	return lht
}

// TODO figure out tip capacity
func makeCyBioTipwaste(idGen *id.IDGenerator) *wtype.LHTipwaste {
	shp := wtype.NewShape("box", "mm", 90.5, 171.0, 90.0)
	w := wtype.NewLHWell(idGen, "ul", 800000.0, 800000.0, shp, 0, 90.5, 171.0, 90.0, 0.0, "mm")
	lht := wtype.NewLHTipwaste(idGen, 700, "CyBiotipwaste", "cybio", wtype.Coordinates{X: tipboxes.SbsX, Y: tipboxes.SbsY, Z: 90.5}, w, 85.5+tipboxes.XOffset, 45.0+tipboxes.YOffset, 0.0)
	return lht
}

// TODO figure out tip capacity
func makeManualTipwaste(idGen *id.IDGenerator) *wtype.LHTipwaste {
	shp := wtype.NewShape("box", "mm", 90.5, 171.0, 90.0)
	w := wtype.NewLHWell(idGen, "ul", 800000.0, 800000.0, shp, 0, 90.5, 171.0, 90.0, 0.0, "mm")
	lht := wtype.NewLHTipwaste(idGen, 1000000, "Manualtipwaste", "ACMEBagsInc", wtype.Coordinates{X: tipboxes.SbsX, Y: tipboxes.SbsY, Z: 90.5}, w, 85.5+tipboxes.XOffset, 45.0+tipboxes.YOffset, 0.0)
	return lht
}

func makeTecanTipwaste(idGen *id.IDGenerator) *wtype.LHTipwaste {
	shp := wtype.NewShape("box", "mm", 90.5, 171.0, 90.0)
	w := wtype.NewLHWell(idGen, "ul", 800000.0, 800000.0, shp, 0, 90.5, 171.0, 90.0, 0.0, "mm")
	lht := wtype.NewLHTipwaste(idGen, 2000, "Tecantipwaste", "Tecan", wtype.Coordinates{X: tipboxes.SbsX, Y: tipboxes.SbsY, Z: 90.5}, w, 85.5+tipboxes.XOffset, 45.0+tipboxes.YOffset, 0.0)
	return lht
}
