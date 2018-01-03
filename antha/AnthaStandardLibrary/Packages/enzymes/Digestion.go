// antha/AnthaStandardLibrary/Packages/enzymes/Digestion.go: Part of the Antha language
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

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/search"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/sequences"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/text"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
)

const (
	blunt string = "blunt"
)

// RestrictionSites holds information on restriction sites found in a DNA sequence
// for a specified RestrictionEnzyme.
// todo: refactor
type RestrictionSites struct {
	Enzyme              wtype.RestrictionEnzyme
	RecognitionSequence string
	SiteFound           bool
	NumberOfSites       int
	ForwardPositions    []int
	ReversePositions    []int
}

/*
func (sites RestrictionSites) SiteFound() bool {
	if len(sites.Positions) > 0 {
		return true
	}
	return false
}

func (sites RestrictionSites) SiteFound() bool {
	if len(sites.Positions) > 0 {
		return true
	}
	return false
}

func (sites RestrictionSites) SiteFound() bool {
	if len(sites.Positions) > 0 {
		return true
	}
	return false
}

func (sites RestrictionSites) SiteFound() bool {
	if len(sites.Positions) > 0 {
		return true
	}
	return false
}
*/

// Positions returns a set of restriction site positions.
// Valid arguments which may be specified are "FWD", "REV or "ALL" to return the correct type of positions.
// todo: refactor
func (sites *RestrictionSites) Positions(fwdRevorNil string) (positions []int) {
	if strings.ToUpper(fwdRevorNil) == strings.ToUpper("FWD") {
		positions = sites.ForwardPositions
	} else if strings.ToUpper(fwdRevorNil) == strings.ToUpper("REV") {
		positions = sites.ReversePositions
	} else if strings.ToUpper(fwdRevorNil) == strings.ToUpper("") ||
		strings.ToUpper(fwdRevorNil) == strings.ToUpper("ALL") {
		positions = make([]int, 0)
		for _, pos := range sites.ForwardPositions {
			positions = append(positions, pos)
		}
		for _, pos := range sites.ReversePositions {
			positions = append(positions, pos)
		}
	}
	return
}

// SitepositionString returns a report of restriction sites found as a string
// todo: deprecate
func SitepositionString(sitesperpart RestrictionSites) (sitepositions string) {
	Num := make([]string, 0)

	for _, site := range sitesperpart.ForwardPositions {
		Num = append(Num, strconv.Itoa(site))
	}
	for _, site := range sitesperpart.ReversePositions {
		Num = append(Num, strconv.Itoa(site))
	}

	sort.Strings(Num)
	sitepositions = strings.Join(Num, ", ")
	return
}

// Restrictionsitefinder finds restriction sites of specified restriction enzymes in a sequence and return the information as a set of ResrictionSites.
func Restrictionsitefinder(sequence wtype.DNASequence, enzymelist []wtype.RestrictionEnzyme) (sites []RestrictionSites) {

	sites = make([]RestrictionSites, 0)

	for _, enzyme := range enzymelist {
		var enzymesite RestrictionSites
		//var siteafterwobble Restrictionsites
		enzymesite.Enzyme = enzyme
		enzymesite.RecognitionSequence = strings.ToUpper(enzyme.RecognitionSequence)
		sequence.Seq = strings.ToUpper(sequence.Seq)

		wobbleproofrecognitionoptions := sequences.Wobble(enzymesite.RecognitionSequence)

		for _, wobbleoption := range wobbleproofrecognitionoptions {

			options := search.FindAll(sequence.Seq, wobbleoption)
			for _, option := range options {
				if option != 0 {
					enzymesite.ForwardPositions = append(enzymesite.ForwardPositions, option)
				}
			}
			if enzyme.RecognitionSequence != strings.ToUpper(sequences.RevComp(wobbleoption)) {
				revoptions := search.FindAll(sequence.Seq, sequences.RevComp(wobbleoption))
				for _, option := range revoptions {
					if option != 0 {
						enzymesite.ReversePositions = append(enzymesite.ReversePositions, option)
					}
				}

			}
			enzymesite.NumberOfSites = len(enzymesite.ForwardPositions) + len(enzymesite.ReversePositions)
			if enzymesite.NumberOfSites > 0 {
				enzymesite.SiteFound = true
			}

		}

		sites = append(sites, enzymesite)
	}

	return sites
}

// Digestedfragment object carrying info on a fragment following digestion
type Digestedfragment struct {
	// Sequence of top strand
	Topstrand string
	// Sequence of bottom strand
	Bottomstrand string

	// FivePrimeTopStrandStickyend is any over hang at the 5' end of the coding strand
	// This will be greater than "" if the fragment was generated by cutting with an enzyme which leaves a 5' overhang.
	// The complementary sticky end will be on the FivePrimeBottomStrandStickyend.
	// If this is not "" this will be the first part of the TopStrand sequence.
	FivePrimeTopStrandStickyend string

	// FivePrimeBottomStrandStickyend is any over hang at the 5' end of the complementary strand
	// This may be referred to as an underhang.
	// This will be greater than "" if the fragment was generated by cutting with an enzyme which leaves a 5' overhang.
	// The complementary sticky end will be on the FivePrimeTopStrandStickyend.
	// This sequence or it's reverse complement will not be present on the end of the Top strand.
	FivePrimeBottomStrandStickyend string

	// TopStickyend_3prime is any over hang at the 3' end of the coding strand.
	// This may be referred to as an underhang.
	// This will be greater than "" if the fragment was generated by cutting with an enzyme which leaves a 3' overhang.
	// The complementary sticky end will be on the BottomStickyend_3prime.
	// This sequence will be present on the end of the Top strand.
	ThreePrimeTopStrandStickyend string

	// BottomStickyend_3prime is any over hang at the 3' end of the complementary strand.
	// This may be referred to as an underhang.
	// This will be greater than "" if the fragment was generated by cutting with an enzyme which leaves a 3' overhang.
	// The complementary sticky end will be on the TopStickyend_3prime.
	// This sequence or it's reverse complement will not be present at the front of the Top strand.
	ThreePrimeBottomStrandStickyEnd string
}

// Ends returns a string description of the 5' and 3' ends of the DigestedFragment
func (fragment Digestedfragment) Ends() string {

	dnaSeq, err := fragment.ToDNASequence("fragment")

	if err != nil {
		panic(err)
	}

	return fmt.Sprintf(`5' end: %s; 3' end: %s`, dnaSeq.Overhang5prime.ToString(), dnaSeq.Overhang3prime.ToString())
}

func toDigestedFragment(seq wtype.DNASequence) (fragment Digestedfragment) {
	return Digestedfragment{
		Topstrand:                       seq.Sequence(),
		Bottomstrand:                    wtype.RevComp(seq.Sequence()),
		FivePrimeTopStrandStickyend:     seq.Overhang5prime.OverHangAt5PrimeEnd(),
		ThreePrimeTopStrandStickyend:    seq.Overhang3prime.OverHangAt5PrimeEnd(),
		FivePrimeBottomStrandStickyend:  seq.Overhang3prime.OverHangAt3PrimeEnd(),
		ThreePrimeBottomStrandStickyEnd: seq.Overhang5prime.OverHangAt3PrimeEnd(),
	}
}

// ToDNASequence assumes phosphorylation since result of digestion.
// todo:  Check and fix the construction of the digested fragment...
// This may be produced incorrectly so the error capture steps have been commented out to ensure the Insert function returns the expected result!
func (fragment Digestedfragment) ToDNASequence(name string) (seq wtype.DNASequence, err error) {

	seq = wtype.MakeLinearDNASequence(name, fragment.Topstrand)

	var overhangstr string
	var overhangtype int

	/* //
	if len(fragment.BottomStickyend_5prime) > 0 && len(fragment.TopStickyend_5prime) > 0 {
		return seq, fmt.Errorf("Cannot have 5' top %s and bottom %s strand overhangs on same sequence: ", fragment.BottomStickyend_5prime, fragment.TopStickyend_5prime)
	}
	*/

	if len(fragment.ThreePrimeTopStrandStickyend) > 0 && len(fragment.ThreePrimeBottomStrandStickyEnd) > 0 {
		return seq, fmt.Errorf("Cannot have 3' top %s and bottom %s strand overhangs on same sequence: ", fragment.ThreePrimeTopStrandStickyend, fragment.ThreePrimeBottomStrandStickyEnd)
	}

	if len(fragment.FivePrimeTopStrandStickyend) > 0 /*&& len(fragment.BottomStickyend_5prime) == 0*/ {
		overhangstr = fragment.FivePrimeTopStrandStickyend
		overhangtype = wtype.OVERHANG
	} else if len(fragment.FivePrimeTopStrandStickyend) == 0 && len(fragment.FivePrimeBottomStrandStickyend) == 0 {
		overhangstr = fragment.FivePrimeTopStrandStickyend
		overhangtype = wtype.BLUNT
	} else if len(fragment.FivePrimeBottomStrandStickyend) > 0 && len(fragment.FivePrimeTopStrandStickyend) == 0 {
		overhangstr = fragment.FivePrimeBottomStrandStickyend
		overhangtype = wtype.UNDERHANG
	} else {
		return seq, fmt.Errorf("Cannot make valid combination of overhangs with this fragment: %+v", fragment)

	}

	var overhang5 = wtype.Overhang{
		End:             5,
		Type:            overhangtype,
		Length:          len(overhangstr),
		Sequence:        overhangstr,
		Phosphorylation: true,
	}

	seq.Overhang5prime = overhang5

	if len(fragment.ThreePrimeTopStrandStickyend) > 0 && len(fragment.ThreePrimeBottomStrandStickyEnd) == 0 {
		overhangstr = fragment.ThreePrimeTopStrandStickyend
		overhangtype = wtype.OVERHANG
	} else if len(fragment.ThreePrimeTopStrandStickyend) == 0 && len(fragment.ThreePrimeBottomStrandStickyEnd) == 0 {
		overhangstr = fragment.ThreePrimeTopStrandStickyend
		overhangtype = wtype.BLUNT
	} else if len(fragment.ThreePrimeBottomStrandStickyEnd) > 0 && len(fragment.ThreePrimeTopStrandStickyend) == 0 {
		overhangstr = fragment.ThreePrimeBottomStrandStickyEnd
		overhangtype = wtype.UNDERHANG
	} else {
		return seq, fmt.Errorf("Cannot make valid combination of overhangs with this fragment: %+v", fragment)

	}

	var overhang3 = wtype.Overhang{
		End:             3,
		Type:            overhangtype,
		Length:          len(overhangstr),
		Sequence:        overhangstr,
		Phosphorylation: true,
	}

	seq.Overhang3prime = overhang3

	return
}

// Digest will simulate digestion of a DNA sequence with a chosen restriction enzyme; returns string arrays of fragments and 5' and 3' sticky ends
// todo: refactor
func Digest(sequence wtype.DNASequence, typeIIenzyme wtype.RestrictionEnzyme) (finalFragments []string, fivePrimeOverhangs []string, threePrimeOverhangs []string) {
	if typeIIenzyme.Class == "TypeII" {
		finalFragments, fivePrimeOverhangs, threePrimeOverhangs = TypeIIDigest(sequence, typeIIenzyme)
	}
	if typeIIenzyme.Class == "TypeIIs" {

		var typeIIsenz = wtype.TypeIIs{RestrictionEnzyme: typeIIenzyme}

		finalFragments, fivePrimeOverhangs, threePrimeOverhangs = TypeIIsdigest(sequence, typeIIsenz)
	}
	return
}

// RestrictionMapper returns a set of fragment sizes expected by digesting a DNA sequence with a restriction enzyme.
func RestrictionMapper(seq wtype.DNASequence, enzyme wtype.RestrictionEnzyme) (fraglengths []int) {
	enzlist := []wtype.RestrictionEnzyme{enzyme}
	frags, _, _ := Digest(seq, enzlist[0]) // doesn't handle non cutters well - returns 1 seq string, blunt, blunt therefore inaccurate representation
	fraglengths = make([]int, 0)
	for _, frag := range frags {
		fraglengths = append(fraglengths, len(frag))
	}
	fragslice := sort.IntSlice(fraglengths)
	fragslice.Sort()

	return fraglengths
}

// utility function
func searchandCutRev(typeIIenzyme wtype.RestrictionEnzyme, topstranddigestproducts []string, fivePrimeTopStrandStickyends []string, threePrimeTopStrandStickyEnds []string) (finalFragments []string, fivePrimeOverhangs []string, threePrimeOverhangs []string) {
	finalFragments = make([]string, 0)
	reverseenzymeseq := sequences.RevComp(strings.ToUpper(typeIIenzyme.RecognitionSequence))

	if reverseenzymeseq == strings.ToUpper(typeIIenzyme.RecognitionSequence) {
		finalFragments = topstranddigestproducts
		fivePrimeOverhangs = fivePrimeTopStrandStickyends
		threePrimeOverhangs = threePrimeTopStrandStickyEnds
	} else {
		originalfwdsequence := strings.Join(topstranddigestproducts, "")
		sites := search.FindAll(originalfwdsequence, reverseenzymeseq)
		// step 2. Search for recognition site on top strand, if it's there then we start processing according to the enzyme cutting properties
		if len(sites) == 0 {
			finalFragments = topstranddigestproducts
		} else {
			finaldigestproducts := make([]string, 0)
			finaltopstrandstickyends5prime := make([]string, 0)
			finaltopstrandstickyends3prime := make([]string, 0)
			for _, fragment := range topstranddigestproducts {
				cuttopstrand := strings.Split(fragment, reverseenzymeseq)
				// reversed
				recognitionsiteup := sequences.Prefix(reverseenzymeseq, (-1 * typeIIenzyme.Bottomstrand5primedistancefromend))
				recognitionsitedown := sequences.Suffix(reverseenzymeseq, (-1 * typeIIenzyme.Topstrand3primedistancefromend))
				firstfrag := strings.Join([]string{cuttopstrand[0], recognitionsiteup}, "")
				finaldigestproducts = append(finaldigestproducts, firstfrag)
				for i := 1; i < len(cuttopstrand); i++ {
					joineddownstream := strings.Join([]string{recognitionsitedown, cuttopstrand[i]}, "")
					if i != len(cuttopstrand)-1 {
						joineddownstream = strings.Join([]string{joineddownstream, recognitionsiteup}, "")
					}
					finaldigestproducts = append(finaldigestproducts, joineddownstream)
				}
				frag2topStickyend5prime := ""
				frag2topStickyend3prime := ""
				// cut with 5prime overhang
				if len(recognitionsitedown) > len(recognitionsiteup) {
					for i := 1; i < len(cuttopstrand); i++ {
						frag2topStickyend5prime = recognitionsitedown[:typeIIenzyme.EndLength]
						finaltopstrandstickyends5prime = append(finaltopstrandstickyends5prime, frag2topStickyend5prime)
						if i != len(cuttopstrand)-1 {
							frag2topStickyend3prime = ""
						} else {
							frag2topStickyend3prime = blunt
						}
						finaltopstrandstickyends3prime = append(finaltopstrandstickyends3prime, frag2topStickyend3prime)
					}
				}
				// blunt cut
				if len(recognitionsitedown) == len(recognitionsiteup) {
					for i := 1; i < len(cuttopstrand); i++ {
						frag2topStickyend5prime = blunt
						finaltopstrandstickyends5prime = append(finaltopstrandstickyends5prime, frag2topStickyend5prime)
						frag2topStickyend3prime = blunt
						finaltopstrandstickyends3prime = append(finaltopstrandstickyends3prime, frag2topStickyend3prime)
					}
				}
				// cut with 3prime overhang
				if len(recognitionsitedown) < len(recognitionsiteup) {

					for i := 1; i < len(cuttopstrand); i++ {
						frag2topStickyend5prime = ""
						finaltopstrandstickyends5prime = append(finaltopstrandstickyends5prime, frag2topStickyend5prime)
						if i != len(cuttopstrand)-1 {
							frag2topStickyend3prime = recognitionsiteup[typeIIenzyme.EndLength:]
						} else {
							frag2topStickyend3prime = blunt
						}
						finaltopstrandstickyends3prime = append(finaltopstrandstickyends3prime, frag2topStickyend3prime)
					}
				}
				for _, strand5 := range finaltopstrandstickyends5prime {
					fivePrimeTopStrandStickyends = append(fivePrimeTopStrandStickyends, strand5)
				}
				for _, strand3 := range finaltopstrandstickyends3prime {
					threePrimeTopStrandStickyEnds = append(threePrimeTopStrandStickyEnds, strand3)
				}
				finalFragments = finaldigestproducts
				fivePrimeOverhangs = fivePrimeTopStrandStickyends
				threePrimeOverhangs = threePrimeTopStrandStickyEnds
			}
		}
	}
	return
}

// utility function to correct number and order of fragments if digested sequence was a plasmid; (e.g. cutting once in plasmid dna creates one fragment; cutting once in linear dna creates 2 fragments.
func lineartoPlasmid(fragmentsiflinearstart []string) (fragmentsifplasmidstart []string) {

	// make linear plasmid part by joining last part to first part
	plasmidcutproducts := make([]string, 0)
	plasmidcutproducts = append(plasmidcutproducts, fragmentsiflinearstart[len(fragmentsiflinearstart)-1])
	plasmidcutproducts = append(plasmidcutproducts, fragmentsiflinearstart[0])
	linearpartfromplasmid := strings.Join(plasmidcutproducts, "")

	// fix order of final fragments
	fragmentsifplasmidstart = make([]string, 0)
	fragmentsifplasmidstart = append(fragmentsifplasmidstart, linearpartfromplasmid)
	for i := 1; i < (len(fragmentsiflinearstart) - 1); i++ {
		fragmentsifplasmidstart = append(fragmentsifplasmidstart, fragmentsiflinearstart[i])
	}

	return
}

// TypeIIDigest digests a DNA sequence using a restriction enzyme and returns 3 string arrays: fragments after digestion, 5prime sticky ends, 3prime sticky ends
// todo: refactor
func TypeIIDigest(sequence wtype.DNASequence, typeIIenzyme wtype.RestrictionEnzyme) (finalFragments []string, fivePrimeOverhangs []string, threePrimeOverhangs []string) {
	// step 1. get sequence in string format from DNASequence, make sure all spaces are removed and all upper case

	if typeIIenzyme.Class != "TypeII" {
		panic("This is not the function you are looking for! Wrong enzyme class for this function")
	}

	originalfwdsequence := strings.TrimSpace(strings.ToUpper(sequence.Seq))
	//originalreversesequence := strings.TrimSpace(strings.ToUpper(RevComp(sequence.Seq)))
	sites := search.FindAll(originalfwdsequence, strings.ToUpper(typeIIenzyme.RecognitionSequence))

	// step 2. Search for recognition site on top strand, if it's there then we start processing according to the enzyme cutting properties
	topstranddigestproducts := make([]string, 0)
	topStrandStickyEnds5prime := make([]string, 0)
	topStrandStickyEnds3Prime := make([]string, 0)

	if len(sites) != 0 {

		cuttopstrand := strings.Split(originalfwdsequence, strings.ToUpper(typeIIenzyme.RecognitionSequence))
		recognitionsitedown := sequences.Suffix(typeIIenzyme.RecognitionSequence, (-1 * typeIIenzyme.Topstrand3primedistancefromend))
		recognitionsiteup := sequences.Prefix(typeIIenzyme.RecognitionSequence, (-1 * typeIIenzyme.Bottomstrand5primedistancefromend))

		//repairedfrag := ""
		//repairedfrags := make([]string,0)

		//if sequence.Plasmid != true{

		firstfrag := strings.Join([]string{cuttopstrand[0], recognitionsiteup}, "")
		topstranddigestproducts = append(topstranddigestproducts, firstfrag)

		for i := 1; i < len(cuttopstrand); i++ {
			joineddownstream := strings.Join([]string{recognitionsitedown, cuttopstrand[i]}, "")
			if i != len(cuttopstrand)-1 {
				joineddownstream = strings.Join([]string{joineddownstream, recognitionsiteup}, "")
			}
			topstranddigestproducts = append(topstranddigestproducts, joineddownstream)

		}

		frag2topStickyend5prime := ""
		frag2topStickyend3prime := ""
		// cut with 5prime overhang
		if len(recognitionsitedown) > len(recognitionsiteup) {
			frag2topStickyend5prime = blunt
			topStrandStickyEnds5prime = append(topStrandStickyEnds5prime, frag2topStickyend5prime)
			frag2topStickyend3prime := ""
			topStrandStickyEnds3Prime = append(topStrandStickyEnds3Prime, frag2topStickyend3prime)
			for i := 1; i < len(cuttopstrand); i++ {
				frag2topStickyend5prime = recognitionsitedown[:typeIIenzyme.EndLength]
				topStrandStickyEnds5prime = append(topStrandStickyEnds5prime, frag2topStickyend5prime)
				if i != len(cuttopstrand)-1 {
					frag2topStickyend3prime = ""
				} else {
					frag2topStickyend3prime = blunt
				}
				topStrandStickyEnds3Prime = append(topStrandStickyEnds3Prime, frag2topStickyend3prime)

			}

		}
		// blunt cut
		if len(recognitionsitedown) == len(recognitionsiteup) {
			for i := 0; i < len(cuttopstrand); i++ {
				frag2topStickyend5prime = blunt
				topStrandStickyEnds5prime = append(topStrandStickyEnds5prime, frag2topStickyend5prime)
				frag2topStickyend3prime = blunt
				topStrandStickyEnds3Prime = append(topStrandStickyEnds3Prime, frag2topStickyend3prime)
			}
		}
		// cut with 3prime overhang
		if len(recognitionsitedown) < len(recognitionsiteup) {
			frag2topStickyend5prime = blunt
			topStrandStickyEnds5prime = append(topStrandStickyEnds5prime, frag2topStickyend5prime)

			frag2topStickyend3prime = sequences.Suffix(recognitionsiteup, typeIIenzyme.EndLength)
			topStrandStickyEnds3Prime = append(topStrandStickyEnds3Prime, frag2topStickyend3prime)

			for i := 1; i < len(cuttopstrand); i++ {
				frag2topStickyend5prime = ""
				topStrandStickyEnds5prime = append(topStrandStickyEnds5prime, frag2topStickyend5prime)
				if i != len(cuttopstrand)-1 {
					frag2topStickyend3prime = recognitionsiteup[typeIIenzyme.EndLength:]
				} else {
					frag2topStickyend3prime = blunt
				}
				topStrandStickyEnds3Prime = append(topStrandStickyEnds3Prime, frag2topStickyend3prime)

			}
		}
	} else {
		topstranddigestproducts = []string{originalfwdsequence}
		topStrandStickyEnds5prime = []string{blunt}
		topStrandStickyEnds3Prime = []string{blunt}
	}

	finalFragments, topStrandStickyEnds5prime, topStrandStickyEnds3Prime = searchandCutRev(typeIIenzyme, topstranddigestproducts, topStrandStickyEnds5prime, topStrandStickyEnds3Prime)

	if len(finalFragments) == 1 && sequence.Plasmid == true {
		// TODO
		// need to really return an uncut plasmid, maybe an error?
		//	// fmt.Println("uncut plasmid returned with no sticky ends!")

	}
	if len(finalFragments) > 1 && sequence.Plasmid == true {
		ifplasmidfinalfragments := lineartoPlasmid(finalFragments)
		finalFragments = ifplasmidfinalfragments
		// now change order of sticky ends
		//5'
		ifplasmidsticky5prime := make([]string, 0)
		ifplasmidsticky5prime = append(ifplasmidsticky5prime, topStrandStickyEnds5prime[len(topStrandStickyEnds5prime)-1])
		for i := 1; i < (len(finalFragments)); i++ {
			ifplasmidsticky5prime = append(ifplasmidsticky5prime, topStrandStickyEnds5prime[i])
		}
		topStrandStickyEnds5prime = ifplasmidsticky5prime
		//hack to fix wrong sticky end assignment in certain cases
		reverseenzymeseq := sequences.RevComp(typeIIenzyme.RecognitionSequence)
		if strings.Index(originalfwdsequence, strings.ToUpper(typeIIenzyme.RecognitionSequence)) > strings.Index(originalfwdsequence, reverseenzymeseq) {
			topStrandStickyEnds5prime = sequences.RevArrayOrder(topStrandStickyEnds5prime)
		}
		//3'
		ifplasmidsticky3prime := make([]string, 0)
		ifplasmidsticky3prime = append(ifplasmidsticky3prime, topStrandStickyEnds3Prime[0])
		for i := 1; i < (len(finalFragments)); i++ {
			ifplasmidsticky3prime = append(ifplasmidsticky3prime, topStrandStickyEnds3Prime[i])
		}
		topStrandStickyEnds3Prime = ifplasmidsticky3prime
	}
	fivePrimeOverhangs = topStrandStickyEnds5prime
	// deal with this later
	threePrimeOverhangs = topStrandStickyEnds3Prime
	return finalFragments, fivePrimeOverhangs, threePrimeOverhangs
}

// TypeIIsdigest returns slices of fragments, 5 prime overhangs and 3 prime underhangs generated from cutting with a typeIIs enzyme which leaves a 5 prime overhang.
func TypeIIsdigest(sequence wtype.DNASequence, typeIIsenzyme wtype.TypeIIs) (finalFragments []string, fivePrimeOverhangs []string, threePrimeUnderhangs []string) {

	restrictionSites := sequences.FindAll(&sequence, &wtype.DNASequence{Nm: typeIIsenzyme.Name(), Seq: typeIIsenzyme.RecognitionSequence})

	seqs, err := makeFragments(typeIIsenzyme, restrictionSites.Positions, sequence)

	if err != nil {
		fmt.Println(err.Error())
	}

	for _, seq := range seqs {
		finalFragments = append(finalFragments, seq.Sequence())
		fivePrimeOverhangs = append(fivePrimeOverhangs, seq.Overhang5prime.OverHangAt5PrimeEnd())
		threePrimeUnderhangs = append(threePrimeUnderhangs, seq.Overhang3prime.UnderHangAt3PrimeEnd())
	}

	return
}

// TypeIIsDigestToFragments returns slices of fragments generated from cutting with a typeIIs enzyme which leaves a 5 prime overhang.
func TypeIIsDigestToFragments(sequence wtype.DNASequence, typeIIsenzyme wtype.TypeIIs) (finalFragments []Digestedfragment) {

	restrictionSites := sequences.FindAll(&sequence, &wtype.DNASequence{Nm: typeIIsenzyme.Name(), Seq: typeIIsenzyme.RecognitionSequence})

	seqs, err := makeFragments(typeIIsenzyme, restrictionSites.Positions, sequence)

	if err != nil {
		fmt.Println(err.Error())
	}

	for _, seq := range seqs {
		finalFragments = append(finalFragments, toDigestedFragment(seq))
	}
	return
}

// calculates teh specific cut position based on an enzyme type and specified position pair.
// NB: since the sequence is not specified and typeIIs enzymes cut remotely to the recognition site,
//  it's possible to calculate a cut position beyond the range of the sequence.
// The logic to prevent this is in the companion functions SequenceBetweenPositions and MakeOverhangs.
func correctTypeIIsCutPosition(enzyme wtype.TypeIIs, recognitionSitePosition sequences.PositionPair) (fragmentStart int) {

	_, endOfRestrictionSite := recognitionSitePosition.Coordinates(wtype.CODEFRIENDLY)

	if !recognitionSitePosition.Reverse {
		fragmentStart = endOfRestrictionSite + 1 + enzyme.Topstrand3primedistancefromend
	} else {
		fragmentStart = endOfRestrictionSite - enzyme.Topstrand3primedistancefromend - enzyme.EndLength //+ 1
	}

	return
}

var nulPosition sequences.PositionPair

// makeFragment generates a fragment from a sequence cut at two specified positions with a specified enzyme.
// if either position is nul then the other end will be blunt
// if both positions are the same the sequence will be cut once at that position and if the sequence is a plasmid a single fragement is returned joining the end sequence to the beginning.
func makeFragment(enzyme wtype.TypeIIs, upstreamCutPosition, downstreamCutPosition sequences.PositionPair, originalSequence wtype.DNASequence) (fragment wtype.DNASequence, err error) {

	if upstreamCutPosition == nulPosition {
		fragment.Append(originalSequence.Sequence()[:correctTypeIIsCutPosition(enzyme, downstreamCutPosition)])

		fragment.Overhang3prime, _, err = makeOverhangs(enzyme, downstreamCutPosition, originalSequence)

		if err != nil {
			return fragment, err
		}

		fragment.Overhang5prime, err = wtype.MakeOverHang(fragment, 5, wtype.NEITHER, 0, true)
		if err != nil {
			return fragment, err
		}

		return
	} else if downstreamCutPosition == nulPosition {
		fragment.Append(originalSequence.Sequence()[correctTypeIIsCutPosition(enzyme, upstreamCutPosition):])

		_, fragment.Overhang5prime, err = makeOverhangs(enzyme, upstreamCutPosition, originalSequence)

		if err != nil {
			return fragment, err
		}

		fragment.Overhang3prime, err = wtype.MakeOverHang(fragment, 3, wtype.NEITHER, 0, true)
		if err != nil {
			return fragment, err
		}
		return
	}

	fragment.Seq, err = seqBetweenPositions(originalSequence, correctTypeIIsCutPosition(enzyme, upstreamCutPosition), correctTypeIIsCutPosition(enzyme, downstreamCutPosition))
	if err != nil {
		return fragment, err
	}
	fragment.Overhang3prime, _, err = makeOverhangs(enzyme, downstreamCutPosition, originalSequence)

	if err != nil {
		return fragment, err
	}
	_, fragment.Overhang5prime, err = makeOverhangs(enzyme, upstreamCutPosition, originalSequence)

	if err != nil {
		return fragment, err
	}
	return

}

// makeOverhangs makes the overhangs for a fragment cut at a specified position with a specified sequence.
func makeOverhangs(enzyme wtype.TypeIIs, recognitionSitePosition sequences.PositionPair, sequence wtype.DNASequence) (upStreamThreePrime, downstreamFivePrime wtype.Overhang, err error) {

	junkSequence := sequence.Dup()

	junkSequence.Plasmid = false

	_, endOfRestrictionSite := recognitionSitePosition.Coordinates(wtype.CODEFRIENDLY)

	var fragmentStart, fragmentEnd int

	if recognitionSitePosition.Reverse {

		fragmentStart = endOfRestrictionSite - enzyme.Topstrand3primedistancefromend - enzyme.EndLength //+ 1

	} else {

		fragmentStart = endOfRestrictionSite + 1 + enzyme.Topstrand3primedistancefromend

	}

	if fragmentStart > len(sequence.Sequence()) {
		fragmentStart = fragmentStart - len(sequence.Sequence())
	}

	fragmentEnd = fragmentStart + enzyme.EndLength

	if fragmentEnd > len(sequence.Sequence()) {
		fragmentEnd = fragmentEnd - len(sequence.Sequence())
	}

	overhangSeq, err := seqBetweenPositions(sequence, fragmentStart, fragmentEnd)
	if err != nil {
		return upStreamThreePrime, downstreamFivePrime, err
	}
	downstreamFivePrime, err = wtype.MakeOverHang(junkSequence, 5, wtype.TOP, enzyme.EndLength, true)
	if err != nil {
		return upStreamThreePrime, downstreamFivePrime, err
	}
	upStreamThreePrime, err = wtype.MakeOverHang(junkSequence, 3, wtype.BOTTOM, enzyme.EndLength, true)
	if err != nil {
		return upStreamThreePrime, downstreamFivePrime, err
	}
	downstreamFivePrime.Sequence = overhangSeq
	upStreamThreePrime.Sequence = wtype.RevComp(overhangSeq)

	return upStreamThreePrime, downstreamFivePrime, nil
}

// code friendly positions and forward orientation only
func seqBetweenPositions(sequence wtype.DNASequence, start, end int) (string, error) {

	if end > len(sequence.Sequence()) {
		// if sequence is a plasmid, wrap around automatically
		if sequence.Plasmid {
			end = end - len(sequence.Sequence())
		} else {
			return sequence.Sequence(), fmt.Errorf("sequence end position out of range: start %d; end %d, sequence %s", start, end, sequence.Sequence())
		}
	}

	if start > len(sequence.Sequence()) {
		// if sequence is a plasmid, wrap around automatically
		if sequence.Plasmid {
			start = start - len(sequence.Sequence())
		} else {
			return sequence.Sequence(), fmt.Errorf("sequence start position out of range: start %d; end %d, sequence %s", start, end, sequence.Sequence())
		}
	}

	if start >= end {

		return sequence.Sequence()[start:] + sequence.Sequence()[:end], nil
	}

	if start < 0 {
		//panic("no")
		if sequence.Plasmid {
			return sequence.Sequence()[len(sequence.Sequence())+start:] + sequence.Sequence()[:end], nil
		}
		// if this happens it's probably a non-cutter
		return sequence.Sequence()[:end], fmt.Errorf("sequence start position %d negative and not plasmid for sequence %s", start, sequence.Name())
	}

	return sequence.Sequence()[start:end], nil
}

// makeFragments will correct the position assignment of a search on a plasmid sequence.
// A search on a plasmid sequence may need to find matches which overlap the end of a plasmid.
// findSeq will therefore first concatenates the plasmid sequence with a duplicate and then perform a search.
// correctPositions will correct the position assignment of any matches which are found that overlap the end of a plasmid sequence.
func makeFragments(enzyme wtype.TypeIIs, positionPairs []sequences.PositionPair, originalSequence wtype.DNASequence) (fragments []wtype.DNASequence, err error) {

	sortedPairs := sequences.ByPositionPairStartPosition(positionPairs)

	sort.Sort(sortedPairs)

	if len(sortedPairs) == 0 {
		return []wtype.DNASequence{originalSequence}, fmt.Errorf("no positions specified to cut sequence for %s", originalSequence.Name())
	}

	last := len(sortedPairs) - 1

	for i := range sortedPairs {
		// first fragment
		var fragment wtype.DNASequence

		var upstreamCutPosition, downstreamCutPosition sequences.PositionPair

		if i == 0 {

			// if plasmid add last part of sequence after final cut position to fragment
			if originalSequence.Plasmid {

				upstreamCutPosition = sortedPairs[last]

			} else {

				upstreamCutPosition = nulPosition

			}

		} else {

			upstreamCutPosition = sortedPairs[i-1]

		}

		downstreamCutPosition = sortedPairs[i]

		fragment, err = makeFragment(enzyme, upstreamCutPosition, downstreamCutPosition, originalSequence)
		if err != nil {
			return fragments, err
		}
		fragments = append(fragments, fragment)

		// an extra fragment needs to be added if in last position and plasmid
		if i == last && !originalSequence.Plasmid {

			nextUpStreamCutPosition := downstreamCutPosition

			fragment, err = makeFragment(enzyme, nextUpStreamCutPosition, nulPosition, originalSequence)
			if err != nil {
				return fragments, err
			}
			fragments = append(fragments, fragment)
		}

	}

	// if plasmid sequence move first fragment to back
	if originalSequence.Plasmid && len(fragments) > 2 {
		var newFragments []wtype.DNASequence
		newFragments = append(newFragments, fragments[1:]...)
		newFragments = append(newFragments, fragments[0])

		for i := range newFragments {
			newFragments[i].Nm = "fragment" + strconv.Itoa(i+1)
		}

		if len(newFragments) != len(fragments) {
			panic("Ahhhh")
		}

		return newFragments, nil
	}

	for i := range fragments {
		fragments[i].Nm = "fragment" + strconv.Itoa(i+1)
	}

	return
}

// EndReport returns a report of all ends expected from digesting a vector sequence and a set of parts as a string.
// Intended to aid the user in trouble shooting unsuccessful assemblies.
func EndReport(restrictionenzyme wtype.TypeIIs, vectordata wtype.DNASequence, parts []wtype.DNASequence) (endreport string) {
	_, stickyends5, stickyends3 := TypeIIsdigest(vectordata, restrictionenzyme)

	allends := make([]string, 0)
	ends := ""

	ends = text.Print(vectordata.Nm+" 5 Prime end: ", stickyends5)
	allends = append(allends, ends)
	ends = text.Print(vectordata.Nm+" 3 Prime end: ", stickyends3)
	allends = append(allends, ends)

	for _, part := range parts {
		_, stickyends5, stickyends3 = TypeIIsdigest(part, restrictionenzyme)
		ends = text.Print(part.Nm+" 5 Prime end: ", stickyends5)
		allends = append(allends, ends)
		ends = text.Print(part.Nm+" 3 Prime end: ", stickyends3)
		allends = append(allends, ends)
	}
	endreport = strings.Join(allends, " ")
	return
}
