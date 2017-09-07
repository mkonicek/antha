package graph

// QGraph is a quotient graph
type QGraph interface {
	Graph
	OrigGraph() Graph    // Original graph that this graph was constructed from
	NumOrigs(Node) int   // Number of nodes in the original graph that this node represents
	Orig(Node, int) Node // The ith node in the original graph that this node represents
}

// Internal representation for a quotient graph and for generic graph return values
type qgraph struct {
	Graph Graph
	Nodes []Node
	Outs  map[Node][]Node
	Origs map[Node][]Node
}

func (a *qgraph) NumNodes() int {
	return len(a.Nodes)
}

func (a *qgraph) Node(i int) Node {
	return a.Nodes[i]
}

func (a *qgraph) NumOuts(n Node) int {
	return len(a.Outs[n])
}

func (a *qgraph) Out(n Node, i int) Node {
	return a.Outs[n][i]
}

func (a *qgraph) OrigGraph() Graph {
	return a.Graph
}

func (a *qgraph) NumOrigs(n Node) int {
	if a.Origs != nil {
		return len(a.Origs[n])
	}
	return 1
}

func (a *qgraph) Orig(n Node, i int) Node {
	if a.Origs != nil {
		return a.Origs[n][i]
	}
	return n
}

// MakeQuotientOpt are options to MakeQuotient
type MakeQuotientOpt struct {
	Graph         Graph
	Colorer       func(Node) interface{}
	HasColor      func(Node) bool
	Present       func(Node) bool // Should n be included at all
	KeepSelfEdges bool
}

// MakeQuotient returns a quotient graph. Nodes with the same color merged into
// a single node. A colorless node is treated as having a color distinct from
// any other node.
func MakeQuotient(opt MakeQuotientOpt) QGraph {
	ret := &qgraph{
		Graph: opt.Graph,
		Outs:  make(map[Node][]Node),
		Origs: make(map[Node][]Node),
	}

	cnodes := make(map[interface{}]Node)
	newNodes := make(map[Node]Node)

	for i, inum := 0, opt.Graph.NumNodes(); i < inum; i++ {
		node := opt.Graph.Node(i)

		switch {
		case opt.Present != nil && !opt.Present(node):
			continue
		case opt.Colorer == nil:
			fallthrough
		case opt.HasColor != nil && !opt.HasColor(node):
			// Add node as itself
			newNode := i
			newNodes[node] = newNode
			ret.Origs[newNode] = append(ret.Origs[newNode], node)
		default:
			// Coarsen node with those of the same color
			c := opt.Colorer(node)
			newNode, ok := cnodes[c]
			if !ok {
				newNode = i
				cnodes[c] = newNode
			}
			newNodes[node] = newNode
			ret.Origs[newNode] = append(ret.Origs[newNode], node)
		}
	}

	for node, newNode := range newNodes {
		neighs := make(map[Node]bool)
		for i, inum := 0, opt.Graph.NumOuts(node); i < inum; i++ {
			n := opt.Graph.Out(node, i)
			// Filter out not present neighbors
			if k, ok := newNodes[n]; ok {
				neighs[k] = true
			}
		}

		for n := range neighs {
			if !opt.KeepSelfEdges && n == newNode {
				continue
			}
			ret.Outs[newNode] = append(ret.Outs[newNode], n)
		}
	}

	for newNode := range ret.Origs {
		ret.Nodes = append(ret.Nodes, newNode)
	}

	return ret
}
