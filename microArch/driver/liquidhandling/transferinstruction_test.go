package liquidhandling

import (
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/inventory"
	"testing"
)

// Conceived as a test of Antha-2357
func TestOrdering(t *testing.T) {
	ctx := GetContextForTest()
	dstp, _ := inventory.NewPlate(ctx, "DWST12")
	rbt := getTestRobot(ctx, dstp, "pcrplate_skirted_riser40")
	pol, err := wtype.GetLHPolicyForTest()

	if err != nil {
		t.Errorf(err.Error())
	}

	ins := getTransferMulti()

	ris := NewRobotInstructionSet(ins)

	var inss []RobotInstruction

	inss, err = ris.Generate(ctx, pol, rbt)

	if err != nil {
		t.Errorf(err.Error())
	}

	if len(inss) != 58 {
		t.Errorf("Expected 58 instructions, got %d", len(inss))
	}

	if inss[4].Type() != MOV {
		t.Errorf("Wrong type for instruction 5: %s", inss[4].Type())
	}

	wls := inss[4].GetParameter(WELLTO)

	if wls == nil {
		t.Errorf("Instruction 5 returned nil for WELLTO")
	}

	strWells, ok := wls.([]string)

	if !ok {
		t.Errorf("Expected []string type for WELLTO")
	}

	if strWells[0] != "A1" {
		t.Errorf("Expected first transfer from well A1 (water), instead got %s", strWells[0])
	}
}

// the idea here is that water should be moved before wine but
// wine multichannels and water doesn't.
func getTransferMulti() RobotInstruction {
	vol := wunit.NewVolume(100.0, "ul")
	v2 := wunit.NewVolume(5000.0, "ul")
	v3 := wunit.NewVolume(0.0, "ul")
	tfr := NewTransferInstruction(
		[]string{"water", "water", "water"},
		[]string{"position_4", "position_4", "position_4"},
		[]string{"position_8", "position_8", "position_8"},
		[]string{"A1", "A2", "A3"},
		[]string{"A1", "B1", "C1"},
		[]string{"DWST12", "DWST12", "DWST12"},
		[]string{"DSW96", "DSW96", "DSW96"},
		[]wunit.Volume{vol.Dup(), vol.Dup(), vol.Dup()},
		[]wunit.Volume{v2.Dup(), v2.Dup(), v2.Dup()},
		[]wunit.Volume{v3.Dup(), v3.Dup(), v3.Dup()},
		[]int{8, 8, 8},
		[]int{12, 12, 12},
		[]int{8, 8, 8},
		[]int{12, 12, 12},
		[]string{"water", "water", "water"},
		[]wtype.LHPolicy{nil, nil, nil},
	)

	tfr.Add(MTPFromArrays(
		[]string{"solvent", "solvent", "solvent"},
		[]string{"position_4", "position_4", "position_4"},
		[]string{"position_8", "position_8", "position_8"},
		[]string{"A4", "A4", "A4"},
		[]string{"A1", "B1", "C1"},
		[]string{"DWST12", "DWST12", "DWST12"},
		[]string{"DSW96", "DSW96", "DSW96"},
		[]wunit.Volume{vol.Dup(), vol.Dup(), vol.Dup()},
		[]wunit.Volume{v2.Dup(), v2.Dup(), v2.Dup()},
		[]wunit.Volume{v3.Dup(), v3.Dup(), v3.Dup()},
		[]int{8, 8, 8},
		[]int{12, 12, 12},
		[]int{8, 8, 8},
		[]int{12, 12, 12},
		[]string{"wine", "wine", "wine"},
		[]wtype.LHPolicy{nil, nil, nil},
	))

	return tfr
}
