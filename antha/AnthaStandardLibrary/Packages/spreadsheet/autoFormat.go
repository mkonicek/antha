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

import "github.com/tealeg/xlsx"

// AutoFormat applies basic default formatting to a spreadsheet.
// Column widths will be estimated and set to the estimated width for each column of a sheet.
// A thick border will be added to the base of the first row of the each sheet.
func AutoFormat(xlsxfile *xlsx.File) {
	for _, sheet := range xlsxfile.Sheets {
		autoFormatHeader(sheet)
		autoColWidth(sheet)
	}
}

const estimatedCharWidth = 9.5 / 10.0

const thick = "thick"

func autoColWidth(sheet *xlsx.Sheet) {
	for _, row := range sheet.Rows {
		for c, cell := range row.Cells {
			v, err := cell.FormattedValue()
			if err != nil {
				continue
			}
			w := float64(len(v)+1) * estimatedCharWidth
			if w > sheet.Col(c).Width {
				sheet.Col(c).Width = w
			}
		}
	}
}

func autoFormatHeader(sheet *xlsx.Sheet) {
	if len(sheet.Rows) < 1 {
		return
	}
	row := sheet.Rows[0]
	for _, cell := range row.Cells {
		style := cell.GetStyle()
		style.Border.Bottom = thick
	}
}
