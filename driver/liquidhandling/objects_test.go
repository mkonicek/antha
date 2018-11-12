package main

import (
	"context"
	"fmt"

	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/inventory"
	"github.com/antha-lang/antha/inventory/testinventory"
	"github.com/antha-lang/antha/microArch/driver/liquidhandling"
)

func makePlateForTest() *wtype.Plate {
	swshp := wtype.NewShape("box", "mm", 8.2, 8.2, 41.3)
	welltype := wtype.NewLHWell("ul", 200, 10, swshp, wtype.VWellBottom, 8.2, 8.2, 41.3, 4.7, "mm")
	p := wtype.NewLHPlate("DSW96", "none", 8, 12, wtype.Coordinates{127.76, 85.48, 43.1}, welltype, 9.0, 9.0, 0.5, 0.5, 0.5)
	return p
}

func makeTipForTest() *wtype.LHTip {
	shp := wtype.NewShape("cylinder", "mm", 7.3, 7.3, 51.2)
	return wtype.NewLHTip("me", "mytype", 0.5, 1000.0, "ul", false, shp, 44.7)
}

func makeTipboxForTest() *wtype.LHTipbox {
	shp := wtype.NewShape("cylinder", "mm", 7.3, 7.3, 51.2)
	w := wtype.NewLHWell("ul", 250.0, 10.0, shp, wtype.FlatWellBottom, 7.3, 7.3, 51.2, 0.0, "mm")
	tiptype := makeTipForTest()
	tb := wtype.NewLHTipbox(8, 12, wtype.Coordinates{127.76, 85.48, 120.0}, "me", "mytype", tiptype, w, 9.0, 9.0, 0.5, 0.5, 0.0)
	return tb
}

func makeTipwasteForTest() *wtype.LHTipwaste {
	shp := wtype.NewShape("box", "mm", 123.0, 80.0, 92.0)
	w := wtype.NewLHWell("ul", 800000.0, 800000.0, shp, 0, 123.0, 80.0, 92.0, 0.0, "mm")
	lht := wtype.NewLHTipwaste(6000, "TipwasteForTest", "ACME Corp.", wtype.Coordinates{X: 127.76, Y: 85.48, Z: 92.0}, w, 49.5, 31.5, 0.0)
	return lht
}

const (
	HVMinRate = 0.225
	HVMaxRate = 37.5
	LVMinRate = 0.0225
	LVMaxRate = 3.75
)

func MakeGilsonForTest(tipList []string) *liquidhandling.LHProperties { //nolint
	ctx := testinventory.NewContext(context.Background())
	return makeGilsonForTest(ctx, tipList)
}

func MakeGilsonWithPlatesAndTipboxesForTest(inputPlateType string) *liquidhandling.LHProperties { //nolint
	ctx := testinventory.NewContext(context.Background())
	ret, err := makeGilsonWithPlatesAndTipboxesForTest(ctx, inputPlateType)
	if err != nil {
		panic(err)
	}
	return ret
}

func MakeGilsonWithTipboxesForTest() *liquidhandling.LHProperties { //nolint
	ctx := testinventory.NewContext(context.Background())
	ret, err := makeGilsonWithTipboxesForTest(ctx)
	if err != nil {
		panic(err)
	}
	return ret
}

func getHVConfig() *wtype.LHChannelParameter {
	minvol := wunit.NewVolume(20, "ul")
	maxvol := wunit.NewVolume(200, "ul")
	minspd := wunit.NewFlowRate(HVMinRate, "ml/min")
	maxspd := wunit.NewFlowRate(HVMaxRate, "ml/min")

	return wtype.NewLHChannelParameter("HVconfig", "GilsonPipetmax", minvol, maxvol, minspd, maxspd, 8, false, wtype.LHVChannel, 0)
}

func getLVConfig() *wtype.LHChannelParameter {
	newminvol := wunit.NewVolume(0.5, "ul")
	newmaxvol := wunit.NewVolume(20, "ul")
	newminspd := wunit.NewFlowRate(LVMinRate, "ml/min")
	newmaxspd := wunit.NewFlowRate(LVMaxRate, "ml/min")

	return wtype.NewLHChannelParameter("LVconfig", "GilsonPipetmax", newminvol, newmaxvol, newminspd, newmaxspd, 8, false, wtype.LHVChannel, 1)
}

func makeGilsonForTest(ctx context.Context, tipList []string) *liquidhandling.LHProperties {
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
	lhp := liquidhandling.NewLHProperties(9, "Pipetmax", "Gilson", liquidhandling.LLLiquidHandler, liquidhandling.DisposableTips, layout)
	// get tips permissible from the factory
	SetUpTipsFor(ctx, lhp, tipList)

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

	hvconfig := getHVConfig()
	hvadaptor := wtype.NewLHAdaptor("DummyAdaptor", "Gilson", hvconfig)
	hvhead := wtype.NewLHHead("HVHead", "Gilson", hvconfig)
	hvhead.Adaptor = hvadaptor

	lvconfig := getLVConfig()
	lvadaptor := wtype.NewLHAdaptor("DummyAdaptor", "Gilson", lvconfig)
	lvhead := wtype.NewLHHead("LVHead", "Gilson", lvconfig)
	lvhead.Adaptor = lvadaptor

	ha := wtype.NewLHHeadAssembly(nil)
	ha.AddPosition(wtype.Coordinates{X: 0, Y: -18.08, Z: 0})
	ha.AddPosition(wtype.Coordinates{X: 0, Y: 0, Z: 0})
	ha.LoadHead(hvhead)
	ha.LoadHead(lvhead)
	lhp.Heads = append(lhp.Heads, hvhead, lvhead)
	lhp.Adaptors = append(lhp.Adaptors, hvadaptor, lvadaptor)
	lhp.HeadAssemblies = append(lhp.HeadAssemblies, ha)

	return lhp
}

func SetUpTipsFor(ctx context.Context, lhp *liquidhandling.LHProperties, tipList []string) *liquidhandling.LHProperties {
	inList := func(s string, sa []string) bool {
		for _, ss := range sa {
			if s == ss {
				return true
			}
		}
		return false
	}

	seen := make(map[string]bool)

	for _, tb := range testinventory.GetTipboxes(ctx) {
		if tb.Mnfr == lhp.Mnfr || lhp.Mnfr == "MotherNature" {
			//ignore filter tips and the hacky "low volume high volume" ones
			//		if tb.Tiptype.Filtered || tb.Tiptype.Type == "LVGilson200" {
			//			continue
			//		}

			// ignore tips not in the list

			if !inList(tb.Tiptype.Type, tipList) {
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

func makeGilsonWithTipboxesForTest(ctx context.Context) (*liquidhandling.LHProperties, error) {
	params := makeGilsonForTest(ctx, []string{"Gilson20", "Gilson200"})

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

func makeGilsonWithPlatesAndTipboxesForTest(ctx context.Context, inputPlateType string) (*liquidhandling.LHProperties, error) {
	params, err := makeGilsonWithTipboxesForTest(ctx)
	if err != nil {
		return nil, err
	}

	inputPlate, err := makeTestInputPlate(ctx, inputPlateType)

	if err != nil {
		return nil, err
	}

	err = params.AddInputPlate(inputPlate)

	if err != nil {
		return nil, err
	}

	outputPlate, err := makeTestOutputPlate(ctx)

	if err != nil {
		return nil, err
	}

	err = params.AddOutputPlate(outputPlate)

	if err != nil {
		return nil, err
	}
	return params, nil
}

func makeTestInputPlate(ctx context.Context, inputPlateType string) (*wtype.Plate, error) {
	if inputPlateType == "" {
		inputPlateType = "DWST12"
	}

	p, err := inventory.NewPlate(ctx, inputPlateType)

	if err != nil {
		return nil, err
	}

	c, err := inventory.NewComponent(ctx, "water")

	if err != nil {
		return nil, err
	}

	c.Vol = 5000.0 // ul

	p.AddComponent(c, true)

	return p, nil
}

func makeTestOutputPlate(ctx context.Context) (*wtype.Plate, error) {
	p, err := inventory.NewPlate(ctx, "DSW96")

	if err != nil {
		return nil, err
	}

	return p, nil
}
