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
	"text/template"
	"time"

	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/enzymes"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/enzymes/lookup"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/sequences"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wutil"
	"github.com/antha-lang/antha/laboratory"
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

// GenbankSerial exports multiple sequences into a multi-record Genbank format file
func GenbankSerial(lab *laboratory.Laboratory, filename string, seqs []wtype.DNASequence) (*wtype.File, error) {

	// Template for multi-record Genbank file
	// https://www.ncbi.nlm.nih.gov/Sitemap/samplerecord.html
	// http://www.insdc.org/documents/feature_table.html
	tmplStr := `{{ range $i, $s := . -}}
LOCUS       {{ $s.Nm }}               {{ length $s }} bp ds-DNA     {{ if $s.Plasmid }}circular{{ else }}linear{{ end }} SYN {{ date }}
DEFINITION  Exported from Antha OS
ACCESSION   
VERSION     
KEYWORDS    
SOURCE      synthetic DNA construct
FEATURES             Location/Qualifiers
     source          1..{{ length . }}
                     /organism="synthetic DNA construct"
                     /mol_type="other DNA"
{{- range $j, $f := $s.Features }}
     {{  keyf $f  }} {{ location $s $f }}
                     /label="{{ $f.Name }}"
{{- end }}
ORIGIN
{{- range $j, $line := origin $s }}
{{ $line }}
{{- end }}
//
{{ end }}`

	tmpl, err := template.New("genbank").Funcs(template.FuncMap{
		"date": func() string {
			t := time.Now()
			return t.Format("2-JAN-2006")
		},
		"length": func(seq wtype.DNASequence) int {
			return len(seq.Sequence())
		},
		// Formatted feature key
		// http://www.insdc.org/documents/feature_table.html#3.1
		"keyf": func(feat wtype.Feature) string {
			return fmt.Sprintf("%-15s", strings.TrimSpace(feat.Class))
		},
		"location": func(seq wtype.DNASequence, feat wtype.Feature) string {
			var ret string
			if seq.Plasmid {
				if feat.Start() > feat.End() {
					if feat.Reverse {
						ret = fmt.Sprintf("complement(%d..%d)", feat.End(), feat.Start())
					} else {
						ret = fmt.Sprintf("join(%d..%d,%d..%d)", feat.Start(), len(seq.Sequence()), 1, feat.End())
					}
				} else {
					if feat.Reverse {
						ret = fmt.Sprintf("complement(join(%d..%d,%d..%d))", feat.End(), len(seq.Sequence()), 1, feat.Start())
					} else {
						ret = fmt.Sprintf("%d..%d", feat.Start(), feat.End())
					}
				}
			} else {
				if feat.Start() > feat.End() {
					if feat.Reverse {
						ret = fmt.Sprintf("complement(%d..%d)", feat.End(), feat.Start())
					} else {
						ret = "ERROR"
					}
				} else {
					if feat.Reverse {
						ret = "ERROR"
					} else {
						ret = fmt.Sprintf("%d..%d", feat.Start(), feat.End())
					}
				}
			}
			return ret
		},
		"origin": func(seq wtype.DNASequence) []string {
			bases := seq.Sequence()
			lines := []string{}
			// Format per https://www.ncbi.nlm.nih.gov/Sitemap/samplerecord.html
			const BASES_PER_LINE = 60
			const BASES_PER_BLOC = 10
			for pos := 0; pos < len(bases); pos += BASES_PER_LINE {
				frags := []string{}
				frags = append(frags, fmt.Sprintf("%9d ", pos+1))
				for inner := pos; inner < pos+BASES_PER_LINE && inner < len(bases); inner += BASES_PER_BLOC {
					last := inner + BASES_PER_BLOC
					if last > len(bases) {
						last = len(bases)
					}
					frags = append(frags, bases[inner:last])
				}
				joined := strings.Join(frags, " ")
				lines = append(lines, joined)
			}
			return lines
		},
	}).Parse(tmplStr)

	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer

	if err := tmpl.Execute(&buf, seqs); err != nil {
		return nil, err
	} else {
		return lab.FileManager.WriteAll(buf.Bytes(), filename)
	}
}

// TextFile exports data in the format of a set of strings to a file.
// Each entry in the set of strings represents a line.
func TextFile(lab *laboratory.Laboratory, filename string, lines []string) (*wtype.File, error) {
	var sb strings.Builder
	for idx, line := range lines {
		if idx != 0 {
			if _, err := sb.WriteRune('\n'); err != nil {
				return nil, err
			}
		}
		if _, err := sb.WriteString(line); err != nil {
			return nil, err
		}
	}

	return lab.FileManager.WriteString(sb.String(), filename)
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
