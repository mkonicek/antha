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

package wtype

import (
	"strings"
)

// feature class options based on genbank conventions
const (
	ORF           = "ORF"
	CDS           = "CDS"
	GENE          = "gene"
	MISC_FEATURE  = "misc_feature"
	PROMOTER      = "promoter"
	TRNA          = "tRNA"
	RRNA          = "rRNA"
	NCRNA         = "ncRNA"
	REGULATORY    = "regulatory"
	REPEAT_REGION = "repeat_region"

	/*
		Promoter
		Ribosomebindingsite
		TranslationInitSite
		Origin
		Marker
		Misc*/
)

// Feature describes a feature within a sequence, it's position in the sequence, in human friendly format,
// a protein sequence if applicablae and a class.
// Use the MakeFeature and AddFeature functions from the sequences packages.
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
type Feature struct {
	Name          string `json:"name"`
	Class         string `json:"class	"` //  defined by constants above
	Reverse       bool   `json:"reverse"`
	StartPosition int    `json:"start_position"` // in human friendly format
	EndPosition   int    `json:"end_position"`   // in human friendly format
	DNASeq        string `json:"dna_seq"`
	Protseq       string `json:"prot_seq"`
}

// DNASequence returns the linear DNA sequence of the feature.
func (f Feature) DNASequence() DNASequence {
	return DNASequence{Nm: f.Name, Seq: f.DNASeq}
}

const (
	// Option to feed into coordinates method.
	// HUMANFRIENDLY returns a sequence PositionPair's start and end positions in a human friendly format
	// i.e. in a Sequence "ATGTGTTG" position 1 is A, 2 is T.
	HUMANFRIENDLY = "humanFriendly"

	// Option to feed into coordinates method.
	// CODEFRIENDLY returns a sequence PositionPair's start and end positions in a code friendly format
	// i.e. in a Sequence "ATGTGTTG" position 0 is A, 1 is T.
	CODEFRIENDLY = "codeFriendly"

	// Option to feed into coordinates method.
	// IGNOREDIRECTION is a constant to specify that direction of a feature position
	// should be ignored when returning start and end positions of a feature.
	// If selected, the start position will be the first position at which the feature is encountered regardless of orientation.
	IGNOREDIRECTION = "ignoreDirection"
)

// ignores case
func containsString(slice []string, testString string) bool {

	upper := func(s string) string { return strings.ToUpper(s) }

	for _, str := range slice {
		if upper(str) == upper(testString) {
			return true
		}
	}
	return false
}

// Coordinates returns the start and end positions of the feature
// by default this will return the start position followed by the end position in human friendly format
// Available options are:
// HUMANFRIENDLY returns a sequence PositionPair's start and end positions in a human friendly format
// i.e. in a Sequence "ATGTGTTG" position 1 is A, 2 is T.
// CODEFRIENDLY returns a sequence PositionPair's start and end positions in a code friendly format
// i.e. in a Sequence "ATGTGTTG" position 0 is A, 1 is T.
// IGNOREDIRECTION is a constant to specify that direction of a feature position
// should be ignored when returning start and end positions of a feature.
// If selected, the start position will be the first position at which the feature is encountered regardless of orientation.
func (feat *Feature) Coordinates(options ...string) (start, end int) {
	start, end = feat.StartPosition, feat.EndPosition
	if containsString(options, CODEFRIENDLY) {
		start--
		end--
	}
	if containsString(options, IGNOREDIRECTION) {
		if start > end {
			return end, start
		}
	}
	return start, end
}

// Start returns the start position of the Feature
// by default this will return a directional human friendly position
func (f *Feature) Start(options ...string) int {
	start, _ := f.Coordinates(options...)
	return start
}

// End returns the end position of the Feature
// by default this will return a directional human friendly position
func (f *Feature) End(options ...string) int {
	_, end := f.Coordinates(options...)
	return end
}

// FeatureNames returns a list of all feature names in the sequence
func (annotated DNASequence) FeatureNames() (featurenames []string) {

	featurenames = make([]string, 0)
	for _, feature := range annotated.Features {
		featurenames = append(featurenames, feature.Name)
	}
	return
}

// GetFeatureByName returns all features found which match the specified feature name.
// Searches are not case sensitive.
func (annotated DNASequence) GetFeatureByName(featureName string) (returnedFeatures []Feature) {

	for _, feature := range annotated.Features {
		if strings.Contains(strings.ToUpper(feature.Name), strings.ToUpper(featureName)) {
			returnedFeatures = append(returnedFeatures, feature)
		}

	}
	return
}

// delete this and use sequences.AddFeatures()

func Annotate(dnaseq DNASequence, features []Feature) (annotated DNASequence) {
	annotated = dnaseq
	annotated.Features = features
	return
}
