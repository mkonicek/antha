package liquidhandling

import (
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"reflect"
	"testing"
)

func makeChannel() *wtype.LHChannelParameter {
	vMin := wunit.NewVolume(0.5, "ul")
	vMax := wunit.NewVolume(1000, "ul")
	sMin := wunit.NewFlowRate(0.1, "ml/min")
	sMax := wunit.NewFlowRate(10, "ml/min")
	prm := wtype.NewLHChannelParameter("tecanHead", "tecanEVO", vMin, vMax, sMin, sMax, 8, true, wtype.LHVChannel, 0)
	return prm
}

func compareOutput(t *testing.T, got, want []TipSubset) {
	if len(got) != len(want) {
		t.Errorf("Expected %d subsets, got %d", len(want), len(got))
	}

	for i := 0; i < len(want); i++ {
		if !reflect.DeepEqual(got[i], want[i]) {
			t.Errorf("Expected %v got %v", want[i], got[i])
		}
	}
}

// func makeChannelSubsets(tiptypes []string, channels []*wtype.LHChannelParameter) ([]TipSubset, error)
func TestMakeChannelSubsetOneSubset(t *testing.T) {
	prm := makeChannel()
	tiptypes := []string{"tip", "tip", "tip", "tip", "", "tip", "tip", ""}

	prms := []*wtype.LHChannelParameter{prm, prm, prm, prm, nil, prm, prm, nil}

	ss, err := makeChannelSubsets(tiptypes, prms)

	if err != nil {
		t.Error(err.Error())
	}

	expected := []TipSubset{{Mask: []bool{true, true, true, true, false, true, true, false}, Channel: prm, TipType: "tip"}}

	compareOutput(t, ss, expected)
}

func TestMakeChannelSubsetTwoSubsets(t *testing.T) {
	prm := makeChannel()
	tiptypes := []string{"tip", "tip", "tip2", "tip", "", "tip", "tip", ""}

	prms := []*wtype.LHChannelParameter{prm, prm, prm, prm, nil, prm, prm, nil}

	ss, err := makeChannelSubsets(tiptypes, prms)

	if err != nil {
		t.Error(err.Error())
	}

	expected := []TipSubset{{Mask: []bool{true, true, false, true, false, true, true, false}, Channel: prm, TipType: "tip"}, {Mask: []bool{false, false, true, false, false, false, false, false}, Channel: prm, TipType: "tip2"}}

	compareOutput(t, ss, expected)
}

func TestMakeChannelSubsetTwoSubsetsParams(t *testing.T) {
	prm := makeChannel()
	prm2 := makeChannel()
	tiptypes := []string{"tip", "tip", "tip", "tip", "", "tip", "tip", ""}

	prms := []*wtype.LHChannelParameter{prm, prm, prm2, prm, nil, prm, prm, nil}

	ss, err := makeChannelSubsets(tiptypes, prms)

	if err != nil {
		t.Error(err.Error())
	}

	expected := []TipSubset{{Mask: []bool{true, true, false, true, false, true, true, false}, Channel: prm, TipType: "tip"}, {Mask: []bool{false, false, true, false, false, false, false, false}, Channel: prm2, TipType: "tip"}}

	compareOutput(t, ss, expected)
}
