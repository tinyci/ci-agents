// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.15.6
// source: github.com/tinyci/ci-agents/ci-gen/grpc/types/user.proto

package types

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// User is ... a user record
type User struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id               int64                  `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`                            // ID of user
	Username         string                 `protobuf:"bytes,2,opt,name=username,proto3" json:"username,omitempty"`                 // Username -- retrieved from github
	LastScannedRepos *timestamppb.Timestamp `protobuf:"bytes,3,opt,name=lastScannedRepos,proto3" json:"lastScannedRepos,omitempty"` // This flag is used to described when someone last scanned repositories for adding to CI.
	Errors           []*UserError           `protobuf:"bytes,4,rep,name=errors,proto3" json:"errors,omitempty"`                     // Errors for the user. See types.Errors
	// JSON corresponding to the oauth2 response from github when first signing
	// up; this contains an access and refresh token. Encrypted with the token
	// key.
	TokenJSON []byte `protobuf:"bytes,5,opt,name=tokenJSON,proto3" json:"tokenJSON,omitempty"`
}

func (x *User) Reset() {
	*x = User{}
	if protoimpl.UnsafeEnabled {
		mi := &file_github_com_tinyci_ci_agents_ci_gen_grpc_types_user_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *User) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*User) ProtoMessage() {}

func (x *User) ProtoReflect() protoreflect.Message {
	mi := &file_github_com_tinyci_ci_agents_ci_gen_grpc_types_user_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use User.ProtoReflect.Descriptor instead.
func (*User) Descriptor() ([]byte, []int) {
	return file_github_com_tinyci_ci_agents_ci_gen_grpc_types_user_proto_rawDescGZIP(), []int{0}
}

func (x *User) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *User) GetUsername() string {
	if x != nil {
		return x.Username
	}
	return ""
}

func (x *User) GetLastScannedRepos() *timestamppb.Timestamp {
	if x != nil {
		return x.LastScannedRepos
	}
	return nil
}

func (x *User) GetErrors() []*UserError {
	if x != nil {
		return x.Errors
	}
	return nil
}

func (x *User) GetTokenJSON() []byte {
	if x != nil {
		return x.TokenJSON
	}
	return nil
}

// UserError is the pre-converted UserError record. It is later converted to a types.Error.
type UserError struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id     int64  `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	UserID int64  `protobuf:"varint,2,opt,name=userID,proto3" json:"userID,omitempty"`
	Error  string `protobuf:"bytes,3,opt,name=error,proto3" json:"error,omitempty"`
}

func (x *UserError) Reset() {
	*x = UserError{}
	if protoimpl.UnsafeEnabled {
		mi := &file_github_com_tinyci_ci_agents_ci_gen_grpc_types_user_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UserError) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UserError) ProtoMessage() {}

func (x *UserError) ProtoReflect() protoreflect.Message {
	mi := &file_github_com_tinyci_ci_agents_ci_gen_grpc_types_user_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UserError.ProtoReflect.Descriptor instead.
func (*UserError) Descriptor() ([]byte, []int) {
	return file_github_com_tinyci_ci_agents_ci_gen_grpc_types_user_proto_rawDescGZIP(), []int{1}
}

func (x *UserError) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *UserError) GetUserID() int64 {
	if x != nil {
		return x.UserID
	}
	return 0
}

func (x *UserError) GetError() string {
	if x != nil {
		return x.Error
	}
	return ""
}

// List of UserError
type UserErrors struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Errors []*UserError `protobuf:"bytes,1,rep,name=errors,proto3" json:"errors,omitempty"` // the list!
}

func (x *UserErrors) Reset() {
	*x = UserErrors{}
	if protoimpl.UnsafeEnabled {
		mi := &file_github_com_tinyci_ci_agents_ci_gen_grpc_types_user_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UserErrors) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UserErrors) ProtoMessage() {}

func (x *UserErrors) ProtoReflect() protoreflect.Message {
	mi := &file_github_com_tinyci_ci_agents_ci_gen_grpc_types_user_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UserErrors.ProtoReflect.Descriptor instead.
func (*UserErrors) Descriptor() ([]byte, []int) {
	return file_github_com_tinyci_ci_agents_ci_gen_grpc_types_user_proto_rawDescGZIP(), []int{2}
}

func (x *UserErrors) GetErrors() []*UserError {
	if x != nil {
		return x.Errors
	}
	return nil
}

// List of Users
type UserList struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Users []*User `protobuf:"bytes,1,rep,name=users,proto3" json:"users,omitempty"` // the list!
}

func (x *UserList) Reset() {
	*x = UserList{}
	if protoimpl.UnsafeEnabled {
		mi := &file_github_com_tinyci_ci_agents_ci_gen_grpc_types_user_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UserList) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UserList) ProtoMessage() {}

func (x *UserList) ProtoReflect() protoreflect.Message {
	mi := &file_github_com_tinyci_ci_agents_ci_gen_grpc_types_user_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UserList.ProtoReflect.Descriptor instead.
func (*UserList) Descriptor() ([]byte, []int) {
	return file_github_com_tinyci_ci_agents_ci_gen_grpc_types_user_proto_rawDescGZIP(), []int{3}
}

func (x *UserList) GetUsers() []*User {
	if x != nil {
		return x.Users
	}
	return nil
}

var File_github_com_tinyci_ci_agents_ci_gen_grpc_types_user_proto protoreflect.FileDescriptor

var file_github_com_tinyci_ci_agents_ci_gen_grpc_types_user_proto_rawDesc = []byte{
	0x0a, 0x38, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x74, 0x69, 0x6e,
	0x79, 0x63, 0x69, 0x2f, 0x63, 0x69, 0x2d, 0x61, 0x67, 0x65, 0x6e, 0x74, 0x73, 0x2f, 0x63, 0x69,
	0x2d, 0x67, 0x65, 0x6e, 0x2f, 0x67, 0x72, 0x70, 0x63, 0x2f, 0x74, 0x79, 0x70, 0x65, 0x73, 0x2f,
	0x75, 0x73, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x05, 0x74, 0x79, 0x70, 0x65,
	0x73, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62,
	0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x22, 0xc2, 0x01, 0x0a, 0x04, 0x55, 0x73, 0x65, 0x72, 0x12, 0x0e, 0x0a, 0x02, 0x69,
	0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x02, 0x69, 0x64, 0x12, 0x1a, 0x0a, 0x08, 0x75,
	0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x75,
	0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x46, 0x0a, 0x10, 0x6c, 0x61, 0x73, 0x74, 0x53,
	0x63, 0x61, 0x6e, 0x6e, 0x65, 0x64, 0x52, 0x65, 0x70, 0x6f, 0x73, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x10, 0x6c,
	0x61, 0x73, 0x74, 0x53, 0x63, 0x61, 0x6e, 0x6e, 0x65, 0x64, 0x52, 0x65, 0x70, 0x6f, 0x73, 0x12,
	0x28, 0x0a, 0x06, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x73, 0x18, 0x04, 0x20, 0x03, 0x28, 0x0b, 0x32,
	0x10, 0x2e, 0x74, 0x79, 0x70, 0x65, 0x73, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x45, 0x72, 0x72, 0x6f,
	0x72, 0x52, 0x06, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x73, 0x12, 0x1c, 0x0a, 0x09, 0x74, 0x6f, 0x6b,
	0x65, 0x6e, 0x4a, 0x53, 0x4f, 0x4e, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x09, 0x74, 0x6f,
	0x6b, 0x65, 0x6e, 0x4a, 0x53, 0x4f, 0x4e, 0x22, 0x49, 0x0a, 0x09, 0x55, 0x73, 0x65, 0x72, 0x45,
	0x72, 0x72, 0x6f, 0x72, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03,
	0x52, 0x02, 0x69, 0x64, 0x12, 0x16, 0x0a, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x44, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x03, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x44, 0x12, 0x14, 0x0a, 0x05,
	0x65, 0x72, 0x72, 0x6f, 0x72, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x65, 0x72, 0x72,
	0x6f, 0x72, 0x22, 0x36, 0x0a, 0x0a, 0x55, 0x73, 0x65, 0x72, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x73,
	0x12, 0x28, 0x0a, 0x06, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b,
	0x32, 0x10, 0x2e, 0x74, 0x79, 0x70, 0x65, 0x73, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x45, 0x72, 0x72,
	0x6f, 0x72, 0x52, 0x06, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x73, 0x22, 0x2d, 0x0a, 0x08, 0x55, 0x73,
	0x65, 0x72, 0x4c, 0x69, 0x73, 0x74, 0x12, 0x21, 0x0a, 0x05, 0x75, 0x73, 0x65, 0x72, 0x73, 0x18,
	0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0b, 0x2e, 0x74, 0x79, 0x70, 0x65, 0x73, 0x2e, 0x55, 0x73,
	0x65, 0x72, 0x52, 0x05, 0x75, 0x73, 0x65, 0x72, 0x73, 0x42, 0x2f, 0x5a, 0x2d, 0x67, 0x69, 0x74,
	0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x74, 0x69, 0x6e, 0x79, 0x63, 0x69, 0x2f, 0x63,
	0x69, 0x2d, 0x61, 0x67, 0x65, 0x6e, 0x74, 0x73, 0x2f, 0x63, 0x69, 0x2d, 0x67, 0x65, 0x6e, 0x2f,
	0x67, 0x72, 0x70, 0x63, 0x2f, 0x74, 0x79, 0x70, 0x65, 0x73, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x33,
}

var (
	file_github_com_tinyci_ci_agents_ci_gen_grpc_types_user_proto_rawDescOnce sync.Once
	file_github_com_tinyci_ci_agents_ci_gen_grpc_types_user_proto_rawDescData = file_github_com_tinyci_ci_agents_ci_gen_grpc_types_user_proto_rawDesc
)

func file_github_com_tinyci_ci_agents_ci_gen_grpc_types_user_proto_rawDescGZIP() []byte {
	file_github_com_tinyci_ci_agents_ci_gen_grpc_types_user_proto_rawDescOnce.Do(func() {
		file_github_com_tinyci_ci_agents_ci_gen_grpc_types_user_proto_rawDescData = protoimpl.X.CompressGZIP(file_github_com_tinyci_ci_agents_ci_gen_grpc_types_user_proto_rawDescData)
	})
	return file_github_com_tinyci_ci_agents_ci_gen_grpc_types_user_proto_rawDescData
}

var file_github_com_tinyci_ci_agents_ci_gen_grpc_types_user_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_github_com_tinyci_ci_agents_ci_gen_grpc_types_user_proto_goTypes = []interface{}{
	(*User)(nil),                  // 0: types.User
	(*UserError)(nil),             // 1: types.UserError
	(*UserErrors)(nil),            // 2: types.UserErrors
	(*UserList)(nil),              // 3: types.UserList
	(*timestamppb.Timestamp)(nil), // 4: google.protobuf.Timestamp
}
var file_github_com_tinyci_ci_agents_ci_gen_grpc_types_user_proto_depIdxs = []int32{
	4, // 0: types.User.lastScannedRepos:type_name -> google.protobuf.Timestamp
	1, // 1: types.User.errors:type_name -> types.UserError
	1, // 2: types.UserErrors.errors:type_name -> types.UserError
	0, // 3: types.UserList.users:type_name -> types.User
	4, // [4:4] is the sub-list for method output_type
	4, // [4:4] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_github_com_tinyci_ci_agents_ci_gen_grpc_types_user_proto_init() }
func file_github_com_tinyci_ci_agents_ci_gen_grpc_types_user_proto_init() {
	if File_github_com_tinyci_ci_agents_ci_gen_grpc_types_user_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_github_com_tinyci_ci_agents_ci_gen_grpc_types_user_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*User); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_github_com_tinyci_ci_agents_ci_gen_grpc_types_user_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UserError); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_github_com_tinyci_ci_agents_ci_gen_grpc_types_user_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UserErrors); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_github_com_tinyci_ci_agents_ci_gen_grpc_types_user_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UserList); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_github_com_tinyci_ci_agents_ci_gen_grpc_types_user_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_github_com_tinyci_ci_agents_ci_gen_grpc_types_user_proto_goTypes,
		DependencyIndexes: file_github_com_tinyci_ci_agents_ci_gen_grpc_types_user_proto_depIdxs,
		MessageInfos:      file_github_com_tinyci_ci_agents_ci_gen_grpc_types_user_proto_msgTypes,
	}.Build()
	File_github_com_tinyci_ci_agents_ci_gen_grpc_types_user_proto = out.File
	file_github_com_tinyci_ci_agents_ci_gen_grpc_types_user_proto_rawDesc = nil
	file_github_com_tinyci_ci_agents_ci_gen_grpc_types_user_proto_goTypes = nil
	file_github_com_tinyci_ci_agents_ci_gen_grpc_types_user_proto_depIdxs = nil
}
