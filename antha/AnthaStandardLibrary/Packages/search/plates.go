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

// NextFreeWell checks for the next well which is empty in a plate.
// The user can also specify wells to avoid, preffered wells and
// whether to search through the well positions by row. The default is by column.
func NextFreeWell(plate *wtype.Plate, avoidWells []string, preferredWells []string, byRow bool) (well string, err error) {

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

	for _, well := range allWellPositions {
		// If a well position is found to already have been used then add one to our counter that specifies the next well to use. See step 2 of the following comments.
		if plate.WellMap()[well].IsEmpty() && !InStrings(avoidWells, well) {
			return well, nil
		}
	}
	return "", fmt.Errorf("no empty wells on plate %s of type %s", plate.Name(), plate.Type)
}

func checkWellValidity(plate *wtype.Plate, well string) error {

	if well != "" {
		wc := wtype.MakeWellCoords(well)
		if wc.X >= len(plate.Cols) {
			return fmt.Errorf("well (%s) specified is out of range of available wells for plate type %s", well, plate.Type)
		}
		if wc.Y >= len(plate.Cols[wc.X]) {
			return fmt.Errorf("well (%s) specified is out of range of available wells for plate type %s", well, plate.Type)
		}

	}
	return nil
}
