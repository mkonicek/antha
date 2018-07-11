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

// utility package for working with solutions
package solutions

import (
	"fmt"

	"github.com/antha-lang/antha/antha/anthalib/wtype"
)

// FindComponentByName is a utility function which looks for a component matching on name only.
// If more than one component present the first component will be returned with no error
// This will ignore concentrations if only a unit is specified.
func FindComponentByName(components []*wtype.Liquid, componentName string) (component *wtype.Liquid, err error) {
	for _, comp := range components {
		if comp.CName == componentName || NormaliseName(componentName) == NormaliseName(comp.CName) {
			return comp, nil
		}
	}
	return component, fmt.Errorf("No component found with name %s in component list", componentName)
}
