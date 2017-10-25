package wtype

import (
	"testing"
)

func TestALHPI(t *testing.T) {
	consequences := GetPolicyConsequents()

	alhpi, ok := consequences["PRE_MIX_Z"]

	if !ok {
		t.Error("PRE_MIX_Z not defined")
	} else if alhpi.TypeName() != "float64" {
		t.Errorf("Type of PRE_MIX_Z not as expected: want float64 got %s", alhpi.TypeName())
	}

}
