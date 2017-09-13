// align
package align

import (
	"fmt"

	"testing"

	"github.com/antha-lang/antha/antha/anthalib/wtype"
)

type alignmentTest struct {
	Name          string
	Seq1, Seq2    wtype.DNASequence
	Alignment     string
	Identity      float64
	ScoringMatrix ScoringMatrix
}

var (
	tests []alignmentTest = []alignmentTest{
		alignmentTest{
			Name: "Test1",
			Seq1: wtype.DNASequence{
				Nm:  "Seq1",
				Seq: "GTTGACAGACTAGATTCACG",
			},
			Seq2: wtype.DNASequence{
				Nm:  "Seq2",
				Seq: "GACAGACGA",
			},
			Alignment:     fmt.Sprintf("%s\n%s\n", "GACAGACTAGA", "GACAGAC--GA"),
			Identity:      0.8181818181818182,
			ScoringMatrix: Fitted,
		},
		alignmentTest{
			Name: "Test2",
			Seq1: wtype.DNASequence{
				Nm:  "Seq3",
				Seq: "GTTGACAGACTAGATTCACG",
			},
			Seq2: wtype.DNASequence{
				Nm:  "Seq4",
				Seq: "GTTGACA",
			},
			Identity:      1,
			Alignment:     fmt.Sprintf("%s\n%s\n", "GTTGACA", "GTTGACA"),
			ScoringMatrix: Fitted,
		},
		alignmentTest{
			Name: "Test3",
			Seq1: wtype.DNASequence{
				Nm:  "Seq5",
				Seq: "GTTGACAGACTTTATTCACG",
			},
			Seq2: wtype.DNASequence{
				Nm:  "Seq6",
				Seq: "GTTGAGACAAGATTCACG",
			},
			Alignment:     fmt.Sprintf("%s\n%s\n", "GTTGACAGACtttATTCACG", "GTTGA--GACaagATTCACG"),
			Identity:      0.75,
			ScoringMatrix: FittedAffine,
		},
		alignmentTest{
			Name: "TerminatorAlignmentCorrect",
			Seq1: wtype.DNASequence{
				Nm:  "SequencingResult",
				Seq: "CGCGacgtAaTACGACTcaCTATAGGgCGAATTGGCGGAAGGCCGTCAAGGCCGCATGAAGGGCGCGCCAGGTCTCAGGTACCAGGCATCAAATAAAACGAAAGGCTCAGTCGAAAGACTGGGCCTTTCGTTTTATCTGTTGTTTGTCGGTGAACGCTCTCTACTAGAGTCACACTGGCTCACCTTCGGGTGGGCCTTTCTGCGTTTATATGTAGAAGAGCTAATAGAGACCATTAATTAACAACCTGGGCCTCATGGGCCTTCCGCTCACTGCCCGCTTTCCAGTCGGGAAACCTGTCGTGCCAGCTGCATTAACATGGTCATAGCTGTTTCCTTGCGTATTGGGCGCTCTCCGCTTCCTCGCTCACTGACTCGCTGCGCTCGGTCGTTCGGGTAAAGCCTGGGGTGCCTAATGAGCAAAAGGCCAGCAAAAGGCCAGGAACCGTAAAAAGGCCGCGTTGCTGGCGTTTTTCCATAGGCTCCGCCCCCCTGACGAGCATCACAAAAATCGACGCTCAAGTCAGAGGTGGCGAAACCCGACAGGACTATAAAGATACCAGGCGTTTCCCCCTGGAAGCTCCCTCGTGCGCTCTCCTGTTCCGACCCTGCCGCTTACCGGATACCTGTCCGCCTTTCTCCCTTCGGGAAGCGTGGCGCTTTCTCATAGCTCACGCTGTAGGTATCTCAGTTCGGTGTAGGTCGTTCGCTCCAAGCTGGGCTGTGTGCACGAACCCCCcGTTCAGCCCGACCGCTGCGCCTTATCCGGTAACTATCGTCTTGAGTCCAACCCGGTAAGACACGACTTATCGCCACTGGCAGCAGCCACTGGTAACAGGATTAGCAGAGCGAGGTATGTaggCGGTGCTACAGAGTTCTTGAAGTGGTGGCCTAACTACGGCTACACTagAAGAACAGTATTTGGtaTCTGCGCTCTGCtgAaGCCAGTTaccttcggAAAAanagTTggAaGCTCTTGATCCGGcnaacAAaCcac",
			},
			Seq2: wtype.DNASequence{
				Nm:  "TerSequence",
				Seq: "GGTCTCAGGTACCAGGCATCAAATAAAACGAAAGGCTCAGTCGAAAGACTGGGCCTTTCGTTTTATCTGTTGTTTGTCGGTGAACGCTCTCTACTAGAGTCACACTGGCTCACCTTCGGGTGGGCCTTTCTGCGTTTATATGTAGAAGAGCTAATAGAGACC",
			},
			Identity:      1,
			Alignment:     fmt.Sprintf("%s\n%s\n", "GGTCTCAGGTACCAGGCATCAAATAAAACGAAAGGCTCAGTCGAAAGACTGGGCCTTTCGTTTTATCTGTTGTTTGTCGGTGAACGCTCTCTACTAGAGTCACACTGGCTCACCTTCGGGTGGGCCTTTCTGCGTTTATATGTAGAAGAGCTAATAGAGACC", "GGTCTCAGGTACCAGGCATCAAATAAAACGAAAGGCTCAGTCGAAAGACTGGGCCTTTCGTTTTATCTGTTGTTTGTCGGTGAACGCTCTCTACTAGAGTCACACTGGCTCACCTTCGGGTGGGCCTTTCTGCGTTTATATGTAGAAGAGCTAATAGAGACC"),
			ScoringMatrix: Fitted,
		},
		alignmentTest{
			Name: "MismatchingAlignmentReverse",
			Seq1: wtype.DNASequence{
				Nm:  "TemplateSequence",
				Seq: "ctcatgaccaaaatcccttaacgtgagttacgcgcgcgtcgttccactgagcgtcagaccccgtagaaaagatcaaaggatcttcttgagatcctttttttctgcgcgtaatctgctgcttgcaaacaaaaaaaccaccgctaccagcggtggtttgtttgccggatcaagagctaccaactctttttccgaaggtaactggcttcagcagagcgcagataccaaatactgttcttctagtgtagccgtagttagcccaccacttcaagaactctgtagcaccgcctacatacctcgctctgctaatcctgttaccagtggctgctgccagtggcgataagtcgtgtcttaccgggttggactcaagacgatagttaccggataaggcgcagcggtcgggctgaacggggggttcgtgcacacagcccagcttggagcgaacgacctacaccgaactgagatacctacagcgtgagctatgagaaagcgccacgcttcccgaagggagaaaggcggacaggtatccggtaagcggcagggtcggaacaggagagcgcacgagggagcttccagggggaaacgcctggtatctttatagtcctgtcgggtttcgccacctctgacttgagcgtcgatttttgtgatgctcgtcaggggggcggagcctatggaaaaacgccagcaacgcggcctttttacggttcctggccttttgctggccttttgctcacatgttctttcctgcgttatcccctgattctgtggataaccgtattaccgcctttgagtgagctgataccgctcgccgcagccgaacgaccgagcgcagcgagtcagtgagcgaggaagcggaaggcgagagtagggaactgccaggcatcaaactaagcagaaggcccctgacggatggcctttttgcgtttctacaaactctttctgtgttgtaaaacgacggccagtcttaagctcgggccccctgggcggttctgataacgagtaatcgttaatccgcaaataacgtaaaaacccgcttcggcgggtttttttatggggggagtttagggaaagagcatttgtcagaatatttaagggcgcctgtcactttgcttgatatatgagaattatttaaccttataaatgagaaaaaagcaacgcactttaaataagatacgttgctttttcgattgatgaacacctataattaaactattcatctattatttatgattttttgtatatacaatatttctagtttgttaaagagaattaagaaaataaatctcgaaaataataaagggaaaatcagtttttgatatcaaaattatacatgtcaacgataatacaaaatataatacaaactataagatgttatcagtatttattatgcatttagaataaattttgtGTCGGCTCTTCAACCTGTGTCGCCCTTAATTGTGAGCGGATAACAATTACGAGCTTCATGCACAGTGAAATCATGAAAAATTTATTTGCTTTGTGAGCGGATAACAATTATAATATGTGGAATTGTGAGCGCTCACAATTCCACAACGGTTTCCCTCTAGAAATAATTTTGTTTAACTACTAATCTCATATATCAAATATAGGGTGGATCATATGACGGCATTGACGGAAGGTGCAAAACTGTTTGAGAAAGAGATCCCGTATATCACCGAACTGGAAGGCGACGTCGAAGGTATGAAATTTATCATTAAAGGCGAGGGTACCGGTGACGCGACCACGGGTACCATTAAAGCGAAATACATCTGCACTACGGGCGACCTGCCGGTCCCGTGGGCAACCCTGGTGAGCACCCTGAGCTACGGTGTTCAGTGTTTCGCCAAGTACCCGAGCCACATCAAGGATTTCTTTAAGAGCGCCATGCCGGAAGGTTATACCCAAGAGCGTACCATCAGCTTCGAAGGCGACGGCGTGTACAAGACGCGTGCTATGGTTACCTACGAACGCGGTTCTATCTACAATCGTGTCACGCTGACTGGTGAGAACTTTAAGAAAGACGGTCACATTCTGCGTAAGAACGTTGCATTCCAATGCCCGCCAAGCATTCTGTATATTCTGCCTGACACCGTTAACAATGGCATCCGCGTTGAGTTCAACCAGGCGTACGATATTGAAGGTGTGACCGAAAAACTGGTTACCAAATGCAGCCAAATGAATCGTCCGTTGGCGGGCTCCGCGGCAGTGCATATCCCGCGTTATCATCACATTACCTACCACACCAAACTGAGCAAAGACCGCGACGAGCGCCGTGATCACATGTGTCTGGTAGAGGTCGTGAAAGCGGTTGATCTGGACACGTATCAGTAATGAGAATTCTGTACACTCGAGGGTACCAGGCATCAAATAAAACGAAAGGCTCAGTCGAAAGACTGGGCCTTTCGTTTTATCTGTTGTTTGTCGGTGAACGCTCTCTACTAGAGTCACACTGGCTCACCTTCGGGTGGGCCTTTCTGCGTTTATATGTAGAAGAGCtaatgaccccaagggcgacaccccctaattagcccgggcgaaaggcccagtctttcgactgagcctttcgttttatttgatgcctggcagttccctactctcgcatggggagtccccacactaccatcggcgctacggcgtttcacttctgagttcggcatggggtcaggtgggaccaccgcgctactgccgccaggcaaacaaggggtgttatgagccatattcaggtataaatgggctcgcgataatgttcagaattggttaattggttgtaacactgacccctatttgtttatttttctaaatacattcaaatatgtatccgctcatgagacaataaccctgataaatgcttcaataatattgaaaaaggaagaatatgagccatattcaacgggaaacgtcgaggccgcgattaaattccaacatggatgctgatttatatgggtataaatgggctcgcgataatgtcgggcaatcaggtgcgacaatctatcgcttgtatgggaagcccgatgcgccagagttgtttctgaaacatggcaaaggtagcgttgccaatgatgttacagatgagatggtcagactaaactggctgacggaatttatgccacttccgaccatcaagcattttatccgtactcctgatgatgcatggttactcaccactgcgatccccggaaaaacagcgttccaggtattagaagaatatcctgattcaggtgaaaatattgttgatgcgctggcagtgttcctgcgccggttgcactcgattcctgtttgtaattgtccttttaacagcgatcgcgtatttcgcctcgctcaggcgcaatcacgaatgaataacggtttggttgatgcgagtgattttgatgacgagcgtaatggctggcctgttgaacaagtctggaaagaaatgcataaacttttgccattctcaccggattcagtcgtcactcatggtgatttctcacttgataaccttatttttgacgaggggaaattaataggttgtattgatgttggacgagtcggaatcgcagaccgataccaggatcttgccatcctatggaactgcctcggtgagttttctccttcattacagaaacggctttttcaaaaatatggtattgataatcctgatatgaataaattgcagtttcatttgatgctcgatgagtttttctaagcggcgcgccatcgaatggcgcaaaacctttcgcggtatggcatgatagcgcccggaagagagtcaattcagggtggtgaatatgaaaccagtaacgttatacgatgtcgcagagtatgccggtgtctcttatcagaccgtttcccgcgtggtgaaccaggccagccacgtttctgcgaaaacgcgggaaaaagtggaagcggcgatggcggagctgaattacattcccaaccgcgtggcacaacaactggcgggcaaacagtcgttgctgattggcgttgccacctccagtctggccctgcacgcgccgtcgcaaattgtcgcggcgattaaatctcgcgccgatcaactgggtgccagcgtggtggtgtcgatggtagaacgaagcggcgtcgaagcctgtaaagcggcggtgcacaatcttctcgcgcaacgcgtcagtgggctgatcattaactatccgctggatgaccaggatgccattgctgtggaagctgcctgcactaatgttccggcgttatttcttgatgtctctgaccagacacccatcaacagtattattttctcccatgaggacggtacgcgactgggcgtggagcatctggtcgcattgggtcaccagcaaatcgcgctgttagcgggcccattaagttctgtctcggcgcgtctgcgtctggctggctggcataaatatctcactcgcaatcaaattcagccgatagcggaacgggaaggcgactggagtgccatgtccggttttcaacaaaccatgcaaatgctgaatgagggcatcgttcccactgcgatgctggttgccaacgatcagatggcgctgggcgcaatgcgcgccattaccgagtccgggctgcgcgttggtgcggatatctcggtagtgggatacgacgataccgaagatagctcatgttatatcccgccgttaaccaccatcaaacaggattttcgcctgctggggcaaaccagcgtggaccgcttgctgcaactctctcagggccaggcggtgaagggcaatcagctgttgccagtctcactggtgaaaagaaaaaccaccctggcgcccaatacgcaaaccgcctctccccgcgcgttggccgattcattaatgcagctggcacgacaggtttcccgactggaaagcgggcagtga",
			},
			Seq2: wtype.DNASequence{
				Nm:  "ReverseSequencingResult",
				Seq: wtype.RevComp("CCGGGGATACGCGGTGGTCCACCTGACCCCATGCCGAACTCAGAAGTGAAACGCCGTAGCGCCGATGGTAGTGTGGGGACTCCCCATGCGAGAGTAGGGAACTGCCAGGCATCAAATAAAACGAAAGGCTCAGTCGAAAGACTGGGCCTTTCGCCCGGGCTAATTAGGGGGTGTCGCCCTTGGGGTCATTAGCTCTTCTACATATAAACGCAGAAAGGCCCACCCGAAGGTGAGCCAGTGTGACTCTAGTAGAGAGCGTTCACCGACAAACAACAGATAAAACGAAAGGCCCAGTCTTTCGACTGAGCCTTTCGTTTTATTTGATGCCTGGGGGGCTCGAGTGTACAGAATTCTCATTACTGATACGTGTCCAGATCAACCGCTTTCACGACCTCTACCAGACACATGTGATCACGGCGCTCGTCGCGGTCTTTGCTCAGTTTGGTGTGGTAGGTAATGTGATGATAACGCGGGATATGCACTGCCGCGGAGCCCGCCAACGGACGATTCATTTGGCTGCATTTGGTAACCAGTTTTTCGGTCACACCTTCAATATCGTACGCCTGGTTGAACTCAACGCGGATGCCATTGTTAACGGTGTCAGGCAGAATATACAGAATGCTTGGCGGGCATTGGAATGCAACGTTCTTACGCAGAATGTGACCGTCTTTCTTAAAGTTCTCACCAGTCAGCGTGACACGATTGTAGATAGAACCGCGTTCGTAGGTAACCATAGCACGCGTCTTGTACACGCCGTCGCCTTCGAAGCTGATGGTACGCTCTTGGGTATAACCTTCCGGCATGGCGCTCTTAAAGAAATCCTTGATGTGGCTCGGGTACTTGGCGAAACACTGAACACCGTAGCTCAGGGTGCTCACCAGGGTTGCCCACGGGACCGGCAGGTCGCCCGTAGTGCAGATGTATTTTCGCTTTAATGGTACCCGTGGTCGCGTCACGGGTACCCTCGCCTTTTATGATAAATTTCATACCTTCGACGTCGCCTTTCCAGGTTCGGTGATATACGGGATCTCTTTCTCAACCAGTTTTGCACCTTCCGTCATGCGGTCATATGATCCACCCTTATATTTGAATATATGGAGATTAAAAGGTTAAACAAAATTAATTTCTAGAAGGGAAACCGTTTGTGGGAATTGTGGAACGCCTCACATCCACATATTATAATGGTATCCGCTCCAAAGGCAATTAATTTCATGGAATACTGTCATGAAGACTCGAATGTTATCGCTTCACAATTAAGGGGGCGACACATCGGTGTGTAGA"),
			},
			Identity:      0.9561875480399693,
			Alignment:     fmt.Sprintf("%s\n%s\n", "TCTtCA-ACC--TGTGTCG--CCCTTAATTGTG-AGCGGATAACAATTACGAG-CTTCATGCACAGTGAAaT-CATGAAAAATTTAtTTG-CTTTGTGAGCGGATAACaATTATAATATGTGGAATTGTGA-GCGcT-CACAATT-CCAC-AACGGTTTCCC-TCTAGAAA-TAATTTTGTTTAA-CTacTAATCT-CATATA-TCAAATAT-AGGGTGGATCATATGACgGCATTGACGGAAGGTGCAAAACTGtTTGAGAAAGAGATCCCGTATATCACCGAA-CTGG-AAGGCGACGTCGAAGGTATGAAATTTATCATtAAAGGCGAGGGTACCgGTGACGCGACCACGGGTACCATTAAAGCG-AAATACATCTGCACTACGGGCGACCTGCCGGTCCCGTGGGCAACCCTGGTGAGCACCCTGAGCTACGGTGTTCAGTGTTTCGCCAAGTACCCGAGCCACATCAAGGATTTCTTTAAGAGCGCCATGCCGGAAGGTTATACCCAAGAGCGTACCATCAGCTTCGAAGGCGACGGCGTGTACAAGACGCGTGCTATGGTTACCTACGAACGCGGTTCTATCTACAATCGTGTCACGCTGACTGGTGAGAACTTTAAGAAAGACGGTCACATTCTGCGTAAGAACGTTGCATTCCAATGCCCGCCAAGCATTCTGTATATTCTGCCTGACACCGTTAACAATGGCATCCGCGTTGAGTTCAACCAGGCGTACGATATTGAAGGTGTGACCGAAAAACTGGTTACCAAATGCAGCCAAATGAATCGTCCGTTGGCGGGCTCCGCGGCAGTGCATATCCCGCGTTATCATCACATTACCTACCACACCAAACTGAGCAAAGACCGCGACGAGCGCCGTGATCACATGTGTCTGGTAGAGGTCGTGAAAGCGGTTGATCTGGACACGTATCAGTAATGAGAATTCTGTACACTCGAGggtaCCAGGCATCAAATAAAACGAAAGGCTCAGTCGAAAGACTGGGCCTTTCGTTTTATCTGTTGTTTGTCGGTGAACGCTCTCTACTAGAGTCACACTGGCTCACCTTCGGGTGGGCCTTTCTGCGTTTATATGTAGAAGAGCtaatgaccccaagggcgacaccccctaattagcccgggcgaaaggcccagtctttcgactgagcctttcgttttatttgatgcctggcagttccctactctcgcatggggagtccccacactaccatcggcgctacggcgtttcacttctgagttcggcatggggtcaggtgggaccaccgcgctactgccgccagg", "TCTaCACACCGATGTGTCGCCCCCTTAATTGTGAAGC-GATAAC-ATT-CGAGTCTTCATG-ACAGT--AtTCCATG--AAA-TTAaTTGCCTTTG-GAGCGGAT-ACcATTATAATATGTGG-A-TGTGAGGCGtTCCACAATTCCCACAAACGGTTTCCCTTCTAGAAATTAATTTTGTTTAACCTttTAATCTCCATATATTCAAATATAAGGGTGGATCATATGACcGCA-TGACGGAAGGTGCAAAACTGgTTGAGAAAGAGATCCCGTATATCACCGAACCTGGAAAGGCGACGTCGAAGGTATGAAATTTATCATaAAAGGCGAGGGTACCcGTGACGCGACCACGGGTACCATTAAAGCGAAAATACATCTGCACTACGGGCGACCTGCCGGTCCCGTGGGCAACCCTGGTGAGCACCCTGAGCTACGGTGTTCAGTGTTTCGCCAAGTACCCGAGCCACATCAAGGATTTCTTTAAGAGCGCCATGCCGGAAGGTTATACCCAAGAGCGTACCATCAGCTTCGAAGGCGACGGCGTGTACAAGACGCGTGCTATGGTTACCTACGAACGCGGTTCTATCTACAATCGTGTCACGCTGACTGGTGAGAACTTTAAGAAAGACGGTCACATTCTGCGTAAGAACGTTGCATTCCAATGCCCGCCAAGCATTCTGTATATTCTGCCTGACACCGTTAACAATGGCATCCGCGTTGAGTTCAACCAGGCGTACGATATTGAAGGTGTGACCGAAAAACTGGTTACCAAATGCAGCCAAATGAATCGTCCGTTGGCGGGCTCCGCGGCAGTGCATATCCCGCGTTATCATCACATTACCTACCACACCAAACTGAGCAAAGACCGCGACGAGCGCCGTGATCACATGTGTCTGGTAGAGGTCGTGAAAGCGGTTGATCTGGACACGTATCAGTAATGAGAATTCTGTACACTCGAGccccCCAGGCATCAAATAAAACGAAAGGCTCAGTCGAAAGACTGGGCCTTTCGTTTTATCTGTTGTTTGTCGGTGAACGCTCTCTACTAGAGTCACACTGGCTCACCTTCGGGTGGGCCTTTCTGCGTTTATATGTAGAAGAGCtaatgaccccaagggcgacaccccctaattagcccgggcgaaaggcccagtctttcgactgagcctttcgttttatttgatgcctggcagttccctactctcgcatggggagtccccacactaccatcggcgctacggcgtttcacttctgagttcggcatggggtcaggt-ggaccaccgcg-ta-t-cc-cc-gg"),
			ScoringMatrix: Fitted,
		},
		alignmentTest{
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
			Identity:      1,
			Alignment:     fmt.Sprintf("%s\n%s\n", "CACGGTTGACA", "CACGGTTGACA"),
			ScoringMatrix: Fitted,
		},
	}
)

// Align two dna sequences based on a specified scoring matrix
func TestAlign(t *testing.T) {
	for _, test := range tests {
		alignment, err := DNA(test.Seq1, test.Seq2, test.ScoringMatrix)

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
		fmt.Println(test.Name, "Coverage:", alignment.Coverage(),
			"Positions in template: ", alignment.Positions().Positions,
			"Matches: ", alignment.Matches(),
			"Mismatches: ", alignment.Mismatches(),
			"Gaps: ", alignment.Gaps(),
			"Longest Matching Sequence:", alignment.LongestContinuousSequence(),
		)
	}
}
