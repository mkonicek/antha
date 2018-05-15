package liquidhandling_test

import (
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	lh "github.com/antha-lang/antha/microArch/simulator/liquidhandling"
	"strings"
	"testing"
)

type tipCoordsTest struct {
	testName     string
	tipsMissing  []string
	num          int
	tipBehaviour wtype.TipLoadingBehaviour
	orientation  int
	expected     string
}

func LTRTipBehaviour() wtype.TipLoadingBehaviour {
	return wtype.TipLoadingBehaviour{
		OverrideLoadTipsCommand:    true,
		AutoRefillTipboxes:         true,
		LoadingOrder:               wtype.ColumnWise,
		VerticalLoadingDirection:   wtype.BottomToTop,
		HorizontalLoadingDirection: wtype.LeftToRight,
		ChunkingBehaviour:          wtype.ReverseSequentialTipLoading,
	}
}

func RTLTipBehaviour() wtype.TipLoadingBehaviour {
	return wtype.TipLoadingBehaviour{
		OverrideLoadTipsCommand:    true,
		AutoRefillTipboxes:         true,
		LoadingOrder:               wtype.ColumnWise,
		VerticalLoadingDirection:   wtype.BottomToTop,
		HorizontalLoadingDirection: wtype.RightToLeft,
		ChunkingBehaviour:          wtype.ReverseSequentialTipLoading,
	}
}

func (self *tipCoordsTest) run(t *testing.T) {
	tb := default_lhtipbox("testbox")

	for _, tipAddrS := range self.tipsMissing {
		wc := wtype.MakeWellCoords(tipAddrS)
		tb.RemoveTip(wc)
	}

	adaptor := lh.NewAdaptorState("", false, 8, wtype.Coordinates{}, &wtype.LHChannelParameter{Orientation: self.orientation}, self.tipBehaviour)

	out, err := adaptor.GetTipCoordsToLoad(tb, self.num)
	if err != nil {
		t.Fatal(err)
	}

	var gotS []string
	for _, wcS := range out {
		var wcsS []string
		for _, wc := range wcS {
			wcsS = append(wcsS, wc.FormatA1())
		}
		gotS = append(gotS, "["+strings.Join(wcsS, ",")+"]")
	}
	got := "[" + strings.Join(gotS, ",") + "]"

	if got != self.expected {
		t.Errorf("In test %s:\n  e: \"%s\",\n  g: \"%s\"", self.testName, self.expected, got)
	}

}

func TestGetTipCoordsToLoad(t *testing.T) {

	tests := []tipCoordsTest{
		{
			testName:     "single channel RTL",
			tipsMissing:  []string{},
			num:          1,
			tipBehaviour: RTLTipBehaviour(),
			orientation:  wtype.LHVChannel,
			expected:     "[[H12]]",
		},
		{
			testName:     "single channel LTR",
			tipsMissing:  []string{},
			num:          1,
			tipBehaviour: LTRTipBehaviour(),
			orientation:  wtype.LHVChannel,
			expected:     "[[H1]]",
		},
		{
			testName:     "single channel RTL - missing tip",
			tipsMissing:  []string{"H12"},
			num:          1,
			tipBehaviour: RTLTipBehaviour(),
			orientation:  wtype.LHVChannel,
			expected:     "[[G12]]",
		},
		{
			testName:     "single channel LTR - missing tip",
			tipsMissing:  []string{"H1"},
			num:          1,
			tipBehaviour: LTRTipBehaviour(),
			orientation:  wtype.LHVChannel,
			expected:     "[[G1]]",
		},
		{
			testName:     "multi RTL",
			tipsMissing:  []string{},
			num:          8,
			tipBehaviour: RTLTipBehaviour(),
			orientation:  wtype.LHVChannel,
			expected:     "[[A12,B12,C12,D12,E12,F12,G12,H12]]",
		},
		{
			testName:     "multi LTR",
			tipsMissing:  []string{},
			num:          8,
			tipBehaviour: LTRTipBehaviour(),
			orientation:  wtype.LHVChannel,
			expected:     "[[A1,B1,C1,D1,E1,F1,G1,H1]]",
		},
		{
			testName:     "multi RTL - chunking",
			tipsMissing:  []string{"E12", "F12", "G12", "H12"},
			num:          8,
			tipBehaviour: RTLTipBehaviour(),
			orientation:  wtype.LHVChannel,
			expected:     "[[E11,F11,G11,H11],[A12,B12,C12,D12]]",
		},
		{
			testName:     "multi LTR - chunking",
			tipsMissing:  []string{"E1", "F1", "G1", "H1"},
			num:          8,
			tipBehaviour: LTRTipBehaviour(),
			orientation:  wtype.LHVChannel,
			expected:     "[[E2,F2,G2,H2],[A1,B1,C1,D1]]",
		},
	}

	for _, test := range tests {
		test.run(t)
	}
}
