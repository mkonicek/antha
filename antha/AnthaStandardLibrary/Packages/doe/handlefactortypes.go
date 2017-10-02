// Part of the Antha language
// Copyright (C) 2015 The Antha authors. All rights reserved.
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

// Package for facilitating DOE methodology in antha
package doe

import (
	"context"
	"fmt"
	"strconv"

	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/antha/anthalib/wutil"
	"github.com/antha-lang/antha/inventory"
)

// parses a factor name and value and returns an antha concentration.
// If the value cannot be converted to a valid concentration an error is returned.
// If the header contains a valid concentration unit a number can be specified as the value.
func HandleConcFactor(header string, value interface{}) (anthaConc wunit.Concentration, err error) {

	var floatValue float64
	var floatFound bool

	rawconcfloat, found := value.(float64)

	if found {
		floatValue = rawconcfloat
		floatFound = true
	} else {
		rawconcstring, found := value.(string)
		var floatParseErr error
		if floatValue, floatParseErr = strconv.ParseFloat(rawconcstring, 64); found && floatParseErr == nil {
			floatFound = true
		}
	}

	if floatFound {

		// handle floating point imprecision
		floatValue, err = wutil.Roundto(floatValue, 6)

		if err != nil {
			return anthaConc, err
		}
		containsconc, conc, _ := wunit.ParseConcentration(header)

		if containsconc {

			concunit := conc.Unit().PrefixedSymbol()

			anthaConc = wunit.NewConcentration(floatValue, concunit)
		} else {
			err = fmt.Errorf("No valid conc found in component %s so can't assign a concentration unit to value", header)
			return anthaConc, err
		}

	} else if rawconcstring, found := value.(string); found {

		containsconc, conc, _ := wunit.ParseConcentration(rawconcstring)

		if containsconc {
			anthaConc = conc
		} else {
			err = fmt.Errorf("No valid conc found in %s", rawconcstring)
			return anthaConc, err
		}

		// if float use conc unit from header component
	} else {
		err = fmt.Errorf("problem with type of %T expected string or float", value)
		return anthaConc, err
	}

	return
}

// parses a factor name and value and returns an antha Volume.
// If the value cannot be converted to a valid Volume an error is returned.
func HandleVolumeFactor(header string, value interface{}) (anthaVolume wunit.Volume, err error) {

	if rawVolString, found := value.(string); found {

		vol, err := wunit.ParseVolume(rawVolString)

		if err == nil {
			anthaVolume = vol
		} else {
			err = fmt.Errorf("No valid Volume found in ", rawVolString)
			return anthaVolume, err
		}

		// if float use vol unit from header component
	} else if rawVolFloat, found := value.(float64); found {

		// handle floating point imprecision
		rawVolFloat, err = wutil.Roundto(rawVolFloat, 6)

		if err != nil {
			return anthaVolume, err
		}
		vol, err := wunit.ParseVolume(header)

		if err == nil {

			volUnit := vol.Unit().PrefixedSymbol()

			anthaVolume = wunit.NewVolume(rawVolFloat, volUnit)
		} else {
			err = fmt.Errorf("No valid Volume found in component %s so can't assign a Volume unit to value", header)
			return anthaVolume, err
		}

	} else {
		err = fmt.Errorf("problem with type of ", value, " expected string or float")
		return anthaVolume, err
	}

	return
}

// HandleLHComponentFactory parses a factor name and value and returns an
// LHComponent.
//
// If the value cannot be converted to a valid component an error is returned.
func HandleLHComponentFactor(ctx context.Context, header string, value interface{}) (*wtype.LHComponent, error) {
	str, found := value.(string)
	if !found {
		return nil, fmt.Errorf("value %T is not a string", value)
	}

	component, err := inventory.NewComponent(ctx, str)
	if err == nil {
		return component, nil
	}

	if err == inventory.ErrUnknownType {
		component, err = inventory.NewComponent(ctx, inventory.WaterType)
		component.CName = str
		return component, err
	}

	return nil, err
}

// HandleLHPlateFactor parses a factor name and value and returns an
// LHComponent.
//
// If the value cannot be converted to a valid component an error is returned.
func HandleLHPlateFactor(ctx context.Context, header string, value interface{}) (*wtype.LHPlate, error) {
	str, found := value.(string)
	if !found {
		return nil, fmt.Errorf("value %T is not a string", value)
	}

	return inventory.NewPlate(ctx, str)
}
