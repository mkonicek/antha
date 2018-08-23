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

// Package wunit is a core Antha package for dealing with units in Antha
package wunit

import (
	"fmt"
)

// MasstoVolume divides a mass (in kg) by a density (in kg/m^3) and returns the volume (in L).
func MasstoVolume(m Mass, d Density) (v Volume) {

	mass := m.SIValue()

	if m.Unit().BaseSIUnit() == "g" {
		// work out mass in kg
		mass = mass / 1000
	}

	density := d.SIValue()
	fmt.Println(mass, density)
	volume := mass / density // in m^3
	volume = volume * 1000   // in l
	v = NewVolume(volume, "l")

	return v
}

// VolumetoMass multiplies a volume (in L) by a density (in kg/m^3) and returns the mass (in kg).
func VolumetoMass(v Volume, d Density) (m Mass) {
	density := d.SIValue()

	volume := v.SIValue() / 1000 // convert m^3 to l

	mass := volume * density // in m^3

	m = NewMass(mass, "kg")
	return m
}

// VolumeForTargetMass returns the volume required to convert a starting stock concentration to a solution containing a target mass.
// returns an error if the concentration units are not in g/l.
// If the stock concentration is zero a volume of 0ul will be returned with an error.
// if the target mass is zero a volume of 0ul will be returned with no error.
func VolumeForTargetMass(targetmass Mass, stockConc Concentration) (v Volume, err error) {

	if stockConc.RawValue() == 0.0 {
		v = NewVolume(0.0, "ul")
		return v, fmt.Errorf("Zero value found when converting concentration and mass to new volume so new volume set to zero: target mass: %s; starting concentration: %s", targetmass.ToString(), stockConc.ToString())
	}

	if targetmass.RawValue() == 0.0 {
		return NewVolume(0.0, "ul"), nil
	}

	if stockConc.Unit().PrefixedSymbol() == "ng/ul" && targetmass.Unit().PrefixedSymbol() == "ng" {
		v = NewVolume(float64((targetmass.RawValue() / stockConc.RawValue())), "ul")
		fmt.Println("starting conc SI ", stockConc.SIValue(), " and target mass SI: ", targetmass.SIValue())

	} else if stockConc.Unit().PrefixedSymbol() == "mg/l" && targetmass.Unit().PrefixedSymbol() == "ng" {
		v = NewVolume(float64((targetmass.RawValue() / stockConc.RawValue())), "ul")
		fmt.Println("starting conc SI ", stockConc.SIValue(), " and target mass SI: ", targetmass.SIValue())

	} else if stockConc.Unit().BaseSIUnit() == "kg/l" && targetmass.Unit().BaseSIUnit() == "kg" {
		v = NewVolume(float64((targetmass.SIValue()/stockConc.SIValue())*1000000), "ul")
		fmt.Println("starting conc SI ", stockConc.SIValue(), " and target mass SI: ", targetmass.SIValue())

	} else if stockConc.Unit().BaseSIUnit() == "g/l" && targetmass.Unit().BaseSIUnit() == "g" {
		v = NewVolume(float64((targetmass.SIValue()/stockConc.SIValue())*1000000), "ul")
		fmt.Println("starting conc SI ", stockConc.SIValue(), " and target mass SI: ", targetmass.SIValue())
	} else {
		fmt.Println("Base units ", stockConc.Unit().BaseSIUnit(), " and ", targetmass.Unit().BaseSIUnit(), " not compatible with this function")
		err = fmt.Errorf("Convert %s to g and %s to g/l", targetmass.ToString(), stockConc.ToString())
	}

	return
}

// VolumeForTargetConcentration returns the volume required to convert a starting stock concentration to a target concentration of volume total volume
// returns an error if the concentration units are incompatible (M/l and g/L) or if the target concentration is higher than the stock concentration
// If the stock concentration is zero a volume of 0ul will be returned with an error.
// if the target concetnration or total volume are set to zero a volume of 0ul will be returned with no error.
func VolumeForTargetConcentration(targetConc Concentration, stockConc Concentration, totalVol Volume) (v Volume, err error) {

	if stockConc.RawValue() == 0.0 {
		return NewVolume(0.0, "ul"), fmt.Errorf("Zero value found when converting concentrations to new volume so new volume set to zero: starting concentration: %s; final concentration: %s; volume set point: %s", stockConc.ToString(), targetConc.ToString(), totalVol.ToString())
	}

	if targetConc.RawValue() == 0.0 || totalVol.RawValue() == 0.0 {
		return NewVolume(0.0, "ul"), nil
	}

	factor, err := DivideConcentrations(targetConc, stockConc)

	if err != nil {
		return NewVolume(0.0, "ul"), fmt.Errorf("Error converting concentrations to new volume so new volume set to zero: starting concentration: %s; final concentration: %s; volume set point: %s. Error: %s", stockConc.ToString(), targetConc.ToString(), totalVol.ToString(), err.Error())
	}

	v = MultiplyVolume(totalVol, factor)

	if v.GreaterThan(totalVol) {
		err = fmt.Errorf(fmt.Sprint("Target concentration, ", targetConc.ToString(), " is higher than stock concentration ", stockConc.ToString(), " so volume calculated ", v.ToString(), " is larger than total volume ", totalVol.ToString()))
	}

	return
}

// MassForTargetConcentration multiplies a concentration (in g/l) by a volume (in l) to return the mass (in g).
// if a concentration is not in a form convertable to g/l an error is returned.
func MassForTargetConcentration(targetConc Concentration, totalVol Volume) (Mass, error) {
	if targetConc.Unit().BaseSIUnit() != "g/l" {
		return Mass{nil}, fmt.Errorf("cannot convert %v to mass: unit not based on g/l", targetConc)
	} else if litres, err := GetGlobalUnitRegistry().GetUnit("l"); err != nil {
		return Mass{nil}, err
	} else if gramsPerLitre, err := GetGlobalUnitRegistry().GetUnit("g/l"); err != nil {
		return Mass{nil}, err
	} else {
		return NewMass(targetConc.ConvertTo(gramsPerLitre)*totalVol.ConvertTo(litres), "g"), nil
	}
}
