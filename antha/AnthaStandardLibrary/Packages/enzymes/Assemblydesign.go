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

// Package for working with enzymes; in particular restriction enzymes
package enzymes

import (
	"fmt"
	"strings"

	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/enzymes/lookup"
	. "github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/search"
	. "github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/sequences"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/text"
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

// Key general function to design parts for assembly based on type IIs enzyme, parts in order, fixed vector sequence (containing sites for the corresponding enzyme).
func MakeScarfreeCustomTypeIIsassemblyParts(parts []wtype.DNASequence, vector wtype.DNASequence, enzyme wtype.TypeIIs) (partswithends []wtype.DNASequence) {

	partswithends = make([]wtype.DNASequence, 0)
	var partwithends wtype.DNASequence

	// find sticky ends from cutting vector with enzyme

	fragments, stickyends5, _ := TypeIIsdigest(vector, enzyme)

	//initialise

	desiredstickyend5prime := ""

	vector3primestickyend := ""

	// add better logic for the scenarios where the vector is cut more than twice or we want to add fragment in either direction
	// picks largest fragment

	for i := 0; i < len(stickyends5)-1; i++ {

		currentlargestfragment := ""

		if stickyends5[i] != "" && len(fragments[i]) > len(currentlargestfragment) {

			currentlargestfragment = fragments[i]
			// RevComp() // fill in later
			vector3primestickyend = stickyends5[i]
			desiredstickyend5prime = stickyends5[i+1]
			/*{
				break
			}*/
		}
	} // fill in later

	// declare as blank so no end added
	desiredstickyend3prime := ""

	for i := 0; i < len(parts); i++ {
		if i == (len(parts) - 1) {
			desiredstickyend3prime = vector3primestickyend
		}
		partwithends = AddCustomEnds(parts[i], enzyme, desiredstickyend5prime, desiredstickyend3prime)
		partwithends.Nm = parts[i].Nm
		partswithends = append(partswithends, partwithends)

		desiredstickyend5prime = Suffix(parts[i].Seq, enzyme.RestrictionEnzyme.EndLength)

	}

	return partswithends
}

// Adds sticky ends to a dna part based upon an assembly standard (e.g. MoClo) and position of part within an array.
func AddStandardStickyEnds(part wtype.DNASequence, assemblystandard string, level string, numberofparts int, position int) (Partwithends wtype.DNASequence) {

	if part.Seq != "" {
		if position > numberofparts {
			panic("Cannot have position greater than number of parts")
		}
		/*if position == 0 {
			("1st position = 1, not 0")
		}
		*/
		vectorends := Vectorends[assemblystandard][level] // this could also look up [assemblystandard][level][numberofparts][0]

		enzyme := Enzymelookup[assemblystandard][level]

		bittoadd := Endlinks[assemblystandard][level][numberofparts][position]
		if strings.HasPrefix(part.Seq, bittoadd) == true {
			bittoadd = ""
		}

		// This code will look for subparts of a standard overhang to add the minimum number of additional nucleotides with a partial match e.g. AATG contains ATG only so we just add A

		truncated := bittoadd[1:]
		// fmt.Println("truncated", truncated)
		if strings.HasPrefix(part.Seq, truncated) == true {
			//truncated = bittoadd[:1]
			//// fmt.Println("truncated", truncated)
			bittoadd = bittoadd[:1]
		}

		bittoadd5prime := Makeoverhang(enzyme, "5prime", bittoadd, ChooseSpacer(enzyme.Topstrand3primedistancefromend, "", []string{}))
		// fmt.Println("bittoadd5prime", bittoadd5prime)
		partwith5primeend := Addoverhang(part.Seq, bittoadd5prime, "5prime")

		bittoadd3 := Endlinks[assemblystandard][level][numberofparts][position+1]
		// fmt.Println("bittoadd3", bittoadd3)

		if numberofparts == position {
			bittoadd3 = RevComp(vectorends[0])
		}
		if strings.HasSuffix(part.Seq, bittoadd3) == true {
			bittoadd3 = ""
		}
		//// fmt.Println("bittoadd3", bittoadd3)
		bittoadd3prime := Makeoverhang(enzyme, "3prime", bittoadd3, ChooseSpacer(enzyme.Topstrand3primedistancefromend, "", []string{}))
		// fmt.Println("bittoadd3prime", bittoadd3prime)
		partwithends := Addoverhang(partwith5primeend, bittoadd3prime, "3prime")

		Partwithends.Nm = part.Nm
		Partwithends.Plasmid = part.Plasmid
		Partwithends.Seq = partwithends
	}
	return Partwithends
}

// This func adds an upstream (5') adaptor for making a level 0 part compatible for Level 1 hierarchical assembly, specifying the desired level 1 class the level0 part should be made into.
// TypeIIs recognition site + spacer + correct overhang in correct orientation will be added according to the correct enzyme at a specified level according to a specified
// assembly standard.
// If reverseOrientation is set to true the adaptor will be added such that the level 1 part will bind in the reverse orientation.
func AddL1UAdaptor(part wtype.DNASequence, assemblyStandard string, level string, class string, reverseOrientation bool) (newpart wtype.DNASequence, err error) {

	enzyme, err := lookUpEnzyme(assemblyStandard, level)

	if err != nil {
		return newpart, err
	}

	bitToAdd, _, err := lookUpOverhangs(assemblyStandard, level, class)

	if err != nil {
		return newpart, err
	}

	if reverseOrientation {
		bitToAdd = wtype.RevComp(bitToAdd)
	}

	bitToAdd5prime := Makeoverhang(enzyme, "5prime", bitToAdd, ChooseSpacer(enzyme.Topstrand3primedistancefromend, "", []string{}))
	partWithEnds := Addoverhang(part.Seq, bitToAdd5prime, "5prime")

	newpart = part.Dup()
	newpart.Seq = partWithEnds
	return
}

// This func adds a downstream (3') adaptor for making a level 0 part compatible for Level 1 hierarchical assembly, specifying the desired level 1 class the level0 part should be made into.
// TypeIIs recognition site + spacer + correct overhang in correct orientation will be added according to the correct enzyme at a specified level according to a specified
// assembly standard.
// If reverseOrientation is set to true the adaptor will be added such that the level 1 part will bind in the reverse orientation.
func AddL1DAdaptor(part wtype.DNASequence, assemblyStandard string, level string, class string, reverseOrientation bool) (newpart wtype.DNASequence, err error) {

	enzyme, err := lookUpEnzyme(assemblyStandard, level)

	if err != nil {
		return newpart, err
	}
	_, bitToAdd, err := lookUpOverhangs(assemblyStandard, level, class)

	if err != nil {
		return newpart, err
	}

	if reverseOrientation {
		bitToAdd = wtype.RevComp(bitToAdd)
	}

	bitToAdd3prime := Makeoverhang(enzyme, "3prime", bitToAdd, ChooseSpacer(enzyme.Topstrand3primedistancefromend, "", []string{}))
	partWithEnds := Addoverhang(part.Seq, bitToAdd3prime, "3prime")

	newpart = part.Dup()
	newpart.Seq = partWithEnds
	return
}

// look up enzyme according to assembly standard and level.
// Errors will be returned if an entry is missing
func lookUpEnzyme(assemblyStandard string, level string) (enzyme wtype.TypeIIs, err error) {
	enzymeStandard, found := Enzymelookup[assemblyStandard]

	if !found {
		return enzyme, fmt.Errorf("No enzyme found for standard %s", assemblyStandard)
	}

	enzyme, found = enzymeStandard[level]

	if !found {
		return enzyme, fmt.Errorf("No level found for enzyme for standard %s at level %s", assemblyStandard, level)
	}
	return enzyme, nil
}

// look up overhangs according to assembly standard, level and part class.
// Errors will be returned if an entry is missing or overhangs are found to be empty
func lookUpOverhangs(assemblyStandard string, level string, class string) (upstream string, downstream string, err error) {
	standard, found := EndlinksString[assemblyStandard]

	if !found {
		return "", "", fmt.Errorf("No assembly standard %s found", assemblyStandard)
	}

	level1, found := standard[level]

	if !found {
		return "", "", fmt.Errorf("No Level: %s found for standard %s", level, assemblyStandard)
	}

	overhangs, found := level1[class]

	if !found {
		return "", "", fmt.Errorf("No overhangs found for %s in standard %s", class, assemblyStandard)
	}

	if len(overhangs) != 2 {
		return "", "", fmt.Errorf("found %d overhangs for %s in standard %s, expecting %d: overhangs: %s", len(overhangs), class, assemblyStandard, 2, overhangs)

	}

	if overhangs[0] == "" {
		return overhangs[0], overhangs[1], fmt.Errorf("blunt 5' overhang found for %s in standard %s, expecting %d", class, assemblyStandard, 2)
	}

	if overhangs[1] == "" {
		return overhangs[0], overhangs[1], fmt.Errorf("blunt 3' overhang found for %s in standard %s, expecting %d", class, assemblyStandard, 2)
	}

	return overhangs[0], overhangs[1], nil
}

// Adds sticky ends to dna part according to the class identifier (e.g. PRO, 5U, CDS)
func AddStandardStickyEndsfromClass(part wtype.DNASequence, assemblystandard string, level string, class string) (Partwithends wtype.DNASequence, err error) {

	enzyme := Enzymelookup[assemblystandard][level]

	bitstoadd, found := EndlinksString[assemblystandard][level][class]

	if !found {
		err = fmt.Errorf("Class " + class + " not found in Assmbly standard map of " + assemblystandard + " level " + level)
		return Partwithends, err
	}

	bittoadd := bitstoadd[0]
	if strings.HasPrefix(part.Seq, bittoadd) == true {
		bittoadd = ""
	}

	// This code will look for subparts of a standard overhang to add the minimum number of additional nucleotides with a partial match e.g. AATG contains ATG only so we just add A

	truncated := bittoadd[1:]
	// fmt.Println("truncated", truncated)
	if strings.HasPrefix(part.Seq, truncated) == true {
		//truncated = bittoadd[:1]
		//// fmt.Println("truncated", truncated)
		bittoadd = bittoadd[:1]
	}

	bittoadd5prime := Makeoverhang(enzyme, "5prime", bittoadd, ChooseSpacer(enzyme.Topstrand3primedistancefromend, "", []string{}))
	// fmt.Println("bittoadd5prime", bittoadd5prime)
	partwith5primeend := Addoverhang(part.Seq, bittoadd5prime, "5prime")

	bittoadd3 := bitstoadd[1]
	// fmt.Println("bittoadd3", bittoadd3)

	if strings.HasSuffix(part.Seq, bittoadd3) == true {
		bittoadd3 = ""
	}
	//// fmt.Println("bittoadd3", bittoadd3)
	bittoadd3prime := Makeoverhang(enzyme, "3prime", bittoadd3, ChooseSpacer(enzyme.Topstrand3primedistancefromend, "", []string{}))
	// fmt.Println("bittoadd3prime", bittoadd3prime)
	partwithends := Addoverhang(partwith5primeend, bittoadd3prime, "3prime")

	Partwithends.Nm = part.Nm
	Partwithends.Plasmid = part.Plasmid
	Partwithends.Seq = partwithends

	return Partwithends, err
}

// Adds ends to the part sequence based upon enzyme chosen and the desired overhangs after digestion
func AddCustomEnds(part wtype.DNASequence, enzyme wtype.TypeIIs, desiredstickyend5prime string, desiredstickyend3prime string) (Partwithends wtype.DNASequence) {

	bittoadd := desiredstickyend5prime
	if strings.HasPrefix(part.Seq, bittoadd) == true {
		bittoadd = ""
	}
	bittoadd5prime := Makeoverhang(enzyme, "5prime", bittoadd, ChooseSpacer(enzyme.Topstrand3primedistancefromend, "", []string{}))

	partwith5primeend := Addoverhang(part.Seq, bittoadd5prime, "5prime")

	bittoadd3 := desiredstickyend3prime

	if strings.HasSuffix(part.Seq, bittoadd3) == true {
		bittoadd3 = ""
	}

	bittoadd3prime := Makeoverhang(enzyme, "3prime", bittoadd3, ChooseSpacer(enzyme.Topstrand3primedistancefromend, "", []string{}))

	partwithends := Addoverhang(partwith5primeend, bittoadd3prime, "3prime")

	Partwithends.Nm = part.Nm
	Partwithends.Plasmid = part.Plasmid
	Partwithends.Seq = partwithends
	return Partwithends
}

// Add compatible ends to an array of parts based on the rules of a typeIIS assembly standard
func MakeStandardTypeIIsassemblyParts(parts []wtype.DNASequence, assemblystandard string, level string, optionalpartclasses []string) (partswithends []wtype.DNASequence, err error) {

	partswithends = make([]wtype.DNASequence, 0)
	var partwithends wtype.DNASequence

	if len(optionalpartclasses) != 0 {
		if len(optionalpartclasses) == len(parts) {
			for i := 0; i < len(parts); i++ {
				partwithends, err = AddStandardStickyEndsfromClass(parts[i], assemblystandard, level, optionalpartclasses[i])
				if err != nil {
					return []wtype.DNASequence{}, err
				}
				partswithends = append(partswithends, partwithends)
			}
		}
	} else {

		for i := 0; i < len(parts); i++ {
			partwithends = AddStandardStickyEnds(parts[i], assemblystandard, level, (len(parts)), i+1)
			partswithends = append(partswithends, partwithends)
		}
	}
	return partswithends, err
}

// Utility function to check whether a part already has typeIIs ends added
func CheckForExistingTypeIISEnds(part wtype.DNASequence, enzyme wtype.TypeIIs) (numberofsitesfound int, stickyends5 []string, stickyends3 []string) {

	enz, err := lookup.EnzymeLookup(enzyme.Name)
	if err != nil {
		panic(err.Error())
	}

	sites := Restrictionsitefinder(part, []wtype.RestrictionEnzyme{enz})

	numberofsitesfound = sites[0].Numberofsites
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

// Lowest level function to add an overhang to a sequence as a string
func Addoverhang(seq string, bittoadd string, end string) (seqwithoverhang string) {

	bittoadd = text.Annotate(bittoadd, "blue")

	if end == "5prime" {
		seqwithoverhang = strings.Join([]string{bittoadd, seq}, "")
	}
	if end == "3prime" {
		seqwithoverhang = strings.Join([]string{seq, bittoadd}, "")
	}
	return seqwithoverhang
}

// Returns an array of all sequence possibilities for a spacer based upon length
func Makeallspaceroptions(spacerlength int) (finalarray []string) {
	// only works for spacer length 1 or 2

	// new better code, but untested! test and replace code below
	newarray := make([][]string, 0)
	for i := 0; i < spacerlength; i++ {
		newarray = append(newarray, nucleotides)
	}

	finalarray = AllCombinations(newarray)

	return finalarray
}

// Picks first valid spacer which avoids all sequences to avoid
func ChooseSpacer(spacerlength int, seq string, seqstoavoid []string) (spacer string) {
	// very simple case to start with

	possibilities := Makeallspaceroptions(spacerlength)

	if len(seqstoavoid) == 0 {
		spacer = possibilities[0]
	} else {
		for _, possibility := range possibilities {
			if len(Findallthings(strings.Join([]string{seq, possibility}, ""), seqstoavoid)) == 0 &&
				len(Findallthings(strings.Join([]string{possibility, seq}, ""), seqstoavoid)) == 0 &&
				len(Findallthings(RevComp(strings.Join([]string{possibility, seq}, "")), seqstoavoid)) == 0 &&
				len(Findallthings(RevComp(strings.Join([]string{seq, possibility}, "")), seqstoavoid)) == 0 {
				spacer = possibility
			}
		}
	}
	return spacer
}

var nucleotides = []string{"A", "T", "C", "G"}

// for a dna sequence as a string as input; the function will return an array of 4 sequences appended with either A, T, C or G
func Addnucleotide(s string) (splus1array []string) {

	splus1 := s
	splus1array = make([]string, 0)
	for _, nucleotide := range nucleotides {
		splus1 = strings.Join([]string{s, nucleotide}, "")
		splus1array = append(splus1array, splus1)
	}
	return splus1array
}

// Function to add an overhang based upon the enzyme chosen, the choice of end ("5Prime" or "3Prime")
func Makeoverhang(enzyme wtype.TypeIIs, end string, stickyendseq string, spacer string) (seqwithoverhang string) {
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
		seqwithoverhang = strings.Join([]string{stickyendseq, spacer, RevComp(enzyme.RestrictionEnzyme.RecognitionSequence)}, "")
	}
	return seqwithoverhang

}

// map of 5' sticky ends required for parts in an assembly standard based purely on number of parts for assembly and position of each part in the array
var Endlinks = map[string]map[string]map[int]map[int]string{
	//map["assembly strategy"]map[number of parts]map[part number]"sticky end to add"
	"MoClo_Raven": map[string]map[int]map[int]string{
		"Level0": map[int]map[int]string{
			4: map[int]string{ // overall number of parts in assembly
				0: "AAGC", // position of part in assembly used as key
				1: "GAGG",
				2: "TACT",
				3: "AATG",
				4: "AGGT",
			},
		},
	},
	"MoClo": map[string]map[int]map[int]string{
		"Level0": map[int]map[int]string{
			4: map[int]string{
				0: "AAGC",
				1: "GGTA",
				2: "",
				3: "",
				4: "",
			},
		},
	},
	"Electra": map[string]map[int]map[int]string{
		"Level0": map[int]map[int]string{
			1: map[int]string{
				0: "ATG",
			},
		},
	},
}

// Map describing the sticky ends required for each class of an assembly standard at  a particular level.
var EndlinksString = map[string]map[string]map[string][]string{
	"MoClo": map[string]map[string][]string{
		"Level0": map[string][]string{
			"Pro":         []string{"GGAG", "TACT"},
			"5U":          []string{"TACT", "CCAT"},
			"5U(f)":       []string{"TACT", "CCAT"},
			"Pro + 5U(f)": []string{"GGAG", "CCAT"},
			"Pro + 5U":    []string{"GGAG", "AATG"},
			"NT1":         []string{"CCAT", "AATG"},
			"5U + NT1":    []string{"TACT", "AATG"},
			"CDS1":        []string{"AATG", "GCTT"},
			"CDS1 ns":     []string{"AATG", "TTCG"},
			"NT2":         []string{"AATG", "AGGT"},
			"SP":          []string{"AATG", "AGGT"},
			"CDS2 ns":     []string{"AGGT", "TTCG"},
			"CDS2":        []string{"AGGT", "GCTT"},
			"CT":          []string{"TTCG", "GCTT"},
			"3U":          []string{"GCTT", "GGTA"},
			"Ter":         []string{"GGTA", "CGCT"},
			"3U + Ter":    []string{"GCTT", "CGCT"},
		},
	},
	"Custom": map[string]map[string][]string{
		"Level0": map[string][]string{
			"L1Uadaptor":                  []string{"GTCG", "GGAG"}, // adaptor to add SapI sites to clone into level 1 vector
			"L1Uadaptor + Pro":            []string{"GTCG", "TTTT"}, // adaptor to add SapI sites to clone into level 1 vector
			"L1Uadaptor + Pro + 5U":       []string{"GTCG", "CCAT"}, // adaptor to add SapI sites to clone into level 1 vector
			"L1Uadaptor + Pro + 5U + NT1": []string{"GTCG", "TATG"}, // adaptor to add SapI sites to clone into level 1 vector
			"Pro":                               []string{"GGAG", "TTTT"},
			"5U":                                []string{"TTTT", "CCAT"}, // 5' untranslated, e.g. rbs // changed from MoClo TACT to TTTT to conform with Protein Paintbox??
			"5U(f)":                             []string{"TTTT", "CCAT"},
			"Pro + 5U(f)":                       []string{"GGAG", "CCAT"},
			"Pro + 5U":                          []string{"GGAG", "TATG"}, //changed AATG to TATG to work with Kosuri paper RBSs
			"NT1":                               []string{"CCAT", "TATG"}, //changed AATG to TATG to work with Kosuri paper RBSs
			"5U + NT1":                          []string{"TTTT", "TATG"}, //changed AATG to TATG to work with Kosuri paper RBSs
			"5U + NT1 + CDS1":                   []string{"TTTT", "GCTT"}, //changed AATG to TATG to work with Kosuri paper RBSs
			"5U + NT1 + CDS1 + 3U":              []string{"TTTT", "CCCC"}, //changed AATG to TATG to work with Kosuri paper RBSs
			"CDS1":                              []string{"TATG", "GCTT"}, //changed AATG to TATG to work with Kosuri paper RBSs
			"CDS1 + 3U":                         []string{"TATG", "CCCC"}, //changed AATG to TATG to work with Kosuri paper RBSs
			"CDS1 ns":                           []string{"TATG", "TTCG"}, //changed AATG to TATG to work with Kosuri paper RBSs
			"CDS1 + CT + 3U + Ter + L1Dadaptor": []string{"TATG", "TAAT"}, //changed AATG to TATG to work with Kosuri paper RBSs
			"NT2":                        []string{"TATG", "AGGT"}, //changed AATG to TATG to work with Kosuri paper RBSs
			"SP":                         []string{"TATG", "AGGT"}, //changed AATG to TATG to work with Kosuri paper RBSs
			"CDS2 ns":                    []string{"AGGT", "TTCG"},
			"CDS2":                       []string{"AGGT", "GCTT"},
			"CT":                         []string{"TTCG", "GCTT"},
			"3U":                         []string{"GCTT", "CCCC"}, // should we cahnge this from GGTA to CCCC to conform with Protein Paintbox??
			"Ter":                        []string{"CCCC", "CGCT"},
			"3U + Ter":                   []string{"GCTT", "CGCT"},
			"3U + Ter + L1Dadaptor":      []string{"GCTT", "TAAT"},
			"CT + 3U + Ter + L1Dadaptor": []string{"TTCG", "TAAT"},
			"L1Dadaptor":                 []string{"CGCT", "TAAT"},
			"Ter + L1Dadaptor":           []string{"CCCC", "TAAT"},
		},
		"Level1": map[string][]string{
			"Device1": []string{"GAA", "ACC"},
			"Device2": []string{"ACC", "CTG"},
			"Device3": []string{"CTG", "GGT"},
		},
	},
	"Antibody": map[string]map[string][]string{
		"Heavy": map[string][]string{
			"Part1": []string{"GCG", "TCG"},
			"Part2": []string{"TGG", "CTG"},
			"Part3": []string{"CTG", "AAG"},
		},
		"Light": map[string][]string{
			"Part1": []string{"GCG", "TCG"},
			"Part2": []string{"TGG", "CTG"},
			"Part3": []string{"CTG", "AAG"},
		},
	},
	"MoClo_Raven": map[string]map[string][]string{
		"Level0": map[string][]string{
			"Pro":         []string{"GAGG", "TACT"},
			"5U":          []string{"TACT", "CCAT"},
			"5U(f)":       []string{"TACT", "CCAT"},
			"Pro + 5U(f)": []string{"GGAG", "CCAT"},
			"Pro + 5U":    []string{"GGAG", "AATG"},
			"NT1":         []string{"CCAT", "AATG"},
			"5U + NT1":    []string{"TACT", "AATG"},
			"CDS1":        []string{"AATG", "GCTT"},
			"CDS1 ns":     []string{"AATG", "TTCG"},
			"NT2":         []string{"AATG", "AGGT"},
			"SP":          []string{"AATG", "AGGT"},
			"CDS2 ns":     []string{"AGGT", "TTCG"},
			"CDS2":        []string{"AGGT", "GCTT"},
			"CT":          []string{"TTCG", "GCTT"},
			"3U":          []string{"GCTT", "GGTA"},
			"Ter":         []string{"GGTA", "CGCT"},
			"3U + Ter":    []string{"GCTT", "GCTT"}, // both same ! look into this
		},
	},
}

const (
	// for indexing part position based on part name/class
	VECTOR = iota
	PROMOTER
	RBS
	CDS
	TERMINATOR
)

// map of standard vector ends for various assembly standards
var Vectorends = map[string]map[string][]string{
	// array of strings returned correspond to [3'overhang and 5'overhang]
	"MoClo_Raven": map[string][]string{
		"Level0": []string{"AAGC", "CCTC"}, //
		"Level1": []string{"", ""},
	},
	"MoClo": map[string][]string{
		"Level0": []string{"CGCT", "GGAG"}, //
		"Level1": []string{"", ""},
	},
	"Synthace": map[string][]string{
		"Level0": []string{"GGT", "GAA"},
		"Level1": []string{"", ""},
	},
	"Custom": map[string][]string{
		"Level0": []string{"TAAT", "GTCG"},
		"Level1": []string{"GGT", "GAA"},
	},
	"Antibody": map[string][]string{
		"Heavy": []string{"GCG", "AAG"},
		"Light": []string{"", ""},
	},
	"Electra": map[string][]string{
		"Level0": []string{"GGT", "ATG"},
		"Level1": []string{"", ""},
	},
}

// map of enzymes used at each level of an assembly standard
var Enzymelookup = map[string]map[string]wtype.TypeIIs{
	// array of strings returned correspond to 5'overhang and 3'overhang
	"MoClo_Raven": map[string]wtype.TypeIIs{
		"Level0": BsaIenz,
		"Level1": BpiIenz,
	},
	"MoClo": map[string]wtype.TypeIIs{
		"Level0": BsaIenz,
		"Level1": BpiIenz,
	},
	"Custom": map[string]wtype.TypeIIs{
		"Level0": BsaIenz,
		"Level1": SapIenz,
	},
	"Antibody": map[string]wtype.TypeIIs{
		"Heavy": SapIenz,
		"Light": SapIenz,
	},
	"Electra": map[string]wtype.TypeIIs{
		"Level0": SapIenz,
	},
}

/*
var MoClo AssemblyStandard{
	[]AssemblyStandardLevel{{BsaIenz,"Level0"},{BpiIenz,"Level1"}}"Moclo"}

*/
type AssemblyStandardLevel struct {
	Enzyme    wtype.TypeIIs
	Levelname string
}

type AssemblyStandard struct {
	Endstable       map[string]map[string]map[int]map[int]string
	EnzymeTable     map[string]map[string]wtype.TypeIIs
	VectorEndstable map[string]map[string][]string // Vector 5prime can also be found in Endstable position 0
}
