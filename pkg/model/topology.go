package model

import (
	"sync"
	"sync/atomic"
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
	epoch    uint64
	nodes    int
	edges    [][]AdjacencyEntry
	idToIdx  map[SatelliteID]int
	idxToID  []SatelliteID
	frozen   uint32
}

type TopologySnapshot struct {
	Matrix *SparseAdjacencyMatrix
	Epoch  uint64
}

var adjacencyPool = sync.Pool{
	New: func() interface{} {
		return &SparseAdjacencyMatrix{}
	},
}

func GetAdjacencyMatrix() *SparseAdjacencyMatrix {
	m := adjacencyPool.Get().(*SparseAdjacencyMatrix)
	atomic.StoreUint32(&m.frozen, 0)
	return m
}

func PutAdjacencyMatrix(m *SparseAdjacencyMatrix) {
	for i := range m.edges {
		m.edges[i] = m.edges[i][:0]
	}
	atomic.StoreUint32(&m.frozen, 0)
	adjacencyPool.Put(m)
}

func NewSparseAdjacencyMatrix(capacity int) *SparseAdjacencyMatrix {
	m := &SparseAdjacencyMatrix{
		nodes:   0,
		edges:   make([][]AdjacencyEntry, capacity),
		idToIdx: make(map[SatelliteID]int, capacity),
		idxToID: make([]SatelliteID, 0, capacity),
	}
	for i := range m.edges {
		m.edges[i] = make([]AdjacencyEntry, 0, MaxLaserLinksPerSatellite*2)
	}
	return m
}

func (m *SparseAdjacencyMatrix) GetEpoch() uint64 {
	return atomic.LoadUint64(&m.epoch)
}

func (m *SparseAdjacencyMatrix) SetEpoch(e uint64) {
	atomic.StoreUint64(&m.epoch, e)
}

func (m *SparseAdjacencyMatrix) Freeze() {
	atomic.StoreUint32(&m.frozen, 1)
}

func (m *SparseAdjacencyMatrix) IsFrozen() bool {
	return atomic.LoadUint32(&m.frozen) == 1
}

func (m *SparseAdjacencyMatrix) Nodes() int {
	return m.nodes
}

func (m *SparseAdjacencyMatrix) AddNode(id SatelliteID) int {
	if m.IsFrozen() {
		if idx, ok := m.idToIdx[id]; ok {
			return idx
		}
		return -1
	}
	if idx, ok := m.idToIdx[id]; ok {
		return idx
	}
	idx := m.nodes
	m.idToIdx[id] = idx
	m.idxToID = append(m.idxToID, id)
	m.nodes++
	return idx
}

func (m *SparseAdjacencyMatrix) GetNodeIndex(id SatelliteID) (int, bool) {
	idx, ok := m.idToIdx[id]
	return idx, ok
}

func (m *SparseAdjacencyMatrix) GetSatelliteID(idx int) (SatelliteID, bool) {
	if idx < 0 || idx >= m.nodes {
		return 0, false
	}
	return m.idxToID[idx], true
}

func (m *SparseAdjacencyMatrix) AddEdge(fromIdx, toIdx int, weight float64) {
	if m.IsFrozen() {
		return
	}
	if fromIdx >= m.nodes || toIdx >= m.nodes {
		return
	}
	m.edges[fromIdx] = append(m.edges[fromIdx], AdjacencyEntry{
		Target: SatelliteID(toIdx),
		Weight: weight,
	})
}

func (m *SparseAdjacencyMatrix) Neighbors(idx int) []AdjacencyEntry {
	if idx < 0 || idx >= m.nodes {
		return nil
	}
	return m.edges[idx]
}

func (m *SparseAdjacencyMatrix) Reset() {
	if m.IsFrozen() {
		return
	}
	for i := 0; i < m.nodes; i++ {
		m.edges[i] = m.edges[i][:0]
	}
	m.nodes = 0
	clear(m.idToIdx)
	m.idxToID = m.idxToID[:0]
}

func clear(m map[SatelliteID]int) {
	for k := range m {
		delete(m, k)
	}
}

func (m *SparseAdjacencyMatrix) CheckConsistency() bool {
	if len(m.idxToID) != m.nodes {
		return false
	}
	for i := 0; i < m.nodes; i++ {
		id := m.idxToID[i]
		if idx, ok := m.idToIdx[id]; !ok || idx != i {
			return false
		}
	}
	return true
}
