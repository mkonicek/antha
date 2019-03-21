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
// by default this will return the start position followed by the end position in human friendly format
// Availabe options are:
// HUMANFRIENDLY returns a sequence PositionPair's start and end positions in a human friendly format
// i.e. in a Sequence "ATGTGTTG" position 1 is A, 2 is T.
// CODEFRIENDLY returns a sequence PositionPair's start and end positions in a code friendly format
// i.e. in a Sequence "ATGTGTTG" position 0 is A, 1 is T.
// IGNOREDIRECTION is a constant to specify that direction of a feature position
// should be ignored when returning start and end positions of a feature.
// If selected, the start position will be the first position at which the feature is encountered regardless of orientation.
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

// ByPositionPairStartPosition obeys the sort interface making the position pairs to be sorted
// in ascending start position.
// Direction is ignored during sorting.
type ByPositionPairStartPosition []PositionPair

// Len returns the number of PositionPairs in PositionPairSet
func (p ByPositionPairStartPosition) Len() int {
	return len(p)
}

// Swap changes positions of two entries in a PositionPairSet
func (p ByPositionPairStartPosition) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

// Less evaluates whether the entry of PositionPairSet with index i is less than entry with index j
// the directionless start position is used to assess this.
// If the start positions are the same the end position is used.
func (p ByPositionPairStartPosition) Less(i, j int) bool {
	starti, endi := p[i].HumanFriendly(IGNOREDIRECTION)

	startj, endj := p[j].HumanFriendly(IGNOREDIRECTION)

	if startj == starti {
		return endi < endj
	}

	return starti < startj

}

// IGNOREDIRECTION is a boolean constant to specify direction of a feature position
// should be ignored when returning start and end positions of a feature.
// If selected, the start position will be the first position at which the feature is encountered regardless of orientation.
const IGNOREDIRECTION bool = true

// HumanFriendly returns a sequence PositionPair's start and end positions in a human friendly format
// i.e. in a Sequence "ATGTGTTG" position 1 is A and there is no position zero.
// If ignoredirection is used as an argument and set to true, the start position will be the first position at which the feature is encountered regardless of orientation.
func (p PositionPair) HumanFriendly(ignoredirection ...bool) (start, end int) {
	if len(ignoredirection) > 0 && ignoredirection[0] {
		if p.Reverse {
			return p.EndPosition, p.StartPosition
		}
	}
	return p.StartPosition, p.EndPosition
}

// CodeFriendly returns a sequence PositionPair's start and end positions in a code friendly format
// i.e. in a Sequence "ATGTGTTG" position 0 is A.
// If ignoredirection is used as an argument and set to true, the start position will be the first position at which the feature is encountered regardless of orientation.
func (p PositionPair) CodeFriendly(ignoredirection ...bool) (start, end int) {
	if len(ignoredirection) > 0 && ignoredirection[0] {
		if p.Reverse {
			return p.EndPosition - 1, p.StartPosition - 1
		}
	}
	return p.StartPosition - 1, p.EndPosition - 1
}

func upper(seq wtype.DNASequence) wtype.DNASequence {
	newSeq := seq
	newSeq.Seq = strings.ToUpper(seq.Seq)
	return newSeq
}

// correctPositions will correct the position assignment of a search on a plasmid sequence.
// A search on a plasmid sequence may need to fins matches which overlap the end of a plasmid.
// findSeq will therefore first concatenates the plasmid sequence with a duplicate and then perform a search.
// correctPositions will correct the position assignment of any matches which are found that overlap the end of a plasmid sequence.
func correctPositions(positionPair PositionPair, originalSequence wtype.DNASequence) (start int, end int, skip bool) {
	start, end = positionPair.Coordinates()

	if start > len(originalSequence.Sequence()) {
		if end > len(originalSequence.Sequence()) {
			return start - len(originalSequence.Sequence()), end - len(originalSequence.Sequence()), true
		}
		if end <= len(originalSequence.Sequence()) {
			return start - len(originalSequence.Sequence()), end, false
		}
		return -1, -1, true
	} else if end > len(originalSequence.Sequence()) {
		return start, end - len(originalSequence.Sequence()), false
	}
	return start, end, false
}

// EqualFold compares whether two sequences are equivalent to each other.
//
// The comparison will be performed in a case insensitive manner with respect to
// the actual sequence.
// The orientation is not important;
// i.e. a sequence and it's reverse complement will be classsified as equal.
// The two sequences must have the same circularisation status (i.e. both plasmid or both linear).
// If the sequences are plasmids then the rotation of the sequences is not important.
// Feature Annotations, double or single stranded status and overhang information
// are not taken into consideration.
func EqualFold(a, b *wtype.DNASequence) bool {
	if len(a.Seq) != len(b.Seq) {
		return false
	}
	if b.Plasmid != a.Plasmid {
		return false
	}

	if len(FindAll(a, b).Positions) == 1 {
		return true
	}

	return false
}

// InSequences evaluates whether a query is present in a set of DNASequences
// using the same criteria as the EqualFold function.
//
// The comparison will be performed in a case insensitive manner with respect to
// the actual sequence.
// The orientation is not important;
// i.e. a sequence and it's reverse complement will be classsified as equal.
// The two sequences must have the same circularisation status (i.e. both plasmid or both linear).
// If the sequences are plasmids then the rotation of the sequences is not important.
// Feature Annotations, double or single stranded status and overhang information
// are not taken into consideration.
func InSequences(seqs []*wtype.DNASequence, query *wtype.DNASequence) bool {
	for _, seq := range seqs {
		if EqualFold(seq, query) {
			return true
		}
	}
	return false
}

// FindAll searches for a DNA sequence within a larger DNA sequence and returns all matches on both coding and complimentary strands.
func FindAll(bigSequence, smallSequence *wtype.DNASequence) (seqsFound SearchResult) {
	if len(smallSequence.Sequence()) > len(bigSequence.Sequence()) {
		seqsFound = SearchResult{
			Template: bigSequence,
			Query:    smallSequence,
		}
		return
	}

	var newPairs []PositionPair

	// if a vector, attempt rotation of bigsequence vector.
	if bigSequence.Plasmid {
		//rotationSize := len(smallSequence.Seq)
		var tempSequence wtype.DNASequence

		err := tempSequence.Append(bigSequence.Sequence())
		if err != nil {
			panic(err)
		}
		err = tempSequence.Append(bigSequence.Sequence())
		if err != nil {
			panic(err)
		}

		tempSeqsFound := findSeq(&tempSequence, smallSequence)

		for _, positionPair := range tempSeqsFound.Positions {

			newStart, newEnd, skip := correctPositions(positionPair, *bigSequence)

			// if no skip set add to pairs
			if !skip {
				newPairs = append(newPairs, PositionPair{
					StartPosition: newStart,
					EndPosition:   newEnd,
					Reverse:       positionPair.Reverse,
				})
			}
		}

	} else {
		seqsFound = findSeq(bigSequence, smallSequence)

		newPairs = append(newPairs, seqsFound.Positions...)
	}

	seqsFound.Positions = newPairs

	return seqsFound
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
		positions := search.FindAll(bigseq, seq)
		if len(positions) > 0 {
			for _, position := range positions {
				var positionFound = PositionPair{
					StartPosition: position,
					EndPosition:   position + len(seq) - 1,
					Reverse:       false,
				}
				positionsFound = append(positionsFound, positionFound)
			}
		}
	}

	revseq := strings.ToUpper(wtype.RevComp(seq))
	if strings.Contains(bigseq, revseq) {
		positions := search.FindAll(bigseq, revseq)
		if len(positions) > 0 {
			for _, position := range positions {
				var positionFound = PositionPair{
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

// FindSeqsinSeqs searches for small sequences (as strings) in a big sequence.
// The sequence is considered to be linear and matches will not be found if the sequence is circular and the sequence overlaps the end of the sequence.
// In this case, FindSeqs should be used.
func FindSeqsinSeqs(bigseq string, smallseqs []string) (seqsfound []search.Result) {

	bigseq = strings.ToUpper(bigseq)

	var seqfound search.Result
	seqsfound = make([]search.Result, 0)
	for _, seq := range smallseqs {
		seq = strings.ToUpper(seq)
		if strings.Contains(bigseq, seq) {
			seqfound.Thing = seq
			seqfound.Positions = search.FindAll(bigseq, seq)
			seqsfound = append(seqsfound, seqfound)
		}
	}
	for _, seq := range smallseqs {
		revseq := strings.ToUpper(wtype.RevComp(seq))
		if strings.Contains(bigseq, revseq) {
			seqfound.Thing = revseq
			seqfound.Positions = search.FindAll(bigseq, revseq)
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

	seqsfound := FindAll(&largeSequence, &smallSequence)

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
	seqsfound := FindAll(&largeSequence, &smallSequence)

	if len(seqsfound.Positions) != 1 {
		errstr := fmt.Sprint(strconv.Itoa(len(seqsfound.Positions)), " sequences of ", smallSequence.Nm, " ", smallSequence.Seq, " found in ", largeSequence.Nm, " ", largeSequence.Seq)
		err = fmt.Errorf(errstr)
		return
	}
	start, end = seqsfound.Positions[0].HumanFriendly()
	return
}
