package graph

// A Visitor is a function over graph nodes
type Visitor func(Node) error

// VisitOpt is a set of options to Visit
type VisitOpt struct {
	Root         Node    // Root of traversal
	Graph        Graph   // Graph to traverse
	Visitor      Visitor // Function to apply when node first encountered
	Seen         Visitor // Function applied on subsequent visits to a node
	BreadthFirst bool    // Visit nodes breadth first
}

// A VisitResult is the result of Visit
type VisitResult struct {
	Seen      NodeSet
	Frontiers []NodeSet // If VisitOpt.BreadthFirst, successive frontiers are placed here
}

type dists map[Node]int

func (a dists) Has(n Node) bool {
	_, seen := a[n]
	return seen
}

func (a dists) Len() int {
	return len(a)
}

func (a dists) Values() (ret []Node) {
	for k := range a {
		ret = append(ret, k)
	}
	return
}

// Visit applies a visitor to each node reachable from root in some order.
// Returns nodes visited. If visitor returns an error, stop traversal early and
// pass returned error.
func Visit(opt VisitOpt) (res *VisitResult, err error) {
	apply := func(v Visitor, n Node) error {
		if v != nil {
			return v(n)
		}
		return nil
	}

	type pair struct {
		Node Node
		Dist int
	}

	dists := make(dists)

	maxDist := 0
	wl := []pair{{opt.Root, maxDist}}
	for l := len(wl); l > 0; l = len(wl) {
		var p pair
		if opt.BreadthFirst {
			p = wl[0]
			wl = wl[1:]
		} else {
			p = wl[l-1]
			wl = wl[:l-1]
		}

		if _, seen := dists[p.Node]; seen {
			if err = apply(opt.Seen, p.Node); err != nil {
				break
			}
			continue
		}

		dists[p.Node] = p.Dist

		if p.Dist > maxDist {
			maxDist = p.Dist
		}

		if err = apply(opt.Visitor, p.Node); err != nil {
			if err == ErrTraversalDone {
				break
			} else {
				continue
			}
		}

		nextDist := p.Dist + 1

		for i, num := 0, opt.Graph.NumOuts(p.Node); i < num; i++ {
			wl = append(wl, pair{opt.Graph.Out(p.Node, i), nextDist})
		}
	}

	var frontiers []NodeSet
	if opt.BreadthFirst {
		fs := make([][]Node, maxDist+1)
		for n, d := range dists {
			fs[d] = append(fs[d], n)
		}
		for _, f := range fs {
			ns := make(nodeSet)
			for _, v := range f {
				ns[v] = true
			}
			frontiers = append(frontiers, ns)
		}
	}

	return &VisitResult{Seen: dists, Frontiers: frontiers}, err
}
