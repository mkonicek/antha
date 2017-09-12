// align
package align

import (
	"fmt"

	"testing"

	"github.com/antha-lang/antha/antha/anthalib/wtype"
)

type alignMentTest struct {
	Name          string
	Seq1, Seq2    wtype.DNASequence
	Alignment     string
	ScoringMatrix ScoringMatrix
}

var (
	tests []alignMentTest = []alignMentTest{
		alignMentTest{
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
			ScoringMatrix: Fitted,
		},
		alignMentTest{
			Name: "Test2",
			Seq1: wtype.DNASequence{
				Nm:  "Seq3",
				Seq: "GTTGACAGACTAGATTCACG",
			},
			Seq2: wtype.DNASequence{
				Nm:  "Seq4",
				Seq: "GTTGACA",
			},
			Alignment:     fmt.Sprintf("%s\n%s\n", "GTTGACA", "GTTGACA"),
			ScoringMatrix: Fitted,
		},
		alignMentTest{
			Name: "Test3",
			Seq1: wtype.DNASequence{
				Nm:  "Seq5",
				Seq: "GTTGACAGACTTTATTCACG",
			},
			Seq2: wtype.DNASequence{
				Nm:  "Seq6",
				Seq: "GTTGAGACAAGATTCACG",
			},
			Alignment:     fmt.Sprintf("%s\n%s\n", "GTTGACAGACTTTATTCACG", "GTTGA--GACAAGATTCACG"),
			ScoringMatrix: FittedAffine,
		},
	}
)

// Align two dna sequences based on a specified scoring matrix
func TestAlign(t *testing.T) {
	for _, test := range tests {
		alignment, err := Align(test.Seq1, test.Seq2, test.ScoringMatrix)

		if err != nil {
			t.Error(
				"For", test.Name, "\n",
				"got error:", err.Error(), "\n",
			)
		}
		if alignment != test.Alignment {
			t.Error(
				"For", test.Name, "\n",
				"expected:", test.Alignment, "\n",
				"got:", alignment, "\n",
			)
		}
	}
}
