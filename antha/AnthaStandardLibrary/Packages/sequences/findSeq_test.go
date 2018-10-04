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
	"testing"

	"github.com/antha-lang/antha/antha/anthalib/wtype"
)

type regionTest struct {
	LargeSeq   *wtype.DNASequence
	SmallSeq   *wtype.DNASequence
	Positions  []PositionPair
	errMessage string
}

var regionTests = []regionTest{

	{
		LargeSeq: &wtype.DNASequence{
			Nm:      "Test1",
			Seq:     "ATCGTAGTGTG",
			Plasmid: false,
		},
		SmallSeq: &wtype.DNASequence{
			Nm:      "TGA",
			Seq:     "TAG",
			Plasmid: false,
		},
		Positions: []PositionPair{
			{
				StartPosition: 5,
				EndPosition:   7,
				Reverse:       false,
			},
		},
	},
	{
		LargeSeq: &wtype.DNASequence{
			Nm:      "wrapAroundFWD",
			Seq:     "ATCGTAGTGTG",
			Plasmid: true,
		},
		SmallSeq: &wtype.DNASequence{
			Nm:      "TGA",
			Seq:     "TGA",
			Plasmid: false,
		},
		Positions: []PositionPair{
			{
				StartPosition: 10,
				EndPosition:   1,
				Reverse:       false,
			},
		},
	},
	{
		LargeSeq: &wtype.DNASequence{
			Nm:      "ReverseHit",
			Seq:     "ATCGTAGTGTG",
			Plasmid: true,
		},
		SmallSeq: &wtype.DNASequence{
			Nm:      "CTA",
			Seq:     "CTA",
			Plasmid: false,
		},
		Positions: []PositionPair{
			{
				StartPosition: 7,
				EndPosition:   5,
				Reverse:       true,
			},
		},
	},
	{
		LargeSeq: &wtype.DNASequence{
			Nm:      "multipleHits",
			Seq:     "ATCGATGTGTG",
			Plasmid: true,
		},
		SmallSeq: &wtype.DNASequence{
			Nm:      "GAT",
			Seq:     "GAT",
			Plasmid: false,
		},
		Positions: []PositionPair{
			{
				StartPosition: 4,
				EndPosition:   6,
				Reverse:       false,
			},
			{
				StartPosition: 11,
				EndPosition:   2,
				Reverse:       false,
			},
			{
				StartPosition: 3,
				EndPosition:   1,
				Reverse:       true,
			},
		},
		errMessage: "3 sequences of GAT GAT found in multipleHits ATCGATGTGTG",
	},
	{
		LargeSeq: &wtype.DNASequence{
			Nm:      "fullSeq",
			Seq:     "ATCGTAGTGTG",
			Plasmid: true,
		},
		SmallSeq: &wtype.DNASequence{
			Nm:      "ATCGTAGTGTG",
			Seq:     "ATCGTAGTGTG",
			Plasmid: false,
		},
		Positions: []PositionPair{
			{
				StartPosition: 1,
				EndPosition:   11,
				Reverse:       false,
			},
		},
	},
	{
		LargeSeq: &wtype.DNASequence{
			Nm:      "PresentAsRotatedSeq",
			Seq:     "ATCGTAGTGTG",
			Plasmid: true,
		},
		SmallSeq: &wtype.DNASequence{
			Nm:      "TAGTGTGATCG",
			Seq:     "TAGTGTGATCG",
			Plasmid: false,
		},
		Positions: []PositionPair{
			{
				StartPosition: 5,
				EndPosition:   4,
				Reverse:       false,
			},
		},
	},
	{
		LargeSeq: &wtype.DNASequence{
			Nm:      "PresentAsRotatedSeq2",
			Seq:     "ATCGTAGTGTG",
			Plasmid: true,
		},
		SmallSeq: &wtype.DNASequence{
			Nm:      "TAGTGTGATC",
			Seq:     "TAGTGTGATC",
			Plasmid: false,
		},
		Positions: []PositionPair{
			{
				StartPosition: 5,
				EndPosition:   3,
				Reverse:       false,
			},
		},
	},
	{
		LargeSeq: &wtype.DNASequence{
			Nm:      "PresentAsRotatedSeq3",
			Seq:     "ATCGTAGTGTG",
			Plasmid: true,
		},
		SmallSeq: &wtype.DNASequence{
			Nm:      "TCGTAGTGTGA",
			Seq:     "TCGTAGTGTGA",
			Plasmid: false,
		},
		Positions: []PositionPair{
			{
				StartPosition: 2,
				EndPosition:   1,
				Reverse:       false,
			},
		},
	},
	{
		LargeSeq: &wtype.DNASequence{
			Nm:      "PresentAsRotatedSeq4",
			Seq:     "ATCGTAGTGTG",
			Plasmid: true,
		},
		SmallSeq: &wtype.DNASequence{
			Nm:      "GATCGTAGTGT",
			Seq:     "GATCGTAGTGT",
			Plasmid: false,
		},
		Positions: []PositionPair{
			{
				StartPosition: 11,
				EndPosition:   10,
				Reverse:       false,
			},
		},
	},
}

func equalPositionPairs(pair1, pair2 PositionPair) bool {
	if pair1.StartPosition == pair2.StartPosition && pair1.EndPosition == pair2.EndPosition && pair1.Reverse == pair2.Reverse {
		return true
	}
	return false
}

func equalPositionPairSets(positionSet1, positionSet2 []PositionPair) bool {
	if len(positionSet1) != len(positionSet2) {
		return false
	}

	for _, pos1 := range positionSet1 {
		var found bool
		for _, pos2 := range positionSet2 {
			if equalPositionPairs(pos1, pos2) {
				found = true
			}
		}
		if !found {
			return false
		}
	}
	return true
}

func TestFindSeq(t *testing.T) {
	for _, test := range regionTests {
		result := FindAll(test.LargeSeq, test.SmallSeq)
		if !equalPositionPairSets(result.Positions, test.Positions) {
			t.Error(
				"For", test.LargeSeq.Nm, "\n",
				"and", test.SmallSeq.Nm, "\n",
				"expected positions", fmt.Sprint(test.Positions), "\n",
				"got", fmt.Sprint(result.Positions), "\n",
			)
		}
	}
}

func TestFindPositioninSequence(t *testing.T) {
	for _, test := range regionTests {
		for _, position := range test.Positions {
			start, end, err := FindPositionInSequence(*test.LargeSeq, *test.SmallSeq)
			testStart, testEnd := position.HumanFriendly(true)
			if err == nil && start != testStart {
				t.Error(
					"For", test.LargeSeq.Nm, "\n",
					"expected Start:", testStart, "\n",
					"got", start, "\n",
				)
			}

			if err == nil && end != testEnd {
				t.Error(
					"For", test.LargeSeq.Nm, "\n",
					"expected End:", testEnd, "\n",
					"got", end, "\n",
				)
			}

			if err != nil {
				if test.errMessage != err.Error() {
					t.Error(
						"For", test.LargeSeq.Nm, "\n",
						"got error", err.Error(), "\n",
					)
				}
			}

		}
	}
}
