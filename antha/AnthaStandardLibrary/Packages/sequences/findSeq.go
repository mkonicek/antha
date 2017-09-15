// antha/AnthaStandardLibrary/Packages/sequences/findSeq.go: Part of the Antha language
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
	"strconv"
	"strings"

	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/search"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
)

// SearchResult stores the results of a search query against a template sequence.
type SearchResult struct {
	Template  wtype.BioSequence
	Query     wtype.BioSequence
	Positions []PositionPair
}

// PositionPair stores the Start and Endposition of feature in a sequence in human friendly format rather than code format.
// i.e. in a Sequence "ATGTGTTG" position 1 is A and there is no position zero.
// To convert the format, the methods HumanFriendly() and CodeFriendly() return the positions in the corresponding formats.
type PositionPair struct {
	StartPosition int
	EndPosition   int
	Reverse       bool
}

// Start returns the start position of the PositionPair
// by default this will return a directional human friendly position
func (p PositionPair) Start(options ...string) int {
	start, _ := p.Coordinates(options...)
	return start
}

// End returns the end position of the PositionPair
// by default this will return a directional human friendly position
func (p PositionPair) End(options ...string) int {
	_, end := p.Coordinates(options...)
	return end
}

// ignores case
func containsString(slice []string, testString string) bool {

	upper := func(s string) string { return strings.ToUpper(s) }

	for _, str := range slice {
		if upper(str) == upper(testString) {
			return true
		}
	}
	return false
}

// Coordinates returns the start and end positions of the feature
// by default this will return the start position followed by the end position
func (p *PositionPair) Coordinates(options ...string) (start, end int) {
	start, end = p.StartPosition, p.EndPosition
	if containsString(options, wtype.CODEFRIENDLY) {
		start--
		end--
	}
	if containsString(options, wtype.IGNOREDIRECTION) {
		if start > end {
			return end, start
		}
	}
	return start, end
}

// PositionPairSet obeys the sort interface making the position pairs to be sorted
// in ascending start position.
// Direction is ignored during sorting.
type PositionPairSet []PositionPair

func (p PositionPairSet) Len() int {
	return len(p)
}

func (p PositionPairSet) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func (p PositionPairSet) Less(i, j int) bool {
	starti, endi := p[i].HumanFriendly(IGNOREDIRECTION)

	startj, endj := p[j].HumanFriendly(IGNOREDIRECTION)

	return starti < startj && endi < endj

}

// IGNOREDIRECTION is a boolean constant to specify direction of a feature position
// should be ignored when returning start and end positions of a feature.
// If selected, the start position will be the first position at which the feature is encountered regardless of orientation.
const IGNOREDIRECTION bool = true

// HumanFriendly returns a sequence PositionPair's start and end positions in a human friendly format
// i.e. in a Sequence "ATGTGTTG" position 1 is A and there is no position zero.
// If ignoredirection is used as an argument and set to true, the start position will be the first position at which the feature is encountered regardless of orientation.
func (pair PositionPair) HumanFriendly(ignoredirection ...bool) (start, end int) {
	if len(ignoredirection) > 0 && ignoredirection[0] {
		if pair.Reverse {
			return pair.EndPosition, pair.StartPosition
		} else {
			return pair.StartPosition, pair.EndPosition
		}
	}
	return pair.StartPosition, pair.EndPosition
}

// CodeFriendly returns a sequence PositionPair's start and end positions in a code friendly format
// i.e. in a Sequence "ATGTGTTG" position 0 is A.
// If ignoredirection is used as an argument and set to true, the start position will be the first position at which the feature is encountered regardless of orientation.
func (pair PositionPair) CodeFriendly(ignoredirection ...bool) (start, end int) {
	if len(ignoredirection) > 0 && ignoredirection[0] {
		if pair.Reverse {
			return pair.EndPosition - 1, pair.StartPosition - 1
		} else {
			return pair.StartPosition - 1, pair.EndPosition - 1
		}
	}
	return pair.StartPosition - 1, pair.EndPosition - 1
}

// FindSeq searches for a DNA sequence within a larger DNA sequence and returns all matches on both coding and complimentary strands.
func FindSeq(bigSequence, smallSequence *wtype.DNASequence) (seqsFound SearchResult) {
	if len(smallSequence.Sequence()) > len(bigSequence.Sequence()) {
		seqsFound = SearchResult{
			Template: bigSequence,
			Query:    smallSequence,
		}
		return
	}
	seqsFound = findSeq(bigSequence, smallSequence)

	originalPairs := seqsFound.Positions

	var newPairs []PositionPair

	newPairs = append(newPairs, originalPairs...)

	// if a vector, attempt rotation of bigsequence vector index 1 position at a time.
	if bigSequence.Plasmid && !smallSequence.Plasmid {
		rotationSize := len(smallSequence.Seq) - 1
		var tempSequence wtype.DNASequence
		tempSequence.Nm = "test"
		var tempseq string

		tempseq += bigSequence.Seq[rotationSize:]
		tempseq += bigSequence.Seq[:rotationSize]
		tempSequence.Seq = tempseq

		tempSeqsFound := findSeq(&tempSequence, smallSequence)

		for j, positionPair := range tempSeqsFound.Positions {

			var skip bool

			if (positionPair.EndPosition + rotationSize) > len(bigSequence.Seq) {
				positionPair.EndPosition = positionPair.EndPosition + rotationSize - len(bigSequence.Seq)
			} else {
				positionPair.EndPosition = positionPair.EndPosition + rotationSize
			}

			if (positionPair.StartPosition + rotationSize) > len(bigSequence.Seq) {
				// correct position offset
				positionPair.StartPosition = positionPair.StartPosition + rotationSize - len(bigSequence.Seq)
			} else {
				positionPair.StartPosition = positionPair.StartPosition + rotationSize
			}
			tempSeqsFound.Positions[j] = positionPair
			// check if any new positions found
			for _, oldPosition := range newPairs {
				// if already present set skip to true
				if equalPositionPairs(positionPair, oldPosition) {
					skip = true
				}
			}
			// if no skip set add to pairs
			if !skip {
				newPairs = append(newPairs, positionPair)
			}
		}

	}

	seqsFound.Positions = newPairs

	return seqsFound
}

func equalPositionPairs(pair1, pair2 PositionPair) bool {
	if pair1.StartPosition == pair2.StartPosition && pair1.EndPosition == pair2.EndPosition && pair1.Reverse == pair2.Reverse {
		return true
	}
	return false
}

func equalPositionPairSets(positionSet1, positionSet2 []PositionPair) bool {
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

func findSeq(bigSequence, smallSequence *wtype.DNASequence) (seqsfound SearchResult) {
	bigseq := strings.ToUpper(bigSequence.Sequence())

	seqsfound = SearchResult{
		Template: bigSequence,
		Query:    smallSequence,
	}

	var positionsFound []PositionPair

	seq := strings.ToUpper(smallSequence.Sequence())
	if strings.Contains(bigseq, seq) {
		positions := search.Findall(bigseq, seq)
		if len(positions) > 0 {
			for _, position := range positions {
				var positionFound PositionPair = PositionPair{
					StartPosition: position,
					EndPosition:   position + len(seq) - 1,
					Reverse:       false,
				}
				positionsFound = append(positionsFound, positionFound)
			}
		}
	}

	revseq := strings.ToUpper(RevComp(seq))
	if strings.Contains(bigseq, revseq) {
		positions := search.Findall(bigseq, revseq)
		if len(positions) > 0 {
			for _, position := range positions {
				var positionFound PositionPair = PositionPair{
					EndPosition:   position,
					StartPosition: position + len(seq) - 1,
					Reverse:       true,
				}
				positionsFound = append(positionsFound, positionFound)
			}
		}
	}

	seqsfound.Positions = positionsFound

	return
}

// FindSeqsInSeqs searches for small sequences (as strings) in a big sequence.
// The sequence is considered to be linear and matches will not be found if the sequence is circular and the sequence overlaps the end of the sequence.
// In this case, FindSeqs should be used.
func FindSeqsinSeqs(bigseq string, smallseqs []string) (seqsfound []search.Thingfound) {

	bigseq = strings.ToUpper(bigseq)

	var seqfound search.Thingfound
	seqsfound = make([]search.Thingfound, 0)
	for _, seq := range smallseqs {
		seq = strings.ToUpper(seq)
		if strings.Contains(bigseq, seq) {
			seqfound.Thing = seq
			seqfound.Positions = search.Findall(bigseq, seq)
			seqsfound = append(seqsfound, seqfound)
		}
	}
	for _, seq := range smallseqs {
		revseq := strings.ToUpper(RevComp(seq))
		if strings.Contains(bigseq, revseq) {
			seqfound.Thing = revseq
			seqfound.Positions = search.Findall(bigseq, revseq)
			seqfound.Reverse = true
			seqsfound = append(seqsfound, seqfound)
		}
	}

	return seqsfound
}

// FindPositionInSequence returns directionless Positions; if a feature is found in the reverse orientation the first position found
// in the sequence will be returned rather than the start of the feature.
// If more than one matching feature is found an error will be returned.
func FindPositionInSequence(largeSequence wtype.DNASequence, smallSequence wtype.DNASequence) (start int, end int, err error) {

	seqsfound := FindSeq(&largeSequence, &smallSequence)

	if len(seqsfound.Positions) != 1 {
		errstr := fmt.Sprint(strconv.Itoa(len(seqsfound.Positions)), " sequences of ", smallSequence.Nm, " ", smallSequence.Seq, " found in ", largeSequence.Nm, " ", largeSequence.Seq)
		err = fmt.Errorf(errstr)
		return
	}
	start, end = seqsfound.Positions[0].HumanFriendly(IGNOREDIRECTION)
	return
}

// FindDirectionalPositionInSequence returns the directional Positions of the feature.
// If more than one matching feature is found an error will be returned.
func FindDirectionalPositionInSequence(largeSequence wtype.DNASequence, smallSequence wtype.DNASequence) (start int, end int, err error) {
	seqsfound := FindSeq(&largeSequence, &smallSequence)

	if len(seqsfound.Positions) != 1 {
		errstr := fmt.Sprint(strconv.Itoa(len(seqsfound.Positions)), " sequences of ", smallSequence.Nm, " ", smallSequence.Seq, " found in ", largeSequence.Nm, " ", largeSequence.Seq)
		err = fmt.Errorf(errstr)
		return
	}
	start, end = seqsfound.Positions[0].HumanFriendly()
	return
}
