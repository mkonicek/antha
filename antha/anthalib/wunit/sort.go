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

import (
	"fmt"
	"sort"

	"github.com/pkg/errors"
)

var (
	noValues error = fmt.Errorf("empty slice specified as argument to sort function")
)

// SortOption is an option to add to
type SortOption string

// DescendingOrder is an optional argument which can be used in SortConcentrations and
// SortVolumes in order to sort by descending order.
const DescendingOrder SortOption = "DescendingOrder"

// SortConcentrations sorts a set of Concentration values.
// An error will be returned if no values are specified or the base units of any of the concentrations are incompatible,
// e.g. units of X and g/l would not be compatible.
// An optional second argument DescendingOrder can be specified to indicate that the sorting will be in descending order.
// By default sorting will be in ascending order.
func SortConcentrations(concs []Concentration, descending ...SortOption) (sorted []Concentration, err error) {
	if len(concs) == 0 {
		return concs, noValues
	}
	for i := range concs {
		if i > 0 {
			if err = sameUnit(concs[i].Unit(), concs[i-1].Unit()); err != nil {
				return concs, err
			}
		}
		sorted = append(sorted, concs[i])
	}

	if len(descending) > 0 {
		if descending[0] == DescendingOrder {
			sort.Sort(sort.Reverse(concentrationSet(sorted)))
			return
		} else {
			err = fmt.Errorf("unrecognised sort option %s only %s is a valid option at present", descending, DescendingOrder)
		}
	}

	sort.Sort(concentrationSet(sorted))

	return
}

// MinConcentration returns the lowest concentration value from a set of concentration values.
// An error will be returned if no values are specified or the base units of any of the concentrations are incompatible,
// e.g. units of X and g/l would not be compatible.
func MinConcentration(concs []Concentration) (min Concentration, err error) {
	sorted, err := SortConcentrations(concs)
	if err != nil {
		err = fmt.Errorf("Cannot return minimum concentration: %s", err.Error())
		return
	}
	return sorted[0], nil
}

// MaxConcentration returns the highest concentration value from a set of concentration values.
// An error will be returned if no values are specified or the base units of any of the concentrations are incompatible,
// e.g. units of X and g/l would not be compatible.
func MaxConcentration(concs []Concentration) (max Concentration, err error) {
	sorted, err := SortConcentrations(concs)
	if err != nil {
		err = fmt.Errorf("Cannot return maximum concentration: %s", err.Error())
		return
	}
	return sorted[len(sorted)-1], nil
}

// sameUnit returns an error if we cannot convert between unitA and unitB
func sameUnit(unitA, unitB PrefixedUnit) error {
	if a, ok := unitA.(*Unit); !ok {
		return errors.Errorf("unsupported PrefixedUnit type %T", unitA)
	} else if b, ok := unitB.(*Unit); !ok {
		return errors.Errorf("unsupported PrefixedUnit type %T", unitB)
	} else {
		_, err := a.getConversionFactor(b)
		return err
	}
}

type concentrationSet []Concentration

func (cs concentrationSet) Len() int {
	return len(cs)
}

func (cs concentrationSet) Swap(i, j int) {
	cs[i], cs[j] = cs[j], cs[i]
}

func (cs concentrationSet) Less(i, j int) bool {
	return cs[i].SIValue() < cs[j].SIValue()
}

// SortVolumes sorts a set of Volume values.
// An error will be returned if no values are specified or the base units of any of the volumes are incompatible,
// An optional second argument DescendingOrder can be specified to indicate that the sorting will be in descending order.
// By default sorting will be in ascending order.
func SortVolumes(volumes []Volume, descending ...SortOption) (sorted []Volume, err error) {
	if len(volumes) == 0 {
		return volumes, noValues
	}
	for i := range volumes {
		if i > 0 {
			if err = sameUnit(volumes[i].Unit(), volumes[i-1].Unit()); err != nil {
				return volumes, err
			}
		}
		sorted = append(sorted, volumes[i])
	}

	if len(descending) > 0 {
		if descending[0] == DescendingOrder {
			sort.Sort(sort.Reverse(volumeSet(sorted)))
			return
		} else {
			err = fmt.Errorf("unrecognised sort option %s only %s is a valid option at present", descending, DescendingOrder)
		}
	}

	sort.Sort(volumeSet(sorted))

	return
}

// MinVolume returns the lowest Volume value from a set of volume values.
// An error will be returned if no values are specified or the base units of any of the volumes are incompatible,
func MinVolume(volumes []Volume) (min Volume, err error) {
	sorted, err := SortVolumes(volumes)
	if err != nil {
		err = fmt.Errorf("Cannot return minimum concentration: %s", err.Error())
		return
	}
	return sorted[0], nil
}

// MaxVolume returns the highest Volume value from a set of volume values.
// An error will be returned if no values are specified or the base units of any of the volumes are incompatible,
func MaxVolume(volumes []Volume) (max Volume, err error) {
	sorted, err := SortVolumes(volumes)
	if err != nil {
		err = fmt.Errorf("Cannot return maximum volume: %s", err.Error())
		return
	}
	return sorted[len(sorted)-1], nil
}

type volumeSet []Volume

func (cs volumeSet) Len() int {
	return len(cs)
}

func (cs volumeSet) Swap(i, j int) {
	cs[i], cs[j] = cs[j], cs[i]
}

func (cs volumeSet) Less(i, j int) bool {
	return cs[i].SIValue() < cs[j].SIValue()
}
