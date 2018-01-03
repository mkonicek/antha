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

package enzymes

import "github.com/antha-lang/antha/antha/anthalib/wtype"

var sapI = wtype.RestrictionEnzyme{
	Enzyme: wtype.Enzyme{
		Nm: "SapI",
	},
	RecognitionSequence:               "GCTCTTC",
	EndLength:                         3,
	Topstrand3primedistancefromend:    1,
	Bottomstrand5primedistancefromend: 4,
	CommercialSource:                  []string{"N"},
	References:                        []int{91, 1109, 1919, 1920},
	Class:                             "TypeIIs",
	Isoschizomers:                     []string{"BspQI", "LguI", "PciSI", "VpaK32I"},
}

var bsaI = wtype.RestrictionEnzyme{
	Enzyme: wtype.Enzyme{
		Nm: "BsaI",
	},
	RecognitionSequence:               "GGTCTC",
	EndLength:                         4,
	Topstrand3primedistancefromend:    1,
	Bottomstrand5primedistancefromend: 5,
	MethylationSite:                   "?(5)",
	CommercialSource:                  []string{"N"},
	References:                        []int{814, 1109, 1912, 1995, 1996},
	Class:                             "TypeIIs",
	Isoschizomers:                     []string{""},
}

var bpiI = wtype.RestrictionEnzyme{
	Enzyme: wtype.Enzyme{
		Nm: "BpiI",
	},
	Prototype:                         "BbvII",
	RecognitionSequence:               "GAAGAC",
	EndLength:                         4,
	Topstrand3primedistancefromend:    2,
	Bottomstrand5primedistancefromend: 6,
	MethylationSite:                   "",
	CommercialSource:                  []string{"B"},
	References:                        []int{718},
	Class:                             "TypeIIs",
	Isoschizomers:                     []string{"BbvII", "BbsI", "BpuAI", "BSTV2I"},
}

// Example TypeIIs enzymes.
var (
	// SapI is a TypeIIs enzyme.
	SapI = wtype.TypeIIs{RestrictionEnzyme: sapI}
	// BsaI is a TypeIIs enzyme
	BsaI = wtype.TypeIIs{RestrictionEnzyme: bsaI}
	// BpiI is a TypeIIs enzyme
	BpiI = wtype.TypeIIs{RestrictionEnzyme: bpiI}
)

// TypeIIsEnzymeproperties carries a map of example TypeIIs enzymes.
var TypeIIsEnzymeproperties = map[string]wtype.TypeIIs{
	"SAPI": SapI,
	"BSAI": BsaI,
	"BPII": BpiI,
}
