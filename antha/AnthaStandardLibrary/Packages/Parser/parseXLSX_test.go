// antha/AnthaStandardLibrary/Packages/Parser/RebaseParser.go: Part of the Antha language
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
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/enzymes"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	//"github.com/antha-lang/antha/antha/anthalib/wunit"
	"testing"
	"reflect"
)

type parseXLSXTest struct {
    designFileName    string
	expectedAssemblies []enzymes.Assemblyparameters
	expectedErrMessage error
}

var tests = []parseXLSXTest{
	parseXLSXTest{
		designFileName:  "xlsxParserTestFile.xlsx",
		expectedAssemblies: []enzymes.Assemblyparameters{
			enzymes.Assemblyparameters{
                Constructname: "Assembly1",
				Enzymename: "SapI",
				Vector: wtype.DNASequence{
					Nm: "vector",
					Seq: "AAACTGT",
					Plasmid: true,
				},
				Partsinorder: []wtype.DNASequence{
					wtype.DNASequence{
						Nm: "part1",
						Seq: "AAACTGT",
						Plasmid: true,
					},
					wtype.DNASequence{
						Nm: "part2",
						Seq: "AAACTGT",
						Plasmid: true,
					},
					wtype.DNASequence{
						Nm: "part3",
						Seq: "AAACTGT",
						Plasmid: true,
					},
				}, 
            },
		},
	},
}

func TestParseExcel(t *testing.T) {

	for _, test := range tests {

		assemblies, err := ParseExcel(test.designFileName)

		if err != nil {
			if err != test.expectedErrMessage {
				t.Error(
					err.Error(),
				)
			}
		} else if test.expectedErrMessage != nil {
			t.Error(
					"For", test.designFileName, "\n",
					"expected Error message:", test.expectedErrMessage, "\n",
					"got no error",
				)
		}

		if !reflect.DeepEqual(assemblies,test.expectedAssemblies) {
			t.Error(
				"for test", test.designFileName, "\n",
				"expected: ", test.expectedAssemblies, "\n",
				"got", fmt.Sprintf("%+v", assemblies), "\n",
			)
		}
	}
}