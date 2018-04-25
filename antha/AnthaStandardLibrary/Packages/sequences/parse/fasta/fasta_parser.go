// antha/AnthaStandardLibrary/Packages/Parser/fasta_parser.go: Part of the Antha language
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

// package fasta converts DNA sequence files in FASTA format into a set of DNA sequences.
package fasta

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"errors"
	"io/ioutil"
	"os"
	"strings"

	"github.com/antha-lang/antha/antha/anthalib/wtype"
)

// Fasta represents the intermediate structure containing contents of Fasta file as strings
type Fasta struct {
	Id   string
	Desc string
	Seq  string
}

// This will retrieve seq from FASTA file using the sequence ID found in the file
func RetrieveSeqFromFASTA(id string, fastaFile wtype.File) (seq wtype.DNASequence, err error) {

	var nofeatures []wtype.Feature

	allparts, err := fastaFile.ReadAll()
	if err != nil {
		return
	}

	// then retrieve the particular record
	for _, record := range fastaParse(allparts) {
		if strings.Contains(record.Id, id) {
			seq = wtype.DNASequence{
				Nm:             record.Id,
				Seq:            record.Seq,
				Plasmid:        false,
				Singlestranded: false,
				Overhang5prime: wtype.Overhang{End: 0, Type: 0, Seq: "", Phosphorylation: false},
				Overhang3prime: wtype.Overhang{End: 0, Type: 0, Seq: "", Phosphorylation: false},
				Methylation:    "",
				Features:       nofeatures,
			}
			return
		}
	}

	seq = wtype.DNASequence{} // blank seq
	if seq.Seq == "" {
		err = errors.New("Record not found in file")
		return
	}
	return
}

// This will retrieve seq from FASTA file
func FASTAtoLinearDNASeqs(fastaFile wtype.File) (seqs []wtype.DNASequence, err error) {

	var nofeatures []wtype.Feature

	seqs = make([]wtype.DNASequence, 0)

	var seq wtype.DNASequence

	allparts, err := fastaFile.ReadAll()
	if err != nil {
		return
	}

	// then retrieve the particular record
	for _, record := range fastaParse(allparts) {
		seq = wtype.DNASequence{
			Nm:             record.Id,
			Seq:            record.Seq,
			Plasmid:        false,
			Singlestranded: false,
			Overhang5prime: wtype.Overhang{End: 0, Type: 0, Seq: "", Phosphorylation: false},
			Overhang3prime: wtype.Overhang{End: 0, Type: 0, Seq: "", Phosphorylation: false},
			Methylation:    "",
			Features:       nofeatures,
		}
		seqs = append(seqs, seq)
	}
	return

}

// This will retrieve sequences from a FASTA file of type wtype.File and set all sequences to plasmids
func FASTAtoPlasmidDNASeqs(file wtype.File) (seqs []wtype.DNASequence, err error) {

	var nofeatures []wtype.Feature

	seqs = make([]wtype.DNASequence, 0)

	var seq wtype.DNASequence

	allparts, err := file.ReadAll()
	if err != nil {
		return
	}

	// then retrieve the particular record
	for _, record := range fastaParse(allparts) {
		seq = wtype.DNASequence{
			Nm:             record.Id,
			Seq:            record.Seq,
			Plasmid:        false,
			Singlestranded: false,
			Overhang5prime: wtype.Overhang{End: 0, Type: 0, Seq: "", Phosphorylation: false},
			Overhang3prime: wtype.Overhang{End: 0, Type: 0, Seq: "", Phosphorylation: false},
			Methylation:    "",
			Features:       nofeatures,
		}
		seqs = append(seqs, seq)

	}
	return

}

// Convert a sequence file in Fasta format to an array of DNASequence.
// If the header does not contain the key words PLASMID, CIRCULAR or VECTOR the sequence will be assumed to be linear.
func FastaToDNASequences(sequenceFile wtype.File) (seqs []wtype.DNASequence, err error) {
	data, err := sequenceFile.ReadAll()
	if err != nil {
		return
	}
	return FastaContentstoDNASequences(data)
}

// Convert the contents of a sequence file in Fasta format to an array of DNASequence.
// If the header does not contain the key words PLASMID, CIRCULAR or VECTOR the sequence will be assumed to be linear.
func FastaContentstoDNASequences(data []byte) (seqs []wtype.DNASequence, err error) {

	for _, record := range fastaParse(data) {
		plasmidstatus := ""

		if strings.Contains(strings.ToUpper(record.Desc), "PLASMID") || strings.Contains(strings.ToUpper(record.Desc), "CIRCULAR") || strings.Contains(strings.ToUpper(record.Desc), "VECTOR") {
			plasmidstatus = "PLASMID"
		}

		seq, err := wtype.MakeDNASequence(record.Id, record.Seq, []string{plasmidstatus})
		if err != nil {
			return seqs, err
		}
		seqs = append(seqs, seq)
	}

	return
}

func build_fasta(header string, seq bytes.Buffer) (Record Fasta) {
	fields := strings.SplitN(header, " ", 2)

	var record Fasta

	if len(fields) > 1 {
		record.Id = fields[0]
		record.Desc = "`" + fields[1] + "`"
	} else {
		record.Id = fields[0]
		record.Desc = ""
	}

	record.Seq = seq.String()

	Record = record

	return Record
}

func fastaParse(fastaFh []byte) []Fasta {
	var outputs []Fasta
	buffer := bytes.NewBuffer(fastaFh)

	scanner := bufio.NewScanner(buffer)
	header := ""
	var seq bytes.Buffer

	// Loop over the letters in inputString
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		line = strings.Replace(line, "\t", "", -1)
		sublines := strings.Split(line, "\r")
		for _, line := range sublines {

			if len(line) == 0 {
				continue
			}

			if line[0] == '>' {
				// If we stored a previous identifier, get the DNA string and map to the
				// identifier and clear the string
				if header != "" {
					// outputChannel <- build_fasta(header, seq.String())
					outputs = append(outputs, build_fasta(header, seq))
					seq.Reset()
				}

				// Standard FASTA identifiers look like: ">id desc"
				header = line[1:]

			} else {
				// Append here since multi-line DNA strings are possible
				seq.WriteString(line)
			}
		}
	}

	outputs = append(outputs, build_fasta(header, seq))

	return outputs
}

func Fastatocsv(inputfilename wtype.File, outputfileprefix string) (csvfile *os.File, err error) {
	fastaFh, err := inputfilename.ReadAll()
	if err != nil {
		return
	}

	csvfile, err = ioutil.TempFile(os.TempDir(), "csv")
	if err != nil {
		return
	}

	records := make([][]string, 0)
	seq := []string{"#Name", "Sequence", "Plasmid?", "Seq Type", "Class"}
	records = append(records, seq)
	for _, record := range fastaParse(fastaFh) {
		plasmidstatus := "FALSE"
		seqtype := "DNA"
		class := "not specified"
		if strings.Contains(record.Id, "Plasmid") || strings.Contains(record.Id, "Circular") || strings.Contains(record.Id, "Vector") {
			plasmidstatus = "TRUE"
		}
		if strings.Contains(record.Desc, "Amino acid") || strings.Contains(record.Id, "aa") {
			seqtype = "AA"
		}

		if strings.Contains(record.Desc, "Class:") {
			uptoclass := strings.Index(record.Desc, "Class:")
			prefix := uptoclass + len("class:")
			class = record.Desc[prefix:]
		}
		seq = []string{record.Id, record.Seq, plasmidstatus, seqtype, class}
		records = append(records, seq)
	}

	writer := csv.NewWriter(csvfile)
	for _, record := range records {
		err = writer.Write(record)
		if err != nil {
			return
		}
	}

	writer.Flush()
	return
}
