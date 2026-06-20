package server

import (
	"context"
	"time"

	routingpb "github.com/aerospace/leo-routing-mesh/api/proto"
	"github.com/aerospace/leo-routing-mesh/internal/topology"
	"github.com/aerospace/leo-routing-mesh/pkg/model"
)

type EphemerisServer struct {
	routingpb.UnimplementedEphemerisServiceServer

	topoMgr *topology.TopologyManager
}

func NewEphemerisServer(topoMgr *topology.TopologyManager) *EphemerisServer {
	return &EphemerisServer{
		topoMgr: topoMgr,
	}
}

func (s *EphemerisServer) UpdateEphemeris(ctx context.Context, req *routingpb.UpdateEphemerisRequest) (*routingpb.UpdateEphemerisResponse, error) {
	ephemProto := req.GetEphemeris()
	if ephemProto == nil {
		return &routingpb.UpdateEphemerisResponse{Success: false}, nil
	}

	ephem := &model.Ephemeris{
		ID:        model.SatelliteID(ephemProto.GetSatelliteId()),
		Timestamp: ephemProto.GetTimestamp(),
	}

	if pos := ephemProto.GetPosition(); pos != nil {
		ephem.Position = model.Vec3{
			X: pos.GetX(),
			Y: pos.GetY(),
			Z: pos.GetZ(),
		}
	}

	if vel := ephemProto.GetVelocity(); vel != nil {
		ephem.Velocity = model.Vec3{
			X: vel.GetX(),
			Y: vel.GetY(),
			Z: vel.GetZ(),
		}
	}

	s.topoMgr.UpdateEphemeris(ephem)

	return &routingpb.UpdateEphemerisResponse{Success: true}, nil
}

func (s *EphemerisServer) BatchUpdateEphemerides(ctx context.Context, req *routingpb.BatchUpdateEphemeridesRequest) (*routingpb.BatchUpdateEphemeridesResponse, error) {
	ephemsProto := req.GetEphemerides()
	if len(ephemsProto) == 0 {
		return &routingpb.BatchUpdateEphemeridesResponse{Success: true, UpdatedCount: 0}, nil
	}

	ephems := make([]model.Ephemeris, 0, len(ephemsProto))
	for _, ep := range ephemsProto {
		ephem := model.Ephemeris{
			ID:        model.SatelliteID(ep.GetSatelliteId()),
			Timestamp: ep.GetTimestamp(),
		}

		if pos := ep.GetPosition(); pos != nil {
			ephem.Position = model.Vec3{
				X: pos.GetX(),
				Y: pos.GetY(),
				Z: pos.GetZ(),
			}
		}

		if vel := ep.GetVelocity(); vel != nil {
			ephem.Velocity = model.Vec3{
				X: vel.GetX(),
				Y: vel.GetY(),
				Z: vel.GetZ(),
			}
		}

		ephems = append(ephems, ephem)
	}

	s.topoMgr.BatchUpdateEphemerides(ephems)

	return &routingpb.BatchUpdateEphemeridesResponse{
		Success:      true,
		UpdatedCount: int32(len(ephems)),
	}, nil
}

func (s *EphemerisServer) RebuildTopology(ctx context.Context, req *routingpb.RebuildTopologyRequest) (*routingpb.RebuildTopologyResponse, error) {
	start := time.Now()

	s.topoMgr.RebuildTopology()

	elapsed := time.Since(start)

	matrix := s.topoMgr.GetMatrix()
	linkCount := 0
	for i := 0; i < matrix.Nodes; i++ {
		linkCount += len(matrix.Neighbors(i))
	}

	return &routingpb.RebuildTopologyResponse{
		Success:        true,
		RebuildTimeNs:  elapsed.Nanoseconds(),
		SatelliteCount: int32(matrix.Nodes),
		LinkCount:      int32(linkCount / 2),
	}, nil
}
