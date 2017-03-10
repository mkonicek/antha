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

// utility function to parse concentration from a component name
func ParseConcentration(componentname string) (containsconc bool, conc Concentration, componentNameOnly string) {

	approvedunits := UnitMap["Concentration"]

	fields := strings.Fields(componentname)
	var unitmatchlength int
	var longestmatchedunit string
	var valueandunit string

	for key, _ := range approvedunits {
		for _, field := range fields {
			if strings.Contains(field, key) {
				if len(key) > unitmatchlength {
					longestmatchedunit = key
					unitmatchlength = len(key)
					valueandunit = field
				}
			}
		}
	}

	for _, field := range fields {
		if len(fields) == 2 && field != longestmatchedunit {
			componentNameOnly = field
		}
	}

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
				panic(fmt.Sprint("error parsing componentname: ", componentname, ": ", err.Error()))
				return false, conc, componentNameOnly
			}
		}
	}

	conc = NewConcentration(value, longestmatchedunit)
	containsconc = true
	return containsconc, conc, componentNameOnly
}

// currently only parses ul; handles cases where the volume is split with a space
func ParseVolume(volstring string) (volume Volume, err error) {
	var volandunit []string
	/*
		approvedunits := wunit.UnitMap["Volume"]

		fields := strings.Fields(volstring)
		var unitmatchlength int
		var longestmatchedunit string
		var valueandunit string

		for key, _ := range approvedunits {
			for _,field := range fields {
			if strings.Contains(field,key){
				if len(key) > unitmatchlength {
					longestmatchedunit = key
					unitmatchlength = len(key)
					valueandunit = field
					}
				}
			}
		}
	*/

	//for _, unit := range approvedunits {
	if strings.Count(volstring, " ") == 1 {
		volandunit = strings.Split(volstring, " ")
	} else if strings.Count(volstring, "ul") == 1 && strings.HasSuffix(volstring, "ul") {
		volandunit = []string{strings.Trim(volstring, "ul"), "ul"}
	}

	//}

	vol, err := strconv.ParseFloat(strings.TrimSpace(volandunit[0]), 64)

	if err != nil {
		return
	}

	volume = NewVolume(vol, strings.TrimSpace(volandunit[1]))
	return
}

/*

func parseVol(volstring string) (volume Volume, err error) {
	approvedunits := wunit.UnitMap["Volume"]

	fields := strings.Fields(volstring)
	var unitmatchlength int
	var longestmatchedunit string
	var valueandunit string

	for key, _ := range approvedunits {
		for _,field := range fields {
		if strings.Contains(field,key){
			if len(key) > unitmatchlength {
				longestmatchedunit = key
				unitmatchlength = len(key)
				valueandunit = field
				}
			}
		}
	}

	for _, field := range fields {
		if len(fields)== 2 && field !=  longestmatchedunit {
			componentNameOnly = field
		}
	}

	// if no match, return original component name
	if unitmatchlength == 0 {
		return false, conc, componentname
	}

	concfields := strings.Split(valueandunit,longestmatchedunit)

	value, err := strconv.ParseFloat(concfields[0],64)
	if err != nil{
		concfields[0] = strings.Trim(concfields[0], "()")
		value, err = strconv.ParseFloat(concfields[0], 64)
		if err != nil {
			if concfields[0] == ""{
				value = 0.0
			}else{
			panic(fmt.Sprint("error parsing componentname: ", componentname,": ",err.Error()))
			return false, conc, componentNameOnly
			}
		}
	}




	conc = wunit.NewConcentration(value,longestmatchedunit)
	containsconc = true
	return
}
*/
