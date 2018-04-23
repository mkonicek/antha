// wunit/genericunit.go: Part of the Antha language
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
)

// structure for defining a generic unit
type GenericUnit struct {
	StrName             string
	StrSymbol           string
	FltConversionfactor float64
	StrBaseUnit         string
}

func (gu *GenericUnit) Name() string {
	return gu.StrName
}
func (gu *GenericUnit) Symbol() string {
	return gu.StrSymbol
}
func (gu *GenericUnit) BaseSIConversionFactor() float64 {
	return gu.FltConversionfactor
}
func (gu *GenericUnit) BaseSIUnit() string {
	return gu.StrBaseUnit
}

func (gu *GenericUnit) ToString() string {
	return fmt.Sprintf("Name: %s Symbol: %s Conversion: %-4g BaseUnit: %s", gu.StrName, gu.StrSymbol, gu.FltConversionfactor, gu.StrBaseUnit)
}

// the generic prefixed unit structure
type GenericPrefixedUnit struct {
	GenericUnit
	SPrefix SIPrefix
}

func (gpu *GenericPrefixedUnit) Prefix() SIPrefix {
	return gpu.SPrefix
}

// multiplier to convert to SI base unit... for composites this is the
// ratio of the base units for the dimensions in question e.g. kg/l for concentration
func (gpu *GenericPrefixedUnit) BaseSIConversionFactor() float64 {
	return gpu.Prefix().Value * gpu.GenericUnit.BaseSIConversionFactor()
}

// symbol without prefix
func (gpu *GenericPrefixedUnit) RawSymbol() string {
	return gpu.GenericUnit.Symbol()
}

// symbol with prefix
func (gpu *GenericPrefixedUnit) PrefixedSymbol() string {
	if gpu == nil {
		return ""
	}
	if gpu.SPrefix.Name == "" || gpu.SPrefix.Name == " " {
		return gpu.GenericUnit.Symbol()
	}
	return fmt.Sprintf("%s%s", gpu.SPrefix.Name, gpu.GenericUnit.Symbol())
}

// symbol for unit after conversion to base si unit
func (gpu *GenericPrefixedUnit) BaseSISymbol() string {
	return gpu.GenericUnit.BaseSIUnit()
}

// symbol with prefix
func (gpu *GenericPrefixedUnit) Symbol() string {
	return gpu.PrefixedSymbol()
}

// gives the conversion factor from one prefixed unit to another
func (gpu *GenericPrefixedUnit) ConvertTo(p2 PrefixedUnit) float64 {
	return gpu.BaseSIConversionFactor() / p2.BaseSIConversionFactor()
}

// convenience

func NewGenericPrefixedUnit() *GenericPrefixedUnit {
	gpu := GenericPrefixedUnit{}
	return &gpu
}
