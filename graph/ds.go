package graph

// A DisjointSet efficiently stores sets of nodes
type DisjointSet struct {
	parent map[Node]Node
	nodes  []Node
}

// NewDisjointSet creates a new disjoint set
func NewDisjointSet() *DisjointSet {
	return &DisjointSet{
		parent: make(map[Node]Node),
	}
}

// NumNodes implements a Graph
func (a *DisjointSet) NumNodes() int {
	return len(a.nodes)
}

// Node implements a Graph
func (a *DisjointSet) Node(i int) Node {
	return a.nodes[i]
}

// NumOuts implements a Graph
func (a *DisjointSet) NumOuts(n Node) int {
	nr := a.Find(n)
	if nr == n {
		return 0
	}
	return 1
}

// Out implements a Graph
func (a *DisjointSet) Out(n Node, i int) Node {
	return a.Find(n)
}

// Union merges to sets
func (a *DisjointSet) Union(x, y Node) {
	xr := a.Find(x)
	yr := a.Find(y)

	a.parent[xr] = yr
}

// Find returns the representative of a set
func (a *DisjointSet) Find(n Node) Node {
	p := a.parent[n]
	if p == nil {
		a.parent[n] = n
		a.nodes = append(a.nodes, n)
		p = n
	}

	if p != n {
		a.parent[n] = a.Find(p)
	}
	return a.parent[n]
}
