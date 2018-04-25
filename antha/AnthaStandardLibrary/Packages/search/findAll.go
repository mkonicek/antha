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
	"strconv"
	"strings"
)

// Result is the positions found of a Thing in a query string.
// Reverse is given as an option to specify if the Result corresponds to a reverse of the original query,
// or reverse compliment in the case of a DNA sequence.
// For searching in BioSequences the sequences.SearchResult type is preferred.
// The sequences.FindAll and sequences.Replace functions are the recommended mechanisms to use with
// the sequences.SearchResult type.
// Positions are recorded in Human friendly form, i.e. the position of A in ABC is 1 and not 0.
type Result struct {
	Thing     string
	Positions []int
	Reverse   bool
}

// HumanFriendlyPositions returns all positions in human friendly format where the first position is 1.
// i.e. the position of A in ABC is 1 and not 0.
func (r Result) HumanFriendlyPositions() []int {
	return r.Positions
}

// CodeFriendlyPositions returns all positions in code friendly format where the first position is 0.
// i.e. the position of A in ABC is 0 and not 1.
func (r Result) CodeFriendlyPositions() []int {
	var codeFriendlyPositions []int
	for _, position := range r.Positions {
		codeFriendlyPositions = append(codeFriendlyPositions, position-1)
	}
	return codeFriendlyPositions
}

// ToString returns a string description of the Result.
func (r Result) ToString() (descriptions string) {
	things := make([]string, 0)
	var reverse string
	for i := range r.Positions {
		if r.Reverse {
			reverse = " in reverse direction"
		} else {
			reverse = " in forward direction"
		}
		things = append(things, r.Thing, " found at position ", strconv.Itoa(r.Positions[i]), reverse, "; ")
	}
	descriptions = strings.Join(things, "")
	return
}

func upper(s string) string {
	return strings.ToUpper(s)
}

// FindAll searches for all instances of a target string in a template string.
// Not perfect yet! issue with byte conversion of certain characters!
// This returns positions in "HumanFriendly" format (i.e. the first position of the sequence will be 1 not 0)
// If the IgnoreCase option is specified the strings will be compared ignoring case.
//
func FindAll(template string, target string, options ...Option) (positions []int) {

	if containsIgnoreCase(options...) {
		template = upper(template)
		target = upper(target)
	}

	positions = make([]int, 0)
	count := strings.Count(template, target)

	if target == "" {
		return
	}
	if count != 0 {

		pos := (strings.Index(template, target))
		restofbigthing := template[(pos + 1):]

		for i := 0; i < count; i++ {
			positions = append(positions, (pos + 1))
			pos = pos + (strings.Index(restofbigthing, target) + 1)
			restofbigthing = template[(pos + 1):]
		}
	}
	return positions
}

// FindAllStrings searches for all instances of a slice of target strings in a template string.
// Not perfect yet! issue with byte conversion of certain characters!
// This returns positions in "HumanFriendly" format (i.e. the first position of the sequence will be 1 not 0)
// If the IgnoreCase option is specified the strings will be compared ignoring case.
//
func FindAllStrings(template string, targets []string, options ...Option) (thingsfound []Result) {

	var thingfound Result
	thingsfound = make([]Result, 0)

	for _, target := range targets {

		if containsIgnoreCase(options...) {
			template = upper(template)
			target = upper(target)
		}

		if strings.Contains(template, target) {
			thingfound.Thing = target
			thingfound.Positions = FindAll(template, target)
			thingsfound = append(thingsfound, thingfound)
		}
	}
	return thingsfound
}

// ContainsAllStrings searches for all instances of a slice of target strings in a template string.
// Not perfect yet! issue with byte conversion of certain characters!
// If all items are present, true is returned.
// If the IgnoreCase option is specified the strings will be compared ignoring case.
//
func ContainsAllStrings(template string, targets []string, options ...Option) (trueornot bool) {
	var i int
	for _, thing := range targets {

		if containsIgnoreCase(options...) {
			if strings.Contains(strings.ToUpper(template), strings.ToUpper(thing)) {
				i++
			}
		} else if strings.Contains(template, thing) {
			i++
		}
	}
	return i == len(targets)
}
