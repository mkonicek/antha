package liquidhandling

import (
	"fmt"
	"testing"

	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wutil"
	"github.com/antha-lang/antha/microArch/factory"
)

func getTransferBlock() (TransferBlockInstruction, *wtype.LHPlate) {
	inss := make([]*wtype.LHInstruction, 8)
	dstp := factory.GetPlateByType("pcrplate_skirted_riser40")

	for i := 0; i < 8; i++ {
		c := factory.GetComponentByType("water")
		c.Vol = 100.0
		c.Vunit = "ul"
		c.SetSample(true)
		ins := wtype.NewLHInstruction()
		ins.Components = append(ins.Components, c)

		c2 := factory.GetComponentByType("tartrazine")
		c2.Vol = 24.0
		c2.Vunit = "ul"
		c2.SetSample(true)

		ins.Components = append(ins.Components, c2)

		c3 := c.Dup()
		c3.Mix(c2)
		ins.Result = c3

		ins.Platetype = "pcrplate_skirted_riser40"
		ins.Welladdress = wutil.NumToAlpha(i+1) + "1"
		ins.SetPlateID(dstp.ID)
		inss[i] = ins
	}

	tb := NewTransferBlockInstruction(inss)

	return tb, dstp
}

func getTestRobot(dstp *wtype.LHPlate) *LHProperties {
	rbt := makeTestGilson()

	// make a couple of plates

	// src

	p := factory.GetPlateByType("pcrplate_skirted_riser40")

	c := factory.GetComponentByType("water")
	c.Vol = 1600
	c.Vunit = "ul"
	p.AddComponent(c, true)

	c = factory.GetComponentByType("tartrazine")
	c.Vol = 1600
	c.Vunit = "ul"

	p.AddComponent(c, true)

	rbt.AddPlate("position_4", p)

	// dst
	rbt.AddPlate("position_8", dstp)

	return rbt

}

func TestMultichannelFailPolicy(t *testing.T) {
	// policy disallows
	tb, dstp := getTransferBlock()
	rbt := getTestRobot(dstp)
	pol, err := GetLHPolicyForTest()
	ris, err := tb.Generate(pol, rbt)
	if err != nil {
		t.Error(err)
	}

	testNegative(ris, pol, rbt, t)
}
func TestMultichannelSucceedSubset(t *testing.T) {
	// can do 7
	tb, dstp := getTransferBlock()

	tb.Inss[0].Welladdress = "B1"

	rbt := getTestRobot(dstp)
	pol, err := GetLHPolicyForTest()
	pol.Policies["water"]["CAN_MULTI"] = true
	ris, err := tb.Generate(pol, rbt)
	if err != nil {
		t.Error(err)
	}

	//testNegative(ris, pol, rbt, t)
	testPositive(ris, pol, rbt, t)
}

func TestMultichannelSucceedPair(t *testing.T) {
	// can do 7
	tb, dstp := getTransferBlock()

	tb.Inss[0].Welladdress = "A1"
	tb.Inss[1].Welladdress = "B1"

	tb.Inss[2].Welladdress = "C2"
	tb.Inss[3].Welladdress = "D2"

	tb.Inss[4].Welladdress = "E3"
	tb.Inss[5].Welladdress = "F3"

	tb.Inss[6].Welladdress = "G4"
	tb.Inss[7].Welladdress = "H4"

	rbt := getTestRobot(dstp)
	pol, err := GetLHPolicyForTest()
	pol.Policies["water"]["CAN_MULTI"] = true
	ris, err := tb.Generate(pol, rbt)
	if err != nil {
		t.Error(err)
	}

	testPositive(ris, pol, rbt, t)
}

func TestMultichannelFailDest(t *testing.T) {
	tb, dstp := getTransferBlock()

	/*
		for i := 0; i < len(tb.Inss); i++ {
			tb.Inss[i].Welladdress = "A1"
		}
	*/

	rbt := getTestRobot(dstp)
	pol, err := GetLHPolicyForTest()
	pol.Policies["water"]["CAN_MULTI"] = true
	ris, err := tb.Generate(pol, rbt)
	if err != nil {
		t.Error(err)
	}

	if len(ris) < 1 {
		t.Errorf("No Transfers made")
	}

	for x := 0; x < len(ris); x++ {
		ris[x].(*TransferInstruction).WellTo[1] = "A1"
		ris[x].(*TransferInstruction).WellTo[3] = "A1"
		ris[x].(*TransferInstruction).WellTo[5] = "A1"
		ris[x].(*TransferInstruction).WellTo[7] = "A1"
	}

	testNegative(ris, pol, rbt, t)
}
func TestMultiChannelFailSrc(t *testing.T) {
	fmt.Println("FAIL SRC")
	// sources not aligned
	tb, dstp := getTransferBlock()
	rbt := getTestRobot(dstp)

	// fix the plate

	for i := 0; i < rbt.Plates["position_4"].WellsX(); i++ {
		rbt.Plates["position_4"].Cols[i][0].Clear()
	}

	pol, err := GetLHPolicyForTest()
	pol.Policies["water"]["CAN_MULTI"] = true
	ris, err := tb.Generate(pol, rbt)
	if err != nil {
		t.Error(err)
	}

	for _, ins := range ris {
		fmt.Println(ins.(*TransferInstruction).FPlateType)
	}

	testNegative(ris, pol, rbt, t)
}

func TestMultiChannelFailComponent(t *testing.T) {
	// components not same liquid type
	tb, dstp := getTransferBlock()
	rbt := getTestRobot(dstp)
	pol, err := GetLHPolicyForTest()
	pol.Policies["water"]["CAN_MULTI"] = true
	ris, err := tb.Generate(pol, rbt)
	if err != nil {
		t.Error(err)
	}

	if len(ris) < 1 {
		t.Errorf("No Transfers made")
	}

	ris[0].(*TransferInstruction).What[3] = "lemonade"

	testNegative(ris, pol, rbt, t)
}

func testNegative(ris []RobotInstruction, pol *wtype.LHPolicyRuleSet, rbt *LHProperties, t *testing.T) {

	if len(ris) == 0 {
		t.Errorf("Error: No transfers generated")
	}

	for _, ins := range ris {
		ri2, err := ins.Generate(pol, rbt)

		if err != nil {
			t.Errorf(err.Error())
		}

		for _, ri := range ri2 {
			fmt.Println(ri.InstructionType(), " ", ri.GetParameter("LIQUIDCLASS"), ri.GetParameter("WELLFROM"), ri.GetParameter("WELLTO"))
			if ri.InstructionType() != SCB {
				t.Errorf("Multichannel block generated without permission: %v %v %v", ri.GetParameter("LIQUIDCLASS"), ri.GetParameter("WELLFROM"), ri.GetParameter("WELLTO"))
			}
		}

	}
}

func TestMultichannelPositive(t *testing.T) {
	tb, dstp := getTransferBlock()
	rbt := getTestRobot(dstp)
	pol, err := GetLHPolicyForTest()

	// allow multi

	pol.Policies["water"]["CAN_MULTI"] = true

	ris, err := tb.Generate(pol, rbt)

	if err != nil {
		t.Error(err)
	}

	if len(ris) != 2 {
		t.Errorf("Error: Expected 2 transfers got %d", len(ris))
	}

	testPositive(ris, pol, rbt, t)
}

func testPositive(ris []RobotInstruction, pol *wtype.LHPolicyRuleSet, rbt *LHProperties, t *testing.T) {
	ins := ris[0]

	ri2, err := ins.Generate(pol, rbt)

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
