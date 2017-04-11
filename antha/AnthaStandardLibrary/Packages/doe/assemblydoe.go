// multipleassemblies.go Part of the Antha language
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

// Package for facilitating DOE methodology in antha
package doe

import (
	"bufio"
	"fmt"
	//"io/ioutil"
	"strconv"
	"strings"

	//"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/doe"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/enzymes"
	//"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/enzymes/lookup"
	//"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/export"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/Parser"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
)

func ReadGeneFiles(genefiles []wtype.File) {

	geneseqs := make([]interface{}, 0)

	for _, genefile := range genefiles {
		genes, err := parser.DNAFileToDNASequence(genefile, false)

		file, err := parser.Fastatocsv(genefile, "output.csv")
		if err != nil {
			panic(err.Error())
		}
		defer file.Close()

		var lines []string
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			lines = append(lines, scanner.Text())
		}
		fmt.Println(lines)

		for _, gene := range genes {
			fmt.Println("genefile:", genefile, "GENE name:", gene.Nm, "plasmid?", gene.Plasmid)
		}
		if err != nil {
			panic(err.Error())
		}
		vals := make([]interface{}, len(genes))
		for i, v := range genes {
			vals[i] = v
		}
		for _, gene := range vals {
			geneseqs = append(geneseqs, gene)
		}

	}

	// fmt.Println("number of sequences", len(geneseqs))

}

func AssemblyparametersfromRuns(runs []Run, enzymename string) (assemblyparameters []enzymes.Assemblyparameters) {

	assemblyparameters = make([]enzymes.Assemblyparameters, 0)
	parts := make([]wtype.DNASequence, 0)
	var parameters enzymes.Assemblyparameters

	for j, run := range runs {

		parameters.Constructname = "run" + strconv.Itoa(j)
		parameters.Enzymename = enzymename

		for i, _ := range run.Setpoints {

			if strings.Contains(run.Factordescriptors[i], "Vector") {

				parameters.Vector = run.Setpoints[i].(wtype.DNASequence)
			} else if i < len(runs)-1 {
				parts = append(parts, run.Setpoints[i].(wtype.DNASequence))
			}

		}
		parameters.Partsinorder = parts
		assemblyparameters = append(assemblyparameters, parameters)
	}

	return
}

func DNASequencetoInterface(genes []wtype.DNASequence) (geneseqs []interface{}) {

	vals := make([]interface{}, len(genes))
	for i, v := range genes {
		vals[i] = v
	}
	for _, gene := range vals {
		geneseqs = append(geneseqs, gene)
	}
	return
}

/*
// inputs
var (
	//vectors       []string = []string{"Rewiring_DaughterVector_I.gb"} // put vector sequence files here
	//promoterfiles []string = []string{"All_proms.fasta"}              // put promoter sequence files here
	genefiles []string = []string{"genes.fasta", "notworking.fasta", "working.fasta"} //genes    []string
	//part3files    []string = []string("")                             // put other part sequence files here if it's not in the right order you'll have to check the code and swap the order
//	enzyme  string = "bbsI"
//	dirname string = "assembly_export"
)
*/
