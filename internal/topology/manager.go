package topology

import (
	"math"
	"sort"
	"sync"
	"sync/atomic"
	"time"

	"github.com/aerospace/leo-routing-mesh/pkg/model"
)

const speedOfLight = 299792458

type TopologyManager struct {
	satellites    []model.Satellite
	satellitesLen int
	ephemerides   []model.Ephemeris
	ephemValid    []bool

	activeMatrix atomic.Value
	buildMatrix  *model.SparseAdjacencyMatrix

	buildLock sync.Mutex

	maxLinkDist float64
	maxLasers   int

	lastUpdate atomic.Int64
}

func NewTopologyManager(maxSatellites int) *TopologyManager {
	tm := &TopologyManager{
		satellites:  make([]model.Satellite, maxSatellites),
		ephemerides: make([]model.Ephemeris, maxSatellites),
		ephemValid:  make([]bool, maxSatellites),
		maxLinkDist: model.MaxLinkDistanceMeters,
		maxLasers:   model.MaxLaserLinksPerSatellite,
	}

	matrixA := model.NewSparseAdjacencyMatrix(maxSatellites)
	matrixB := model.NewSparseAdjacencyMatrix(maxSatellites)

	tm.activeMatrix.Store(matrixA)
	tm.buildMatrix = matrixB

	return tm
}

func (tm *TopologyManager) GetMatrix() *model.SparseAdjacencyMatrix {
	return tm.activeMatrix.Load().(*model.SparseAdjacencyMatrix)
}

func (tm *TopologyManager) UpdateEphemeris(ephem *model.Ephemeris) {
	idx := int(ephem.ID)
	if idx >= len(tm.ephemerides) {
		return
	}
	tm.ephemerides[idx] = *ephem
	tm.ephemValid[idx] = true
}

func (tm *TopologyManager) BatchUpdateEphemerides(ephems []model.Ephemeris) {
	for _, e := range ephems {
		idx := int(e.ID)
		if idx < len(tm.ephemerides) {
			tm.ephemerides[idx] = e
			tm.ephemValid[idx] = true
		}
	}
}

func (tm *TopologyManager) RebuildTopology() {
	tm.buildLock.Lock()
	defer tm.buildLock.Unlock()

	matrix := tm.buildMatrix
	matrix.Reset()

	validCount := 0
	for i := 0; i < len(tm.ephemValid); i++ {
		if tm.ephemValid[i] {
			matrix.AddNode(model.SatelliteID(i))
			validCount++
		}
	}

	if validCount < 2 {
		tm.swapMatrix()
		return
	}

	tm.buildEdges(matrix)
	tm.swapMatrix()
	tm.lastUpdate.Store(time.Now().UnixNano())
}

func (tm *TopologyManager) swapMatrix() {
	old := tm.activeMatrix.Load().(*model.SparseAdjacencyMatrix)
	tm.activeMatrix.Store(tm.buildMatrix)
	tm.buildMatrix = old
}

type candidateLink struct {
	to     int
	weight float64
}

func (tm *TopologyManager) buildEdges(matrix *model.SparseAdjacencyMatrix) {
	n := matrix.Nodes
	if n == 0 {
		return
	}

	candidates := make([][]candidateLink, n)
	for i := 0; i < n; i++ {
		candidates[i] = make([]candidateLink, 0, n/10)
	}

	ids := make([]model.SatelliteID, n)
	for i := 0; i < n; i++ {
		id, _ := matrix.GetSatelliteID(i)
		ids[i] = id
	}

	for i := 0; i < n; i++ {
		posI := tm.ephemerides[ids[i]].Position
		for j := i + 1; j < n; j++ {
			posJ := tm.ephemerides[ids[j]].Position
			distSq := posI.DistanceSquared(posJ)

			if distSq > tm.maxLinkDist*tm.maxLinkDist {
				continue
			}

			dist := math.Sqrt(distSq)
			weight := dist / speedOfLight

			candidates[i] = append(candidates[i], candidateLink{to: j, weight: weight})
			candidates[j] = append(candidates[j], candidateLink{to: i, weight: weight})
		}
	}

	for i := 0; i < n; i++ {
		sort.Slice(candidates[i], func(a, b int) bool {
			return candidates[i][a].weight < candidates[i][b].weight
		})

		limit := tm.maxLasers
		if len(candidates[i]) < limit {
			limit = len(candidates[i])
		}

		for k := 0; k < limit; k++ {
			c := candidates[i][k]
			matrix.AddEdge(i, c.to, c.weight)
		}
	}
}

func (tm *TopologyManager) LastUpdateTime() int64 {
	return tm.lastUpdate.Load()
}

func (tm *TopologyManager) GetSatelliteCount() int {
	count := 0
	for _, v := range tm.ephemValid {
		if v {
			count++
		}
	}
	return count
}

func (tm *TopologyManager) GetSatellite(id model.SatelliteID) (model.Satellite, bool) {
	idx := int(id)
	if idx < 0 || idx >= len(tm.satellites) || !tm.ephemValid[idx] {
		return model.Satellite{}, false
	}
	return tm.satellites[idx], true
}

func (tm *TopologyManager) GetEphemeris(id model.SatelliteID) model.Ephemeris {
	idx := int(id)
	if idx < 0 || idx >= len(tm.ephemerides) {
		return model.Ephemeris{}
	}
	return tm.ephemerides[idx]
}

func (tm *TopologyManager) AddSatellite(sat model.Satellite) {
	idx := int(sat.ID)
	if idx >= len(tm.satellites) {
		return
	}
	tm.satellites[idx] = sat
}
