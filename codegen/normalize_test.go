package codegen

import (
	"testing"

	"github.com/antha-lang/antha/instructions"
)

func equals(as, bs []*instructions.UseComp) bool {
	mas := make(map[*instructions.UseComp]bool)
	bas := make(map[*instructions.UseComp]bool)
	for _, v := range as {
		mas[v] = true
	}
	for _, v := range bs {
		bas[v] = true
		if !mas[v] {
			return false
		}
	}

	return len(mas) == len(bas)
}

func TestReachingUsesChain(t *testing.T) {
	u1 := &instructions.UseComp{}
	i1 := &instructions.Command{
		From: []instructions.Node{u1},
	}
	u2 := &instructions.UseComp{
		From: []instructions.Node{i1},
	}
	i2 := &instructions.Command{
		From: []instructions.Node{u2},
	}
	u3 := &instructions.UseComp{
		From: []instructions.Node{i2},
	}
	i3 := &instructions.Command{
		From: []instructions.Node{u3},
	}
	ir, err := build(i3)
	if err != nil {
		t.Fatal(err)
	}
	if es, fs := []*instructions.UseComp{u3}, ir.reachingUses[i3]; !equals(es, fs) {
		t.Errorf("expected %v found %v", es, fs)
	} else if es, fs := []*instructions.UseComp{u1}, ir.reachingUses[i1]; !equals(es, fs) {
		t.Errorf("expected %v found %v", es, fs)
	}
}

func TestReachingUsesMultiple(t *testing.T) {
	u1 := &instructions.UseComp{}
	i1 := &instructions.Command{
		From: []instructions.Node{u1},
	}
	u2a := &instructions.UseComp{}
	u2b := &instructions.UseComp{
		From: []instructions.Node{u2a},
	}
	u2c := &instructions.UseComp{
		From: []instructions.Node{i1},
	}
	i2 := &instructions.Command{
		From: []instructions.Node{u2b, u2c},
	}
	ir, err := build(i2)
	if err != nil {
		t.Fatal(err)
	}
	if es, fs := []*instructions.UseComp{u2a, u2b, u2c}, ir.reachingUses[i2]; !equals(es, fs) {
		t.Errorf("expected %v found %v", es, fs)
	}
}
