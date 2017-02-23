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

// Package for exporting to files
package export

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/AnthaPath"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/enzymes"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/enzymes/lookup"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/sequences"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wutil"
)

const (
	ANTHAPATH bool = true
	LOCAL     bool = false
)

// function to export a standard report of sequence properties to a txt file
func SequenceReport(dir string, seq wtype.BioSequence) (wtype.File, string, error) {

	var anthafile wtype.File
	filename := filepath.Join(anthapath.Path(), fmt.Sprintf("%s_%s.txt", dir, seq.Name()))
	if err := os.MkdirAll(filepath.Dir(filename), 0777); err != nil {
		return anthafile, "", err
	}

	f, err := os.Create(filename)
	if err != nil {
		return anthafile, "", err
	}
	defer f.Close()

	var buf bytes.Buffer

	// GC content
	GC := sequences.GCcontent(seq.Sequence())

	// Find all orfs:
	orfs := sequences.DoublestrandedORFS(seq.Sequence())

	fmt.Fprintln(&buf, ">", dir[2:]+"_"+seq.Name())
	fmt.Fprintln(&buf, seq.Sequence())

	fmt.Fprintln(&buf, "Sequence length:", len(seq.Sequence()))
	fmt.Fprintln(&buf, "Molecular weight:", wutil.RoundInt(sequences.MassDNA(seq.Sequence(), false, true)), "g/mol")
	fmt.Fprintln(&buf, "GC Content:", wutil.RoundInt((GC * 100)), "%")

	fmt.Fprintln(&buf, (len(orfs.TopstrandORFS) + len(orfs.BottomstrandORFS)), "Potential Open reading frames found:")
	for _, strandorf := range orfs.TopstrandORFS {
		fmt.Fprintln(&buf, "Topstrand")
		fmt.Fprintln(&buf, "Position:", strandorf.StartPosition, "..", strandorf.EndPosition)

		fmt.Fprintln(&buf, " DNA Sequence:", strandorf.DNASeq)

		fmt.Fprintln(&buf, "Translated Amino Acid Sequence:", strandorf.ProtSeq)
		fmt.Fprintln(&buf, "Length of Amino acid sequence:", len(strandorf.ProtSeq)-1)
		fmt.Fprintln(&buf, "molecular weight:", sequences.Molecularweight(strandorf), "kDA")
	}
	for _, strandorf := range orfs.BottomstrandORFS {
		fmt.Fprintln(&buf, "Bottom strand")
		fmt.Fprintln(&buf, "Position:", strandorf.StartPosition, "..", strandorf.EndPosition)

		fmt.Fprintln(&buf, " DNA Sequence:", strandorf.DNASeq)

		fmt.Fprintln(&buf, "Translated Amino Acid Sequence:", strandorf.ProtSeq)
		fmt.Fprintln(&buf, "Length of Amino acid sequence:", len(strandorf.ProtSeq)-1)
		fmt.Fprintln(&buf, "molecular weight:", sequences.Molecularweight(strandorf), "kDA")
	}

	_, err = io.Copy(f, &buf)

	allbytes := streamToByte(f)

	anthafile.Name = filename
	anthafile.WriteAll(allbytes)

	return anthafile, filename, err
}

// function to export a sequence to a txt file
func Fasta(dir string, seq wtype.BioSequence) (wtype.File, string, error) {
	var anthafile wtype.File
	filename := filepath.Join(anthapath.Path(), fmt.Sprintf("%s_%s.fasta", dir, seq.Name()))
	if err := os.MkdirAll(filepath.Dir(filename), 0777); err != nil {
		return anthafile, "", err
	}

	f, err := os.Create(filename)
	if err != nil {
		return anthafile, "", err
	}
	defer f.Close()

	_, err = fmt.Fprintf(f, ">%s\n%s\n", seq.Name(), seq.Sequence())

	allbytes := streamToByte(f)

	anthafile.Name = filename
	anthafile.WriteAll(allbytes)

	return anthafile, filename, err
}

// function to export multiple sequences in fasta format into a specified directory
// specify whether to save locally or to the anthapath in a specified sub directory dir.
func FastaSerial(makeinanthapath bool, dir string, seqs []wtype.DNASequence) (wtype.File, string, error) {

	var anthafile wtype.File
	var filename string
	if makeinanthapath {
		filename = filepath.Join(anthapath.Path(), fmt.Sprintf("%s.fasta", dir))
	} else {
		filename = filepath.Join(fmt.Sprintf("%s.fasta", dir))
	}
	if err := os.MkdirAll(filepath.Dir(filename), 0777); err != nil {
		return anthafile, "", err
	}

	f, err := os.Create(filename)
	if err != nil {
		return anthafile, "", err
	}
	defer f.Close()

	for _, seq := range seqs {
		if _, err := fmt.Fprintf(f, ">%s\n%s\n", seq.Name(), seq.Sequence()); err != nil {
			return anthafile, "", err
		}
	}

	allbytes := streamToByte(f)

	anthafile.Name = filename
	anthafile.WriteAll(allbytes)

	return anthafile, filename, nil
}

// Simultaneously export multiple Fasta files and summary files for a series of assembly products
func FastaAndSeqReports(assemblyparameters enzymes.Assemblyparameters) (fastafiles []wtype.File, summaryfiles []wtype.File, err error) {

	enzymename := strings.ToUpper(assemblyparameters.Enzymename)

	// should change this to rebase lookup; what happens if this fails?
	//enzyme := TypeIIsEnzymeproperties[enzymename]
	enzyme, err := lookup.TypeIIsLookup(enzymename)
	if err != nil {
		return fastafiles, summaryfiles, err
	}
	//assemble (note that sapIenz is found in package enzymes)
	_, plasmidproductsfromXprimaryseq, err := enzymes.JoinXnumberofparts(assemblyparameters.Vector, assemblyparameters.Partsinorder, enzyme)

	if err != nil {
		return fastafiles, summaryfiles, err
	}

	for _, assemblyproduct := range plasmidproductsfromXprimaryseq {
		filename := filepath.Join(anthapath.Path(), assemblyparameters.Constructname)
		if summary, _, err := SequenceReport(filename, &assemblyproduct); err != nil {
			return fastafiles, summaryfiles, err
		} else {
			summaryfiles = append(summaryfiles, summary)

		}
		if fasta, _, err := Fasta(filename, &assemblyproduct); err != nil {
			return fastafiles, summaryfiles, err
		} else {
			fastafiles = append(fastafiles, fasta)
		}
	}

	return fastafiles, summaryfiles, nil
}

// Simultaneously export a single Fasta file containing the assembled sequences for a series of assembly products
func FastaSerialfromMultipleAssemblies(dirname string, multipleassemblyparameters []enzymes.Assemblyparameters) (wtype.File, string, error) {
	var anthafile wtype.File
	seqs := make([]wtype.DNASequence, 0)

	for _, assemblyparameters := range multipleassemblyparameters {

		enzymename := strings.ToUpper(assemblyparameters.Enzymename)

		// should change this to rebase lookup; what happens if this fails?
		enzyme, err := lookup.TypeIIsLookup(enzymename)
		if err != nil {
			return anthafile, "", err
		}
		//assemble
		_, plasmidproductsfromXprimaryseq, err := enzymes.JoinXnumberofparts(assemblyparameters.Vector, assemblyparameters.Partsinorder, enzyme)
		if err != nil {
			return anthafile, "", err
		}

		for _, assemblyproduct := range plasmidproductsfromXprimaryseq {
			seqs = append(seqs, assemblyproduct)
		}

	}

	return FastaSerial(ANTHAPATH, dirname, seqs)
}

// export data in the format of an array of strings to a file
func TextFile(filename string, data []string) (wtype.File, error) {

	var anthafile wtype.File

	f, err := os.Create(filename)
	if err != nil {
		return anthafile, err
	}
	defer f.Close()

	for _, str := range data {

		if _, err := fmt.Fprintln(f, str); err != nil {
			return anthafile, err
		}
	}
	alldata := stringsToBytes(data)
	anthafile.Name = filename

	anthafile.WriteAll(alldata)

	return anthafile, nil
}

// Export any data as a json object in  a file
func JSON(data interface{}, filename string) (anthafile wtype.File, err error) {
	bytes, err := json.Marshal(data)

	if err != nil {
		return anthafile, err
	}

	ioutil.WriteFile(filename, bytes, 0644)

	anthafile.Name = filename
	anthafile.WriteAll(bytes)
	return anthafile, nil
}

// Export a 2D array of string data as a csv file
func CSV(records [][]string, filename string) (wtype.File, error) {
	var anthafile wtype.File
	var buf bytes.Buffer

	/// use the buffer to create a csv writer
	w := csv.NewWriter(&buf)

	// write all records to the buffer
	w.WriteAll(records) // calls Flush internally

	if err := w.Error(); err != nil {
		return anthafile, fmt.Errorf("error writing csv: %s", err.Error())
	}

	//This code shows how to create an antha File from this buffer which can be downloaded through the UI:
	var SequencingResultsFile wtype.File

	SequencingResultsFile.Name = filename

	SequencingResultsFile.WriteAll(buf.Bytes())

	///// to write this to a file on the command line this is what we'd do (or something similar)

	// also create a file on os
	file, _ := os.Create(filename)
	defer file.Close()

	// this time we'll use the file to create the writer instead of a buffer (anything which fulfils the writer interface can be used here ... checkout golang io.Writer and io.Reader)
	fw := csv.NewWriter(file)

	// same as before ...
	fw.WriteAll(records)
	return anthafile, nil
}

// export bytes into a file
func Binary(data []byte, filename string) (wtype.File, error) {
	var anthafile wtype.File
	if len(data) == 0 {
		return anthafile, fmt.Errorf("No data to export into file")
	}
	anthafile.Name = filename
	anthafile.WriteAll(data)
	return anthafile, nil
}

// export a stream into a file
func Stream(stream io.Reader, filename string) (wtype.File, error) {
	return Binary(streamToByte(stream), filename)
}

func stringsToBytes(data []string) []byte {
	var alldata []byte

	for _, str := range data {
		bts := []byte(str)
		for i := range bts {
			alldata = append(alldata, bts[i])
		}
	}
	return alldata
}

func streamToByte(stream io.Reader) []byte {
	buf := new(bytes.Buffer)
	buf.ReadFrom(stream)
	return buf.Bytes()
}
