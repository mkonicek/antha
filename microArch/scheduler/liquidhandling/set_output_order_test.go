package liquidhandling

import (
	"github.com/antha-lang/antha/antha/anthalib/mixer"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"testing"
)

func TestSetOutputOrder(t *testing.T) {
	// we go mix, prompt, split, mix

	rq := GetLHRequestForTest()

	c1 := wtype.NewLHComponent()
	c1.CName = "water"
	c1.Vol = 10

	c2 := wtype.NewLHComponent()
	c2.CName = "washBuffer"
	c2.Vol = 20

	c3 := c1.Dup()
	c3.Mix(c2)

	ins := wtype.NewLHMixInstruction()
	ins.AddComponent(c1)
	ins.AddComponent(c2)
	ins.AddProduct(c3)

	rq.LHInstructions[ins.ID] = ins

	c4 := c3.Cp()
	c4.ParentID = c3.ID
	c3.DaughterID = c4.ID

	prm := wtype.NewLHPromptInstruction()
	prm.Message = "Wait for 10 minutes"
	prm.AddComponent(c3)
	prm.AddResult(c4)

	rq.LHInstructions[prm.ID] = prm

	split := wtype.NewLHSplitInstruction()
	c5, c4a := mixer.SplitSample(c4, wunit.NewVolume(10, "ul"))

	split.AddComponent(c4)
	split.AddProduct(c5)
	split.AddProduct(c4a)

	rq.LHInstructions[split.ID] = split

	mix2 := wtype.NewLHMixInstruction()
	// put them back together again!

	c6 := c4a.Dup()
	c6.Mix(c5)

	mix2.AddComponent(c4a)
	mix2.AddComponent(c5)
	mix2.AddProduct(c6)

	rq.LHInstructions[mix2.ID] = mix2

	set_output_order(rq)

	if rq.InstructionChain.Height() != 4 {
		t.Errorf("Expected instruction chain of length 4, instead got %d", rq.InstructionChain.Height())
	}
}
