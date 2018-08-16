package graph

import (
	"fmt"
	"math"
	"math/rand"
	"testing"
)

func makeTestColors(m map[string][]int) func(n Node) []int {
	return func(n Node) []int {
		return m[n.(string)]
	}
}

func checkPartition(opt PartitionTreeOpt) (*TreePartition, error) {
	contains := func(xs []int, c int) bool {
		for _, x := range xs {
			if x == c {
				return true
			}
		}
		return false
	}

	r, err := PartitionTree(opt)
	if err != nil {
		return r, err
	}

	if lr, lt := len(r.Parts), opt.Tree.NumNodes(); lr != lt {
		return r, fmt.Errorf("expected %d found %d", lt, lr)
	}

	for i, inum := 0, opt.Tree.NumNodes(); i < inum; i++ {
		n := opt.Tree.Node(i)
		if c, hasColor := r.Parts[n]; !hasColor {
			return r, fmt.Errorf("expecting color for node %q but found none", n)
		} else if !contains(opt.Colors(n), c) {
			return r, fmt.Errorf("node %q assigned wrong color %d", n, c)
		}
	}

	sum := int64(0)
	if err := VisitTree(VisitTreeOpt{
		Tree: opt.Tree,
		Root: opt.Root,
		PreOrder: func(n, parent Node, err error) error {
			if parent != nil {
				sum += opt.EdgeWeight(r.Parts[parent], r.Parts[n])
			}
			return nil
		},
	}); err != nil {
		return nil, err
	}

	if sum != r.Weight {
		return r, fmt.Errorf("expected weight %d found %d", sum, r.Weight)
	}
	return r, nil
}

// Use benchmark infrastructure to figure out the largest problem solvable in
// go test -benchtime TIME (TIME is 1s by default)
func BenchmarkMinWeightTree(b *testing.B) {
	logN := int(math.Log(float64(b.N)))
	if logN < 1 {
		logN = 1
	}

	N := struct{ Tree, Color int }{Tree: logN, Color: (logN+1)/2 + 1}

	var colors []int
	for i := 0; i < N.Color; i++ {
		colors = append(colors, i+1)
	}

	weights := make(map[struct{ A, B int }]int64)
	// Make weight space sparse enough for interesting patterns to occur
	maxWeight := N.Color * N.Color * N.Color * N.Color
	for i := 0; i < N.Color; i++ {
		for j := 0; j < N.Color; j++ {
			k := struct{ A, B int }{A: i + 1, B: j + 1}
			weights[k] = rand.Int63n(int64(maxWeight)) + 1
		}
	}

	makeTree := func(n int) (Graph, Node) {
		m := make(map[string][]string)
		for i := 0; i < n; i++ {
			id := fmt.Sprintf("v%d", i)
			kid1 := fmt.Sprintf("v%d", 2*i+1)
			kid2 := fmt.Sprintf("v%d", 2*i+2)
			m[id] = append(m[id], kid1, kid2)
		}
		return MakeTestGraph(m), "v0"
	}

	graph, root := makeTree(N.Tree)

	opt := PartitionTreeOpt{
		exact: true,
		Tree:  graph,
		Root:  root,
		Colors: func(Node) []int {
			return colors
		},
		EdgeWeight: func(x, y int) int64 {
			if w, seen := weights[struct{ A, B int }{A: x, B: y}]; seen {
				return w
			}

			b.Fatal("unexpected color")
			return 0
		},
	}

	r1, err := checkPartition(opt)
	if err != nil {
		b.Fatal(err)
	}
	opt.exact = false
	r2, err := checkPartition(opt)
	if err != nil {
		b.Fatal(err)
	}

	factor := float64(r2.Weight) / float64(r1.Weight)
	if factor < 1.0 {
		b.Errorf("non-exact solution better than exact")
	}
	b.Logf("Largest problem solvable in benchtime is %+v; SP factor is %fX (min: %d)\n", N, factor, r1.Weight)
}

func TestMinWeightTree(t *testing.T) {
	testMinWeightTree(t, true)
	testMinWeightTree(t, false)
}

func testMinWeightTree(t *testing.T, exact bool) {
	if r, err := checkPartition(PartitionTreeOpt{
		exact: exact,
		Tree: MakeTestGraph(map[string][]string{
			"root": {},
		}),
		Root: "root",
		Colors: makeTestColors(map[string][]int{
			"root": {1},
		}),
		EdgeWeight: func(a, b int) int64 {
			return 1
		},
	}); err != nil {
		t.Fatal(err)
	} else if root := r.Parts["root"]; root != 1 {
		t.Errorf("expected %d but found %d", 1, root)
	}

	if r, err := checkPartition(PartitionTreeOpt{
		exact: exact,
		Tree: MakeTestGraph(map[string][]string{
			"root": {"a", "b"},
			"a":    {"c"},
		}),
		Root: "root",
		Colors: makeTestColors(map[string][]int{
			"root": {1, 3},
			"a":    {2, 3},
			"b":    {1, 3},
			"c":    {4},
		}),
		EdgeWeight: func(a, b int) int64 {
			if a == b {
				return 0
			}
			return 1
		},
	}); err != nil {
		t.Fatal(err)
	} else if !exact {
	} else if root, a, b := r.Parts["root"], r.Parts["a"], r.Parts["b"]; root != 3 || a != 3 || b != 3 {
		t.Errorf("expected %d but found %d %d %d", 3, root, a, b)
	} else if c := r.Parts["c"]; c != 4 {
		t.Errorf("expected %d found %d", 4, c)
	}

	if r, err := checkPartition(PartitionTreeOpt{
		exact: exact,
		Tree: MakeTestGraph(map[string][]string{
			"root": {"a", "b"},
			"a":    {"c", "d"},
			"b":    {"e", "f"},
		}),
		Root: "root",
		Colors: makeTestColors(map[string][]int{
			"root": {1, 2, 3},
			"a":    {1, 2, 3},
			"b":    {1, 2, 3},
			"c":    {1, 2, 3},
			"d":    {1, 2, 3},
			"e":    {1, 2, 3},
			"f":    {1, 2, 3},
		}),
		EdgeWeight: func(a, b int) int64 {
			return 100 - int64(a+b)
		},
	}); err != nil {
		t.Fatal(err)
	} else if !exact {
	} else {
		for n, v := range r.Parts {
			if v != 3 {
				t.Errorf("expected %d found %d for node %q", 3, v, n)
			}
		}
	}
}
