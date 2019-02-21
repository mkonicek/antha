// wtype/biology.go: Part of the Antha language
// Copyright (C) 2014 the Antha authors. All rights reserved.
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

package wtype

import (
	"encoding/json"
	"fmt"
	"strings"

	. "github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/sequences/biogo/ncbi/blast"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/sequences/blast"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
)

// the following are all physical things; we need a way to separate
// out just the logical part

// structure which defines an enzyme -- solutions containing
// enzymes need careful handling as they can be quite delicate
type Enzyme struct {
	Properties map[string]wunit.Measurement
	Nm         string
}

func (enzyme Enzyme) Name() string {
	return enzyme.Nm
}

// RestrictionEnzyme is an enzyme which cleaves DNA
type RestrictionEnzyme struct {
	Enzyme
	// sequence
	RecognitionSequence               string
	EndLength                         int
	Prototype                         string
	Topstrand3primedistancefromend    int
	Bottomstrand5primedistancefromend int
	MethylationSite                   string   //"attr, <4>"
	CommercialSource                  []string //string "attr, <5>"
	References                        []int
	Class                             string
	Isoschizomers                     []string
}

type TypeIIs struct {
	RestrictionEnzyme
}

func ToTypeIIs(typeIIenzyme RestrictionEnzyme) (typeIIsenz TypeIIs, err error) {
	if typeIIenzyme.Class == "TypeII" {
		err = fmt.Errorf("You can't do this, enzyme is not a type IIs")
		return
	}
	if typeIIenzyme.Class == "TypeIIs" {

		typeIIsenz = TypeIIs{RestrictionEnzyme: typeIIenzyme}

	}
	return
}

// structure which defines an organism. These need specific handling
// -- some detail is derived using the TOL structure
type Organism struct {
	Species *TOL // position on the TOL
}

// a set of organisms, can be mixed or homogeneous
type Population struct {
}

// defines a plasmid
type Plasmid struct {
}

// defines things which have biosequences... useful for operations
// valid on biosequences such as BLASTing / other alignment methods
type BioSequence interface {
	Name() string
	Sequence() string
	Append(string) error
	Prepend(string) error
	Blast() ([]Hit, error)
	MolecularWeight() float64
}

// defines something as physical DNA
// hence it is physical and has a DNASequence
type DNA struct {
	Seq DNASequence
}

// DNAsequence is a type of Biosequence
type DNASequence struct {
	Nm             string    `json:"nm"`
	Seq            string    `json:"seq"`
	Plasmid        bool      `json:"plasmid"`
	Singlestranded bool      `json:"single_stranded"`
	Overhang5prime Overhang  `json:"overhang_5_prime"`
	Overhang3prime Overhang  `json:"overhang_3_prime"`
	Methylation    string    `json:"methylation"` // add histones etc?
	Features       []Feature `json:"features"`
}

func (seq DNASequence) Dup() DNASequence {
	var ret DNASequence

	d, err := json.Marshal(seq)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(d, &ret)
	if err != nil {
		panic(err)
	}

	return ret
}

// AddOverhang adds an overhang to a specified end.
// Valid options are either 5 (for 5 prime/upstream end) or 3 (for 3 prime/upstream end).
func (seq *DNASequence) AddOverhang(end int, sequence string) (err error) {

	var overHang Overhang

	if end == 5 {

		overHang, err = MakeOverHang(sequence, 5, TOP, true)

		if err != nil {
			return
		}

		err = seq.Set5PrimeEnd(overHang)

		return err

	} else if end == 3 {

		overHang, err = MakeOverHang(sequence, 3, TOP, true)

		if err != nil {
			return
		}

		err = seq.Set3PrimeEnd(overHang)

		return err
	}
	return fmt.Errorf("cannot add overhang to end %d. Please choose either 5 (for 5 prime/upstream end) or 3 (for 3 prime/upstream end) ", end)
}

// AddUnderhang adds an underhang to a specified end.
// Valid options are either 5 (for 5 prime/upstream end) or 3 (for 3 prime/upstream end).
func (seq *DNASequence) AddUnderhang(end int, sequence string) (err error) {

	var overHang Overhang

	if end == 5 {

		overHang, err = MakeOverHang(sequence, 5, BOTTOM, true)

		if err != nil {
			return
		}

		err = seq.Set5PrimeEnd(overHang)

		return err
	} else if end == 3 {

		seq.Overhang3prime, err = MakeOverHang(sequence, 3, BOTTOM, true)

		if err != nil {
			return
		}

		err = seq.Set3PrimeEnd(overHang)

		return err
	}
	return fmt.Errorf("cannot add overhang to end %d. Please choose either 5 (for 5 prime/upstream end) or 3 (for 3 prime/upstream end) ", end)
}

// AddBluntOverhang adds a blunt overhang to a specified end.
// Valid options are either 5 (for 5 prime/upstream end) or 3 (for 3 prime/upstream end).
func (seq *DNASequence) AddBluntEnd(end int) (err error) {
	if end == 5 {
		seq.Overhang5prime, err = MakeOverHang("", 5, NEITHER, true)
		return err
	} else if end == 3 {
		seq.Overhang3prime, err = MakeOverHang("", 3, NEITHER, true)
		return err
	}

	return fmt.Errorf("cannot add blunt end to end %d. Please choose either 5 (for 5 prime/upstream end) or 3 (for 3 prime/upstream end) ", end)
}

// Set5PrimeEnd adds a 5 prime overhang to a sequence.
// Validation is performed on the compatibility of the overhang with the sequence.
func (seq *DNASequence) Set5PrimeEnd(overhang Overhang) (err error) {
	if seq.Singlestranded {
		err = fmt.Errorf("Can't have overhang on single stranded dna")
		return
	}
	if seq.Plasmid {
		err = fmt.Errorf("Can't have overhang on Plasmid(circular) dna")
		return
	}

	if overhang.Type == OVERHANG {

		expectedOverhang := Prefix(seq.Seq, len(overhang.Sequence()))

		if !strings.EqualFold(expectedOverhang, overhang.Sequence()) {
			return fmt.Errorf("specified overhang %s to add to 5' end of sequence %s is not equal to sequence prefix %s", overhang.Sequence(), seq.Name(), expectedOverhang)
		}
	}

	seq.Overhang5prime = overhang
	return nil
}

// Set3PrimeEnd adds a 3 prime overhang to a sequence.
// Validation is performed on the compatibility of the overhang with the sequence.
func (seq *DNASequence) Set3PrimeEnd(overhang Overhang) (err error) {
	if seq.Singlestranded {
		err = fmt.Errorf("Can't have overhang on single stranded dna")
		return
	}
	if seq.Plasmid {
		err = fmt.Errorf("Can't have overhang on Plasmid(circular) dna")
		return
	}

	if overhang.Type == OVERHANG {

		expectedOverhang := Prefix(RevComp(seq.Seq), len(overhang.Sequence()))

		if !strings.EqualFold(expectedOverhang, overhang.Sequence()) {
			return fmt.Errorf("specified underhang %s to add to 3' end of sequence %s is not equal to sequence suffix %s", overhang.Sequence(), seq.Name(), expectedOverhang)
		}
	}

	seq.Overhang3prime = overhang
	return nil
}

func MakeDNASequence(name string, seqstring string, properties []string) (seq DNASequence, err error) {
	seq.Nm = name
	seq.Seq = seqstring
	for _, property := range properties {
		property = strings.ToUpper(property)

		if strings.Contains(property, "DCM") || strings.Contains(property, "DAM") || strings.Contains(property, "CPG") {
			seq.Methylation = property
		}

		if strings.Contains(property, "PLASMID") || strings.Contains(property, "CIRCULAR") || strings.Contains(property, "VECTOR") {
			seq.Plasmid = true
			break
		}
		if strings.Contains(property, "SS") || strings.Contains(property, "SINGLE STRANDED") {
			seq.Singlestranded = true
			break
		}
	}
	return
}
func MakeLinearDNASequence(name string, seqstring string) (seq DNASequence) {
	seq.Nm = name
	seq.Seq = seqstring

	return
}
func MakePlasmidDNASequence(name string, seqstring string) (seq DNASequence) {
	seq.Nm = name
	seq.Seq = seqstring
	seq.Plasmid = true
	return
}
func MakeSingleStrandedDNASequence(name string, seqstring string) (seq DNASequence) {
	seq.Nm = name
	seq.Seq = seqstring
	seq.Singlestranded = true
	return
}

// MakeOverHang is used to create an overhang.
func MakeOverHang(overhangSequence string, end int, toporbottom int, phosphorylated bool) (overhang Overhang, err error) {

	length := len(overhangSequence)

	if end == 0 {
		err = fmt.Errorf("if end = 0, all fields are returned empty")
		return
	}

	if end == 5 || end == 3 || end == 0 {
		overhang.End = end
	} else {
		err = fmt.Errorf("invalid entry for end: 5PRIME = 5, 3PRIME = 3, NA = 0")
		return
	}
	if toporbottom == NEITHER && length == 0 {
		overhang.Type = BLUNT
		return
	}
	if toporbottom == NEITHER && length != 0 {
		err = fmt.Errorf("If length of overhang is not 0, toporbottom must be 0")
		return
	}
	if toporbottom != NEITHER && length == 0 {
		err = fmt.Errorf("If length of overhang is 0, toporbottom must be 0")
		return
	}
	if toporbottom > 2 {
		err = fmt.Errorf("invalid entry for toporbottom: NEITHER = 0, TOP = 1, BOTTOM = 2")
		return
	}
	if end == 5 {
		if toporbottom == TOP {
			overhang.Type = OVERHANG
		}
		if toporbottom == BOTTOM {
			overhang.Type = UNDERHANG
		}
	} else if end == 3 {
		if toporbottom == TOP {
			overhang.Type = OVERHANG
		}
		if toporbottom == BOTTOM {
			overhang.Type = UNDERHANG
		}
	}

	overhang.Seq = overhangSequence
	overhang.Phosphorylation = phosphorylated
	return
}

func Phosphorylate(dnaseq DNASequence) (phosphorylateddna DNASequence, err error) {
	if dnaseq.Plasmid {
		err = fmt.Errorf("Can't phosphorylate circular dna")
		phosphorylateddna = dnaseq
		return
	}
	if dnaseq.Overhang5prime.Type != 0 {
		dnaseq.Overhang5prime.Phosphorylation = true
	}
	if dnaseq.Overhang3prime.Type != 0 {
		dnaseq.Overhang3prime.Phosphorylation = true
	}
	if dnaseq.Overhang3prime.Type == 0 && dnaseq.Overhang5prime.Type == 0 {
		err = fmt.Errorf("No ends available, but not plasmid! This doesn't seem possible!")
		phosphorylateddna = dnaseq
	}
	return
}

// OverHangType represents the type of an overhang.
// Valid options are
// 	FALSE     OverHangType = 0
//	BLUNT     OverHangType = 1
//	OVERHANG  OverHangType = 2
//	UNDERHANG OverHangType = -1
type OverHangType int

// Valid overhang types
const (
	// no overhang
	FALSE OverHangType = 0
	// A blunt overhang
	BLUNT OverHangType = 1
	// An overhang (5' sequence overhangs complementary strand)
	OVERHANG OverHangType = 2
	// an underhang (5' sequence underhangs complementary strand)
	UNDERHANG OverHangType = -1
)

// Options for Strand choice
const (
	NEITHER = 0
	// Top strand, or coding strand
	TOP = 1
	// Bottom strand, or complimentary strand.
	BOTTOM = 2
)

// Overhang represents an end of a DNASequence.
type Overhang struct {
	// Valid options are 5 (5 Prime end), 3 (3 prime end) or 0 (nul)
	End int `json:"end"`
	// Valid options are FALSE, BLUNT, OVERHANG, UNDERHANG
	Type OverHangType `json:"type"`
	// Overhang sequence
	Seq string `json:"sequence"`
	// Whether the overhang is phosphorylated.
	Phosphorylation bool `json:"phosphorylation"`
}

// Sequence returns the sequence of the overhang.
func (oh Overhang) Sequence() string {
	return oh.Seq
}

// Length returns the length of the overhang.
func (oh Overhang) Length() int {
	return len(oh.Sequence())
}

// ToString returns a string summary of the overhang.
func (oh Overhang) ToString() string {
	if oh.End == 5 {
		if oh.Type == OVERHANG {
			return `5' overhang: ` + oh.Sequence()
		}
		if oh.Type == BLUNT || oh.Type == FALSE {
			return `5' Blunt`
		}
		if oh.Type == UNDERHANG {
			return `5' underhang: ` + oh.Sequence()
		}

	}

	if oh.End == 3 {
		if oh.Type == OVERHANG {
			return `3' overhang: ` + oh.Sequence()
		}
		if oh.Type == BLUNT || oh.Type == FALSE {
			return `3' Blunt`
		}
		if oh.Type == UNDERHANG {
			return `3' underhang: ` + oh.Sequence()
		}

	}
	return ""
}

// TypeName returns the name of the overhang type as a string.
func (oh Overhang) TypeName() string {
	if oh.Type == OVERHANG {
		return "Overhang"
	} else if oh.Type == UNDERHANG {
		return "Underhang"
	} else if oh.Type == BLUNT {
		return "blunt"
	}
	return "no overhang"
}

// Overhang returns the sequence if the overhang is of type OVERHANG.
func (oh Overhang) OverHang() (sequence string) {
	if oh.Type == OVERHANG {
		return oh.Sequence()
	} else if oh.Type == BLUNT {
		return "blunt"
	}
	return ""
}

// Overhang returns any sequence if the underhang is of type UNDERHANG.
func (oh Overhang) UnderHang() (sequence string) {
	if oh.Type == UNDERHANG {
		return oh.Sequence()
	} else if oh.Type == BLUNT {
		return "blunt"
	}

	return ""
}

func valid(seq, validOptions string) error {
	var errs []string

	for i, character := range seq {
		if !strings.Contains(validOptions, strings.ToUpper(string(character))) {
			errs = append(errs, fmt.Sprint(string(character), " found at position ", i+1, "; "))
		}
	}
	if len(errs) > 0 {
		return fmt.Errorf("Invalid characters found %v", errs)
	}
	return nil
}

// ValidDNA checks a sequence given as a string for validity as a DNASequence.
// Any IUPAC nucleotide is considered valid, not just ACTG.
func ValidDNA(seq string) error {
	validNucleotides := "ACTGNXBHVDMKSWRYU"

	return valid(seq, validNucleotides)
}

// ValidRNA checks a sequence given as a string for validity as an RNASequence.
// ACGU are valid entries.
func ValidRNA(seq string) error {
	validRNA := "ACGU"

	return valid(seq, validRNA)
}

// ValidAA checks a sequence given as a string for validity as a ProteinSequence.
// All standard single letter AminoAcids are valid as well as * indicating stop.
func ValidAA(seq string) error {

	var aminoAcids []string

	for key := range aa_mw {
		aminoAcids = append(aminoAcids, key)
	}

	// add stop
	aminoAcids = append(aminoAcids, "*")

	validAminoAcids := strings.Join(aminoAcids, "")

	return valid(seq, validAminoAcids)
}

func upper(s string) string {
	return trimString(strings.ToUpper(s))
}

func trimString(str string) string {
	return strings.TrimSpace(str)
}

// Sequence returns the sequence of the DNA Sequence
func (dna *DNASequence) Sequence() string {
	return dna.Seq
}

// Name returns the name of the DNASequence
func (dna *DNASequence) Name() string {
	return dna.Nm
}

// SetName sets the names of the dna sequence
func (dna *DNASequence) SetName(name string) {
	dna.Nm = trimString(name)
}

// RevComp returns the reverse complement of the DNASequence.
func (dna *DNASequence) RevComp() string {
	return RevComp(dna.Seq)
}

// SetSequence checks the validity of sequence given as an argument and if all characters are present in ValidNucleotides
// If invalid characters are found an error returned listing all invalid characters and their positions in human friendly form. i.e. the first position is 1 and not 0.
func (dna *DNASequence) SetSequence(seq string) error {

	dna.Seq = seq

	return ValidDNA(seq)
}

// Append appends the existing dna sequence with the upper case of the string added
func (dna *DNASequence) Append(s string) error {
	err := ValidDNA(s)
	if err != nil {
		return fmt.Errorf("invalid characters requested for Append: %s", err.Error())
	}

	dna.Seq = dna.Seq + s

	return nil
}

// Preprend adds the requested sequence to the beginning of the existing sequence.
func (dna *DNASequence) Prepend(s string) error {

	err := ValidDNA(s)

	if err != nil {
		return fmt.Errorf("invalid characters requested for Prepend: %s", err.Error())
	}

	dna.Seq = s + dna.Seq

	return nil
}

// Blast performs a blast search on the sequence and returns any hits found.
// An error is returned if a problem interacting with the blast server occurs.
func (seq *DNASequence) Blast() (hits []Hit, err error) {
	hits, err = blast.MegaBlastN(seq.Seq)
	return
}

var nucleotidegpermol = map[string]float64{
	"A":    313.2,
	"T":    304.2,
	"C":    289.2,
	"G":    329.2,
	"N":    303.7,
	"dATP": 491.2,
	"dCTP": 467.2,
	"dGTP": 507.2,
	"dTTP": 482.2,
	"dNTP": 487.0,
}

// MolecularWeight calculates the molecular weight of the specified DNA sequence.
// For accuracy it is important to specify if the DNA is single stranded or doublestranded along with phosphorylation.
func (seq *DNASequence) MolecularWeight() float64 {
	//Calculate Molecular weight of DNA

	// need to add effect of methylation on molecular weight
	fwdsequence := seq.Seq
	phosphate5prime := seq.Overhang5prime.Phosphorylation
	phosphate3prime := seq.Overhang3prime.Phosphorylation
	singlestranded := seq.Singlestranded

	upperCase := func(s string) string { return strings.ToUpper(s) }

	numberofAs := strings.Count(upperCase(fwdsequence), "A")
	numberofTs := strings.Count(upperCase(fwdsequence), "T")
	numberofCs := strings.Count(upperCase(fwdsequence), "C")
	numberofGs := strings.Count(upperCase(fwdsequence), "G")
	massofAs := (float64(numberofAs) * nucleotidegpermol["A"])
	massofTs := (float64(numberofTs) * nucleotidegpermol["T"])
	massofCs := (float64(numberofCs) * nucleotidegpermol["C"])
	massofGs := (float64(numberofGs) * nucleotidegpermol["G"])
	mw := (massofAs + massofTs + massofCs + massofGs)
	if phosphate5prime {
		mw = mw + 79.0 // extra for phosphate left at 5' end following digestion, not relevant for primer extension
	}
	if phosphate3prime {
		mw = mw + 79.0 // extra for phosphate left at 3' end following digestion, not relevant for primer extension
	}
	if !singlestranded {
		mw = 2 * mw
	}
	return mw
}

// RNA sample: physical RNA, has an RNASequence object
type RNA struct {
	Seq RNASequence
}

// RNASequence object is a type of Biosequence
type RNASequence struct {
	Nm  string
	Seq string
}

func (rna *RNASequence) Sequence() string {
	return rna.Seq
}

func (rna *RNASequence) SetSequence(seq string) error {
	rna.Seq = upper(seq)
	return ValidRNA(seq)
}

func (rna *RNASequence) Name() string {
	return rna.Nm
}

func (rna *RNASequence) SetName(name string) {
	rna.Nm = trimString(name)
}

func (rna *RNASequence) Append(s string) error {
	err := ValidRNA(s)

	if err != nil {
		return fmt.Errorf("invalid characters requested for Append: %s", err.Error())
	}

	rna.Seq = rna.Seq + s
	return nil
}

func (rna *RNASequence) Prepend(s string) error {

	err := ValidRNA(s)

	if err != nil {
		return fmt.Errorf("invalid characters requested for Prepend: %s", err.Error())
	}

	rna.Seq = s + rna.Seq
	return nil
}

func (seq *RNASequence) Blast() (hits []Hit, err error) {
	hits, err = blast.MegaBlastN(seq.Seq)
	return
}

// physical protein sample
// has a ProteinSequence
type Protein struct {
	Seq ProteinSequence
}

// AminoAcid is a single letter format amino acid in string form.
// It can be validated as a valid AminoAcid using the SetAminoAcid function.
type AminoAcid string

// SetAminoAcid creates an AminoAcid from a string input and returns an error
// if the string is not a valid amino acid.
func SetAminoAcid(aa string) (AminoAcid, error) {

	if len(aa) != 1 {
		return "", fmt.Errorf("amino acid \"%s\" not valid. Please use single letter code.", aa)
	}

	if err := ValidAA(aa); err != nil {
		return "", fmt.Errorf("amino acid \"%s\" not valid: %s", aa, err.Error())
	}
	return AminoAcid(strings.ToUpper(strings.TrimSpace(aa))), nil
}

// Codon is a triplet of valid nucleotides which encodes an amino acid or stop codon.
// It can be validated using the SetCodon function.
type Codon string

// SetCodon creates a Codon from a string input and returns an error
// if the string is not a valid codon.
func SetCodon(dna string) (Codon, error) {

	if len(dna) != 3 {
		return "", fmt.Errorf("codon \"%s\" not valid. must be three nucleotides.", dna)
	}

	if err := ValidDNA(dna); err != nil {
		return "", fmt.Errorf("codon \"%s\" not valid: %s", dna, err.Error())
	}
	return Codon(strings.ToUpper(strings.TrimSpace(dna))), nil
}

// ProteinSequence object is a type of Biosequence
type ProteinSequence struct {
	Nm  string
	Seq string
}

func (prot *ProteinSequence) Sequence() string {
	return prot.Seq
}

func (prot *ProteinSequence) SetSequence(seq string) error {
	prot.Seq = upper(seq)
	return ValidAA(seq)
}

func (prot *ProteinSequence) Name() string {
	return prot.Nm
}

func (prot *ProteinSequence) SetName(name string) {
	prot.Nm = trimString(name)
}

func (prot *ProteinSequence) Append(s string) error {
	err := ValidAA(s)

	if err != nil {
		return fmt.Errorf("invalid characters requested for Append: %s", err.Error())
	}

	prot.Seq = prot.Seq + s
	return nil
}

func (prot *ProteinSequence) Prepend(s string) error {

	err := ValidAA(s)

	if err != nil {
		return fmt.Errorf("invalid characters requested for Prepend: %s", err.Error())
	}

	prot.Seq = s + prot.Seq
	return nil
}

func (seq *ProteinSequence) Blast() (hits []Hit, err error) {
	hits, err = blast.MegaBlastP(seq.Seq)
	return
}

// Estimate molecular weight of protein product
func (seq *ProteinSequence) Molecularweight() (daltons float64) {
	aaarray := strings.Split(seq.Seq, "")
	array := make([]float64, len(aaarray))
	for i := 0; i < len(aaarray); i++ {
		array = append(array, (aa_mw[aaarray[i]] - 18.0))
	}
	sum := 0.0
	for j := 0; j < len(array); j++ {
		sum += array[j]
	}
	daltons = sum
	//kDa = sum / 1000
	return

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

func Prefix(seq string, lengthofprefix int) (prefix string) {
	prefix = seq[:lengthofprefix]
	return prefix
}
func Suffix(seq string, lengthofsuffix int) (suffix string) {
	suffix = seq[(len(seq) - lengthofsuffix):]
	return suffix
}
func Rev(s string) string {
	r := ""

	for i := len(s) - 1; i >= 0; i-- {
		r += string(s[i])
	}

	return r
}
func Comp(s string) string {
	r := ""

	m := map[string]string{
		"A": "T",
		"T": "A",
		"U": "A",
		"C": "G",
		"G": "C",
		"Y": "R",
		"R": "Y",
		"W": "W",
		"S": "S",
		"K": "M",
		"M": "K",
		"D": "H",
		"V": "B",
		"H": "D",
		"B": "V",
		"N": "N",
		"X": "X",
	}

	for _, c := range s {
		if rc, ok := m[string(c)]; ok {
			r += rc
		} else {
			r += string(c)
		}
	}

	return r
}

// Reverse Complement
func RevComp(s string) string {
	s = strings.ToUpper(s)
	return Comp(Rev(s))
}

type DNASeqSet []*DNASequence

func (dss DNASeqSet) AsBioSequences() []BioSequence {
	r := make([]BioSequence, len(dss))

	for i := 0; i < len(dss); i++ {
		r[i] = BioSequence(dss[i])
	}

	return r
}
