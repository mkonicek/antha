package liquidhandling

import (
	"fmt"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"testing"
)

func getComponentWithNameVolume(name string, volume float64) *wtype.LHComponent {
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
	ins.Components = []*wtype.LHComponent{c1, c2}
	ins.Result = c3
	ins.ProductID = ins.Result.ID

	req.LHInstructions[ins.ID] = ins

	/*
	   Parent *IChain
	   Child  *IChain
	   Values []*wtype.LHInstruction
	   Depth  int
	*/

	ic := &IChain{
		Parent: nil,
		Child:  nil,
		Values: []*wtype.LHInstruction{ins},
		Depth:  0,
	}

	req.InstructionChain = ic

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
		ins.Components = []*wtype.LHComponent{smp}
		res := getComponentWithNameVolume("water+milk", 15.0)
		res.ParentID = ins.Components[0].ID
		ins.Result = res
		ins.Result.DeclareInstance()
		ins.ProductID = ins.Result.ID
		req.LHInstructions[ins.ID] = ins
		inss = append(inss, ins)
	}

	ic.Child = &IChain{Parent: ic, Child: nil, Values: inss, Depth: 1}

	// try fixing the volumes

	req, err := FixVolumes(req)

	if err != nil {
		t.Errorf(err.Error())
	}

	// check to see if the result of the first mix is now 150.0 ul

	mix1 := req.InstructionChain.Values[0]

	if mix1.Result.Vol != 155.0 {
		t.Errorf(fmt.Sprintf("Expected 155.0 got volume %s", mix1.Result.Volume()))
	}
}

func TestFixVolumes2(t *testing.T) {
	// findUpdateInstructionVolumes(ch *IChain, wanted map[string]wunit.Volume) (map[string]wunit.Volume, error)

	c1 := getComponentWithNameVolume("water", 50.0)
	c2 := getComponentWithNameVolume("milk", 50.0)

	c3 := c1.Dup()
	c3.Mix(c2)
	c3.DeclareInstance()

	ins := wtype.NewLHMixInstruction()
	ins.Components = []*wtype.LHComponent{c1, c2}
	ins.Result = c3
	ins.ProductID = ins.Result.ID

	inss := []*wtype.LHInstruction{ins}

	ch := &IChain{Parent: nil, Child: nil, Values: inss, Depth: 0}

	want := make(map[string]wunit.Volume, 1)

	want[c3.FullyQualifiedName()] = wunit.NewVolume(150.0, "ul")

	newWant, _ := findUpdateInstructionVolumes(ch, want, make(map[string]*wtype.LHPlate))

	for _, v := range newWant {
		if !v.EqualTo(wunit.NewVolume(75.5, "ul")) {
			t.Errorf(fmt.Sprintf("Expected 75.5 ul got %s", v))
		}
	}

}

func TestFixVolumes3(t *testing.T) {
	t.Skip()
	req := NewLHRequest()

	c1 := getComponentWithNameVolume("water", 50.0)
	c2 := getComponentWithNameVolume("milk", 50.0)

	c3 := c1.Dup()
	c3.Mix(c2)
	c3.DeclareInstance()

	ins := wtype.NewLHMixInstruction()
	ins.Components = []*wtype.LHComponent{c1, c2}
	ins.Result = c3
	ins.ProductID = ins.Result.ID

	req.LHInstructions[ins.ID] = ins

	ic := &IChain{
		Parent: nil,
		Child:  nil,
		Values: []*wtype.LHInstruction{ins},
		Depth:  0,
	}

	req.InstructionChain = ic

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
		ins.Components = []*wtype.LHComponent{smp}
		res := getComponentWithNameVolume("water+milk", 15.0)
		res.ParentID = ins.Components[0].ID
		ins.Result = res
		ins.Result.DeclareInstance()
		ins.ProductID = ins.Result.ID
		req.LHInstructions[ins.ID] = ins
		inss = append(inss, ins)
	}

	ic.Child = &IChain{Parent: ic, Child: nil, Values: inss, Depth: 1}

	// try fixing the volumes

	req, err := FixVolumes(req)

	if err != nil {
		t.Errorf(err.Error())
	}

	// check to see if the result of the first mix is now 155.0 ul (10 * 15.0 +  0.5)

	mix1 := req.InstructionChain.Values[0]

	if mix1.Result.Vol != 155.0 {
		t.Errorf(fmt.Sprintf("Expected 155.0 got volume %s", mix1.Result.Volume()))
	}
}
