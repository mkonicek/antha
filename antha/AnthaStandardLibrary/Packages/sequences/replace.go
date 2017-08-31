// antha/AnthaStandardLibrary/Packages/enzymes/Translation.go: Part of the Antha language
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

// Package for interacting with and manipulating dna sequences in extension to methods available in wtype
package sequences

import (
	"fmt"
	"strings"

	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/search"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/text"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wutil"
)

/*
func Siteinorfs(Features features, site string) bool {

}
*/

// FindSeq searches for a DNA sequence within a larger DNA sequence
func FindSeq(bigsequence, smallsequence wtype.DNASequence) (seqsfound []search.Thingfound) {

	bigseq := strings.ToUpper(bigsequence.Sequence())

	var seqfound search.Thingfound

	seq := strings.ToUpper(smallsequence.Sequence())
	if strings.Contains(bigseq, seq) {
		seqfound.Thing = seq
		seqfound.Positions = search.Findall(bigseq, seq)
		seqsfound = append(seqsfound, seqfound)
	}

	revseq := strings.ToUpper(RevComp(seq))
	if strings.Contains(bigseq, revseq) {
		seqfound.Thing = revseq
		seqfound.Positions = search.Findall(bigseq, revseq)
		seqfound.Reverse = true
		seqsfound = append(seqsfound, seqfound)
	}

	return seqsfound
}

// FindSeq searches for a DNA sequence within a larger DNA sequence
// Will look for sequences which partially match based on mismatch tolerance
/*
Work in progress
func FindSeqWithMisMatches(bigsequence, smallsequence wtype.DNASequence, mismatchTolerance int) (seqsfound []search.Thingfound) {

	bigseq := strings.ToUpper(bigsequence.Sequence())

	var seqfound search.Thingfound

	seq := strings.ToUpper(smallsequence.Sequence())
	if strings.Contains(bigseq, seq) {
		seqfound.Thing = seq
		seqfound.Positions = search.Findall(bigseq, seq)
		seqsfound = append(seqsfound, seqfound)
	}

	revseq := strings.ToUpper(RevComp(seq))
	if strings.Contains(bigseq, revseq) {
		seqfound.Thing = revseq
		seqfound.Positions = search.Findall(bigseq, revseq)
		seqfound.Reverse = true
		seqsfound = append(seqsfound, seqfound)
	}

	return seqsfound
}
*/

// refactor strings to DNASequences to enable handling plasmid sequences
func FindSeqsinSeqs(bigseq string, smallseqs []string) (seqsfound []search.Thingfound) {

	bigseq = strings.ToUpper(bigseq)

	var seqfound search.Thingfound
	seqsfound = make([]search.Thingfound, 0)
	// fmt.Println("looking for ", smallseqs)
	for _, seq := range smallseqs {
		seq = strings.ToUpper(seq)
		if strings.Contains(bigseq, seq) {
			// fmt.Println("Fwd seq found")
			seqfound.Thing = seq
			seqfound.Positions = search.Findall(bigseq, seq)
			seqsfound = append(seqsfound, seqfound)
		}
	}
	for _, seq := range smallseqs {
		revseq := strings.ToUpper(RevComp(seq))
		if strings.Contains(bigseq, revseq) {
			// fmt.Println("rev seq found")
			seqfound.Thing = revseq
			seqfound.Positions = search.Findall(bigseq, revseq)
			seqfound.Reverse = true
			seqsfound = append(seqsfound, seqfound)
		}
	}

	return seqsfound
}

var Algorithmlookuptable = map[string]ReplacementAlgorithm{
	"ReplacebyComplement": ReplaceBycomplement,
}

// will potentially be generalisable for codon optimisation
type ReplacementAlgorithm func(sequence, thingtoreplace string, otherseqstoavoid []string) (replacement string, err error)

func ReplaceBycomplement(sequence, thingtoreplace string, otherseqstoavoid []string) (replacement string, err error) {

	seqsfound := FindSeqsinSeqs(sequence, []string{thingtoreplace})
	if len(seqsfound) == 1 {
		for _, instance := range seqsfound {
			if instance.Reverse == true {
				thingtoreplace = RevComp(thingtoreplace)
			}
		}

		allthingstoavoid := make([]string, len(otherseqstoavoid))
		allthingstoavoid = otherseqstoavoid
		allthingstoavoid = append(otherseqstoavoid, thingtoreplace)
		allthingstoavoid = search.RemoveDuplicates(allthingstoavoid)

		for i, _ := range thingtoreplace {

			replacementnucleotide := Comp(string(thingtoreplace[i]))
			replacement := strings.Replace(thingtoreplace, string(thingtoreplace[i]), replacementnucleotide, 1)
			newseq := strings.Replace(sequence, thingtoreplace, replacement, -1)
			checksitesfoundagain := FindSeqsinSeqs(newseq, allthingstoavoid)
			if len(checksitesfoundagain) == 0 {
				// fmt.Println("all things removed")
				return replacement, err
			}

		}

		for i, _ := range thingtoreplace {

			replacementnucleotide := Comp(thingtoreplace[i : i+1])
			replacement := strings.Replace(thingtoreplace, thingtoreplace[i:i+1], replacementnucleotide, 1)
			newseq := strings.Replace(sequence, thingtoreplace, replacement, -1)
			checksitesfoundagain := search.Findallthings(newseq, allthingstoavoid)
			if len(checksitesfoundagain) == 0 {
				// fmt.Println("all things removed, second try")
				return replacement, err
			}
			if i+2 == len(thingtoreplace) {
				specificseqs := text.Print("Specific Sequences", allthingstoavoid)
				err = fmt.Errorf("Not possible to remove site from sequence without avoiding the sequences to avoid using this algorithm; check specific sequences and adapt algorithm: ", specificseqs)
				break
			}
		}

	}
	return
}

// iterates through each position of a restriction site and replaces with the complementary base and then removes these from the main sequence
// if that fails the algorithm will attempt to find the complements of two adjacent positions. The algorithm needs improvement
func RemoveSiteOnestrand(sequence wtype.DNASequence, enzymeseq string, otherseqstoavoid []string) (newseq wtype.DNASequence, err error) {

	allthingstoavoid := make([]string, len(otherseqstoavoid))
	allthingstoavoid = otherseqstoavoid
	allthingstoavoid = append(otherseqstoavoid, enzymeseq)
	allthingstoavoid = append(otherseqstoavoid, RevComp(enzymeseq))

	for i, _ := range enzymeseq {

		replacementnucleotide := Comp(string(enzymeseq[i]))
		replacement := strings.Replace(enzymeseq, string(enzymeseq[i]), replacementnucleotide, 1)
		newseq.Seq = strings.Replace(sequence.Seq, enzymeseq, replacement, -1)
		checksitesfoundagain := FindSeqsinSeqs(newseq.Seq, allthingstoavoid)
		if len(checksitesfoundagain) == 0 {
			// fmt.Println("all things removed, first try")
			return
		}
	}

	for i, _ := range enzymeseq {

		replacementnucleotide := Comp(enzymeseq[i : i+1])
		replacement := strings.Replace(enzymeseq, enzymeseq[i:i+1], replacementnucleotide, 1)
		newseq.Seq = strings.Replace(sequence.Seq, enzymeseq, replacement, -1)
		checksitesfoundagain := search.Findallthings(newseq.Seq, allthingstoavoid)
		if len(checksitesfoundagain) == 0 {
			// fmt.Println("all things removed, second try")
			return
		}
		if i+2 == len(enzymeseq) {
			specificseqs := text.Print("Specific Sequences", allthingstoavoid)
			err = fmt.Errorf("Not possible to remove site from sequence without avoiding the sequences to avoid using this algorithm; check specific sequences and adapt algorithm: ", specificseqs)
			break
		}
	}

	return
}

// todo: fix this func
func RemoveSite(sequence wtype.DNASequence, enzyme wtype.RestrictionEnzyme, otherseqstoavoid []string) (newseq wtype.DNASequence, err error) {

	var tempseq wtype.DNASequence

	allthingstoavoid := make([]string, len(otherseqstoavoid))
	allthingstoavoid = otherseqstoavoid
	allthingstoavoid = append(allthingstoavoid, enzyme.RecognitionSequence)
	allthingstoavoid = append(allthingstoavoid, RevComp(enzyme.RecognitionSequence))

	seqsfound := FindSeqsinSeqs(sequence.Seq, []string{enzyme.RecognitionSequence})
	// fmt.Println("RemoveSite: ", seqsfound)
	if len(seqsfound) == 0 {
		return
	}

	thingtoreplace := enzyme.RecognitionSequence

	if len(seqsfound) == 1 {

		for _, instance := range seqsfound {
			if instance.Reverse == true {
				thingtoreplace = RevComp(enzyme.RecognitionSequence)
			}
		}

		tempseq, err = RemoveSiteOnestrand(sequence, thingtoreplace, allthingstoavoid)
		if err != nil {
			return tempseq, err
		}
		if tempseq.Seq != sequence.Seq {
			return tempseq, fmt.Errorf("New sequence is the same as old sequence")
		}
		newseq = sequence.Dup()
		newseq.Seq = tempseq.Seq
		return newseq, nil
	}

	if len(seqsfound) == 2 {

		tempseq, err := RemoveSiteOnestrand(sequence, thingtoreplace, allthingstoavoid)

		for _, instance := range seqsfound {
			if instance.Reverse == true {
				thingtoreplace = RevComp(enzyme.RecognitionSequence)
			}
		}

		tempseq, err = RemoveSiteOnestrand(tempseq, thingtoreplace, allthingstoavoid)
		if err != nil {
			return newseq, err
		}
		if tempseq.Seq != sequence.Seq {
			return newseq, fmt.Errorf("New sequence is the same as old sequence")
		}
		newseq = sequence.Dup()
		newseq.Seq = tempseq.Seq
		return newseq, nil

	}

	newseq = sequence.Dup()
	newseq.Seq = tempseq.Seq
	return
}

// this replaces all instances but this is not what we want
func ReplaceString(sequence string, seq string, otherseqstoavoid []string) (newseq string, err error) {

	allthingstoavoid := make([]string, len(otherseqstoavoid))
	allthingstoavoid = otherseqstoavoid
	allthingstoavoid = append(otherseqstoavoid, seq)

	for i, _ := range seq {

		replacementnucleotide := Comp(string(seq[i]))
		replacement := strings.Replace(seq, string(seq[i]), replacementnucleotide, 1)
		newseq = strings.Replace(sequence, seq, replacement, -1)
		checksitesfoundagain := search.Findallthings(newseq, allthingstoavoid)
		if len(checksitesfoundagain) == 0 {
			return
		}
	}

	for i, _ := range seq {

		replacementnucleotide := Comp(seq[i : i+1])
		replacement := strings.Replace(seq, seq[i:i+1], replacementnucleotide, 1)
		newseq = strings.Replace(sequence, seq, replacement, -1)
		checksitesfoundagain := search.Findallthings(newseq, allthingstoavoid)
		if len(checksitesfoundagain) == 0 {
			return
		}
		if i+2 == len(seq) {
			specificseqs := text.Print("Specific Sequences", allthingstoavoid)
			err = fmt.Errorf("Not possible to remove site from sequence without avoiding the sequences to avoid using this algorithm; check specific sequences and adapt algorithm: ", specificseqs)
			break
		}
	}
	return
}

/*
// working on this
func RemoveSiteFromORF(orf ORF, enzyme wtype.LogicalRestrictionEnzyme, otherseqstoavoid []string) (newseq wtype.DNASequence, err error) {

	allthingstoavoid := make([]string, len(otherseqstoavoid))
	allthingstoavoid = otherseqstoavoid
	allthingstoavoid = append(otherseqstoavoid, enzyme.RecognitionSequence)

	sitesfound := search.Findallthings(orf.DNASeq, allthingstoavoid)

	if len(sitesfound) == 0 {
		err = fmt.Errorf("no sites found in this Orf!", orf, enzyme, otherseqstoavoid)
		return
	}

	allpositions := search.Findall(orf.DNASeq, enzyme.RecognitionSequence)

	if len(allpositions) != 0 {
		for _, position := range allpositions {
			if orf.Direction != "Reverse" {

				_, _ = Codonfromposition(orf.DNASeq, position)

			}
		}
	}
	return
}
*/
/*
func RemoveSiteFromSeq(annotated AnnotatedSeq, enzyme wtype.LogicalRestrictionEnzyme, otherseqstoavoid []string) (newseq AnnotatedSeq, err error) {

	allthingstoavoid := make([]string, len(otherseqstoavoid))
	allthingstoavoid = otherseqstoavoid
	allthingstoavoid = append(otherseqstoavoid, enzyme.RecognitionSequence)

	sitesfound := search.Findallthings(orf.DNASeq, allthingstoavoid)

	if len(sitesfound) == 0 {
		err = fmt.Errorf("no sites found in this Orf!", orf, enzyme, otherseqstoavoid)
		return
	}

	allpositions := search.Findall(orf.DNASeq, enzyme.RecognitionSequence)

	if len(allpositions) != 0 {
		for _, position := range allpositions {
			if orf.Direction != "Reverse" {

				_ = Codonfromposition(orf.DNASeq, position)

			}
		}
	}
	return
}
*/
func RemoveSitesOutsideofFeatures(dnaseq wtype.DNASequence, site string, algorithm ReplacementAlgorithm, featurelisttoavoid []wtype.Feature) (newseq wtype.DNASequence, err error) {

	newseq = dnaseq

	pairs := make([]StartEndPair, 2)
	var pair StartEndPair

	for _, feature := range featurelisttoavoid {
		pair[0] = feature.StartPosition
		pair[1] = feature.EndPosition
		pairs = append(pairs, pair)
	}

	var otherseqstoavoid = []string{}

	replacement, err := algorithm(dnaseq.Seq, site, otherseqstoavoid)
	if err != nil {
		panic("choose different replacement choice func or change parameters")
	}

	newseq.Seq = ReplaceAvoidingPositionPairs(dnaseq.Seq, pairs, site, replacement)

	return
}

func ReplacefrombetweenPositions(seq string, start int, end int, original string, replacement string) (newseq string) {

	newseq = strings.Replace(seq[start-1:end-1], original, replacement, -1)
	return
}

func ReplaceAvoidingPositionPairs(seq string, positionpairs []StartEndPair, original string, replacement string) (newseq string) {

	temp := "£££££££££££"
	newseq = ""
	for _, pair := range positionpairs {
		if pair[0] < pair[1] {
			newseq = strings.Replace(seq[pair[0]-1:pair[1]-1], original, temp, -1)
		}
	}

	newseq = strings.Replace(newseq, original, replacement, -1)

	newseq = strings.Replace(newseq, temp, original, -1)

	// now look for reverse
	for _, pair := range positionpairs {
		if pair[0] > pair[1] {

			newseq = strings.Replace(seq[pair[1]+1:pair[0]+1], RevComp(original), temp, -1)
		}
	}

	newseq = strings.Replace(newseq, RevComp(original), RevComp(replacement), -1)

	newseq = strings.Replace(newseq, temp, RevComp(original), -1)
	return
}

type StartEndPair [2]int

func MakeStartendPair(start, end int) (pair StartEndPair) {

	pair[0] = start
	pair[1] = end
	return
}

func AAPosition(dnaposition int) (aaposition int) {

	remainder := dnaposition % 3
	aaposition = wutil.RoundInt(float64(dnaposition/3) + float64(remainder/3))

	return
}

func CodonOptions(codon string) (replacementoptions []string) {

	aa := DNAtoAASeq([]string{codon})

	replacementoptions = RevCodonTable[aa]
	return
}

func SwapCodon(codon string, position int) (replacement string) {

	replacementarray := CodonOptions(codon)

	replacement = replacementarray[position]

	return
}

func ReplaceCodoninORF(sequence wtype.DNASequence, startandendoforf StartEndPair, position int, seqstoavoid []string) (newseq wtype.DNASequence, codontochange string, option string, err error) {

	sequence.Seq = strings.ToUpper(sequence.Seq)

	// only handling cases where orf is not in reverse strand currently
	if startandendoforf[0] < startandendoforf[1] {

		if position < startandendoforf[0] || position > startandendoforf[1] {
			return sequence, codontochange, option, fmt.Errorf("position %d specified is out of range of orf start and finish specified %+v for %s", position, startandendoforf, sequence.Nm)

		}
		fmt.Println(1)
		seqslice := sequence.Seq[startandendoforf[0]-1 : startandendoforf[1]]
		orf, orftrue := FindORF(seqslice)
		fmt.Println(2, orf)
		if orftrue /*&& len(orf.DNASeq) == len(seqslice)*/ {
			codontochange, pair, err := Codonfromposition(seqslice, (position - startandendoforf[0]))
			if err != nil {
				return sequence, codontochange, option, err
			}

			options := CodonOptions(codontochange)

			for _, option := range options {
				tempseq := ReplacePosition(seqslice, pair, option)
				temporf, _ := FindORF(tempseq)

				sitesfound := search.Findallthings(tempseq, seqstoavoid)

				if temporf.ProtSeq == orf.ProtSeq && len(sitesfound) == 0 {
					newseq := sequence
					newseq.Seq = tempseq
					return newseq, codontochange, option, err
				}

			}
			err = fmt.Errorf("No satisfactory alternative codon options found to replace codon: %+v in options %+v", codontochange, options)
			return sequence, codontochange, option, err
		} else {
			err = fmt.Errorf("No orf found in sequence %s positions %d to %d", sequence.Nm, startandendoforf[0], startandendoforf[1])
			return sequence, codontochange, option, err
		}
	} else {
		newseq = sequence
		err = fmt.Errorf("orf in reverse direction, fix ReplaceCodoninORF func to handle this")
	}
	return
}

func ReplacePosition(sequence string, position StartEndPair, replacement string) (newseq string) {

	if position[0] < position[1] {
		one := sequence[0:position[0]]
		_ = sequence[position[0] : position[1]-1]
		three := sequence[position[1]:]

		newseq = one + replacement + three
	}
	return
}

func Codonfromposition(sequence string, dnaposition int) (codontoreturn string, position StartEndPair, err error) {

	if dnaposition > len(sequence) {
		return codontoreturn, position, fmt.Errorf("dnaposition %d is out of range of sequence length: %d", dnaposition, len(sequence))
	}

	nucleotides := []rune(sequence)
	res := ""
	aas := make([]string, 0)
	codon := ""
	for i, r := range nucleotides {
		res = res + string(r)

		if i > 0 && (i+1)%3 == 0 {
			codon = res
			aas = append(aas, res)
			res = ""
		}
		if i+1 > dnaposition && i > 0 && (i+1)%3 == 0 {
			if strings.ToUpper(aas[0]) != "ATG" {
				err = fmt.Errorf("sequence does not start with start codon ATG")
			}

			codontoreturn = codon
			position[1] = i + 1
			position[0] = i - 2

			return
		}
	}
	return codontoreturn, position, fmt.Errorf("No replacement codon found at position %d in sequence %s length %d", dnaposition, sequence, len(sequence))
}
