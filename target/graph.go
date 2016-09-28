package target

import "github.com/antha-lang/antha/graph"

// View instructions as a graph
type Graph struct {
	Insts []Inst
}

func (a *Graph) NumNodes() int {
	return len(a.Insts)
}

func (a *Graph) Node(i int) graph.Node {
	return a.Insts[i].(graph.Node)
}

func (a *Graph) NumOuts(n graph.Node) int {
	return len(n.(Inst).DependsOn())
}

func (a *Graph) Out(n graph.Node, i int) graph.Node {
	return n.(Inst).DependsOn()[i].(graph.Node)
}
