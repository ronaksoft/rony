// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0
// 	protoc        v3.9.2
// source: test.proto

package packageName

import (
	proto "github.com/golang/protobuf/proto"
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

type Message1 struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Param1 int64  `protobuf:"varint,1,opt,name=Param1,proto3" json:"Param1,omitempty"`
	Param2 string `protobuf:"bytes,2,opt,name=Param2,proto3" json:"Param2,omitempty"`
}

func (x *Message1) Reset() {
	*x = Message1{}
	if protoimpl.UnsafeEnabled {
		mi := &file_test_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Message1) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Message1) ProtoMessage() {}

func (x *Message1) ProtoReflect() protoreflect.Message {
	mi := &file_test_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Message1.ProtoReflect.Descriptor instead.
func (*Message1) Descriptor() ([]byte, []int) {
	return file_test_proto_rawDescGZIP(), []int{0}
}

func (x *Message1) GetParam1() int64 {
	if x != nil {
		return x.Param1
	}
	return 0
}

func (x *Message1) GetParam2() string {
	if x != nil {
		return x.Param2
	}
	return ""
}

type Message2 struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Param1 uint32  `protobuf:"fixed32,1,opt,name=Param1,proto3" json:"Param1,omitempty"`
	P2     []byte  `protobuf:"bytes,2,opt,name=P2,proto3" json:"P2,omitempty"`
	P3     []int32 `protobuf:"varint,3,rep,packed,name=P3,proto3" json:"P3,omitempty"`
}

func (x *Message2) Reset() {
	*x = Message2{}
	if protoimpl.UnsafeEnabled {
		mi := &file_test_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Message2) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Message2) ProtoMessage() {}

func (x *Message2) ProtoReflect() protoreflect.Message {
	mi := &file_test_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Message2.ProtoReflect.Descriptor instead.
func (*Message2) Descriptor() ([]byte, []int) {
	return file_test_proto_rawDescGZIP(), []int{1}
}

func (x *Message2) GetParam1() uint32 {
	if x != nil {
		return x.Param1
	}
	return 0
}

func (x *Message2) GetP2() []byte {
	if x != nil {
		return x.P2
	}
	return nil
}

func (x *Message2) GetP3() []int32 {
	if x != nil {
		return x.P3
	}
	return nil
}

// {{@model cql}}
// {{@tab (ID, ShardKey)}}
// {{@view (ShardKey, ID)}}
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
		mi := &file_test_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Model1) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Model1) ProtoMessage() {}

func (x *Model1) ProtoReflect() protoreflect.Message {
	mi := &file_test_proto_msgTypes[2]
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
	return file_test_proto_rawDescGZIP(), []int{2}
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

// {{@model cql}}
// {{@tab ((ID, ShardKey), -P1)}}
// {{@view (P1, ShardKey, ID)}}
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
		mi := &file_test_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Model2) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Model2) ProtoMessage() {}

func (x *Model2) ProtoReflect() protoreflect.Message {
	mi := &file_test_proto_msgTypes[3]
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
	return file_test_proto_rawDescGZIP(), []int{3}
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

var File_test_proto protoreflect.FileDescriptor

var file_test_proto_rawDesc = []byte{
	0x0a, 0x0a, 0x74, 0x65, 0x73, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0b, 0x70, 0x61,
	0x63, 0x6b, 0x61, 0x67, 0x65, 0x4e, 0x61, 0x6d, 0x65, 0x22, 0x3a, 0x0a, 0x08, 0x4d, 0x65, 0x73,
	0x73, 0x61, 0x67, 0x65, 0x31, 0x12, 0x16, 0x0a, 0x06, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x31, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x06, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x31, 0x12, 0x16, 0x0a,
	0x06, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x32, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x50,
	0x61, 0x72, 0x61, 0x6d, 0x32, 0x22, 0x42, 0x0a, 0x08, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65,
	0x32, 0x12, 0x16, 0x0a, 0x06, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x31, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x07, 0x52, 0x06, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x31, 0x12, 0x0e, 0x0a, 0x02, 0x50, 0x32, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x02, 0x50, 0x32, 0x12, 0x0e, 0x0a, 0x02, 0x50, 0x33, 0x18,
	0x03, 0x20, 0x03, 0x28, 0x05, 0x52, 0x02, 0x50, 0x33, 0x22, 0x64, 0x0a, 0x06, 0x4d, 0x6f, 0x64,
	0x65, 0x6c, 0x31, 0x12, 0x0e, 0x0a, 0x02, 0x49, 0x44, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52,
	0x02, 0x49, 0x44, 0x12, 0x1a, 0x0a, 0x08, 0x53, 0x68, 0x61, 0x72, 0x64, 0x4b, 0x65, 0x79, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x08, 0x53, 0x68, 0x61, 0x72, 0x64, 0x4b, 0x65, 0x79, 0x12,
	0x0e, 0x0a, 0x02, 0x50, 0x31, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x50, 0x31, 0x12,
	0x0e, 0x0a, 0x02, 0x50, 0x32, 0x18, 0x04, 0x20, 0x03, 0x28, 0x09, 0x52, 0x02, 0x50, 0x32, 0x12,
	0x0e, 0x0a, 0x02, 0x50, 0x35, 0x18, 0x05, 0x20, 0x01, 0x28, 0x04, 0x52, 0x02, 0x50, 0x35, 0x22,
	0x64, 0x0a, 0x06, 0x4d, 0x6f, 0x64, 0x65, 0x6c, 0x32, 0x12, 0x0e, 0x0a, 0x02, 0x49, 0x44, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x02, 0x49, 0x44, 0x12, 0x1a, 0x0a, 0x08, 0x53, 0x68, 0x61,
	0x72, 0x64, 0x4b, 0x65, 0x79, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x08, 0x53, 0x68, 0x61,
	0x72, 0x64, 0x4b, 0x65, 0x79, 0x12, 0x0e, 0x0a, 0x02, 0x50, 0x31, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x02, 0x50, 0x31, 0x12, 0x0e, 0x0a, 0x02, 0x50, 0x32, 0x18, 0x04, 0x20, 0x03, 0x28,
	0x09, 0x52, 0x02, 0x50, 0x32, 0x12, 0x0e, 0x0a, 0x02, 0x50, 0x35, 0x18, 0x05, 0x20, 0x01, 0x28,
	0x04, 0x52, 0x02, 0x50, 0x35, 0x42, 0x0f, 0x5a, 0x0d, 0x2e, 0x3b, 0x70, 0x61, 0x63, 0x6b, 0x61,
	0x67, 0x65, 0x4e, 0x61, 0x6d, 0x65, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_test_proto_rawDescOnce sync.Once
	file_test_proto_rawDescData = file_test_proto_rawDesc
)

func file_test_proto_rawDescGZIP() []byte {
	file_test_proto_rawDescOnce.Do(func() {
		file_test_proto_rawDescData = protoimpl.X.CompressGZIP(file_test_proto_rawDescData)
	})
	return file_test_proto_rawDescData
}

var file_test_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_test_proto_goTypes = []interface{}{
	(*Message1)(nil), // 0: packageName.Message1
	(*Message2)(nil), // 1: packageName.Message2
	(*Model1)(nil),   // 2: packageName.Model1
	(*Model2)(nil),   // 3: packageName.Model2
}
var file_test_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_test_proto_init() }
func file_test_proto_init() {
	if File_test_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_test_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Message1); i {
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
		file_test_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Message2); i {
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
		file_test_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
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
		file_test_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
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
			RawDescriptor: file_test_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_test_proto_goTypes,
		DependencyIndexes: file_test_proto_depIdxs,
		MessageInfos:      file_test_proto_msgTypes,
	}.Build()
	File_test_proto = out.File
	file_test_proto_rawDesc = nil
	file_test_proto_goTypes = nil
	file_test_proto_depIdxs = nil
}
