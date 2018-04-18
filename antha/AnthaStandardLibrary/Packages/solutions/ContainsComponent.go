// Part of the Antha language
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

// solutions is a utility package for working with solutions of LHComponents
package solutions

import (
	"fmt"
	"strings"

	"github.com/antha-lang/antha/antha/anthalib/wtype"
)

func equalFold(a, b string) bool {
	return strings.EqualFold(strings.TrimSpace(a), strings.TrimSpace(b))
}

// Looks for components matching name, concentration and all sub components (including their concentrations).
// A position of -1 is returned if no match found.
// If the component does not contain a concentration, the name will be matched only
// if multiple matches are found the first will be returned
func ContainsComponent(components []*wtype.LHComponent, component *wtype.LHComponent, lookForSubComponents bool) (found bool, position int, err error) {

	var errs []string

	// normalise names for more robust evaluation
	//var normalisedComponentName string

	//normalisedComponentName = NormaliseName(component.CName)
	//_,_, normalisedComponentName = wunit.ParseConcentration(component.CName)

	for i, comp := range components {

		if !comp.HasConcentration() {
			errs = append(errs, fmt.Sprintf("cannot compare component in list %s without a concentration", comp.CName))
		}

		// normalise names for more robust evaluation
		//var normalisedCompName string

		//normalisedCompName = NormaliseName(comp.CName)
		//_,_, normalisedCompName = wunit.ParseConcentration(comp.CName)

		if equalFold(comp.CName, component.CName) {

			if component.HasConcentration() && comp.HasConcentration() {

				if comp.Concentration().EqualTo(component.Concentration()) {

					if lookForSubComponents {

						compsubcomponents, err := GetSubComponents(comp)
						if err != nil {
							return false, -1, err
						}

						componentSubcomponents, err := GetSubComponents(component)
						if err != nil {
							return false, -1, err
						}

						err = EqualLists(compsubcomponents, componentSubcomponents)
						if err == nil {
							return true, i, nil
						} else {
							errs = append(errs, fmt.Sprintf("Subcomponents lists not equal for %s and %s: %s", comp.CName, component.CName, err.Error()))
						}
					} else {
						return true, i, nil
					}
				} else {
					errs = append(errs, comp.CName+"concentration "+comp.Concentration().ToString()+" not equal to "+component.CName+" "+component.Concentration().ToString())
				}
			} else {
				if lookForSubComponents {

					compsubcomponents, err := GetSubComponents(comp)
					if err != nil {
						return false, -1, err
					}

					componentSubcomponents, err := GetSubComponents(component)
					if err != nil {
						return false, -1, err
					}
					err = EqualLists(compsubcomponents, componentSubcomponents)
					if err == nil {
						return true, i, nil
					} else {
						errs = append(errs, fmt.Sprintf("Subcomponents lists not equal for %s and %s: %s", comp.CName, component.CName, err.Error()))
					}
				} else {
					return true, i, nil
				}
			}
		} else {
			errs = append(errs, comp.CName+" name not equal to "+component.CName)
		}
	}

	return false, -1, fmt.Errorf("component %s not found in list: %s. : Errors for each: %s", componentSummary(component), componentNames(components), strings.Join(errs, "\n"))
}

func nonZeroComponents(compList ComponentList) int {
	var nonZero int
	for _, conc := range compList.Components {
		if conc.RawValue() > 0 {
			nonZero++
		}
	}
	return nonZero
}

// EqualLists compares two ComponentLists and returns an error if the lists are not identical.
func EqualLists(list1, list2 ComponentList) error {
	var notEqual []string

	if nonZeroComponents(list1) == 0 && nonZeroComponents(list2) == 0 {
		return nil
	}

	if nonZeroComponents(list1) != nonZeroComponents(list2) {
		return fmt.Errorf("componentlists unequal length: %d, %d", nonZeroComponents(list1), nonZeroComponents(list2))
	}

	for key, value1 := range list1.Components {
		if value2, found := list2.Components[key]; found {
			if fmt.Sprintf("%.2e", value1.SIValue()) != fmt.Sprintf("%.2e", value2.SIValue()) {
				notEqual = append(notEqual, key+" "+fmt.Sprint(value1)+" in list 1 and "+fmt.Sprint(value2)+" in list 2.")
			}
		} else {
			notEqual = append(notEqual, key+" not found in list2. ")
		}
	}
	if len(notEqual) > 0 {
		return fmt.Errorf(strings.Join(notEqual, ". \n"))
	}
	return nil
}

func componentSummary(component *wtype.LHComponent) string {
	subComps, err := GetSubComponents(component)
	var message string
	if err != nil {
		message = err.Error()
	} else {
		message = subComps.List(true)
	}

	conc := "No concentration found"

	if component.HasConcentration() {
		conc = component.Concentration().ToString()
	}

	return fmt.Sprint("Component Name: ", component.CName, "Conc: ", conc, ". SubComponents: ", message)
}

func componentNames(components []*wtype.LHComponent) (names []string) {
	for _, component := range components {
		names = append(names, component.CName)
	}
	return
}
