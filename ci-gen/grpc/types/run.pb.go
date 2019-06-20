// Code generated by protoc-gen-go. DO NOT EDIT.
// source: github.com/tinyci/ci-agents/ci-gen/grpc/types/run.proto

package types

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	timestamp "github.com/golang/protobuf/ptypes/timestamp"
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

// Run is a single CI run, intended to be sent to a runner.
type Run struct {
	Id                   int64                `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Name                 string               `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	CreatedAt            *timestamp.Timestamp `protobuf:"bytes,3,opt,name=createdAt,proto3" json:"createdAt,omitempty"`
	StartedAt            *timestamp.Timestamp `protobuf:"bytes,4,opt,name=startedAt,proto3" json:"startedAt,omitempty"`
	FinishedAt           *timestamp.Timestamp `protobuf:"bytes,5,opt,name=finishedAt,proto3" json:"finishedAt,omitempty"`
	Status               bool                 `protobuf:"varint,6,opt,name=status,proto3" json:"status,omitempty"`
	StatusSet            bool                 `protobuf:"varint,7,opt,name=statusSet,proto3" json:"statusSet,omitempty"`
	Settings             *RunSettings         `protobuf:"bytes,8,opt,name=settings,proto3" json:"settings,omitempty"`
	Task                 *Task                `protobuf:"bytes,9,opt,name=task,proto3" json:"task,omitempty"`
	XXX_NoUnkeyedLiteral struct{}             `json:"-"`
	XXX_unrecognized     []byte               `json:"-"`
	XXX_sizecache        int32                `json:"-"`
}

func (m *Run) Reset()         { *m = Run{} }
func (m *Run) String() string { return proto.CompactTextString(m) }
func (*Run) ProtoMessage()    {}
func (*Run) Descriptor() ([]byte, []int) {
	return fileDescriptor_b8e94412e957c3d8, []int{0}
}

func (m *Run) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Run.Unmarshal(m, b)
}
func (m *Run) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Run.Marshal(b, m, deterministic)
}
func (m *Run) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Run.Merge(m, src)
}
func (m *Run) XXX_Size() int {
	return xxx_messageInfo_Run.Size(m)
}
func (m *Run) XXX_DiscardUnknown() {
	xxx_messageInfo_Run.DiscardUnknown(m)
}

var xxx_messageInfo_Run proto.InternalMessageInfo

func (m *Run) GetId() int64 {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *Run) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *Run) GetCreatedAt() *timestamp.Timestamp {
	if m != nil {
		return m.CreatedAt
	}
	return nil
}

func (m *Run) GetStartedAt() *timestamp.Timestamp {
	if m != nil {
		return m.StartedAt
	}
	return nil
}

func (m *Run) GetFinishedAt() *timestamp.Timestamp {
	if m != nil {
		return m.FinishedAt
	}
	return nil
}

func (m *Run) GetStatus() bool {
	if m != nil {
		return m.Status
	}
	return false
}

func (m *Run) GetStatusSet() bool {
	if m != nil {
		return m.StatusSet
	}
	return false
}

func (m *Run) GetSettings() *RunSettings {
	if m != nil {
		return m.Settings
	}
	return nil
}

func (m *Run) GetTask() *Task {
	if m != nil {
		return m.Task
	}
	return nil
}

// RunList is just an array of runs
type RunList struct {
	List                 []*Run   `protobuf:"bytes,1,rep,name=list,proto3" json:"list,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RunList) Reset()         { *m = RunList{} }
func (m *RunList) String() string { return proto.CompactTextString(m) }
func (*RunList) ProtoMessage()    {}
func (*RunList) Descriptor() ([]byte, []int) {
	return fileDescriptor_b8e94412e957c3d8, []int{1}
}

func (m *RunList) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RunList.Unmarshal(m, b)
}
func (m *RunList) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RunList.Marshal(b, m, deterministic)
}
func (m *RunList) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RunList.Merge(m, src)
}
func (m *RunList) XXX_Size() int {
	return xxx_messageInfo_RunList.Size(m)
}
func (m *RunList) XXX_DiscardUnknown() {
	xxx_messageInfo_RunList.DiscardUnknown(m)
}

var xxx_messageInfo_RunList proto.InternalMessageInfo

func (m *RunList) GetList() []*Run {
	if m != nil {
		return m.List
	}
	return nil
}

func init() {
	proto.RegisterType((*Run)(nil), "types.Run")
	proto.RegisterType((*RunList)(nil), "types.RunList")
}

func init() {
	proto.RegisterFile("github.com/tinyci/ci-agents/ci-gen/grpc/types/run.proto", fileDescriptor_b8e94412e957c3d8)
}

var fileDescriptor_b8e94412e957c3d8 = []byte{
	// 341 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x94, 0x91, 0x3b, 0x6f, 0xc2, 0x30,
	0x10, 0xc7, 0x95, 0x07, 0x2f, 0x23, 0x75, 0xf0, 0x50, 0x59, 0xa8, 0x6a, 0x23, 0xa6, 0x74, 0xc0,
	0x96, 0xe8, 0x50, 0xd4, 0x2e, 0x6d, 0xe7, 0x4e, 0x86, 0xa9, 0x4b, 0x65, 0x82, 0x31, 0x16, 0xc4,
	0x89, 0xe2, 0xf3, 0xc0, 0xb7, 0xed, 0x47, 0xa9, 0x70, 0x02, 0x61, 0x43, 0xd9, 0xce, 0xf7, 0xff,
	0xff, 0xce, 0xf7, 0x40, 0xaf, 0x4a, 0xc3, 0xce, 0xad, 0x69, 0x56, 0xe4, 0x0c, 0xb4, 0x39, 0x66,
	0x9a, 0x65, 0x7a, 0x26, 0x94, 0x34, 0x60, 0x4f, 0x91, 0x92, 0x86, 0xa9, 0xaa, 0xcc, 0x18, 0x1c,
	0x4b, 0x69, 0x59, 0xe5, 0x0c, 0x2d, 0xab, 0x02, 0x0a, 0xdc, 0xf3, 0x89, 0xc9, 0xfb, 0x15, 0xaf,
	0x8a, 0x83, 0x30, 0x8a, 0x79, 0x7d, 0xed, 0xb6, 0xac, 0xac, 0x19, 0xd0, 0xb9, 0xb4, 0x20, 0xf2,
	0xb2, 0x8d, 0xea, 0x1a, 0x93, 0x45, 0xb7, 0xcf, 0x41, 0xd8, 0x7d, 0x43, 0x7e, 0x74, 0x6e, 0xfb,
	0xd7, 0x4a, 0x00, 0x6d, 0x94, 0xad, 0x2b, 0x4c, 0xff, 0x42, 0x14, 0x71, 0x67, 0xf0, 0x1d, 0x0a,
	0xf5, 0x86, 0x04, 0x49, 0x90, 0x46, 0x3c, 0xd4, 0x1b, 0x8c, 0x51, 0x6c, 0x44, 0x2e, 0x49, 0x98,
	0x04, 0xe9, 0x88, 0xfb, 0x18, 0x2f, 0xd0, 0x28, 0xab, 0xa4, 0x00, 0xb9, 0xf9, 0x04, 0x12, 0x25,
	0x41, 0x3a, 0x9e, 0x4f, 0xa8, 0x2a, 0x0a, 0x75, 0x90, 0xf4, 0x3c, 0x2d, 0x5d, 0x9d, 0x87, 0xe3,
	0xad, 0xf9, 0x44, 0x5a, 0x10, 0x55, 0x4d, 0xc6, 0xb7, 0xc9, 0x8b, 0x19, 0xbf, 0x21, 0xb4, 0xd5,
	0x46, 0xdb, 0x9d, 0x47, 0x7b, 0x37, 0xd1, 0x2b, 0x37, 0xbe, 0x47, 0x7d, 0x0b, 0x02, 0x9c, 0x25,
	0xfd, 0x24, 0x48, 0x87, 0xbc, 0x79, 0xe1, 0x07, 0xdf, 0x0d, 0x38, 0xbb, 0x94, 0x40, 0x06, 0x5e,
	0x6a, 0x13, 0x98, 0xa2, 0xe1, 0x79, 0x47, 0x64, 0xe8, 0xff, 0xc3, 0xd4, 0xaf, 0x8f, 0x72, 0x67,
	0x96, 0x8d, 0xc2, 0x2f, 0x1e, 0xfc, 0x84, 0xe2, 0xd3, 0x45, 0xc8, 0xc8, 0x7b, 0xc7, 0x8d, 0x77,
	0x25, 0xec, 0x9e, 0x7b, 0x61, 0xfa, 0x8c, 0x06, 0xdc, 0x99, 0x6f, 0x6d, 0x01, 0x3f, 0xa2, 0xf8,
	0xa0, 0x2d, 0x90, 0x20, 0x89, 0xd2, 0xf1, 0x1c, 0xb5, 0x75, 0xb9, 0xcf, 0x7f, 0xb1, 0x9f, 0x59,
	0xa7, 0x8b, 0xae, 0xfb, 0x7e, 0x05, 0x2f, 0xff, 0x01, 0x00, 0x00, 0xff, 0xff, 0x8c, 0x17, 0x4c,
	0x98, 0xc0, 0x02, 0x00, 0x00,
}