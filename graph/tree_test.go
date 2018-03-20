package graph

import (
	"testing"
)

func TestIsTree(t *testing.T) {
	g := MakeTestGraph(map[string][]string{
		"root": {"a", "b"},
	})
	if err := IsTree(g, "root"); err != nil {
		t.Fatal(err)
	}
}

func TestIsNotTree(t *testing.T) {
	g := MakeTestGraph(map[string][]string{
		"root": {"a", "b"},
		"a":    {"c", "d"},
		"b":    {"d"},
	})
	if err := IsTree(g, "root"); err == nil {
		t.Fatal("found tree but did not expect one")
	}
}

func TestIsNotTreeWhenHasCycle(t *testing.T) {
	g := MakeTestGraph(map[string][]string{
		"root": {"a", "b"},
		"a":    {"c", "root"},
	})
	if err := IsTree(g, "root"); err == nil {
		t.Fatal("found tree but did not expect one")
	}
}

func TestTreeVisit(t *testing.T) {
	g := MakeTestGraph(map[string][]string{
		"root": {"a", "b"},
		"a":    {"c", "d"},
		"b":    {"e"},
	})
	v := map[string]int{
		"c": 4,
		"d": 5,
		"e": 1,
	}
	if err := VisitTree(VisitTreeOpt{
		Tree: g,
		Root: "root",
		PostOrder: func(n, parent Node, err error) error {
			if g.NumOuts(n) == 0 {
				return nil
			}
			sum := 0
			for i, inum := 0, g.NumOuts(n); i < inum; i++ {
				k := g.Out(n, i).(string)
				sum += v[k]
			}
			if sum > 0 {
				v[n.(string)] = sum
			}
			return nil
		},
	}); err != nil {
		t.Fatal(err)
	} else if vv := v["c"]; vv != 4 {
		t.Errorf("leaf value changed: expected %d found %d", 4, vv)
	} else if vv := v["d"]; vv != 5 {
		t.Errorf("leaf value changed: expected %d found %d", 5, vv)
	} else if vv := v["e"]; vv != 1 {
		t.Errorf("leaf value changed: expected %d found %d", 1, vv)
	} else if vv := v["a"]; vv != 9 {
		t.Errorf("expected %d found %d", 9, vv)
	} else if vv := v["b"]; vv != 1 {
		t.Errorf("expected %d found %d", 1, vv)
	} else if vv := v["root"]; vv != 10 {
		t.Errorf("expected %d found %d", 10, vv)
	}
}
