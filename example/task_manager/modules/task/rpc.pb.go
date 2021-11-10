// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        v3.17.3
// source: rpc.proto

package task

import (
	reflect "reflect"
	sync "sync"

	_ "github.com/ronaksoft/rony"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type CreateRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Title   string   `protobuf:"bytes,1,opt,name=Title,proto3" json:"Title,omitempty"`
	TODOs   []string `protobuf:"bytes,2,rep,name=TODOs,proto3" json:"TODOs,omitempty"`
	DueDate int64    `protobuf:"varint,3,opt,name=DueDate,proto3" json:"DueDate,omitempty"`
}

func (x *CreateRequest) Reset() {
	*x = CreateRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_rpc_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateRequest) ProtoMessage() {}

func (x *CreateRequest) ProtoReflect() protoreflect.Message {
	mi := &file_rpc_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateRequest.ProtoReflect.Descriptor instead.
func (*CreateRequest) Descriptor() ([]byte, []int) {
	return file_rpc_proto_rawDescGZIP(), []int{0}
}

func (x *CreateRequest) GetTitle() string {
	if x != nil {
		return x.Title
	}
	return ""
}

func (x *CreateRequest) GetTODOs() []string {
	if x != nil {
		return x.TODOs
	}
	return nil
}

func (x *CreateRequest) GetDueDate() int64 {
	if x != nil {
		return x.DueDate
	}
	return 0
}

type GetRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	TaskID int64 `protobuf:"varint,1,opt,name=TaskID,proto3" json:"TaskID,omitempty"`
}

func (x *GetRequest) Reset() {
	*x = GetRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_rpc_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetRequest) ProtoMessage() {}

func (x *GetRequest) ProtoReflect() protoreflect.Message {
	mi := &file_rpc_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetRequest.ProtoReflect.Descriptor instead.
func (*GetRequest) Descriptor() ([]byte, []int) {
	return file_rpc_proto_rawDescGZIP(), []int{1}
}

func (x *GetRequest) GetTaskID() int64 {
	if x != nil {
		return x.TaskID
	}
	return 0
}

type DeleteRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	TaskID int64 `protobuf:"varint,1,opt,name=TaskID,proto3" json:"TaskID,omitempty"`
}

func (x *DeleteRequest) Reset() {
	*x = DeleteRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_rpc_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeleteRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteRequest) ProtoMessage() {}

func (x *DeleteRequest) ProtoReflect() protoreflect.Message {
	mi := &file_rpc_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteRequest.ProtoReflect.Descriptor instead.
func (*DeleteRequest) Descriptor() ([]byte, []int) {
	return file_rpc_proto_rawDescGZIP(), []int{2}
}

func (x *DeleteRequest) GetTaskID() int64 {
	if x != nil {
		return x.TaskID
	}
	return 0
}

type ListRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Offset int32 `protobuf:"varint,1,opt,name=Offset,proto3" json:"Offset,omitempty"`
	Limit  int32 `protobuf:"varint,2,opt,name=Limit,proto3" json:"Limit,omitempty"`
}

func (x *ListRequest) Reset() {
	*x = ListRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_rpc_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListRequest) ProtoMessage() {}

func (x *ListRequest) ProtoReflect() protoreflect.Message {
	mi := &file_rpc_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListRequest.ProtoReflect.Descriptor instead.
func (*ListRequest) Descriptor() ([]byte, []int) {
	return file_rpc_proto_rawDescGZIP(), []int{3}
}

func (x *ListRequest) GetOffset() int32 {
	if x != nil {
		return x.Offset
	}
	return 0
}

func (x *ListRequest) GetLimit() int32 {
	if x != nil {
		return x.Limit
	}
	return 0
}

type TaskView struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ID      int64    `protobuf:"varint,1,opt,name=ID,proto3" json:"ID,omitempty"`
	Title   string   `protobuf:"bytes,2,opt,name=Title,proto3" json:"Title,omitempty"`
	TODOs   []string `protobuf:"bytes,3,rep,name=TODOs,proto3" json:"TODOs,omitempty"`
	DueDate int64    `protobuf:"varint,4,opt,name=DueDate,proto3" json:"DueDate,omitempty"`
}

func (x *TaskView) Reset() {
	*x = TaskView{}
	if protoimpl.UnsafeEnabled {
		mi := &file_rpc_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TaskView) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TaskView) ProtoMessage() {}

func (x *TaskView) ProtoReflect() protoreflect.Message {
	mi := &file_rpc_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TaskView.ProtoReflect.Descriptor instead.
func (*TaskView) Descriptor() ([]byte, []int) {
	return file_rpc_proto_rawDescGZIP(), []int{4}
}

func (x *TaskView) GetID() int64 {
	if x != nil {
		return x.ID
	}
	return 0
}

func (x *TaskView) GetTitle() string {
	if x != nil {
		return x.Title
	}
	return ""
}

func (x *TaskView) GetTODOs() []string {
	if x != nil {
		return x.TODOs
	}
	return nil
}

func (x *TaskView) GetDueDate() int64 {
	if x != nil {
		return x.DueDate
	}
	return 0
}

type TaskViewMany struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Tasks []*TaskView `protobuf:"bytes,1,rep,name=Tasks,proto3" json:"Tasks,omitempty"`
}

func (x *TaskViewMany) Reset() {
	*x = TaskViewMany{}
	if protoimpl.UnsafeEnabled {
		mi := &file_rpc_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TaskViewMany) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TaskViewMany) ProtoMessage() {}

func (x *TaskViewMany) ProtoReflect() protoreflect.Message {
	mi := &file_rpc_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TaskViewMany.ProtoReflect.Descriptor instead.
func (*TaskViewMany) Descriptor() ([]byte, []int) {
	return file_rpc_proto_rawDescGZIP(), []int{5}
}

func (x *TaskViewMany) GetTasks() []*TaskView {
	if x != nil {
		return x.Tasks
	}
	return nil
}

type Bool struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Result bool `protobuf:"varint,1,opt,name=Result,proto3" json:"Result,omitempty"`
}

func (x *Bool) Reset() {
	*x = Bool{}
	if protoimpl.UnsafeEnabled {
		mi := &file_rpc_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Bool) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Bool) ProtoMessage() {}

func (x *Bool) ProtoReflect() protoreflect.Message {
	mi := &file_rpc_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Bool.ProtoReflect.Descriptor instead.
func (*Bool) Descriptor() ([]byte, []int) {
	return file_rpc_proto_rawDescGZIP(), []int{6}
}

func (x *Bool) GetResult() bool {
	if x != nil {
		return x.Result
	}
	return false
}

var File_rpc_proto protoreflect.FileDescriptor

var file_rpc_proto_rawDesc = []byte{
	0x0a, 0x09, 0x72, 0x70, 0x63, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x04, 0x74, 0x61, 0x73,
	0x6b, 0x1a, 0x0d, 0x6f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x22, 0x55, 0x0a, 0x0d, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x12, 0x14, 0x0a, 0x05, 0x54, 0x69, 0x74, 0x6c, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x05, 0x54, 0x69, 0x74, 0x6c, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x54, 0x4f, 0x44, 0x4f, 0x73,
	0x18, 0x02, 0x20, 0x03, 0x28, 0x09, 0x52, 0x05, 0x54, 0x4f, 0x44, 0x4f, 0x73, 0x12, 0x18, 0x0a,
	0x07, 0x44, 0x75, 0x65, 0x44, 0x61, 0x74, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x52, 0x07,
	0x44, 0x75, 0x65, 0x44, 0x61, 0x74, 0x65, 0x22, 0x24, 0x0a, 0x0a, 0x47, 0x65, 0x74, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x16, 0x0a, 0x06, 0x54, 0x61, 0x73, 0x6b, 0x49, 0x44, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x06, 0x54, 0x61, 0x73, 0x6b, 0x49, 0x44, 0x22, 0x27, 0x0a,
	0x0d, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x16,
	0x0a, 0x06, 0x54, 0x61, 0x73, 0x6b, 0x49, 0x44, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x06,
	0x54, 0x61, 0x73, 0x6b, 0x49, 0x44, 0x22, 0x3b, 0x0a, 0x0b, 0x4c, 0x69, 0x73, 0x74, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x16, 0x0a, 0x06, 0x4f, 0x66, 0x66, 0x73, 0x65, 0x74, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x06, 0x4f, 0x66, 0x66, 0x73, 0x65, 0x74, 0x12, 0x14, 0x0a,
	0x05, 0x4c, 0x69, 0x6d, 0x69, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x05, 0x4c, 0x69,
	0x6d, 0x69, 0x74, 0x22, 0x60, 0x0a, 0x08, 0x54, 0x61, 0x73, 0x6b, 0x56, 0x69, 0x65, 0x77, 0x12,
	0x0e, 0x0a, 0x02, 0x49, 0x44, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x02, 0x49, 0x44, 0x12,
	0x14, 0x0a, 0x05, 0x54, 0x69, 0x74, 0x6c, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05,
	0x54, 0x69, 0x74, 0x6c, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x54, 0x4f, 0x44, 0x4f, 0x73, 0x18, 0x03,
	0x20, 0x03, 0x28, 0x09, 0x52, 0x05, 0x54, 0x4f, 0x44, 0x4f, 0x73, 0x12, 0x18, 0x0a, 0x07, 0x44,
	0x75, 0x65, 0x44, 0x61, 0x74, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x03, 0x52, 0x07, 0x44, 0x75,
	0x65, 0x44, 0x61, 0x74, 0x65, 0x22, 0x34, 0x0a, 0x0c, 0x54, 0x61, 0x73, 0x6b, 0x56, 0x69, 0x65,
	0x77, 0x4d, 0x61, 0x6e, 0x79, 0x12, 0x24, 0x0a, 0x05, 0x54, 0x61, 0x73, 0x6b, 0x73, 0x18, 0x01,
	0x20, 0x03, 0x28, 0x0b, 0x32, 0x0e, 0x2e, 0x74, 0x61, 0x73, 0x6b, 0x2e, 0x54, 0x61, 0x73, 0x6b,
	0x56, 0x69, 0x65, 0x77, 0x52, 0x05, 0x54, 0x61, 0x73, 0x6b, 0x73, 0x22, 0x1e, 0x0a, 0x04, 0x42,
	0x6f, 0x6f, 0x6c, 0x12, 0x16, 0x0a, 0x06, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x08, 0x52, 0x06, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x32, 0xc5, 0x01, 0x0a, 0x0b,
	0x54, 0x61, 0x73, 0x6b, 0x4d, 0x61, 0x6e, 0x61, 0x67, 0x65, 0x72, 0x12, 0x2d, 0x0a, 0x06, 0x43,
	0x72, 0x65, 0x61, 0x74, 0x65, 0x12, 0x13, 0x2e, 0x74, 0x61, 0x73, 0x6b, 0x2e, 0x43, 0x72, 0x65,
	0x61, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x0e, 0x2e, 0x74, 0x61, 0x73,
	0x6b, 0x2e, 0x54, 0x61, 0x73, 0x6b, 0x56, 0x69, 0x65, 0x77, 0x12, 0x27, 0x0a, 0x03, 0x47, 0x65,
	0x74, 0x12, 0x10, 0x2e, 0x74, 0x61, 0x73, 0x6b, 0x2e, 0x47, 0x65, 0x74, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x1a, 0x0e, 0x2e, 0x74, 0x61, 0x73, 0x6b, 0x2e, 0x54, 0x61, 0x73, 0x6b, 0x56,
	0x69, 0x65, 0x77, 0x12, 0x29, 0x0a, 0x06, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x12, 0x13, 0x2e,
	0x74, 0x61, 0x73, 0x6b, 0x2e, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x1a, 0x0a, 0x2e, 0x74, 0x61, 0x73, 0x6b, 0x2e, 0x42, 0x6f, 0x6f, 0x6c, 0x12, 0x2d,
	0x0a, 0x04, 0x4c, 0x69, 0x73, 0x74, 0x12, 0x11, 0x2e, 0x74, 0x61, 0x73, 0x6b, 0x2e, 0x4c, 0x69,
	0x73, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x12, 0x2e, 0x74, 0x61, 0x73, 0x6b,
	0x2e, 0x54, 0x61, 0x73, 0x6b, 0x56, 0x69, 0x65, 0x77, 0x4d, 0x61, 0x6e, 0x79, 0x1a, 0x04, 0x8a,
	0xb5, 0x18, 0x00, 0x42, 0x3d, 0x5a, 0x3b, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f,
	0x6d, 0x2f, 0x72, 0x6f, 0x6e, 0x61, 0x6b, 0x73, 0x6f, 0x66, 0x74, 0x2f, 0x72, 0x6f, 0x6e, 0x79,
	0x2f, 0x65, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x2f, 0x74, 0x61, 0x73, 0x6b, 0x5f, 0x6d, 0x61,
	0x6e, 0x61, 0x67, 0x65, 0x72, 0x2f, 0x6d, 0x6f, 0x64, 0x75, 0x6c, 0x65, 0x73, 0x2f, 0x74, 0x61,
	0x73, 0x6b, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_rpc_proto_rawDescOnce sync.Once
	file_rpc_proto_rawDescData = file_rpc_proto_rawDesc
)

func file_rpc_proto_rawDescGZIP() []byte {
	file_rpc_proto_rawDescOnce.Do(func() {
		file_rpc_proto_rawDescData = protoimpl.X.CompressGZIP(file_rpc_proto_rawDescData)
	})
	return file_rpc_proto_rawDescData
}

var file_rpc_proto_msgTypes = make([]protoimpl.MessageInfo, 7)
var file_rpc_proto_goTypes = []interface{}{
	(*CreateRequest)(nil), // 0: task.CreateRequest
	(*GetRequest)(nil),    // 1: task.GetRequest
	(*DeleteRequest)(nil), // 2: task.DeleteRequest
	(*ListRequest)(nil),   // 3: task.ListRequest
	(*TaskView)(nil),      // 4: task.TaskView
	(*TaskViewMany)(nil),  // 5: task.TaskViewMany
	(*Bool)(nil),          // 6: task.Bool
}
var file_rpc_proto_depIdxs = []int32{
	4, // 0: task.TaskViewMany.Tasks:type_name -> task.TaskView
	0, // 1: task.TaskManager.Create:input_type -> task.CreateRequest
	1, // 2: task.TaskManager.Get:input_type -> task.GetRequest
	2, // 3: task.TaskManager.Delete:input_type -> task.DeleteRequest
	3, // 4: task.TaskManager.List:input_type -> task.ListRequest
	4, // 5: task.TaskManager.Create:output_type -> task.TaskView
	4, // 6: task.TaskManager.Get:output_type -> task.TaskView
	6, // 7: task.TaskManager.Delete:output_type -> task.Bool
	5, // 8: task.TaskManager.List:output_type -> task.TaskViewMany
	5, // [5:9] is the sub-list for method output_type
	1, // [1:5] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_rpc_proto_init() }
func file_rpc_proto_init() {
	if File_rpc_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_rpc_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CreateRequest); i {
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
		file_rpc_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetRequest); i {
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
		file_rpc_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DeleteRequest); i {
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
		file_rpc_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ListRequest); i {
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
		file_rpc_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TaskView); i {
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
		file_rpc_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TaskViewMany); i {
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
		file_rpc_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Bool); i {
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
			RawDescriptor: file_rpc_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   7,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_rpc_proto_goTypes,
		DependencyIndexes: file_rpc_proto_depIdxs,
		MessageInfos:      file_rpc_proto_msgTypes,
	}.Build()
	File_rpc_proto = out.File
	file_rpc_proto_rawDesc = nil
	file_rpc_proto_goTypes = nil
	file_rpc_proto_depIdxs = nil
}
