// antha/AnthaStandardLibrary/Packages/Inventory/Inventory.go: Part of the Antha language
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

package Inventory

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"

	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/sequences/parse/fasta"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/sequences/parse/genbank"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
)

// PartSource is any type which has a GetSequences method
type PartSource interface {
	GetSequences() ([]wtype.DNASequence, error)
}

// SequenceMap is a map of DNA Sequences which obeys the PartSource interface.
type SequenceMap map[string]wtype.DNASequence

// GetSequences returns all sequences and any errors which occur.
func (ex SequenceMap) GetSequences() (partslist []wtype.DNASequence, err error) {
	var sortedKeys []string

	for key := range ex {
		sortedKeys = append(sortedKeys, key)
	}

	sort.Strings(sortedKeys)

	for _, key := range sortedKeys {
		partslist = append(partslist, ex[key])
	}
	return partslist, nil
}

// SequenceSet is a set of DNA Sequences which obeys the PartSource interface.
type SequenceSet []wtype.DNASequence

// GetSequences returns all sequences and any errors which occur.
func (seqs SequenceSet) GetSequences() (partslist []wtype.DNASequence, err error) {
	for _, value := range seqs {
		partslist = append(partslist, value)
	}
	return partslist, nil
}

// FileSet is a set of files which obeys the PartSource interface.
// Any fasta or Genbank files will return DNASequences.
type FileSet []wtype.File

// GetSequences returns all sequences and any errors which occur.
// Any fasta or Genbank files will return DNASequences.
func (fs FileSet) GetSequences() (partslist []wtype.DNASequence, err error) {
	var errs []string

	for _, file := range fs {

		filename := file.Name
		if strings.EqualFold(filepath.Ext(filename), ".fasta") || strings.EqualFold(filepath.Ext(filename), ".fa") {
			sequences, err := fasta.FastaToDNASequences(file)
			if err != nil {
				errs = append(errs, err.Error())
			}
			partslist = append(partslist, sequences...)
		} else if strings.EqualFold(filepath.Ext(filename), ".gb") || strings.EqualFold(filepath.Ext(filename), ".gbk") {
			seq, err := genbank.GenbankToAnnotatedSeq(file)
			if err != nil {
				errs = append(errs, err.Error())
			}
			partslist = append(partslist, seq)

		} else {
			errs = append(errs, fmt.Sprintf("cannot return DNA Sequences from file %s. Only Fasta (.fasta) and Genbank (.gb) files are valid", filename))
		}
	}

	if len(errs) > 0 {
		return partslist, fmt.Errorf(strings.Join(errs, "\n"))
	}

	return partslist, nil
}
