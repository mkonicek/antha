// wunit/unitfromstring.go: Part of the Antha language
// Copyright (C) 2014 the Antha authors. All rights reserved.
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

package wunit

import (
	"strconv"
	"strings"
)

// extractFloat extract the longest valid float from the left hand side of
// the string and returns the float and the remaining string.
// If no float is found, returns zero and the entire string.
func extractFloat(s string) (float64, string) {
	var longest int
	var ret float64
	for i := range s {
		if f, err := strconv.ParseFloat(s[:i+1], 64); err == nil {
			ret = f
			longest = i + 1
		}
	}

	return ret, s[longest:]
}

// extractLastFloat extract the longest valid float from the right hand side of
// the string and return the float and the remaining string.
// If no float is found, returns zero and the entire string.
func extractLastFloat(s string) (float64, string) {
	startAt := len(s)
	var ret float64
	for i := len(s) - 1; i >= 0; i-- {
		if f, err := strconv.ParseFloat(s[i:], 64); err == nil {
			ret = f
			startAt = i
		}
	}

	return ret, s[:startAt]
}

// extractSymbol extract the longest valid unit symbol from the left hand side
// of the given string, returning the symbol and the remaining string.
// The symbol must either be terminated with a space (which will not be included
// in the remainder) or the end of the string.
// If no valid symbol is found, return "" and the entire string.
// Examples
//   extractSymbol("M Glucose", []string{"M"}) -> ("M", "Glucose")
//   extractSymbol("M-Glucose", []string{"M"}) -> ("", "M-Glucose")
//   extractSymbol("M", []string{"M"}) -> ("M", "")
func extractSymbol(s string, validSymbols []string) (string, string) {
	longest := ""
	remainder := s
	for _, v := range validSymbols {
		if len(v) > len(longest) && strings.HasPrefix(s, v) {
			if r := s[len(v):]; strings.HasPrefix(r, " ") {
				remainder = r
				longest = v
			}
		}
	}
	return longest, remainder
}

// extractLastSymbol extract the longest valid unit symbol from the right hand side
// of the given string, returning the unit and the remaining string.
// If no valid units are found, return "" and the entire string
func extractLastSymbol(s string, validSymbols []string) (string, string) {
	longest := ""
	for _, v := range validSymbols {
		if len(v) > len(longest) && strings.HasSuffix(s, v) {
			longest = v
		}
	}
	return longest, s[:len(s)-len(longest)]
}
