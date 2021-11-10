// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        v3.17.3
// source: model.proto

package model

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

type Enum int32

const (
	Enum_None      Enum = 0
	Enum_Something Enum = 1
	Enum_Else      Enum = 2
)

// Enum value maps for Enum.
var (
	Enum_name = map[int32]string{
		0: "None",
		1: "Something",
		2: "Else",
	}
	Enum_value = map[string]int32{
		"None":      0,
		"Something": 1,
		"Else":      2,
	}
)

func (x Enum) Enum() *Enum {
	p := new(Enum)
	*p = x
	return p
}

func (x Enum) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Enum) Descriptor() protoreflect.EnumDescriptor {
	return file_model_proto_enumTypes[0].Descriptor()
}

func (Enum) Type() protoreflect.EnumType {
	return &file_model_proto_enumTypes[0]
}

func (x Enum) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Enum.Descriptor instead.
func (Enum) EnumDescriptor() ([]byte, []int) {
	return file_model_proto_rawDescGZIP(), []int{0}
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
	Enum     Enum     `protobuf:"varint,7,opt,name=Enum,proto3,enum=model.Enum" json:"Enum,omitempty"`
}

func (x *Model1) Reset() {
	*x = Model1{}
	if protoimpl.UnsafeEnabled {
		mi := &file_model_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Model1) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Model1) ProtoMessage() {}

func (x *Model1) ProtoReflect() protoreflect.Message {
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

// Deprecated: Use Model1.ProtoReflect.Descriptor instead.
func (*Model1) Descriptor() ([]byte, []int) {
	return file_model_proto_rawDescGZIP(), []int{0}
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

func (x *Model1) GetEnum() Enum {
	if x != nil {
		return x.Enum
	}
	return Enum_None
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
		mi := &file_model_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Model2) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Model2) ProtoMessage() {}

func (x *Model2) ProtoReflect() protoreflect.Message {
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

// Deprecated: Use Model2.ProtoReflect.Descriptor instead.
func (*Model2) Descriptor() ([]byte, []int) {
	return file_model_proto_rawDescGZIP(), []int{1}
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

type Model3 struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ID       int64    `protobuf:"varint,1,opt,name=ID,proto3" json:"ID,omitempty"`
	ShardKey int32    `protobuf:"varint,2,opt,name=ShardKey,proto3" json:"ShardKey,omitempty"`
	P1       []byte   `protobuf:"bytes,3,opt,name=P1,proto3" json:"P1,omitempty"`
	P2       []string `protobuf:"bytes,4,rep,name=P2,proto3" json:"P2,omitempty"`
	P5       [][]byte `protobuf:"bytes,5,rep,name=P5,proto3" json:"P5,omitempty"`
}

func (x *Model3) Reset() {
	*x = Model3{}
	if protoimpl.UnsafeEnabled {
		mi := &file_model_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Model3) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Model3) ProtoMessage() {}

func (x *Model3) ProtoReflect() protoreflect.Message {
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

// Deprecated: Use Model3.ProtoReflect.Descriptor instead.
func (*Model3) Descriptor() ([]byte, []int) {
	return file_model_proto_rawDescGZIP(), []int{2}
}

func (x *Model3) GetID() int64 {
	if x != nil {
		return x.ID
	}
	return 0
}

func (x *Model3) GetShardKey() int32 {
	if x != nil {
		return x.ShardKey
	}
	return 0
}

func (x *Model3) GetP1() []byte {
	if x != nil {
		return x.P1
	}
	return nil
}

func (x *Model3) GetP2() []string {
	if x != nil {
		return x.P2
	}
	return nil
}

func (x *Model3) GetP5() [][]byte {
	if x != nil {
		return x.P5
	}
	return nil
}

var File_model_proto protoreflect.FileDescriptor

var file_model_proto_rawDesc = []byte{
	0x0a, 0x0b, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x05, 0x6d,
	0x6f, 0x64, 0x65, 0x6c, 0x1a, 0x0d, 0x6f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x22, 0xdd, 0x01, 0x0a, 0x06, 0x4d, 0x6f, 0x64, 0x65, 0x6c, 0x31, 0x12, 0x0e,
	0x0a, 0x02, 0x49, 0x44, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x02, 0x49, 0x44, 0x12, 0x1a,
	0x0a, 0x08, 0x53, 0x68, 0x61, 0x72, 0x64, 0x4b, 0x65, 0x79, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05,
	0x52, 0x08, 0x53, 0x68, 0x61, 0x72, 0x64, 0x4b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x02, 0x50, 0x31,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x42, 0x04, 0x88, 0xb5, 0x18, 0x01, 0x52, 0x02, 0x50, 0x31,
	0x12, 0x14, 0x0a, 0x02, 0x50, 0x32, 0x18, 0x04, 0x20, 0x03, 0x28, 0x09, 0x42, 0x04, 0x88, 0xb5,
	0x18, 0x01, 0x52, 0x02, 0x50, 0x32, 0x12, 0x0e, 0x0a, 0x02, 0x50, 0x35, 0x18, 0x05, 0x20, 0x01,
	0x28, 0x04, 0x52, 0x02, 0x50, 0x35, 0x12, 0x1f, 0x0a, 0x04, 0x45, 0x6e, 0x75, 0x6d, 0x18, 0x07,
	0x20, 0x01, 0x28, 0x0e, 0x32, 0x0b, 0x2e, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x2e, 0x45, 0x6e, 0x75,
	0x6d, 0x52, 0x04, 0x45, 0x6e, 0x75, 0x6d, 0x3a, 0x4a, 0x8a, 0xb5, 0x18, 0x46, 0x12, 0x03, 0x63,
	0x71, 0x6c, 0x1a, 0x05, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x22, 0x14, 0x0a, 0x02, 0x49, 0x44, 0x12,
	0x08, 0x53, 0x68, 0x61, 0x72, 0x64, 0x4b, 0x65, 0x79, 0x12, 0x04, 0x45, 0x6e, 0x75, 0x6d, 0x2a,
	0x22, 0x0a, 0x04, 0x45, 0x6e, 0x75, 0x6d, 0x12, 0x08, 0x53, 0x68, 0x61, 0x72, 0x64, 0x4b, 0x65,
	0x79, 0x12, 0x02, 0x49, 0x44, 0x22, 0x0c, 0x43, 0x75, 0x73, 0x74, 0x6f, 0x6d, 0x65, 0x72, 0x53,
	0x6f, 0x72, 0x74, 0x22, 0x9f, 0x01, 0x0a, 0x06, 0x4d, 0x6f, 0x64, 0x65, 0x6c, 0x32, 0x12, 0x0e,
	0x0a, 0x02, 0x49, 0x44, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x02, 0x49, 0x44, 0x12, 0x1a,
	0x0a, 0x08, 0x53, 0x68, 0x61, 0x72, 0x64, 0x4b, 0x65, 0x79, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05,
	0x52, 0x08, 0x53, 0x68, 0x61, 0x72, 0x64, 0x4b, 0x65, 0x79, 0x12, 0x0e, 0x0a, 0x02, 0x50, 0x31,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x50, 0x31, 0x12, 0x0e, 0x0a, 0x02, 0x50, 0x32,
	0x18, 0x04, 0x20, 0x03, 0x28, 0x09, 0x52, 0x02, 0x50, 0x32, 0x12, 0x0e, 0x0a, 0x02, 0x50, 0x35,
	0x18, 0x05, 0x20, 0x01, 0x28, 0x04, 0x52, 0x02, 0x50, 0x35, 0x3a, 0x39, 0x8a, 0xb5, 0x18, 0x35,
	0x12, 0x03, 0x63, 0x71, 0x6c, 0x1a, 0x05, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x22, 0x13, 0x0a, 0x02,
	0x49, 0x44, 0x0a, 0x08, 0x53, 0x68, 0x61, 0x72, 0x64, 0x4b, 0x65, 0x79, 0x12, 0x03, 0x2d, 0x50,
	0x31, 0x2a, 0x12, 0x0a, 0x02, 0x50, 0x31, 0x12, 0x08, 0x53, 0x68, 0x61, 0x72, 0x64, 0x4b, 0x65,
	0x79, 0x12, 0x02, 0x49, 0x44, 0x22, 0xb9, 0x01, 0x0a, 0x06, 0x4d, 0x6f, 0x64, 0x65, 0x6c, 0x33,
	0x12, 0x0e, 0x0a, 0x02, 0x49, 0x44, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x02, 0x49, 0x44,
	0x12, 0x1a, 0x0a, 0x08, 0x53, 0x68, 0x61, 0x72, 0x64, 0x4b, 0x65, 0x79, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x05, 0x52, 0x08, 0x53, 0x68, 0x61, 0x72, 0x64, 0x4b, 0x65, 0x79, 0x12, 0x0e, 0x0a, 0x02,
	0x50, 0x31, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x02, 0x50, 0x31, 0x12, 0x0e, 0x0a, 0x02,
	0x50, 0x32, 0x18, 0x04, 0x20, 0x03, 0x28, 0x09, 0x52, 0x02, 0x50, 0x32, 0x12, 0x14, 0x0a, 0x02,
	0x50, 0x35, 0x18, 0x05, 0x20, 0x03, 0x28, 0x0c, 0x42, 0x04, 0x88, 0xb5, 0x18, 0x01, 0x52, 0x02,
	0x50, 0x35, 0x3a, 0x4d, 0x8a, 0xb5, 0x18, 0x49, 0x12, 0x03, 0x63, 0x71, 0x6c, 0x1a, 0x05, 0x73,
	0x74, 0x6f, 0x72, 0x65, 0x22, 0x13, 0x0a, 0x02, 0x49, 0x44, 0x0a, 0x08, 0x53, 0x68, 0x61, 0x72,
	0x64, 0x4b, 0x65, 0x79, 0x12, 0x03, 0x2d, 0x50, 0x31, 0x2a, 0x12, 0x0a, 0x02, 0x50, 0x31, 0x12,
	0x08, 0x53, 0x68, 0x61, 0x72, 0x64, 0x4b, 0x65, 0x79, 0x12, 0x02, 0x49, 0x44, 0x2a, 0x12, 0x0a,
	0x02, 0x50, 0x31, 0x12, 0x02, 0x49, 0x44, 0x12, 0x08, 0x53, 0x68, 0x61, 0x72, 0x64, 0x4b, 0x65,
	0x79, 0x2a, 0x29, 0x0a, 0x04, 0x45, 0x6e, 0x75, 0x6d, 0x12, 0x08, 0x0a, 0x04, 0x4e, 0x6f, 0x6e,
	0x65, 0x10, 0x00, 0x12, 0x0d, 0x0a, 0x09, 0x53, 0x6f, 0x6d, 0x65, 0x74, 0x68, 0x69, 0x6e, 0x67,
	0x10, 0x01, 0x12, 0x08, 0x0a, 0x04, 0x45, 0x6c, 0x73, 0x65, 0x10, 0x02, 0x42, 0x0a, 0x5a, 0x08,
	0x2e, 0x2f, 0x3b, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
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

var file_model_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_model_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_model_proto_goTypes = []interface{}{
	(Enum)(0),      // 0: model.Enum
	(*Model1)(nil), // 1: model.Model1
	(*Model2)(nil), // 2: model.Model2
	(*Model3)(nil), // 3: model.Model3
}
var file_model_proto_depIdxs = []int32{
	0, // 0: model.Model1.Enum:type_name -> model.Enum
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_model_proto_init() }
func file_model_proto_init() {
	if File_model_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_model_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
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
		file_model_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
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
		file_model_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Model3); i {
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
			NumEnums:      1,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_model_proto_goTypes,
		DependencyIndexes: file_model_proto_depIdxs,
		EnumInfos:         file_model_proto_enumTypes,
		MessageInfos:      file_model_proto_msgTypes,
	}.Build()
	File_model_proto = out.File
	file_model_proto_rawDesc = nil
	file_model_proto_goTypes = nil
	file_model_proto_depIdxs = nil
}
