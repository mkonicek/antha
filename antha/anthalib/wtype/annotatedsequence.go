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
	"fmt"
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
	StartPosition int    `json:"start_position"`
	EndPosition   int    `json:"end_position"`
	DNASeq        string `json:"dna_seq"`
	Protseq       string `json:"prot_seq"`
	//Status        string
}

func (feat *Feature) Coordinates() (pair []int) {
	pair[0] = feat.StartPosition
	pair[1] = feat.EndPosition
	return
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

func MakeAnnotatedSeq(name string, seq string, circular bool, features []Feature) (annotated DNASequence, err error) {
	annotated.Nm = name
	//annotated.Seq.Nm = name
	annotated.Seq = seq //.Seq = seq
	annotated.Plasmid = circular

	for _, feature := range features {
		if strings.Contains(seq, feature.DNASeq) {
			feature.StartPosition = strings.Index(seq, feature.DNASeq)
			feature.EndPosition = feature.EndPosition + feature.StartPosition
		} else if strings.Contains(seq, RevComp(feature.DNASeq)) {
			feature.StartPosition = strings.Index(seq, feature.DNASeq)
			feature.EndPosition = feature.EndPosition + feature.StartPosition
			err = fmt.Errorf(feature.Name, " Feature only found in reverse direction")
		} else {
			err = fmt.Errorf(feature.Name, " not found in sequence")
		}
	}
	annotated.Features = features
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
