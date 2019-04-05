package graph

import "errors"

var (
	errReachesSeen = errors.New("reaches already seen")
)

type Reachability map[Node]map[Node]bool

// NewReachability compute the reachability matrix for the graph
func NewReachability(g Graph) Reachability {
	reaches := make(Reachability, g.NumNodes())

	// Compute reachability of nodes in turn; reuse reachability of previously
	// processed nodes
	for i, inum := 0, g.NumNodes(); i < inum; i++ {
		sameAs := make(map[Node]bool, g.NumNodes())
		root := g.Node(i)
		// Visiting always includes root, differentiate between cyclic and
		// acyclic cases
		rootSeen := false

		vr, _ := Visit(VisitOpt{
			Root:  root,
			Graph: g,
			Seen: func(n Node) error {
				if n == root {
					rootSeen = true
				}
				return nil
			},
			Visitor: func(n Node) error {
				_, seen := reaches[n]
				if seen {
					sameAs[n] = true
					return errReachesSeen
				}
				return nil
			},
		})

		rs := make(map[Node]bool, g.NumNodes())
		for _, v := range vr.Seen.Values() {
			rs[v] = true
		}
		if !rootSeen {
			delete(rs, root)
		}

		for same := range sameAs {
			for v := range reaches[same] {
				rs[v] = true
			}
		}

		reaches[root] = rs
	}

	return reaches
}

// Reaches computes reachability over graph. A reaches B if there is a path
// from A to B.
func Reaches(g Graph) Graph {
	reaches := NewReachability(g)

	ret := &qgraph{
		Outs: make(map[Node][]Node),
	}

	for i, inum := 0, g.NumNodes(); i < inum; i++ {
		node := g.Node(i)
		ret.Nodes = append(ret.Nodes, node)
		for r := range reaches[node] {
			ret.Outs[node] = append(ret.Outs[node], r)
		}
	}

	return ret
}
