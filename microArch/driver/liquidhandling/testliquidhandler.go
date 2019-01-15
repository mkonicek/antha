package liquidhandling

import (
	"fmt"

	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/inventory/testinventory"
)

const (
	HVMinRate = 0.225
	HVMaxRate = 37.5
	LVMinRate = 0.0225
	LVMaxRate = 3.75
)

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

func MakeLHForTest(tipList []string) (*LHProperties, error) {
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
			layout[posname] = wtype.NewLHPosition(posname, wtype.Coordinates{X: xp, Y: yp, Z: zp})
			i += 1
			xp += xi
		}
		yp += yi
	}
	lhp := NewLHProperties("Pipetmax", "Gilson", LLLiquidHandler, DisposableTips, layout)
	// get tips permissible from the factory
	if _, err := SetUpTipsFor(lhp, tipList); err != nil {
		return nil, err
	}

	lhp.Preferences = &LayoutOpt{
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
	ha.AddPosition(wtype.Coordinates{X: 0, Y: -18.08, Z: 0})
	ha.AddPosition(wtype.Coordinates{X: 0, Y: 0, Z: 0})
	ha.LoadHead(hvhead)
	ha.LoadHead(lvhead)
	lhp.Heads = append(lhp.Heads, hvhead, lvhead)
	lhp.Adaptors = append(lhp.Adaptors, hvadaptor, lvadaptor)
	lhp.HeadAssemblies = append(lhp.HeadAssemblies, ha)

	return lhp, nil
}

func SetUpTipsFor(lhp *LHProperties, tipList []string) (*LHProperties, error) {
	allTipboxes := make([]*wtype.LHTipbox, 0, len(tipList))
	inv := testinventory.GetInventoryForTest()
	for _, tipName := range tipList {
		if tb, err := inv.NewTipbox(tipName); err != nil {
			return nil, err
		} else {
			allTipboxes = append(allTipboxes, tb)
		}
	}

	tipboxes := make([]*wtype.LHTipbox, 0, len(allTipboxes))
	for _, tb := range allTipboxes {
		if tb.Mnfr == lhp.Mnfr || lhp.Mnfr == "MotherNature" {
			tipboxes = append(tipboxes, tb)
		}
	}

	lhp.TipFactory = NewTipFactory(tipboxes, []*wtype.LHTipwaste{testinventory.MakeTestTipWaste()})
	return lhp, nil
}

func MakeLHWithTipboxesForTest() (*LHProperties, error) {
	params, err := MakeLHForTest([]string{"Gilson20", "Gilson200"})
	if err != nil {
		return nil, err
	}

	if tw, err := params.TipFactory.NewTipwaste("Gilsontipwaste"); err != nil {
		return nil, err
	} else {
		params.AddTipwaste(tw)
	}

	if tb, err := params.TipFactory.NewTipbox("DL10 Tip Rack (PIPETMAX 8x20)"); err != nil {
		return nil, err
	} else {
		params.AddTipBox(tb)
	}

	if tb, err := params.TipFactory.NewTipbox("D200 Tip Rack (PIPETMAX 8x200)"); err != nil {
		return nil, err
	} else {
		params.AddTipBox(tb)
	}

	return params, nil
}

func MakeLHWithPlatesAndTipboxesForTest(inputPlateType string) (*LHProperties, error) {
	params, err := MakeLHWithTipboxesForTest()
	if err != nil {
		return nil, err
	}

	if inputPlate, err := makeTestInputPlate(inputPlateType); err != nil {
		return nil, err
	} else if err := params.AddInputPlate(inputPlate); err != nil {
		return nil, err
	}

	if outputPlate, err := testinventory.GetInventoryForTest().NewPlate("DSW96"); err != nil {
		return nil, err
	} else if err := params.AddOutputPlate(outputPlate); err != nil {
		return nil, err
	}
	return params, nil
}

func makeTestInputPlate(inputPlateType string) (*wtype.Plate, error) {
	if inputPlateType == "" {
		inputPlateType = "DWST12"
	}

	if p, err := testinventory.GetInventoryForTest().NewPlate(inputPlateType); err != nil {
		return nil, err
	} else {
		c := wtype.NewLHComponent()
		c.CName = "water"
		c.Type = wtype.LTWater
		c.Smax = 9999
		c.Vol = 5000.0 // ul
		p.AddComponent(c, true)
		return p, nil
	}
}
