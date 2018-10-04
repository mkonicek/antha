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
	"math"
	"strings"

	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/pkg/errors"
)

// CorrectionType is a label given to describe the nature of the correction of an absorbance reading.
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

// ReferencePathlength is 10mm and is the pathlength used to normalise an absorbance measurement to.
// The value is the pathlength of a standard cuvette
var ReferencePathlength = wunit.NewLength(10, "mm")

// AbsorbanceCorrection stores the details of how an Absorbance reading has ben corrected.
type AbsorbanceCorrection struct {
	Type              CorrectionType
	CorrectionReading *Absorbance
}

// Absorbance stores the key properties of an absorbance reading.
type Absorbance struct {
	WellLocation WellCoords             `json:"WellCoords"`
	Reading      float64                `json:"Reading"`
	Wavelength   wunit.Length           `json:"Wavelength"`
	Pathlength   wunit.Length           `json:"Pathlength"`
	Corrections  []AbsorbanceCorrection `json:"Corrections"`
	Reader       string                 `json:"Reader"`
	ID           string                 `json:"ID"`
	// Annotations is a field to add custom user labels
	Annotations []string `json:"Annotations"`
}

// WavelengthToNearestNm will return the Wavelength field as an int.
// Whilst it is possible that the wavelength used may be a decimal,
// Wavelength would typically be expected to be in the form of an integer of the wavelength in nm.
// In some platereader data sets this is stored as a float so this method is
// intended to take the safest representation, as a float, and return the more
// common representation, as an int, for parsers where it is known that the wavelength
// is stored as an int.
// This method would therefore not be safe to use for situations
// where the wavelenWavelengthToNearestNm be represented by a decimal.
//
func (a Absorbance) WavelengthToNearestNm() int {
	if uint64(toNM(a.Wavelength)) > math.MaxUint64 {
		panic(errors.Errorf("the value for wavelength %v cannot be safely converted to an integer value", a.Wavelength))
	}
	return int(toNM(a.Wavelength))
}

// IsBlankCorrected returns true if the absorbance reading has been blank corrected
func (a Absorbance) IsBlankCorrected() bool {
	for _, correction := range a.Corrections {
		if correction.Type == BlankCorrected {
			return true
		}
	}
	return false
}

// IsPathLengthCorrected returns true if the absorbance reading has been pathlength corrected
func (a Absorbance) IsPathLengthCorrected() bool {
	for _, correction := range a.Corrections {
		if correction.Type == PathLengthCorrected {
			return true
		}
	}
	return false
}

func toNM(l wunit.Length) float64 {
	return l.SIValue() * wunit.Nano.Value
}

type Reading interface {
	BlankCorrect(blank Absorbance) error
	PathLengthCorrect(pathlength wunit.Length)
	NormaliseTo(target Absorbance)
	CorrecttoRefStandard()
}

// Dup creates a duplicate of the absorbance reading, with exact equality for all values.
func (sample *Absorbance) Dup() Absorbance {
	return Absorbance{
		WellLocation: sample.WellLocation,
		Reading:      sample.Reading,
		Wavelength:   sample.Wavelength,
		Pathlength:   sample.Pathlength,
		Corrections:  append([]AbsorbanceCorrection{}, sample.Corrections...),
		Reader:       sample.Reader,
		ID:           sample.ID,
		Annotations:  append([]string{}, sample.Annotations...),
	}
}

// BlankCorrect subtracts the blank reading from the sample absorbance.
// If the blank sample is not equivalent to the sample, based on wavelength and pathlength, an error is returned.
func (sample *Absorbance) BlankCorrect(blanks ...Absorbance) error {
	var errs []string
	for _, blank := range blanks {
		if sample.Wavelength.EqualToRounded(blank.Wavelength, 9); sample.Pathlength.EqualToRounded(blank.Pathlength, 4) && sample.Reader == blank.Reader {
			sample.Reading = sample.Reading - blank.Reading
			sample.Corrections = append(sample.Corrections,
				AbsorbanceCorrection{
					Type:              BlankCorrected,
					CorrectionReading: &blank,
				},
			)
		} else {
			errs = append(errs,
				fmt.Sprintf(
					`cannot pathlength correct as Absorbance readings for 
			sample (%+v) and blank (%+v) are incompatible due to 
			either wavelength, pathlength or reader differences.`,
					sample,
					blank,
				),
			)
		}
	}

	if len(errs) > 0 {
		return errors.Errorf(strings.Join(errs, ";"))
	}
	return nil
}

// PathLengthCorrect normalises an absorbance reading
// to a standard reference pathlength of 1cm.
// 1cm is the pathlength used to normalise absorbance readings to OD.
func (sample *Absorbance) PathLengthCorrect(pathlength wunit.Length) error {

	if sample.IsPathLengthCorrected() {
		return errors.Errorf("absorbance sample %+v has already been pathlength corrected", sample)
	}

	sample.Reading = sample.Reading * ReferencePathlength.RawValue() / pathlength.RawValue()

	sample.Corrections = append(sample.Corrections, AbsorbanceCorrection{Type: PathLengthCorrected, CorrectionReading: nil})
	return nil
}
