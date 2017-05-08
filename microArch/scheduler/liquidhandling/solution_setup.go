// anthalib//liquidhandling/solution_setup.go: Part of the Antha language
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

package liquidhandling

import (
	"fmt"

	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/antha/anthalib/wutil"
	"github.com/antha-lang/antha/microArch/driver/liquidhandling"
)

// determines how to
// fulfil the requirements for making
// instructions to specifications

func solution_setup(request *LHRequest, prms *liquidhandling.LHProperties) (map[string]*wtype.LHInstruction, map[string]float64, error) {
	instructions := request.LHInstructions

	// index of components used to make up to a total volume, along with the required total
	mtvols := make(map[string][]float64, 10)
	// index of components with concentration targets, along with the target concentrations
	mconcs := make(map[string][]float64, 10)
	// keep a list of components which have fixed stock concentrations
	fixconcs := make([]*wtype.LHComponent, 0)
	// maximum solubilities of each component
	Smax := make(map[string]float64, 10)
	// maximum total volume of any instruction containing each component
	hshTVol := make(map[string]float64)

	// find the minimum and maximum required concentrations
	// across all the instructions

	// -- migrate this to chains of dependent instructions
	for _, instruction := range instructions {
		components := instruction.Components

		// we need to identify the concentration components
		// and the total volume components, if we have
		// concentrations but no tvols we have to return
		// an error

		arrCncs := make([]*wtype.LHComponent, 0, len(components))
		arrTvol := make([]*wtype.LHComponent, 0, len(components))
		cmpvol := 0.0
		totalvol := 0.0

		for _, component := range components {
			// what sort of component is it?
			conc := component.Conc
			tvol := component.Tvol
			if conc != 0.0 {
				arrCncs = append(arrCncs, component)
			} else if tvol != 0.0 {
				tv := component.Tvol
				if totalvol == 0.0 || totalvol == tv {
					totalvol = tv
				} else {
					// error
					wtype.LHError(wtype.LH_ERR_CONC, fmt.Sprintf("Inconsistent total volumes %-6.4f and %-6.4f at component %s", totalvol, tv, component.Name))
				}
			} else {
				cmpvol += component.Vol
			}
		}

		// add everything to the maps

		for _, cmp := range arrCncs {
			nm := cmp.CName
			cnc := cmp.Conc

			_, ok := Smax[nm]

			if !ok {
				Smax[nm] = cmp.Smax
			}

			if cmp.StockConcentration != 0.0 {
				fixconcs = append(fixconcs, cmp)
				continue
			}

			var cncslc []float64

			cncslc, ok = mconcs[nm]

			if !ok {
				cncslc = make([]float64, 0, 10)
			}

			cncslc = append(cncslc, cnc)

			mconcs[nm] = cncslc
			_, ok = hshTVol[nm]
			if !ok || hshTVol[nm] > totalvol {
				hshTVol[nm] = totalvol
			}
		}

		// now the total volumes

		for _, cmp := range arrTvol {
			nm := cmp.CName
			tvol := cmp.Tvol

			var tvslc []float64

			tvslc, ok := mtvols[nm]

			if !ok {
				tvslc = make([]float64, 0, 10)
			}

			tvslc = append(tvslc, tvol)

			mtvols[nm] = tvslc
		}

	} // end instructions
	// so now we should be able to make stock concentrations
	// first we need the min and max for each

	minrequired := make(map[string]float64, len(mconcs))
	maxrequired := make(map[string]float64, len(mconcs))

	//TODO this needs to be migrated elsewhere
	var vmin wunit.Volume = wunit.NewVolume(1.0, "ul")

	//	fmt.Println("PRMS: ", prms)

	if prms.CurrConf != nil && !prms.CurrConf.Minvol.LessThanFloat(0.00000001) {
		vmin = prms.CurrConf.Minvol
	}

	for cmp, arr := range mconcs {
		min, _ := wutil.FMin(arr)
		max, _ := wutil.FMax(arr)
		minrequired[cmp] = min
		maxrequired[cmp] = max
		// if smax undefined we need to deal  - we assume infinite solubility!!

		_, ok := Smax[cmp]

		if !ok {
			Smax[cmp] = 9999999
			wutil.Warn(fmt.Sprintf("Max solubility undefined for component %s -- assuming infinite solubility!", cmp))
		}

	}

	stockconcs := choose_stock_concentrations(minrequired, maxrequired, Smax, vmin.RawValue(), hshTVol)

	// handle any errors here

	// add the fixed concentrations into stockconcs

	for _, cmp := range fixconcs {
		stockconcs[cmp.CName] = cmp.StockConcentration
	}

	// nearly there now! Need to turn all the components into volumes, then we're done

	// make an array for the new instructions

	newInstructions := make(map[string]*wtype.LHInstruction, len(instructions))

	for _, instruction := range instructions {
		components := instruction.Components
		arrCncs := make([]*wtype.LHComponent, 0, len(components))
		arrTvol := make([]*wtype.LHComponent, 0, len(components))
		arrSvol := make([]*wtype.LHComponent, 0, len(components))
		cmpvol := 0.0
		totalvol := 0.0
		totalvolunit := ""

		for _, component := range components {
			// what sort of component is it?
			// what is the total volume ?
			if component.Conc != 0.0 {
				arrCncs = append(arrCncs, component)
			} else if component.Tvol != 0.0 {
				arrTvol = append(arrTvol, component)
				tv := component.Tvol
				totalvolunit = component.Vunit
				if totalvol == 0.0 || totalvol == tv {
					totalvol = tv
				} else {
					// error
					wtype.LHError(wtype.LH_ERR_CONC, fmt.Sprintf("Inconsistent total volumes %-6.4f and %-6.4f at component %s", totalvol, tv, component.Name))
				}
			} else {
				// need to add in the volume taken up by any volume components
				cmpvol += component.Vol
				arrSvol = append(arrSvol, component)
			}
		}

		// first we add the volumes to the concentration components

		arrFinalComponents := make([]*wtype.LHComponent, 0, len(components))

		for _, component := range arrCncs {
			name := component.CName
			cnc := component.Conc
			vol := totalvol * cnc / stockconcs[name]
			cmpvol += vol
			component.Vol = vol
			component.Vunit = totalvolunit
			component.StockConcentration = stockconcs[name]
			arrFinalComponents = append(arrFinalComponents, component)
		}

		// next we get the final volume for total volume components

		for _, component := range arrTvol {
			vol := totalvol - cmpvol
			component.Vol = vol
			arrFinalComponents = append(arrFinalComponents, component)
		}

		// then we add the rest

		arrFinalComponents = append(arrFinalComponents, arrSvol...)

		// finally we replace the components in this instruction

		instruction.Components = arrFinalComponents

		// and put the new instruction in the array

		newInstructions[instruction.ID] = instruction
	}

	return newInstructions, stockconcs, nil
}
