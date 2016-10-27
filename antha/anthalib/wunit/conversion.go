// Part of the Antha language
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

// Core Antha package for dealing with units in Antha
package wunit

import "fmt"

/*
type

func Splitunit(unit string)(numerators[]string, denominators[]string)

var conversiontable = map[string]map[string]float64{
	"density":map[string]float64{
		"g/L":
	}
}
*/

func MasstoVolume(m Mass, d Density) (v Volume) {
	mass := m.SIValue()

	if m.Unit().BaseSIUnit() == "g" {
		// work out mass in kg
		mass = mass / 1000
	}
	density := d.SIValue()
	volume := mass / density // in m^3
	volume = volume * 1000   // in l
	v = NewVolume(mass, "l")
	return v
}

func VolumetoMass(v Volume, d Density) (m Mass) {
	//mass := m.SIValue()
	density := d.SIValue()

	volume := v.SIValue() //* 1000 // convert m^3 to l

	mass := volume * density // in m^3

	m = NewMass(mass, "kg")
	return m
}

func VolumeForTargetMass(targetmass Mass, startingconc Concentration) (v Volume, err error) {
	fmt.Println("Base units ", startingconc.Unit().BaseSIUnit(), " and ", targetmass.Unit().BaseSIUnit())

	if startingconc.Unit().PrefixedSymbol() == "ng/ul" && targetmass.Unit().PrefixedSymbol() == "ng" {
		v = NewVolume(float64((targetmass.RawValue() / startingconc.RawValue())), "ul")
		fmt.Println("starting conc SI ", startingconc.SIValue(), " and target mass SI: ", targetmass.SIValue())

	} else if startingconc.Unit().PrefixedSymbol() == "mg/l" && targetmass.Unit().PrefixedSymbol() == "ng" {
		v = NewVolume(float64((targetmass.RawValue() / startingconc.RawValue())), "ul")
		fmt.Println("starting conc SI ", startingconc.SIValue(), " and target mass SI: ", targetmass.SIValue())

	} else if startingconc.Unit().BaseSIUnit() == "kg/l" && targetmass.Unit().BaseSIUnit() == "kg" {
		v = NewVolume(float64((targetmass.SIValue()/startingconc.SIValue())*1000000), "ul")
		fmt.Println("starting conc SI ", startingconc.SIValue(), " and target mass SI: ", targetmass.SIValue())

	} else if startingconc.Unit().BaseSIUnit() == "g/l" && targetmass.Unit().BaseSIUnit() == "g" {
		v = NewVolume(float64((targetmass.SIValue()/startingconc.SIValue())*1000000), "ul")
		fmt.Println("starting conc SI ", startingconc.SIValue(), " and target mass SI: ", targetmass.SIValue())
	} else {
		fmt.Println("Base units ", startingconc.Unit().BaseSIUnit(), " and ", targetmass.Unit().BaseSIUnit(), " not compatible with this function")
		err = fmt.Errorf("Convert ", targetmass.ToString(), " to g and ", startingconc.ToString(), " to g/l")
	}

	return
}

func VolumeForTargetConcentration(targetconc Concentration, startingconc Concentration, totalvol Volume) (v Volume, err error) {

	var factor float64

	if startingconc.Unit().BaseSIUnit() == targetconc.Unit().BaseSIUnit() {
		factor = targetconc.SIValue() / startingconc.SIValue()
	} else {
		err = fmt.Errorf("incompatible units of ", targetconc.ToString(), " and ", startingconc.ToString())
	}

	v = MultiplyVolume(totalvol, factor)

	//v = NewVolume(float64((targetconc.SIValue()/startingconc.SIValue())*1000000)*totalvol.SIValue(), "ul")

	return
}

func MassForTargetConcentration(targetconc Concentration, totalvol Volume) (m Mass, err error) {

	litre := NewVolume(1.0, "l")

	var multiplier float64 = 1
	var unit string

	if targetconc.Unit().PrefixedSymbol() == "kg/l" {
		multiplier = 1000
		unit = "g"
		//fmt.Println("targetconc.Unit().BaseSISymbol() == kg/l")
	} else if targetconc.Unit().PrefixedSymbol() == "g/l" {
		multiplier = 1
		unit = "g"
		//fmt.Println("targetconc.Unit().BaseSISymbol() == g/l")
	} else if targetconc.Unit().PrefixedSymbol() == "mg/l" {
		multiplier = 1
		unit = "mg"
	} else if targetconc.Unit().PrefixedSymbol() == "ng/ul" {
		multiplier = 1
		unit = "mg"
	} else {
		err = fmt.Errorf("Convert conc ", targetconc, " to g/l first")
	}

	m = NewMass(float64((targetconc.RawValue()*multiplier)*(totalvol.SIValue()/litre.SIValue())), unit)

	return
}
