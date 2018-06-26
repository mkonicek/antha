package liquidhandling

import (
	"fmt"
	"testing"

	"github.com/antha-lang/antha/antha/anthalib/wunit"
)

func getVols() []wunit.Volume {
	// a selection of volumes
	vols := make([]wunit.Volume, 0, 1)
	for _, v := range []float64{0.5, 1.0, 2.0, 5.0, 10.0, 20.0, 30.0, 50.0, 100.0, 200.0} {
		vol := wunit.NewVolume(v, "ul")
		vols = append(vols, vol)
	}
	return vols
}

// answers to test

func getMinvols1() []wunit.Volume {
	v1 := wunit.NewVolume(0.5, "ul")
	v2 := wunit.NewVolume(10.0, "ul")

	ret := []wunit.Volume{v1, v1, v1, v1, v1, v1, v2, v2, v2, v2}

	return ret
}

func getMaxvols1() []wunit.Volume {
	v1 := wunit.NewVolume(20.0, "ul")
	v2 := wunit.NewVolume(250.0, "ul")

	ret := []wunit.Volume{v1, v1, v1, v1, v1, v1, v2, v2, v2, v2}

	return ret
}

/*

 */
func getTypes1() []string {
	ret := []string{"Gilson20", "Gilson20", "Gilson20", "Gilson20", "Gilson20", "Gilson20", "Gilson200", "Gilson200", "Gilson200", "Gilson200"}

	return ret
}

func TestDefaultChooser(t *testing.T) {
	vols := getVols()
	lhp := MakeGilsonForTest()
	minvols := getMinvols1()
	maxvols := getMaxvols1()
	types := getTypes1()

	for i, vol := range vols {
		prm, tip, err := ChooseChannel(vol, lhp)
		if err != nil {
			t.Error(err)
		}

		tiptype := ""

		if tip != nil {
			tiptype = tip.Type
		}

		mxr := maxvols[i]
		mnr := minvols[i]
		tpr := types[i]

		if prm == nil {
			if !mxr.IsZero() || !mnr.IsZero() || tpr != tiptype {
				t.Errorf(fmt.Sprint("Incorrect channel choice for volume ", vol.ToString(), " Got nil want: ", mnr.ToString(), " ", tpr))
			}

		} else if !prm.Maxvol.EqualTo(mxr) || !prm.Minvol.EqualTo(mnr) || tiptype != tpr {
			t.Errorf(fmt.Sprint("Incorrect channel choice for volume ", vol.ToString(), "\n\tGot ", prm.Minvol.ToString(), " ", prm.Maxvol.ToString(), " ", tiptype, " \n\tWANT: ", mnr.ToString(), " ", mxr.ToString(), " ", tpr))
		}
	}
}

func TestSmallVolumeError(t *testing.T) {
	lhp := MakeGilsonForTest()

	vol := wunit.NewVolume(0.47, "ul")

	prm, tip, err := ChooseChannel(vol, lhp)

	if prm != nil {
		t.Error("channel was not nil for small volume")
	}
	if tip != nil {
		t.Error("tip was not nil for small volume")
	}
	if err == nil {
		t.Error("error not generated for small volume")
	}
}
