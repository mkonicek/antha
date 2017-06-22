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
	"bufio"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"time"
	"unicode"
	"unicode/utf8"

	"github.com/Synthace/antha/antha/anthalib/wunit"
)

const timeFormat = "15:04 01/02/2006"

type customTime struct {
	time.Time
}

//XMLExperiment is exported so requires a comment
type SpectraMaxData struct {
	Name       xml.Name           ` xml:"Experiment"`
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
	ReadMode           ReadMode           `xml:"ReadMode,attr"`
	ReadType           ReadType           `xml:"ReadType,attr"`
	PlateType          string             `xml:"PlateType,attr"` // may need to change to a string for now since it's unlikely the plate names in the platereader software will correspond to those in antha
	AutoMix            bool               `xml:"AutoMix"`
	MoreSettings       MoreSettings       `xml:"MoreSettings"`
	WavelengthSettings WavelengthSettings `xml:"WavelengthSettings"`
}

// ReadMode  is exported so requires a comment
type ReadMode string // or could just make this a string

// Absorbance is exported so requires a comment
const (
	Absorbance ReadMode = "Absorbance"
)

// ReadType is exported so requires a comment
type ReadType string // or could just make this a string

//Endpoint  is exported so requires a comment
const (
	Endpoint ReadType = "Endpoint"
)

//WavelengthSettings is exported so requires a comment
type WavelengthSettings struct {
	NumberOfWavelengths int      `xml:"NumberOfWavelengths,attr"`
	Wavelength          []string `xml:"Wavelength"`
}

//Wavelength is exported so requires a comment
type Wavelength struct {
	Index int     `xml:"WavelengthIndex,attr"`
	Wells []Wells `xml: "Wells"`
}

//MoreSettings is exported so requires a comment
type MoreSettings struct {
	Calibrate     string `xml:"Calibrate"`
	CarriageSpeed string `xml:"CarriageSpeed"`
	ReadOrder     string `xml:"ReadOrder"`
}

//Reading is exported so requires a comment
type Reading struct {
	Wavelength Wavelength `xml:"Wavelength"`
	Wells      []Well     `xml:"Wells"`
}

//Wells is exported so requires a comment
type Wells struct {
	Wells []Well `xml:"Well"`
}

//Well is exported so requires a comment
type Well struct {
	ID      string  `xml:"ID,attr"`
	Name    string  `xml:"Name,attr"`
	Row     int     `xml:"Row,attr"`
	Column  int     `xml:"Col,attr"`
	RawData float64 `xml:"RawData"`
}

func (c *customTime) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	const shortForm = timeFormat // yyyymmdd date format
	var v string
	d.DecodeElement(&v, &start)
	parse, err := time.Parse(shortForm, v)
	if err != nil {
		return err
	}
	*c = customTime{parse}
	return nil
}

func (c *customTime) UnmarshalXMLAttr(attr xml.Attr) error {
	parse, _ := time.Parse(timeFormat, attr.Value)
	*c = customTime{parse}
	return nil
}

//
func ParseSpectraMaxData(xmlFileContents []byte) (dataOutput SpectraMaxData, err error) {

	/*
		buff := bytes.NewBuffer(xmlFileContents)

		decoder := xml.NewDecoder(NewValidUTF8Reader(buff))

		err = decoder.Decode(&dataOutput)

	*/

	// add header
	xmlFileContents = []byte(xml.Header + string(xmlFileContents))

	err = xml.Unmarshal(xmlFileContents, &dataOutput)

	if err != nil {
		fmt.Println("error:", err)
	}
	pretty, err := json.MarshalIndent(dataOutput, "", "  ")

	if err != nil {
		fmt.Println("error:", err)
	}

	fmt.Printf("%s", string(pretty))
	return
}

func readExperiment(reader io.Reader) ([]XMLPlateSections, error) {

	var spectraMaxData SpectraMaxData
	if err := xml.NewDecoder(reader).Decode(&spectraMaxData); err != nil {
		return nil, err
	}

	return spectraMaxData.Experiment, nil
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

//
// ValidUTF8Reader implements a Reader which reads only bytes that constitute valid UTF-8
type ValidUTF8Reader struct {
	buffer *bufio.Reader
}

// Function Read reads bytes in the byte array b. n is the number of bytes read.
func (rd ValidUTF8Reader) Read(b []byte) (n int, err error) {
	for {
		var r rune
		var size int
		r, size, err = rd.buffer.ReadRune()
		if err != nil {
			return
		}
		if r == unicode.ReplacementChar && size == 1 {
			continue
		} else if n+size < len(b) {
			fmt.Println("replacing: ", string(r))
			utf8.EncodeRune(b[n:], r)
			n += size
		} else {
			rd.buffer.UnreadRune()
			break
		}
	}
	return
}

// NewValidUTF8Reader constructs a new ValidUTF8Reader that wraps an existing io.Reader
func NewValidUTF8Reader(rd io.Reader) ValidUTF8Reader {
	return ValidUTF8Reader{bufio.NewReader(rd)}
}
