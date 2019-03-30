package align

import (
	"testing"
)

func TestAlignmentMethods(t *testing.T) {

	type sectionData struct {
		TStart, TEnd, QStart, QEnd int
		TResult, QResult, QTMatch  string
	}

	type testData struct {
		maxSectionLength int
		sections         []sectionData
	}

	testAln := Alignment{
		TemplateResult:    "GCTTTTTTATAATGCCAACTTTG-----AAAAG",
		QueryResult:       "GGG-TTTTATAATGCCAACTTTGTACATTAAAG",
		TemplatePositions: []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 23, 23, 23, 23, 23, 24, 25, 26, 27, 28},
		QueryPositions:    []int{33, 32, 31, 31, 30, 29, 28, 27, 26, 25, 24, 23, 22, 21, 20, 19, 18, 17, 16, 15, 14, 13, 12, 11, 10, 9, 8, 7, 6, 5, 4, 3, 2},
	}

	for _, test := range []testData{
		{
			maxSectionLength: 10,
			sections: []sectionData{
				{
					TStart: 1, TEnd: 10, QStart: 33, QEnd: 25,
					TResult: "GCTTTTTTAT",
					QTMatch: "|   ||||||",
					QResult: "GGG-TTTTAT",
				},
				{
					TStart: 11, TEnd: 20, QStart: 24, QEnd: 15,
					TResult: "AATGCCAACT",
					QTMatch: "||||||||||",
					QResult: "AATGCCAACT",
				},
				{
					TStart: 21, TEnd: 25, QStart: 14, QEnd: 5,
					TResult: "TTG-----AA",
					QTMatch: "|||      |",
					QResult: "TTGTACATTA",
				},
				{
					TStart: 26, TEnd: 28, QStart: 4, QEnd: 2,
					TResult: "AAG",
					QTMatch: "|||",
					QResult: "AAG",
				},
			},
		},
		{
			maxSectionLength: 20,
			sections: []sectionData{
				{
					TStart: 1, TEnd: 20, QStart: 33, QEnd: 15,
					TResult: "GCTTTTTTATAATGCCAACT",
					QTMatch: "|   ||||||||||||||||",
					QResult: "GGG-TTTTATAATGCCAACT",
				},
				{
					TStart: 21, TEnd: 28, QStart: 14, QEnd: 2,
					TResult: "TTG-----AAAAG",
					QTMatch: "|||      ||||",
					QResult: "TTGTACATTAAAG",
				},
			},
		},
		{
			maxSectionLength: 500,
			sections: []sectionData{
				{
					TStart: 1, TEnd: 28, QStart: 33, QEnd: 2,
					TResult: "GCTTTTTTATAATGCCAACTTTG-----AAAAG",
					QTMatch: "|   |||||||||||||||||||      ||||",
					QResult: "GGG-TTTTATAATGCCAACTTTGTACATTAAAG",
				},
			},
		},
	} {

		got, err := testAln.Split(test.maxSectionLength)
		if err != nil {
			t.Fatalf(err.Error())
		}
		if len(got) != len(test.sections) {
			t.Errorf("unexpected number of sections: got %d, want %d\n", len(got), len(test.sections))
		}
		for i := 0; i < len(test.sections); i++ {
			want := test.sections[i]
			if want.TStart != got[i].TemplateStart() {
				t.Errorf("Error, TemplateStart: got %d, want %d\n", want.TStart, got[i].TemplateStart())
			}
			if want.TEnd != got[i].TemplateEnd() {
				t.Errorf("Error, TEnd: got %d, want %d\n", want.TEnd, got[i].TemplateEnd())
			}
			if want.QStart != got[i].QueryStart() {
				t.Errorf("Error, QueryStart: got %d, want %d\n", want.QStart, got[i].QueryStart())
			}
			if want.QEnd != got[i].QueryEnd() {
				t.Errorf("Error, QueryEnd: got %d, want %d\n", want.QEnd, got[i].QueryEnd())
			}
			if want.QTMatch != got[i].Match() {
				t.Errorf("Error, Match: got %s, want %s\n", want.QTMatch, got[i].Match())
			}
		}
	}

}

func TestInferFrame(t *testing.T) {

	type testData struct {
		positions []int
		frame     int
		name      string
	}

	for _, test := range []testData{
		{[]int{}, 1, "inconclusive, no data"},
		{[]int{1}, 1, "inconclusive, too short"},
		{[]int{1, 1, 1, 1}, 1, "inconclusive"},
		{[]int{1, 2}, 1, "ungapped fwd"},
		{[]int{2, 1}, -1, "ungapped rev"},
		{[]int{2, 2, 3}, 1, "gapped fwd"},
		{[]int{2, 2, 1}, -1, "gapped rev"},
		{[]int{5, 6, 7, 8, 9, 10, 1, 2, 3}, 1, "plasmid fwd"},
		{[]int{5, 4, 3, 2, 1, 10, 9, 8, 7}, -1, "plasmid rev"},
	} {
		got := inferFrame(test.positions)
		if got != test.frame {
			t.Errorf("Error, %s: got %d, want %d\n", test.name, got, test.frame)
		}
	}
}
