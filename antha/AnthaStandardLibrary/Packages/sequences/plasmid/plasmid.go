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

// Package plasmid checks for common plasmid features in a test DNA sequence.
package plasmid

import (
	"strings"

	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/sequences"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/sequences/parse/fasta"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
)

// Based on Plasmapper annotation system

// fasta header format:
// > Name(Abbr)[Type]{X},Length, Y
// AATCTCT....

var (
	plasmapperTypeCodes = map[string]string{
		"ORIGIN":    "[ORI]",
		"SELECTION": "[SEL]",
	}
)

var (
	plasmapperfile []byte = []byte(commonFeatures)
)

// FeatureMap stores a map of features using feature type as key.
type FeatureMap map[string][]wtype.DNASequence

func makeFeatureMap(contents []byte) (featuremap FeatureMap, err error) {

	featuremap = make(FeatureMap)

	matchingseqs := make([]wtype.DNASequence, 0)

	seqs, err := fasta.FastaContentstoDNASequences(contents)
	if err != nil {
		return featuremap, err
	}

	for key, value := range plasmapperTypeCodes {
		for _, seq := range seqs {
			if strings.Contains(seq.Nm, value) {
				matchingseqs = append(matchingseqs, seq)
			}
		}
		// add to map
		featuremap[key] = matchingseqs
		//reset
		matchingseqs = make([]wtype.DNASequence, 0)
	}
	return
}

// MakePlasmapperFeatures initialises a FeatureMap according to common plasmid features, mostly from plasmapper.
func MakePlasmapperFeatures() (featuremap FeatureMap, err error) {

	featuremap, err = makeFeatureMap(plasmapperfile)
	return
}

// ValidPlasmid evaluates whether a test sequence is circular, contains any origins of replications and selection markers.
// The features are evaluated for exact matches against a restricted list of common features defined as the variable commonfeatures.
func ValidPlasmid(sequence wtype.DNASequence) (plasmid bool, oris []string, selectionmarkers []string, err error) {
	if sequence.Plasmid {
		plasmid = true
	}
	featuremap, err := MakePlasmapperFeatures()
	if err != nil {
		return plasmid, []string{}, []string{}, err
	}

	seqfeatures := sequence.Features

	for _, feature := range seqfeatures {
		if feature.Class == "Origin" {
			oris = append(oris, feature.Name)
		}
		if feature.Class == "Marker" {
			selectionmarkers = append(selectionmarkers, feature.Name)
		}
	}

	var oriseqs []wtype.DNASequence
	oriseqs = append(oriseqs, featuremap["ORIGIN"]...)
	for _, oriseq := range oriseqs {
		if len(sequence.Sequence()) >= len(oriseq.Sequence()) {
			if len(sequences.FindAll(&sequence, &oriseq).Positions) > 0 {
				oris = append(oris, oriseq.Name())
			}
		}
	}

	var markerseqs []wtype.DNASequence
	markerseqs = append(markerseqs, featuremap["SELECTION"]...)

	for _, markerseq := range markerseqs {
		if len(sequence.Sequence()) >= len(markerseq.Sequence()) {
			if len(sequences.FindAll(&sequence, &markerseq).Positions) > 0 {
				selectionmarkers = append(selectionmarkers, markerseq.Name())
			}
		}
	}

	return
}
