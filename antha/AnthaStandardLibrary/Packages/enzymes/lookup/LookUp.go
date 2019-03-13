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

// Package lookup enables looking up restriction enzyme properties from name.
package lookup

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/asset"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/rebase"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
)

// TypeIIs looks up a TypeIIs enzyme and returns the result as a TypeIIs enzyme object.
// An error is returned if no enzyme is found with the specified name.
func TypeIIs(name string) (enzyme wtype.TypeIIs, err error) {
	enz, err := RestrictionEnzyme(name)
	if err != nil {
		return enzyme, err
	}
	enzyme, err = wtype.ToTypeIIs(enz)
	return
}

// RestrictionEnzyme looks up a Restriction enzyme and returns the result as a RestrictionEnzyme object.
// An error is returned if no enzyme is found with the specified name.
func RestrictionEnzyme(name string) (enzyme wtype.RestrictionEnzyme, err error) {

	if name == "" {
		return enzyme, fmt.Errorf(`Error! Enzyme has been specified as "", check the enzymes listed in parameters for "" `)
	}

	enzymes, err := asset.Asset("rebase/type2.txt")
	if err != nil {
		return
	}

	rebaseFh := bytes.NewReader(enzymes)

	for _, record := range rebase.Parse(rebaseFh) {
		/*plasmidstatus := "FALSE"
		seqtype := "DNA"
		class := "not specified"*/

		if strings.EqualFold(strings.TrimSpace(record.Name()), strings.TrimSpace(name)) {
			enzyme = record
			return enzyme, nil
		}

	}

	return enzyme, fmt.Errorf("No enzyme %s found", name)
}

// FindEnzymesofClass returns a list of all RestrictionEnzymes belonging to the requested class.
// Example class arguments are typeII and typeIIs.
// If an invalid class is specified an empty list will be returned.
func FindEnzymesofClass(class string) (enzymelist []wtype.RestrictionEnzyme) {
	enzymes, err := asset.Asset("rebase/type2.txt")
	if err != nil {
		return
	}

	rebaseFh := bytes.NewReader(enzymes)

	for _, record := range rebase.Parse(rebaseFh) {
		if strings.EqualFold(record.Class, class) {
			//RecognitionSeqs = append(RecognitionSeqs, record)
			enzymelist = append(enzymelist, record)
		}
	}
	return enzymelist
}

// FindEnzymeNamesofClass returns a list of all RestrictionEnzyme names belonging to the requested class.
// Example class arguments are typeII and typeIIs.
// If an invalid class is specified an empty list will be returned.
func FindEnzymeNamesofClass(class string) (enzymelist []string) {
	for _, enzyme := range FindEnzymesofClass(class) {
		enzymelist = append(enzymelist, enzyme.Nm)
	}
	return
}
