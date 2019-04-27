package tests

import (
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/antha/anthalib/wutil"
	"github.com/antha-lang/antha/inventory/components"
	"github.com/antha-lang/antha/laboratory"
	"github.com/antha-lang/antha/laboratory/effects"
	"github.com/antha-lang/antha/laboratory/testlab"
	"github.com/antha-lang/antha/microArch/driver/liquidhandling"
	"github.com/antha-lang/antha/utils"
)

func getComponent(lab *laboratory.Laboratory, name string, volume float64) (*wtype.Liquid, error) {
	c, err := lab.Inventory.Components.NewComponent(name)
	if err != nil {
		return nil, err
	}
	c.Vol = volume
	c.Vunit = "ul"
	c.SetSample(true)
	return c, nil
}

func getMixInstructions(lab *laboratory.Laboratory, numInstructions int, componentNames []string, componentVolumes []float64) ([]*wtype.LHInstruction, error) {
	numComponents := len(componentNames)
	if len(componentVolumes) != numComponents {
		return nil, fmt.Errorf("componentNames and componentVolumes should be the same length")
	}

	ret := make([]*wtype.LHInstruction, 0, numComponents)

	for i := 0; i < numInstructions; i++ {

		components := make([]*wtype.Liquid, 0, numComponents)
		for j := 0; j < numComponents; j++ {
			c, err := getComponent(lab, componentNames[j], componentVolumes[j])
			if err != nil {
				return nil, err
			}
			components = append(components, c)
		}

		ins := wtype.NewLHMixInstruction(lab.IDGenerator)
		ins.Inputs = append(ins.Inputs, components...)

		result := components[0].Dup(lab.IDGenerator)
		for j, c := range components {
			if j == 0 {
				continue
			}
			result.Mix(lab.IDGenerator, c)
		}
		ins.AddOutput(result)

		ret = append(ret, ins)
	}

	return ret, nil
}

func getTransferBlock(lab *laboratory.Laboratory, inss []*wtype.LHInstruction, destPlateType wtype.PlateTypeName) (*liquidhandling.TransferBlockInstruction, *wtype.Plate) {
	if destPlateType == "" {
		destPlateType = "pcrplate_skirted_riser40"
	}

	dstp, err := lab.Inventory.Plates.NewPlate(destPlateType)
	if err != nil {
		panic(err)
	}

	for i, ins := range inss {
		ins.SetPlateID(dstp.ID)
		ins.Platetype = destPlateType
		ins.Welladdress = fmt.Sprintf("%s%d", wutil.NumToAlpha(i%8+1), i/8+1)
	}

	tb := liquidhandling.NewTransferBlockInstruction(inss)

	return tb, dstp
}

func getTransferBlock2Component(lab *laboratory.Laboratory) (*liquidhandling.TransferBlockInstruction, *wtype.Plate) {
	inss, err := getMixInstructions(lab, 8, []string{components.WaterType, "tartrazine"}, []float64{100.0, 64.0})
	if err != nil {
		panic(err)
	}

	return getTransferBlock(lab, inss, "pcrplate_skirted_riser40")
}

func getTestRobot(lab *laboratory.Laboratory, dstp *wtype.Plate, platetype wtype.PlateTypeName) *liquidhandling.LHProperties {
	rbt, err := makeGilsonWithTipboxesForTest(lab)
	if err != nil {
		panic(err)
	}

	// make a couple of plates

	// src

	p, err := lab.Inventory.Plates.NewPlate(platetype)
	if err != nil {
		panic(err)
	}

	c, err := lab.Inventory.Components.NewComponent(components.WaterType)
	if err != nil {
		panic(err)
	}

	// add a columnsw'th

	v := p.ColVol().ConvertToString("ul")

	c.Vol = v
	c.Vunit = "ul"
	if _, err := p.AddComponent(lab.IDGenerator, c, true); err != nil {
		panic(err)
	}

	c, err = lab.Inventory.Components.NewComponent("tartrazine")
	if err != nil {
		panic(err)
	}
	c.Vol = v
	c.Vunit = "ul"

	if _, err := p.AddComponent(lab.IDGenerator, c, true); err != nil {
		panic(err)
	}

	c, err = lab.Inventory.Components.NewComponent("ethanol")
	if err != nil {
		panic(err)
	}
	c.Vol = v
	c.Vunit = "ul"

	if _, err := p.AddComponent(lab.IDGenerator, c, true); err != nil {
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
	testlab.WithTestLab(t, "", &testlab.TestElementCallbacks{
		Steps: func(lab *laboratory.Laboratory) error {

			// policy disallows
			tb, dstp := getTransferBlock2Component(lab)
			rbt := getTestRobot(lab, dstp, "pcrplate_skirted_riser40")

			pol, err := wtype.GetLHPolicyForTest()
			if err != nil {
				return err
			}
			ris, err := tb.Generate(lab.LaboratoryEffects, pol, rbt)
			if err != nil {
				return err
			}

			return testNegative(lab.LaboratoryEffects, ris, pol, rbt)
		},
	})
}

func TestMultichannelSucceedSubset(t *testing.T) {
	testlab.WithTestLab(t, "", &testlab.TestElementCallbacks{
		Steps: func(lab *laboratory.Laboratory) error {

			// can do 7
			tb, dstp := getTransferBlock2Component(lab)

			tb.Inss[0].Welladdress = "B2"

			rbt := getTestRobot(lab, dstp, "pcrplate_skirted_riser40")
			pol, err := wtype.GetLHPolicyForTest()
			if err != nil {
				return err
			}
			pol.Policies["water"]["CAN_MULTI"] = true

			ris, err := tb.Generate(lab.LaboratoryEffects, pol, rbt)
			if err != nil {
				return err
			}

			//testNegative(lab ,ris, pol, rbt)
			return testPositive(lab.LaboratoryEffects, ris, pol, rbt)
		},
	})
}

func TestMultichannelSucceedPair(t *testing.T) {
	testlab.WithTestLab(t, "", &testlab.TestElementCallbacks{
		Steps: func(lab *laboratory.Laboratory) error {

			// can do 7
			tb, dstp := getTransferBlock2Component(lab)

			tb.Inss[0].Welladdress = "A1"
			tb.Inss[1].Welladdress = "B1"

			tb.Inss[2].Welladdress = "C2"
			tb.Inss[3].Welladdress = "D2"

			tb.Inss[4].Welladdress = "E3"
			tb.Inss[5].Welladdress = "F3"

			tb.Inss[6].Welladdress = "G4"
			tb.Inss[7].Welladdress = "H4"

			rbt := getTestRobot(lab, dstp, "pcrplate_skirted_riser40")
			pol, err := wtype.GetLHPolicyForTest()
			if err != nil {
				return err
			}
			pol.Policies["water"]["CAN_MULTI"] = true
			ris, err := tb.Generate(lab.LaboratoryEffects, pol, rbt)
			if err != nil {
				return err
			}

			return testPositive(lab.LaboratoryEffects, ris, pol, rbt)
		},
	})
}

func TestMultichannelFailDest(t *testing.T) {
	testlab.WithTestLab(t, "", &testlab.TestElementCallbacks{
		Steps: func(lab *laboratory.Laboratory) error {

			tb, dstp := getTransferBlock2Component(lab)

			/*
				for i := 0; i < len(tb.Inss); i++ {
					tb.Inss[i].Welladdress = "A1"
				}
			*/

			rbt := getTestRobot(lab, dstp, "pcrplate_skirted_riser40")
			pol, err := wtype.GetLHPolicyForTest()
			if err != nil {
				return err
			}
			pol.Policies["water"]["CAN_MULTI"] = true
			ris, err := tb.Generate(lab.LaboratoryEffects, pol, rbt)
			if err != nil {
				return err
			}

			if len(ris) < 1 {
				return errors.New("No Transfers made")
			}

			for x := 0; x < len(ris[0].(*liquidhandling.TransferInstruction).Transfers); x++ {
				ris[0].(*liquidhandling.TransferInstruction).Transfers[x].Transfers[1].WellTo = "A1"
				ris[0].(*liquidhandling.TransferInstruction).Transfers[x].Transfers[3].WellTo = "A1"
				ris[0].(*liquidhandling.TransferInstruction).Transfers[x].Transfers[5].WellTo = "A1"
				ris[0].(*liquidhandling.TransferInstruction).Transfers[x].Transfers[7].WellTo = "A1"
			}

			return testNegative(lab.LaboratoryEffects, ris, pol, rbt)
		},
	})
}

func TestMultiChannelFailComponent(t *testing.T) {
	testlab.WithTestLab(t, "", &testlab.TestElementCallbacks{
		Steps: func(lab *laboratory.Laboratory) error {
			// components not same liquid type
			tb, dstp := getTransferBlock2Component(lab)
			rbt := getTestRobot(lab, dstp, "pcrplate_skirted_riser40")
			pol, err := wtype.GetLHPolicyForTest()
			if err != nil {
				return err
			}
			// swap CAN_MULTI parameter of water and multiwater
			pol.Policies["water"]["CAN_MULTI"] = true
			pol.Policies["multiwater"]["CAN_MULTI"] = false
			ris, err := tb.Generate(lab.LaboratoryEffects, pol, rbt)
			if err != nil {
				return err
			}

			if len(ris) != 1 {
				return fmt.Errorf("Expected 1 transfer got %d", len(ris))
			}

			tf := ris[0].(*liquidhandling.TransferInstruction)

			if len(tf.Transfers) != 2 {
				return fmt.Errorf("Error expected 2 transfers got %d", len(tf.Transfers))
			}

			ris[0].(*liquidhandling.TransferInstruction).Transfers[0].Transfers[3].What = "multiwater"
			ris[0].(*liquidhandling.TransferInstruction).Transfers[1].Transfers[3].What = "multiwater"

			return testNegative(lab.LaboratoryEffects, ris, pol, rbt)
		},
	})
}

func TestMultichannelPositive(t *testing.T) {
	testlab.WithTestLab(t, "", &testlab.TestElementCallbacks{
		Steps: func(lab *laboratory.Laboratory) error {
			tb, dstp := getTransferBlock2Component(lab)
			rbt := getTestRobot(lab, dstp, "pcrplate_skirted_riser40")
			pol, err := wtype.GetLHPolicyForTest()
			if err != nil {
				return err
			}

			// allow multi
			pol.Policies["water"]["CAN_MULTI"] = true

			ris, err := tb.Generate(lab.LaboratoryEffects, pol, rbt)

			if err != nil {
				return err
			}

			if len(ris) != 1 {
				return fmt.Errorf("Error: Expected 1 transfer got %d", len(ris))
			}

			tf := ris[0].(*liquidhandling.TransferInstruction)

			if len(tf.Transfers) != 2 {
				return fmt.Errorf("Error: Expected 2 transfers got %d", len(tf.Transfers))
			}

			return testPositive(lab.LaboratoryEffects, ris, pol, rbt)
		},
	})
}

func TestIndependentMultichannelPositive(t *testing.T) {
	testlab.WithTestLab(t, "", &testlab.TestElementCallbacks{
		Steps: func(lab *laboratory.Laboratory) error {
			tb, dstp := getTransferBlock2Component(lab)

			ins := make([]*wtype.LHInstruction, 0, len(tb.Inss)-1)

			for i := 0; i < len(tb.Inss); i++ {
				// make one hole
				if i == 4 {
					continue
				}
				ins = append(ins, tb.Inss[i])
			}

			tb.Inss = ins

			rbt := getTestRobot(lab, dstp, "pcrplate_skirted_riser40")

			// allow independent multichannel activity
			headsLoaded := rbt.GetLoadedHeads()
			headsLoaded[0].Params.Independent = true

			pol, err := wtype.GetLHPolicyForTest()
			if err != nil {
				return err
			}

			// allow multi
			pol.Policies["water"]["CAN_MULTI"] = true

			ris, err := tb.Generate(lab.LaboratoryEffects, pol, rbt)

			if err != nil {
				return err
			}

			if len(ris) != 1 {
				return fmt.Errorf("Error: Expected 1 transfer got %d", len(ris))
			}

			tf := ris[0].(*liquidhandling.TransferInstruction)

			if len(tf.Transfers) != 2 {
				return fmt.Errorf("Error: Expected 2 transfers got %d", len(tf.Transfers))
			}

			return testPositive(lab.LaboratoryEffects, ris, pol, rbt)
		},
	})
}

func TestTroughMultichannelPositive(t *testing.T) {
	testlab.WithTestLab(t, "", &testlab.TestElementCallbacks{
		Steps: func(lab *laboratory.Laboratory) error {
			tb, dstp := getTransferBlock2Component(lab)

			rbt := getTestRobot(lab, dstp, "DWST12_riser40")

			pol, err := wtype.GetLHPolicyForTest()
			if err != nil {
				return err
			}

			// allow multi
			pol.Policies["water"]["CAN_MULTI"] = true

			ris, err := tb.Generate(lab.LaboratoryEffects, pol, rbt)

			if err != nil {
				return err
			}

			if len(ris) != 1 {
				return fmt.Errorf("Error: Expected 1 transfer got %d", len(ris))
			}

			tf := ris[0].(*liquidhandling.TransferInstruction)

			if len(tf.Transfers) != 2 {
				return fmt.Errorf("Error: Expected 2 transfers got %d", len(tf.Transfers))
			}

			return testPositive(lab.LaboratoryEffects, ris, pol, rbt)
		},
	})
}

/*
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

	tf := ris[0].(*liquidhandling.TransferInstruction)

	if len(tf.Transfers) != 2 {
		t.Errorf("Error: Expected 2 transfers got %d", len(tf.Transfers))
	}

	testPositive(ctx, ris, pol, rbt, t)
}

// TODO --> Create new version of the below
func TestTransferMerge(t *testing.T) {
	policy, _ := wtype.GetLHPolicyForTest()
	ins1 := getMeATransfer("milk")

	toMerge := []*liquidhandling.TransferInstruction{ins1, ins1}

	ins3 := ins1.Dup()
	ins3 = ins3.MergeWith(ins1)

	ins4 := mergeTransfers(toMerge, policy)[0]

	if !reflect.DeepEqual(ins3, ins4) {
		t.Errorf("Must merge transfers with same components")
	}

	// negative case

	ins2 := getMeATransfer("bile")

	toMerge = []*liquidhandling.TransferInstruction{ins1, ins2}

	merged := mergeTransfers(toMerge, policy)

	if len(merged) == 1 {
		t.Errorf("Must not merge transfers with different components")
	}

}
*/

func testPositive(labEffects *effects.LaboratoryEffects, ris []liquidhandling.RobotInstruction, pol *wtype.LHPolicyRuleSet, rbt *liquidhandling.LHProperties) error {
	if len(ris) < 1 {
		return fmt.Errorf("No instructions to test positive")
	}
	ins := ris[0]

	ri2, err := ins.Generate(labEffects, pol, rbt)

	if err != nil {
		return err
	}

	multi := 0
	single := 0
	for _, ri := range ri2 {
		switch ri.Type() {
		case liquidhandling.CBI:
			ma := ri.GetParameter(liquidhandling.MULTI).([]int)

			for _, m := range ma {
				if m > 1 {
					multi += 1
				} else {
					single += 1
				}
			}
		case liquidhandling.TFR:
			return errors.New("ERROR: Transfer generated from Transfer")
		}
	}

	if multi == 0 {
		return errors.New("No multichannel transfers generated")
	}

	return nil
}

func testNegative(labEffects *effects.LaboratoryEffects, ris []liquidhandling.RobotInstruction, pol *wtype.LHPolicyRuleSet, rbt *liquidhandling.LHProperties) error {
	if len(ris) == 0 {
		return errors.New("Error: No transfers generated")
	}

	for _, ins := range ris {
		ri2, err := ins.Generate(labEffects, pol, rbt)

		if err != nil {
			return err
		}

		for _, ri := range ri2 {
			if mcb, ok := ri.(*liquidhandling.ChannelBlockInstruction); ok {
				if mcb.MaxMulti() > 1 {
					return fmt.Errorf("Multichannel transfer(s) generated without permission: %v %v %v %v", ri.GetParameter(liquidhandling.MULTI), ri.GetParameter(liquidhandling.LIQUIDCLASS), ri.GetParameter(liquidhandling.WELLFROM), ri.GetParameter(liquidhandling.WELLTO))
				}

			}
		}

	}
	return nil
}

func generateRobotInstructions(lab *laboratory.Laboratory, inss []*wtype.LHInstruction, pol *wtype.LHPolicyRuleSet) ([]liquidhandling.TerminalRobotInstruction, error) {
	tb, dstp := getTransferBlock(lab, inss, "pcrplate_skirted_riser40")

	rbt := getTestRobot(lab, dstp, "pcrplate_skirted_riser40")
	var err error
	if pol == nil {
		pol, err = wtype.GetLHPolicyForTest()
		if err != nil {
			return nil, err
		}
		// allow multi
		pol.Policies["water"]["CAN_MULTI"] = true
	}

	//generate the low level instructions

	iTree := liquidhandling.NewITree(tb)
	if _, err := iTree.Build(lab.LaboratoryEffects, pol, rbt); err != nil {
		return nil, err
	}
	return iTree.Leaves(), nil
}

func assertNumTipsUsed(instructions []liquidhandling.TerminalRobotInstruction, expectedTips int) error {
	var loaded, unloaded int
	for _, instruction := range instructions {
		switch ins := instruction.(type) {
		case *liquidhandling.LoadTipsInstruction:
			loaded += ins.Multi
		case *liquidhandling.UnloadTipsInstruction:
			unloaded += ins.Multi
		}
	}

	if loaded != unloaded {
		return fmt.Errorf("Loaded %d and Unloaded %d tips in instructions", loaded, unloaded)
	}

	if e, g := expectedTips, loaded; e != g {
		return fmt.Errorf("Used %d tips, should have used %d", g, e)
	}
	return nil
}

func assertNumLoadUnloadInstructions(instructions []liquidhandling.TerminalRobotInstruction, expected int) error {
	var loads, unloads int

	for _, instruction := range instructions {
		switch instruction.(type) {
		case *liquidhandling.LoadTipsInstruction:
			loads += 1
		case *liquidhandling.UnloadTipsInstruction:
			unloads += 1
		}
	}

	if e, g := expected, loads; e != g {
		return fmt.Errorf("Generated %d load tips instructions, expected %d", g, e)
	}

	if e, g := expected, unloads; e != g {
		return fmt.Errorf("Generated %d unload tips instructions, expected %d", g, e)
	}
	return nil
}

//TestChannelTipReuseGood Move water to two columns of wells - shouldn't need to change tips in between
func TestChannelTipReuseGood(t *testing.T) {
	testlab.WithTestLab(t, "", &testlab.TestElementCallbacks{
		Steps: func(lab *laboratory.Laboratory) error {
			inss, err := getMixInstructions(lab, 16, []string{components.WaterType}, []float64{50.0})
			if err != nil {
				return err
			}

			if ris, err := generateRobotInstructions(lab, inss, nil); err != nil {
				return err
			} else {

				return utils.ErrorSlice{
					assertNumTipsUsed(ris, 8),
					assertNumLoadUnloadInstructions(ris, 1),
				}.Pack()
			}
		},
	})
}

//TestChannelTipReuseDisabled identical to good, except disable tip reuse
func TestChannelTipReuseDisabled(t *testing.T) {
	testlab.WithTestLab(t, "", &testlab.TestElementCallbacks{
		Steps: func(lab *laboratory.Laboratory) error {
			inss, err := getMixInstructions(lab, 16, []string{components.WaterType}, []float64{50.0})
			if err != nil {
				return err
			}

			pol, err := wtype.GetLHPolicyForTest()
			if err != nil {
				return err
			}
			// allow multi
			pol.Policies["water"]["CAN_MULTI"] = true
			pol.Policies["water"]["TIP_REUSE_LIMIT"] = 0

			if ris, err := generateRobotInstructions(lab, inss, pol); err != nil {
				return err
			} else {
				return utils.ErrorSlice{
					assertNumTipsUsed(ris, 16),
					assertNumLoadUnloadInstructions(ris, 2),
				}.Pack()
			}
		},
	})
}

//TestSingleChannelTipReuse -- based same as above but with multichannel disabled
//and allowing tip reuse
func TestSingleChannelTipReuse(t *testing.T) {
	testlab.WithTestLab(t, "", &testlab.TestElementCallbacks{
		Steps: func(lab *laboratory.Laboratory) error {
			inss, err := getMixInstructions(lab, 16, []string{components.WaterType}, []float64{50.0})
			if err != nil {
				return err
			}

			pol, err := wtype.GetLHPolicyForTest()
			if err != nil {
				return err
			}
			// allow multi
			pol.Policies["water"]["CAN_MULTI"] = false

			if ris, err := generateRobotInstructions(lab, inss, pol); err != nil {
				return err
			} else {
				return utils.ErrorSlice{
					assertNumTipsUsed(ris, 1),
					assertNumLoadUnloadInstructions(ris, 1),
				}.Pack()
			}
		},
	})
}

//TestSingleChannelTipReuse2 -- now we move two things
func TestSingleChannelTipReuse2(t *testing.T) {
	testlab.WithTestLab(t, "", &testlab.TestElementCallbacks{
		Steps: func(lab *laboratory.Laboratory) error {
			inss, err := getMixInstructions(lab, 16, []string{components.WaterType}, []float64{50.0})
			if err != nil {
				return err
			}

			ins2, err := getMixInstructions(lab, 8, []string{"ethanol"}, []float64{50.0})
			if err != nil {
				return err
			}

			inss = append(inss, ins2...)

			pol, err := wtype.GetLHPolicyForTest()
			if err != nil {
				return err
			}
			// allow multi
			pol.Policies["water"]["CAN_MULTI"] = false

			if ris, err := generateRobotInstructions(lab, inss, pol); err != nil {
				return err
			} else {
				return utils.ErrorSlice{
					assertNumTipsUsed(ris, 2),
					assertNumLoadUnloadInstructions(ris, 2),
				}.Pack()
			}
		},
	})
}

//TestChannelTipReuseBad Move water and ethanol to two separate columns of wells - should change tips in between
func TestChannelTipReuseBad(t *testing.T) {
	testlab.WithTestLab(t, "", &testlab.TestElementCallbacks{
		Steps: func(lab *laboratory.Laboratory) error {
			inss, err := getMixInstructions(lab, 8, []string{components.WaterType}, []float64{50.0})
			if err != nil {
				return err
			}

			ins2, err := getMixInstructions(lab, 8, []string{"ethanol"}, []float64{50.0})
			if err != nil {
				return err
			}

			inss = append(inss, ins2...)

			if ris, err := generateRobotInstructions(lab, inss, nil); err != nil {
				return err
			} else {
				return utils.ErrorSlice{
					assertNumTipsUsed(ris, 16),
					assertNumLoadUnloadInstructions(ris, 2),
				}.Pack()
			}
		},
	})
}

//TestChannelTipReuseUgly Move water and ethanol to the same columns of wells - should change tips in between
func TestChannelTipReuseUgly(t *testing.T) {
	testlab.WithTestLab(t, "", &testlab.TestElementCallbacks{
		Steps: func(lab *laboratory.Laboratory) error {
			inss, err := getMixInstructions(lab, 8, []string{components.WaterType, "ethanol"}, []float64{50.0, 50.0})
			if err != nil {
				return err
			}

			if ris, err := generateRobotInstructions(lab, inss, nil); err != nil {
				return err
			} else {
				return utils.ErrorSlice{
					assertNumTipsUsed(ris, 16),
					assertNumLoadUnloadInstructions(ris, 2),
				}.Pack()
			}
		},
	})
}

/*
//TestChannelTipReuseUgly Move water and ethanol to the same columns of wells - should change tips in between
func BenchmarkChannelTipReuseUgly(b *testing.B) {
	for i := 0; i < b.N; i++ {
		testlab.WithTestLab(b, "", &testlab.TestElementCallbacks{
			Name: fmt.Sprint(i),
			Steps: func(lab *laboratory.Laboratory) error {
				inss, err := getMixInstructions(lab, 8, []string{components.WaterType, "ethanol"}, []float64{50.0, 50.0})
				if err != nil {
					return err
				}

				return generateRobotInstructions2(lab, inss, nil)
			},
		})
	}
}
*/

func generateRobotInstructions2(lab *laboratory.Laboratory, inss []*wtype.LHInstruction, pol *wtype.LHPolicyRuleSet) ([]liquidhandling.TerminalRobotInstruction, error) {
	tb, dstp := getTransferBlock(lab, inss, "pcrplate_skirted_riser40")

	rbt := getTestRobot(lab, dstp, "pcrplate_skirted_riser40")
	if pol == nil {
		pol, _ = wtype.GetLHPolicyForTest()
		// allow multi
		pol.Policies["water"]["CAN_MULTI"] = true
	}

	//generate the low level instructions
	iTree := liquidhandling.NewITree(tb)

	if _, err := iTree.Build(lab.LaboratoryEffects, pol, rbt); err != nil {
		return nil, err
	}

	return iTree.Leaves(), nil
}

// regression test for issue with additional transfers being
// generated with sequential, different-length multichannel
// operations and non-independent heads
func TestMultiTransferError(t *testing.T) {
	testlab.WithTestLab(t, "", &testlab.TestElementCallbacks{
		Steps: func(lab *laboratory.Laboratory) error {
			// transfer two waters at volumes 50.0, 40.0
			inss, err := getMixInstructions(lab, 2, []string{components.WaterType}, []float64{50.0})

			if err != nil {
				return err
			}

			inss[1].Inputs[0].Vol = 40.0

			tb, p := getTransferBlock(lab, inss, "pcrplate_skirted_riser18")

			rbt := getTestRobot(lab, p, "pcrplate_skirted_riser40")

			pol, err := wtype.GetLHPolicyForTest()
			if err != nil {
				return err
			}

			pol.Policies["water"]["CAN_MULTI"] = true

			ris, err := tb.Generate(lab.LaboratoryEffects, pol, rbt)
			if err != nil {
				return err
			}

			rr, err := ris[0].Generate(lab.LaboratoryEffects, pol, rbt)

			if err != nil {
				return err
			}

			if len(rr) != 1 {
				return fmt.Errorf("Expected 1 instruction got %d", len(rr))
			}

			mcb, ok := rr[0].(*liquidhandling.ChannelBlockInstruction)

			if !ok {
				return fmt.Errorf("Expected *liquidhandling.ChannelBlockInstruction, got %T", rr)
			}

			if len(mcb.What) != 2 {
				return fmt.Errorf("Expected 2 transfers, got %d", len(mcb.What))
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
				return fmt.Errorf("Volumes inconsistent: expected %v got %v", []wunit.Volume{fiftyul, fortyul}, volSums)
			}
			return nil
		},
	})
}

// ensure that gapped transfers are correctly handled
// i.e. [A1:50, B1:40, C1:50] should be done as
//      [A1:40, B1:40, C1:40], [A1:10], [C1:10]
// on a non-independent liquid handler
func TestGappedTransfer(t *testing.T) {
	testlab.WithTestLab(t, "", &testlab.TestElementCallbacks{
		Steps: func(lab *laboratory.Laboratory) error {
			// transfer three waters at volumes 50.0, 40.0, 50.0
			inss, err := getMixInstructions(lab, 3, []string{components.WaterType}, []float64{50.0})

			if err != nil {
				return err
			}

			inss[1].Inputs[0].Vol = 40.0

			tb, p := getTransferBlock(lab, inss, "pcrplate_skirted_riser18")

			rbt := getTestRobot(lab, p, "pcrplate_skirted_riser40")

			pol, err := wtype.GetLHPolicyForTest()
			if err != nil {
				return err
			}

			pol.Policies["water"]["CAN_MULTI"] = true

			ris, err := tb.Generate(lab.LaboratoryEffects, pol, rbt)
			if err != nil {
				return err
			}

			rr, err := ris[0].Generate(lab.LaboratoryEffects, pol, rbt)

			if err != nil {
				return err
			}

			mcb, ok := rr[0].(*liquidhandling.ChannelBlockInstruction)

			if !ok {
				return fmt.Errorf("Expected *liquidhandling.ChannelBlockInstruction, got %T", rr)
			}

			if len(mcb.What) != 3 {
				return fmt.Errorf("Expected 3 transfers in ChannelBlock, got %d", len(mcb.What))
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
				return fmt.Errorf("Volumes inconsistent: expected %v got %v", expect, volSums)
			}
			return nil
		},
	})
}

// another regression test for 2357
// ensure that convertinstructions does not
// modify execution order

// -- the issue here is how getComponents works when getting
// singles, it takes the highest volume first, rather than
// preserving mix order. Feeding in components one at a time should fix this
func TestTransferBlockMixOrdering(t *testing.T) {
	testlab.WithTestLab(t, "", &testlab.TestElementCallbacks{
		Steps: func(lab *laboratory.Laboratory) error {
			// transfer three things at volumes 30.0, 40.0, 50.0 in one instruction
			inss, err := getMixInstructions(lab, 1, []string{"ethanol", "tartrazine", components.WaterType}, []float64{30.0, 40.0, 50.0})

			if err != nil {
				return err
			}

			tb, p := getTransferBlock(lab, inss, "eppendorfrack424_1.5ml_lidholder")

			rbt := getTestRobot(lab, p, "pcrplate_skirted_riser40")

			pol, err := wtype.GetLHPolicyForTest()
			if err != nil {
				return err
			}

			pol.Policies["water"]["CAN_MULTI"] = true
			pol.Policies["ethanol"] = map[string]interface{}{"CAN_MULTI": true}
			pol.Policies["tartrazine"] = map[string]interface{}{"CAN_MULTI": true}

			ris, err := tb.Generate(lab.LaboratoryEffects, pol, rbt)
			if err != nil {
				return err
			}

			rr, err := ris[0].Generate(lab.LaboratoryEffects, pol, rbt)
			if err != nil {
				return err
			}

			mcb, ok := rr[0].(*liquidhandling.ChannelBlockInstruction)

			if !ok {
				return fmt.Errorf("Expected *liquidhandling.ChannelBlockInstruction, got %T", rr)
			}

			if len(mcb.What) != 3 {
				return fmt.Errorf("Expected 3 transfers in ChannelBlock, got %d", len(mcb.What))
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
				return fmt.Errorf("Volumes inconsistent: expected %v got %v", expect, volSums)
			}

			// check ordering, must be preserved

			expectOrder := []string{"ethanol", "tartrazine", components.WaterType}

			gotOrder := make([]string, 3)

			for i := 0; i < len(mcb.What); i++ {
				gotOrder[i] = mcb.Component[i][0]
			}

			if !reflect.DeepEqual(expectOrder, gotOrder) {
				return fmt.Errorf("Order inconsistent: Expected %v got %v", expectOrder, gotOrder)
			}
			return nil
		},
	})
}
