// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v3.12.4
// source: proto/logic/logic.proto

package logic

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
	Logic_Single_FullMethodName       = "/lane.logic.logic/Single"
	Logic_Brodcast_FullMethodName     = "/lane.logic.logic/Brodcast"
	Logic_BrodcastRoom_FullMethodName = "/lane.logic.logic/BrodcastRoom"
	Logic_SetOnline_FullMethodName    = "/lane.logic.logic/SetOnline"
	Logic_SetOffline_FullMethodName   = "/lane.logic.logic/SetOffline"
	Logic_Room_FullMethodName         = "/lane.logic.logic/Room"
)

// LogicClient is the client API for Logic service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type LogicClient interface {
	Single(ctx context.Context, in *SingleReq, opts ...grpc.CallOption) (*NoResp, error)
	Brodcast(ctx context.Context, in *BrodcastReq, opts ...grpc.CallOption) (*NoResp, error)
	BrodcastRoom(ctx context.Context, in *BrodcastRoomReq, opts ...grpc.CallOption) (*NoResp, error)
	SetOnline(ctx context.Context, in *SetOnlineReq, opts ...grpc.CallOption) (*NoResp, error)
	SetOffline(ctx context.Context, in *SetOfflineReq, opts ...grpc.CallOption) (*NoResp, error)
	Room(ctx context.Context, in *RoomReq, opts ...grpc.CallOption) (*RoomResp, error)
}

type logicClient struct {
	cc grpc.ClientConnInterface
}

func NewLogicClient(cc grpc.ClientConnInterface) LogicClient {
	return &logicClient{cc}
}

func (c *logicClient) Single(ctx context.Context, in *SingleReq, opts ...grpc.CallOption) (*NoResp, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(NoResp)
	err := c.cc.Invoke(ctx, Logic_Single_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *logicClient) Brodcast(ctx context.Context, in *BrodcastReq, opts ...grpc.CallOption) (*NoResp, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(NoResp)
	err := c.cc.Invoke(ctx, Logic_Brodcast_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *logicClient) BrodcastRoom(ctx context.Context, in *BrodcastRoomReq, opts ...grpc.CallOption) (*NoResp, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(NoResp)
	err := c.cc.Invoke(ctx, Logic_BrodcastRoom_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *logicClient) SetOnline(ctx context.Context, in *SetOnlineReq, opts ...grpc.CallOption) (*NoResp, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(NoResp)
	err := c.cc.Invoke(ctx, Logic_SetOnline_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *logicClient) SetOffline(ctx context.Context, in *SetOfflineReq, opts ...grpc.CallOption) (*NoResp, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(NoResp)
	err := c.cc.Invoke(ctx, Logic_SetOffline_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *logicClient) Room(ctx context.Context, in *RoomReq, opts ...grpc.CallOption) (*RoomResp, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(RoomResp)
	err := c.cc.Invoke(ctx, Logic_Room_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// LogicServer is the server API for Logic service.
// All implementations must embed UnimplementedLogicServer
// for forward compatibility.
type LogicServer interface {
	Single(context.Context, *SingleReq) (*NoResp, error)
	Brodcast(context.Context, *BrodcastReq) (*NoResp, error)
	BrodcastRoom(context.Context, *BrodcastRoomReq) (*NoResp, error)
	SetOnline(context.Context, *SetOnlineReq) (*NoResp, error)
	SetOffline(context.Context, *SetOfflineReq) (*NoResp, error)
	Room(context.Context, *RoomReq) (*RoomResp, error)
	mustEmbedUnimplementedLogicServer()
}

// UnimplementedLogicServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedLogicServer struct{}

func (UnimplementedLogicServer) Single(context.Context, *SingleReq) (*NoResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Single not implemented")
}
func (UnimplementedLogicServer) Brodcast(context.Context, *BrodcastReq) (*NoResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Brodcast not implemented")
}
func (UnimplementedLogicServer) BrodcastRoom(context.Context, *BrodcastRoomReq) (*NoResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method BrodcastRoom not implemented")
}
func (UnimplementedLogicServer) SetOnline(context.Context, *SetOnlineReq) (*NoResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SetOnline not implemented")
}
func (UnimplementedLogicServer) SetOffline(context.Context, *SetOfflineReq) (*NoResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SetOffline not implemented")
}
func (UnimplementedLogicServer) Room(context.Context, *RoomReq) (*RoomResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Room not implemented")
}
func (UnimplementedLogicServer) mustEmbedUnimplementedLogicServer() {}
func (UnimplementedLogicServer) testEmbeddedByValue()               {}

// UnsafeLogicServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to LogicServer will
// result in compilation errors.
type UnsafeLogicServer interface {
	mustEmbedUnimplementedLogicServer()
}

func RegisterLogicServer(s grpc.ServiceRegistrar, srv LogicServer) {
	// If the following call pancis, it indicates UnimplementedLogicServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&Logic_ServiceDesc, srv)
}

func _Logic_Single_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SingleReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LogicServer).Single(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Logic_Single_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LogicServer).Single(ctx, req.(*SingleReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Logic_Brodcast_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(BrodcastReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LogicServer).Brodcast(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Logic_Brodcast_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LogicServer).Brodcast(ctx, req.(*BrodcastReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Logic_BrodcastRoom_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(BrodcastRoomReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LogicServer).BrodcastRoom(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Logic_BrodcastRoom_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LogicServer).BrodcastRoom(ctx, req.(*BrodcastRoomReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Logic_SetOnline_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SetOnlineReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LogicServer).SetOnline(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Logic_SetOnline_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LogicServer).SetOnline(ctx, req.(*SetOnlineReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Logic_SetOffline_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SetOfflineReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LogicServer).SetOffline(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Logic_SetOffline_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LogicServer).SetOffline(ctx, req.(*SetOfflineReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Logic_Room_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RoomReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LogicServer).Room(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Logic_Room_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LogicServer).Room(ctx, req.(*RoomReq))
	}
	return interceptor(ctx, in, info, handler)
}

// Logic_ServiceDesc is the grpc.ServiceDesc for Logic service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Logic_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "lane.logic.logic",
	HandlerType: (*LogicServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Single",
			Handler:    _Logic_Single_Handler,
		},
		{
			MethodName: "Brodcast",
			Handler:    _Logic_Brodcast_Handler,
		},
		{
			MethodName: "BrodcastRoom",
			Handler:    _Logic_BrodcastRoom_Handler,
		},
		{
			MethodName: "SetOnline",
			Handler:    _Logic_SetOnline_Handler,
		},
		{
			MethodName: "SetOffline",
			Handler:    _Logic_SetOffline_Handler,
		},
		{
			MethodName: "Room",
			Handler:    _Logic_Room_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/logic/logic.proto",
}
