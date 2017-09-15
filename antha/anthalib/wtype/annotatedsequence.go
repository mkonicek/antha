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

/*
const (
	orf = iota
	Promoter
	Ribosomebindingsite
	TranslationInitSite
	Origin
	Marker
	Misc
)
*/

type Feature struct {
	Name          string `json:"name"`
	Class         string `json:"class	"` //int // defined by constants above
	Reverse       bool   `json:"reverse"`
	StartPosition int    `json:"start_position"` // in human friendly format
	EndPosition   int    `json:"end_position"`   // in human friendly format
	DNASeq        string `json:"dna_seq"`
	Protseq       string `json:"prot_seq"`
	//Synonyms      map[string]string `json:"synonyms"`
	//Status        string
}

func (f Feature) DNASequence() DNASequence {
	return DNASequence{Nm: f.Name, Seq: f.DNASeq}
}

const (
	HUMANFRIENDLY   = "humanFriendly"
	CODEFRIENDLY    = "codeFriendly"
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
// by default this will return the start position followed by the end position
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

func (annotated DNASequence) FeatureNames() (featurenames []string) {

	featurenames = make([]string, 0)
	for _, feature := range annotated.Features {
		featurenames = append(featurenames, feature.Name)
	}
	return
}

func (annotated DNASequence) FeatureStart(featurename string) (featureStart int) {

	for _, feature := range annotated.Features {
		if feature.Name == featurename {
			featureStart = feature.StartPosition
			return
		}

	}
	return
}

func (annotated DNASequence) FeatureEnd(featurename string) (featureEnd int) {

	for _, feature := range annotated.Features {
		if feature.Name == featurename {
			featureEnd = feature.EndPosition
			return
		}

	}
	return
}

func (annotated DNASequence) GetFeatureByName(featurename string) (returnedfeature *Feature) {

	for _, feature := range annotated.Features {
		if strings.Contains(strings.ToUpper(feature.Name), strings.ToUpper(featurename)) {
			returnedfeature = &feature
			return
		}

	}
	return
}

func Annotate(dnaseq DNASequence, features []Feature) (annotated DNASequence) {
	annotated = dnaseq
	annotated.Features = features
	return
}

func AddFeatures(annotated DNASequence, features []Feature) (updated DNASequence) {

	for _, feature := range features {
		annotated.Features = append(annotated.Features, feature)
	}
	return
}

func ConcatenateFeatures(name string, featuresinorder []Feature) (annotated DNASequence) {

	annotated.Nm = name
	//annotated.Seq.Nm = name
	annotated.Seq = featuresinorder[0].DNASeq
	annotated.Features = make([]Feature, 0)
	annotated.Features = append(annotated.Features, featuresinorder[0])
	for i := 1; i < len(featuresinorder); i++ {
		nextfeature := featuresinorder[i]
		nextfeature.StartPosition = nextfeature.StartPosition + annotated.Features[i-1].EndPosition
		nextfeature.EndPosition = nextfeature.EndPosition + annotated.Features[i-1].EndPosition
		annotated.Seq = annotated.Seq + featuresinorder[i].DNASeq
		annotated.Features = append(annotated.Features, nextfeature)
	}
	return annotated
}
