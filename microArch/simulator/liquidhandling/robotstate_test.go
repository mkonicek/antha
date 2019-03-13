package liquidhandling

import (
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"strings"
	"testing"
)

type tipCoordsTest struct {
	testName     string
	tipsMissing  []string
	num          int
	tipBehaviour wtype.TipLoadingBehaviour
	orientation  wtype.ChannelOrientation
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
	tb := defaultLHTipbox("testbox")

	for _, tipAddrS := range self.tipsMissing {
		wc := wtype.MakeWellCoords(tipAddrS)
		tb.RemoveTip(wc)
	}

	adaptor := NewAdaptorState("", false, 8, wtype.Coordinates3D{}, 0.0, &wtype.LHChannelParameter{Orientation: self.orientation}, self.tipBehaviour)

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
			testName: "single channel LTR - one tip remaining",
			tipsMissing: []string{
				"A1", "A2", "A3", "A4", "A5", "A6", "A7", "A8", "A9", "A10", "A11", //"A12",
				"B1", "B2", "B3", "B4", "B5", "B6", "B7", "B8", "B9", "B10", "B11", "B12",
				"C1", "C2", "C3", "C4", "C5", "C6", "C7", "C8", "C9", "C10", "C11", "C12",
				"D1", "D2", "D3", "D4", "D5", "D6", "D7", "D8", "D9", "D10", "D11", "D12",
				"E1", "E2", "E3", "E4", "E5", "E6", "E7", "E8", "E9", "E10", "E11", "E12",
				"F1", "F2", "F3", "F4", "F5", "F6", "F7", "F8", "F9", "F10", "F11", "F12",
				"G1", "G2", "G3", "G4", "G5", "G6", "G7", "G8", "G9", "G10", "G11", "G12",
				"H1", "H2", "H3", "H4", "H5", "H6", "H7", "H8", "H9", "H10", "H11", "H12",
			},
			num:          1,
			tipBehaviour: LTRTipBehaviour(),
			orientation:  wtype.LHVChannel,
			expected:     "[[A12]]",
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
