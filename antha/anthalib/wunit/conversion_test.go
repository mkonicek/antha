// make_plate_library_test.go
package wunit

import (
	"fmt"
	"testing"
)

type concConversionTest struct {
	StockConc    Concentration
	TargetConc   Concentration
	TotalVolume  Volume
	VolumeNeeded Volume
}

var (
	// Some Concentrations
	nilConc Concentration
	uMPer0  = NewConcentration(0, "uM")

	x2   = NewConcentration(2, "X")
	x10  = NewConcentration(10, "X")
	x100 = NewConcentration(100, "X")

	kgPerL01 = NewConcentration(0.1, "kg/L")
	gPerL1   = NewConcentration(1, "g/L")
	gPerL10  = NewConcentration(10, "g/L")
	gPerL100 = NewConcentration(100, "g/L")

	mPerL01  = NewConcentration(0.1, "M/L")
	mMPerL1  = NewConcentration(1, "mM/L")
	uMPerL10 = NewConcentration(10, "uM")

	// Some volumes
	nilVol Volume
	l0     = NewVolume(0, "l")
	l01    = NewVolume(0.1, "l")
	ul1    = NewVolume(1, "ul")
	ul10   = NewVolume(10, "ul")
	ul100  = NewVolume(100, "ul")
	ml1    = NewVolume(1, "ml")
)

var tests1 = []concConversionTest{
	concConversionTest{StockConc: kgPerL01, TargetConc: gPerL1, TotalVolume: ul100, VolumeNeeded: ul1},
	concConversionTest{StockConc: gPerL100, TargetConc: gPerL1, TotalVolume: ul100, VolumeNeeded: ul1},
	concConversionTest{StockConc: x100, TargetConc: x2, TotalVolume: ul100, VolumeNeeded: NewVolume(2.0, "ul")},
	concConversionTest{StockConc: mMPerL1, TargetConc: uMPerL10, TotalVolume: ul100, VolumeNeeded: NewVolume(1.0, "ul")},
	concConversionTest{StockConc: mPerL01, TargetConc: mMPerL1, TotalVolume: ul100, VolumeNeeded: NewVolume(1.0, "ul")},
}

func TestVolumeForTargetConcentration(t *testing.T) {

	for _, test := range tests1 {

		vol, err := VolumeForTargetConcentration(test.TargetConc, test.StockConc, test.TotalVolume)

		if err != nil {
			t.Error(
				"for", test, "\n",
				"got error:", err.Error(), "\n",
			)
		}

		if !vol.EqualTo(test.VolumeNeeded) {
			t.Error(
				"for", fmt.Sprintf("%+v", test), "\n",
				"Expected Vol:", test.VolumeNeeded.ToString(), "\n",
				"Got Vol:", vol.ToString(), "\n",
			)
		}

	}
}
