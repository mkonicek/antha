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
package search

// Check the wells that have already been used in this plate to be sure not to pipette on top of a previous solution. This is predominantly used for when two aliquot elements are used in serial to make sure the second set of aliquots are not pipette onto the first set.
func NextFreeWell(allWellPositions []string, wellsUsed []string) (well string, err error) {

	for _, well := range allWellPositions {
		// If a well position is found to already have been used then add one to our counter that specifies the next well to use. See step 2 of the following comments.
		if !InSlice(well, wellsUsed) {
			return well, nil
		}
	}
	return
}
