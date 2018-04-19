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

// Package sequences is for interacting with and manipulating biological sequences; in extension to methods available in wtype
package sequences

import (
	"fmt"
	"strings"

	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/search"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/text"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wutil"
)

// Replace takes a PositionPair and replaces the sequence between the pair with the replaeWith sequence.
// Features will be deleted if part of the feature is replaced.
// Note, if used to delete sections from a plasmid, the sequence returned will be in plasmid form and it will be attempted to maintion the original orientation.
// In this case it may be necessary to rotate the sequence if looking to generate a linear sequence of interest.
func Replace(sequence wtype.DNASequence, position PositionPair, replaceWith wtype.DNASequence) (newSeq wtype.DNASequence, err error) {
	newSeq = sequence

	originalFeatures := sequence.Features

	start, end := position.Coordinates(wtype.CODEFRIENDLY)

	// check position pair is valid for this sequence
	if start > end && !sequence.Plasmid && !position.Reverse {
		return newSeq, fmt.Errorf("invalid position %+v to replace in sequence. Start position must be lower than end position unless position is reverse or sequence is a plasmid. Sequence %s: Plasmid = %v", position, sequence.Name(), sequence.Plasmid)
	}

	// convert sequence to replace into  string
	var replaceSeqString string = upper(replaceWith).Seq

	// if reverse, convert to reverse complement
	if position.Reverse {
		replaceSeqString = wtype.RevComp(upper(replaceWith).Seq)
	}

	// if it does not strand the end of the plasmid
	if !sequence.Plasmid {

		fmt.Println(len(newSeq.Seq), start, end)
		newSeq.Seq = newSeq.Seq[:start] + replaceSeqString + newSeq.Seq[end+1:]
		// if it does strand the end of the plasmid
	} else if sequence.Plasmid {

		if end > start {

			if !position.Reverse {

				newSeq.Seq = newSeq.Seq[:start] + replaceSeqString + newSeq.Seq[end+1:]

			}
		} else {

			if position.Reverse {
				newSeq.Seq = newSeq.Seq[:end] + replaceSeqString + newSeq.Seq[start+1:]
			} else {
				newSeq.Seq = replaceSeqString + newSeq.Seq[end+1:start]
			}
		}
	}

	SetFeatures(&newSeq, originalFeatures)

	return
}

// ReplaceAll searches for a sequence within a sequence and replaces all instances with the replaceWith sequence. Features will be deleted if part of the feature is replaced.
// Note, if used to delete sections from a plasmid, the sequence returned will be in plasmid form and it will be attempted to maintion the original orientation.
// In this case it may be necessary to rotate the sequence if looking to generate a linear sequence of interest.
func ReplaceAll(sequence, seqToReplace, replaceWith wtype.DNASequence) (newSeq wtype.DNASequence, err error) {
	searchResult := FindAll(&sequence, &seqToReplace)

	if len(seqToReplace.Seq) == 0 {
		return sequence, fmt.Errorf("no sequence to replace specified of %s to replace in %s ", seqToReplace.Name(), sequence.Name())

	}

	if len(searchResult.Positions) == 0 {
		return sequence, fmt.Errorf("no sequences of %s to replace in %s ", seqToReplace.Name(), sequence.Name())
	}

	newSeq = sequence

	originalFeatures := sequence.Features

	for _, position := range returnAllOrientationsOnly(searchResult) {

		start, end := position.Coordinates(wtype.CODEFRIENDLY)

		var replaceSeqString string = upper(seqToReplace).Seq

		// if reverse
		if position.Reverse {
			replaceSeqString = wtype.RevComp(upper(seqToReplace).Seq)
		}

		// if it does not strand the end of the plasmid
		if !sequence.Plasmid {

			if replaceSeqString != "" {

				newSeq.Seq = strings.Replace(upper(newSeq).Seq, replaceSeqString, upper(replaceWith).Seq, -1)

			}
			// if it does strand the end of the plasmid
		} else if sequence.Plasmid {

			if replaceSeqString != "" {

				if end > start {

					if replaceSeqString != "" {

						newSeq.Seq = strings.Replace(upper(newSeq).Seq, replaceSeqString, upper(replaceWith).Seq, -1)

					}
				} else {
					rotationSize := len(replaceSeqString) - 1

					replacedSeq := strings.Replace(Rotate(upper(newSeq), rotationSize, false).Seq, replaceSeqString, upper(replaceWith).Seq, -1)
					newSeq.Seq = replacedSeq
					newSeq = Rotate(newSeq, rotationSize, true)
				}
			}
		}

	}

	SetFeatures(&newSeq, originalFeatures)

	return
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
			if instance.Reverse {
				thingtoreplace = wtype.RevComp(thingtoreplace)
			}
		}

		allthingstoavoid := append(otherseqstoavoid, thingtoreplace)
		allthingstoavoid = search.RemoveDuplicateStrings(allthingstoavoid)

		for i := range thingtoreplace {

			replacementnucleotide := wtype.Comp(string(thingtoreplace[i]))
			replacement := strings.Replace(thingtoreplace, string(thingtoreplace[i]), replacementnucleotide, 1)
			newseq := strings.Replace(sequence, thingtoreplace, replacement, -1)
			checksitesfoundagain := FindSeqsinSeqs(newseq, allthingstoavoid)
			if len(checksitesfoundagain) == 0 {
				// fmt.Println("all things removed")
				return replacement, err
			}

		}

		for i := range thingtoreplace {

			replacementnucleotide := wtype.Comp(thingtoreplace[i : i+1])
			replacement := strings.Replace(thingtoreplace, thingtoreplace[i:i+1], replacementnucleotide, 1)
			newseq := strings.Replace(sequence, thingtoreplace, replacement, -1)
			checksitesfoundagain := search.FindAllStrings(newseq, allthingstoavoid)
			if len(checksitesfoundagain) == 0 {
				// fmt.Println("all things removed, second try")
				return replacement, err
			}
			if i+2 == len(thingtoreplace) {
				specificseqs := text.Sprint("Specific Sequences", allthingstoavoid)
				err = fmt.Errorf("Not possible to remove site from sequence without avoiding the sequences to avoid using this algorithm; check specific sequences and adapt algorithm: %v", specificseqs)
				break
			}
		}

	}
	return
}

// iterates through each position of a restriction site and replaces with the complementary base and then removes these from the main sequence
// if that fails the algorithm will attempt to find the complements of two adjacent positions. The algorithm needs improvement
func removeSiteOnestrand(sequence wtype.DNASequence, enzymeseq string, otherseqstoavoid []string) (newseq wtype.DNASequence, err error) {
	// XXX: Should probably add enzymeseq to allthingstoavoid as well, but to keep it
	// functionally equivalent to what it was don't do this at this time.

	//allthingstoavoid := append(otherseqstoavoid, enzymeseq)
	allthingstoavoid := append(otherseqstoavoid, wtype.RevComp(enzymeseq))

	for i := range enzymeseq {

		replacementnucleotide := wtype.Comp(string(enzymeseq[i]))
		replacement := strings.Replace(enzymeseq, string(enzymeseq[i]), replacementnucleotide, 1)
		newseq.Seq = strings.Replace(sequence.Seq, enzymeseq, replacement, -1)
		checksitesfoundagain := FindSeqsinSeqs(newseq.Seq, allthingstoavoid)
		if len(checksitesfoundagain) == 0 {
			// fmt.Println("all things removed, first try")
			return
		}
	}

	for i := range enzymeseq {

		replacementnucleotide := wtype.Comp(enzymeseq[i : i+1])
		replacement := strings.Replace(enzymeseq, enzymeseq[i:i+1], replacementnucleotide, 1)
		newseq.Seq = strings.Replace(sequence.Seq, enzymeseq, replacement, -1)
		checksitesfoundagain := search.FindAllStrings(newseq.Seq, allthingstoavoid)
		if len(checksitesfoundagain) == 0 {
			// fmt.Println("all things removed, second try")
			return
		}
		if i+2 == len(enzymeseq) {
			specificseqs := text.Sprint("Specific Sequences", allthingstoavoid)
			err = fmt.Errorf("Not possible to remove site from sequence without avoiding the sequences to avoid using this algorithm; check specific sequences and adapt algorithm: %v", specificseqs)
			break
		}
	}

	return
}

// todo: fix this func
func RemoveSite(sequence wtype.DNASequence, enzyme wtype.RestrictionEnzyme, otherseqstoavoid []string) (newseq wtype.DNASequence, err error) {

	var tempseq wtype.DNASequence

	allthingstoavoid := otherseqstoavoid
	allthingstoavoid = append(allthingstoavoid, enzyme.RecognitionSequence)
	allthingstoavoid = append(allthingstoavoid, wtype.RevComp(enzyme.RecognitionSequence))

	seqsfound := FindSeqsinSeqs(sequence.Seq, []string{enzyme.RecognitionSequence})
	// fmt.Println("RemoveSite: ", seqsfound)
	if len(seqsfound) == 0 {
		return
	}

	thingtoreplace := enzyme.RecognitionSequence

	if len(seqsfound) == 1 {

		for _, instance := range seqsfound {
			if instance.Reverse {
				thingtoreplace = wtype.RevComp(enzyme.RecognitionSequence)
			}
		}

		tempseq, err = removeSiteOnestrand(sequence, thingtoreplace, allthingstoavoid)
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

		tempseq, err := removeSiteOnestrand(sequence, thingtoreplace, allthingstoavoid)
		if err != nil {
			return newseq, err
		}

		for _, instance := range seqsfound {
			if instance.Reverse {
				thingtoreplace = wtype.RevComp(enzyme.RecognitionSequence)
			}
		}

		tempseq, err = removeSiteOnestrand(tempseq, thingtoreplace, allthingstoavoid)
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

			newseq = strings.Replace(seq[pair[1]+1:pair[0]+1], wtype.RevComp(original), temp, -1)
		}
	}

	newseq = strings.Replace(newseq, wtype.RevComp(original), wtype.RevComp(replacement), -1)

	newseq = strings.Replace(newseq, temp, wtype.RevComp(original), -1)
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

	aa := dNAtoAASeq([]string{codon})

	replacementoptions = RevCodonTable[aa]
	return
}

func ReplaceCodoninORF(sequence wtype.DNASequence, startandendoforf StartEndPair, position int, seqstoavoid []string) (newseq wtype.DNASequence, codontochange string, option string, err error) {

	sequence.Seq = strings.ToUpper(sequence.Seq)

	// only handling cases where orf is not in reverse strand currently
	if startandendoforf[0] < startandendoforf[1] {

		if position < startandendoforf[0] || position > startandendoforf[1] {
			return sequence, codontochange, option, fmt.Errorf("position %d specified is out of range of orf start and finish specified %+v for %s", position, startandendoforf, sequence.Nm)

		}
		seqslice := sequence.Seq[startandendoforf[0]-1 : startandendoforf[1]]
		orf, orftrue := FindORF(seqslice)
		if orftrue /*&& len(orf.DNASeq) == len(seqslice)*/ {
			codontochange, pair, err := Codonfromposition(seqslice, (position - startandendoforf[0]))
			if err != nil {
				return sequence, codontochange, option, err
			}

			options := CodonOptions(codontochange)

			for _, option := range options {
				tempseq := ReplacePosition(seqslice, pair, option)
				temporf, _ := FindORF(tempseq)

				sitesfound := search.FindAllStrings(tempseq, seqstoavoid)

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

func returnAllOrientationsOnly(searchResult SearchResult) (positions []PositionPair) {

	for _, position := range searchResult.Positions {
		if len(positions) == 2 {
			return positions
		} else if len(positions) == 1 {
			if positions[0].Reverse != position.Reverse {
				positions = append(positions, position)
			}
		} else {
			positions = append(positions, position)
		}
	}
	return
}
