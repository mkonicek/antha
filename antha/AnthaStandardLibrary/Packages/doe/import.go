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
	"strings"

	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/search"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/spreadsheet"
	"github.com/tealeg/xlsx"
)

func parseRunWellPair(pair string, nameappendage string) (runnumber int, well string, err error) {
	split := strings.Split(pair, ":")

	numberstring := strings.SplitAfter(split[0], nameappendage)

	runnumber, err = strconv.Atoi(string(numberstring[1]))
	if err != nil {
		err = fmt.Errorf(err.Error(), "+ Failed at", pair, nameappendage)
	}
	well = split[1]
	return
}

func rowEmpty(sheet *xlsx.Sheet, rownumber int) bool {
	for i := 0; i < sheet.MaxCol; i++ {
		cell := sheet.Cell(rownumber, i)
		str := cell.String()

		if len(strings.TrimSpace(str)) > 0 {
			return false
		}
	}
	return true
}

func AddWelllocations(DXORJMP string, xlsxfile string, oldsheet int, runnumbertowellcombos []string, nameappendage string, pathtosave string, extracolumnheaders []string, extracolumnvalues []interface{}) error {

	var xlsxcell *xlsx.Cell

	file, err := spreadsheet.OpenXLSXFromFileName(xlsxfile)
	if err != nil {
		return err
	}

	sheet, err := spreadsheet.Sheet(file, oldsheet)
	if err != nil {
		return err
	}

	_, _ = file.AddSheet("hello")

	//extracolumn := sheet.MaxCol + 1

	// add extra column headers first
	for _, extracolumnheader := range extracolumnheaders {
		xlsxcell = sheet.Rows[0].AddCell()

		xlsxcell.Value = "Extra column added"
		// fmt.Println("CEllll added succesfully", sheet.Cell(0, extracolumn).String())
		xlsxcell = sheet.Rows[1].AddCell()
		xlsxcell.Value = extracolumnheader
	}

	// now add well position column
	xlsxcell = sheet.Rows[0].AddCell()

	xlsxcell.Value = "Location"
	// fmt.Println("CEllll added succesfully", sheet.Cell(0, extracolumn).String())
	xlsxcell = sheet.Rows[1].AddCell()
	xlsxcell.Value = "Well ID"

	for i := 3; i < sheet.MaxRow; i++ {
		for _, pair := range runnumbertowellcombos {
			runnumber, well, err := parseRunWellPair(pair, nameappendage)
			if err != nil {
				return err
			}
			xlsxrunmumber, err := sheet.Cell(i, 1).Int()
			if err != nil {
				return err
			}
			if xlsxrunmumber == runnumber {
				for _, extracolumnvalue := range extracolumnvalues {
					xlsxcell = sheet.Rows[i].AddCell()
					xlsxcell.SetValue(extracolumnvalue)
				}
				xlsxcell = sheet.Rows[i].AddCell()
				xlsxcell.Value = well

			}
		}
	}

	err = file.Save(pathtosave)

	return err
}

func RunsFromDXDesignContents(bytes []byte, intfactors []string) (runs []Run, err error) {
	file, err := spreadsheet.OpenXLSXBinary(bytes)
	if err != nil {
		return runs, err
	}
	sheet, err := spreadsheet.Sheet(file, 0)
	if err != nil {
		return runs, err
	}
	runs = make([]Run, 0)
	var run Run

	var setpoint interface{}
	var descriptor string

	for i := 3; i < sheet.MaxRow; i++ {

		if !rowEmpty(sheet, i) {

			factordescriptors := make([]string, 0)
			responsedescriptors := make([]string, 0)
			setpoints := make([]interface{}, 0)
			responsevalues := make([]interface{}, 0)
			otherheaders := make([]string, 0)
			othersubheaders := make([]string, 0)
			otherresponsevalues := make([]interface{}, 0)

			run.RunNumber, err = sheet.Cell(i, 1).Int()
			if err != nil {
				return runs, err
			}
			run.StdNumber, err = sheet.Cell(i, 0).Int()
			if err != nil {
				return runs, err
			}

			for j := 2; j < sheet.MaxCol; j++ {
				factororresponse := sheet.Cell(0, j).String()

				if strings.Contains(factororresponse, "Factor") {

					desc := sheet.Cell(1, j).String()

					descriptor = strings.Split(desc, ":")[1]
					factrodescriptor := descriptor
					//fmt.Println(i, j, descriptor)

					cell := sheet.Cell(i, j)

					celltype := cell.Type()

					_, err = cell.Float()

					if strings.ToUpper(cell.Value) == "TRUE" {
						setpoint = true //cell.SetBool(true)
					} else if strings.ToUpper(cell.Value) == "FALSE" {
						setpoint = false //cell.SetBool(false)
					} else if celltype == 3 {
						setpoint = cell.Bool()
					} else if err == nil || celltype == 1 {
						setpoint, _ = cell.Float()
						if search.InStrings(intfactors, descriptor) {
							setpoint, err = cell.Int()
							if err != nil {
								return runs, err
							}
						}
					} else {
						setpoint = cell.String()
					}

					factordescriptors = append(factordescriptors, factrodescriptor)
					setpoints = append(setpoints, setpoint)

				} else if strings.Contains(factororresponse, "Response") {
					descriptor = sheet.Cell(1, j).String()
					responsedescriptor := descriptor
					//// fmt.Println("response", i, j, descriptor)
					responsedescriptors = append(responsedescriptors, responsedescriptor)

					cell := sheet.Cell(i, j)

					if cell == nil {

						break
					}

					celltype := cell.Type()

					if celltype == 1 {
						responsevalue, err := cell.Float()
						if err != nil {
							return runs, err
						}
						responsevalues = append(responsevalues, responsevalue)
					} else {
						responsevalue := cell.String()
						responsevalues = append(responsevalues, responsevalue)
					}

				} else {
					descriptor = sheet.Cell(1, j).String()
					responsedescriptor := descriptor

					otherheaders = append(otherheaders, factororresponse)
					othersubheaders = append(othersubheaders, responsedescriptor)

					cell := sheet.Cell(i, j)

					if cell == nil {

						break
					}

					celltype := cell.Type()

					if celltype == 1 {
						responsevalue, err := cell.Float()
						if err != nil {
							return runs, err
						}
						otherresponsevalues = append(otherresponsevalues, responsevalue)
					} else {
						responsevalue := cell.String()
						otherresponsevalues = append(otherresponsevalues, responsevalue)
					}

				}
			}
			run.Factordescriptors = factordescriptors
			run.Responsedescriptors = responsedescriptors
			run.Setpoints = setpoints
			run.ResponseValues = responsevalues
			run.AdditionalHeaders = otherheaders
			run.AdditionalSubheaders = othersubheaders
			run.AdditionalValues = otherresponsevalues

			runs = append(runs, run)
		}
	}

	return
}

func RunsFromJMPDesignContents(bytes []byte, factorcolumns []int, responsecolumns []int, intfactors []string) (runs []Run, err error) {
	file, err := spreadsheet.OpenXLSXBinary(bytes)
	if err != nil {
		return runs, err
	}
	sheet, err := spreadsheet.Sheet(file, 0)
	if err != nil {
		return runs, err
	}
	runs = make([]Run, 0)
	var run Run

	var setpoint interface{}
	var descriptor string
	for i := 1; i < sheet.MaxRow; i++ {

		if !rowEmpty(sheet, i) {

			factordescriptors := make([]string, 0)
			responsedescriptors := make([]string, 0)
			setpoints := make([]interface{}, 0)
			responsevalues := make([]interface{}, 0)
			otherheaders := make([]string, 0)
			othersubheaders := make([]string, 0)
			otherresponsevalues := make([]interface{}, 0)

			run.RunNumber = i //sheet.Cell(i, 1).Int()

			run.StdNumber = i //sheet.Cell(i, 0).Int()

			for j := 0; j < sheet.MaxCol; j++ {

				var factororresponse string

				if search.InInts(factorcolumns, j) {
					factororresponse = "Factor"
				} else if search.InInts(responsecolumns, j) {
					factororresponse = "Response"
				}

				if strings.Contains(factororresponse, "Factor") {

					descriptor = sheet.Cell(0, j).String()
					factrodescriptor := descriptor
					cell := sheet.Cell(i, j)

					celltype := cell.Type()

					_, err := cell.Float()

					if strings.ToUpper(cell.Value) == "TRUE" {
						setpoint = true //cell.SetBool(true)
					} else if strings.ToUpper(cell.Value) == "FALSE" {
						setpoint = false //cell.SetBool(false)
					} else if celltype == 3 {
						setpoint = cell.Bool()
					} else if err == nil || celltype == 1 {
						setpoint, _ = cell.Float()
						if search.InStrings(intfactors, descriptor) {
							setpoint, err = cell.Int()
							if err != nil {
								return runs, err
							}
						}
					} else {
						setpoint = cell.String()
					}
					factordescriptors = append(factordescriptors, factrodescriptor)
					setpoints = append(setpoints, setpoint)

				} else if strings.Contains(factororresponse, "Response") {
					descriptor = sheet.Cell(0, j).String()
					responsedescriptor := descriptor

					responsedescriptors = append(responsedescriptors, responsedescriptor)

					cell := sheet.Cell(i, j)

					if cell == nil {

						break
					}

					celltype := cell.Type()

					if celltype == 1 {
						responsevalue, err := cell.Float()
						if err != nil {
							return runs, err
						}
						responsevalues = append(responsevalues, responsevalue)
					} else {
						responsevalue := cell.String()
						responsevalues = append(responsevalues, responsevalue)
					}

				} else /*if j != patterncolumn*/ {
					descriptor = sheet.Cell(0, j).String()
					responsedescriptor := descriptor

					otherheaders = append(otherheaders, factororresponse)
					othersubheaders = append(othersubheaders, responsedescriptor)

					cell := sheet.Cell(i, j)

					if cell == nil {

						break
					}

					celltype := cell.Type()

					if responsevalue, err := cell.Float(); celltype == 1 && err == nil {

						// formulas are also given the celltype 1
						if err != nil {
							return runs, err
						}
						otherresponsevalues = append(otherresponsevalues, responsevalue)
					} else {
						responsevalue := cell.String()
						otherresponsevalues = append(otherresponsevalues, responsevalue)
					}

				}
			}
			run.Factordescriptors = factordescriptors
			run.Responsedescriptors = responsedescriptors
			run.Setpoints = setpoints
			run.ResponseValues = responsevalues
			run.AdditionalHeaders = otherheaders
			run.AdditionalSubheaders = othersubheaders
			run.AdditionalValues = otherresponsevalues
			runs = append(runs, run)
		}
	}

	return
}

func RunsFromDesign(contents []byte, intfactors []string, dxorjmp string) (runs []Run, err error) {

	if dxorjmp == "DX" {

		runs, err = RunsFromDXDesignContents(contents, intfactors)
		if err != nil {
			return runs, err
		}

	} else if dxorjmp == "JMP" {

		factorcolumns, responsecolumns, _ := findJMPFactorandResponseColumnsinEmptyDesignContents(contents)

		runs, err = RunsFromJMPDesignContents(contents, factorcolumns, responsecolumns, intfactors)
		if err != nil {
			return runs, err
		}
	} else {
		err = fmt.Errorf("Unknown design file format. Please specify File type as JMP or DX (Design Expert)")
	}
	return
}

func RunsFromDesignPreResponsesContents(designfileContents []byte, intfactors []string, dxorjmp string) (runs []Run, err error) {

	if dxorjmp == "DX" {

		runs, err = RunsFromDXDesignContents(designfileContents, intfactors)
		if err != nil {
			return runs, err
		}

	} else if dxorjmp == "JMP" {

		factorcolumns, responsecolumns, _ := findJMPFactorandResponseColumnsinEmptyDesignContents(designfileContents)

		runs, err = RunsFromJMPDesignContents(designfileContents, factorcolumns, responsecolumns, intfactors)
		if err != nil {
			return runs, err
		}
	} else {
		err = fmt.Errorf("Unknown design file format. Please specify File type as JMP or DX (Design Expert)")
	}
	return

}

// add func to auto check for Response and factor status based on empty entries implying Response column
func findJMPFactorandResponseColumnsinEmptyDesignContents(bytes []byte) (factorcolumns []int, responsecolumns []int, PatternColumn int) {
	var patternfound bool
	factorcolumns = make([]int, 0)
	responsecolumns = make([]int, 0)

	file, err := spreadsheet.OpenXLSXBinary(bytes)
	if err != nil {
		return
	}
	sheet, err := spreadsheet.Sheet(file, 0)
	if err != nil {
		return
	}
	for j := 0; j < sheet.MaxCol; j++ {

		descriptor := sheet.Cell(0, j).String()
		if strings.ToUpper(descriptor) == "PATTERN" {
			PatternColumn = j
			patternfound = true
		}
	}
	// iterate through every run of the design sheet (row) and if all values for that row == "", the column is interpreted as a response
	for i := 1; i < sheet.MaxRow; i++ {
		//maxfactorcol := 2
		for j := 0; j < sheet.MaxCol; j++ {
			cellstr := sheet.Cell(i, j).String()

			if patternfound && j != PatternColumn && cellstr != "" {
				factorcolumns = append(factorcolumns, j)
			} else if !patternfound && cellstr != "" {
				factorcolumns = append(factorcolumns, j)
			} else if cellstr == "" {

				responsecolumns = append(responsecolumns, j)
			}

		}

	}

	factorcolumns = search.RemoveDuplicateInts(factorcolumns)
	responsecolumns = search.RemoveDuplicateInts(responsecolumns)

	return
}

/////////

// DEPRECATE THESE FUNCS

func RunsFromDesignPreResponses(designfile string, intfactors []string, dxorjmp string) (runs []Run, err error) {

	if dxorjmp == "DX" {

		runs, err = RunsFromDXDesign(designfile, intfactors)
		if err != nil {
			return runs, err
		}

	} else if dxorjmp == "JMP" {

		factorcolumns, responsecolumns, _ := findJMPFactorandResponseColumnsinEmptyDesign(designfile)

		runs, err = RunsFromJMPDesign(designfile, factorcolumns, responsecolumns, intfactors)
		if err != nil {
			return runs, err
		}
	} else {
		err = fmt.Errorf("Unknown design file format. Please specify File type as JMP or DX (Design Expert)")
	}
	return

}

func RunsFromDXDesign(filename string, intfactors []string) (runs []Run, err error) {
	file, err := spreadsheet.OpenXLSXFromFileName(filename)
	if err != nil {
		return runs, err
	}
	sheet, err := spreadsheet.Sheet(file, 0)
	if err != nil {
		return runs, err
	}
	runs = make([]Run, 0)
	var run Run

	var setpoint interface{}
	var descriptor string
	for i := 3; i < sheet.MaxRow; i++ {
		if !rowEmpty(sheet, i) {

			factordescriptors := make([]string, 0)
			responsedescriptors := make([]string, 0)
			setpoints := make([]interface{}, 0)
			responsevalues := make([]interface{}, 0)
			otherheaders := make([]string, 0)
			othersubheaders := make([]string, 0)
			otherresponsevalues := make([]interface{}, 0)

			run.RunNumber, err = sheet.Cell(i, 1).Int()
			if err != nil {
				return runs, err
			}
			run.StdNumber, err = sheet.Cell(i, 0).Int()
			if err != nil {
				return runs, err
			}

			for j := 2; j < sheet.MaxCol; j++ {
				factororresponse := sheet.Cell(0, j).String()

				if strings.Contains(factororresponse, "Factor") {

					desc := sheet.Cell(1, j).String()

					descriptor = strings.Split(desc, ":")[1]
					factrodescriptor := descriptor
					//fmt.Println(i, j, descriptor)

					cell := sheet.Cell(i, j)

					celltype := cell.Type()

					_, err = cell.Float()

					if strings.ToUpper(cell.Value) == "TRUE" {
						setpoint = true //cell.SetBool(true)
					} else if strings.ToUpper(cell.Value) == "FALSE" {
						setpoint = false //cell.SetBool(false)
					} else if celltype == 3 {
						setpoint = cell.Bool()
					} else if err == nil || celltype == 1 {
						setpoint, _ = cell.Float()
						if search.InStrings(intfactors, descriptor) {
							setpoint, err = cell.Int()
							if err != nil {
								return runs, err
							}
						}
					} else {

						setpoint = cell.String()
					}

					factordescriptors = append(factordescriptors, factrodescriptor)
					setpoints = append(setpoints, setpoint)

				} else if strings.Contains(factororresponse, "Response") {
					descriptor = sheet.Cell(1, j).String()
					responsedescriptor := descriptor
					//// fmt.Println("response", i, j, descriptor)
					responsedescriptors = append(responsedescriptors, responsedescriptor)

					cell := sheet.Cell(i, j)

					if cell == nil {

						break
					}

					celltype := cell.Type()

					if celltype == 1 {
						responsevalue, err := cell.Float()
						if err != nil {
							return runs, err
						}
						responsevalues = append(responsevalues, responsevalue)
					} else {
						responsevalue := cell.String()
						responsevalues = append(responsevalues, responsevalue)
					}

				} else {
					descriptor = sheet.Cell(1, j).String()
					responsedescriptor := descriptor

					otherheaders = append(otherheaders, factororresponse)
					othersubheaders = append(othersubheaders, responsedescriptor)

					cell := sheet.Cell(i, j)

					if cell == nil {

						break
					}

					celltype := cell.Type()

					if celltype == 1 {
						responsevalue, err := cell.Float()
						if err != nil {
							return runs, err
						}
						otherresponsevalues = append(otherresponsevalues, responsevalue)
					} else {
						responsevalue := cell.String()
						otherresponsevalues = append(otherresponsevalues, responsevalue)
					}

				}
			}
			run.Factordescriptors = factordescriptors
			run.Responsedescriptors = responsedescriptors
			run.Setpoints = setpoints
			run.ResponseValues = responsevalues
			run.AdditionalHeaders = otherheaders
			run.AdditionalSubheaders = othersubheaders
			run.AdditionalValues = otherresponsevalues

			runs = append(runs, run)
		}
	}

	return
}

func RunsFromJMPDesign(xlsx string, factorcolumns []int, responsecolumns []int, intfactors []string) (runs []Run, err error) {
	file, err := spreadsheet.OpenXLSXFromFileName(xlsx)
	if err != nil {
		return runs, err
	}
	sheet, err := spreadsheet.Sheet(file, 0)
	if err != nil {
		return runs, err
	}
	runs = make([]Run, 0)
	var run Run

	var setpoint interface{}
	var descriptor string
	for i := 1; i < sheet.MaxRow; i++ {

		if !rowEmpty(sheet, i) {

			factordescriptors := make([]string, 0)
			responsedescriptors := make([]string, 0)
			setpoints := make([]interface{}, 0)
			responsevalues := make([]interface{}, 0)
			otherheaders := make([]string, 0)
			othersubheaders := make([]string, 0)
			otherresponsevalues := make([]interface{}, 0)

			run.RunNumber = i //sheet.Cell(i, 1).Int()

			run.StdNumber = i //sheet.Cell(i, 0).Int()

			for j := 0; j < sheet.MaxCol; j++ {

				var factororresponse string

				if search.InInts(factorcolumns, j) {
					factororresponse = "Factor"
				} else if search.InInts(responsecolumns, j) {
					factororresponse = "Response"
				}

				if strings.Contains(factororresponse, "Factor") {

					descriptor = sheet.Cell(0, j).String()
					factrodescriptor := descriptor
					cell := sheet.Cell(i, j)

					celltype := cell.Type()

					_, err := cell.Float()

					if strings.ToUpper(cell.Value) == "TRUE" {
						setpoint = true //cell.SetBool(true)
					} else if strings.ToUpper(cell.Value) == "FALSE" {
						setpoint = false //cell.SetBool(false)
					} else if celltype == 3 {
						setpoint = cell.Bool()
					} else if err == nil || celltype == 1 {
						setpoint, _ = cell.Float()
						if search.InStrings(intfactors, descriptor) {
							setpoint, err = cell.Int()
							if err != nil {
								return runs, err
							}
						}
					} else {
						setpoint = cell.String()
					}
					factordescriptors = append(factordescriptors, factrodescriptor)
					setpoints = append(setpoints, setpoint)

				} else if strings.Contains(factororresponse, "Response") {
					descriptor = sheet.Cell(0, j).String()
					responsedescriptor := descriptor

					responsedescriptors = append(responsedescriptors, responsedescriptor)

					cell := sheet.Cell(i, j)

					if cell == nil {

						break
					}

					celltype := cell.Type()

					if celltype == 1 {
						responsevalue, err := cell.Float()
						if err != nil {
							return runs, err
						}
						responsevalues = append(responsevalues, responsevalue)
					} else {
						responsevalue := cell.String()
						responsevalues = append(responsevalues, responsevalue)
					}

				} else /*if j != patterncolumn*/ {
					descriptor = sheet.Cell(0, j).String()
					responsedescriptor := descriptor

					otherheaders = append(otherheaders, factororresponse)
					othersubheaders = append(othersubheaders, responsedescriptor)

					cell := sheet.Cell(i, j)

					if cell == nil {

						break
					}

					celltype := cell.Type()

					if celltype == 1 {
						responsevalue, err := cell.Float()
						if err != nil {
							return runs, err
						}
						otherresponsevalues = append(otherresponsevalues, responsevalue)
					} else {
						responsevalue := cell.String()
						otherresponsevalues = append(otherresponsevalues, responsevalue)
					}

				}
			}
			run.Factordescriptors = factordescriptors
			run.Responsedescriptors = responsedescriptors
			run.Setpoints = setpoints
			run.ResponseValues = responsevalues
			run.AdditionalHeaders = otherheaders
			run.AdditionalSubheaders = othersubheaders
			run.AdditionalValues = otherresponsevalues

			runs = append(runs, run)
		}
	}

	return
}

// add func to auto check for Response and factor status based on empty entries implying Response column
func findJMPFactorandResponseColumnsinEmptyDesign(xlsx string) (factorcolumns []int, responsecolumns []int, PatternColumn int) {
	var patternfound bool
	factorcolumns = make([]int, 0)
	responsecolumns = make([]int, 0)

	file, err := spreadsheet.OpenXLSXFromFileName(xlsx)
	if err != nil {
		return
	}
	sheet, err := spreadsheet.Sheet(file, 0)
	if err != nil {
		return
	}
	for j := 0; j < sheet.MaxCol; j++ {

		descriptor := sheet.Cell(0, j).String()
		//	descriptors = append(descriptors,descriptor)
		if strings.ToUpper(descriptor) == "PATTERN" {
			PatternColumn = j
			patternfound = true
		}
	}
	// iterate through every run of the design sheet (row) and if all values for that row == "", the column is interpreted as a response
	for i := 1; i < sheet.MaxRow; i++ {
		//maxfactorcol := 2
		for j := 0; j < sheet.MaxCol; j++ {

			cellstr := sheet.Cell(i, j).String()

			if patternfound && j != PatternColumn && cellstr != "" {
				factorcolumns = append(factorcolumns, j)
			} else if !patternfound && cellstr != "" {
				factorcolumns = append(factorcolumns, j)
			} else if cellstr == "" {

				responsecolumns = append(responsecolumns, j)
			}

		}

	}

	factorcolumns = search.RemoveDuplicateInts(factorcolumns)
	responsecolumns = search.RemoveDuplicateInts(responsecolumns)

	return
}
