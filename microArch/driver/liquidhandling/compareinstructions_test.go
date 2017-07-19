package liquidhandling

import (
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"testing"
)

func testInstructions1() []RobotInstruction {
	insS := make([]RobotInstruction, 0, 1)
	var ins RobotInstruction
	ins = NewMoveInstruction()
	ins.(*MoveInstruction).OffsetZ = append(ins.(*MoveInstruction).OffsetZ, 0.5)
	ins.(*MoveInstruction).Reference = append(ins.(*MoveInstruction).Reference, 0)
	ins.(*MoveInstruction).Pos = append(ins.(*MoveInstruction).Pos, "position_1")
	ins.(*MoveInstruction).Plt = append(ins.(*MoveInstruction).Plt, "DWST12")
	ins.(*MoveInstruction).Well = append(ins.(*MoveInstruction).Well, "A1")
	insS = append(insS, ins)
	ins = NewAspirateInstruction()
	ins.(*AspirateInstruction).What = append(ins.(*AspirateInstruction).What, "string")
	ins.(*AspirateInstruction).Volume = append(ins.(*AspirateInstruction).Volume, wunit.NewVolume(100.0, "ul"))
	insS = append(insS, ins)
	ins = NewMoveInstruction()
	ins.(*MoveInstruction).OffsetX = append(ins.(*MoveInstruction).OffsetZ, 1.5)
	ins.(*MoveInstruction).Reference = append(ins.(*MoveInstruction).Reference, 1)
	ins.(*MoveInstruction).Pos = append(ins.(*MoveInstruction).Pos, "position_2")
	ins.(*MoveInstruction).Plt = append(ins.(*MoveInstruction).Plt, "DSW96")
	ins.(*MoveInstruction).Well = append(ins.(*MoveInstruction).Well, "C3")
	insS = append(insS, ins)
	ins = NewDispenseInstruction()
	ins.(*DispenseInstruction).What = append(ins.(*DispenseInstruction).What, "string")
	ins.(*DispenseInstruction).Volume = append(ins.(*DispenseInstruction).Volume, wunit.NewVolume(100.0, "ul"))
	insS = append(insS, ins)
	ins = NewMoveInstruction()
	ins.(*MoveInstruction).OffsetX = append(ins.(*MoveInstruction).OffsetZ, 0.5)
	ins.(*MoveInstruction).Reference = append(ins.(*MoveInstruction).Reference, 1)
	ins.(*MoveInstruction).Pos = append(ins.(*MoveInstruction).Pos, "position_2")
	ins.(*MoveInstruction).Plt = append(ins.(*MoveInstruction).Plt, "DSW96")
	ins.(*MoveInstruction).Well = append(ins.(*MoveInstruction).Well, "C3")
	insS = append(insS, ins)
	ins = NewMixInstruction()
	ins.(*MixInstruction).What = append(ins.(*MixInstruction).What, "string")
	ins.(*MixInstruction).Volume = append(ins.(*MixInstruction).Volume, wunit.NewVolume(100.0, "ul"))
	insS = append(insS, ins)
	ins = NewMoveInstruction()
	insS = append(insS, ins)
	ins = NewMoveInstruction()
	insS = append(insS, ins)
	return insS
}

func testInstructions2() []RobotInstruction {
	insS := make([]RobotInstruction, 0, 1)
	var ins RobotInstruction
	ins = NewMoveInstruction()
	ins.(*MoveInstruction).OffsetZ = append(ins.(*MoveInstruction).OffsetZ, 0.5)
	ins.(*MoveInstruction).Reference = append(ins.(*MoveInstruction).Reference, 0)
	ins.(*MoveInstruction).Pos = append(ins.(*MoveInstruction).Pos, "position_2")
	ins.(*MoveInstruction).Plt = append(ins.(*MoveInstruction).Plt, "DWST12")
	ins.(*MoveInstruction).Well = append(ins.(*MoveInstruction).Well, "A1")
	insS = append(insS, ins)
	ins = NewAspirateInstruction()
	ins.(*AspirateInstruction).What = append(ins.(*AspirateInstruction).What, "string")
	ins.(*AspirateInstruction).Volume = append(ins.(*AspirateInstruction).Volume, wunit.NewVolume(50.0, "ul"))
	insS = append(insS, ins)
	ins = NewMoveInstruction()
	ins.(*MoveInstruction).OffsetX = append(ins.(*MoveInstruction).OffsetZ, 1.5)
	ins.(*MoveInstruction).Reference = append(ins.(*MoveInstruction).Reference, 1)
	ins.(*MoveInstruction).Pos = append(ins.(*MoveInstruction).Pos, "position_1")
	ins.(*MoveInstruction).Plt = append(ins.(*MoveInstruction).Plt, "DSW96")
	ins.(*MoveInstruction).Well = append(ins.(*MoveInstruction).Well, "C3")
	insS = append(insS, ins)
	ins = NewDispenseInstruction()
	ins.(*DispenseInstruction).What = append(ins.(*DispenseInstruction).What, "string")
	ins.(*DispenseInstruction).Volume = append(ins.(*DispenseInstruction).Volume, wunit.NewVolume(50.0, "ul"))
	insS = append(insS, ins)
	ins = NewMoveInstruction()
	ins.(*MoveInstruction).OffsetX = append(ins.(*MoveInstruction).OffsetZ, 0.5)
	ins.(*MoveInstruction).Reference = append(ins.(*MoveInstruction).Reference, 1)
	ins.(*MoveInstruction).Pos = append(ins.(*MoveInstruction).Pos, "position_1")
	ins.(*MoveInstruction).Plt = append(ins.(*MoveInstruction).Plt, "DSW96")
	ins.(*MoveInstruction).Well = append(ins.(*MoveInstruction).Well, "C3")
	insS = append(insS, ins)
	ins = NewMixInstruction()
	ins.(*MixInstruction).What = append(ins.(*MixInstruction).What, "string")
	ins.(*MixInstruction).Volume = append(ins.(*MixInstruction).Volume, wunit.NewVolume(50.0, "ul"))
	insS = append(insS, ins)
	ins = NewMoveInstruction()
	insS = append(insS, ins)
	ins = NewMoveInstruction()
	insS = append(insS, ins)
	return insS
}

func stringArrsSame(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}

	for i := 0; i < len(a); i++ {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func insTypeArr(arr []RobotInstruction) []string {
	ret := make([]string, len(arr))
	for i := 0; i < len(arr); i++ {
		ret[i] = InstructionTypeName(arr[i])
	}
	return ret
}

func TestMerger(t *testing.T) {
	inss := testInstructions1()
	merged := mergeMovs(inss)

	expected := []string{"MOVASP", "MOVDSP", "MOVMIX", "MOV", "MOV"}
	got := insTypeArr(merged)

	if !stringArrsSame(expected, got) {
		t.Errorf("Merge fail: expected %v got %v", expected, got)
	}
}

func TestBasicComparison(t *testing.T) {
	inss1 := testInstructions1()
	inss2 := testInstructions1()

	ret := CompareInstructionSets(inss1, inss2, ComparisonOpt{InstructionParameters: CompareAll()})

	if len(ret.Errors) != 0 {
		t.Errorf("Expected 0 errors, got %d: %v", len(ret.Errors), ret.Errors)
	}

	inss2 = testInstructions2()

	ret = CompareInstructionSets(inss1, inss2, ComparisonOpt{InstructionParameters: CompareAll()})

	if len(ret.Errors) != 5 {
		t.Errorf("Expected 5 errors, got %d: %v", len(ret.Errors), ret.Errors)
	}
}
