package liquidhandling

import (
	"github.com/antha-lang/antha/antha/anthalib/mixer"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"testing"
)

func GetMixForTest(id string, input ...*wtype.Liquid) (*wtype.LHInstruction, *wtype.Liquid) {
	output := wtype.NewLHComponent()
	output.Type = wtype.LTNIL
	for _, ip := range input {
		output.Mix(ip)
	}

	mix := wtype.NewLHMixInstruction()
	mix.ID = id
	for _, ip := range input {
		mix.AddComponent(ip)
	}
	mix.AddProduct(output)

	return mix, output
}

func GetSplitForTest(id string, input *wtype.Liquid, volume float64) (*wtype.LHInstruction, *wtype.Liquid, *wtype.Liquid) {
	split := wtype.NewLHSplitInstruction()
	split.ID = id
	moving, remaining := mixer.SplitSample(input, wunit.NewVolume(volume, "ul"))

	split.AddComponent(input)
	split.AddProduct(moving)
	split.AddProduct(remaining)

	return split, moving, remaining
}

func GetPromptForTest(message string, inputs ...*wtype.Liquid) (*wtype.LHInstruction, []*wtype.Liquid) {
	ret := wtype.NewLHPromptInstruction()
	ret.ID = message //ID will be overwritten, set the message as well for testing
	ret.Message = message
	for _, input := range inputs {
		output := input.Cp()
		output.ParentID = input.ID
		input.DaughterID = output.ID
		ret.AddComponent(input)
		ret.AddResult(output)
	}

	return ret, ret.Results
}

func GetLiquidForTest(name string, volume float64) *wtype.Liquid {
	ret := wtype.NewLHComponent()
	ret.CName = name
	ret.Vol = volume

	return ret
}

type setOutputOrderTest struct {
	Instructions  []*wtype.LHInstruction
	OutputSort    bool
	ExpectedOrder []string
	ChainHeight   int
}

func (self *setOutputOrderTest) Run(t *testing.T) {
	rq := GetLHRequestForTest()

	for _, ins := range self.Instructions {
		rq.LHInstructions[ins.ID] = ins
	}

	rq.Options.OutputSort = self.OutputSort

	if err := setOutputOrder(rq); err != nil {
		t.Fatal(err)
	}

	if e, g := self.ChainHeight, rq.InstructionChain.Height(); e != g {
		t.Fatalf("Instruction chain length mismatch, e: %d, g: %d", e, g)
	}
	if e, g := len(self.ExpectedOrder), len(rq.Output_order); e != g {
		t.Fatalf("Expected Order length mismatch:\n\te: %v\n\tg: %v", self.ExpectedOrder, rq.Output_order)
	}

	outputOrder := make([]string, 0, len(rq.Output_order))
	for _, id := range rq.Output_order {
		//for promts check the message as the ID is overwritten
		if ins, ok := rq.LHInstructions[id]; ok && ins.Type == wtype.LHIPRM { //LHIPRM == prompt instruction
			outputOrder = append(outputOrder, ins.Message)
		} else {
			outputOrder = append(outputOrder, id)
		}
	}

	for i := range self.ExpectedOrder {
		if e, g := self.ExpectedOrder[i], outputOrder[i]; e != g {
			t.Fatalf("Expected Order mismatch in item %d:\n\te: %v\n\tg: %v", i, self.ExpectedOrder, outputOrder)
		}
	}

}

func TestSetOutputOrdering_Splits(t *testing.T) {

	water := GetLiquidForTest("water", 50.0)
	lemonJuice := GetLiquidForTest("lemonJuce", 10.0)

	split, waterSample, _ := GetSplitForTest("theSplit", water, 20.0)

	mix, _ := GetMixForTest("firstMix", waterSample, lemonJuice)

	test := setOutputOrderTest{
		Instructions:  []*wtype.LHInstruction{split, mix},
		OutputSort:    true,
		ExpectedOrder: []string{"firstMix", "theSplit"},
		ChainHeight:   2,
	}

	test.Run(t)
}

func TestSetOutputOrdering_SplitMixes(t *testing.T) {

	water := GetLiquidForTest("water", 250.0)
	concentrate := GetLiquidForTest("concentratedSquash", 5000.0)
	vodka := GetLiquidForTest("vodka", 50.0)

	split, concentrateSample, _ := GetSplitForTest("theSplit", concentrate, 25.0)
	mixSquash, squash := GetMixForTest("mixSquash", concentrateSample, water)
	mixShot, _ := GetMixForTest("mixShot", vodka, squash)

	test := setOutputOrderTest{
		Instructions:  []*wtype.LHInstruction{split, mixShot, mixSquash},
		OutputSort:    true,
		ExpectedOrder: []string{"mixSquash", "mixShot", "theSplit"},
		ChainHeight:   3,
	}

	test.Run(t)
}

func TestSetOutputOrdering_SplitMixes2(t *testing.T) {

	water := GetLiquidForTest("water", 250.0)
	concentrate := GetLiquidForTest("concentratedSquash", 5000.0)
	vodka := GetLiquidForTest("vodka", 50.0)
	milk := GetLiquidForTest("milk", 200.0)

	split, concentrateSample, concentrateRemainder := GetSplitForTest("theSplit", concentrate, 25.0)
	mixSquash, squash := GetMixForTest("mixSquash", concentrateSample, water)
	mixShot, _ := GetMixForTest("mixShot", vodka, squash)
	mixLumpy, _ := GetMixForTest("mixLumpy", milk, concentrateRemainder)

	test := setOutputOrderTest{
		Instructions:  []*wtype.LHInstruction{split, mixShot, mixSquash, mixLumpy},
		OutputSort:    true,
		ExpectedOrder: []string{"mixSquash", "mixShot", "theSplit", "mixLumpy"},
		ChainHeight:   4,
	}

	test.Run(t)
}

func TestSetOutputOrder_Prompt(t *testing.T) {
	// we go mix, prompt, split, mix
	var instructions []*wtype.LHInstruction
	var expectedIDOrder []string

	c1 := GetLiquidForTest("water", 10.0)
	c2 := GetLiquidForTest("washBuffer", 20.0)
	mix1, c3 := GetMixForTest("mix1", c1, c2)

	expectedIDOrder = append(expectedIDOrder, mix1.ID)
	instructions = append(instructions, mix1)

	prompt, promptResults := GetPromptForTest("prompt1", c3)
	c4 := promptResults[0]

	expectedIDOrder = append(expectedIDOrder, prompt.ID)
	instructions = append(instructions, prompt)

	split, c5, c4a := GetSplitForTest("split", c4, 10.0)
	instructions = append(instructions, split)

	c6 := GetLiquidForTest("turps", 100.0)
	mix2, _ := GetMixForTest("mix2", c6, c5)

	expectedIDOrder = append(expectedIDOrder, mix2.ID)
	expectedIDOrder = append(expectedIDOrder, split.ID)
	instructions = append(instructions, mix2)

	// mix the static component with some more water
	c7 := GetLiquidForTest("water", 200.0)
	mix3, _ := GetMixForTest("mix3", c4a, c7)

	expectedIDOrder = append(expectedIDOrder, mix3.ID)
	instructions = append(instructions, mix3)

	test := &setOutputOrderTest{
		Instructions:  instructions,
		ExpectedOrder: expectedIDOrder,
		ChainHeight:   5,
	}

	test.Run(t)
}
