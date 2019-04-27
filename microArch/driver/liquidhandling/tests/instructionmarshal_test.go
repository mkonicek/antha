package tests

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/antha-lang/antha/microArch/driver/liquidhandling"
)

func TestInstructionMarshal(t *testing.T) {
	arr := []liquidhandling.RobotInstruction{
		liquidhandling.NewMoveInstruction(),
		liquidhandling.NewAspirateInstruction(),
		liquidhandling.NewMoveInstruction(),
		liquidhandling.NewDispenseInstruction(),
		liquidhandling.NewMoveInstruction(),
	}

	marshalled, err := json.Marshal(liquidhandling.SetOfRobotInstructions{RobotInstructions: arr})

	if err != nil {
		t.Error(err)
	}

	unmarshalled := liquidhandling.SetOfRobotInstructions{}

	err = json.Unmarshal(marshalled, &unmarshalled)

	if err != nil {
		t.Error(err)
	}

	expected := []*liquidhandling.InstructionType{liquidhandling.MOV, liquidhandling.ASP, liquidhandling.MOV, liquidhandling.DSP, liquidhandling.MOV}

	got := insTypeArr(unmarshalled.RobotInstructions)

	if !reflect.DeepEqual(expected, got) {
		t.Errorf("Expected  %v got %v", expected, got)
	}
}
