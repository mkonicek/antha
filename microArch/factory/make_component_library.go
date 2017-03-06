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
	"sort"

	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/image"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
)

type alreadyAdded struct {
	Name string
}

func (a *alreadyAdded) Error() string {
	return "component " + a.Name + " already added"
}

type notFound struct {
	Name string
}

func (a *notFound) Error() string {
	return "component " + a.Name + " not found"
}

func makeComponentLibrary() (map[string]*wtype.LHComponent, error) {
	var components []*wtype.LHComponent

	add := func(name string, typ wtype.LiquidType, smax float64) {
		c := wtype.NewLHComponent()
		c.CName = name
		c.Type = typ
		c.Smax = smax
		components = append(components, c)
	}

	add("water", wtype.LTWater, 9999)
	add("multiwater", wtype.LTMultiWater, 9999)
	add("Overlay", wtype.LTPLATEOUT, 9999)
	add("PEG", wtype.LTPEG, 9999)
	add("protoplasts", wtype.LTProtoplasts, 9999)
	add("fluorescein", wtype.LTWater, 9999)
	add("ethanol", wtype.LTWater, 9999)
	add("whiteFabricDye", wtype.LTGlycerol, 9999)
	add("blackFabricDye", wtype.LTGlycerol, 9999)
	add("Some component in factory", wtype.LTWater, 9999)
	add("neb5compcells", wtype.LTCulture, 1.0)
	add("mediaonculture", wtype.LTNeedToMix, 1.0)
	add("10x_M9Salts", wtype.LTWater, 9999)
	add("100x_MEMVitamins", wtype.LTWater, 9999)
	add("Yeast extract", wtype.LTWater, 9999)
	add("Tryptone", wtype.LTWater, 9999)
	add("Glycerol", wtype.LTPostMix, 9999)
	add("culture", wtype.LTCulture, 9999)
	// the pubchem name for tartrazine
	add("Acid yellow 23", wtype.LTWater, 9999)
	add("tartrazine", wtype.LTWater, 9999)
	add("tartrazinePostMix", wtype.LTPostMix, 9999)
	add("tartrazineNeedtoMix", wtype.LTNeedToMix, 9999)
	add("tartrazine_DNA", wtype.LTDNA, 9999)
	add("tartrazine_Glycerol", wtype.LTGlycerol, 9999)
	add("Yellow_ink", wtype.LTPAINT, 9999)
	add("Cyan", wtype.LTPAINT, 9999)
	add("Magenta", wtype.LTPAINT, 9999)
	add("transparent", wtype.LTWater, 9999)
	add("Black", wtype.LTPAINT, 9999)
	add("Paint", wtype.LTPostMix, 9999)
	add("yellow", wtype.LTWater, 9999)
	add("blue", wtype.LTWater, 9999)
	add("darkblue", wtype.LTWater, 9999)
	add("grey", wtype.LTWater, 9999)
	add("green", wtype.LTWater, 9999)
	add("red", wtype.LTWater, 9999)
	add("white", wtype.LTWater, 9999)
	add("black", wtype.LTWater, 9999)
	add("purple", wtype.LTWater, 9999)
	add("pink", wtype.LTWater, 9999)
	add("orange", wtype.LTWater, 9999)
	add("DNAsolution", wtype.LTDNA, 1.0)
	add("1kb DNA Ladder", wtype.LTDNA, 10.0)
	add("restrictionenzyme", wtype.LTGlycerol, 1.0)
	add("bsa", wtype.LTWater, 100)
	add("dna_part", wtype.LTDNA, 1.0)
	add("dna", wtype.LTDNA, 1.0)
	add("SapI", wtype.LTGlycerol, 1.0)
	add("T4Ligase", wtype.LTGlycerol, 1.0)
	add("EcoRI", wtype.LTGlycerol, 1.0)
	add("EnzMastermix: 1/2 SapI; 1/2 T4 Ligase", wtype.LTGlycerol, 1.0)
	add("TypeIIsbuffer: 2/11 10xCutsmart; 1/11 1mM ATP; 8/11 Water", wtype.LTWater, 9999)
	add("CutsmartBuffer", wtype.LTWater, 1.0)
	add("ATP", wtype.LTWater, 5.0)
	add("mastermix_sapI", wtype.LTWater, 1.0)
	add("standard_cloning_vector_mark_1", wtype.LTDNA, 1.0)
	add("Q5Polymerase", wtype.LTGlycerol, 1.0)
	add("GoTaq_ green 2x mastermix", wtype.LTGlycerol, 9999.0)
	add("DMSO", wtype.LTWater, 1.0)
	add("pET_GFP", wtype.LTWater, 1.0)
	add("HC", wtype.LTWater, 1.0)
	add("GCenhancer", wtype.LTWater, 9999.0)
	add("Q5buffer", wtype.LTWater, 1.0)
	add("Q5mastermix", wtype.LTWater, 1.0)
	add("PrimerFw", wtype.LTDNA, 1.0)
	add("PrimerRev", wtype.LTDNA, 1.0)
	add("template_part", wtype.LTDNA, 1.0)
	add("DNTPs", wtype.LTWater, 1.0)
	add("ProteinMarker", wtype.LTProtein, 1.0)
	add("ProteinFraction", wtype.LTProtein, 1.0)
	add("EColiLysate", wtype.LTProtein, 1.0)
	add("SDSbuffer", wtype.LTDetergent, 1.0)
	add("Load", wtype.LTload, 1.0)
	add("LB", wtype.LTWater, 1.0)
	add("TB", wtype.LTWater, 1.0)
	add("Kanamycin", wtype.LTWater, 1.0)
	add("Glucose", wtype.LTPostMix, 1.0)
	add("IPTG", wtype.LTPostMix, 1.0)
	add("Lactose", wtype.LTWater, 1.0)
	add("colony", wtype.LTCOLONY, 1.0)
	add("LB_autoinduction_Amp", wtype.LTWater, 1.0)
	add("LB_Kan", wtype.LTWater, 1.0)
	add("Apramycin", wtype.LTWater, 1.0)
	add("Agar", wtype.LTWater, 1.0)
	add("X-glc", wtype.LTWater, 1.0)
	add("X-Glucuro", wtype.LTWater, 1.0)
	add("BaseGrowthMedium", wtype.LTWater, 1.0)
	add("SterileWater", wtype.LTWater, 1.0)
	add("100mMPhosphate", wtype.LTWater, 1.0)
	add("100g/LGlucose", wtype.LTWater, 1.0)
	add("10g/LGlucose", wtype.LTWater, 1.0)
	add("1g/LGlucose", wtype.LTWater, 1.0)
	add("0.1g/Lglucose", wtype.LTWater, 1.0)
	add("100g/Lglycerol", wtype.LTWater, 1.0)
	add("10g/Lglycerol", wtype.LTWater, 1.0)
	add("1g/Lglycerol", wtype.LTWater, 1.0)
	add("0.1g/Lglycerol", wtype.LTWater, 1.0)
	add("100g/Lpeptone", wtype.LTWater, 1.0)
	add("100g/LYeastExtract", wtype.LTWater, 1.0)
	add("10g/LYeastExtract", wtype.LTWater, 1.0)
	add("100g/L Glucose", wtype.LTWater, 1.0)
	add("10g/L Glucose", wtype.LTWater, 1.0)
	add("1g/L Glucose", wtype.LTWater, 1.0)
	add("0.1g/L glucose", wtype.LTWater, 1.0)
	add("100g/L glycerol", wtype.LTWater, 1.0)
	add("10g/L glycerol", wtype.LTWater, 1.0)
	add("1g/L glycerol", wtype.LTWater, 1.0)
	add("0.1g/L glycerol", wtype.LTWater, 1.0)
	add("100g/L peptone", wtype.LTWater, 1.0)
	add("100g/L YeastExtract", wtype.LTWater, 1.0)
	add("10g/L YeastExtract", wtype.LTWater, 1.0)
	add("1000ng/ml ATC", wtype.LTWater, 1.0)
	add("ATC", wtype.LTWater, 1.0)
	add("C6", wtype.LTWater, 1.0)
	add("C12", wtype.LTWater, 1.0)
	add("250uM C6", wtype.LTWater, 1.0)
	add("25uM C6", wtype.LTWater, 1.0)
	add("2.5uM C6", wtype.LTWater, 1.0)
	add("0.25uM C6", wtype.LTWater, 1.0)
	add("0.025uM C6", wtype.LTWater, 1.0)
	add("250g/L C6", wtype.LTWater, 1.0)
	add("25g/L C6", wtype.LTWater, 1.0)
	add("2.5g/L C6", wtype.LTWater, 1.0)
	add("0.25g/L C6", wtype.LTWater, 1.0)
	add("0.025g/L C6", wtype.LTWater, 1.0)
	add("IPTG 1mM", wtype.LTWater, 1.0)
	add("Glucose 100g/L", wtype.LTWater, 1.0)
	add("Glucose 1g/L", wtype.LTWater, 1.0)
	add("Glycerol 100g/L", wtype.LTWater, 1.0)
	add("M9", wtype.LTWater, 1.0)
	add("HYYest412", wtype.LTNSrc, 1.0)
	add("HYYest503", wtype.LTNSrc, 1.0)
	add("HYYest504", wtype.LTNSrc, 1.0)
	add("PeaPeptone", wtype.LTNSrc, 1.0)
	add("WheatPeptone", wtype.LTNSrc, 1.0)
	add("VegPeptone", wtype.LTNSrc, 1.0)
	add("SoyPeptone", wtype.LTNSrc, 1.0)
	add("VegExtract", wtype.LTNSrc, 1.0)
	add("CSL", wtype.LTNSrc, 1.0)
	add("NH42SO4", wtype.LTNSrc, 1.0)
	add("Gluc", wtype.LTNSrc, 1.0)
	add("Suc", wtype.LTNSrc, 1.0)
	add("Fruc", wtype.LTNSrc, 1.0)
	add("Malt", wtype.LTNSrc, 1.0)
	add("water2", wtype.LTNSrc, 1.0)

	// protein paintbox
	for _, value := range image.ProteinPaintboxmap {
		add(value, wtype.LTPostMix, 1.0)
	}

	cmap := make(map[string]*wtype.LHComponent)
	for _, c := range components {
		if _, seen := cmap[c.CName]; seen {
			return nil, &alreadyAdded{Name: c.CName}
		}
		cmap[c.CName] = c
	}

	return cmap, nil
}

type componentLibrary struct {
	lib map[string]*wtype.LHComponent
}

var defaultComponentLibrary *componentLibrary

func init() {
	lib, err := makeComponentLibrary()
	if err != nil {
		panic(err)
	}

	defaultComponentLibrary = &componentLibrary{
		lib: lib,
	}
}

func (i *componentLibrary) GetComponentByType(typ string) *wtype.LHComponent {
	c, ok := i.lib[typ]
	if !ok {
		panic(&notFound{Name: typ})
	}
	return c.Dup()
}

func (i *componentLibrary) GetComponent(typ string) (*wtype.LHComponent, error) {
	c, ok := i.lib[typ]
	if !ok {
		return nil, &notFound{Name: typ}
	}
	return c.Dup(), nil
}

func ComponentInFactory(typ string) bool {
	_, ok := defaultComponentLibrary.lib[typ]
	return ok
}

func GetComponents() []*wtype.LHComponent {
	var comps []*wtype.LHComponent
	for _, c := range defaultComponentLibrary.lib {
		comps = append(comps, c)
	}

	return wtype.CopyComponentArray(comps)
}

// TODO: deprecate
func GetComponentList() []string {
	comps := GetComponents()
	var names []string
	for _, c := range comps {
		names = append(names, c.CName)
	}

	sort.Strings(names)

	return names
}

func GetComponentByType(typ string) *wtype.LHComponent {
	return defaultComponentLibrary.GetComponentByType(typ)
}
