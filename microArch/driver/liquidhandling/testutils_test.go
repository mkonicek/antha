package liquidhandling

import (
	"context"
	"fmt"

	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/inventory"
	"github.com/antha-lang/antha/inventory/testinventory"
)

func MakeGilsonForTest() *LHProperties {
	// gilson pipetmax

	layout := make(map[string]wtype.Coordinates)
	i := 0
	x0 := 3.886
	y0 := 3.513
	z0 := -82.035
	xi := 149.86
	yi := 95.25
	xp := x0 // nolint
	yp := y0
	zp := z0
	for y := 0; y < 3; y++ {
		xp = x0
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

	lhp.Tip_preferences = []string{"position_2", "position_3", "position_6", "position_9", "position_8", "position_5", "position_4", "position_7"}
	//lhp.Tip_preferences = []string{"position_2", "position_3", "position_6", "position_9", "position_8", "position_5", "position_7"}

	//lhp.Tip_preferences = []string{"position_9", "position_6", "position_3", "position_5", "position_2"} //jmanart i cut it down to 5, as it was hardcoded in the liquidhandler getInputs call before

	// original preferences
	lhp.Input_preferences = []string{"position_4", "position_5", "position_6", "position_9", "position_8", "position_3"}
	lhp.Output_preferences = []string{"position_8", "position_9", "position_6", "position_5", "position_3", "position_1"}

	// use these new preferences for gel loading: this is needed because outplate overlaps inplate otherwise so move inplate to position 5 rather than 4 (pos 4 deleted)
	//lhp.Input_preferences = []string{"position_5", "position_6", "position_9", "position_8", "position_3"}
	//lhp.Output_preferences = []string{"position_9", "position_8", "position_7", "position_6", "position_5", "position_3"}

	lhp.Wash_preferences = []string{"position_8"}
	lhp.Tipwaste_preferences = []string{"position_1", "position_7"}
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

	ha := wtype.NewLHHeadAssembly(nil)
	ha.AddPosition(wtype.Coordinates{0, -18.08, 0})
	ha.AddPosition(wtype.Coordinates{0, 0, 0})
	ha.LoadHead(hvhead)
	ha.LoadHead(lvhead)
	lhp.HeadAssemblies = append(lhp.HeadAssemblies, ha)

	return lhp
}

func SetUpTipsFor(lhp *LHProperties) *LHProperties {
	ctx := testinventory.NewContext(context.Background())

	seen := make(map[string]bool)

	for _, tb := range testinventory.GetTipboxes(ctx) {
		if tb.Mnfr == lhp.Mnfr || lhp.Mnfr == "MotherNature" {
			//ignore filter tips and the hacky "low volume high volume" ones
			if tb.Tiptype.Filtered || tb.Tiptype.Type == "LVGilson200" {
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

func MakeGilsonWithPlatesForTest(ctx context.Context) (*LHProperties, error) {
	params := MakeGilsonForTest()

	tw, err := inventory.NewTipwaste(ctx, "Gilsontipwaste")
	if err != nil {
		return nil, err
	}
	params.AddTipWaste(tw)

	tb, err := inventory.NewTipbox(ctx, "DL10 Tip Rack (PIPETMAX 8x20)")
	if err != nil {
		return nil, err
	}
	params.AddTipBox(tb)

	tb, err = inventory.NewTipbox(ctx, "DF200 Tip Rack (PIPETMAX 8x200)")
	if err != nil {
		return nil, err
	}
	params.AddTipBox(tb)

	return params, nil
}
