// wunit/siprefix.go: Part of the Antha language
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
	"math"
)

// prefix library
var prefices map[string]SIPrefix

var allPrefixes = []string{
	"y",
	"z",
	"a",
	"f",
	"p",
	"n",
	"u",
	"m",
	"c",
	"d",
	"h",
	"k",
	"M",
	"G",
	"T",
	"P",
	"E",
	"Z",
	"Y",
}

var longPrefixNames = map[string]string{
	"y": "yocto",
	"z": "zepto",
	"a": "atto",
	"f": "femto",
	"p": "pico",
	"n": "nano",
	"u": "micro",
	"m": "milli",
	"c": "centi",
	"d": "deci",
	"":  "",
	" ": "",
	"h": "hecto",
	"k": "kilo",
	"M": "mega",
	"G": "giga",
	"T": "tera",
	"P": "peta",
	"E": "exa",
	"Z": "zetta",
	"Y": "yotta",
}

// structure defining an SI prefix
type SIPrefix struct {
	// prefix name
	Name string
	// meaning in base 10
	Value float64
}

// LongName get the long name of the prefix (e.g. "mega" instead of "m")
func (self SIPrefix) LongName() string {
	return longPrefixNames[self.Name]
}

// ListSIPrefixSymbols returns a list of all valid SI prefixes
func SIPrefixSymbols() []string {
	return allPrefixes
}

// helper function to allow lookup of prefix
func SIPrefixBySymbol(symbol string) SIPrefix {
	if prefices == nil {
		prefices = MakePrefices()
	}
	// sugar to allow using empty prefix
	if symbol == "" {
		symbol = " "
	}

	return prefices[symbol]
}

// make the prefix structure
func MakePrefices() map[string]SIPrefix {
	pref_map := make(map[string]SIPrefix, 20)
	exponent := -24
	pfcs := "yzafpnum"

	for _, rune := range pfcs {
		prefix := SIPrefix{string(rune), math.Pow10(exponent)}
		//	logger.Debug(fmt.Sprintln(prefix))
		pref_map[string(rune)] = prefix
		exponent += 3
	}

	pfcs = "cd h"

	exponent = -2

	for _, rune := range pfcs {
		prefix := SIPrefix{string(rune), math.Pow10(exponent)}
		pref_map[string(rune)] = prefix
		exponent += 1
	}

	exponent = 3

	pfcs = "kMGTPEZY"

	for _, rune := range pfcs {
		prefix := SIPrefix{string(rune), math.Pow10(exponent)}
		//	logger.Debug(fmt.Sprintln(prefix))
		pref_map[string(rune)] = prefix
		exponent += 3
	}

	return pref_map
}
