// antha/AnthaStandardLibrary/Packages/Platereader/Platereader.go: Part of the Antha language
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

//Package platereader contains functions for manipulating absorbance readings and platereader data.
package platereader

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
)

// BlankCorrect subtracts the blank reading from the sample absorbance.
// If the blank sample is not equivalent to the sample, based on wavelength and pathlength, an error is returned.
func Blankcorrect(sample *wtype.Absorbance, blank *wtype.Absorbance) (blankcorrected *wtype.Absorbance, err error) {

	blankcorrected = sample.Dup()

	err = blankcorrected.BlankCorrect(blank)

	return
}

func EstimatePathLength(plate *wtype.Plate, volume wunit.Volume) (pathlength wunit.Length, err error) {

	if plate.Welltype.Bottom == 0 /* i.e. flat */ && plate.Welltype.Shape().LengthUnit == "mm" {
		wellarea, err := plate.Welltype.CalculateMaxCrossSectionArea()
		if err != nil {

			return pathlength, err
		}
		wellvol, err := plate.Welltype.CalculateMaxVolume()
		if err != nil {
			return pathlength, err
		}

		if volume.Unit().PrefixedSymbol() == "ul" && wellvol.Unit().PrefixedSymbol() == "ul" && wellarea.Unit().PrefixedSymbol() == "mm^2" || wellarea.Unit().PrefixedSymbol() == "mm" /* mm generated previously - wrong and needs fixing */ {
			ratio := volume.RawValue() / wellvol.RawValue()

			wellheightinmm := wellvol.RawValue() / wellarea.RawValue()

			pathlengthinmm := wellheightinmm * ratio

			pathlength = wunit.NewLength(pathlengthinmm, "mm")

		} else {
			fmt.Println(volume.Unit().PrefixedSymbol(), wellvol.Unit().PrefixedSymbol(), wellarea.Unit().PrefixedSymbol(), wellarea.ToString())
		}
	} else {
		err = errors.New(fmt.Sprint("Can't yet estimate pathlength for this welltype shape unit ", plate.Welltype.Shape().LengthUnit, "or non flat bottom type"))
	}

	return
}

// PathLengthCorrect normalises an absorbance reading
// to a standard reference pathlength of 1cm.
// 1cm is the pathlength used to normalise absorbance readings to OD.
func PathlengthCorrect(pathlength wunit.Length, reading *wtype.Absorbance) (pathlengthcorrected *wtype.Absorbance) {
	pathlengthcorrected = reading.Dup()
	err := pathlengthcorrected.PathLengthCorrect(pathlength)
	if err != nil {
		panic(err)
	}
	return pathlengthcorrected
}

// based on Beer Lambert law A = ε l c
/*
Limitations of the Beer-Lambert law

The linearity of the Beer-Lambert law is limited by chemical and instrumental factors. Causes of nonlinearity include:
deviations in absorptivity coefficients at high concentrations (>0.01M) due to electrostatic interactions between molecules in close proximity
scattering of light due to particulates in the sample
fluoresecence or phosphorescence of the sample
changes in refractive index at high analyte concentration
shifts in chemical equilibria as a function of concentration
non-monochromatic radiation, deviations can be minimized by using a relatively flat part of the absorption spectrum such as the maximum of an absorption band
stray light
*/
func Concentration(pathlengthcorrected wtype.Absorbance, molarabsorbtivityatwavelengthLpermolpercm float64) (conc wunit.Concentration, err error) {

	if !pathlengthcorrected.IsPathLengthCorrected() {
		return wunit.Concentration{}, errors.Errorf("absorbance reading (%+v) has not been pathlength corrected, please use PathlengthCorrect method on the Absorbance value.", pathlengthcorrected)
	}

	A := pathlengthcorrected
	l := 1                                         // 1cm if pathlengthcorrected add logic to use pathlength of absorbance reading input
	ε := molarabsorbtivityatwavelengthLpermolpercm // l/Mol/cm

	concfloat := A.Reading / (float64(l) * ε)
	return wunit.NewConcentration(concfloat, "Mol/l"), nil
}
