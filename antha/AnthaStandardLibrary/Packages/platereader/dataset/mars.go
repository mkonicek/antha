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
	"strconv"
	"strings"
	"time"

	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/platereader"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/search"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/montanaflynn/stats"
	"github.com/pkg/errors"
)

const (
	absorbanceSpectrumHeader = "(Abs Spectrum)"
	emissionSpectrumHeader   = "(Em Spectrum)"
	excitationSpectrumHeader = "(Ex Spectrum)"
	absorbanceHeader         = "(A-"
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

	readingDescriptions = search.RemoveDuplicateStrings(readingDescriptions)

	return
}

// TimeCourse returns either a fluorescence timecourse or an Absorbance timecourse data series.
// If Absorbance, the excitation and emmission wavelengths should both be set to the Absorbance wavelength.
// scriptnumber allows the possibility to distinguish runs with the same excitation and emmission wavelengths but run as different scripts, e.g. if ran with different gains.
// If scriptnumber is set to 0 this will not be used.
func (data MarsData) TimeCourse(wellname string, exWavelength int, emWavelength int, scriptnumber int) (xaxis []time.Duration, yaxis []float64, err error) {

	xaxis = make([]time.Duration, 0)
	yaxis = make([]float64, 0)
	var emfound bool
	var exfound bool
	if _, found := data.Dataforeachwell[wellname]; !found {
		return xaxis, yaxis, fmt.Errorf(fmt.Sprint("No data found for wellname ", wellname))
	} else if len(data.Dataforeachwell[wellname].Data.Readings[0]) == 0 {
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
	if !emfound && !exfound {
		return xaxis, yaxis, fmt.Errorf(fmt.Sprint("No values found for emWavelength ", emWavelength, " and/or exWavelength ", exWavelength, ". ", "Available Values found: ", data.AvailableReadings(wellname)))
	}
	return
}

// AllAbsorbanceData returns all absorbance readings using the well location as key.
func (data MarsData) AllAbsorbanceData() (readings map[string][]wtype.Absorbance, err error) {

	readings = make(map[string][]wtype.Absorbance, len(data.Dataforeachwell))

	for wellName, wellData := range data.Dataforeachwell {

		var wellReadings = make([]wtype.Absorbance, len(wellData.Data.Readings[0]))

		for readingIndex, measurement := range wellData.Data.Readings[0] {
			wellReadings[readingIndex] = wtype.Absorbance{
				Wavelength:   float64(measurement.RWavelength),
				Reading:      measurement.Reading,
				WellLocation: wtype.MakeWellCoordsA1(wellName),
				Annotations:  []string{measurement.ReadingType},
			}
		}

		readings[wellName] = wellReadings
	}

	return readings, nil
}

// AbsScanData returns all wavelengths and readings for a specified well.
func (data MarsData) AbsScanData(well string) (wavelengths []int, Readings []float64) {
	wavelengths = make([]int, 0)
	Readings = make([]float64, 0)
	for _, measurement := range data.Dataforeachwell[well].Data.Readings[0] {

		if strings.Contains(measurement.ReadingType, absorbanceSpectrumHeader) {

			wavelengths = append(wavelengths, measurement.RWavelength)
			Readings = append(Readings, measurement.Reading)

		}
	}

	return
}

// EMScanData returns all emmission wavelengths and readings for a specified well and excitation wavelength.
func (data MarsData) EMScanData(well string, exWavelength int) (wavelengths []int, Readings []float64) {
	wavelengths = make([]int, 0)
	Readings = make([]float64, 0)
	for _, measurement := range data.Dataforeachwell[well].Data.Readings[0] {

		if measurement.EWavelength == exWavelength && strings.Contains(measurement.ReadingType, emissionSpectrumHeader) {

			wavelengths = append(wavelengths, measurement.RWavelength)
			Readings = append(Readings, measurement.Reading)

		}

	}

	return
}

// EMScanData returns all excitation wavelengths and readings for a specified well and emmission wavelength.
func (data MarsData) EXScanData(well string, emWavelength int) (wavelengths []int, Readings []float64) {
	wavelengths = make([]int, 0)
	Readings = make([]float64, 0)
	for _, measurement := range data.Dataforeachwell[well].Data.Readings[0] {

		if measurement.RWavelength == emWavelength && strings.Contains(measurement.ReadingType, excitationSpectrumHeader) {

			wavelengths = append(wavelengths, measurement.EWavelength)
			Readings = append(Readings, measurement.Reading)

		}
	}

	return
}

// WelltoDataMap returns a map of well location (in format A1) to plate reader data for that well.
func (data MarsData) WelltoDataMap() map[string]WellData {
	return data.Dataforeachwell
}

// Readings returns all measurements found for the sepcified well (wells should be specified in A1 format).
func (data MarsData) Readings(well string) []PRMeasurement {
	return data.Dataforeachwell[well].Data.Readings[0]
}

// ReadingsAsAverage returnwwellnameamee data for the specified well matching ReadingType with appropriate fieldvalue.
// field value is the value which the data is to be filtered by,
// e.g. if filtering by time, this would be the time at which to return readings for;
// if filtering by excitation wavelength, this would be the wavelength at which to return readings for.
// readingtypekeyword corresponds to key words found in the header of a data column.
// Examples:
/*
//const (
	absorbanceSpectrumHeader = "(Abs Spectrum)"
	emissionSpectrumHeader   = "(Em Spectrum)"
	excitationSpectrumHeader = "(Ex Spectrum)"
	absorbanceHeader         = "(A-"
	rawDataHeader            = "Raw Data"
) */
func (data MarsData) ReadingsAsAverage(well string, emexortime platereader.FilterOption, fieldvalue interface{}, readingtypekeyword string) (average float64, err error) {

	readings := make([]float64, 0)
	readingtypes := make([]string, 0)
	readingsforaverage := make([]float64, 0)

	well = strings.TrimSpace(well)

	if _, ok := data.Dataforeachwell[well]; !ok {
		return 0.0, fmt.Errorf(fmt.Sprint("no data for well, ", well))
	}
	for _, measurement := range data.Dataforeachwell[well].Data.Readings[0] {

		if emexortime == platereader.TIME {
			if str, ok := fieldvalue.(string); ok {

				gotime, err := time.ParseDuration(str)
				if err != nil {
					return average, err
				}
				if measurement.Timestamp == gotime {
					readings = append(readings, measurement.Reading)
					readingtypes = append(readingtypes, measurement.ReadingType)
				}
			}
		} else if emexortime == platereader.EMWAVELENGTH && measurement.RWavelength == fieldvalue {
			readings = append(readings, measurement.Reading)
			readingtypes = append(readingtypes, measurement.ReadingType)
		} else if emexortime == platereader.EXWAVELENGTH && measurement.EWavelength == fieldvalue {
			readings = append(readings, measurement.Reading)
			readingtypes = append(readingtypes, measurement.ReadingType)
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

// Absorbance returns the average of all readings at a specified wavelength.
// First the exact absorbance reading is searched for, failing that a scan will be searched for.
// If a value for options is declared, this can be used as the header to look for when matching in cases where multiple headers are present for a sample ... e.g. "Blank corrected based on Raw Data (Abs Spectrum)" and "Raw Data (Abs Spectrum)"
func (data MarsData) Absorbance(well string, wavelength int, options ...interface{}) (average wtype.Absorbance, err error) {
	var errs []string

	if len(options) > 1 {
		return wtype.Absorbance{}, errors.Errorf("Only one option is permitted as an argument to the Absorbance method for MarsData")
	}

	if len(options) == 1 {
		result, err := data.ReadingsAsAverage(well, platereader.EMWAVELENGTH, wavelength, fmt.Sprint(options[0]))
		if err == nil {
			return wtype.Absorbance{
				Reading:     result,
				Wavelength:  float64(wavelength),
				Annotations: []string{fmt.Sprint(options[0])},
			}, nil
		} else {
			errs = append(errs, err.Error())
		}
	}
	result, err := data.ReadingsAsAverage(well, platereader.EMWAVELENGTH, wavelength, absorbanceHeader)
	if err == nil {
		return wtype.Absorbance{
			Reading:    result,
			Wavelength: float64(wavelength),
		}, nil
	} else {
		errs = append(errs, err.Error())
	}
	result, err = data.ReadingsAsAverage(well, platereader.EMWAVELENGTH, wavelength, absorbanceSpectrumHeader)
	if err == nil {
		return wtype.Absorbance{
			Reading:    result,
			Wavelength: float64(wavelength),
		}, nil
	} else {
		errs = append(errs, err.Error())
	}
	result, err = data.ReadingsAsAverage(well, platereader.EMWAVELENGTH, wavelength, strings.Join([]string{"(", strconv.Itoa(wavelength), ""}, ""))

	if err == nil {
		return wtype.Absorbance{
			Reading:    result,
			Wavelength: float64(wavelength),
		}, nil
	}

	errs = append(errs, err.Error())

	return wtype.Absorbance{
		Reading:    0.0,
		Wavelength: float64(wavelength),
	}, fmt.Errorf(strings.Join(errs, "\n"))
}

// FindOptimalAbsorbanceWavelength returns the wavelength for which the difference in signal between the sample and blank is greatest.
func (data MarsData) FindOptimalAbsorbanceWavelength(well string, blankname string) (wavelength int, err error) {

	if _, ok := data.Dataforeachwell[well]; !ok {
		return 0, fmt.Errorf("no data found for well, %s", well)
	}
	biggestdifferenceindex := 0
	biggestdifference := 0.0

	wavelengths, readings := data.AbsScanData(well)
	blankwavelengths, blankreadings := data.AbsScanData(blankname)

	for i, reading := range readings {

		difference := reading - blankreadings[i]

		if difference > biggestdifference && wavelengths[i] == blankwavelengths[i] {
			biggestdifferenceindex = i
		}

	}

	wavelength = wavelengths[biggestdifferenceindex]

	return
}

// MarsData represents the contents of a parsed platereader data file exported from Mars.
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

// WellData represents the details of a reading for a well.
type WellData struct {
	Well            string // in a1 format
	Name            string
	Data            PROutput
	Injected        bool
	InjectionVolume float64
}

// PROutput is a modified version of the platereeader data type from antha/microArch/driver/platereading.
type PROutput struct {
	Readings []PRMeasurementSet
}

// PRMeasurementSet is a set of Plate reader measurements
type PRMeasurementSet []PRMeasurement

// PRMeasurement contains the details of a plate reader measurement.
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
	// ReadingType is the annotation found in the column header in the exported mars excel file
	// e.g. Raw Data (Abs Scan)
	ReadingType string
}
