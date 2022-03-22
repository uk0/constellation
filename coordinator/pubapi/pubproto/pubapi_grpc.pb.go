// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package pubproto

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// APIClient is the client API for API service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type APIClient interface {
	GetState(ctx context.Context, in *GetStateRequest, opts ...grpc.CallOption) (*GetStateResponse, error)
	ActivateAsCoordinator(ctx context.Context, in *ActivateAsCoordinatorRequest, opts ...grpc.CallOption) (API_ActivateAsCoordinatorClient, error)
	ActivateAsNode(ctx context.Context, in *ActivateAsNodeRequest, opts ...grpc.CallOption) (*ActivateAsNodeResponse, error)
	ActivateAdditionalNodes(ctx context.Context, in *ActivateAdditionalNodesRequest, opts ...grpc.CallOption) (API_ActivateAdditionalNodesClient, error)
	JoinCluster(ctx context.Context, in *JoinClusterRequest, opts ...grpc.CallOption) (*JoinClusterResponse, error)
	TriggerNodeUpdate(ctx context.Context, in *TriggerNodeUpdateRequest, opts ...grpc.CallOption) (*TriggerNodeUpdateResponse, error)
}

type aPIClient struct {
	cc grpc.ClientConnInterface
}

func NewAPIClient(cc grpc.ClientConnInterface) APIClient {
	return &aPIClient{cc}
}

func (c *aPIClient) GetState(ctx context.Context, in *GetStateRequest, opts ...grpc.CallOption) (*GetStateResponse, error) {
	out := new(GetStateResponse)
	err := c.cc.Invoke(ctx, "/pubapi.API/GetState", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *aPIClient) ActivateAsCoordinator(ctx context.Context, in *ActivateAsCoordinatorRequest, opts ...grpc.CallOption) (API_ActivateAsCoordinatorClient, error) {
	stream, err := c.cc.NewStream(ctx, &API_ServiceDesc.Streams[0], "/pubapi.API/ActivateAsCoordinator", opts...)
	if err != nil {
		return nil, err
	}
	x := &aPIActivateAsCoordinatorClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type API_ActivateAsCoordinatorClient interface {
	Recv() (*ActivateAsCoordinatorResponse, error)
	grpc.ClientStream
}

type aPIActivateAsCoordinatorClient struct {
	grpc.ClientStream
}

func (x *aPIActivateAsCoordinatorClient) Recv() (*ActivateAsCoordinatorResponse, error) {
	m := new(ActivateAsCoordinatorResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *aPIClient) ActivateAsNode(ctx context.Context, in *ActivateAsNodeRequest, opts ...grpc.CallOption) (*ActivateAsNodeResponse, error) {
	out := new(ActivateAsNodeResponse)
	err := c.cc.Invoke(ctx, "/pubapi.API/ActivateAsNode", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *aPIClient) ActivateAdditionalNodes(ctx context.Context, in *ActivateAdditionalNodesRequest, opts ...grpc.CallOption) (API_ActivateAdditionalNodesClient, error) {
	stream, err := c.cc.NewStream(ctx, &API_ServiceDesc.Streams[1], "/pubapi.API/ActivateAdditionalNodes", opts...)
	if err != nil {
		return nil, err
	}
	x := &aPIActivateAdditionalNodesClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type API_ActivateAdditionalNodesClient interface {
	Recv() (*ActivateAdditionalNodesResponse, error)
	grpc.ClientStream
}

type aPIActivateAdditionalNodesClient struct {
	grpc.ClientStream
}

func (x *aPIActivateAdditionalNodesClient) Recv() (*ActivateAdditionalNodesResponse, error) {
	m := new(ActivateAdditionalNodesResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *aPIClient) JoinCluster(ctx context.Context, in *JoinClusterRequest, opts ...grpc.CallOption) (*JoinClusterResponse, error) {
	out := new(JoinClusterResponse)
	err := c.cc.Invoke(ctx, "/pubapi.API/JoinCluster", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *aPIClient) TriggerNodeUpdate(ctx context.Context, in *TriggerNodeUpdateRequest, opts ...grpc.CallOption) (*TriggerNodeUpdateResponse, error) {
	out := new(TriggerNodeUpdateResponse)
	err := c.cc.Invoke(ctx, "/pubapi.API/TriggerNodeUpdate", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// APIServer is the server API for API service.
// All implementations must embed UnimplementedAPIServer
// for forward compatibility
type APIServer interface {
	GetState(context.Context, *GetStateRequest) (*GetStateResponse, error)
	ActivateAsCoordinator(*ActivateAsCoordinatorRequest, API_ActivateAsCoordinatorServer) error
	ActivateAsNode(context.Context, *ActivateAsNodeRequest) (*ActivateAsNodeResponse, error)
	ActivateAdditionalNodes(*ActivateAdditionalNodesRequest, API_ActivateAdditionalNodesServer) error
	JoinCluster(context.Context, *JoinClusterRequest) (*JoinClusterResponse, error)
	TriggerNodeUpdate(context.Context, *TriggerNodeUpdateRequest) (*TriggerNodeUpdateResponse, error)
	mustEmbedUnimplementedAPIServer()
}

// UnimplementedAPIServer must be embedded to have forward compatible implementations.
type UnimplementedAPIServer struct {
}

func (UnimplementedAPIServer) GetState(context.Context, *GetStateRequest) (*GetStateResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetState not implemented")
}
func (UnimplementedAPIServer) ActivateAsCoordinator(*ActivateAsCoordinatorRequest, API_ActivateAsCoordinatorServer) error {
	return status.Errorf(codes.Unimplemented, "method ActivateAsCoordinator not implemented")
}
func (UnimplementedAPIServer) ActivateAsNode(context.Context, *ActivateAsNodeRequest) (*ActivateAsNodeResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ActivateAsNode not implemented")
}
func (UnimplementedAPIServer) ActivateAdditionalNodes(*ActivateAdditionalNodesRequest, API_ActivateAdditionalNodesServer) error {
	return status.Errorf(codes.Unimplemented, "method ActivateAdditionalNodes not implemented")
}
func (UnimplementedAPIServer) JoinCluster(context.Context, *JoinClusterRequest) (*JoinClusterResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method JoinCluster not implemented")
}
func (UnimplementedAPIServer) TriggerNodeUpdate(context.Context, *TriggerNodeUpdateRequest) (*TriggerNodeUpdateResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method TriggerNodeUpdate not implemented")
}
func (UnimplementedAPIServer) mustEmbedUnimplementedAPIServer() {}

// UnsafeAPIServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to APIServer will
// result in compilation errors.
type UnsafeAPIServer interface {
	mustEmbedUnimplementedAPIServer()
}

func RegisterAPIServer(s grpc.ServiceRegistrar, srv APIServer) {
	s.RegisterService(&API_ServiceDesc, srv)
}

func _API_GetState_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetStateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(APIServer).GetState(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pubapi.API/GetState",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(APIServer).GetState(ctx, req.(*GetStateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _API_ActivateAsCoordinator_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(ActivateAsCoordinatorRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(APIServer).ActivateAsCoordinator(m, &aPIActivateAsCoordinatorServer{stream})
}

type API_ActivateAsCoordinatorServer interface {
	Send(*ActivateAsCoordinatorResponse) error
	grpc.ServerStream
}

type aPIActivateAsCoordinatorServer struct {
	grpc.ServerStream
}

func (x *aPIActivateAsCoordinatorServer) Send(m *ActivateAsCoordinatorResponse) error {
	return x.ServerStream.SendMsg(m)
}

func _API_ActivateAsNode_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ActivateAsNodeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(APIServer).ActivateAsNode(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pubapi.API/ActivateAsNode",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(APIServer).ActivateAsNode(ctx, req.(*ActivateAsNodeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _API_ActivateAdditionalNodes_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(ActivateAdditionalNodesRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(APIServer).ActivateAdditionalNodes(m, &aPIActivateAdditionalNodesServer{stream})
}

type API_ActivateAdditionalNodesServer interface {
	Send(*ActivateAdditionalNodesResponse) error
	grpc.ServerStream
}

type aPIActivateAdditionalNodesServer struct {
	grpc.ServerStream
}

func (x *aPIActivateAdditionalNodesServer) Send(m *ActivateAdditionalNodesResponse) error {
	return x.ServerStream.SendMsg(m)
}

func _API_JoinCluster_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(JoinClusterRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(APIServer).JoinCluster(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pubapi.API/JoinCluster",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(APIServer).JoinCluster(ctx, req.(*JoinClusterRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _API_TriggerNodeUpdate_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TriggerNodeUpdateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(APIServer).TriggerNodeUpdate(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pubapi.API/TriggerNodeUpdate",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(APIServer).TriggerNodeUpdate(ctx, req.(*TriggerNodeUpdateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// API_ServiceDesc is the grpc.ServiceDesc for API service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var API_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "pubapi.API",
	HandlerType: (*APIServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetState",
			Handler:    _API_GetState_Handler,
		},
		{
			MethodName: "ActivateAsNode",
			Handler:    _API_ActivateAsNode_Handler,
		},
		{
			MethodName: "JoinCluster",
			Handler:    _API_JoinCluster_Handler,
		},
		{
			MethodName: "TriggerNodeUpdate",
			Handler:    _API_TriggerNodeUpdate_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "ActivateAsCoordinator",
			Handler:       _API_ActivateAsCoordinator_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "ActivateAdditionalNodes",
			Handler:       _API_ActivateAdditionalNodes_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "pubapi.proto",
}
