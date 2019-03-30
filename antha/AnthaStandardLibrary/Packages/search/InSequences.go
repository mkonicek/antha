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
	"reflect"
	"strings"

	"github.com/antha-lang/antha/antha/anthalib/wtype"
)

func trimmedEqual(a, b string) bool {
	return strings.TrimSpace(a) == strings.TrimSpace(b)
}

// InSequences searches the positions of any matching instances of a sequence in a slice of sequences.
// By default, the name and the sequence will be checked for equality;
// If IgnoreName is included as an option, matching by name is excluded.
// If IgnoreSequences is included as an option, matching by sequence is excluded.
// If IgnoreCase is added as an option the case will be ignored.
// Ignores any sequence annotations and overhang information.
// Circularisation is taken into account. i.e. i.e a plasmid will not match with a linear sequence.
// To require matches to be strict and identical (annotations, overhand info and all)
// use the ExactMatch option. If ExactMatch is specified, all other Options
// are ignored.
//
// Important Note: this function will not detect sequence equality if the sequence is:
// (A) a reverse complement
// (B) is a plasmid sequence which has been rotated.
// To handle these cases please use sequences.EqualFold and sequences.InSet:
// github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/sequences/findSeq.go
func InSequences(seqs []wtype.DNASequence, seq wtype.DNASequence, options ...Option) (bool, []int) {

	var positionsFound []int

	// ExactMatch takes priority and overrides other options.
	if containsExactMatch(options...) {
		for i := range seqs {
			if reflect.DeepEqual(seqs[i], seq) {
				positionsFound = append(
					positionsFound,
					i,
				)
			}
		}
		if len(positionsFound) > 0 {
			return true, positionsFound
		}
	}

	caseInsensitive := containsIgnoreCase(options...)

	ignoreNames := containsIgnoreName(options...)

	ignoreSequences := containsIgnoreSequence(options...)

	for i := range seqs {
		var nameMatches, seqMatches bool

		if !ignoreNames {
			if caseInsensitive && EqualFold(seqs[i].Name(), seq.Name()) {
				nameMatches = true
			} else if trimmedEqual(seqs[i].Name(), seq.Name()) {
				nameMatches = true
			}
		}

		if !ignoreSequences {
			if caseInsensitive && EqualFold(seqs[i].Sequence(), seq.Sequence()) && seqs[i].Plasmid == seq.Plasmid {
				seqMatches = true
			} else if trimmedEqual(seqs[i].Sequence(), seq.Sequence()) && seqs[i].Plasmid == seq.Plasmid {
				seqMatches = true
			}
		}

		if !ignoreNames && nameMatches {
			positionsFound = append(positionsFound, i)
		} else if !ignoreSequences && seqMatches {
			positionsFound = append(positionsFound, i)
		} else {
			if nameMatches && seqMatches {
				positionsFound = append(positionsFound, i)
			}
		}
	}

	if len(positionsFound) > 0 {
		return true, positionsFound
	}

	return false, positionsFound
}
