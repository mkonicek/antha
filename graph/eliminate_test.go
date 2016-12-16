package graph

import "testing"

func TestTreeEliminate(t *testing.T) {
	g := MakeTestGraph(map[string][]string{
		"root": []string{"a", "b"},
		"a":    []string{"c", "d"},
		"b":    []string{"e"},
		"e":    []string{"f", "g"},
	})

	in := map[string]bool{
		"root": true,
		"b":    true,
		"f":    true,
		"g":    true,
	}

	gnext := Eliminate(EliminateOpt{
		Graph: g,
		In: func(n Node) bool {
			return in[n.(string)]
		},
	})

	if l := gnext.NumNodes(); l != 4 {
		t.Errorf("expected %d nodes found %d", 4, l)
	} else if l := gnext.NumOuts("root"); l != 1 {
		t.Errorf("expected %d nodes found %d", 1, l)
	} else if n := gnext.Out("root", 0).(string); n != "b" {
		t.Errorf("expected %q found %q", "b", n)
	} else if l := gnext.NumOuts("b"); l != 2 {
		t.Errorf("expected %d nodes found %d", 2, l)
	} else if n := gnext.Out("b", 0).(string); n != "f" && n != "g" {
		t.Errorf("expected %q or %q found %q", "f", "g", n)
	} else if n := gnext.Out("b", 1).(string); n != "f" && n != "g" {
		t.Errorf("expected %q or %q found %q", "f", "g", n)
	}
}

func TestGraphEliminate(t *testing.T) {
	g := MakeTestGraph(map[string][]string{
		"a": []string{"c"},
		"b": []string{"c"},
		"c": []string{"d"},
		"d": []string{"e"},
		"e": []string{"f"},
		"f": []string{"g", "h"},
	})

	out := map[string]bool{
		"c": true,
		"d": true,
		"e": true,
		"f": true,
	}

	gnext := Eliminate(EliminateOpt{
		Graph: g,
		In: func(n Node) bool {
			return !out[n.(string)]
		},
	})

	if l := gnext.NumNodes(); l != 4 {
		t.Errorf("expected %d nodes found %d", 4, l)
	} else if l := gnext.NumOuts("a"); l != 2 {
		t.Errorf("expected %d nodes found %d", 2, l)
	} else if l := gnext.NumOuts("b"); l != 2 {
		t.Errorf("expected %d nodes found %d", 2, l)
	}
}

func TestHarderGraphEliminate(t *testing.T) {
	notIns := []string{
		"v1",
		"v2",
		"v3",
		"v5",
		"v6",
		"v7",
		"v12",
		"v13",
		"v14",
		"v15",
		"v19",
		"v20",
		"v23",
		"v24",
		"v25",
		"v26",
		"v27",
		"v28",
		"v31",
		"v32",
		"v33",
		"v34",
		"v35",
		"v40",
		"v41",
		"v42",
		"v43",
		"v44",
		"v45",
		"v47",
	}
	notInMap := make(map[string]bool)
	for _, v := range notIns {
		notInMap[v] = true
	}
	g := MakeTestGraph(map[string][]string{
		"v0":  []string{"v47", "v23", "v7", "v35"},
		"v1":  []string{"v38", "v5", "v28"},
		"v2":  []string{"v44"},
		"v3":  []string{"v26"},
		"v4":  []string{"v42", "v15"},
		"v5":  []string{"v8", "v2"},
		"v6":  []string{"v29", "v2", "v45"},
		"v7":  []string{"v17", "v34"},
		"v8":  []string{"v2"},
		"v9":  []string{"v27"},
		"v10": []string{"v20"},
		"v11": []string{"v44"},
		"v12": []string{"v16", "v31"},
		"v13": []string{"v33"},
		"v14": []string{"v43"},
		"v15": []string{"v28"},
		"v16": []string{"v40"},
		"v17": []string{"v34"},
		"v18": []string{"v31"},
		"v19": []string{"v9", "v27"},
		"v20": []string{"v32"},
		"v21": []string{"v1", "v14"},
		"v22": []string{"v13"},
		"v23": []string{"v37", "v24"},
		"v24": []string{"v18", "v31"},
		"v25": []string{"v36", "v31", "v3"},
		"v26": []string{"v46", "v3", "v41", "v14"},
		"v27": []string{"v20"},
		"v28": []string{"v15"},
		"v29": []string{"v45"},
		"v30": []string{"v12"},
		"v31": []string{"v40"},
		"v32": []string{"v22", "v27", "v13"},
		"v33": []string{"v21", "v13", "v1", "v14"},
		"v34": []string{"v10", "v27"},
		"v35": []string{"v39", "v19"},
		"v36": []string{"v3"},
		"v37": []string{"v24"},
		"v38": []string{"v5", "v28"},
		"v39": []string{"v19"},
		"v40": []string{"v25"},
		"v41": []string{"v4", "v42", "v28"},
		"v42": []string{"v11", "v2"},
		"v43": []string{"v14"},
		"v44": []string{"v6"},
		"v45": []string{},
		"v46": []string{"v41", "v43"},
		"v47": []string{"v30", "v12"},
	})
	eg := Eliminate(EliminateOpt{
		Graph: g,
		In: func(x Node) bool {
			return !notInMap[x.(string)]
		},
	})

	tr, err := TransitiveReduction(eg)
	if err != nil {
		t.Fatal(err)
	}

	root := "v0"
	depth := 7
	dist := ShortestPath(ShortestPathOpt{
		Graph:   tr,
		Sources: []Node{root},
		Weight: func(x, y Node) int {
			return 1
		},
	})
	for i, inum := 0, eg.NumNodes(); i < inum; i += 1 {
		n := eg.Node(i)
		if eg.NumOuts(n) != 0 {
			continue
		}

		if dist[n] != depth {
			t.Errorf("%s -> %s length expected %d found %d", root, n, depth, dist[n])
		}
	}
}
