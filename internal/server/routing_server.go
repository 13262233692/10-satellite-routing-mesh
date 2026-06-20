package server

import (
	"context"
	"strings"

	routingpb "github.com/aerospace/leo-routing-mesh/api/proto"
	"github.com/aerospace/leo-routing-mesh/internal/routing"
	"github.com/aerospace/leo-routing-mesh/internal/topology"
	"github.com/aerospace/leo-routing-mesh/pkg/model"
)

type RoutingServer struct {
	routingpb.UnimplementedRoutingServiceServer

	topoMgr *topology.TopologyManager
	router  *routing.Router
}

func NewRoutingServer(topoMgr *topology.TopologyManager, router *routing.Router) *RoutingServer {
	return &RoutingServer{
		topoMgr: topoMgr,
		router:  router,
	}
}

func (s *RoutingServer) FindRoute(ctx context.Context, req *routingpb.FindRouteRequest) (*routingpb.FindRouteResponse, error) {
	matrix := s.topoMgr.GetMatrix()

	fromID := model.SatelliteID(req.GetFromSatelliteId())
	toID := model.SatelliteID(req.GetToSatelliteId())

	algo := strings.ToLower(req.GetAlgorithm())
	if algo == "" {
		algo = "dijkstra"
	}

	var result model.RouteResult

	switch algo {
	case "astar", "a*", "a_star":
		positions := s.getPositionArray(matrix)
		result = s.router.AStar(matrix, fromID, toID, positions)
	default:
		result = s.router.Dijkstra(matrix, fromID, toID)
	}

	resp := &routingpb.FindRouteResponse{
		Found:              result.Found,
		TotalDistanceMeters: result.TotalDist,
		TotalLatencySeconds: result.Latency,
		Hops:               int32(result.Hops),
		TopologyVersion:    s.topoMgr.LastUpdateTime(),
	}

	if result.Found {
		resp.Path = make([]uint32, len(result.Path))
		for i, id := range result.Path {
			resp.Path[i] = uint32(id)
		}

		resp.HopsDetail = s.buildHopsDetail(result.Path, matrix)
	}

	return resp, nil
}

func (s *RoutingServer) BatchFindRoute(ctx context.Context, req *routingpb.BatchFindRouteRequest) (*routingpb.BatchFindRouteResponse, error) {
	requests := req.GetRequests()
	responses := make([]*routingpb.FindRouteResponse, len(requests))

	for i, r := range requests {
		resp, err := s.FindRoute(ctx, r)
		if err != nil {
			responses[i] = &routingpb.FindRouteResponse{Found: false}
			continue
		}
		responses[i] = resp
	}

	return &routingpb.BatchFindRouteResponse{
		Responses: responses,
	}, nil
}

func (s *RoutingServer) buildHopsDetail(path []model.SatelliteID, matrix *model.SparseAdjacencyMatrix) []*routingpb.RouteHop {
	if len(path) < 2 {
		return nil
	}

	hops := make([]*routingpb.RouteHop, 0, len(path)-1)

	for i := 0; i < len(path)-1; i++ {
		fromIdx, ok := matrix.GetNodeIndex(path[i])
		if !ok {
			continue
		}

		toID := path[i+1]
		var weight float64

		neighbors := matrix.Neighbors(fromIdx)
		for _, n := range neighbors {
			if model.SatelliteID(n.Target) == toID {
				weight = n.Weight
				break
			}
		}

		hops = append(hops, &routingpb.RouteHop{
			FromSatelliteId: uint32(path[i]),
			ToSatelliteId:   uint32(toID),
			DistanceMeters:  weight * 299792458,
			LatencySeconds:  weight,
		})
	}

	return hops
}

func (s *RoutingServer) getPositionArray(matrix *model.SparseAdjacencyMatrix) []model.Vec3 {
	n := matrix.Nodes
	positions := make([]model.Vec3, n)

	for i := 0; i < n; i++ {
		id, _ := matrix.GetSatelliteID(i)
		positions[i] = s.topoMgr.GetEphemeris(id).Position
	}

	return positions
}
