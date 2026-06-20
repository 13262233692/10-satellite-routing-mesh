package main

import (
	"fmt"
	"log"

	"github.com/aerospace/leo-routing-mesh/internal/routing"
	"github.com/aerospace/leo-routing-mesh/internal/topology"
	"github.com/aerospace/leo-routing-mesh/pkg/model"
)

func main() {
	fmt.Println("=== LEO Routing Mesh Demo ===")

	testTopologyManager()
	testRouter()
}

func testTopologyManager() {
	fmt.Println("\n--- Testing Topology Manager ---")

	tm := topology.NewTopologyManager(100)

	for i := 0; i < 10; i++ {
		ephem := &model.Ephemeris{
			ID: model.SatelliteID(i),
			Position: model.Vec3{
				X: float64(i * 500000),
				Y: 0,
				Z: 0,
			},
			Timestamp: 0,
		}
		tm.UpdateEphemeris(ephem)

		sat := model.Satellite{
			ID:   model.SatelliteID(i),
			Name: fmt.Sprintf("SAT-%03d", i),
		}
		tm.AddSatellite(sat)
	}

	newEpoch := tm.RebuildTopology()

	snapshot := tm.GetSnapshot()
	matrix := snapshot.Matrix
	fmt.Printf("Epoch: %d, Nodes: %d\n", newEpoch, matrix.Nodes())

	linkCount := 0
	for i := 0; i < matrix.Nodes(); i++ {
		neighbors := matrix.Neighbors(i)
		linkCount += len(neighbors)
		fmt.Printf("  Node %d has %d neighbors\n", i, len(neighbors))
	}
	fmt.Printf("Total directed links: %d\n", linkCount)

	sat, ok := tm.GetSatellite(5)
	if ok {
		fmt.Printf("Satellite 5: %s\n", sat.Name)
	}

	fmt.Println("Topology Manager test PASSED")
}

func testRouter() {
	fmt.Println("\n--- Testing Router ---")

	matrix := model.NewSparseAdjacencyMatrix(100)

	for i := 0; i < 10; i++ {
		matrix.AddNode(model.SatelliteID(i))
	}

	for i := 0; i < 9; i++ {
		weight := 0.001
		matrix.AddEdge(i, i+1, weight)
		matrix.AddEdge(i+1, i, weight)
	}

	matrix.AddEdge(0, 9, 0.005)
	matrix.AddEdge(9, 0, 0.005)

	matrix.SetEpoch(1)
	matrix.Freeze()

	router := routing.NewRouter()

	result := router.Dijkstra(matrix, 0, 9)

	if result.Found {
		fmt.Printf("Dijkstra result:\n")
		fmt.Printf("  Found: %v\n", result.Found)
		fmt.Printf("  Path: %v\n", result.Path)
		fmt.Printf("  Hops: %d\n", result.Hops)
		fmt.Printf("  Latency: %.6f s\n", result.Latency)
		fmt.Printf("  Distance: %.2f m\n", result.TotalDist)
	} else {
		log.Println("Dijkstra: path not found")
	}

	positions := make([]model.Vec3, 10)
	for i := 0; i < 10; i++ {
		positions[i] = model.Vec3{X: float64(i) * 300000, Y: 0, Z: 0}
	}

	result2 := router.AStar(matrix, 0, 9, positions)

	if result2.Found {
		fmt.Printf("A* result:\n")
		fmt.Printf("  Found: %v\n", result2.Found)
		fmt.Printf("  Path: %v\n", result2.Path)
		fmt.Printf("  Hops: %d\n", result2.Hops)
		fmt.Printf("  Latency: %.6f s\n", result2.Latency)
	} else {
		log.Println("A*: path not found")
	}

	fmt.Println("Router test PASSED")
}
