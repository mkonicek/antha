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

// SIPrefix
type SIPrefix struct {
	Symbol string  // short version of the prefix
	Name   string  // long name of the prefix
	Value  float64 // multiplier that the exponent applies to the value
}

// SIPrefixes a list containing all valid SI prefixes
var SIPrefixes []SIPrefix

func newPrefix(symbol, name string, value float64) SIPrefix {
	ret := SIPrefix{
		Symbol: symbol,
		Name:   name,
		Value:  value,
	}
	SIPrefixes = append(SIPrefixes, ret)
	return ret
}

var ( //all supported SI prefixes, in smallest to largest order as they appear in SIPrefices
	Yocto = newPrefix("y", "yocto", 1e-24)
	Zepto = newPrefix("z", "zepto", 1e-21)
	Atto  = newPrefix("a", "atto", 1e-18)
	Femto = newPrefix("f", "femto", 1e-15)
	Pico  = newPrefix("p", "pico", 1e-12)
	Nano  = newPrefix("n", "nano", 1e-9)
	Micro = newPrefix("u", "micro", 1e-6)
	Milli = newPrefix("m", "milli", 1e-3)
	Centi = newPrefix("c", "centi", 1e-2)
	Deci  = newPrefix("d", "deci", 1e-1)
	Deca  = newPrefix("da", "deca", 1e1)
	Hecto = newPrefix("h", "hecto", 1e2)
	Kilo  = newPrefix("k", "kilo", 1e3)
	Mega  = newPrefix("M", "mega", 1e6)
	Giga  = newPrefix("G", "giga", 1e9)
	Tera  = newPrefix("T", "tera", 1e12)
	Peta  = newPrefix("P", "peta", 1e15)
	Exa   = newPrefix("E", "exa", 1e18)
	Zetta = newPrefix("Z", "zetta", 1e21)
	Yotta = newPrefix("Y", "yotta", 1e24)
	None  = SIPrefix{Symbol: "", Name: "", Value: 1.0} // not a valid SIPrefix, hence not in SIPrefixes, but used for non-prefixed units
)

// SIPrefixSymbols returns a list of all supported SI prefixes
func SIPrefixSymbols() []string {
	ret := make([]string, 0, len(SIPrefixes))
	for _, prefix := range SIPrefixes {
		ret = append(ret, prefix.Symbol)
	}
	return ret
}
