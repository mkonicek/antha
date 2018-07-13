package liquidhandling

import (
	"github.com/antha-lang/antha/antha/anthalib/mixer"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"testing"
)

func GetMixForTest(id string, input ...*wtype.Liquid) (*wtype.LHInstruction, *wtype.Liquid) {
	output := input[0].Dup()
	for _, ip := range input[1:] {
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

	if err := set_output_order(rq); err != nil {
		t.Fatal(err)
	}

	if e, g := self.ChainHeight, rq.InstructionChain.Height(); e != g {
		t.Fatalf("Instruction chain length mismatch, e: %d, g: %d", e, g)
	}
	if e, g := len(self.ExpectedOrder), len(rq.Output_order); e != g {
		t.Fatalf("Expected Order length mismatch:\n\te: %v\n\tg: %v", self.ExpectedOrder, rq.Output_order)
	}
	for i := range self.ExpectedOrder {
		if e, g := self.ExpectedOrder[i], rq.Output_order[i]; e != g {
			t.Fatalf("Expected Order mismatch in item %d:\n\te: %v\n\tg: %v", i, self.ExpectedOrder, rq.Output_order)
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
