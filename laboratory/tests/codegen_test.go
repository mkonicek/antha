package tests

import (
	"fmt"
	"testing"

	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/codegen"
	"github.com/antha-lang/antha/laboratory/effects"
	"github.com/antha-lang/antha/laboratory/testlab"
	"github.com/antha-lang/antha/target"
	"github.com/antha-lang/antha/target/human"
	"github.com/antha-lang/antha/workflow"
)

type lowLevelTestInst struct {
	effects.DependsMixin
	effects.IdMixin
	effects.NoDeviceMixin
}

type highLevelTestInst struct{}

type testDriver struct{}

var testDriverSelector = effects.NameValue{
	Name:  target.DriverSelectorV1Name,
	Value: "antha.test.v0",
}

func (a *testDriver) CanCompile(req effects.Request) bool {
	can := effects.Request{
		Selector: []effects.NameValue{
			testDriverSelector,
		},
	}
	return can.Contains(req)
}

func (a *testDriver) Compile(labEffects *effects.LaboratoryEffects, dir string, cmds []effects.Node) (effects.Insts, error) {
	for _, n := range cmds {
		if c, ok := n.(*effects.Command); !ok {
			return nil, fmt.Errorf("unexpected node %T", n)
		} else if _, ok := c.Inst.(*highLevelTestInst); !ok {
			return nil, fmt.Errorf("unexpected inst %T", c.Inst)
		}
	}
	return effects.Insts{&lowLevelTestInst{}}, nil
}

func (a *testDriver) Connect(wf *workflow.Workflow) error {
	return nil
}

func (a *testDriver) Close() {}

func (a *testDriver) Id() workflow.DeviceInstanceID {
	return workflow.DeviceInstanceID("testDevice")
}

func TestWellFormed(t *testing.T) {
	labEffects := testlab.NewTestLabEffects(nil, nil)

	nodes := make([]effects.Node, 4)
	for idx := 0; idx < len(nodes); idx++ {
		m := &effects.Command{
			Request: effects.Request{
				Selector: []effects.NameValue{
					target.DriverSelectorV1Mixer,
				},
			},
			Inst: &wtype.LHInstruction{},
			From: []effects.Node{
				&effects.UseComp{},
				&effects.UseComp{},
				&effects.UseComp{},
			},
		}
		u := &effects.UseComp{
			From: []effects.Node{m},
		}

		t := &effects.Command{
			Request: effects.Request{
				Selector: []effects.NameValue{
					testDriverSelector,
				},
			},
			Inst: &highLevelTestInst{},
			From: []effects.Node{u},
		}

		nodes[idx] = t
	}

	tgt := target.New()
	tgt.AddDevice(&testDriver{})
	human.New(labEffects.IDGenerator).DetermineRole(tgt)

	if insts, err := codegen.Compile(labEffects, "", tgt, nodes); err != nil {
		t.Fatal(err)
	} else if l := len(insts); l == 0 {
		t.Errorf("expected > %d instructions found %d", 0, l)
	} else if last, ok := insts[l-1].(*lowLevelTestInst); !ok {
		t.Errorf("expected testInst found %T", insts[l-1])
	} else if n := len(last.Depends); n != 1 {
		t.Errorf("expected %d dependencies found %d", 1, n)
	}
}
