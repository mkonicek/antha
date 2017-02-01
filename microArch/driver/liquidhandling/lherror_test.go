package liquidhandling

import (
	"fmt"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"testing"
)

func TestPolicyError(t *testing.T) {
	rule := wtype.NewLHPolicyRule("test_rule")
	err := rule.AddNumericConditionOn("VOLUME", 32.5, 32.3)
	e1 := "6 (LH_ERR_POLICY) : liquid handling policy error : Numeric condition requested with lower limit (32.500000) greater than upper limit (32.300000), which is not allowed"
	if err == nil {
		t.Fatal(fmt.Sprintf("Expected error %s\nGot nil", e1))
	} else if err.Error() != e1 {
		t.Fatal(fmt.Sprintf("Expected error --%s--\nGot error --%s--\n", e1, err.Error()))
	}
	e2 := "Parameter VOLUME is not categoric"
	err = rule.AddCategoryConditionOn("VOLUME", "")
	if err == nil {
		t.Fatal(fmt.Sprintf("Expected error %s\nGot nil", e2))
	} else if err.Error() != e2 {
		t.Fatal(fmt.Sprintf("Expected error %s\nGot error %s\n", e2, err.Error()))
	}

	e3 := "6 (LH_ERR_POLICY) : liquid handling policy error : Categoric condition  has an empty category, which is not allowed"
	err = rule.AddCategoryConditionOn("LIQUIDCLASS", "")
	if err == nil {
		t.Fatal(fmt.Sprintf("Expected error %s\nGot nil", e2))
	} else if err.Error() != e3 {
		t.Fatal(fmt.Sprintf("Expected error %s\nGot error %s\n", e2, err.Error()))
	}
}
