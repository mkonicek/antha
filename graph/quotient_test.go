package graph

import (
	"testing"
)

func TestQuotient(t *testing.T) {
	g := MakeTestGraph(map[string][]string{
		"root": {"a", "b"},
		"a":    {"c", "d"},
		"e":    {},
	})
	cg := MakeQuotient(MakeQuotientOpt{
		Graph: g,
		Colorer: func(n Node) interface{} {
			return 1
		},
	})
	if n := cg.NumNodes(); n != 1 {
		t.Errorf("expected 1 node but found %d", n)
	}
}

func TestQuotienWithHasColor(t *testing.T) {
	g := MakeTestGraph(map[string][]string{
		"root": {".a", ".b"},
		".a":   {"c", ".d"},
		"e":    {},
	})
	cg := MakeQuotient(MakeQuotientOpt{
		Graph: g,
		Colorer: func(n Node) interface{} {
			return 1
		},
		HasColor: func(n Node) bool {
			s := n.(string)
			return s[0] == '.'
		},
	})
	if n := cg.NumNodes(); n != 4 {
		t.Errorf("expected 4 nodes but found %d", n)
	}
}
