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
	"encoding/xml"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/platereader"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/montanaflynn/stats"
)

const timeFormat = "15:04 01/02/2006"

type customTime struct {
	time.Time
}

//XMLExperiment is exported so requires a comment
type SpectraMaxData struct {
	Name       xml.Name           `xml:"Experiment"`
	Experiment []XMLPlateSections `xml:"PlateSections"`
}

//XMLPlateSections is exported so requires a comment
type XMLPlateSections struct {
	PlateSections []XMLPlateSection `xml:"PlateSection"`
}

//XMLPlateSection is exported so requires a comment
type XMLPlateSection struct {
	Name               string                `xml:"Name,attr"`
	InstrumentInfo     string                `xml:"InstrumentInfo,attr"`
	ReadTime           customTime            `xml:"ReadTime,attr"`
	Barcode            string                `xml:"Barcode"`
	InstrumentSettings XMLInstrumentSettings `xml:"InstrumentSettings"`
	Wavelengths        []Reading             `xml:"Wavelengths"`
	TemperatureData    wunit.Temperature     `xml:"TemperatureData"`
}

//XMLInstrumentSettings is exported so requires a comment
type XMLInstrumentSettings struct {
	ReadMode           platereader.ReadMode `xml:"ReadMode,attr"`
	ReadType           platereader.ReadType `xml:"ReadType,attr"`
	PlateType          string               `xml:"PlateType,attr"` // may need to change to a string for now since it's unlikely the plate names in the platereader software will correspond to those in antha
	AutoMix            bool                 `xml:"AutoMix"`
	MoreSettings       MoreSettings         `xml:"MoreSettings"`
	WavelengthSettings WavelengthSettings   `xml:"WavelengthSettings"`
}

// WavelengthSettings is exported so requires a comment
type WavelengthSettings struct {
	NumberOfWavelengths int      `xml:"NumberOfWavelengths,attr"`
	Wavelength          []string `xml:"Wavelength"`
}

// Wavelength is exported so requires a comment
type Wavelength struct {
	Index int     `xml:"WavelengthIndex,attr"`
	Wells []Wells `xml:"Wells"`
}

// MoreSettings is exported so requires a comment
type MoreSettings struct {
	Calibrate     string `xml:"Calibrate"`
	CarriageSpeed string `xml:"CarriageSpeed"`
	ReadOrder     string `xml:"ReadOrder"`
}

// Reading is exported so requires a comment
type Reading struct {
	Wavelength Wavelength `xml:"Wavelength"`
	Wells      []Well     `xml:"Wells"`
}

// Wells is exported so requires a comment
type Wells struct {
	Wells []Well `xml:"Well"`
}

// Well is exported so requires a comment
type Well struct {
	ID       string `xml:"ID,attr"`     // Single reading
	WellID   string `xml:"WellID,attr"` // Scan data
	Name     string `xml:"Name,attr"`
	Row      int    `xml:"Row,attr"`
	Column   int    `xml:"Col,attr"`
	RawData  string `xml:"RawData"`
	WaveData string `xml:"WaveData"` // Scan data
}

type WavelengthReading struct {
	Wavelength int
	Reading    float64
}

func (w Well) IsScanData() bool {
	return len(w.WaveData) > 0
}

func readingAtWavelength(readings []WavelengthReading, wavelength int) (reading float64, err error) {
	for _, reading := range readings {
		if reading.Wavelength == wavelength {
			return reading.Reading, nil
		}
	}
	return 0.0, fmt.Errorf("No reading found for wavelength %d: found: %+v", wavelength, readings)
}

// GetDataByWell returns all readings for that well.
func (s SpectraMaxData) GetDataByWell(wellName string) (readings []WavelengthReading, err error) {

	wells := s.Experiment[0].PlateSections[0].Wavelengths[0].Wavelength.Wells[0].Wells
	var w Well
	var wellFound bool

	for _, well := range wells {
		if well.Name == wellName {
			w = well
			wellFound = true
			break
		}
	}

	if !wellFound {
		return readings, fmt.Errorf("No readings found for well %s: found: %+v", wellName, s.Experiment[0].PlateSections[0].Wavelengths[0])
	}

	if w.IsScanData() {
		dataStrings := strings.Fields(w.RawData)
		wavelengthsStrings := strings.Fields(w.WaveData)

		for i := range wavelengthsStrings {
			wavelength, err := strconv.Atoi(wavelengthsStrings[i])
			if err != nil {
				return readings, err
			}
			data, err := strconv.ParseFloat(dataStrings[i], 64)
			if err != nil {
				return readings, err
			}

			var reading WavelengthReading
			reading.Wavelength = wavelength
			reading.Reading = data

			readings = append(readings, reading)

		}

		if len(readings) == 0 {
			return readings, fmt.Errorf("No readings found for well %s: found: %+v", wellName, wells)
		}

	} else {
		return readings, fmt.Errorf("Only Spectramax data in scan format is currently supported. Please run Absorbance reading as scan")
	}
	return
}

func (c *customTime) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	const shortForm = timeFormat // yyyymmdd date format
	var v string
	err := d.DecodeElement(&v, &start)
	if err != nil {
		return err
	}
	parse, err := time.Parse(shortForm, v)
	if err != nil {
		return err
	}
	*c = customTime{parse}
	return nil
}

func (c *customTime) UnmarshalXMLAttr(attr xml.Attr) error {
	parse, err := time.Parse(timeFormat, attr.Value)
	if err != nil {
		return err
	}

	*c = customTime{parse}
	return nil
}

// BlankCorrect subtracts the mean of the values matching the specified wavelength of the blank wells specified by the sample wells.
func (s SpectraMaxData) BlankCorrect(wellnames []string, blanknames []string, wavelength int) (blankcorrectedaverage float64, err error) {
	var data []float64
	var blankdata []float64

	// replace with Readings method
	for _, well := range wellnames {

		reading, err := s.ReadingsAsAverage(well, platereader.EMWAVELENGTH, wavelength)
		if err != nil {
			return blankcorrectedaverage, err
		}
		data = append(data, reading)

	}

	// replace with Readings method
	for _, blankWell := range blanknames {

		reading, err := s.ReadingsAsAverage(blankWell, platereader.EMWAVELENGTH, wavelength)
		if err != nil {
			return blankcorrectedaverage, err
		}
		blankdata = append(blankdata, reading)

	}

	mean, err := stats.Mean(data)
	if err != nil {
		return blankcorrectedaverage, err
	}
	blankmean, err := stats.Mean(blankdata)
	if err != nil {
		return blankcorrectedaverage, err
	}

	blankcorrectedaverage = mean - blankmean

	return blankcorrectedaverage, err
}

// ReadingsAsAverage returns the data for the specified well matching ReadingType with appropriate fieldvalue.
// Currently only Absorbance data as endpoint scans are supported.
// Currently the only valid FilterOptions are platereader.EMWAVELENGTH or platereader.EXWAVELENGTH with the field value being the wavelength as an int.
func (s SpectraMaxData) ReadingsAsAverage(wellname string, emexortime platereader.FilterOption, fieldvalue interface{}) (average float64, err error) {
	var data []float64
	var wavelength int

	if emexortime == platereader.EMWAVELENGTH || emexortime == platereader.EXWAVELENGTH {
		var ok bool

		wavelength, ok = fieldvalue.(int)

		if !ok {
			return average, fmt.Errorf("fieldvalue must be a wavelength if emexortime is set to EMWAVELENGTH or EXWAVELENGTH")
		}
	} else {
		return average, fmt.Errorf("currently spectramax data is only supported as an endpoint absorbance reading. Hence, FilterOption must be set to platereader.EMWAVELENGTH or platereader.EXWAVELENGTH and fieldvalue must be an int.")
	}

	wellData, err := s.GetDataByWell(wellname)

	if err != nil {
		return average, err
	}

	reading, err := readingAtWavelength(wellData, wavelength)

	if err != nil {
		return average, err
	}

	data = append(data, reading)

	return stats.Mean(data)
}

// Absorbance returns the absorbance reading of the specified well at the specified wavelength.
// currently no additional options are supported.
func (s SpectraMaxData) Absorbance(wellname string, wavelength int, options ...interface{}) (average wtype.Absorbance, err error) {
	raw, err := s.ReadingsAsAverage(wellname, platereader.EMWAVELENGTH, wavelength)

	return wtype.Absorbance{
		Reading:    raw,
		Wavelength: float64(wavelength),
	}, err
}

// FindOptimalAbsorbanceWavelength returns the wavelength for which the difference in signal between the sample and blank is greatest.
func (s SpectraMaxData) FindOptimalAbsorbanceWavelength(wellname string, blankname string) (wavelength int, err error) {

	wellData, err := s.GetDataByWell(wellname)

	if err != nil {
		return wavelength, err
	}

	blankData, err := s.GetDataByWell(blankname)

	if err != nil {
		return wavelength, err
	}

	biggestdifferenceindex := 0
	biggestdifference := 0.0

	for i, reading := range wellData {

		difference := reading.Reading - blankData[i].Reading

		if difference > biggestdifference && reading.Wavelength == blankData[i].Wavelength {
			biggestdifferenceindex = i
		}

	}

	wavelength = wellData[biggestdifferenceindex].Wavelength
	return wavelength, nil
}
