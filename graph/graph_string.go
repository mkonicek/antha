package graph

// A StringGraph is a graph where nodes are strings
type StringGraph struct {
	Nodes []string
	Outs  map[string][]string
}

// NumNodes implements a Graph
func (a *StringGraph) NumNodes() int {
	return len(a.Nodes)
}

// Node implements a Graph
func (a *StringGraph) Node(i int) Node {
	return a.Nodes[i]
}

// NumOuts implements a Graph
func (a *StringGraph) NumOuts(n Node) int {
	return len(a.Outs[n.(string)])
}

// Out implements a Graph
func (a *StringGraph) Out(n Node, i int) Node {
	return a.Outs[n.(string)][i]
}
