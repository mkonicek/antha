package liquidhandling

import (
	"fmt"
	"testing"
)

func TestIsValidLiquidHandlerType(t *testing.T) {

	validTypes := []LiquidHandlerLevel{LLLiquidHandler, HLLiquidHandler}

	for _, c := range validTypes {
		if !c.IsValid() {
			t.Errorf(fmt.Sprintf("error: Type %s must return as valid, returns invalid", c))
		}
	}

	invalidTypes := []LiquidHandlerLevel{157, -65}

	for _, c := range invalidTypes {
		if c.IsValid() {
			t.Errorf(fmt.Sprintf("error: Type %d must return as invalid, returns valid", c))

		}
	}

}

func TestIsValidTipType(t *testing.T) {
	validTypes := []TipType{FixedTips, DisposableTips, MixedDisposableAndFixedTips, NoTips}

	for _, c := range validTypes {
		if !c.IsValid() {
			t.Errorf(fmt.Sprintf("error: Type %s must return as valid, returns invalid", c))
		}
	}

	invalidTypes := []TipType{-2, 78}

	for _, c := range invalidTypes {
		if c.IsValid() {
			t.Errorf(fmt.Sprintf("error: Type %d must return as invalid, returns valid", c))

		}
	}
}
