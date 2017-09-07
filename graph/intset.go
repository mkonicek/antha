package graph

import (
	"sort"
)

type intpair struct {
	A, B *intpair
}

// An IntSet is a set of integers
type IntSet struct {
	nodes  map[intpair]*intpair
	leaves map[int]*intpair
}

// NewIntSet returns a new IntSet
func NewIntSet() *IntSet {
	return &IntSet{
		nodes:  make(map[intpair]*intpair),
		leaves: make(map[int]*intpair),
	}
}

// Add returns a unique identifier for a set of ints
func (a *IntSet) Add(xs ...int) interface{} {
	sort.Ints(xs)

	var leaves []*intpair
	for _, x := range xs {
		l, seen := a.leaves[x]
		if !seen {
			l = &intpair{}
			a.leaves[x] = l
		}
		leaves = append(leaves, l)
	}

	for num := len(leaves); num > 1; num = len(leaves) {
		var next []*intpair
		for i := 1; i < num; i += 2 {
			pp := &intpair{A: leaves[i-1], B: leaves[i]}
			p, seen := a.nodes[*pp]
			if !seen {
				p = pp
				a.nodes[*pp] = pp
			}
			next = append(next, p)
		}

		if num&1 == 1 {
			// Odd
			next = append(next, leaves[num-1])
		}
		leaves = next
	}

	switch len(leaves) {
	default:
		fallthrough
	case 0:
		return nil
	case 1:
		return leaves[0]
	}
}
