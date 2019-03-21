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

// Package search is a utility package providing functions useful for:
// Searching for a target entry in a slice;
// Removing duplicate values from a slice;
// Comparing the Name of two entries of any type with a Name() method returning a string.
// FindAll instances of a target string within a template string.
package search

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/antha-lang/antha/antha/anthalib/wtype"
)

// RemoveDuplicateStrings removes all duplicate values in a slice of strings.
// If IgnoreCase is included as an option, case insensitive comparison will be used.
// Any leading or traling space is always removed before comparing elements.
func RemoveDuplicateStrings(elements []string, options ...Option) []string {
	// Use map to record duplicates as we find them.
	encountered := map[string]bool{}
	result := []string{}

	for v := range elements {

		var key string

		if containsIgnoreCase(options...) {
			key = strings.ToUpper(strings.TrimSpace(elements[v]))
		} else {
			key = strings.TrimSpace(elements[v])
		}

		if encountered[key] {
			// Do not add duplicate.
		} else {
			// Record this element as an encountered element.
			encountered[key] = true
			// Append to result slice.
			result = append(result, elements[v])
		}
	}
	// Return the new slice.
	return result
}

// RemoveDuplicateInts removes all duplicate values in a slice of ints.
func RemoveDuplicateInts(elements []int) []int {
	// Use map to record duplicates as we find them.
	encountered := map[int]bool{}
	result := []int{}

	for v := range elements {
		if encountered[elements[v]] {
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

// RemoveDuplicateFloats removes all duplicate values encountered in a slice of floats.
// No precision is specified to floats must be exact matches in order to be considered duplicates.
func RemoveDuplicateFloats(elements []float64) []float64 {
	// Use map to record duplicates as we find them.
	encountered := map[float64]bool{}
	result := []float64{}

	for v := range elements {
		if encountered[elements[v]] {
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

// RemoveDuplicateSequences checks for sequence and name matches only
// and will remove any duplicates based on only these criteria.
// To remove duplicates based on absolute equality use the RemoveDuplicateValues function.
//
//
// If ignoreName option is added, only duplicate sequences with duplicate plasmid status will be removed.
// If ignoreSequences option is added, sequences with duplicate names will be removed.
// If IgnoreCase is included as an option, case insensitive name comparison will be used.
// Sequence comparison will always be case insensitive
// Any leading or trailing space is always removed before comparing elements.
// Ignores any sequence annotations, overhang information and distinction between linear and plasmid DNA.
// To require matches to be strict and identical (annotations, overhand info and all)
// use the ExactMatch option. If ExactMatch is specified, all other Options
// are ignored.
func RemoveDuplicateSequences(elements []wtype.DNASequence, options ...Option) []wtype.DNASequence {

	// Exact Match takes priority and overrides other options.
	if containsExactMatch(options...) {

		// must cast to slice of interface{}
		var elems = make([]interface{}, len(elements))

		for i := range elems {
			elems[i] = interface{}(elements[i])
		}

		noDuplicates, err := RemoveDuplicateValues(elems)

		if err != nil {
			panic(err)
		}

		noDuplicateSeqs := make([]wtype.DNASequence, len(noDuplicates))

		for i := range noDuplicateSeqs {
			noDuplicateSeqs[i] = noDuplicates[i].(wtype.DNASequence)
		}
		return noDuplicateSeqs
	}

	// Use map to record duplicates as we find them.
	encountered := map[string]bool{}
	result := []wtype.DNASequence{}

	for v := range elements {

		var query string

		if !containsIgnoreName(options...) {
			if containsIgnoreCase(options...) {
				query = strings.ToUpper(strings.TrimSpace(elements[v].Name()))
			} else {
				query = strings.TrimSpace(elements[v].Name())
			}
		}

		if !containsIgnoreSequence(options...) {
			if containsIgnoreCase(options...) {
				query += "_" + strings.ToUpper(strings.TrimSpace(elements[v].Sequence()))
			} else {
				query += "_" + strings.TrimSpace(elements[v].Sequence())
			}
			query += "_" + strconv.FormatBool(elements[v].Plasmid)
		}

		if encountered[query] {
			// Do not add duplicate.
		} else {
			// Record this element as an encountered element.
			encountered[query] = true
			// Append to result slice.
			result = append(result, elements[v])
		}
	}
	// Return the new slice.
	return result
}

func removeDuplicateInterface(sliceOfElements []interface{}) []interface{} {

	length := len(sliceOfElements) - 1

	for i := 0; i < length; i++ {
		for j := i + 1; j <= length; j++ {
			if reflect.DeepEqual(sliceOfElements[i], sliceOfElements[j]) {
				sliceOfElements[j] = sliceOfElements[length]
				sliceOfElements = sliceOfElements[0:length]
				length--
				j--
			}
		}
	}

	return sliceOfElements
}

// CheckArrayType reflects the type name as a string of all elements in a slice of interface values.
// If all values in the slice are not of the same type an error is returned.
func CheckArrayType(elements []interface{}) (typeName string, err error) {
	var foundthistype string
	for i, element := range elements {
		typeName = reflect.TypeOf(element).Name()

		if typeName != foundthistype && i != 0 {
			return "mixed types", fmt.Errorf("found different types in []interface{} at %s: %s and %s ", element, typeName, foundthistype)
		}

		foundthistype = typeName
	}

	return
}

// RemoveDuplicateValues removes duplicate values of an input slice of interface values and returns
// all unique entries in slice, preserving the order of the original slice
// an error will be returned if length of elements is 0
func RemoveDuplicateValues(elements []interface{}) ([]interface{}, error) {

	var unique []interface{}

	if len(elements) == 0 {
		return unique, fmt.Errorf("No entries in slice! ")
	}

	t, err := CheckArrayType(elements)

	if err != nil {
		return unique, err
	}

	if t == "int" {
		var intelements []int
		for _, element := range elements {
			intelements = append(intelements, element.(int))
		}

		intelements = RemoveDuplicateInts(intelements)

		for j := range intelements {
			u := intelements[j]
			unique = append(unique, u)
		}
		return unique, nil
	} else if t == "float64" {
		var values []float64
		for _, element := range elements {
			values = append(values, element.(float64))
		}
		values = RemoveDuplicateFloats(values)
		for j := range values {
			u := values[j]
			unique = append(unique, u)
		}
		return unique, nil
	} else if t == "string" {
		var values []string
		for _, element := range elements {
			values = append(values, element.(string))
		}
		values = RemoveDuplicateStrings(values)

		for j := range values {
			u := values[j]
			unique = append(unique, u)
		}
		return unique, nil
	}

	unique = removeDuplicateInterface(elements)
	if len(unique) == 0 {
		return unique, fmt.Errorf("[]interface{} conversion has gone wrong!, length of output differs to input: %d and %d: Array type: %s", len(unique), len(elements), t)
	}

	// Return the new slice.
	return unique, nil
}
