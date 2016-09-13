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

// Package for exporting to file
package export

import (
	"bytes"
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
func Exporttofile(dir string, seq wtype.BioSequence) (string, error) {
	filename := filepath.Join(anthapath.Path(), fmt.Sprintf("%s_%s.txt", dir, seq.Name()))
	if err := os.MkdirAll(filepath.Dir(filename), 0777); err != nil {
		return "", err
	}

	f, err := os.Create(filename)
	if err != nil {
		return "", err
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

	return filename, err
}

// function to export a sequence to a txt file
func ExportFasta(dir string, seq wtype.BioSequence) (string, error) {
	filename := filepath.Join(anthapath.Path(), fmt.Sprintf("%s_%s.fasta", dir, seq.Name()))
	if err := os.MkdirAll(filepath.Dir(filename), 0777); err != nil {
		return "", err
	}

	f, err := os.Create(filename)
	if err != nil {
		return "", err
	}
	defer f.Close()

	_, err = fmt.Fprintf(f, ">%s\n%s\n", seq.Name(), seq.Sequence())

	return filename, err
}

// function to export multiple sequences in fasta format into a single txt file
// Modify this for the more general case
func Makefastaserial(dir string, seqs []*wtype.DNASequence) (string, error) {
	filename := filepath.Join(anthapath.Path(), fmt.Sprintf("%s.fasta", dir))
	if err := os.MkdirAll(filepath.Dir(filename), 0777); err != nil {
		return "", err
	}

	f, err := os.Create(filename)
	if err != nil {
		return "", err
	}
	defer f.Close()

	for _, seq := range seqs {
		if _, err := fmt.Fprintf(f, ">%s\n%s\n", seq.Name(), seq.Sequence()); err != nil {
			return "", err
		}
	}

	return filename, nil
}

func Makefastaserial2(makeinanthapath bool, dir string, seqs []wtype.DNASequence) (string, error) {

	var filename string
	if makeinanthapath {
		filename = filepath.Join(anthapath.Path(), fmt.Sprintf("%s.fasta", dir))
	} else {
		filename = filepath.Join(fmt.Sprintf("%s.fasta", dir))
	}
	if err := os.MkdirAll(filepath.Dir(filename), 0777); err != nil {
		return "", err
	}

	f, err := os.Create(filename)
	if err != nil {
		return "", err
	}
	defer f.Close()

	for _, seq := range seqs {
		if _, err := fmt.Fprintf(f, ">%s\n%s\n", seq.Name(), seq.Sequence()); err != nil {
			return "", err
		}
	}

	return filename, nil
}

func ExportFastaandSummaryforEachSeq(assemblyparameters enzymes.Assemblyparameters) (err error) {

	enzymename := strings.ToUpper(assemblyparameters.Enzymename)

	// should change this to rebase lookup; what happens if this fails?
	//enzyme := TypeIIsEnzymeproperties[enzymename]
	enzyme, err := lookup.TypeIIsLookup(enzymename)
	if err != nil {
		return err
	}
	//assemble (note that sapIenz is found in package enzymes)
	_, plasmidproductsfromXprimaryseq, err := enzymes.JoinXnumberofparts(assemblyparameters.Vector, assemblyparameters.Partsinorder, enzyme)

	if err != nil {
		return err
	}

	for _, assemblyproduct := range plasmidproductsfromXprimaryseq {
		filename := filepath.Join(anthapath.Path(), assemblyparameters.Constructname)
		if _, err := Exporttofile(filename, &assemblyproduct); err != nil {
			return err
		}

		if _, err := ExportFasta(filename, &assemblyproduct); err != nil {
			return err
		}
	}

	return nil
}

func ExportFastaSerialfromMultipleAssemblies(dirname string, multipleassemblyparameters []enzymes.Assemblyparameters) (string, error) {

	seqs := make([]wtype.DNASequence, 0)

	for _, assemblyparameters := range multipleassemblyparameters {

		enzymename := strings.ToUpper(assemblyparameters.Enzymename)

		// should change this to rebase lookup; what happens if this fails?
		//enzyme := TypeIIsEnzymeproperties[enzymename]
		enzyme, err := lookup.TypeIIsLookup(enzymename)
		if err != nil {
			return "", err
		}
		//assemble (note that sapIenz is found in package enzymes)
		_, plasmidproductsfromXprimaryseq, err := enzymes.JoinXnumberofparts(assemblyparameters.Vector, assemblyparameters.Partsinorder, enzyme)
		if err != nil {
			return "", err
		}

		for _, assemblyproduct := range plasmidproductsfromXprimaryseq {

			/*	fileprefix := anthapath.Dirpath() + "/"
				tojoin := make([]string, 0)
				tojoin = append(tojoin, fileprefix, assemblyparameters.Constructname)
				filename := strings.Join(tojoin, "")
				Exporttofile(filename, &assemblyproduct)
				ExportFasta(filename, &assemblyproduct)*/

			seqs = append(seqs, assemblyproduct)
		}

	}

	return Makefastaserial2(ANTHAPATH, dirname, seqs)
}

func ExporttoTextFile(filename string, data []string) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	for _, str := range data {
		if _, err := fmt.Fprintln(f, str); err != nil {
			return err
		}
	}

	return nil
}

func ExporttoJSON(data interface{}, filename string) (err error) {
	bytes, err := json.Marshal(data)

	if err != nil {
		return err
	}

	ioutil.WriteFile(filename, bytes, 0644)
	return nil
}
