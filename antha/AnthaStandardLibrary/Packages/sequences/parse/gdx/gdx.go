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

// package gdx converts DNA sequence files in .gdx format into a set of DNA sequences.
package gdx

import (
	"encoding/xml"

	parse "github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/Parser"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
)

// GDXToDNASequence parses DNA sequence files in .gdx format into a set of DNA sequences of type []wtype.DNASequence
func GDXToDNASequence(sequenceFile wtype.File) (partsList []wtype.DNASequence, err error) {
	data, err := sequenceFile.ReadAll()
	if err != nil {
		return nil, err
	}

	var gdx parse.Project
	err = xml.Unmarshal(data, &gdx)

	if err != nil {
		return partsList, err
	}

	partsList = make([]wtype.DNASequence, 0)

	for _, construct := range gdx.DesignConstruct {
		for _, dnaElement := range construct.DNAElements {
			var newseq wtype.DNASequence
			for counter := 0; counter < len(construct.DNAElements); counter++ {
				newseq.Nm = dnaElement.Label
				newseq.Seq = dnaElement.Sequence
				if construct.Plasmid == "true" {
					newseq.Plasmid = true
				}
				partsList = append(partsList, newseq)
			}
		}
	}

	return partsList, err
}
