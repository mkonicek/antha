package graph

import (
	"testing"
)

func TestShortestPaths(t *testing.T) {
	g := MakeTestGraph(map[string][]string{
		"a": {"b", "c"},
		"b": {"d"},
		"c": {"d"},
		"d": {"e", "f"},
		"e": {"g"},
		"f": {"g"},
	})
	type edge struct{ A, B string }
	weights := map[edge]int64{
		{A: "a", B: "b"}: 1,
		{A: "a", B: "c"}: 10,
		{A: "b", B: "d"}: 20,
		{A: "c", B: "d"}: 1,
		{A: "d", B: "e"}: 1,
		{A: "d", B: "f"}: 1,
		{A: "e", B: "g"}: 10,
		{A: "f", B: "g"}: 1,
	}

	edist := map[string]int64{
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
		Weight: func(x, y Node) int64 {
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
