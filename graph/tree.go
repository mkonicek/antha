package graph

import (
	"fmt"
)

// IsTree returns nil if graph is a rooted tree. If not, it returns error
// containing a counterexample.
func IsTree(g Graph, root Node) error {
	_, err := Visit(VisitOpt{
		Root:  root,
		Graph: g,
		Seen: func(n Node) error {
			return fmt.Errorf("not tree: at least two paths to %s", n)
		},
	})
	return err
}

// A TreeVisitor is a function that can be applied to nodes in a tree
type TreeVisitor func(n, parent Node, err error) error

// VisitTreeOpt are a set of options to VisitTree
type VisitTreeOpt struct {
	Tree      Graph
	Root      Node
	PreOrder  TreeVisitor // if err != TraversalDone, propagate error
	PostOrder TreeVisitor // if err != TraversalDone, propagate error
}

// VisitTree applies a tree visitor.
func VisitTree(opt VisitTreeOpt) error {
	apply := func(v TreeVisitor, n, parent Node, err error) error {
		if v == nil {
			return nil
		}
		return v(n, parent, err)
	}

	type frame struct {
		Parent, Node Node
		Post         bool
	}

	stack := []frame{{Node: opt.Root}}

	var lastError error
	for l := len(stack); l > 0; l = len(stack) {
		f := stack[l-1]
		stack = stack[:l-1]

		if f.Post {
			if err := apply(opt.PostOrder, f.Node, f.Parent, lastError); err != nil {
				lastError = err
				if err == ErrTraversalDone {
					break
				}
			}
			continue
		} else {
			if err := apply(opt.PreOrder, f.Node, f.Parent, lastError); err != nil {
				lastError = err
				if err == ErrTraversalDone {
					break
				}
				continue
			}
		}

		stack = append(stack, frame{Node: f.Node, Parent: f.Parent, Post: true})

		for i, inum := 0, opt.Tree.NumOuts(f.Node); i < inum; i++ {
			stack = append(stack, frame{Node: opt.Tree.Out(f.Node, i), Parent: f.Node})
		}
	}

	return lastError
}
