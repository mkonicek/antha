package liquidhandling

import (
	"context"
	"fmt"
	"testing"

	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/inventory/testinventory"
)

func getTestBlowout(robot *LHProperties) RobotInstruction {
	v := wunit.NewVolume(10.0, "ul")
	ch, _, _ := ChooseChannel(v, robot)
	bi := NewBlowInstruction()
	bi.Multi = 1
	bi.What = append(bi.What, "soup")
	bi.PltTo = append(bi.PltTo, "position_4")
	bi.WellTo = append(bi.WellTo, "A1")
	bi.Volume = append(bi.Volume, v)
	bi.TPlateType = append(bi.TPlateType, "pcrplate_skirted_riser40")
	bi.TVolume = append(bi.TVolume, wunit.NewVolume(5.0, "ul"))
	bi.Prms = ch
	bi.Head = ch.Head
	return bi
}

func TestBlowWithTipChange(t *testing.T) {
	t.Skip()

	ctx := testinventory.NewContext(context.Background())
	robot, err := MakeGilsonWithPlatesForTest(ctx)
	if err != nil {
		t.Fatal(err)
	}

	bi := getTestBlowout(robot)
	pol, _ := wtype.GetLHPolicyForTest()

	rule := wtype.NewLHPolicyRule("TESTRULE1")
	rule.AddCategoryConditionOn("LIQUIDCLASS", "soup")
	pols := make(wtype.LHPolicy, 2)
	pols["POST_MIX"] = 5
	pols["POST_MIX_VOLUME"] = 100.0
	pol.AddRule(rule, pols)
	set := NewRobotInstructionSet(bi)

	ris, err := set.Generate(ctx, pol, robot)

	if err != nil {
		t.Fatal(err)
	}

	expectedIns := []int{MOV, DSP, MOV, ULD, MOV, LOD, MOV, MMX, MOV, BLO}

	if len(ris) != len(expectedIns) {
		t.Fatal(fmt.Sprintf("Error: Expected %d instructions, got %d", len(expectedIns), len(ris)))
	}

	for i, ins := range ris {
		if ins.InstructionType() != expectedIns[i] {
			t.Fatal(fmt.Sprintf("Error generating high mix volume blow: expected %s got %s", Robotinstructionnames[expectedIns[i]], Robotinstructionnames[ins.InstructionType()]))
		}
	}
}

func TestBlowNoTipChange(t *testing.T) {
	ctx := testinventory.NewContext(context.Background())

	robot, err := MakeGilsonWithPlatesForTest(ctx)
	if err != nil {
		t.Fatal(err)
	}

	bi := getTestBlowout(robot)
	pol, _ := wtype.GetLHPolicyForTest()

	rule := wtype.NewLHPolicyRule("TESTRULE1")
	rule.AddCategoryConditionOn("LIQUIDCLASS", "soup")
	soupPolicy := make(wtype.LHPolicy, 2)
	soupPolicy["POST_MIX"] = 5
	soupPolicy["POST_MIX_VOLUME"] = 10.0
	pol.Policies["soup"] = soupPolicy
	pol.AddRule(rule, soupPolicy)
	set := NewRobotInstructionSet(bi)

	ris, err := set.Generate(ctx, pol, robot)

	if err != nil {
		t.Fatal(err)
	}
	expectedIns := []int{MOV, DSP, MOV, MIX, MOV, BLO}

	if len(ris) != len(expectedIns) {
		t.Fatal(fmt.Sprintf("Error: Expected %d instructions, got %d", len(expectedIns), len(ris)))
	}

	for i, ins := range ris {
		if ins.InstructionType() != expectedIns[i] {
			t.Fatal(fmt.Sprintf("Error generating low mix volume blow: expected %s got %s", Robotinstructionnames[expectedIns[i]], Robotinstructionnames[ins.InstructionType()]))
		}
	}
}
