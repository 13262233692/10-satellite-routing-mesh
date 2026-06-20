package server

import (
	"context"

	routingpb "github.com/aerospace/leo-routing-mesh/api/proto"
	"github.com/aerospace/leo-routing-mesh/internal/topology"
	"github.com/aerospace/leo-routing-mesh/pkg/model"
)

type TopologyServer struct {
	routingpb.UnimplementedTopologyServiceServer

	topoMgr *topology.TopologyManager
}

func NewTopologyServer(topoMgr *topology.TopologyManager) *TopologyServer {
	return &TopologyServer{
		topoMgr: topoMgr,
	}
}

func (s *TopologyServer) GetTopologyInfo(ctx context.Context, req *routingpb.GetTopologyInfoRequest) (*routingpb.GetTopologyInfoResponse, error) {
	snapshot := s.topoMgr.GetSnapshot()
	matrix := snapshot.Matrix

	linkCount := 0
	for i := 0; i < matrix.Nodes(); i++ {
		linkCount += len(matrix.Neighbors(i))
	}

	return &routingpb.GetTopologyInfoResponse{
		SatelliteCount:     int32(matrix.Nodes()),
		LinkCount:          int32(linkCount / 2),
		LastUpdateTimestamp: s.topoMgr.LastUpdateTime(),
		TopologyVersion:    snapshot.Epoch,
	}, nil
}

func (s *TopologyServer) GetSatelliteInfo(ctx context.Context, req *routingpb.GetSatelliteInfoRequest) (*routingpb.GetSatelliteInfoResponse, error) {
	satID := model.SatelliteID(req.GetSatelliteId())

	sat, exists := s.topoMgr.GetSatellite(satID)
	if !exists {
		return &routingpb.GetSatelliteInfoResponse{Exists: false}, nil
	}

	ephem := s.topoMgr.GetEphemeris(satID)

	matrix := s.topoMgr.GetMatrix()
	neighborCount := 0
	if idx, ok := matrix.GetNodeIndex(satID); ok {
		neighborCount = len(matrix.Neighbors(idx))
	}

	resp := &routingpb.GetSatelliteInfoResponse{
		Exists:        true,
		SatelliteId:   uint32(sat.ID),
		Name:          sat.Name,
		NeighborCount: int32(neighborCount),
	}

	resp.Position = &routingpb.Vec3{
		X: ephem.Position.X,
		Y: ephem.Position.Y,
		Z: ephem.Position.Z,
	}

	resp.Velocity = &routingpb.Vec3{
		X: ephem.Velocity.X,
		Y: ephem.Velocity.Y,
		Z: ephem.Velocity.Z,
	}

	return resp, nil
}

func (s *TopologyServer) GetNeighbors(ctx context.Context, req *routingpb.GetNeighborsRequest) (*routingpb.GetNeighborsResponse, error) {
	satID := model.SatelliteID(req.GetSatelliteId())

	matrix := s.topoMgr.GetMatrix()

	idx, ok := matrix.GetNodeIndex(satID)
	if !ok {
		return &routingpb.GetNeighborsResponse{Found: false}, nil
	}

	neighbors := matrix.Neighbors(idx)
	result := make([]*routingpb.NeighborInfo, 0, len(neighbors))

	for _, n := range neighbors {
		neighborID := model.SatelliteID(n.Target)
		result = append(result, &routingpb.NeighborInfo{
			SatelliteId:    uint32(neighborID),
			DistanceMeters: n.Weight * 299792458,
			LatencySeconds: n.Weight,
		})
	}

	return &routingpb.GetNeighborsResponse{
		Found:     true,
		Neighbors: result,
	}, nil
}
