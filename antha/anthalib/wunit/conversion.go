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
	"github.com/pkg/errors"
)

// MassToVolume divides a mass (in kg) by a density (in kg/m^3) and returns the volume (in L).
func MassToVolume(m Mass, d Density) (Volume, error) {
	if mass, err := m.InStringUnit("g"); err != nil {
		return Volume{}, err
	} else if density, err := d.InStringUnit("kg/m^3"); err != nil {
		return Volume{}, err
	} else {
		return NewVolume(mass.RawValue()/density.RawValue(), "l"), nil
	}
}

// VolumeToMass multiplies a volume (in L) by a density (in kg/m^3) and returns the mass (in kg).
func VolumeToMass(v Volume, d Density) (Mass, error) {
	if volume, err := v.InStringUnit("m^3"); err != nil {
		return Mass{}, err
	} else if density, err := d.InStringUnit("kg/m^3"); err != nil {
		return Mass{}, err
	} else {
		return NewMass(volume.RawValue()*density.RawValue(), "kg"), nil
	}
}

// MasstoVolume deprecated, please use MassToVolume instead
func MasstoVolume(m Mass, d Density) Volume {
	if ret, err := MassToVolume(m, d); err != nil {
		panic(err)
	} else {
		return ret
	}
}

// VolumetoMass deprecated, pelase use VolumeToMass
func VolumetoMass(v Volume, d Density) Mass {
	if ret, err := VolumeToMass(v, d); err != nil {
		panic(err)
	} else {
		return ret
	}
}

// VolumeForTargetMass returns the volume required to convert a starting stock concentration to a solution containing a target mass.
// returns an error if the concentration units are not based on g/l.
// If the stock concentration is zero a volume of 0ul will be returned with an error.
// if the target mass is zero a volume of 0ul will be returned with no error.
func VolumeForTargetMass(targetMass Mass, stockConc Concentration) (Volume, error) {
	if stockConc.IsZero() {
		return Volume{}, errors.New("stock concentration cannot be zero")
	} else if targetMass.IsZero() {
		return NewVolume(0.0, "ul"), nil
	} else if concInGramsPerULitre, err := stockConc.InStringUnit("g/ul"); err != nil {
		return Volume{}, err
	} else if massInGrams, err := targetMass.InStringUnit("g"); err != nil {
		return Volume{}, err
	} else {
		return NewVolume(massInGrams.RawValue()/concInGramsPerULitre.RawValue(), "ul"), nil
	}
}

// VolumeForTargetConcentration returns the volume required to convert a starting stock concentration to a target concentration of volume total volume
// returns an error if the concentration units are incompatible (M/l and g/L) or if the target concentration is higher than the stock concentration
// unless the total volume is zero
func VolumeForTargetConcentration(targetConc Concentration, stockConc Concentration, totalVol Volume) (Volume, error) {
	if totalVol.IsZero() {
		return NewVolume(0.0, "ul"), nil
	} else if stockConcInTargetUnits, err := stockConc.InUnit(targetConc.Unit()); err != nil {
		return Volume{}, err
	} else if stockConc.LessThan(targetConc) {
		return Volume{}, errors.Errorf("cannot dilute stock at %v to higher concentration %v", stockConc, targetConc)
	} else {
		return NewVolume(totalVol.RawValue()*targetConc.RawValue()/stockConcInTargetUnits.RawValue(), totalVol.Unit().PrefixedSymbol()), nil
	}
}

// MassForTargetConcentration multiplies a concentration (in g/l) by a volume (in l) to return the mass (in g).
// if a concentration is not in a form convertable to g/l an error is returned.
func MassForTargetConcentration(targetConc Concentration, totalVol Volume) (Mass, error) {
	if volumeInLitres, err := totalVol.InStringUnit("l"); err != nil {
		return Mass{nil}, err
	} else if concInGramsPerLitre, err := targetConc.InStringUnit("g/l"); err != nil {
		return Mass{nil}, err
	} else {
		return NewMass(concInGramsPerLitre.RawValue()*volumeInLitres.RawValue(), "g"), nil
	}
}
