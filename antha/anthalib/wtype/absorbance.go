// wunit/dimensionset.go: Part of the Antha language
// Copyright (C) 2014 the Antha authors. All rights reserved.
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
// 1 Royal College St, London NW1 0NH UK

package wtype

import (
	"github.com/antha-lang/antha/antha/anthalib/wunit"
)

type Absorbance struct {
	WellLocation WellCoords
	Reading      float64
	Wavelength   float64
	Pathlength   wunit.Length
	Status       []string
	Reader       string
}

type Reading interface {
	BlankCorrect(blank Absorbance)
	PathlengthCorrect(pathlength wunit.Length)
	NormaliseTo(target Absorbance)
	CorrecttoRefStandard()
}

func (sample *Absorbance) BlankCorrect(blank Absorbance) {
	if sample.Wavelength == blank.Wavelength &&
		sample.Pathlength == blank.Pathlength &&
		sample.Reader == blank.Reader {
		sample.Reading = sample.Reading - blank.Reading

		sample.Status = append(sample.Status, "Blank Corrected")
	}
}

func (sample *Absorbance) PathlengthCorrect(pathlength wunit.Length) {

	referencepathlength := wunit.NewLength(0.01, "m")

	sample.Reading = sample.Reading * referencepathlength.SIValue() / pathlength.SIValue()
}
