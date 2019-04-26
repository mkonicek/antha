package tests

import (
	"fmt"
	"testing"

	"github.com/antha-lang/antha/antha/anthalib/mixer"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/laboratory"
	"github.com/antha-lang/antha/laboratory/effects/id"
	"github.com/antha-lang/antha/laboratory/testlab"
	lh "github.com/antha-lang/antha/microArch/scheduler/liquidhandling"
)

func getComponentWithNameVolume(idGen *id.IDGenerator, name string, volume float64) *wtype.Liquid {
	c := wtype.NewLHComponent(idGen)
	c.CName = name
	c.Vol = volume
	c.Type = wtype.LTWater
	return c
}

func TestFixVolumes(t *testing.T) {
	idGen := id.NewIDGenerator(t.Name())
	req := lh.NewLHRequest(idGen)

	c1 := getComponentWithNameVolume(idGen, "water", 50.0)
	c2 := getComponentWithNameVolume(idGen, "milk", 50.0)

	c3 := c1.Dup(idGen)
	c3.Mix(idGen, c2)
	c3.DeclareInstance()

	ins := wtype.NewLHMixInstruction(idGen)
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
		ins = wtype.NewLHMixInstruction(idGen)
		smp, err := c3.Sample(idGen, wunit.NewVolume(15.0, "ul"))
		smp.SetSample(true)
		smp.DeclareInstance()
		smp.ParentID = c3.ID
		if err != nil {
			t.Error(err)
		}
		c3.Vol = 100.0
		ins.Inputs = []*wtype.Liquid{smp}
		res := getComponentWithNameVolume(idGen, "water+milk", 15.0)
		res.ParentID = ins.Inputs[0].ID
		res.DeclareInstance()
		ins.AddOutput(res)
		req.LHInstructions[ins.ID] = ins
		inss = append(inss, ins)
	}

	ic.Child = &wtype.IChain{Parent: ic, Child: nil, Values: inss, Depth: 1}

	// try fixing the volumes

	err := lh.FixVolumes(req, wtype.GLOBALCARRYVOLUME)

	if err != nil {
		t.Error(err)
	}

	// check to see if the result of the first mix is now 155.0 ul

	mix1 := req.InstructionChain.Values[0]

	if mix1.Outputs[0].Vol != 155.0 {
		t.Errorf("Expected 155.0 got volume %s", mix1.Outputs[0].Volume())
	}
}

func TestFixVolumes2(t *testing.T) {
	idGen := id.NewIDGenerator(t.Name())
	c1 := getComponentWithNameVolume(idGen, "water", 50.0)
	c2 := getComponentWithNameVolume(idGen, "milk", 50.0)

	c3 := c1.Dup(idGen)
	c3.Mix(idGen, c2)
	c3.DeclareInstance()

	ins := wtype.NewLHMixInstruction(idGen)
	ins.Inputs = []*wtype.Liquid{c1, c2}
	ins.AddOutput(c3)

	inss := []*wtype.LHInstruction{ins}

	ch := &wtype.IChain{Parent: nil, Child: nil, Values: inss, Depth: 0}

	want := make(map[string]wunit.Volume, 1)

	want[c3.FullyQualifiedName()] = wunit.NewVolume(150.0, "ul")

	newWant, _ := lh.FindUpdateInstructionVolumes(ch, want, make(map[string]*wtype.Plate), wunit.NewVolume(0.5, "ul"))

	v := newWant["water"+wtype.InPlaceMarker]
	if !v.EqualTo(wunit.NewVolume(75.0, "ul")) {
		t.Errorf("Expected 75.0 ul got %s", v)
	}
	v = newWant["milk"]
	if !v.EqualTo(wunit.NewVolume(75.5, "ul")) {
		t.Errorf("Expected 75.5 ul got %s", v)
	}

}

func TestFixVolumes3(t *testing.T) {
	//	t.Skip()
	idGen := id.NewIDGenerator(t.Name())
	req := lh.NewLHRequest(idGen)

	c1 := getComponentWithNameVolume(idGen, "water", 50.0)
	c3 := c1.Dup(idGen)

	c3.DeclareInstance()

	ins := wtype.NewLHMixInstruction(idGen)
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

	c2 := getComponentWithNameVolume(idGen, "milk", 50.0)
	ins = wtype.NewLHMixInstruction(idGen)
	ins.Inputs = []*wtype.Liquid{c3, c2}
	c4 := c3.Dup(idGen)
	c3.Mix(idGen, c2)
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
		ins = wtype.NewLHMixInstruction(idGen)
		smp, err := c4.Sample(idGen, wunit.NewVolume(15.0, "ul"))
		smp.SetSample(true)
		smp.DeclareInstance()
		smp.ParentID = c4.ID
		if err != nil {
			t.Error(err)
		}
		c4.Vol = 100.0
		ins.Inputs = []*wtype.Liquid{smp}
		res := getComponentWithNameVolume(idGen, "water+milk", 15.0)
		res.ParentID = ins.Inputs[0].ID
		res.DeclareInstance()
		ins.AddOutput(res)
		req.LHInstructions[ins.ID] = ins
		inss = append(inss, ins)
	}

	ic.Child = &wtype.IChain{Parent: ic, Child: nil, Values: inss, Depth: 2}

	// try fixing the volumes

	err := lh.FixVolumes(req, wtype.GLOBALCARRYVOLUME)

	if err != nil {
		t.Error(err)
	}

	// check to see if the result of the first mix is now 155.0 ul (10 * 15.0 +  0.5)

	mix1 := req.InstructionChain.Values[0]

	if mix1.Outputs[0].Vol != 155.0 {
		t.Errorf("Expected 155.0 got volume %s", mix1.Outputs[0].Volume())
	}
}

func TestFixVolumes4(t *testing.T) {
	idGen := id.NewIDGenerator(t.Name())
	req := lh.NewLHRequest(idGen)

	c1 := getComponentWithNameVolume(idGen, "water", 50.0)
	c3 := c1.Cp(idGen)

	c3.DeclareInstance()

	ins := wtype.NewLHMixInstruction(idGen)
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

	ins = wtype.NewLHPromptInstruction(idGen)

	c4 := c3.Cp(idGen)

	ins.PassThrough[c3.ID] = c4

	ic.Child = &wtype.IChain{
		Parent: ic,
		Child:  nil,
		Values: []*wtype.LHInstruction{ins},
		Depth:  1,
	}

	ins = wtype.NewLHMixInstruction(idGen)

	c5 := c4.Cp(idGen)
	c5.Vol = 200.0

	ins.Inputs = []*wtype.Liquid{c4}
	ins.AddOutput(c5)

	ic.Child.Child = &wtype.IChain{
		Parent: ic.Child,
		Child:  nil,
		Values: []*wtype.LHInstruction{ins},
		Depth:  2,
	}

	err := lh.FixVolumes(req, wtype.GLOBALCARRYVOLUME)

	if err != nil {
		t.Error(err)
	}
}

// test for splitsample fixing
func TestFixVolumesSplitSample(t *testing.T) {
	testlab.WithTestLab(t, "", &testlab.TestElementCallbacks{
		Steps: func(lab *laboratory.Laboratory) error {

			req := lh.NewLHRequest(lab.IDGenerator)

			c1 := getComponentWithNameVolume(lab.IDGenerator, "water", 50.0)
			c2 := getComponentWithNameVolume(lab.IDGenerator, "milk", 50.0)

			c3 := c1.Dup(lab.IDGenerator)
			c3.Mix(lab.IDGenerator, c2)
			c3.DeclareInstance()

			ins := wtype.NewLHMixInstruction(lab.IDGenerator)
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
				ins = wtype.NewLHMixInstruction(lab.IDGenerator)
				//smp, err := c3.Sample(wunit.NewVolume(15.0, "ul"))
				//smp.SetSample(true)
				//smp.DeclareInstance()
				//smp.ParentID = c3.ID

				smp, newC3 := mixer.SplitSample(lab, c3, wunit.NewVolume(15.0, "ul"))

				ins.Inputs = []*wtype.Liquid{smp}
				res := getComponentWithNameVolume(lab.IDGenerator, "water+milk", 15.0)
				res.ParentID = ins.Inputs[0].ID
				res.DeclareInstance()
				ins.AddOutput(res)
				req.LHInstructions[ins.ID] = ins

				// make the split instruction
				splitIns := wtype.NewLHSplitInstruction(lab.IDGenerator)
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

			err := lh.FixVolumes(req, wtype.GLOBALCARRYVOLUME)

			if err != nil {
				return err
			}

			// check to see if the result of the first mix is now 155.0 ul

			mix1 := req.InstructionChain.Values[0]

			if mix1.Outputs[0].Vol != 155.0 {
				return fmt.Errorf("Expected 155.0 got volume %s", mix1.Outputs[0].Volume())
			}
			return nil
		},
	})
}
