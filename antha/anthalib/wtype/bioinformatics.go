// anthalib//wtype/bioinformatics.go: Part of the Antha language
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

package wtype

import (
	"fmt"
	"os"
)

type AlignedBioSequence struct {
	Query   string
	Subject string
	Score   float64
}

type SequenceDatabase struct {
	Name      string
	Filename  string
	Type      string
	Sequences []BioSequence
}

// struct for holding results of a blast search
type BlastResults struct {
	Program       string
	DBname        string
	DBSizeSeqs    int
	DBSizeLetters int
	Query         string
	Hits          []BlastHit
}

// constructor, makes an empty BlastResults structure
func NewBlastResults() BlastResults {
	return BlastResults{"", "", -1, -1, "", make([]BlastHit, 0, 1)}
}

// struct for holding a particular hit
type BlastHit struct {
	Name       string
	Score      float64
	Eval       float64
	Alignments []AlignedSequence
}

// constructor, makes an empty BlastHit structure
func NewBlastHit() BlastHit {
	return BlastHit{"", 0.0, 0.0, make([]AlignedSequence, 0, 2)}
}

// struct for holding an aligned sequence
type AlignedSequence struct {
	Qstrand string
	Sstrand string
	Qstart  int
	Qend    int
	Sstart  int
	Send    int
	Qseq    string
	Sseq    string
	ID      float64
}

// constructor for an AlignedSequence object, makes an empty structure
func NewAlignedSequence() AlignedSequence {
	return AlignedSequence{"", "", -1, -1, -1, -1, "", "", 0.0}
}

// struct for holding BLAST parameters

type BLASTSearchParameters struct {
	Evalthreshold float64
	Matrix        string
	Filter        bool
	Open          int
	Extend        int
	DBSeqs        int
	DBAlns        int
	GCode         int
}

func DefaultBLASTSearchParameters() BLASTSearchParameters {
	return BLASTSearchParameters{10.0, "BLOSUM62", true, -1, -1, 250, 250, 1}
}

// creates a fasta file containing the sequence
func Makeseq(dir string, seq BioSequence) string {
	filename := dir + "/" + seq.Name() + ".fasta"
	f, e := os.Create(filename)
	if e != nil {
		panic(e)
	}
	defer f.Close() //nolint

	if _, err := fmt.Fprintf(f, ">%s\n%s\n", seq.Name(), seq.Sequence()); err != nil {
		panic(err)
	}

	return filename
}
func (res BlastResults) QueryCentredAlignment() []AlignedBioSequence {
	ret := make([]AlignedBioSequence, 0, len(res.Hits))

	for _, h := range res.Hits {
		qry, sbj := h.Alignments[0].CentreToQuery(res.Query)
		ret = append(ret, AlignedBioSequence{Query: qry, Subject: sbj, Score: h.Score})
	}

	return ret
}

type SimpleAlignment []AlignedBioSequence

// no guarantees... it's just some strings
type ReallySimpleAlignment []string

func (aln ReallySimpleAlignment) Column(i int) string {
	if i < 0 || i >= len(aln[0]) {
		panic(fmt.Sprintf("Error: Cannot take column %d in alignment of length %d", i, len(aln[0])))
	}

	r := ""

	for s := 0; s < len(aln); s++ {
		c := string(aln[s][i])

		if c != "-" {
			r += c
		}
	}

	return r
}

// find column of length j slices at pos i
func (aln ReallySimpleAlignment) MultiColumn(i, j int) []string {
	if i < 0 || i >= len(aln[0])-j+1 {
		panic(fmt.Sprintf("Error: Cannot take column %d of size %d in alignment of length %d", i, j, len(aln[0])))
	}

	r := make([]string, 0, len(aln))

	for s := 0; s < len(aln); s++ {
		c := string(aln[s][i : i+j])

		r = append(r, c)
	}

	return r
}

func (aln ReallySimpleAlignment) TrimToFrame(frame int) ReallySimpleAlignment {
	trim1 := aln
	if frame != 0 {
		trim1 = ReallySimpleAlignment(aln.MultiColumn(3-frame, len(aln[0])+frame-3))
	}
	endFrame := len(trim1[0]) % 3

	return trim1.MultiColumn(0, len(trim1[0])-endFrame)
}

func (aln SimpleAlignment) Column(i int) string {
	if i < 0 || i >= len(aln[0].Subject) {
		panic(fmt.Sprintf("Error: Cannot take column %d in alignment of length %d", i, len(aln[0].Subject)))
	}

	r := ""

	for s := 0; s < len(aln); s++ {
		c := string(aln[s].Subject[i])

		if c != "-" {
			r += c
		}
	}

	return r
}

// CentreToQuery trims aligned (subject) sequences to only those
// residues/nucleotides aligned to those in the query. This removes inserts
// (with respect to the query) from aligned sequences. Gaps within the interior
// of the aligned sequence are already represented by '-' characters, however
// gaps at the ends of the aligned sequence must be added here as well. The
// resulting aligned sequences will have the same length as the query.
func (aln AlignedSequence) CentreToQuery(q string) (string, string) {

	s := "" // query
	r := "" // aligned

	// Add any gaps to the start of the aligned sequence.
	for i := 1; i < aln.Qstart; i++ {
		r += "-"
		s += string(q[i-1])
	}

	// Remove any inserts from the aligned sequence.
	for i := 0; i < len(aln.Qseq); i++ {
		if aln.Qseq[i] != '-' {
			r += string(aln.Sseq[i])
			s += string(aln.Qseq[i])
		}
	}

	// Add any gaps to the end of the aligned sequence.
	for i := aln.Qend + 1; i <= len(q); i++ {
		r += "-"
		s += string(q[i-1])
	}

	// Return query, aligned.
	return s, r
}
