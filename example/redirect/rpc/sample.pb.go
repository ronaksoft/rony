// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        v3.17.3
// source: sample.proto

package rpc

import (
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

// InfoRequest
type InfoRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ReplicaSet uint64 `protobuf:"varint,1,opt,name=ReplicaSet,proto3" json:"ReplicaSet,omitempty"`
	RandomText string `protobuf:"bytes,2,opt,name=RandomText,proto3" json:"RandomText,omitempty"`
}

func (x *InfoRequest) Reset() {
	*x = InfoRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_sample_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *InfoRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*InfoRequest) ProtoMessage() {}

func (x *InfoRequest) ProtoReflect() protoreflect.Message {
	mi := &file_sample_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use InfoRequest.ProtoReflect.Descriptor instead.
func (*InfoRequest) Descriptor() ([]byte, []int) {
	return file_sample_proto_rawDescGZIP(), []int{0}
}

func (x *InfoRequest) GetReplicaSet() uint64 {
	if x != nil {
		return x.ReplicaSet
	}
	return 0
}

func (x *InfoRequest) GetRandomText() string {
	if x != nil {
		return x.RandomText
	}
	return ""
}

// InfoResponse
type InfoResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ServerID   string `protobuf:"bytes,1,opt,name=ServerID,proto3" json:"ServerID,omitempty"`
	RandomText string `protobuf:"bytes,2,opt,name=RandomText,proto3" json:"RandomText,omitempty"`
}

func (x *InfoResponse) Reset() {
	*x = InfoResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_sample_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *InfoResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*InfoResponse) ProtoMessage() {}

func (x *InfoResponse) ProtoReflect() protoreflect.Message {
	mi := &file_sample_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use InfoResponse.ProtoReflect.Descriptor instead.
func (*InfoResponse) Descriptor() ([]byte, []int) {
	return file_sample_proto_rawDescGZIP(), []int{1}
}

func (x *InfoResponse) GetServerID() string {
	if x != nil {
		return x.ServerID
	}
	return ""
}

func (x *InfoResponse) GetRandomText() string {
	if x != nil {
		return x.RandomText
	}
	return ""
}

var File_sample_proto protoreflect.FileDescriptor

var file_sample_proto_rawDesc = []byte{
	0x0a, 0x0c, 0x73, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x03,
	0x72, 0x70, 0x63, 0x1a, 0x0d, 0x6f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x22, 0x4d, 0x0a, 0x0b, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x12, 0x1e, 0x0a, 0x0a, 0x52, 0x65, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x53, 0x65, 0x74, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x0a, 0x52, 0x65, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x53, 0x65,
	0x74, 0x12, 0x1e, 0x0a, 0x0a, 0x52, 0x61, 0x6e, 0x64, 0x6f, 0x6d, 0x54, 0x65, 0x78, 0x74, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x52, 0x61, 0x6e, 0x64, 0x6f, 0x6d, 0x54, 0x65, 0x78,
	0x74, 0x22, 0x4a, 0x0a, 0x0c, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x12, 0x1a, 0x0a, 0x08, 0x53, 0x65, 0x72, 0x76, 0x65, 0x72, 0x49, 0x44, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x08, 0x53, 0x65, 0x72, 0x76, 0x65, 0x72, 0x49, 0x44, 0x12, 0x1e, 0x0a,
	0x0a, 0x52, 0x61, 0x6e, 0x64, 0x6f, 0x6d, 0x54, 0x65, 0x78, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x0a, 0x52, 0x61, 0x6e, 0x64, 0x6f, 0x6d, 0x54, 0x65, 0x78, 0x74, 0x32, 0xfe, 0x01,
	0x0a, 0x06, 0x53, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x12, 0x79, 0x0a, 0x16, 0x49, 0x6e, 0x66, 0x6f,
	0x57, 0x69, 0x74, 0x68, 0x43, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x64, 0x69, 0x72, 0x65,
	0x63, 0x74, 0x12, 0x10, 0x2e, 0x72, 0x70, 0x63, 0x2e, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x1a, 0x11, 0x2e, 0x72, 0x70, 0x63, 0x2e, 0x49, 0x6e, 0x66, 0x6f, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x3a, 0x9a, 0xb5, 0x18, 0x36, 0x0a, 0x03, 0x67,
	0x65, 0x74, 0x12, 0x2d, 0x2f, 0x69, 0x6e, 0x66, 0x6f, 0x2f, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74,
	0x2d, 0x72, 0x65, 0x64, 0x69, 0x72, 0x65, 0x63, 0x74, 0x2f, 0x3a, 0x52, 0x65, 0x70, 0x6c, 0x69,
	0x63, 0x61, 0x53, 0x65, 0x74, 0x2f, 0x3a, 0x52, 0x61, 0x6e, 0x64, 0x6f, 0x6d, 0x54, 0x65, 0x78,
	0x74, 0x28, 0x01, 0x12, 0x79, 0x0a, 0x16, 0x49, 0x6e, 0x66, 0x6f, 0x57, 0x69, 0x74, 0x68, 0x53,
	0x65, 0x72, 0x76, 0x65, 0x72, 0x52, 0x65, 0x64, 0x69, 0x72, 0x65, 0x63, 0x74, 0x12, 0x10, 0x2e,
	0x72, 0x70, 0x63, 0x2e, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a,
	0x11, 0x2e, 0x72, 0x70, 0x63, 0x2e, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x22, 0x3a, 0x9a, 0xb5, 0x18, 0x36, 0x0a, 0x03, 0x67, 0x65, 0x74, 0x12, 0x2d, 0x2f,
	0x69, 0x6e, 0x66, 0x6f, 0x2f, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2d, 0x72, 0x65, 0x64, 0x69,
	0x72, 0x65, 0x63, 0x74, 0x2f, 0x3a, 0x52, 0x65, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x53, 0x65, 0x74,
	0x2f, 0x3a, 0x52, 0x61, 0x6e, 0x64, 0x6f, 0x6d, 0x54, 0x65, 0x78, 0x74, 0x28, 0x01, 0x42, 0x30,
	0x5a, 0x2e, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x72, 0x6f, 0x6e,
	0x61, 0x6b, 0x73, 0x6f, 0x66, 0x74, 0x2f, 0x72, 0x6f, 0x6e, 0x79, 0x2f, 0x65, 0x78, 0x61, 0x6d,
	0x70, 0x6c, 0x65, 0x2f, 0x72, 0x65, 0x64, 0x69, 0x72, 0x65, 0x63, 0x74, 0x2f, 0x72, 0x70, 0x63,
	0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_sample_proto_rawDescOnce sync.Once
	file_sample_proto_rawDescData = file_sample_proto_rawDesc
)

func file_sample_proto_rawDescGZIP() []byte {
	file_sample_proto_rawDescOnce.Do(func() {
		file_sample_proto_rawDescData = protoimpl.X.CompressGZIP(file_sample_proto_rawDescData)
	})
	return file_sample_proto_rawDescData
}

var file_sample_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_sample_proto_goTypes = []interface{}{
	(*InfoRequest)(nil),  // 0: rpc.InfoRequest
	(*InfoResponse)(nil), // 1: rpc.InfoResponse
}
var file_sample_proto_depIdxs = []int32{
	0, // 0: rpc.Sample.InfoWithClientRedirect:input_type -> rpc.InfoRequest
	0, // 1: rpc.Sample.InfoWithServerRedirect:input_type -> rpc.InfoRequest
	1, // 2: rpc.Sample.InfoWithClientRedirect:output_type -> rpc.InfoResponse
	1, // 3: rpc.Sample.InfoWithServerRedirect:output_type -> rpc.InfoResponse
	2, // [2:4] is the sub-list for method output_type
	0, // [0:2] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_sample_proto_init() }
func file_sample_proto_init() {
	if File_sample_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_sample_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*InfoRequest); i {
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
		file_sample_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*InfoResponse); i {
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
			RawDescriptor: file_sample_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_sample_proto_goTypes,
		DependencyIndexes: file_sample_proto_depIdxs,
		MessageInfos:      file_sample_proto_msgTypes,
	}.Build()
	File_sample_proto = out.File
	file_sample_proto_rawDesc = nil
	file_sample_proto_goTypes = nil
	file_sample_proto_depIdxs = nil
}
