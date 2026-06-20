package routing

import (
	"testing"

	"github.com/aerospace/leo-routing-mesh/pkg/model"
)

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

	router := NewRouter()
	result := router.Dijkstra(matrix, 0, 4)

	t.Logf("Found: %v", result.Found)
	t.Logf("Hops: %d", result.Hops)
	t.Logf("Path length: %d", len(result.Path))
}
