// antha/AnthaStandardLibrary/Packages/Parser/gdxparser.go: Part of the Antha language
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
	"io/ioutil"
	"strconv"
	"strings"

	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/enzymes"

	"github.com/antha-lang/antha/antha/anthalib/wtype"
)

type Project struct {
	DesignConstruct []DesignConstruct `xml:"DesignConstruct"`
}

type DesignConstruct struct {
	Label       string       `xml:"label,attr"`
	Plasmid     string       `xml:"circular,attr"`
	Rev         string       `xml:"reverseComplement,attr"`
	Notes       string       `xml:"notes"`
	DNAElements []DNAElement `xml:"DNAElement"`
	AAElements  []AAElement  `xml:"AAElement"`
}

type DNAElement struct {
	Label    string `xml:"label,attr"`
	Sequence string `xml:"sequence"`
	Notes    string `xml:"notes"`
}

type AAElement struct {
	Label    string `xml:"label,attr"`
	Sequence string `xml:"sequence"`
	Notes    string `xml:"notes"`
}

func Parse(filename string) (parts_list []string, err error) {

	str, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var gdx Project

	err = xml.Unmarshal(str, &gdx)
	if err != nil {
		return parts_list, err
	}

	parts_list = make([]string, len(gdx.DesignConstruct))

	nconstructs := 0
	for _, c := range gdx.DesignConstruct {
		parts_list[nconstructs] = "Construct: " + strconv.Itoa(nconstructs) + " n parts: " + strconv.Itoa(len(c.DNAElements)+len(c.AAElements))
		nconstructs++
	}

	return parts_list, err
}

func ParseToAssemblyParameters(filename string) ([]enzymes.Assemblyparameters, error) {
	str, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return gdxToAssemblyParameters(str)
}

// ParsesGDX parses a typeIIs assembly design file in gdx format. An example
// file is provided: "Assembly_Input_Controls.gdx" The output will be
// []enzymes.AssemblyParameters which can be used in the
// enzymes.Assemblysimulator() and enzymes.Digestionsimulator() functions.  The
// design file is expected to follow a format as shown in the provided example
// files An error will be returned if no data is found within the .gdx design
// file or if the file is not in the expected format.
func ParseGDX(file wtype.File) ([]enzymes.Assemblyparameters, error) {
	data, err := file.ReadAll()
	if err != nil {
		return nil, err
	}
	return gdxToAssemblyParameters(data)
}

// ParseGDXBinary parses the contents of a typeIIs assembly design file in gdx
// format. An example file is provided: "Assembly_Input_Controls.gdx" The
// output will be []enzymes.AssemblyParameters which can be used in the
// enzymes.Assemblysimulator() and enzymes.Digestionsimulator() functions.  The
// design file is expected to follow a format as shown in the provided example
// files An error will be returned if no data is found within the .gdx design
// file or if the file is not in the expected format.
func ParseGDXBinary(data []byte) ([]enzymes.Assemblyparameters, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("no data found")
	}
	return gdxToAssemblyParameters(data)
}

func gdxToAssemblyParameters(data []byte) ([]enzymes.Assemblyparameters, error) {

	var gdx Project

	construct_list := make([]enzymes.Assemblyparameters, 0)
	err := xml.Unmarshal(data, &gdx)
	if err != nil {
		return construct_list, err
	}

	if len(gdx.DesignConstruct) == 0 {
		return construct_list, fmt.Errorf("Empty design construct in gdx file")
	}
	for _, a := range gdx.DesignConstruct {
		var newconstruct enzymes.Assemblyparameters
		newconstruct.Constructname = a.Label
		if strings.Contains(a.Notes, "Enzyme:") {
			newconstruct.Enzymename = strings.TrimSpace(strings.TrimPrefix(a.Notes, "Enzyme:")) // add trim function to trim after space
		}
		for _, b := range a.DNAElements {
			var newseq wtype.DNASequence
			if strings.Contains(strings.ToUpper(b.Notes), "VECTOR") {
				newseq.Nm = b.Label
				newseq.Seq = b.Sequence
				if strings.Contains(strings.ToUpper(a.Notes), "PLASMID") || strings.Contains(strings.ToUpper(a.Notes), "CIRCULAR") {
					newseq.Plasmid = true
				}
				newconstruct.Vector = newseq
			} else {
				newseq.Nm = b.Label
				newseq.Seq = b.Sequence
				if strings.Contains(a.Notes, "Plasmid") {
					newseq.Plasmid = true
				}
				newconstruct.Partsinorder = append(newconstruct.Partsinorder, newseq)
			}
		}
		construct_list = append(construct_list, newconstruct)
	}
	return construct_list, nil
}
