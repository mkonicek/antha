// antha/AnthaStandardLibrary/Packages/enzymes/Annotatedseq.go: Part of the Antha language
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

// Package sequences is for interacting with and manipulating biological sequences; in extension to methods available in wtype
package sequences

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/antha-lang/antha/antha/anthalib/wtype"
)

// AddFeature adds a feature to a DNASequence
// The positions will be added automatically; if more than one matching sequence is found, multiple features will be added.
func AddFeature(annotated *wtype.DNASequence, newFeature wtype.Feature) {

	featureSeq := newFeature.DNASequence()

	positions := findSeq(annotated, &featureSeq)

	var features []wtype.Feature
	features = append(features, annotated.Features...)
	for _, featPosition := range positions.Positions {
		newFeature.StartPosition, newFeature.EndPosition = featPosition.Coordinates()
		features = append(features, newFeature)
	}
	annotated.Features = features
}

// RemoveFeatures clears all existing feature annotations from a sequence.
func RemoveFeatures(dnaSeq *wtype.DNASequence) {
	var noFeatures []wtype.Feature
	dnaSeq.Features = noFeatures
}

// SetFeatures replaces any existing feature annotations of the DNASequence with the
// features specified.
func SetFeatures(dnaSeq *wtype.DNASequence, features []wtype.Feature) {

	RemoveFeatures(dnaSeq)

	for _, feature := range features {
		AddFeature(dnaSeq, feature)
	}

	dnaSeq.Features = features
}

// AddFeatures adds features to the existing features of the DNASequence.
func AddFeatures(dnaSeq *wtype.DNASequence, features []wtype.Feature) {

	for _, feature := range features {
		AddFeature(dnaSeq, feature)
	}

	dnaSeq.Features = features
}

// ORFs2Features converts a set of ORFs into a set of features
func ORFs2Features(orfs []ORF) (features []wtype.Feature) {

	features = make([]wtype.Feature, 0)

	for i, orf := range orfs {
		// currently just names each orf + number of orf. Add Rename orf function and sort by struct field function to run first to put orfs in order
		reverse := false
		if strings.EqualFold(orf.Direction, "REVERSE") {
			reverse = true
		}

		feature := wtype.Feature{
			Name:          "orf" + strconv.Itoa(i),
			Class:         "orf",
			Reverse:       reverse,
			StartPosition: orf.StartPosition,
			EndPosition:   orf.EndPosition,
			DNASeq:        orf.DNASeq,
			Protseq:       orf.ProtSeq,
		}
		features = append(features, feature)
	}
	return
}

// MakeFeature constructs an annotated feature to be added to a sequence.
// The feature will be defined by it's class and it's position in the sequence once added to a sequence using AddFeature.
// A protein sequence can be specified if appropriate.
// valid class fields are:
/*
	ORF = "ORF"
	CDS = "CDS"
	GENE = "gene"
	MISC_FEATURE = "misc_feature"
	PROMOTER = "promoter"
	TRNA = "tRNA"
	RRNA = "rRNA"
	NCRNA = "ncRNA"
	REGULATORY = "regulatory"
	REPEAT_REGION = "repeat_region"
*/
// valid sequence types entries are:
// "aa" = amino acid/ protein sequence
// "dna" = DNA sequence
// "rna = "RNA sequence
// Use the AddFeature function to add the feature to a DNASequence such that the positions are added correctly.
func MakeFeature(name string, seq string, start int, end int, sequencetype string, class string, reverse string) (feature wtype.Feature) {
	feature.Name = name
	feature.DNASeq = strings.ToUpper(seq)
	feature.Class = class
	if strings.ToLower(reverse) == "reverse" {
		feature.Reverse = true
	}

	if strings.ToLower(sequencetype) == "aa" {
		feature.DNASeq = RevTranslatetoNstring(seq)
		feature.Protseq = seq
		feature.StartPosition = start
		feature.EndPosition = end
	} else {
		if strings.ToLower(sequencetype) == "rna" {
			seq = rnaToDNA(seq)
		}

		if !feature.Reverse {
			feature.DNASeq = seq
		}
		if feature.Reverse {
			seq = wtype.RevComp(seq)
			feature.DNASeq = seq
		}
		feature.StartPosition = start
		feature.EndPosition = end

		if feature.Class == "gene" || feature.Class == "CDS" {
			orf, orftrue := FindORF(seq)
			if orftrue {
				feature.Protseq = orf.ProtSeq

			}
		}
	}

	if feature.Class == "ORF" || feature.Class == "orf" {
		orf, orftrue := FindORF(seq)
		if orftrue {
			feature.Protseq = orf.ProtSeq
		}
	}
	return feature
}

func rnaToDNA(rnaSeq string) (dnaSeq string) {

	var newSeq []string

	for _, letter := range rnaSeq {
		if strings.ToUpper(string(letter)) == "U" {
			letter = rune('T')
		}
		newSeq = append(newSeq, string(letter))
	}

	dnaSeq = strings.Join(newSeq, "")

	return
}

// MakeAnnotatedSeq makes a DNA sequence adding the specified features with their correct positions in the sequence specified in human friendly format.
func MakeAnnotatedSeq(name string, seq string, circular bool, features []wtype.Feature) (annotated wtype.DNASequence, err error) {
	annotated.Nm = name
	annotated.Seq = seq
	annotated.Plasmid = circular

	var newFeatures []wtype.Feature

	for _, feature := range features {
		featureDNASequence := feature.DNASequence()
		featurePositionsFound := FindAll(&annotated, &featureDNASequence)

		if len(featurePositionsFound.Positions) > 0 {
			for _, positionPair := range featurePositionsFound.Positions {
				feature.StartPosition, feature.EndPosition = positionPair.HumanFriendly()
				feature.Reverse = positionPair.Reverse
				newFeatures = append(newFeatures, feature)
			}
		} else if len(featurePositionsFound.Positions) == 0 {
			err = fmt.Errorf("%s sequence %s not found in sequence %s ", feature.Name, feature.DNASeq, seq)
		}

	}
	annotated.Features = newFeatures
	return
}
