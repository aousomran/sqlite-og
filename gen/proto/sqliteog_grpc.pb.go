// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.20.3
// source: proto/sqliteog.proto

package sqlite_og

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

// SqliteOGClient is the client API for SqliteOG service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type SqliteOGClient interface {
	Query(ctx context.Context, in *Statement, opts ...grpc.CallOption) (*QueryResult, error)
	Execute(ctx context.Context, in *Statement, opts ...grpc.CallOption) (*ExecuteResult, error)
	ExecuteOrQuery(ctx context.Context, in *Statement, opts ...grpc.CallOption) (*ExecuteOrQueryResult, error)
	Callback(ctx context.Context, opts ...grpc.CallOption) (SqliteOG_CallbackClient, error)
	Connection(ctx context.Context, in *ConnectionRequest, opts ...grpc.CallOption) (*ConnectionId, error)
	Close(ctx context.Context, in *ConnectionId, opts ...grpc.CallOption) (*Empty, error)
	IsValid(ctx context.Context, in *ConnectionId, opts ...grpc.CallOption) (*Empty, error)
	Ping(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Empty, error)
	ResetSession(ctx context.Context, in *ConnectionId, opts ...grpc.CallOption) (*ConnectionId, error)
}

type sqliteOGClient struct {
	cc grpc.ClientConnInterface
}

func NewSqliteOGClient(cc grpc.ClientConnInterface) SqliteOGClient {
	return &sqliteOGClient{cc}
}

func (c *sqliteOGClient) Query(ctx context.Context, in *Statement, opts ...grpc.CallOption) (*QueryResult, error) {
	out := new(QueryResult)
	err := c.cc.Invoke(ctx, "/SqliteOG/Query", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *sqliteOGClient) Execute(ctx context.Context, in *Statement, opts ...grpc.CallOption) (*ExecuteResult, error) {
	out := new(ExecuteResult)
	err := c.cc.Invoke(ctx, "/SqliteOG/Execute", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *sqliteOGClient) ExecuteOrQuery(ctx context.Context, in *Statement, opts ...grpc.CallOption) (*ExecuteOrQueryResult, error) {
	out := new(ExecuteOrQueryResult)
	err := c.cc.Invoke(ctx, "/SqliteOG/ExecuteOrQuery", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *sqliteOGClient) Callback(ctx context.Context, opts ...grpc.CallOption) (SqliteOG_CallbackClient, error) {
	stream, err := c.cc.NewStream(ctx, &SqliteOG_ServiceDesc.Streams[0], "/SqliteOG/Callback", opts...)
	if err != nil {
		return nil, err
	}
	x := &sqliteOGCallbackClient{stream}
	return x, nil
}

type SqliteOG_CallbackClient interface {
	Send(*InvocationResult) error
	Recv() (*Invoke, error)
	grpc.ClientStream
}

type sqliteOGCallbackClient struct {
	grpc.ClientStream
}

func (x *sqliteOGCallbackClient) Send(m *InvocationResult) error {
	return x.ClientStream.SendMsg(m)
}

func (x *sqliteOGCallbackClient) Recv() (*Invoke, error) {
	m := new(Invoke)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *sqliteOGClient) Connection(ctx context.Context, in *ConnectionRequest, opts ...grpc.CallOption) (*ConnectionId, error) {
	out := new(ConnectionId)
	err := c.cc.Invoke(ctx, "/SqliteOG/Connection", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *sqliteOGClient) Close(ctx context.Context, in *ConnectionId, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/SqliteOG/Close", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *sqliteOGClient) IsValid(ctx context.Context, in *ConnectionId, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/SqliteOG/IsValid", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *sqliteOGClient) Ping(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/SqliteOG/Ping", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *sqliteOGClient) ResetSession(ctx context.Context, in *ConnectionId, opts ...grpc.CallOption) (*ConnectionId, error) {
	out := new(ConnectionId)
	err := c.cc.Invoke(ctx, "/SqliteOG/ResetSession", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// SqliteOGServer is the server API for SqliteOG service.
// All implementations must embed UnimplementedSqliteOGServer
// for forward compatibility
type SqliteOGServer interface {
	Query(context.Context, *Statement) (*QueryResult, error)
	Execute(context.Context, *Statement) (*ExecuteResult, error)
	ExecuteOrQuery(context.Context, *Statement) (*ExecuteOrQueryResult, error)
	Callback(SqliteOG_CallbackServer) error
	Connection(context.Context, *ConnectionRequest) (*ConnectionId, error)
	Close(context.Context, *ConnectionId) (*Empty, error)
	IsValid(context.Context, *ConnectionId) (*Empty, error)
	Ping(context.Context, *Empty) (*Empty, error)
	ResetSession(context.Context, *ConnectionId) (*ConnectionId, error)
	mustEmbedUnimplementedSqliteOGServer()
}

// UnimplementedSqliteOGServer must be embedded to have forward compatible implementations.
type UnimplementedSqliteOGServer struct {
}

func (UnimplementedSqliteOGServer) Query(context.Context, *Statement) (*QueryResult, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Query not implemented")
}
func (UnimplementedSqliteOGServer) Execute(context.Context, *Statement) (*ExecuteResult, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Execute not implemented")
}
func (UnimplementedSqliteOGServer) ExecuteOrQuery(context.Context, *Statement) (*ExecuteOrQueryResult, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ExecuteOrQuery not implemented")
}
func (UnimplementedSqliteOGServer) Callback(SqliteOG_CallbackServer) error {
	return status.Errorf(codes.Unimplemented, "method Callback not implemented")
}
func (UnimplementedSqliteOGServer) Connection(context.Context, *ConnectionRequest) (*ConnectionId, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Connection not implemented")
}
func (UnimplementedSqliteOGServer) Close(context.Context, *ConnectionId) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Close not implemented")
}
func (UnimplementedSqliteOGServer) IsValid(context.Context, *ConnectionId) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method IsValid not implemented")
}
func (UnimplementedSqliteOGServer) Ping(context.Context, *Empty) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Ping not implemented")
}
func (UnimplementedSqliteOGServer) ResetSession(context.Context, *ConnectionId) (*ConnectionId, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ResetSession not implemented")
}
func (UnimplementedSqliteOGServer) mustEmbedUnimplementedSqliteOGServer() {}

// UnsafeSqliteOGServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to SqliteOGServer will
// result in compilation errors.
type UnsafeSqliteOGServer interface {
	mustEmbedUnimplementedSqliteOGServer()
}

func RegisterSqliteOGServer(s grpc.ServiceRegistrar, srv SqliteOGServer) {
	s.RegisterService(&SqliteOG_ServiceDesc, srv)
}

func _SqliteOG_Query_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Statement)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SqliteOGServer).Query(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/SqliteOG/Query",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SqliteOGServer).Query(ctx, req.(*Statement))
	}
	return interceptor(ctx, in, info, handler)
}

func _SqliteOG_Execute_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Statement)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SqliteOGServer).Execute(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/SqliteOG/Execute",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SqliteOGServer).Execute(ctx, req.(*Statement))
	}
	return interceptor(ctx, in, info, handler)
}

func _SqliteOG_ExecuteOrQuery_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Statement)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SqliteOGServer).ExecuteOrQuery(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/SqliteOG/ExecuteOrQuery",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SqliteOGServer).ExecuteOrQuery(ctx, req.(*Statement))
	}
	return interceptor(ctx, in, info, handler)
}

func _SqliteOG_Callback_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(SqliteOGServer).Callback(&sqliteOGCallbackServer{stream})
}

type SqliteOG_CallbackServer interface {
	Send(*Invoke) error
	Recv() (*InvocationResult, error)
	grpc.ServerStream
}

type sqliteOGCallbackServer struct {
	grpc.ServerStream
}

func (x *sqliteOGCallbackServer) Send(m *Invoke) error {
	return x.ServerStream.SendMsg(m)
}

func (x *sqliteOGCallbackServer) Recv() (*InvocationResult, error) {
	m := new(InvocationResult)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _SqliteOG_Connection_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ConnectionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SqliteOGServer).Connection(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/SqliteOG/Connection",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SqliteOGServer).Connection(ctx, req.(*ConnectionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _SqliteOG_Close_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ConnectionId)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SqliteOGServer).Close(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/SqliteOG/Close",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SqliteOGServer).Close(ctx, req.(*ConnectionId))
	}
	return interceptor(ctx, in, info, handler)
}

func _SqliteOG_IsValid_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ConnectionId)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SqliteOGServer).IsValid(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/SqliteOG/IsValid",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SqliteOGServer).IsValid(ctx, req.(*ConnectionId))
	}
	return interceptor(ctx, in, info, handler)
}

func _SqliteOG_Ping_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SqliteOGServer).Ping(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/SqliteOG/Ping",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SqliteOGServer).Ping(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _SqliteOG_ResetSession_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ConnectionId)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SqliteOGServer).ResetSession(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/SqliteOG/ResetSession",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SqliteOGServer).ResetSession(ctx, req.(*ConnectionId))
	}
	return interceptor(ctx, in, info, handler)
}

// SqliteOG_ServiceDesc is the grpc.ServiceDesc for SqliteOG service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var SqliteOG_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "SqliteOG",
	HandlerType: (*SqliteOGServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Query",
			Handler:    _SqliteOG_Query_Handler,
		},
		{
			MethodName: "Execute",
			Handler:    _SqliteOG_Execute_Handler,
		},
		{
			MethodName: "ExecuteOrQuery",
			Handler:    _SqliteOG_ExecuteOrQuery_Handler,
		},
		{
			MethodName: "Connection",
			Handler:    _SqliteOG_Connection_Handler,
		},
		{
			MethodName: "Close",
			Handler:    _SqliteOG_Close_Handler,
		},
		{
			MethodName: "IsValid",
			Handler:    _SqliteOG_IsValid_Handler,
		},
		{
			MethodName: "Ping",
			Handler:    _SqliteOG_Ping_Handler,
		},
		{
			MethodName: "ResetSession",
			Handler:    _SqliteOG_ResetSession_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Callback",
			Handler:       _SqliteOG_Callback_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "proto/sqliteog.proto",
}
