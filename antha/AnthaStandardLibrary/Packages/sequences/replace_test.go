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

type replaceAllTest struct {
	LargeSeq       *wtype.DNASequence
	ReplaceSeq     *wtype.DNASequence
	ReplaceWithSeq *wtype.DNASequence
	ExpectedResult *wtype.DNASequence
	ErrMessage     string
}

type replaceTest struct {
	LargeSeq       *wtype.DNASequence
	ReplaceWithSeq *wtype.DNASequence
	Position       PositionPair
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
	{
		seq:            "REVERSE",
		rotateBy:       0,
		reverse:        false,
		expectedResult: "REVERSE",
	},
	{
		seq:            "REVERSE",
		rotateBy:       1,
		reverse:        false,
		expectedResult: "EVERSER",
	},
	{
		seq:            "REVERSE",
		rotateBy:       1,
		reverse:        true,
		expectedResult: "EREVERS",
	},
}

var replaceTests = []replaceAllTest{

	{
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
	{
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
	{
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
	{
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
	{
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
	{
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
	{
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

var replacePositionTests = []replaceTest{
	{
		LargeSeq: &wtype.DNASequence{
			Nm:      "Test1",
			Seq:     "ATCGTAGTGTG",
			Plasmid: false,
		},
		ReplaceWithSeq: &wtype.DNASequence{
			Nm:      "nothing",
			Seq:     "",
			Plasmid: false,
		},
		Position: PositionPair{
			StartPosition: 1,
			EndPosition:   3,
			Reverse:       false,
		},
		ExpectedResult: &wtype.DNASequence{
			Nm:      "Test1",
			Seq:     "GTAGTGTG",
			Plasmid: false,
		},
		ErrMessage: "",
	},
	{
		LargeSeq: &wtype.DNASequence{
			Nm:      "Test2",
			Seq:     "ATCGTAGTGTG",
			Plasmid: false,
		},
		ReplaceWithSeq: &wtype.DNASequence{
			Nm:      "TTT",
			Seq:     "TTT",
			Plasmid: false,
		},
		Position: PositionPair{
			StartPosition: 1,
			EndPosition:   3,
			Reverse:       false,
		},
		ExpectedResult: &wtype.DNASequence{
			Nm:      "Test2",
			Seq:     "TTTGTAGTGTG",
			Plasmid: false,
		},
		ErrMessage: "",
	},
	{
		LargeSeq: &wtype.DNASequence{
			Nm:      "Test3",
			Seq:     "ATCGTAGTGTG",
			Plasmid: false,
		},
		ReplaceWithSeq: &wtype.DNASequence{
			Nm:      "TTT",
			Seq:     "TTTTTT",
			Plasmid: false,
		},
		Position: PositionPair{
			StartPosition: 1,
			EndPosition:   3,
			Reverse:       false,
		},
		ExpectedResult: &wtype.DNASequence{
			Nm:      "Test3",
			Seq:     "TTTTTTGTAGTGTG",
			Plasmid: false,
		},
		ErrMessage: "",
	},
	{
		LargeSeq: &wtype.DNASequence{
			Nm:      "Test4",
			Seq:     "ATCGTAGTGTG",
			Plasmid: false,
		},
		ReplaceWithSeq: &wtype.DNASequence{
			Nm:      "TTT",
			Seq:     "TTTTTT",
			Plasmid: false,
		},
		Position: PositionPair{
			StartPosition: 1,
			EndPosition:   3,
			Reverse:       true,
		},
		ExpectedResult: &wtype.DNASequence{
			Nm:      "Test4",
			Seq:     "AAAAAAGTAGTGTG",
			Plasmid: false,
		},
		ErrMessage: "",
	},
	{
		LargeSeq: &wtype.DNASequence{
			Nm:      "Test5",
			Seq:     "ACCGTAGTGTG",
			Plasmid: true,
		},
		ReplaceWithSeq: &wtype.DNASequence{
			Nm:      "nothing",
			Seq:     "",
			Plasmid: false,
		},
		Position: PositionPair{
			StartPosition: 4,
			EndPosition:   1,
			Reverse:       false,
		},
		ExpectedResult: &wtype.DNASequence{
			Nm:      "Test5",
			Seq:     "CC",
			Plasmid: true,
		},
		ErrMessage: "",
	},
	{
		LargeSeq: &wtype.DNASequence{
			Nm:      "Test6",
			Seq:     "ACCGTAGTGTG",
			Plasmid: true,
		},
		ReplaceWithSeq: &wtype.DNASequence{
			Nm:      "TTT",
			Seq:     "TTT",
			Plasmid: false,
		},
		Position: PositionPair{
			StartPosition: 3,
			EndPosition:   1,
			Reverse:       true,
		},
		ExpectedResult: &wtype.DNASequence{
			Nm:      "Test6",
			Seq:     "AAAGTAGTGTG",
			Plasmid: true,
		},
		ErrMessage: "",
	},
	{
		LargeSeq: &wtype.DNASequence{
			Nm:      "Test7",
			Seq:     "ACCGTAGTGTG",
			Plasmid: false,
		},
		ReplaceWithSeq: &wtype.DNASequence{
			Nm:      "TTT",
			Seq:     "TTT",
			Plasmid: false,
		},
		Position: PositionPair{
			StartPosition: 3,
			EndPosition:   1,
			Reverse:       false,
		},
		ExpectedResult: &wtype.DNASequence{
			Nm:      "Test7",
			Seq:     "ACCGTAGTGTG",
			Plasmid: false,
		},
		ErrMessage: "invalid position {StartPosition:3 EndPosition:1 Reverse:false} to replace in sequence. Start position must be lower than end position unless position is reverse or sequence is a plasmid. Sequence Test7: Plasmid = false",
	},
}

func TestReplaceAll(t *testing.T) {
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

func TestReplace(t *testing.T) {
	for _, test := range replacePositionTests {
		result, err := Replace(*test.LargeSeq, test.Position, *test.ReplaceWithSeq)
		if !reflect.DeepEqual(&result, test.ExpectedResult) {
			t.Error(
				"For", test.LargeSeq.Nm, test.LargeSeq.Seq, "\n",
				"replace position", test.Position, "\n",
				"with ", test.ReplaceWithSeq.Seq, "\n",
				"expected ", fmt.Sprintf("%+v", test.ExpectedResult.Seq), "\n",
				"got", fmt.Sprintf("%+v", result.Seq), "\n",
			)
		}
		if err != nil {
			if err.Error() != test.ErrMessage {
				t.Error(
					"For", test, "\n",
					"got error ", test.ErrMessage, "\n",
					"got", err.Error(), "\n",
				)
			}
		}
	}
}

func TestRotate(t *testing.T) {
	for _, test := range rotateTests {
		result := Rotate(wtype.DNASequence{Seq: test.seq}, test.rotateBy, test.reverse)
		if result.Seq != test.expectedResult {
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
