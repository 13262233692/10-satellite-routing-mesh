package routing

import (
	"container/heap"
	"math"
	"sync"

	"github.com/aerospace/leo-routing-mesh/pkg/model"
)

const (
	inf = math.MaxFloat64
)

type Router struct {
	distPool  sync.Pool
	prevPool  sync.Pool
	heapPool  sync.Pool
	visitedPool sync.Pool
}

func NewRouter() *Router {
	return &Router{
		distPool: sync.Pool{
			New: func() interface{} {
				return make([]float64, 0, model.MaxSatellites)
			},
		},
		prevPool: sync.Pool{
			New: func() interface{} {
				return make([]int, 0, model.MaxSatellites)
			},
		},
		heapPool: sync.Pool{
			New: func() interface{} {
				h := make(PriorityQueue, 0, model.MaxSatellites)
				return &h
			},
		},
		visitedPool: sync.Pool{
			New: func() interface{} {
				return make([]bool, 0, model.MaxSatellites)
			},
		},
	}
}

type Item struct {
	index    int
	priority float64
}

type PriorityQueue []Item

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].priority < pq[j].priority
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

func (pq *PriorityQueue) Push(x interface{}) {
	item := x.(Item)
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	*pq = old[0 : n-1]
	return item
}

func (r *Router) Dijkstra(matrix *model.SparseAdjacencyMatrix, fromID, toID model.SatelliteID) model.RouteResult {
	fromIdx, ok := matrix.GetNodeIndex(fromID)
	if !ok {
		return model.RouteResult{Found: false}
	}

	toIdx, ok := matrix.GetNodeIndex(toID)
	if !ok {
		return model.RouteResult{Found: false}
	}

	n := matrix.Nodes

	dist := r.getDistSlice(n)
	prev := r.getPrevSlice(n)
	visited := r.getVisitedSlice(n)

	defer r.putDistSlice(dist)
	defer r.putPrevSlice(prev)
	defer r.putVisitedSlice(visited)

	for i := 0; i < n; i++ {
		dist[i] = inf
		prev[i] = -1
		visited[i] = false
	}
	dist[fromIdx] = 0

	pq := r.getHeap()
	*pq = (*pq)[:0]
	defer r.putHeap(pq)

	heap.Push(pq, Item{index: fromIdx, priority: 0})

	for pq.Len() > 0 {
		u := heap.Pop(pq).(Item)

		if visited[u.index] {
			continue
		}
		visited[u.index] = true

		if u.index == toIdx {
			break
		}

		if u.priority > dist[u.index] {
			continue
		}

		neighbors := matrix.Neighbors(u.index)
		for _, edge := range neighbors {
			v := int(edge.Target)
			alt := dist[u.index] + edge.Weight

			if alt < dist[v] {
				dist[v] = alt
				prev[v] = u.index
				heap.Push(pq, Item{index: v, priority: alt})
			}
		}
	}

	if dist[toIdx] == inf {
		return model.RouteResult{Found: false}
	}

	path := r.reconstructPath(prev, fromIdx, toIdx, matrix)

	return model.RouteResult{
		Path:      path,
		TotalDist: dist[toIdx] * 299792458,
		Latency:   dist[toIdx],
		Hops:      len(path) - 1,
		Found:     true,
	}
}

func (r *Router) AStar(matrix *model.SparseAdjacencyMatrix, fromID, toID model.SatelliteID, positions []model.Vec3) model.RouteResult {
	fromIdx, ok := matrix.GetNodeIndex(fromID)
	if !ok {
		return model.RouteResult{Found: false}
	}

	toIdx, ok := matrix.GetNodeIndex(toID)
	if !ok {
		return model.RouteResult{Found: false}
	}

	n := matrix.Nodes

	gScore := r.getDistSlice(n)
	fScore := r.getDistSlice(n)
	prev := r.getPrevSlice(n)
	visited := r.getVisitedSlice(n)

	defer r.putDistSlice(gScore)
	defer r.putDistSlice(fScore)
	defer r.putPrevSlice(prev)
	defer r.putVisitedSlice(visited)

	for i := 0; i < n; i++ {
		gScore[i] = inf
		fScore[i] = inf
		prev[i] = -1
		visited[i] = false
	}
	gScore[fromIdx] = 0

	goalPos := positions[toIdx]
	startPos := positions[fromIdx]
	fScore[fromIdx] = startPos.Distance(goalPos) / 299792458

	pq := r.getHeap()
	*pq = (*pq)[:0]
	defer r.putHeap(pq)

	heap.Push(pq, Item{index: fromIdx, priority: fScore[fromIdx]})

	for pq.Len() > 0 {
		current := heap.Pop(pq).(Item)

		if visited[current.index] {
			continue
		}
		visited[current.index] = true

		if current.index == toIdx {
			break
		}

		if current.priority > fScore[current.index] {
			continue
		}

		neighbors := matrix.Neighbors(current.index)
		for _, edge := range neighbors {
			v := int(edge.Target)

			if visited[v] {
				continue
			}

			tentativeG := gScore[current.index] + edge.Weight

			if tentativeG < gScore[v] {
				gScore[v] = tentativeG

				heuristic := positions[v].Distance(goalPos) / 299792458
				fScore[v] = tentativeG + heuristic
				prev[v] = current.index

				heap.Push(pq, Item{index: v, priority: fScore[v]})
			}
		}
	}

	if gScore[toIdx] == inf {
		return model.RouteResult{Found: false}
	}

	path := r.reconstructPath(prev, fromIdx, toIdx, matrix)

	return model.RouteResult{
		Path:      path,
		TotalDist: gScore[toIdx] * 299792458,
		Latency:   gScore[toIdx],
		Hops:      len(path) - 1,
		Found:     true,
	}
}

func (r *Router) reconstructPath(prev []int, fromIdx, toIdx int, matrix *model.SparseAdjacencyMatrix) []model.SatelliteID {
	pathIndices := make([]int, 0, 32)
	current := toIdx

	for current != -1 {
		pathIndices = append(pathIndices, current)
		if current == fromIdx {
			break
		}
		current = prev[current]
	}

	for i, j := 0, len(pathIndices)-1; i < j; i, j = i+1, j-1 {
		pathIndices[i], pathIndices[j] = pathIndices[j], pathIndices[i]
	}

	path := make([]model.SatelliteID, len(pathIndices))
	for i, idx := range pathIndices {
		id, _ := matrix.GetSatelliteID(idx)
		path[i] = id
	}

	return path
}

func (r *Router) getDistSlice(n int) []float64 {
	s := r.distPool.Get().([]float64)
	if cap(s) < n {
		s = make([]float64, n)
	} else {
		s = s[:n]
	}
	return s
}

func (r *Router) putDistSlice(s []float64) {
	r.distPool.Put(s)
}

func (r *Router) getPrevSlice(n int) []int {
	s := r.prevPool.Get().([]int)
	if cap(s) < n {
		s = make([]int, n)
	} else {
		s = s[:n]
	}
	return s
}

func (r *Router) putPrevSlice(s []int) {
	r.prevPool.Put(s)
}

func (r *Router) getVisitedSlice(n int) []bool {
	s := r.visitedPool.Get().([]bool)
	if cap(s) < n {
		s = make([]bool, n)
	} else {
		s = s[:n]
	}
	return s
}

func (r *Router) putVisitedSlice(s []bool) {
	r.visitedPool.Put(s)
}

func (r *Router) getHeap() *PriorityQueue {
	return r.heapPool.Get().(*PriorityQueue)
}

func (r *Router) putHeap(pq *PriorityQueue) {
	r.heapPool.Put(pq)
}
