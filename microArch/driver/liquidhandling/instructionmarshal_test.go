package liquidhandling

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestInstructionMarshal(t *testing.T) {
	arr := []RobotInstruction{
		NewMoveInstruction(),
		NewAspirateInstruction(),
		NewMoveInstruction(),
		NewDispenseInstruction(),
		NewMoveInstruction(),
	}

	marshalled, err := json.Marshal(SetOfRobotInstructions{RobotInstructions: arr})

	if err != nil {
		t.Error(err)
	}

	unmarshalled := SetOfRobotInstructions{}

	err = json.Unmarshal(marshalled, &unmarshalled)

	if err != nil {
		t.Error(err)
	}

	expected := []*InstructionType{MOV, ASP, MOV, DSP, MOV}

	got := insTypeArr(unmarshalled.RobotInstructions)

	if !reflect.DeepEqual(expected, got) {
		t.Errorf("Expected  %v got %v", expected, got)
	}
}
