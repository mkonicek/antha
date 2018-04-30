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

	"github.com/antha-lang/antha/antha/anthalib/wtype"
)

// RevCodonTable describes the mapping between an amino acid in single letter format and the codons which encode it.
var RevCodonTable = map[string][]string{

	"N": {"AAC", "AAT"},
	"K": {"AAA", "AAG"},
	"T": {"ACC", "ACT", "ACA", "ACG"},
	"I": {"ATC", "ATT", "ATA"},
	"M": {"ATG"},
	"R": {"AGA", "AGG", "CGC", "CGT", "CGA", "CGG"},
	"Y": {"TAC", "TAT"},
	"*": {"TAA", "TAG", "TGA"},
	"S": {"AGC", "AGT", "TCC", "TCT", "TCA", "TCG"},
	"F": {"TTC", "TTT"},
	"L": {"TTA", "TTG", "CTC", "CTT", "CTA", "CTG"},
	"C": {"TGC", "TGT"},
	"W": {"TGG"},
	"D": {"GAC", "GAT"},
	"E": {"GAA", "GAG"},
	"V": {"GTC", "GTT", "GTA", "GTG"},
	"A": {"GCA", "GCC", "GCG", "GCT"},
	"G": {"GGC", "GGT", "GGA", "GGG"},
	"H": {"CAC", "CAT"},
	"Q": {"CAA", "CAG"},
	"P": {"CCC", "CCT", "CCA", "CCG"},
}

// RevTranslate converts an amino acid sequence into a dna sequence according the codon usage table specified.
// A CodonUsageTable is an interface for any type which has a ChooseCodon method.
// Examples of these are SimpleUsageTable, FrequencyTable and NTable
func RevTranslate(aaSeq wtype.ProteinSequence, codonUsageTable CodonUsageTable) (dnaSeq wtype.DNASequence, err error) {

	dnaSeq.SetName(aaSeq.Name())

	if err = wtype.ValidAA(aaSeq.Sequence()); err != nil {
		return
	}

	for _, aminoAcid := range aaSeq.Sequence() {

		aa, _ := wtype.SetAminoAcid(string(aminoAcid))

		nextCodon, err := codonUsageTable.ChooseCodon(aa)

		if err != nil {
			return dnaSeq, err
		}

		err = dnaSeq.Append(string(nextCodon))
		if err != nil {
			return dnaSeq, err
		}
	}
	return dnaSeq, nil
}

// CodonUsageTable is an interface for any type which can convert an amino acid into a codon and error.
type CodonUsageTable interface {
	// ChooseCodon converts an amino acid into a codon.
	// A nil error is returned if this is done successfully.
	ChooseCodon(aminoAcid wtype.AminoAcid) (wtype.Codon, error)
}

// SimpleUsageTable contains a reverse translation table mapping of amino acid to all codon options.
// The first codon option for a specified Amino Acid is always chosen.
type SimpleUsageTable struct {
	// Table is a mapping between the amino acid and all codon options for that amino acid.
	Table map[string][]string
}

// ChooseCodon converts an amino acid into a codon.
// An error is returned if no value for the amino acid is found.
func (table SimpleUsageTable) ChooseCodon(aa wtype.AminoAcid) (codon wtype.Codon, err error) {
	codons, found := table.Table[string(aa)]

	if !found {
		return "", fmt.Errorf("%s not found in codon usage table", string(aa))
	}

	if len(codons) == 0 {
		return "", fmt.Errorf("0 codons found in codon usage table for %s", string(aa))
	}

	return wtype.Codon(codons[0]), nil
}

// type NTable converts each amino acid to NNN.
// This may be useful when a sequence is left to a DNA synthesis provider to codon optimise.
type NTable struct {
}

// ChooseCodon converts an amino acid into a codon.
// All amino acids will be converted to NNN; all stop codons to ***
func (table NTable) ChooseCodon(aa wtype.AminoAcid) (codon wtype.Codon, err error) {
	if aa == "*" {
		return wtype.Codon("***"), nil
	}

	return wtype.Codon("NNN"), nil
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
func (table FrequencyTable) ChooseCodon(aa wtype.AminoAcid) (codon wtype.Codon, err error) {

	codonTable := wtype.CodonTable(table)

	codonSeq := codonTable.ChooseWeighted(string(aa))

	if codonSeq == "" {
		return "", fmt.Errorf("codon not found in table for %v. Please set up Frequency Table first.", aa)
	}

	return wtype.Codon(codonSeq), nil
}

// Some example CodonUsageTables.
var (

	// Convert all amino acids to NNN; all stop codons to ***
	ConvertToNNN NTable = NTable{}

	// Return the first Codon value in the RevCodonTable for any amino acid.
	UseAnyCodon = SimpleUsageTable{Table: RevCodonTable}

	// EcoliTable is an example of a frequency table for E.Coli.
	// A codon for a specific amino acid will be returned with the probability set by the CodonSet
	//
	EColiTable = FrequencyTable{
		TaxID: "E.Coli",
		CodonByAA: map[string]wtype.CodonSet{
			"F": {
				"TTT": 0.58,
				"TTC": 0.42,
			},
			"L": {
				"TTA": 0.14,
				"TTG": 0.13,
				"CTT": 0.12,
				"CTC": 0.1,
				"CTA": 0.04,
				"CTG": 0.47,
			},
			"Y": {
				"TAT": 0.59,
				"TAC": 0.41,
			},
			"*": {
				"TAA": 0.61,
				"TAG": 0.09,
				"TGA": 0.3,
			},
			"H": {
				"CAT": 0.57,
				"CAC": 0.43,
			},
			"Q": {
				"CAA": 0.34,
				"CAG": 0.66,
			},
			"I": {
				"ATT": 0.49,
				"ATC": 0.39,
				"ATA": 0.11,
			},
			"M": {
				"ATG": 1.0,
			},
			"N": {
				"AAT": 0.49,
				"AAC": 0.51,
			},
			"K": {
				"AAA": 0.74,
				"AAG": 0.26,
			},
			"V": {
				"GTT": 0.28,
				"GTC": 0.2,
				"GTA": 0.17,
				"GTG": 0.35,
			},
			"D": {
				"GAT": 0.63,
				"GAC": 0.37,
			},
			"E": {
				"GAA": 0.68,
				"GAG": 0.32,
			},
			"S": {
				"TCT": 0.17,
				"TCC": 0.15,
				"TCA": 0.14,
				"TCG": 0.14,
				"AGT": 0.16,
				"AGC": 0.25,
			},
			"C": {
				"TGT": 0.46,
				"TGC": 0.54,
			},
			"W": {
				"TGG": 1,
			},
			"P": {
				"CCT": 0.18,
				"CCC": 0.13,
				"CCA": 0.2,
				"CCG": 0.49,
			},
			"R": {
				"CGT": 0.36,
				"CGC": 0.36,
				"CGA": 0.07,
				"CGG": 0.11,
				"AGA": 0.07,
				"AGG": 0.04,
			},
			"T": {
				"ACT": 0.19,
				"ACC": 0.4,
				"ACA": 0.17,
				"ACG": 0.25,
			},
			"A": {
				"GCT": 0.18,
				"GCC": 0.26,
				"GCA": 0.23,
				"GCG": 0.33,
			},
			"G": {
				"GGT": 0.35,
				"GGC": 0.37,
				"GGA": 0.13,
				"GGG": 0.15,
			},
		},
		AAByCodon: Codontable,
	}
)

// RevTranslatetoNstring converts a string amino acid sequence to a sequence of NNN codons.
func RevTranslatetoNstring(aaSeq string) (NNN string) {

	var codonSeq wtype.DNASequence

	for _, aa := range aaSeq {
		aminoAcid, err := wtype.SetAminoAcid(string(aa))

		if err != nil {
			panic(err)
		}

		codon, err := ConvertToNNN.ChooseCodon(aminoAcid)

		if err != nil {
			panic(err)
		}

		err = codonSeq.Append(string(codon))

		if err != nil {
			panic(err)
		}

	}

	return codonSeq.Sequence()

}
