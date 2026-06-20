package routingpb

import (
	proto "google.golang.org/protobuf/proto"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
)

type Vec3 struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	X float64 `protobuf:"fixed64,1,opt,name=x,proto3" json:"x,omitempty"`
	Y float64 `protobuf:"fixed64,2,opt,name=y,proto3" json:"y,omitempty"`
	Z float64 `protobuf:"fixed64,3,opt,name=z,proto3" json:"z,omitempty"`
}

func (x *Vec3) Reset()         { *x = Vec3{} }
func (x *Vec3) String() string { return protoimpl.X.MessageStringOf(x) }
func (*Vec3) ProtoMessage()    {}
func (x *Vec3) ProtoReflect() protoreflect.Message {
	return nil
}
func (x *Vec3) GetX() float64 {
	if x != nil {
		return x.X
	}
	return 0
}
func (x *Vec3) GetY() float64 {
	if x != nil {
		return x.Y
	}
	return 0
}
func (x *Vec3) GetZ() float64 {
	if x != nil {
		return x.Z
	}
	return 0
}

type Ephemeris struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	SatelliteId uint32 `protobuf:"varint,1,opt,name=satellite_id,json=satelliteId,proto3" json:"satellite_id,omitempty"`
	Position    *Vec3  `protobuf:"bytes,2,opt,name=position,proto3" json:"position,omitempty"`
	Velocity    *Vec3  `protobuf:"bytes,3,opt,name=velocity,proto3" json:"velocity,omitempty"`
	Timestamp   int64  `protobuf:"varint,4,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
}

func (x *Ephemeris) Reset()         { *x = Ephemeris{} }
func (x *Ephemeris) String() string { return protoimpl.X.MessageStringOf(x) }
func (*Ephemeris) ProtoMessage()    {}
func (x *Ephemeris) ProtoReflect() protoreflect.Message {
	return nil
}
func (x *Ephemeris) GetSatelliteId() uint32 {
	if x != nil {
		return x.SatelliteId
	}
	return 0
}
func (x *Ephemeris) GetPosition() *Vec3 {
	if x != nil {
		return x.Position
	}
	return nil
}
func (x *Ephemeris) GetVelocity() *Vec3 {
	if x != nil {
		return x.Velocity
	}
	return nil
}
func (x *Ephemeris) GetTimestamp() int64 {
	if x != nil {
		return x.Timestamp
	}
	return 0
}

type FindRouteRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	FromSatelliteId uint32 `protobuf:"varint,1,opt,name=from_satellite_id,json=fromSatelliteId,proto3" json:"from_satellite_id,omitempty"`
	ToSatelliteId   uint32 `protobuf:"varint,2,opt,name=to_satellite_id,json=toSatelliteId,proto3" json:"to_satellite_id,omitempty"`
	Timestamp       int64  `protobuf:"varint,3,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	Algorithm       string `protobuf:"bytes,4,opt,name=algorithm,proto3" json:"algorithm,omitempty"`
}

func (x *FindRouteRequest) Reset()         { *x = FindRouteRequest{} }
func (x *FindRouteRequest) String() string { return protoimpl.X.MessageStringOf(x) }
func (*FindRouteRequest) ProtoMessage()    {}
func (x *FindRouteRequest) ProtoReflect() protoreflect.Message {
	return nil
}
func (x *FindRouteRequest) GetFromSatelliteId() uint32 {
	if x != nil {
		return x.FromSatelliteId
	}
	return 0
}
func (x *FindRouteRequest) GetToSatelliteId() uint32 {
	if x != nil {
		return x.ToSatelliteId
	}
	return 0
}
func (x *FindRouteRequest) GetTimestamp() int64 {
	if x != nil {
		return x.Timestamp
	}
	return 0
}
func (x *FindRouteRequest) GetAlgorithm() string {
	if x != nil {
		return x.Algorithm
	}
	return ""
}

type RouteHop struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	FromSatelliteId uint32  `protobuf:"varint,1,opt,name=from_satellite_id,json=fromSatelliteId,proto3" json:"from_satellite_id,omitempty"`
	ToSatelliteId   uint32  `protobuf:"varint,2,opt,name=to_satellite_id,json=toSatelliteId,proto3" json:"to_satellite_id,omitempty"`
	DistanceMeters  float64 `protobuf:"fixed64,3,opt,name=distance_meters,json=distanceMeters,proto3" json:"distance_meters,omitempty"`
	LatencySeconds  float64 `protobuf:"fixed64,4,opt,name=latency_seconds,json=latencySeconds,proto3" json:"latency_seconds,omitempty"`
}

func (x *RouteHop) Reset()         { *x = RouteHop{} }
func (x *RouteHop) String() string { return protoimpl.X.MessageStringOf(x) }
func (*RouteHop) ProtoMessage()    {}
func (x *RouteHop) ProtoReflect() protoreflect.Message {
	return nil
}
func (x *RouteHop) GetFromSatelliteId() uint32 {
	if x != nil {
		return x.FromSatelliteId
	}
	return 0
}
func (x *RouteHop) GetToSatelliteId() uint32 {
	if x != nil {
		return x.ToSatelliteId
	}
	return 0
}
func (x *RouteHop) GetDistanceMeters() float64 {
	if x != nil {
		return x.DistanceMeters
	}
	return 0
}
func (x *RouteHop) GetLatencySeconds() float64 {
	if x != nil {
		return x.LatencySeconds
	}
	return 0
}

type FindRouteResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Found              bool       `protobuf:"varint,1,opt,name=found,proto3" json:"found,omitempty"`
	Path               []uint32   `protobuf:"varint,2,rep,packed,name=path,proto3" json:"path,omitempty"`
	TotalDistanceMeters float64   `protobuf:"fixed64,3,opt,name=total_distance_meters,json=totalDistanceMeters,proto3" json:"total_distance_meters,omitempty"`
	TotalLatencySeconds float64   `protobuf:"fixed64,4,opt,name=total_latency_seconds,json=totalLatencySeconds,proto3" json:"total_latency_seconds,omitempty"`
	Hops               int32      `protobuf:"varint,5,opt,name=hops,proto3" json:"hops,omitempty"`
	HopsDetail         []*RouteHop `protobuf:"bytes,6,rep,name=hops_detail,json=hopsDetail,proto3" json:"hops_detail,omitempty"`
	TopologyVersion    int64      `protobuf:"varint,7,opt,name=topology_version,json=topologyVersion,proto3" json:"topology_version,omitempty"`
}

func (x *FindRouteResponse) Reset()         { *x = FindRouteResponse{} }
func (x *FindRouteResponse) String() string { return protoimpl.X.MessageStringOf(x) }
func (*FindRouteResponse) ProtoMessage()    {}
func (x *FindRouteResponse) ProtoReflect() protoreflect.Message {
	return nil
}
func (x *FindRouteResponse) GetFound() bool {
	if x != nil {
		return x.Found
	}
	return false
}
func (x *FindRouteResponse) GetPath() []uint32 {
	if x != nil {
		return x.Path
	}
	return nil
}
func (x *FindRouteResponse) GetTotalDistanceMeters() float64 {
	if x != nil {
		return x.TotalDistanceMeters
	}
	return 0
}
func (x *FindRouteResponse) GetTotalLatencySeconds() float64 {
	if x != nil {
		return x.TotalLatencySeconds
	}
	return 0
}
func (x *FindRouteResponse) GetHops() int32 {
	if x != nil {
		return x.Hops
	}
	return 0
}
func (x *FindRouteResponse) GetHopsDetail() []*RouteHop {
	if x != nil {
		return x.HopsDetail
	}
	return nil
}
func (x *FindRouteResponse) GetTopologyVersion() int64 {
	if x != nil {
		return x.TopologyVersion
	}
	return 0
}

type BatchFindRouteRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Requests []*FindRouteRequest `protobuf:"bytes,1,rep,name=requests,proto3" json:"requests,omitempty"`
}

func (x *BatchFindRouteRequest) Reset()         { *x = BatchFindRouteRequest{} }
func (x *BatchFindRouteRequest) String() string { return protoimpl.X.MessageStringOf(x) }
func (*BatchFindRouteRequest) ProtoMessage()    {}
func (x *BatchFindRouteRequest) ProtoReflect() protoreflect.Message {
	return nil
}
func (x *BatchFindRouteRequest) GetRequests() []*FindRouteRequest {
	if x != nil {
		return x.Requests
	}
	return nil
}

type BatchFindRouteResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Responses []*FindRouteResponse `protobuf:"bytes,1,rep,name=responses,proto3" json:"responses,omitempty"`
}

func (x *BatchFindRouteResponse) Reset()         { *x = BatchFindRouteResponse{} }
func (x *BatchFindRouteResponse) String() string { return protoimpl.X.MessageStringOf(x) }
func (*BatchFindRouteResponse) ProtoMessage()    {}
func (x *BatchFindRouteResponse) ProtoReflect() protoreflect.Message {
	return nil
}
func (x *BatchFindRouteResponse) GetResponses() []*FindRouteResponse {
	if x != nil {
		return x.Responses
	}
	return nil
}

type UpdateEphemerisRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Ephemeris *Ephemeris `protobuf:"bytes,1,opt,name=ephemeris,proto3" json:"ephemeris,omitempty"`
}

func (x *UpdateEphemerisRequest) Reset()         { *x = UpdateEphemerisRequest{} }
func (x *UpdateEphemerisRequest) String() string { return protoimpl.X.MessageStringOf(x) }
func (*UpdateEphemerisRequest) ProtoMessage()    {}
func (x *UpdateEphemerisRequest) ProtoReflect() protoreflect.Message {
	return nil
}
func (x *UpdateEphemerisRequest) GetEphemeris() *Ephemeris {
	if x != nil {
		return x.Ephemeris
	}
	return nil
}

type UpdateEphemerisResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Success bool `protobuf:"varint,1,opt,name=success,proto3" json:"success,omitempty"`
}

func (x *UpdateEphemerisResponse) Reset()         { *x = UpdateEphemerisResponse{} }
func (x *UpdateEphemerisResponse) String() string { return protoimpl.X.MessageStringOf(x) }
func (*UpdateEphemerisResponse) ProtoMessage()    {}
func (x *UpdateEphemerisResponse) ProtoReflect() protoreflect.Message {
	return nil
}
func (x *UpdateEphemerisResponse) GetSuccess() bool {
	if x != nil {
		return x.Success
	}
	return false
}

type BatchUpdateEphemeridesRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Ephemerides []*Ephemeris `protobuf:"bytes,1,rep,name=ephemerides,proto3" json:"ephemerides,omitempty"`
}

func (x *BatchUpdateEphemeridesRequest) Reset()         { *x = BatchUpdateEphemeridesRequest{} }
func (x *BatchUpdateEphemeridesRequest) String() string { return protoimpl.X.MessageStringOf(x) }
func (*BatchUpdateEphemeridesRequest) ProtoMessage()    {}
func (x *BatchUpdateEphemeridesRequest) ProtoReflect() protoreflect.Message {
	return nil
}
func (x *BatchUpdateEphemeridesRequest) GetEphemerides() []*Ephemeris {
	if x != nil {
		return x.Ephemerides
	}
	return nil
}

type BatchUpdateEphemeridesResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Success      bool  `protobuf:"varint,1,opt,name=success,proto3" json:"success,omitempty"`
	UpdatedCount int32 `protobuf:"varint,2,opt,name=updated_count,json=updatedCount,proto3" json:"updated_count,omitempty"`
}

func (x *BatchUpdateEphemeridesResponse) Reset()         { *x = BatchUpdateEphemeridesResponse{} }
func (x *BatchUpdateEphemeridesResponse) String() string { return protoimpl.X.MessageStringOf(x) }
func (*BatchUpdateEphemeridesResponse) ProtoMessage()    {}
func (x *BatchUpdateEphemeridesResponse) ProtoReflect() protoreflect.Message {
	return nil
}
func (x *BatchUpdateEphemeridesResponse) GetSuccess() bool {
	if x != nil {
		return x.Success
	}
	return false
}
func (x *BatchUpdateEphemeridesResponse) GetUpdatedCount() int32 {
	if x != nil {
		return x.UpdatedCount
	}
	return 0
}

type RebuildTopologyRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *RebuildTopologyRequest) Reset()         { *x = RebuildTopologyRequest{} }
func (x *RebuildTopologyRequest) String() string { return protoimpl.X.MessageStringOf(x) }
func (*RebuildTopologyRequest) ProtoMessage()    {}
func (x *RebuildTopologyRequest) ProtoReflect() protoreflect.Message {
	return nil
}

type RebuildTopologyResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Success         bool  `protobuf:"varint,1,opt,name=success,proto3" json:"success,omitempty"`
	RebuildTimeNs   int64 `protobuf:"varint,2,opt,name=rebuild_time_ns,json=rebuildTimeNs,proto3" json:"rebuild_time_ns,omitempty"`
	SatelliteCount  int32 `protobuf:"varint,3,opt,name=satellite_count,json=satelliteCount,proto3" json:"satellite_count,omitempty"`
	LinkCount       int32 `protobuf:"varint,4,opt,name=link_count,json=linkCount,proto3" json:"link_count,omitempty"`
}

func (x *RebuildTopologyResponse) Reset()         { *x = RebuildTopologyResponse{} }
func (x *RebuildTopologyResponse) String() string { return protoimpl.X.MessageStringOf(x) }
func (*RebuildTopologyResponse) ProtoMessage()    {}
func (x *RebuildTopologyResponse) ProtoReflect() protoreflect.Message {
	return nil
}
func (x *RebuildTopologyResponse) GetSuccess() bool {
	if x != nil {
		return x.Success
	}
	return false
}
func (x *RebuildTopologyResponse) GetRebuildTimeNs() int64 {
	if x != nil {
		return x.RebuildTimeNs
	}
	return 0
}
func (x *RebuildTopologyResponse) GetSatelliteCount() int32 {
	if x != nil {
		return x.SatelliteCount
	}
	return 0
}
func (x *RebuildTopologyResponse) GetLinkCount() int32 {
	if x != nil {
		return x.LinkCount
	}
	return 0
}

type GetTopologyInfoRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *GetTopologyInfoRequest) Reset()         { *x = GetTopologyInfoRequest{} }
func (x *GetTopologyInfoRequest) String() string { return protoimpl.X.MessageStringOf(x) }
func (*GetTopologyInfoRequest) ProtoMessage()    {}
func (x *GetTopologyInfoRequest) ProtoReflect() protoreflect.Message {
	return nil
}

type GetTopologyInfoResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	SatelliteCount    int32  `protobuf:"varint,1,opt,name=satellite_count,json=satelliteCount,proto3" json:"satellite_count,omitempty"`
	LinkCount         int32  `protobuf:"varint,2,opt,name=link_count,json=linkCount,proto3" json:"link_count,omitempty"`
	LastUpdateTimestamp int64 `protobuf:"varint,3,opt,name=last_update_timestamp,json=lastUpdateTimestamp,proto3" json:"last_update_timestamp,omitempty"`
	TopologyVersion   uint64 `protobuf:"varint,4,opt,name=topology_version,json=topologyVersion,proto3" json:"topology_version,omitempty"`
}

func (x *GetTopologyInfoResponse) Reset()         { *x = GetTopologyInfoResponse{} }
func (x *GetTopologyInfoResponse) String() string { return protoimpl.X.MessageStringOf(x) }
func (*GetTopologyInfoResponse) ProtoMessage()    {}
func (x *GetTopologyInfoResponse) ProtoReflect() protoreflect.Message {
	return nil
}
func (x *GetTopologyInfoResponse) GetSatelliteCount() int32 {
	if x != nil {
		return x.SatelliteCount
	}
	return 0
}
func (x *GetTopologyInfoResponse) GetLinkCount() int32 {
	if x != nil {
		return x.LinkCount
	}
	return 0
}
func (x *GetTopologyInfoResponse) GetLastUpdateTimestamp() int64 {
	if x != nil {
		return x.LastUpdateTimestamp
	}
	return 0
}
func (x *GetTopologyInfoResponse) GetTopologyVersion() uint64 {
	if x != nil {
		return x.TopologyVersion
	}
	return 0
}

type GetSatelliteInfoRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	SatelliteId uint32 `protobuf:"varint,1,opt,name=satellite_id,json=satelliteId,proto3" json:"satellite_id,omitempty"`
}

func (x *GetSatelliteInfoRequest) Reset()         { *x = GetSatelliteInfoRequest{} }
func (x *GetSatelliteInfoRequest) String() string { return protoimpl.X.MessageStringOf(x) }
func (*GetSatelliteInfoRequest) ProtoMessage()    {}
func (x *GetSatelliteInfoRequest) ProtoReflect() protoreflect.Message {
	return nil
}
func (x *GetSatelliteInfoRequest) GetSatelliteId() uint32 {
	if x != nil {
		return x.SatelliteId
	}
	return 0
}

type GetSatelliteInfoResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Exists         bool   `protobuf:"varint,1,opt,name=exists,proto3" json:"exists,omitempty"`
	SatelliteId    uint32 `protobuf:"varint,2,opt,name=satellite_id,json=satelliteId,proto3" json:"satellite_id,omitempty"`
	Name           string `protobuf:"bytes,3,opt,name=name,proto3" json:"name,omitempty"`
	Position       *Vec3  `protobuf:"bytes,4,opt,name=position,proto3" json:"position,omitempty"`
	Velocity       *Vec3  `protobuf:"bytes,5,opt,name=velocity,proto3" json:"velocity,omitempty"`
	NeighborCount  int32  `protobuf:"varint,6,opt,name=neighbor_count,json=neighborCount,proto3" json:"neighbor_count,omitempty"`
}

func (x *GetSatelliteInfoResponse) Reset()         { *x = GetSatelliteInfoResponse{} }
func (x *GetSatelliteInfoResponse) String() string { return protoimpl.X.MessageStringOf(x) }
func (*GetSatelliteInfoResponse) ProtoMessage()    {}
func (x *GetSatelliteInfoResponse) ProtoReflect() protoreflect.Message {
	return nil
}
func (x *GetSatelliteInfoResponse) GetExists() bool {
	if x != nil {
		return x.Exists
	}
	return false
}
func (x *GetSatelliteInfoResponse) GetSatelliteId() uint32 {
	if x != nil {
		return x.SatelliteId
	}
	return 0
}
func (x *GetSatelliteInfoResponse) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}
func (x *GetSatelliteInfoResponse) GetPosition() *Vec3 {
	if x != nil {
		return x.Position
	}
	return nil
}
func (x *GetSatelliteInfoResponse) GetVelocity() *Vec3 {
	if x != nil {
		return x.Velocity
	}
	return nil
}
func (x *GetSatelliteInfoResponse) GetNeighborCount() int32 {
	if x != nil {
		return x.NeighborCount
	}
	return 0
}

type GetNeighborsRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	SatelliteId uint32 `protobuf:"varint,1,opt,name=satellite_id,json=satelliteId,proto3" json:"satellite_id,omitempty"`
}

func (x *GetNeighborsRequest) Reset()         { *x = GetNeighborsRequest{} }
func (x *GetNeighborsRequest) String() string { return protoimpl.X.MessageStringOf(x) }
func (*GetNeighborsRequest) ProtoMessage()    {}
func (x *GetNeighborsRequest) ProtoReflect() protoreflect.Message {
	return nil
}
func (x *GetNeighborsRequest) GetSatelliteId() uint32 {
	if x != nil {
		return x.SatelliteId
	}
	return 0
}

type NeighborInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	SatelliteId    uint32  `protobuf:"varint,1,opt,name=satellite_id,json=satelliteId,proto3" json:"satellite_id,omitempty"`
	DistanceMeters float64 `protobuf:"fixed64,2,opt,name=distance_meters,json=distanceMeters,proto3" json:"distance_meters,omitempty"`
	LatencySeconds float64 `protobuf:"fixed64,3,opt,name=latency_seconds,json=latencySeconds,proto3" json:"latency_seconds,omitempty"`
}

func (x *NeighborInfo) Reset()         { *x = NeighborInfo{} }
func (x *NeighborInfo) String() string { return protoimpl.X.MessageStringOf(x) }
func (*NeighborInfo) ProtoMessage()    {}
func (x *NeighborInfo) ProtoReflect() protoreflect.Message {
	return nil
}
func (x *NeighborInfo) GetSatelliteId() uint32 {
	if x != nil {
		return x.SatelliteId
	}
	return 0
}
func (x *NeighborInfo) GetDistanceMeters() float64 {
	if x != nil {
		return x.DistanceMeters
	}
	return 0
}
func (x *NeighborInfo) GetLatencySeconds() float64 {
	if x != nil {
		return x.LatencySeconds
	}
	return 0
}

type GetNeighborsResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Found     bool           `protobuf:"varint,1,opt,name=found,proto3" json:"found,omitempty"`
	Neighbors []*NeighborInfo `protobuf:"bytes,2,rep,name=neighbors,proto3" json:"neighbors,omitempty"`
}

func (x *GetNeighborsResponse) Reset()         { *x = GetNeighborsResponse{} }
func (x *GetNeighborsResponse) String() string { return protoimpl.X.MessageStringOf(x) }
func (*GetNeighborsResponse) ProtoMessage()    {}
func (x *GetNeighborsResponse) ProtoReflect() protoreflect.Message {
	return nil
}
func (x *GetNeighborsResponse) GetFound() bool {
	if x != nil {
		return x.Found
	}
	return false
}
func (x *GetNeighborsResponse) GetNeighbors() []*NeighborInfo {
	if x != nil {
		return x.Neighbors
	}
	return nil
}

var _ = proto.Marshal
