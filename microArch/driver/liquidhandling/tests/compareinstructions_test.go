package tests

import (
	"reflect"
	"testing"

	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/microArch/driver/liquidhandling"
)

func testInstructions1() []liquidhandling.RobotInstruction {
	_1_mov := liquidhandling.NewMoveInstruction()
	_1_mov.OffsetZ = append(_1_mov.OffsetZ, 0.5)
	_1_mov.Reference = append(_1_mov.Reference, 0)
	_1_mov.Pos = append(_1_mov.Pos, "position_1")
	_1_mov.Plt = append(_1_mov.Plt, "DWST12")
	_1_mov.Well = append(_1_mov.Well, "A1")

	_2_asp := liquidhandling.NewAspirateInstruction()
	_2_asp.What = append(_2_asp.What, "string")
	_2_asp.Volume = append(_2_asp.Volume, wunit.NewVolume(100.0, "ul"))

	_3_mov := liquidhandling.NewMoveInstruction()
	_3_mov.OffsetX = append(_3_mov.OffsetZ, 1.5)
	_3_mov.Reference = append(_3_mov.Reference, 1)
	_3_mov.Pos = append(_3_mov.Pos, "position_2")
	_3_mov.Plt = append(_3_mov.Plt, "DSW96")
	_3_mov.Well = append(_3_mov.Well, "C3")

	_4_dsp := liquidhandling.NewDispenseInstruction()
	_4_dsp.What = append(_4_dsp.What, "string")
	_4_dsp.Volume = append(_4_dsp.Volume, wunit.NewVolume(100.0, "ul"))

	_5_mov := liquidhandling.NewMoveInstruction()
	_5_mov.OffsetX = append(_5_mov.OffsetZ, 0.5)
	_5_mov.Reference = append(_5_mov.Reference, 0)
	_5_mov.Pos = append(_5_mov.Pos, "position_2")
	_5_mov.Plt = append(_5_mov.Plt, "DSW96")
	_5_mov.Well = append(_5_mov.Well, "C3")

	_6_mix := liquidhandling.NewMixInstruction()
	_6_mix.What = append(_6_mix.What, "string")
	_6_mix.Volume = append(_6_mix.Volume, wunit.NewVolume(100.0, "ul"))

	_7_mov := liquidhandling.NewMoveInstruction()
	_8_mov := liquidhandling.NewMoveInstruction()

	return []liquidhandling.RobotInstruction{_1_mov, _2_asp, _3_mov, _4_dsp, _5_mov, _6_mix, _7_mov, _8_mov}
}

func testInstructions2() []liquidhandling.RobotInstruction {

	_1_mov := liquidhandling.NewMoveInstruction()
	_1_mov.OffsetZ = append(_1_mov.OffsetZ, 0.5)
	_1_mov.Reference = append(_1_mov.Reference, 0)
	_1_mov.Pos = append(_1_mov.Pos, "position_2")
	_1_mov.Plt = append(_1_mov.Plt, "DWST12")
	_1_mov.Well = append(_1_mov.Well, "A1")

	_2_asp := liquidhandling.NewAspirateInstruction()
	_2_asp.What = append(_2_asp.What, "string")
	_2_asp.Volume = append(_2_asp.Volume, wunit.NewVolume(50.0, "ul"))

	_3_mov := liquidhandling.NewMoveInstruction()
	_3_mov.OffsetX = append(_3_mov.OffsetZ, 1.5)
	_3_mov.Reference = append(_3_mov.Reference, 1)
	_3_mov.Pos = append(_3_mov.Pos, "position_1")
	_3_mov.Plt = append(_3_mov.Plt, "DSW96")
	_3_mov.Well = append(_3_mov.Well, "C3")

	_4_dsp := liquidhandling.NewDispenseInstruction()
	_4_dsp.What = append(_4_dsp.What, "string")
	_4_dsp.Volume = append(_4_dsp.Volume, wunit.NewVolume(50.0, "ul"))

	_5_mov := liquidhandling.NewMoveInstruction()
	_5_mov.OffsetX = append(_5_mov.OffsetZ, 0.5)
	_5_mov.Reference = append(_5_mov.Reference, 0)
	_5_mov.Pos = append(_5_mov.Pos, "position_1")
	_5_mov.Plt = append(_5_mov.Plt, "DSW96")
	_5_mov.Well = append(_5_mov.Well, "C3")

	_6_mix := liquidhandling.NewMixInstruction()
	_6_mix.What = append(_6_mix.What, "string")
	_6_mix.Volume = append(_6_mix.Volume, wunit.NewVolume(50.0, "ul"))

	_7_mov := liquidhandling.NewMoveInstruction()

	_8_mov := liquidhandling.NewMoveInstruction()

	return []liquidhandling.RobotInstruction{_1_mov, _2_asp, _3_mov, _4_dsp, _5_mov, _6_mix, _7_mov, _8_mov}
}

func insTypeArr(arr []liquidhandling.RobotInstruction) []*liquidhandling.InstructionType {
	ret := make([]*liquidhandling.InstructionType, len(arr))
	for i, inst := range arr {
		ret[i] = inst.Type()
	}
	return ret
}

func TestMerger(t *testing.T) {
	inss := testInstructions1()
	merged := liquidhandling.MergeMovs(inss)

	expected := []*liquidhandling.InstructionType{liquidhandling.MAS, liquidhandling.MDS, liquidhandling.MVM, liquidhandling.MOV, liquidhandling.MOV}
	got := insTypeArr(merged)

	if !reflect.DeepEqual(expected, got) {
		t.Errorf("Merge fail: expected %v got %v", expected, got)
	}
}

func TestErroneousComparisons(t *testing.T) {
	nilIns := []liquidhandling.RobotInstruction{nil}
	if err := liquidhandling.CompareInstructionSets(nilIns, nilIns); len(err) == 0 {
		t.Error("Expected error when comparing nil instruction sets")
	}

	inss1 := testInstructions1()
	if err := liquidhandling.CompareInstructionSets(inss1, nil); len(err) == 0 {
		t.Error("Expected error when comparing non-nil with nil instruction sets")
	}

	if err := liquidhandling.CompareInstructionSets(nil, inss1); len(err) == 0 {
		t.Error("Expected error when comparing nil with non-nil instruction sets")
	}
}

func TestDifferentTypes(t *testing.T) {
	as := []liquidhandling.RobotInstruction{liquidhandling.NewMoveInstruction(), liquidhandling.NewMixInstruction(), liquidhandling.NewMixInstruction()}
	bs := []liquidhandling.RobotInstruction{liquidhandling.NewMixInstruction(), liquidhandling.NewMixInstruction(), liquidhandling.NewMoveInstruction()}
	if err := liquidhandling.CompareInstructionSets(as, bs); len(err) != 2 {
		t.Errorf("Expected two errors when comparing %v with %v but got %v instead.", as, bs, err)
	}
}

func TestBasicComparison(t *testing.T) {
	inss1 := testInstructions1()
	inss2 := testInstructions1()

	ret := liquidhandling.CompareInstructionSets(inss1, inss2, liquidhandling.CompareAllParameters...)

	if len(ret) != 0 {
		t.Errorf("Expected 0 errors, got %d: %v", len(ret), ret)
	}

	inss2 = testInstructions2()

	ret = liquidhandling.CompareInstructionSets(inss1, inss2, liquidhandling.CompareAllParameters...)

	if len(ret) != 6 {
		t.Errorf("Expected 6 errors, got %d: %v", len(ret), ret)
	}
}

func TestComparisonoptions(t *testing.T) {
	inss1 := testInstructions1()
	inss2 := testInstructions2()

	expected := []int{0, 0, 0, 0, 3, 3}
	prms := [][]liquidhandling.RobotInstructionComparatorFunc{
		liquidhandling.CompareReferences,
		liquidhandling.CompareOffsets,
		liquidhandling.CompareWells,
		liquidhandling.ComparePlateTypes,
		{liquidhandling.CompareVolumes},
		liquidhandling.ComparePositions,
	}
	names := []string{"References", "Offsets", "Wells", "PlateTypes", "Volumes", "Positions"}

	for i := 0; i < len(expected); i++ {
		ret := liquidhandling.CompareInstructionSets(inss1, inss2, prms[i]...)
		if len(ret) != expected[i] {
			t.Errorf("Comparison %s expected %d errors got %d (%v)", names[i], expected[i], len(ret), ret)
		}
	}

}
