package graph

import (
	"container/heap"
	"errors"
	"math"
)

var (
	errMissingColors = errors.New("missing color")
)

// PartitionTreeOpt is a set of options for PartitionTree
type PartitionTreeOpt struct {
	Tree       Graph
	Root       Node
	Colors     func(n Node) []int   // Possible colors n may take; all nodes in graph must have at least one color
	EdgeWeight func(a, b int) int64 // Weight between a and b (non-negative); depth(a) < depth(b)
	NodeWeight func(a int) int      // Weight of a
	exact      bool
}

// A TreePartition is a partition of the nodes of a tree
type TreePartition struct {
	Parts  map[Node]int
	Weight int64
}

// PartitionTree partitions a tree. Given a tree and a set of colors that each
// node may take, select a color for each node that minimizes the weight of the
// tree.
func PartitionTree(opt PartitionTreeOpt) (*TreePartition, error) {
	// Cache colors
	colors := make(map[Node][]int)
	for i, inum := 0, opt.Tree.NumNodes(); i < inum; i++ {
		n := opt.Tree.Node(i)
		cs := opt.Colors(n)
		if len(cs) == 0 {
			return nil, errMissingColors
		}
		colors[n] = cs
	}

	algo := &partitionTree{
		colors: colors,
	}

	if opt.exact {
		return algo.runExact(opt), nil
	}
	return algo.runSP(opt), nil
}

type partitionTree struct {
	edgeWeight map[struct{ A, B int }]int64
	colors     map[Node][]int
}

// Cache results of weight function. We'll be calling it a lot.
func (a *partitionTree) getW(opt PartitionTreeOpt, x, y int) int64 {
	if a.edgeWeight == nil {
		a.edgeWeight = make(map[struct{ A, B int }]int64)
	}
	p := struct{ A, B int }{A: x, B: y}
	v, seen := a.edgeWeight[p]
	if !seen {
		v = opt.EdgeWeight(x, y)
		a.edgeWeight[p] = v
	}
	return v
}

// Non-optional solution by greedily extending optimal, shortest-paths based
// solution on line graphs.
func (a *partitionTree) runSP(opt PartitionTreeOpt) *TreePartition {
	// Sketch: If the input graph were a line, the optimal solution could be
	// found with shortest paths.
	//
	// Example: Graph where labels are possible colors:
	//
	//   [c, d] - [d, e] - [f, g]
	//
	// Finding the shortest path in this graph, gives the minimium assignment
	// (where xx is the product of weights from one column to another):
	//
	//   (c) x (d) x (f)
	//   (d) x (e) x (g)
	//
	// This algorithm finds the shortest path on a tree and uses that to
	// greedily assign colors for that path. It then recursively solves the
	// remaining subproblems.

	type node struct {
		Node  Node
		Index int
	}

	nodes := make(map[Node][]*node)
	item := make(map[*node]*nodeItem)
	dist := make(map[*node]int64) // Best distance so far (zero means unvisited)
	best := make(map[*node]*node) // Corresponding best parent
	ret := &TreePartition{Parts: make(map[Node]int)}

	var pq priorityQueue

	enqueue := func(n *node, priority int64) {
		if ni, seen := item[n]; !seen {
			ni := &nodeItem{
				Priority: priority,
				Value:    n,
			}
			item[n] = ni
			heap.Push(&pq, ni)
		} else {
			pq.Update(ni, priority)
		}
	}

	// Node relaxation for shortest paths: update shortest path to here (kid)
	// and update priority queue
	visitKid := func(kid Node, parent *node, c int, d int64) {
		for idx, kc := range a.colors[kid] {
			kn := nodes[kid][idx]
			kd := dist[kn]
			newD := d + a.getW(opt, c, kc)
			if kd != 0 && kd <= newD {
				continue
			}
			dist[kn] = newD
			best[kn] = parent
			enqueue(kn, newD)
		}
	}

	// Node relaxation for shortest paths
	visit := func(n *node) {
		c := a.colors[n.Node][n.Index]
		d := dist[n]
		for i, inum := 0, opt.Tree.NumOuts(n.Node); i < inum; i++ {
			kid := opt.Tree.Out(n.Node, i)
			if _, seen := ret.Parts[kid]; seen {
				continue
			}
			visitKid(kid, n, c, d)
		}
	}

	dequeue := func(n *node) {
		ni := item[n]
		if ni != nil {
			delete(item, n)
			heap.Remove(&pq, ni.Index)
		}
	}

	// Reset whole Node subtrees at a time.
	resetSubtree := func(root Node) {
		if err := VisitTree(VisitTreeOpt{
			Tree: opt.Tree,
			Root: root,
			PreOrder: func(n, parent Node, err error) error {
				for _, n := range nodes[n] {
					best[n] = nil
					dist[n] = 0
					dequeue(n)
				}
				return nil
			},
		}); err != nil {
			panic(err)
		}
	}

	sameBest := func(kid Node, b *node) bool {
		for _, n := range nodes[kid] {
			if best[n] != nil && best[n] != b {
				return false
			}
		}
		return true
	}

	// Greedily assign colors up to root; reset distances of invalidated subtrees
	assignParents := func(n *node) {
		for ; n != nil; n = best[n] {
			if _, seen := ret.Parts[n.Node]; seen {
				return
			}
			// Assign color
			ret.Parts[n.Node] = a.colors[n.Node][n.Index]
			if p := best[n]; p != nil {
				ret.Weight += dist[n] - dist[p]
			}
			for _, x := range nodes[n.Node] {
				// Squelch other paths
				if x != n {
					delete(best, x)
					delete(dist, x)
					dequeue(x)
				}
			}
			// Kick off other subtrees
			for i, inum := 0, opt.Tree.NumOuts(n.Node); i < inum; i++ {
				kid := opt.Tree.Out(n.Node, i)
				if _, seen := ret.Parts[kid]; seen {
					continue
				} else if sameBest(kid, n) {
					continue
				}

				resetSubtree(kid)
			}
			visit(n)
		}
	}

	for i, inum := 0, opt.Tree.NumNodes(); i < inum; i++ {
		n := opt.Tree.Node(i)
		for idx := range a.colors[n] {
			nodes[n] = append(nodes[n], &node{Node: n, Index: idx})
		}
	}

	for _, n := range nodes[opt.Root] {
		dist[n] = 1 // start at 1 so zero means unvisited
		enqueue(n, 1)
	}

	for pq.Len() > 0 {
		ni := heap.Pop(&pq).(*nodeItem)
		n := ni.Value.(*node)
		delete(item, n)

		if opt.Tree.NumOuts(n.Node) == 0 {
			// Reached leaf
			assignParents(n)
		} else {
			visit(n)
		}
	}

	return ret
}

// Uses an exact but exponential algorithm, so only use when number of colors
// (C) and size of tree (T) is small (C, T < 15).
func (a *partitionTree) runExact(opt PartitionTreeOpt) *TreePartition {
	// Sketch: Expand out (exponential) search tree.
	//
	// Example: Graph where labels are possible colors:
	//
	//   [c, d]
	//     |   \
	//   [d, e] [f, g]
	//
	// Partial of search tree [tree nodes], (choice nodes):
	//
	//   [troot]
	//      |   \
	//     (c)  (d)
	//      | \    \...
	//    [n1] [n2]
	//    / |   | \
	//  (d)(e) (f)(g)
	//
	//  n1 will eventually contain the least weight subtree w(cd) + w(...) or w(ce) + w(...).

	// Node in search tree. There are two types "tree" nodes, which have
	// non-nil Node fields, and "choice" nodes.  Tree nodes have choice node as
	// children and vice versa.
	type node struct {
		Parent *node
		Node   Node  // Corresponding node in input tree
		Weight int64 // Tree Node: least weight subtree below here
		Choice int   // Choice this represents
		Kids   []*node
	}

	weighCNode := func(cnode *node) (sum int64) {
		for _, tnode := range cnode.Kids {
			sum += tnode.Weight
		}
		return
	}

	// Compute lighest weight subtree rooted at a tree node
	getBestColor := func(n *node) (idx int, weight int64) {
		parent := n.Parent
		minW := int64(math.MaxInt64)
		minKid := -1
		for i, kid := range n.Kids {
			w := weighCNode(kid)
			if parent != nil {
				w += a.getW(opt, parent.Choice, kid.Choice)
			}
			if w < minW {
				minW = w
				minKid = i
			}
		}
		idx = minKid
		weight = minW
		return
	}

	// Walk tree top down, following best edges to extract solution
	extractSolution := func(root *node) *TreePartition {
		r := &TreePartition{
			Parts:  make(map[Node]int),
			Weight: root.Weight,
		}
		stack := []*node{root}

		for num := len(stack); num > 0; num = len(stack) {
			n := stack[num-1]
			stack = stack[:num-1]
			if n.Node != nil {
				// Best has been place in front
				r.Parts[n.Node] = n.Kids[0].Choice
				stack = append(stack, n.Kids[0])
			} else {
				stack = append(stack, n.Kids...)
			}
		}
		return r
	}

	var stack []*node
	troot := &node{Node: opt.Root}

	stack = append(stack, troot)

	for _, v := range a.colors[opt.Root] {
		cnode := &node{
			Parent: troot,
			Choice: v,
		}
		troot.Kids = append(troot.Kids, cnode)
		stack = append(stack, cnode)
	}

	for num := len(stack); num > 0; num = len(stack) {
		snode := stack[num-1]
		stack = stack[:num-1]

		if snode.Node == nil {
			// Color node: preorder
			onode := snode.Parent.Node
			for i, inum := 0, opt.Tree.NumOuts(onode); i < inum; i++ {
				kid := opt.Tree.Out(onode, i)
				tnode := &node{
					Parent: snode,
					Node:   kid,
				}
				snode.Kids = append(snode.Kids, tnode)

				stack = append(stack, tnode)
				for _, v := range a.colors[kid] {
					cnode := &node{
						Parent: tnode,
						Choice: v,
					}
					tnode.Kids = append(tnode.Kids, cnode)
					stack = append(stack, cnode)
				}
			}
		} else {
			// Tree node: postorder
			idx, w := getBestColor(snode)
			// Place best in front
			snode.Kids[idx], snode.Kids[0] = snode.Kids[0], snode.Kids[idx]
			snode.Choice = snode.Kids[0].Choice
			snode.Weight = w
		}
	}

	return extractSolution(troot)
}
