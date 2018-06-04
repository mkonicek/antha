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

// Package for designing oligos
package oligos

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/search"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/sequences"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
)

type Primer struct {
	wtype.DNASequence
	Length      int
	GCContent   float64
	Reverse     bool
	MeltingTemp wunit.Temperature
}

// checks for overlap between sequences (not including mismatches)
func OverlapCheck(seq1 string, seq2 string) (maxpercentOverlapofsmallest float64, maxnumberofbpOverlap int, overlappingseq string) {

	if len(seq1) > 0 && len(seq2) > 0 {

		var refseq string
		var testseq string
		var partoftestseq string
		var biggestsofar int
		var overlapsofar string

		if len(seq1) >= len(seq2) {
			refseq = seq1
			testseq = seq2
		} else {
			refseq = seq2
			testseq = seq1
		}

		for i := 0; i < len(testseq); i++ {
			partoftestseq = testseq[i:]

			if strings.Contains(refseq, partoftestseq) && len(partoftestseq) > biggestsofar {
				biggestsofar = len(partoftestseq)
				overlapsofar = partoftestseq
			}
		}

		for j := 0; j < len(testseq); j++ {
			partoftestseq = testseq[:len(testseq)-j]

			if strings.Contains(refseq, partoftestseq) && len(partoftestseq) > biggestsofar {
				biggestsofar = len(partoftestseq)
				overlapsofar = partoftestseq
			}
		}

		maxnumberofbpOverlap = biggestsofar
		overlappingseq = overlapsofar

		maxpercentOverlapofsmallest = float64(maxnumberofbpOverlap) / float64(len(testseq))
	}
	return
}

/*
BasicMeltingTemp calculates the expected melting temperature of a DNASequence.

The two standard approximation calculations are used.
For sequences less than 14 nucleotides the formula is
Tm= (wA+xT) * 2 + (yG+zC) * 4

where w,x,y,z are the number of the bases A,T,G,C in the sequence, respectively
(from Marmur,J., and Doty,P. (1962) J Mol Biol 5:109-118 [PubMed]).

For sequences longer than 13 nucleotides, the equation used is
Tm= 64.9 +41*(yG+zC-16.4)/(wA+xT+yG+zC)

See Wallace,R.B., Shaffer,J., Murphy,R.F., Bonner,J., Hirose,T., and Itakura,K. (1979) Nucleic Acids Res 6:3543-3557 (Abstract)
and Sambrook,J., and Russell,D.W. (2001) Molecular Cloning: A Laboratory Manual. Cold Spring Harbor Laboratory Press; Cold Spring Harbor, NY. (CHSL Press)

ASSUMPTIONS:
Both equations assume that the annealing occurs under the standard conditions of 50 nM primer, 50 mM Na+, and pH 7.0.

*/
func BasicMeltingTemp(primersequence wtype.DNASequence) (meltingtemp wunit.Temperature) {

	primerseq := primersequence.Sequence()

	primerseq = strings.ToUpper(primerseq)

	g := strings.Count(primerseq, "G")
	c := strings.Count(primerseq, "C")
	a := strings.Count(primerseq, "A")
	t := strings.Count(primerseq, "T")

	var mt float64

	if len(primerseq) < 14 {

		mt = float64((a+t)*2 + (g+c)*4)

	} else {
		mt = 64.9 + 41.0*(float64(g+c)-16.4)/float64(a+t+c+g)
	}

	meltingtemp = wunit.NewTemperature(mt, "℃")
	return
}

// define region in DNA sequence
// this is directionless and does not check for reverse complement
// assumes positions stated are in user format (i.e. first position is 1 and not 0)
func DNAregion(sequence wtype.DNASequence, startposition int, endposition int) (region wtype.DNASequence) {

	dnaseq := sequence.Sequence()

	//define region in sequence to create primer.
	if endposition > startposition {
		// fix properly!!!
		if startposition == 0 {
			startposition = 1
		}

		if endposition > len(dnaseq) {
			message := fmt.Sprint("endposition ", endposition, " exceeds length of sequence ", sequence.Name(), " Length: ", len(dnaseq))
			panic(message)
			//endposition = len(dnaseq) - 1
		}
		region = wtype.MakeLinearDNASequence("region"+strconv.Itoa(startposition)+":"+strconv.Itoa(endposition), dnaseq[startposition-1:endposition])
	} else if endposition < startposition && sequence.Plasmid {
		region = wtype.MakeLinearDNASequence("region"+strconv.Itoa(startposition)+":"+strconv.Itoa(endposition), dnaseq[startposition-1:]+dnaseq[:endposition])
	} else if endposition < startposition && !sequence.Plasmid {
		message := fmt.Sprint("DNA Region start position cannot be larger than end position for linear dna fragments. startposition: ", startposition, "endposition", endposition, "Sequence:", sequence)
		panic(message)
	}
	return

}

// Takes defined region and makes an oligosequence between a defined minimum and maximum length
// with a melting temperature between a defined minimum and maximum and a maximum GC content ( between 0 and 1).
// function finds oligo by starting at position 0 and making sequence of the minimum length, calculating parameters
// and if they do not match then adds one basepair to end of sequence until the maximum length is reached.
// if still unsuccessful, the function begins again at position 1 and cycles through until a matching oligo sequence is found.
// overlapthresholdwithseqstoavoid allows maximum permissable partial overlap to be specified by the user, if set to -1 any overlap is tolerated
func FWDOligoSeq(seq wtype.DNASequence, maxGCcontent float64, minlength int, maxlength int, minmeltingtemp wunit.Temperature, maxmeltingtemp wunit.Temperature, seqstoavoid []string, overlapthresholdwithseqstoavoid int) (oligoseq Primer, err error) {

	if maxlength > len(seq.Sequence()) {
		return oligoseq, fmt.Errorf("Sequence %s %s too small to design primer for or max length of primer %d too long", seq.Nm, seq.Seq, maxlength)
	}

	region := strings.ToUpper(seq.Sequence())

	for start := 0; start < maxlength; start++ {

		for end := minlength + start; end <= start+maxlength; end++ {

			var overlapthresholdfail bool

			tempoligoseq := region[start:end]

			ssoligo := wtype.MakeSingleStrandedDNASequence("oligo", tempoligoseq)

			temppercentage := sequences.GCcontent(tempoligoseq)

			meltingtemp := BasicMeltingTemp(ssoligo)

			bindingsites := CheckNonSpecificBinding(seq, ssoligo)

			if len(seqstoavoid) > 0 && overlapthresholdwithseqstoavoid > 0 {
				for _, seq := range seqstoavoid {
					_, overlap, _ := OverlapCheck(tempoligoseq, seq)

					if overlap > overlapthresholdwithseqstoavoid {
						overlapthresholdfail = true
					}
				}
			}

			if temppercentage <= maxGCcontent && minmeltingtemp.SIValue() < meltingtemp.SIValue() && maxmeltingtemp.SIValue() > meltingtemp.SIValue() && bindingsites == 1 && !search.PartialInStrings(seqstoavoid, tempoligoseq, search.IgnoreCase) && !overlapthresholdfail {
				oligoseq.DNASequence = wtype.MakeSingleStrandedDNASequence("Primer", tempoligoseq)
				oligoseq.GCContent = temppercentage
				oligoseq.Length = len(tempoligoseq)
				oligoseq.MeltingTemp = meltingtemp
				err = nil
				return

			} else {
				var combinedErrors []string
				var i bool = true
				var primerType, seqName string = "FORWARD", seq.Name()
				switch i {
				case temppercentage >= maxGCcontent:
					combinedErrors = append(combinedErrors, gcContentErrorString(primerType, seqName, temppercentage, maxGCcontent))
				case minmeltingtemp.SIValue() > meltingtemp.SIValue():
					combinedErrors = append(combinedErrors, meltingTempErrorString(primerType, seqName, "minimum", meltingtemp.SIValue(), minmeltingtemp.SIValue()))
				case maxmeltingtemp.SIValue() < meltingtemp.SIValue():
					combinedErrors = append(combinedErrors, meltingTempErrorString(primerType, seqName, "maximum", meltingtemp.SIValue(), maxmeltingtemp.SIValue()))
				case bindingsites > 1:
					combinedErrors = append(combinedErrors, bindingSiteErrorString(primerType, seqName, bindingsites))
				case search.PartialInStrings(seqstoavoid, tempoligoseq, search.IgnoreCase):
					combinedErrors = append(combinedErrors, primerErrorString(primerType, seqName, " that contain the specified sequences to avoid.", " removing these from the parameters."))
				case overlapthresholdfail:
					combinedErrors = append(combinedErrors, primerErrorString(primerType, seqName, " that violate the overlap threshold.", " adjusting this parameter."))
				}
				err = fmt.Errorf(strings.Join(combinedErrors, "\n"))
			}
		}
	}
	return
}

func REVOligoSeq(seq wtype.DNASequence, maxGCcontent float64, minlength int, maxlength int, minmeltingtemp wunit.Temperature, maxmeltingtemp wunit.Temperature, seqstoavoid []string, overlapthresholdwithseqstoavoid int) (oligoseq Primer, err error) {

	//var start int
	//var end int

	// get the reverse complement of the region

	if maxlength > len(seq.Sequence()) {
		return oligoseq, fmt.Errorf("Sequence %s %s too small to design primer for or max length of primer %d too long", seq.Nm, seq.Seq, maxlength)
	}

	region := seq.Sequence()
	revregion := sequences.RevComp(region)

	for start := 0; start < maxlength; start++ {

		for end := minlength + start; end <= start+maxlength; end++ {

			var overlapthresholdfail bool

			tempoligoseq := revregion[start:end]

			ssoligo := wtype.MakeSingleStrandedDNASequence("oligo", tempoligoseq)

			temppercentage := sequences.GCcontent(tempoligoseq)

			meltingtemp := BasicMeltingTemp(ssoligo)

			bindingsites := CheckNonSpecificBinding(seq, ssoligo)

			if len(seqstoavoid) > 0 && overlapthresholdwithseqstoavoid > 0 {
				for _, seq := range seqstoavoid {
					_, overlap, _ := OverlapCheck(tempoligoseq, seq)

					if overlap > overlapthresholdwithseqstoavoid {
						overlapthresholdfail = true
					}
				}
			}

			if temppercentage <= maxGCcontent && minmeltingtemp.SIValue() < meltingtemp.SIValue() && maxmeltingtemp.SIValue() > meltingtemp.SIValue() && bindingsites == 1 && !search.PartialInStrings(seqstoavoid, tempoligoseq, search.IgnoreCase) && !overlapthresholdfail {
				oligoseq.DNASequence = wtype.MakeSingleStrandedDNASequence("Primer", tempoligoseq)
				oligoseq.GCContent = temppercentage
				oligoseq.Length = len(tempoligoseq)
				oligoseq.MeltingTemp = meltingtemp
				err = nil
				return

			} else {
				var combinedErrors []string
				var i bool = true
				var primerType, seqName string = "REVERSE", seq.Name()
				switch i {
				case temppercentage >= maxGCcontent:
					combinedErrors = append(combinedErrors, gcContentErrorString(primerType, seqName, temppercentage, maxGCcontent))
				case minmeltingtemp.SIValue() > meltingtemp.SIValue():
					combinedErrors = append(combinedErrors, meltingTempErrorString(primerType, seqName, "minimum", meltingtemp.SIValue(), minmeltingtemp.SIValue()))
				case maxmeltingtemp.SIValue() < meltingtemp.SIValue():
					combinedErrors = append(combinedErrors, meltingTempErrorString(primerType, seqName, "maximum", meltingtemp.SIValue(), maxmeltingtemp.SIValue()))
				case bindingsites > 1:
					combinedErrors = append(combinedErrors, bindingSiteErrorString(primerType, seqName, bindingsites))
				case search.PartialInStrings(seqstoavoid, tempoligoseq, search.IgnoreCase):
					combinedErrors = append(combinedErrors, primerErrorString(primerType, seqName, " that contain the specified sequences to avoid.", " removing these from the parameters."))
				case overlapthresholdfail:
					combinedErrors = append(combinedErrors, primerErrorString(primerType, seqName, " that violate the overlap threshold.", " adjusting this parameter."))
				}
				err = fmt.Errorf(strings.Join(combinedErrors, "\n"))
			}
		}
	}

	return
}

// primerErrorString formats a textual message in the form of an erorr related to primerdesign.
// primerorientation and sequence name are inputted, along with the error message i.e. "X not found in Y".
// and a proposed resolution i.e. "Please put X in Y".
func primerErrorString(primerOrientation, sequenceName, erorrString, resolutionString string) string {
	errorBeginning := fmt.Sprintf(" For %s could only generate %s primers", sequenceName, primerOrientation)
	resolutionBeginning := " Please try"
	return fmt.Sprintf("%s%s%s%s", errorBeginning, erorrString, resolutionBeginning, resolutionString)
}

// gcContentErrorString formats a textual message in the form of an error related to a primer exceeding the maximum specified GC content.
func gcContentErrorString(primerOrientation, sequenceName string, gcContent, maxGC float64) string {
	return primerErrorString(primerOrientation, sequenceName, fmt.Sprintf(" with GC Content (%f) greater than the maximum GC Content specified (%f).", gcContent, maxGC), " lowering this parameter, or selecting a less-GC rich region.")
}

// meltingTempErrorString formats a textual message in the form of an error related to a primer being higher or lower than the temptype ("minimum" or "maximum")
// values specified (mt and specifiedmt).
func meltingTempErrorString(primerOrientation, sequenceName, tempType string, mt, specifiedmt float64) string {
	var sign, indication string
	if mt > specifiedmt {
		sign, indication = "greater", "increasing"
	} else {
		sign, indication = "lower", "lowering"
	}
	return primerErrorString(primerOrientation, sequenceName, fmt.Sprintf(" with melting temperature (%f) %s than the %s melting temperature specified (%f).", mt, sign, tempType, specifiedmt), fmt.Sprintf(" %s this parameter.", indication))
}

// meltingTempErrorString formats a textual message if a primer has more than one binding site.
func bindingSiteErrorString(primerOrientation, sequenceName string, bindingSites int) string {
	return primerErrorString(primerOrientation, sequenceName, fmt.Sprintf(" with more than one (%d) binding sites.", bindingSites), " selecting another region.")
}

func DesignFWDPRimerstoCoverFullSequence(seq wtype.DNASequence, sequenceinterval int, maxGCcontent float64, minlength int, maxlength int, minmeltingtemp wunit.Temperature, maxmeltingtemp wunit.Temperature, seqstoavoid []string, overlapthresholdwithseqstoavoid int) (primers []Primer) {
	primers = make([]Primer, 0)
	avoidthese := make([]string, 0)
	avoidthese = append(avoidthese, seqstoavoid...)

	for i := 1; i < len(seq.Sequence()); i = i + sequenceinterval {
		region := DNAregion(seq, i, len(seq.Sequence()))

		primer, err := FWDOligoSeq(region, maxGCcontent, minlength, maxlength, minmeltingtemp, maxmeltingtemp, avoidthese, overlapthresholdwithseqstoavoid)

		if err != nil {
			panic(err.Error() + " for " + region.Nm)
		}

		primer.Nm = "primer_" + seq.Nm + "_" + strconv.Itoa(i) + ":" + strconv.Itoa(i-1+sequenceinterval)

		primers = append(primers, primer)

		avoidthese = append(avoidthese, primer.Sequence())
	}
	return
}

func DesignFWDPRimerstoCoverRegion(seq wtype.DNASequence, regionstart, regionend, sequenceinterval int, maxGCcontent float64, minlength int, maxlength int, minmeltingtemp wunit.Temperature, maxmeltingtemp wunit.Temperature, seqstoavoid []string, overlapthresholdwithseqstoavoid int) (primers []Primer) {
	primers = make([]Primer, 0)
	avoidthese := make([]string, 0)
	avoidthese = append(avoidthese, seqstoavoid...)

	if regionstart-100 > 0 {
		regionstart = regionstart - 100
	} else {
		regionstart = 0
	}

	for i := regionstart; i < regionend; i = i + sequenceinterval {

		region := DNAregion(seq, i, len(seq.Sequence()))

		primer, err := FWDOligoSeq(region, maxGCcontent, minlength, maxlength, minmeltingtemp, maxmeltingtemp, avoidthese, overlapthresholdwithseqstoavoid)

		if err != nil {
			panic(err.Error() + " for " + region.Nm)
		}

		primer.Nm = "primer_" + seq.Nm + "_" + strconv.Itoa(i) + ":" + strconv.Itoa(i-1+sequenceinterval)

		primers = append(primers, primer)

		avoidthese = append(avoidthese, primer.Sequence())
	}
	return
}

/*
func DesignPrimerstoFlankRegion(seq wtype.DNASequence, regionstart, regionend int, maxGCcontent float64, minlength int, maxlength int, minmeltingtemp wunit.Temperature, maxmeltingtemp wunit.Temperature, seqstoavoid []string, overlapthresholdwithseqstoavoid int) (primers [2]Primer) {

	primers = make([]Primer, 2)

	avoidthese := make([]string, 0)

	if len(seqstoavoid) != 0 {
		for _, seq := range seqstoavoid {
			avoidthese = append(avoidthese, seq)
		}
	}

	if regionstart-100 > 0 {
		regionstart = regionstart - 100
	} else {
		regionstart = 0
	}

	for i := regionstart; i < regionend; i = i + sequenceinterval {

		region := DNAregion(seq, i, len(seq.Sequence()))

		primer, err := FWDOligoSeq(region, maxGCcontent, minlength, maxlength, minmeltingtemp, maxmeltingtemp, avoidthese, overlapthresholdwithseqstoavoid)

		if err != nil {
			panic(err.Error() + " for " + region.Nm)
		}

		primer.Nm = "primer_" + seq.Nm + "_" + strconv.Itoa(i) + ":" + strconv.Itoa(i-1+sequenceinterval)

		primers = append(primers, primer)

		avoidthese = append(avoidthese, primer.Sequence())
	}
	return
}
*/

func DesignFWDPRimerstoCoverSequence(seq wtype.DNASequence, targetseq string, sequenceinterval int, maxGCcontent float64, minlength int, maxlength int, minmeltingtemp wunit.Temperature, maxmeltingtemp wunit.Temperature, seqstoavoid []string, overlapthresholdwithseqstoavoid int) (primers []Primer) {
	primers = make([]Primer, 0)
	avoidthese := make([]string, 0)
	avoidthese = append(avoidthese, seqstoavoid...)

	seqsfound := sequences.FindSeqsinSeqs(seq.Sequence(), []string{targetseq})

	if len(seqsfound) != 1 {
		panicstatement := fmt.Sprintln("found ", len(seqsfound), " instances of ", targetseq, " in ", seq)
		panic(panicstatement)
	}
	/*
		if len(seqsfound[0].Positions) != 2 {
			panicstatement := fmt.Sprintln("positions found == ", len(seqsfound[0].Positions))
			panic(panicstatement)
		}*/

	regionstart := seqsfound[0].Positions[0]
	regionend := regionstart + len(targetseq)

	if regionstart-100 > 0 {
		regionstart = regionstart - 100
	} else {
		regionstart = 0
	}

	/*if regionstart > regionend {
		temp := regionstart
		regionstart = regionend
		regionend = temp
	}*/

	for i := regionstart; i < regionend; i = i + sequenceinterval {

		region := DNAregion(seq, i, len(seq.Sequence()))

		primer, err := FWDOligoSeq(region, maxGCcontent, minlength, maxlength, minmeltingtemp, maxmeltingtemp, avoidthese, overlapthresholdwithseqstoavoid)

		if err != nil {
			panic(err.Error() + " for " + region.Nm)
		}

		primer.Nm = "primer_" + seq.Nm + "_" + strconv.Itoa(i) + ":" + strconv.Itoa(i-1+sequenceinterval)

		primers = append(primers, primer)

		avoidthese = append(avoidthese, primer.Sequence())
	}
	return
}

func DesignFWDPRimerstoCoverFeature(seq wtype.DNASequence, targetfeaturename string, sequenceinterval int, maxGCcontent float64, minlength int, maxlength int, minmeltingtemp wunit.Temperature, maxmeltingtemp wunit.Temperature, seqstoavoid []string, overlapthresholdwithseqstoavoid int) (primers []Primer) {
	primers = make([]Primer, 0)
	avoidthese := make([]string, 0)
	avoidthese = append(avoidthese, seqstoavoid...)

	features := seq.GetFeatureByName(targetfeaturename)
	if len(features) == 0 {
		panicstatement := fmt.Sprintf("feature: %s not found amongst features: %+v", targetfeaturename, seq.Features)
		panic(panicstatement)
	} else if len(features) > 1 {
		panicstatement := fmt.Sprintf("feature: %s found %d times amongst features: %+v", targetfeaturename, len(features), seq.Features)
		panic(panicstatement)
	}

	targetseq := features[0].DNASeq

	seqsfound := sequences.FindSeqsinSeqs(seq.Sequence(), []string{targetseq})

	if len(seqsfound) != 1 {
		panicstatement := fmt.Sprintln("found ", len(seqsfound), " instances of ", targetseq, " in ", seq)
		panic(panicstatement)
	}

	regionstart := seqsfound[0].Positions[0]
	regionend := regionstart + len(targetseq)

	if regionstart-100 > 0 {
		regionstart = regionstart - 100
	} else {
		regionstart = 0
	}

	for i := regionstart; i < regionend; i = i + sequenceinterval {

		region := DNAregion(seq, i, len(seq.Sequence()))

		primer, err := FWDOligoSeq(region, maxGCcontent, minlength, maxlength, minmeltingtemp, maxmeltingtemp, avoidthese, overlapthresholdwithseqstoavoid)

		if err != nil {
			panic(err.Error() + " for " + region.Nm)
		}

		primer.Nm = "primer_" + seq.Nm + "_" + strconv.Itoa(i) + ":" + strconv.Itoa(i-1+sequenceinterval)

		primers = append(primers, primer)

		avoidthese = append(avoidthese, primer.Sequence())
	}
	return
}

// CheckNonSpecificBinding checks for number of EXACT matches of a primerSeq in a fullSeq.
func CheckNonSpecificBinding(fullSeq, primerSeq wtype.DNASequence) (count int) {
	seqsfound := sequences.FindAll(&fullSeq, &primerSeq)
	count = len(seqsfound.Positions)
	return
}

func MakeOutwardFacingPrimers(sequence wtype.DNASequence, maxGCcontent float64, minlength int, maxlength int, minmeltingtemp wunit.Temperature, maxmeltingtemp wunit.Temperature, seqstoavoid []string, overlapthresholdwithseqstoavoid int) (oligoforpartsafter Primer, oligoforpartsbefore Primer) {

	endstartingpoint := wtype.MakeLinearDNASequence("endprimer", sequence.Sequence()[len(sequence.Sequence())-100:len(sequence.Sequence())-1])

	oligoforpartsafter, _ = FWDOligoSeq(endstartingpoint, maxGCcontent, minlength, maxlength, minmeltingtemp, maxmeltingtemp, seqstoavoid, overlapthresholdwithseqstoavoid) // nolint

	// now reverse
	reversesequence := wtype.RevComp(sequence.Sequence())

	endstartingpoint = wtype.MakeLinearDNASequence("endprimer", reversesequence[len(reversesequence)-100:len(reversesequence)-1])

	oligoforpartsbefore, _ = FWDOligoSeq(endstartingpoint, maxGCcontent, minlength, maxlength, minmeltingtemp, maxmeltingtemp, seqstoavoid, overlapthresholdwithseqstoavoid) // nolint

	oligoforpartsbefore.Reverse = true

	return
}
