package liquidhandling

import (
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"testing"
)

type testCase struct {
	Name     string
	Cmps     wtype.ComponentVector
	Match    wtype.Match
	Expected wtype.ComponentVector
}

func (tc testCase) Run(t *testing.T) {

	got := updateDests(tc.Cmps, tc.Match)

	if !got.Equal(tc.Expected) {
		t.Errorf("%s: Expected %v got %v", tc.Name, tc.Expected, got)
	}
}

func TestUpdateDests(t *testing.T) {
	// func updateDests(dst wtype.ComponentVector, match wtype.Match) wtype.ComponentVector

	testCases := []testCase{
		{
			Name:     "Regression1 - Don't lose relatively small transfers",
			Cmps:     wtype.ComponentVector{&wtype.Liquid{Vol: 0.01, Vunit: "El"}},
			Match:    wtype.Match{M: []int{0}, Vols: []wunit.Volume{wunit.NewVolume(0.00999, "El")}},
			Expected: wtype.ComponentVector{&wtype.Liquid{Vol: 0.00001, Vunit: "El"}},
		},
		{
			Name:     "Positive - round very small volumes down to zero",
			Cmps:     wtype.ComponentVector{&wtype.Liquid{Vol: 0.1, Vunit: "ul"}},
			Match:    wtype.Match{M: []int{0}, Vols: []wunit.Volume{wunit.NewVolume(0.9999999, "ul")}},
			Expected: wtype.ComponentVector{&wtype.Liquid{Vol: 0.0, Vunit: "ul"}},
		},
	}

	for _, tc := range testCases {
		tc.Run(t)
	}
}
