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

package solutions

import (
	"fmt"
	"math"

	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/search"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
)

// Equal evaluates whether two components are equal.
// First the component names are normalised and evaluated for equality.
// If equal, component lists are retrieved for both components;
// if no sub components are found, concentrations are evaluated foe equality.
// If component lists are found the component lists are evaluated for equality.
// A discrepency of 1% difference in concentration is permitted.
func Equal(component1, component2 *wtype.LHComponent) error {

	if component1 == nil && component2 != nil {
		return fmt.Errorf("%s compared to nil component", component2.Name())
	}

	if component1 != nil && component2 == nil {
		return fmt.Errorf("%s compared to nil component", component1.Name())
	}

	if component1 == nil && component2 == nil {
		return nil
	}

	precision := 0.01

	if search.EqualFold(NormaliseName(component1.Name()), NormaliseName(component2.Name())) {

		compList1, _ := wtype.GetSubComponents(component1)
		compList2, _ := wtype.GetSubComponents(component2)

		if nonZeroComponents(compList1) == 0 && nonZeroComponents(compList2) == 0 {

			if component1.Concentration().RawValue() == 0 && component2.Concentration().RawValue() == 0 {
				return nil
			}

			factor, err := wunit.DivideConcentrations(component1.Concentration(), component2.Concentration())
			if err != nil {
				return err
			}
			if math.Abs(factor-1) <= precision {
				return nil
			}

			return fmt.Errorf("dilution factor of components not equal to 1 +/- 1%%: %f", factor)
		}

		if _, factor, err := DilutableComponentLists(compList1, compList2); math.Abs(factor-1) <= precision && err == nil {
			return nil
		} else if err != nil {
			return err
		} else if math.Abs(factor-1) > precision {
			return fmt.Errorf("dilution factor of components not equal to 1 +/- 1%%: %f", factor)
		}

	}
	return fmt.Errorf("Component names not equal after normalisation: %s and %s", NormaliseName(component1.Name()), NormaliseName(component2.Name()))
}

// Equivalent is a utility function to evaluate if components are equivalent based upon name and evaluating if both components contain a means of acquiring the concentration,
// either from parsing the component name or from the component properties.
// The concentration is not evaluated.
// Using ATP as an example component name, the function should work for components which are specified in the protocol in the following forms:
// 1uM ATP, ATP 1uM, 1g/L ATP, 1mMol/l ATP
func Equivalent(sourceComponent *wtype.LHComponent, targetComponent *wtype.LHComponent, lookForSubComponents bool) error {

	containsconc1, _, sourceComponentName := wunit.ParseConcentration(sourceComponent.CName)
	containsconc2, _, targetComponentName := wunit.ParseConcentration(targetComponent.CName)
	sourceCompList, _ := wtype.GetSubComponents(sourceComponent)
	targetCompList, _ := wtype.GetSubComponents(targetComponent)
	if containsconc1 && containsconc2 && search.EqualFold(sourceComponentName, targetComponentName) {

		if lookForSubComponents {

			equal, _, err := DilutableComponentLists(sourceCompList, targetCompList)
			if equal {
				return nil
			}
			return err
		}

		return nil

	} else if search.EqualFold(sourceComponentName, targetComponentName) && sourceComponent.HasConcentration() && targetComponent.HasConcentration() {
		if lookForSubComponents {

			equal, _, err := DilutableComponentLists(sourceCompList, targetCompList)
			if equal {
				return nil
			}

			return err

		}
		return nil
	}

	// don't look at name of component if sub component list is found for both: evaluate component lists
	if lookForSubComponents && nonZeroComponents(sourceCompList) > 0 && nonZeroComponents(targetCompList) > 0 {

		equal, _, err := DilutableComponentLists(sourceCompList, targetCompList)

		if equal {
			return nil
		}
		return err

	}
	if search.EqualFold(sourceComponentName, targetComponentName) {
		return nil
	}
	return fmt.Errorf("Components %s and %s not equivalent", ReturnNormalisedComponentName(sourceComponent), ReturnNormalisedComponentName(targetComponent))
}

// DilutableComponentLists evaluates whether sourceComponentList is dilutable to become equal to targetComponentList.
func DilutableComponentLists(sourceComponentList, targetComponentList wtype.ComponentList) (equal bool, dilutionFactor float64, err error) {

	if nonZeroComponents(sourceComponentList) == 0 && nonZeroComponents(targetComponentList) == 0 {
		return true, 1.0, nil
	}

	if nonZeroComponents(sourceComponentList) != nonZeroComponents(targetComponentList) {
		return false, -1, fmt.Errorf("componentlists unequal: %v, %v", sourceComponentList, targetComponentList)
	}

	//var dilutionFactorSoFar float64 = -1.0
	precision := 0.025
	var factors []float64
	for _, componentName := range sourceComponentList.AllComponents() {
		sourceConc, foundInSource := sourceComponentList.Components[componentName]
		targetConc, foundInTarget := targetComponentList.Components[componentName]

		if foundInSource != foundInTarget {
			return false, -1, fmt.Errorf("%s: found in source component list: %v; found in target component list: %v", componentName, foundInSource, foundInTarget)
		}

		var factor float64
		var err error

		if sourceConc.RawValue() == 0.0 && targetConc.RawValue() == 0.0 {
			// skip
			factor = 0.0
		} else {
			factor, err = wunit.DivideConcentrations(sourceConc, targetConc)
			if err != nil {
				return false, -1, fmt.Errorf("%s: %s. TargetComponentList: %v", componentName, err.Error(), targetComponentList)
			}

			if math.IsInf(factor, 0) {
				factor = 0.0
			}
		}
		factors = append(factors, factor)

	}

	allFactors := removeDuplicateFloats(factors, precision)
	if len(removeZeros(allFactors)) == 1 {
		return true, allFactors[0], nil
	} else if len(allFactors) == 1 {
		return true, 0, nil
	}

	return false, -1, fmt.Errorf("componentList %+v cannot be diluted to make componentList %+v. Factors needed for each subcomponent: %v", sourceComponentList, targetComponentList, allFactors)
}

func removeDuplicateFloats(elements []float64, precision float64) []float64 {
	// Use slice to record duplicates as we find them.
	var encountered []float64

	for v := range elements {

		for i := range encountered {
			if math.Abs(elements[v]-encountered[i]) < precision {
				// Do not add duplicate.
			} else {
				// Record this element as an encountered element.
				encountered = append(encountered, elements[v])

			}
		}
		if len(encountered) == 0 {
			encountered = append(encountered, elements[v])
		}
	}
	// Return the new slice.
	return encountered
}

func removeZeros(elements []float64) []float64 {
	var nonZero []float64
	for _, element := range elements {
		if element > 0.0 {
			nonZero = append(nonZero, element)
		}
	}
	return nonZero
}

func nonZeroComponents(compList wtype.ComponentList) int {
	var nonZero int
	for _, conc := range compList.Components {
		if conc.RawValue() > 0 {
			nonZero++
		}
	}
	return nonZero
}
