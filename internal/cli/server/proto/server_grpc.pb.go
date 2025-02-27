// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.12
// source: internal/cli/server/proto/server.proto

package proto

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

// EireneClient is the client API for Eirene service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type EireneClient interface {
	PeersAdd(ctx context.Context, in *PeersAddRequest, opts ...grpc.CallOption) (*PeersAddResponse, error)
	PeersRemove(ctx context.Context, in *PeersRemoveRequest, opts ...grpc.CallOption) (*PeersRemoveResponse, error)
	PeersList(ctx context.Context, in *PeersListRequest, opts ...grpc.CallOption) (*PeersListResponse, error)
	PeersStatus(ctx context.Context, in *PeersStatusRequest, opts ...grpc.CallOption) (*PeersStatusResponse, error)
	ChainSetHead(ctx context.Context, in *ChainSetHeadRequest, opts ...grpc.CallOption) (*ChainSetHeadResponse, error)
	Status(ctx context.Context, in *StatusRequest, opts ...grpc.CallOption) (*StatusResponse, error)
	ChainWatch(ctx context.Context, in *ChainWatchRequest, opts ...grpc.CallOption) (Eirene_ChainWatchClient, error)
	DebugPprof(ctx context.Context, in *DebugPprofRequest, opts ...grpc.CallOption) (Eirene_DebugPprofClient, error)
	DebugBlock(ctx context.Context, in *DebugBlockRequest, opts ...grpc.CallOption) (Eirene_DebugBlockClient, error)
}

type eireneClient struct {
	cc grpc.ClientConnInterface
}

func NewEireneClient(cc grpc.ClientConnInterface) EireneClient {
	return &eireneClient{cc}
}

func (c *eireneClient) PeersAdd(ctx context.Context, in *PeersAddRequest, opts ...grpc.CallOption) (*PeersAddResponse, error) {
	out := new(PeersAddResponse)

	err := c.cc.Invoke(ctx, "/proto.Eirene/PeersAdd", in, out, opts...)
	if err != nil {
		return nil, err
	}

	return out, nil
}

func (c *eireneClient) PeersRemove(ctx context.Context, in *PeersRemoveRequest, opts ...grpc.CallOption) (*PeersRemoveResponse, error) {
	out := new(PeersRemoveResponse)

	err := c.cc.Invoke(ctx, "/proto.Eirene/PeersRemove", in, out, opts...)
	if err != nil {
		return nil, err
	}

	return out, nil
}

func (c *eireneClient) PeersList(ctx context.Context, in *PeersListRequest, opts ...grpc.CallOption) (*PeersListResponse, error) {
	out := new(PeersListResponse)

	err := c.cc.Invoke(ctx, "/proto.Eirene/PeersList", in, out, opts...)
	if err != nil {
		return nil, err
	}

	return out, nil
}

func (c *eireneClient) PeersStatus(ctx context.Context, in *PeersStatusRequest, opts ...grpc.CallOption) (*PeersStatusResponse, error) {
	out := new(PeersStatusResponse)

	err := c.cc.Invoke(ctx, "/proto.Eirene/PeersStatus", in, out, opts...)
	if err != nil {
		return nil, err
	}

	return out, nil
}

func (c *eireneClient) ChainSetHead(ctx context.Context, in *ChainSetHeadRequest, opts ...grpc.CallOption) (*ChainSetHeadResponse, error) {
	out := new(ChainSetHeadResponse)

	err := c.cc.Invoke(ctx, "/proto.Eirene/ChainSetHead", in, out, opts...)
	if err != nil {
		return nil, err
	}

	return out, nil
}

func (c *eireneClient) Status(ctx context.Context, in *StatusRequest, opts ...grpc.CallOption) (*StatusResponse, error) {
	out := new(StatusResponse)

	err := c.cc.Invoke(ctx, "/proto.Eirene/Status", in, out, opts...)
	if err != nil {
		return nil, err
	}

	return out, nil
}

func (c *eireneClient) ChainWatch(ctx context.Context, in *ChainWatchRequest, opts ...grpc.CallOption) (Eirene_ChainWatchClient, error) {
	stream, err := c.cc.NewStream(ctx, &Eirene_ServiceDesc.Streams[0], "/proto.Eirene/ChainWatch", opts...)
	if err != nil {
		return nil, err
	}

	x := &eireneChainWatchClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}

	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}

	return x, nil
}

type Eirene_ChainWatchClient interface {
	Recv() (*ChainWatchResponse, error)
	grpc.ClientStream
}

type eireneChainWatchClient struct {
	grpc.ClientStream
}

func (x *eireneChainWatchClient) Recv() (*ChainWatchResponse, error) {
	m := new(ChainWatchResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}

	return m, nil
}

func (c *eireneClient) DebugPprof(ctx context.Context, in *DebugPprofRequest, opts ...grpc.CallOption) (Eirene_DebugPprofClient, error) {
	stream, err := c.cc.NewStream(ctx, &Eirene_ServiceDesc.Streams[1], "/proto.Eirene/DebugPprof", opts...)
	if err != nil {
		return nil, err
	}

	x := &eireneDebugPprofClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}

	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}

	return x, nil
}

type Eirene_DebugPprofClient interface {
	Recv() (*DebugFileResponse, error)
	grpc.ClientStream
}

type eireneDebugPprofClient struct {
	grpc.ClientStream
}

func (x *eireneDebugPprofClient) Recv() (*DebugFileResponse, error) {
	m := new(DebugFileResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}

	return m, nil
}

func (c *eireneClient) DebugBlock(ctx context.Context, in *DebugBlockRequest, opts ...grpc.CallOption) (Eirene_DebugBlockClient, error) {
	stream, err := c.cc.NewStream(ctx, &Eirene_ServiceDesc.Streams[2], "/proto.Eirene/DebugBlock", opts...)
	if err != nil {
		return nil, err
	}

	x := &eireneDebugBlockClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}

	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}

	return x, nil
}

type Eirene_DebugBlockClient interface {
	Recv() (*DebugFileResponse, error)
	grpc.ClientStream
}

type eireneDebugBlockClient struct {
	grpc.ClientStream
}

func (x *eireneDebugBlockClient) Recv() (*DebugFileResponse, error) {
	m := new(DebugFileResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}

	return m, nil
}

// EireneServer is the server API for Eirene service.
// All implementations must embed UnimplementedEireneServer
// for forward compatibility
type EireneServer interface {
	PeersAdd(context.Context, *PeersAddRequest) (*PeersAddResponse, error)
	PeersRemove(context.Context, *PeersRemoveRequest) (*PeersRemoveResponse, error)
	PeersList(context.Context, *PeersListRequest) (*PeersListResponse, error)
	PeersStatus(context.Context, *PeersStatusRequest) (*PeersStatusResponse, error)
	ChainSetHead(context.Context, *ChainSetHeadRequest) (*ChainSetHeadResponse, error)
	Status(context.Context, *StatusRequest) (*StatusResponse, error)
	ChainWatch(*ChainWatchRequest, Eirene_ChainWatchServer) error
	DebugPprof(*DebugPprofRequest, Eirene_DebugPprofServer) error
	DebugBlock(*DebugBlockRequest, Eirene_DebugBlockServer) error
	mustEmbedUnimplementedEireneServer()
}

// UnimplementedEireneServer must be embedded to have forward compatible implementations.
type UnimplementedEireneServer struct {
}

func (UnimplementedEireneServer) PeersAdd(context.Context, *PeersAddRequest) (*PeersAddResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PeersAdd not implemented")
}
func (UnimplementedEireneServer) PeersRemove(context.Context, *PeersRemoveRequest) (*PeersRemoveResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PeersRemove not implemented")
}
func (UnimplementedEireneServer) PeersList(context.Context, *PeersListRequest) (*PeersListResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PeersList not implemented")
}
func (UnimplementedEireneServer) PeersStatus(context.Context, *PeersStatusRequest) (*PeersStatusResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PeersStatus not implemented")
}
func (UnimplementedEireneServer) ChainSetHead(context.Context, *ChainSetHeadRequest) (*ChainSetHeadResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ChainSetHead not implemented")
}
func (UnimplementedEireneServer) Status(context.Context, *StatusRequest) (*StatusResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Status not implemented")
}
func (UnimplementedEireneServer) ChainWatch(*ChainWatchRequest, Eirene_ChainWatchServer) error {
	return status.Errorf(codes.Unimplemented, "method ChainWatch not implemented")
}
func (UnimplementedEireneServer) DebugPprof(*DebugPprofRequest, Eirene_DebugPprofServer) error {
	return status.Errorf(codes.Unimplemented, "method DebugPprof not implemented")
}
func (UnimplementedEireneServer) DebugBlock(*DebugBlockRequest, Eirene_DebugBlockServer) error {
	return status.Errorf(codes.Unimplemented, "method DebugBlock not implemented")
}
func (UnimplementedEireneServer) mustEmbedUnimplementedEireneServer() {}

// UnsafeEireneServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to EireneServer will
// result in compilation errors.
type UnsafeEireneServer interface {
	mustEmbedUnimplementedEireneServer()
}

func RegisterEireneServer(s grpc.ServiceRegistrar, srv EireneServer) {
	s.RegisterService(&Eirene_ServiceDesc, srv)
}

func _Eirene_PeersAdd_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PeersAddRequest)
	if err := dec(in); err != nil {
		return nil, err
	}

	if interceptor == nil {
		return srv.(EireneServer).PeersAdd(ctx, in)
	}

	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.Eirene/PeersAdd",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EireneServer).PeersAdd(ctx, req.(*PeersAddRequest))
	}

	return interceptor(ctx, in, info, handler)
}

func _Eirene_PeersRemove_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PeersRemoveRequest)
	if err := dec(in); err != nil {
		return nil, err
	}

	if interceptor == nil {
		return srv.(EireneServer).PeersRemove(ctx, in)
	}

	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.Eirene/PeersRemove",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EireneServer).PeersRemove(ctx, req.(*PeersRemoveRequest))
	}

	return interceptor(ctx, in, info, handler)
}

func _Eirene_PeersList_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PeersListRequest)
	if err := dec(in); err != nil {
		return nil, err
	}

	if interceptor == nil {
		return srv.(EireneServer).PeersList(ctx, in)
	}

	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.Eirene/PeersList",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EireneServer).PeersList(ctx, req.(*PeersListRequest))
	}

	return interceptor(ctx, in, info, handler)
}

func _Eirene_PeersStatus_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PeersStatusRequest)
	if err := dec(in); err != nil {
		return nil, err
	}

	if interceptor == nil {
		return srv.(EireneServer).PeersStatus(ctx, in)
	}

	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.Eirene/PeersStatus",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EireneServer).PeersStatus(ctx, req.(*PeersStatusRequest))
	}

	return interceptor(ctx, in, info, handler)
}

func _Eirene_ChainSetHead_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ChainSetHeadRequest)
	if err := dec(in); err != nil {
		return nil, err
	}

	if interceptor == nil {
		return srv.(EireneServer).ChainSetHead(ctx, in)
	}

	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.Eirene/ChainSetHead",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EireneServer).ChainSetHead(ctx, req.(*ChainSetHeadRequest))
	}

	return interceptor(ctx, in, info, handler)
}

func _Eirene_Status_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(StatusRequest)
	if err := dec(in); err != nil {
		return nil, err
	}

	if interceptor == nil {
		return srv.(EireneServer).Status(ctx, in)
	}

	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.Eirene/Status",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EireneServer).Status(ctx, req.(*StatusRequest))
	}

	return interceptor(ctx, in, info, handler)
}

func _Eirene_ChainWatch_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(ChainWatchRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}

	return srv.(EireneServer).ChainWatch(m, &eireneChainWatchServer{stream})
}

type Eirene_ChainWatchServer interface {
	Send(*ChainWatchResponse) error
	grpc.ServerStream
}

type eireneChainWatchServer struct {
	grpc.ServerStream
}

func (x *eireneChainWatchServer) Send(m *ChainWatchResponse) error {
	return x.ServerStream.SendMsg(m)
}

func _Eirene_DebugPprof_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(DebugPprofRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}

	return srv.(EireneServer).DebugPprof(m, &eireneDebugPprofServer{stream})
}

type Eirene_DebugPprofServer interface {
	Send(*DebugFileResponse) error
	grpc.ServerStream
}

type eireneDebugPprofServer struct {
	grpc.ServerStream
}

func (x *eireneDebugPprofServer) Send(m *DebugFileResponse) error {
	return x.ServerStream.SendMsg(m)
}

func _Eirene_DebugBlock_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(DebugBlockRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}

	return srv.(EireneServer).DebugBlock(m, &eireneDebugBlockServer{stream})
}

type Eirene_DebugBlockServer interface {
	Send(*DebugFileResponse) error
	grpc.ServerStream
}

type eireneDebugBlockServer struct {
	grpc.ServerStream
}

func (x *eireneDebugBlockServer) Send(m *DebugFileResponse) error {
	return x.ServerStream.SendMsg(m)
}

// Eirene_ServiceDesc is the grpc.ServiceDesc for Eirene service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Eirene_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "proto.Eirene",
	HandlerType: (*EireneServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "PeersAdd",
			Handler:    _Eirene_PeersAdd_Handler,
		},
		{
			MethodName: "PeersRemove",
			Handler:    _Eirene_PeersRemove_Handler,
		},
		{
			MethodName: "PeersList",
			Handler:    _Eirene_PeersList_Handler,
		},
		{
			MethodName: "PeersStatus",
			Handler:    _Eirene_PeersStatus_Handler,
		},
		{
			MethodName: "ChainSetHead",
			Handler:    _Eirene_ChainSetHead_Handler,
		},
		{
			MethodName: "Status",
			Handler:    _Eirene_Status_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "ChainWatch",
			Handler:       _Eirene_ChainWatch_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "DebugPprof",
			Handler:       _Eirene_DebugPprof_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "DebugBlock",
			Handler:       _Eirene_DebugBlock_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "internal/cli/server/proto/server.proto",
}
