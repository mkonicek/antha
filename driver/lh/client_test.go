package lh

import (
	"context"
	"fmt"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/inventory"
	"github.com/antha-lang/antha/inventory/testinventory"
	"github.com/antha-lang/antha/microArch/driver/liquidhandling"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestPlateSerializeDeserialize(t *testing.T) {
	p := makeplatefortest()

	validatePlate(t, p)

	encodedP := Encodeinterface(p)

	decodedP, err := DecodeGenericPlate(string(encodedP.Arg_1))

	if err != nil {
		t.Errorf(err.Error())
	}

	plate, ok := decodedP.(*wtype.LHPlate)

	if !ok {
		t.Errorf("WANT *wtype.LHPlate got %T", decodedP)
	}

	validatePlate(t, plate)

}

func TestTipboxSerializeDeserialize(t *testing.T) {
	tb := makeTipboxForTest()

	//validateTipbox(t, tb)

	encodedTB := Encodeinterface(tb)

	decodedTB, err := DecodeGenericPlate(encodedTB.GetArg_1())
	if err != nil {
		t.Error(err)
	}

	tb2, ok := decodedTB.(*wtype.LHTipbox)
	if !ok {
		t.Errorf("Expected *wtype.LHTipbox, got %T", decodedTB)
	}

	assertTipBoxesEqual(t, tb, tb2)

}

func TestLHPropertiesSerialiseDeserialise(t *testing.T) {
	ctx := testinventory.NewContext(context.Background())
	lhp, err := MakeGilsonWithPlatesForTest(ctx)
	if err != nil {
		t.Fatal(err)
	}

	s := EncodePtrToLHProperties(lhp)

	dec := DecodePtrToLHProperties(s)

	AssertLHPropertiesEqual(t, lhp, dec, "LHProperties")
}

func makeplatefortest() *wtype.LHPlate {
	swshp := wtype.NewShape("box", "mm", 8.2, 8.2, 41.3)
	welltype := wtype.NewLHWell("ul", 200, 10, swshp, wtype.VWellBottom, 8.2, 8.2, 41.3, 4.7, "mm")
	p := wtype.NewLHPlate("DSW96", "none", 8, 12, wtype.Coordinates{127.76, 85.48, 43.1}, welltype, 0.5, 0.5, 0.5, 0.5, 0.5)
	p.Welltype.SetWellTargets("Magick", []wtype.Coordinates{{0.0, -10.0, 0.0}, {0.0, 10.0, 0.0}})
	return p
}

func makeTipboxForTest() *wtype.LHTipbox {

	shp := wtype.NewShape("cylinder", "mm", 7.3, 7.3, 51.2)
	w := wtype.NewLHWell("ul", 20.0, 1.0, shp, 0, 7.3, 7.3, 46.0, 0.0, "mm")
	w.Extra["InnerL"] = 5.5
	w.Extra["InnerW"] = 5.5
	w.Extra["Tipeffectiveheight"] = 34.6
	tip := wtype.NewLHTip("gilson", "Gilson20", 0.5, 20.0, "ul", false, shp)
	tb := wtype.NewLHTipbox(8, 12, wtype.Coordinates{127.76, 85.48, 60.13}, "Gilson", "DL10 Tip Rack (PIPETMAX 8x20)", tip, w, 9.0, 9.0, 0.0, 0.0, 28.93)

	return tb

}

func MakeGilsonForTest() *liquidhandling.LHProperties {
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
	lhp.Heads = append(lhp.Heads, hvhead, lvhead)
	lhp.Adaptors = append(lhp.Adaptors, hvadaptor, lvadaptor)
	lhp.HeadAssemblies = append(lhp.HeadAssemblies, ha)

	return lhp
}

func SetUpTipsFor(lhp *liquidhandling.LHProperties) *liquidhandling.LHProperties {
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

func MakeGilsonWithPlatesForTest(ctx context.Context) (*liquidhandling.LHProperties, error) {
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

func validatePlate(t *testing.T, plate *wtype.LHPlate) {
	assertWellsEqual := func(what string, as, bs []*wtype.LHWell) {
		seen := make(map[*wtype.LHWell]int)
		for _, w := range as {
			seen[w] += 1
		}
		for _, w := range bs {
			seen[w] += 1
		}
		for w, count := range seen {
			if count != 2 {
				t.Errorf("%s: no matching well found (%d != %d) for %p %s:%s", what, count, 2, w, w.ID, w.Crds.FormatA1())
			}
		}
	}

	var ws1, ws2, ws3, ws4 []*wtype.LHWell

	for _, w := range plate.HWells {
		ws1 = append(ws1, w)
	}
	for crds, w := range plate.Wellcoords {
		ws2 = append(ws2, w)

		if w.Crds.FormatA1() != crds {
			t.Fatal(fmt.Sprintf("ERROR: Well coords not consistent -- %s != %s", w.Crds.FormatA1(), crds))
		}

		if w.WContents.Loc == "" {
			t.Fatal(fmt.Sprintf("ERROR: Well contents do not have loc set"))
		}

		ltx := strings.Split(w.WContents.Loc, ":")

		if ltx[0] != plate.ID {
			t.Fatal(fmt.Sprintf("ERROR: Plate ID for component not consistent -- %s != %s", ltx[0], plate.ID))
		}

		if w.Plate != nil && ltx[0] != w.Plate.(*wtype.LHPlate).ID {
			t.Fatal(fmt.Sprintf("ERROR: Plate ID for component not consistent with well -- %s != %s", ltx[0], w.Plate.(*wtype.LHPlate).ID))
		}

		if ltx[1] != crds {
			t.Fatal(fmt.Sprintf("ERROR: Coords for component not consistent: -- %s != %s", ltx[1], crds))
		}

	}

	for _, ws := range plate.Rows {
		ws3 = append(ws3, ws...)
	}
	for _, ws := range plate.Cols {
		ws4 = append(ws4, ws...)
	}
	assertWellsEqual("HWells != Rows", ws1, ws2)
	assertWellsEqual("Rows != Cols", ws2, ws3)
	assertWellsEqual("Cols != Wellcoords", ws3, ws4)

	// Check pointer-ID equality
	comp := make(map[string]*wtype.LHComponent)
	for _, w := range append(append(ws1, ws2...), ws3...) {
		c := w.WContents
		if c == nil || c.Vol == 0.0 {
			continue
		}
		if co, seen := comp[c.ID]; seen && co != c {
			t.Errorf("component %s duplicated as %+v and %+v", c.ID, c, co)
		} else if !seen {
			comp[c.ID] = c
		}
	}

	targets := []wtype.Coordinates{{0.0, -10.0, 0.0}, {0.0, 10.0, 0.0}}
	assert.Equal(t, targets, plate.GetTargets("Magick"), "Well Targets")
	assert.Equal(t, []wtype.Coordinates{}, plate.GetTargets("Muggles"), "Well Targets")
}

func assertTipsEqual(t *testing.T, a, b *wtype.LHTip, message string) {
	if a == nil && b == nil {
		fmt.Printf("Comparing nil tips at %s\n", message)
		return
	} else if a == nil {
		t.Errorf("Expected nil tip %s", message)
	} else if b == nil {
		t.Errorf("Unexpected nil tip %s", message)
	}

	assert.Equal(t, a.ID, b.ID, "Tip(%s) ID", message)
	assert.Equal(t, a.Mnfr, b.Mnfr, "Tip(%s) Mnfr", message)
	assert.Equal(t, a.Dirty, b.Dirty, "Tip(%s) Dirty", message)
	assert.Equal(t, a.MaxVol, b.MaxVol, "Tip(%s) MaxVol", message)
	assert.Equal(t, a.MinVol, b.MinVol, "Tip(%s) MinVol", message)
	assert.Equal(t, a.Bounds, b.Bounds, "Tip(%s) Bounds", message)
	//assuming contents are the same - IDs are changed though
	//assert.Equal(t, a.Contents(), b.Contents(), "Tip(%s) Contents", message)
	assert.Equal(t, a.Shape, b.Shape, "Tip(%s) Shape", message)
}

func assertTipBoxesEqual(t *testing.T, a, b *wtype.LHTipbox) {

	assert.Equal(t, a.ID, b.ID, "Tipbox ID")
	assert.Equal(t, a.Boxname, b.Boxname, "Tipbox Boxname")
	assert.Equal(t, a.Type, b.Type, "Tipbox Type")
	assert.Equal(t, a.Mnfr, b.Mnfr, "Tipbox Mnfr")
	assert.Equal(t, a.Nrows, b.Nrows, "Tipbox Nrows")
	assert.Equal(t, a.Ncols, b.Ncols, "Tipbox Ncols")
	assert.Equal(t, a.Height, b.Height, "Tipbox Height")
	assert.Equal(t, a.NTips, b.NTips, "Tipbox NTips")
	assert.Equal(t, a.TipXOffset, b.TipXOffset, "Tipbox TipXOffset")
	assert.Equal(t, a.TipYOffset, b.TipYOffset, "Tipbox TipYOffset")
	assert.Equal(t, a.TipXStart, b.TipXStart, "Tipbox TipXStart")
	assert.Equal(t, a.TipYStart, b.TipYStart, "Tipbox TipYStart")
	assert.Equal(t, a.TipZStart, b.TipZStart, "Tipbox TipZStart")

	assertTipsEqual(t, a.Tiptype, b.Tiptype, "Tipbox.Tiptype")

	for i := 0; i < a.Nrows; i++ {
		for j := 0; j < b.Ncols; j++ {
			assertTipsEqual(t, a.Tips[j][i], b.Tips[j][i], fmt.Sprintf("Tipbox.Tips[%d][%d]", j, i))
		}
	}
}

func AssertLHPropertiesEqual(t *testing.T, e, g *liquidhandling.LHProperties, msg string) {
	assert.Equalf(t, e.ID, g.ID, "%s: ID", msg)
	assert.Equalf(t, e.Nposns, g.Nposns, "%s: Nposns", msg)
	assert.Equalf(t, e.Positions, g.Positions, "%s: Positions", msg)
	assert.Equalf(t, e.PosLookup, g.PosLookup, "%s: PosLookup", msg)
	assert.Equalf(t, e.PlateIDLookup, g.PlateIDLookup, "%s: PlateIDLookup", msg)
	assert.Equalf(t, e.Devices, g.Devices, "%s: Devices", msg)
	assert.Equalf(t, e.Model, g.Model, "%s: Model", msg)
	assert.Equalf(t, e.Mnfr, g.Mnfr, "%s: Mnfr", msg)
	assert.Equalf(t, e.LHType, g.LHType, "%s: LHType", msg)
	assert.Equalf(t, e.TipType, g.TipType, "%s: TipType", msg)
	assert.Equalf(t, e.Tip_preferences, g.Tip_preferences, "%s: Tip_preferences", msg)
	assert.Equalf(t, e.Input_preferences, g.Input_preferences, "%s: Input_preferences", msg)
	assert.Equalf(t, e.Output_preferences, g.Output_preferences, "%s: Output_preferences", msg)
	assert.Equalf(t, e.Tipwaste_preferences, g.Tipwaste_preferences, "%s: Tipwaste_preferences", msg)
	assert.Equalf(t, e.Waste_preferences, g.Waste_preferences, "%s: Waste_preferences", msg)
	assert.Equalf(t, e.Wash_preferences, g.Wash_preferences, "%s: Wash_preferences", msg)
	assert.Equalf(t, e.CurrConf, g.CurrConf, "%s: CurrConf", msg)
	//commented out as decoding doesn't return uninitialised arrays
	//assert.Equalf(t, e.Cnfvol, g.Cnfvol, "%s: Cnfvol", msg)
	assert.Equalf(t, e.Layout, g.Layout, "%s: Layout", msg)
	assert.Equalf(t, e.MaterialType, g.MaterialType, "%s: MaterialType", msg)
	assert.Equalf(t, e.Heads, g.Heads, "%s: Heads", msg)
	assert.Equalf(t, e.Adaptors, g.Adaptors, "%s: Adaptors", msg)
	assert.Equalf(t, e.HeadAssemblies, g.HeadAssemblies, "%s: HeadAssemblies", msg)
}
