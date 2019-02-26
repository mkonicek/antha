// antha/AnthaStandardLibrary/Packages/enzymes/plates.go: Part of the Antha language
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

// Package search is a utility package providing functions useful for:
// Searching for a target entry in a slice;
// Removing duplicate values from a slice;
// Comparing the Name of two entries of any type with a Name() method returning a string.
// FindAll instances of a target string within a template string.
package search

import (
	"fmt"

	"github.com/antha-lang/antha/antha/anthalib/wtype"
)

// SkipAlternateWells is an option which can be used in the NextFreeWell function
// to skip each other well and once the end of the plate is reached go back to the
// beginning and fill the skipped wells.
// This is designed to facilitate multichannel when using a 384 well plate and
// a fixed 8 channel pipette head.
//
const SkipAlternateWells Option = "SkipAlternateWells"

// NextFreeWell checks for the next well which is empty in a plate.
// The user can also specify wells to avoid, preffered wells and
// whether to search through the well positions by row. The default is by column.
// If the SkipAlternateWells option is used we'll skip every other well.
// This is designed to support multichannelling when a plate has 16 rows (i.e. 384 well plate)
// when using a fixed 8 channel pipette head.
//
func NextFreeWell(plate *wtype.Plate, avoidWells []string, preferredWells []string, byRow bool, options ...Option) (well string, err error) {

	if plate == nil {
		return "", fmt.Errorf("no plate specified as argument to NextFreeWell function")
	}

	if len(preferredWells) > 0 {
		for _, well := range preferredWells {
			if err := checkWellValidity(plate, well); err != nil {
				return "", err
			}
			if plate.WellMap()[well].IsEmpty() && !InStrings(avoidWells, well) {
				return well, nil
			}
		}
	}

	allWellPositions := plate.AllWellPositions(byRow)

	if containsOption(options, SkipAlternateWells) {
		newWellPositions := make([]string, 0)
		// odd numbers first
		for index := range allWellPositions {
			if index%2 == 0 {
				newWellPositions = append(newWellPositions, allWellPositions[index])
			}
		}
		// then even numbers
		for index := range allWellPositions {
			if index%2 != 0 {
				newWellPositions = append(newWellPositions, allWellPositions[index])
			}
		}
		allWellPositions = newWellPositions
	}

	for _, well := range allWellPositions {
		// If a well position is found to already have been used then add one to our counter that specifies the next well to use. See step 2 of the following comments.
		if plate.WellMap()[well].IsEmpty() && !InStrings(avoidWells, well) {
			return well, nil
		}
	}
	return "", fmt.Errorf("no empty wells on plate %s of type %s", plate.Name(), plate.Type)
}

// InvalidWell is an error type for when an well is requested from a plate which is invalid.
type InvalidWell string

// Error returns an error message.
func (err InvalidWell) Error() string {
	return string(err)
}

func checkWellValidity(plate *wtype.Plate, well string) error {

	if well != "" {
		wc := wtype.MakeWellCoords(well)
		if wc.X >= len(plate.Cols) {
			return InvalidWell(fmt.Sprintf("well (%s) specified is out of range of available wells for plate type %s", well, plate.Type))
		}
		if wc.Y >= len(plate.Cols[wc.X]) {
			return InvalidWell(fmt.Sprintf("well (%s) specified is out of range of available wells for plate type %s", well, plate.Type))
		}

	}
	return nil
}

// IsFreeWell checks for whether a well on a plate is free.
// An error is returned if the well is not found on the plate or is occupied.
func IsFreeWell(plate *wtype.Plate, well string) error {
	if err := checkWellValidity(plate, well); err != nil {
		return err
	}
	if plate.WellMap()[well].IsEmpty() {
		return nil
	}
	return fmt.Errorf("well %s not free on plate %s %s. Contains %s", well, plate.Name(), plate.Type, plate.WellMap()[well].WContents.Name())
}
