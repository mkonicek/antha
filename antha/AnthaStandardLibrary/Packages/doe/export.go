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

package doe

import (
	"fmt"

	"strconv"

	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/spreadsheet"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wutil"
	"github.com/tealeg/xlsx"
)

func JMPXLSXFilefromRuns(runs []Run, outputfilename string) (xlsxfile *xlsx.File) {

	// if output is a struct look for a sensible field to print

	//var file *xlsx.File
	var sheet *xlsx.Sheet
	var row *xlsx.Row
	var cell *xlsx.Cell
	var err error

	xlsxfile = xlsx.NewFile()
	sheet, err = xlsxfile.AddSheet("Sheet1")
	if err != nil {
		panic(err.Error())
	}
	// new row
	row = sheet.AddRow()

	// then add subheadings and descriptors
	if len(runs) == 0 {
		return
	}
	for _, descriptor := range runs[0].Factordescriptors {

		cell = row.AddCell()
		cell.Value = descriptor

	}
	for _, descriptor := range runs[0].Responsedescriptors {
		cell = row.AddCell()
		cell.Value = descriptor

	}
	for _, descriptor := range runs[0].AdditionalSubheaders {
		cell = row.AddCell()
		cell.Value = descriptor

	}
	//add data 1 row per run
	for _, run := range runs {

		row = sheet.AddRow()

		// factors
		for _, factor := range run.Setpoints {

			cell = row.AddCell()

			dna, amIdna := factor.(wtype.DNASequence)
			if amIdna {
				cell.SetValue(dna.Nm)
			} else {
				cell.SetValue(factor) //= factor.(string)
			}

		}

		// responses
		for _, response := range run.ResponseValues {
			cell = row.AddCell()
			cell.SetValue(response)
		}

		// additional
		for _, additional := range run.AdditionalValues {
			cell = row.AddCell()
			cell.SetValue(additional)
		}
	}
	err = xlsxfile.Save(outputfilename)
	if err != nil {
		fmt.Print(err.Error())
	}
	return
}

func DXXLSXFilefromRuns(runs []Run, outputfilename string) (xlsxfile *xlsx.File) {

	// if output is a struct look for a sensible field to print

	var sheet *xlsx.Sheet
	var row *xlsx.Row
	var cell *xlsx.Cell
	var err error

	xlsxfile = xlsx.NewFile()
	sheet, err = xlsxfile.AddSheet("Sheet1")
	if err != nil {
		panic(err.Error())
	}
	// add headers
	row = sheet.AddRow()

	// 2 blank cells
	cell = row.AddCell()
	cell.Value = ""
	cell = row.AddCell()
	cell.Value = ""

	// take factor and run descriptors from first run (assuming they're all the same)
	for i := range runs[0].Factordescriptors {
		cell = row.AddCell()
		cell.Value = "Factor " + strconv.Itoa(i+1)

	}
	for i := range runs[0].Responsedescriptors {
		cell = row.AddCell()
		cell.Value = "Response " + strconv.Itoa(i+1)

	}
	for _, additionalheader := range runs[0].AdditionalHeaders {
		cell = row.AddCell()
		cell.Value = additionalheader

	}
	// new row
	row = sheet.AddRow()

	// add Std and Run number headers
	cell = row.AddCell()
	cell.Value = "Std"
	cell = row.AddCell()
	cell.Value = "Run"

	// then add subheadings and descriptors
	for i, descriptor := range runs[0].Factordescriptors {
		letter := wutil.NumToAlpha(i + 1)
		cell = row.AddCell()
		cell.Value = letter + ":" + descriptor

	}
	for _, descriptor := range runs[0].Responsedescriptors {
		cell = row.AddCell()
		cell.Value = descriptor

	}
	for _, descriptor := range runs[0].AdditionalSubheaders {
		cell = row.AddCell()
		cell.Value = descriptor

	}

	// add blank row
	sheet.AddRow()

	//add data 1 row per run
	for _, run := range runs {

		row = sheet.AddRow()
		// Std
		cell = row.AddCell()
		cell.SetValue(run.StdNumber)

		// Run
		cell = row.AddCell()
		cell.SetValue(run.RunNumber)

		// factors
		for _, factor := range run.Setpoints {

			cell = row.AddCell()

			dna, amIdna := factor.(wtype.DNASequence)
			if amIdna {
				cell.SetValue(dna.Nm)
			} else {
				cell.SetValue(factor) //= factor.(string)
			}

		}

		// responses
		for _, response := range run.ResponseValues {
			cell = row.AddCell()
			cell.SetValue(response)
		}

		// additional
		for _, additional := range run.AdditionalValues {
			cell = row.AddCell()
			cell.SetValue(additional)
		}
	}
	err = xlsxfile.Save(outputfilename)
	if err != nil {
		fmt.Print(err.Error())
	}
	return
}

func XLSXFileFromRuns(runs []Run, outputfilename string, dxorjmp string) (xlsxfile *xlsx.File) {
	if dxorjmp == "DX" {
		xlsxfile = DXXLSXFilefromRuns(runs, outputfilename)
	} else if dxorjmp == "JMP" {
		xlsxfile = JMPXLSXFilefromRuns(runs, outputfilename)
	} else {
		panic(fmt.Sprintf("Unknown design file format %s when exporting design to XLSX file. Please specify File type as JMP or DX (Design Expert)", dxorjmp))
	}
	spreadsheet.AutoFormat(xlsxfile)
	return
}
