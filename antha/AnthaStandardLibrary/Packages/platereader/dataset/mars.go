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
package dataset

import (
	"fmt"
	"strings"
	"time"

	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/search"

	"github.com/montanaflynn/stats"
)

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
