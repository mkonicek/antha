package liquidhandling

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/antha/anthalib/wutil"
	"github.com/antha-lang/antha/inventory"
	"github.com/antha-lang/antha/inventory/cache/plateCache"
	"github.com/antha-lang/antha/inventory/testinventory"
)

func GetContextForTest() context.Context {
	ctx := testinventory.NewContext(context.Background())
	//also need to add a plateCache as we're not using the liquidhandler.Plan interface
	ctx = plateCache.NewContext(ctx)
	return ctx
}

func getComponent(ctx context.Context, name string, volume float64) (*wtype.Liquid, error) {
	c, err := inventory.NewComponent(ctx, name)
	if err != nil {
		return nil, err
	}
	c.Vol = volume
	c.Vunit = "ul"
	c.SetSample(true)
	return c, nil
}

func getMixInstructions(ctx context.Context, numInstructions int, componentNames []string, componentVolumes []float64) ([]*wtype.LHInstruction, error) {
	numComponents := len(componentNames)
	if len(componentVolumes) != numComponents {
		return nil, fmt.Errorf("componentNames and componentVolumes should be the same length")
	}

	ret := make([]*wtype.LHInstruction, 0, numComponents)

	for i := 0; i < numInstructions; i++ {

		components := make([]*wtype.Liquid, 0, numComponents)
		for j := 0; j < numComponents; j++ {
			c, err := getComponent(ctx, componentNames[j], componentVolumes[j])
			if err != nil {
				return nil, err
			}
			components = append(components, c)
		}

		ins := wtype.NewLHMixInstruction()
		ins.Inputs = append(ins.Inputs, components...)

		result := components[0].Dup()
		for j, c := range components {
			if j == 0 {
				continue
			}
			result.Mix(c)
		}
		ins.AddOutput(result)

		ret = append(ret, ins)
	}

	return ret, nil
}

func getTransferBlock(ctx context.Context, inss []*wtype.LHInstruction, destPlateType string) (*TransferBlockInstruction, *wtype.Plate) {
	if destPlateType == "" {
		destPlateType = "pcrplate_skirted_riser40"
	}

	dstp, err := inventory.NewPlate(ctx, destPlateType)
	if err != nil {
		panic(err)
	}

	for i, ins := range inss {
		ins.SetPlateID(dstp.ID)
		ins.Platetype = destPlateType
		ins.Welladdress = fmt.Sprintf("%s%d", wutil.NumToAlpha(i%8+1), i/8+1)
	}

	tb := NewTransferBlockInstruction(inss)

	return tb, dstp
}

func getTransferBlock2Component(ctx context.Context) (*TransferBlockInstruction, *wtype.Plate) {
	inss, err := getMixInstructions(ctx, 8, []string{inventory.WaterType, "tartrazine"}, []float64{100.0, 64.0})
	if err != nil {
		panic(err)
	}

	return getTransferBlock(ctx, inss, "pcrplate_skirted_riser40")
}

func getTestRobot(ctx context.Context, dstp *wtype.Plate, platetype string) *LHProperties {
	rbt, err := makeGilsonWithTipboxesForTest(ctx)
	if err != nil {
		panic(err)
	}

	// make a couple of plates

	// src

	p, err := inventory.NewPlate(ctx, platetype)
	if err != nil {
		panic(err)
	}

	c, err := inventory.NewComponent(ctx, inventory.WaterType)
	if err != nil {
		panic(err)
	}

	// add a columnsw'th

	v := p.ColVol().ConvertToString("ul")

	c.Vol = v
	c.Vunit = "ul"
	if _, err := p.AddComponent(c, true); err != nil {
		panic(err)
	}

	c, err = inventory.NewComponent(ctx, "tartrazine")
	if err != nil {
		panic(err)
	}
	c.Vol = v
	c.Vunit = "ul"

	if _, err := p.AddComponent(c, true); err != nil {
		panic(err)
	}

	c, err = inventory.NewComponent(ctx, "ethanol")
	if err != nil {
		panic(err)
	}
	c.Vol = v
	c.Vunit = "ul"

	if _, err := p.AddComponent(c, true); err != nil {
		panic(err)
	}

	if err := rbt.AddPlateTo("position_4", p); err != nil {
		panic(err)
	}

	// dst
	if err := rbt.AddPlateTo("position_8", dstp); err != nil {
		panic(err)
	}

	return rbt

}

func TestMultichannelFailPolicy(t *testing.T) {
	ctx := GetContextForTest()

	// policy disallows
	tb, dstp := getTransferBlock2Component(ctx)
	rbt := getTestRobot(ctx, dstp, "pcrplate_skirted_riser40")

	pol, err := wtype.GetLHPolicyForTest()
	if err != nil {
		t.Error(err)
	}
	ris, err := tb.Generate(ctx, pol, rbt)
	if err != nil {
		t.Error(err)
	}

	testNegative(ctx, ris, pol, rbt, t)
}
func TestMultichannelSucceedSubset(t *testing.T) {
	ctx := GetContextForTest()

	// can do 7
	tb, dstp := getTransferBlock2Component(ctx)

	tb.Inss[0].Welladdress = "B2"

	rbt := getTestRobot(ctx, dstp, "pcrplate_skirted_riser40")
	pol, err := wtype.GetLHPolicyForTest()
	if err != nil {
		t.Error(err)
	}
	pol.Policies["water"]["CAN_MULTI"] = true

	ris, err := tb.Generate(ctx, pol, rbt)
	if err != nil {
		t.Error(err)
	}

	//testNegative(ris, pol, rbt, t)
	testPositive(ctx, ris, pol, rbt, t)
}

func TestMultichannelSucceedPair(t *testing.T) {
	ctx := GetContextForTest()

	// can do 7
	tb, dstp := getTransferBlock2Component(ctx)

	tb.Inss[0].Welladdress = "A1"
	tb.Inss[1].Welladdress = "B1"

	tb.Inss[2].Welladdress = "C2"
	tb.Inss[3].Welladdress = "D2"

	tb.Inss[4].Welladdress = "E3"
	tb.Inss[5].Welladdress = "F3"

	tb.Inss[6].Welladdress = "G4"
	tb.Inss[7].Welladdress = "H4"

	rbt := getTestRobot(ctx, dstp, "pcrplate_skirted_riser40")
	pol, err := wtype.GetLHPolicyForTest()
	if err != nil {
		t.Error(err)
	}
	pol.Policies["water"]["CAN_MULTI"] = true
	ris, err := tb.Generate(ctx, pol, rbt)
	if err != nil {
		t.Error(err)
	}

	testPositive(ctx, ris, pol, rbt, t)
}

func TestMultichannelFailDest(t *testing.T) {
	ctx := GetContextForTest()
	tb, dstp := getTransferBlock2Component(ctx)

	/*
		for i := 0; i < len(tb.Inss); i++ {
			tb.Inss[i].Welladdress = "A1"
		}
	*/

	rbt := getTestRobot(ctx, dstp, "pcrplate_skirted_riser40")
	pol, err := wtype.GetLHPolicyForTest()
	if err != nil {
		t.Error(err)
	}
	pol.Policies["water"]["CAN_MULTI"] = true
	ris, err := tb.Generate(ctx, pol, rbt)
	if err != nil {
		t.Error(err)
	}

	if len(ris) < 1 {
		t.Errorf("No Transfers made")
	}

	for x := 0; x < len(ris[0].(*TransferInstruction).Transfers); x++ {
		ris[0].(*TransferInstruction).Transfers[x].Transfers[1].WellTo = "A1"
		ris[0].(*TransferInstruction).Transfers[x].Transfers[3].WellTo = "A1"
		ris[0].(*TransferInstruction).Transfers[x].Transfers[5].WellTo = "A1"
		ris[0].(*TransferInstruction).Transfers[x].Transfers[7].WellTo = "A1"
	}

	testNegative(ctx, ris, pol, rbt, t)
}
func TestMultiChannelFailSrc(t *testing.T) {
	// this actually works
	t.Skip()
	ctx := GetContextForTest()

	// sources not aligned
	tb, dstp := getTransferBlock2Component(ctx)
	rbt := getTestRobot(ctx, dstp, "pcrplate_skirted_riser40")

	// fix the plate

	for i := 0; i < rbt.Plates["position_4"].WellsX(); i++ {
		rbt.Plates["position_4"].Cols[i][0].Clear()
	}

	pol, err := wtype.GetLHPolicyForTest()
	if err != nil {
		t.Error(err)
	}
	pol.Policies["water"]["CAN_MULTI"] = true
	ris, err := tb.Generate(ctx, pol, rbt)
	if err != nil {
		t.Error(err)
	}

	testNegative(ctx, ris, pol, rbt, t)
}

func TestMultiChannelFailComponent(t *testing.T) {
	ctx := GetContextForTest()

	// components not same liquid type
	tb, dstp := getTransferBlock2Component(ctx)
	rbt := getTestRobot(ctx, dstp, "pcrplate_skirted_riser40")
	pol, err := wtype.GetLHPolicyForTest()
	if err != nil {
		t.Error(err)
	}
	// swap CAN_MULTI parameter of water and multiwater
	pol.Policies["water"]["CAN_MULTI"] = true
	pol.Policies["multiwater"]["CAN_MULTI"] = false
	ris, err := tb.Generate(ctx, pol, rbt)
	if err != nil {
		t.Error(err)
	}

	if len(ris) != 1 {
		t.Errorf("Expected 1 transfer got %d", len(ris))
	}

	tf := ris[0].(*TransferInstruction)

	if len(tf.Transfers) != 2 {
		t.Errorf("Error expected 2 transfers got %d", len(tf.Transfers))
	}

	ris[0].(*TransferInstruction).Transfers[0].Transfers[3].What = "multiwater"
	ris[0].(*TransferInstruction).Transfers[1].Transfers[3].What = "multiwater"

	testNegative(ctx, ris, pol, rbt, t)
}

func TestMultichannelPositive(t *testing.T) {
	ctx := GetContextForTest()

	tb, dstp := getTransferBlock2Component(ctx)
	rbt := getTestRobot(ctx, dstp, "pcrplate_skirted_riser40")
	pol, err := wtype.GetLHPolicyForTest()
	if err != nil {
		t.Error(err)
	}

	// allow multi
	pol.Policies["water"]["CAN_MULTI"] = true

	ris, err := tb.Generate(ctx, pol, rbt)

	if err != nil {
		t.Error(err)
	}

	if len(ris) != 1 {
		t.Errorf("Error: Expected 1 transfer got %d", len(ris))
	}

	tf := ris[0].(*TransferInstruction)

	if len(tf.Transfers) != 2 {
		t.Errorf("Error: Expected 2 transfers got %d", len(tf.Transfers))
	}

	testPositive(ctx, ris, pol, rbt, t)
}

func TestIndependentMultichannelPositive(t *testing.T) {
	ctx := GetContextForTest()

	tb, dstp := getTransferBlock2Component(ctx)

	ins := make([]*wtype.LHInstruction, 0, len(tb.Inss)-1)

	for i := 0; i < len(tb.Inss); i++ {
		// make one hole
		if i == 4 {
			continue
		}
		ins = append(ins, tb.Inss[i])
	}

	tb.Inss = ins

	rbt := getTestRobot(ctx, dstp, "pcrplate_skirted_riser40")

	// allow independent multichannel activity
	headsLoaded := rbt.GetLoadedHeads()
	headsLoaded[0].Params.Independent = true

	pol, err := wtype.GetLHPolicyForTest()
	if err != nil {
		t.Error(err)
	}

	// allow multi
	pol.Policies["water"]["CAN_MULTI"] = true

	ris, err := tb.Generate(ctx, pol, rbt)

	if err != nil {
		t.Error(err)
	}

	if len(ris) != 1 {
		t.Errorf("Error: Expected 1 transfer got %d", len(ris))
	}

	tf := ris[0].(*TransferInstruction)

	if len(tf.Transfers) != 2 {
		t.Errorf("Error: Expected 2 transfers got %d", len(tf.Transfers))
	}

	testPositive(ctx, ris, pol, rbt, t)
}

func TestTroughMultichannelPositive(t *testing.T) {
	ctx := GetContextForTest()

	tb, dstp := getTransferBlock2Component(ctx)

	rbt := getTestRobot(ctx, dstp, "DWST12_riser40")

	pol, err := wtype.GetLHPolicyForTest()
	if err != nil {
		t.Error(err)
	}

	// allow multi
	pol.Policies["water"]["CAN_MULTI"] = true

	ris, err := tb.Generate(ctx, pol, rbt)

	if err != nil {
		t.Error(err)
	}

	if len(ris) != 1 {
		t.Errorf("Error: Expected 1 transfer got %d", len(ris))
	}

	tf := ris[0].(*TransferInstruction)

	if len(tf.Transfers) != 2 {
		t.Errorf("Error: Expected 2 transfers got %d", len(tf.Transfers))
	}

	testPositive(ctx, ris, pol, rbt, t)
}

func TestBigWellMultichannelPositive(t *testing.T) {
	t.Skip() // pending revisions
	ctx := GetContextForTest()

	tb, dstp := getTransferBlock2Component(ctx)

	rbt := getTestRobot(ctx, dstp, "falcon6wellAgar_riser40")

	pol, err := wtype.GetLHPolicyForTest()
	if err != nil {
		t.Error(err)
	}

	// allow multi
	pol.Policies["water"]["CAN_MULTI"] = true

	ris, err := tb.Generate(ctx, pol, rbt)

	if err != nil {
		t.Error(err)
	}

	if len(ris) != 1 {
		t.Errorf("Error: Expected 1 transfer got %d", len(ris))
	}

	tf := ris[0].(*TransferInstruction)

	if len(tf.Transfers) != 2 {
		t.Errorf("Error: Expected 2 transfers got %d", len(tf.Transfers))
	}

	testPositive(ctx, ris, pol, rbt, t)
}

// TODO --> Create new version of the below
/*
func TestTransferMerge(t *testing.T) {
	policy, _ := wtype.GetLHPolicyForTest()
	ins1 := getMeATransfer("milk")

	toMerge := []*TransferInstruction{ins1, ins1}

	ins3 := ins1.Dup()
	ins3 = ins3.MergeWith(ins1)

	ins4 := mergeTransfers(toMerge, policy)[0]

	if !reflect.DeepEqual(ins3, ins4) {
		t.Errorf("Must merge transfers with same components")
	}

	// negative case

	ins2 := getMeATransfer("bile")

	toMerge = []*TransferInstruction{ins1, ins2}

	merged := mergeTransfers(toMerge, policy)

	if len(merged) == 1 {
		t.Errorf("Must not merge transfers with different components")
	}

}
*/

func testPositive(ctx context.Context, ris []RobotInstruction, pol *wtype.LHPolicyRuleSet, rbt *LHProperties, t *testing.T) []RobotInstruction {
	if len(ris) < 1 {
		t.Errorf("No instructions to test positive")
		return []RobotInstruction{}
	}
	ins := ris[0]

	ri2, err := ins.Generate(ctx, pol, rbt)

	if err != nil {
		t.Errorf(err.Error())
	}

	multi := 0
	single := 0
	for _, ri := range ri2 {
		switch ri.Type() {
		case CBI:
			ma := ri.GetParameter(MULTI).([]int)

			for _, m := range ma {
				if m > 1 {
					multi += 1
				} else {
					single += 1
				}
			}
		case TFR:
			t.Error("ERROR: Transfer generated from Transfer")
		}
	}

	if multi == 0 {
		t.Errorf("No multichannel transfers generated")
	}

	return ri2
}

func testNegative(ctx context.Context, ris []RobotInstruction, pol *wtype.LHPolicyRuleSet, rbt *LHProperties, t *testing.T) {

	if len(ris) == 0 {
		t.Errorf("Error: No transfers generated")
	}

	for _, ins := range ris {
		ri2, err := ins.Generate(ctx, pol, rbt)

		if err != nil {
			t.Errorf(err.Error())
		}

		for _, ri := range ri2 {
			if mcb, ok := ri.(*ChannelBlockInstruction); ok {
				if mcb.MaxMulti() > 1 {
					t.Errorf("Multichannel transfer(s) generated without permission: %v %v %v %v", ri.GetParameter(MULTI), ri.GetParameter(LIQUIDCLASS), ri.GetParameter(WELLFROM), ri.GetParameter(WELLTO))
				}

			}
		}

	}
}

func generateRobotInstructions(t *testing.T, ctx context.Context, inss []*wtype.LHInstruction, pol *wtype.LHPolicyRuleSet) []TerminalRobotInstruction {

	tb, dstp := getTransferBlock(ctx, inss, "pcrplate_skirted_riser40")

	rbt := getTestRobot(ctx, dstp, "pcrplate_skirted_riser40")
	var err error
	if pol == nil {
		pol, err = wtype.GetLHPolicyForTest()
		if err != nil {
			t.Fatal(err)
		}
		// allow multi
		pol.Policies["water"]["CAN_MULTI"] = true
	}

	//generate the low level instructions

	iTree := NewITree(tb)
	if _, err := iTree.Build(ctx, pol, rbt); err != nil {
		t.Fatal(err)
	}
	return iTree.Leaves()
}

func assertNumTipsUsed(t *testing.T, instructions []TerminalRobotInstruction, expectedTips int) {
	var loaded, unloaded int
	for _, instruction := range instructions {
		switch ins := instruction.(type) {
		case *LoadTipsInstruction:
			loaded += ins.Multi
		case *UnloadTipsInstruction:
			unloaded += ins.Multi
		}
	}

	if loaded != unloaded {
		t.Errorf("Loaded %d and Unloaded %d tips in instructions", loaded, unloaded)
	}

	if e, g := expectedTips, loaded; e != g {
		t.Errorf("Used %d tips, should have used %d", g, e)
	}

}

func assertNumLoadUnloadInstructions(t *testing.T, instructions []TerminalRobotInstruction, expected int) {
	var loads, unloads int

	for _, instruction := range instructions {
		switch instruction.(type) {
		case *LoadTipsInstruction:
			loads += 1
		case *UnloadTipsInstruction:
			unloads += 1
		}
	}

	if e, g := expected, loads; e != g {
		t.Errorf("Generated %d load tips instructions, expected %d", g, e)
	}

	if e, g := expected, unloads; e != g {
		t.Errorf("Generated %d unload tips instructions, expected %d", g, e)
	}
}

//TestChannelTipReuseGood Move water to two columns of wells - shouldn't need to change tips in between
func TestChannelTipReuseGood(t *testing.T) {
	ctx := GetContextForTest()

	inss, err := getMixInstructions(ctx, 16, []string{inventory.WaterType}, []float64{50.0})
	if err != nil {
		panic(err)
	}

	ris := generateRobotInstructions(t, ctx, inss, nil)

	assertNumTipsUsed(t, ris, 8)

	assertNumLoadUnloadInstructions(t, ris, 1)
}

//TestChannelTipReuseDisabled identical to good, except disable tip reuse
func TestChannelTipReuseDisabled(t *testing.T) {
	ctx := GetContextForTest()

	inss, err := getMixInstructions(ctx, 16, []string{inventory.WaterType}, []float64{50.0})
	if err != nil {
		panic(err)
	}

	pol, err := wtype.GetLHPolicyForTest()
	if err != nil {
		t.Fatal(err)
	}
	// allow multi
	pol.Policies["water"]["CAN_MULTI"] = true
	pol.Policies["water"]["TIP_REUSE_LIMIT"] = 0

	ris := generateRobotInstructions(t, ctx, inss, pol)

	assertNumTipsUsed(t, ris, 16)

	assertNumLoadUnloadInstructions(t, ris, 2)
}

//TestSingleChannelTipReuse -- based same as above but with multichannel disabled
//and allowing tip reuse
func TestSingleChannelTipReuse(t *testing.T) {
	ctx := GetContextForTest()

	inss, err := getMixInstructions(ctx, 16, []string{inventory.WaterType}, []float64{50.0})
	if err != nil {
		panic(err)
	}

	pol, err := wtype.GetLHPolicyForTest()
	if err != nil {
		t.Fatal(err)
	}
	// allow multi
	pol.Policies["water"]["CAN_MULTI"] = false

	ris := generateRobotInstructions(t, ctx, inss, pol)

	assertNumTipsUsed(t, ris, 1)

	assertNumLoadUnloadInstructions(t, ris, 1)
}

//TestSingleChannelTipReuse2 -- now we move two things
func TestSingleChannelTipReuse2(t *testing.T) {
	ctx := GetContextForTest()

	inss, err := getMixInstructions(ctx, 16, []string{inventory.WaterType}, []float64{50.0})
	if err != nil {
		panic(err)
	}

	ins2, err := getMixInstructions(ctx, 8, []string{"ethanol"}, []float64{50.0})
	if err != nil {
		panic(err)
	}

	inss = append(inss, ins2...)

	pol, err := wtype.GetLHPolicyForTest()
	if err != nil {
		t.Fatal(err)
	}
	// allow multi
	pol.Policies["water"]["CAN_MULTI"] = false

	ris := generateRobotInstructions(t, ctx, inss, pol)

	assertNumTipsUsed(t, ris, 2)

	assertNumLoadUnloadInstructions(t, ris, 2)
}

//TestChannelTipReuseBad Move water and ethanol to two separate columns of wells - should change tips in between
func TestChannelTipReuseBad(t *testing.T) {
	ctx := GetContextForTest()

	inss, err := getMixInstructions(ctx, 8, []string{inventory.WaterType}, []float64{50.0})
	if err != nil {
		panic(err)
	}

	ins2, err := getMixInstructions(ctx, 8, []string{"ethanol"}, []float64{50.0})
	if err != nil {
		panic(err)
	}

	inss = append(inss, ins2...)

	ris := generateRobotInstructions(t, ctx, inss, nil)

	assertNumTipsUsed(t, ris, 16)

	assertNumLoadUnloadInstructions(t, ris, 2)
}

//TestChannelTipReuseUgly Move water and ethanol to the same columns of wells - should change tips in between
func TestChannelTipReuseUgly(t *testing.T) {
	ctx := GetContextForTest()

	inss, err := getMixInstructions(ctx, 8, []string{inventory.WaterType, "ethanol"}, []float64{50.0, 50.0})
	if err != nil {
		panic(err)
	}

	ris := generateRobotInstructions(t, ctx, inss, nil)

	assertNumTipsUsed(t, ris, 16)

	assertNumLoadUnloadInstructions(t, ris, 2)
}

//TestChannelTipReuseUgly Move water and ethanol to the same columns of wells - should change tips in between
func BenchmarkChannelTipReuseUgly(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ctx := GetContextForTest()

		inss, err := getMixInstructions(ctx, 8, []string{inventory.WaterType, "ethanol"}, []float64{50.0, 50.0})
		if err != nil {
			panic(err)
		}

		generateRobotInstructions2(ctx, inss, nil)
	}
}

func generateRobotInstructions2(ctx context.Context, inss []*wtype.LHInstruction, pol *wtype.LHPolicyRuleSet) []TerminalRobotInstruction {

	tb, dstp := getTransferBlock(ctx, inss, "pcrplate_skirted_riser40")

	rbt := getTestRobot(ctx, dstp, "pcrplate_skirted_riser40")
	if pol == nil {
		pol, _ = wtype.GetLHPolicyForTest()
		// allow multi
		pol.Policies["water"]["CAN_MULTI"] = true
	}

	//generate the low level instructions
	iTree := NewITree(tb)

	if _, err := iTree.Build(ctx, pol, rbt); err != nil {
		panic(err)
	}

	return iTree.Leaves()
}

// regression test for issue with additional transfers being
// generated with sequential, different-length multichannel
// operations and non-independent heads
func TestMultiTransferError(t *testing.T) {
	ctx := GetContextForTest()

	// transfer two waters at volumes 50.0, 40.0
	inss, err := getMixInstructions(ctx, 2, []string{inventory.WaterType}, []float64{50.0})

	if err != nil {
		t.Errorf(err.Error())
	}

	inss[1].Inputs[0].Vol = 40.0

	tb, p := getTransferBlock(ctx, inss, "pcrplate_skirted_riser18")

	rbt := getTestRobot(ctx, p, "pcrplate_skirted_riser40")

	pol, err := wtype.GetLHPolicyForTest()
	if err != nil {
		t.Error(err)
	}

	pol.Policies["water"]["CAN_MULTI"] = true

	ris, err := tb.Generate(ctx, pol, rbt)
	if err != nil {
		t.Error(err)
	}

	rr, err := ris[0].Generate(ctx, pol, rbt)

	if err != nil {
		t.Error(err)
	}

	if len(rr) != 1 {
		t.Errorf("Expected 1 instruction got %d", len(rr))
	}

	mcb, ok := rr[0].(*ChannelBlockInstruction)

	if !ok {
		t.Errorf("Expected *ChannelBlockInstruction, got %T", rr)
	}

	if len(mcb.What) != 2 {
		t.Errorf("Expected 2 transfers, got %d", len(mcb.What))
	}

	// we expect 50ul, 40ul to be transferred

	volSums := []wunit.Volume{wunit.ZeroVolume(), wunit.ZeroVolume()}

	for i := 0; i < len(mcb.What); i++ {
		for j := 0; j < len(mcb.What[i]); j++ {
			if mcb.What[i][j] != "" {
				volSums[j].Add(mcb.Volume[i][j])
			}
		}
	}

	fiftyul := wunit.NewVolume(50.0, "ul")
	fortyul := wunit.NewVolume(40.0, "ul")

	if !reflect.DeepEqual(volSums, []wunit.Volume{fiftyul, fortyul}) {
		t.Errorf("Volumes inconsistent: expected %v got %v", []wunit.Volume{fiftyul, fortyul}, volSums)
	}

}

// ensure that gapped transfers are correctly handled
// i.e. [A1:50, B1:40, C1:50] should be done as
//      [A1:40, B1:40, C1:40], [A1:10], [C1:10]
// on a non-independent liquid handler
func TestGappedTransfer(t *testing.T) {
	ctx := GetContextForTest()

	// transfer three waters at volumes 50.0, 40.0, 50.0
	inss, err := getMixInstructions(ctx, 3, []string{inventory.WaterType}, []float64{50.0})

	if err != nil {
		t.Errorf(err.Error())
	}

	inss[1].Inputs[0].Vol = 40.0

	tb, p := getTransferBlock(ctx, inss, "pcrplate_skirted_riser18")

	rbt := getTestRobot(ctx, p, "pcrplate_skirted_riser40")

	pol, err := wtype.GetLHPolicyForTest()
	if err != nil {
		t.Error(err)
	}

	pol.Policies["water"]["CAN_MULTI"] = true

	ris, err := tb.Generate(ctx, pol, rbt)
	if err != nil {
		t.Error(err)
	}

	rr, err := ris[0].Generate(ctx, pol, rbt)

	if err != nil {
		t.Error(err)
	}

	mcb, ok := rr[0].(*ChannelBlockInstruction)

	if !ok {
		t.Errorf("Expected *ChannelBlockInstruction, got %T", rr)
	}

	if len(mcb.What) != 3 {
		t.Errorf("Expected 3 transfers in ChannelBlock, got %d", len(mcb.What))
	}

	// we expect 50ul, 40ul, 50ul to be transferred

	volSums := map[string]wunit.Volume{"A1": wunit.ZeroVolume(), "B1": wunit.ZeroVolume(), "C1": wunit.ZeroVolume()}

	for i := 0; i < len(mcb.What); i++ {
		for j := 0; j < len(mcb.What[i]); j++ {
			if !mcb.Volume[i][j].IsZero() {
				volSums[mcb.WellTo[i][j]].Add(mcb.Volume[i][j])
			}
		}
	}

	fiftyul := wunit.NewVolume(50.0, "ul")
	fortyul := wunit.NewVolume(40.0, "ul")

	expect := map[string]wunit.Volume{"A1": fiftyul, "B1": fortyul, "C1": fiftyul}
	if !reflect.DeepEqual(volSums, expect) {
		t.Errorf("Volumes inconsistent: expected %v got %v", expect, volSums)
	}
}

// another regression test for 2357
// ensure that convertinstructions does not
// modify execution order

// -- the issue here is how getComponents works when getting
// singles, it takes the highest volume first, rather than
// preserving mix order. Feeding in components one at a time should fix this
func TestTransferBlockMixOrdering(t *testing.T) {
	ctx := GetContextForTest()

	// transfer three things at volumes 30.0, 40.0, 50.0 in one instruction
	inss, err := getMixInstructions(ctx, 1, []string{"ethanol", "tartrazine", inventory.WaterType}, []float64{30.0, 40.0, 50.0})

	if err != nil {
		t.Errorf(err.Error())
	}

	tb, p := getTransferBlock(ctx, inss, "eppendorfrack424_1.5ml_lidholder")

	rbt := getTestRobot(ctx, p, "pcrplate_skirted_riser40")

	pol, err := wtype.GetLHPolicyForTest()
	if err != nil {
		t.Error(err)
	}

	pol.Policies["water"]["CAN_MULTI"] = true
	pol.Policies["ethanol"] = map[string]interface{}{"CAN_MULTI": true}
	pol.Policies["tartrazine"] = map[string]interface{}{"CAN_MULTI": true}

	ris, err := tb.Generate(ctx, pol, rbt)
	if err != nil {
		t.Error(err)
	}

	rr, err := ris[0].Generate(ctx, pol, rbt)

	if err != nil {
		t.Error(err)
	}

	mcb, ok := rr[0].(*ChannelBlockInstruction)

	if !ok {
		t.Errorf("Expected *ChannelBlockInstruction, got %T", rr)
	}

	if len(mcb.What) != 3 {
		t.Errorf("Expected 3 transfers in ChannelBlock, got %d", len(mcb.What))
	}

	// we expect 30ul, 40ul, 50ul to be transferred to A1

	volSums := map[string]wunit.Volume{"A1": wunit.ZeroVolume()}

	for i := 0; i < len(mcb.What); i++ {
		for j := 0; j < len(mcb.What[i]); j++ {
			if !mcb.Volume[i][j].IsZero() {
				volSums[mcb.WellTo[i][j]].Add(mcb.Volume[i][j])
			}
		}
	}

	onetwentyul := wunit.NewVolume(120.0, "ul")

	expect := map[string]wunit.Volume{"A1": onetwentyul}
	if !reflect.DeepEqual(volSums, expect) {
		t.Errorf("Volumes inconsistent: expected %v got %v", expect, volSums)
	}

	// check ordering, must be preserved

	expectOrder := []string{"ethanol", "tartrazine", inventory.WaterType}

	gotOrder := make([]string, 3)

	for i := 0; i < len(mcb.What); i++ {
		gotOrder[i] = mcb.Component[i][0]
	}

	if !reflect.DeepEqual(expectOrder, gotOrder) {
		t.Errorf("Order inconsistent: Expected %v got %v", expectOrder, gotOrder)
	}
}
