// replace_test.go
package sequences

import (
	"testing"

	"github.com/antha-lang/antha/antha/anthalib/wtype"
)

type regionTest struct {
	LargeSeq wtype.DNASequence
	SmallSeq wtype.DNASequence
	Start    int
	End      int
}

var regionTests = []regionTest{

	regionTest{
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
	/*regionTest{
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
	},*/
	regionTest{
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
	regionTest{
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
		result := FindSeqsinSeqs(test.LargeSeq.Seq, []string{test.SmallSeq.Seq})
		if len(result) > 0 {
			if result[0].Positions[0] != test.Start {
				t.Error(
					"For", test.LargeSeq.Nm, "\n",
					"expected", test.Start, "\n",
					"got", result[0].Positions[0], "\n",
				)
			}
		} else {
			t.Error(
				"For", test.LargeSeq.Nm, "\n",
				"expected", test.Start, "\n",
				"got", "no results", "\n",
			)
		}

	}
}
