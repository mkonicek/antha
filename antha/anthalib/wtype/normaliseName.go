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
package wtype

import (
	"strings"

	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/antha/anthalib/wutil"
)

// NormaliseName evaluates whether a string contains a concentration and attempts to normalise the name in a standardised format.
// e.g. 10ng/ul glucose will be normalised to 10 mg/l glucose or 10mM glucose to 10 mM/l glucose or 10mM/l glucose to 10 mM/l glucose or glucose 10mM/l to 10 mM/l glucose
// A concatanenated name such as 10g/L glucose + 10g/L yeast extract will be returned with no modifications
func NormaliseName(name string) (normalised string) {

	if strings.Contains(name, wutil.MIXDELIMITER) {
		return name
	}

	containsConc, conc, nameonly := wunit.ParseConcentration(strings.TrimSpace(name))

	if containsConc {
		if conc.RawValue() > 0 {
			return conc.ToString() + " " + nameonly
		} else {
			return nameonly
		}
	} else {
		return name
	}
}

func removeConcUnitFromName(name string) string {

	for n := 0; n < 9; n++ {
		_, _, name = wunit.ParseConcentration(strings.TrimSpace(name))
	}
	newName := name
	return newName
}

// ReturnNormalisedComponentName will return the component name in a normalised form.
// If sub components exist the name will be changed to the list of sub components with concentrations.
// e.g. a solution called LB with a concentration of 10X and components 10g/L Yeast Extract and 5g/L Tryptone will be normalised to 10g/L Yeast Extract + 5g/L Tryptone.
// An LB solution with concentration 1 X and no components is returned as 1X LB.
// An LB solution with no concentration and no components is returned as LB.
func ReturnNormalisedComponentName(component *Liquid) string {
	originalcompList, _ := GetSubComponents(component)

	compList := originalcompList.RemoveConcsFromSubComponentNames()

	if len(compList.Components) == 0 {
		if component.HasConcentration() {
			name := component.Concentration().ToString() + " " + removeConcUnitFromName(component.Name()) + " " + compList.List(false)
			return name
		} else {
			return removeConcUnitFromName(component.Name())
		}
	}

	name := compList.List(false, true)

	return name
}

// NormaliseComponentName will change the name of the component to the normalised form returned by ReturnNormalisedComponentName.
// If sub components exist the name will be changed to the list of sub components with concentrations.
// e.g. a solution called LB with a concentration of 10X and components 10g/L Yeast Extract and 5g/L Tryptone will be normalised to 10g/L Yeast Extract + 5g/L Tryptone.
// An LB solution with concentration 1 X and no components is returned as 1X LB.
// An LB solution with no concentration and no components is returned as LB.
func NormaliseComponentName(component *Liquid) error {

	newName := ReturnNormalisedComponentName(component)

	component.SetName(newName)

	return nil
}
