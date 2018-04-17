package liquidhandling

import (
	"context"
	"fmt"
	"testing"

	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/inventory/testinventory"
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
	lhp := makeTestLH()
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
				t.Fatal(fmt.Sprint("Incorrect channel choice for volume ", vol.ToString(), " Got nil want: ", mnr.ToString(), " ", tpr))
			}

		} else if !prm.Maxvol.EqualTo(mxr) || !prm.Minvol.EqualTo(mnr) || tiptype != tpr {
			t.Fatal(fmt.Sprint("Incorrect channel choice for volume ", vol.ToString(), "\n\tGot ", prm.Minvol.ToString(), " ", prm.Maxvol.ToString(), " ", tiptype, " \n\tWANT: ", mnr.ToString(), " ", mxr.ToString(), " ", tpr))
		}
		/*
			if prm == nil {
				fmt.Println("V: ", vol.ToString(), " NO SOLUTION")
			} else {
				fmt.Println("V: ", vol.ToString(), " Mn: ", prm.Minvol.ToString(), " Mx: ", prm.Maxvol.ToString(), " TIP: ", tiptype)
			}
		*/

	}
}

func TestSmallVolumeError(t *testing.T) {
	lhp := makeTestLH()

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

func makeTestLH() *LHProperties {
	// gilson pipetmax

	layout := make(map[string]wtype.Coordinates)
	i := 0
	x0 := 3.886
	y0 := 3.513
	z0 := -82.035
	xi := 149.86
	yi := 95.25
	xp := x0
	yp := y0
	zp := z0
	for y := 0; y < 3; y++ {
		for x := 0; x < 3; x++ {
			posname := fmt.Sprintf("position_%d", i+1)
			crds := wtype.Coordinates{X: xp, Y: yp, Z: zp}
			layout[posname] = crds
			i += 1
			xp += xi
		}
		yp += yi
	}
	lhp := NewLHProperties(9, "Pipetmax", "Gilson", LLLiquidHandler, DisposableTips, layout)
	// get tips permissible from the factory
	SetUpTipsFor(lhp)

	//lhp.Tip_preferences = []string{"position_2", "position_3", "position_6", "position_9", "position_8", "position_5", "position_4", "position_7"}
	//lhp.Tip_preferences = []string{"position_2", "position_3", "position_6", "position_9", "position_8", "position_5", "position_7"}

	lhp.Tip_preferences = []string{"position_9", "position_6", "position_3", "position_5", "position_2"} //jmanart i cut it down to 5, as it was hardcoded in the liquidhandler getInputs call before

	// original preferences
	lhp.Input_preferences = []string{"position_4", "position_5", "position_6", "position_9", "position_8", "position_3"}
	lhp.Output_preferences = []string{"position_7", "position_8", "position_9", "position_6", "position_5", "position_3"}

	// use these new preferences for gel loading: this is needed because outplate overlaps inplate otherwise so move inplate to position 5 rather than 4 (pos 4 deleted)
	//lhp.Input_preferences = []string{"position_5", "position_6", "position_9", "position_8", "position_3"}
	//lhp.Output_preferences = []string{"position_9", "position_8", "position_7", "position_6", "position_5", "position_3"}

	lhp.Wash_preferences = []string{"position_8"}
	lhp.Tipwaste_preferences = []string{"position_1"}
	lhp.Waste_preferences = []string{"position_9"}
	//	lhp.Tip_preferences = []int{2, 3, 6, 9, 5, 8, 4, 7}
	//	lhp.Input_preferences = []int{24, 25, 26, 29, 28, 23}
	//	lhp.Output_preferences = []int{10, 11, 12, 13, 14, 15}
	minvol := wunit.NewVolume(10, "ul")
	maxvol := wunit.NewVolume(250, "ul")
	minspd := wunit.NewFlowRate(0.5, "ml/min")
	maxspd := wunit.NewFlowRate(2, "ml/min")

	hvconfig := wtype.NewLHChannelParameter("HVconfig", "GilsonPipetmax", minvol, maxvol, minspd, maxspd, 8, false, wtype.LHVChannel, 0)
	hvadaptor := wtype.NewLHAdaptor("DummyAdaptor", "Gilson", hvconfig)
	hvhead := wtype.NewLHHead("HVHead", "Gilson", hvconfig)
	hvhead.Adaptor = hvadaptor
	newminvol := wunit.NewVolume(0.5, "ul")
	newmaxvol := wunit.NewVolume(20, "ul")
	newminspd := wunit.NewFlowRate(0.1, "ml/min")
	newmaxspd := wunit.NewFlowRate(0.5, "ml/min")

	lvconfig := wtype.NewLHChannelParameter("LVconfig", "GilsonPipetmax", newminvol, newmaxvol, newminspd, newmaxspd, 8, false, wtype.LHVChannel, 1)
	lvadaptor := wtype.NewLHAdaptor("DummyAdaptor", "Gilson", lvconfig)
	lvhead := wtype.NewLHHead("LVHead", "Gilson", lvconfig)
	lvhead.Adaptor = lvadaptor

	lhp.Heads = append(lhp.Heads, hvhead)
	lhp.Heads = append(lhp.Heads, lvhead)
	lhp.HeadsLoaded = append(lhp.HeadsLoaded, hvhead)
	lhp.HeadsLoaded = append(lhp.HeadsLoaded, lvhead)

	return lhp
}

func SetUpTipsFor(lhp *LHProperties) *LHProperties {
	ctx := testinventory.NewContext(context.Background())

	seen := make(map[string]bool)

	for _, tb := range testinventory.GetTipboxes(ctx) {
		if tb.Mnfr == lhp.Mnfr || lhp.Mnfr == "MotherNature" {
			//ignore filter tips
			if tb.Tiptype.Filtered {
				continue
			}
			tip := tb.Tips[0][0]
			str := tip.Mnfr + tip.Type + tip.MinVol.ToString() + tip.MaxVol.ToString()
			if seen[str] {
				continue
			}

			seen[str] = true
			lhp.Tips = append(lhp.Tips, tb.Tips[0][0])
		}
	}
	return lhp
}
