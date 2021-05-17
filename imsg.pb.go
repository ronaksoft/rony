// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.15.8
// source: imsg.proto

package rony

import (
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

// GetPage
// @Function
// @Return: Page
type GetPage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	PageID     uint32 `protobuf:"varint,1,opt,name=PageID,proto3" json:"PageID,omitempty"`
	ReplicaSet uint64 `protobuf:"varint,2,opt,name=ReplicaSet,proto3" json:"ReplicaSet,omitempty"`
}

func (x *GetPage) Reset() {
	*x = GetPage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_imsg_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetPage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetPage) ProtoMessage() {}

func (x *GetPage) ProtoReflect() protoreflect.Message {
	mi := &file_imsg_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetPage.ProtoReflect.Descriptor instead.
func (*GetPage) Descriptor() ([]byte, []int) {
	return file_imsg_proto_rawDescGZIP(), []int{0}
}

func (x *GetPage) GetPageID() uint32 {
	if x != nil {
		return x.PageID
	}
	return 0
}

func (x *GetPage) GetReplicaSet() uint64 {
	if x != nil {
		return x.ReplicaSet
	}
	return 0
}

// StoreMessage
// This is dummy message, we only use its constructor. The actual message is
// raftpb.Message
type StoreMessage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Dummy bool `protobuf:"varint,1,opt,name=Dummy,proto3" json:"Dummy,omitempty"`
}

func (x *StoreMessage) Reset() {
	*x = StoreMessage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_imsg_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StoreMessage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StoreMessage) ProtoMessage() {}

func (x *StoreMessage) ProtoReflect() protoreflect.Message {
	mi := &file_imsg_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StoreMessage.ProtoReflect.Descriptor instead.
func (*StoreMessage) Descriptor() ([]byte, []int) {
	return file_imsg_proto_rawDescGZIP(), []int{1}
}

func (x *StoreMessage) GetDummy() bool {
	if x != nil {
		return x.Dummy
	}
	return false
}

// TunnelMessage
type TunnelMessage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	SenderID         []byte           `protobuf:"bytes,1,opt,name=SenderID,proto3" json:"SenderID,omitempty"`
	SenderReplicaSet uint64           `protobuf:"varint,2,opt,name=SenderReplicaSet,proto3" json:"SenderReplicaSet,omitempty"`
	Store            []*KeyValue      `protobuf:"bytes,3,rep,name=Store,proto3" json:"Store,omitempty"`
	Envelope         *MessageEnvelope `protobuf:"bytes,4,opt,name=Envelope,proto3" json:"Envelope,omitempty"`
}

func (x *TunnelMessage) Reset() {
	*x = TunnelMessage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_imsg_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TunnelMessage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TunnelMessage) ProtoMessage() {}

func (x *TunnelMessage) ProtoReflect() protoreflect.Message {
	mi := &file_imsg_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TunnelMessage.ProtoReflect.Descriptor instead.
func (*TunnelMessage) Descriptor() ([]byte, []int) {
	return file_imsg_proto_rawDescGZIP(), []int{2}
}

func (x *TunnelMessage) GetSenderID() []byte {
	if x != nil {
		return x.SenderID
	}
	return nil
}

func (x *TunnelMessage) GetSenderReplicaSet() uint64 {
	if x != nil {
		return x.SenderReplicaSet
	}
	return 0
}

func (x *TunnelMessage) GetStore() []*KeyValue {
	if x != nil {
		return x.Store
	}
	return nil
}

func (x *TunnelMessage) GetEnvelope() *MessageEnvelope {
	if x != nil {
		return x.Envelope
	}
	return nil
}

// EdgeNode
type EdgeNode struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ServerID    []byte   `protobuf:"bytes,1,opt,name=ServerID,proto3" json:"ServerID,omitempty"`
	ReplicaSet  uint64   `protobuf:"varint,2,opt,name=ReplicaSet,proto3" json:"ReplicaSet,omitempty"`
	Hash        uint64   `protobuf:"varint,3,opt,name=Hash,proto3" json:"Hash,omitempty"`
	GatewayAddr []string `protobuf:"bytes,7,rep,name=GatewayAddr,proto3" json:"GatewayAddr,omitempty"`
	TunnelAddr  []string `protobuf:"bytes,8,rep,name=TunnelAddr,proto3" json:"TunnelAddr,omitempty"`
}

func (x *EdgeNode) Reset() {
	*x = EdgeNode{}
	if protoimpl.UnsafeEnabled {
		mi := &file_imsg_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *EdgeNode) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EdgeNode) ProtoMessage() {}

func (x *EdgeNode) ProtoReflect() protoreflect.Message {
	mi := &file_imsg_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EdgeNode.ProtoReflect.Descriptor instead.
func (*EdgeNode) Descriptor() ([]byte, []int) {
	return file_imsg_proto_rawDescGZIP(), []int{3}
}

func (x *EdgeNode) GetServerID() []byte {
	if x != nil {
		return x.ServerID
	}
	return nil
}

func (x *EdgeNode) GetReplicaSet() uint64 {
	if x != nil {
		return x.ReplicaSet
	}
	return 0
}

func (x *EdgeNode) GetHash() uint64 {
	if x != nil {
		return x.Hash
	}
	return 0
}

func (x *EdgeNode) GetGatewayAddr() []string {
	if x != nil {
		return x.GatewayAddr
	}
	return nil
}

func (x *EdgeNode) GetTunnelAddr() []string {
	if x != nil {
		return x.TunnelAddr
	}
	return nil
}

// Page
type Page struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ID         uint32 `protobuf:"varint,1,opt,name=ID,proto3" json:"ID,omitempty"`
	ReplicaSet uint64 `protobuf:"varint,2,opt,name=ReplicaSet,proto3" json:"ReplicaSet,omitempty"`
}

func (x *Page) Reset() {
	*x = Page{}
	if protoimpl.UnsafeEnabled {
		mi := &file_imsg_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Page) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Page) ProtoMessage() {}

func (x *Page) ProtoReflect() protoreflect.Message {
	mi := &file_imsg_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Page.ProtoReflect.Descriptor instead.
func (*Page) Descriptor() ([]byte, []int) {
	return file_imsg_proto_rawDescGZIP(), []int{4}
}

func (x *Page) GetID() uint32 {
	if x != nil {
		return x.ID
	}
	return 0
}

func (x *Page) GetReplicaSet() uint64 {
	if x != nil {
		return x.ReplicaSet
	}
	return 0
}

var File_imsg_proto protoreflect.FileDescriptor

var file_imsg_proto_rawDesc = []byte{
	0x0a, 0x0a, 0x69, 0x6d, 0x73, 0x67, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x04, 0x72, 0x6f,
	0x6e, 0x79, 0x1a, 0x0d, 0x6f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x1a, 0x09, 0x6d, 0x73, 0x67, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x41, 0x0a, 0x07,
	0x47, 0x65, 0x74, 0x50, 0x61, 0x67, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x50, 0x61, 0x67, 0x65, 0x49,
	0x44, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x06, 0x50, 0x61, 0x67, 0x65, 0x49, 0x44, 0x12,
	0x1e, 0x0a, 0x0a, 0x52, 0x65, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x53, 0x65, 0x74, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x04, 0x52, 0x0a, 0x52, 0x65, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x53, 0x65, 0x74, 0x22,
	0x24, 0x0a, 0x0c, 0x53, 0x74, 0x6f, 0x72, 0x65, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12,
	0x14, 0x0a, 0x05, 0x44, 0x75, 0x6d, 0x6d, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x05,
	0x44, 0x75, 0x6d, 0x6d, 0x79, 0x22, 0xb0, 0x01, 0x0a, 0x0d, 0x54, 0x75, 0x6e, 0x6e, 0x65, 0x6c,
	0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x53, 0x65, 0x6e, 0x64, 0x65,
	0x72, 0x49, 0x44, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x08, 0x53, 0x65, 0x6e, 0x64, 0x65,
	0x72, 0x49, 0x44, 0x12, 0x2a, 0x0a, 0x10, 0x53, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x52, 0x65, 0x70,
	0x6c, 0x69, 0x63, 0x61, 0x53, 0x65, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x04, 0x52, 0x10, 0x53,
	0x65, 0x6e, 0x64, 0x65, 0x72, 0x52, 0x65, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x53, 0x65, 0x74, 0x12,
	0x24, 0x0a, 0x05, 0x53, 0x74, 0x6f, 0x72, 0x65, 0x18, 0x03, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0e,
	0x2e, 0x72, 0x6f, 0x6e, 0x79, 0x2e, 0x4b, 0x65, 0x79, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x52, 0x05,
	0x53, 0x74, 0x6f, 0x72, 0x65, 0x12, 0x31, 0x0a, 0x08, 0x45, 0x6e, 0x76, 0x65, 0x6c, 0x6f, 0x70,
	0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x15, 0x2e, 0x72, 0x6f, 0x6e, 0x79, 0x2e, 0x4d,
	0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x45, 0x6e, 0x76, 0x65, 0x6c, 0x6f, 0x70, 0x65, 0x52, 0x08,
	0x45, 0x6e, 0x76, 0x65, 0x6c, 0x6f, 0x70, 0x65, 0x22, 0x9c, 0x01, 0x0a, 0x08, 0x45, 0x64, 0x67,
	0x65, 0x4e, 0x6f, 0x64, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x53, 0x65, 0x72, 0x76, 0x65, 0x72, 0x49,
	0x44, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x08, 0x53, 0x65, 0x72, 0x76, 0x65, 0x72, 0x49,
	0x44, 0x12, 0x1e, 0x0a, 0x0a, 0x52, 0x65, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x53, 0x65, 0x74, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x04, 0x52, 0x0a, 0x52, 0x65, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x53, 0x65,
	0x74, 0x12, 0x12, 0x0a, 0x04, 0x48, 0x61, 0x73, 0x68, 0x18, 0x03, 0x20, 0x01, 0x28, 0x04, 0x52,
	0x04, 0x48, 0x61, 0x73, 0x68, 0x12, 0x20, 0x0a, 0x0b, 0x47, 0x61, 0x74, 0x65, 0x77, 0x61, 0x79,
	0x41, 0x64, 0x64, 0x72, 0x18, 0x07, 0x20, 0x03, 0x28, 0x09, 0x52, 0x0b, 0x47, 0x61, 0x74, 0x65,
	0x77, 0x61, 0x79, 0x41, 0x64, 0x64, 0x72, 0x12, 0x1e, 0x0a, 0x0a, 0x54, 0x75, 0x6e, 0x6e, 0x65,
	0x6c, 0x41, 0x64, 0x64, 0x72, 0x18, 0x08, 0x20, 0x03, 0x28, 0x09, 0x52, 0x0a, 0x54, 0x75, 0x6e,
	0x6e, 0x65, 0x6c, 0x41, 0x64, 0x64, 0x72, 0x22, 0x60, 0x0a, 0x04, 0x50, 0x61, 0x67, 0x65, 0x12,
	0x0e, 0x0a, 0x02, 0x49, 0x44, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x02, 0x49, 0x44, 0x12,
	0x1e, 0x0a, 0x0a, 0x52, 0x65, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x53, 0x65, 0x74, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x04, 0x52, 0x0a, 0x52, 0x65, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x53, 0x65, 0x74, 0x3a,
	0x28, 0x88, 0xb5, 0x18, 0x01, 0xd2, 0xb5, 0x18, 0x04, 0x63, 0x72, 0x75, 0x64, 0xf2, 0xb5, 0x18,
	0x04, 0x28, 0x49, 0x44, 0x29, 0xfa, 0xb5, 0x18, 0x10, 0x28, 0x52, 0x65, 0x70, 0x6c, 0x69, 0x63,
	0x61, 0x53, 0x65, 0x74, 0x2c, 0x20, 0x49, 0x44, 0x29, 0x42, 0x1b, 0x5a, 0x19, 0x67, 0x69, 0x74,
	0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x72, 0x6f, 0x6e, 0x61, 0x6b, 0x73, 0x6f, 0x66,
	0x74, 0x2f, 0x72, 0x6f, 0x6e, 0x79, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_imsg_proto_rawDescOnce sync.Once
	file_imsg_proto_rawDescData = file_imsg_proto_rawDesc
)

func file_imsg_proto_rawDescGZIP() []byte {
	file_imsg_proto_rawDescOnce.Do(func() {
		file_imsg_proto_rawDescData = protoimpl.X.CompressGZIP(file_imsg_proto_rawDescData)
	})
	return file_imsg_proto_rawDescData
}

var file_imsg_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_imsg_proto_goTypes = []interface{}{
	(*GetPage)(nil),         // 0: rony.GetPage
	(*StoreMessage)(nil),    // 1: rony.StoreMessage
	(*TunnelMessage)(nil),   // 2: rony.TunnelMessage
	(*EdgeNode)(nil),        // 3: rony.EdgeNode
	(*Page)(nil),            // 4: rony.Page
	(*KeyValue)(nil),        // 5: rony.KeyValue
	(*MessageEnvelope)(nil), // 6: rony.MessageEnvelope
}
var file_imsg_proto_depIdxs = []int32{
	5, // 0: rony.TunnelMessage.Store:type_name -> rony.KeyValue
	6, // 1: rony.TunnelMessage.Envelope:type_name -> rony.MessageEnvelope
	2, // [2:2] is the sub-list for method output_type
	2, // [2:2] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_imsg_proto_init() }
func file_imsg_proto_init() {
	if File_imsg_proto != nil {
		return
	}
	file_options_proto_init()
	file_msg_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_imsg_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetPage); i {
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
		file_imsg_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*StoreMessage); i {
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
		file_imsg_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TunnelMessage); i {
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
		file_imsg_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*EdgeNode); i {
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
		file_imsg_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Page); i {
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
			RawDescriptor: file_imsg_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_imsg_proto_goTypes,
		DependencyIndexes: file_imsg_proto_depIdxs,
		MessageInfos:      file_imsg_proto_msgTypes,
	}.Build()
	File_imsg_proto = out.File
	file_imsg_proto_rawDesc = nil
	file_imsg_proto_goTypes = nil
	file_imsg_proto_depIdxs = nil
}
