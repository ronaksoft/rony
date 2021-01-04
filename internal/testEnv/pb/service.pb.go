// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0
// 	protoc        v3.14.0
// source: service.proto

package pb

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

type EchoRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Int       int64 `protobuf:"varint,1,opt,name=Int,proto3" json:"Int,omitempty"`
	Timestamp int64 `protobuf:"varint,3,opt,name=Timestamp,proto3" json:"Timestamp,omitempty"`
}

func (x *EchoRequest) Reset() {
	*x = EchoRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_service_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *EchoRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EchoRequest) ProtoMessage() {}

func (x *EchoRequest) ProtoReflect() protoreflect.Message {
	mi := &file_service_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EchoRequest.ProtoReflect.Descriptor instead.
func (*EchoRequest) Descriptor() ([]byte, []int) {
	return file_service_proto_rawDescGZIP(), []int{0}
}

func (x *EchoRequest) GetInt() int64 {
	if x != nil {
		return x.Int
	}
	return 0
}

func (x *EchoRequest) GetTimestamp() int64 {
	if x != nil {
		return x.Timestamp
	}
	return 0
}

type EchoResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Int       int64  `protobuf:"varint,1,opt,name=Int,proto3" json:"Int,omitempty"`
	Responder string `protobuf:"bytes,2,opt,name=Responder,proto3" json:"Responder,omitempty"`
	Timestamp int64  `protobuf:"varint,4,opt,name=Timestamp,proto3" json:"Timestamp,omitempty"`
	Delay     int64  `protobuf:"varint,5,opt,name=Delay,proto3" json:"Delay,omitempty"`
	ServerID  string `protobuf:"bytes,6,opt,name=ServerID,proto3" json:"ServerID,omitempty"`
}

func (x *EchoResponse) Reset() {
	*x = EchoResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_service_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *EchoResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EchoResponse) ProtoMessage() {}

func (x *EchoResponse) ProtoReflect() protoreflect.Message {
	mi := &file_service_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EchoResponse.ProtoReflect.Descriptor instead.
func (*EchoResponse) Descriptor() ([]byte, []int) {
	return file_service_proto_rawDescGZIP(), []int{1}
}

func (x *EchoResponse) GetInt() int64 {
	if x != nil {
		return x.Int
	}
	return 0
}

func (x *EchoResponse) GetResponder() string {
	if x != nil {
		return x.Responder
	}
	return ""
}

func (x *EchoResponse) GetTimestamp() int64 {
	if x != nil {
		return x.Timestamp
	}
	return 0
}

func (x *EchoResponse) GetDelay() int64 {
	if x != nil {
		return x.Delay
	}
	return 0
}

func (x *EchoResponse) GetServerID() string {
	if x != nil {
		return x.ServerID
	}
	return ""
}

var File_service_proto protoreflect.FileDescriptor

var file_service_proto_rawDesc = []byte{
	0x0a, 0x0d, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12,
	0x02, 0x70, 0x62, 0x1a, 0x0d, 0x6f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x22, 0x3d, 0x0a, 0x0b, 0x45, 0x63, 0x68, 0x6f, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x12, 0x10, 0x0a, 0x03, 0x49, 0x6e, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x03,
	0x49, 0x6e, 0x74, 0x12, 0x1c, 0x0a, 0x09, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d,
	0x70, 0x22, 0x8e, 0x01, 0x0a, 0x0c, 0x45, 0x63, 0x68, 0x6f, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x12, 0x10, 0x0a, 0x03, 0x49, 0x6e, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52,
	0x03, 0x49, 0x6e, 0x74, 0x12, 0x1c, 0x0a, 0x09, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x64, 0x65,
	0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x64,
	0x65, 0x72, 0x12, 0x1c, 0x0a, 0x09, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x18,
	0x04, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70,
	0x12, 0x14, 0x0a, 0x05, 0x44, 0x65, 0x6c, 0x61, 0x79, 0x18, 0x05, 0x20, 0x01, 0x28, 0x03, 0x52,
	0x05, 0x44, 0x65, 0x6c, 0x61, 0x79, 0x12, 0x1a, 0x0a, 0x08, 0x53, 0x65, 0x72, 0x76, 0x65, 0x72,
	0x49, 0x44, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x53, 0x65, 0x72, 0x76, 0x65, 0x72,
	0x49, 0x44, 0x32, 0x9f, 0x01, 0x0a, 0x06, 0x53, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x12, 0x29, 0x0a,
	0x04, 0x45, 0x63, 0x68, 0x6f, 0x12, 0x0f, 0x2e, 0x70, 0x62, 0x2e, 0x45, 0x63, 0x68, 0x6f, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x10, 0x2e, 0x70, 0x62, 0x2e, 0x45, 0x63, 0x68, 0x6f,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x39, 0x0a, 0x0e, 0x45, 0x63, 0x68, 0x6f,
	0x4c, 0x65, 0x61, 0x64, 0x65, 0x72, 0x4f, 0x6e, 0x6c, 0x79, 0x12, 0x0f, 0x2e, 0x70, 0x62, 0x2e,
	0x45, 0x63, 0x68, 0x6f, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x10, 0x2e, 0x70, 0x62,
	0x2e, 0x45, 0x63, 0x68, 0x6f, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x04, 0x88,
	0xb5, 0x18, 0x01, 0x12, 0x2f, 0x0a, 0x0a, 0x45, 0x63, 0x68, 0x6f, 0x54, 0x75, 0x6e, 0x6e, 0x65,
	0x6c, 0x12, 0x0f, 0x2e, 0x70, 0x62, 0x2e, 0x45, 0x63, 0x68, 0x6f, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x1a, 0x10, 0x2e, 0x70, 0x62, 0x2e, 0x45, 0x63, 0x68, 0x6f, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x42, 0x06, 0x5a, 0x04, 0x2e, 0x3b, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_service_proto_rawDescOnce sync.Once
	file_service_proto_rawDescData = file_service_proto_rawDesc
)

func file_service_proto_rawDescGZIP() []byte {
	file_service_proto_rawDescOnce.Do(func() {
		file_service_proto_rawDescData = protoimpl.X.CompressGZIP(file_service_proto_rawDescData)
	})
	return file_service_proto_rawDescData
}

var file_service_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_service_proto_goTypes = []interface{}{
	(*EchoRequest)(nil),  // 0: pb.EchoRequest
	(*EchoResponse)(nil), // 1: pb.EchoResponse
}
var file_service_proto_depIdxs = []int32{
	0, // 0: pb.Sample.Echo:input_type -> pb.EchoRequest
	0, // 1: pb.Sample.EchoLeaderOnly:input_type -> pb.EchoRequest
	0, // 2: pb.Sample.EchoTunnel:input_type -> pb.EchoRequest
	1, // 3: pb.Sample.Echo:output_type -> pb.EchoResponse
	1, // 4: pb.Sample.EchoLeaderOnly:output_type -> pb.EchoResponse
	1, // 5: pb.Sample.EchoTunnel:output_type -> pb.EchoResponse
	3, // [3:6] is the sub-list for method output_type
	0, // [0:3] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_service_proto_init() }
func file_service_proto_init() {
	if File_service_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_service_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*EchoRequest); i {
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
		file_service_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*EchoResponse); i {
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
			RawDescriptor: file_service_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_service_proto_goTypes,
		DependencyIndexes: file_service_proto_depIdxs,
		MessageInfos:      file_service_proto_msgTypes,
	}.Build()
	File_service_proto = out.File
	file_service_proto_rawDesc = nil
	file_service_proto_goTypes = nil
	file_service_proto_depIdxs = nil
}
