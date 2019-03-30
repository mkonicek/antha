// align
package align

import (
	"fmt"
	"strings"

	"testing"

	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/sequences"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
)

type alignmentTest struct {
	Name              string
	Seq1, Seq2        wtype.DNASequence
	Reverse           bool // if the query should align to the reverse strand of the template
	Alignment         string
	Identity          float64
	AlignmentStartPos int
	AlignmentEndPos   int
	ScoringMatrix     ScoringMatrix
	Score             int
}

var (
	tests = []alignmentTest{
		{
			Name: "Test1",
			Seq1: wtype.DNASequence{
				Nm:  "Seq1",
				Seq: "GTTGACAGACTAGATTCACG",
			},
			Seq2: wtype.DNASequence{
				Nm:  "Seq2",
				Seq: "GACAGACGA",
			},
			Reverse:           false,
			Alignment:         fmt.Sprintf("%s\n%s\n", "GACAGACTAGA", "GACAGAC--GA"),
			Identity:          0.8181818181818182,
			ScoringMatrix:     Fitted,
			AlignmentStartPos: 4,
			AlignmentEndPos:   10,
			Score:             69, // 7 + 10 + 9 + 10 + 7 + 10 + 9 - 5 - 5 + 7 + 10
		},
		{
			Name: "Test2",
			Seq1: wtype.DNASequence{
				Nm:  "Seq3",
				Seq: "GTTGACAGACTAGATTCACG",
			},
			Seq2: wtype.DNASequence{
				Nm: "Seq4",

				Seq: "GTTGACA",
			},
			Reverse:           false,
			Identity:          1,
			Alignment:         fmt.Sprintf("%s\n%s\n", "GTTGACA", "GTTGACA"),
			ScoringMatrix:     Fitted,
			AlignmentStartPos: 1,
			AlignmentEndPos:   7,
			Score:             59, // 7 + 8 + 8 + 7 + 10 + 9 + 10
		},
		{
			Name: "Test3",
			Seq1: wtype.DNASequence{
				Nm:  "Seq5",
				Seq: "GTTGACAGACTTTATTCACG",
			},
			Seq2: wtype.DNASequence{
				Nm:  "Seq6",
				Seq: "GTTGAGACAAGATTCACG",
			},
			Reverse:           false,
			Alignment:         fmt.Sprintf("%s\n%s\n", "GTTGACAGACtttATTCACG", "GTTGA--GACaagATTCACG"),
			Identity:          0.75,
			ScoringMatrix:     FittedAffine,
			AlignmentStartPos: 14,
			AlignmentEndPos:   20,
			Score:             5, // 5 * 1 - 5 (gap open) + 2 * -1 (gap extend) + 3 * 1 + 3 * -1 + 7 * 1
		},
		{
			Name: "TerminatorAlignmentCorrect",
			Seq1: wtype.DNASequence{
				Nm:  "SequencingResult",
				Seq: "CGCGacgtAaTACGACTcaCTATAGGgCGAATTGGCGGAAGGCCGTCAAGGCCGCATGAAGGGCGCGCCAGGTCTCAGGTACCAGGCATCAAATAAAACGAAAGGCTCAGTCGAAAGACTGGGCCTTTCGTTTTATCTGTTGTTTGTCGGTGAACGCTCTCTACTAGAGTCACACTGGCTCACCTTCGGGTGGGCCTTTCTGCGTTTATATGTAGAAGAGCTAATAGAGACCATTAATTAACAACCTGGGCCTCATGGGCCTTCCGCTCACTGCCCGCTTTCCAGTCGGGAAACCTGTCGTGCCAGCTGCATTAACATGGTCATAGCTGTTTCCTTGCGTATTGGGCGCTCTCCGCTTCCTCGCTCACTGACTCGCTGCGCTCGGTCGTTCGGGTAAAGCCTGGGGTGCCTAATGAGCAAAAGGCCAGCAAAAGGCCAGGAACCGTAAAAAGGCCGCGTTGCTGGCGTTTTTCCATAGGCTCCGCCCCCCTGACGAGCATCACAAAAATCGACGCTCAAGTCAGAGGTGGCGAAACCCGACAGGACTATAAAGATACCAGGCGTTTCCCCCTGGAAGCTCCCTCGTGCGCTCTCCTGTTCCGACCCTGCCGCTTACCGGATACCTGTCCGCCTTTCTCCCTTCGGGAAGCGTGGCGCTTTCTCATAGCTCACGCTGTAGGTATCTCAGTTCGGTGTAGGTCGTTCGCTCCAAGCTGGGCTGTGTGCACGAACCCCCcGTTCAGCCCGACCGCTGCGCCTTATCCGGTAACTATCGTCTTGAGTCCAACCCGGTAAGACACGACTTATCGCCACTGGCAGCAGCCACTGGTAACAGGATTAGCAGAGCGAGGTATGTaggCGGTGCTACAGAGTTCTTGAAGTGGTGGCCTAACTACGGCTACACTagAAGAACAGTATTTGGtaTCTGCGCTCTGCtgAaGCCAGTTaccttcggAAAAanagTTggAaGCTCTTGATCCGGcnaacAAaCcac",
			},
			Seq2: wtype.DNASequence{
				Nm:  "TerSequence",
				Seq: "GGTCTCAGGTACCAGGCATCAAATAAAACGAAAGGCTCAGTCGAAAGACTGGGCCTTTCGTTTTATCTGTTGTTTGTCGGTGAACGCTCTCTACTAGAGTCACACTGGCTCACCTTCGGGTGGGCCTTTCTGCGTTTATATGTAGAAGAGCTAATAGAGACC",
			},
			Reverse:           false,
			Identity:          1,
			Alignment:         fmt.Sprintf("%s\n%s\n", "GGTCTCAGGTACCAGGCATCAAATAAAACGAAAGGCTCAGTCGAAAGACTGGGCCTTTCGTTTTATCTGTTGTTTGTCGGTGAACGCTCTCTACTAGAGTCACACTGGCTCACCTTCGGGTGGGCCTTTCTGCGTTTATATGTAGAAGAGCTAATAGAGACC", "GGTCTCAGGTACCAGGCATCAAATAAAACGAAAGGCTCAGTCGAAAGACTGGGCCTTTCGTTTTATCTGTTGTTTGTCGGTGAACGCTCTCTACTAGAGTCACACTGGCTCACCTTCGGGTGGGCCTTTCTGCGTTTATATGTAGAAGAGCTAATAGAGACC"),
			ScoringMatrix:     Fitted,
			AlignmentStartPos: 0,
			AlignmentEndPos:   0,
			Score:             1369, // 39 * 10 (A) + 36 * 9 (C) + 41 * 7 (G) + 46 * 8 (T)
		},
		{
			Name: "MismatchingAlignmentReverse",
			Seq1: wtype.DNASequence{
				Nm:  "TemplateSequence",
				Seq: "ctcatgaccaaaatcccttaacgtgagttacgcgcgcgtcgttccactgagcgtcagaccccgtagaaaagatcaaaggatcttcttgagatcctttttttctgcgcgtaatctgctgcttgcaaacaaaaaaaccaccgctaccagcggtggtttgtttgccggatcaagagctaccaactctttttccgaaggtaactggcttcagcagagcgcagataccaaatactgttcttctagtgtagccgtagttagcccaccacttcaagaactctgtagcaccgcctacatacctcgctctgctaatcctgttaccagtggctgctgccagtggcgataagtcgtgtcttaccgggttggactcaagacgatagttaccggataaggcgcagcggtcgggctgaacggggggttcgtgcacacagcccagcttggagcgaacgacctacaccgaactgagatacctacagcgtgagctatgagaaagcgccacgcttcccgaagggagaaaggcggacaggtatccggtaagcggcagggtcggaacaggagagcgcacgagggagcttccagggggaaacgcctggtatctttatagtcctgtcgggtttcgccacctctgacttgagcgtcgatttttgtgatgctcgtcaggggggcggagcctatggaaaaacgccagcaacgcggcctttttacggttcctggccttttgctggccttttgctcacatgttctttcctgcgttatcccctgattctgtggataaccgtattaccgcctttgagtgagctgataccgctcgccgcagccgaacgaccgagcgcagcgagtcagtgagcgaggaagcggaaggcgagagtagggaactgccaggcatcaaactaagcagaaggcccctgacggatggcctttttgcgtttctacaaactctttctgtgttgtaaaacgacggccagtcttaagctcgggccccctgggcggttctgataacgagtaatcgttaatccgcaaataacgtaaaaacccgcttcggcgggtttttttatggggggagtttagggaaagagcatttgtcagaatatttaagggcgcctgtcactttgcttgatatatgagaattatttaaccttataaatgagaaaaaagcaacgcactttaaataagatacgttgctttttcgattgatgaacacctataattaaactattcatctattatttatgattttttgtatatacaatatttctagtttgttaaagagaattaagaaaataaatctcgaaaataataaagggaaaatcagtttttgatatcaaaattatacatgtcaacgataatacaaaatataatacaaactataagatgttatcagtatttattatgcatttagaataaattttgtGTCGGCTCTTCAACCTGTGTCGCCCTTAATTGTGAGCGGATAACAATTACGAGCTTCATGCACAGTGAAATCATGAAAAATTTATTTGCTTTGTGAGCGGATAACAATTATAATATGTGGAATTGTGAGCGCTCACAATTCCACAACGGTTTCCCTCTAGAAATAATTTTGTTTAACTACTAATCTCATATATCAAATATAGGGTGGATCATATGACGGCATTGACGGAAGGTGCAAAACTGTTTGAGAAAGAGATCCCGTATATCACCGAACTGGAAGGCGACGTCGAAGGTATGAAATTTATCATTAAAGGCGAGGGTACCGGTGACGCGACCACGGGTACCATTAAAGCGAAATACATCTGCACTACGGGCGACCTGCCGGTCCCGTGGGCAACCCTGGTGAGCACCCTGAGCTACGGTGTTCAGTGTTTCGCCAAGTACCCGAGCCACATCAAGGATTTCTTTAAGAGCGCCATGCCGGAAGGTTATACCCAAGAGCGTACCATCAGCTTCGAAGGCGACGGCGTGTACAAGACGCGTGCTATGGTTACCTACGAACGCGGTTCTATCTACAATCGTGTCACGCTGACTGGTGAGAACTTTAAGAAAGACGGTCACATTCTGCGTAAGAACGTTGCATTCCAATGCCCGCCAAGCATTCTGTATATTCTGCCTGACACCGTTAACAATGGCATCCGCGTTGAGTTCAACCAGGCGTACGATATTGAAGGTGTGACCGAAAAACTGGTTACCAAATGCAGCCAAATGAATCGTCCGTTGGCGGGCTCCGCGGCAGTGCATATCCCGCGTTATCATCACATTACCTACCACACCAAACTGAGCAAAGACCGCGACGAGCGCCGTGATCACATGTGTCTGGTAGAGGTCGTGAAAGCGGTTGATCTGGACACGTATCAGTAATGAGAATTCTGTACACTCGAGGGTACCAGGCATCAAATAAAACGAAAGGCTCAGTCGAAAGACTGGGCCTTTCGTTTTATCTGTTGTTTGTCGGTGAACGCTCTCTACTAGAGTCACACTGGCTCACCTTCGGGTGGGCCTTTCTGCGTTTATATGTAGAAGAGCtaatgaccccaagggcgacaccccctaattagcccgggcgaaaggcccagtctttcgactgagcctttcgttttatttgatgcctggcagttccctactctcgcatggggagtccccacactaccatcggcgctacggcgtttcacttctgagttcggcatggggtcaggtgggaccaccgcgctactgccgccaggcaaacaaggggtgttatgagccatattcaggtataaatgggctcgcgataatgttcagaattggttaattggttgtaacactgacccctatttgtttatttttctaaatacattcaaatatgtatccgctcatgagacaataaccctgataaatgcttcaataatattgaaaaaggaagaatatgagccatattcaacgggaaacgtcgaggccgcgattaaattccaacatggatgctgatttatatgggtataaatgggctcgcgataatgtcgggcaatcaggtgcgacaatctatcgcttgtatgggaagcccgatgcgccagagttgtttctgaaacatggcaaaggtagcgttgccaatgatgttacagatgagatggtcagactaaactggctgacggaatttatgccacttccgaccatcaagcattttatccgtactcctgatgatgcatggttactcaccactgcgatccccggaaaaacagcgttccaggtattagaagaatatcctgattcaggtgaaaatattgttgatgcgctggcagtgttcctgcgccggttgcactcgattcctgtttgtaattgtccttttaacagcgatcgcgtatttcgcctcgctcaggcgcaatcacgaatgaataacggtttggttgatgcgagtgattttgatgacgagcgtaatggctggcctgttgaacaagtctggaaagaaatgcataaacttttgccattctcaccggattcagtcgtcactcatggtgatttctcacttgataaccttatttttgacgaggggaaattaataggttgtattgatgttggacgagtcggaatcgcagaccgataccaggatcttgccatcctatggaactgcctcggtgagttttctccttcattacagaaacggctttttcaaaaatatggtattgataatcctgatatgaataaattgcagtttcatttgatgctcgatgagtttttctaagcggcgcgccatcgaatggcgcaaaacctttcgcggtatggcatgatagcgcccggaagagagtcaattcagggtggtgaatatgaaaccagtaacgttatacgatgtcgcagagtatgccggtgtctcttatcagaccgtttcccgcgtggtgaaccaggccagccacgtttctgcgaaaacgcgggaaaaagtggaagcggcgatggcggagctgaattacattcccaaccgcgtggcacaacaactggcgggcaaacagtcgttgctgattggcgttgccacctccagtctggccctgcacgcgccgtcgcaaattgtcgcggcgattaaatctcgcgccgatcaactgggtgccagcgtggtggtgtcgatggtagaacgaagcggcgtcgaagcctgtaaagcggcggtgcacaatcttctcgcgcaacgcgtcagtgggctgatcattaactatccgctggatgaccaggatgccattgctgtggaagctgcctgcactaatgttccggcgttatttcttgatgtctctgaccagacacccatcaacagtattattttctcccatgaggacggtacgcgactgggcgtggagcatctggtcgcattgggtcaccagcaaatcgcgctgttagcgggcccattaagttctgtctcggcgcgtctgcgtctggctggctggcataaatatctcactcgcaatcaaattcagccgatagcggaacgggaaggcgactggagtgccatgtccggttttcaacaaaccatgcaaatgctgaatgagggcatcgttcccactgcgatgctggttgccaacgatcagatggcgctgggcgcaatgcgcgccattaccgagtccgggctgcgcgttggtgcggatatctcggtagtgggatacgacgataccgaagatagctcatgttatatcccgccgttaaccaccatcaaacaggattttcgcctgctggggcaaaccagcgtggaccgcttgctgcaactctctcagggccaggcggtgaagggcaatcagctgttgccagtctcactggtgaaaagaaaaaccaccctggcgcccaatacgcaaaccgcctctccccgcgcgttggccgattcattaatgcagctggcacgacaggtttcccgactggaaagcgggcagtga",
			},
			Seq2: wtype.DNASequence{
				Nm:  "ReverseSequencingResult",
				Seq: wtype.RevComp("CCGGGGATACGCGGTGGTCCACCTGACCCCATGCCGAACTCAGAAGTGAAACGCCGTAGCGCCGATGGTAGTGTGGGGACTCCCCATGCGAGAGTAGGGAACTGCCAGGCATCAAATAAAACGAAAGGCTCAGTCGAAAGACTGGGCCTTTCGCCCGGGCTAATTAGGGGGTGTCGCCCTTGGGGTCATTAGCTCTTCTACATATAAACGCAGAAAGGCCCACCCGAAGGTGAGCCAGTGTGACTCTAGTAGAGAGCGTTCACCGACAAACAACAGATAAAACGAAAGGCCCAGTCTTTCGACTGAGCCTTTCGTTTTATTTGATGCCTGGGGGGCTCGAGTGTACAGAATTCTCATTACTGATACGTGTCCAGATCAACCGCTTTCACGACCTCTACCAGACACATGTGATCACGGCGCTCGTCGCGGTCTTTGCTCAGTTTGGTGTGGTAGGTAATGTGATGATAACGCGGGATATGCACTGCCGCGGAGCCCGCCAACGGACGATTCATTTGGCTGCATTTGGTAACCAGTTTTTCGGTCACACCTTCAATATCGTACGCCTGGTTGAACTCAACGCGGATGCCATTGTTAACGGTGTCAGGCAGAATATACAGAATGCTTGGCGGGCATTGGAATGCAACGTTCTTACGCAGAATGTGACCGTCTTTCTTAAAGTTCTCACCAGTCAGCGTGACACGATTGTAGATAGAACCGCGTTCGTAGGTAACCATAGCACGCGTCTTGTACACGCCGTCGCCTTCGAAGCTGATGGTACGCTCTTGGGTATAACCTTCCGGCATGGCGCTCTTAAAGAAATCCTTGATGTGGCTCGGGTACTTGGCGAAACACTGAACACCGTAGCTCAGGGTGCTCACCAGGGTTGCCCACGGGACCGGCAGGTCGCCCGTAGTGCAGATGTATTTTCGCTTTAATGGTACCCGTGGTCGCGTCACGGGTACCCTCGCCTTTTATGATAAATTTCATACCTTCGACGTCGCCTTTCCAGGTTCGGTGATATACGGGATCTCTTTCTCAACCAGTTTTGCACCTTCCGTCATGCGGTCATATGATCCACCCTTATATTTGAATATATGGAGATTAAAAGGTTAAACAAAATTAATTTCTAGAAGGGAAACCGTTTGTGGGAATTGTGGAACGCCTCACATCCACATATTATAATGGTATCCGCTCCAAAGGCAATTAATTTCATGGAATACTGTCATGAAGACTCGAATGTTATCGCTTCACAATTAAGGGGGCGACACATCGGTGTGTAGA"),
			},
			Reverse:           false,
			Identity:          0.9561875480399693,
			Alignment:         fmt.Sprintf("%s\n%s\n", "TCTtCA-ACC--TGTGTCG--CCCTTAATTGTG-AGCGGATAACAATTACGAG-CTTCATGCACAGTGAAaT-CATGAAAAATTTAtTTG-CTTTGTGAGCGGATAACaATTATAATATGTGGAATTGTGA-GCGcT-CACAATT-CCAC-AACGGTTTCCC-TCTAGAAA-TAATTTTGTTTAA-CTacTAATCT-CATATA-TCAAATAT-AGGGTGGATCATATGACgGCATTGACGGAAGGTGCAAAACTGtTTGAGAAAGAGATCCCGTATATCACCGAA-CTGG-AAGGCGACGTCGAAGGTATGAAATTTATCATtAAAGGCGAGGGTACCgGTGACGCGACCACGGGTACCATTAAAGCG-AAATACATCTGCACTACGGGCGACCTGCCGGTCCCGTGGGCAACCCTGGTGAGCACCCTGAGCTACGGTGTTCAGTGTTTCGCCAAGTACCCGAGCCACATCAAGGATTTCTTTAAGAGCGCCATGCCGGAAGGTTATACCCAAGAGCGTACCATCAGCTTCGAAGGCGACGGCGTGTACAAGACGCGTGCTATGGTTACCTACGAACGCGGTTCTATCTACAATCGTGTCACGCTGACTGGTGAGAACTTTAAGAAAGACGGTCACATTCTGCGTAAGAACGTTGCATTCCAATGCCCGCCAAGCATTCTGTATATTCTGCCTGACACCGTTAACAATGGCATCCGCGTTGAGTTCAACCAGGCGTACGATATTGAAGGTGTGACCGAAAAACTGGTTACCAAATGCAGCCAAATGAATCGTCCGTTGGCGGGCTCCGCGGCAGTGCATATCCCGCGTTATCATCACATTACCTACCACACCAAACTGAGCAAAGACCGCGACGAGCGCCGTGATCACATGTGTCTGGTAGAGGTCGTGAAAGCGGTTGATCTGGACACGTATCAGTAATGAGAATTCTGTACACTCGAGggtaCCAGGCATCAAATAAAACGAAAGGCTCAGTCGAAAGACTGGGCCTTTCGTTTTATCTGTTGTTTGTCGGTGAACGCTCTCTACTAGAGTCACACTGGCTCACCTTCGGGTGGGCCTTTCTGCGTTTATATGTAGAAGAGCTAATGACCCCAAGGGCGACACCCCCTAATTAGCCCGGGCGAAAGGCCCAGTCTTTCGACTGAGCCTTTCGTTTTATTTGATGCCTGGCAGTTCCCTACTCTCGCATGGGGAGTCCCCACACTACCATCGGCGCTACGGCGTTTCACTTCTGAGTTCGGCATGGGGTCAGGTGGGACCACCGCGCTACTGCCGCCAGG", "TCTaCACACCGATGTGTCGCCCCCTTAATTGTGAAGC-GATAAC-ATT-CGAGTCTTCATG-ACAGT--AtTCCATG--AAA-TTAaTTGCCTTTG-GAGCGGAT-ACcATTATAATATGTGG-A-TGTGAGGCGtTCCACAATTCCCACAAACGGTTTCCCTTCTAGAAATTAATTTTGTTTAACCTttTAATCTCCATATATTCAAATATAAGGGTGGATCATATGACcGCA-TGACGGAAGGTGCAAAACTGgTTGAGAAAGAGATCCCGTATATCACCGAACCTGGAAAGGCGACGTCGAAGGTATGAAATTTATCATaAAAGGCGAGGGTACCcGTGACGCGACCACGGGTACCATTAAAGCGAAAATACATCTGCACTACGGGCGACCTGCCGGTCCCGTGGGCAACCCTGGTGAGCACCCTGAGCTACGGTGTTCAGTGTTTCGCCAAGTACCCGAGCCACATCAAGGATTTCTTTAAGAGCGCCATGCCGGAAGGTTATACCCAAGAGCGTACCATCAGCTTCGAAGGCGACGGCGTGTACAAGACGCGTGCTATGGTTACCTACGAACGCGGTTCTATCTACAATCGTGTCACGCTGACTGGTGAGAACTTTAAGAAAGACGGTCACATTCTGCGTAAGAACGTTGCATTCCAATGCCCGCCAAGCATTCTGTATATTCTGCCTGACACCGTTAACAATGGCATCCGCGTTGAGTTCAACCAGGCGTACGATATTGAAGGTGTGACCGAAAAACTGGTTACCAAATGCAGCCAAATGAATCGTCCGTTGGCGGGCTCCGCGGCAGTGCATATCCCGCGTTATCATCACATTACCTACCACACCAAACTGAGCAAAGACCGCGACGAGCGCCGTGATCACATGTGTCTGGTAGAGGTCGTGAAAGCGGTTGATCTGGACACGTATCAGTAATGAGAATTCTGTACACTCGAGccccCCAGGCATCAAATAAAACGAAAGGCTCAGTCGAAAGACTGGGCCTTTCGTTTTATCTGTTGTTTGTCGGTGAACGCTCTCTACTAGAGTCACACTGGCTCACCTTCGGGTGGGCCTTTCTGCGTTTATATGTAGAAGAGCTAATGACCCCAAGGGCGACACCCCCTAATTAGCCCGGGCGAAAGGCCCAGTCTTTCGACTGAGCCTTTCGTTTTATTTGATGCCTGGCAGTTCCCTACTCTCGCATGGGGAGTCCCCACACTACCATCGGCGCTACGGCGTTTCACTTCTGAGTTCGGCATGGGGTCAGGT-GGACCACCGCG-TA-T-CC-CC-GG"),
			ScoringMatrix:     Fitted,
			AlignmentStartPos: 0,
			AlignmentEndPos:   0,
			Score:             10333,
		},
		{
			Name: "plasmidAlignmentTest",
			Seq1: wtype.DNASequence{
				Nm:      "Seq3Plasmid",
				Seq:     "GTTGACAGACTAGATTCACG",
				Plasmid: true,
			},
			Seq2: wtype.DNASequence{
				Nm:  "Seq4",
				Seq: "CACGGTTGACA",
			},
			Reverse:           false,
			Identity:          1,
			Alignment:         fmt.Sprintf("%s\n%s\n", "CACGGTTGACA", "CACGGTTGACA"),
			ScoringMatrix:     Fitted,
			AlignmentStartPos: 17,
			AlignmentEndPos:   7,
			Score:             94, // 3 * 10 (A) + 3 * 9 (C) + 3 * 7 (G) + 2 * 8 (T)
		},
		{
			Name: "plasmidAlignmentTest2",
			Seq1: wtype.DNASequence{
				Nm:      "Seq3Plasmid",
				Seq:     "GTTGACAGACTAGATTCACG",
				Plasmid: true,
			},
			Seq2: wtype.DNASequence{
				Nm:  "Seq4",
				Seq: "GTTGACAGA",
			},
			Reverse:           false,
			Identity:          1,
			Alignment:         fmt.Sprintf("%s\n%s\n", "GTTGACAGA", "GTTGACAGA"),
			ScoringMatrix:     Fitted,
			AlignmentStartPos: 1,
			AlignmentEndPos:   9,
			Score:             76, // 3 * 10 (A) + 1 * 9 (C) + 3 * 7 (G) + 2 * 8 (T)
		},
		{
			Name: "revTest",
			Seq1: wtype.DNASequence{
				Nm:      "Seq3Plasmid",
				Seq:     "GTTGACAGACTAGATTCACG",
				Plasmid: true,
			},
			Seq2: wtype.DNASequence{
				Nm:  "Seq4",
				Seq: wtype.RevComp("GTTGACAGA"),
			},
			Reverse:           true,
			Identity:          1,
			Alignment:         fmt.Sprintf("%s\n%s\n", wtype.RevComp("GTTGACAGA"), wtype.RevComp("GTTGACAGA")),
			ScoringMatrix:     Fitted,
			AlignmentStartPos: 9,
			AlignmentEndPos:   1,
			Score:             76, // 3 * 10 (A) + 1 * 9 (C) + 3 * 7 (G) + 2 * 8 (T)
		},
		{
			Name: "revTestGapped",
			Seq1: wtype.DNASequence{
				Nm:      "Other",
				Seq:     "GTTGCCACAGACTAGATTCACG",
				Plasmid: false,
			},
			Seq2: wtype.DNASequence{
				Nm:  "Seq4",
				Seq: wtype.RevComp("GTTGACAGACTAGATTCACG"), // missing CC at position [5, 6]
			},
			Reverse:           true,
			Identity:          0.9090909090909091,
			Alignment:         fmt.Sprintf("%s\n%s\n", wtype.RevComp("GTTGCCACAGACTAGATTCACG"), wtype.RevComp("GTTG--ACAGACTAGATTCACG")),
			ScoringMatrix:     Fitted,
			AlignmentStartPos: 22,
			AlignmentEndPos:   7,   // longest continuous match start, not alignment start
			Score:             161, // 6 * 10 (A) + 4 * 9 (C) + 5 * 7 (G) + 5 * 8 (T) - 2 * 5 (aligned gap)
		},
		{
			Name: "plasmidRevTest",
			Seq1: wtype.DNASequence{
				Nm:      "Seq3Plasmid",
				Seq:     "GTTGACAGACTAGATTCACG",
				Plasmid: true,
			},
			Seq2: wtype.DNASequence{
				Nm:  "Seq4",
				Seq: wtype.RevComp("CACGGTTGACA"),
			},
			Reverse:           true,
			Identity:          1,
			Alignment:         fmt.Sprintf("%s\n%s\n", wtype.RevComp("CACGGTTGACA"), wtype.RevComp("CACGGTTGACA")),
			ScoringMatrix:     Fitted,
			AlignmentStartPos: 7,
			AlignmentEndPos:   17,
			Score:             94, // 3 * 10 (A) + 3 * 9 (C) + 3 * 7 (G) + 2 * 8 (T)
		}, /*
			alignmentTest{
				Name: "bigQueryagainstSmallTemplate",
				Seq1: wtype.DNASequence{
					Nm:      "Seq3Plasmid",
					Seq:     "GTTGACAGA",
					Plasmid: false,
				},
				Seq2: wtype.DNASequence{
					Nm:  "Seq4",
					Seq: "GTTGACAGACTAGATTCACG",
				},
				Identity:      1,
				Alignment:     fmt.Sprintf("%s\n%s\n", "GTTGACAGA", "GTTGACAGA"),
				ScoringMatrix: Fitted,
			},*/
	}
)

// the biogo implementation of alignment requires the N nucleotides to be replaced with -
func replaceN(seq wtype.DNASequence) wtype.DNASequence {

	var newSeq []string

	for _, letter := range seq.Seq {
		if strings.ToUpper(string(letter)) == "N" {
			letter = rune('-')
		}
		newSeq = append(newSeq, string(letter))
	}

	seq.Seq = strings.Join(newSeq, "")

	return seq
}

// Align two dna sequences based on a specified scoring matrix
func TestAlign(t *testing.T) {
	for _, test := range tests {
		alignment, err := DNA(replaceN(test.Seq1), replaceN(test.Seq2), test.ScoringMatrix)

		if err != nil {
			t.Error(
				"For", test.Name, "\n",
				"got error:", err.Error(), "\n",
			)
		}
		if alignment.String() != test.Alignment {
			t.Error(
				"For", test.Name, "\n",
				"expected:", "\n",
				test.Alignment,
				"got:", "\n",
				alignment,
			)
		}
		if alignment.Identity() != test.Identity {
			t.Error(
				"For", test.Name, "\n",
				"expected Identity:", test.Identity, "\n",
				"got:", alignment.Identity(), "\n",
			)
		}
		if alignment.Score() != test.Score {
			t.Error(
				"For", test.Name, "\n",
				"expected Score:", test.Score, "\n",
				"got:", alignment.Score(), "\n",
			)
		}

		longestMatch := alignment.LongestContinuousSequence()

		positions := sequences.FindAll(&test.Seq1, &longestMatch)

		// longestMatch reverse complements if TemplateResult is not in forward frame
		// (see func looksLikeReverseSequence); the following reverses the positions
		// found, where appropriate
		posInTemplate := [2]int{positions.Positions[0].Start(), positions.Positions[0].End()}
		if test.Reverse {
			posInTemplate[1], posInTemplate[0] = posInTemplate[0], posInTemplate[1]
		}

		if test.AlignmentStartPos > 0 && test.AlignmentEndPos > 0 {
			if test.AlignmentStartPos != posInTemplate[0] || test.AlignmentEndPos != posInTemplate[1] {
				t.Error(
					"For", test.Name, "\n",
					"longest match found: ", longestMatch, "\n",
					"expected start, end:", test.AlignmentStartPos, test.AlignmentEndPos, "\n",
					"got:", fmt.Sprint(posInTemplate), "\n",
				)
			}
		}

		fmt.Println(test.Name, "Coverage:", alignment.Coverage(),
			"Positions in template: ", alignment.Positions().Positions,
			"Matches: ", alignment.Matches(),
			"Mismatches: ", alignment.Mismatches(),
			"Gaps: ", alignment.Gaps(),
			"Longest Matching Sequence:", alignment.LongestContinuousSequence(),
		)
	}
}

func TestAlignPositions(t *testing.T) {

	seq1 := wtype.DNASequence{
		Seq: "GGGGGGGGGGGGGGGGGATGGTACAGG",
	}
	seq2 := wtype.DNASequence{
		Seq: "GATTACA",
	}

	alignment, err := DNA(replaceN(seq1), replaceN(seq2), Fitted)

	if err != nil {
		t.Error(err.Error())
	}

	// GATGGTACA
	// GAT--TACA
	// [17 18 19 20 21 22 23 24 25]
	// [1 2 3 20 20 4 5 6 7] wrong
	// [1 2 3 4 4 4 5 6 7] better

	gotPositions := alignment.Alignment.QueryPositions
	wantPositions := []int{1, 2, 3, 4, 4, 4, 5, 6, 7}

	alnLen := len(gotPositions)

	for i := 0; i < alnLen; i++ {
		if gotPositions[i] != wantPositions[i] {
			t.Error(
				"Expected position:", wantPositions[i], "\n",
				"got:", gotPositions[i], "\n",
			)
		}
	}

}
