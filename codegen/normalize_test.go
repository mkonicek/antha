package codegen

import (
	"testing"

	"github.com/antha-lang/antha/ast"
)

func equals(as, bs []*ast.UseComp) bool {
	mas := make(map[*ast.UseComp]bool)
	bas := make(map[*ast.UseComp]bool)
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
	u1 := &ast.UseComp{}
	i1 := &ast.Command{
		From: []ast.Node{u1},
	}
	u2 := &ast.UseComp{
		From: []ast.Node{i1},
	}
	i2 := &ast.Command{
		From: []ast.Node{u2},
	}
	u3 := &ast.UseComp{
		From: []ast.Node{i2},
	}
	i3 := &ast.Command{
		From: []ast.Node{u3},
	}
	ir, err := build(i3)
	if err != nil {
		t.Fatal(err)
	}
	if es, fs := []*ast.UseComp{u3}, ir.reachingUses[i3]; !equals(es, fs) {
		t.Errorf("expected %q found %q", es, fs)
	} else if es, fs := []*ast.UseComp{u1}, ir.reachingUses[i1]; !equals(es, fs) {
		t.Errorf("expected %q found %q", es, fs)
	}
}

func TestReachingUsesMultiple(t *testing.T) {
	u1 := &ast.UseComp{}
	i1 := &ast.Command{
		From: []ast.Node{u1},
	}
	u2a := &ast.UseComp{}
	u2b := &ast.UseComp{
		From: []ast.Node{u2a},
	}
	u2c := &ast.UseComp{
		From: []ast.Node{i1},
	}
	i2 := &ast.Command{
		From: []ast.Node{u2b, u2c},
	}
	ir, err := build(i2)
	if err != nil {
		t.Fatal(err)
	}
	if es, fs := []*ast.UseComp{u2a, u2b, u2c}, ir.reachingUses[i2]; !equals(es, fs) {
		t.Errorf("expected %q found %q", es, fs)
	}
}
