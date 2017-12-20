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

// type AminoAcid is a single letter format amino acid in string form.
// It can be validated as a valid AminoAcid using the SetAminoAcid function.
type AminoAcid string

// SetAminoAcid creates an AminoAcid from a string input and returns an error
// if the string is not a valid amino acid.
func SetAminoAcid(aa string) (AminoAcid, error) {

	if len(aa) != 1 {
		return "", fmt.Errorf("amino acid %s not valid. Please use single letter code.")
	}

	if err := wtype.ValidAA(aa); err != nil {
		return "", fmt.Errorf("amino acid %s not valid: %s", aa, err.Error())
	}
	return AminoAcid(strings.ToUpper(strings.TrimSpace(aa))), nil
}

// type Codon is a triplet of valid nucleotides which encodes an amino acid or stop codon.
// It can be validated using the SetCodon function.
type Codon string

// SetCodon creates a Codon from a string input and returns an error
// if the string is not a valid codon.
func SetCodon(dna string) (Codon, error) {

	if len(dna) != 3 {
		return "", fmt.Errorf("codon %s not valid. must be three nucleotides.")
	}

	if err := wtype.ValidDNA(dna); err != nil {
		return "", fmt.Errorf("codon %s not valid: %s", dna, err.Error())
	}
	return Codon(strings.ToUpper(strings.TrimSpace(dna))), nil
}

// A CodonUsageTable is an interface for any type which can convert an amino acid into a codon and error.
type CodonUsageTable interface {
	// ChooseCodon converts an amino acid into a codon.
	// A nil error is returned if this is done successfully.
	ChooseCodon(aminoAcid AminoAcid) (Codon, error)
}

// type SimpleUsageTable chooses the next codon as the first codon option from the Table field.
type SimpleUsageTable struct {
	// Table is a mapping between the amino acid and all codon options for that amino acid.
	Table map[string][]string
}

// ChooseCodon converts an amino acid into a codon.
// An error is returned if no value for the amino acid is found.
func (table SimpleUsageTable) ChooseCodon(aa AminoAcid) (codon Codon, err error) {
	codons, found := table.Table[string(aa)]

	if !found {
		return "", fmt.Errorf("%s not found in codon usage table", string(aa))
	}

	if len(codons) == 0 {
		return "", fmt.Errorf("0 codons found in codon usage table for %s", string(aa))
	}

	return Codon(codons[0]), nil
}

// type NTable converts each amino acid to NNN.
// This may be useful when a sequence is left to a DNA synthesis provider to codon optimise.
type NTable struct {
}

// ChooseCodon converts an amino acid into a codon.
// All amino acids will be converted to NNN; all stop codons to ***
func (table NTable) ChooseCodon(aa AminoAcid) (codon Codon, err error) {
	if aa == "*" {
		return Codon("***"), nil
	}

	return Codon("NNN"), nil
}

// type FrequencyTable chooses the next codon based on the frequency of the codon
// for that amino acid in the specified organism.
// for example:
// in Ecoli, F is encoded by TTT and TTC.
// The relative frequency of each is:
// TTT 0.58
// TTC 0.42
// The ChooseCodon method run on F would therefore return TTT 58% of the time and TTC 42%.
//
type FrequencyTable wtype.CodonTable

// ChooseCodon converts an amino acid into a codon.
// A nil error is returned if this is done successfully.
func (table FrequencyTable) ChooseCodon(aa AminoAcid) (codon Codon, err error) {

	codonTable := wtype.CodonTable(table)

	codonSeq := codonTable.ChooseWeighted(string(aa))

	if codonSeq == "" {
		return "", fmt.Errorf("codon not found in table for %v. Please set up Frequency Table first.", aa)
	}

	return Codon(codonSeq), nil
}

// Some example CodonUsageTables.
var (

	// Convert all amino acids to NNN; all stop codons to ***
	ConvertToN NTable = NTable{}

	// Return the first Codon value in the RevCodonTable for any amino acid.
	UseAnyCodon = SimpleUsageTable{Table: RevCodonTable}

	// EcoliTable is an example of a frequency table for E.Coli.
	// A codon for a specific amino acid will be returned with the probability set by the CodonSet
	//
	EColiTable = FrequencyTable{
		TaxID: "E.Coli",
		CodonByAA: map[string]wtype.CodonSet{
			"F": wtype.CodonSet{
				"TTT": 0.58,
				"TTC": 0.42,
			},
			"L": wtype.CodonSet{
				"TTA": 0.14,
				"TTG": 0.13,
				"CTT": 0.12,
				"CTC": 0.1,
				"CTA": 0.04,
				"CTG": 0.47,
			},
			"Y": wtype.CodonSet{
				"TAT": 0.59,
				"TAC": 0.41,
			},
			"*": wtype.CodonSet{
				"TAA": 0.61,
				"TAG": 0.09,
				"TGA": 0.3,
			},
			"H": wtype.CodonSet{
				"CAT": 0.57,
				"CAC": 0.43,
			},
			"Q": wtype.CodonSet{
				"CAA": 0.34,
				"CAG": 0.66,
			},
			"I": wtype.CodonSet{
				"ATT": 0.49,
				"ATC": 0.39,
				"ATA": 0.11,
			},
			"M": wtype.CodonSet{
				"ATG": 1.0,
			},
			"N": wtype.CodonSet{
				"AAT": 0.49,
				"AAC": 0.51,
			},
			"K": wtype.CodonSet{
				"AAA": 0.74,
				"AAG": 0.26,
			},
			"V": wtype.CodonSet{
				"GTT": 0.28,
				"GTC": 0.2,
				"GTA": 0.17,
				"GTG": 0.35,
			},
			"D": wtype.CodonSet{
				"GAT": 0.63,
				"GAC": 0.37,
			},
			"E": wtype.CodonSet{
				"GAA": 0.68,
				"GAG": 0.32,
			},
			"S": wtype.CodonSet{
				"TCT": 0.17,
				"TCC": 0.15,
				"TCA": 0.14,
				"TCG": 0.14,
				"AGT": 0.16,
				"AGC": 0.25,
			},
			"C": wtype.CodonSet{
				"TGT": 0.46,
				"TGC": 0.54,
			},
			"W": wtype.CodonSet{
				"TGG": 1,
			},
			"P": wtype.CodonSet{
				"CCT": 0.18,
				"CCC": 0.13,
				"CCA": 0.2,
				"CCG": 0.49,
			},
			"R": wtype.CodonSet{
				"CGT": 0.36,
				"CGC": 0.36,
				"CGA": 0.07,
				"CGG": 0.11,
				"AGA": 0.07,
				"AGG": 0.04,
			},
			"T": wtype.CodonSet{
				"ACT": 0.19,
				"ACC": 0.4,
				"ACA": 0.17,
				"ACG": 0.25,
			},
			"A": wtype.CodonSet{
				"GCT": 0.18,
				"GCC": 0.26,
				"GCA": 0.23,
				"GCG": 0.33,
			},
			"G": wtype.CodonSet{
				"GGT": 0.35,
				"GGC": 0.37,
				"GGA": 0.13,
				"GGG": 0.15,
			},
		},
		AAByCodon: Codontable,
	}
)

// RevTranslate converts an amino acid sequence into a dna sequence according the codon usage table specified.
// A CodonUsageTable is an interface for any type which has a ChooseCodon method.
// Examples of these are SimpleUsageTable, FrequencyTable and NTable
func RevTranslate(aaSeq wtype.ProteinSequence, codonUsageTable CodonUsageTable) (dnaSeq wtype.DNASequence, err error) {

	dnaSeq.SetName(aaSeq.Name())

	for _, aminoAcid := range aaSeq.Sequence() {

		aa, err := SetAminoAcid(string(aminoAcid))

		if err != nil {
			return dnaSeq, err
		}

		nextCodon, err := codonUsageTable.ChooseCodon(aa)

		if err != nil {
			return dnaSeq, err
		}

		dnaSeq.Append(string(nextCodon))
	}
	return dnaSeq, nil
}

// RevCodonTable describes the mapping between an amino acid in single letter format and the codons which encode it.
var RevCodonTable = map[string][]string{

	"N": []string{"AAC", "AAT"},
	"K": []string{"AAA", "AAG"},
	"T": []string{"ACC", "ACT", "ACA", "ACG"},
	"I": []string{"ATC", "ATT", "ATA"},
	"M": []string{"ATG"},
	"R": []string{"AGA", "AGG", "CGC", "CGT", "CGA", "CGG"},
	"Y": []string{"TAC", "TAT"},
	"*": []string{"TAA", "TAG", "TGA"},
	"S": []string{"AGC", "AGT", "TCC", "TCT", "TCA", "TCG"},
	"F": []string{"TTC", "TTT"},
	"L": []string{"TTA", "TTG", "CTC", "CTT", "CTA", "CTG"},
	"C": []string{"TGC", "TGT"},
	"W": []string{"TGG"},
	"D": []string{"GAC", "GAT"},
	"E": []string{"GAA", "GAG"},
	"V": []string{"GTC", "GTT", "GTA", "GTG"},
	"A": []string{"GCA", "GCC", "GCG", "GCT"},
	"G": []string{"GGC", "GGT", "GGA", "GGG"},
	"H": []string{"CAC", "CAT"},
	"Q": []string{"CAA", "CAG"},
	"P": []string{"CCC", "CCT", "CCA", "CCG"},
}

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

var (
	startcodons = []string{"ATG", "CTG", "GTG"}
)

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

	if reverse == false {
		orf, orftrue = FindORF(seq)
		orf.Direction = "Forward"
	}
	if reverse == true {
		revseq := RevComp(seq)
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
	if orftrue == false {
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
	if orftrue == false {
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
		if orftrue == false {
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

	return orfs
}

func FindNonOverlappingORFsinstrand(seq string) (orfs []ORF) {

	orfs = make([]ORF, 0)
	neworf, orftrue := FindORF(seq)
	if orftrue == false {
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
		if orftrue == false {
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
	revseq := RevComp(strings.ToUpper(seq))
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
	revseq := RevComp(strings.ToUpper(seq))
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

func MergeORFs(feats features) (orfs []ORF) {
	orfs = make([]ORF, 0)
	for _, top := range feats.TopstrandORFS {
		orfs = append(orfs, top)
	}
	// fmt.Println("TopStrandORFS: ", orfs)
	for _, bottom := range feats.BottomstrandORFS {
		orfs = append(orfs, bottom)
	}
	return
}

// should make this an interface
type features struct {
	TopstrandORFS    []ORF
	BottomstrandORFS []ORF
}
