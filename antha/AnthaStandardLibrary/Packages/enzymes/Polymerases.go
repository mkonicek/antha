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

// Package for working with enzymes; in particular restriction enzymes
package enzymes

import (
	"fmt"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"sort"
	"strings"
)

type invalidProperty struct{}

var anInvalidProperty = &invalidProperty{}

func (ip *invalidProperty) Error() string {
	var props []string
	for key, _ := range dnaPolymeraseProperties {
		props = append(props, key)
	}
	sort.Strings(props)
	return "Valid options are: " + strings.Join(props, ",")
}

var (
	//

	dnaPolymeraseProperties = map[string]map[string]float64{
		"Q5Polymerase": map[string]float64{
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
		"Taq": map[string]float64{
			"activity_U":          1.0,
			"SecperKb_upper":      90,
			"SecperKb_lower":      60,
			"KBperSecuncertainty": 0.01,
			"Fidelity":            0.0000001,
			"stockconc":           0.01,
			"workingconc":         0.0005,
		},
	}

	DNApolymerasetemps = map[string]map[string]wunit.Temperature{
		"Q5Polymerase": map[string]wunit.Temperature{
			"extensiontemp": wunit.NewTemperature(72, "C"),
			"meltingtemp":   wunit.NewTemperature(98, "C"),
		},
		"Taq": map[string]wunit.Temperature{
			"extensiontemp": wunit.NewTemperature(68, "C"),
			"meltingtemp":   wunit.NewTemperature(95, "C"),
		},
	}
)

func CalculateExtensionTime(polymerase *wtype.LHComponent, targetSequence wtype.DNASequence) (wunit.Time, error) {

	polymerasename := polymerase.CName

	polymeraseproperties, polymerasefound := dnaPolymeraseProperties[polymerasename]

	if !polymerasefound {

		return wunit.Time{}, anInvalidProperty
	}

	sperkblower, found := polymeraseproperties["SperKb_lower"]
	if !found {
		return wunit.Time{}, fmt.Errorf("no property, SperKb_lower found for", polymerase.CName)
	}

	return wunit.NewTime(float64(len(targetSequence.Sequence()))/sperkblower, "s"), anInvalidProperty
}

/*
type assayconds struct {
	Buffer buffermixture
	Temp   wunit.Temperature
}*/

/*
type struct {
	processmodel[]
	processfactors[]
	processfactorcoefficients[]
	Ffactor
	pvalue
}
type buffermixture struct{
	25 mM TAPS-HCl
	(pH 9.3 @ 25°C),
	50 mM KCl,
	2 mM MgCl2,
	1 mM β-mercaptoethanol,
	200 μM dNTPs including [3H]-dTTP and
	400 μg/ml activated Calf Thymus DNA.
	}
*/

type Polymerase struct {
	*wtype.LHComponent
	Uperml              float64
	Rate_sperBP         float64
	Fidelity_errorrate  float64 // could dictate how many colonies are checked in validation!
	Extensiontemp       wunit.Temperature
	Hotstart            bool
	StockConcentration  wunit.Concentration // this is normally in U?
	TargetConcentration wunit.Concentration
	heatinactivation    bool

	// this is also a glycerol solution rather than a watersolution!
} /*
func makePolymerasestatsLibrary() map[string]*Polymerase {

	cmap := make(map[string]*Polymerase)

	poly := NewPolymerase()
	A.CName = "water"
	A.Type = "water"
	A.Smax = 9999
	cmap[A.CName] = A
}

func NewPolymerase() *Polymerase {
	var poly Polymerase
	poly.Rate_sperBP = &poly.Fidelity_errorrate // could dictate how many colonies are checked in validation!
	poly.Extensiontemp
	poly.Hotstart
	poly.StockConcentration // this is normally in U?
	poly.TargetConcentration
	return &poly
}
*/
