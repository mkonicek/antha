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

	"github.com/antha-lang/antha/antha/anthalib/wtype"
)

// Codontable describes the mapping between a Codon and the amino acid which it encodes.
var Codontable = map[string]string{

	"AAC": "N",
	"AAT": "N",
	"AAA": "K",
	"AAG": "K",

	"ACC": "T",
	"ACT": "T",
	"ACA": "T",
	"ACG": "T",

	"ATC": "I",
	"ATT": "I",
	"ATA": "I",
	"ATG": "M",

	"AGC": "S",
	"AGT": "S",
	"AGA": "R",
	"AGG": "R",

	"TAC": "Y",
	"TAT": "Y",
	"TAA": "*",
	"TAG": "*",

	"TCC": "S",
	"TCT": "S",
	"TCA": "S",
	"TCG": "S",

	"TTC": "F",
	"TTT": "F",
	"TTA": "L",
	"TTG": "L",

	"TGC": "C",
	"TGT": "C",
	"TGA": "*",
	"TGG": "W",

	"GAC": "D",
	"GAT": "D",
	"GAA": "E",
	"GAG": "E",

	"GTC": "V",
	"GTT": "V",
	"GTA": "V",
	"GTG": "V",

	"GCA": "A",
	"GCC": "A",
	"GCG": "A",
	"GCT": "A",

	"GGC": "G",
	"GGT": "G",
	"GGA": "G",
	"GGG": "G",

	"CAC": "H",
	"CAT": "H",
	"CAA": "Q",
	"CAG": "Q",

	"CCC": "P",
	"CCT": "P",
	"CCA": "P",
	"CCG": "P",

	"CTC": "L",
	"CTT": "L",
	"CTA": "L",
	"CTG": "L",

	"CGC": "R",
	"CGT": "R",
	"CGA": "R",
	"CGG": "R",
}

func dNAtoAASeq(s []string) string {
	r := make([]string, 0)

	for _, c := range s {
		r = append(r, Codontable[string(c)])
	}
	rstring := strings.Join(r, "")
	return rstring
}

// type ORF is an open reading frame
type ORF struct {
	StartPosition int
	EndPosition   int
	DNASeq        string
	ProtSeq       string
	Direction     string
}

// molecular weight in g/mol
/*Molecular Weight notes:
The molecular weights above are those of the free acid and not the residue , which is used in the claculations performed by the Peptide Properties Calculator.
Subtracting an the weight of a mole of water (18g/mol) yields the molecular weight of the residue.
The weights used for Glx and Asx are averages.

http://www.basic.northwestern.edu/biotools/proteincalc.html

*/

// Molecularweight estimates molecular weight of a protein product.
func Molecularweight(orf ORF) (kDa float64) {
	aaarray := strings.Split(orf.ProtSeq, "")
	array := make([]float64, len(aaarray))
	for i := 0; i < len(aaarray); i++ {
		array = append(array, (aa_mw[aaarray[i]] - 18.0))
	}
	sum := 0.0
	for j := 0; j < len(array); j++ {
		sum += array[j]
	}
	kDa = sum / 1000
	return kDa

}

var aa_mw = map[string]float64{
	//1-letter Code	Molecular Weight (g/mol)
	"A": 89.09,
	"R": 174.2,
	"N": 132.12,
	"D": 133.1,
	"C": 121.16,
	"E": 147.13,
	"Q": 146.15,
	"G": 75.07,
	"H": 155.16,
	"I": 131.18,
	"L": 131.18,
	"K": 146.19,
	"M": 149.21,
	"F": 165.19,
	"P": 115.13,
	"S": 105.09,
	"T": 119.12,
	"W": 204.23,
	"Y": 181.19,
	"V": 117.15,
}

/*
type Promoter struct {
	StartPosition int
	EndPosition   int
	DNASeq        string
}

func FindPromoter (seq string) promoter Promoter {

	seq = strings.ToUpper(seq)



	if strings.Contains(seq,"TTGACA") {
		index := strings.Index(seq,"TTGACA")
		if strings.Index(seq+25,restofsequence := seq[index:]
		if
	}


}
*/
func FindStarts(seq string) (atgs int) {
	atgs = strings.Count(seq, "ATG") // extend later to include ctg, gtg etc
	return atgs
}

func FindDirectionalORF(seq string, reverse bool) (orf ORF, orftrue bool) {

	if !reverse {
		orf, orftrue = FindORF(seq)
		orf.Direction = "Forward"
	}
	if reverse {
		revseq := wtype.RevComp(seq)
		orf, orftrue = FindORF(revseq)
		orf.Direction = "Reverse"
		//tempend := orf.EndPosition
		//orf.DNASeq = RevComp(orf.DNASeq)
		orf.EndPosition = (len(seq) + 1 - orf.EndPosition)
		orf.StartPosition = (len(seq) + 1 - orf.StartPosition)
	}
	return orf, orftrue
}

func Translate(dna wtype.DNASequence) (aa wtype.ProteinSequence, err error) {
	orf, orftrue := FindORF(dna.Seq)
	if !orftrue {
		err = fmt.Errorf("Cannot translate this! no open reading frame detected")
		return
	} else {
		aa.Nm = dna.Nm + "Translated"
		aa.Seq = orf.ProtSeq
	}
	return
}

func FindORF(seq string) (orf ORF, orftrue bool) { // finds an orf in the forward direction only
	orftrue = false
	seq = strings.ToUpper(seq)

	if strings.Contains(seq, "ATG") {
		index := strings.Index(seq, "ATG")
		tempstart := index + 1
		//// fmt.Println("index=", index)
		restofsequence := seq[index:]
		codons := []rune(restofsequence)
		//// fmt.Println("codons=", string(codons))
		res := ""
		aas := make([]string, 0)
		for i, r := range codons {
			res = res + string(r)
			//fmt.Printf("i%d r %c\n", i, r)

			if i > 0 && (i+1)%3 == 0 {
				//fmt.Printf("=>(%d) '%v'\n", i, res)
				codon := res
				aas = append(aas, res)
				res = ""
				//// fmt.Println("codon=", codon)
				if codon == "TAA" {
					ORFcodons := aas
					//	// fmt.Println("orfcodons", ORFcodons)
					orf.StartPosition = tempstart
					orf.DNASeq = strings.Join(ORFcodons, "")
					orf.ProtSeq = dNAtoAASeq(ORFcodons)
					orf.EndPosition = orf.StartPosition + len(orf.DNASeq) - 1
					//// fmt.Println("translated=", translated)
				}
				if codon == "TGA" {
					ORFcodons := aas
					//	// fmt.Println("orfcodons", ORFcodons)
					orf.StartPosition = tempstart
					orf.DNASeq = strings.Join(ORFcodons, "")
					orf.ProtSeq = dNAtoAASeq(ORFcodons)
					orf.EndPosition = orf.StartPosition + len(orf.DNASeq) - 1
					//// fmt.Println("translated=", translated)
				}
				if codon == "TAG" {
					ORFcodons := aas
					//	// fmt.Println("orfcodons", ORFcodons)
					orf.StartPosition = tempstart
					orf.DNASeq = strings.Join(ORFcodons, "")
					orf.ProtSeq = dNAtoAASeq(ORFcodons)
					orf.EndPosition = orf.StartPosition + len(orf.DNASeq) - 1
					//// fmt.Println("translated=", translated)
				}
				if codon == "TAA" {
					orftrue = true
					return
				}
				if codon == "TGA" {
					orftrue = true
					return
				}
				if codon == "TAG" {
					orftrue = true
					return
				}
			}

		}

	}

	return orf, orftrue
}

func FindBiggestORF(seq string) (finalorf ORF, orftrue bool) { // finds an orf in the forward direction only

	var orf ORF
	orftrue = false
	seq = strings.ToUpper(seq)

	if strings.Contains(seq, "ATG") {
		index := strings.Index(seq, "ATG")
		tempstart := index + 1
		//// fmt.Println("index=", index)
		restofsequence := seq[index:]
		codons := []rune(restofsequence)
		//// fmt.Println("codons=", string(codons))
		res := ""
		aas := make([]string, 0)
		for i, r := range codons {
			res = res + string(r)
			//fmt.Printf("i%d r %c\n", i, r)

			if i > 0 && (i+1)%3 == 0 {
				//fmt.Printf("=>(%d) '%v'\n", i, res)
				codon := res
				aas = append(aas, res)
				res = ""
				//// fmt.Println("codon=", codon)
				if codon == "TAA" {
					ORFcodons := aas
					//	// fmt.Println("orfcodons", ORFcodons)
					orf.StartPosition = tempstart
					orf.DNASeq = strings.Join(ORFcodons, "")
					orf.ProtSeq = dNAtoAASeq(ORFcodons)
					orf.EndPosition = orf.StartPosition + len(orf.DNASeq) - 1
					//// fmt.Println("translated=", translated)
				}
				if codon == "TGA" {
					ORFcodons := aas
					//	// fmt.Println("orfcodons", ORFcodons)
					orf.StartPosition = tempstart
					orf.DNASeq = strings.Join(ORFcodons, "")
					orf.ProtSeq = dNAtoAASeq(ORFcodons)
					orf.EndPosition = orf.StartPosition + len(orf.DNASeq) - 1
					//// fmt.Println("translated=", translated)
				}
				if codon == "TAG" {
					ORFcodons := aas
					//	// fmt.Println("orfcodons", ORFcodons)
					orf.StartPosition = tempstart
					orf.DNASeq = strings.Join(ORFcodons, "")
					orf.ProtSeq = dNAtoAASeq(ORFcodons)
					orf.EndPosition = orf.StartPosition + len(orf.DNASeq) - 1
					//// fmt.Println("translated=", translated)
				}
				if codon == "TAA" {
					orftrue = true

				}
				if codon == "TGA" {
					orftrue = true

				}
				if codon == "TAG" {
					orftrue = true

				}
				if i == len(codons)-1 {
					finalorf = orf
					return
				}
			}

		}

	}

	return orf, orftrue
}

// finds all orfs and if they're greater than 20 amino acids (the smallest known protein) in length adds them to an array of orfs to be returned
func Findorfsinstrand(seq string) (orfs []ORF) {

	orfs = make([]ORF, 0)
	neworf, orftrue := FindORF(seq)
	if !orftrue {
		return
	}
	if len(neworf.ProtSeq) > 20 {
		orfs = append(orfs, neworf)
	}
	newseq := seq[(neworf.StartPosition):]
	orf1 := neworf
	i := 0
	for {

		neworf, orftrue := FindORF(newseq)
		if !orftrue {
			return
		}
		newseq = newseq[(neworf.StartPosition):]
		neworf.StartPosition = (neworf.StartPosition + orf1.StartPosition)
		neworf.EndPosition = (neworf.EndPosition + orf1.StartPosition)
		orf1 = neworf

		if len(neworf.ProtSeq) > 20 {
			orfs = append(orfs, neworf)
		}
		i++
	}
}

func FindNonOverlappingORFsinstrand(seq string) (orfs []ORF) {

	orfs = make([]ORF, 0)
	neworf, orftrue := FindORF(seq)
	if !orftrue {
		return
	}
	if len(neworf.ProtSeq) > 20 {
		orfs = append(orfs, neworf)
	}

	newseq := seq[(neworf.StartPosition):]
	orf1 := neworf
	i := 0
	for {
		neworf, orftrue := FindORF(newseq)
		if !orftrue {
			break
		}
		newseq = newseq[(neworf.EndPosition):]
		neworf.StartPosition = (neworf.StartPosition + orf1.StartPosition)
		neworf.EndPosition = (neworf.EndPosition + orf1.StartPosition)
		orf1 = neworf
		if len(neworf.ProtSeq) > 20 {
			orfs = append(orfs, neworf)
		}
		i++
	}

	return orfs
}

func LookforSpecificORF(seq string, targetAASeq string) (present bool) {
	ORFS := DoublestrandedORFS(seq)
	present = false
	for _, orf := range ORFS.TopstrandORFS {
		if strings.Contains(orf.ProtSeq, targetAASeq) {
			present = true
			return present
		}
	}
	for _, revorf := range ORFS.BottomstrandORFS {
		if strings.Contains(revorf.ProtSeq, targetAASeq) {
			present = true
		}
	}
	return present
}

// all orfs above 20 amino acids
func FindallORFs(seq string) []ORF {
	return MergeORFs(DoublestrandedORFS(seq))
}

func FindallNonOverlappingORFS(seq string) []ORF {
	return MergeORFs(DoublestrandedNonOverlappingORFS(seq))
}

func DoublestrandedORFS(seq string) (features features) {
	forwardorfs := Findorfsinstrand(seq)
	// fmt.Println("features.Top: ", features.TopstrandORFS)
	revseq := wtype.RevComp(strings.ToUpper(seq))
	reverseorfs := Findorfsinstrand(revseq)
	// fmt.Println("features.Bottom: ", reverseorfs)
	revORFpositionsreassigned := make([]ORF, 0)
	for _, orf := range reverseorfs {
		orf.Direction = "Reverse"
		orf.EndPosition = (len(seq) + 1 - orf.EndPosition)
		orf.StartPosition = (len(seq) + 1 - orf.StartPosition)
		revORFpositionsreassigned = append(revORFpositionsreassigned, orf)
	}
	features.BottomstrandORFS = revORFpositionsreassigned
	features.TopstrandORFS = forwardorfs
	return features
}

func DoublestrandedNonOverlappingORFS(seq string) (features features) {
	features.TopstrandORFS = FindNonOverlappingORFsinstrand(seq)
	revseq := wtype.RevComp(strings.ToUpper(seq))
	reverseorfs := FindNonOverlappingORFsinstrand(revseq)
	revORFpositionsreassigned := make([]ORF, 0)
	for _, orf := range reverseorfs {
		orf.Direction = "Reverse"
		orf.EndPosition = (len(seq) + 1 - orf.EndPosition)
		orf.StartPosition = (len(seq) + 1 - orf.StartPosition)
		revORFpositionsreassigned = append(revORFpositionsreassigned, orf)
	}
	features.BottomstrandORFS = revORFpositionsreassigned
	return features
}

func MergeORFs(feats features) []ORF {
	return append(feats.TopstrandORFS, feats.BottomstrandORFS...)
}

// should make this an interface
type features struct {
	TopstrandORFS    []ORF
	BottomstrandORFS []ORF
}
