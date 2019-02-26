package graph

import (
	"testing"
)

func TestVisit(t *testing.T) {
	g := MakeTestGraph(map[string][]string{
		"root": {"a", "b"},
		"a":    {"c", "d"},
		"e":    {},
	})
	if res, err := Visit(VisitOpt{
		Graph: g,
		Root:  "root",
	}); err != nil {
		t.Fatal(err)
	} else if !res.Seen.Has("d") {
		t.Error("did not visit node d")
	} else if res.Seen.Has("e") {
		t.Error("visited node e")
	}
}

func TestVisitBfs(t *testing.T) {
	g := MakeTestGraph(map[string][]string{
		"root": {"a", "b", "c"},
		"a":    {"c"},
		"b":    {"c"},
	})
	if res, err := Visit(VisitOpt{
		Graph:        g,
		Root:         "root",
		BreadthFirst: true,
	}); err != nil {
		t.Fatal(err)
	} else if l := len(res.Frontiers); l != 2 {
		t.Errorf("expecting 2 levels found %d", l)
	} else if l := res.Frontiers[0].Len(); l != 1 {
		t.Errorf("expecting 1 node at level 0 found %d", l)
	} else if l := res.Frontiers[1].Len(); l != 3 {
		t.Errorf("expecting 3 nodes at level 1 found %d", l)
	}
}
