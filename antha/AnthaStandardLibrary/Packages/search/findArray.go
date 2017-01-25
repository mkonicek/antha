// antha/AnthaStandardLibrary/Packages/enzymes/Find.go: Part of the Antha language
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

// Utility package providing functions useful for searches
package search

import (
	"strings"

	"github.com/antha-lang/antha/antha/anthalib/wtype"
)

func InSlice(slice string, list []string) bool {
	for _, b := range list {
		if b == slice {
			return true
		}
	}
	return false
}

func Position(slice []string, value string) int {
	for p, v := range slice {
		if v == value {
			return p
		}
	}
	return -1
}

func RemoveDuplicates(elements []string) []string {
	// Use map to record duplicates as we find them.
	encountered := map[string]bool{}
	result := []string{}

	for v := range elements {
		if encountered[elements[v]] == true {
			// Do not add duplicate.
		} else {
			// Record this element as an encountered element.
			encountered[elements[v]] = true
			// Append to result slice.
			result = append(result, elements[v])
		}
	}
	// Return the new slice.
	return result
}

func RemoveDuplicateInts(elements []int) []int {
	// Use map to record duplicates as we find them.
	encountered := map[int]bool{}
	result := []int{}

	for v := range elements {
		if encountered[elements[v]] == true {
			// Do not add duplicate.
		} else {
			// Record this element as an encountered element.
			encountered[elements[v]] = true
			// Append to result slice.
			result = append(result, elements[v])
		}
	}
	// Return the new slice.
	return result
}

func RemoveDuplicateInterface(elements []interface{}) []interface{} {
	// Use map to record duplicates as we find them.
	encountered := map[interface{}]bool{}
	result := []interface{}{}

	for v := range elements {
		if encountered[elements[v]] == true {
			// Do not add duplicate.
		} else {
			// Record this element as an encountered element.
			encountered[elements[v]] = true
			// Append to result slice.
			result = append(result, elements[v])
		}
	}
	// Return the new slice.
	return result
}

func RemoveDuplicatesKeysfromMap(elements map[interface{}]interface{}) map[interface{}]interface{} {
	// Use map to record duplicates as we find them.
	encountered := map[interface{}]bool{}
	result := make(map[interface{}]interface{}, 0)

	for key, v := range elements {

		if encountered[key] == true {
			// Do not add duplicate.
		} else {
			// Record this element as an encountered element.
			encountered[key] = true
			// Append to result slice.
			result[key] = v
		}
	}
	// Return the new slice.
	return result
}

func RemoveDuplicatesValuesfromMap(elements map[interface{}]interface{}) map[interface{}]interface{} {
	// Use map to record duplicates as we find them.
	encountered := map[interface{}]bool{}
	result := make(map[interface{}]interface{}, 0)

	for key, v := range elements {

		if encountered[v] == true {
			// Do not add duplicate.
		} else {
			// Record this element as an encountered element.
			encountered[v] = true
			// Append to result slice.
			result[key] = v
		}
	}
	// Return the new slice.
	return result
}

// based on exact sequence matches only; ignores name
func RemoveDuplicateSequences(elements []wtype.DNASequence) []wtype.DNASequence {
	// Use map to record duplicates as we find them.
	encountered := map[string]bool{}
	result := []wtype.DNASequence{}

	for v := range elements {
		if encountered[strings.ToUpper(elements[v].Seq)] == true {
			// Do not add duplicate.
		} else {
			// Record this element as an encountered element.
			encountered[strings.ToUpper(elements[v].Seq)] = true
			// Append to result slice.
			result = append(result, elements[v])
		}
	}
	// Return the new slice.
	return result
}
