// antha/AnthaStandardLibrary/Packages/sequences/findSeq_test.go: Part of the Antha language
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

// Package sequences is for interacting with and manipulating biological sequences; in extension to methods available in wtype
package sequences

import (
	"fmt"
	"strings"
	"testing"

	"github.com/antha-lang/antha/antha/anthalib/wtype"
)

type translateTest struct {
	AminoAcids      *wtype.ProteinSequence
	DNA             *wtype.DNASequence
	CodonUsageTable CodonUsageTable
}

func makeLSequence(size int) string {
	var lll []string
	lll = append(lll, "M")
	for i := 0; i < size; i++ {
		lll = append(lll, "L")
	}
	lll = append(lll, "*")
	return strings.Join(lll, "")
}

var revTranslateTests = []translateTest{

	{
		AminoAcids: &wtype.ProteinSequence{
			Nm:  "Test1",
			Seq: makeLSequence(10),
		},
		DNA: &wtype.DNASequence{
			Nm:  "Test1",
			Seq: "ATGCTCCTTCTTCTTTTACTGCTCCTACTGTTATGA",
		},
		CodonUsageTable: EColiTable,
	},
}

func TestRevTranslate(t *testing.T) {
	for _, test := range revTranslateTests {
		dnaSeq, err := RevTranslate(*test.AminoAcids, test.CodonUsageTable)
		if dnaSeq.Sequence() != test.DNA.Sequence() {
			t.Error(
				"For", test, "\n",
				"expected ", fmt.Sprintf("%v", test.DNA), "\n",
				"got", fmt.Sprintf("%v", dnaSeq), "\n",
			)
		}
		if err != nil {
			t.Error(
				"For", test, "\n",
				"got error: ", err.Error(), "\n",
			)
		}
		// now translate back
		protSeq, err := Translate(*test.DNA)
		if protSeq.Sequence() != test.AminoAcids.Sequence() {
			t.Error(
				"For", test, "\n",
				"expected ", fmt.Sprintf("%v", test.AminoAcids), "\n",
				"got", fmt.Sprintf("%v", protSeq), "\n",
			)
		}
		if err != nil {
			t.Error(
				"For", test, "\n",
				"got error: ", err.Error(), "\n",
			)
		}
	}

}
