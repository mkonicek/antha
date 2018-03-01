package liquidhandling

import (
	"fmt"
	"github.com/antha-lang/antha/antha/anthalib/mixer"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"reflect"
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

	expectedIDOrder := make([]string, 0, 5)

	mix1 := wtype.NewLHMixInstruction()
	mix1.AddComponent(c1)
	mix1.AddComponent(c2)
	mix1.AddProduct(c3)

	expectedIDOrder = append(expectedIDOrder, mix1.ID)

	rq.LHInstructions[mix1.ID] = mix1

	c4 := c3.Cp()
	c4.ParentID = c3.ID
	c3.DaughterID = c4.ID

	prm := wtype.NewLHPromptInstruction()
	prm.Message = "Wait for 10 minutes"
	prm.AddComponent(c3)
	prm.AddResult(c4)

	expectedIDOrder = append(expectedIDOrder, "PROMPT") // ID changes but there's only one

	rq.LHInstructions[prm.ID] = prm

	split := wtype.NewLHSplitInstruction()
	c5, c4a := mixer.SplitSample(c4, wunit.NewVolume(10, "ul"))

	split.AddComponent(c4)
	split.AddProduct(c5)
	split.AddProduct(c4a)

	rq.LHInstructions[split.ID] = split

	mix2 := wtype.NewLHMixInstruction()

	c6 := wtype.NewLHComponent()
	c6.CName = "turps"
	c6.Vol = 100

	// mix the sample with c6

	c7 := c6.Dup()
	c7.Mix(c5)

	mix2.AddComponent(c6)
	mix2.AddComponent(c5)
	mix2.AddProduct(c7)

	expectedIDOrder = append(expectedIDOrder, mix2.ID)
	expectedIDOrder = append(expectedIDOrder, split.ID)

	rq.LHInstructions[mix2.ID] = mix2

	// mix the static component with some more water

	c8 := wtype.NewLHComponent()
	c8.CName = "water"
	c8.Vol = 200

	c9 := c4a.Dup()
	c9.Mix(c8)

	mix3 := wtype.NewLHMixInstruction()
	mix3.AddComponent(c4a)
	mix3.AddComponent(c8)
	mix3.AddProduct(c9)

	expectedIDOrder = append(expectedIDOrder, mix3.ID)

	rq.LHInstructions[mix3.ID] = mix3

	set_output_order(rq)

	if rq.InstructionChain.Height() != 5 {
		t.Errorf("Expected instruction chain of length 5, instead got %d", rq.InstructionChain.Height())
	}

	gotIDOrder := flatten(rq.InstructionChain)

	if !reflect.DeepEqual(expectedIDOrder, gotIDOrder) {
		t.Errorf(fmt.Sprintf("Expected %v got %v", expectedIDOrder, gotIDOrder))
	}

}

func flatten(ic *IChain) []string {
	if ic == nil {
		return []string{}
	}

	r := make([]string, 0, len(ic.Values))

	for _, v := range ic.Values {
		if v.Type == wtype.LHIPRM {
			r = append(r, "PROMPT")
		} else {
			r = append(r, v.ID)
		}
	}

	return append(r, flatten(ic.Child)...)
}
