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
	"fmt"

	"github.com/antha-lang/antha/antha/anthalib/wunit"
)

// CorrectionType is a label given to the nature of the correction of an absorbance reading.
type CorrectionType string

const (
	// BlankCorrected refers to when an absorbance reading has been corrected
	// by substracting the absorbance at equivalent conditions of a blank sample.
	BlankCorrected CorrectionType = "Blank Corrected"

	// PathLengthCorrected refers to when an absorbance reading has been normalised
	// to a standard reference pathlength of 1cm.
	// 1cm is the pathlength used to normalise absorbance readings to OD.
	PathLengthCorrected CorrectionType = "Pathlength Corrected"

	// ReferenceStandardCorrected refers to when an absorbance reading is corrected based
	// on a reference sample. Placeholder: Not yet implemented.
	ReferenceStandardCorrected CorrectionType = "Reference Standard Corrected"
)

// AbsorbanceCorrection stores the details of how an Absorbance reading has ben corrected.
type AbsorbanceCorrection struct {
	Type              CorrectionType
	CorrectionReading *Absorbance
}

// Absorbance stores the key properties of an absorbance reading.
type Absorbance struct {
	Reading     float64                `json:"Reading"`
	Wavelength  wunit.Length           `json:"Wavelength"`
	Pathlength  wunit.Length           `json:"Pathlength"`
	Corrections []AbsorbanceCorrection `json:"Corrections"`
	Reader      string                 `json:"Reader"`
	ID          string                 `json:"ID"`
}

type Reading interface {
	BlankCorrect(blank Absorbance) error
	PathLengthCorrect(pathlength wunit.Length)
	NormaliseTo(target Absorbance)
	CorrecttoRefStandard()
}

// BlankCorrect subtracts the blank reading from the sample absorbance.
// If the blank sample is not equivalent to the sample, based on wavelength and pathlength, an error is returned.
func (sample *Absorbance) BlankCorrect(blank *Absorbance) error {

	if sample.Wavelength.EqualToRounded(blank.Wavelength, 9); sample.Pathlength.EqualToRounded(blank.Pathlength, 4) &&
		sample.Reader == blank.Reader {
		sample.Reading = sample.Reading - blank.Reading

		sample.Corrections = append(sample.Corrections, AbsorbanceCorrection{Type: BlankCorrected, CorrectionReading: blank})
		return nil
	}
	return fmt.Errorf("Cannot pathlength correct as Absorbance readings for sample (%+v) and blank (%+v) are incompatible due to either wavelength, pathlength or reader differences. ", sample, blank)
}

// PathLengthCorrect normalises an absorbance reading
// to a standard reference pathlength of 1cm.
// 1cm is the pathlength used to normalise absorbance readings to OD.
func (sample *Absorbance) PathLengthCorrect(pathlength wunit.Length) {

	referencepathlength := wunit.NewLength(10, "mm")

	sample.Reading = sample.Reading * referencepathlength.RawValue() / pathlength.RawValue()

	sample.Corrections = append(sample.Corrections, AbsorbanceCorrection{Type: PathLengthCorrected, CorrectionReading: nil})

	return
}
