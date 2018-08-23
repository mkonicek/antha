package wunit

import (
	"math"
	"testing"
)

func TestMetersCubed(t *testing.T) {
	reg := makeGlobalUnitRegistry()

	//test that we support m^3 and correct prefix exponent
	if volume, err := reg.NewMeasurement(1.0, "mm^3"); err != nil {
		t.Error(err)
	} else if ul, err := reg.GetUnit("ul"); err != nil {
		t.Error(err)
	} else if g, e := volume.ConvertTo(ul), 1.0; math.Abs(g-e) > 1.0e-6 {
		t.Errorf("converting %v to ul: expected %g, got %g", volume, e, g)
	}
}
