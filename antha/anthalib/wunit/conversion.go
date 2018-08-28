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
	"github.com/pkg/errors"
)

// MasstoVolume divides a mass (in kg) by a density (in kg/m^3) and returns the volume (in L).
func MasstoVolume(m Mass, d Density) Volume {
	return NewVolume(m.ConvertToString("g")/d.ConvertToString("kg/m^3"), "l")
}

// VolumetoMass multiplies a volume (in L) by a density (in kg/m^3) and returns the mass (in kg).
func VolumetoMass(v Volume, d Density) Mass {
	return NewMass(v.ConvertToString("m^3")*d.ConvertToString("kg/m^3"), "kg")
}

// VolumeForTargetMass returns the volume required to convert a starting stock concentration to a solution containing a target mass.
// returns an error if the concentration units are not based on g/l.
// If the stock concentration is zero a volume of 0ul will be returned with an error.
// if the target mass is zero a volume of 0ul will be returned with no error.
func VolumeForTargetMass(targetMass Mass, stockConc Concentration) (Volume, error) {

	if stockConc.IsZero() {
		return Volume{}, errors.New("stock concentration cannot be zero")
	}

	if targetMass.IsZero() {
		return NewVolume(0.0, "ul"), nil
	}

	if gramsPerULitre, err := GetGlobalUnitRegistry().GetUnit("g/ul"); err != nil {
		return Volume{}, err
	} else if !stockConc.Unit().CompatibleWith(gramsPerULitre) {
		return Volume{}, errors.Errorf("invalid stock concentration units %v: must be based on %v", stockConc.Unit(), gramsPerULitre)
	} else if grams, err := GetGlobalUnitRegistry().GetUnit("g"); err != nil {
		return Volume{}, err
	} else if !targetMass.Unit().CompatibleWith(grams) {
		return Volume{}, errors.Errorf("invalid target mass units %v: must be based on %v", targetMass.Unit(), grams)
	} else {
		return NewVolume(targetMass.ConvertTo(grams)/stockConc.ConvertTo(gramsPerULitre), "ul"), nil
	}
}

// VolumeForTargetConcentration returns the volume required to convert a starting stock concentration to a target concentration of volume total volume
// returns an error if the concentration units are incompatible (M/l and g/L) or if the target concentration is higher than the stock concentration
// unless the total volume is zero
func VolumeForTargetConcentration(targetConc Concentration, stockConc Concentration, totalVol Volume) (Volume, error) {
	if !targetConc.Unit().CompatibleWith(stockConc.Unit()) {
		return Volume{}, errors.Errorf("incompatible units %v and %v", targetConc.Unit(), stockConc.Unit())
	} else if totalVol.IsZero() {
		return NewVolume(0.0, "ul"), nil
	} else if stockConc.LessThan(targetConc) {
		return Volume{}, errors.Errorf("cannot dilute stock at %v to higher concentration %v", stockConc, targetConc)
	} else {
		return NewVolume(totalVol.RawValue()*targetConc.RawValue()/stockConc.ConvertTo(targetConc.Unit()), totalVol.Unit().PrefixedSymbol()), nil
	}
}

// MassForTargetConcentration multiplies a concentration (in g/l) by a volume (in l) to return the mass (in g).
// if a concentration is not in a form convertable to g/l an error is returned.
func MassForTargetConcentration(targetConc Concentration, totalVol Volume) (Mass, error) {
	if litres, err := GetGlobalUnitRegistry().GetUnit("l"); err != nil {
		return Mass{nil}, err
	} else if gramsPerLitre, err := GetGlobalUnitRegistry().GetUnit("g/l"); err != nil {
		return Mass{nil}, err
	} else if !targetConc.Unit().CompatibleWith(gramsPerLitre) {
		return Mass{nil}, fmt.Errorf("cannot convert %v to %v: incomptible units", targetConc.Unit(), gramsPerLitre)
	} else if !totalVol.Unit().CompatibleWith(litres) {
		return Mass{nil}, fmt.Errorf("cannot convert %v to %v: incomptible units", totalVol.Unit(), litres)
	} else {
		return NewMass(targetConc.ConvertTo(gramsPerLitre)*totalVol.ConvertTo(litres), "g"), nil
	}
}
