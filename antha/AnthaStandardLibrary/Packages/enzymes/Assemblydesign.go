// antha/AnthaStandardLibrary/Packages/enzymes/Assemblydesign.go: Part of the Antha language
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

// Package enzymes for working with enzymes; in particular restriction enzymes
package enzymes

import (
	"fmt"
	"sort"
	"strings"

	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/enzymes/lookup"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/search"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/sequences"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
)

/*
not finished
func LengthofPrefixOverlap(seq string, subseq string) (number int, end string) { // add end string


	i:=0;i<len(subseq);i++{
	truncated := subseq[i:]
	// fmt.Println("truncated", truncated)
	if strings.HasPrefix(part.Seq, truncated) == true {
		number = i
		end = "end"
	}
	/*start := subseq[:i]
	// fmt.Println("start", start)
	if strings.HasPrefix(part.Seq, start) == true {
		number = i
		end = "start"
	}
	return number
}
*/

/*
// Checks for duplicate ends in the list of parts to assemble
// This code is completely wrong!!! it needs to digest the fragements first!
func CheckEndCompatibility(fragments []wtype.DNASequence)err error{

	type partEnd struct{
		Parts []string
		End string
	}

	// Check that parts have unique overhangs
	var endMap = make(map[string]partEnd)
	var errs []string
	for _, part := range PartswithOverhangs {
		prefix := wtype.Prefix(part.OverHang(),restrictionenzyme.EndLength)
		suffix := wtype.Prefix(part.Sequence(),restrictionenzyme.EndLength)

		if prefix == suffix {
			errs = append(errs,fmt.Sprintf("5 prime end %s of part %s same as 3 prime end",prefix, part.Name()))
		}

		if end, found := endMap[prefix];found{
			end.Parts = append(end.Parts,part.Name())
			if len(end.Parts)>2{
				errs = append(errs,fmt.Sprintf("5 prime end %s of part %s already found in more than one other part %s ",end.End, part.Name(),strings.Join(end.Parts,";")))
			}
		}else{
			endMap[prefix]= partEnd{Parts:[]string{part.Name()},End: prefix}
		}

		if end, found := endMap[suffix];found{
			end.Parts = append(end.Parts,part.Name())
			if len(end.Parts)>2{
				errs = append(errs,fmt.Sprintf("3 prime end %s of part %s already found in more than one other part %s ",end.End, part.Name(),strings.Join(end.Parts,";")))
			}

		}else{
			endMap[suffix]= partEnd{Parts:[]string{part.Name()},End: suffix}
		}
	}
	if len(errs)> 0{
		err = fmt.Errorf(strings.Join(errs,";"))
	}
	return
}
*/

// VectorEnds returns the 5' and 3' sticky ends found from cutting a vector DNASequence with TypeIIs enzyme.
// If the vector is cut more than once the first two sticky ends will be returned.
// If the vector is not cut "" and "" will be returned.
func VectorEnds(vector wtype.DNASequence, enzyme wtype.TypeIIs) (desiredstickyend5prime string, vector3primestickyend string) {
	// find sticky ends from cutting vector with enzyme

	fragments, stickyends5, _ := TypeIIsdigest(vector, enzyme)

	// add better logic for the scenarios where the vector is cut more than twice or we want to add fragment in either direction
	// picks largest fragment

	for i := 0; i < len(stickyends5)-1; i++ {

		if stickyends5[i] != "" && len(fragments[i]) > 0 {

			vector3primestickyend = stickyends5[i]
			desiredstickyend5prime = stickyends5[i+1]

		}
	}
	return
}

// MakeScarfreeCustomTypeIIsassemblyParts adds typeIIs assembly ends to a set of parts.
// The ends will be added in order to enable correct assembly of the parts in the order specified leaving no scar sequence between parts.
// The ends will be added based on a specified typeIIs enzyme and vector which contains sites for that typeIIs enzyme.
// If the parts already contain typeIIs sites for the specified enzyme the ends will be checked for compatibility.
func MakeScarfreeCustomTypeIIsassemblyParts(parts []wtype.DNASequence, vector wtype.DNASequence, enzyme wtype.TypeIIs) (partswithends []wtype.DNASequence) {

	partswithends = make([]wtype.DNASequence, 0)

	// find sticky ends from cutting vector with enzyme

	fragments, stickyends5, _ := TypeIIsdigest(vector, enzyme)

	//initialise

	desiredstickyend5prime := ""

	vector3primestickyend := ""

	// add better logic for the scenarios where the vector is cut more than twice or we want to add fragment in either direction
	// picks largest fragment

	for i := 0; i < len(stickyends5)-1; i++ {

		if stickyends5[i] != "" && len(fragments[i]) > 0 {

			vector3primestickyend = stickyends5[i]
			desiredstickyend5prime = stickyends5[i+1]

		}
	} // fill in later

	// declare as blank so no end added
	desiredstickyend3prime := ""

	for i := 0; i < len(parts); i++ {
		var partwithends wtype.DNASequence

		if i == (len(parts) - 1) {
			desiredstickyend3prime = vector3primestickyend
		}

		sites, sticky5s, sticky3s := CheckForExistingTypeIISEnds(parts[i], enzyme)

		if sites == 0 {
			partwithends = AddCustomEnds(parts[i], enzyme, desiredstickyend5prime, desiredstickyend3prime)
		} else if sites == 2 && search.InStrings(sticky5s, desiredstickyend5prime) && search.InStrings(sticky3s, desiredstickyend3prime) {
			partwithends = parts[i]
		} else {
			panic(fmt.Sprint("cutting part ", parts[i], " with", enzyme, " results in ", sites, "cut sites with 5 prime fragment overhangs: ", sticky5s, " and 3 prime fragment overhangs: ", sticky3s, ". Wanted: 5prime: ", desiredstickyend5prime, " 3prime: ", desiredstickyend3prime))
		}
		partwithends.Nm = parts[i].Nm

		partswithends = append(partswithends, partwithends)

		desiredstickyend5prime = sequences.Suffix(parts[i].Seq, enzyme.RestrictionEnzyme.EndLength)

	}

	return partswithends
}

// AddL1UAdaptor adds an upstream (5') adaptor for making a level 0 part compatible for Level 1 hierarchical assembly, specifying the desired level 1 class the level0 part should be made into.
// TypeIIs recognition site + spacer + correct overhang in correct orientation will be added according to the correct enzyme at a specified level according to a specified
// assembly standard.
// If reverseOrientation is set to true the adaptor will be added such that the level 1 part will bind in the reverse orientation.
func AddL1UAdaptor(part wtype.DNASequence, assemblyStandard AssemblyStandard, level string, class string, reverseOrientation bool) (newpart wtype.DNASequence, err error) {

	enzyme, err := lookUpEnzyme(assemblyStandard, level)

	if err != nil {
		return newpart, err
	}

	bitToAdd, bitToAdd3, err := lookUpOverhangs(assemblyStandard, level, class)

	if err != nil {
		return newpart, err
	}

	if reverseOrientation {
		bitToAdd = wtype.RevComp(bitToAdd3)
	}

	bitToAdd5prime := MakeOverhang(enzyme, "5prime", bitToAdd, ChooseSpacer(enzyme.Topstrand3primedistancefromend, "", []string{}))
	partWithEnds := AddOverhang(part.Seq, bitToAdd5prime, "5prime")

	newpart = part.Dup()
	newpart.Seq = partWithEnds
	return
}

// AddL1DAdaptor adds a downstream (3') adaptor for making a level 0 part compatible for Level 1 hierarchical assembly, specifying the desired level 1 class the level0 part should be made into.
// TypeIIs recognition site + spacer + correct overhang in correct orientation will be added according to the correct enzyme at a specified level according to a specified
// assembly standard.
// If reverseOrientation is set to true the adaptor will be added such that the level 1 part will bind in the reverse orientation.
func AddL1DAdaptor(part wtype.DNASequence, assemblyStandard AssemblyStandard, level string, class string, reverseOrientation bool) (newpart wtype.DNASequence, err error) {

	enzyme, err := lookUpEnzyme(assemblyStandard, level)

	if err != nil {
		return newpart, err
	}
	bitToAdd5, bitToAdd, err := lookUpOverhangs(assemblyStandard, level, class)

	if err != nil {
		return newpart, err
	}

	if reverseOrientation {
		bitToAdd = wtype.RevComp(bitToAdd5)
	}

	bitToAdd3prime := MakeOverhang(enzyme, "3prime", bitToAdd, ChooseSpacer(enzyme.Topstrand3primedistancefromend, "", []string{}))
	partWithEnds := AddOverhang(part.Seq, bitToAdd3prime, "3prime")

	newpart = part.Dup()
	newpart.Seq = partWithEnds
	return
}

// lookUpEnzyme looks up enzyme according to assembly standard and level.
// Errors will be returned if an entry is missing
func lookUpEnzyme(assemblyStandard AssemblyStandard, level string) (enzyme wtype.TypeIIs, err error) {

	assemblyLevel, err := assemblyStandard.GetLevel(level)

	if err != nil {
		return enzyme, err
	}

	enzyme = assemblyLevel.GetEnzyme()

	return enzyme, nil
}

// lookUpOverhangs looks up overhangs according to assembly standard, level and part class.
// Errors will be returned if an entry is missing or overhangs are found to be empty
func lookUpOverhangs(assemblyStandard AssemblyStandard, level string, class string) (upstream string, downstream string, err error) {

	assemblyLevel, err := assemblyStandard.GetLevel(level)

	if err != nil {
		return "", "", err
	}

	ends, err := assemblyLevel.GetPartOverhangs(class)

	return ends.Upstream, ends.Downstream, err
}

// AddStandardStickyEndsfromClass adds sticky ends to a DNA part according to the class identifier (e.g. PRO, 5U, CDS).
// An error will be returned if any invalid class or level is requested.
func AddStandardStickyEndsfromClass(part wtype.DNASequence, assemblyStandard AssemblyStandard, level string, class string) (partWithEnds wtype.DNASequence, err error) {

	enzyme, err := lookUpEnzyme(assemblyStandard, level)

	if err != nil {
		return partWithEnds, err
	}
	bittoadd, bittoadd3, err := lookUpOverhangs(assemblyStandard, level, class)

	if err != nil {
		return partWithEnds, err
	}

	// This code will find the minimal additional overhang to add
	// with this commented out the full overhang is added
	//bittoadd = findMinimumAdditional5PrimeAddition(bittoadd, part)

	bittoadd5prime := MakeOverhang(enzyme, "5prime", bittoadd, ChooseSpacer(enzyme.Topstrand3primedistancefromend, "", []string{}))

	partwith5primeend := AddOverhang(part.Seq, bittoadd5prime, "5prime")

	// This code will find the minimal additional overhang to add
	// with this commented out the full overhang is added
	//bittoadd3 = findMinimumAdditional3PrimeAddition(bittoadd3, part)

	bittoadd3prime := MakeOverhang(enzyme, "3prime", bittoadd3, ChooseSpacer(enzyme.Topstrand3primedistancefromend, "", []string{}))

	partwithends := AddOverhang(partwith5primeend, bittoadd3prime, "3prime")

	partWithEnds.Nm = part.Nm
	partWithEnds.Plasmid = part.Plasmid
	partWithEnds.Seq = partwithends

	return partWithEnds, err
}

func upper(s string) string {
	return strings.ToUpper(s)
}

func findMinimumAdditional3PrimeAddition(desiredstickyend3prime string, part wtype.DNASequence) (bittoadd string) {

	var present string

	// This code will look for subparts of a standard overhang to add the minimum number of additional nucleotides with a partial match e.g. AATG contains ATG only so we just add A
	for i := 0; i < len(desiredstickyend3prime)+1; i++ {
		truncated := desiredstickyend3prime[:i]
		if strings.HasSuffix(upper(part.Seq), upper(truncated)) {
			present = truncated
		}
	}
	if len(present) == len(desiredstickyend3prime) {
		bittoadd = ""
	} else if len(present) > 0 {
		bittoadd = desiredstickyend3prime[len(present):]
	} else {
		bittoadd = desiredstickyend3prime
	}
	return
}

func findMinimumAdditional5PrimeAddition(desiredstickyend5prime string, part wtype.DNASequence) (bittoadd string) {
	// This code will look for subparts of a standard overhang to add the minimum number of additional nucleotides with a partial match e.g. AATG contains ATG only so we just add A

	var present string

	for i := len(desiredstickyend5prime) - 1; i >= 0; i-- {
		truncated := desiredstickyend5prime[i:]
		if strings.HasPrefix(upper(part.Seq), upper(truncated)) {
			present = truncated
		}
	}
	if len(present) == len(desiredstickyend5prime) {
		bittoadd = ""
	} else if len(present) > 0 {
		bittoadd = desiredstickyend5prime[:len(present)+1]
	} else {
		bittoadd = desiredstickyend5prime
	}
	return
}

// AddCustomEnds adds specified ends to the part sequence based upon enzyme chosen and the desired overhangs after digestion
func AddCustomEnds(part wtype.DNASequence, enzyme wtype.TypeIIs, desiredstickyend5prime string, desiredstickyend3prime string) (Partwithends wtype.DNASequence) {

	bittoadd := findMinimumAdditional5PrimeAddition(desiredstickyend5prime, part)

	bittoadd5prime := MakeOverhang(enzyme, "5prime", bittoadd, ChooseSpacer(enzyme.Topstrand3primedistancefromend, "", []string{}))
	partwith5primeend := AddOverhang(part.Seq, bittoadd5prime, "5prime")

	bittoadd3 := findMinimumAdditional3PrimeAddition(desiredstickyend3prime, part)

	bittoadd3prime := MakeOverhang(enzyme, "3prime", bittoadd3, ChooseSpacer(enzyme.Topstrand3primedistancefromend, "", []string{}))
	partwithends := AddOverhang(partwith5primeend, bittoadd3prime, "3prime")

	Partwithends.Nm = part.Nm
	Partwithends.Plasmid = part.Plasmid
	Partwithends.Seq = partwithends
	return Partwithends
}

// MakeStandardTypeIIsassemblyParts adds compatible ends to a set of parts based on the rules of a typeIIS assembly standard.
func MakeStandardTypeIIsassemblyParts(parts []wtype.DNASequence, assemblystandard AssemblyStandard, level string, partClasses []string) (partswithends []wtype.DNASequence, err error) {

	if len(partClasses) == len(parts) {
		for i := 0; i < len(parts); i++ {
			var partwithends wtype.DNASequence

			partwithends, err = AddStandardStickyEndsfromClass(parts[i], assemblystandard, level, partClasses[i])
			if err != nil {
				return []wtype.DNASequence{}, err
			}
			partswithends = append(partswithends, partwithends)
		}
	} else {
		return partswithends, fmt.Errorf("Number of parts %d (%+v) does not match number of classes specified %d (%+v)", len(parts), partNames(parts), len(partClasses), partClasses)
	}

	return partswithends, err
}

// CheckForExistingTypeIISEnds checks for whether a part already has typeIIs ends added.
func CheckForExistingTypeIISEnds(part wtype.DNASequence, enzyme wtype.TypeIIs) (numberofsitesfound int, stickyends5 []string, stickyends3 []string) {

	enz, err := lookup.RestrictionEnzyme(enzyme.Name())
	if err != nil {
		panic(err.Error())
	}

	sites := RestrictionSiteFinder(part, enz)

	numberofsitesfound = sites[0].NumberOfSites()
	_, stickyends5, stickyends3 = TypeIIsdigest(part, enzyme)

	return
}

/*

func HandleExistingEnds (parts []wtype.DNASequence, enzymewtype.RestrictionEnzyme)(partswithoverhangs []wtype.DNASequence {
	partswithexistingsites := make([]RestrictionSites, 0)

	for _, part := range parts {
		sites := Restrictionsitefinder(part, wtype.RestrictionEnzyme{enzyme})
		if len(sites) != 0 {
			partswithexistingsites = append(partswithexistingsites, sites)
		}

	}
	return
}

func AddStandardVectorEnds (vector wtype.DNASequence, standard, level string) (vectrowithends wtype.DNASequence) {

	}
*/

// AddOverhang is the lowest level function to add an overhang to a sequence as a string
func AddOverhang(seq string, bittoadd string, end string) (seqwithoverhang string) {

	if end == "5prime" {
		seqwithoverhang = strings.Join([]string{bittoadd, seq}, "")
	}
	if end == "3prime" {
		seqwithoverhang = strings.Join([]string{seq, bittoadd}, "")
	}
	return seqwithoverhang
}

func allCombinations(arr [][]string) []string {
	if len(arr) == 1 {
		return arr[0]
	}

	results := make([]string, 0)
	allRem := allCombinations(arr[1:])
	for i := 0; i < len(allRem); i++ {
		for j := 0; j < len(arr[0]); j++ {
			x := arr[0][j] + allRem[i]
			results = append(results, x)
		}
	}
	return results
}

// MakeAllSpacerOptions returns an array of all sequence possibilities for a spacer based upon length.
func MakeAllSpacerOptions(spacerlength int) (finalarray []string) {
	// only works for spacer length 1 or 2

	// new better code, but untested! test and replace code below
	newarray := make([][]string, 0)
	for i := 0; i < spacerlength; i++ {
		newarray = append(newarray, nucleotides)
	}

	finalarray = allCombinations(newarray)

	return finalarray
}

// ChooseSpacer picks the first valid spacer which avoids all sequences to avoid.
func ChooseSpacer(spacerlength int, seq string, seqstoavoid []string) (spacer string) {
	// very simple case to start with

	possibilities := MakeAllSpacerOptions(spacerlength)

	if len(seqstoavoid) == 0 {
		spacer = possibilities[0]
	} else {
		for _, possibility := range possibilities {
			if len(search.FindAllStrings(strings.Join([]string{seq, possibility}, ""), seqstoavoid)) == 0 &&
				len(search.FindAllStrings(strings.Join([]string{possibility, seq}, ""), seqstoavoid)) == 0 &&
				len(search.FindAllStrings(sequences.RevComp(strings.Join([]string{possibility, seq}, "")), seqstoavoid)) == 0 &&
				len(search.FindAllStrings(sequences.RevComp(strings.Join([]string{seq, possibility}, "")), seqstoavoid)) == 0 {
				spacer = possibility
			}
		}
	}
	return spacer
}

var nucleotides = []string{"A", "T", "C", "G"}

// MakeOverhang adds an overhang based upon the enzyme chosen, the choice of end ("5Prime" or "3Prime"), the desired sticky end and the desired spacer.
func MakeOverhang(enzyme wtype.TypeIIs, end string, stickyendseq string, spacer string) (seqwithoverhang string) {
	if end == "5prime" {
		if enzyme.Topstrand3primedistancefromend < 0 {
			panic("Unlikely to work with this enzyme in making a 5'prime spacer")
		}

		if len(spacer) != enzyme.Topstrand3primedistancefromend {
			panic("length of spacer will lead to cutting at run position! change length to match enzyme NN region length")
		}
		seqwithoverhang = strings.Join([]string{enzyme.RestrictionEnzyme.RecognitionSequence, spacer, stickyendseq}, "")
	}

	// This case needs work, but may not appear in reality so is a place holder for now until a real scenario becomes apparent
	if end == "3prime" {
		/*if enzyme.Topstrand3primedistancefromend < 0 && len(spacer) == enzyme.Bottomstrand5primedistancefromend {
			seqwithoverhang = strings.Join([]string{stickyendseq, spacer, enzyme.RestrictionEnzyme.RecognitionSequence}, "")
		}*/
		seqwithoverhang = strings.Join([]string{stickyendseq, spacer, sequences.RevComp(enzyme.RestrictionEnzyme.RecognitionSequence)}, "")
	}
	return seqwithoverhang

}

// Assembly standards
var availableStandards = map[string]AssemblyStandard{
	"Custom":      customStandard,
	"MoClo":       mocloStandard,
	"MoClo_Raven": mocloRavenStandard,
	"Antibody":    antibodyStandard,
}

func allStandards() (standards []string) {
	for k := range availableStandards {
		standards = append(standards, k)
	}
	sort.Strings(standards)
	return
}

// LookupAssemblyStandard looks up a TypeIIS Assembly Standard by name.
// An error will be returned if no Assembly standard is found for the requested name.
func LookupAssemblyStandard(name string) (standard AssemblyStandard, err error) {
	var found bool
	standard, found = availableStandards[name]
	if !found {
		err = fmt.Errorf("assembly standard %s not found. Vslid options are: %s", name, strings.Join(allStandards(), ";"))
	}
	return
}

var mocloStandard = AssemblyStandard{
	Name: "MoClo",
	Levels: map[string]AssemblyLevel{
		"Level0": {
			Enzyme: BsaI,
			PartOverhangs: map[string]StandardOverhangs{
				"Pro":         {"GGAG", "TACT"},
				"5U":          {"TACT", "CCAT"},
				"5U(f)":       {"TACT", "CCAT"},
				"Pro + 5U(f)": {"GGAG", "CCAT"},
				"Pro + 5U":    {"GGAG", "AATG"},
				"NT1":         {"CCAT", "AATG"},
				"5U + NT1":    {"TACT", "AATG"},
				"CDS1":        {"AATG", "GCTT"},
				"CDS1 ns":     {"AATG", "TTCG"},
				"NT2":         {"AATG", "AGGT"},
				"SP":          {"AATG", "AGGT"},
				"CDS2 ns":     {"AGGT", "TTCG"},
				"CDS2":        {"AGGT", "GCTT"},
				"CT":          {"TTCG", "GCTT"},
				"3U":          {"GCTT", "GGTA"},
				"Ter":         {"GGTA", "CGCT"},
				"3U + Ter":    {"GCTT", "CGCT"},
			},
			EntryVectorEnds: StandardOverhangs{"TAAT", "GTCG"},
		},
		"Level1": {
			Enzyme:          BpiI,
			PartOverhangs:   map[string]StandardOverhangs{},
			EntryVectorEnds: StandardOverhangs{"", ""},
		},
	},
}

var mocloRavenStandard = AssemblyStandard{
	Name: "MoClo_Raven",
	Levels: map[string]AssemblyLevel{
		"Level0": {
			Enzyme: BsaI,
			PartOverhangs: map[string]StandardOverhangs{
				"Pro":         {"GAGG", "TACT"},
				"5U":          {"TACT", "CCAT"},
				"5U(f)":       {"TACT", "CCAT"},
				"Pro + 5U(f)": {"GGAG", "CCAT"},
				"Pro + 5U":    {"GGAG", "AATG"},
				"NT1":         {"CCAT", "AATG"},
				"5U + NT1":    {"TACT", "AATG"},
				"CDS1":        {"AATG", "GCTT"},
				"CDS1 ns":     {"AATG", "TTCG"},
				"NT2":         {"AATG", "AGGT"},
				"SP":          {"AATG", "AGGT"},
				"CDS2 ns":     {"AGGT", "TTCG"},
				"CDS2":        {"AGGT", "GCTT"},
				"CT":          {"TTCG", "GCTT"},
				"3U":          {"GCTT", "GGTA"},
				"Ter":         {"GGTA", "CGCT"},
				"3U + Ter":    {"GCTT", "GCTT"}, // both same ! look into this
			},
			EntryVectorEnds: StandardOverhangs{"AAGC", "CCTC"},
		},
		"Level1": {
			Enzyme:          BpiI,
			PartOverhangs:   map[string]StandardOverhangs{},
			EntryVectorEnds: StandardOverhangs{"", ""},
		},
	},
}

var customStandard = AssemblyStandard{
	Name: "Custom",
	Levels: map[string]AssemblyLevel{
		"Level0": {
			Enzyme: BsaI,
			PartOverhangs: map[string]StandardOverhangs{
				"L1Uadaptor":                        {"GTCG", "GGAG"}, // adaptor to add SapI sites to clone into level 1 vector
				"L1Uadaptor + Pro":                  {"GTCG", "TTTT"}, // adaptor to add SapI sites to clone into level 1 vector
				"L1Uadaptor + Pro(MoClo)":           {"GTCG", "TACT"}, // original MoClo overhang of TACT
				"L1Uadaptor + Pro + 5U":             {"GTCG", "CCAT"}, // adaptor to add SapI sites to clone into level 1 vector
				"L1Uadaptor + Pro + 5U + NT1":       {"GTCG", "TATG"}, // adaptor to add SapI sites to clone into level 1 vector
				"Pro":                               {"GGAG", "TTTT"}, //changed from MoClo TACT to TTTT to conform with Protein Paintbox
				"5U":                                {"TTTT", "CCAT"}, //changed from MoClo TACT to TTTT to conform with Protein Paintbox
				"5U(f)":                             {"TTTT", "CCAT"}, //changed from MoClo TACT to TTTT to conform with Protein Paintbox
				"Pro + 5U(f)":                       {"GGAG", "CCAT"},
				"Pro + 5U":                          {"GGAG", "CCAT"},
				"Pro + 5U + NT1":                    {"GGAG", "TATG"}, //changed AATG to TATG to work with Kosuri paper RBSs
				"NT1":                               {"CCAT", "TATG"}, //changed AATG to TATG to work with Kosuri paper RBSs
				"5U + NT1":                          {"TTTT", "TATG"}, //changed AATG to TATG to work with Kosuri paper RBSs
				"5U(MoClo) + NT1":                   {"TACT", "TATG"}, //original MoClo overhang of TACT
				"5U + NT1 + CDS1":                   {"TTTT", "GCTT"}, //changed from MoClo TACT to TTTT to conform with Protein Paintbox
				"5U + NT1 + CDS1 + 3U":              {"TTTT", "CCCC"}, //changed from MoClo TACT to TTTT to conform with Protein Paintbox and changed GGTA to CCCC to conform with Protein Paintbox
				"CDS1":                              {"TATG", "GCTT"}, //changed AATG to TATG to work with Kosuri paper RBSs
				"CDS1 + 3U":                         {"TATG", "CCCC"}, //changed AATG to TATG to work with Kosuri paper RBSs and changed GGTA to CCCC to conform with Protein Paintbox
				"CDS1 + 3U(MoClo)":                  {"TATG", "GGTA"}, //original MoClo overhang of GGTA
				"CDS1 ns":                           {"TATG", "TTCG"}, //changed AATG to TATG to work with Kosuri paper RBSs
				"CDS1 + CT + 3U + Ter + L1Dadaptor": {"TATG", "TAAT"}, //changed AATG to TATG to work with Kosuri paper RBSs
				"NT2":                               {"TATG", "AGGT"}, //changed AATG to TATG to work with Kosuri paper RBSs
				"SP":                                {"TATG", "AGGT"}, //changed AATG to TATG to work with Kosuri paper RBSs
				"CDS2 ns":                           {"AGGT", "TTCG"},
				"CDS2":                              {"AGGT", "GCTT"},
				"CT":                                {"TTCG", "GCTT"},
				"3U":                                {"GCTT", "CCCC"}, //changed GGTA to CCCC to conform with Protein Paintbox
				"Ter":                               {"CCCC", "CGCT"},
				"3U + Ter":                          {"GCTT", "CGCT"},
				"3U + Ter + L1Dadaptor":             {"GCTT", "TAAT"},
				"CT + 3U + Ter + L1Dadaptor":        {"TTCG", "TAAT"},
				"L1Dadaptor":                        {"CGCT", "TAAT"},
				"Ter + L1Dadaptor":                  {"CCCC", "TAAT"},
				"Ter(MoClo) + L1Dadaptor":           {"GGTA", "TAAT"},
			},
			EntryVectorEnds: StandardOverhangs{"TAAT", "GTCG"},
		},
		"Level1": {
			Enzyme: SapI,
			PartOverhangs: map[string]StandardOverhangs{
				"Device1": {"GAA", "ACC"},
				"Device2": {"ACC", "CTG"},
				"Device3": {"CTG", "GGT"},
			},
			EntryVectorEnds: StandardOverhangs{"GGT", "GAA"},
		},
	},
}

var antibodyStandard = AssemblyStandard{
	Name: "Antibody",
	Levels: map[string]AssemblyLevel{
		"Heavy": {
			Enzyme: SapI,
			PartOverhangs: map[string]StandardOverhangs{
				"Part1": {"GCG", "TCG"},
				"Part2": {"TGG", "CTG"},
				"Part3": {"CTG", "AAG"},
			},
			EntryVectorEnds: StandardOverhangs{"GCG", "AAG"},
		},
		"Light": {
			Enzyme: SapI,
			PartOverhangs: map[string]StandardOverhangs{
				"Part1": {"GCG", "TCG"},
				"Part2": {"TGG", "CTG"},
				"Part3": {"CTG", "AAG"},
			},
			EntryVectorEnds: StandardOverhangs{"GCG", "AAG"},
		},
	},
}

// AssemblyStandard is an assembly standard for modular assembly of DNA parts using TypeIIs enzyme assembly.
// The AssemblyStandard may consist of a number of assembly levels which may use a different enzyme, set of standard overhangs or both.
type AssemblyStandard struct {
	// Name of the Assembly Standard.
	Name string
	// The AssemblyStandard may consist of a number of assembly levels which may use a different enzyme, set of standard overhangs or both.
	// The name of the level is used as the key to calling an AssemblyLevel object.
	Levels map[string]AssemblyLevel
}

// Enzyme returns the typeIIs enzyme for a specified named AssemblyLevel.
// An error is returned if an invalid level is requested for the AssemblyStandard.
func (l AssemblyStandard) Enzyme(level string) (enz wtype.TypeIIs, err error) {

	assemblyLevel, err := l.GetLevel(level)
	if err != nil {
		return enz, err
	}
	enz = assemblyLevel.GetEnzyme()
	return
}

// LevelNames returns the names of all valid Assembly Levels for an AssemblyStandard.
func (l AssemblyStandard) LevelNames() []string {

	var levels []string
	for level := range l.Levels {
		levels = append(levels, level)
	}

	sort.Strings(levels)

	return levels
}

// GetLevel returns an AssemblyLevel for a specified named AssemblyLevel.
// An error is returned if an invalid level is requested for the AssemblyStandard.
func (l AssemblyStandard) GetLevel(level string) (assemblyLevel AssemblyLevel, err error) {

	assemblyLevel, found := l.Levels[level]
	if !found {
		return assemblyLevel, fmt.Errorf("No level %s found for assembly standard %s, found %+v", level, l.Name, l.LevelNames())
	}
	return
}

// AssemblyLevel is a specified TypeIIs standard for assembly of a series of labelled parts.
type AssemblyLevel struct {
	//Enzyme used for assembly
	Enzyme wtype.TypeIIs
	// map of part labels to standard overhangs.
	PartOverhangs map[string]StandardOverhangs
	// Expected overhangs in the entry vector.
	EntryVectorEnds StandardOverhangs // Vector 5prime can also be found in Endstable position 0
}

// StandardOverhangs represents the upstream and downstream overhangs expected for a part or vector in a
// level of an Assembly Standard.
type StandardOverhangs struct {
	// Overhang expected to be left at the upstream end of a part after typeIIs digestion.
	Upstream string
	// Overhang expected to be left at the downstream end of a part after typeIIs digestion.
	Downstream string
}

// AnnotationOptions returns all valid part labels for a specified AssemblyLevel.
func (l AssemblyLevel) AnnotationOptions() []string {
	var ls []string
	for class := range l.PartOverhangs {
		ls = append(ls, class)
	}

	sort.Strings(ls)
	return ls
}

// GetEnzyme returns the typeIIs enzyme for the AssemblyLevel.
func (l AssemblyLevel) GetEnzyme() wtype.TypeIIs {
	return l.Enzyme
}

// GetVectorEnds returns the expected vector overhangs for the AssemblyLevel.
func (l AssemblyLevel) GetVectorEnds() StandardOverhangs {
	return l.EntryVectorEnds
}

// GetPartOverhangs returns the expected part overhangs for the specified part class.
func (l AssemblyLevel) GetPartOverhangs(class string) (overhangs StandardOverhangs, err error) {

	overhangs, found := l.PartOverhangs[class]

	if !found {
		return overhangs, fmt.Errorf("No overhangs found for %s in assembly standard: found: %+v", class, l.AnnotationOptions())
	}

	if overhangs.Upstream == "" {
		return overhangs, fmt.Errorf("blunt 5' overhang found for %s", class)
	}

	if overhangs.Downstream == "" {
		return overhangs, fmt.Errorf("blunt 3' overhang found for %s", class)
	}

	return
}

func partNames(parts []wtype.DNASequence) []string {
	var names []string
	for _, part := range parts {
		names = append(names, part.Name())
	}
	return names
}
