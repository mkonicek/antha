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
	"reflect"
	"testing"

	"github.com/antha-lang/antha/antha/anthalib/wtype"
)

type replaceTest struct {
	LargeSeq       *wtype.DNASequence
	ReplaceSeq     *wtype.DNASequence
	ReplaceWithSeq *wtype.DNASequence
	ExpectedResult *wtype.DNASequence
	ErrMessage     string
}

type rotateTest struct {
	seq            string
	rotateBy       int
	reverse        bool
	expectedResult string
}

var rotateTests = []rotateTest{
	rotateTest{
		seq:            "REVERSE",
		rotateBy:       0,
		reverse:        false,
		expectedResult: "REVERSE",
	},
	rotateTest{
		seq:            "REVERSE",
		rotateBy:       1,
		reverse:        false,
		expectedResult: "EVERSER",
	},
	rotateTest{
		seq:            "REVERSE",
		rotateBy:       1,
		reverse:        true,
		expectedResult: "EREVERS",
	},
}

var replaceTests = []replaceTest{

	replaceTest{
		LargeSeq: &wtype.DNASequence{
			Nm:      "Test1",
			Seq:     "ATCGTAGTGTG",
			Plasmid: false,
		},
		ReplaceSeq: &wtype.DNASequence{
			Nm:      "TAG",
			Seq:     "TAG",
			Plasmid: false,
		},
		ReplaceWithSeq: &wtype.DNASequence{
			Nm:      "AAA",
			Seq:     "AAA",
			Plasmid: false,
		},
		ExpectedResult: &wtype.DNASequence{
			Nm:      "Test1",
			Seq:     "ATCGAAATGTG",
			Plasmid: false,
		},
		ErrMessage: "",
	},
	replaceTest{
		LargeSeq: &wtype.DNASequence{
			Nm:      "Test2", // reverse complement
			Seq:     "ATCGTAGTGTG",
			Plasmid: false,
		},
		ReplaceSeq: &wtype.DNASequence{
			Nm:      "TAC",
			Seq:     "TAC",
			Plasmid: false,
		},
		ReplaceWithSeq: &wtype.DNASequence{
			Nm:      "AAA",
			Seq:     "AAA",
			Plasmid: false,
		},
		ExpectedResult: &wtype.DNASequence{
			Nm:      "Test2",
			Seq:     "ATCAAAGTGTG",
			Plasmid: false,
		},
		ErrMessage: "",
	},
	replaceTest{
		LargeSeq: &wtype.DNASequence{
			Nm:      "Test3", // reverse complement two replace regions
			Seq:     "ATCGTAGTGTGTAC",
			Plasmid: false,
		},
		ReplaceSeq: &wtype.DNASequence{
			Nm:      "TAC",
			Seq:     "TAC",
			Plasmid: false,
		},
		ReplaceWithSeq: &wtype.DNASequence{
			Nm:      "AAA",
			Seq:     "AAA",
			Plasmid: false,
		},
		ExpectedResult: &wtype.DNASequence{
			Nm:      "Test3",
			Seq:     "ATCAAAGTGTGAAA",
			Plasmid: false,
		},
		ErrMessage: "",
	},
	replaceTest{
		LargeSeq: &wtype.DNASequence{
			Nm:      "Test4", // plasmid
			Seq:     "ATCGTAGTGTGTAC",
			Plasmid: true,
		},
		ReplaceSeq: &wtype.DNASequence{
			Nm:      "TAC",
			Seq:     "TAC",
			Plasmid: false,
		},
		ReplaceWithSeq: &wtype.DNASequence{
			Nm:      "AAA",
			Seq:     "AAA",
			Plasmid: false,
		},
		ExpectedResult: &wtype.DNASequence{
			Nm:      "Test4",
			Seq:     "ATCAAAGTGTGAAA",
			Plasmid: true,
		},
		ErrMessage: "",
	},
	replaceTest{
		LargeSeq: &wtype.DNASequence{
			Nm:      "Test5", // plasmid overlapping end of plasmid
			Seq:     "ATCGTAGTGTGTAC",
			Plasmid: true,
		},
		ReplaceSeq: &wtype.DNASequence{
			Nm:      "TACATCG",
			Seq:     "TACATCG",
			Plasmid: false,
		},
		ReplaceWithSeq: &wtype.DNASequence{
			Nm:      "AAAAAAA",
			Seq:     "AAAAAAA",
			Plasmid: false,
		},
		ExpectedResult: &wtype.DNASequence{
			Nm:      "Test5",
			Seq:     "AAAATAGTGTGAAA",
			Plasmid: true,
		},
		ErrMessage: "",
	},
	replaceTest{
		LargeSeq: &wtype.DNASequence{
			Nm:      "Test6", // replace with nothing
			Seq:     "ATCGTAGTGTGTAC",
			Plasmid: true,
		},
		ReplaceSeq: &wtype.DNASequence{
			Nm:      "TACATCG",
			Seq:     "TACATCG",
			Plasmid: false,
		},
		ReplaceWithSeq: &wtype.DNASequence{
			Nm:      "nothing",
			Seq:     "",
			Plasmid: false,
		},
		ExpectedResult: &wtype.DNASequence{
			Nm:      "Test6",
			Seq:     "TGTGTAG",
			Plasmid: true,
		},
		ErrMessage: "",
	},
	replaceTest{
		LargeSeq: &wtype.DNASequence{
			Nm:      "Test7", // replace with nothing
			Seq:     "ATCGTAGTGTGTAC",
			Plasmid: true,
		},
		ReplaceSeq: &wtype.DNASequence{
			Nm:      "GTAGTGTG",
			Seq:     "GTAGTGTG",
			Plasmid: false,
		},
		ReplaceWithSeq: &wtype.DNASequence{
			Nm:      "nothing",
			Seq:     "",
			Plasmid: false,
		},
		ExpectedResult: &wtype.DNASequence{
			Nm:      "Test7",
			Seq:     "ATCTAC",
			Plasmid: true,
		},
		ErrMessage: "",
	},
}

func TestReplace(t *testing.T) {
	for _, test := range replaceTests {
		result, err := ReplaceAll(*test.LargeSeq, *test.ReplaceSeq, *test.ReplaceWithSeq)
		if !reflect.DeepEqual(&result, test.ExpectedResult) {
			t.Error(
				"For", test.LargeSeq.Nm, "\n",
				"replacing ", test.ReplaceSeq.Seq, "\n",
				"with ", test.ReplaceWithSeq.Seq, "\n",
				"expected ", fmt.Sprint(test.ExpectedResult), "\n",
				"got", fmt.Sprint(result), "\n",
			)
		}
		if err != nil {
			if err.Error() != test.ErrMessage {
				t.Error(
					"For", test.LargeSeq.Nm, "\n",
					"replacing ", test.ReplaceSeq.Seq, "\n",
					"with ", test.ReplaceWithSeq.Seq, "\n",
					"got error ", test.ErrMessage, "\n",
					"got", err.Error(), "\n",
				)
			}
		}
	}
}

func TestRotate(t *testing.T) {
	for _, test := range rotateTests {
		result := Rotate(test.seq, test.rotateBy, test.reverse)
		if reflect.DeepEqual(&result, test.expectedResult) {
			t.Error(
				"For", test.seq, "\n",
				"rotating by ", test.rotateBy, "\n",
				"reverse ", test.reverse, "\n",
				"expected ", test.expectedResult, "\n",
				"got", result, "\n",
			)
		}
	}
}

/*
func seqsEqual(seq1, seq2 wtype.DNASequence)bool{

}
*/
