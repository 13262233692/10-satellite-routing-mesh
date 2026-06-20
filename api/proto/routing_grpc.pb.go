package routingpb

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

const RoutingService_FindRoute_FullMethodName = "/leo_routing.RoutingService/FindRoute"
const RoutingService_BatchFindRoute_FullMethodName = "/leo_routing.RoutingService/BatchFindRoute"

type RoutingServiceClient interface {
	FindRoute(ctx context.Context, in *FindRouteRequest, opts ...grpc.CallOption) (*FindRouteResponse, error)
	BatchFindRoute(ctx context.Context, in *BatchFindRouteRequest, opts ...grpc.CallOption) (*BatchFindRouteResponse, error)
}

type routingServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewRoutingServiceClient(cc grpc.ClientConnInterface) RoutingServiceClient {
	return &routingServiceClient{cc}
}

func (c *routingServiceClient) FindRoute(ctx context.Context, in *FindRouteRequest, opts ...grpc.CallOption) (*FindRouteResponse, error) {
	out := new(FindRouteResponse)
	err := c.cc.Invoke(ctx, RoutingService_FindRoute_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *routingServiceClient) BatchFindRoute(ctx context.Context, in *BatchFindRouteRequest, opts ...grpc.CallOption) (*BatchFindRouteResponse, error) {
	out := new(BatchFindRouteResponse)
	err := c.cc.Invoke(ctx, RoutingService_BatchFindRoute_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

type RoutingServiceServer interface {
	FindRoute(context.Context, *FindRouteRequest) (*FindRouteResponse, error)
	BatchFindRoute(context.Context, *BatchFindRouteRequest) (*BatchFindRouteResponse, error)
	mustEmbedUnimplementedRoutingServiceServer()
}

type UnimplementedRoutingServiceServer struct {
}

func (UnimplementedRoutingServiceServer) FindRoute(context.Context, *FindRouteRequest) (*FindRouteResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method FindRoute not implemented")
}
func (UnimplementedRoutingServiceServer) BatchFindRoute(context.Context, *BatchFindRouteRequest) (*BatchFindRouteResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method BatchFindRoute not implemented")
}
func (UnimplementedRoutingServiceServer) mustEmbedUnimplementedRoutingServiceServer() {}

type UnsafeRoutingServiceServer interface {
	mustEmbedUnimplementedRoutingServiceServer()
}

func RegisterRoutingServiceServer(s grpc.ServiceRegistrar, srv RoutingServiceServer) {
	s.RegisterService(&RoutingService_ServiceDesc, srv)
}

func _RoutingService_FindRoute_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(FindRouteRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RoutingServiceServer).FindRoute(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RoutingService_FindRoute_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RoutingServiceServer).FindRoute(ctx, req.(*FindRouteRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RoutingService_BatchFindRoute_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(BatchFindRouteRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RoutingServiceServer).BatchFindRoute(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RoutingService_BatchFindRoute_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RoutingServiceServer).BatchFindRoute(ctx, req.(*BatchFindRouteRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var RoutingService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "leo_routing.RoutingService",
	HandlerType: (*RoutingServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "FindRoute",
			Handler:    _RoutingService_FindRoute_Handler,
		},
		{
			MethodName: "BatchFindRoute",
			Handler:    _RoutingService_BatchFindRoute_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "routing.proto",
}

const EphemerisService_UpdateEphemeris_FullMethodName = "/leo_routing.EphemerisService/UpdateEphemeris"
const EphemerisService_BatchUpdateEphemerides_FullMethodName = "/leo_routing.EphemerisService/BatchUpdateEphemerides"
const EphemerisService_RebuildTopology_FullMethodName = "/leo_routing.EphemerisService/RebuildTopology"

type EphemerisServiceClient interface {
	UpdateEphemeris(ctx context.Context, in *UpdateEphemerisRequest, opts ...grpc.CallOption) (*UpdateEphemerisResponse, error)
	BatchUpdateEphemerides(ctx context.Context, in *BatchUpdateEphemeridesRequest, opts ...grpc.CallOption) (*BatchUpdateEphemeridesResponse, error)
	RebuildTopology(ctx context.Context, in *RebuildTopologyRequest, opts ...grpc.CallOption) (*RebuildTopologyResponse, error)
}

type ephemerisServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewEphemerisServiceClient(cc grpc.ClientConnInterface) EphemerisServiceClient {
	return &ephemerisServiceClient{cc}
}

func (c *ephemerisServiceClient) UpdateEphemeris(ctx context.Context, in *UpdateEphemerisRequest, opts ...grpc.CallOption) (*UpdateEphemerisResponse, error) {
	out := new(UpdateEphemerisResponse)
	err := c.cc.Invoke(ctx, EphemerisService_UpdateEphemeris_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *ephemerisServiceClient) BatchUpdateEphemerides(ctx context.Context, in *BatchUpdateEphemeridesRequest, opts ...grpc.CallOption) (*BatchUpdateEphemeridesResponse, error) {
	out := new(BatchUpdateEphemeridesResponse)
	err := c.cc.Invoke(ctx, EphemerisService_BatchUpdateEphemerides_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *ephemerisServiceClient) RebuildTopology(ctx context.Context, in *RebuildTopologyRequest, opts ...grpc.CallOption) (*RebuildTopologyResponse, error) {
	out := new(RebuildTopologyResponse)
	err := c.cc.Invoke(ctx, EphemerisService_RebuildTopology_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

type EphemerisServiceServer interface {
	UpdateEphemeris(context.Context, *UpdateEphemerisRequest) (*UpdateEphemerisResponse, error)
	BatchUpdateEphemerides(context.Context, *BatchUpdateEphemeridesRequest) (*BatchUpdateEphemeridesResponse, error)
	RebuildTopology(context.Context, *RebuildTopologyRequest) (*RebuildTopologyResponse, error)
	mustEmbedUnimplementedEphemerisServiceServer()
}

type UnimplementedEphemerisServiceServer struct {
}

func (UnimplementedEphemerisServiceServer) UpdateEphemeris(context.Context, *UpdateEphemerisRequest) (*UpdateEphemerisResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateEphemeris not implemented")
}
func (UnimplementedEphemerisServiceServer) BatchUpdateEphemerides(context.Context, *BatchUpdateEphemeridesRequest) (*BatchUpdateEphemeridesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method BatchUpdateEphemerides not implemented")
}
func (UnimplementedEphemerisServiceServer) RebuildTopology(context.Context, *RebuildTopologyRequest) (*RebuildTopologyResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RebuildTopology not implemented")
}
func (UnimplementedEphemerisServiceServer) mustEmbedUnimplementedEphemerisServiceServer() {}

type UnsafeEphemerisServiceServer interface {
	mustEmbedUnimplementedEphemerisServiceServer()
}

func RegisterEphemerisServiceServer(s grpc.ServiceRegistrar, srv EphemerisServiceServer) {
	s.RegisterService(&EphemerisService_ServiceDesc, srv)
}

func _EphemerisService_UpdateEphemeris_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateEphemerisRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EphemerisServiceServer).UpdateEphemeris(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: EphemerisService_UpdateEphemeris_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EphemerisServiceServer).UpdateEphemeris(ctx, req.(*UpdateEphemerisRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _EphemerisService_BatchUpdateEphemerides_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(BatchUpdateEphemeridesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EphemerisServiceServer).BatchUpdateEphemerides(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: EphemerisService_BatchUpdateEphemerides_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EphemerisServiceServer).BatchUpdateEphemerides(ctx, req.(*BatchUpdateEphemeridesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _EphemerisService_RebuildTopology_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RebuildTopologyRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EphemerisServiceServer).RebuildTopology(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: EphemerisService_RebuildTopology_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EphemerisServiceServer).RebuildTopology(ctx, req.(*RebuildTopologyRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var EphemerisService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "leo_routing.EphemerisService",
	HandlerType: (*EphemerisServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "UpdateEphemeris",
			Handler:    _EphemerisService_UpdateEphemeris_Handler,
		},
		{
			MethodName: "BatchUpdateEphemerides",
			Handler:    _EphemerisService_BatchUpdateEphemerides_Handler,
		},
		{
			MethodName: "RebuildTopology",
			Handler:    _EphemerisService_RebuildTopology_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "routing.proto",
}

const TopologyService_GetTopologyInfo_FullMethodName = "/leo_routing.TopologyService/GetTopologyInfo"
const TopologyService_GetSatelliteInfo_FullMethodName = "/leo_routing.TopologyService/GetSatelliteInfo"
const TopologyService_GetNeighbors_FullMethodName = "/leo_routing.TopologyService/GetNeighbors"

type TopologyServiceClient interface {
	GetTopologyInfo(ctx context.Context, in *GetTopologyInfoRequest, opts ...grpc.CallOption) (*GetTopologyInfoResponse, error)
	GetSatelliteInfo(ctx context.Context, in *GetSatelliteInfoRequest, opts ...grpc.CallOption) (*GetSatelliteInfoResponse, error)
	GetNeighbors(ctx context.Context, in *GetNeighborsRequest, opts ...grpc.CallOption) (*GetNeighborsResponse, error)
}

type topologyServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewTopologyServiceClient(cc grpc.ClientConnInterface) TopologyServiceClient {
	return &topologyServiceClient{cc}
}

func (c *topologyServiceClient) GetTopologyInfo(ctx context.Context, in *GetTopologyInfoRequest, opts ...grpc.CallOption) (*GetTopologyInfoResponse, error) {
	out := new(GetTopologyInfoResponse)
	err := c.cc.Invoke(ctx, TopologyService_GetTopologyInfo_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *topologyServiceClient) GetSatelliteInfo(ctx context.Context, in *GetSatelliteInfoRequest, opts ...grpc.CallOption) (*GetSatelliteInfoResponse, error) {
	out := new(GetSatelliteInfoResponse)
	err := c.cc.Invoke(ctx, TopologyService_GetSatelliteInfo_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *topologyServiceClient) GetNeighbors(ctx context.Context, in *GetNeighborsRequest, opts ...grpc.CallOption) (*GetNeighborsResponse, error) {
	out := new(GetNeighborsResponse)
	err := c.cc.Invoke(ctx, TopologyService_GetNeighbors_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

type TopologyServiceServer interface {
	GetTopologyInfo(context.Context, *GetTopologyInfoRequest) (*GetTopologyInfoResponse, error)
	GetSatelliteInfo(context.Context, *GetSatelliteInfoRequest) (*GetSatelliteInfoResponse, error)
	GetNeighbors(context.Context, *GetNeighborsRequest) (*GetNeighborsResponse, error)
	mustEmbedUnimplementedTopologyServiceServer()
}

type UnimplementedTopologyServiceServer struct {
}

func (UnimplementedTopologyServiceServer) GetTopologyInfo(context.Context, *GetTopologyInfoRequest) (*GetTopologyInfoResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetTopologyInfo not implemented")
}
func (UnimplementedTopologyServiceServer) GetSatelliteInfo(context.Context, *GetSatelliteInfoRequest) (*GetSatelliteInfoResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetSatelliteInfo not implemented")
}
func (UnimplementedTopologyServiceServer) GetNeighbors(context.Context, *GetNeighborsRequest) (*GetNeighborsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetNeighbors not implemented")
}
func (UnimplementedTopologyServiceServer) mustEmbedUnimplementedTopologyServiceServer() {}

type UnsafeTopologyServiceServer interface {
	mustEmbedUnimplementedTopologyServiceServer()
}

func RegisterTopologyServiceServer(s grpc.ServiceRegistrar, srv TopologyServiceServer) {
	s.RegisterService(&TopologyService_ServiceDesc, srv)
}

func _TopologyService_GetTopologyInfo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetTopologyInfoRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TopologyServiceServer).GetTopologyInfo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: TopologyService_GetTopologyInfo_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TopologyServiceServer).GetTopologyInfo(ctx, req.(*GetTopologyInfoRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _TopologyService_GetSatelliteInfo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetSatelliteInfoRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TopologyServiceServer).GetSatelliteInfo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: TopologyService_GetSatelliteInfo_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TopologyServiceServer).GetSatelliteInfo(ctx, req.(*GetSatelliteInfoRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _TopologyService_GetNeighbors_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetNeighborsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TopologyServiceServer).GetNeighbors(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: TopologyService_GetNeighbors_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TopologyServiceServer).GetNeighbors(ctx, req.(*GetNeighborsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var TopologyService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "leo_routing.TopologyService",
	HandlerType: (*TopologyServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetTopologyInfo",
			Handler:    _TopologyService_GetTopologyInfo_Handler,
		},
		{
			MethodName: "GetSatelliteInfo",
			Handler:    _TopologyService_GetSatelliteInfo_Handler,
		},
		{
			MethodName: "GetNeighbors",
			Handler:    _TopologyService_GetNeighbors_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "routing.proto",
}

const SpaceWeatherService_IngestSpaceWeather_FullMethodName = "/leo_routing.SpaceWeatherService/IngestSpaceWeather"
const SpaceWeatherService_BatchIngestSpaceWeather_FullMethodName = "/leo_routing.SpaceWeatherService/BatchIngestSpaceWeather"
const SpaceWeatherService_GetLinkQualityPrediction_FullMethodName = "/leo_routing.SpaceWeatherService/GetLinkQualityPrediction"
const SpaceWeatherService_GetAffectedLinks_FullMethodName = "/leo_routing.SpaceWeatherService/GetAffectedLinks"

type SpaceWeatherServiceClient interface {
	IngestSpaceWeather(ctx context.Context, in *IngestSpaceWeatherRequest, opts ...grpc.CallOption) (*IngestSpaceWeatherResponse, error)
	BatchIngestSpaceWeather(ctx context.Context, in *BatchIngestSpaceWeatherRequest, opts ...grpc.CallOption) (*BatchIngestSpaceWeatherResponse, error)
	GetLinkQualityPrediction(ctx context.Context, in *GetLinkQualityPredictionRequest, opts ...grpc.CallOption) (*GetLinkQualityPredictionResponse, error)
	GetAffectedLinks(ctx context.Context, in *GetAffectedLinksRequest, opts ...grpc.CallOption) (*GetAffectedLinksResponse, error)
}

type spaceWeatherServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewSpaceWeatherServiceClient(cc grpc.ClientConnInterface) SpaceWeatherServiceClient {
	return &spaceWeatherServiceClient{cc}
}

func (c *spaceWeatherServiceClient) IngestSpaceWeather(ctx context.Context, in *IngestSpaceWeatherRequest, opts ...grpc.CallOption) (*IngestSpaceWeatherResponse, error) {
	out := new(IngestSpaceWeatherResponse)
	err := c.cc.Invoke(ctx, SpaceWeatherService_IngestSpaceWeather_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *spaceWeatherServiceClient) BatchIngestSpaceWeather(ctx context.Context, in *BatchIngestSpaceWeatherRequest, opts ...grpc.CallOption) (*BatchIngestSpaceWeatherResponse, error) {
	out := new(BatchIngestSpaceWeatherResponse)
	err := c.cc.Invoke(ctx, SpaceWeatherService_BatchIngestSpaceWeather_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *spaceWeatherServiceClient) GetLinkQualityPrediction(ctx context.Context, in *GetLinkQualityPredictionRequest, opts ...grpc.CallOption) (*GetLinkQualityPredictionResponse, error) {
	out := new(GetLinkQualityPredictionResponse)
	err := c.cc.Invoke(ctx, SpaceWeatherService_GetLinkQualityPrediction_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *spaceWeatherServiceClient) GetAffectedLinks(ctx context.Context, in *GetAffectedLinksRequest, opts ...grpc.CallOption) (*GetAffectedLinksResponse, error) {
	out := new(GetAffectedLinksResponse)
	err := c.cc.Invoke(ctx, SpaceWeatherService_GetAffectedLinks_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

type SpaceWeatherServiceServer interface {
	IngestSpaceWeather(context.Context, *IngestSpaceWeatherRequest) (*IngestSpaceWeatherResponse, error)
	BatchIngestSpaceWeather(context.Context, *BatchIngestSpaceWeatherRequest) (*BatchIngestSpaceWeatherResponse, error)
	GetLinkQualityPrediction(context.Context, *GetLinkQualityPredictionRequest) (*GetLinkQualityPredictionResponse, error)
	GetAffectedLinks(context.Context, *GetAffectedLinksRequest) (*GetAffectedLinksResponse, error)
	mustEmbedUnimplementedSpaceWeatherServiceServer()
}

type UnimplementedSpaceWeatherServiceServer struct {
}

func (UnimplementedSpaceWeatherServiceServer) IngestSpaceWeather(context.Context, *IngestSpaceWeatherRequest) (*IngestSpaceWeatherResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method IngestSpaceWeather not implemented")
}
func (UnimplementedSpaceWeatherServiceServer) BatchIngestSpaceWeather(context.Context, *BatchIngestSpaceWeatherRequest) (*BatchIngestSpaceWeatherResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method BatchIngestSpaceWeather not implemented")
}
func (UnimplementedSpaceWeatherServiceServer) GetLinkQualityPrediction(context.Context, *GetLinkQualityPredictionRequest) (*GetLinkQualityPredictionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetLinkQualityPrediction not implemented")
}
func (UnimplementedSpaceWeatherServiceServer) GetAffectedLinks(context.Context, *GetAffectedLinksRequest) (*GetAffectedLinksResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAffectedLinks not implemented")
}
func (UnimplementedSpaceWeatherServiceServer) mustEmbedUnimplementedSpaceWeatherServiceServer() {}

type UnsafeSpaceWeatherServiceServer interface {
	mustEmbedUnimplementedSpaceWeatherServiceServer()
}

func RegisterSpaceWeatherServiceServer(s grpc.ServiceRegistrar, srv SpaceWeatherServiceServer) {
	s.RegisterService(&SpaceWeatherService_ServiceDesc, srv)
}

func _SpaceWeatherService_IngestSpaceWeather_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(IngestSpaceWeatherRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SpaceWeatherServiceServer).IngestSpaceWeather(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: SpaceWeatherService_IngestSpaceWeather_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SpaceWeatherServiceServer).IngestSpaceWeather(ctx, req.(*IngestSpaceWeatherRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _SpaceWeatherService_BatchIngestSpaceWeather_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(BatchIngestSpaceWeatherRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SpaceWeatherServiceServer).BatchIngestSpaceWeather(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: SpaceWeatherService_BatchIngestSpaceWeather_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SpaceWeatherServiceServer).BatchIngestSpaceWeather(ctx, req.(*BatchIngestSpaceWeatherRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _SpaceWeatherService_GetLinkQualityPrediction_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetLinkQualityPredictionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SpaceWeatherServiceServer).GetLinkQualityPrediction(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: SpaceWeatherService_GetLinkQualityPrediction_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SpaceWeatherServiceServer).GetLinkQualityPrediction(ctx, req.(*GetLinkQualityPredictionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _SpaceWeatherService_GetAffectedLinks_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetAffectedLinksRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SpaceWeatherServiceServer).GetAffectedLinks(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: SpaceWeatherService_GetAffectedLinks_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SpaceWeatherServiceServer).GetAffectedLinks(ctx, req.(*GetAffectedLinksRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var SpaceWeatherService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "leo_routing.SpaceWeatherService",
	HandlerType: (*SpaceWeatherServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "IngestSpaceWeather",
			Handler:    _SpaceWeatherService_IngestSpaceWeather_Handler,
		},
		{
			MethodName: "BatchIngestSpaceWeather",
			Handler:    _SpaceWeatherService_BatchIngestSpaceWeather_Handler,
		},
		{
			MethodName: "GetLinkQualityPrediction",
			Handler:    _SpaceWeatherService_GetLinkQualityPrediction_Handler,
		},
		{
			MethodName: "GetAffectedLinks",
			Handler:    _SpaceWeatherService_GetAffectedLinks_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "routing.proto",
}
