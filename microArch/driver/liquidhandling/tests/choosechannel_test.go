package tests

import (
	"errors"
	"fmt"
	"testing"

	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/laboratory"
	"github.com/antha-lang/antha/laboratory/testlab"
	"github.com/antha-lang/antha/microArch/driver/liquidhandling"
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
	v2 := wunit.NewVolume(20.0, "ul")

	ret := []wunit.Volume{v1, v1, v1, v1, v1, v1, v2, v2, v2, v2}

	return ret
}

func getMaxvols1() []wunit.Volume {
	v1 := wunit.NewVolume(20.0, "ul")
	v2 := wunit.NewVolume(200.0, "ul")

	ret := []wunit.Volume{v1, v1, v1, v1, v1, v1, v2, v2, v2, v2}

	return ret
}

/*

 */
func getTypes1() []string {
	ret := []string{"Gilson20", "Gilson20", "Gilson20", "Gilson20", "Gilson20", "Gilson20", "Gilson200", "Gilson200", "Gilson200", "Gilson200"}

	return ret
}

func getVols2() []wunit.Volume {
	// a selection of volumes
	vols := make([]wunit.Volume, 0, 1)
	for _, v := range []float64{1.0, 2.0, 5.0, 10.0, 20.0, 30.0, 50.0, 100.0, 200.0} {
		vol := wunit.NewVolume(v, "ul")
		vols = append(vols, vol)
	}
	return vols
}

// answers to test

func getMinvols2() []wunit.Volume {
	v1 := wunit.NewVolume(0.5, "ul")

	ret := []wunit.Volume{v1, v1, v1, v1, v1, v1, v1, v1, v1}

	return ret
}

func getMaxvols2() []wunit.Volume {
	v1 := wunit.NewVolume(20.0, "ul")

	ret := []wunit.Volume{v1, v1, v1, v1, v1, v1, v1, v1, v1}

	return ret
}

/*

 */
func getTypes2() []string {
	ret := []string{"LVGilson200", "LVGilson200", "LVGilson200", "LVGilson200", "LVGilson200", "LVGilson200", "LVGilson200", "LVGilson200", "LVGilson200"}

	return ret
}

func defaultTipList() []string {
	return []string{"Gilson20", "Gilson200"}
}

func TestDefaultChooser(t *testing.T) {
	testlab.WithTestLab(t, "", &testlab.TestElementCallbacks{
		Steps: func(lab *laboratory.Laboratory) error {
			vols := getVols()
			lhp := MakeGilsonForTest(lab, defaultTipList())
			minvols := getMinvols1()
			maxvols := getMaxvols1()
			types := getTypes1()

			for i, vol := range vols {
				prm, tip, err := liquidhandling.ChooseChannel(vol, lhp)
				if err != nil {
					return err
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
						return fmt.Errorf("Incorrect channel choice for volume %v\n\tGot nil want: %v %v ", vol.ToString(), mnr.ToString(), tpr)
					}

				} else if !prm.Maxvol.EqualTo(mxr) || !prm.Minvol.EqualTo(mnr) || tiptype != tpr {
					return fmt.Errorf("Incorrect channel choice for volume %v\n\tGot %v %v %v\n\tWant %v %v %v",
						vol.ToString(), prm.Minvol.ToString(), prm.Maxvol.ToString(), tiptype, mnr.ToString(), mxr.ToString(), tpr)
				}
			}
			return nil
		},
	})
}

func TestHVHVHVLVChooser(t *testing.T) {
	testlab.WithTestLab(t, "", &testlab.TestElementCallbacks{
		Steps: func(lab *laboratory.Laboratory) error {
			vols := getVols2()
			lhp := MakeGilsonForTest(lab, []string{"LVGilson200"})
			minvols := getMinvols2()
			maxvols := getMaxvols2()
			types := getTypes2()

			for i, vol := range vols {
				prm, tip, err := liquidhandling.ChooseChannel(vol, lhp)
				if err != nil {
					return err
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
						return fmt.Errorf("Incorrect channel choice for volume %v\n\tGot nit want: %v %v", vol.ToString(), mnr.ToString(), tpr)
					}

				} else if !prm.Maxvol.EqualTo(mxr) || !prm.Minvol.EqualTo(mnr) || tiptype != tpr {
					return fmt.Errorf("Incorrect channel choice for volume %v\n\tGot %v %v %v\n\tWant %v %v %v",
						vol.ToString(), prm.Minvol.ToString(), prm.Maxvol.ToString(), tiptype, mnr.ToString(), mxr.ToString(), tpr)
				}
			}
			return nil
		},
	})
}

func TestSmallVolumeError(t *testing.T) {
	testlab.WithTestLab(t, "", &testlab.TestElementCallbacks{
		Steps: func(lab *laboratory.Laboratory) error {
			lhp := MakeGilsonForTest(lab, defaultTipList())

			vol := wunit.NewVolume(0.47, "ul")

			prm, tip, err := liquidhandling.ChooseChannel(vol, lhp)

			if prm != nil {
				return errors.New("channel was not nil for small volume")
			}
			if tip != nil {
				return errors.New("tip was not nil for small volume")
			}
			if err == nil {
				return errors.New("error not generated for small volume")
			}
			return nil
		},
	})
}
