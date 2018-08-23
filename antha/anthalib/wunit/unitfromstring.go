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
	"sort"
	"strconv"
	"strings"

	"github.com/antha-lang/antha/antha/anthalib/wutil"
)

// SplitValueAndUnit splits a joined value and unit in string format into seperate typed value and unit fields.
// If the string input is not in the valid format of value followed by unit it will not be parsed correctly.
// If a value on its own is given the unit will be returned blank, if the unit is given alone the value will be 0.0
// valid: 10s, 10 s, 10.5s, 2.16e+04 s, 10, s
// invalid: s 10 s10
func SplitValueAndUnit(str string) (value float64, unit string) {

	fields := strings.Fields(str)

	if len(fields) == 2 {
		if value, err := strconv.ParseFloat(fields[0], 64); err == nil {
			return value, fields[1]
		}
		if value, err := strconv.Atoi(fields[0]); err == nil {
			return float64(value), fields[1]
		}
	} else if len(fields) == 1 {
		for i := 0; i < len(str); i++ {
			if value, err := strconv.ParseFloat(str[:len(str)-i], 64); err == nil {
				return value, str[len(str)-i:]
			}
		}
	}
	return value, str
}

// Utility function to parse concentration from a component name.
// Not currently robust to situations where the component name (without the concentration) is more than one field (e.g. ammonium sulphate) or if the component name is a concatenation of component names (e.g. 1mM Glucose + 10mM Glycerol).
func ParseConcentration(componentname string) (containsconc bool, conc Concentration, componentNameOnly string) {

	approvedunits := UnitMap["Concentration"]

	var sortedKeys []string

	for k := range approvedunits {
		sortedKeys = append(sortedKeys, k)
	}

	sort.Strings(sortedKeys)

	fields := strings.Fields(componentname)

	if len(fields) == 1 {
		trimmed := strings.Trim(componentname, "()")
		value, unit := SplitValueAndUnit(trimmed)
		if unit == componentname {
			return false, conc, componentname
		}
		if err := ValidMeasurementUnit("Concentration", unit); err != nil {
			return false, conc, componentname
		}
		return true, NewConcentration(value, unit), componentname
	}
	var unitmatchlength int
	var longestmatchedunit string
	var valueandunit string
	var unit string
	var valueString string
	var notConcFields []string
	for _, key := range sortedKeys {
		for i, field := range fields {

			/// if value and unit are separate fields
			if field == key && i != 0 {
				f, err := strconv.ParseFloat(fields[i-1], 64)
				if err == nil && f != 0 {
					if len(key) > unitmatchlength {
						notConcFields = make([]string, 0)
						longestmatchedunit = key
						unitmatchlength = len(key)
						valueString = fields[i-1]
						unit = field
						valueandunit = valueString + unit
						if (i - 1) > 0 {
							notConcFields = append(notConcFields, fields[:i-1]...)
						}
						if len(fields) > i+1 {
							notConcFields = append(notConcFields, fields[i+1:]...)
						}
						//break
					}
					// support for cases where concentration unit is given but no value
				} else if trimmed := strings.Trim(field, "()"); trimmed == key || field == key {
					if len(key) > unitmatchlength {
						notConcFields = make([]string, 0)
						longestmatchedunit = key
						unitmatchlength = len(key)
						valueandunit = field
						if i > 0 {
							notConcFields = append(notConcFields, fields[:i]...)
						}
						if len(fields) > i+1 {
							notConcFields = append(notConcFields, fields[i+1:]...)
						}
						//break
					}
				}
				// if value and unit are one joined field
				// change this to separate number and match rest of valueandunit
			} else if trimmed := strings.Trim(field, "()"); trimmed == key || field == key {
				if len(key) > unitmatchlength {
					notConcFields = make([]string, 0)
					longestmatchedunit = key
					unitmatchlength = len(key)
					valueandunit = field
					if i > 0 {
						notConcFields = append(notConcFields, fields[:i]...)
					}
					if len(fields) > i+1 {
						notConcFields = append(notConcFields, fields[i+1:]...)
					}
					//break
				}
			} else if trimmed := strings.Trim(field, "()"); looksLikeNumberAndUnit(field, key) || looksLikeNumberAndUnit(trimmed, key) {
				if len(key) > unitmatchlength {
					notConcFields = make([]string, 0)
					longestmatchedunit = key
					unitmatchlength = len(key)
					valueandunit = field
					if i > 0 {
						notConcFields = append(notConcFields, fields[:i]...)
					}
					if len(fields) > i+1 {
						notConcFields = append(notConcFields, fields[i+1:]...)
					}
					//break
				}
			}
		}
	}

	componentNameOnly = strings.Join(notConcFields, " ")

	// if no match, return original component name
	if unitmatchlength == 0 {
		return false, conc, componentname
	}

	concfields := strings.Split(valueandunit, longestmatchedunit)

	value, err := strconv.ParseFloat(concfields[0], 64)
	if err != nil {
		concfields[0] = strings.Trim(concfields[0], "()")
		value, err = strconv.ParseFloat(concfields[0], 64)
		if err != nil {
			if concfields[0] == "" {
				value = 0.0
			} else {
				if strings.Contains(componentname, wutil.MIXDELIMITER) {
					return false, conc, componentname
				} else {
					fmt.Println("warning parsing componentname: ", componentname, ": ", err.Error())
					return false, conc, componentname
				}
			}
		}
	}

	conc = NewConcentration(value, longestmatchedunit)
	containsconc = true

	return containsconc, conc, componentNameOnly
}
func looksLikeNumberAndUnit(testString string, targetUnit string) bool {
	if strings.HasSuffix(testString, targetUnit) {
		trimmed := strings.Split(testString, targetUnit)
		if len(trimmed) == 0 {
			return false
		}
		_, err := strconv.ParseFloat(trimmed[0], 64)
		if err == nil {
			return true
		}
		_, err = strconv.Atoi(trimmed[0])
		if err == nil {
			return true
		}
	}
	return false
}

// ParseVolume parses a volume and valid unit (nl, ul, ml, l) in string format; handles cases where the volume is split with a space.
func ParseVolume(volstring string) (volume Volume, err error) {
	var volandunit []string

	approvedunits := UnitMap["Volume"]

	var sortedKeys []string

	for k := range approvedunits {
		sortedKeys = append(sortedKeys, k)
	}

	sort.Strings(sortedKeys)

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
