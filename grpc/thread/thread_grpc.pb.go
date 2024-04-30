// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v5.26.1
// source: grpc/thread/pb/thread.proto

package thread

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

const (
	ThreadRpc_CreateThread_FullMethodName = "/thread.ThreadRpc/CreateThread"
)

// ThreadRpcClient is the client API for ThreadRpc service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ThreadRpcClient interface {
	CreateThread(ctx context.Context, in *ReqCreateThread, opts ...grpc.CallOption) (*Thread, error)
}

type threadRpcClient struct {
	cc grpc.ClientConnInterface
}

func NewThreadRpcClient(cc grpc.ClientConnInterface) ThreadRpcClient {
	return &threadRpcClient{cc}
}

func (c *threadRpcClient) CreateThread(ctx context.Context, in *ReqCreateThread, opts ...grpc.CallOption) (*Thread, error) {
	out := new(Thread)
	err := c.cc.Invoke(ctx, ThreadRpc_CreateThread_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ThreadRpcServer is the server API for ThreadRpc service.
// All implementations must embed UnimplementedThreadRpcServer
// for forward compatibility
type ThreadRpcServer interface {
	CreateThread(context.Context, *ReqCreateThread) (*Thread, error)
	mustEmbedUnimplementedThreadRpcServer()
}

// UnimplementedThreadRpcServer must be embedded to have forward compatible implementations.
type UnimplementedThreadRpcServer struct {
}

func (UnimplementedThreadRpcServer) CreateThread(context.Context, *ReqCreateThread) (*Thread, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateThread not implemented")
}
func (UnimplementedThreadRpcServer) mustEmbedUnimplementedThreadRpcServer() {}

// UnsafeThreadRpcServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ThreadRpcServer will
// result in compilation errors.
type UnsafeThreadRpcServer interface {
	mustEmbedUnimplementedThreadRpcServer()
}

func RegisterThreadRpcServer(s grpc.ServiceRegistrar, srv ThreadRpcServer) {
	s.RegisterService(&ThreadRpc_ServiceDesc, srv)
}

func _ThreadRpc_CreateThread_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ReqCreateThread)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ThreadRpcServer).CreateThread(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ThreadRpc_CreateThread_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ThreadRpcServer).CreateThread(ctx, req.(*ReqCreateThread))
	}
	return interceptor(ctx, in, info, handler)
}

// ThreadRpc_ServiceDesc is the grpc.ServiceDesc for ThreadRpc service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ThreadRpc_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "thread.ThreadRpc",
	HandlerType: (*ThreadRpcServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateThread",
			Handler:    _ThreadRpc_CreateThread_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "grpc/thread/pb/thread.proto",
}
