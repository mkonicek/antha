package tests

import (
	"fmt"
	"testing"

	"github.com/antha-lang/antha/antha/anthalib/mixer"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/laboratory"
	"github.com/antha-lang/antha/laboratory/testlab"
	lh "github.com/antha-lang/antha/microArch/scheduler/liquidhandling"
)

func GetMixForTest(lab *laboratory.Laboratory, id string, input ...*wtype.Liquid) (*wtype.LHInstruction, *wtype.Liquid) {
	output := wtype.NewLHComponent(lab.IDGenerator)
	output.Type = wtype.LTNIL
	for _, ip := range input {
		output.Mix(lab.IDGenerator, ip)
	}

	mix := wtype.NewLHMixInstruction(lab.IDGenerator)
	mix.ID = id
	for _, ip := range input {
		mix.AddInput(ip)
	}
	mix.AddOutput(output)

	return mix, output
}

func GetSplitForTest(lab *laboratory.Laboratory, id string, input *wtype.Liquid, volume float64) (*wtype.LHInstruction, *wtype.Liquid, *wtype.Liquid) {
	split := wtype.NewLHSplitInstruction(lab.IDGenerator)
	split.ID = id
	moving, remaining := mixer.SplitSample(lab, input, wunit.NewVolume(volume, "ul"))

	split.AddInput(input)
	split.AddOutput(moving)
	split.AddOutput(remaining)

	return split, moving, remaining
}

func GetPromptForTest(lab *laboratory.Laboratory, message string, inputs ...*wtype.Liquid) (*wtype.LHInstruction, []*wtype.Liquid) {
	ret := wtype.NewLHPromptInstruction(lab.IDGenerator)
	ret.ID = message //ID will be overwritten, set the message as well for testing
	ret.Message = message
	for _, input := range inputs {
		output := input.Cp(lab.IDGenerator)
		output.ParentID = input.ID
		input.DaughtersID = map[string]struct{}{output.ID: {}}
		ret.AddInput(input)
		ret.AddOutput(output)
	}

	return ret, ret.Outputs
}

func GetLiquidForTest(lab *laboratory.Laboratory, name string, volume float64) *wtype.Liquid {
	ret := wtype.NewLHComponent(lab.IDGenerator)
	ret.CName = name
	ret.Vol = volume

	return ret
}

type setOutputOrderTest struct {
	Instructions   []*wtype.LHInstruction
	OutputSort     bool
	ExpectedOrder  []string
	ChainHeight    int
	ExpectingError bool
}

func (self *setOutputOrderTest) Run(lab *laboratory.Laboratory) error {
	insMap := make(map[string]*wtype.LHInstruction, len(self.Instructions))
	for _, instruction := range self.Instructions {
		insMap[instruction.ID] = instruction
	}

	ichain, err := lh.BuildInstructionChain(lab.IDGenerator, insMap)
	if encounteredError := err != nil; self.ExpectingError != encounteredError {
		return fmt.Errorf("ExpectingError: %t, Encountered Error: %v", self.ExpectingError, err)
	}

	//sort the instructions within each link of the chain
	ichain.SortInstructions(self.OutputSort)

	if e, g := self.ChainHeight, ichain.Height(); e != g {
		return fmt.Errorf("Instruction chain length mismatch, e: %d, g: %d", e, g)
	}
	if e, g := len(self.ExpectedOrder), len(ichain.FlattenInstructionIDs()); e != g {
		return fmt.Errorf("Expected Order length mismatch:\n\te: %v\n\tg: %v", e, g)
	}

	sorted := ichain.GetOrderedLHInstructions()
	outputOrder := make([]string, 0, len(sorted))
	for _, ins := range sorted {
		//for prompts check the message as the ID is overwritten
		if ins.Type == wtype.LHIPRM { //LHIPRM == prompt instruction
			outputOrder = append(outputOrder, ins.Message)
		} else {
			outputOrder = append(outputOrder, ins.ID)
		}
	}

	for i := range self.ExpectedOrder {
		if e, g := self.ExpectedOrder[i], outputOrder[i]; e != g {
			return fmt.Errorf("Expected Order mismatch in item %d:\n\te: %v\n\tg: %v", i, self.ExpectedOrder, outputOrder)
		}
	}
	return nil
}

func TestSetOutputOrdering_Splits(t *testing.T) {
	testlab.WithTestLab(t, "", &testlab.TestElementCallbacks{
		Steps: func(lab *laboratory.Laboratory) error {

			water := GetLiquidForTest(lab, "water", 50.0)
			lemonJuice := GetLiquidForTest(lab, "lemonJuce", 10.0)

			split, waterSample, _ := GetSplitForTest(lab, "theSplit", water, 20.0)

			mix, _ := GetMixForTest(lab, "firstMix", waterSample, lemonJuice)

			test := setOutputOrderTest{
				Instructions:  []*wtype.LHInstruction{split, mix},
				OutputSort:    true,
				ExpectedOrder: []string{"firstMix", "theSplit"},
				ChainHeight:   2,
			}

			return test.Run(lab)
		},
	})
}

func TestSetOutputOrdering_SplitUnused(t *testing.T) {
	testlab.WithTestLab(t, "", &testlab.TestElementCallbacks{
		Steps: func(lab *laboratory.Laboratory) error {
			water := GetLiquidForTest(lab, "water", 50.0)
			split, _, _ := GetSplitForTest(lab, "theSplit", water, 20.0)

			test := setOutputOrderTest{
				Instructions:   []*wtype.LHInstruction{split},
				OutputSort:     true,
				ExpectedOrder:  []string{},
				ChainHeight:    0,
				ExpectingError: true,
			}

			return test.Run(lab)
		},
	})
}

func TestSetOutputOrdering_SplitMixes(t *testing.T) {
	testlab.WithTestLab(t, "", &testlab.TestElementCallbacks{
		Steps: func(lab *laboratory.Laboratory) error {

			water := GetLiquidForTest(lab, "water", 250.0)
			concentrate := GetLiquidForTest(lab, "concentratedSquash", 5000.0)
			vodka := GetLiquidForTest(lab, "vodka", 50.0)

			split, concentrateSample, _ := GetSplitForTest(lab, "theSplit", concentrate, 25.0)
			mixSquash, squash := GetMixForTest(lab, "mixSquash", concentrateSample, water)
			mixShot, _ := GetMixForTest(lab, "mixShot", vodka, squash)

			test := setOutputOrderTest{
				Instructions:  []*wtype.LHInstruction{split, mixShot, mixSquash},
				OutputSort:    true,
				ExpectedOrder: []string{"mixSquash", "mixShot", "theSplit"},
				ChainHeight:   3,
			}

			return test.Run(lab)
		},
	})
}

func TestSetOutputOrdering_SplitMixes2(t *testing.T) {
	testlab.WithTestLab(t, "", &testlab.TestElementCallbacks{
		Steps: func(lab *laboratory.Laboratory) error {

			water := GetLiquidForTest(lab, "water", 250.0)
			concentrate := GetLiquidForTest(lab, "concentratedSquash", 5000.0)
			vodka := GetLiquidForTest(lab, "vodka", 50.0)
			milk := GetLiquidForTest(lab, "milk", 200.0)

			split, concentrateSample, concentrateRemainder := GetSplitForTest(lab, "theSplit", concentrate, 25.0)
			mixSquash, squash := GetMixForTest(lab, "mixSquash", concentrateSample, water)
			mixShot, _ := GetMixForTest(lab, "mixShot", vodka, squash)
			mixLumpy, _ := GetMixForTest(lab, "mixLumpy", milk, concentrateRemainder)

			test := setOutputOrderTest{
				Instructions:  []*wtype.LHInstruction{split, mixShot, mixSquash, mixLumpy},
				OutputSort:    true,
				ExpectedOrder: []string{"mixSquash", "mixShot", "theSplit", "mixLumpy"},
				ChainHeight:   4,
			}

			return test.Run(lab)
		},
	})
}

func TestSetOutputOrder_Prompt(t *testing.T) {
	testlab.WithTestLab(t, "", &testlab.TestElementCallbacks{
		Steps: func(lab *laboratory.Laboratory) error {
			// we go mix, prompt, split, mix

			c1 := GetLiquidForTest(lab, "water", 10.0)
			c2 := GetLiquidForTest(lab, "washBuffer", 20.0)
			mix1, c3 := GetMixForTest(lab, "mix1", c1, c2)

			prompt, promptResults := GetPromptForTest(lab, "prompt1", c3)

			c4 := promptResults[0]
			split, c5, c4a := GetSplitForTest(lab, "split", c4, 10.0)

			c6 := GetLiquidForTest(lab, "turps", 100.0)
			mix2, _ := GetMixForTest(lab, "mix2", c6, c5)

			// mix the static component with some more water
			c7 := GetLiquidForTest(lab, "water", 200.0)
			mix3, _ := GetMixForTest(lab, "mix3", c4a, c7)

			test := &setOutputOrderTest{
				Instructions:  []*wtype.LHInstruction{mix1, prompt, split, mix2, mix3},
				ExpectedOrder: []string{mix1.ID, prompt.ID, mix2.ID, split.ID, mix3.ID},
				ChainHeight:   5,
			}

			return test.Run(lab)
		},
	})
}
