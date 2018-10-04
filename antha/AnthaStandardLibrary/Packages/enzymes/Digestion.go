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

	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/sequences"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
)

const (
	typeIIs string = "TypeIIs"
)

// RestrictionSites holds information on restriction sites found in a DNA sequence
// for a specified RestrictionEnzyme.
type RestrictionSites struct {
	Enzyme    wtype.RestrictionEnzyme
	Positions []sequences.PositionPair
}

// RecognitionSequence returns the recognition sequence of the enzyme as a string.
func (sites RestrictionSites) RecognitionSequence() string {
	return sites.Enzyme.RecognitionSequence
}

// SiteFound evaluates whether at least one site has been found.
func (sites RestrictionSites) SiteFound() bool {
	return len(sites.Positions) > 0
}

// NumberOfSites returns the number of restriction sites found.
func (sites RestrictionSites) NumberOfSites() int {
	return len(sites.Positions)
}

// ForwardPositions returns the recognition site positions
func (sites RestrictionSites) ForwardPositions() []int {
	var forwardPositions []int
	for _, position := range sites.Positions {
		if !position.Reverse {
			forwardPositions = append(forwardPositions, position.StartPosition)
		}
	}
	return forwardPositions
}

// ReversePositions returns the recognition site positions
func (sites RestrictionSites) ReversePositions() []int {
	var reversePositions []int
	for _, position := range sites.Positions {
		if position.Reverse {
			start, _ := position.Coordinates(wtype.IGNOREDIRECTION)
			reversePositions = append(reversePositions, start)
		}
	}
	return reversePositions
}

// AllPositions returns all forward and reverse restriction site positions.
func (sites *RestrictionSites) AllPositions() []int {
	var allPositions []int
	for _, position := range sites.Positions {
		start, _ := position.Coordinates(wtype.IGNOREDIRECTION)
		allPositions = append(allPositions, start)
	}
	return allPositions
}

// SitepositionString returns a report of restriction sites found as a string
// todo: deprecate
func SitepositionString(sitesperpart RestrictionSites) (sitepositions string) {
	Num := make([]string, 0)

	for _, site := range sitesperpart.ForwardPositions() {
		Num = append(Num, strconv.Itoa(site))
	}
	for _, site := range sitesperpart.ReversePositions() {
		Num = append(Num, strconv.Itoa(site))
	}

	sort.Strings(Num)
	sitepositions = strings.Join(Num, ", ")
	return
}

func isEqualSite(position1, position2 sequences.PositionPair) bool {
	if position1.StartPosition == position2.StartPosition && position1.EndPosition == position2.EndPosition && position1.Reverse == position2.Reverse {
		return true
	}
	if position1.StartPosition == position2.EndPosition && position1.EndPosition == position2.StartPosition && position1.Reverse != position2.Reverse {
		return true
	}
	return false
}

func inPositions(positions []sequences.PositionPair, target sequences.PositionPair) bool {
	for _, position := range positions {
		if isEqualSite(position, target) {
			return true
		}
	}
	return false
}

func removePalindromic(positions []sequences.PositionPair) []sequences.PositionPair {
	var nonPalindromic []sequences.PositionPair
	for _, position := range positions {
		if !inPositions(nonPalindromic, position) {
			nonPalindromic = append(nonPalindromic, position)
		}
	}
	return nonPalindromic
}

// RestrictionSiteFinder finds restriction sites of specified restriction enzymes in a sequence and return the information as a set of ResrictionSites.
func RestrictionSiteFinder(sequence wtype.DNASequence, enzymelist ...wtype.RestrictionEnzyme) (sites []RestrictionSites) {

	sites = make([]RestrictionSites, 0)

	for _, enzyme := range enzymelist {
		var enzymesite RestrictionSites
		enzymesite.Enzyme = enzyme
		recognitionSite := strings.ToUpper(enzyme.RecognitionSequence)
		sequence.Seq = strings.ToUpper(sequence.Seq)

		wobbleproofrecognitionoptions := sequences.Wobble(recognitionSite)

		for _, wobbleoption := range wobbleproofrecognitionoptions {
			options := sequences.FindAll(&sequence, &wtype.DNASequence{Nm: wobbleoption, Seq: wobbleoption})
			enzymesite.Positions = append(enzymesite.Positions, options.Positions...)
		}

		enzymesite.Positions = removePalindromic(enzymesite.Positions)

		sites = append(sites, enzymesite)
	}

	return sites
}

// DigestedFragment object carrying info on a fragment following digestion
type DigestedFragment struct {

	// Sequence of top strand
	TopStrand string

	// Sequence of bottom strand
	BottomStrand string

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
func (fragment DigestedFragment) Ends() string {

	dnaSeq, err := fragment.ToDNASequence("fragment")

	if err != nil {
		panic(err)
	}

	return fmt.Sprintf(`5' end: %s; 3' end: %s`, dnaSeq.Overhang5prime.ToString(), dnaSeq.Overhang3prime.ToString())
}

func toDigestedFragment(seq wtype.DNASequence) (fragment DigestedFragment) {
	return DigestedFragment{
		TopStrand:                       seq.Sequence(),
		BottomStrand:                    wtype.RevComp(seq.Sequence()),
		FivePrimeTopStrandStickyend:     seq.Overhang5prime.OverHang(),
		ThreePrimeTopStrandStickyend:    seq.Overhang3prime.OverHang(),
		FivePrimeBottomStrandStickyend:  seq.Overhang3prime.UnderHang(),
		ThreePrimeBottomStrandStickyEnd: seq.Overhang5prime.UnderHang(),
	}
}

// ToDNASequence assumes phosphorylation since result of digestion.
// todo:  Check and fix the construction of the digested fragment...
// This may be produced incorrectly so the error capture steps have been commented out to ensure the Insert function returns the expected result!
func (fragment DigestedFragment) ToDNASequence(name string) (seq wtype.DNASequence, err error) {

	seq = wtype.MakeLinearDNASequence(name, fragment.TopStrand)

	var overhangstr string
	var overhangtype wtype.OverHangType

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
		Seq:             overhangstr,
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
		Seq:             overhangstr,
		Phosphorylation: true,
	}

	seq.Overhang3prime = overhang3

	return
}

// DigestToFragments will simulate digestion of a DNA sequence with one or more restriction enzymes;
// returns the products of the digestion in the form of a set of DigestedFragment.
func DigestToFragments(sequence wtype.DNASequence, typeIIenzymes ...wtype.RestrictionEnzyme) (finalFragments []DigestedFragment, err error) {

	seqs, err := Digest(sequence, typeIIenzymes...)

	for _, seq := range seqs {
		finalFragments = append(finalFragments, toDigestedFragment(seq))
	}

	return
}

// RestrictionMapper returns a set of fragment sizes expected by digesting a DNA sequence with a series of restriction enzymes.
func RestrictionMapper(seq wtype.DNASequence, enzymes ...wtype.RestrictionEnzyme) (fraglengths []int) {
	frags, err := DigestToFragments(seq, enzymes...) // doesn't handle non cutters well - returns 1 seq string, blunt, blunt therefore inaccurate representation
	if err != nil {
		panic(err.Error())
	}
	fraglengths = make([]int, 0)
	for _, frag := range frags {
		fraglengths = append(fraglengths, len(frag.TopStrand))
	}
	fragslice := sort.IntSlice(fraglengths)
	fragslice.Sort()

	return fraglengths
}

// TypeIIsdigest returns slices of fragments, 5 prime overhangs and 3 prime underhangs generated from cutting with a typeIIs enzyme which leaves a 5 prime overhang.
func TypeIIsdigest(sequence wtype.DNASequence, typeIIsenzyme wtype.TypeIIs) (finalFragments []string, fivePrimeOverhangs []string, threePrimeUnderhangs []string) {

	restrictionSites := RestrictionSiteFinder(sequence, typeIIsenzyme.RestrictionEnzyme)

	seqs, err := makeFragments(typeIIsenzyme.RestrictionEnzyme, restrictionSites[0].Positions, sequence)

	if err != nil {
		fmt.Println(err.Error())
	}

	for _, seq := range seqs {
		finalFragments = append(finalFragments, seq.Sequence())
		fivePrimeOverhangs = append(fivePrimeOverhangs, seq.Overhang5prime.OverHang())
		threePrimeUnderhangs = append(threePrimeUnderhangs, seq.Overhang3prime.UnderHang())
	}

	return
}

// TypeIIsDigestToFragments returns slices of fragments generated from cutting with a typeIIs enzyme which leaves a 5 prime overhang.
func typeIIsDigestToFragments(sequence wtype.DNASequence, typeIIsenzymes ...wtype.TypeIIs) (finalFragments []DigestedFragment, err error) {

	var enzymes []wtype.RestrictionEnzyme

	for _, typeIIs := range typeIIsenzymes {
		enzymes = append(enzymes, typeIIs.RestrictionEnzyme)
	}

	seqs, err := Digest(sequence, enzymes...)

	for _, seq := range seqs {
		finalFragments = append(finalFragments, toDigestedFragment(seq))
	}
	return
}

// calculates teh specific cut position based on an enzyme type and specified position pair.
// NB: since the sequence is not specified and typeIIs enzymes cut remotely to the recognition site,
//  it's possible to calculate a cut position beyond the range of the sequence.
// The logic to prevent this is in the companion functions SequenceBetweenPositions and MakeOverhangs.
func correctTypeIIsCutPosition(enzyme wtype.RestrictionEnzyme, recognitionSitePosition sequences.PositionPair) (fragmentStart int) {

	startOfRestrictionSite, endOfRestrictionSite := recognitionSitePosition.Coordinates(wtype.CODEFRIENDLY)

	switch class := enzyme.Class; class {

	case typeIIs:

		if !recognitionSitePosition.Reverse {
			fragmentStart = endOfRestrictionSite + 1 + enzyme.Topstrand3primedistancefromend
		} else {
			fragmentStart = endOfRestrictionSite - enzyme.Bottomstrand5primedistancefromend // - enzyme.EndLength //+ 1
		}

	default:

		// any TypeII enzymes coming in should only be coming in forwards
		if !recognitionSitePosition.Reverse {
			fragmentStart = startOfRestrictionSite - enzyme.Bottomstrand5primedistancefromend
		} else {
			// not verified
			fragmentStart = startOfRestrictionSite - enzyme.Bottomstrand5primedistancefromend
		}

	}
	return
}

var nulPosition sequences.PositionPair

// makeFragment generates a fragment from a sequence cut at two specified positions with a specified enzyme.
// if either position is nul then the other end will be blunt
// if both positions are the same the sequence will be cut once at that position and if the sequence is a plasmid a single fragement is returned joining the end sequence to the beginning.
func makeFragment(enzyme wtype.RestrictionEnzyme, upstreamCutPosition, downstreamCutPosition sequences.PositionPair, originalSequence wtype.DNASequence) (fragment wtype.DNASequence, err error) {

	if upstreamCutPosition == nulPosition {
		err = fragment.Append(originalSequence.Sequence()[:correctTypeIIsCutPosition(enzyme, downstreamCutPosition)])

		if err != nil {
			return fragment, err
		}

		threePrimeEnd, _, err := makeOverhangs(enzyme, downstreamCutPosition, originalSequence)

		if err != nil {
			return fragment, err
		}

		err = fragment.Set3PrimeEnd(threePrimeEnd)

		if err != nil {
			return fragment, err
		}

		fragment.Overhang5prime = originalSequence.Overhang5prime

		return fragment, nil
	} else if downstreamCutPosition == nulPosition {
		err = fragment.Append(originalSequence.Sequence()[correctTypeIIsCutPosition(enzyme, upstreamCutPosition):])

		if err != nil {
			return fragment, err
		}

		_, fivePrimeEnd, err := makeOverhangs(enzyme, upstreamCutPosition, originalSequence)

		if err != nil {
			return fragment, err
		}

		err = fragment.Set5PrimeEnd(fivePrimeEnd)

		if err != nil {
			return fragment, err
		}

		fragment.Overhang3prime = originalSequence.Overhang3prime

		return fragment, nil
	}

	fragment.Seq, err = seqBetweenPositions(originalSequence, correctTypeIIsCutPosition(enzyme, upstreamCutPosition), correctTypeIIsCutPosition(enzyme, downstreamCutPosition))
	if err != nil {
		return fragment, err
	}
	threePrimeEnd, _, err := makeOverhangs(enzyme, downstreamCutPosition, originalSequence)

	if err != nil {
		return fragment, err
	}

	err = fragment.Set3PrimeEnd(threePrimeEnd)

	if err != nil {
		return fragment, err
	}

	_, fivePrimeEnd, err := makeOverhangs(enzyme, upstreamCutPosition, originalSequence)

	if err != nil {
		return fragment, err
	}

	err = fragment.Set5PrimeEnd(fivePrimeEnd)

	if err != nil {
		return fragment, err
	}

	return

}

func bottomCut(enzyme wtype.RestrictionEnzyme, recognitionSitePosition sequences.PositionPair) int {

	topStrandCut := correctTypeIIsCutPosition(enzyme, recognitionSitePosition)

	switch class := enzyme.Class; class {

	case typeIIs:

		return topStrandCut + enzyme.EndLength

	default:
		return topStrandCut + (enzyme.Bottomstrand5primedistancefromend - enzyme.Topstrand3primedistancefromend)

	}
}

// makeOverhangs makes the overhangs for a fragment cut at a specified position with a specified sequence.
func makeOverhangs(enzyme wtype.RestrictionEnzyme, recognitionSitePosition sequences.PositionPair, sequence wtype.DNASequence) (upStreamThreePrime, downstreamFivePrime wtype.Overhang, err error) {

	var fragmentStart, fragmentEnd int

	fragmentStart = correctTypeIIsCutPosition(enzyme, recognitionSitePosition)

	switch class := enzyme.Class; class {

	case typeIIs:

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
		downstreamFivePrime, err = wtype.MakeOverHang(overhangSeq, 5, wtype.TOP, true)
		if err != nil {
			return upStreamThreePrime, downstreamFivePrime, err
		}
		upStreamThreePrime, err = wtype.MakeOverHang(wtype.RevComp(overhangSeq), 3, wtype.BOTTOM, true)
		if err != nil {
			return upStreamThreePrime, downstreamFivePrime, err
		}

	default:

		fragmentEnd = bottomCut(enzyme, recognitionSitePosition)

		if fragmentEnd == fragmentStart {

			downstreamFivePrime, err = wtype.MakeOverHang("", 5, wtype.NEITHER, true)
			if err != nil {
				return upStreamThreePrime, downstreamFivePrime, err
			}
			upStreamThreePrime, err = wtype.MakeOverHang("", 3, wtype.NEITHER, true)
			if err != nil {
				return upStreamThreePrime, downstreamFivePrime, err
			}

		} else if fragmentEnd > fragmentStart {

			overhangSeq, err := seqBetweenPositions(sequence, fragmentStart, fragmentEnd)
			if err != nil {
				return upStreamThreePrime, downstreamFivePrime, err
			}
			downstreamFivePrime, err = wtype.MakeOverHang(overhangSeq, 5, wtype.TOP, true)
			if err != nil {
				return upStreamThreePrime, downstreamFivePrime, err
			}
			upStreamThreePrime, err = wtype.MakeOverHang(wtype.RevComp(overhangSeq), 3, wtype.BOTTOM, true)
			if err != nil {
				return upStreamThreePrime, downstreamFivePrime, err
			}

		} else if fragmentEnd < fragmentStart {
			overhangSeq, err := seqBetweenPositions(sequence, fragmentEnd, fragmentStart)
			if err != nil {
				return upStreamThreePrime, downstreamFivePrime, err
			}
			downstreamFivePrime, err = wtype.MakeOverHang(wtype.RevComp(overhangSeq), 5, wtype.BOTTOM, true)
			if err != nil {
				return upStreamThreePrime, downstreamFivePrime, err
			}
			upStreamThreePrime, err = wtype.MakeOverHang(overhangSeq, 3, wtype.TOP, true)
			if err != nil {
				return upStreamThreePrime, downstreamFivePrime, err
			}
		}

	}

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
func makeFragments(enzyme wtype.RestrictionEnzyme, positionPairs []sequences.PositionPair, originalSequence wtype.DNASequence) (fragments []wtype.DNASequence, err error) {

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

		fragments = append(fragments, fragment)
		if err != nil {
			return fragments, err
		}

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
			panic("something's gone wrong: new fragments are different length to original")
		}

		return newFragments, nil
	}

	for i := range fragments {
		fragments[i].Nm = "fragment" + strconv.Itoa(i+1)
	}

	return
}

// Digest will simulate digestion of a DNA sequence with one or more restriction enzymes;
// returns the products of the digestion in the form of a set of DNASequence.
func Digest(originalSequence wtype.DNASequence, enzymes ...wtype.RestrictionEnzyme) (fragments []wtype.DNASequence, err error) {

	if len(enzymes) == 0 {
		return []wtype.DNASequence{}, fmt.Errorf("No enzymes specified to make fragments")
	}

	restrictionSites := RestrictionSiteFinder(originalSequence, enzymes...)

	var someEnzymeSitesFound bool
	// return original sequence if no positions found
	for _, enzySitesFound := range restrictionSites {
		if len(enzySitesFound.Positions) > 0 {
			someEnzymeSitesFound = true
			break
		}
	}

	if !someEnzymeSitesFound {
		var names []string
		for _, enz := range enzymes {
			names = append(names, enz.Name())
		}
		return []wtype.DNASequence{originalSequence}, fmt.Errorf("no enzyme positions for %v found in sequence %s", names, originalSequence.Name())
	}

	var errs []string

	restrictionSitesForFirstEnzyme := restrictionSites[0]

	fragments, err = makeFragments(enzymes[0], restrictionSitesForFirstEnzyme.Positions, originalSequence)
	if err != nil {
		errs = append(errs, err.Error())
	}

	if len(fragments) == 0 {
		fragments = []wtype.DNASequence{originalSequence}
	}
	for i := 1; i < len(enzymes); i++ {

		latestFragments := fragments
		fragments = []wtype.DNASequence{}

		for _, fragment := range latestFragments {
			restrictionSites := RestrictionSiteFinder(fragment, enzymes[i])
			newFragments, err := makeFragments(enzymes[i], restrictionSites[0].Positions, fragment)
			if err != nil {
				errs = append(errs, err.Error())
			}
			fragments = append(fragments, newFragments...)
		}
	}

	if len(fragments) == 0 && len(errs) > 0 {
		return []wtype.DNASequence{originalSequence}, fmt.Errorf("digestion errors: %s", strings.Join(errs, "\n"))
	}
	return fragments, nil
}

func fragmentEnds(fragments []DigestedFragment) string {
	var summaries []string
	for i, fragment := range fragments {
		summaries = append(summaries, fmt.Sprintf("fragment %d: %s", i, fragment.Ends()))
	}
	return strings.Join(summaries, "\n")
}

// EndReport returns a report of all ends expected from digesting a vector sequence and a set of parts as a string.
// Intended to aid the user in trouble shooting unsuccessful assemblies.
func EndReport(restrictionenzyme wtype.TypeIIs, vector wtype.DNASequence, parts []wtype.DNASequence) (endreport string) {

	allends := make([]string, 0)

	vectorFragments, err := DigestToFragments(vector, restrictionenzyme.RestrictionEnzyme)

	if err != nil {
		panic(err.Error())
	}

	allends = append(allends, vector.Name()+" cut with "+restrictionenzyme.Name()+";  Fragment Ends :", fragmentEnds(vectorFragments))

	for _, part := range parts {
		partFragments, err := DigestToFragments(part, restrictionenzyme.RestrictionEnzyme)

		if err != nil {
			panic(err.Error())
		}

		allends = append(allends, part.Name()+" cut with "+restrictionenzyme.Name()+"; Fragment  Ends :", fragmentEnds(partFragments))
	}
	endreport = strings.Join(allends, "\n")
	return
}
