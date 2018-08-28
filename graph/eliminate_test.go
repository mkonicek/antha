package graph

import "testing"

func TestTreeEliminate(t *testing.T) {
	g := MakeTestGraph(map[string][]string{
		"root": {"a", "b"},
		"a":    {"c", "d"},
		"b":    {"e"},
		"e":    {"f", "g"},
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
		"a": {"c"},
		"b": {"c"},
		"c": {"d"},
		"d": {"e"},
		"e": {"f"},
		"f": {"g", "h"},
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
		"v0":  {"v47", "v23", "v7", "v35"},
		"v1":  {"v38", "v5", "v28"},
		"v2":  {"v44"},
		"v3":  {"v26"},
		"v4":  {"v42", "v15"},
		"v5":  {"v8", "v2"},
		"v6":  {"v29", "v2", "v45"},
		"v7":  {"v17", "v34"},
		"v8":  {"v2"},
		"v9":  {"v27"},
		"v10": {"v20"},
		"v11": {"v44"},
		"v12": {"v16", "v31"},
		"v13": {"v33"},
		"v14": {"v43"},
		"v15": {"v28"},
		"v16": {"v40"},
		"v17": {"v34"},
		"v18": {"v31"},
		"v19": {"v9", "v27"},
		"v20": {"v32"},
		"v21": {"v1", "v14"},
		"v22": {"v13"},
		"v23": {"v37", "v24"},
		"v24": {"v18", "v31"},
		"v25": {"v36", "v31", "v3"},
		"v26": {"v46", "v3", "v41", "v14"},
		"v27": {"v20"},
		"v28": {"v15"},
		"v29": {"v45"},
		"v30": {"v12"},
		"v31": {"v40"},
		"v32": {"v22", "v27", "v13"},
		"v33": {"v21", "v13", "v1", "v14"},
		"v34": {"v10", "v27"},
		"v35": {"v39", "v19"},
		"v36": {"v3"},
		"v37": {"v24"},
		"v38": {"v5", "v28"},
		"v39": {"v19"},
		"v40": {"v25"},
		"v41": {"v4", "v42", "v28"},
		"v42": {"v11", "v2"},
		"v43": {"v14"},
		"v44": {"v6"},
		"v45": {},
		"v46": {"v41", "v43"},
		"v47": {"v30", "v12"},
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
	depth := int64(7)
	dist := ShortestPath(ShortestPathOpt{
		Graph:   tr,
		Sources: []Node{root},
		Weight: func(x, y Node) int64 {
			return 1
		},
	})
	for i, inum := 0, eg.NumNodes(); i < inum; i++ {
		n := eg.Node(i)
		if eg.NumOuts(n) != 0 {
			continue
		}

		if dist[n] != depth {
			t.Errorf("%s -> %s length expected %d found %d", root, n, depth, dist[n])
		}
	}
}
