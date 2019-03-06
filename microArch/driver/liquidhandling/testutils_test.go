package liquidhandling

import (
	"context"
	"fmt"

	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/inventory"
	"github.com/antha-lang/antha/inventory/testinventory"
)

const (
	HVMinRate = 0.225
	HVMaxRate = 37.5
	LVMinRate = 0.0225
	LVMaxRate = 3.75
)

func MakeGilsonForTest(tipList []string) *LHProperties { //nolint
	ctx := testinventory.NewContext(context.Background())
	return makeGilsonForTest(ctx, tipList)
}

func MakeGilsonWithPlatesAndTipboxesForTest(inputPlateType string) *LHProperties { //nolint
	ctx := testinventory.NewContext(context.Background())
	ret, err := makeGilsonWithPlatesAndTipboxesForTest(ctx, inputPlateType)
	if err != nil {
		panic(err)
	}
	return ret
}

func MakeGilsonWithTipboxesForTest() *LHProperties { //nolint
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

func makeGilsonForTest(ctx context.Context, tipList []string) *LHProperties {
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
			pos := wtype.NewLHPosition(fmt.Sprintf("position_%d", i+1), wtype.Coordinates3D{X: xp, Y: yp, Z: zp}, wtype.SBSFootprint)
			layout[pos.Name] = pos
			i += 1
			xp += xi
		}
		yp += yi
	}
	lhp := NewLHProperties("Pipetmax", "Gilson", LLLiquidHandler, DisposableTips, layout)
	// get tips permissible from the factory
	SetUpTipsFor(ctx, lhp, tipList)

	lhp.Preferences = &LayoutOpt{
		Tipboxes:  Addresses{"position_2", "position_3", "position_6", "position_9", "position_8", "position_5", "position_4", "position_7"},
		Inputs:    Addresses{"position_4", "position_5", "position_6", "position_9", "position_8", "position_3"},
		Outputs:   Addresses{"position_8", "position_9", "position_6", "position_5", "position_3", "position_1"},
		Washes:    Addresses{"position_8"},
		Tipwastes: Addresses{"position_1", "position_7"},
		Wastes:    Addresses{"position_9"},
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

func SetUpTipsFor(ctx context.Context, lhp *LHProperties, tipList []string) *LHProperties {
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

func makeGilsonWithTipboxesForTest(ctx context.Context) (*LHProperties, error) {
	params := makeGilsonForTest(ctx, defaultTipList())

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

func makeGilsonWithPlatesAndTipboxesForTest(ctx context.Context, inputPlateType string) (*LHProperties, error) {
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
