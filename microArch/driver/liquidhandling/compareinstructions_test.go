package liquidhandling

import (
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"reflect"
	"testing"
)

func testInstructions1() []RobotInstruction {
	_1_mov := NewMoveInstruction()
	_1_mov.OffsetZ = append(_1_mov.OffsetZ, 0.5)
	_1_mov.Reference = append(_1_mov.Reference, 0)
	_1_mov.Pos = append(_1_mov.Pos, "position_1")
	_1_mov.Plt = append(_1_mov.Plt, "DWST12")
	_1_mov.Well = append(_1_mov.Well, "A1")

	_2_asp := NewAspirateInstruction()
	_2_asp.What = append(_2_asp.What, "string")
	_2_asp.Volume = append(_2_asp.Volume, wunit.NewVolume(100.0, "ul"))

	_3_mov := NewMoveInstruction()
	_3_mov.OffsetX = append(_3_mov.OffsetZ, 1.5)
	_3_mov.Reference = append(_3_mov.Reference, 1)
	_3_mov.Pos = append(_3_mov.Pos, "position_2")
	_3_mov.Plt = append(_3_mov.Plt, "DSW96")
	_3_mov.Well = append(_3_mov.Well, "C3")

	_4_dsp := NewDispenseInstruction()
	_4_dsp.What = append(_4_dsp.What, "string")
	_4_dsp.Volume = append(_4_dsp.Volume, wunit.NewVolume(100.0, "ul"))

	_5_mov := NewMoveInstruction()
	_5_mov.OffsetX = append(_5_mov.OffsetZ, 0.5)
	_5_mov.Reference = append(_5_mov.Reference, 0)
	_5_mov.Pos = append(_5_mov.Pos, "position_2")
	_5_mov.Plt = append(_5_mov.Plt, "DSW96")
	_5_mov.Well = append(_5_mov.Well, "C3")

	_6_mix := NewMixInstruction()
	_6_mix.What = append(_6_mix.What, "string")
	_6_mix.Volume = append(_6_mix.Volume, wunit.NewVolume(100.0, "ul"))

	_7_mov := NewMoveInstruction()
	_8_mov := NewMoveInstruction()

	return []RobotInstruction{_1_mov, _2_asp, _3_mov, _4_dsp, _5_mov, _6_mix, _7_mov, _8_mov}
}

func testInstructions2() []RobotInstruction {

	_1_mov := NewMoveInstruction()
	_1_mov.OffsetZ = append(_1_mov.OffsetZ, 0.5)
	_1_mov.Reference = append(_1_mov.Reference, 0)
	_1_mov.Pos = append(_1_mov.Pos, "position_2")
	_1_mov.Plt = append(_1_mov.Plt, "DWST12")
	_1_mov.Well = append(_1_mov.Well, "A1")

	_2_asp := NewAspirateInstruction()
	_2_asp.What = append(_2_asp.What, "string")
	_2_asp.Volume = append(_2_asp.Volume, wunit.NewVolume(50.0, "ul"))

	_3_mov := NewMoveInstruction()
	_3_mov.OffsetX = append(_3_mov.OffsetZ, 1.5)
	_3_mov.Reference = append(_3_mov.Reference, 1)
	_3_mov.Pos = append(_3_mov.Pos, "position_1")
	_3_mov.Plt = append(_3_mov.Plt, "DSW96")
	_3_mov.Well = append(_3_mov.Well, "C3")

	_4_dsp := NewDispenseInstruction()
	_4_dsp.What = append(_4_dsp.What, "string")
	_4_dsp.Volume = append(_4_dsp.Volume, wunit.NewVolume(50.0, "ul"))

	_5_mov := NewMoveInstruction()
	_5_mov.OffsetX = append(_5_mov.OffsetZ, 0.5)
	_5_mov.Reference = append(_5_mov.Reference, 0)
	_5_mov.Pos = append(_5_mov.Pos, "position_1")
	_5_mov.Plt = append(_5_mov.Plt, "DSW96")
	_5_mov.Well = append(_5_mov.Well, "C3")

	_6_mix := NewMixInstruction()
	_6_mix.What = append(_6_mix.What, "string")
	_6_mix.Volume = append(_6_mix.Volume, wunit.NewVolume(50.0, "ul"))

	_7_mov := NewMoveInstruction()

	_8_mov := NewMoveInstruction()

	return []RobotInstruction{_1_mov, _2_asp, _3_mov, _4_dsp, _5_mov, _6_mix, _7_mov, _8_mov}
}

func insTypeArr(arr []RobotInstruction) []*InstructionType {
	ret := make([]*InstructionType, len(arr))
	for i, inst := range arr {
		ret[i] = inst.Type()
	}
	return ret
}

func TestMerger(t *testing.T) {
	inss := testInstructions1()
	merged := mergeMovs(inss)

	expected := []*InstructionType{MAS, MDS, MVM, MOV, MOV}
	got := insTypeArr(merged)

	if !reflect.DeepEqual(expected, got) {
		t.Errorf("Merge fail: expected %v got %v", expected, got)
	}
}

func TestErroneousComparisons(t *testing.T) {
	nilIns := []RobotInstruction{nil}
	if err := CompareInstructionSets(nilIns, nilIns); len(err) == 0 {
		t.Error("Expected error when comparing nil instruction sets")
	}

	inss1 := testInstructions1()
	if err := CompareInstructionSets(inss1, nil); len(err) == 0 {
		t.Error("Expected error when comparing non-nil with nil instruction sets")
	}

	if err := CompareInstructionSets(nil, inss1); len(err) == 0 {
		t.Error("Expected error when comparing nil with non-nil instruction sets")
	}
}

func TestDifferentTypes(t *testing.T) {
	as := []RobotInstruction{NewMoveInstruction(), NewMixInstruction(), NewMixInstruction()}
	bs := []RobotInstruction{NewMixInstruction(), NewMixInstruction(), NewMoveInstruction()}
	if err := CompareInstructionSets(as, bs); len(err) != 2 {
		t.Errorf("Expected two errors when comparing %v with %v but got %v instead.", as, bs, err)
	}
}

func TestBasicComparison(t *testing.T) {
	inss1 := testInstructions1()
	inss2 := testInstructions1()

	ret := CompareInstructionSets(inss1, inss2, CompareAllParameters...)

	if len(ret) != 0 {
		t.Errorf("Expected 0 errors, got %d: %v", len(ret), ret)
	}

	inss2 = testInstructions2()

	ret = CompareInstructionSets(inss1, inss2, CompareAllParameters...)

	if len(ret) != 6 {
		t.Errorf("Expected 6 errors, got %d: %v", len(ret), ret)
	}
}

func TestComparisonoptions(t *testing.T) {
	inss1 := testInstructions1()
	inss2 := testInstructions2()

	expected := []int{0, 0, 0, 0, 3, 3}
	prms := [][]RobotInstructionComparatorFunc{CompareReferences, CompareOffsets, CompareWells, ComparePlateTypes, {CompareVolumes}, ComparePositions}
	names := []string{"References", "Offsets", "Wells", "PlateTypes", "Volumes", "Positions"}

	for i := 0; i < len(expected); i++ {
		ret := CompareInstructionSets(inss1, inss2, prms[i]...)
		if len(ret) != expected[i] {
			t.Errorf("Comparison %s expected %d errors got %d (%v)", names[i], expected[i], len(ret), ret)
		}
	}

}
