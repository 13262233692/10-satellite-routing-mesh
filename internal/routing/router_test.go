package routing

import (
	"sync"
	"sync/atomic"
	"testing"

	"github.com/aerospace/leo-routing-mesh/pkg/model"
)

func buildTestMatrix(n int, epoch uint64) *model.SparseAdjacencyMatrix {
	matrix := model.NewSparseAdjacencyMatrix(n + 10)
	for i := 0; i < n; i++ {
		matrix.AddNode(model.SatelliteID(i))
	}
	for i := 0; i < n; i++ {
		for j := 1; j <= 4; j++ {
			if i+j < n {
				weight := float64(j) * 0.001
				matrix.AddEdge(i, i+j, weight)
				matrix.AddEdge(i+j, i, weight)
			}
		}
	}
	matrix.SetEpoch(epoch)
	matrix.Freeze()
	return matrix
}

func TestDijkstraSimple(t *testing.T) {
	matrix := model.NewSparseAdjacencyMatrix(10)

	for i := 0; i < 5; i++ {
		matrix.AddNode(model.SatelliteID(i))
	}

	matrix.AddEdge(0, 1, 1.0)
	matrix.AddEdge(1, 0, 1.0)
	matrix.AddEdge(1, 2, 2.0)
	matrix.AddEdge(2, 1, 2.0)
	matrix.AddEdge(0, 2, 5.0)
	matrix.AddEdge(2, 0, 5.0)
	matrix.AddEdge(2, 3, 1.0)
	matrix.AddEdge(3, 2, 1.0)
	matrix.AddEdge(3, 4, 3.0)
	matrix.AddEdge(4, 3, 3.0)

	matrix.SetEpoch(1)
	matrix.Freeze()

	router := NewRouter()

	result := router.Dijkstra(matrix, 0, 4)

	if !result.Found {
		t.Fatal("Expected path to be found")
	}

	if result.Hops != 4 {
		t.Fatalf("Expected 4 hops, got %d", result.Hops)
	}

	expectedPath := []model.SatelliteID{0, 1, 2, 3, 4}
	if len(result.Path) != len(expectedPath) {
		t.Fatalf("Expected path length %d, got %d", len(expectedPath), len(result.Path))
	}

	for i, id := range expectedPath {
		if result.Path[i] != id {
			t.Fatalf("Expected path[%d] = %d, got %d", i, id, result.Path[i])
		}
	}

	expectedLatency := 1.0 + 2.0 + 1.0 + 3.0
	if result.Latency != expectedLatency {
		t.Fatalf("Expected latency %f, got %f", expectedLatency, result.Latency)
	}

	t.Logf("Path: %v, Latency: %f, Hops: %d", result.Path, result.Latency, result.Hops)
}

func TestDijkstraRejectsUnfrozen(t *testing.T) {
	matrix := model.NewSparseAdjacencyMatrix(10)
	for i := 0; i < 5; i++ {
		matrix.AddNode(model.SatelliteID(i))
	}
	matrix.AddEdge(0, 1, 1.0)
	matrix.AddEdge(1, 0, 1.0)

	router := NewRouter()

	result := router.Dijkstra(matrix, 0, 1)
	if result.Found {
		t.Fatal("Expected Dijkstra to reject unfrozen matrix")
	}
}

func TestDijkstraRejectsZeroEpoch(t *testing.T) {
	matrix := model.NewSparseAdjacencyMatrix(10)
	for i := 0; i < 5; i++ {
		matrix.AddNode(model.SatelliteID(i))
	}
	matrix.AddEdge(0, 1, 1.0)
	matrix.AddEdge(1, 0, 1.0)
	matrix.Freeze()

	router := NewRouter()

	result := router.Dijkstra(matrix, 0, 1)
	if result.Found {
		t.Fatal("Expected Dijkstra to reject zero-epoch matrix")
	}
}

func TestNoRoutingLoops(t *testing.T) {
	n := 100
	matrix := buildTestMatrix(n, 1)
	router := NewRouter()

	for from := 0; from < n; from += 3 {
		for to := 0; to < n; to += 5 {
			if from == to {
				continue
			}
			result := router.Dijkstra(matrix, model.SatelliteID(from), model.SatelliteID(to))
			if !result.Found {
				continue
			}

			seen := make(map[model.SatelliteID]struct{})
			for _, id := range result.Path {
				if _, exists := seen[id]; exists {
					t.Fatalf("Routing loop detected! Path: %v contains duplicate %d", result.Path, id)
				}
				seen[id] = struct{}{}
			}
		}
	}
}

func TestConcurrentRoutingNoLoops(t *testing.T) {
	n := 200
	matrix := buildTestMatrix(n, 42)
	router := NewRouter()

	var loopCount atomic.Int64
	var completed atomic.Int64

	var wg sync.WaitGroup
	numWorkers := 16
	queriesPerWorker := 500

	for w := 0; w < numWorkers; w++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for q := 0; q < queriesPerWorker; q++ {
				from := (workerID*queriesPerWorker + q) % n
				to := (workerID*queriesPerWorker + q*7 + 3) % n
				if from == to {
					continue
				}

				result := router.Dijkstra(matrix, model.SatelliteID(from), model.SatelliteID(to))
				if !result.Found {
					continue
				}

				seen := make(map[model.SatelliteID]struct{})
				hasLoop := false
				for _, id := range result.Path {
					if _, exists := seen[id]; exists {
						hasLoop = true
						break
					}
					seen[id] = struct{}{}
				}
				if hasLoop {
					loopCount.Add(1)
					t.Errorf("Worker %d: routing loop detected in path %v", workerID, result.Path)
				}
				completed.Add(1)
			}
		}(w)
	}

	wg.Wait()

	t.Logf("Completed %d queries, loops detected: %d", completed.Load(), loopCount.Load())
	if loopCount.Load() > 0 {
		t.Fatalf("Found %d routing loops under concurrent access", loopCount.Load())
	}
}

func TestEpochCacheInvalidation(t *testing.T) {
	cache := NewRouteCache(100)

	n := 50
	matrix1 := buildTestMatrix(n, 1)
	matrix2 := buildTestMatrix(n, 2)

	router := NewRouter()

	for i := 0; i < 10; i++ {
		result := router.Dijkstra(matrix1, model.SatelliteID(i), model.SatelliteID(i+5))
		if result.Found {
			cache.Put(model.SatelliteID(i), model.SatelliteID(i+5), 1, result)
		}
	}

	if cache.Len() == 0 {
		t.Fatal("Expected some cache entries")
	}

	cache.InvalidateEpoch(1)

	if cache.Len() != 0 {
		t.Fatalf("Expected all epoch 1 entries invalidated, got %d", cache.Len())
	}

	for i := 0; i < 10; i++ {
		result := router.Dijkstra(matrix2, model.SatelliteID(i), model.SatelliteID(i+5))
		if result.Found {
			cache.Put(model.SatelliteID(i), model.SatelliteID(i+5), 2, result)
		}
	}

	if _, ok := cache.Get(0, 5, 1); ok {
		t.Fatal("Should not find stale epoch 1 entry after invalidation")
	}

	if _, ok := cache.Get(0, 5, 2); !ok {
		t.Fatal("Should find epoch 2 entry")
	}

	hits, misses, _ := cache.Stats()
	t.Logf("Cache stats: hits=%d, misses=%d, size=%d", hits, misses, cache.Len())
}

func TestCacheEpochMismatch(t *testing.T) {
	cache := NewRouteCache(100)

	result := model.RouteResult{
		Found: true,
		Hops:  3,
		Path:  []model.SatelliteID{0, 1, 2, 3},
	}
	cache.Put(0, 3, 10, result)

	if _, ok := cache.Get(0, 3, 10); !ok {
		t.Fatal("Should find cache entry for matching epoch")
	}

	if _, ok := cache.Get(0, 3, 9); ok {
		t.Fatal("Should NOT find cache entry for older epoch")
	}

	if _, ok := cache.Get(0, 3, 11); ok {
		t.Fatal("Should NOT find cache entry for newer epoch")
	}
}

func TestAStarSimple(t *testing.T) {
	matrix := model.NewSparseAdjacencyMatrix(10)

	for i := 0; i < 5; i++ {
		matrix.AddNode(model.SatelliteID(i))
	}

	matrix.AddEdge(0, 1, 1.0)
	matrix.AddEdge(1, 0, 1.0)
	matrix.AddEdge(1, 2, 2.0)
	matrix.AddEdge(2, 1, 2.0)
	matrix.AddEdge(0, 2, 5.0)
	matrix.AddEdge(2, 0, 5.0)
	matrix.AddEdge(2, 3, 1.0)
	matrix.AddEdge(3, 2, 1.0)
	matrix.AddEdge(3, 4, 3.0)
	matrix.AddEdge(4, 3, 3.0)

	matrix.SetEpoch(1)
	matrix.Freeze()

	positions := make([]model.Vec3, 5)
	for i := 0; i < 5; i++ {
		positions[i] = model.Vec3{X: float64(i * 1000000), Y: 0, Z: 0}
	}

	router := NewRouter()

	result := router.AStar(matrix, 0, 4, positions)

	if !result.Found {
		t.Fatal("Expected path to be found")
	}

	t.Logf("A* Path: %v, Latency: %f, Hops: %d", result.Path, result.Latency, result.Hops)
}

func BenchmarkDijkstra(b *testing.B) {
	n := 1000
	matrix := buildTestMatrix(n, 1)
	router := NewRouter()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		router.Dijkstra(matrix, 0, model.SatelliteID(n-1))
	}
}

func BenchmarkAStar(b *testing.B) {
	n := 1000
	matrix := model.NewSparseAdjacencyMatrix(n + 10)

	for i := 0; i < n; i++ {
		matrix.AddNode(model.SatelliteID(i))
	}

	positions := make([]model.Vec3, n)
	for i := 0; i < n; i++ {
		positions[i] = model.Vec3{X: float64(i) * 100000, Y: 0, Z: 0}

		for j := 1; j <= 4; j++ {
			if i+j < n {
				weight := float64(j) * 100000 / 299792458
				matrix.AddEdge(i, i+j, weight)
				matrix.AddEdge(i+j, i, weight)
			}
		}
	}
	matrix.SetEpoch(1)
	matrix.Freeze()

	router := NewRouter()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		router.AStar(matrix, 0, model.SatelliteID(n-1), positions)
	}
}

func TestDijkstraExample(t *testing.T) {
	matrix := model.NewSparseAdjacencyMatrix(10)

	for i := 0; i < 5; i++ {
		matrix.AddNode(model.SatelliteID(i))
	}

	matrix.AddEdge(0, 1, 0.001)
	matrix.AddEdge(1, 0, 0.001)
	matrix.AddEdge(1, 2, 0.002)
	matrix.AddEdge(2, 1, 0.002)
	matrix.AddEdge(2, 3, 0.001)
	matrix.AddEdge(3, 2, 0.001)
	matrix.AddEdge(3, 4, 0.003)
	matrix.AddEdge(4, 3, 0.003)

	matrix.SetEpoch(1)
	matrix.Freeze()

	router := NewRouter()
	result := router.Dijkstra(matrix, 0, 4)

	t.Logf("Found: %v", result.Found)
	t.Logf("Hops: %d", result.Hops)
	t.Logf("Path length: %d", len(result.Path))
}
