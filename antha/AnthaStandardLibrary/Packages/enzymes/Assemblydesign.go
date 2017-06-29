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
	"sort"
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

func VectorEnds(vector wtype.DNASequence, enzyme wtype.TypeIIs) (desiredstickyend5prime string, vector3primestickyend string) {
	// find sticky ends from cutting vector with enzyme

	fragments, stickyends5, _ := TypeIIsdigest(vector, enzyme)

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
	}
	return
}

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

		sites, sticky5s, sticky3s := CheckForExistingTypeIISEnds(parts[i], enzyme)
		//sites, _, _ := CheckForExistingTypeIISEnds(parts[i], enzyme)

		if sites == 0 {
			partwithends = AddCustomEnds(parts[i], enzyme, desiredstickyend5prime, desiredstickyend3prime)
		} else if sites == 2 && InSlice(desiredstickyend5prime, sticky5s) && InSlice(desiredstickyend3prime, sticky3s) {
			partwithends = parts[i]
		} else {
			panic(fmt.Sprint("cutting part ", parts[i], " with", enzyme, " results in ", sites, "cut sites with 5 prime fragment overhangs: ", sticky5s, " and 3 prime fragment overhangs: ", sticky3s, ". Wanted: 5prime: ", desiredstickyend5prime, " 3prime: ", desiredstickyend3prime))
		}
		partwithends.Nm = parts[i].Nm
		partswithends = append(partswithends, partwithends)

		desiredstickyend5prime = Suffix(partwithends.Seq, enzyme.RestrictionEnzyme.EndLength)

	}

	return partswithends
}

// This func adds an upstream (5') adaptor for making a level 0 part compatible for Level 1 hierarchical assembly, specifying the desired level 1 class the level0 part should be made into.
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

	bitToAdd3prime := Makeoverhang(enzyme, "3prime", bitToAdd, ChooseSpacer(enzyme.Topstrand3primedistancefromend, "", []string{}))
	partWithEnds := Addoverhang(part.Seq, bitToAdd3prime, "3prime")

	newpart = part.Dup()
	newpart.Seq = partWithEnds
	return
}

// look up enzyme according to assembly standard and level.
// Errors will be returned if an entry is missing
func lookUpEnzyme(assemblyStandard AssemblyStandard, level string) (enzyme wtype.TypeIIs, err error) {

	assemblyLevel, err := assemblyStandard.GetLevel(level)

	if err != nil {
		return enzyme, err
	}

	enzyme = assemblyLevel.GetEnzyme()

	return enzyme, nil
}

// look up overhangs according to assembly standard, level and part class.
// Errors will be returned if an entry is missing or overhangs are found to be empty
func lookUpOverhangs(assemblyStandard AssemblyStandard, level string, class string) (upstream string, downstream string, err error) {

	assemblyLevel, err := assemblyStandard.GetLevel(level)

	if err != nil {
		return "", "", err
	}

	ends, err := assemblyLevel.GetPartOverhangs(class)

	return ends.Upstream, ends.Downstream, err
}

// Adds sticky ends to dna part according to the class identifier (e.g. PRO, 5U, CDS)
func AddStandardStickyEndsfromClass(part wtype.DNASequence, assemblyStandard AssemblyStandard, level string, class string) (partWithEnds wtype.DNASequence, err error) {

	enzyme, err := lookUpEnzyme(assemblyStandard, level)

	if err != nil {
		return partWithEnds, err
	}
	bittoadd, bittoadd3, err := lookUpOverhangs(assemblyStandard, level, class)

	if err != nil {
		return partWithEnds, err
	}

	bittoadd = findMinimumAdditional5PrimeAddition(bittoadd, part)

	bittoadd5prime := Makeoverhang(enzyme, "5prime", bittoadd, ChooseSpacer(enzyme.Topstrand3primedistancefromend, "", []string{}))

	partwith5primeend := Addoverhang(part.Seq, bittoadd5prime, "5prime")

	bittoadd3 = findMinimumAdditional3PrimeAddition(bittoadd3, part)

	bittoadd3prime := Makeoverhang(enzyme, "3prime", bittoadd3, ChooseSpacer(enzyme.Topstrand3primedistancefromend, "", []string{}))

	partwithends := Addoverhang(partwith5primeend, bittoadd3prime, "3prime")

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
		fmt.Println("Part: ", part.Name(), "Desired 3 prime:", desiredstickyend3prime, "Truncated: ", truncated)
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
	fmt.Println("Part: ", part.Name(), "Desired 3 prime:", desiredstickyend3prime, "Final: ", bittoadd)
	return
}

func findMinimumAdditional5PrimeAddition(desiredstickyend5prime string, part wtype.DNASequence) (bittoadd string) {
	// This code will look for subparts of a standard overhang to add the minimum number of additional nucleotides with a partial match e.g. AATG contains ATG only so we just add A

	var present string

	for i := len(desiredstickyend5prime) - 1; i >= 0; i-- {
		truncated := desiredstickyend5prime[i:]
		fmt.Println("Part: ", part.Name(), "Desired 5 prime:", desiredstickyend5prime, "Truncated: ", truncated)
		if strings.HasPrefix(upper(part.Seq), upper(truncated)) {
			present = truncated
		}
	}
	if len(present) == len(desiredstickyend5prime) {
		bittoadd = ""
	} else if len(present) > 0 {
		fmt.Println(desiredstickyend5prime, present)
		bittoadd = desiredstickyend5prime[:len(present)+1]
	} else {
		bittoadd = desiredstickyend5prime
	}
	fmt.Println("Part: ", part.Name(), "Desired 5 prime:", desiredstickyend5prime, "Final: ", bittoadd)
	return
}

// Adds ends to the part sequence based upon enzyme chosen and the desired overhangs after digestion
func AddCustomEnds(part wtype.DNASequence, enzyme wtype.TypeIIs, desiredstickyend5prime string, desiredstickyend3prime string) (Partwithends wtype.DNASequence) {

	fmt.Println("Old part: ", Partwithends)

	///

	bittoadd := findMinimumAdditional5PrimeAddition(desiredstickyend5prime, part)

	bittoadd5prime := Makeoverhang(enzyme, "5prime", bittoadd, ChooseSpacer(enzyme.Topstrand3primedistancefromend, "", []string{}))
	partwith5primeend := Addoverhang(part.Seq, bittoadd5prime, "5prime")

	bittoadd3 := findMinimumAdditional3PrimeAddition(desiredstickyend3prime, part)

	bittoadd3prime := Makeoverhang(enzyme, "3prime", bittoadd3, ChooseSpacer(enzyme.Topstrand3primedistancefromend, "", []string{}))
	partwithends := Addoverhang(partwith5primeend, bittoadd3prime, "3prime")

	Partwithends.Nm = part.Nm
	Partwithends.Plasmid = part.Plasmid
	Partwithends.Seq = partwithends
	fmt.Println("New part: ", Partwithends)
	return Partwithends
}

// Add compatible ends to an array of parts based on the rules of a typeIIS assembly standard
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

// Assembly standards
var availableStandards = map[string]AssemblyStandard{
	"Custom":      customStandard,
	"MoClo":       customStandard,
	"MoClo_Raven": customStandard,
	"Antibody":    customStandard,
}

func allStandards() (standards []string) {
	for k, _ := range availableStandards {
		standards = append(standards, k)
	}
	sort.Strings(standards)
	return
}

// Lookup TypeIIs assembly standard by name
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
		"Level0": AssemblyLevel{
			Enzyme: BsaIenz,
			PartOverhangs: map[string]StandardOverhangs{
				"Pro":         StandardOverhangs{"GGAG", "TACT"},
				"5U":          StandardOverhangs{"TACT", "CCAT"},
				"5U(f)":       StandardOverhangs{"TACT", "CCAT"},
				"Pro + 5U(f)": StandardOverhangs{"GGAG", "CCAT"},
				"Pro + 5U":    StandardOverhangs{"GGAG", "AATG"},
				"NT1":         StandardOverhangs{"CCAT", "AATG"},
				"5U + NT1":    StandardOverhangs{"TACT", "AATG"},
				"CDS1":        StandardOverhangs{"AATG", "GCTT"},
				"CDS1 ns":     StandardOverhangs{"AATG", "TTCG"},
				"NT2":         StandardOverhangs{"AATG", "AGGT"},
				"SP":          StandardOverhangs{"AATG", "AGGT"},
				"CDS2 ns":     StandardOverhangs{"AGGT", "TTCG"},
				"CDS2":        StandardOverhangs{"AGGT", "GCTT"},
				"CT":          StandardOverhangs{"TTCG", "GCTT"},
				"3U":          StandardOverhangs{"GCTT", "GGTA"},
				"Ter":         StandardOverhangs{"GGTA", "CGCT"},
				"3U + Ter":    StandardOverhangs{"GCTT", "CGCT"},
			},
			EntryVectorEnds: StandardOverhangs{"TAAT", "GTCG"},
		},
		"Level1": AssemblyLevel{
			Enzyme:          BpiIenz,
			PartOverhangs:   map[string]StandardOverhangs{},
			EntryVectorEnds: StandardOverhangs{"", ""},
		},
	},
}

var mocloRavenStandard = AssemblyStandard{
	Name: "MoClo_Raven",
	Levels: map[string]AssemblyLevel{
		"Level0": AssemblyLevel{
			Enzyme: BsaIenz,
			PartOverhangs: map[string]StandardOverhangs{
				"Pro":         StandardOverhangs{"GAGG", "TACT"},
				"5U":          StandardOverhangs{"TACT", "CCAT"},
				"5U(f)":       StandardOverhangs{"TACT", "CCAT"},
				"Pro + 5U(f)": StandardOverhangs{"GGAG", "CCAT"},
				"Pro + 5U":    StandardOverhangs{"GGAG", "AATG"},
				"NT1":         StandardOverhangs{"CCAT", "AATG"},
				"5U + NT1":    StandardOverhangs{"TACT", "AATG"},
				"CDS1":        StandardOverhangs{"AATG", "GCTT"},
				"CDS1 ns":     StandardOverhangs{"AATG", "TTCG"},
				"NT2":         StandardOverhangs{"AATG", "AGGT"},
				"SP":          StandardOverhangs{"AATG", "AGGT"},
				"CDS2 ns":     StandardOverhangs{"AGGT", "TTCG"},
				"CDS2":        StandardOverhangs{"AGGT", "GCTT"},
				"CT":          StandardOverhangs{"TTCG", "GCTT"},
				"3U":          StandardOverhangs{"GCTT", "GGTA"},
				"Ter":         StandardOverhangs{"GGTA", "CGCT"},
				"3U + Ter":    StandardOverhangs{"GCTT", "GCTT"}, // both same ! look into this
			},
			EntryVectorEnds: StandardOverhangs{"AAGC", "CCTC"},
		},
		"Level1": AssemblyLevel{
			Enzyme:          BpiIenz,
			PartOverhangs:   map[string]StandardOverhangs{},
			EntryVectorEnds: StandardOverhangs{"", ""},
		},
	},
}

var customStandard = AssemblyStandard{
	Name: "Custom",
	Levels: map[string]AssemblyLevel{
		"Level0": AssemblyLevel{
			Enzyme: BsaIenz,
			PartOverhangs: map[string]StandardOverhangs{
				"L1Uadaptor":                  StandardOverhangs{"GTCG", "GGAG"}, // adaptor to add SapI sites to clone into level 1 vector
				"L1Uadaptor + Pro":            StandardOverhangs{"GTCG", "TTTT"}, // adaptor to add SapI sites to clone into level 1 vector
				"L1Uadaptor + Pro(MoClo)":     StandardOverhangs{"GTCG", "TACT"}, // original MoClo overhang of TACT
				"L1Uadaptor + Pro + 5U":       StandardOverhangs{"GTCG", "CCAT"}, // adaptor to add SapI sites to clone into level 1 vector
				"L1Uadaptor + Pro + 5U + NT1": StandardOverhangs{"GTCG", "TATG"}, // adaptor to add SapI sites to clone into level 1 vector
				"Pro":                               StandardOverhangs{"GGAG", "TTTT"}, //changed from MoClo TACT to TTTT to conform with Protein Paintbox
				"5U":                                StandardOverhangs{"TTTT", "CCAT"}, //changed from MoClo TACT to TTTT to conform with Protein Paintbox
				"5U(f)":                             StandardOverhangs{"TTTT", "CCAT"}, //changed from MoClo TACT to TTTT to conform with Protein Paintbox
				"Pro + 5U(f)":                       StandardOverhangs{"GGAG", "CCAT"},
				"Pro + 5U":                          StandardOverhangs{"GGAG", "CCAT"},
				"Pro + 5U + NT1":                    StandardOverhangs{"GGAG", "TATG"}, //changed AATG to TATG to work with Kosuri paper RBSs
				"NT1":                               StandardOverhangs{"CCAT", "TATG"}, //changed AATG to TATG to work with Kosuri paper RBSs
				"5U + NT1":                          StandardOverhangs{"TTTT", "TATG"}, //changed AATG to TATG to work with Kosuri paper RBSs
				"5U(MoClo) + NT1":                   StandardOverhangs{"TACT", "TATG"}, //original MoClo overhang of TACT
				"5U + NT1 + CDS1":                   StandardOverhangs{"TTTT", "GCTT"}, //changed from MoClo TACT to TTTT to conform with Protein Paintbox
				"5U + NT1 + CDS1 + 3U":              StandardOverhangs{"TTTT", "CCCC"}, //changed from MoClo TACT to TTTT to conform with Protein Paintbox and changed GGTA to CCCC to conform with Protein Paintbox
				"CDS1":                              StandardOverhangs{"TATG", "GCTT"}, //changed AATG to TATG to work with Kosuri paper RBSs
				"CDS1 + 3U":                         StandardOverhangs{"TATG", "CCCC"}, //changed AATG to TATG to work with Kosuri paper RBSs and changed GGTA to CCCC to conform with Protein Paintbox
				"CDS1 + 3U(MoClo)":                  StandardOverhangs{"TATG", "GGTA"}, //original MoClo overhang of GGTA
				"CDS1 ns":                           StandardOverhangs{"TATG", "TTCG"}, //changed AATG to TATG to work with Kosuri paper RBSs
				"CDS1 + CT + 3U + Ter + L1Dadaptor": StandardOverhangs{"TATG", "TAAT"}, //changed AATG to TATG to work with Kosuri paper RBSs
				"NT2":                        StandardOverhangs{"TATG", "AGGT"}, //changed AATG to TATG to work with Kosuri paper RBSs
				"SP":                         StandardOverhangs{"TATG", "AGGT"}, //changed AATG to TATG to work with Kosuri paper RBSs
				"CDS2 ns":                    StandardOverhangs{"AGGT", "TTCG"},
				"CDS2":                       StandardOverhangs{"AGGT", "GCTT"},
				"CT":                         StandardOverhangs{"TTCG", "GCTT"},
				"3U":                         StandardOverhangs{"GCTT", "CCCC"}, //changed GGTA to CCCC to conform with Protein Paintbox
				"Ter":                        StandardOverhangs{"CCCC", "CGCT"},
				"3U + Ter":                   StandardOverhangs{"GCTT", "CGCT"},
				"3U + Ter + L1Dadaptor":      StandardOverhangs{"GCTT", "TAAT"},
				"CT + 3U + Ter + L1Dadaptor": StandardOverhangs{"TTCG", "TAAT"},
				"L1Dadaptor":                 StandardOverhangs{"CGCT", "TAAT"},
				"Ter + L1Dadaptor":           StandardOverhangs{"CCCC", "TAAT"},
				"Ter(MoClo) + L1Dadaptor":    StandardOverhangs{"GGTA", "TAAT"},
			},
			EntryVectorEnds: StandardOverhangs{"TAAT", "GTCG"},
		},
		"Level1": AssemblyLevel{
			Enzyme: SapIenz,
			PartOverhangs: map[string]StandardOverhangs{
				"Device1": StandardOverhangs{"GAA", "ACC"},
				"Device2": StandardOverhangs{"ACC", "CTG"},
				"Device3": StandardOverhangs{"CTG", "GGT"},
			},
			EntryVectorEnds: StandardOverhangs{"GGT", "GAA"},
		},
	},
}

var antibodyStandard = AssemblyStandard{
	Name: "Antibody",
	Levels: map[string]AssemblyLevel{
		"Heavy": AssemblyLevel{
			Enzyme: SapIenz,
			PartOverhangs: map[string]StandardOverhangs{
				"Part1": StandardOverhangs{"GCG", "TCG"},
				"Part2": StandardOverhangs{"TGG", "CTG"},
				"Part3": StandardOverhangs{"CTG", "AAG"},
			},
			EntryVectorEnds: StandardOverhangs{"GCG", "AAG"},
		},
		"Light": AssemblyLevel{
			Enzyme: SapIenz,
			PartOverhangs: map[string]StandardOverhangs{
				"Part1": StandardOverhangs{"GCG", "TCG"},
				"Part2": StandardOverhangs{"TGG", "CTG"},
				"Part3": StandardOverhangs{"CTG", "AAG"},
			},
			EntryVectorEnds: StandardOverhangs{"GCG", "AAG"},
		},
	},
}

type AssemblyStandard struct {
	Name   string
	Levels map[string]AssemblyLevel
}

func (l AssemblyStandard) Enzyme(level string) (enz wtype.TypeIIs, err error) {

	assemblyLevel, err := l.GetLevel(level)
	if err != nil {
		return enz, err
	}
	enz = assemblyLevel.GetEnzyme()
	return
}

func (l AssemblyStandard) LevelNames() []string {

	var ls []string
	for lev, _ := range l.Levels {
		ls = append(ls, lev)
	}

	sort.Strings(ls)

	return ls
}

func (l AssemblyStandard) GetLevel(level string) (assemblyLevel AssemblyLevel, err error) {

	assemblyLevel, found := l.Levels[level]
	if !found {
		return assemblyLevel, fmt.Errorf("No level %s found for assembly standard %s, found %+v", level, l.Name, l.LevelNames())
	}
	return
}

type AssemblyLevel struct {
	Enzyme          wtype.TypeIIs
	PartOverhangs   map[string]StandardOverhangs
	EntryVectorEnds StandardOverhangs // Vector 5prime can also be found in Endstable position 0
}

type StandardOverhangs struct {
	Upstream   string
	Downstream string
}

func (l AssemblyLevel) AnnotationOptions() []string {
	var ls []string
	for class, _ := range l.PartOverhangs {
		ls = append(ls, class)
	}

	sort.Strings(ls)
	return ls
}

func (l AssemblyLevel) GetEnzyme() wtype.TypeIIs {
	return l.Enzyme
}

func (l AssemblyLevel) GetVectorEnds() StandardOverhangs {
	return l.EntryVectorEnds
}

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
