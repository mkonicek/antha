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
	"github.com/antha-lang/antha/antha/anthalib/wutil/text"
)

// solutionSetup determines how to fulfil the requirements for making instructions to specifications
func (rq *LHRequest) solutionSetup() (wtype.LHInstructions, map[string]wunit.Concentration, error) {

	// set this from extra or calculate, but skip for now
	var skipSampleForConcentrationCalc bool = true

	instructions := rq.LHInstructions

	// index of components used to make up to a total volume, along with the required total
	mtvols := make(map[string][]wunit.Volume, 10)
	// index of components with concentration targets, along with the target concentrations
	mconcs := make(map[string][]wunit.Concentration, 10)
	// keep a list of components which have fixed stock concentrations
	fixconcs := make([]*wtype.Liquid, 0)
	// maximum solubilities of each component
	Smax := make(map[string]float64, 10)
	// maximum total volume of any instruction containing each component
	hshTVol := make(map[string]wunit.Volume)

	// find the minimum and maximum required concentrations
	// across all the instructions

	// -- migrate this to chains of dependent instructions
	for _, instruction := range instructions {
		components := instruction.Inputs

		// we need to identify the concentration components
		// and the total volume components, if we have
		// concentrations but no tvols we have to return
		// an error

		arrCncs := make([]*wtype.Liquid, 0, len(components))
		arrTvol := make([]*wtype.Liquid, 0, len(components))
		cmpvol := wunit.NewVolume(0.0, "ul")
		totalvol := wunit.NewVolume(0.0, "ul")

		for _, component := range components {

			// what sort of component is it?
			conc := component.Conc
			tvol := component.Tvol
			if conc != 0.0 && !skipSampleForConcentrationCalc {
				arrCncs = append(arrCncs, component)
			} else if tvol != 0.0 {
				tv := component.TotalVolume()
				if totalvol.IsZero() || totalvol.EqualTo(tv) {
					totalvol = tv // not needed
				} else {
					return nil, nil, wtype.LHErrorf(wtype.LH_ERR_CONC, "Inconsistent total volumes %s and %s at component %s", totalvol, tv, component.CName)
				}
			} else {
				cmpvol.Add(component.Volume())
			}
		}

		if len(arrCncs) > 0 {
			fmt.Println(text.Blue(fmt.Sprintf("arrCncs: %+v", arrCncs)))
		}

		// add everything to the maps

		for _, cmp := range arrCncs {
			nm := cmp.CName
			cnc := wunit.NewConcentration(cmp.Conc, cmp.Cunit)

			_, ok := Smax[nm]

			if !ok {
				Smax[nm] = cmp.Smax
			}

			if cmp.StockConcentration != 0.0 {
				fixconcs = append(fixconcs, cmp)
				continue
			}

			var cncslc []wunit.Concentration

			cncslc, ok = mconcs[nm]

			if !ok {
				cncslc = make([]wunit.Concentration, 0, 10)
			}

			cncslc = append(cncslc, cnc)

			mconcs[nm] = cncslc
			_, ok = hshTVol[nm]
			if !ok || hshTVol[nm].GreaterThan(totalvol) {
				hshTVol[nm] = totalvol
			}
		}

		// now the total volumes

		for _, cmp := range arrTvol {
			nm := cmp.CName
			tvol := cmp.TotalVolume()

			var tvslc []wunit.Volume

			tvslc, ok := mtvols[nm]

			if !ok {
				tvslc = make([]wunit.Volume, 0, 10)
			}

			tvslc = append(tvslc, tvol)

			mtvols[nm] = tvslc
		}

	} // end instructions
	// so now we should be able to make stock concentrations
	// first we need the min and max for each

	minrequired := make(map[string]wunit.Concentration, len(mconcs))
	maxrequired := make(map[string]wunit.Concentration, len(mconcs))

	if len(mconcs) > 0 {
		fmt.Println(text.Green(fmt.Sprintf("mconcs: %+v", mconcs)))
	}

	for cmp, arr := range mconcs {
		min, err := wunit.MinConcentration(arr)
		if err != nil {
			return nil, nil, err
		}
		max, err := wunit.MaxConcentration(arr)
		if err != nil {
			return nil, nil, err
		}
		minrequired[cmp] = min
		maxrequired[cmp] = max
		// if smax undefined we need to deal  - we assume infinite solubility!!

		_, ok := Smax[cmp]

		if !ok {
			Smax[cmp] = 9999999
			fmt.Printf("Max solubility undefined for component %s -- assuming infinite solubility!\n", cmp)
		}

	}

	_, minUnit, err := convertToSIValues(minrequired)

	if err != nil && len(minrequired) > 0 {
		return nil, nil, err
	}

	_, maxUnit, err := convertToSIValues(maxrequired)

	if err != nil && len(maxrequired) > 0 {
		return nil, nil, err
	}

	if minUnit != maxUnit {
		return nil, nil, fmt.Errorf("min unit %s not equal to max unit %s ", minUnit, maxUnit)
	}

	stockconcs := make(map[string]float64)

	// handle any errors here

	// add the fixed concentrations into stockconcs

	for _, cmp := range fixconcs {
		stockconcs[cmp.CName] = cmp.StockConcentration
	}

	// nearly there now! Need to turn all the components into volumes, then we're done

	// make an array for the new instructions

	newInstructions := make(wtype.LHInstructions, len(instructions))

	for _, instruction := range instructions {
		components := instruction.Inputs
		arrCncs := make([]*wtype.Liquid, 0, len(components))
		arrTvol := make([]*wtype.Liquid, 0, len(components))
		arrSvol := make([]*wtype.Liquid, 0, len(components))
		cmpvol := wunit.NewVolume(0.0, "ul")
		totalvol := wunit.NewVolume(0.0, "ul")

		for _, component := range components {
			// what sort of component is it?
			// what is the total volume ?
			if component.Conc != 0.0 && !skipSampleForConcentrationCalc {
				arrCncs = append(arrCncs, component)
			} else if component.Tvol != 0.0 {
				arrTvol = append(arrTvol, component)
				tv := component.TotalVolume()
				if totalvol.IsZero() || totalvol.EqualTo(tv) {
					totalvol = tv
				} else {
					return nil, nil, wtype.LHErrorf(wtype.LH_ERR_CONC, "Inconsistent total volumes %s and %s at component %s", totalvol, tv, component.CName)
				}
			} else {
				// need to add in the volume taken up by any volume components
				cmpvol.Add(component.Volume())
				arrSvol = append(arrSvol, component)
			}
		}

		// first we add the volumes to the concentration components

		arrFinalComponents := make([]*wtype.Liquid, 0, len(components))

		for _, component := range arrCncs {
			name := component.CName
			cnc := component.Conc
			//vol := totalvol * cnc / stockconcs[name]
			vol := wunit.MultiplyVolume(totalvol, cnc/stockconcs[name])
			cmpvol.Add(vol)
			component.Vol = vol.RawValue()
			component.Vunit = totalvol.Unit().PrefixedSymbol()
			component.StockConcentration = stockconcs[name]
			arrFinalComponents = append(arrFinalComponents, component)
		}

		// next we get the final volume for total volume components

		for _, component := range arrTvol {
			vol := wunit.SubtractVolumes(totalvol, cmpvol)
			if vol.IsNegative() {
				return nil, nil, wtype.LHErrorf(wtype.LH_ERR_VOL, "invalid total volume for component %q in instruction:\n%s", component.CName, instruction.Summarize(1))
			}
			component.SetVolume(vol)
			component.Tvol = 0.0 // reset Tvol
			arrFinalComponents = append(arrFinalComponents, component)
		}

		// then we add the rest

		arrFinalComponents = append(arrFinalComponents, arrSvol...)

		// finally we replace the components in this instruction

		instruction.Inputs = arrFinalComponents

		// and put the new instruction in the array

		newInstructions[instruction.ID] = instruction
	}

	if len(fixconcs) > 0 {
		fmt.Println(text.Red(fmt.Sprintf("fixconcs: %+v", fixconcs)))
	}

	stockConcs, err := convertFloatsToConc(stockconcs, minUnit)

	if err != nil && len(stockconcs) > 0 {
		return newInstructions, stockConcs, err
	}

	return newInstructions, stockConcs, nil
}

// converts all to SI Values, all entries must have the same SI base unit or an error will be returned.
func convertToSIValues(concMap map[string]wunit.Concentration) (floats map[string]float64, unit string, err error) {
	floats = make(map[string]float64, len(concMap))
	var concSlice []wunit.Concentration

	for _, conc := range concMap {
		concSlice = append(concSlice, conc)
	}

	_, err = wunit.SortConcentrations(concSlice)
	if err != nil {
		err = fmt.Errorf("Cannot convert concentration map to floats: %s", err.Error())
		return
	}

	for key, concValue := range concMap {
		floats[key] = concValue.SIValue()
	}

	unit = concSlice[0].Unit().PrefixedSymbol()

	return
}

// converts all float values to concentration values with specified unit
func convertFloatsToConc(floatMap map[string]float64, unit string) (map[string]wunit.Concentration, error) {

	reg := wunit.GetGlobalUnitRegistry()

	if !reg.ValidUnitForType("Concentration", unit) {
		return nil, fmt.Errorf("unapproved concentration unit %q, approved units are %v", unit, reg.ListValidUnitsForType("Concentration"))
	}

	concMap := make(map[string]wunit.Concentration, len(floatMap))

	for key, concValue := range floatMap {
		concMap[key] = wunit.NewConcentration(concValue, unit)
	}
	return concMap, nil
}
