// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package node

import (
	context "context"
	msg "github.com/moratsam/cluster/api/v1/msg"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion7

// NodeClient is the client API for Node service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type NodeClient interface {
	Subscribe(ctx context.Context, in *SubscriptionRequest, opts ...grpc.CallOption) (Node_SubscribeClient, error)
	Publish(ctx context.Context, in *msg.Msg, opts ...grpc.CallOption) (*Ack, error)
}

type nodeClient struct {
	cc grpc.ClientConnInterface
}

func NewNodeClient(cc grpc.ClientConnInterface) NodeClient {
	return &nodeClient{cc}
}

func (c *nodeClient) Subscribe(ctx context.Context, in *SubscriptionRequest, opts ...grpc.CallOption) (Node_SubscribeClient, error) {
	stream, err := c.cc.NewStream(ctx, &_Node_serviceDesc.Streams[0], "/node.v1.Node/Subscribe", opts...)
	if err != nil {
		return nil, err
	}
	x := &nodeSubscribeClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Node_SubscribeClient interface {
	Recv() (*msg.Msg, error)
	grpc.ClientStream
}

type nodeSubscribeClient struct {
	grpc.ClientStream
}

func (x *nodeSubscribeClient) Recv() (*msg.Msg, error) {
	m := new(msg.Msg)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *nodeClient) Publish(ctx context.Context, in *msg.Msg, opts ...grpc.CallOption) (*Ack, error) {
	out := new(Ack)
	err := c.cc.Invoke(ctx, "/node.v1.Node/Publish", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// NodeServer is the server API for Node service.
// All implementations must embed UnimplementedNodeServer
// for forward compatibility
type NodeServer interface {
	Subscribe(*SubscriptionRequest, Node_SubscribeServer) error
	Publish(context.Context, *msg.Msg) (*Ack, error)
	mustEmbedUnimplementedNodeServer()
}

// UnimplementedNodeServer must be embedded to have forward compatible implementations.
type UnimplementedNodeServer struct {
}

func (UnimplementedNodeServer) Subscribe(*SubscriptionRequest, Node_SubscribeServer) error {
	return status.Errorf(codes.Unimplemented, "method Subscribe not implemented")
}
func (UnimplementedNodeServer) Publish(context.Context, *msg.Msg) (*Ack, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Publish not implemented")
}
func (UnimplementedNodeServer) mustEmbedUnimplementedNodeServer() {}

// UnsafeNodeServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to NodeServer will
// result in compilation errors.
type UnsafeNodeServer interface {
	mustEmbedUnimplementedNodeServer()
}

func RegisterNodeServer(s *grpc.Server, srv NodeServer) {
	s.RegisterService(&_Node_serviceDesc, srv)
}

func _Node_Subscribe_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(SubscriptionRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(NodeServer).Subscribe(m, &nodeSubscribeServer{stream})
}

type Node_SubscribeServer interface {
	Send(*msg.Msg) error
	grpc.ServerStream
}

type nodeSubscribeServer struct {
	grpc.ServerStream
}

func (x *nodeSubscribeServer) Send(m *msg.Msg) error {
	return x.ServerStream.SendMsg(m)
}

func _Node_Publish_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(msg.Msg)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NodeServer).Publish(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/node.v1.Node/Publish",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NodeServer).Publish(ctx, req.(*msg.Msg))
	}
	return interceptor(ctx, in, info, handler)
}

var _Node_serviceDesc = grpc.ServiceDesc{
	ServiceName: "node.v1.Node",
	HandlerType: (*NodeServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Publish",
			Handler:    _Node_Publish_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Subscribe",
			Handler:       _Node_Subscribe_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "api/v1/node/node.proto",
}