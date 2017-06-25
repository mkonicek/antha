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
	"io/ioutil"
	"strconv"
	"strings"
	"time"

	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/search"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/spreadsheet"
	"github.com/antha-lang/antha/antha/anthalib/wutil"

	"github.com/montanaflynn/stats"
	"github.com/tealeg/xlsx"
)

// minimal interface to support existing antha elements which use plate reader data (AddPlateReder_Results and AddGFPODPlateReaderResults)
type PlateReaderData interface {
	BlankCorrect(wellnames []string, blanknames []string, wavelength int, readingtypekeyword string) (blankcorrectedaverage float64, err error)
	ReadingsAsAverage(wellname string, emexortime int, fieldvalue interface{}, readingtypekeyword string) (average float64, err error)
	FindOptimalWavelength(wellname string, blankname string, readingtypekeyword string) (wavelength int, err error)
	TimeCourse(wellname string, exWavelength int, emWavelength int, scriptnumber int) (xaxis []time.Duration, yaxis []float64, err error)
}

///////
type AbsorbanceData interface {
	BlankCorrect(wellnames []string, blanknames []string, wavelength int, readingtypekeyword string) (blankcorrectedaverage float64, err error)
	AbsorbanceReading(wellname string, wavelength int) (average float64, err error)
	FindOptimalWavelength(wellname string, blankname string, readingtypekeyword string) (wavelength int, err error)
	TimeCourse(wellname string, exWavelength int, emWavelength int, scriptnumber int) (xaxis []time.Duration, yaxis []float64, err error)
}

type TimeCourseData interface {
	TimeCourse(wellname string, exWavelength int, emWavelength int, scriptnumber int) (xaxis []time.Duration, yaxis []float64, err error)
}

//////////////////////////

// parse mars data from excel filename
func ParseMarsXLSXOutput(xlsxname string, sheet int) (dataoutput MarsData, err error) {

	bytes, err := ioutil.ReadFile(xlsxname)

	if err != nil {
		return
	}

	clario, headerrowcount, err := parseHeadLines(bytes, sheet)
	if err != nil {
		return
	}

	wellmap, err := parseWellData(bytes, sheet, headerrowcount)
	if err != nil {
		return
	}
	clario.Dataforeachwell = wellmap
	dataoutput = clario
	return
}

// parse mars data from excel filename
func ParseMarsXLSXBinary(xlsxContents []byte, sheet int) (dataoutput MarsData, err error) {

	clario, headerrowcount, err := parseHeadLines(xlsxContents, sheet)
	if err != nil {
		return
	}

	wellmap, err := parseWellData(xlsxContents, sheet, headerrowcount)
	if err != nil {
		return
	}
	clario.Dataforeachwell = wellmap
	dataoutput = clario
	return
}

func parseHeadLines(xlsxBinary []byte, sheet int) (dataoutput MarsData, headerrowcount int, err error) {
	xlsx, err := spreadsheet.OpenBinary(xlsxBinary)

	if err != nil {
		return
	}

	sheet1 := xlsx.Sheets[sheet]

	for i := 0; i < sheet1.MaxRow; i++ {
		str, err := sheet1.Cell(i, 0).String()

		if err != nil {
			return dataoutput, headerrowcount, err
		}

		if str == "" {
			headerrowcount = i //+ 1
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

	dataoutput.Description, err = cells[len(cells)-1].String()
	if err != nil {
		return
	}

	for _, cell := range cells {

		cellstr, err := cell.String()

		if err != nil {
			return dataoutput, headerrowcount, err
		}

		if cellstr[0:4] == "User" {
			dataoutput.User = strings.Split(cellstr, ": ")[1]
		}
		if cellstr[0:4] == "Path" {
			dataoutput.Path = strings.Split(cellstr, ": ")[1]
		}
		if cellstr[0:7] == "Test ID" {
			id, err := strconv.Atoi(strings.Split(cellstr, ": ")[1])
			if err != nil {
				return dataoutput, headerrowcount, err
			}
			dataoutput.TestID = id
		}
		if cellstr[0:9] == "Test Name" {
			dataoutput.Testname = strings.Split(cellstr, ": ")[1]
		}
		if cellstr[0:4] == "Date" {
			date := strings.Split(cellstr, ": ")[1]
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
				return dataoutput, headerrowcount, err
			}
			dataoutput.Date = godate

		}
		if cellstr[0:4] == "Time" {
			stringtime := strings.Split(cellstr, ": ")[1]
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
		}
		if cellstr[0:3] == "ID1" {
			dataoutput.ID1 = strings.Split(cellstr, ": ")[1]
		}
		if cellstr[0:3] == "ID2" {
			dataoutput.ID2 = strings.Split(cellstr, ": ")[1]
		}
		if cellstr[0:3] == "ID3" {
			dataoutput.ID3 = strings.Split(cellstr, ": ")[1]
		}

	}
	return
}

func parseWellData(xlsxBinary []byte, sheet int, headerrows int) (welldatamap map[string]WellData, err error) {
	welldatamap = make(map[string]WellData)
	var welldata WellData
	var wavelengthstring string
	var wavelength int
	var timestring string
	var timestamp time.Duration
	xlsx, err := spreadsheet.OpenBinary(xlsxBinary)

	if err != nil {
		return
	}

	sheet1 := xlsx.Sheets[sheet]
	if sheet1.MaxRow == 0 {
		return welldatamap, fmt.Errorf("No well data found in Mars data file")
	}
	wellrowstart := 0
	headerrow := headerrows + 2
	timerow := 0
	wavelengthrow := 0

	for i := 0; i < sheet1.MaxRow; i++ {

		cell := sheet1.Cell(i, 0)

		cellstr, err := cell.String()

		if err != nil {
			return welldatamap, err
		}

		if cellstr == "A" {
			wellrowstart = i
			break
		}

	}
	wavelengths := make([]int, 0)

	times := make([]time.Duration, 0)

	if wellrowstart-headerrow > 0 {
		for i := 0; i < wellrowstart-headerrow; i++ {

			rowabove, err := spreadsheet.Getdatafromrowcol(sheet1, wellrowstart-(i+1), 2).String()
			if err != nil {
				return welldatamap, err
			}
			if strings.Contains(rowabove, "Time") {
				timerow = wellrowstart - (i + 1)
			} else if strings.Contains(rowabove, "Wavelength") {
				wavelengthrow = wellrowstart - (i + 1)
			}

		}
	}
	// check other row names in case the row labels are not in order (this can happen)
	for i := wellrowstart; i < sheet1.MaxRow; i++ {

		rowname, err := spreadsheet.Getdatafromrowcol(sheet1, i, 2).String()
		if err != nil {
			return welldatamap, err
		}

		if strings.Contains(rowname, "Time") {
			timerow = i
		} else if strings.Contains(rowname, "Wavelength") {
			wavelengthrow = i
		}

	}

	// find special columns
	tempcolumn := 0
	injectionvoumecolumn := 0

	for m := 3; m < sheet1.MaxCol; m++ {

		columnheader, err := spreadsheet.Getdatafromrowcol(sheet1, headerrow, m).String()
		if err != nil {
			return welldatamap, err
		}

		if strings.Contains(columnheader, "Temperature") {
			tempcolumn = m
		}
		if strings.Contains(columnheader, "Volume") {
			injectionvoumecolumn = m
		}

	}

	for j := wellrowstart; j < sheet1.MaxRow; j++ {

		if j != timerow && j != wavelengthrow {

			welldata.Name, err = spreadsheet.Getdatafromrowcol(sheet1, j, 2).String()
			if err != nil {
				return welldatamap, err
			}

			part1, err := spreadsheet.Getdatafromrowcol(sheet1, (j), 0).String()
			if err != nil {
				return welldatamap, err
			}
			part2, err := spreadsheet.Getdatafromrowcol(sheet1, j, 1).String()
			if err != nil {
				return welldatamap, err
			}

			welldata.Well = part1 + part2

			for k := 3; k < sheet1.MaxCol; k++ {
				if k != tempcolumn && k != injectionvoumecolumn {

					readingtype, err := spreadsheet.Getdatafromrowcol(sheet1, headerrow, k).String()

					if err != nil {
						return welldatamap, err
					}

					welldata.ReadingType = readingtype

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
							timelabel, err := spreadsheet.Getdatafromrowcol(sheet1, timerow, 2).String()

							if err != nil {
								return welldatamap, err
							}

							if strings.Contains(timelabel, "[s]") {
								time, err := spreadsheet.Getdatafromrowcol(sheet1, timerow, m).String()

								if err != nil {
									return welldatamap, err
								}
								timeplusseconds := time + "s"
								gotime, err := ParseTime(timeplusseconds)

								if err != nil {
									return welldatamap, err
								}
								times = append(times, gotime)
							} else {
								time, err := spreadsheet.Getdatafromrowcol(sheet1, timerow, m).String()
								if err != nil {
									return welldatamap, err
								}
								gotime, err := ParseTime(time)
								if err != nil {
									return welldatamap, err
								}
								times = append(times, gotime)
							}

						}
					}
				}
			}
			var measurement PRMeasurement

			var Measurements = make([]PRMeasurement, 0)

			maxcol := sheet1.MaxCol

			for m := 3; m < maxcol; m++ {

				//check header
				header, err := spreadsheet.Getdatafromrowcol(sheet1, headerrow, m).String()
				if err != nil {
					return welldatamap, err
				}
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
						timelabel, err := spreadsheet.Getdatafromrowcol(sheet1, timerow, 2).String()
						if err != nil {
							return welldatamap, err
						}

						timecellcontents, err := spreadsheet.Getdatafromrowcol(sheet1, timerow, m).String()

						if err != nil {
							return welldatamap, err
						}

						if strings.Contains(timelabel, "[s]") && timecellcontents != "" {
							timestring = timecellcontents + "s"

							if err != nil {
								return welldatamap, err
							}
						} else if timecellcontents != "" {
							timestring = timecellcontents
						}
						timestamp, err = ParseTime(timestring)
						if err != nil {
							fmt.Println(timestring, timestamp)
							return welldatamap, err
						}

						measurement.Timestamp = timestamp
					}
					// need to have some different options here for handling different types
					// Ex Spectrum, Absorbance reading etc.. Abs spectrum, ex spectrum

					welldata.ReadingType, err = spreadsheet.Getdatafromrowcol(sheet1, headerrow, m).String()
					if err != nil {
						return welldatamap, err
					}
					parsedatatype := strings.Split(welldata.ReadingType, `(`)

					parsedatatype = strings.Split(parsedatatype[1], `)`)

					if strings.Contains(header, "Temperature") == false && strings.Contains(header, "Volume") == false {

						// handle case of absorbance (may need to add others.. if contains Ex, Ex else number = abs
						ex, exband, em, emband, scriptposition, err := parseBracketedColumnHeader(welldata.ReadingType)

						if err == nil {
							measurement.RWavelength = em
							measurement.EWavelength = ex
							measurement.EBand = exband
							measurement.RBand = emband
							measurement.Script = scriptposition

						} else if strings.Contains(parsedatatype[0], "A-") {

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

						} else if headerContainsWavelength(sheet1, headerrow, j) {
							if wavelengthrow == 0 && headerContainsWavelength(sheet1, headerrow, j) {
								_, wavelength, err := headerWavelength(sheet1, headerrow, j)

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
		}
	}

	return
}

func bracketed(header string) bool {

	header = strings.TrimSpace(header)

	if strings.Contains(header, "(") && strings.Contains(header, ")") {
		return true
	}
	return false
}

func parseBracketedColumnHeader(header string) (ex int, exband int, em int, emband int, scriptposition int, err error) {
	if bracketed(header) {
		start := strings.Index(header, "(")
		header = header[start+1:]
		header = strings.TrimRight(header, ")")

		if integer, er := strconv.Atoi(header); er == nil {
			ex = integer
			em = integer
			scriptposition = 0
			return
		} else if len(strings.Fields(header)) == 2 {

			fields := strings.Fields(header)

			// handle wavelength part
			if integer, er := strconv.Atoi(fields[0]); er == nil { // single wavelength
				ex = integer
				em = integer
			} else if strings.Count(fields[0], "/") == 1 { // ex and em

				subfields := strings.Split(fields[0], "/")
				if len(subfields) == 2 {

					// excitation
					if integer, er := strconv.Atoi(subfields[0]); er == nil {
						ex = integer
					} else if strings.Count(subfields[0], "-") == 1 {
						exfields := strings.Split(subfields[0], "-")

						// excitation
						if integer, er := strconv.Atoi(exfields[0]); er == nil {
							ex = integer
						} else if f, er := strconv.ParseFloat(exfields[0], 64); er == nil {
							ex = wutil.RoundInt(f)
						}

						// band
						if integer, er := strconv.Atoi(exfields[1]); er == nil {
							exband = integer
						}
					}
					// and emission
					if integer, er := strconv.Atoi(subfields[1]); er == nil {
						em = integer
					} else if strings.Count(subfields[1], "-") == 1 {
						emfields := strings.Split(subfields[1], "-")

						// emission
						if integer, er := strconv.Atoi(emfields[0]); er == nil {
							em = integer
						} else if f, er := strconv.ParseFloat(emfields[0], 64); er == nil {
							em = wutil.RoundInt(f)
						}

						// band
						if integer, er := strconv.Atoi(emfields[1]); er == nil {
							emband = integer
						}
					}
				}

			} else {
				err = fmt.Errorf("Unknown header type, %s ,found in Mars data file, problem with %s", header, fields[0])
				return
			}

			// handle scriptnumber part
			if integer, er := strconv.Atoi(fields[1]); er == nil {
				scriptposition = integer
				return
			} else {
				err = fmt.Errorf("Unknown header type, %s ,found in Mars data file, problem with %s", header, fields[1])
				return
			}

		}
	}
	err = fmt.Errorf("Error with header %s found in Mars data file", header)
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

func headerContainsWavelength(sheet *xlsx.Sheet, cellrow, cellcolumn int) (yesno bool) {
	headercell, err := spreadsheet.Getdatafromrowcol(sheet, cellrow, cellcolumn).String()

	if err != nil {
		panic(err.Error())

	}

	if strings.Contains(headercell, "(") && strings.Contains(headercell, ")") {
		start := strings.Index(headercell, "(")
		finish := strings.Index(headercell, ")")

		isthisanumber := headercell[start+1 : finish]

		_, err := strconv.Atoi(isthisanumber)

		if err == nil {
			yesno = true
		}
	}

	return
}

func headerWavelength(sheet *xlsx.Sheet, cellrow, cellcolumn int) (yesno bool, number int, err error) {
	headercell, err := spreadsheet.Getdatafromrowcol(sheet, cellrow, cellcolumn).String()

	if err != nil {
		return
	}

	if strings.Contains(headercell, "(") && strings.Contains(headercell, ")") {
		start := strings.Index(headercell, "(")
		finish := strings.Index(headercell, ")")

		isthisanumber := headercell[start+1 : finish]

		number, err = strconv.Atoi(isthisanumber)

		if err == nil {
			yesno = true
		}
	} else {
		err = fmt.Errorf("no (  ) found in header")
	}

	return
}

func (data MarsData) AvailableReadings(wellname string) (readingDescriptions []string) {

	for _, measurement := range data.Dataforeachwell[wellname].Data.Readings[0] {

		var description string

		if measurement.EWavelength == measurement.RWavelength {

			if measurement.Script > 0 {
				description = fmt.Sprintln("Absorbance: ", measurement.EWavelength, "nm. ", "Script position: ", measurement.Script)
			} else {
				if measurement.Script > 0 {
					description = fmt.Sprintln("Absorbance: ", measurement.EWavelength, "nm. ")
				}
			}

		} else {

			if measurement.Script > 0 {
				description = fmt.Sprintln("Excitation: ", measurement.EWavelength, "nm. ", "Emission: ", measurement.RWavelength, "nm. ", "Script position: ", measurement.Script)
			} else {
				if measurement.Script > 0 {
					description = fmt.Sprintln("Excitation: ", measurement.EWavelength, "nm. ", "Emission: ", measurement.RWavelength)
				}
			}

		}

		readingDescriptions = append(readingDescriptions, description)
	}

	readingDescriptions = search.RemoveDuplicates(readingDescriptions)

	return
}

func (data MarsData) TimeCourse(wellname string, exWavelength int, emWavelength int, scriptnumber int) (xaxis []time.Duration, yaxis []float64, err error) {

	xaxis = make([]time.Duration, 0)
	yaxis = make([]float64, 0)
	var emfound bool
	var exfound bool
	if _, found := data.Dataforeachwell[wellname]; !found {
		return xaxis, yaxis, fmt.Errorf(fmt.Sprint("No data found for wellname ", wellname))
	}

	if len(data.Dataforeachwell[wellname].Data.Readings[0]) == 0 {
		return xaxis, yaxis, fmt.Errorf(fmt.Sprint("No readings found for wellname ", wellname))
	}
	for _, measurement := range data.Dataforeachwell[wellname].Data.Readings[0] {

		var checkscriptnumber bool

		if scriptnumber > 0 {
			checkscriptnumber = true
		}

		if measurement.EWavelength == exWavelength && measurement.RWavelength == emWavelength && checkscriptnumber && measurement.Script == scriptnumber {
			emfound = true
			exfound = true
			xaxis = append(xaxis, measurement.Timestamp)
			yaxis = append(yaxis, measurement.Reading)

		} else if measurement.EWavelength == exWavelength && measurement.RWavelength == emWavelength && !checkscriptnumber {

			emfound = true
			exfound = true
			xaxis = append(xaxis, measurement.Timestamp)
			yaxis = append(yaxis, measurement.Reading)

		}

	}
	if emfound != true && exfound != true {
		return xaxis, yaxis, fmt.Errorf(fmt.Sprint("No values found for emWavelength ", emWavelength, " and/or exWavelength ", exWavelength, ". ", "Available Values found: ", data.AvailableReadings(wellname)))
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

// field value is the value which the data is to be filtered by,
// e.g. if filtering by time, this would be the time at which to return readings for;
// if filtering by excitation wavelength, this would be the wavelength at which to return readings for
func (data MarsData) ReadingsAsAverage(wellname string, emexortime int, fieldvalue interface{}, readingtypekeyword string) (average float64, err error) {
	readings := make([]float64, 0)
	readingtypes := make([]string, 0)
	readingsforaverage := make([]float64, 0)
	if _, ok := data.Dataforeachwell[wellname]; !ok {
		return 0.0, fmt.Errorf(fmt.Sprint("no data for well, ", wellname))
	}
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

func (data MarsData) FindOptimalWavelength(wellname string, blankname string, readingtypekeyword string) (wavelength int, err error) {

	if _, ok := data.Dataforeachwell[wellname]; !ok {
		return 0, fmt.Errorf("no data found for well, %s", wellname)
	}
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
	EBand       int
	RBand       int
	Script      int
	Gain        int
}
