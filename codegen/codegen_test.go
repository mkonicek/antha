package codegen

import (
	"context"
	"fmt"
	"testing"

	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/ast"
	"github.com/antha-lang/antha/target"
	"github.com/antha-lang/antha/target/human"
)

type incubateInst struct {
	Depends []ast.Inst
}

func (a *incubateInst) Device() ast.Device {
	return nil
}

func (a *incubateInst) DependsOn() []ast.Inst {
	return a.Depends
}

func (a *incubateInst) SetDependsOn(xs ...ast.Inst) {
	a.Depends = xs
}

func (a *incubateInst) AppendDependsOn(xs ...ast.Inst) {
	a.Depends = append(a.Depends, xs...)
}

type incubator struct{}

func (a *incubator) CanCompile(req ast.Request) bool {
	can := ast.Request{}
	can.Selector = append(can.Selector, target.DriverSelectorV1ShakerIncubator)
	return can.Contains(req)
}

func (a *incubator) Compile(ctx context.Context, nodes []ast.Node) ([]ast.Inst, error) {
	for _, n := range nodes {
		if c, ok := n.(*ast.Command); !ok {
			return nil, fmt.Errorf("unexpected node %T", n)
		} else if _, ok := c.Inst.(*ast.IncubateInst); !ok {
			return nil, fmt.Errorf("unexpected inst %T", c.Inst)
		}
	}
	return []ast.Inst{&incubateInst{}}, nil
}

func TestWellFormed(t *testing.T) {
	ctx := context.Background()

	var nodes []ast.Node
	for idx := 0; idx < 4; idx++ {
		m := &ast.Command{
			Request: ast.Request{
				Selector: []ast.NameValue{
					target.DriverSelectorV1Mixer,
				},
			},
			Inst: &wtype.LHInstruction{},
			From: []ast.Node{
				&ast.UseComp{},
				&ast.UseComp{},
				&ast.UseComp{},
			},
		}
		u := &ast.UseComp{}
		u.From = append(u.From, m)

		i := &ast.Command{
			Request: ast.Request{
				Selector: []ast.NameValue{
					target.DriverSelectorV1ShakerIncubator,
				},
			},
			Inst: &ast.IncubateInst{},
			From: []ast.Node{u},
		}

		nodes = append(nodes, i)
	}

	machine := target.New()
	machine.AddDevice(human.New(human.Opt{CanMix: true}))
	machine.AddDevice(&incubator{})

	if insts, err := Compile(ctx, machine, nodes); err != nil {
		t.Fatal(err)
	} else if l := len(insts); l == 0 {
		t.Errorf("expected > %d instructions found %d", 0, l)
	} else if last, ok := insts[l-1].(*incubateInst); !ok {
		t.Errorf("expected incubateInst found %T", insts[l-1])
	} else if n := len(last.Depends); n != 1 {
		t.Errorf("expected %d dependencies found %d", 1, n)
	}
}
