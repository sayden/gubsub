// Code generated by protoc-gen-go.
// source: grpc/gprc-message.proto
// DO NOT EDIT!

/*
Package grpcservice is a generated protocol buffer package.

It is generated from these files:
	grpc/gprc-message.proto

It has these top-level messages:
	GubsubMessage
	GubsubReply
*/
package grpcservice

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
const _ = proto.ProtoPackageIsVersion1

type GubsubMessage struct {
	M []byte `protobuf:"bytes,1,opt,name=m,proto3" json:"m,omitempty"`
	T string `protobuf:"bytes,2,opt,name=t" json:"t,omitempty"`
}

func (m *GubsubMessage) Reset()                    { *m = GubsubMessage{} }
func (m *GubsubMessage) String() string            { return proto.CompactTextString(m) }
func (*GubsubMessage) ProtoMessage()               {}
func (*GubsubMessage) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

type GubsubReply struct {
	StatusCode int32 `protobuf:"varint,1,opt,name=statusCode" json:"statusCode,omitempty"`
}

func (m *GubsubReply) Reset()                    { *m = GubsubReply{} }
func (m *GubsubReply) String() string            { return proto.CompactTextString(m) }
func (*GubsubReply) ProtoMessage()               {}
func (*GubsubReply) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func init() {
	proto.RegisterType((*GubsubMessage)(nil), "grpcservice.GubsubMessage")
	proto.RegisterType((*GubsubReply)(nil), "grpcservice.GubsubReply")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// Client API for MessageService service

type MessageServiceClient interface {
	NewMessage(ctx context.Context, in *GubsubMessage, opts ...grpc.CallOption) (*GubsubReply, error)
}

type messageServiceClient struct {
	cc *grpc.ClientConn
}

func NewMessageServiceClient(cc *grpc.ClientConn) MessageServiceClient {
	return &messageServiceClient{cc}
}

func (c *messageServiceClient) NewMessage(ctx context.Context, in *GubsubMessage, opts ...grpc.CallOption) (*GubsubReply, error) {
	out := new(GubsubReply)
	err := grpc.Invoke(ctx, "/grpcservice.MessageService/NewMessage", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for MessageService service

type MessageServiceServer interface {
	NewMessage(context.Context, *GubsubMessage) (*GubsubReply, error)
}

func RegisterMessageServiceServer(s *grpc.Server, srv MessageServiceServer) {
	s.RegisterService(&_MessageService_serviceDesc, srv)
}

func _MessageService_NewMessage_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error) (interface{}, error) {
	in := new(GubsubMessage)
	if err := dec(in); err != nil {
		return nil, err
	}
	out, err := srv.(MessageServiceServer).NewMessage(ctx, in)
	if err != nil {
		return nil, err
	}
	return out, nil
}

var _MessageService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "grpcservice.MessageService",
	HandlerType: (*MessageServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "NewMessage",
			Handler:    _MessageService_NewMessage_Handler,
		},
	},
	Streams: []grpc.StreamDesc{},
}

var fileDescriptor0 = []byte{
	// 168 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0xe2, 0x12, 0x4f, 0x2f, 0x2a, 0x48,
	0xd6, 0x4f, 0x2f, 0x28, 0x4a, 0xd6, 0xcd, 0x4d, 0x2d, 0x2e, 0x4e, 0x4c, 0x4f, 0xd5, 0x2b, 0x28,
	0xca, 0x2f, 0xc9, 0x17, 0xe2, 0x06, 0x49, 0x14, 0xa7, 0x16, 0x95, 0x65, 0x26, 0xa7, 0x2a, 0x69,
	0x73, 0xf1, 0xba, 0x97, 0x26, 0x15, 0x97, 0x26, 0xf9, 0x42, 0xd4, 0x08, 0xf1, 0x70, 0x31, 0xe6,
	0x4a, 0x30, 0x2a, 0x30, 0x6a, 0xf0, 0x04, 0x31, 0xe6, 0x82, 0x78, 0x25, 0x12, 0x4c, 0x40, 0x1e,
	0x67, 0x10, 0x63, 0x89, 0x92, 0x2e, 0x17, 0x37, 0x44, 0x71, 0x50, 0x6a, 0x41, 0x4e, 0xa5, 0x90,
	0x1c, 0x17, 0x57, 0x71, 0x49, 0x62, 0x49, 0x69, 0xb1, 0x73, 0x7e, 0x4a, 0x2a, 0x58, 0x0f, 0x6b,
	0x10, 0x92, 0x88, 0x51, 0x18, 0x17, 0x1f, 0xd4, 0xd4, 0x60, 0x88, 0x6d, 0x42, 0x2e, 0x5c, 0x5c,
	0x7e, 0xa9, 0xe5, 0x30, 0xab, 0xa4, 0xf4, 0x90, 0x5c, 0xa2, 0x87, 0xe2, 0x0c, 0x29, 0x09, 0x2c,
	0x72, 0x60, 0x5b, 0x95, 0x18, 0x92, 0xd8, 0xc0, 0xfe, 0x30, 0x06, 0x04, 0x00, 0x00, 0xff, 0xff,
	0x50, 0x2c, 0x38, 0xf0, 0xe2, 0x00, 0x00, 0x00,
}
