package tests

import (
	"reflect"
	"testing"

	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/laboratory/effects/id"
	"github.com/antha-lang/antha/microArch/driver/liquidhandling"
)

func makeChannel(idGen *id.IDGenerator) *wtype.LHChannelParameter {
	vMin := wunit.NewVolume(0.5, "ul")
	vMax := wunit.NewVolume(1000, "ul")
	sMin := wunit.NewFlowRate(0.1, "ml/min")
	sMax := wunit.NewFlowRate(10, "ml/min")
	prm := wtype.NewLHChannelParameter(idGen, "tecanHead", "tecanEVO", vMin, vMax, sMin, sMax, 8, true, wtype.LHVChannel, 0)
	return prm
}

func compareOutput(t *testing.T, got, want []liquidhandling.TipSubset) {
	if len(got) != len(want) {
		t.Errorf("Expected %d subsets, got %d", len(want), len(got))
	}

	for i := 0; i < len(want); i++ {
		if !reflect.DeepEqual(got[i], want[i]) {
			t.Errorf("Expected %v got %v", want[i], got[i])
		}
	}
}

// func MakeChannelSubsets(tiptypes []string, channels []*wtype.LHChannelParameter) ([]TipSubset, error)
func TestMakeChannelSubsetOneSubset(t *testing.T) {
	idGen := id.NewIDGenerator(t.Name())
	prm := makeChannel(idGen)
	tiptypes := []string{"tip", "tip", "tip", "tip", "", "tip", "tip", ""}

	prms := []*wtype.LHChannelParameter{prm, prm, prm, prm, nil, prm, prm, nil}

	ss, err := liquidhandling.MakeChannelSubsets(tiptypes, prms)

	if err != nil {
		t.Error(err.Error())
	}

	expected := []liquidhandling.TipSubset{{Mask: []bool{true, true, true, true, false, true, true, false}, Channel: prm, TipType: "tip"}}

	compareOutput(t, ss, expected)
}

func TestMakeChannelSubsetTwoSubsets(t *testing.T) {
	idGen := id.NewIDGenerator(t.Name())
	prm := makeChannel(idGen)
	tiptypes := []string{"tip", "tip", "tip2", "tip", "", "tip", "tip", ""}

	prms := []*wtype.LHChannelParameter{prm, prm, prm, prm, nil, prm, prm, nil}

	ss, err := liquidhandling.MakeChannelSubsets(tiptypes, prms)

	if err != nil {
		t.Error(err.Error())
	}

	expected := []liquidhandling.TipSubset{{Mask: []bool{true, true, false, true, false, true, true, false}, Channel: prm, TipType: "tip"}, {Mask: []bool{false, false, true, false, false, false, false, false}, Channel: prm, TipType: "tip2"}}

	compareOutput(t, ss, expected)
}

func TestMakeChannelSubsetTwoSubsetsParams(t *testing.T) {
	idGen := id.NewIDGenerator(t.Name())
	prm := makeChannel(idGen)
	prm2 := makeChannel(idGen)
	tiptypes := []string{"tip", "tip", "tip", "tip", "", "tip", "tip", ""}

	prms := []*wtype.LHChannelParameter{prm, prm, prm2, prm, nil, prm, prm, nil}

	ss, err := liquidhandling.MakeChannelSubsets(tiptypes, prms)

	if err != nil {
		t.Error(err.Error())
	}

	expected := []liquidhandling.TipSubset{{Mask: []bool{true, true, false, true, false, true, true, false}, Channel: prm, TipType: "tip"}, {Mask: []bool{false, false, true, false, false, false, false, false}, Channel: prm2, TipType: "tip"}}

	compareOutput(t, ss, expected)
}
