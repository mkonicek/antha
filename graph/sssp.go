package graph

import (
	"container/heap"
)

// ShortestPathOpt is a set of options to ShortestPath
type ShortestPathOpt struct {
	Graph   Graph
	Sources []Node
	Weight  func(x, y Node) int64
}

// ShortestPath implements Dijkstra's algorithm
func ShortestPath(opt ShortestPathOpt) map[Node]int64 {
	dist := make(map[Node]int64)
	item := make(map[Node]*nodeItem)

	var pq priorityQueue

	enqueue := func(n Node, priority int64) {
		if ni, seen := item[n]; seen {
			pq.Update(ni, priority)
			return
		}
		ni := &nodeItem{
			Node:     n,
			Priority: 0,
		}
		item[n] = ni
		heap.Push(&pq, ni)
	}

	for _, n := range opt.Sources {
		dist[n] = 0
		enqueue(n, 0)
	}

	for pq.Len() > 0 {
		ni := heap.Pop(&pq).(*nodeItem)
		n := ni.Node
		delete(item, n)
		dn := dist[n]

		for i, inum := 0, opt.Graph.NumOuts(n); i < inum; i++ {
			out := opt.Graph.Out(n, i)
			w := dn + opt.Weight(n, out)
			if dout, seen := dist[out]; !seen || w < dout {
				dist[out] = w
				enqueue(out, w)
			}
		}
	}

	return dist
}
