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

// Package spreadsheet for interacting with spreadsheets
package spreadsheet

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"

	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wutil"
	"github.com/tealeg/xlsx"
)

// OpenXLSXFromFileName will open an xlsx file from a filename.
func OpenXLSXFromFileName(filename string) (file *xlsx.File, err error) {

	bytes, err := ioutil.ReadFile(filename)

	if err != nil {
		return
	}

	file, err = xlsx.OpenBinary(bytes)

	return
}

// OpenXLSX opens an xlsx file and returns the xlsx.File data structure.
func OpenXLSX(xlsx wtype.File) (file *xlsx.File, err error) {
	fileContents, err := xlsx.ReadAll()
	if err != nil {
		return nil, err
	}

	return OpenXLSXBinary(fileContents)
}

// OpenXLSXBinary parses the contents of an xlsx file into the xlsx.File data structure.
func OpenXLSXBinary(bytes []byte) (file *xlsx.File, err error) {

	file, err = xlsx.OpenBinary(bytes)

	return
}

// Sheet returns the xlsx.Sheet in the xlsx file according to the sheet number.
// An error is returned if an invalid sheet number is specified.
// The sheet object looks like this:
/*
type Sheet struct {
	Name        string
	File        *File
	Rows        []*Row
	Cols        []*Col
	MaxRow      int
	MaxCol      int
	Hidden      bool
	Selected    bool
	SheetViews  []SheetView
	SheetFormat SheetFormat
	AutoFilter  *AutoFilter
}
*/
func Sheet(file *xlsx.File, sheetnum int) (sheet *xlsx.Sheet, err error) {

	if sheetnum < 0 {
		return nil, fmt.Errorf("sheet %d is invalid. The first sheet should be 0, not 1", sheetnum)
	}

	if sheetnum >= len(file.Sheets) {
		var sheets []int
		for key := range file.Sheets {
			sheets = append(sheets, key)
		}
		return nil, fmt.Errorf("sheet %d not found in xlsx file. Found these: %v", sheetnum, sheets)
	}

	return file.Sheets[sheetnum], nil
}

// SheetToCSV returns a matrix of string values for the contents of each cell in the sheet.
func SheetToCSV(sheet *xlsx.Sheet) (records [][]string) {
	for _, row := range sheet.Rows {
		var cellsForRow []string
		for _, cell := range row.Cells {
			cellsForRow = append(cellsForRow, cell.String())
		}
		records = append(records, cellsForRow)
	}
	return
}

// GetDataFromRowCol returns the cell contents at the specified row and column number in an xlsx sheet.
// An error is returned if the column aor row number specified is beyond the range available in the sheet.
// Counting starts from zero. i.e. the cell at the first row  and first column would be called by
// cell, err := GetDataFromRowCol(sheet, 0,0)
func GetDataFromRowCol(sheet *xlsx.Sheet, col int, row int) (cell *xlsx.Cell, err error) {
	if col >= len(sheet.Rows) || col < 0 {
		var cols []int
		for key := range sheet.Rows {
			cols = append(cols, key)
		}
		return nil, fmt.Errorf("column %d not found in xlsx sheet. Found these: %v", col, cols)
	}

	column := sheet.Rows[col]
	if row >= len(column.Cells) || row < 0 {
		var rows []int
		for key := range column.Cells {
			rows = append(rows, key)
		}
		return nil, fmt.Errorf("row %d not found in xlsx sheet. Found these: %v", row, rows)
	}
	return column.Cells[row], nil
}

// GetDataFromCell returns the cell contents at the specified cell position in an xlsx sheet.
// The cellPositionInA1Formatshould be specified according to standard xslx nomenclature.
// i.e. a letter corresponding to column followed by a number corresponding to row (starting from 1).
// i.e. the first cell would be A1.
// An error is returned if no value at the requested well position is found.
func GetDataFromCell(sheet *xlsx.Sheet, cellPositionInA1Format string) (cell *xlsx.Cell, err error) {
	row, col, err := A1FormatToRowColumn(cellPositionInA1Format)
	if err != nil {
		return
	}

	cell, err = GetDataFromRowCol(sheet, col, row)

	if err != nil {
		return cell, fmt.Errorf("error getting data for cell %s", cellPositionInA1Format)
	}

	return cell, nil
}

// GetDataFromCells returns the cell contents at the specified cell positions in an xlsx sheet.
// The cellcoords be specified according to standard xslx nomenclature.
// i.e. a letter corresponding to column followed by a number corresponding to row (starting from 1).
// i.e. the first cell would be A1.
// An error is returned if no value at the requested well position is found.
func GetDataFromCells(sheet *xlsx.Sheet, cellcoords []string) (cells []*xlsx.Cell, err error) {

	cells = make([]*xlsx.Cell, 0)
	for _, a1 := range cellcoords {
		cell, err := GetDataFromCell(sheet, a1)
		if err != nil {
			return cells, err
		}
		cells = append(cells, cell)
	}

	return cells, err
}

// Column returns all cells for a column. The index of the column should be used.
func Column(sheet *xlsx.Sheet, column int) (cells []*xlsx.Cell, err error) {

	if column < 0 {
		return nil, fmt.Errorf("sheet column %d is invalid. The first column should be 0, not 1", column)
	}

	colabcformat := wutil.NumToAlpha(column + 1)

	cellcoords, err := ConvertMinMaxtoArray([]string{(colabcformat + strconv.Itoa(1)), (colabcformat + strconv.Itoa(sheet.MaxRow))})
	if err != nil {
		return cells, err
	}
	cells, err = GetDataFromCells(sheet, cellcoords)

	return cells, err
}

// Row returns all cells for a row. The index of the row should be used.
func Row(sheet *xlsx.Sheet, rowNumber int) (cells []*xlsx.Cell, err error) {
	if rowNumber < 0 {
		return nil, fmt.Errorf("sheet row %d is invalid. The first row should be 0, not 1", rowNumber)
	}

	if rowNumber >= len(sheet.Rows) {
		var rows []int
		for key := range sheet.Rows {
			rows = append(rows, key)
		}

		return nil, fmt.Errorf("row %d not found in xlsx sheet. Found these: %v", rowNumber, rows)
	}

	cellsInrow := sheet.Rows[rowNumber]

	maxColLetter := wutil.NumToAlpha(len(cellsInrow.Cells))

	cellcoords, err := ConvertMinMaxtoArray([]string{("A" + strconv.Itoa(rowNumber+1)), (maxColLetter + strconv.Itoa(rowNumber+1))})
	if err != nil {
		return cells, err
	}

	cells, err = GetDataFromCells(sheet, cellcoords)

	return cells, err
}

// ToHeaderDataMap returns a map of all column headers with corresponding cells.
// useHeaderRow corresponds to the row to use for the headers.
// If this is the first row, it should be set to 0.
func ToHeaderDataMap(sheet *xlsx.Sheet, useHeaderRow int) (headerdatamap map[string][]*xlsx.Cell, err error) {

	headerdatamap = make(map[string][]*xlsx.Cell)

	var columnNumberToHeader = make(map[int]string)

	rows := sheet.Rows

	if useHeaderRow >= len(rows) {
		return nil, fmt.Errorf("specified header row %d is not available in specified sheet", useHeaderRow)
	}

	// get values of first non empty row
	headerRow := rows[useHeaderRow]

	if len(headerRow.Cells) == 0 {
		return nil, fmt.Errorf("specified header row %d contains no values in any cells", useHeaderRow)
	}

	for column, headerCell := range headerRow.Cells {
		if _, found := columnNumberToHeader[column]; !found {
			columnNumberToHeader[column] = headerCell.String()
			if _, found := headerdatamap[headerCell.String()]; !found {
				headerdatamap[headerCell.String()] = []*xlsx.Cell{}
			} else {
				return nil, fmt.Errorf(`duplicate header found "%s"`, headerCell.String())
			}
		} else {
			return nil, fmt.Errorf(`duplicate header found "%d"`, column)
		}
	}

	for i, row := range rows {
		if i != useHeaderRow {
			for columnNumber, cell := range row.Cells {
				header := columnNumberToHeader[columnNumber]
				headerdatamap[header] = append(headerdatamap[header], cell)
			}
		}
	}

	return
}

// A1FormatToRowColumn parses an A1 style excel cell coordinate into ints for row and column for use by plotinum library
// note that 1 is subtracted from the column number in accordance with the go convention of counting from 0
func A1FormatToRowColumn(a1 string) (row, column int, err error) {
	a1 = strings.ToUpper(a1)

	column, err = strconv.Atoi(a1[1:])
	column = column - 1
	if err == nil {
		rowcoord := string(a1[0])
		row := wutil.AlphaToNum(rowcoord) - 1
		return row, column, err
	}
	column, err = strconv.Atoi(a1[2:])
	column = column - 1
	if err == nil {
		rowcoord := a1[0:2]
		row := wutil.AlphaToNum(rowcoord) - 1
		return row, column, err
	}

	column, err = strconv.Atoi(a1[3:])
	column = column - 1
	if err == nil {
		rowcoord := a1[0:3]
		row := wutil.AlphaToNum(rowcoord) - 1
		return row, column, err
	}

	newerr := fmt.Errorf(err.Error() + "more than first three letters of coordinate not int! seems unlikely")
	err = newerr
	return

}

// ConvertMinMaxtoArray converts a pair of cell positions an array of all entries between the pair will be returned (e.g. a1:a12 or a1:e1)
func ConvertMinMaxtoArray(minmax []string) (array []string, err error) {
	if len(minmax) != 2 {
		err = fmt.Errorf("can only make array from a pair of values")
		return
	}

	minrow, mincol, err := A1FormatToRowColumn(minmax[0])
	if err != nil {
		return
	}
	maxrow, maxcol, err := A1FormatToRowColumn(minmax[1])
	if err != nil {
		return
	}

	if minrow == maxrow {
		// fill by column
		array = make([]string, 0)
		for i := mincol; i < maxcol+1; i++ {
			rowstring := wutil.NumToAlpha(minrow + 1)
			colstring := strconv.Itoa(i + 1)

			array = append(array, rowstring+colstring)
		}

	} else if mincol == maxcol {
		// fill by row
		array = make([]string, 0)
		for i := minrow; i < maxrow+1; i++ {
			colstring := strconv.Itoa(mincol + 1)
			rowstring := wutil.NumToAlpha(i + 1)

			array = append(array, rowstring+colstring)
		}
	} else {
		err = fmt.Errorf("either column or row needs to be the same to make an array from two cordinates")
	}
	return

}
