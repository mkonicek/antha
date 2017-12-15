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
	"math/rand"
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
}

type RestrictionEnzyme struct {
	// other fields required but for now the main things are...
	RecognitionSequence               string
	EndLength                         int
	Name                              string
	Prototype                         string
	Topstrand3primedistancefromend    int
	Bottomstrand5primedistancefromend int
	MethylationSite                   string   //"attr, <4>"
	CommercialSource                  []string //string "attr, <5>"
	References                        []int
	Class                             string
}

type TypeIIs struct {
	RestrictionEnzyme
	Name                              string
	Isoschizomers                     []string
	Topstrand3primedistancefromend    int
	Bottomstrand5primedistancefromend int
}

func ToTypeIIs(typeIIenzyme RestrictionEnzyme) (typeIIsenz TypeIIs, err error) {
	if typeIIenzyme.Class == "TypeII" {
		err = fmt.Errorf("You can't do this, enzyme is not a type IIs")
	}
	if typeIIenzyme.Class == "TypeIIs" {

		var isoschizomers = make([]string, 0)
		/*for _, lookup := range ...
		add code to lookup isoschizers from rebase
		*/
		typeIIsenz = TypeIIs{typeIIenzyme, typeIIenzyme.Name, isoschizomers, typeIIenzyme.Topstrand3primedistancefromend, typeIIenzyme.Bottomstrand5primedistancefromend}

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

	d, _ := json.Marshal(seq)
	json.Unmarshal(d, &ret)

	return ret
}

func MakeDNASequence(name string, seqstring string, properties []string) (seq DNASequence, err error) {
	seq.Nm = name
	seq.Seq = upper(seqstring)
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
		/*
		   // deal with overhangs separately
		   if strings.Contains(property,"5'") {
		   	seq.Overhang5prime.End = 5
		   	seq.Overhang5prime.Type =
		   }
		*/
	}
	return
}
func MakeLinearDNASequence(name string, seqstring string) (seq DNASequence) {
	seq.Nm = name
	seq.Seq = upper(seqstring)

	return
}
func MakePlasmidDNASequence(name string, seqstring string) (seq DNASequence) {
	seq.Nm = name
	seq.Seq = upper(seqstring)
	seq.Plasmid = true
	return
}
func MakeSingleStrandedDNASequence(name string, seqstring string) (seq DNASequence) {
	seq.Nm = name
	seq.Seq = upper(seqstring)
	seq.Singlestranded = true
	return
}

func MakeOverhang(sequence DNASequence, end int, toporbottom int, length int, phosphorylated bool) (overhang Overhang, err error) {

	if sequence.Singlestranded {
		err = fmt.Errorf("Can't have overhang on single stranded dna")
		return
	}
	if sequence.Plasmid {
		err = fmt.Errorf("Can't have overhang on Plasmid(circular) dna")
		return
	}
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
	if toporbottom == 0 && length == 0 {
		overhang.Type = 1
		return
	}
	if toporbottom == 0 && length != 0 {
		err = fmt.Errorf("If length of overhang is not 0, toporbottom must be 0")
		return
	}
	if toporbottom != 0 && length == 0 {
		err = fmt.Errorf("If length of overhang is not 0, toporbottom must be 0")
		return
	}
	if toporbottom > 2 {
		err = fmt.Errorf("invalid entry for toporbottom: NEITHER = 0, TOP    = 1, BOTTOM = 2")
		return
	}
	if toporbottom == 1 {
		overhang.Type = 2
		overhang.Sequence = Prefix(sequence.Seq, length)
	}
	if toporbottom == 2 {
		overhang.Type = -1
		overhang.Sequence = Suffix(RevComp(sequence.Seq), length)
	}
	overhang.Phosphorylation = phosphorylated
	return
}

func Phosphorylate(dnaseq DNASequence) (phosphorylateddna DNASequence, err error) {
	if dnaseq.Plasmid == true {
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

const (
	FALSE     = 0
	BLUNT     = 1
	OVERHANG  = 2
	UNDERHANG = -1
)

const (
	NEITHER = 0
	TOP     = 1
	BOTTOM  = 2
)

type Overhang struct {
	//Strand          int // i.e. 1 or 2 (top or bottom
	End             int    `json:"end"`  // i.e. 5 or 3 or 0
	Type            int    `json:"type"` //as contants above
	Length          int    `json:"length"`
	Sequence        string `json:"sequence"`
	Phosphorylation bool   `json:"phosphorylation"`
}

func (oh Overhang) OverHangAt5PrimeEnd() (sequence string) {
	if oh.End == 5 {
		if oh.Type == OVERHANG {
			return oh.Sequence
		}

	}
	return ""
}

func (oh Overhang) OverHangAt3PrimeEnd() (sequence string) {
	if oh.End == 3 {
		if oh.Type == OVERHANG {
			return oh.Sequence
		}

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

func validDNA(seq string) error {
	validNucleotides := "ACTGNXBHVDMKSWRYU"

	return valid(seq, validNucleotides)
}

func validRNA(seq string) error {
	validRNA := "ACGU"

	return valid(seq, validRNA)
}

func validAA(seq string) error {

	var aminoAcids []string

	for key := range aa_mw {
		aminoAcids = append(aminoAcids, key)
	}

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
// after conversion to upper case the upper case sequence will be added with no error returned.
// If invalid characters are found an error returned listing all invalid characters and their positions in human friendly form. i.e. the first position is 1 and not 0.
func (dna *DNASequence) SetSequence(seq string) error {

	dna.Seq = upper(seq)

	return validDNA(seq)
}

// Append appends the existing dna sequence with the upper case of the string added
func (dna *DNASequence) Append(s string) error {

	err := validDNA(s)

	if err != nil {
		return fmt.Errorf("invalid characters requested for Append: %s", err.Error())
	}

	dna.Seq = dna.Seq + upper(s)

	return nil
}

// Preprend adds the requested sequence to the beginning of the existing sequence.
func (dna *DNASequence) Prepend(s string) error {

	err := validDNA(s)

	if err != nil {
		return fmt.Errorf("invalid characters requested for Prepend: %s", err.Error())
	}

	dna.Seq = upper(s) + dna.Seq

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
	if phosphate5prime == true {
		mw = mw + 79.0 // extra for phosphate left at 5' end following digestion, not relevant for primer extension
	}
	if phosphate3prime == true {
		mw = mw + 79.0 // extra for phosphate left at 3' end following digestion, not relevant for primer extension
	}
	if singlestranded != true {
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
	return validRNA(seq)
}

func (rna *RNASequence) Name() string {
	return rna.Nm
}

func (rna *RNASequence) SetName(name string) {
	rna.Nm = trimString(name)
}

func (rna *RNASequence) Append(s string) error {
	err := validRNA(s)

	if err != nil {
		return fmt.Errorf("invalid characters requested for Append: %s", err.Error())
	}

	rna.Seq = rna.Seq + upper(s)
	return nil
}

func (rna *RNASequence) Prepend(s string) error {

	err := validRNA(s)

	if err != nil {
		return fmt.Errorf("invalid characters requested for Prepend: %s", err.Error())
	}

	rna.Seq = upper(s) + rna.Seq
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
	return validAA(seq)
}

func (prot *ProteinSequence) Name() string {
	return prot.Nm
}

func (prot *ProteinSequence) SetName(name string) {
	prot.Nm = trimString(name)
}

func (prot *ProteinSequence) Append(s string) error {
	err := validAA(s)

	if err != nil {
		return fmt.Errorf("invalid characters requested for Append: %s", err.Error())
	}

	prot.Seq = prot.Seq + upper(s)
	return nil
}

func (prot *ProteinSequence) Prepend(s string) error {

	err := validAA(s)

	if err != nil {
		return fmt.Errorf("invalid characters requested for Prepend: %s", err.Error())
	}

	prot.Seq = upper(s) + prot.Seq
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

func random_dna_seq(leng int) string {
	s := ""
	for i := 0; i < leng; i++ {
		s += random_char("ACTG")
	}
	return s
}

func random_char(chars string) string {
	return string(chars[rand.Intn(len(chars))])
}

func makeABunchaRandomSeqs(n_seq_sets, seqs_per_set, min_len, len_var int) [][]DNASequence {
	var seqs [][]DNASequence
	var features []Feature

	seqs = make([][]DNASequence, n_seq_sets)

	for i := 0; i < n_seq_sets; i++ {
		seqs[i] = make([]DNASequence, seqs_per_set)
		for j := 0; j < seqs_per_set; j++ {
			seqs[i][j] = DNASequence{fmt.Sprintf("SEQ%04d", i*seqs_per_set+j+1), random_dna_seq(rand.Intn(len_var) + min_len), false, false, Overhang{0, 0, 0, "", false}, Overhang{0, 0, 0, "", false}, "", features}
		}
	}
	return seqs
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
		r += m[string(c)]
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
