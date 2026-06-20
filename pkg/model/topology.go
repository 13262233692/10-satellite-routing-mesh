package model

import (
	"sync"
)

type Link struct {
	From     SatelliteID
	To       SatelliteID
	Distance float64
	Latency  float64
	Weight   float64
	Active   bool
}

type AdjacencyEntry struct {
	Target SatelliteID
	Weight float64
}

type SparseAdjacencyMatrix struct {
	Nodes    int
	Edges    [][]AdjacencyEntry
	idToIdx  map[SatelliteID]int
	idxToID  []SatelliteID
	version  uint64
}

var adjacencyPool = sync.Pool{
	New: func() interface{} {
		return &SparseAdjacencyMatrix{}
	},
}

func GetAdjacencyMatrix() *SparseAdjacencyMatrix {
	m := adjacencyPool.Get().(*SparseAdjacencyMatrix)
	m.version++
	return m
}

func PutAdjacencyMatrix(m *SparseAdjacencyMatrix) {
	for i := range m.Edges {
		m.Edges[i] = m.Edges[i][:0]
	}
	adjacencyPool.Put(m)
}

func NewSparseAdjacencyMatrix(capacity int) *SparseAdjacencyMatrix {
	m := &SparseAdjacencyMatrix{
		Nodes:   0,
		Edges:   make([][]AdjacencyEntry, capacity),
		idToIdx: make(map[SatelliteID]int, capacity),
		idxToID: make([]SatelliteID, 0, capacity),
	}
	for i := range m.Edges {
		m.Edges[i] = make([]AdjacencyEntry, 0, MaxLaserLinksPerSatellite*2)
	}
	return m
}

func (m *SparseAdjacencyMatrix) AddNode(id SatelliteID) int {
	if idx, ok := m.idToIdx[id]; ok {
		return idx
	}
	idx := m.Nodes
	m.idToIdx[id] = idx
	m.idxToID = append(m.idxToID, id)
	m.Nodes++
	return idx
}

func (m *SparseAdjacencyMatrix) GetNodeIndex(id SatelliteID) (int, bool) {
	idx, ok := m.idToIdx[id]
	return idx, ok
}

func (m *SparseAdjacencyMatrix) GetSatelliteID(idx int) (SatelliteID, bool) {
	if idx < 0 || idx >= m.Nodes {
		return 0, false
	}
	return m.idxToID[idx], true
}

func (m *SparseAdjacencyMatrix) AddEdge(fromIdx, toIdx int, weight float64) {
	if fromIdx >= m.Nodes || toIdx >= m.Nodes {
		return
	}
	m.Edges[fromIdx] = append(m.Edges[fromIdx], AdjacencyEntry{
		Target: SatelliteID(toIdx),
		Weight: weight,
	})
}

func (m *SparseAdjacencyMatrix) Neighbors(idx int) []AdjacencyEntry {
	if idx < 0 || idx >= m.Nodes {
		return nil
	}
	return m.Edges[idx]
}

func (m *SparseAdjacencyMatrix) Reset() {
	for i := 0; i < m.Nodes; i++ {
		m.Edges[i] = m.Edges[i][:0]
	}
	m.Nodes = 0
	clear(m.idToIdx)
	m.idxToID = m.idxToID[:0]
}

func clear(m map[SatelliteID]int) {
	for k := range m {
		delete(m, k)
	}
}
