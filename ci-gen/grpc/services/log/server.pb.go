// Code generated by protoc-gen-go. DO NOT EDIT.
// source: grpc/services/log/server.proto

package log

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	empty "github.com/golang/protobuf/ptypes/empty"
	_struct "github.com/golang/protobuf/ptypes/struct"
	timestamp "github.com/golang/protobuf/ptypes/timestamp"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

// LogMessage is a log message in struct form.
type LogMessage struct {
	At                   *timestamp.Timestamp `protobuf:"bytes,1,opt,name=at,proto3" json:"at,omitempty"`
	Level                string               `protobuf:"bytes,2,opt,name=level,proto3" json:"level,omitempty"`
	Fields               *_struct.Struct      `protobuf:"bytes,3,opt,name=fields,proto3" json:"fields,omitempty"`
	Service              string               `protobuf:"bytes,4,opt,name=service,proto3" json:"service,omitempty"`
	Message              string               `protobuf:"bytes,5,opt,name=message,proto3" json:"message,omitempty"`
	XXX_NoUnkeyedLiteral struct{}             `json:"-"`
	XXX_unrecognized     []byte               `json:"-"`
	XXX_sizecache        int32                `json:"-"`
}

func (m *LogMessage) Reset()         { *m = LogMessage{} }
func (m *LogMessage) String() string { return proto.CompactTextString(m) }
func (*LogMessage) ProtoMessage()    {}
func (*LogMessage) Descriptor() ([]byte, []int) {
	return fileDescriptor_02c5c7a49b3d9bb5, []int{0}
}

func (m *LogMessage) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_LogMessage.Unmarshal(m, b)
}
func (m *LogMessage) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_LogMessage.Marshal(b, m, deterministic)
}
func (m *LogMessage) XXX_Merge(src proto.Message) {
	xxx_messageInfo_LogMessage.Merge(m, src)
}
func (m *LogMessage) XXX_Size() int {
	return xxx_messageInfo_LogMessage.Size(m)
}
func (m *LogMessage) XXX_DiscardUnknown() {
	xxx_messageInfo_LogMessage.DiscardUnknown(m)
}

var xxx_messageInfo_LogMessage proto.InternalMessageInfo

func (m *LogMessage) GetAt() *timestamp.Timestamp {
	if m != nil {
		return m.At
	}
	return nil
}

func (m *LogMessage) GetLevel() string {
	if m != nil {
		return m.Level
	}
	return ""
}

func (m *LogMessage) GetFields() *_struct.Struct {
	if m != nil {
		return m.Fields
	}
	return nil
}

func (m *LogMessage) GetService() string {
	if m != nil {
		return m.Service
	}
	return ""
}

func (m *LogMessage) GetMessage() string {
	if m != nil {
		return m.Message
	}
	return ""
}

func init() {
	proto.RegisterType((*LogMessage)(nil), "log.LogMessage")
}

func init() { proto.RegisterFile("grpc/services/log/server.proto", fileDescriptor_02c5c7a49b3d9bb5) }

var fileDescriptor_02c5c7a49b3d9bb5 = []byte{
	// 292 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x84, 0x91, 0x41, 0x4b, 0x3b, 0x31,
	0x10, 0xc5, 0xd9, 0xee, 0xbf, 0xfd, 0x63, 0x3c, 0x08, 0x41, 0x34, 0xec, 0x41, 0x8a, 0xa7, 0x22,
	0x98, 0x40, 0x6b, 0xbd, 0x78, 0x13, 0xbc, 0x55, 0x90, 0xea, 0xc9, 0x5b, 0x1a, 0xa7, 0x63, 0x20,
	0xbb, 0x09, 0x9b, 0xd9, 0x42, 0x3f, 0x97, 0x5f, 0x50, 0x9a, 0x6c, 0xa9, 0xd0, 0x43, 0x2f, 0xbb,
	0xf3, 0x98, 0xbc, 0xdf, 0x4c, 0x5e, 0xd8, 0x0d, 0xb6, 0xc1, 0xa8, 0x08, 0xed, 0xc6, 0x1a, 0x88,
	0xca, 0x79, 0x4c, 0x02, 0x5a, 0x19, 0x5a, 0x4f, 0x9e, 0x97, 0xce, 0x63, 0xf5, 0x84, 0x96, 0xbe,
	0xbb, 0x95, 0x34, 0xbe, 0x56, 0xe8, 0x9d, 0x6e, 0x50, 0xa5, 0xee, 0xaa, 0x5b, 0xab, 0x40, 0xdb,
	0x00, 0x51, 0x91, 0xad, 0x21, 0x92, 0xae, 0xc3, 0xa1, 0xca, 0x84, 0x6a, 0x76, 0xda, 0x0c, 0x75,
	0xa0, 0x6d, 0xfe, 0xf6, 0xa6, 0xf9, 0x69, 0x53, 0xa4, 0xb6, 0x33, 0xd4, 0xff, 0xb2, 0xed, 0xf6,
	0xa7, 0x60, 0x6c, 0xe1, 0xf1, 0x15, 0x62, 0xd4, 0x08, 0xfc, 0x8e, 0x0d, 0x34, 0x89, 0x62, 0x5c,
	0x4c, 0xce, 0xa7, 0x95, 0x44, 0xef, 0xd1, 0x81, 0xdc, 0x73, 0xe4, 0xc7, 0x7e, 0xd1, 0xe5, 0x40,
	0x13, 0xbf, 0x64, 0x43, 0x07, 0x1b, 0x70, 0x62, 0x30, 0x2e, 0x26, 0x67, 0xcb, 0x2c, 0xb8, 0x62,
	0xa3, 0xb5, 0x05, 0xf7, 0x15, 0x45, 0x99, 0x28, 0xd7, 0x47, 0x94, 0xf7, 0x34, 0x7f, 0xd9, 0x1f,
	0xe3, 0x82, 0xfd, 0xef, 0xc3, 0x14, 0xff, 0x12, 0x68, 0x2f, 0x77, 0x9d, 0x3a, 0xef, 0x25, 0x86,
	0xb9, 0xd3, 0xcb, 0xe9, 0x9c, 0x95, 0x0b, 0x8f, 0x5c, 0xb2, 0xf2, 0xad, 0x23, 0x7e, 0x21, 0x9d,
	0x47, 0x79, 0xb8, 0x45, 0x75, 0x75, 0x34, 0xf3, 0x65, 0x97, 0xd4, 0xf3, 0xe3, 0xe7, 0xc3, 0x9f,
	0x94, 0xc8, 0x36, 0x5b, 0x63, 0x95, 0xb1, 0xf7, 0x1a, 0xa1, 0xa1, 0xb8, 0xab, 0x10, 0x1a, 0x75,
	0xf4, 0xbe, 0xab, 0x51, 0xe2, 0xcc, 0x7e, 0x03, 0x00, 0x00, 0xff, 0xff, 0xa1, 0xe8, 0x53, 0xfd,
	0xfb, 0x01, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// LogClient is the client API for Log service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type LogClient interface {
	Put(ctx context.Context, in *LogMessage, opts ...grpc.CallOption) (*empty.Empty, error)
}

type logClient struct {
	cc *grpc.ClientConn
}

func NewLogClient(cc *grpc.ClientConn) LogClient {
	return &logClient{cc}
}

func (c *logClient) Put(ctx context.Context, in *LogMessage, opts ...grpc.CallOption) (*empty.Empty, error) {
	out := new(empty.Empty)
	err := c.cc.Invoke(ctx, "/log.Log/Put", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// LogServer is the server API for Log service.
type LogServer interface {
	Put(context.Context, *LogMessage) (*empty.Empty, error)
}

// UnimplementedLogServer can be embedded to have forward compatible implementations.
type UnimplementedLogServer struct {
}

func (*UnimplementedLogServer) Put(ctx context.Context, req *LogMessage) (*empty.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Put not implemented")
}

func RegisterLogServer(s *grpc.Server, srv LogServer) {
	s.RegisterService(&_Log_serviceDesc, srv)
}

func _Log_Put_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LogMessage)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LogServer).Put(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/log.Log/Put",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LogServer).Put(ctx, req.(*LogMessage))
	}
	return interceptor(ctx, in, info, handler)
}

var _Log_serviceDesc = grpc.ServiceDesc{
	ServiceName: "log.Log",
	HandlerType: (*LogServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Put",
			Handler:    _Log_Put_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "grpc/services/log/server.proto",
}
