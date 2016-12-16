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

	fmt.Fprintf(f, ">%s\n%s\n", seq.Name(), seq.Sequence())

	f.Close()

	return filename
}
