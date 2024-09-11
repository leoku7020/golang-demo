// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.12
// source: example/example.proto

package example

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

// ExampleClient is the client API for Example service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ExampleClient interface {
	Login(ctx context.Context, in *LoginReq, opts ...grpc.CallOption) (*LoginResp, error)
	ListItems(ctx context.Context, in *ListItemsReq, opts ...grpc.CallOption) (*ListItemsResp, error)
}

type exampleClient struct {
	cc grpc.ClientConnInterface
}

func NewExampleClient(cc grpc.ClientConnInterface) ExampleClient {
	return &exampleClient{cc}
}

func (c *exampleClient) Login(ctx context.Context, in *LoginReq, opts ...grpc.CallOption) (*LoginResp, error) {
	out := new(LoginResp)
	err := c.cc.Invoke(ctx, "/example.Example/Login", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *exampleClient) ListItems(ctx context.Context, in *ListItemsReq, opts ...grpc.CallOption) (*ListItemsResp, error) {
	out := new(ListItemsResp)
	err := c.cc.Invoke(ctx, "/example.Example/ListItems", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ExampleServer is the server API for Example service.
// All implementations must embed UnimplementedExampleServer
// for forward compatibility
type ExampleServer interface {
	Login(context.Context, *LoginReq) (*LoginResp, error)
	ListItems(context.Context, *ListItemsReq) (*ListItemsResp, error)
	mustEmbedUnimplementedExampleServer()
}

// UnimplementedExampleServer must be embedded to have forward compatible implementations.
type UnimplementedExampleServer struct {
}

func (UnimplementedExampleServer) Login(context.Context, *LoginReq) (*LoginResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Login not implemented")
}
func (UnimplementedExampleServer) ListItems(context.Context, *ListItemsReq) (*ListItemsResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListItems not implemented")
}
func (UnimplementedExampleServer) mustEmbedUnimplementedExampleServer() {}

// UnsafeExampleServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ExampleServer will
// result in compilation errors.
type UnsafeExampleServer interface {
	mustEmbedUnimplementedExampleServer()
}

func RegisterExampleServer(s grpc.ServiceRegistrar, srv ExampleServer) {
	s.RegisterService(&Example_ServiceDesc, srv)
}

func _Example_Login_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LoginReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ExampleServer).Login(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/example.Example/Login",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ExampleServer).Login(ctx, req.(*LoginReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Example_ListItems_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListItemsReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ExampleServer).ListItems(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/example.Example/ListItems",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ExampleServer).ListItems(ctx, req.(*ListItemsReq))
	}
	return interceptor(ctx, in, info, handler)
}

// Example_ServiceDesc is the grpc.ServiceDesc for Example service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Example_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "example.Example",
	HandlerType: (*ExampleServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Login",
			Handler:    _Example_Login_Handler,
		},
		{
			MethodName: "ListItems",
			Handler:    _Example_ListItems_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "example/example.proto",
}
