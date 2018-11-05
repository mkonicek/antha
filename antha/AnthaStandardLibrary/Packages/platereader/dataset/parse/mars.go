package parse

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"time"

	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/platereader/dataset"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/search"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/spreadsheet"
	"github.com/antha-lang/antha/antha/anthalib/wutil"
	"github.com/tealeg/xlsx"
)

// ParseMarsXLSXOutput parses mars data from excel filename
func ParseMarsXLSXOutput(xlsxname string, sheet int) (dataoutput dataset.MarsData, err error) {
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

// ParseMarsXLSXBinary parses mars data from excel filename
func ParseMarsXLSXBinary(xlsxContents []byte, sheet int) (dataoutput dataset.MarsData, err error) {

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

func parseHeadLines(xlsxBinary []byte, sheet int) (dataoutput dataset.MarsData, headerrowcount int, err error) {
	xlsx, err := spreadsheet.OpenXLSXBinary(xlsxBinary)

	if err != nil {
		return
	}

	if sheet > len(xlsx.Sheets)-1 {
		err = fmt.Errorf("sheet number %d specified does not exist in file, only found %d sheets. remember that sheet position starts from 0 (i.e., the first sheet position is 0 not 1)", sheet, len(xlsx.Sheets))
		return
	}
	sheet1 := xlsx.Sheets[sheet]

	for i := 0; i < sheet1.MaxRow; i++ {
		str := sheet1.Cell(i, 0).String()

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

	cells, err := spreadsheet.GetDataFromCells(sheet1, cellnames)
	if err != nil {
		return
	}

	dataoutput.Description = cells[len(cells)-1].String()

	for _, cell := range cells {

		cellstr := cell.String()

		if strings.HasPrefix(cellstr, "User") {
			dataoutput.User = strings.Split(cellstr, ": ")[1]
		}
		if strings.HasPrefix(cellstr, "Path") {
			dataoutput.Path = strings.Split(cellstr, ": ")[1]
		}
		if strings.HasPrefix(cellstr, "Test ID") {
			id, err := strconv.Atoi(strings.Split(cellstr, ": ")[1])
			if err != nil {
				return dataoutput, headerrowcount, err
			}
			dataoutput.TestID = id
		}
		if strings.HasPrefix(cellstr, "Test Name") {
			dataoutput.Testname = strings.Split(cellstr, ": ")[1]
		}
		if strings.HasPrefix(cellstr, "Date") {
			date := strings.Split(cellstr, ": ")[1]
			godate, err := time.Parse("02/01/2006", date)
			if err != nil {
				return dataoutput, headerrowcount, err
			}

			dataoutput.Date = godate
		}
		if strings.HasPrefix(cellstr, "Time") {
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
		if strings.HasPrefix(cellstr, "ID1") {
			dataoutput.ID1 = strings.Split(cellstr, ": ")[1]
		}
		if strings.HasPrefix(cellstr, "ID2") {
			dataoutput.ID2 = strings.Split(cellstr, ": ")[1]
		}
		if strings.HasPrefix(cellstr, "ID3") {
			dataoutput.ID3 = strings.Split(cellstr, ": ")[1]
		}

	}
	return
}

func parseWellData(xlsxBinary []byte, sheet int, headerrows int) (welldatamap map[string]dataset.WellData, err error) {
	welldatamap = make(map[string]dataset.WellData)
	var welldata dataset.WellData
	var wavelengthstring string
	var wavelength int
	var timestring string
	var timestamp time.Duration
	xlsx, err := spreadsheet.OpenXLSXBinary(xlsxBinary)

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

		cellstr := cell.String()

		if cellstr == "A" {
			wellrowstart = i
			break
		}

	}
	wavelengths := make([]int, 0)

	if wellrowstart-headerrow > 0 {
		for i := 0; i < wellrowstart-headerrow; i++ {

			rowabove, err := spreadsheet.GetDataFromRowCol(sheet1, wellrowstart-(i+1), 2)
			if err != nil {
				return welldatamap, err
			}
			if strings.Contains(rowabove.String(), "Time") {
				timerow = wellrowstart - (i + 1)
			} else if strings.Contains(rowabove.String(), "Wavelength") {
				wavelengthrow = wellrowstart - (i + 1)
			}

		}
	}
	// check other row names in case the row labels are not in order (this can happen)
	for i := wellrowstart; i < sheet1.MaxRow; i++ {

		rowname, err := spreadsheet.GetDataFromRowCol(sheet1, i, 2)
		if err != nil {
			return welldatamap, err
		}
		if strings.Contains(rowname.String(), "Time") {
			timerow = i
		} else if strings.Contains(rowname.String(), "Wavelength") {
			wavelengthrow = i
		}

	}

	// find special columns
	tempcolumn := 0
	injectionvoumecolumn := 0

	for m := 3; m < sheet1.MaxCol; m++ {

		columnheader, err := spreadsheet.GetDataFromRowCol(sheet1, headerrow, m)
		if err != nil {
			return welldatamap, err
		}
		if strings.Contains(columnheader.String(), "Temperature") {
			tempcolumn = m
		}
		if strings.Contains(columnheader.String(), "Volume") {
			injectionvoumecolumn = m
		}

	}

	for j := wellrowstart; j < sheet1.MaxRow; j++ {

		if j != timerow && j != wavelengthrow {

			cellData, err := spreadsheet.GetDataFromRowCol(sheet1, j, 2)
			if err != nil {
				return welldatamap, err
			}
			welldata.Name = cellData.String()

			part1, err := spreadsheet.GetDataFromRowCol(sheet1, (j), 0)
			if err != nil {
				return welldatamap, err
			}
			part2, err := spreadsheet.GetDataFromRowCol(sheet1, j, 1)
			if err != nil {
				return welldatamap, err
			}
			welldata.Well = part1.String() + part2.String()

			for k := 3; k < sheet1.MaxCol; k++ {
				if k != tempcolumn && k != injectionvoumecolumn {
					if wavelengthrow != 0 {
						wavelength, err := spreadsheet.GetDataFromRowCol(sheet1, wavelengthrow, k)
						if err != nil {
							return welldatamap, err
						}

						wavelengthInt, err := wavelength.Int()
						if err != nil {
							return welldatamap, err
						}
						if len(wavelengths) == 0 {
							wavelengths = append(wavelengths, wavelengthInt)
						} else if !search.InInts(wavelengths, wavelengthInt) {
							wavelengths = append(wavelengths, wavelengthInt)
						}
					}
				}
			}

			for m := 3; m < sheet1.MaxCol; m++ {

				if wavelengthrow != 0 {
					wavelength, err := spreadsheet.GetDataFromRowCol(sheet1, wavelengthrow, m)
					if err != nil {
						return welldatamap, err
					}
					wavelengthInt, err := wavelength.Int()
					if err != nil {
						return welldatamap, err
					}
					if wavelengths[0] != wavelengthInt {
						break
					}
				}
			}
			var measurement dataset.PRMeasurement

			var Measurements = make([]dataset.PRMeasurement, 0)

			maxcol := sheet1.MaxCol

			for m := 3; m < maxcol; m++ {

				readingtype, err := spreadsheet.GetDataFromRowCol(sheet1, headerrow, m)
				if err != nil {
					return welldatamap, err
				}
				measurement.ReadingType = readingtype.String()

				//check header
				header, err := spreadsheet.GetDataFromRowCol(sheet1, headerrow, m)
				if err != nil {
					return welldatamap, err
				}
				// the measurement itself (if not a special column - e.g. volume injection or temp)
				if !strings.Contains(header.String(), "Temperature") && !strings.Contains(header.String(), "Volume") {

					reading, err := spreadsheet.GetDataFromRowCol(sheet1, j, m)
					if err != nil {
						return welldatamap, err
					}
					measurement.Reading, err = reading.Float()
					if err != nil {
						if reading.String() == "" {
							measurement.Reading = 0.0
						} else {
							return welldatamap, err
						}
					}
				}

				// add logic to check column heading
				// add similar for volume (injection)
				if strings.Contains(header.String(), "Temperature") {
					temp, err := spreadsheet.GetDataFromRowCol(sheet1, j, tempcolumn)
					if err != nil {
						return welldatamap, err
					}
					measurement.Temp, err = temp.Float()
					if err != nil {
						return welldatamap, err
					}
				} else if strings.Contains(header.String(), "Volume") {
					injVol, err := spreadsheet.GetDataFromRowCol(sheet1, j, injectionvoumecolumn)
					if err != nil {
						return welldatamap, err
					}
					welldata.InjectionVolume, err = injVol.Float()
					if err != nil {
						return welldatamap, err
					}
					welldata.Injected = true
				} else {

					// add time row and wavelength row calculators
					if timerow != 0 {
						//gotime, err := ParseTime(spreadsheet.Getdatafromrowcol(sheet1, timerow, m).String())
						timelabel, err := spreadsheet.GetDataFromRowCol(sheet1, timerow, 2)
						if err != nil {
							return welldatamap, err
						}
						timecellcontents, err := spreadsheet.GetDataFromRowCol(sheet1, timerow, m)
						if err != nil {
							return welldatamap, err
						}
						if strings.Contains(timelabel.String(), "[s]") && timecellcontents.String() != "" {
							timestring = timecellcontents.String() + "s"

							if err != nil {
								return welldatamap, err
							}
						} else if timecellcontents.String() != "" {
							timestring = timecellcontents.String()
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

					readingType, err := spreadsheet.GetDataFromRowCol(sheet1, headerrow, m)
					if err != nil {
						return welldatamap, err
					}
					measurement.ReadingType = readingType.String()

					parsedatatype := strings.Split(measurement.ReadingType, `(`)

					parsedatatype = strings.Split(parsedatatype[1], `)`)

					if !strings.Contains(header.String(), "Temperature") && !strings.Contains(header.String(), "Volume") {

						// handle case of absorbance (may need to add others.. if contains Ex, Ex else number = abs
						ex, exband, em, emband, scriptposition, err := parseBracketedColumnHeader(measurement.ReadingType)

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
						} else if strings.Contains(measurement.ReadingType, "Em Spectrum") {
							if wavelengthrow != 0 {
								emWavelength, err := spreadsheet.GetDataFromRowCol(sheet1, wavelengthrow, m)
								if err != nil {
									return welldatamap, err
								}
								measurement.RWavelength, err = emWavelength.Int()
								if err != nil {
									return welldatamap, err
								}
							}
						} else if strings.Contains(measurement.ReadingType, "Ex Spectrum") {
							if wavelengthrow != 0 {
								exWavelength, err := spreadsheet.GetDataFromRowCol(sheet1, wavelengthrow, m)
								if err != nil {
									return welldatamap, err
								}
								measurement.EWavelength, err = exWavelength.Int()
								if err != nil {
									return welldatamap, err
								}
							}
						} else if strings.Contains(measurement.ReadingType, "Abs Spectrum") {
							if wavelengthrow == 0 {
								wavelengthstring = parsedatatype[0]

								wavelength, err = strconv.Atoi(wavelengthstring)
								if err != nil {
									return welldatamap, err
								}
								measurement.RWavelength = wavelength
								measurement.EWavelength = wavelength
							} else {
								wavelengthCell, err := spreadsheet.GetDataFromRowCol(sheet1, wavelengthrow, m)
								if err != nil {
									return welldatamap, err
								}

								wavelength, err := wavelengthCell.Int()
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
								wavelengthCell, err := spreadsheet.GetDataFromRowCol(sheet1, wavelengthrow, m)
								if err != nil {
									return welldatamap, err
								}

								wavelength, err := wavelengthCell.Int()
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

			var output dataset.PROutput
			set := Measurements

			output.Readings = make([]dataset.PRMeasurementSet, 1)
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
			}
			err = fmt.Errorf("unknown header type, %s ,found in Mars data file, problem with %s", header, fields[1])
			return

		}
	}
	err = fmt.Errorf("Error with header %s found in Mars data file", header)
	return
}

// ParseTime parses a plate reader string into a Duration
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
	headercell, err := spreadsheet.GetDataFromRowCol(sheet, cellrow, cellcolumn)
	if err != nil {
		panic(err)
	}
	if strings.Contains(headercell.String(), "(") && strings.Contains(headercell.String(), ")") {
		start := strings.Index(headercell.String(), "(")
		finish := strings.Index(headercell.String(), ")")

		isthisanumber := headercell.String()[start+1 : finish]

		_, err := strconv.Atoi(isthisanumber)

		if err == nil {
			yesno = true
		}
	}

	return
}

func headerWavelength(sheet *xlsx.Sheet, cellrow, cellcolumn int) (yesno bool, number int, err error) {
	headercell, err := spreadsheet.GetDataFromRowCol(sheet, cellrow, cellcolumn)
	if err != nil {
		return
	}
	if strings.Contains(headercell.String(), "(") && strings.Contains(headercell.String(), ")") {
		start := strings.Index(headercell.String(), "(")
		finish := strings.Index(headercell.String(), ")")

		isthisanumber := headercell.String()[start+1 : finish]

		number, err = strconv.Atoi(isthisanumber)

		if err == nil {
			yesno = true
		}
	} else {
		err = fmt.Errorf("no (  ) found in header")
	}

	return
}
