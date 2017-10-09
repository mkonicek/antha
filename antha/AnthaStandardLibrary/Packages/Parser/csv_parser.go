// antha/AnthaStandardLibrary/Packages/Parser/csv_parser.go: Part of the Antha language
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

// Package for reading file formats, in particular focused toward dna sequence parsing
package parser

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"

	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/doe"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/enzymes"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/sequences"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
)

func ReadDesign(filename string) [][]string {

	var constructs [][]string

	csvfile, err := os.Open(filename)

	if err != nil {
		fmt.Println(err)
		return constructs
	}

	defer csvfile.Close()

	reader := csv.NewReader(csvfile)

	reader.FieldsPerRecord = -1 // see the Reader struct information below

	rawCSVdata, err := reader.ReadAll()

	if err != nil {
		panic(err)
	}

	// sanity check, display to standard output
	for _, each := range rawCSVdata {
		var parts []string

		if len(each[0]) > 1 {
			if string(strings.TrimSpace(each[0])[0]) != "#" {
				for _, p := range each {
					if p != "" {
						parts = append(parts, strings.TrimSpace(p))
					}
				}
				constructs = append(constructs, parts)
			}
		}
	}

	return constructs
}

func trim(s string) string {
	return strings.TrimSpace(s)
}

func readPartConcentrations(fileName string) (partNamesInOrder []string, concMap map[string]wunit.Concentration, err error) {

	concMap = make(map[string]wunit.Concentration)

	csvfile, err := os.Open(fileName)

	if err != nil {
		return
	}

	defer csvfile.Close()

	reader := csv.NewReader(csvfile)

	reader.FieldsPerRecord = -1 // see the Reader struct information below

	rawCSVdata, err := reader.ReadAll()

	if err != nil {
		return
	}
	var header []string

	var errs []string

	var nameColumn int = -1
	var concColumn int = -1

	for i, row := range rawCSVdata {

		// first row is a header
		if i == 0 {
			header = row

			for k, cell := range header {

				if strings.Contains(strings.ToLower(trim(cell)), "name") {
					nameColumn = k
				} else if strings.Contains(strings.ToLower(trim(cell)), "concentration") {
					concColumn = k
				}
			}

			if nameColumn < 0 {
				err := fmt.Errorf(`No column header found containing part "Name". Please add this as the first column. Found %v`, header)
				return partNamesInOrder, concMap, err
			}

			if concColumn < 0 {
				errs = append(errs, fmt.Sprintf(`No column header found containing part "Concentration". Please add this. Found %v`, header))
			}

		} else if !rowEmpty(row) {

			var partName string
			var partConc wunit.Concentration

			if nameColumn < len(row)-1 {
				partName = row[nameColumn]
			}

			partNamesInOrder = append(partNamesInOrder, partName)

			if strings.TrimSpace(partName) == "" {
				break
			}

			if concColumn < len(row) && concColumn >= 0 {

				partConc, err = doe.HandleConcFactor(header[concColumn], interface{}(row[concColumn]))
				if err != nil {
					errs = append(errs, err.Error())
				}
			} else {
				errs = append(errs, fmt.Sprintf("no concentration given for %s: concColumn %d, cells with date in row %d", partName, concColumn, len(row)))
			}

			if _, found := concMap[partName]; found {
				err = fmt.Errorf("duplicate part specified in parts list %s ", partName)
				errs = append(errs, err.Error())
			} else {
				concMap[partName] = partConc
			}
		}
	}
	if len(errs) > 0 {
		err = fmt.Errorf("Errors encountered parsing part concentrations: %s", strings.Join(errs, "; "))
	}
	return

}

func rowEmpty(row []string) bool {
	for _, cell := range row {
		if len(strings.TrimSpace(cell)) > 0 {
			return false
		}
	}
	return true
}

func ReadParts(filename string) map[string]wtype.DNASequence {

	m := make(map[string]wtype.DNASequence)

	var parts []wtype.DNASequence

	csvfile, err := os.Open(filename)

	if err != nil {
		fmt.Println(err)
		return m
	}

	defer csvfile.Close()

	reader := csv.NewReader(csvfile)

	reader.FieldsPerRecord = -1 // see the Reader struct information below

	rawCSVdata, err := reader.ReadAll()

	if err != nil {
		panic(err)
	}

	// sanity check, display to standard output
	for _, each := range rawCSVdata {

		var part wtype.DNASequence

		if len(each[0]) > 1 {
			if string(strings.TrimSpace(each[0])[0]) != "#" && len(each) > 1 {

				part.Nm = each[0]

				if strings.ToUpper(each[3]) == "AA" || strings.ToUpper(each[3]) == "Protein" || strings.ToUpper(each[3]) == "Amino Acid" {

					part.Seq = sequences.RevTranslatetoNstring(each[1])
				} else if each[3] == "DNA" || each[3] == "RNA" {
					part.Seq = strings.ToUpper(each[1])
				}

				part.Plasmid = false

				if len(each) > 2 {
					if strings.ToUpper(each[2]) == "TRUE" {
						part.Plasmid = true
					}
					if strings.ToUpper(each[2]) == "1" {
						part.Plasmid = true
					}
					if strings.ToUpper(each[2]) == "PLASMID" {
						part.Plasmid = true
					}
					if strings.ToUpper(each[2]) == "YES" {
						part.Plasmid = true
					}
					if strings.ToUpper(each[2]) == "FALSE" {
						part.Plasmid = false
					}
					if strings.ToUpper(each[2]) == "LINEAR" {
						part.Plasmid = false
					}
					if strings.ToUpper(each[2]) == "NO" {
						part.Plasmid = false
					}
					if strings.ToUpper(each[2]) == "0" {
						part.Plasmid = false
					}
				}

				parts = append(parts, part)
				m[part.Nm] = part
			}
		}
	}

	return m

}

func Assemblyfromcsv(designfile string, partsfile string) (assemblyparameters []enzymes.Assemblyparameters) {

	var designedconstructs [][]string

	designedconstructs = ReadDesign(designfile)

	var definedparts map[string]wtype.DNASequence

	definedparts = ReadParts(partsfile)

	assemblyparameters = make([]enzymes.Assemblyparameters, 0)

	for _, c := range designedconstructs {
		var newassemblyparameters enzymes.Assemblyparameters
		newassemblyparameters.Constructname = c[0]
		newassemblyparameters.Enzymename = c[1]
		newassemblyparameters.Vector = definedparts[c[2]]
		var nextpart wtype.DNASequence

		partsinorder := make([]wtype.DNASequence, 0)
		for k := 3; k < len(c); k++ {
			nextpart = definedparts[c[k]]
			if nextpart.Nm != "" {
				partsinorder = append(partsinorder, nextpart)
			}
		}
		newassemblyparameters.Partsinorder = partsinorder

		assemblyparameters = append(assemblyparameters, newassemblyparameters)

	}

	return assemblyparameters
}
