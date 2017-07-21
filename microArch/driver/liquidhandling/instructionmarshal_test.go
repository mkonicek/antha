package liquidhandling

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
)

func TestInstructionMarshal(t *testing.T) {
	arr := make([]RobotInstruction, 0, 1)
	arr = append(arr, NewMoveInstruction())
	arr = append(arr, NewAspirateInstruction())
	arr = append(arr, NewMoveInstruction())
	arr = append(arr, NewDispenseInstruction())
	arr = append(arr, NewMoveInstruction())

	marshalled, err := json.Marshal(SetOfRobotInstructions{Instructions: arr})

	if err != nil {
		t.Error(err)
	}

	unmarshalled := SetOfRobotInstructions{}

	err = json.Unmarshal(marshalled, &unmarshalled)

	if err != nil {
		t.Error(err)
	}

	expected := []string{"MOV", "ASP", "MOV", "DSP", "MOV"}

	got := insTypeArr(unmarshalled.Instructions)

	if !reflect.DeepEqual(expected, got) {
		fmt.Errorf("Expected  %v got %v", expected, got)
	}
}
