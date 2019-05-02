package liquidhandling

import (
	"fmt"
	"testing"

	"github.com/antha-lang/antha/antha/anthalib/mixer"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
)

func getComponentWithNameVolume(name string, volume float64) *wtype.Liquid {
	c := wtype.NewLHComponent()
	c.CName = name
	c.Vol = volume
	c.Type = wtype.LTWater
	return c
}

func TestFixVolumes(t *testing.T) {
	req := NewLHRequest()

	c1 := getComponentWithNameVolume("water", 50.0)
	c2 := getComponentWithNameVolume("milk", 50.0)

	c3 := c1.Dup()
	c3.Mix(c2)
	c3.DeclareInstance()

	ins := wtype.NewLHMixInstruction()
	ins.Inputs = []*wtype.Liquid{c1, c2}
	ins.AddOutput(c3)

	req.LHInstructions[ins.ID] = ins

	ic := &wtype.IChain{
		Parent: nil,
		Child:  nil,
		Values: []*wtype.LHInstruction{ins},
		Depth:  0,
	}

	req.InstructionChain = ic

	//now take lots of samples
	inss := make([]*wtype.LHInstruction, 0, 10)

	for i := 0; i < 10; i++ {
		ins = wtype.NewLHMixInstruction()
		smp, err := c3.Sample(wunit.NewVolume(15.0, "ul"))
		smp.SetSample(true)
		smp.DeclareInstance()
		smp.ParentID = c3.ID
		if err != nil {
			t.Errorf(err.Error())
		}
		c3.Vol = 100.0
		ins.Inputs = []*wtype.Liquid{smp}
		res := getComponentWithNameVolume("water+milk", 15.0)
		res.ParentID = ins.Inputs[0].ID
		res.DeclareInstance()
		ins.AddOutput(res)
		req.LHInstructions[ins.ID] = ins
		inss = append(inss, ins)
	}

	ic.Child = &wtype.IChain{Parent: ic, Child: nil, Values: inss, Depth: 1}

	// try fixing the volumes

	err := FixVolumes(req, wtype.GLOBALCARRYVOLUME)

	if err != nil {
		t.Errorf(err.Error())
	}

	// check to see if the result of the first mix is now 155.0 ul

	mix1 := req.InstructionChain.Values[0]

	if mix1.Outputs[0].Vol != 155.0 {
		t.Errorf(fmt.Sprintf("Expected 155.0 got volume %s", mix1.Outputs[0].Volume()))
	}
}

func TestFixVolumes2(t *testing.T) {

	c1 := getComponentWithNameVolume("water", 50.0)
	c2 := getComponentWithNameVolume("milk", 50.0)

	c3 := c1.Dup()
	c3.Mix(c2)
	c3.DeclareInstance()

	ins := wtype.NewLHMixInstruction()
	ins.Inputs = []*wtype.Liquid{c1, c2}
	ins.AddOutput(c3)

	inss := []*wtype.LHInstruction{ins}

	ch := &wtype.IChain{Parent: nil, Child: nil, Values: inss, Depth: 0}

	want := make(map[string]wunit.Volume, 1)

	want[c3.FullyQualifiedName()] = wunit.NewVolume(150.0, "ul")

	newWant, _ := findUpdateInstructionVolumes(ch, want, make(map[string]*wtype.Plate), wunit.NewVolume(0.5, "ul"))

	v := newWant["water"+wtype.InPlaceMarker]
	if !v.EqualTo(wunit.NewVolume(75.0, "ul")) {
		t.Errorf(fmt.Sprintf("Expected 75.0 ul got %s", v))
	}
	v = newWant["milk"]
	if !v.EqualTo(wunit.NewVolume(75.5, "ul")) {
		t.Errorf(fmt.Sprintf("Expected 75.5 ul got %s", v))
	}

}

func TestFixVolumes3(t *testing.T) {
	//	t.Skip()
	req := NewLHRequest()

	c1 := getComponentWithNameVolume("water", 50.0)
	c3 := c1.Dup()

	c3.DeclareInstance()

	ins := wtype.NewLHMixInstruction()
	ins.Inputs = []*wtype.Liquid{c1}

	ins.AddOutput(c3)
	req.LHInstructions[ins.ID] = ins

	ic := &wtype.IChain{
		Parent: nil,
		Child:  nil,
		Values: []*wtype.LHInstruction{ins},
		Depth:  0,
	}

	req.InstructionChain = ic

	// mix-in-place

	c2 := getComponentWithNameVolume("milk", 50.0)
	ins = wtype.NewLHMixInstruction()
	ins.Inputs = []*wtype.Liquid{c3, c2}
	c4 := c3.Dup()
	c3.Mix(c2)
	ins.AddOutput(c4)
	req.LHInstructions[ins.ID] = ins

	ic = &wtype.IChain{
		Parent: req.InstructionChain,
		Child:  nil,
		Values: []*wtype.LHInstruction{ins},
		Depth:  1,
	}

	req.InstructionChain.Child = ic

	//now take lots of samples
	inss := make([]*wtype.LHInstruction, 0, 10)

	for i := 0; i < 10; i++ {
		ins = wtype.NewLHMixInstruction()
		smp, err := c4.Sample(wunit.NewVolume(15.0, "ul"))
		smp.SetSample(true)
		smp.DeclareInstance()
		smp.ParentID = c4.ID
		if err != nil {
			t.Errorf(err.Error())
		}
		c4.Vol = 100.0
		ins.Inputs = []*wtype.Liquid{smp}
		res := getComponentWithNameVolume("water+milk", 15.0)
		res.ParentID = ins.Inputs[0].ID
		res.DeclareInstance()
		ins.AddOutput(res)
		req.LHInstructions[ins.ID] = ins
		inss = append(inss, ins)
	}

	ic.Child = &wtype.IChain{Parent: ic, Child: nil, Values: inss, Depth: 2}

	// try fixing the volumes

	err := FixVolumes(req, wtype.GLOBALCARRYVOLUME)

	if err != nil {
		t.Errorf(err.Error())
	}

	// check to see if the result of the first mix is now 155.0 ul (10 * 15.0 +  0.5)

	mix1 := req.InstructionChain.Values[0]

	if mix1.Outputs[0].Vol != 155.0 {
		t.Errorf(fmt.Sprintf("Expected 155.0 got volume %s", mix1.Outputs[0].Volume()))
	}
}

func TestFixVolumes4(t *testing.T) {
	req := NewLHRequest()

	c1 := getComponentWithNameVolume("water", 50.0)
	c3 := c1.Cp()

	c3.DeclareInstance()

	ins := wtype.NewLHMixInstruction()
	ins.Inputs = []*wtype.Liquid{c1}

	ins.AddOutput(c3)
	req.LHInstructions[ins.ID] = ins

	ic := &wtype.IChain{
		Parent: nil,
		Child:  nil,
		Values: []*wtype.LHInstruction{ins},
		Depth:  0,
	}

	req.InstructionChain = ic

	ins = wtype.NewLHPromptInstruction()

	c4 := c3.Cp()

	ins.AddInput(c3)
	ins.AddOutput(c4)

	ic.Child = &wtype.IChain{
		Parent: ic,
		Child:  nil,
		Values: []*wtype.LHInstruction{ins},
		Depth:  1,
	}

	ins = wtype.NewLHMixInstruction()

	c5 := c4.Cp()
	c5.Vol = 200.0

	ins.Inputs = []*wtype.Liquid{c4}
	ins.AddOutput(c5)

	ic.Child.Child = &wtype.IChain{
		Parent: ic.Child,
		Child:  nil,
		Values: []*wtype.LHInstruction{ins},
		Depth:  2,
	}

	err := FixVolumes(req, wtype.GLOBALCARRYVOLUME)

	if err != nil {
		t.Errorf(err.Error())
	}
}

// test for splitsample fixing
func TestFixVolumesSplitSample(t *testing.T) {
	req := NewLHRequest()

	c1 := getComponentWithNameVolume("water", 50.0)
	c2 := getComponentWithNameVolume("milk", 50.0)

	c3 := c1.Dup()
	c3.Mix(c2)
	c3.DeclareInstance()

	ins := wtype.NewLHMixInstruction()
	ins.Inputs = []*wtype.Liquid{c1, c2}
	ins.AddOutput(c3)

	req.LHInstructions[ins.ID] = ins

	ic := &wtype.IChain{
		Parent: nil,
		Child:  nil,
		Values: []*wtype.LHInstruction{ins},
		Depth:  0,
	}

	req.InstructionChain = ic

	//now take lots of split samples
	mixInss := make([]*wtype.LHInstruction, 0, 10)
	splInss := make([]*wtype.LHInstruction, 0, 10)

	for i := 0; i < 10; i++ {
		ins = wtype.NewLHMixInstruction()
		//smp, err := c3.Sample(wunit.NewVolume(15.0, "ul"))
		//smp.SetSample(true)
		//smp.DeclareInstance()
		//smp.ParentID = c3.ID

		smp, newC3 := mixer.SplitSample(c3, wunit.NewVolume(15.0, "ul"))

		ins.Inputs = []*wtype.Liquid{smp}
		res := getComponentWithNameVolume("water+milk", 15.0)
		res.ParentID = ins.Inputs[0].ID
		res.DeclareInstance()
		ins.AddOutput(res)
		req.LHInstructions[ins.ID] = ins

		// make the split instruction
		splitIns := wtype.NewLHSplitInstruction()
		splitIns.AddInput(c3)
		splitIns.AddOutput(smp)
		splitIns.AddOutput(newC3)
		c3.Vol = 100.0
		c3 = newC3

		mixInss = append(mixInss, ins)
		splInss = append(splInss, splitIns)
	}

	ic.Child = &wtype.IChain{Parent: ic, Child: nil, Values: splInss, Depth: 1}
	ic.Child.Child = &wtype.IChain{Parent: ic.Child, Child: nil, Values: mixInss, Depth: 2}

	// try fixing the volumes

	err := FixVolumes(req, wtype.GLOBALCARRYVOLUME)

	if err != nil {
		t.Errorf(err.Error())
	}

	// check to see if the result of the first mix is now 155.0 ul

	mix1 := req.InstructionChain.Values[0]

	if mix1.Outputs[0].Vol != 155.0 {
		t.Errorf(fmt.Sprintf("Expected 155.0 got volume %s", mix1.Outputs[0].Volume()))
	}
}
