package graph

import "testing"

func TestTreeReaches(t *testing.T) {
	g := MakeTestGraph(map[string][]string{
		"root": []string{"a", "b"},
		"a":    []string{"c", "d"},
		"b":    []string{"e"},
		"e":    []string{"f", "g"},
	})

	expected := map[string]int{
		"root": 7,
		"a":    2,
		"b":    3,
		"e":    2,
	}

	res := Reaches(g)

	if e, f := g.NumNodes()-1, expected["root"]; e != f {
		t.Errorf("expected %q to reach %d nodes found %d in testdata", "root", e, f)
	}

	for i, inum := 0, res.NumNodes(); i < inum; i++ {
		node := res.Node(i)
		e := expected[node.(string)]
		f := res.NumOuts(node)

		if e != f {
			t.Errorf("expected %q to reach %d nodes found %d", node, e, f)
		}
	}
}

func TestLoopReaches(t *testing.T) {
	g := MakeTestGraph(map[string][]string{
		"a": []string{"b"},
		"b": []string{"c"},
		"c": []string{"a"},
	})

	expected := map[string]int{
		"a": 3,
		"b": 3,
		"c": 3,
	}

	res := Reaches(g)

	if e, f := g.NumNodes(), expected["a"]; e != f {
		t.Errorf("expected %q to reach %d nodes found %d in testdata", "a", e, f)
	}

	for i, inum := 0, res.NumNodes(); i < inum; i++ {
		node := res.Node(i)
		e := expected[node.(string)]
		f := res.NumOuts(node)

		if e != f {
			t.Errorf("expected %q to reach %d nodes found %d", node, e, f)
		}
	}
}
