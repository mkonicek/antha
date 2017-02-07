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
	"github.com/tealeg/xlsx"
)

var delimiter = ","

type outputer func(s string)

func generateCSVFromXLSXsheet(excelFileName string, sheetIndex int, outputf outputer) error {
	xlFile, error := xlsx.OpenFile(excelFileName)
	if error != nil {
		return error
	}
	sheetLen := len(xlFile.Sheets)
	switch {
	case sheetLen == 0:
		return errors.New("This XLSX file contains no sheets.")
	case sheetIndex >= sheetLen:
		return fmt.Errorf("No sheet %d available, please select a sheet between 0 and %d\n", sheetIndex, sheetLen-1)
	}
	sheet := xlFile.Sheets[sheetIndex]
	for _, row := range sheet.Rows {
		var vals []string
		if row != nil {
			for _, cell := range row.Cells {
				cellstr, err := cell.String()
				if err != nil {
					return err
				}
				vals = append(vals, fmt.Sprintf("%q", cellstr))
			}
			outputf(strings.Join(vals, delimiter) + "\n")
		}
	}
	return nil
}

func generateCSVFromspecificXLSXsheet(excelFileName string, sheetname string, outputf outputer) error {
	xlFile, error := xlsx.OpenFile(excelFileName)
	if error != nil {
		return error
	}
	sheetLen := len(xlFile.Sheets)
	switch {
	case sheetLen == 0:
		return errors.New("This XLSX file contains no sheets.")
		for _, sheet := range xlFile.Sheets {
			for _, row := range sheet.Rows {
				var vals []string
				if row != nil {
					for _, cell := range row.Cells {
						cellstr, err := cell.String()
						if err != nil {
							return err
						}
						vals = append(vals, fmt.Sprintf("%q", cellstr))
					}
					outputf(strings.Join(vals, delimiter) + "\n")
				}
			}
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
		_, _ = f.WriteString(s)
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
		return Assemblyfromcsv(dl.Name(), pl.Name()), nil
	}
}
