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
	"encoding/xml"
	"fmt"
	"time"

	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
)

type SpectraMaxData struct {
	Experiment Experiment `xml:"Experiment"`
}

type Experiment struct {
	PlateSections `xml:"PlateSections"`
}

type PlateSections struct {
	PlateSections []PlateSection `xml:"PlateSection"`
}

type PlateSection struct {
	Name               string             `xml:"Name,attr"`
	InstrumentInfo     string             `xml:"InstrumentInfo,attr"`
	ReadTime           time.Duration      `xml:"ReadTime,attr"`
	Barcode            string             `xml:"Barcode"`
	InstrumentSettings InstrumentSettings `xml:"InstrumentSettings"`
	Wavelengths        []Reading          `xml:"Wavelengths"`
	TemperatureData    wunit.Temperature  `xml:"TemperatureData"`
}

type InstrumentSettings struct {
	ReadMode           ReadMode           `xml:"ReadMode,attr"`
	ReadType           ReadType           `xml:"ReadType,attr"`
	PlateType          wtype.LHPlate      `xml:"PlateType,attr"` // may need to change to a string for now since it's unlikely the plate names in the platereader software will correspond to those in antha
	AutoMix            bool               `xml:"AutoMix"`
	MoreSettings       MoreSettings       `xml:"MoreSettings"`
	WavelengthSettings WavelengthSettings `xml:"WavelengthSettings"`
}

type ReadMode string // or could just make this a string

const (
	Absorbance ReadMode = "Absorbance"
)

type ReadType string // or could just make this a string

const (
	Endpoint ReadType = "Endpoint"
)

type WavelengthSettings struct {
	NumberOfWavelengths int          `xml: "NumberOfWavelengths,attr"`
	Wavelengths         []Wavelength `xml:"Wavelength"`
}

type Wavelength struct {
	Index   int `xml:"Index,attr"`
	WLength int `xml:"Wavelength"`
	//WavelengthIndex int `xml: "WavelengthIndex,attr"`
}

type MoreSettings struct {
	Calibrate     bool   `xml:"Calibrate"`
	CarriageSpeed string `xml:"CarriageSpeed"`
	ReadOrder     string `xml:"ReadOrder"`
}

type Reading struct {
	Wavelength Wavelength `xml:"Wavelength"`
	Wells      Wells      `xml:"Wells"`
}

type Wells struct {
	Wells []Well `xml:"Well"`
}

type Well struct {
	ID      string  `xml:"ID,attr"`
	Name    string  `xml:"Name,attr"`
	Row     int     `xml:"Row,attr"`
	Column  int     `xml:"Column,attr"`
	RawData float64 `xml:"RawData"`
}

func ParseSpectraMaxData(xmlFileContents []byte) (dataOutput SpectraMaxData, err error) {

	err = xml.Unmarshal(xmlFileContents, &dataOutput)

	if err != nil {
		fmt.Println("error:", err)
	}

	return
}

// readingtypekeyword is irrelevant for this data set but needed to conform to the current interface!
func (s SpectraMaxData) BlankCorrect(wellnames []string, blanknames []string, wavelength int, readingtypekeyword string) (blankcorrectedaverage float64, err error) {

	return
}

// emexortime is selected from the constants above
//const (
//	TIME = iota
//	EMWAVELENGTH
//	EXWAVELENGTH
//)
/*
could make this a type

type ReadingType int

const (
	TIME ReadingType = iota
	EMWAVELENGTH
	EXWAVELENGTH
)

func (r ReadingType) ValidTypes() string{

return
}

*/
// readingtypekeyword is irrelevant for this data set but needed to conform to the current interface!
// field value is the value which the data is to be filtered by,
// e.g. if filtering by time, this would be the time at which to return readings for;
// if filtering by excitation wavelength, this would be the wavelength at which to return readings for
func (s SpectraMaxData) ReadingsAsAverage(wellname string, emexortime int, fieldvalue interface{}, readingtypekeyword string) (average float64, err error) {

	return
}

// readingtypekeyword is irrelevant for this data set but needed to conform to the current interface!
func (s SpectraMaxData) FindOptimalWavelength(wellname string, blankname string, readingtypekeyword string) (wavelength int, err error) {

	return
}

// scriptnumber is irrelevant for this data set but needed to conform to the current interface!
func (s SpectraMaxData) TimeCourse(wellname string, exWavelength int, emWavelength int, scriptnumber int) (xaxis []time.Duration, yaxis []float64, err error) {
	return
}
