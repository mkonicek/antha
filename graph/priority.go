package graph

import (
	"container/heap"
)

type nodeItem struct {
	Node     Node
	Priority int64
	Value    interface{}
	Index    int
}

type priorityQueue []*nodeItem

func (a priorityQueue) Len() int {
	return len(a)
}

func (a priorityQueue) Swap(x, y int) {
	a[x], a[y] = a[y], a[x]
	a[x].Index = x
	a[y].Index = y
}

func (a priorityQueue) Less(x, y int) bool {
	return a[x].Priority < a[y].Priority
}

func (a *priorityQueue) Push(x interface{}) {
	n := len(*a)
	item := x.(*nodeItem)
	item.Index = n
	*a = append(*a, item)
}

func (a *priorityQueue) Pop() interface{} {
	n := len(*a)
	item := (*a)[n-1]
	item.Index = -1
	*a = (*a)[:n-1]
	return item
}

func (a *priorityQueue) Update(item *nodeItem, priority int64) {
	item.Priority = priority
	heap.Fix(a, item.Index)
}
