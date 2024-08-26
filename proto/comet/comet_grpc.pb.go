// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v3.12.4
// source: proto/comet/comet.proto

package comet

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	Comet_Single_FullMethodName   = "/lane.comet.comet/Single"
	Comet_Brodcast_FullMethodName = "/lane.comet.comet/Brodcast"
	Comet_Room_FullMethodName     = "/lane.comet.comet/Room"
)

// CometClient is the client API for Comet service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type CometClient interface {
	Single(ctx context.Context, in *SingleReq, opts ...grpc.CallOption) (*NoResp, error)
	Brodcast(ctx context.Context, in *BrodcastReq, opts ...grpc.CallOption) (*NoResp, error)
	Room(ctx context.Context, in *RoomReq, opts ...grpc.CallOption) (*NoResp, error)
}

type cometClient struct {
	cc grpc.ClientConnInterface
}

func NewCometClient(cc grpc.ClientConnInterface) CometClient {
	return &cometClient{cc}
}

func (c *cometClient) Single(ctx context.Context, in *SingleReq, opts ...grpc.CallOption) (*NoResp, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(NoResp)
	err := c.cc.Invoke(ctx, Comet_Single_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *cometClient) Brodcast(ctx context.Context, in *BrodcastReq, opts ...grpc.CallOption) (*NoResp, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(NoResp)
	err := c.cc.Invoke(ctx, Comet_Brodcast_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *cometClient) Room(ctx context.Context, in *RoomReq, opts ...grpc.CallOption) (*NoResp, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(NoResp)
	err := c.cc.Invoke(ctx, Comet_Room_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// CometServer is the server API for Comet service.
// All implementations should embed UnimplementedCometServer
// for forward compatibility.
type CometServer interface {
	Single(context.Context, *SingleReq) (*NoResp, error)
	Brodcast(context.Context, *BrodcastReq) (*NoResp, error)
	Room(context.Context, *RoomReq) (*NoResp, error)
}

// UnimplementedCometServer should be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedCometServer struct{}

func (UnimplementedCometServer) Single(context.Context, *SingleReq) (*NoResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Single not implemented")
}
func (UnimplementedCometServer) Brodcast(context.Context, *BrodcastReq) (*NoResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Brodcast not implemented")
}
func (UnimplementedCometServer) Room(context.Context, *RoomReq) (*NoResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Room not implemented")
}
func (UnimplementedCometServer) testEmbeddedByValue() {}

// UnsafeCometServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to CometServer will
// result in compilation errors.
type UnsafeCometServer interface {
	mustEmbedUnimplementedCometServer()
}

func RegisterCometServer(s grpc.ServiceRegistrar, srv CometServer) {
	// If the following call pancis, it indicates UnimplementedCometServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&Comet_ServiceDesc, srv)
}

func _Comet_Single_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SingleReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CometServer).Single(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Comet_Single_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CometServer).Single(ctx, req.(*SingleReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Comet_Brodcast_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(BrodcastReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CometServer).Brodcast(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Comet_Brodcast_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CometServer).Brodcast(ctx, req.(*BrodcastReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Comet_Room_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RoomReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CometServer).Room(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Comet_Room_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CometServer).Room(ctx, req.(*RoomReq))
	}
	return interceptor(ctx, in, info, handler)
}

// Comet_ServiceDesc is the grpc.ServiceDesc for Comet service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Comet_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "lane.comet.comet",
	HandlerType: (*CometServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Single",
			Handler:    _Comet_Single_Handler,
		},
		{
			MethodName: "Brodcast",
			Handler:    _Comet_Brodcast_Handler,
		},
		{
			MethodName: "Room",
			Handler:    _Comet_Room_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/comet/comet.proto",
}
