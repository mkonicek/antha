// platereaderparse.go
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
	"fmt"
	//"os"
	"strconv"
	"strings"
	"time"

	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/search"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/spreadsheet"
	"github.com/montanaflynn/stats"
	"github.com/tealeg/xlsx"
)

/*
// Inputs (parameters)
var (
	//filename string = "platereaderexamplefiles/simpleabsorbance.xlsx"
	//filename string = "platereaderexamplefiles/platereaderxlsxexample.xlsx"
	filename string = "platereaderexamplefiles/jmresultsinj.xlsx"
	//filename string = "platereaderexamplefiles/AishaAbsEndpointmethodexample.xlsx"
	status string
)

// outputs (Data)
var (
	dataoutput MarsData
)
*/
func ParseMarsXLSXOutput(xlsxname string, sheet int) (dataoutput MarsData, err error) {

	clario, headerrowcount, err := ParseHeadLines(xlsxname, sheet)
	if err != nil {
		return
	}
	//cell := sheet1.Cell(0, 0)

	//dataoutput.User = strings.Split(cell.String(), ": ")[1]

	wellmap, err := ParseWellData(xlsxname, sheet, headerrowcount)
	if err != nil {
		return
	}
	clario.Dataforeachwell = wellmap
	dataoutput = clario
	return
}

func ParseHeadLines(xlsxname string, sheet int) (dataoutput MarsData, headerrowcount int, err error) {
	xlsx, err := spreadsheet.OpenFile(xlsxname)

	//file, err := os.Open(filename)
	if err != nil {
		return
	}

	sheet1 := xlsx.Sheets[sheet]

	for i := 0; i < sheet1.MaxRow; i++ {
		if sheet1.Cell(i, 0).String() == "" {
			headerrowcount = i //+ 1
			// fmt.Println("headerrowcount", headerrowcount)
			break
		}
	}

	maxcell := "a" + strconv.Itoa(headerrowcount)
	// fix this! variable number of IDs leads to range of 8 to 10 header rows
	cellnames, err := spreadsheet.ConvertMinMaxtoArray([]string{"a1", maxcell})
	if err != nil {
		return
	}

	cells, err := spreadsheet.Getdatafromcells(sheet1, cellnames)
	if err != nil {
		return
	}

	dataoutput.Description = cells[len(cells)-1].String()
	for _, cell := range cells {
		if cell.String()[0:4] == "User" {
			dataoutput.User = strings.Split(cell.String(), ": ")[1]
		}
		if cell.String()[0:4] == "Path" {
			dataoutput.Path = strings.Split(cell.String(), ": ")[1]
		}
		if cell.String()[0:7] == "Test ID" {
			id, err := strconv.Atoi(strings.Split(cell.String(), ": ")[1])
			if err != nil {
				return dataoutput, headerrowcount, err
			}
			dataoutput.TestID = id
		}
		if cell.String()[0:9] == "Test Name" {
			dataoutput.Testname = strings.Split(cell.String(), ": ")[1]
		}
		if cell.String()[0:4] == "Date" {
			date := strings.Split(cell.String(), ": ")[1]
			//fmt.Println(date)
			dateparts := strings.Split(date, `/`)
			dateints := make([]int, 0)
			for _, part := range dateparts {
				dateint, err := strconv.Atoi(part)
				if err != nil {
					return dataoutput, headerrowcount, err
				}
				dateints = append(dateints, dateint)
			}

			var godate time.Time

			godate, err = time.Parse("02/01/2006", date)
			if err != nil {
				return
			}
			dataoutput.Date = godate
			//fmt.Println(godate)
			//dataoutput.Date.AddDate(dateints[2], dateints[1], dateints[0])

		}
		if cell.String()[0:4] == "Time" {
			stringtime := strings.Split(cell.String(), ": ")[1]
			if strings.Contains(stringtime, " AM") {
				stringtime = stringtime[0:strings.Index(stringtime, " AM")]
			}
			if strings.Contains(stringtime, " PM") {
				stringtime = stringtime[0:strings.Index(stringtime, " PM")]
				// add something here to correct for 12 hours to add on
			}
			gotime, err := time.Parse("15:4:5", stringtime)
			if err != nil {
				return dataoutput, headerrowcount, err
			}
			dataoutput.Time = gotime
			//fmt.Println(gotime)
		}
		if cell.String()[0:3] == "ID1" {
			dataoutput.ID1 = strings.Split(cell.String(), ": ")[1]
		}
		if cell.String()[0:3] == "ID2" {
			dataoutput.ID2 = strings.Split(cell.String(), ": ")[1]
		}
		if cell.String()[0:3] == "ID3" {
			dataoutput.ID3 = strings.Split(cell.String(), ": ")[1]
		}

	}
	return
}

func ParseWellData(xlsxname string, sheet int, headerrows int) (welldatamap map[string]WellData, err error) {
	welldatamap = make(map[string]WellData)
	var welldata WellData
	var wavelengthstring string
	var wavelength int
	var timestring string
	var timestamp time.Duration
	xlsx, err := spreadsheet.OpenFile(xlsxname)

	//file, err := os.Open(filename)
	if err != nil {
		return
	}

	sheet1 := xlsx.Sheets[sheet]

	//column1 := sheet1.Col(1)
	wellrowstart := 0
	headerrow := headerrows + 2
	timerow := 0
	wavelengthrow := 0

	for i := 0; i < sheet1.MaxRow; i++ {

		cell := sheet1.Cell(i, 0)
		if cell.String() == "A" {
			wellrowstart = i
			//fmt.Println(wellrowstart)
			break
		}

	}
	//if wellrowstart == 15 {
	wavelengths := make([]int, 0)

	times := make([]time.Duration, 0)

	//find special rows
	// wavelength or time rows?
	//fmt.Println(headerrows, headerrow)
	if wellrowstart-headerrow > 0 {
		for i := 0; i < wellrowstart-headerrow; i++ {

			rowabove := spreadsheet.Getdatafromrowcol(sheet1, wellrowstart-(i+1), 2).String()
			if strings.Contains(rowabove, "Time") {
				timerow = wellrowstart - (i + 1)
				// fmt.Println("timerow:", timerow)
			} else if strings.Contains(rowabove, "Wavelength") {
				wavelengthrow = wellrowstart - (i + 1)
			}

		}
	}
	// check other row names in case the row labels are not in order (this can happen)
	for i := wellrowstart; i < sheet1.MaxRow; i++ {

		rowname := spreadsheet.Getdatafromrowcol(sheet1, i, 2).String()
		// fmt.Println("rowname", rowname)
		if strings.Contains(rowname, "Time") {
			timerow = i
			// fmt.Println("timerow new:", timerow)
		} else if strings.Contains(rowname, "Wavelength") {
			wavelengthrow = i
		}

	}

	// find special columns
	tempcolumn := 0
	injectionvoumecolumn := 0

	for m := 3; m < sheet1.MaxCol; m++ {

		columnheader := spreadsheet.Getdatafromrowcol(sheet1, headerrow, m).String()

		if strings.Contains(columnheader, "Temperature") {
			tempcolumn = m
		}
		if strings.Contains(columnheader, "Volume") {
			injectionvoumecolumn = m
		}

	}

	for j := wellrowstart; j < sheet1.MaxRow; j++ {

		if j != timerow && j != wavelengthrow {

			welldata.Name = spreadsheet.Getdatafromrowcol(sheet1, j, 2).String()
			welldata.Well = spreadsheet.Getdatafromrowcol(sheet1, (j), 0).String() + spreadsheet.Getdatafromrowcol(sheet1, j, 1).String()

			//if j == wellrowstart {
			for k := 3; k < sheet1.MaxCol; k++ {
				if k != tempcolumn && k != injectionvoumecolumn {
					welldata.ReadingType = spreadsheet.Getdatafromrowcol(sheet1, headerrow, k).String()

					if wavelengthrow != 0 {
						wavelength, err := spreadsheet.Getdatafromrowcol(sheet1, wavelengthrow, k).Int()
						if err != nil {
							return welldatamap, err
						}
						if len(wavelengths) == 0 {
							wavelengths = append(wavelengths, wavelength)
						} else if search.Contains(wavelengths, wavelength) == false {
							wavelengths = append(wavelengths, wavelength)
						}
					}
				}
			}
			//	}

			for m := 3; m < sheet1.MaxCol; m++ {

				if wavelengthrow != 0 {
					wavelength, err := spreadsheet.Getdatafromrowcol(sheet1, wavelengthrow, m).Int()
					if err != nil {
						return welldatamap, err
					}
					if wavelengths[0] != wavelength {
						break
					} else {
						if timerow != 0 {
							timelabel := spreadsheet.Getdatafromrowcol(sheet1, timerow, 2).String()
							if strings.Contains(timelabel, "[s]") {
								timeplusseconds := spreadsheet.Getdatafromrowcol(sheet1, timerow, m).String() + "s"
								gotime, err := ParseTime(timeplusseconds)
								// fmt.Println("added s", timeplusseconds)

								if err != nil {
									return welldatamap, err
								}
								times = append(times, gotime)
								fmt.Println(times)
							} else {
								gotime, err := ParseTime(spreadsheet.Getdatafromrowcol(sheet1, timerow, m).String())
								if err != nil {
									return welldatamap, err
								}
								times = append(times, gotime)
								fmt.Println(times)
							}

						}
					}
				}
			}
			var measurement PRMeasurement

			var Measurements = make([]PRMeasurement, 0)

			maxcol := sheet1.MaxCol

			// this should probably be deleted to avoid cases where last column is not temp
			/*
				tempheadercheck := spreadsheet.Getdatafromrowcol(sheet1, headerrow, sheet1.MaxCol-1).String()

				if strings.Contains(tempheadercheck, "Temperature") {
					maxcol = sheet1.MaxCol - 1
				}
			*/
			for m := 3; m < maxcol; m++ {

				//check header
				header := spreadsheet.Getdatafromrowcol(sheet1, headerrow, m).String()

				// the measurement itself (if not a special column - e.g. volume injection or temp)
				if strings.Contains(header, "Temperature") == false && strings.Contains(header, "Volume") == false {
					measurement.Reading, err = spreadsheet.Getdatafromrowcol(sheet1, j, m).Float()
				}

				// add logic to check column heading
				// add similar for volume (injection)
				if strings.Contains(header, "Temperature") {
					measurement.Temp, err = spreadsheet.Getdatafromrowcol(sheet1, j, tempcolumn).Float()
				} else if strings.Contains(header, "Volume") {
					welldata.InjectionVolume, err = spreadsheet.Getdatafromrowcol(sheet1, j, injectionvoumecolumn).Float()
					welldata.Injected = true
				} else {

					// add time row and wavelength row calculators
					if timerow != 0 {
						//gotime, err := ParseTime(spreadsheet.Getdatafromrowcol(sheet1, timerow, m).String())
						timelabel := spreadsheet.Getdatafromrowcol(sheet1, timerow, 2).String()

						if strings.Contains(timelabel, "[s]") && spreadsheet.Getdatafromrowcol(sheet1, timerow, m).String() != "" {
							timestring = spreadsheet.Getdatafromrowcol(sheet1, timerow, m).String() + "s"

							// fmt.Println("added s", timestring)

							if err != nil {
								return welldatamap, err
							}
						} else if spreadsheet.Getdatafromrowcol(sheet1, timerow, m).String() != "" {
							timestring = spreadsheet.Getdatafromrowcol(sheet1, timerow, m).String()
						}
						timestamp, err = ParseTime(timestring)
						if err != nil {
							fmt.Println(timestring, timestamp)
							return welldatamap, err
						}

						measurement.Timestamp = timestamp
						// fmt.Println("timestamp:", timestamp)
					}
					// need to have some different options here for handling different types
					// Ex Spectrum, Absorbance reading etc.. Abs spectrum, ex spectrum
					welldata.ReadingType = spreadsheet.Getdatafromrowcol(sheet1, headerrow, m).String()

					parsedatatype := strings.Split(welldata.ReadingType, `(`)
					// fmt.Println("parsed data", parsedatatype)
					parsedatatype = strings.Split(parsedatatype[1], `)`)

					if strings.Contains(header, "Temperature") == false && strings.Contains(header, "Volume") == false {

						// handle case of absorbance (may need to add others.. if contains Ex, Ex else number = abs
						if strings.Contains(parsedatatype[0], "A-") {

							wavelengthstring = parsedatatype[0][strings.Index(parsedatatype[0], `-`)+1:]

							wavelength, err = strconv.Atoi(wavelengthstring)
							if err != nil {
								return welldatamap, err
							}
							measurement.RWavelength = wavelength
							measurement.EWavelength = wavelength
						} else if strings.Contains(welldata.ReadingType, "Em Spectrum") {
							if wavelengthrow != 0 {
								emWavelength, err := spreadsheet.Getdatafromrowcol(sheet1, wavelengthrow, m).Int()
								if err != nil {
									return welldatamap, err
								}
								measurement.RWavelength = emWavelength
							}
						} else if strings.Contains(welldata.ReadingType, "Ex Spectrum") {
							if wavelengthrow != 0 {
								exWavelength, err := spreadsheet.Getdatafromrowcol(sheet1, wavelengthrow, m).Int()
								if err != nil {
									return welldatamap, err
								}
								measurement.EWavelength = exWavelength
							}
						} else if strings.Contains(welldata.ReadingType, "Abs Spectrum") {
							if wavelengthrow == 0 {
								wavelengthstring = parsedatatype[0]

								wavelength, err = strconv.Atoi(wavelengthstring)
								if err != nil {
									return welldatamap, err
								}
								measurement.RWavelength = wavelength
								measurement.EWavelength = wavelength
							} else {
								wavelength, err := spreadsheet.Getdatafromrowcol(sheet1, wavelengthrow, m).Int()
								if err != nil {
									return welldatamap, err
								}
								measurement.RWavelength = wavelength
								measurement.EWavelength = wavelength
							}

						} else if strings.Contains(welldata.ReadingType, "(484-8/517.5-9 4)") {
							measurement.RWavelength = 517
							measurement.EWavelength = 484
						} else if strings.Contains(welldata.ReadingType, "(484-8 2)") {
							measurement.RWavelength = 517
							measurement.EWavelength = 484
						} else if strings.Contains(welldata.ReadingType, "(600 2)") {
							measurement.RWavelength = 600
							measurement.EWavelength = 600
							/*} else if strings.Contains(welldata.ReadingType, "(484-8/517.5-9 3)") {
							measurement.RWavelength = 517
							measurement.EWavelength = 484
							*/
						} else if strings.Contains(welldata.ReadingType, "(600 1)") {
							measurement.RWavelength = 600
							measurement.EWavelength = 600
						} else if HeaderContainsWavelength(sheet1, headerrow, j) {
							if wavelengthrow == 0 && HeaderContainsWavelength(sheet1, headerrow, j) {
								_, wavelength, err := HeaderWavelength(sheet1, headerrow, j)

								if err != nil {
									return welldatamap, err
								}
								measurement.RWavelength = wavelength
								measurement.EWavelength = wavelength

								/*wavelengthstring = parsedatatype[0]

								wavelength, err = strconv.Atoi(wavelengthstring)
								if err != nil {
									return welldatamap, err
								}
								measurement.RWavelength = wavelength
								measurement.EWavelength = wavelength
								*/
							} else {
								wavelength, err := spreadsheet.Getdatafromrowcol(sheet1, wavelengthrow, m).Int()
								if err != nil {
									return welldatamap, err
								}
								measurement.RWavelength = wavelength
								measurement.EWavelength = wavelength
							}
						}
					}
					Measurements = append(Measurements, measurement)
				}
			}

			times = make([]time.Duration, 0)
			var output PROutput
			var set PRMeasurementSet
			set = Measurements

			output.Readings = make([]PRMeasurementSet, 1)
			output.Readings[0] = set
			welldata.Data = output
			welldatamap[welldata.Well] = welldata
			//fmt.Println(welldata)
			//fmt.Println(output)

		}
	}

	return
}

func ParseTime(timestring string) (gotime time.Duration, err error) {

	fields := strings.Fields(timestring)

	newfields := make([]string, 0)
	for _, field := range fields {
		if strings.Contains(field, "min") {

			newfields = append(newfields, "m")
		} else {
			newfields = append(newfields, field)
		}
	}

	parsethis := strings.Join(newfields, "")

	gotime, err = time.ParseDuration(parsethis)

	return
}

func HeaderContainsWavelength(sheet *xlsx.Sheet, cellrow, cellcolumn int) (yesno bool) {
	headercell := spreadsheet.Getdatafromrowcol(sheet, cellrow, cellcolumn).String()
	//fmt.Println(headercell)

	if strings.Contains(headercell, "(") && strings.Contains(headercell, ")") {
		start := strings.Index(headercell, "(")
		finish := strings.Index(headercell, ")")

		isthisanumber := headercell[start+1 : finish]
		//	fmt.Println(isthisanumber)
		_, err := strconv.Atoi(isthisanumber)

		if err == nil {
			yesno = true
		}
	}

	return
}

func HeaderWavelength(sheet *xlsx.Sheet, cellrow, cellcolumn int) (yesno bool, number int, err error) {
	headercell := spreadsheet.Getdatafromrowcol(sheet, cellrow, cellcolumn).String()
	//fmt.Println(headercell)

	if strings.Contains(headercell, "(") && strings.Contains(headercell, ")") {
		start := strings.Index(headercell, "(")
		finish := strings.Index(headercell, ")")

		isthisanumber := headercell[start+1 : finish]
		//fmt.Println(isthisanumber)
		number, err = strconv.Atoi(isthisanumber)

		if err == nil {
			yesno = true
		}
	} else {
		err = fmt.Errorf("no (  ) found in header")
	}

	return
}

func (data MarsData) TimeCourse(wellname string, exWavelength int, emWavelength int) (xaxis []time.Duration, yaxis []float64) {
	xaxis = make([]time.Duration, 0)
	yaxis = make([]float64, 0)
	for _, measurement := range data.Dataforeachwell[wellname].Data.Readings[0] {

		if measurement.EWavelength == exWavelength && measurement.RWavelength == emWavelength {

			xaxis = append(xaxis, measurement.Timestamp)
			yaxis = append(yaxis, measurement.Reading)

		}
	}

	return
}

// readingtypekeyword added in case mars used to process data in advance. Example keywords : Raw Data, Em Spectrum, Abs Spectrum, Blank Corrected, Average or "" to capture all
func (data MarsData) AbsScanData(wellname string, readingtypekeyword string) (wavelengths []int, Readings []float64) {
	wavelengths = make([]int, 0)
	Readings = make([]float64, 0)
	for _, measurement := range data.Dataforeachwell[wellname].Data.Readings[0] {

		if strings.Contains(data.Dataforeachwell[wellname].ReadingType, readingtypekeyword) {

			wavelengths = append(wavelengths, measurement.RWavelength)
			Readings = append(Readings, measurement.Reading)

		}
	}

	return
}

func (data MarsData) EMScanData(wellname string, exWavelength int, readingtypekeyword string) (wavelengths []int, Readings []float64) {
	wavelengths = make([]int, 0)
	Readings = make([]float64, 0)
	for _, measurement := range data.Dataforeachwell[wellname].Data.Readings[0] {

		if measurement.EWavelength == exWavelength && strings.Contains(data.Dataforeachwell[wellname].ReadingType, readingtypekeyword) {

			wavelengths = append(wavelengths, measurement.RWavelength)
			Readings = append(Readings, measurement.Reading)

		}

	}

	return
}

func (data MarsData) EXScanData(wellname string, emWavelength int, readingtypekeyword string) (wavelengths []int, Readings []float64) {
	wavelengths = make([]int, 0)
	Readings = make([]float64, 0)
	for _, measurement := range data.Dataforeachwell[wellname].Data.Readings[0] {

		if measurement.RWavelength == emWavelength && strings.Contains(data.Dataforeachwell[wellname].ReadingType, readingtypekeyword) {

			wavelengths = append(wavelengths, measurement.EWavelength)
			Readings = append(Readings, measurement.Reading)

		}
	}

	return
}

func (data MarsData) WelltoDataMap() map[string]WellData {
	return data.Dataforeachwell
}

func (data MarsData) Readings(wellname string) []PRMeasurement {
	return data.Dataforeachwell[wellname].Data.Readings[0]
}

func (data MarsData) ReadingsThat(wellname string, emexortime int, fieldvalue interface{}) ([]PRMeasurement, error) {
	newset := make([]PRMeasurement, 0)
	for _, measurement := range data.Dataforeachwell[wellname].Data.Readings[0] {

		if emexortime == 0 {
			if str, ok := fieldvalue.(string); ok {

				gotime, err := time.ParseDuration(str)
				if err != nil {
					return newset, err
				}
				if measurement.Timestamp == gotime {
					newset = append(newset, measurement)
				}
			}
		} else if emexortime == 1 && measurement.RWavelength == fieldvalue {
			newset = append(newset, measurement)
		} else if emexortime == 2 && measurement.EWavelength == fieldvalue {
			newset = append(newset, measurement)
		}
	}

	return newset, nil
}

func (data MarsData) ReadingsAsFloats(wellname string, emexortime int, fieldvalue interface{}) (readings []float64, readingtypes []string, err error) {
	readings = make([]float64, 0)
	readingtypes = make([]string, 0)
	for _, measurement := range data.Dataforeachwell[wellname].Data.Readings[0] {

		if emexortime == 0 {
			if str, ok := fieldvalue.(string); ok {

				gotime, err := time.ParseDuration(str)
				if err != nil {
					return readings, readingtypes, err
				}
				if measurement.Timestamp == gotime {
					readings = append(readings, measurement.Reading)
					readingtypes = append(readingtypes, data.Dataforeachwell[wellname].ReadingType)
				}
			}
		} else if emexortime == 1 && measurement.RWavelength == fieldvalue {
			readings = append(readings, measurement.Reading)
			readingtypes = append(readingtypes, data.Dataforeachwell[wellname].ReadingType)
		} else if emexortime == 2 && measurement.EWavelength == fieldvalue {
			readings = append(readings, measurement.Reading)
			readingtypes = append(readingtypes, data.Dataforeachwell[wellname].ReadingType)
		}
	}

	return
}

func (data MarsData) ReadingsAsAverage(wellname string, emexortime int, fieldvalue interface{}, readingtypekeyword string) (average float64, err error) {
	readings := make([]float64, 0)
	readingtypes := make([]string, 0)
	readingsforaverage := make([]float64, 0)
	for _, measurement := range data.Dataforeachwell[wellname].Data.Readings[0] {

		if emexortime == 0 {
			if str, ok := fieldvalue.(string); ok {

				gotime, err := time.ParseDuration(str)
				if err != nil {
					return average, err
				}
				if measurement.Timestamp == gotime {
					readings = append(readings, measurement.Reading)
					readingtypes = append(readingtypes, data.Dataforeachwell[wellname].ReadingType)
				}
			}
		} else if emexortime == 1 && measurement.RWavelength == fieldvalue {
			readings = append(readings, measurement.Reading)
			readingtypes = append(readingtypes, data.Dataforeachwell[wellname].ReadingType)
		} else if emexortime == 2 && measurement.EWavelength == fieldvalue {
			readings = append(readings, measurement.Reading)
			readingtypes = append(readingtypes, data.Dataforeachwell[wellname].ReadingType)
		}
	}

	for i, readingtype := range readingtypes {
		if strings.Contains(readingtype, readingtypekeyword) {
			readingsforaverage = append(readingsforaverage, readings[i])
		}
	}

	average, err = stats.Mean(readingsforaverage)

	return
}

func (data MarsData) AbsorbanceReading(wellname string, wavelength int, readingtypekeyword string) (average float64, err error) {
	readings := make([]float64, 0)
	readingtypes := make([]string, 0)
	readingsforaverage := make([]float64, 0)
	for _, measurement := range data.Dataforeachwell[wellname].Data.Readings[0] {

		if measurement.RWavelength == wavelength {
			readings = append(readings, measurement.Reading)
			readingtypes = append(readingtypes, data.Dataforeachwell[wellname].ReadingType)

		}

		for i, readingtype := range readingtypes {
			if strings.Contains(readingtype, readingtypekeyword) {
				readingsforaverage = append(readingsforaverage, readings[i])
			}
		}
	}
	average, err = stats.Mean(readingsforaverage)

	return
}

func (data MarsData) FindOptimalWavelength(wellname string, blankname string, readingtypekeyword string) (wavelength int) {

	//differences := make([]float64, 0)
	biggestdifferenceindex := 0
	biggestdifference := 0.0

	wavelengths, readings := data.AbsScanData(wellname, readingtypekeyword)
	blankwavelengths, blankreadings := data.AbsScanData(blankname, readingtypekeyword)

	for i, reading := range readings {

		difference := reading - blankreadings[i]

		if difference > biggestdifference && wavelengths[i] == blankwavelengths[i] {
			biggestdifferenceindex = i
		}

	}

	wavelength = wavelengths[biggestdifferenceindex]

	return
}

func (data MarsData) BlankCorrect(wellnames []string, blanknames []string, wavelength int, readingtypekeyword string) (blankcorrectedaverage float64, err error) {
	readings := make([]float64, 0)
	readingtypes := make([]string, 0)
	readingsforaverage := make([]float64, 0)

	for _, wellname := range blanknames {

		for _, measurement := range data.Dataforeachwell[wellname].Data.Readings[0] {

			if measurement.RWavelength == wavelength {
				readings = append(readings, measurement.Reading)
				readingtypes = append(readingtypes, data.Dataforeachwell[wellname].ReadingType)

			}

			for i, readingtype := range readingtypes {
				if strings.Contains(readingtype, readingtypekeyword) {
					readingsforaverage = append(readingsforaverage, readings[i])
				}
			}
		}
	}

	blankaverage, err := stats.Mean(readingsforaverage)

	readings = make([]float64, 0)
	readingtypes = make([]string, 0)
	readingsforaverage = make([]float64, 0)

	for _, wellname := range wellnames {

		for _, measurement := range data.Dataforeachwell[wellname].Data.Readings[0] {

			if measurement.RWavelength == wavelength {
				readings = append(readings, measurement.Reading)
				readingtypes = append(readingtypes, data.Dataforeachwell[wellname].ReadingType)

			}

			for i, readingtype := range readingtypes {
				if strings.Contains(readingtype, readingtypekeyword) {
					readingsforaverage = append(readingsforaverage, readings[i])
				}
			}
		}

	}
	average, err := stats.Mean(readingsforaverage)

	blankcorrectedaverage = average - blankaverage

	return
}

const (
	TIME = iota
	EMWAVELENGTH
	EXWAVELENGTH
)

/*
func main() {

	//clario, err := ParseHeadLines(filename)
	clario, err := ParseMarsXLSXOutput(filename)
	if err != nil {
		status = err.Error()
	} else {
		status = "no problemo"
	}
	//cell := sheet1.Cell(0, 0)

	//dataoutput.User = strings.Split(cell.String(), ": ")[1]
	/*
		wellmap, err := ParseWellData(filename)
		if err != nil {
			status = err.Error()
		} else {
			status = "no problemo"
		}
*/

//readings, err := clario.ReadingsThat("A1", EMWAVELENGTH, 600)
//	x, y := clario.TimeCourse("A1", 340, 340)
//	fmt.Println(status, clario, x, y /*clario.Readings("A1") ,, readings*/)

//}
//*/
type MarsData struct {
	User            string
	Path            string
	TestID          int
	Testname        string
	Date            time.Time
	Time            time.Time
	ID1             string
	ID2             string
	ID3             string
	Description     string
	Dataforeachwell map[string]WellData
}

type WellData struct {
	Well            string // in a1 format
	Name            string
	Data            PROutput
	ReadingType     string
	Injected        bool
	InjectionVolume float64
}

// from antha/microArch/driver/platereader
type PROutput struct {
	Readings []PRMeasurementSet
}

type PRMeasurementSet []PRMeasurement

type PRMeasurement struct {
	EWavelength int           //	excitation wavelength
	RWavelength int           //	emission wavelength
	Reading     float64       //int           // 	value read
	Xoff        int           //	position - x, relative to well centre
	Yoff        int           //	position - y, relative to well centre
	Zoff        int           // 	position - z, relative to well centre
	Timestamp   time.Duration // instant measurement was taken
	Temp        float64       //int       //   temperature
	O2          int           // o2 conc when measurement was taken
	CO2         int           // co2 conc when measurement was taken
}
