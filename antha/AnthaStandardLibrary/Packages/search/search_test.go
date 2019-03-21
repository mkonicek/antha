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
	"testing"

	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/text"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/go-test/deep"
)

type searchIntsTest struct {
	slice          []int
	value          int
	expectedResult bool
}

func makeMassiveSlice(size int) []int {
	var massiveList []int
	for i := 0; i < size; i++ {
		massiveList = append(massiveList, i)
	}
	return massiveList
}

var intTests = []searchIntsTest{
	{
		slice:          []int{1, 2, 3, 4},
		value:          4,
		expectedResult: true,
	},
	{
		slice:          makeMassiveSlice(100000),
		value:          9999,
		expectedResult: true,
	},
}

type searchStringsTest struct {
	slice          []string
	value          string
	ignoreCase     Option
	expectedResult bool
}

var stringTests = []searchStringsTest{
	{
		slice:          []string{"a", "b"},
		value:          "a",
		ignoreCase:     IgnoreCase,
		expectedResult: true,
	},
	{
		slice:          []string{"a", "b"},
		value:          "A",
		ignoreCase:     IgnoreCase,
		expectedResult: true,
	},
	{
		slice:          []string{"a", "b"},
		value:          "A ",
		ignoreCase:     IgnoreCase,
		expectedResult: true,
	},
	{
		slice:          []string{"a", "b"},
		value:          "A",
		ignoreCase:     "",
		expectedResult: false,
	},
	{
		slice:          []string{},
		value:          "A ",
		ignoreCase:     IgnoreCase,
		expectedResult: false,
	},
}

// InInts searchs for an int in a slice of ints
func TestInInts(t *testing.T) {
	for _, test := range intTests {
		present := InInts(test.slice, test.value)

		if present != test.expectedResult {
			t.Error(
				"For", test.slice, test.value, "\n",
				"expected:", test.expectedResult, "\n",
				"got", present, "\n",
			)
		}
	}
}

func TestBinarySearch(t *testing.T) {
	for _, test := range intTests {
		present := BinarySearchInts(test.slice, test.value)

		if present != test.expectedResult {
			t.Error(
				"For", test.slice, test.value, "\n",
				"expected:", test.expectedResult, "\n",
				"got", present, "\n",
			)
		}
	}
}

type seqSearchTest struct {
	slice           []wtype.DNASequence
	value           wtype.DNASequence
	expectedResult  bool
	ignoreSequences Option
	ignoreCase      Option
}

var seqTests = []seqSearchTest{
	{
		slice: []wtype.DNASequence{
			{Nm: "Bob", Seq: "AACCACACTT"},
		},
		value:          wtype.DNASequence{Nm: "bob", Seq: "AACCACACTT"},
		expectedResult: true,
		ignoreCase:     IgnoreCase,
	},
	{
		slice: []wtype.DNASequence{
			{Nm: "Bob", Seq: "AACCACACTT"},
		},
		value:           wtype.DNASequence{Nm: "bob", Seq: "AACCACACTT"},
		expectedResult:  false,
		ignoreSequences: IgnoreSequence,
		ignoreCase:      "",
	},
}

func TestNamed(t *testing.T) {
	for _, test := range seqTests {

		for _, entry := range test.slice {
			present := EqualName(&entry, &test.value, test.ignoreCase)

			if present != test.expectedResult {
				t.Error(
					"For", test.slice, test.value, "\n",
					"expected:", test.expectedResult, "\n",
					"got", present, "\n",
				)
			}
		}
	}
}

func TestInSequences(t *testing.T) {
	for _, test := range seqTests {

		present, positions := InSequences(test.slice, test.value, test.ignoreSequences, test.ignoreCase)

		if present != test.expectedResult {
			t.Error(
				"For", test.slice, test.value, "\n",
				"expected:", test.expectedResult, "\n",
				"got", present, positions, "\n",
			)
		}

	}
}

func TestInStrings(t *testing.T) {
	for _, test := range stringTests {

		present := InStrings(test.slice, test.value, test.ignoreCase)

		if present != test.expectedResult {
			t.Error(
				"For", test.slice, test.value, "\n",
				"expected:", test.expectedResult, "\n",
				"got", present, "\n",
			)
		}

	}
}

type removalTest struct {
	Values         []interface{}
	ExpectedResult []interface{}
	ExpectedErr    bool
}

func TestRemoveDuplicateValues(t *testing.T) {
	tests := []removalTest{
		{
			Values: []interface{}{
				&wtype.DNASequence{
					Nm:      "GATCGTAGTGT",
					Seq:     "GATCGTAGTGT",
					Plasmid: true,
				},
				&wtype.DNASequence{
					Nm:      "GATCGTAGTGT",
					Seq:     "GATCGTAGTGT",
					Plasmid: true,
				},
			},
			ExpectedResult: []interface{}{
				&wtype.DNASequence{
					Nm:      "GATCGTAGTGT",
					Seq:     "GATCGTAGTGT",
					Plasmid: true,
				},
			},
		},
		{
			Values: []interface{}{
				&wtype.DNASequence{
					Nm:             "GATCGTAGTGT",
					Seq:            "GATCGTAGTGT",
					Plasmid:        true,
					Singlestranded: false,
				},
				&wtype.DNASequence{
					Nm:             "GATCGTAGTGT",
					Seq:            "GATCGTAGTGT",
					Plasmid:        true,
					Singlestranded: true,
				},
			},
			ExpectedResult: []interface{}{
				&wtype.DNASequence{
					Nm:             "GATCGTAGTGT",
					Seq:            "GATCGTAGTGT",
					Plasmid:        true,
					Singlestranded: false,
				},
				&wtype.DNASequence{
					Nm:             "GATCGTAGTGT",
					Seq:            "GATCGTAGTGT",
					Plasmid:        true,
					Singlestranded: true,
				},
			},
		},
	}

	for _, test := range tests {
		result, err := RemoveDuplicateValues(test.Values)

		if (err != nil) && test.ExpectedErr {
			t.Error(
				"Unexpected test result for ", text.PrettyPrint(test), "\n",
				"Got error: ", err, "\n ",
			)
		}

		if diffs := deep.Equal(test.ExpectedResult, result); len(diffs) > 0 {
			t.Error(
				"Unexpected test result for ", text.PrettyPrint(test), "\n",
				"Differences detected: ", diffs, "\n ",
			)
		}
	}
}
