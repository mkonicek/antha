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
	swshp := wtype.NewShape(wtype.BoxShape, "mm", 8.2, 8.2, 41.3)
	welltype := wtype.NewLHWell("ul", 200, 10, swshp, wtype.VWellBottom, 8.2, 8.2, 41.3, 4.7, "mm")
	p := wtype.NewLHPlate("DSW96", "none", 8, 12, wtype.Coordinates3D{X: 127.76, Y: 85.48, Z: 43.1}, welltype, 9.0, 9.0, 0.5, 0.5, 0.5)
	return p
}

func makeTipForTest() *wtype.LHTip {
	shp := wtype.NewShape(wtype.CylinderShape, "mm", 7.3, 7.3, 51.2)
	return wtype.NewLHTip("me", "mytype", 0.5, 1000.0, "ul", false, shp, 44.7)
}

func makeTipboxForTest() *wtype.LHTipbox {
	shp := wtype.NewShape(wtype.CylinderShape, "mm", 7.3, 7.3, 51.2)
	w := wtype.NewLHWell("ul", 250.0, 10.0, shp, wtype.FlatWellBottom, 7.3, 7.3, 51.2, 0.0, "mm")
	tiptype := makeTipForTest()
	tb := wtype.NewLHTipbox(8, 12, wtype.Coordinates3D{X: 127.76, Y: 85.48, Z: 120.0}, "me", "mytype", tiptype, w, 9.0, 9.0, 0.5, 0.5, 0.0)
	return tb
}

func makeTipwasteForTest() *wtype.LHTipwaste {
	shp := wtype.NewShape(wtype.BoxShape, "mm", 123.0, 80.0, 92.0)
	w := wtype.NewLHWell("ul", 800000.0, 800000.0, shp, 0, 123.0, 80.0, 92.0, 0.0, "mm")
	lht := wtype.NewLHTipwaste(6000, "TipwasteForTest", "ACME Corp.", wtype.Coordinates3D{X: 127.76, Y: 85.48, Z: 92.0}, w, 49.5, 31.5, 0.0)
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

	layout := make(map[string]*wtype.LHPosition)
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
			layout[posname] = wtype.NewLHPosition(posname, wtype.Coordinates3D{X: xp, Y: yp, Z: zp}, wtype.SBSFootprint)
			i += 1
			xp += xi
		}
		yp += yi
	}
	lhp := liquidhandling.NewLHProperties("Pipetmax", "Gilson", liquidhandling.LLLiquidHandler, liquidhandling.DisposableTips, layout)
	// get tips permissible from the factory
	SetUpTipsFor(ctx, lhp, tipList)

	lhp.Preferences = &liquidhandling.LayoutOpt{
		Tipboxes:  []string{"position_2", "position_3", "position_6", "position_9", "position_8", "position_5", "position_4", "position_7"},
		Inputs:    []string{"position_4", "position_5", "position_6", "position_9", "position_8", "position_3"},
		Outputs:   []string{"position_8", "position_9", "position_6", "position_5", "position_3", "position_1"},
		Washes:    []string{"position_8"},
		Tipwastes: []string{"position_1", "position_7"},
		Wastes:    []string{"position_9"},
	}

	hvconfig := getHVConfig()
	hvadaptor := wtype.NewLHAdaptor("DummyAdaptor", "Gilson", hvconfig)
	hvhead := wtype.NewLHHead("HVHead", "Gilson", hvconfig)
	hvhead.Adaptor = hvadaptor

	lvconfig := getLVConfig()
	lvadaptor := wtype.NewLHAdaptor("DummyAdaptor", "Gilson", lvconfig)
	lvhead := wtype.NewLHHead("LVHead", "Gilson", lvconfig)
	lvhead.Adaptor = lvadaptor

	ha := wtype.NewLHHeadAssembly(nil)
	ha.AddPosition(wtype.Coordinates3D{X: 0, Y: -18.08, Z: 0})
	ha.AddPosition(wtype.Coordinates3D{X: 0, Y: 0, Z: 0})
	if err := ha.LoadHead(hvhead); err != nil {
		panic(err)
	}
	if err := ha.LoadHead(lvhead); err != nil {
		panic(err)
	}
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

	if tw, err := inventory.NewTipwaste(ctx, "Gilsontipwaste"); err != nil {
		return nil, err
	} else if err := params.AddTipWaste(tw); err != nil {
		return nil, err
	}

	if tb, err := inventory.NewTipbox(ctx, "DL10 Tip Rack (PIPETMAX 8x20)"); err != nil {
		return nil, err
	} else if err := params.AddTipBox(tb); err != nil {
		return nil, err
	}

	if tb, err := inventory.NewTipbox(ctx, "DF200 Tip Rack (PIPETMAX 8x200)"); err != nil {
		return nil, err
	} else if err := params.AddTipBox(tb); err != nil {
		return nil, err
	}

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

	if _, err := p.AddComponent(c, true); err != nil {
		return nil, err
	}

	return p, nil
}

func makeTestOutputPlate(ctx context.Context) (*wtype.Plate, error) {
	p, err := inventory.NewPlate(ctx, "DSW96")

	if err != nil {
		return nil, err
	}

	return p, nil
}
