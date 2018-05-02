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

package text

import (
	"github.com/antha-lang/antha/antha/anthalib/wutil/text"
)

// PrettyPrint will return an indented string of an object.
/* e.g.
[	{
		"Components": {
			"Glucose (g/L)": "150 g/l"
		},
		"Volume": "351 ul"
	},
	{
		"Components": {
			"Magnesium Sulphate (mM)": "50 mM/l"
		},
		"Volume": "210 ul"
	},
]
*/
var PrettyPrint = text.PrettyPrint
