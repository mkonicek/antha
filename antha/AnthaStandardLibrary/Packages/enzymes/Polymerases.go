// antha/AnthaStandardLibrary/Packages/enzymes/Polymerases.go: Part of the Antha language
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

package enzymes

import (
	"fmt"
	"sort"
	"strings"

	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
)

type invalidProperty struct{}

var anInvalidProperty = &invalidProperty{}

func (ip *invalidProperty) Error() string {
	var props []string
	for key := range dnaPolymeraseProperties {
		props = append(props, key)
	}
	sort.Strings(props)
	return "Valid options are: " + strings.Join(props, ",")
}

var (
	//

	dnaPolymeraseProperties = map[string]map[string]float64{
		"Q5Polymerase": {
			"activity_U/ml_assayconds": 50.0,
			"SecperKb_upper":           30,
			"SperKb_lower":             20,
			"KBperSecuncertainty":      0.01,
			"Fidelity":                 0.000000001,
			"stockconc":                0.01,
			"workingconc":              0.0005,
			"extensiontemp":            72.0,
			"meltingtemp":              98.0,
		},
		"Taq": {
			"activity_U":          1.0,
			"SecperKb_upper":      90,
			"SecperKb_lower":      60,
			"KBperSecuncertainty": 0.01,
			"Fidelity":            0.0000001,
			"stockconc":           0.01,
			"workingconc":         0.0005,
		},
	}
	// DNApolymerasetemps contains example cycling properties for two common DNA polymerase enzymes
	DNApolymerasetemps = map[string]map[string]wunit.Temperature{
		"Q5Polymerase": {
			"extensiontemp": wunit.NewTemperature(72, "C"),
			"meltingtemp":   wunit.NewTemperature(98, "C"),
		},
		"Taq": {
			"extensiontemp": wunit.NewTemperature(68, "C"),
			"meltingtemp":   wunit.NewTemperature(95, "C"),
		},
	}
)

// CalculateExtensionTime returns the calculated extension time to amplify a targetSequence with a specified polymerase.
// An error will be returned if the required properties cannot be found for the polymerase.
// Currently the standard valid polymerase options are Taq and Q5Polymerase.
func CalculateExtensionTime(polymeraseName string, targetSequence wtype.DNASequence) (wunit.Time, error) {

	polymeraseproperties, polymerasefound := dnaPolymeraseProperties[polymeraseName]

	if !polymerasefound {

		return wunit.Time{}, anInvalidProperty
	}

	sperkblower, found := polymeraseproperties["SperKb_lower"]
	if !found {
		return wunit.Time{}, fmt.Errorf("no property, SperKb_lower found for %s", polymeraseName)
	}

	return wunit.NewTime(float64(len(targetSequence.Sequence()))/sperkblower, "s"), nil
}
