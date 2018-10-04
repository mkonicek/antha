package liquidhandling

import (
	"fmt"
	"testing"
)

func TestIsValidLiquidHandlerType(t *testing.T) {

	validTypes := []string{LLLiquidHandler, HLLiquidHandler}

	for _, c := range validTypes {
		if !IsValidLiquidHandlerType(c) {
			t.Errorf(fmt.Sprintf("error: Type %s must return as valid, returns invalid", c))
		}
	}

	invalidTypes := []string{"anythingElse", ""}

	for _, c := range invalidTypes {
		if IsValidLiquidHandlerType(c) {
			t.Errorf(fmt.Sprintf("error: Type %s must return as invalid, returns valid", c))

		}
	}

}

func TestIsValidTipType(t *testing.T) {
	validTypes := []string{FixedTips, DisposableTips, MixedDisposableAndFixedTips, NoTips}

	for _, c := range validTypes {
		if !IsValidTipType(c) {
			t.Errorf(fmt.Sprintf("error: Type %s must return as valid, returns invalid", c))
		}
	}

	invalidTypes := []string{"anythingElse", ""}

	for _, c := range invalidTypes {
		if IsValidTipType(c) {
			t.Errorf(fmt.Sprintf("error: Type %s must return as invalid, returns valid", c))

		}
	}
}
