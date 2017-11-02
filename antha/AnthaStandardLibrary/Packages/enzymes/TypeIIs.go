// antha/AnthaStandardLibrary/Packages/enzymes/TypeIIs.go: Part of the Antha language
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

// Package for working with enzymes; in particular restriction enzymes
package enzymes

import "github.com/antha-lang/antha/antha/anthalib/wtype"

var sapI = wtype.RestrictionEnzyme{"GCTCTTC", 3, "SapI", "", 1, 4, "", []string{"N"}, []int{91, 1109, 1919, 1920}, "TypeIIs"}
var isoschizomers = []string{"BspQI", "LguI", "PciSI", "VpaK32I"}

var SapIenz = wtype.TypeIIs{sapI, "SapI", isoschizomers, 1, 4}

var bsaI = wtype.RestrictionEnzyme{"GGTCTC", 4, "BsaI", "Eco31I", 1, 5, "?(5)", []string{"N"}, []int{814, 1109, 1912, 1995, 1996}, "TypeIIs"}

var BsaIenz = wtype.TypeIIs{bsaI, "BsaI", []string{"none"}, 1, 5}

var bpiI = wtype.RestrictionEnzyme{"GAAGAC", 4, "BpiI", "BbvII", 2, 6, "", []string{"B"}, []int{718}, "TypeIIs"}

var BpiIenz = wtype.TypeIIs{bpiI, "BpiI", []string{"BbvII", "BbsI", "BpuAI", "BSTV2I"}, 2, 6}

var TypeIIsEnzymeproperties = map[string]wtype.TypeIIs{
	"SAPI": SapIenz,
	"BSAI": BsaIenz,
	"BPII": BpiIenz,
}
