package liquidhandling

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
)

type DistributeVolumesTest struct {
	Name            string
	Requested       []wunit.Volume   // Volumes to request per channel
	AvailableByWell [][]wunit.Volume // Volumes available by well then channel - total number of volumes should match requested
	ExpectedUl      []float64        // Expected volumes supplied to each channel
}

func (test *DistributeVolumesTest) Run(t *testing.T) {
	requested := make(wtype.ComponentVector, 0, len(test.Requested))
	for _, v := range test.Requested {
		requested = append(requested, &wtype.Liquid{Vol: v.RawValue(), Vunit: v.Unit().PrefixedSymbol()})
	}

	available := make(wtype.ComponentVector, 0, len(test.Requested))
	for w, wv := range test.AvailableByWell {
		loc := fmt.Sprintf("well_%d", w)
		for _, v := range wv {
			available = append(available, &wtype.Liquid{Vol: v.RawValue(), Vunit: v.Unit().PrefixedSymbol(), Loc: loc})
		}
	}

	if len(requested) != len(available) { //test is bad, explode
		t.Fatalf("bad test: len(requested) != len(available): %d != %d", len(requested), len(available))
	}

	got := distributeVolumes(requested, available)
	gotVols := make([]float64, 0, len(got))
	for _, v := range got {
		vol := wunit.NewVolume(v.Vol, v.Vunit)
		gotVols = append(gotVols, vol.MustInStringUnit("ul").RawValue())
	}

	totalAvailable := 0.0
	for _, w := range test.AvailableByWell {
		totalAvailable += w[0].MustInStringUnit("ul").RawValue()
	}
	totalGot := 0.0
	for _, v := range gotVols {
		totalGot += v
	}

	if ta, tg := fmt.Sprintf("%.4g", totalAvailable), fmt.Sprintf("%.4g", totalGot); ta != tg {
		t.Errorf("didn't return all available volume: available %s ul, got %s ul\n", ta, tg)
	}

	if !reflect.DeepEqual(gotVols, test.ExpectedUl) {
		t.Errorf("return didn't match expected:\ne: %v\ng: %v", test.ExpectedUl, gotVols)
	}
}

type DistributeVolumesTests []*DistributeVolumesTest

func (tests DistributeVolumesTests) Run(t *testing.T) {
	for _, test := range tests {
		t.Run(test.Name, test.Run)
	}
}

func TestDistrubuteVolumes(t *testing.T) {
	DistributeVolumesTests{
		{
			Name:            "all equal with excess",
			Requested:       []wunit.Volume{wunit.NewVolume(100, "ul"), wunit.NewVolume(100, "ul"), wunit.NewVolume(100, "ul"), wunit.NewVolume(100, "ul")},
			AvailableByWell: [][]wunit.Volume{{wunit.NewVolume(500, "ul"), wunit.NewVolume(500, "ul"), wunit.NewVolume(500, "ul"), wunit.NewVolume(500, "ul")}},
			ExpectedUl:      []float64{100 + 25, 100 + 25, 100 + 25, 100 + 25}, // value allocate by need + evenly distributed excess
		},
		{
			Name:            "all equal exact match",
			Requested:       []wunit.Volume{wunit.NewVolume(100, "ul"), wunit.NewVolume(100, "ul"), wunit.NewVolume(100, "ul"), wunit.NewVolume(100, "ul")},
			AvailableByWell: [][]wunit.Volume{{wunit.NewVolume(400, "ul"), wunit.NewVolume(400, "ul"), wunit.NewVolume(400, "ul"), wunit.NewVolume(400, "ul")}},
			ExpectedUl:      []float64{100, 100, 100, 100},
		},
		{
			Name:            "all equal with shortfall",
			Requested:       []wunit.Volume{wunit.NewVolume(100, "ul"), wunit.NewVolume(100, "ul"), wunit.NewVolume(100, "ul"), wunit.NewVolume(100, "ul")},
			AvailableByWell: [][]wunit.Volume{{wunit.NewVolume(200, "ul"), wunit.NewVolume(200, "ul"), wunit.NewVolume(200, "ul"), wunit.NewVolume(200, "ul")}},
			ExpectedUl:      []float64{50, 50, 50, 50},
		},
		{
			Name:            "mixed with excess",
			Requested:       []wunit.Volume{wunit.NewVolume(20, "ul"), wunit.NewVolume(100, "ul"), wunit.NewVolume(0, "ul"), wunit.NewVolume(70, "ul")},
			AvailableByWell: [][]wunit.Volume{{wunit.NewVolume(500, "ul"), wunit.NewVolume(500, "ul"), wunit.NewVolume(500, "ul"), wunit.NewVolume(500, "ul")}},
			ExpectedUl:      []float64{20 + 77.5, 100 + 77.5, 0 + 77.5, 70 + 77.5},
		},
		{
			Name:            "mixed exact match",
			Requested:       []wunit.Volume{wunit.NewVolume(20, "ul"), wunit.NewVolume(100, "ul"), wunit.NewVolume(0, "ul"), wunit.NewVolume(70, "ul")},
			AvailableByWell: [][]wunit.Volume{{wunit.NewVolume(190, "ul"), wunit.NewVolume(190, "ul"), wunit.NewVolume(190, "ul"), wunit.NewVolume(190, "ul")}},
			ExpectedUl:      []float64{20, 100, 0, 70},
		},
		{
			Name:            "mixed with shortfall",
			Requested:       []wunit.Volume{wunit.NewVolume(20, "ul"), wunit.NewVolume(100, "ul"), wunit.NewVolume(0, "ul"), wunit.NewVolume(70, "ul")},
			AvailableByWell: [][]wunit.Volume{{wunit.NewVolume(100, "ul"), wunit.NewVolume(100, "ul"), wunit.NewVolume(100, "ul"), wunit.NewVolume(100, "ul")}},
			ExpectedUl:      []float64{20, 40, 0, 40},
		},
		{
			Name:            "mixed with excess multiwell",
			Requested:       []wunit.Volume{wunit.NewVolume(20, "ul"), wunit.NewVolume(100, "ul"), wunit.NewVolume(0, "ul"), wunit.NewVolume(70, "ul")},
			AvailableByWell: [][]wunit.Volume{{wunit.NewVolume(500, "ul"), wunit.NewVolume(500, "ul")}, {wunit.NewVolume(500, "ul"), wunit.NewVolume(500, "ul")}},
			ExpectedUl:      []float64{20 + (500-120)/2.0, 100 + (500-120)/2.0, 0 + (500-70)/2.0, 70 + (500-70)/2.0},
		},
		{
			Name:            "mixed exact match multiwell",
			Requested:       []wunit.Volume{wunit.NewVolume(20, "ul"), wunit.NewVolume(100, "ul"), wunit.NewVolume(0, "ul"), wunit.NewVolume(70, "ul")},
			AvailableByWell: [][]wunit.Volume{{wunit.NewVolume(120, "ul"), wunit.NewVolume(120, "ul")}, {wunit.NewVolume(70, "ul"), wunit.NewVolume(70, "ul")}},
			ExpectedUl:      []float64{20, 100, 0, 70},
		},
		{
			Name:            "mixed with shortfall multiwell",
			Requested:       []wunit.Volume{wunit.NewVolume(20, "ul"), wunit.NewVolume(100, "ul"), wunit.NewVolume(0, "ul"), wunit.NewVolume(70, "ul")},
			AvailableByWell: [][]wunit.Volume{{wunit.NewVolume(50, "ul"), wunit.NewVolume(50, "ul")}, {wunit.NewVolume(50, "ul"), wunit.NewVolume(50, "ul")}},
			ExpectedUl:      []float64{20, 30, 0, 50},
		},
		{
			Name: "requested in litres",
			Requested: []wunit.Volume{
				wunit.NewVolume(2.1099999999999998e-05, "l"),
				wunit.NewVolume(3.1049999999999996e-05, "l"),
				wunit.NewVolume(2.1099999999999998e-05, "l"),
				wunit.NewVolume(2.14e-05, "l"),
				wunit.NewVolume(2.1e-05, "l"),
				wunit.NewVolume(2.1449999999999996e-05, "l"),
				wunit.NewVolume(2.4449999999999995e-05, "l"),
				wunit.NewVolume(2.1299999999999996e-05, "l"),
			},
			AvailableByWell: [][]wunit.Volume{
				{
					wunit.NewVolume(10000.0, "ul"),
					wunit.NewVolume(10000.0, "ul"),
					wunit.NewVolume(10000.0, "ul"),
					wunit.NewVolume(10000.0, "ul"),
					wunit.NewVolume(10000.0, "ul"),
					wunit.NewVolume(10000.0, "ul"),
					wunit.NewVolume(10000.0, "ul"),
					wunit.NewVolume(10000.0, "ul"),
				},
			},
			ExpectedUl: []float64{1248.24375, 1258.1937500000001, 1248.24375, 1248.5437500000003, 1248.1437500000002, 1248.5937500000002, 1251.5937500000002, 1248.4437500000001},
		},
	}.Run(t)
}
