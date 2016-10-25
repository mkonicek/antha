// Package graph provides common graph algorithms
package graph

type Node interface{}

type Graph interface {
	NumNodes() int
	Node(int) Node
	NumOuts(Node) int
	Out(Node, int) Node
}

type LessThan func(Node, Node) bool

type NodeSet interface {
	Has(Node) bool
	Values() []Node
	Len() int
}

type nodeSet map[Node]bool

func (a nodeSet) Has(n Node) bool {
	return a[n]
}

func (a nodeSet) Len() int {
	return len(a)
}

func (a nodeSet) Values() (ret []Node) {
	for k := range a {
		ret = append(ret, k)
	}
	return
}

// Reverse edges
func Reverse(graph Graph) Graph {
	ret := &qgraph{
		Outs: make(map[Node][]Node),
	}
	for i, inum := 0, graph.NumNodes(); i < inum; i += 1 {
		n := graph.Node(i)
		ret.Nodes = append(ret.Nodes, n)
		for j, jnum := 0, graph.NumOuts(n); j < jnum; j += 1 {
			dst := graph.Out(n, j)
			ret.Outs[dst] = append(ret.Outs[dst], n)
		}
	}
	return ret
}

type SimplifyOpt struct {
	Graph            Graph
	RemoveSelfLoops  bool
	RemoveMultiEdges bool
}

// Simplify graph
func Simplify(opt SimplifyOpt) Graph {
	ret := &qgraph{
		Outs: make(map[Node][]Node),
	}
	for i, inum := 0, opt.Graph.NumNodes(); i < inum; i += 1 {
		n := opt.Graph.Node(i)
		ret.Nodes = append(ret.Nodes, n)
		seen := make(map[Node]bool)
		for j, jnum := 0, opt.Graph.NumOuts(n); j < jnum; j += 1 {
			dst := opt.Graph.Out(n, j)

			if opt.RemoveSelfLoops {
				if dst == n {
					continue
				}
			}

			if opt.RemoveMultiEdges {
				if seen[dst] {
					continue
				}
				seen[dst] = true
			}

			ret.Outs[dst] = append(ret.Outs[dst], n)
		}
	}
	return ret
}
