// Package graph provides common graph algorithms
package graph

import "errors"

var (
	errDuplicateNode = errors.New("duplicate node")
	errOutNotInGraph = errors.New("out not in graph")
)

// A Node is a node in a graph
type Node interface{}

// A Graph is a relation between nodes
type Graph interface {
	NumNodes() int
	Node(int) Node
	NumOuts(Node) int
	Out(Node, int) Node
}

// LessThan compares nodes
type LessThan func(Node, Node) bool

// A NodeSet is a unordered collection of nodes
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
	for i, inum := 0, graph.NumNodes(); i < inum; i++ {
		n := graph.Node(i)
		ret.Nodes = append(ret.Nodes, n)
		for j, jnum := 0, graph.NumOuts(n); j < jnum; j++ {
			dst := graph.Out(n, j)
			ret.Outs[dst] = append(ret.Outs[dst], n)
		}
	}
	return ret
}

// A SimplifyOpt are options to Simplify
type SimplifyOpt struct {
	Graph            Graph
	RemoveSelfLoops  bool
	RemoveMultiEdges bool
	RemoveNodes      func(Node) bool // Should node be removed
}

// Simplify graph
func Simplify(opt SimplifyOpt) Graph {
	ret := &qgraph{
		Outs: make(map[Node][]Node),
	}

	remove := make(map[Node]bool)
	if opt.RemoveNodes != nil {
		for i, inum := 0, opt.Graph.NumNodes(); i < inum; i++ {
			n := opt.Graph.Node(i)
			remove[n] = opt.RemoveNodes(n)
		}
	}

	for i, inum := 0, opt.Graph.NumNodes(); i < inum; i++ {
		n := opt.Graph.Node(i)
		if remove[n] {
			continue
		}

		ret.Nodes = append(ret.Nodes, n)
		seen := make(map[Node]bool)
		for j, jnum := 0, opt.Graph.NumOuts(n); j < jnum; j++ {
			dst := opt.Graph.Out(n, j)
			if remove[dst] {
				continue
			}

			if opt.RemoveSelfLoops && dst == n {
				continue
			}

			if opt.RemoveMultiEdges {
				if seen[dst] {
					continue
				}
				seen[dst] = true
			}

			ret.Outs[n] = append(ret.Outs[n], dst)
		}
	}
	return ret
}

// Verify returns an error in a graph doesn't satisfy some basic consistent
// properties.
func Verify(graph Graph) error {
	seen := make(map[Node]bool)
	for i, inum := 0, graph.NumNodes(); i < inum; i++ {
		node := graph.Node(i)
		if seen[node] {
			return errDuplicateNode
		}
		seen[node] = true
	}

	for i, inum := 0, graph.NumNodes(); i < inum; i++ {
		node := graph.Node(i)
		for j, jnum := 0, graph.NumOuts(node); j < jnum; j++ {
			dst := graph.Out(node, j)
			if !seen[dst] {
				return errOutNotInGraph
			}
		}
	}
	return nil
}
