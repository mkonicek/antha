// /anthalib/factory/make_component_library.go: Part of the Antha language
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

package factory

import (
	"errors"
	"fmt"

	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/image"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/microArch/logger"
)

func makeComponentLibrary() map[string]*wtype.LHComponent {
	//	matter := wtype.MakeMatterLib()

	cmap := make(map[string]*wtype.LHComponent)

	A := wtype.NewLHComponent()
	A.CName = "water"
	A.Type = wtype.LTWater
	A.Smax = 9999
	cmap[A.CName] = A

	A = wtype.NewLHComponent()
	A.CName = "PEG"
	A.Type = wtype.LTPEG
	A.Smax = 9999
	cmap[A.CName] = A

	A = wtype.NewLHComponent()
	A.CName = "protoplasts"
	A.Type = wtype.LTProtoplasts
	A.Smax = 9999
	cmap[A.CName] = A

	A = wtype.NewLHComponent()
	//	//A.GenericMatter = matter[wtype.LTWater]
	A.CName = "fluorescein"
	A.Type = wtype.LTWater
	A.Smax = 9999
	cmap[A.CName] = A

	A = wtype.NewLHComponent()
	A.CName = "ethanol"
	A.Type = wtype.LTWater
	A.Smax = 9999
	cmap[A.CName] = A

	A = wtype.NewLHComponent()
	//A.GenericMatter = matter[wtype.LTWater]
	A.CName = "whiteFabricDye"
	A.Type = wtype.LTGlycerol
	A.Smax = 9999
	cmap[A.CName] = A

	A = wtype.NewLHComponent()
	//A.GenericMatter = matter[wtype.LTWater]
	A.CName = "blackFabricDye"
	A.Type = wtype.LTGlycerol
	A.Smax = 9999
	cmap[A.CName] = A

	A = wtype.NewLHComponent()
	//	//A.GenericMatter = matter[wtype.LTWater]
	A.CName = "Some component in factory"
	A.Type = wtype.LTWater
	A.Smax = 9999
	cmap[A.CName] = A

	A = wtype.NewLHComponent()
	//A.GenericMatter = matter[wtype.LTWater]
	A.CName = "neb5compcells"
	A.Type = wtype.LTCulture
	A.Smax = 1.0
	cmap[A.CName] = A

	A = wtype.NewLHComponent()
	//A.GenericMatter = matter[wtype.LTWater]
	A.CName = "mediaonculture"
	A.Type = wtype.LTNeedToMix
	A.Smax = 1.0
	cmap[A.CName] = A

	A = wtype.NewLHComponent()
	//A.GenericMatter = matter[wtype.LTWater]
	A.CName = "10x_M9Salts"
	A.Type = wtype.LTWater
	A.Smax = 9999
	cmap[A.CName] = A

	A = wtype.NewLHComponent()
	//A.GenericMatter = matter[wtype.LTWater]
	A.CName = "100x_MEMVitamins"
	A.Type = wtype.LTWater
	A.Smax = 9999
	cmap[A.CName] = A

	A = wtype.NewLHComponent()
	//A.GenericMatter = matter[wtype.LTWater]
	A.CName = "Yeast extract"
	A.Type = wtype.LTWater
	A.Smax = 9999
	cmap[A.CName] = A

	A = wtype.NewLHComponent()
	//A.GenericMatter = matter[wtype.LTWater]
	A.CName = "Tryptone"
	A.Type = wtype.LTWater
	A.Smax = 9999
	cmap[A.CName] = A

	A = wtype.NewLHComponent()
	//A.GenericMatter = matter[wtype.LTWater]
	A.CName = "Glycerol"
	A.Type = wtype.LTPostMix
	A.Smax = 9999
	cmap[A.CName] = A

	A = wtype.NewLHComponent()
	//A.GenericMatter = matter[wtype.LTWater]
	A.CName = "culture"
	A.Type = wtype.LTCulture
	A.Smax = 9999
	cmap[A.CName] = A

	A = wtype.NewLHComponent()
	A.CName = "Acid yellow 23" // the pubchem name for tartrazine
	A.Type = wtype.LTWater
	A.Smax = 9999
	cmap[A.CName] = A

	A = wtype.NewLHComponent()
	A.CName = "tartrazine"
	A.Type = wtype.LTWater
	A.Smax = 9999
	cmap[A.CName] = A

	A = wtype.NewLHComponent()
	A.CName = "tartrazinePostMix"
	A.Type = wtype.LTPostMix
	A.Smax = 9999
	cmap[A.CName] = A

	A = wtype.NewLHComponent()
	A.CName = "tartrazineNeedtoMix"
	A.Type = wtype.LTNeedToMix
	A.Smax = 9999
	cmap[A.CName] = A

	A = wtype.NewLHComponent()
	A.CName = "tartrazine_DNA"
	A.Type = wtype.LTDNA
	A.Smax = 9999
	cmap[A.CName] = A

	A = wtype.NewLHComponent()
	A.CName = "tartrazine_Glycerol"
	A.Type = wtype.LTGlycerol
	A.Smax = 9999
	cmap[A.CName] = A

	A = wtype.NewLHComponent()
	//A.GenericMatter = matter[wtype.LTWater]
	A.CName = "Yellow_ink"
	A.Type = wtype.LTPAINT
	A.Smax = 9999
	cmap[A.CName] = A

	A = wtype.NewLHComponent()
	//A.GenericMatter = matter[wtype.LTWater]
	A.CName = "Cyan"
	A.Type = wtype.LTPAINT
	A.Smax = 9999
	cmap[A.CName] = A

	A = wtype.NewLHComponent()
	//A.GenericMatter = matter[wtype.LTWater]
	A.CName = "Magenta"
	A.Type = wtype.LTPAINT
	A.Smax = 9999
	cmap[A.CName] = A

	A = wtype.NewLHComponent()
	//A.GenericMatter = matter[wtype.LTWater]
	A.CName = "transparent"
	A.Type = wtype.LTWater
	A.Smax = 9999
	cmap[A.CName] = A

	A = wtype.NewLHComponent()
	//A.GenericMatter = matter[wtype.LTWater]
	A.CName = "Black"
	A.Type = wtype.LTPAINT
	A.Smax = 9999
	cmap[A.CName] = A

	A = wtype.NewLHComponent()
	//A.GenericMatter = matter[wtype.LTWater]
	A.CName = "Paint"
	A.Type = wtype.LTPostMix
	A.Smax = 9999
	cmap[A.CName] = A

	A = wtype.NewLHComponent()
	//A.GenericMatter = matter[wtype.LTWater]
	A.CName = "yellow"
	A.Type = wtype.LTWater
	A.Smax = 9999
	cmap[A.CName] = A

	A = wtype.NewLHComponent()
	//A.GenericMatter = matter[wtype.LTWater]
	A.CName = "blue"
	A.Type = wtype.LTWater
	A.Smax = 9999
	cmap[A.CName] = A

	A = wtype.NewLHComponent()
	//A.GenericMatter = matter[wtype.LTWater]
	A.CName = "darkblue"
	A.Type = wtype.LTWater
	A.Smax = 9999
	cmap[A.CName] = A

	A = wtype.NewLHComponent()
	//A.GenericMatter = matter[wtype.LTWater]
	A.CName = "grey"
	A.Type = wtype.LTWater
	A.Smax = 9999
	cmap[A.CName] = A

	A = wtype.NewLHComponent()
	//A.GenericMatter = matter[wtype.LTWater]
	A.CName = "green"
	A.Type = wtype.LTWater
	A.Smax = 9999
	cmap[A.CName] = A

	A = wtype.NewLHComponent()
	//A.GenericMatter = matter[wtype.LTWater]
	A.CName = "red"
	A.Type = wtype.LTWater
	A.Smax = 9999
	cmap[A.CName] = A

	A = wtype.NewLHComponent()
	//A.GenericMatter = matter[wtype.LTWater]
	A.CName = "white"
	A.Type = wtype.LTWater
	A.Smax = 9999
	cmap[A.CName] = A

	A = wtype.NewLHComponent()
	//A.GenericMatter = matter[wtype.LTWater]
	A.CName = "black"
	A.Type = wtype.LTWater
	A.Smax = 9999
	cmap[A.CName] = A

	A = wtype.NewLHComponent()
	//A.GenericMatter = matter[wtype.LTWater]
	A.CName = "purple"
	A.Type = wtype.LTWater
	A.Smax = 9999
	cmap[A.CName] = A

	A = wtype.NewLHComponent()
	//A.GenericMatter = matter[wtype.LTWater]
	A.CName = "pink"
	A.Type = wtype.LTWater
	A.Smax = 9999
	cmap[A.CName] = A

	A = wtype.NewLHComponent()
	//A.GenericMatter = matter[wtype.LTWater]
	A.CName = "orange"
	A.Type = wtype.LTWater
	A.Smax = 9999
	cmap[A.CName] = A

	A = wtype.NewLHComponent()
	//A.GenericMatter = matter[wtype.LTWater]
	A.CName = "DNAsolution"
	A.Type = wtype.LTDNA
	A.Smax = 1.0
	cmap[A.CName] = A

	A = wtype.NewLHComponent()
	//A.GenericMatter = matter[wtype.LTWater]
	A.CName = "1kb DNA Ladder"
	A.Type = wtype.LTDNA
	A.Smax = 10.0
	cmap[A.CName] = A

	A = wtype.NewLHComponent()
	//A.GenericMatter = matter["glycerol"]
	A.CName = "restrictionenzyme"
	A.Type = wtype.LTGlycerol
	A.Smax = 1.0
	cmap[A.CName] = A

	A = wtype.NewLHComponent()
	//A.GenericMatter = matter[wtype.LTWater]
	A.CName = "bsa"
	A.Type = wtype.LTWater
	A.Smax = 100
	cmap[A.CName] = A

	A = wtype.NewLHComponent()
	//A.GenericMatter = matter[wtype.LTWater]
	A.CName = "dna_part"
	A.Type = wtype.LTDNA
	A.Smax = 1.0
	cmap[A.CName] = A

	A = wtype.NewLHComponent()
	//A.GenericMatter = matter[wtype.LTWater]
	A.CName = "dna"
	A.Type = wtype.LTDNA
	A.Smax = 1.0
	cmap[A.CName] = A

	A = wtype.NewLHComponent()
	//A.GenericMatter = matter["glycerol"]
	A.CName = "SapI"
	A.Type = wtype.LTGlycerol
	A.Smax = 1.0
	cmap[A.CName] = A

	A = wtype.NewLHComponent()
	//A.GenericMatter = matter["glycerol"]
	A.CName = "T4Ligase"
	A.Type = wtype.LTGlycerol
	A.Smax = 1.0
	cmap[A.CName] = A

	A = wtype.NewLHComponent()
	//A.GenericMatter = matter["glycerol"]
	A.CName = "EcoRI"
	A.Type = wtype.LTGlycerol
	A.Smax = 1.0
	cmap[A.CName] = A

	A = wtype.NewLHComponent()
	//A.GenericMatter = matter["glycerol"]
	A.CName = "EnzMastermix: 1/2 SapI; 1/2 T4 Ligase"
	A.Type = wtype.LTGlycerol
	A.Smax = 1.0
	cmap[A.CName] = A

	A = wtype.NewLHComponent()
	//A.GenericMatter = matter[wtype.LTWater]
	A.CName = "TypeIIsbuffer: 2/11 10xCutsmart; 1/11 1mM ATP; 8/11 Water"
	A.Type = wtype.LTWater
	A.Smax = 9999
	cmap[A.CName] = A

	A = wtype.NewLHComponent()
	//A.GenericMatter = matter[wtype.LTWater]
	A.CName = "CutsmartBuffer"
	A.Type = wtype.LTWater
	A.Smax = 1.0
	cmap[A.CName] = A

	A = wtype.NewLHComponent()
	//A.GenericMatter = matter[wtype.LTWater]
	A.CName = "ATP"
	A.Type = wtype.LTWater
	A.Smax = 5.0
	cmap[A.CName] = A

	A = wtype.NewLHComponent()
	//A.GenericMatter = matter[wtype.LTWater]
	//A.CName = "SapI_Mastermix: 1/5 SapI; 1/5 T4 Ligase; 2/5 Cutsmart; 1/5 1mM ATP"
	A.CName = "mastermix_sapI"
	A.Type = wtype.LTWater
	A.Smax = 1.0
	cmap[A.CName] = A

	A = wtype.NewLHComponent()
	//A.GenericMatter = matter[wtype.LTWater]
	A.CName = "standard_cloning_vector_mark_1"
	A.Type = wtype.LTDNA
	A.Smax = 1.0
	cmap[A.CName] = A

	//

	A = wtype.NewLHComponent()
	//A.GenericMatter = matter[wtype.LTWater]
	A.CName = "standard_cloning_vector_mark_1"
	A.Type = wtype.LTDNA
	A.Smax = 1.0
	cmap[A.CName] = A

	// solutions needed for PCR example:

	A = wtype.NewLHComponent()
	//A.GenericMatter = matter["glycerol"]
	A.CName = "Q5Polymerase"
	A.Type = wtype.LTGlycerol
	A.Smax = 1.0 // not sure if this is correct
	cmap[A.CName] = A

	A = wtype.NewLHComponent()
	//A.GenericMatter = matter["glycerol"]
	A.CName = "GoTaq_ green 2x mastermix"
	A.Type = wtype.LTGlycerol
	A.Smax = 9999.0 // not sure if this is correct
	cmap[A.CName] = A

	A = wtype.NewLHComponent()
	//A.GenericMatter = matter[wtype.LTWater]
	A.CName = "DMSO"
	A.Type = wtype.LTWater
	A.Smax = 1.0 // not sure if this is correct
	cmap[A.CName] = A

	A = wtype.NewLHComponent()
	//A.GenericMatter = matter[wtype.LTWater]
	A.CName = "pET_GFP"
	A.Type = wtype.LTWater
	A.Smax = 1.0 // not sure if this is correct
	cmap[A.CName] = A

	A = wtype.NewLHComponent()
	//A.GenericMatter = matter[wtype.LTWater]
	A.CName = "HC"
	A.Type = wtype.LTWater
	A.Smax = 1.0 // not sure if this is correct
	cmap[A.CName] = A

	A = wtype.NewLHComponent()
	//A.GenericMatter = matter[wtype.LTWater]
	A.CName = "GCenhancer"
	A.Type = wtype.LTWater
	A.Smax = 9999.0 // not sure if this is correct
	cmap[A.CName] = A

	A = wtype.NewLHComponent()
	//A.GenericMatter = matter[wtype.LTWater]
	A.CName = "Q5buffer"
	A.Type = wtype.LTWater
	A.Smax = 1.0 // not sure if this is correct
	cmap[A.CName] = A

	A = wtype.NewLHComponent()
	//A.GenericMatter = matter[wtype.LTWater]
	A.CName = "Q5mastermix"
	A.Type = wtype.LTWater
	A.Smax = 1.0 // not sure if this is correct
	cmap[A.CName] = A

	A = wtype.NewLHComponent()
	//A.GenericMatter = matter[wtype.LTWater]
	A.CName = "PrimerFw"
	A.Type = wtype.LTDNA
	A.Smax = 1.0 // not sure if this is correct
	cmap[A.CName] = A

	A = wtype.NewLHComponent()
	//A.GenericMatter = matter[wtype.LTWater]
	A.CName = "PrimerRev"
	A.Type = wtype.LTDNA
	A.Smax = 1.0 // not sure if this is correct
	cmap[A.CName] = A

	A = wtype.NewLHComponent()
	//A.GenericMatter = matter[wtype.LTWater]
	A.CName = "template_part"
	A.Type = wtype.LTDNA
	A.Smax = 1.0 // not sure if this is correct
	cmap[A.CName] = A

	A = wtype.NewLHComponent()
	//A.GenericMatter = matter[wtype.LTWater]
	A.CName = "DNTPs"
	A.Type = wtype.LTWater
	A.Smax = 1.0 // not sure if this is correct
	cmap[A.CName] = A

	A = wtype.NewLHComponent()
	//A.GenericMatter = matter[wtype.LTWater]
	A.CName = "ProteinMarker"
	A.Type = wtype.LTProtein
	A.Smax = 1.0 //not sure if this is correct
	cmap[A.CName] = A

	A = wtype.NewLHComponent()
	//A.GenericMatter = matter[wtype.LTWater]
	A.CName = "ProteinFraction"
	A.Type = wtype.LTProtein
	A.Smax = 1.0 //still not sure
	cmap[A.CName] = A

	A = wtype.NewLHComponent()
	//A.GenericMatter = matter[wtype.LTWater]
	A.CName = "EColiLysate"
	A.Type = wtype.LTProtein
	A.Smax = 1.0 //not sure what this is!
	cmap[A.CName] = A

	A = wtype.NewLHComponent()
	//A.GenericMatter = matter[wtype.LTWater]
	A.CName = "SDSbuffer"
	A.Type = wtype.LTDetergent
	A.Smax = 1.0 //still not sure....
	cmap[A.CName] = A

	A = wtype.NewLHComponent()
	//A.GenericMatter = matter[wtype.LTWater]
	A.CName = "Load"
	A.Type = wtype.LTload
	A.Smax = 1.0 //still not sure....
	cmap[A.CName] = A

	A = wtype.NewLHComponent()
	//A.GenericMatter = matter[wtype.LTWater]
	A.CName = "LB"
	A.Type = wtype.LTWater
	A.Smax = 1.0 //still not sure....
	cmap[A.CName] = A

	A = wtype.NewLHComponent()
	//A.GenericMatter = matter[wtype.LTWater]
	A.CName = "TB"
	A.Type = wtype.LTWater
	A.Smax = 1.0 //still not sure....
	cmap[A.CName] = A

	A = wtype.NewLHComponent()
	//A.GenericMatter = matter[wtype.LTWater]
	A.CName = "Kanamycin"
	A.Type = wtype.LTWater
	A.Smax = 1.0 //still not sure....
	cmap[A.CName] = A

	A = wtype.NewLHComponent()
	//A.GenericMatter = matter[wtype.LTWater]
	A.CName = "Glucose"
	A.Type = wtype.LTPostMix
	A.Smax = 1.0 //still not sure....
	cmap[A.CName] = A

	A = wtype.NewLHComponent()
	//A.GenericMatter = matter[wtype.LTWater]
	A.CName = "IPTG"
	A.Type = wtype.LTPostMix
	A.Smax = 1.0 //still not sure....
	cmap[A.CName] = A

	A = wtype.NewLHComponent()
	//A.GenericMatter = matter[wtype.LTWater]
	A.CName = "Lactose"
	A.Type = wtype.LTWater
	A.Smax = 1.0 //still not sure....
	cmap[A.CName] = A

	A = wtype.NewLHComponent()
	//A.GenericMatter = matter[wtype.LTWater]
	A.CName = "colony"
	A.Type = wtype.LTCOLONY
	A.Smax = 1.0 //still not sure....
	cmap[A.CName] = A
	A = wtype.NewLHComponent()
	//A.GenericMatter = matter[wtype.LTWater]
	A.CName = "LB_autoinduction_Amp"
	A.Type = wtype.LTWater
	A.Smax = 1.0 //still not sure....
	cmap[A.CName] = A

	A = wtype.NewLHComponent()
	//A.GenericMatter = matter[wtype.LTWater]
	A.CName = "LB_Kan"
	A.Type = wtype.LTWater
	A.Smax = 1.0 //still not sure....
	cmap[A.CName] = A

	A = wtype.NewLHComponent()
	//A.GenericMatter = matter[wtype.LTWater]
	A.CName = "Apramycin"
	A.Type = wtype.LTWater
	A.Smax = 1.0 //still not sure....
	cmap[A.CName] = A

	A = wtype.NewLHComponent()
	//A.GenericMatter = matter[wtype.LTWater]
	A.CName = "Agar"
	A.Type = wtype.LTWater
	A.Smax = 1.0 //still not sure....
	cmap[A.CName] = A

	A = wtype.NewLHComponent()
	//A.GenericMatter = matter[wtype.LTWater]
	A.CName = "X-glc"
	A.Type = wtype.LTWater
	A.Smax = 1.0 //still not sure....
	cmap[A.CName] = A

	A = wtype.NewLHComponent()
	//A.GenericMatter = matter[wtype.LTWater]
	A.CName = "X-Glucuro"
	A.Type = wtype.LTWater
	A.Smax = 1.0 //still not sure....
	cmap[A.CName] = A

	A = wtype.NewLHComponent()
	//A.GenericMatter = matter[wtype.LTWater]
	A.CName = "BaseGrowthMedium"
	A.Type = wtype.LTWater
	A.Smax = 1.0 //still not sure....
	cmap[A.CName] = A

	A = wtype.NewLHComponent()
	//A.GenericMatter = matter[wtype.LTWater]
	A.CName = "SterileWater"
	A.Type = wtype.LTWater
	A.Smax = 1.0 //still not sure....
	cmap[A.CName] = A

	A = wtype.NewLHComponent()
	//A.GenericMatter = matter[wtype.LTWater]
	A.CName = "100mMPhosphate"
	A.Type = wtype.LTWater
	A.Smax = 1.0 //still not sure....
	cmap[A.CName] = A
	A = wtype.NewLHComponent()
	//A.GenericMatter = matter[wtype.LTWater]
	A.CName = "100g/LGlucose"
	A.Type = wtype.LTWater
	A.Smax = 1.0 //still not sure....
	cmap[A.CName] = A
	A = wtype.NewLHComponent()
	//A.GenericMatter = matter[wtype.LTWater]
	A.CName = "10g/LGlucose"
	A.Type = wtype.LTWater
	A.Smax = 1.0 //still not sure....
	cmap[A.CName] = A
	A = wtype.NewLHComponent()
	//A.GenericMatter = matter[wtype.LTWater]
	A.CName = "1g/LGlucose"
	A.Type = wtype.LTWater
	A.Smax = 1.0 //still not sure....
	cmap[A.CName] = A
	A = wtype.NewLHComponent()
	//A.GenericMatter = matter[wtype.LTWater]
	A.CName = "0.1g/Lglucose"
	A.Type = wtype.LTWater
	A.Smax = 1.0 //still not sure....
	cmap[A.CName] = A

	A = wtype.NewLHComponent()
	//A.GenericMatter = matter[wtype.LTWater]
	A.CName = "100g/Lglycerol"
	A.Type = wtype.LTWater
	A.Smax = 1.0 //still not sure....
	cmap[A.CName] = A
	A = wtype.NewLHComponent()
	//A.GenericMatter = matter[wtype.LTWater]
	A.CName = "10g/Lglycerol"
	A.Type = wtype.LTWater
	A.Smax = 1.0 //still not sure....
	cmap[A.CName] = A
	A = wtype.NewLHComponent()
	//A.GenericMatter = matter[wtype.LTWater]
	A.CName = "1g/Lglycerol"
	A.Type = wtype.LTWater
	A.Smax = 1.0 //still not sure....
	cmap[A.CName] = A
	A = wtype.NewLHComponent()
	//A.GenericMatter = matter[wtype.LTWater]
	A.CName = "0.1g/Lglycerol"
	A.Type = wtype.LTWater
	A.Smax = 1.0 //still not sure....
	cmap[A.CName] = A

	A = wtype.NewLHComponent()
	//A.GenericMatter = matter[wtype.LTWater]
	A.CName = "100g/Lpeptone"
	A.Type = wtype.LTWater
	A.Smax = 1.0 //still not sure....
	cmap[A.CName] = A

	A = wtype.NewLHComponent()
	//A.GenericMatter = matter[wtype.LTWater]
	A.CName = "100g/LYeastExtract"
	A.Type = wtype.LTWater
	A.Smax = 1.0 //still not sure....
	cmap[A.CName] = A

	A = wtype.NewLHComponent()
	//A.GenericMatter = matter[wtype.LTWater]
	A.CName = "10g/LYeastExtract"
	A.Type = wtype.LTWater
	A.Smax = 1.0 //still not sure....
	cmap[A.CName] = A

	// protein paintbox

	for key, value := range image.ProteinPaintboxmap {

		_, ok := cmap[value]
		if ok == true {
			alreadyinthere := cmap[value]

			err := fmt.Errorf("attempt to add value", key, "for key", value, "to component factory", cmap, "failed due to duplicate entry", alreadyinthere)
			panic(err)
		} else {

			A = wtype.NewLHComponent()
			//A.GenericMatter = matter[wtype.LTWater]
			A.CName = value
			A.Type = wtype.LTPostMix
			A.Smax = 1.0
			cmap[A.CName] = A

		}
	}
	return cmap
}

type Inventory struct {
	inv map[string]*wtype.LHComponent
}

var defaultInventory *Inventory

func init() {
	defaultInventory = &Inventory{
		inv: makeComponentLibrary(),
	}
}

func AddComponent(name string, c *wtype.LHComponent) error {
	return errors.New("tbd")
}
func GetComponentByType(typ string) *wtype.LHComponent {
	return defaultInventory.GetComponentByType(typ)
}

func (i *Inventory) GetComponentByType(typ string) *wtype.LHComponent {
	c, ok := i.inv[typ]
	if !ok {
		logger.Fatal(fmt.Sprintf("Component %s not found", typ))
		panic(fmt.Errorf("Component %s not found", typ)) //TODO refactor to errors
	}
	return c.Dup()
}

func ComponentInFactory(typ string) bool {
	components := makeComponentLibrary()
	c, ok := components[typ]
	if c == nil || ok == false {
		return false
	}
	if ok {
		return true
	}
	return false
}

func GetComponentList() []string {
	components := makeComponentLibrary()
	kz := make([]string, len(components))
	x := 0
	for name, _ := range components {
		kz[x] = name
		x += 1
	}
	return kz

}
