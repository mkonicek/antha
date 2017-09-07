package graph

import (
	"testing"
)

func TestShortestPaths(t *testing.T) {
	g := MakeTestGraph(map[string][]string{
		"a": []string{"b", "c"},
		"b": []string{"d"},
		"c": []string{"d"},
		"d": []string{"e", "f"},
		"e": []string{"g"},
		"f": []string{"g"},
	})
	type edge struct{ A, B string }
	weights := map[edge]int{
		edge{A: "a", B: "b"}: 1,
		edge{A: "a", B: "c"}: 10,
		edge{A: "b", B: "d"}: 20,
		edge{A: "c", B: "d"}: 1,
		edge{A: "d", B: "e"}: 1,
		edge{A: "d", B: "f"}: 1,
		edge{A: "e", B: "g"}: 10,
		edge{A: "f", B: "g"}: 1,
	}

	edist := map[string]int{
		"a": 0,
		"b": 1,
		"c": 10,
		"d": 11,
		"e": 12,
		"f": 12,
		"g": 13,
	}

	dist := ShortestPath(ShortestPathOpt{
		Graph:   g,
		Sources: []Node{"a"},
		Weight: func(x, y Node) int {
			k := edge{A: x.(string), B: y.(string)}
			return weights[k]
		},
	})
	if nd, ne := len(dist), len(edist); nd != ne {
		t.Errorf("expected %d found %d", ne, nd)
	}
	for k, v := range edist {
		if d, ok := dist[k]; !ok {
			t.Errorf("did not find dist for node %q", k)
		} else if d != v {
			t.Errorf("expected %d for node %q found %d instead", v, k, d)
		}
	}
}
