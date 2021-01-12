// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0
// 	protoc        v3.14.0
// source: model.proto

package model

import (
	proto "github.com/golang/protobuf/proto"
	_ "github.com/ronaksoft/rony"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// This is a compile-time assertion that a sufficiently up-to-date version
// of the legacy proto package is being used.
const _ = proto.ProtoPackageIsVersion4

type Hook struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ClientID  string `protobuf:"bytes,1,opt,name=ClientID,proto3" json:"ClientID,omitempty"`
	ID        string `protobuf:"bytes,2,opt,name=ID,proto3" json:"ID,omitempty"`
	Timestamp string `protobuf:"bytes,3,opt,name=Timestamp,proto3" json:"Timestamp,omitempty"` // UTC unix epoch time
	HookUrl   string `protobuf:"bytes,4,opt,name=HookUrl,proto3" json:"HookUrl,omitempty"`
	Fired     bool   `protobuf:"varint,5,opt,name=Fired,proto3" json:"Fired,omitempty"`
	Success   bool   `protobuf:"varint,6,opt,name=Success,proto3" json:"Success,omitempty"`
}

func (x *Hook) Reset() {
	*x = Hook{}
	if protoimpl.UnsafeEnabled {
		mi := &file_model_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Hook) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Hook) ProtoMessage() {}

func (x *Hook) ProtoReflect() protoreflect.Message {
	mi := &file_model_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Hook.ProtoReflect.Descriptor instead.
func (*Hook) Descriptor() ([]byte, []int) {
	return file_model_proto_rawDescGZIP(), []int{0}
}

func (x *Hook) GetClientID() string {
	if x != nil {
		return x.ClientID
	}
	return ""
}

func (x *Hook) GetID() string {
	if x != nil {
		return x.ID
	}
	return ""
}

func (x *Hook) GetTimestamp() string {
	if x != nil {
		return x.Timestamp
	}
	return ""
}

func (x *Hook) GetHookUrl() string {
	if x != nil {
		return x.HookUrl
	}
	return ""
}

func (x *Hook) GetFired() bool {
	if x != nil {
		return x.Fired
	}
	return false
}

func (x *Hook) GetSuccess() bool {
	if x != nil {
		return x.Success
	}
	return false
}

type Model1 struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ID       int32    `protobuf:"varint,1,opt,name=ID,proto3" json:"ID,omitempty"`
	ShardKey int32    `protobuf:"varint,2,opt,name=ShardKey,proto3" json:"ShardKey,omitempty"`
	P1       string   `protobuf:"bytes,3,opt,name=P1,proto3" json:"P1,omitempty"`
	P2       []string `protobuf:"bytes,4,rep,name=P2,proto3" json:"P2,omitempty"`
	P5       uint64   `protobuf:"varint,5,opt,name=P5,proto3" json:"P5,omitempty"`
}

func (x *Model1) Reset() {
	*x = Model1{}
	if protoimpl.UnsafeEnabled {
		mi := &file_model_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Model1) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Model1) ProtoMessage() {}

func (x *Model1) ProtoReflect() protoreflect.Message {
	mi := &file_model_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Model1.ProtoReflect.Descriptor instead.
func (*Model1) Descriptor() ([]byte, []int) {
	return file_model_proto_rawDescGZIP(), []int{1}
}

func (x *Model1) GetID() int32 {
	if x != nil {
		return x.ID
	}
	return 0
}

func (x *Model1) GetShardKey() int32 {
	if x != nil {
		return x.ShardKey
	}
	return 0
}

func (x *Model1) GetP1() string {
	if x != nil {
		return x.P1
	}
	return ""
}

func (x *Model1) GetP2() []string {
	if x != nil {
		return x.P2
	}
	return nil
}

func (x *Model1) GetP5() uint64 {
	if x != nil {
		return x.P5
	}
	return 0
}

type Model2 struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ID       int64    `protobuf:"varint,1,opt,name=ID,proto3" json:"ID,omitempty"`
	ShardKey int32    `protobuf:"varint,2,opt,name=ShardKey,proto3" json:"ShardKey,omitempty"`
	P1       string   `protobuf:"bytes,3,opt,name=P1,proto3" json:"P1,omitempty"`
	P2       []string `protobuf:"bytes,4,rep,name=P2,proto3" json:"P2,omitempty"`
	P5       uint64   `protobuf:"varint,5,opt,name=P5,proto3" json:"P5,omitempty"`
}

func (x *Model2) Reset() {
	*x = Model2{}
	if protoimpl.UnsafeEnabled {
		mi := &file_model_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Model2) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Model2) ProtoMessage() {}

func (x *Model2) ProtoReflect() protoreflect.Message {
	mi := &file_model_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Model2.ProtoReflect.Descriptor instead.
func (*Model2) Descriptor() ([]byte, []int) {
	return file_model_proto_rawDescGZIP(), []int{2}
}

func (x *Model2) GetID() int64 {
	if x != nil {
		return x.ID
	}
	return 0
}

func (x *Model2) GetShardKey() int32 {
	if x != nil {
		return x.ShardKey
	}
	return 0
}

func (x *Model2) GetP1() string {
	if x != nil {
		return x.P1
	}
	return ""
}

func (x *Model2) GetP2() []string {
	if x != nil {
		return x.P2
	}
	return nil
}

func (x *Model2) GetP5() uint64 {
	if x != nil {
		return x.P5
	}
	return 0
}

var File_model_proto protoreflect.FileDescriptor

var file_model_proto_rawDesc = []byte{
	0x0a, 0x0b, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x05, 0x6d,
	0x6f, 0x64, 0x65, 0x6c, 0x1a, 0x0d, 0x6f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x22, 0x9a, 0x01, 0x0a, 0x04, 0x48, 0x6f, 0x6f, 0x6b, 0x12, 0x1a, 0x0a, 0x08,
	0x43, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x49, 0x44, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08,
	0x43, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x49, 0x44, 0x12, 0x0e, 0x0a, 0x02, 0x49, 0x44, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x49, 0x44, 0x12, 0x1c, 0x0a, 0x09, 0x54, 0x69, 0x6d, 0x65,
	0x73, 0x74, 0x61, 0x6d, 0x70, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x54, 0x69, 0x6d,
	0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x12, 0x18, 0x0a, 0x07, 0x48, 0x6f, 0x6f, 0x6b, 0x55, 0x72,
	0x6c, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x48, 0x6f, 0x6f, 0x6b, 0x55, 0x72, 0x6c,
	0x12, 0x14, 0x0a, 0x05, 0x46, 0x69, 0x72, 0x65, 0x64, 0x18, 0x05, 0x20, 0x01, 0x28, 0x08, 0x52,
	0x05, 0x46, 0x69, 0x72, 0x65, 0x64, 0x12, 0x18, 0x0a, 0x07, 0x53, 0x75, 0x63, 0x63, 0x65, 0x73,
	0x73, 0x18, 0x06, 0x20, 0x01, 0x28, 0x08, 0x52, 0x07, 0x53, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73,
	0x22, 0x8a, 0x01, 0x0a, 0x06, 0x4d, 0x6f, 0x64, 0x65, 0x6c, 0x31, 0x12, 0x0e, 0x0a, 0x02, 0x49,
	0x44, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x02, 0x49, 0x44, 0x12, 0x1a, 0x0a, 0x08, 0x53,
	0x68, 0x61, 0x72, 0x64, 0x4b, 0x65, 0x79, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x08, 0x53,
	0x68, 0x61, 0x72, 0x64, 0x4b, 0x65, 0x79, 0x12, 0x0e, 0x0a, 0x02, 0x50, 0x31, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x02, 0x50, 0x31, 0x12, 0x0e, 0x0a, 0x02, 0x50, 0x32, 0x18, 0x04, 0x20,
	0x03, 0x28, 0x09, 0x52, 0x02, 0x50, 0x32, 0x12, 0x0e, 0x0a, 0x02, 0x50, 0x35, 0x18, 0x05, 0x20,
	0x01, 0x28, 0x04, 0x52, 0x02, 0x50, 0x35, 0x3a, 0x24, 0x92, 0xb5, 0x18, 0x0e, 0x28, 0x49, 0x44,
	0x2c, 0x20, 0x53, 0x68, 0x61, 0x72, 0x64, 0x4b, 0x65, 0x79, 0x29, 0x9a, 0xb5, 0x18, 0x0e, 0x28,
	0x53, 0x68, 0x61, 0x72, 0x64, 0x4b, 0x65, 0x79, 0x2c, 0x20, 0x49, 0x44, 0x29, 0x22, 0x95, 0x01,
	0x0a, 0x06, 0x4d, 0x6f, 0x64, 0x65, 0x6c, 0x32, 0x12, 0x0e, 0x0a, 0x02, 0x49, 0x44, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x03, 0x52, 0x02, 0x49, 0x44, 0x12, 0x1a, 0x0a, 0x08, 0x53, 0x68, 0x61, 0x72,
	0x64, 0x4b, 0x65, 0x79, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x08, 0x53, 0x68, 0x61, 0x72,
	0x64, 0x4b, 0x65, 0x79, 0x12, 0x0e, 0x0a, 0x02, 0x50, 0x31, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x02, 0x50, 0x31, 0x12, 0x0e, 0x0a, 0x02, 0x50, 0x32, 0x18, 0x04, 0x20, 0x03, 0x28, 0x09,
	0x52, 0x02, 0x50, 0x32, 0x12, 0x0e, 0x0a, 0x02, 0x50, 0x35, 0x18, 0x05, 0x20, 0x01, 0x28, 0x04,
	0x52, 0x02, 0x50, 0x35, 0x3a, 0x2f, 0x92, 0xb5, 0x18, 0x15, 0x28, 0x28, 0x49, 0x44, 0x2c, 0x20,
	0x53, 0x68, 0x61, 0x72, 0x64, 0x4b, 0x65, 0x79, 0x29, 0x2c, 0x20, 0x2d, 0x50, 0x31, 0x29, 0x9a,
	0xb5, 0x18, 0x12, 0x28, 0x50, 0x31, 0x2c, 0x20, 0x53, 0x68, 0x61, 0x72, 0x64, 0x4b, 0x65, 0x79,
	0x2c, 0x20, 0x49, 0x44, 0x29, 0x42, 0x13, 0x5a, 0x07, 0x2e, 0x3b, 0x6d, 0x6f, 0x64, 0x65, 0x6c,
	0x8a, 0x92, 0xf4, 0x01, 0x05, 0x6c, 0x6f, 0x63, 0x61, 0x6c, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x33,
}

var (
	file_model_proto_rawDescOnce sync.Once
	file_model_proto_rawDescData = file_model_proto_rawDesc
)

func file_model_proto_rawDescGZIP() []byte {
	file_model_proto_rawDescOnce.Do(func() {
		file_model_proto_rawDescData = protoimpl.X.CompressGZIP(file_model_proto_rawDescData)
	})
	return file_model_proto_rawDescData
}

var file_model_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_model_proto_goTypes = []interface{}{
	(*Hook)(nil),   // 0: model.Hook
	(*Model1)(nil), // 1: model.Model1
	(*Model2)(nil), // 2: model.Model2
}
var file_model_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_model_proto_init() }
func file_model_proto_init() {
	if File_model_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_model_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Hook); i {
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
		file_model_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Model1); i {
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
		file_model_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Model2); i {
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
			RawDescriptor: file_model_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_model_proto_goTypes,
		DependencyIndexes: file_model_proto_depIdxs,
		MessageInfos:      file_model_proto_msgTypes,
	}.Build()
	File_model_proto = out.File
	file_model_proto_rawDesc = nil
	file_model_proto_goTypes = nil
	file_model_proto_depIdxs = nil
}
