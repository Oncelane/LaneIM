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
	msg "laneIM/proto/msg"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	Logic_SendMsg_FullMethodName         = "/lane.logic.logic/SendMsg"
	Logic_SendMsgBatch_FullMethodName    = "/lane.logic.logic/SendMsgBatch"
	Logic_SetOnline_FullMethodName       = "/lane.logic.logic/SetOnline"
	Logic_SetOnlineBatch_FullMethodName  = "/lane.logic.logic/SetOnlineBatch"
	Logic_SetOffline_FullMethodName      = "/lane.logic.logic/SetOffline"
	Logic_SetOfflineBatch_FullMethodName = "/lane.logic.logic/SetOfflineBatch"
	Logic_NewUser_FullMethodName         = "/lane.logic.logic/NewUser"
	Logic_NewUserBatch_FullMethodName    = "/lane.logic.logic/NewUserBatch"
	Logic_DelUser_FullMethodName         = "/lane.logic.logic/DelUser"
	Logic_NewRoom_FullMethodName         = "/lane.logic.logic/NewRoom"
	Logic_JoinRoom_FullMethodName        = "/lane.logic.logic/JoinRoom"
	Logic_JoinRoomBatch_FullMethodName   = "/lane.logic.logic/JoinRoomBatch"
	Logic_QuitRoom_FullMethodName        = "/lane.logic.logic/QuitRoom"
	Logic_QueryRoom_FullMethodName       = "/lane.logic.logic/QueryRoom"
	Logic_QueryServer_FullMethodName     = "/lane.logic.logic/QueryServer"
	Logic_Auth_FullMethodName            = "/lane.logic.logic/Auth"
)

// LogicClient is the client API for Logic service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type LogicClient interface {
	SendMsg(ctx context.Context, in *msg.SendMsgReq, opts ...grpc.CallOption) (*NoResp, error)
	SendMsgBatch(ctx context.Context, in *msg.SendMsgBatchReq, opts ...grpc.CallOption) (*NoResp, error)
	SetOnline(ctx context.Context, in *SetOnlineReq, opts ...grpc.CallOption) (*NoResp, error)
	SetOnlineBatch(ctx context.Context, in *SetOnlineBatchReq, opts ...grpc.CallOption) (*NoResp, error)
	SetOffline(ctx context.Context, in *SetOfflineReq, opts ...grpc.CallOption) (*NoResp, error)
	SetOfflineBatch(ctx context.Context, in *SetOfflineBatchReq, opts ...grpc.CallOption) (*NoResp, error)
	NewUser(ctx context.Context, in *NewUserReq, opts ...grpc.CallOption) (*NewUserResp, error)
	NewUserBatch(ctx context.Context, in *NewUserBatchReq, opts ...grpc.CallOption) (*NewUserBatchResp, error)
	DelUser(ctx context.Context, in *DelUserReq, opts ...grpc.CallOption) (*NoResp, error)
	NewRoom(ctx context.Context, in *NewRoomReq, opts ...grpc.CallOption) (*NewRoomResp, error)
	JoinRoom(ctx context.Context, in *JoinRoomReq, opts ...grpc.CallOption) (*NoResp, error)
	JoinRoomBatch(ctx context.Context, in *JoinRoomBatchReq, opts ...grpc.CallOption) (*NoResp, error)
	QuitRoom(ctx context.Context, in *QuitRoomReq, opts ...grpc.CallOption) (*NoResp, error)
	QueryRoom(ctx context.Context, in *QueryRoomReq, opts ...grpc.CallOption) (*QueryRoomResp, error)
	QueryServer(ctx context.Context, in *QueryServerReq, opts ...grpc.CallOption) (*QueryServerResp, error)
	Auth(ctx context.Context, in *AuthReq, opts ...grpc.CallOption) (*AuthResp, error)
}

type logicClient struct {
	cc grpc.ClientConnInterface
}

func NewLogicClient(cc grpc.ClientConnInterface) LogicClient {
	return &logicClient{cc}
}

func (c *logicClient) SendMsg(ctx context.Context, in *msg.SendMsgReq, opts ...grpc.CallOption) (*NoResp, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(NoResp)
	err := c.cc.Invoke(ctx, Logic_SendMsg_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *logicClient) SendMsgBatch(ctx context.Context, in *msg.SendMsgBatchReq, opts ...grpc.CallOption) (*NoResp, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(NoResp)
	err := c.cc.Invoke(ctx, Logic_SendMsgBatch_FullMethodName, in, out, cOpts...)
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

func (c *logicClient) SetOnlineBatch(ctx context.Context, in *SetOnlineBatchReq, opts ...grpc.CallOption) (*NoResp, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(NoResp)
	err := c.cc.Invoke(ctx, Logic_SetOnlineBatch_FullMethodName, in, out, cOpts...)
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

func (c *logicClient) SetOfflineBatch(ctx context.Context, in *SetOfflineBatchReq, opts ...grpc.CallOption) (*NoResp, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(NoResp)
	err := c.cc.Invoke(ctx, Logic_SetOfflineBatch_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *logicClient) NewUser(ctx context.Context, in *NewUserReq, opts ...grpc.CallOption) (*NewUserResp, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(NewUserResp)
	err := c.cc.Invoke(ctx, Logic_NewUser_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *logicClient) NewUserBatch(ctx context.Context, in *NewUserBatchReq, opts ...grpc.CallOption) (*NewUserBatchResp, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(NewUserBatchResp)
	err := c.cc.Invoke(ctx, Logic_NewUserBatch_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *logicClient) DelUser(ctx context.Context, in *DelUserReq, opts ...grpc.CallOption) (*NoResp, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(NoResp)
	err := c.cc.Invoke(ctx, Logic_DelUser_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *logicClient) NewRoom(ctx context.Context, in *NewRoomReq, opts ...grpc.CallOption) (*NewRoomResp, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(NewRoomResp)
	err := c.cc.Invoke(ctx, Logic_NewRoom_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *logicClient) JoinRoom(ctx context.Context, in *JoinRoomReq, opts ...grpc.CallOption) (*NoResp, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(NoResp)
	err := c.cc.Invoke(ctx, Logic_JoinRoom_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *logicClient) JoinRoomBatch(ctx context.Context, in *JoinRoomBatchReq, opts ...grpc.CallOption) (*NoResp, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(NoResp)
	err := c.cc.Invoke(ctx, Logic_JoinRoomBatch_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *logicClient) QuitRoom(ctx context.Context, in *QuitRoomReq, opts ...grpc.CallOption) (*NoResp, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(NoResp)
	err := c.cc.Invoke(ctx, Logic_QuitRoom_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *logicClient) QueryRoom(ctx context.Context, in *QueryRoomReq, opts ...grpc.CallOption) (*QueryRoomResp, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(QueryRoomResp)
	err := c.cc.Invoke(ctx, Logic_QueryRoom_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *logicClient) QueryServer(ctx context.Context, in *QueryServerReq, opts ...grpc.CallOption) (*QueryServerResp, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(QueryServerResp)
	err := c.cc.Invoke(ctx, Logic_QueryServer_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *logicClient) Auth(ctx context.Context, in *AuthReq, opts ...grpc.CallOption) (*AuthResp, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(AuthResp)
	err := c.cc.Invoke(ctx, Logic_Auth_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// LogicServer is the server API for Logic service.
// All implementations should embed UnimplementedLogicServer
// for forward compatibility.
type LogicServer interface {
	SendMsg(context.Context, *msg.SendMsgReq) (*NoResp, error)
	SendMsgBatch(context.Context, *msg.SendMsgBatchReq) (*NoResp, error)
	SetOnline(context.Context, *SetOnlineReq) (*NoResp, error)
	SetOnlineBatch(context.Context, *SetOnlineBatchReq) (*NoResp, error)
	SetOffline(context.Context, *SetOfflineReq) (*NoResp, error)
	SetOfflineBatch(context.Context, *SetOfflineBatchReq) (*NoResp, error)
	NewUser(context.Context, *NewUserReq) (*NewUserResp, error)
	NewUserBatch(context.Context, *NewUserBatchReq) (*NewUserBatchResp, error)
	DelUser(context.Context, *DelUserReq) (*NoResp, error)
	NewRoom(context.Context, *NewRoomReq) (*NewRoomResp, error)
	JoinRoom(context.Context, *JoinRoomReq) (*NoResp, error)
	JoinRoomBatch(context.Context, *JoinRoomBatchReq) (*NoResp, error)
	QuitRoom(context.Context, *QuitRoomReq) (*NoResp, error)
	QueryRoom(context.Context, *QueryRoomReq) (*QueryRoomResp, error)
	QueryServer(context.Context, *QueryServerReq) (*QueryServerResp, error)
	Auth(context.Context, *AuthReq) (*AuthResp, error)
}

// UnimplementedLogicServer should be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedLogicServer struct{}

func (UnimplementedLogicServer) SendMsg(context.Context, *msg.SendMsgReq) (*NoResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendMsg not implemented")
}
func (UnimplementedLogicServer) SendMsgBatch(context.Context, *msg.SendMsgBatchReq) (*NoResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendMsgBatch not implemented")
}
func (UnimplementedLogicServer) SetOnline(context.Context, *SetOnlineReq) (*NoResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SetOnline not implemented")
}
func (UnimplementedLogicServer) SetOnlineBatch(context.Context, *SetOnlineBatchReq) (*NoResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SetOnlineBatch not implemented")
}
func (UnimplementedLogicServer) SetOffline(context.Context, *SetOfflineReq) (*NoResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SetOffline not implemented")
}
func (UnimplementedLogicServer) SetOfflineBatch(context.Context, *SetOfflineBatchReq) (*NoResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SetOfflineBatch not implemented")
}
func (UnimplementedLogicServer) NewUser(context.Context, *NewUserReq) (*NewUserResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method NewUser not implemented")
}
func (UnimplementedLogicServer) NewUserBatch(context.Context, *NewUserBatchReq) (*NewUserBatchResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method NewUserBatch not implemented")
}
func (UnimplementedLogicServer) DelUser(context.Context, *DelUserReq) (*NoResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DelUser not implemented")
}
func (UnimplementedLogicServer) NewRoom(context.Context, *NewRoomReq) (*NewRoomResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method NewRoom not implemented")
}
func (UnimplementedLogicServer) JoinRoom(context.Context, *JoinRoomReq) (*NoResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method JoinRoom not implemented")
}
func (UnimplementedLogicServer) JoinRoomBatch(context.Context, *JoinRoomBatchReq) (*NoResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method JoinRoomBatch not implemented")
}
func (UnimplementedLogicServer) QuitRoom(context.Context, *QuitRoomReq) (*NoResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method QuitRoom not implemented")
}
func (UnimplementedLogicServer) QueryRoom(context.Context, *QueryRoomReq) (*QueryRoomResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method QueryRoom not implemented")
}
func (UnimplementedLogicServer) QueryServer(context.Context, *QueryServerReq) (*QueryServerResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method QueryServer not implemented")
}
func (UnimplementedLogicServer) Auth(context.Context, *AuthReq) (*AuthResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Auth not implemented")
}
func (UnimplementedLogicServer) testEmbeddedByValue() {}

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

func _Logic_SendMsg_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(msg.SendMsgReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LogicServer).SendMsg(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Logic_SendMsg_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LogicServer).SendMsg(ctx, req.(*msg.SendMsgReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Logic_SendMsgBatch_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(msg.SendMsgBatchReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LogicServer).SendMsgBatch(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Logic_SendMsgBatch_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LogicServer).SendMsgBatch(ctx, req.(*msg.SendMsgBatchReq))
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

func _Logic_SetOnlineBatch_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SetOnlineBatchReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LogicServer).SetOnlineBatch(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Logic_SetOnlineBatch_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LogicServer).SetOnlineBatch(ctx, req.(*SetOnlineBatchReq))
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

func _Logic_SetOfflineBatch_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SetOfflineBatchReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LogicServer).SetOfflineBatch(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Logic_SetOfflineBatch_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LogicServer).SetOfflineBatch(ctx, req.(*SetOfflineBatchReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Logic_NewUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(NewUserReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LogicServer).NewUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Logic_NewUser_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LogicServer).NewUser(ctx, req.(*NewUserReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Logic_NewUserBatch_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(NewUserBatchReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LogicServer).NewUserBatch(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Logic_NewUserBatch_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LogicServer).NewUserBatch(ctx, req.(*NewUserBatchReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Logic_DelUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DelUserReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LogicServer).DelUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Logic_DelUser_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LogicServer).DelUser(ctx, req.(*DelUserReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Logic_NewRoom_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(NewRoomReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LogicServer).NewRoom(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Logic_NewRoom_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LogicServer).NewRoom(ctx, req.(*NewRoomReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Logic_JoinRoom_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(JoinRoomReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LogicServer).JoinRoom(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Logic_JoinRoom_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LogicServer).JoinRoom(ctx, req.(*JoinRoomReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Logic_JoinRoomBatch_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(JoinRoomBatchReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LogicServer).JoinRoomBatch(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Logic_JoinRoomBatch_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LogicServer).JoinRoomBatch(ctx, req.(*JoinRoomBatchReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Logic_QuitRoom_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QuitRoomReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LogicServer).QuitRoom(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Logic_QuitRoom_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LogicServer).QuitRoom(ctx, req.(*QuitRoomReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Logic_QueryRoom_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryRoomReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LogicServer).QueryRoom(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Logic_QueryRoom_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LogicServer).QueryRoom(ctx, req.(*QueryRoomReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Logic_QueryServer_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryServerReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LogicServer).QueryServer(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Logic_QueryServer_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LogicServer).QueryServer(ctx, req.(*QueryServerReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Logic_Auth_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AuthReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LogicServer).Auth(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Logic_Auth_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LogicServer).Auth(ctx, req.(*AuthReq))
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
			MethodName: "SendMsg",
			Handler:    _Logic_SendMsg_Handler,
		},
		{
			MethodName: "SendMsgBatch",
			Handler:    _Logic_SendMsgBatch_Handler,
		},
		{
			MethodName: "SetOnline",
			Handler:    _Logic_SetOnline_Handler,
		},
		{
			MethodName: "SetOnlineBatch",
			Handler:    _Logic_SetOnlineBatch_Handler,
		},
		{
			MethodName: "SetOffline",
			Handler:    _Logic_SetOffline_Handler,
		},
		{
			MethodName: "SetOfflineBatch",
			Handler:    _Logic_SetOfflineBatch_Handler,
		},
		{
			MethodName: "NewUser",
			Handler:    _Logic_NewUser_Handler,
		},
		{
			MethodName: "NewUserBatch",
			Handler:    _Logic_NewUserBatch_Handler,
		},
		{
			MethodName: "DelUser",
			Handler:    _Logic_DelUser_Handler,
		},
		{
			MethodName: "NewRoom",
			Handler:    _Logic_NewRoom_Handler,
		},
		{
			MethodName: "JoinRoom",
			Handler:    _Logic_JoinRoom_Handler,
		},
		{
			MethodName: "JoinRoomBatch",
			Handler:    _Logic_JoinRoomBatch_Handler,
		},
		{
			MethodName: "QuitRoom",
			Handler:    _Logic_QuitRoom_Handler,
		},
		{
			MethodName: "QueryRoom",
			Handler:    _Logic_QueryRoom_Handler,
		},
		{
			MethodName: "QueryServer",
			Handler:    _Logic_QueryServer_Handler,
		},
		{
			MethodName: "Auth",
			Handler:    _Logic_Auth_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/logic/logic.proto",
}
