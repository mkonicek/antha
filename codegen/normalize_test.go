package codegen

import (
	"testing"

	"github.com/antha-lang/antha/laboratory/effects"
)

func equals(as, bs []*effects.UseComp) bool {
	mas := make(map[*effects.UseComp]bool)
	bas := make(map[*effects.UseComp]bool)
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
	u1 := &effects.UseComp{}
	i1 := &effects.Command{
		From: []effects.Node{u1},
	}
	u2 := &effects.UseComp{
		From: []effects.Node{i1},
	}
	i2 := &effects.Command{
		From: []effects.Node{u2},
	}
	u3 := &effects.UseComp{
		From: []effects.Node{i2},
	}
	i3 := &effects.Command{
		From: []effects.Node{u3},
	}
	ir, err := build(i3)
	if err != nil {
		t.Fatal(err)
	}
	if es, fs := []*effects.UseComp{u3}, ir.reachingUses[i3]; !equals(es, fs) {
		t.Errorf("expected %v found %v", es, fs)
	} else if es, fs := []*effects.UseComp{u1}, ir.reachingUses[i1]; !equals(es, fs) {
		t.Errorf("expected %v found %v", es, fs)
	}
}

func TestReachingUsesMultiple(t *testing.T) {
	u1 := &effects.UseComp{}
	i1 := &effects.Command{
		From: []effects.Node{u1},
	}
	u2a := &effects.UseComp{}
	u2b := &effects.UseComp{
		From: []effects.Node{u2a},
	}
	u2c := &effects.UseComp{
		From: []effects.Node{i1},
	}
	i2 := &effects.Command{
		From: []effects.Node{u2b, u2c},
	}
	ir, err := build(i2)
	if err != nil {
		t.Fatal(err)
	}
	if es, fs := []*effects.UseComp{u2a, u2b, u2c}, ir.reachingUses[i2]; !equals(es, fs) {
		t.Errorf("expected %v found %v", es, fs)
	}
}
