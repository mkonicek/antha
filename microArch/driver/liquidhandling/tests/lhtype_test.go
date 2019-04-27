package tests

import (
	"fmt"
	"testing"

	"github.com/antha-lang/antha/microArch/driver/liquidhandling"
)

func TestIsValidLiquidHandlerType(t *testing.T) {

	validTypes := []liquidhandling.LiquidHandlerLevel{liquidhandling.LLLiquidHandler, liquidhandling.HLLiquidHandler}

	for _, c := range validTypes {
		if !c.IsValid() {
			t.Errorf(fmt.Sprintf("error: Type %s must return as valid, returns invalid", c))
		}
	}

	invalidTypes := []liquidhandling.LiquidHandlerLevel{157, 65}

	for _, c := range invalidTypes {
		if c.IsValid() {
			t.Errorf(fmt.Sprintf("error: Type %d must return as invalid, returns valid", c))

		}
	}

}

func TestIsValidTipType(t *testing.T) {
	validTypes := []liquidhandling.TipType{
		liquidhandling.FixedTips,
		liquidhandling.DisposableTips,
		liquidhandling.MixedDisposableAndFixedTips,
		liquidhandling.NoTips,
	}

	for _, c := range validTypes {
		if !c.IsValid() {
			t.Errorf(fmt.Sprintf("error: Type %s must return as valid, returns invalid", c))
		}
	}

	invalidTypes := []liquidhandling.TipType{12, 78}

	for _, c := range invalidTypes {
		if c.IsValid() {
			t.Errorf(fmt.Sprintf("error: Type %d must return as invalid, returns valid", c))

		}
	}
}
