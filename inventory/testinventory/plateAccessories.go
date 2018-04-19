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
	"strings"
)

// The height below which an error will be generated
// when attempting to perform transfers with low volume head and tips (0.5 - 20ul) on the Gilson PipetMax.
const MinimumZHeightPermissableForLVPipetMax = 0.636

//var commonwelltypes

var platespecificoffset = map[string]float64{
	"pcrplate_skirted": gilsonoffsetpcrplate,
	"greiner384":       gilsonoffsetgreiner,
	"costar48well":     3.0,
	"Nuncon12well":     11.0, // this must be wrong!! check z start without riser properly
	"Nuncon12wellAgar": 11.0, // this must be wrong!! check z start without riser properly
	"VWR12well":        3.0,
}

// function to check if a platename already contains a riser
func containsRiser(plate *wtype.LHPlate) bool {
	for _, dev := range defaultDevices {
		for _, synonym := range dev.GetSynonyms() {
			if strings.Contains(plate.Type, "_"+synonym) {
				return true
			}
		}
	}

	return false
}

func addRiser(plate *wtype.LHPlate, riser device) (plates []*wtype.LHPlate) {
	if containsRiser(plate) || doNotAddThisRiserToThisPlate(plate, riser) {
		return
	}

	for _, risername := range riser.GetSynonyms() {
		var dontaddrisertothisplate bool

		newplate := plate.Dup()
		riserheight := riser.GetHeightInmm()
		if offset, found := platespecificoffset[plate.Type]; found {
			riserheight = riserheight - offset
		}

		riserheight = riserheight + plateRiserSpecificOffset(plate, riser)

		newplate.WellZStart = plate.WellZStart + riserheight
		newname := plate.Type + "_" + risername
		newplate.Type = newname
		if riser.GetConstraints() != nil {
			// duplicate well before adding constraint to prevent applying
			// constraint to all common &Welltype on other plates

			for device, allowedpositions := range riser.GetConstraints() {
				newwell := newplate.Welltype.Dup()
				newplate.Welltype = newwell
				_, ok := newwell.Extra[device]
				if !ok {
					newplate.SetConstrained(device, allowedpositions)
				} else {
					dontaddrisertothisplate = true
				}
			}
		}

		if !dontaddrisertothisplate {
			plates = append(plates, newplate)
		}
	}

	return
}

func addAllDevices(plates []*wtype.LHPlate) (ret []*wtype.LHPlate) {
	for _, plate := range plates {
		for _, dev := range defaultDevices {
			ret = append(ret, addRiser(plate, dev)...)
		}
	}
	return
}
