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
	"strings"

	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
)

// if the component name contains a concentration the concentration name will be normalised
// e.g. 10ng/ul glucose will be normalised to 10 mg/l glucose or 10mM glucose to 10 mM/l glucose or 10mM/l glucose to 10 mM/l glucose or glucose 10mM/l to 10 mM/l glucose
func NormaliseName(name string) (normalised string) {

	if strings.Contains(name, wtype.MIXDELIMITER) {
		return name
	}

	containsConc, conc, nameonly := wunit.ParseConcentration(name)

	if containsConc {
		return conc.ToString() + " " + nameonly
	} else {
		return nameonly
	}
}
