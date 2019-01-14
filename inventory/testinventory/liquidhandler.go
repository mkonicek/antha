package testinventory

import (
	"fmt"

	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/microArch/driver/liquidhandling"
)

const (
	HVMinRate = 0.225
	HVMaxRate = 37.5
	LVMinRate = 0.0225
	LVMaxRate = 3.75
)

func MakeLHForTest(tipList []string) *liquidhandling.LHProperties { //nolint
	return makeLHForTest(tipList)
}

func MakeLHWithPlatesAndTipboxesForTest(inputPlateType string) *liquidhandling.LHProperties { //nolint
	ret, err := makeLHWithPlatesAndTipboxesForTest(inputPlateType)
	if err != nil {
		panic(err)
	}
	return ret
}

func MakeLHWithTipboxesForTest() *liquidhandling.LHProperties { //nolint
	if ret, err := makeLHWithTipboxesForTest(); err != nil {
		panic(err)
	} else {
		return ret
	}
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

func makeLHForTest(tipList []string) *liquidhandling.LHProperties {
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
	lhp := liquidhandling.NewLHProperties("Pipetmax", "Gilson", liquidhandling.LLLiquidHandler, liquidhandling.DisposableTips, layout)
	// get tips permissible from the factory
	SetUpTipsFor(lhp, tipList)

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
	ha.AddPosition(wtype.Coordinates{X: 0, Y: -18.08, Z: 0})
	ha.AddPosition(wtype.Coordinates{X: 0, Y: 0, Z: 0})
	ha.LoadHead(hvhead)
	ha.LoadHead(lvhead)
	lhp.Heads = append(lhp.Heads, hvhead, lvhead)
	lhp.Adaptors = append(lhp.Adaptors, hvadaptor, lvadaptor)
	lhp.HeadAssemblies = append(lhp.HeadAssemblies, ha)

	return lhp
}

func SetUpTipsFor(lhp *liquidhandling.LHProperties, tipList []string) *liquidhandling.LHProperties {
	inList := make(map[string]bool, len(tipList))
	for _, tipType := range tipList {
		inList[tipType] = true
	}

	seen := make(map[string]bool)

	tipboxes := make([]*wtype.LHTipbox, 0, len(inList))
	for _, tb := range GetTipboxes() {
		if tb.Mnfr == lhp.Mnfr || lhp.Mnfr == "MotherNature" {
			tip := tb.Tips[0][0]
			str := tip.Mnfr + tip.Type + tip.MinVol.ToString() + tip.MaxVol.ToString()

			if inList[tip.Type] && !seen[str] {
				seen[str] = true
				tipboxes = append(tipboxes, tb)
			}
		}
	}

	lhp.TipFactory = liquidhandling.NewTipFactory(tipboxes, []*wtype.LHTipwaste{makeGilsonTipWaste()})
	return lhp
}

func makeLHWithTipboxesForTest() (*liquidhandling.LHProperties, error) {
	params := makeLHForTest([]string{"Gilson20", "Gilson200"})

	if tw, err := params.TipFactory.NewTipwaste("Gilsontipwaste"); err != nil {
		return nil, err
	} else {
		params.AddTipWaste(tw)
	}

	if tb, err := params.TipFactory.NewTipbox("DL10 Tip Rack (PIPETMAX 8x20)"); err != nil {
		return nil, err
	} else {
		params.AddTipBox(tb)
	}

	if tb, err = params.TipFactory.NewTipbox("DF200 Tip Rack (PIPETMAX 8x200)"); err != nil {
		return nil, err
	} else {
		params.AddTipBox(tb)
	}

	return params, nil
}

func makeLHWithPlatesAndTipboxesForTest(inputPlateType string) (*liquidhandling.LHProperties, error) {
	params, err := makeLHWithTipboxesForTest()
	if err != nil {
		return nil, err
	}

	plates := getPlateByType()

	if inputPlate, err := makeTestInputPlate(plates, inputPlateType); err != nil {
		return nil, err
	} else if err := params.AddInputPlate(inputPlate); err != nil {
		return nil, err
	}

	if outputPlate, err := makeTestOutputPlate(plates, ctx); err != nil {
		return nil, err
	} else if err := params.AddOutputPlate(outputPlate); err != nil {
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
