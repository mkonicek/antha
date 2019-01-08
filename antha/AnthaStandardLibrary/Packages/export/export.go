// antha/AnthaStandardLibrary/Packages/enzymes/exporttofile.go: Part of the Antha language
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

// Package export provides functions for exporting common file formats into the Antha File type.
package export

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"strings"

	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/enzymes"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/enzymes/lookup"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/sequences"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wutil"
)

// SequenceReport exports a standard report of sequence properties to a txt file.
func SequenceReport(seq wtype.BioSequence) ([]byte, error) {
	// GC content
	GC := sequences.GCcontent(seq.Sequence())

	// Find all orfs:
	orfs := sequences.DoublestrandedORFS(seq.Sequence())

	lines := []string{
		fmt.Sprintln(">", seq.Name()),
		fmt.Sprintln(seq.Sequence()),
		fmt.Sprintln("Sequence length:", len(seq.Sequence())),
		fmt.Sprintln("Molecular weight:", wutil.RoundInt(sequences.MassDNA(seq.Sequence(), false, true)), "g/mol"),
		fmt.Sprintln("GC Content:", wutil.RoundInt((GC * 100)), "%"),
		fmt.Sprintln((len(orfs.TopstrandORFS) + len(orfs.BottomstrandORFS)), "Potential Open reading frames found:"),
	}

	for _, strandorf := range orfs.TopstrandORFS {
		lines = append(lines,
			fmt.Sprintln("Topstrand"),
			fmt.Sprintln("Position:", strandorf.StartPosition, "..", strandorf.EndPosition),
			fmt.Sprintln(" DNA Sequence:", strandorf.DNASeq),
			fmt.Sprintln("Translated Amino Acid Sequence:", strandorf.ProtSeq),
			fmt.Sprintln("Length of Amino acid sequence:", len(strandorf.ProtSeq)-1),
			fmt.Sprintln("molecular weight:", sequences.Molecularweight(strandorf), "kDA"),
		)

	}
	for _, strandorf := range orfs.BottomstrandORFS {
		lines = append(lines,
			fmt.Sprintln("Bottom strand"),
			fmt.Sprintln("Position:", strandorf.StartPosition, "..", strandorf.EndPosition),
			fmt.Sprintln(" DNA Sequence:", strandorf.DNASeq),
			fmt.Sprintln("Translated Amino Acid Sequence:", strandorf.ProtSeq),
			fmt.Sprintln("Length of Amino acid sequence:", len(strandorf.ProtSeq)-1),
			fmt.Sprintln("molecular weight:", sequences.Molecularweight(strandorf), "kDA"),
		)
	}

	var buf bytes.Buffer

	if _, err := fmt.Fprintf(&buf, strings.Join(lines, "")); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// Fasta exports a sequence to a txt file in Fasta format.
func Fasta(seq wtype.BioSequence) ([]byte, error) {
	var buf bytes.Buffer
	if _, err := fmt.Fprintf(&buf, ">%s\n%s\n", seq.Name(), seq.Sequence()); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// FastaSerial exports multiple sequences in fasta format into a specified txt file.
// The makeinanthapath argument specifies whether a copy of the file should be saved locally or to the anthapath in a specified sub directory directory.
func FastaSerial(seqs []wtype.DNASequence) ([]byte, error) {
	var buf bytes.Buffer

	for _, seq := range seqs {
		_, err := fmt.Fprintf(&buf, ">%s\n%s\n", seq.Name(), seq.Sequence())
		if err != nil {
			return nil, err
		}
	}

	allbytes := buf.Bytes()

	if len(allbytes) == 0 {
		return nil, fmt.Errorf("empty Fasta file created for seqs")
	}

	return allbytes, nil
}

// FastaAndSeqReports simultaneously exports multiple Fasta files and summary files for a TypeIIs assembly design.
func FastaAndSeqReports(assemblyparameters enzymes.Assemblyparameters) (fastafiles [][]byte, summaryfiles [][]byte, err error) {

	enzymename := strings.ToUpper(assemblyparameters.Enzymename)

	// should change this to rebase lookup; what happens if this fails?
	//enzyme := TypeIIsEnzymeproperties[enzymename]
	enzyme, err := lookup.TypeIIs(enzymename)
	if err != nil {
		return fastafiles, summaryfiles, err
	}
	//assemble (note that sapIenz is found in package enzymes)
	_, plasmidproductsfromXprimaryseq, _, err := enzymes.JoinXNumberOfParts(assemblyparameters.Vector, assemblyparameters.Partsinorder, enzyme)

	if err != nil {
		return fastafiles, summaryfiles, err
	}

	for _, assemblyproduct := range plasmidproductsfromXprimaryseq {
		summary, err := SequenceReport(&assemblyproduct)
		if err != nil {
			return nil, nil, err
		}
		summaryfiles = append(summaryfiles, summary)

		fasta, err := Fasta(&assemblyproduct)
		if err != nil {
			return nil, nil, err
		}
		fastafiles = append(fastafiles, fasta)
	}

	return fastafiles, summaryfiles, nil
}

// FastaSerialfromMultipleAssemblies simultaneously export a single Fasta file containing the assembled sequences for a series of TypeIIs assembly designs.
func FastaSerialfromMultipleAssemblies(multipleassemblyparameters []enzymes.Assemblyparameters) ([]byte, error) {
	seqs := make([]wtype.DNASequence, 0)

	for _, assemblyparameters := range multipleassemblyparameters {
		enzymename := strings.ToUpper(assemblyparameters.Enzymename)

		// should change this to rebase lookup; what happens if this fails?
		enzyme, err := lookup.TypeIIs(enzymename)
		if err != nil {
			return nil, err
		}
		//assemble
		_, plasmidproductsfromXprimaryseq, _, err := enzymes.JoinXNumberOfParts(assemblyparameters.Vector, assemblyparameters.Partsinorder, enzyme)
		if err != nil {
			return nil, err
		}

		seqs = append(seqs, plasmidproductsfromXprimaryseq...)
	}

	return FastaSerial(seqs)
}

// CSV exports a matrix of string data as a csv file.
func CSV(records [][]string) ([]byte, error) {
	var buf bytes.Buffer

	/// use the buffer to create a csv writer
	w := csv.NewWriter(&buf)

	// write all records to the buffer
	if err := w.WriteAll(records); err != nil {
		return nil, fmt.Errorf("error writing csv: %s", err.Error())
	}

	if err := w.Error(); err != nil {
		return nil, fmt.Errorf("error writing csv: %s", err.Error())
	}
	return buf.Bytes(), nil
}
