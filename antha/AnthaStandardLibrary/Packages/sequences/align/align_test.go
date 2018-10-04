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
	Alignment         string
	Identity          float64
	AlignmentStartPos int
	AlignmentEndPos   int
	ScoringMatrix     ScoringMatrix
}

var (
	tests []alignmentTest = []alignmentTest{
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
			Alignment:         fmt.Sprintf("%s\n%s\n", "GACAGACTAGA", "GACAGAC--GA"),
			Identity:          0.8181818181818182,
			ScoringMatrix:     Fitted,
			AlignmentStartPos: 4,
			AlignmentEndPos:   10,
		},
		{
			Name: "Test2",
			Seq1: wtype.DNASequence{
				Nm:  "Seq3",
				Seq: "GTTGACAGACTAGATTCACG",
			},
			Seq2: wtype.DNASequence{
				Nm:  "Seq4",
				Seq: "GTTGACA",
			},
			Identity:          1,
			Alignment:         fmt.Sprintf("%s\n%s\n", "GTTGACA", "GTTGACA"),
			ScoringMatrix:     Fitted,
			AlignmentStartPos: 1,
			AlignmentEndPos:   7,
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
			Alignment:         fmt.Sprintf("%s\n%s\n", "GTTGACAGACtttATTCACG", "GTTGA--GACaagATTCACG"),
			Identity:          0.75,
			ScoringMatrix:     FittedAffine,
			AlignmentStartPos: 14,
			AlignmentEndPos:   20,
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
			Identity:          1,
			Alignment:         fmt.Sprintf("%s\n%s\n", "GGTCTCAGGTACCAGGCATCAAATAAAACGAAAGGCTCAGTCGAAAGACTGGGCCTTTCGTTTTATCTGTTGTTTGTCGGTGAACGCTCTCTACTAGAGTCACACTGGCTCACCTTCGGGTGGGCCTTTCTGCGTTTATATGTAGAAGAGCTAATAGAGACC", "GGTCTCAGGTACCAGGCATCAAATAAAACGAAAGGCTCAGTCGAAAGACTGGGCCTTTCGTTTTATCTGTTGTTTGTCGGTGAACGCTCTCTACTAGAGTCACACTGGCTCACCTTCGGGTGGGCCTTTCTGCGTTTATATGTAGAAGAGCTAATAGAGACC"),
			ScoringMatrix:     Fitted,
			AlignmentStartPos: 0,
			AlignmentEndPos:   0,
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
			Identity:          0.9561875480399693,
			Alignment:         fmt.Sprintf("%s\n%s\n", "TCTtCA-ACC--TGTGTCG--CCCTTAATTGTG-AGCGGATAACAATTACGAG-CTTCATGCACAGTGAAaT-CATGAAAAATTTAtTTG-CTTTGTGAGCGGATAACaATTATAATATGTGGAATTGTGA-GCGcT-CACAATT-CCAC-AACGGTTTCCC-TCTAGAAA-TAATTTTGTTTAA-CTacTAATCT-CATATA-TCAAATAT-AGGGTGGATCATATGACgGCATTGACGGAAGGTGCAAAACTGtTTGAGAAAGAGATCCCGTATATCACCGAA-CTGG-AAGGCGACGTCGAAGGTATGAAATTTATCATtAAAGGCGAGGGTACCgGTGACGCGACCACGGGTACCATTAAAGCG-AAATACATCTGCACTACGGGCGACCTGCCGGTCCCGTGGGCAACCCTGGTGAGCACCCTGAGCTACGGTGTTCAGTGTTTCGCCAAGTACCCGAGCCACATCAAGGATTTCTTTAAGAGCGCCATGCCGGAAGGTTATACCCAAGAGCGTACCATCAGCTTCGAAGGCGACGGCGTGTACAAGACGCGTGCTATGGTTACCTACGAACGCGGTTCTATCTACAATCGTGTCACGCTGACTGGTGAGAACTTTAAGAAAGACGGTCACATTCTGCGTAAGAACGTTGCATTCCAATGCCCGCCAAGCATTCTGTATATTCTGCCTGACACCGTTAACAATGGCATCCGCGTTGAGTTCAACCAGGCGTACGATATTGAAGGTGTGACCGAAAAACTGGTTACCAAATGCAGCCAAATGAATCGTCCGTTGGCGGGCTCCGCGGCAGTGCATATCCCGCGTTATCATCACATTACCTACCACACCAAACTGAGCAAAGACCGCGACGAGCGCCGTGATCACATGTGTCTGGTAGAGGTCGTGAAAGCGGTTGATCTGGACACGTATCAGTAATGAGAATTCTGTACACTCGAGggtaCCAGGCATCAAATAAAACGAAAGGCTCAGTCGAAAGACTGGGCCTTTCGTTTTATCTGTTGTTTGTCGGTGAACGCTCTCTACTAGAGTCACACTGGCTCACCTTCGGGTGGGCCTTTCTGCGTTTATATGTAGAAGAGCTAATGACCCCAAGGGCGACACCCCCTAATTAGCCCGGGCGAAAGGCCCAGTCTTTCGACTGAGCCTTTCGTTTTATTTGATGCCTGGCAGTTCCCTACTCTCGCATGGGGAGTCCCCACACTACCATCGGCGCTACGGCGTTTCACTTCTGAGTTCGGCATGGGGTCAGGTGGGACCACCGCGCTACTGCCGCCAGG", "TCTaCACACCGATGTGTCGCCCCCTTAATTGTGAAGC-GATAAC-ATT-CGAGTCTTCATG-ACAGT--AtTCCATG--AAA-TTAaTTGCCTTTG-GAGCGGAT-ACcATTATAATATGTGG-A-TGTGAGGCGtTCCACAATTCCCACAAACGGTTTCCCTTCTAGAAATTAATTTTGTTTAACCTttTAATCTCCATATATTCAAATATAAGGGTGGATCATATGACcGCA-TGACGGAAGGTGCAAAACTGgTTGAGAAAGAGATCCCGTATATCACCGAACCTGGAAAGGCGACGTCGAAGGTATGAAATTTATCATaAAAGGCGAGGGTACCcGTGACGCGACCACGGGTACCATTAAAGCGAAAATACATCTGCACTACGGGCGACCTGCCGGTCCCGTGGGCAACCCTGGTGAGCACCCTGAGCTACGGTGTTCAGTGTTTCGCCAAGTACCCGAGCCACATCAAGGATTTCTTTAAGAGCGCCATGCCGGAAGGTTATACCCAAGAGCGTACCATCAGCTTCGAAGGCGACGGCGTGTACAAGACGCGTGCTATGGTTACCTACGAACGCGGTTCTATCTACAATCGTGTCACGCTGACTGGTGAGAACTTTAAGAAAGACGGTCACATTCTGCGTAAGAACGTTGCATTCCAATGCCCGCCAAGCATTCTGTATATTCTGCCTGACACCGTTAACAATGGCATCCGCGTTGAGTTCAACCAGGCGTACGATATTGAAGGTGTGACCGAAAAACTGGTTACCAAATGCAGCCAAATGAATCGTCCGTTGGCGGGCTCCGCGGCAGTGCATATCCCGCGTTATCATCACATTACCTACCACACCAAACTGAGCAAAGACCGCGACGAGCGCCGTGATCACATGTGTCTGGTAGAGGTCGTGAAAGCGGTTGATCTGGACACGTATCAGTAATGAGAATTCTGTACACTCGAGccccCCAGGCATCAAATAAAACGAAAGGCTCAGTCGAAAGACTGGGCCTTTCGTTTTATCTGTTGTTTGTCGGTGAACGCTCTCTACTAGAGTCACACTGGCTCACCTTCGGGTGGGCCTTTCTGCGTTTATATGTAGAAGAGCTAATGACCCCAAGGGCGACACCCCCTAATTAGCCCGGGCGAAAGGCCCAGTCTTTCGACTGAGCCTTTCGTTTTATTTGATGCCTGGCAGTTCCCTACTCTCGCATGGGGAGTCCCCACACTACCATCGGCGCTACGGCGTTTCACTTCTGAGTTCGGCATGGGGTCAGGT-GGACCACCGCG-TA-T-CC-CC-GG"),
			ScoringMatrix:     Fitted,
			AlignmentStartPos: 0,
			AlignmentEndPos:   0,
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
			Identity:          1,
			Alignment:         fmt.Sprintf("%s\n%s\n", "CACGGTTGACA", "CACGGTTGACA"),
			ScoringMatrix:     Fitted,
			AlignmentStartPos: 17,
			AlignmentEndPos:   7,
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
			Identity:          1,
			Alignment:         fmt.Sprintf("%s\n%s\n", "GTTGACAGA", "GTTGACAGA"),
			ScoringMatrix:     Fitted,
			AlignmentStartPos: 1,
			AlignmentEndPos:   9,
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
			Identity:          1,
			Alignment:         fmt.Sprintf("%s\n%s\n", "GTTGACAGA", "GTTGACAGA"),
			ScoringMatrix:     Fitted,
			AlignmentStartPos: 9,
			AlignmentEndPos:   1,
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
			Identity:          1,
			Alignment:         fmt.Sprintf("%s\n%s\n", "CACGGTTGACA", "CACGGTTGACA"),
			ScoringMatrix:     Fitted,
			AlignmentStartPos: 7,
			AlignmentEndPos:   17,
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

		longestMatch := alignment.LongestContinuousSequence()

		positions := sequences.FindAll(&test.Seq1, &longestMatch)

		if test.AlignmentStartPos > 0 && test.AlignmentEndPos > 0 {
			if test.AlignmentStartPos != positions.Positions[0].Start() || test.AlignmentEndPos != positions.Positions[0].End() {
				t.Error(
					"For", test.Name, "\n",
					"longest match found: ", longestMatch, "\n",
					"expected start, end:", test.AlignmentStartPos, test.AlignmentEndPos, "\n",
					"got:", fmt.Sprint(positions.Positions[0].Coordinates()), "\n",
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
