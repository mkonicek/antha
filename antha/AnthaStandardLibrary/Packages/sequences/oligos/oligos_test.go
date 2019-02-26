// oligos_test.go
package oligos

import (
	"fmt"
	"testing"

	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
)

// simple reverse complement check to test testing methodology initially
type testpair struct {
	//BasicMeltingTemp test
	sequence    wtype.DNASequence
	meltingtemp wunit.Temperature

	// FWDoligoTest
	mintemp             wunit.Temperature
	maxtemp             wunit.Temperature
	maxGCcontent        float64
	minlength           int
	maxlength           int
	seqstoavoid         []string
	overlapthreshold    int
	outputoligoseq      string
	calculatedGCcontent float64

	// tests of OverlapCheck
	primer1        string
	primer2        string
	overlapPercent float64
	overlapNumber  int
	overlapSeq     string

	// errorCheckingTests
	expectedError error
}

var meltingtemptests = []testpair{

	{sequence: wtype.MakeSingleStrandedDNASequence("whatever", "AAAAAAAAAAAAAAAAAAA"),
		meltingtemp: wunit.NewTemperature(29.511, "C")},
}

var oligotests = []testpair{

	{sequence: wtype.MakeSingleStrandedDNASequence("sfGFP", "ATGAGCAAAGGAGAAGAACTTTTCACTGGAGTTGTCCCAATTCTTGTTGAATTAGATGGTGATGTTAATGGGCACAAATTTTCTGTCCGTGGAGAGGGTGAAGGTGATGCTACAAACGGAAAACTCACCCTTAAATTTATTTGCACTACTGGAAAACTACCTGTTCCATGGCCAACACTTGTCACTACTCTGACCTATGGTGTTCAATGCTTTTCCCGTTATCCGGATCACATGAAACGGCATGACTTTTTCAAGAGTGCCATGCCCGAAGGTTATGTACAGGAACGCACTATATCTTTCAAAGATGACGGGACCTACAAGACGCGTGCTGAAGTCAAGTTTGAAGGTGATACCCTTGTTAATCGTATCGAGTTAAAAGGTATTGATTTTAAAGAAGATGGAAACATTCTCGGACACAAACTCGAGTACAACTTTAACTCACACAATGTATACATCACGGCAGACAAACAAAAGAATGGAATCAAAGCTAACTTCAAAATTCGCCACAACGTTGAAGATGGTTCCGTTCAACTAGCAGACCATTATCAACAAAATACTCCAATTGGCGATGGCCCTGTCCTTTTACCAGACAACCATTACCTGTCGACACAATCTGTCCTTTCGAAAGATCCCAACGAAAAGCGTGACCACATGGTCCTTCTTGAGTTTGTAACTGCTGCTGGGATTACACATGGCATGGATGAGCTCTACAAATAA"),
		meltingtemp:         wunit.NewTemperature(52.764, "C"),
		mintemp:             wunit.NewTemperature(50, "C"),
		maxtemp:             wunit.NewTemperature(85, "C"),
		maxGCcontent:        0.6,
		minlength:           25,
		maxlength:           45,
		seqstoavoid:         []string{""},
		overlapthreshold:    45,
		outputoligoseq:      "ATGAGCAAAGGAGAAGAACTTTTCA",
		calculatedGCcontent: 0.36},
}

var primerErrorTests = []testpair{
	{
		sequence:            wtype.MakeSingleStrandedDNASequence("SGAP-004b-001-SacO1_partial", "GTTTAAATCAAAACTGGTGAAACTCACCCAGGGATTGGCTGACACGAAAAACATATTCTCAATAAATCCTTTAGGGAAATAGGCCAGGTTTTCACCGTAACACGCCACATCTTGCGAATATATGTGTAGAAACTGCCGGAAATCGTCGTGGTATTCACTCCAGAGCGATGAAAACGTTTCAGTTTGCTCATGGAAAACGGTGTA"),
		mintemp:             wunit.NewTemperature(50, "C"),
		maxtemp:             wunit.NewTemperature(85, "C"),
		maxGCcontent:        0.1,
		minlength:           21,
		maxlength:           45,
		seqstoavoid:         []string{""},
		overlapthreshold:    45,
		calculatedGCcontent: 0.355556,
		expectedError:       fmt.Errorf(" For %s could only generate FORWARD primers with GC Content (%f) greater than the maximum GC Content specified (%f). Please try lowering this parameter, or selecting a less-GC rich region.", "SGAP-004b-001-SacO1_partial", 0.355556, 0.100000),
	},
	{
		sequence:            wtype.MakeSingleStrandedDNASequence("SGAP-004b-001-SacO1_partial", "GTTTAAATCAAAACTGGTGAAACTCACCCAGGGATTGGCTGACACGAAAAACATATTCTCAATAAATCCTTTAGGGAAATAGGCCAGGTTTTCACCGTAACACGCCACATCTTGCGAATATATGTGTAGAAACTGCCGGAAATCGTCGTGGTATTCACTCCAGAGCGATGAAAACGTTTCAGTTTGCTCATGGAAAACGGTGTA."),
		meltingtemp:         wunit.NewTemperature(28.000000, "C"),
		mintemp:             wunit.NewTemperature(50, "C"),
		maxtemp:             wunit.NewTemperature(85, "C"),
		maxGCcontent:        0.6,
		minlength:           9,
		maxlength:           10,
		seqstoavoid:         []string{""},
		overlapthreshold:    45,
		calculatedGCcontent: 0.36,
		expectedError:       fmt.Errorf(" For %s could only generate FORWARD primers with melting temperature (%f) lower than the minimum melting temperature specified (%f). Please try lowering this parameter.", "SGAP-004b-001-SacO1_partial", 28.000000, 50.000000),
	},
	{
		sequence:            wtype.MakeSingleStrandedDNASequence("MultipleBindingSites", "ACGGGGGCGAAGAAGTTGTCCATATTGGCCACGTTTAAATCAAAAACGGGGGCGAAGAAGTTGTCCATATTGGCCACGTTTAAATCAAAAACGGGGGCGAAGAAGTTGTCCATATTGGCCACGTTTAAATCAAAAACGGGGGCGAAGAAGTTGTCCATATTGGCCACGTTTAAATCAAAATGGGATATATCAACGGTGGTATATCCAGTGATTTTTTTCTCCATATTCTTCCTTTTTCA."),
		mintemp:             wunit.NewTemperature(45, "C"),
		maxtemp:             wunit.NewTemperature(85, "C"),
		maxGCcontent:        0.6,
		minlength:           9,
		maxlength:           21,
		seqstoavoid:         []string{""},
		overlapthreshold:    45,
		calculatedGCcontent: 0.36,
		expectedError:       fmt.Errorf(" For %s could only generate FORWARD primers with more than one (%d) binding sites. Please try selecting another region.", "MultipleBindingSites", 4),
	},
	{
		sequence:            wtype.MakeSingleStrandedDNASequence("SeqsToAvoid", "GTTTAAATCAAAACTGGTGAAACTCACCCAGGGATTGGCTGACACGAAAAACATATTCTCAATAAATCCTTTAGGGAAATAGGCCAGGTTTTCACCGTAACACGCCACATCTTGCGAATATATGTGTAGAAACTGCCGGAAATCGTCGTGGTATTCACTCCAGAGCGATGAAAACGTTTCAGTTTGCTCATGGAAAACGGTGTA"),
		mintemp:             wunit.NewTemperature(45, "C"),
		maxtemp:             wunit.NewTemperature(85, "C"),
		maxGCcontent:        0.6,
		minlength:           9,
		maxlength:           21,
		seqstoavoid:         []string{"GTTTAAATCAAAACTGGTGAAACTCACCCAGGGATTGGCTGACACGAAAAACATATTCTCAATAAATCCTTTAGGGAAATAGGCCAGGTTTTCACCGTAACACGCCACATCTTGCGAATATATGTGTAGAAACTGCCGGAAATCGTCGTGGTATTCACTCCAGAGCGATGAAAACGTTTCAGTTTGCTCATGGAAAACGGTGTA"},
		overlapthreshold:    45,
		calculatedGCcontent: 0.36,
		expectedError:       fmt.Errorf(" For %s could only generate FORWARD primers that contain the specified sequences to avoid. Please try removing these from the parameters.", "SeqsToAvoid"),
	},
}

var overlaptests = []testpair{

	{primer1: "AATACCAGTAGGGTAGAAGAGCACG",
		primer2:        "AATACCAGTAGGGTAGAAGAGCAC",
		overlapPercent: 1.0,
		overlapNumber:  24,
		overlapSeq:     "AATACCAGTAGGGTAGAAGAGCAC"},
	{primer1: "AAAAAAAAAAAAAAAA",
		primer2:        "TTTTTTTTTTTTT",
		overlapPercent: 0.0,
		overlapNumber:  0,
		overlapSeq:     ""},
	{primer1: "AATACCAGTAGGGTAGAAGAGCACG",
		primer2:        "CCAAAAAATAGAAGAGCAC",
		overlapPercent: 0.5789473684210527,
		overlapNumber:  11,
		overlapSeq:     "TAGAAGAGCAC"},
	{primer2: "TAGAAGAGCACGCCCCCCCC",
		primer1:        "CCAAAAAATAGAAGAGCAC",
		overlapPercent: 0.5789473684210527,
		overlapNumber:  11,
		overlapSeq:     "TAGAAGAGCAC"},
	{primer2: "",
		primer1:        "CCAAAAAATAGAAGAGCAC",
		overlapPercent: 0.0,
		overlapNumber:  0,
		overlapSeq:     ""},
}

func TestBasicMeltingTemp(t *testing.T) {
	for _, oligo := range meltingtemptests {
		result := BasicMeltingTemp(oligo.sequence)
		if result.ToString() != oligo.meltingtemp.ToString() {
			t.Error(
				"For", oligo.sequence, "\n",
				"expected", oligo.meltingtemp.ToString(), "\n",
				"got", result.ToString(), "\n",
			)
		}
	}

}

func TestFWDOligoSeq(t *testing.T) {
	for _, oligo := range oligotests {
		oligoseq, err := FWDOligoSeq(oligo.sequence, oligo.maxGCcontent, oligo.minlength, oligo.maxlength, oligo.mintemp, oligo.maxtemp, oligo.seqstoavoid, oligo.overlapthreshold)
		if oligoseq.Sequence() != oligo.outputoligoseq {
			t.Error(
				"For", oligo.sequence, "\n",
				"expected", oligo.outputoligoseq, "\n",
				"got", oligoseq.Sequence(), "\n",
			)
		}
		if oligoseq.GCContent != oligo.calculatedGCcontent {
			t.Error(
				"For", oligo.sequence, "\n",
				"expected", oligo.calculatedGCcontent, "\n",
				"got", oligoseq.GCContent, "\n",
			)
		}
		if err != nil {
			t.Error(
				"errors:", err.Error(), "\n",
			)
		}

	}

	for _, oligo := range primerErrorTests {
		_, err := FWDOligoSeq(oligo.sequence, oligo.maxGCcontent, oligo.minlength, oligo.maxlength, oligo.mintemp, oligo.maxtemp, oligo.seqstoavoid, oligo.overlapthreshold)
		if err.Error() != oligo.expectedError.Error() {
			t.Error(
				"For", oligo.sequence.Name(), "\n",
				"expected", oligo.expectedError.Error(), "\n",
				"got", err.Error(), "\n",
			)
		}
	}
}

func TestOverlapCheck(t *testing.T) {
	for _, test := range overlaptests {
		percent, number, seq := OverlapCheck(test.primer1, test.primer2)

		if percent != test.overlapPercent {
			t.Error(
				"For", test.primer1, " and ", test.primer2, "\n",
				"expected", test.overlapPercent, "\n",
				"got", percent, "\n",
			)
		}

		if number != test.overlapNumber {
			t.Error(
				"For", test.primer1, " and ", test.primer2, "\n",
				"expected", test.overlapNumber, "\n",
				"got", number, "\n",
			)
		}

		if seq != test.overlapSeq {
			t.Error(
				"For", test.primer1, " and ", test.primer2, "\n",
				"expected", test.overlapSeq, "\n",
				"got", seq, "\n",
			)
		}
	}
}

type regionTest struct {
	LargeSeq wtype.DNASequence
	SmallSeq wtype.DNASequence
	Start    int
	End      int
}

var regionTests = []regionTest{

	{
		LargeSeq: wtype.DNASequence{
			Nm:      "Test1",
			Seq:     "ATCGTAGTGTG",
			Plasmid: false,
		},
		SmallSeq: wtype.DNASequence{
			Seq:     "TAG",
			Plasmid: false,
		},
		Start: 5,
		End:   7,
	},
	{
		LargeSeq: wtype.DNASequence{
			Nm:      "Test2",
			Seq:     "ATCGTAGTGTG",
			Plasmid: true,
		},
		SmallSeq: wtype.DNASequence{
			Seq:     "GTGTGATCGT",
			Plasmid: false,
		},
		Start: 7,
		End:   5,
	},
	{
		LargeSeq: wtype.DNASequence{
			Nm:      "Test3",
			Seq:     "ATCGTAGTGTG",
			Plasmid: false,
		},
		SmallSeq: wtype.DNASequence{
			Seq:     "ATCGTAGTGTG",
			Plasmid: false,
		},
		Start: 1,
		End:   11,
	},
	{
		LargeSeq: wtype.DNASequence{
			Nm:      "Test4",
			Seq:     "ATCGTAGTGTG",
			Plasmid: false,
		},
		SmallSeq: wtype.DNASequence{
			Seq:     "ATCGTAGTGT",
			Plasmid: false,
		},
		Start: 1,
		End:   10,
	},
}

func TestDNARegion(t *testing.T) {
	for _, test := range regionTests {
		result := DNAregion(test.LargeSeq, test.Start, test.End)

		if result.Seq != test.SmallSeq.Seq {
			t.Error(
				"For", test.LargeSeq.Nm, "\n",
				"expected", test.SmallSeq.Seq, "\n",
				"got", result.Seq, "\n",
			)
		}

	}
}
