// antha/AnthaStandardLibrary/Packages/enzymes/Ligation.go: Part of the Antha language
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
	"strconv"
	"strings"

	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/enzymes/lookup"
	. "github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/sequences"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/text"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	//features "github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/sequences/features"
)

func jointwoparts(upstreampart []Digestedfragment, downstreampart []Digestedfragment) (assembledfragments []Digestedfragment, plasmidproducts []wtype.DNASequence, err error) {

	sequencestojoin := make([]string, 0)

	assembledfragments = make([]Digestedfragment, 0)
	plasmidproducts = make([]wtype.DNASequence, 0)

	for _, upfragment := range upstreampart {
		for _, downfragment := range downstreampart {
			if RevComp(upfragment.BottomStickyend_5prime) == downfragment.TopStickyend_5prime && RevComp(downfragment.BottomStickyend_5prime) == upfragment.TopStickyend_5prime {
				sequencestojoin = append(sequencestojoin, upfragment.Topstrand, downfragment.Topstrand)
				dnastring := strings.Join(sequencestojoin, "")
				fullyassembledfragment := wtype.DNASequence{Nm: "simulatedassemblysequence", Seq: dnastring, Plasmid: true}
				plasmidproducts = append(plasmidproducts, fullyassembledfragment)
				sequencestojoin = make([]string, 0)
			}
			if upfragment.BottomStickyend_5prime == RevComp(downfragment.BottomStickyend_5prime) && downfragment.TopStickyend_5prime == RevComp(upfragment.TopStickyend_5prime) {
				sequencestojoin = append(sequencestojoin, upfragment.Topstrand, downfragment.Bottomstrand)
				dnastring := strings.Join(sequencestojoin, "")
				fullyassembledfragment := wtype.DNASequence{Nm: "simulatedassemblysequence", Seq: dnastring, Plasmid: true}
				plasmidproducts = append(plasmidproducts, fullyassembledfragment)
				sequencestojoin = make([]string, 0)
			}
			if /*upfragment.BottomStickyend_5prime == RevComp(downfragment.TopStickyend_5prime) ||*/ RevComp(upfragment.BottomStickyend_5prime) == downfragment.TopStickyend_5prime {
				sequencestojoin = append(sequencestojoin, upfragment.Topstrand, downfragment.Topstrand)
				dnastring := strings.Join(sequencestojoin, "")
				assembledfragment := Digestedfragment{dnastring, "", upfragment.TopStickyend_5prime, downfragment.TopStickyend_3prime, downfragment.BottomStickyend_5prime, upfragment.BottomStickyend_3prime}
				assembledfragments = append(assembledfragments, assembledfragment)
				sequencestojoin = make([]string, 0)
			}
			if upfragment.BottomStickyend_5prime == RevComp(downfragment.BottomStickyend_5prime) {
				sequencestojoin = append(sequencestojoin, upfragment.Topstrand, downfragment.Bottomstrand)
				dnastring := strings.Join(sequencestojoin, "")
				assembledfragment := Digestedfragment{dnastring, "", upfragment.TopStickyend_5prime, downfragment.BottomStickyend_3prime, downfragment.TopStickyend_5prime, upfragment.BottomStickyend_3prime}
				assembledfragments = append(assembledfragments, assembledfragment)
				sequencestojoin = make([]string, 0)
			}
		}
	}
	if len(assembledfragments) == 0 && len(plasmidproducts) == 0 {
		errstr := fmt.Sprintln("fragments aren't compatible, check ends",
			text.Print("upstream fragments:", upstreampart),
			text.Print("downstream fragements:", downstreampart))

		err = fmt.Errorf(errstr)
	}
	return assembledfragments, plasmidproducts, err
}

// key function for returning arrays of partially assembled fragments and fully assembled fragments from performing typeIIS assembly on a vector and a part
func Jointwopartsfromsequence(vector wtype.DNASequence, part1 wtype.DNASequence, enzyme wtype.TypeIIs) (assembledfragments []Digestedfragment, plasmidproducts []wtype.DNASequence) {
	doublestrandedpart1 := MakedoublestrandedDNA(part1)
	digestedpart1 := DigestionPairs(doublestrandedpart1, enzyme)

	doublestrandedvector := MakedoublestrandedDNA(vector)
	digestedvector := DigestionPairs(doublestrandedvector, enzyme)

	assembledfragments, plasmidproducts, _ = jointwoparts(digestedvector, digestedpart1)

	return assembledfragments, plasmidproducts
}

func rotate_vector(vector wtype.DNASequence, enzyme wtype.TypeIIs) (wtype.DNASequence, error) {
	ret := vector.Dup()

	// the purpose of this is to ensure the RE sites go ---> xxxx <---

	// we just ensure the first one is first in the sequence... if there's more than one
	// it's not our problem

	ix := strings.Index(strings.ToUpper(ret.Seq), strings.ToUpper(enzyme.RecognitionSequence))

	if ix == -1 {
		err := fmt.Errorf("No restriction sites found in vector - cannot rotate")
		return ret, err
	}

	/*thingsfound := FindSeqsinSeqs(ret.Seq, []string{enzyme.RecognitionSequence})

	if len(thingsfound) == 0 {
		err := fmt.Errorf("No restriction sites found in vector - cannot rotate")
		return ret, err
	}
	if len(thingsfound) != 2 {
		errstr := fmt.Sprint(len(thingsfound), "restriction sites found in vector - cannot rotate")

		err := fmt.Errorf(errstr)
		return ret, err
	}

	if len(thingsfound[0].Positions) > 1 {

		errstr := fmt.Sprint(len(thingsfound[0].Positions), "restriction sites found in vector - cannot rotate")

		err := fmt.Errorf(errstr)
		return ret, err
	}*/
	/*
			if thingsfound[0].Reverse {
				err := fmt.Errorf("first site is reverse")
				return ret, err
			}
			if thingsfound[1].Reverse {
				err := fmt.Errorf("second site is reverse")
				return ret, err
			}

		ix := thingsfound[0].Positions[0]
	*/
	newseq := ""

	newseq += ret.Seq[ix:]
	newseq += ret.Seq[:ix]

	ret.Seq = newseq

	return ret, nil
}

// key function for returning an error message, arrays of partially assembled fragments and fully assembled fragments from performing typeIIS assembly on a vector and array of parts
func JoinXnumberofparts(vector wtype.DNASequence, partsinorder []wtype.DNASequence, enzyme wtype.TypeIIs) (assembledfragments []Digestedfragment, plasmidproducts []wtype.DNASequence, err error) {

	if vector.Seq == "" {
		err = fmt.Errorf("No Vector sequence found")
		return assembledfragments, plasmidproducts, err
	}
	// there are two cases: either the vector comes in same way parts do
	// i.e. SAPI--->xxxx<---IPAS
	// OR it comes in the other way round
	// i.e. xxxx<---IPASyyyySAPI--->zzzz
	// we have either to rotate the vector or tolerate this
	// probably best to rotate first

	rotatedvector, err := rotate_vector(vector, enzyme)

	if err != nil {
		return assembledfragments, plasmidproducts, err
	}

	doublestrandedvector := MakedoublestrandedDNA(rotatedvector)
	digestedvector := DigestionPairs(doublestrandedvector, enzyme)

	if len(partsinorder) == 0 {
		return nil, nil, fmt.Errorf("No parts found")
	}
	if len(partsinorder[0].Seq) == 0 {
		name := partsinorder[0].Nm
		errorstring := name + " has no sequence"
		err = fmt.Errorf(errorstring)
		return assembledfragments, plasmidproducts, err
	}
	doublestrandedpart := MakedoublestrandedDNA(partsinorder[0])
	digestedpart := DigestionPairs(doublestrandedpart, enzyme)

	var newerr error
	assembledfragments, plasmidproducts, newerr = jointwoparts(digestedvector, digestedpart)
	if newerr != nil {
		message := fmt.Sprint(vector.Nm, " and ", partsinorder[0].Nm, ": ", newerr.Error())
		err = fmt.Errorf(message)
		return
	}

	for i := 1; i < len(partsinorder); i++ {
		if len(partsinorder[i].Seq) == 0 {
			name := partsinorder[i].Nm
			errorstring := name + " has no sequence"
			err = fmt.Errorf(errorstring)
			return assembledfragments, plasmidproducts, err
		}

		doublestrandedpart = MakedoublestrandedDNA(partsinorder[i])
		digestedpart := DigestionPairs(doublestrandedpart, enzyme)
		//for _, newfragments := range assembledfragments {

		assembledfragments, plasmidproducts, newerr = jointwoparts(assembledfragments, digestedpart)
		//err = newerr

		if newerr != nil {
			//	if err != nil {
			message := fmt.Sprint(partsinorder[i-1].Nm, " and ", partsinorder[i].Nm, ": ", newerr.Error())
			err = fmt.Errorf(message)
			//	} else {
			//		message := fmt.Sprint(partsinorder[i - 1].Nm, " and ", partsinorder[i].Nm, ": ", newerr.Error())
			//		err = fmt.Errorf(message)
			//	}
			return
		}
		//}
	}

	partnames := make([]string, 0)

	for _, part := range partsinorder {
		partnames = append(partnames, part.Nm)
	}

	for _, plasmidproduct := range plasmidproducts {

		plasmidproduct.Nm = vector.Nm + "_" + strings.Join(partnames, "_")
	}

	return assembledfragments, plasmidproducts, err
}

/*func JoinAnnotatedparts(vector wtype.DNASequence, partsinorder []wtype.DNASequence, enzyme TypeIIs) (assembledfragments []Digestedfragment, plasmidproducts []wtype.DNASequence) {

	doublestrandedvector := MakedoublestrandedDNA(vector)
	digestedvector := DigestionPairs(doublestrandedvector, enzyme)

	doublestrandedpart := MakedoublestrandedDNA(partsinorder[0])
	digestedpart := DigestionPairs(doublestrandedpart, enzyme)
	assembledfragments, plasmidproducts = Jointwoparts(digestedvector, digestedpart)
	//// fmt.Println("vector + part1 product = ", assembledfragments, plasmidproducts)
	for i := 1; i < len(partsinorder); i++ {
		doublestrandedpart = MakedoublestrandedDNA(partsinorder[i])
		digestedpart := DigestionPairs(doublestrandedpart, enzyme)
		//for _, newfragments := range assembledfragments {
		assembledfragments, plasmidproducts = Jointwoparts(assembledfragments, digestedpart)
		//}
	}
	return assembledfragments, plasmidproducts
}
*/

// struct containing all information required to use AssemblySimulator function
type Assemblyparameters struct {
	Constructname string              `json:"construct_name"`
	Enzymename    string              `json:"enzyme_name"`
	Vector        wtype.DNASequence   `json:"vector"`
	Partsinorder  []wtype.DNASequence `json:"parts_in_order"`
}

/*type AA_DNA_Assemblyparameters struct {
	Constructname string
	Enzymename    string
	Vector        wtype.DNASequence
	Partsinorder  []wtype.BioSequence
}*/

// Simulate assembly success; returns status, number of correct assemblies, any sites found
func Assemblysimulator(assemblyparameters Assemblyparameters) (s string, successfulassemblies int, sites []Restrictionsites, newDNASequence wtype.DNASequence, err error) {

	// fetch enzyme properties from map (this is basically a look up table for those who don't know)
	successfulassemblies = 0
	enzymename := strings.ToUpper(assemblyparameters.Enzymename)

	// should change this to rebase lookup; what happens if this fails?
	//enzyme := TypeIIsEnzymeproperties[enzymename]
	enzyme, err := lookup.TypeIIsLookup(enzymename)

	if err != nil {
		return s, successfulassemblies, sites, newDNASequence, err
	}

	// need to expand this to include other enzyme possibilities
	if enzyme.Class != "TypeIIs" { // enzyme.Name != "SapI" && enzyme.Name != "BsaI" && enzyme.Name != "BpiI" {
		s = fmt.Sprint(enzymename, ": Incorrect Enzyme or no enzyme specified")
		err = fmt.Errorf(s)
		return s, successfulassemblies, sites, newDNASequence, err
	}

	//assemble (note that sapIenz is found in package enzymes)
	failedassemblies, plasmidproductsfromXprimaryseq, err := JoinXnumberofparts(assemblyparameters.Vector, assemblyparameters.Partsinorder, enzyme)

	if err != nil {
		//s = "Failure Joining fragments after digestion" //
		err = fmt.Errorf("Failure Joining fragments after digestion: %s", err)
		s = err.Error()
		return s, successfulassemblies, sites, newDNASequence, err
	}

	enzy, err := lookup.EnzymeLookup(enzymename)

	if err != nil {
		return s, successfulassemblies, sites, newDNASequence, err
	}

	if len(plasmidproductsfromXprimaryseq) == 1 {

		sites = Restrictionsitefinder(plasmidproductsfromXprimaryseq[0], []wtype.RestrictionEnzyme{BsaI, SapI, enzy})
		newDNASequence = plasmidproductsfromXprimaryseq[0]
	}
	// returns first plasmid in array! should be changed later!
	if len(plasmidproductsfromXprimaryseq) > 1 {
		sites = make([]Restrictionsites, 0)
		for i := 0; i < len(plasmidproductsfromXprimaryseq); i++ {
			sitesperplasmid := Restrictionsitefinder(plasmidproductsfromXprimaryseq[i], []wtype.RestrictionEnzyme{BsaI, SapI, enzy})
			for _, site := range sitesperplasmid {
				sites = append(sites, site)
			}
		}
		//return first for now
		newDNASequence = plasmidproductsfromXprimaryseq[0]
	}

	s = "hmmm I'm confused, this doesn't seem to make any sense"

	if len(plasmidproductsfromXprimaryseq) == 0 && len(failedassemblies) == 0 {
		err = fmt.Errorf("Nope! this construct won't work: ", err)
		s = err.Error()
	}
	if len(plasmidproductsfromXprimaryseq) == 1 {
		s = "Yay! this should work"
		successfulassemblies = successfulassemblies + 1
	}

	if len(plasmidproductsfromXprimaryseq) > 1 {

		err = fmt.Errorf("Yay! this should work but there seems to be more than one possible plasmid which could form %s", err)
		s = err.Error()
	}

	if len(plasmidproductsfromXprimaryseq) == 0 && len(failedassemblies) != 0 {

		s = fmt.Sprint("Ooh, only partial assembly expected: ", assemblyparameters.Partsinorder[(len(assemblyparameters.Partsinorder)-1)].Nm, " and ", assemblyparameters.Vector.Nm, ": ", "Not compatible, check ends")

		err = fmt.Errorf(s)
		//err = fmt.Errorf(funk, err.Error())
		//s = err.Error()
	}

	if s != "Yay! this should work" {
		err = fmt.Errorf(s)
	}

	return s, successfulassemblies, sites, newDNASequence, err
}

// MultipleAssemblies will perform simulated assemblies on multiple constructs
// and return a description of whether each was successful and how many are
// expected to work
func MultipleAssemblies(parameters []Assemblyparameters) (s string, successfulassemblies int, errors map[string]string, seqs []wtype.DNASequence) {

	seqs = make([]wtype.DNASequence, 0)
	errors = make(map[string]string) // construct to error

	for _, construct := range parameters {
		output, _, _, seq, err := Assemblysimulator(construct)

		seqs = append(seqs, seq)

		if err == nil {
			successfulassemblies += 1
			continue
		} else {

			errors[construct.Constructname] = err.Error()

			if strings.Contains(err.Error(), "Failure Joining fragments after digestion") {
				sitesperpart := make([]Restrictionsites, 0)
				constructsitesstring := make([]string, 0)
				constructsitesstring = append(constructsitesstring, output)
				sitestring := ""
				enzyme, err := lookup.EnzymeLookup(construct.Enzymename)
				if err != nil {

					originalerror := errors[construct.Constructname]

					errors[construct.Constructname] = originalerror + " and " + err.Error()
				}
				sitesperpart = Restrictionsitefinder(construct.Vector, []wtype.RestrictionEnzyme{enzyme})

				if sitesperpart[0].Numberofsites != 2 {
					// need to loop through sitesperpart

					sitepositions := SitepositionString(sitesperpart[0])
					sitestring = "For " + construct.Vector.Nm + ": " + strconv.Itoa(sitesperpart[0].Numberofsites) + " sites found at positions: " + sitepositions
					constructsitesstring = append(constructsitesstring, sitestring)
				}

				for _, part := range construct.Partsinorder {
					sitesperpart = Restrictionsitefinder(part, []wtype.RestrictionEnzyme{enzyme})
					if sitesperpart[0].Numberofsites != 2 {
						sitepositions := SitepositionString(sitesperpart[0])
						positions := ""
						if sitesperpart[0].Numberofsites != 0 {
							positions = fmt.Sprint("at positions:", sitepositions)
						}
						sitestring = fmt.Sprint("For ", part.Nm, ": ", strconv.Itoa(sitesperpart[0].Numberofsites), " sites were found ", positions)
						constructsitesstring = append(constructsitesstring, sitestring)
					}

				}
				if len(constructsitesstring) != 1 {
					message := strings.Join(constructsitesstring, "; ")
					err = fmt.Errorf(message)
				}
			}

			s = err.Error()

			/*if errors == nil {
				errors = make(map[string]string)
			}*/

			if _, ok := errors[construct.Constructname]; !ok {
				errors[construct.Constructname] = s
			}

		}
	}

	if successfulassemblies == len(parameters) {
		s = "success, all assemblies seem to work"
	}
	return
}
