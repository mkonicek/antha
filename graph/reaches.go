package graph

import "errors"

// Merge n and src into dst, modifying dst. Assumes src != dst.
func mergeInto(n Node, src, dst map[Node]bool) {
	dst[n] = true
	for k := range src {
		dst[k] = true
	}
}

//

var (
	reachesSeen = errors.New("reaches already seen")
)

// Compute reachability over graph. A reaches B if there is a path from A to B.
func Reaches(g Graph) Graph {
	reaches := make(map[Node]map[Node]bool)

	// Compute reachability of nodes in turn; reuse reachability of previously
	// processed nodes
	for i, inum := 0, g.NumNodes(); i < inum; i++ {
		sameAs := make(map[Node]bool)
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
					return reachesSeen
				}
				return nil
			},
		})

		rs := make(map[Node]bool)
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
