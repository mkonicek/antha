package codegen

import (
	"fmt"

	"github.com/antha-lang/antha/graph"
	"github.com/antha-lang/antha/laboratory/effects"
)

// Build rooted graph
func makeRoot(nodes []effects.Node) (effects.Node, error) {
	someNode := func(g graph.Graph, m map[graph.Node]bool) graph.Node {
		for i, inum := 0, g.NumNodes(); i < inum; i++ {
			n := g.Node(i)
			if !m[n] {
				return n
			}
		}
		return nil
	}

	g := effects.ToGraph(effects.ToGraphOpt{
		Roots: nodes,
	})

	roots := graph.Schedule(g).Roots
	seen := make(map[graph.Node]bool)
	for _, root := range roots {
		results, _ := graph.Visit(graph.VisitOpt{
			Graph: g,
			Root:  root,
			Visitor: func(n graph.Node) error {
				if seen[n] {
					return graph.ErrNextNode
				}
				return nil
			},
		})
		for _, k := range results.Seen.Values() {
			seen[k] = true
		}
	}

	// If some nodes are not reachable from roots, there must be a cycle
	if len(seen) != g.NumNodes() {
		return nil, fmt.Errorf("cycle containing %T", someNode(g, seen))
	}

	ret := &effects.Bundle{}
	for _, r := range roots {
		ret.From = append(ret.From, r.(effects.Node))
	}
	return ret, nil
}

// What is the set of UseComps that reach each command
func buildReachingUses(g graph.Graph) map[effects.Node][]*effects.UseComp {
	// Simple fixpoint:
	//   Value: set of use comps,
	//   Merge: union
	//   Transfer functions:
	//     - Command c -> { }
	//     - UseComp u -> {u}

	values := make(map[effects.Node][]*effects.UseComp)

	merge := func(n effects.Node) []*effects.UseComp {
		var vs []*effects.UseComp
		for i, inum := 0, g.NumOuts(n); i < inum; i++ {
			pred := g.Out(n, i).(effects.Node)
			switch pred := pred.(type) {
			case *effects.Command:
				// Kill
			case *effects.UseComp:
				vs = append(vs, values[pred]...)
				vs = append(vs, pred)
			default:
				vs = append(vs, values[pred]...)
			}
		}
		return vs
	}

	dag := graph.Schedule(graph.Reverse(g))

	for len(dag.Roots) > 0 {
		var next []graph.Node
		for _, n := range dag.Roots {
			n := n.(effects.Node)
			seen := make(map[*effects.UseComp]bool)

			for _, v := range merge(n) {
				if seen[v] {
					continue
				}
				seen[v] = true
				values[n] = append(values[n], v)
			}

			next = append(next, dag.Visit(n)...)
		}

		dag.Roots = next
	}

	return values
}

// Eliminate nodes while preserving dependency relation
func simplifyWithDeps(g graph.Graph, in func(n graph.Node) bool) (graph.Graph, error) {
	rg := graph.Reaches(graph.Simplify(graph.SimplifyOpt{
		Graph:            g,
		RemoveSelfLoops:  true,
		RemoveMultiEdges: true,
	}))

	rg = graph.Simplify(graph.SimplifyOpt{
		Graph: rg,
		RemoveNodes: func(n graph.Node) bool {
			return !in(n)
		},
	})

	return graph.TransitiveReduction(rg)
}

// Build IR
func build(root effects.Node) (*ir, error) {
	g := effects.ToGraph(effects.ToGraphOpt{
		Roots: []effects.Node{root},
	})

	// Remove UseComps primarily. They may be locally cyclic.
	ct, err := simplifyWithDeps(g, func(n graph.Node) bool {
		c, ok := n.(*effects.Command)
		return (ok && c.Output == nil) || n == root
	})

	if err != nil {
		return nil, err
	}

	// TODO: Add back some validity checks like the same UseComp cannot be used
	// multiple times

	return &ir{
		Root:         root,
		Graph:        g,
		Commands:     ct,
		reachingUses: buildReachingUses(g),
	}, nil
}
