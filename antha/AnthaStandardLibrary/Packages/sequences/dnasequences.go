// antha/AnthaStandardLibrary/Packages/enzymes/Utility.go: Part of the Antha language
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

// Package for interacting with and manipulating dna sequences in extension to methods available in wtype
package sequences

import (
	"strings"
)

func Rev(s string) string {
	r := ""

	for i := len(s) - 1; i >= 0; i-- {
		r += string(s[i])
	}

	return r
}
func Comp(s string) string {
	r := ""

	m := map[string]string{
		"A": "T",
		"T": "A",
		"U": "A",
		"C": "G",
		"G": "C",
		"Y": "R",
		"R": "Y",
		"W": "W",
		"S": "S",
		"K": "M",
		"M": "K",
		"D": "H",
		"V": "B",
		"H": "D",
		"B": "V",
		"N": "N",
		"X": "X",
	}

	for _, c := range s {
		r += m[string(c)]
	}

	return r
}

// Reverse Complement
func RevComp(s string) string {
	s = strings.ToUpper(s)
	return Comp(Rev(s))
}

func AllCombinations(arr [][]string) []string {
	if len(arr) == 1 {
		return arr[0]
	}

	results := make([]string, 0)
	allRem := AllCombinations(arr[1:len(arr)])
	for i := 0; i < len(allRem); i++ {
		for j := 0; j < len(arr[0]); j++ {
			x := arr[0][j] + allRem[i]
			results = append(results, x)
		}
	}
	return results
}

func Prefix(seq string, lengthofprefix int) (prefix string) {
	end := lengthofprefix
	if end > len(seq) {
		end = len(seq)
	}
	prefix = seq[:end]
	return prefix
}

func Suffix(seq string, lengthofsuffix int) (suffix string) {
	start := len(seq) - lengthofsuffix
	if start < 0 {
		start = 0
	}
	suffix = seq[start:]
	return suffix
}
