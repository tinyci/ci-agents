// Code generated by protoc-gen-go. DO NOT EDIT.
// source: github.com/tinyci/ci-agents/ci-gen/grpc/types/task.proto

package types

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	_struct "github.com/golang/protobuf/ptypes/struct"
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

// Task corresponds to directories within the tree that have a `task.yml`
// placed in them. Each task is decomposed into runs, and this record is
// created indicating the group of them, as well as properties they share.
type Task struct {
	Id                   int64                `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Parent               *Repository          `protobuf:"bytes,2,opt,name=parent,proto3" json:"parent,omitempty"`
	Ref                  *Ref                 `protobuf:"bytes,3,opt,name=ref,proto3" json:"ref,omitempty"`
	BaseSHA              string               `protobuf:"bytes,4,opt,name=baseSHA,proto3" json:"baseSHA,omitempty"`
	PullRequestID        int64                `protobuf:"varint,5,opt,name=pullRequestID,proto3" json:"pullRequestID,omitempty"`
	Canceled             bool                 `protobuf:"varint,6,opt,name=canceled,proto3" json:"canceled,omitempty"`
	FinishedAt           *timestamp.Timestamp `protobuf:"bytes,7,opt,name=finishedAt,proto3" json:"finishedAt,omitempty"`
	CreatedAt            *timestamp.Timestamp `protobuf:"bytes,8,opt,name=createdAt,proto3" json:"createdAt,omitempty"`
	StartedAt            *timestamp.Timestamp `protobuf:"bytes,9,opt,name=startedAt,proto3" json:"startedAt,omitempty"`
	Status               bool                 `protobuf:"varint,10,opt,name=status,proto3" json:"status,omitempty"`
	StatusSet            bool                 `protobuf:"varint,11,opt,name=statusSet,proto3" json:"statusSet,omitempty"`
	Settings             *TaskSettings        `protobuf:"bytes,12,opt,name=settings,proto3" json:"settings,omitempty"`
	Path                 string               `protobuf:"bytes,13,opt,name=path,proto3" json:"path,omitempty"`
	Runs                 int64                `protobuf:"varint,14,opt,name=runs,proto3" json:"runs,omitempty"`
	XXX_NoUnkeyedLiteral struct{}             `json:"-"`
	XXX_unrecognized     []byte               `json:"-"`
	XXX_sizecache        int32                `json:"-"`
}

func (m *Task) Reset()         { *m = Task{} }
func (m *Task) String() string { return proto.CompactTextString(m) }
func (*Task) ProtoMessage()    {}
func (*Task) Descriptor() ([]byte, []int) {
	return fileDescriptor_ca2db32a7ab25d95, []int{0}
}

func (m *Task) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Task.Unmarshal(m, b)
}
func (m *Task) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Task.Marshal(b, m, deterministic)
}
func (m *Task) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Task.Merge(m, src)
}
func (m *Task) XXX_Size() int {
	return xxx_messageInfo_Task.Size(m)
}
func (m *Task) XXX_DiscardUnknown() {
	xxx_messageInfo_Task.DiscardUnknown(m)
}

var xxx_messageInfo_Task proto.InternalMessageInfo

func (m *Task) GetId() int64 {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *Task) GetParent() *Repository {
	if m != nil {
		return m.Parent
	}
	return nil
}

func (m *Task) GetRef() *Ref {
	if m != nil {
		return m.Ref
	}
	return nil
}

func (m *Task) GetBaseSHA() string {
	if m != nil {
		return m.BaseSHA
	}
	return ""
}

func (m *Task) GetPullRequestID() int64 {
	if m != nil {
		return m.PullRequestID
	}
	return 0
}

func (m *Task) GetCanceled() bool {
	if m != nil {
		return m.Canceled
	}
	return false
}

func (m *Task) GetFinishedAt() *timestamp.Timestamp {
	if m != nil {
		return m.FinishedAt
	}
	return nil
}

func (m *Task) GetCreatedAt() *timestamp.Timestamp {
	if m != nil {
		return m.CreatedAt
	}
	return nil
}

func (m *Task) GetStartedAt() *timestamp.Timestamp {
	if m != nil {
		return m.StartedAt
	}
	return nil
}

func (m *Task) GetStatus() bool {
	if m != nil {
		return m.Status
	}
	return false
}

func (m *Task) GetStatusSet() bool {
	if m != nil {
		return m.StatusSet
	}
	return false
}

func (m *Task) GetSettings() *TaskSettings {
	if m != nil {
		return m.Settings
	}
	return nil
}

func (m *Task) GetPath() string {
	if m != nil {
		return m.Path
	}
	return ""
}

func (m *Task) GetRuns() int64 {
	if m != nil {
		return m.Runs
	}
	return 0
}

// TaskSettings is the parsed representation to struct of task.yml files.
type TaskSettings struct {
	Mountpoint           string                  `protobuf:"bytes,1,opt,name=mountpoint,proto3" json:"mountpoint,omitempty"`
	Env                  []string                `protobuf:"bytes,2,rep,name=env,proto3" json:"env,omitempty"`
	Workdir              string                  `protobuf:"bytes,3,opt,name=workdir,proto3" json:"workdir,omitempty"`
	Runs                 map[string]*RunSettings `protobuf:"bytes,4,rep,name=runs,proto3" json:"runs,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	DefaultTimeout       int64                   `protobuf:"varint,5,opt,name=defaultTimeout,proto3" json:"defaultTimeout,omitempty"`
	DefaultQueue         string                  `protobuf:"bytes,6,opt,name=defaultQueue,proto3" json:"defaultQueue,omitempty"`
	DefaultImage         string                  `protobuf:"bytes,7,opt,name=defaultImage,proto3" json:"defaultImage,omitempty"`
	Metadata             *_struct.Struct         `protobuf:"bytes,8,opt,name=metadata,proto3" json:"metadata,omitempty"`
	Dependencies         []string                `protobuf:"bytes,9,rep,name=dependencies,proto3" json:"dependencies,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                `json:"-"`
	XXX_unrecognized     []byte                  `json:"-"`
	XXX_sizecache        int32                   `json:"-"`
}

func (m *TaskSettings) Reset()         { *m = TaskSettings{} }
func (m *TaskSettings) String() string { return proto.CompactTextString(m) }
func (*TaskSettings) ProtoMessage()    {}
func (*TaskSettings) Descriptor() ([]byte, []int) {
	return fileDescriptor_ca2db32a7ab25d95, []int{1}
}

func (m *TaskSettings) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_TaskSettings.Unmarshal(m, b)
}
func (m *TaskSettings) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_TaskSettings.Marshal(b, m, deterministic)
}
func (m *TaskSettings) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TaskSettings.Merge(m, src)
}
func (m *TaskSettings) XXX_Size() int {
	return xxx_messageInfo_TaskSettings.Size(m)
}
func (m *TaskSettings) XXX_DiscardUnknown() {
	xxx_messageInfo_TaskSettings.DiscardUnknown(m)
}

var xxx_messageInfo_TaskSettings proto.InternalMessageInfo

func (m *TaskSettings) GetMountpoint() string {
	if m != nil {
		return m.Mountpoint
	}
	return ""
}

func (m *TaskSettings) GetEnv() []string {
	if m != nil {
		return m.Env
	}
	return nil
}

func (m *TaskSettings) GetWorkdir() string {
	if m != nil {
		return m.Workdir
	}
	return ""
}

func (m *TaskSettings) GetRuns() map[string]*RunSettings {
	if m != nil {
		return m.Runs
	}
	return nil
}

func (m *TaskSettings) GetDefaultTimeout() int64 {
	if m != nil {
		return m.DefaultTimeout
	}
	return 0
}

func (m *TaskSettings) GetDefaultQueue() string {
	if m != nil {
		return m.DefaultQueue
	}
	return ""
}

func (m *TaskSettings) GetDefaultImage() string {
	if m != nil {
		return m.DefaultImage
	}
	return ""
}

func (m *TaskSettings) GetMetadata() *_struct.Struct {
	if m != nil {
		return m.Metadata
	}
	return nil
}

func (m *TaskSettings) GetDependencies() []string {
	if m != nil {
		return m.Dependencies
	}
	return nil
}

// TaskList is simply a repeated list of tasks.
type TaskList struct {
	Tasks                []*Task  `protobuf:"bytes,1,rep,name=Tasks,proto3" json:"Tasks,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *TaskList) Reset()         { *m = TaskList{} }
func (m *TaskList) String() string { return proto.CompactTextString(m) }
func (*TaskList) ProtoMessage()    {}
func (*TaskList) Descriptor() ([]byte, []int) {
	return fileDescriptor_ca2db32a7ab25d95, []int{2}
}

func (m *TaskList) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_TaskList.Unmarshal(m, b)
}
func (m *TaskList) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_TaskList.Marshal(b, m, deterministic)
}
func (m *TaskList) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TaskList.Merge(m, src)
}
func (m *TaskList) XXX_Size() int {
	return xxx_messageInfo_TaskList.Size(m)
}
func (m *TaskList) XXX_DiscardUnknown() {
	xxx_messageInfo_TaskList.DiscardUnknown(m)
}

var xxx_messageInfo_TaskList proto.InternalMessageInfo

func (m *TaskList) GetTasks() []*Task {
	if m != nil {
		return m.Tasks
	}
	return nil
}

// CancelPRRequest is used in CancelTasksByPR in the datasvc; can be used to
// cancel all runs for a PR.
type CancelPRRequest struct {
	Id                   int64    `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Repository           string   `protobuf:"bytes,2,opt,name=repository,proto3" json:"repository,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CancelPRRequest) Reset()         { *m = CancelPRRequest{} }
func (m *CancelPRRequest) String() string { return proto.CompactTextString(m) }
func (*CancelPRRequest) ProtoMessage()    {}
func (*CancelPRRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_ca2db32a7ab25d95, []int{3}
}

func (m *CancelPRRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CancelPRRequest.Unmarshal(m, b)
}
func (m *CancelPRRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CancelPRRequest.Marshal(b, m, deterministic)
}
func (m *CancelPRRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CancelPRRequest.Merge(m, src)
}
func (m *CancelPRRequest) XXX_Size() int {
	return xxx_messageInfo_CancelPRRequest.Size(m)
}
func (m *CancelPRRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_CancelPRRequest.DiscardUnknown(m)
}

var xxx_messageInfo_CancelPRRequest proto.InternalMessageInfo

func (m *CancelPRRequest) GetId() int64 {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *CancelPRRequest) GetRepository() string {
	if m != nil {
		return m.Repository
	}
	return ""
}

func init() {
	proto.RegisterType((*Task)(nil), "types.Task")
	proto.RegisterType((*TaskSettings)(nil), "types.TaskSettings")
	proto.RegisterMapType((map[string]*RunSettings)(nil), "types.TaskSettings.RunsEntry")
	proto.RegisterType((*TaskList)(nil), "types.TaskList")
	proto.RegisterType((*CancelPRRequest)(nil), "types.CancelPRRequest")
}

func init() {
	proto.RegisterFile("github.com/tinyci/ci-agents/ci-gen/grpc/types/task.proto", fileDescriptor_ca2db32a7ab25d95)
}

var fileDescriptor_ca2db32a7ab25d95 = []byte{
	// 651 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x94, 0x54, 0x5f, 0x6f, 0xd3, 0x3e,
	0x14, 0x55, 0x9b, 0xb6, 0x6b, 0x6e, 0xb7, 0xfd, 0x7e, 0x18, 0x09, 0xac, 0x6a, 0x4c, 0xa5, 0x42,
	0xa8, 0x3c, 0x2c, 0x11, 0x9b, 0x10, 0xd3, 0x90, 0x10, 0xe3, 0x8f, 0xc4, 0x04, 0x0f, 0xe0, 0xee,
	0x89, 0x17, 0xe4, 0x26, 0x6e, 0x66, 0x35, 0x75, 0x42, 0x7c, 0x3d, 0xd4, 0xcf, 0xc3, 0xb7, 0xe4,
	0x09, 0xc5, 0x71, 0xd3, 0x8c, 0x21, 0x4d, 0x7d, 0xea, 0xf5, 0xf1, 0x39, 0xd7, 0xb7, 0xe7, 0xde,
	0x1b, 0x38, 0x4d, 0x24, 0x5e, 0x99, 0x59, 0x10, 0x65, 0xcb, 0x10, 0xa5, 0x5a, 0x45, 0x32, 0x8c,
	0xe4, 0x11, 0x4f, 0x84, 0x42, 0x5d, 0x46, 0x89, 0x50, 0x61, 0x52, 0xe4, 0x51, 0x88, 0xab, 0x5c,
	0xe8, 0x10, 0xb9, 0x5e, 0x04, 0x79, 0x91, 0x61, 0x46, 0xba, 0x16, 0x19, 0xbe, 0x6a, 0x24, 0x48,
	0xb2, 0x94, 0xab, 0x24, 0xb4, 0xf7, 0x33, 0x33, 0x0f, 0x73, 0x27, 0x92, 0x4b, 0xa1, 0x91, 0x2f,
	0xf3, 0x4d, 0x54, 0xe5, 0x18, 0xbe, 0xb8, 0x5b, 0xac, 0xb1, 0x30, 0x11, 0xba, 0x1f, 0x27, 0x7b,
	0xbd, 0x5d, 0xd1, 0x85, 0xc8, 0x33, 0x2d, 0x31, 0x2b, 0x56, 0x4e, 0xff, 0x72, 0x5b, 0xfd, 0xdc,
	0x09, 0xdf, 0x6c, 0x29, 0x34, 0xea, 0xbb, 0x16, 0x88, 0x52, 0x25, 0xba, 0xca, 0x30, 0xfe, 0xed,
	0x41, 0xe7, 0x92, 0xeb, 0x05, 0xd9, 0x87, 0xb6, 0x8c, 0x69, 0x6b, 0xd4, 0x9a, 0x78, 0xac, 0x2d,
	0x63, 0xf2, 0x0c, 0x7a, 0x39, 0x2f, 0x84, 0x42, 0xda, 0x1e, 0xb5, 0x26, 0x83, 0xe3, 0x7b, 0x81,
	0xcd, 0x11, 0xb0, 0xba, 0x78, 0xe6, 0x08, 0xe4, 0x00, 0xbc, 0x42, 0xcc, 0xa9, 0x67, 0x79, 0x50,
	0xf3, 0xe6, 0xac, 0x84, 0x09, 0x85, 0x9d, 0x19, 0xd7, 0x62, 0xfa, 0xf1, 0x9c, 0x76, 0x46, 0xad,
	0x89, 0xcf, 0xd6, 0x47, 0xf2, 0x04, 0xf6, 0x72, 0x93, 0xa6, 0x4c, 0xfc, 0x30, 0x42, 0xe3, 0xc5,
	0x7b, 0xda, 0xb5, 0xaf, 0xdf, 0x04, 0xc9, 0x10, 0xfa, 0x11, 0x57, 0x91, 0x48, 0x45, 0x4c, 0x7b,
	0xa3, 0xd6, 0xa4, 0xcf, 0xea, 0x33, 0x39, 0x03, 0x98, 0x4b, 0x25, 0xf5, 0x95, 0x88, 0xcf, 0x91,
	0xee, 0xd8, 0x02, 0x86, 0x41, 0x92, 0x65, 0x49, 0x2a, 0x82, 0x75, 0xe7, 0x82, 0xcb, 0x75, 0x97,
	0x59, 0x83, 0x4d, 0x4e, 0xc1, 0x8f, 0x0a, 0xc1, 0xd1, 0x4a, 0xfb, 0x77, 0x4a, 0x37, 0xe4, 0x52,
	0xa9, 0x91, 0x17, 0x95, 0xd2, 0xbf, 0x5b, 0x59, 0x93, 0xc9, 0x03, 0xe8, 0x69, 0xe4, 0x68, 0x34,
	0x05, 0xfb, 0x4f, 0xdc, 0x89, 0x1c, 0xd8, 0x8c, 0x68, 0xf4, 0x54, 0x20, 0x1d, 0xd8, 0xab, 0x0d,
	0x40, 0x42, 0xe8, 0xaf, 0xbb, 0x46, 0x77, 0xed, 0x73, 0xf7, 0x9d, 0xc9, 0x65, 0xe7, 0xa6, 0xee,
	0x8a, 0xd5, 0x24, 0x42, 0xa0, 0x93, 0x73, 0xbc, 0xa2, 0x7b, 0xd6, 0x6f, 0x1b, 0x97, 0x58, 0x61,
	0x94, 0xa6, 0xfb, 0xd6, 0x63, 0x1b, 0x8f, 0x7f, 0x79, 0xb0, 0xdb, 0x4c, 0x41, 0x0e, 0x01, 0x96,
	0x99, 0x51, 0x98, 0x67, 0x52, 0xa1, 0x1d, 0x06, 0x9f, 0x35, 0x10, 0xf2, 0x3f, 0x78, 0x42, 0x5d,
	0xd3, 0xf6, 0xc8, 0x9b, 0xf8, 0xac, 0x0c, 0xcb, 0xee, 0xfe, 0xcc, 0x8a, 0x45, 0x2c, 0x0b, 0xdb,
	0x7f, 0x9f, 0xad, 0x8f, 0xe4, 0xb9, 0x7b, 0xb0, 0x33, 0xf2, 0x26, 0x83, 0xe3, 0x47, 0xff, 0xa8,
	0x38, 0x60, 0x46, 0xe9, 0x0f, 0x0a, 0x8b, 0x55, 0x55, 0x0f, 0x79, 0x0a, 0xfb, 0xb1, 0x98, 0x73,
	0x93, 0x62, 0xe9, 0x5e, 0x66, 0xd0, 0x4d, 0xc4, 0x5f, 0x28, 0x19, 0xc3, 0xae, 0x43, 0xbe, 0x1a,
	0x61, 0x84, 0x1d, 0x0b, 0x9f, 0xdd, 0xc0, 0x1a, 0x9c, 0x8b, 0x25, 0x4f, 0x84, 0x1d, 0x8e, 0x0d,
	0xc7, 0x62, 0xe4, 0x04, 0xfa, 0x4b, 0x81, 0x3c, 0xe6, 0xc8, 0xdd, 0x04, 0x3c, 0xbc, 0xd5, 0xc7,
	0xa9, 0x5d, 0x74, 0x56, 0x13, 0xab, 0xc4, 0xb9, 0x50, 0xb1, 0x50, 0x91, 0x14, 0x9a, 0xfa, 0xd6,
	0x8c, 0x1b, 0xd8, 0xf0, 0x13, 0xf8, 0xf5, 0x7f, 0x2b, 0x4d, 0x5b, 0x88, 0x95, 0x73, 0xb3, 0x0c,
	0xc9, 0x04, 0xba, 0xd7, 0x3c, 0x35, 0xc2, 0xad, 0x16, 0x59, 0xaf, 0x8c, 0x51, 0x75, 0x33, 0x2b,
	0xc2, 0x59, 0xfb, 0xb4, 0x35, 0x3e, 0x82, 0x7e, 0xe9, 0xda, 0x67, 0xa9, 0x91, 0x3c, 0x86, 0x6e,
	0x19, 0x6b, 0xda, 0xb2, 0xae, 0x0e, 0x1a, 0xae, 0xb2, 0xea, 0x66, 0x7c, 0x0e, 0xff, 0xbd, 0xb3,
	0xfb, 0xf1, 0x85, 0xb9, 0x25, 0xba, 0xb5, 0xdb, 0x87, 0x00, 0x9b, 0x6f, 0x90, 0x2d, 0xc2, 0x67,
	0x0d, 0xe4, 0x6d, 0xf8, 0xed, 0x68, 0xab, 0x0f, 0xcb, 0xac, 0x67, 0xed, 0x3a, 0xf9, 0x13, 0x00,
	0x00, 0xff, 0xff, 0xc4, 0x0c, 0x0b, 0x1c, 0xbe, 0x05, 0x00, 0x00,
}