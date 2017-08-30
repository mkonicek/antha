package graph

// An SDag is the current state of DAG scheduling
type SDag struct {
	Graph Graph
	Roots []Node
	Ins   map[Node]int
}

// Visit marks a node as visited and returns nodes that can be visited next
func (a *SDag) Visit(n Node) (r []Node) {
	for i, inum := 0, a.Graph.NumOuts(n); i < inum; i++ {
		succ := a.Graph.Out(n, i)
		a.Ins[succ]--
		if a.Ins[succ] == 0 {
			r = append(r, succ)
		}
	}
	return
}

// Schedule treats directed acyclic graph as a dependency graph and schedules
// nodes to execute
func Schedule(graph Graph) *SDag {
	dag := &SDag{
		Graph: graph,
		Ins:   make(map[Node]int),
	}
	for i, inum := 0, dag.Graph.NumNodes(); i < inum; i++ {
		n := dag.Graph.Node(i)
		for j, jnum := 0, dag.Graph.NumOuts(n); j < jnum; j++ {
			dst := dag.Graph.Out(n, j)
			dag.Ins[dst]++
		}
	}
	for i, inum := 0, dag.Graph.NumNodes(); i < inum; i++ {
		n := dag.Graph.Node(i)
		if dag.Ins[n] == 0 {
			dag.Roots = append(dag.Roots, n)
		}
	}
	return dag
}
