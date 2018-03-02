package liquidhandling

import (
	"context"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/antha/anthalib/wutil"
	"github.com/antha-lang/antha/inventory"
	"github.com/antha-lang/antha/inventory/testinventory"
	"testing"
)

func getTransferBlock2Component(ctx context.Context) (TransferBlockInstruction, *wtype.LHPlate) {
	inss := make([]*wtype.LHInstruction, 8)
	dstp, err := inventory.NewPlate(ctx, "pcrplate_skirted_riser40")
	if err != nil {
		panic(err)
	}

	for i := 0; i < 8; i++ {
		c, err := inventory.NewComponent(ctx, inventory.WaterType)
		if err != nil {
			panic(err)
		}

		c.Vol = 100.0
		c.Vunit = "ul"
		c.SetSample(true)
		ins := wtype.NewLHMixInstruction()
		ins.Components = append(ins.Components, c)

		c2, err := inventory.NewComponent(ctx, "tartrazine")
		if err != nil {
			panic(err)
		}

		c2.Vol = 24.0
		c2.Vunit = "ul"
		c2.SetSample(true)

		ins.Components = append(ins.Components, c2)

		c3 := c.Dup()
		c3.Mix(c2)
		ins.AddResult(c3)

		ins.Platetype = "pcrplate_skirted_riser40"
		ins.Welladdress = wutil.NumToAlpha(i+1) + "1"
		ins.SetPlateID(dstp.ID)
		inss[i] = ins
	}

	tb := NewTransferBlockInstruction(inss)

	return tb, dstp
}

func getTransferBlock3Component(ctx context.Context) (TransferBlockInstruction, *wtype.LHPlate) {
	inss := make([]*wtype.LHInstruction, 8)
	dstp, err := inventory.NewPlate(ctx, "pcrplate_skirted_riser40")
	if err != nil {
		panic(err)
	}

	for i := 0; i < 8; i++ {
		c, err := inventory.NewComponent(ctx, inventory.WaterType)
		if err != nil {
			panic(err)
		}

		c.Vol = 100.0
		c.Vunit = "ul"
		c.SetSample(true)
		ins := wtype.NewLHMixInstruction()
		ins.Components = append(ins.Components, c)

		c2, err := inventory.NewComponent(ctx, "tartrazine")
		if err != nil {
			panic(err)
		}

		c2.Vol = 24.0
		c2.Vunit = "ul"
		c2.SetSample(true)

		ins.Components = append(ins.Components, c2)

		c3, err := inventory.NewComponent(ctx, "ethanol")
		if err != nil {
			panic(err)
		}

		c3.Vol = 12.0
		c3.Vunit = "ul"
		c3.SetSample(true)

		ins.Components = append(ins.Components, c3)

		c4 := c.Dup()
		c4.Mix(c2)
		c4.Mix(c3)
		ins.AddResult(c4)

		ins.Platetype = "pcrplate_skirted_riser40"
		ins.Welladdress = wutil.NumToAlpha(i+1) + "1"
		ins.SetPlateID(dstp.ID)
		inss[i] = ins
	}

	tb := NewTransferBlockInstruction(inss)

	return tb, dstp
}

func getTestRobot(ctx context.Context, dstp *wtype.LHPlate, platetype string) *LHProperties {
	rbt, err := makeTestGilson(ctx)
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
	p.AddComponent(c, true)

	c, err = inventory.NewComponent(ctx, "tartrazine")
	if err != nil {
		panic(err)
	}
	c.Vol = v
	c.Vunit = "ul"

	p.AddComponent(c, true)

	c, err = inventory.NewComponent(ctx, "ethanol")
	if err != nil {
		panic(err)
	}
	c.Vol = v
	c.Vunit = "ul"

	p.AddComponent(c, true)

	rbt.AddPlate("position_4", p)

	// dst
	rbt.AddPlate("position_8", dstp)

	return rbt

}

func TestMultichannelFailPolicy(t *testing.T) {
	ctx := testinventory.NewContext(context.Background())

	// policy disallows
	tb, dstp := getTransferBlock2Component(ctx)
	rbt := getTestRobot(ctx, dstp, "pcrplate_skirted_riser40")
	pol, err := GetLHPolicyForTest()
	ris, err := tb.Generate(ctx, pol, rbt)
	if err != nil {
		t.Error(err)
	}

	testNegative(ctx, ris, pol, rbt, t)
}
func TestMultichannelSucceedSubset(t *testing.T) {
	ctx := testinventory.NewContext(context.Background())

	// can do 7
	tb, dstp := getTransferBlock2Component(ctx)

	tb.Inss[0].Welladdress = "B1"

	rbt := getTestRobot(ctx, dstp, "pcrplate_skirted_riser40")
	pol, err := GetLHPolicyForTest()
	pol.Policies["water"]["CAN_MULTI"] = true
	ris, err := tb.Generate(ctx, pol, rbt)
	if err != nil {
		t.Error(err)
	}

	//testNegative(ris, pol, rbt, t)
	testPositive(ctx, ris, pol, rbt, t)
}

func TestMultichannelSucceedPair(t *testing.T) {
	ctx := testinventory.NewContext(context.Background())

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
	pol, err := GetLHPolicyForTest()
	pol.Policies["water"]["CAN_MULTI"] = true
	ris, err := tb.Generate(ctx, pol, rbt)
	if err != nil {
		t.Error(err)
	}

	testPositive(ctx, ris, pol, rbt, t)
}

func TestMultichannelFailDest(t *testing.T) {
	ctx := testinventory.NewContext(context.Background())
	tb, dstp := getTransferBlock2Component(ctx)

	/*
		for i := 0; i < len(tb.Inss); i++ {
			tb.Inss[i].Welladdress = "A1"
		}
	*/

	rbt := getTestRobot(ctx, dstp, "pcrplate_skirted_riser40")
	pol, err := GetLHPolicyForTest()
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
	ctx := testinventory.NewContext(context.Background())

	// sources not aligned
	tb, dstp := getTransferBlock2Component(ctx)
	rbt := getTestRobot(ctx, dstp, "pcrplate_skirted_riser40")

	// fix the plate

	for i := 0; i < rbt.Plates["position_4"].WellsX(); i++ {
		rbt.Plates["position_4"].Cols[i][0].Clear()
	}

	pol, err := GetLHPolicyForTest()
	pol.Policies["water"]["CAN_MULTI"] = true
	ris, err := tb.Generate(ctx, pol, rbt)
	if err != nil {
		t.Error(err)
	}

	testNegative(ctx, ris, pol, rbt, t)
}

func TestMultiChannelFailComponent(t *testing.T) {
	ctx := testinventory.NewContext(context.Background())

	// components not same liquid type
	tb, dstp := getTransferBlock2Component(ctx)
	rbt := getTestRobot(ctx, dstp, "pcrplate_skirted_riser40")
	pol, err := GetLHPolicyForTest()
	pol.Policies["water"]["CAN_MULTI"] = true
	ris, err := tb.Generate(ctx, pol, rbt)
	if err != nil {
		t.Error(err)
	}

	if len(ris) != 1 {
		t.Errorf("Expected 1 transfer got ", len(ris))
	}

	tf := ris[0].(*TransferInstruction)

	if len(tf.Transfers) != 2 {
		t.Errorf("Error expected 2 transfers got %d", len(tf.Transfers))
	}

	ris[0].(*TransferInstruction).Transfers[0].Transfers[3].What = "lemonade"
	ris[0].(*TransferInstruction).Transfers[1].Transfers[3].What = "lemonade"

	testNegative(ctx, ris, pol, rbt, t)
}

func TestMultichannelPositive(t *testing.T) {
	ctx := testinventory.NewContext(context.Background())

	tb, dstp := getTransferBlock2Component(ctx)
	rbt := getTestRobot(ctx, dstp, "pcrplate_skirted_riser40")
	pol, err := GetLHPolicyForTest()

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
	ctx := testinventory.NewContext(context.Background())

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
	rbt.HeadsLoaded[0].Params.Independent = true

	pol, err := GetLHPolicyForTest()

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
	ctx := testinventory.NewContext(context.Background())

	tb, dstp := getTransferBlock2Component(ctx)

	rbt := getTestRobot(ctx, dstp, "DWST12_riser40")

	pol, err := GetLHPolicyForTest()

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
	ctx := testinventory.NewContext(context.Background())

	tb, dstp := getTransferBlock2Component(ctx)

	rbt := getTestRobot(ctx, dstp, "DSW24_riser40")

	pol, err := GetLHPolicyForTest()

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

func TestInsByInsMixPositiveMultichannel(t *testing.T) {
	ctx := testinventory.NewContext(context.Background())

	tb, dstp := getTransferBlock3Component(ctx)

	rbt := getTestRobot(ctx, dstp, "DWST12_riser40")

	pol, err := GetLHPolicyForTest()

	// allow multi

	pol.Policies["water"]["CAN_MULTI"] = true

	ris, err := tb.Generate(ctx, pol, rbt)

	if err != nil {
		t.Error(err)
	}

	// component-by-component multichanneling should be supported IFF
	// we can do all the solutions in the subset

	if len(ris) != 1 {
		t.Errorf("Error: Expected 1 transfer got %d", len(ris))
	}

	tf := ris[0].(*TransferInstruction)

	if len(tf.Transfers) != 3 {
		t.Errorf("Error: Expected 3 transfers got %d", len(tf.Transfers))
	}

	testPositive(ctx, ris, pol, rbt, t)
}

func TestInsByInsMixNegativeMultichannel(t *testing.T) {
	ctx := testinventory.NewContext(context.Background())

	tb, dstp := getTransferBlock3Component(ctx)

	rbt := getTestRobot(ctx, dstp, "DWST12_riser40")

	pol, err := GetLHPolicyForTest()

	// allow multi

	pol.Policies["water"]["CAN_MULTI"] = false

	ris, err := tb.Generate(ctx, pol, rbt)

	if err != nil {
		t.Error(err)
	}

	// atomic mixes now come through all split up... in future this should revert to the
	// older case of 8 x 3 transfers either by merging or something else

	if len(ris) != 1 {
		t.Errorf("Error: Expected 1 transfers got %d", len(ris))
	}

	tf := ris[0].(*TransferInstruction)

	if len(tf.Transfers) != 24 {
		t.Errorf("Error: Expected 24 transfers got %d", len(tf.Transfers))
	}

	testNegative(ctx, ris, pol, rbt, t)
}

func getMeATransfer(ctype string) *TransferInstruction {
	wh := []string{"a"}
	pltfrom := []string{"pos1"}
	pltto := []string{"pos2"}
	wellfrom := []string{"A1"}
	wellto := []string{"B1"}
	fplatetype := []string{"pcrplate_skirted"}
	tplatetype := []string{"anotherplate_type"}
	volume := []wunit.Volume{wunit.NewVolume(10.0, "ul")}
	fvolume := []wunit.Volume{wunit.ZeroVolume()}
	tvolume := []wunit.Volume{wunit.ZeroVolume()}
	fpwx := []int{12}
	fpwy := []int{8}
	tpwx := []int{12}
	tpwy := []int{8}
	cmps := []string{ctype}

	return NewTransferInstruction(wh, pltfrom, pltto, wellfrom, wellto, fplatetype, tplatetype, volume, fvolume, tvolume, fpwx, fpwy, tpwx, tpwy, cmps)
}

// TODO --> Create new version of the below
/*
func TestTransferMerge(t *testing.T) {
	policy, _ := GetLHPolicyForTest()
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

func testPositive(ctx context.Context, ris []RobotInstruction, pol *wtype.LHPolicyRuleSet, rbt *LHProperties, t *testing.T) {
	if len(ris) < 1 {
		t.Errorf("No instructions to test positive")
		return
	}
	ins := ris[0]

	ri2, err := ins.Generate(ctx, pol, rbt)

	if err != nil {
		t.Errorf(err.Error())
	}

	multi := 0
	single := 0
	for _, ri := range ri2 {
		if ri.InstructionType() == MCB {
			multi += 1
		} else if ri.InstructionType() == SCB {
			single += 1
		} else if ri.InstructionType() == TFR {
			t.Error("ERROR: Transfer generated from Transfer")
		}
	}

	if multi == 0 {
		t.Errorf("Multichannel block not generated")
	}
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
			if ri.InstructionType() != SCB {
				t.Errorf("Multichannel block generated without permission: %v %v %v", ri.GetParameter("LIQUIDCLASS"), ri.GetParameter("WELLFROM"), ri.GetParameter("WELLTO"))
			}
		}

	}
}
