package topology

import (
	"log"
	"math"
	"sort"
	"sync"
	"sync/atomic"
	"time"

	"github.com/aerospace/leo-routing-mesh/pkg/model"
)

const speedOfLight = 299792458

type CacheInvalidator interface {
	InvalidateEpoch(epoch uint64)
	SetEpoch(epoch uint64)
}

type WeightPredictor interface {
	GetWeightMultiplier(satA, satB model.SatelliteID) float64
}

type TopologyManager struct {
	satellites  []model.Satellite
	ephemerides []model.Ephemeris
	ephemValid  []bool

	activeMatrix atomic.Value
	buildMatrix  *model.SparseAdjacencyMatrix

	epochCounter atomic.Uint64
	buildLock    sync.Mutex

	maxLinkDist float64
	maxLasers   int

	lastUpdate atomic.Int64

	cacheInvalidator CacheInvalidator
	cacheMu          sync.RWMutex

	weightPredictor WeightPredictor
	enablePrediction bool
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

	matrixA.SetEpoch(1)
	matrixB.SetEpoch(0)
	matrixA.Freeze()

	tm.epochCounter.Store(1)
	tm.activeMatrix.Store(matrixA)
	tm.buildMatrix = matrixB

	return tm
}

func (tm *TopologyManager) SetCacheInvalidator(ci CacheInvalidator) {
	tm.cacheMu.Lock()
	defer tm.cacheMu.Unlock()
	tm.cacheInvalidator = ci
}

func (tm *TopologyManager) SetWeightPredictor(wp WeightPredictor) {
	tm.buildLock.Lock()
	defer tm.buildLock.Unlock()
	tm.weightPredictor = wp
	tm.enablePrediction = wp != nil
}

func (tm *TopologyManager) GetSnapshot() model.TopologySnapshot {
	matrix := tm.activeMatrix.Load().(*model.SparseAdjacencyMatrix)
	return model.TopologySnapshot{
		Matrix: matrix,
		Epoch:  matrix.GetEpoch(),
	}
}

func (tm *TopologyManager) GetMatrix() *model.SparseAdjacencyMatrix {
	return tm.activeMatrix.Load().(*model.SparseAdjacencyMatrix)
}

func (tm *TopologyManager) GetEpoch() uint64 {
	matrix := tm.activeMatrix.Load().(*model.SparseAdjacencyMatrix)
	return matrix.GetEpoch()
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

func (tm *TopologyManager) RebuildTopology() uint64 {
	tm.buildLock.Lock()
	defer tm.buildLock.Unlock()

	oldActive := tm.activeMatrix.Load().(*model.SparseAdjacencyMatrix)
	oldEpoch := oldActive.GetEpoch()

	newMatrix := tm.buildMatrix
	newMatrix.Reset()

	validCount := 0
	for i := 0; i < len(tm.ephemValid); i++ {
		if tm.ephemValid[i] {
			newMatrix.AddNode(model.SatelliteID(i))
			validCount++
		}
	}

	if validCount >= 2 {
		tm.buildEdges(newMatrix)
	}

	if !newMatrix.CheckConsistency() {
		log.Printf("WARNING: Matrix consistency check failed at epoch %d, aborting swap", oldEpoch+1)
		return oldEpoch
	}

	newEpoch := tm.epochCounter.Add(1)
	newMatrix.SetEpoch(newEpoch)
	newMatrix.Freeze()

	tm.activeMatrix.Store(newMatrix)
	tm.buildMatrix = oldActive

	if tm.cacheInvalidator != nil {
		tm.cacheInvalidator.SetEpoch(newEpoch)
		tm.cacheInvalidator.InvalidateEpoch(oldEpoch)
	}

	tm.lastUpdate.Store(time.Now().UnixNano())

	return newEpoch
}

type candidateLink struct {
	to     int
	weight float64
}

func (tm *TopologyManager) buildEdges(matrix *model.SparseAdjacencyMatrix) {
	n := matrix.Nodes()
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
			edgeWeight := c.weight

			if tm.enablePrediction && tm.weightPredictor != nil {
				mult := tm.weightPredictor.GetWeightMultiplier(ids[i], ids[c.to])
				if mult > 1.0 {
					edgeWeight *= mult
				}
			}

			matrix.AddEdge(i, c.to, edgeWeight)
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
