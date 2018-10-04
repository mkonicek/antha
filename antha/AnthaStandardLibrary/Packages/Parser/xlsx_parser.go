// Part of the Antha language
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

package parser

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/enzymes"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/tealeg/xlsx"
)

const delimiter = ","

type outputer func(s string)

func generateCSVFromXLSXsheet(excelFileName string, sheetIndex int, outputf outputer) error {
	xlFile, error := xlsx.OpenFile(excelFileName)
	if error != nil {
		return error
	}
	sheetLen := len(xlFile.Sheets)
	switch {
	case sheetLen == 0:
		return errors.New("this XLSX file contains no sheets")
	case sheetIndex >= sheetLen:
		return fmt.Errorf("no sheet %d available, please select a sheet between 0 and %d", sheetIndex, sheetLen-1)
	}
	sheet := xlFile.Sheets[sheetIndex]
	for _, row := range sheet.Rows {
		var vals []string
		if row != nil {
			for _, cell := range row.Cells {
				cellstr := cell.String()
				vals = append(vals, fmt.Sprintf("%q", cellstr))
			}
			outputf(strings.Join(vals, delimiter) + "\n")
		}
	}
	return nil
}

func generateCSVFromXLSXsheetBinary(excelFileContents []byte, sheetIndex int, outputf outputer) error {
	xlFile, error := xlsx.OpenBinary(excelFileContents)
	if error != nil {
		return error
	}
	sheetLen := len(xlFile.Sheets)
	switch {
	case sheetLen == 0:
		return errors.New("this XLSX file contains no sheets")
	case sheetIndex >= sheetLen:
		return fmt.Errorf("no sheet %d available, please select a sheet between 0 and %d", sheetIndex, sheetLen-1)
	}
	sheet := xlFile.Sheets[sheetIndex]
	for _, row := range sheet.Rows {
		var vals []string
		if row != nil {
			for _, cell := range row.Cells {

				cellstr := cell.String()
				cellstr = strings.TrimSpace(cellstr)

				vals = append(vals, fmt.Sprintf("%q", cellstr))
			}
			outputf(strings.Join(vals, delimiter) + "\n")
		}
	}
	return nil
}

func Xlsxparser(filename string, sheetIndex int, outputprefix string) (f *os.File, err error) {
	f, err = ioutil.TempFile("", outputprefix)
	if err != nil {
		return
	}

	printer := func(s string) {
		f.WriteString(s) // nolint
	}

	err = generateCSVFromXLSXsheet(filename, sheetIndex, printer)
	return
}

func ParseExcel(filename string) ([]enzymes.Assemblyparameters, error) {
	if pl, err := Xlsxparser(filename, 0, "partslist"); err != nil {
		return nil, err
	} else if dl, err := Xlsxparser(filename, 1, "designlist"); err != nil {
		return nil, err
	} else {
		ap, err := AssemblyFromCsv(dl.Name(), pl.Name())
		return ap, err
	}
}

func xlsxparserBinary(data []byte, sheetIndex int, outputprefix string) (f *os.File, err error) {
	f, err = ioutil.TempFile("", outputprefix)
	if err != nil {
		return
	}
	printer := func(s string) {
		_, _ = f.WriteString(s)
	}
	err = generateCSVFromXLSXsheetBinary(data, sheetIndex, printer)
	return
}

// ParseExcelBinary parses the contents of a typeIIs assembly design file in
// xlsx format. An example file is provided: "Assembly_Input_Controls.xlsx" The
// output will be []enzymes.AssemblyParameters which can be used in the
// enzymes.Assemblysimulator() and enzymes.Digestionsimulator() functions.  The
// design file is expected to follow a format as shown in the provided example
// files An error will be returned if no data is found within the .xlsx design
// file or if the file is not in the expected format.
func ParseExcelBinary(data []byte) ([]enzymes.Assemblyparameters, error) {
	if pl, err := xlsxparserBinary(data, 0, "partslist"); err != nil {
		return nil, err
	} else if dl, err := xlsxparserBinary(data, 1, "designlist"); err != nil {
		return nil, err
	} else {
		ap, err := AssemblyFromCsv(dl.Name(), pl.Name())
		return ap, err
	}
}

// MakePartsFromXLSXPartsList parses the parts in an xlsx format design file
// into a list of LHComponents.  The concentration will be set if a
// concentration column is present in the parts list.  If no concentrations are
// found the parts list will be created with no concentrations and an error
// returned.
func MakePartsFromXLSXPartsList(data []byte) (parts []*wtype.Liquid, concMap map[string]wunit.Concentration, err error) {
	pl, err := xlsxparserBinary(data, 0, "partslist")
	if err != nil {
		return nil, nil, err
	}

	partSeqs := ReadParts(pl.Name())

	partNamesInOrder, concMap, err := readPartConcentrations(pl.Name())

	if err != nil {
		// don't return if the first error is no concentration column found
		if !strings.Contains(err.Error(), `Errors encountered parsing part concentrations: No column header found containing part "Concentration"`) {
			return
		}
	}

	for _, partName := range partNamesInOrder {
		newComponent := wtype.NewLHComponent()
		newComponent.CName = partName
		if concMap[partName].RawValue() != 0 {
			newComponent.SetConcentration(concMap[partName])
		}
		err := newComponent.AddDNASequence(partSeqs[partName])
		if err != nil {
			return nil, nil, err
		}
		parts = append(parts, newComponent)
	}

	return
}
