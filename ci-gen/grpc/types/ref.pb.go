// Code generated by protoc-gen-go. DO NOT EDIT.
// source: github.com/tinyci/ci-agents/ci-gen/grpc/types/ref.proto

package types

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
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

// Ref is the encapsulation of a git ref and communicates repository as well as version information.
type Ref struct {
	Id                   int64       `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Repository           *Repository `protobuf:"bytes,2,opt,name=repository,proto3" json:"repository,omitempty"`
	RefName              string      `protobuf:"bytes,3,opt,name=refName,proto3" json:"refName,omitempty"`
	Sha                  string      `protobuf:"bytes,4,opt,name=sha,proto3" json:"sha,omitempty"`
	XXX_NoUnkeyedLiteral struct{}    `json:"-"`
	XXX_unrecognized     []byte      `json:"-"`
	XXX_sizecache        int32       `json:"-"`
}

func (m *Ref) Reset()         { *m = Ref{} }
func (m *Ref) String() string { return proto.CompactTextString(m) }
func (*Ref) ProtoMessage()    {}
func (*Ref) Descriptor() ([]byte, []int) {
	return fileDescriptor_b85b5c6cf1323432, []int{0}
}

func (m *Ref) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Ref.Unmarshal(m, b)
}
func (m *Ref) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Ref.Marshal(b, m, deterministic)
}
func (m *Ref) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Ref.Merge(m, src)
}
func (m *Ref) XXX_Size() int {
	return xxx_messageInfo_Ref.Size(m)
}
func (m *Ref) XXX_DiscardUnknown() {
	xxx_messageInfo_Ref.DiscardUnknown(m)
}

var xxx_messageInfo_Ref proto.InternalMessageInfo

func (m *Ref) GetId() int64 {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *Ref) GetRepository() *Repository {
	if m != nil {
		return m.Repository
	}
	return nil
}

func (m *Ref) GetRefName() string {
	if m != nil {
		return m.RefName
	}
	return ""
}

func (m *Ref) GetSha() string {
	if m != nil {
		return m.Sha
	}
	return ""
}

func init() {
	proto.RegisterType((*Ref)(nil), "types.Ref")
}

func init() {
	proto.RegisterFile("github.com/tinyci/ci-agents/ci-gen/grpc/types/ref.proto", fileDescriptor_b85b5c6cf1323432)
}

var fileDescriptor_b85b5c6cf1323432 = []byte{
	// 189 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x32, 0x4f, 0xcf, 0x2c, 0xc9,
	0x28, 0x4d, 0xd2, 0x4b, 0xce, 0xcf, 0xd5, 0x2f, 0xc9, 0xcc, 0xab, 0x4c, 0xce, 0xd4, 0x4f, 0xce,
	0xd4, 0x4d, 0x4c, 0x4f, 0xcd, 0x2b, 0x29, 0x06, 0xb1, 0xd2, 0x53, 0xf3, 0xf4, 0xd3, 0x8b, 0x0a,
	0x92, 0xf5, 0x4b, 0x2a, 0x0b, 0x52, 0x8b, 0xf5, 0x8b, 0x52, 0xd3, 0xf4, 0x0a, 0x8a, 0xf2, 0x4b,
	0xf2, 0x85, 0x58, 0xc1, 0x02, 0x52, 0x76, 0xa4, 0xea, 0x2f, 0xc8, 0x2f, 0xce, 0x2c, 0xc9, 0x2f,
	0xaa, 0x84, 0x18, 0xa3, 0x54, 0xc2, 0xc5, 0x1c, 0x94, 0x9a, 0x26, 0xc4, 0xc7, 0xc5, 0x94, 0x99,
	0x22, 0xc1, 0xa8, 0xc0, 0xa8, 0xc1, 0x1c, 0xc4, 0x94, 0x99, 0x22, 0x64, 0xc8, 0xc5, 0x85, 0x50,
	0x2a, 0xc1, 0xa4, 0xc0, 0xa8, 0xc1, 0x6d, 0x24, 0xa8, 0x07, 0x36, 0x43, 0x2f, 0x08, 0x2e, 0x11,
	0x84, 0xa4, 0x48, 0x48, 0x82, 0x8b, 0xbd, 0x28, 0x35, 0xcd, 0x2f, 0x31, 0x37, 0x55, 0x82, 0x59,
	0x81, 0x51, 0x83, 0x33, 0x08, 0xc6, 0x15, 0x12, 0xe0, 0x62, 0x2e, 0xce, 0x48, 0x94, 0x60, 0x01,
	0x8b, 0x82, 0x98, 0x4e, 0xfa, 0x51, 0xba, 0x24, 0xb9, 0x3b, 0x89, 0x0d, 0xec, 0x5a, 0x63, 0x40,
	0x00, 0x00, 0x00, 0xff, 0xff, 0xae, 0x73, 0xcb, 0x1c, 0x2f, 0x01, 0x00, 0x00,
}
