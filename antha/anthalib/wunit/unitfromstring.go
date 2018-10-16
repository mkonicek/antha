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
	"fmt"
	"strconv"
	"strings"
)

// SplitValueAndUnit splits a joined value and unit in string format into seperate typed value and unit fields.
// If the string input is not in the valid format of value followed by unit it will not be parsed correctly.
// If a value on its own is given the unit will be returned blank, if the unit is given alone the value will be 0.0
// valid: 10s, 10 s, 10.5s, 2.16e+04 s, 10, s
// invalid: s 10 s10
func SplitValueAndUnit(str string) (float64, string) {
	value, unit := extractFloat(str)
	return value, strings.TrimSpace(unit)
}

// ParseConcentration utility to extract concentration and component name from a string.
// Valid inputs include
//   - "6M Glucose", "6 M Glucose"
//   - "Glucose 6 M", "Glucose 6M"
//   - "Glucose (M)", "Glucose M", "Glucose (6 M)"
//   - "Glucose"
// returns three values - a boolean which is true if the value of the concentration was set,
// the parsed concentration, and the remaining component name.
// if removing the concentration would leave the empty string, the component name is
// set as the input string
func ParseConcentration(s string) (bool, Concentration, string) {

	reg := GetGlobalUnitRegistry()

	// unit location indicated by parentheses after component Name - e.g. "Glucose (6M)"
	if l, r := strings.LastIndex(s, "("), strings.LastIndex(s, ")"); l >= 0 && l < r && strings.HasSuffix(s, ")") {
		value, unit := extractFloat(strings.TrimSpace(s[l+1 : r]))
		if trimUnit := strings.TrimSpace(unit); reg.ValidUnitForType("Concentration", trimUnit) {
			if componentName := strings.TrimSpace(s[:l] + s[r+1:]); componentName != "" {
				return true, NewConcentration(value, trimUnit), componentName
			} else {
				return true, NewConcentration(value, trimUnit), s
			}
		}
	}

	approvedUnits := reg.ListValidUnitsForType("Concentration")

	// value and unit at left - e.g. "6M Glucose", but not "6 Glucose"
	if value, remainder := SplitValueAndUnit(s); len(remainder) < len(s) { // value must be given for unit at left
		if sym, componentName := extractSymbol(remainder, approvedUnits); sym != "" {
			if componentName := strings.TrimSpace(componentName); componentName != "" {
				return true, NewConcentration(value, sym), componentName
			}
		}
	}

	// unit at right - e.g. "Glucose 6M" or "Glucose M", but not "SolutionX"
	if sym, remainder := extractLastSymbol(s, approvedUnits); sym != "" && remainder != "" {
		trimRemainder := strings.TrimSpace(remainder)
		// units at the right must be preceded by either a number or a space
		if value, componentName := extractLastFloat(trimRemainder); len(componentName) != len(trimRemainder) || strings.HasSuffix(remainder, " ") {
			if componentName := strings.TrimSpace(componentName); componentName != "" {
				return true, NewConcentration(value, sym), componentName
			} else {
				return true, NewConcentration(value, sym), s
			}
		}
	}

	// no unit found
	return false, NewConcentration(0, "g/l"), s
}

// ParseVolume parses a volume and valid unit (nl, ul, ml, l) in string format; handles cases where the volume is split with a space.
func ParseVolume(volstring string) (volume Volume, err error) {
	var volandunit []string

	sortedKeys := GetGlobalUnitRegistry().ListValidUnitsForType("Volume")

	var longestmatchedunit string

	for _, approvedUnit := range sortedKeys {
		if strings.HasSuffix(volstring, approvedUnit) {
			volandunit = []string{strings.TrimSpace(strings.Trim(volstring, approvedUnit)), approvedUnit}
			if len(volandunit[1]) > len(longestmatchedunit) {
				longestmatchedunit = volandunit[1]
			}
		}
	}

	if len(longestmatchedunit) == 0 {
		err = fmt.Errorf("no valid unit found for %s: valid units are: %v", volstring, sortedKeys)
		return
	}

	if len(volandunit) == 0 {
		err = fmt.Errorf("error parsing volume for %s", volstring)
		return
	}

	vol, err := strconv.ParseFloat(strings.TrimSpace(volandunit[0]), 64)

	if err != nil {
		return
	}

	volume = NewVolume(vol, longestmatchedunit)
	return
}
