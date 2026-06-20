package main

import (
	"fmt"
	"os"

	"github.com/aerospace/leo-routing-mesh/internal/routing"
	"github.com/aerospace/leo-routing-mesh/pkg/model"
)

func main() {
	fmt.Fprintln(os.Stderr, "Starting verification test...")

	matrix := model.NewSparseAdjacencyMatrix(100)

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

	fmt.Fprintf(os.Stderr, "Matrix nodes: %d\n", matrix.Nodes)

	router := routing.NewRouter()
	result := router.Dijkstra(matrix, 0, 4)

	fmt.Fprintf(os.Stderr, "Found: %v\n", result.Found)
	fmt.Fprintf(os.Stderr, "Hops: %d\n", result.Hops)
	fmt.Fprintf(os.Stderr, "Path: %v\n", result.Path)
	fmt.Fprintf(os.Stderr, "Latency: %f\n", result.Latency)

	if result.Found && result.Hops == 4 && result.Latency == 7.0 {
		fmt.Fprintln(os.Stderr, "VERIFICATION PASSED")
		os.Exit(0)
	} else {
		fmt.Fprintln(os.Stderr, "VERIFICATION FAILED")
		os.Exit(1)
	}
}
