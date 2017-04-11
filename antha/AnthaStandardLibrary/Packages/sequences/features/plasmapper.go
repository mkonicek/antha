// plasmapper
package features

import (
	"strings"

	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/Parser"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/sequences"

	"github.com/antha-lang/antha/antha/anthalib/wtype"
)

// Based on Plasmapper annotation system

// fasta header format:
// > Name(Abbr)[Type]{X},Length, Y
// AATCTCT....

var (
	PlasmapperTypeCodes = map[string]string{
		"ORIGIN":    "[ORI]",
		"SELECTION": "[SEL]",
	}
)

// add filter registry option
var (
	iGemRegistryCodes = map[string][]string{
		"ORIGIN":    []string{"Origin", "ori"},
		"SELECTION": []string{"resistance"},
	}
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

var commonFeatures string = ``

var (
	plasmapperfile []byte = []byte(commonFeatures)
)

type FeatureMap map[string][]wtype.DNASequence

func MakeFeatureMap(contents []byte) (featuremap FeatureMap, err error) {

	featuremap = make(FeatureMap)

	matchingseqs := make([]wtype.DNASequence, 0)

	seqs, err := parser.FastatoDNASequences(contents)
	if err != nil {
		return featuremap, err
	}

	for key, value := range PlasmapperTypeCodes {
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

// should be initialised to produce a feature map from different sources
func MakePlasmapperMap() (featuremap FeatureMap, err error) {

	featuremap, err = MakeFeatureMap(plasmapperfile) //-----------------------------------------------------------------------------------
	return
}

func ValidPlasmid(sequence wtype.DNASequence) (plasmid bool, ori bool, selectionmarker bool, err error) {
	if sequence.Plasmid == true {
		plasmid = true
	}
	featuremap, err := MakePlasmapperMap()
	if err != nil {
		return plasmid, false, false, err
	}

	seqfeatures := sequence.Features

	for _, feature := range seqfeatures {
		if feature.Class == "Origin" {
			ori = true
		}
		if feature.Class == "Marker" {
			selectionmarker = true
		}
	}

	oriseqs := make([]string, 0)

	for _, oriseq := range featuremap["ORIGIN"] {
		oriseqs = append(oriseqs, oriseq.Seq)
	}

	if len(sequences.FindSeqsinSeqs(sequence.Seq, oriseqs)) > 0 {
		ori = true
	}

	markerseqs := make([]string, 0)

	for _, markerseq := range featuremap["SELECTION"] {
		markerseqs = append(markerseqs, markerseq.Seq)
	}

	if len(sequences.FindSeqsinSeqs(sequence.Seq, markerseqs)) > 0 {
		selectionmarker = true
	}

	return
}

/*
func CompatibleParts(plasmids []wtype.DNASequence) (Compatible bool, err error) {

}
*/
