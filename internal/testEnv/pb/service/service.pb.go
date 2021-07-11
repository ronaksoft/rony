// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.15.8
// source: service.proto

package service

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

type GetRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Key []byte `protobuf:"bytes,1,opt,name=Key,proto3" json:"Key,omitempty"`
}

func (x *GetRequest) Reset() {
	*x = GetRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_service_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetRequest) ProtoMessage() {}

func (x *GetRequest) ProtoReflect() protoreflect.Message {
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

// Deprecated: Use GetRequest.ProtoReflect.Descriptor instead.
func (*GetRequest) Descriptor() ([]byte, []int) {
	return file_service_proto_rawDescGZIP(), []int{0}
}

func (x *GetRequest) GetKey() []byte {
	if x != nil {
		return x.Key
	}
	return nil
}

type GetResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Key   []byte `protobuf:"bytes,1,opt,name=Key,proto3" json:"Key,omitempty"`
	Value []byte `protobuf:"bytes,2,opt,name=Value,proto3" json:"Value,omitempty"`
}

func (x *GetResponse) Reset() {
	*x = GetResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_service_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetResponse) ProtoMessage() {}

func (x *GetResponse) ProtoReflect() protoreflect.Message {
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

// Deprecated: Use GetResponse.ProtoReflect.Descriptor instead.
func (*GetResponse) Descriptor() ([]byte, []int) {
	return file_service_proto_rawDescGZIP(), []int{1}
}

func (x *GetResponse) GetKey() []byte {
	if x != nil {
		return x.Key
	}
	return nil
}

func (x *GetResponse) GetValue() []byte {
	if x != nil {
		return x.Value
	}
	return nil
}

type SetRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Key   []byte `protobuf:"bytes,1,opt,name=Key,proto3" json:"Key,omitempty"`
	Value []byte `protobuf:"bytes,2,opt,name=Value,proto3" json:"Value,omitempty"`
}

func (x *SetRequest) Reset() {
	*x = SetRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_service_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SetRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SetRequest) ProtoMessage() {}

func (x *SetRequest) ProtoReflect() protoreflect.Message {
	mi := &file_service_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SetRequest.ProtoReflect.Descriptor instead.
func (*SetRequest) Descriptor() ([]byte, []int) {
	return file_service_proto_rawDescGZIP(), []int{2}
}

func (x *SetRequest) GetKey() []byte {
	if x != nil {
		return x.Key
	}
	return nil
}

func (x *SetRequest) GetValue() []byte {
	if x != nil {
		return x.Value
	}
	return nil
}

type SetResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	OK bool `protobuf:"varint,1,opt,name=OK,proto3" json:"OK,omitempty"`
}

func (x *SetResponse) Reset() {
	*x = SetResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_service_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SetResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SetResponse) ProtoMessage() {}

func (x *SetResponse) ProtoReflect() protoreflect.Message {
	mi := &file_service_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SetResponse.ProtoReflect.Descriptor instead.
func (*SetResponse) Descriptor() ([]byte, []int) {
	return file_service_proto_rawDescGZIP(), []int{3}
}

func (x *SetResponse) GetOK() bool {
	if x != nil {
		return x.OK
	}
	return false
}

type EchoRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Int        int64  `protobuf:"varint,1,opt,name=Int,proto3" json:"Int,omitempty"`
	Timestamp  int64  `protobuf:"varint,3,opt,name=Timestamp,proto3" json:"Timestamp,omitempty"`
	ReplicaSet uint64 `protobuf:"varint,4,opt,name=ReplicaSet,proto3" json:"ReplicaSet,omitempty"`
}

func (x *EchoRequest) Reset() {
	*x = EchoRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_service_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *EchoRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EchoRequest) ProtoMessage() {}

func (x *EchoRequest) ProtoReflect() protoreflect.Message {
	mi := &file_service_proto_msgTypes[4]
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
	return file_service_proto_rawDescGZIP(), []int{4}
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

func (x *EchoRequest) GetReplicaSet() uint64 {
	if x != nil {
		return x.ReplicaSet
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
		mi := &file_service_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *EchoResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EchoResponse) ProtoMessage() {}

func (x *EchoResponse) ProtoReflect() protoreflect.Message {
	mi := &file_service_proto_msgTypes[5]
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
	return file_service_proto_rawDescGZIP(), []int{5}
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

type Message1 struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Param1 int64       `protobuf:"varint,1,opt,name=Param1,proto3" json:"Param1,omitempty"`
	Param2 string      `protobuf:"bytes,2,opt,name=Param2,proto3" json:"Param2,omitempty"`
	M2     *Message2   `protobuf:"bytes,3,opt,name=M2,proto3" json:"M2,omitempty"`
	M2S    []*Message2 `protobuf:"bytes,4,rep,name=M2S,proto3" json:"M2S,omitempty"`
}

func (x *Message1) Reset() {
	*x = Message1{}
	if protoimpl.UnsafeEnabled {
		mi := &file_service_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Message1) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Message1) ProtoMessage() {}

func (x *Message1) ProtoReflect() protoreflect.Message {
	mi := &file_service_proto_msgTypes[6]
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
	return file_service_proto_rawDescGZIP(), []int{6}
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

func (x *Message1) GetM2() *Message2 {
	if x != nil {
		return x.M2
	}
	return nil
}

func (x *Message1) GetM2S() []*Message2 {
	if x != nil {
		return x.M2S
	}
	return nil
}

type Message2 struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Param1 uint32    `protobuf:"fixed32,1,opt,name=Param1,proto3" json:"Param1,omitempty"`
	P2     []byte    `protobuf:"bytes,2,opt,name=P2,proto3" json:"P2,omitempty"`
	P3     []int32   `protobuf:"varint,3,rep,packed,name=P3,proto3" json:"P3,omitempty"`
	M1     *Message1 `protobuf:"bytes,4,opt,name=M1,proto3" json:"M1,omitempty"`
}

func (x *Message2) Reset() {
	*x = Message2{}
	if protoimpl.UnsafeEnabled {
		mi := &file_service_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Message2) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Message2) ProtoMessage() {}

func (x *Message2) ProtoReflect() protoreflect.Message {
	mi := &file_service_proto_msgTypes[7]
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
	return file_service_proto_rawDescGZIP(), []int{7}
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

func (x *Message2) GetM1() *Message1 {
	if x != nil {
		return x.M1
	}
	return nil
}

var File_service_proto protoreflect.FileDescriptor

var file_service_proto_rawDesc = []byte{
	0x0a, 0x0d, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12,
	0x07, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x1a, 0x0d, 0x6f, 0x70, 0x74, 0x69, 0x6f, 0x6e,
	0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x1e, 0x0a, 0x0a, 0x47, 0x65, 0x74, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x10, 0x0a, 0x03, 0x4b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x0c, 0x52, 0x03, 0x4b, 0x65, 0x79, 0x22, 0x35, 0x0a, 0x0b, 0x47, 0x65, 0x74, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x10, 0x0a, 0x03, 0x4b, 0x65, 0x79, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x0c, 0x52, 0x03, 0x4b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x56, 0x61, 0x6c, 0x75,
	0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x05, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x22, 0x34,
	0x0a, 0x0a, 0x53, 0x65, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x10, 0x0a, 0x03,
	0x4b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x03, 0x4b, 0x65, 0x79, 0x12, 0x14,
	0x0a, 0x05, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x05, 0x56,
	0x61, 0x6c, 0x75, 0x65, 0x22, 0x1d, 0x0a, 0x0b, 0x53, 0x65, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x12, 0x0e, 0x0a, 0x02, 0x4f, 0x4b, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52,
	0x02, 0x4f, 0x4b, 0x22, 0x5d, 0x0a, 0x0b, 0x45, 0x63, 0x68, 0x6f, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x12, 0x10, 0x0a, 0x03, 0x49, 0x6e, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52,
	0x03, 0x49, 0x6e, 0x74, 0x12, 0x1c, 0x0a, 0x09, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d,
	0x70, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61,
	0x6d, 0x70, 0x12, 0x1e, 0x0a, 0x0a, 0x52, 0x65, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x53, 0x65, 0x74,
	0x18, 0x04, 0x20, 0x01, 0x28, 0x04, 0x52, 0x0a, 0x52, 0x65, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x53,
	0x65, 0x74, 0x22, 0x8e, 0x01, 0x0a, 0x0c, 0x45, 0x63, 0x68, 0x6f, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x12, 0x10, 0x0a, 0x03, 0x49, 0x6e, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03,
	0x52, 0x03, 0x49, 0x6e, 0x74, 0x12, 0x1c, 0x0a, 0x09, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x64,
	0x65, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x64, 0x65, 0x72, 0x12, 0x1c, 0x0a, 0x09, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70,
	0x18, 0x04, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d,
	0x70, 0x12, 0x14, 0x0a, 0x05, 0x44, 0x65, 0x6c, 0x61, 0x79, 0x18, 0x05, 0x20, 0x01, 0x28, 0x03,
	0x52, 0x05, 0x44, 0x65, 0x6c, 0x61, 0x79, 0x12, 0x1a, 0x0a, 0x08, 0x53, 0x65, 0x72, 0x76, 0x65,
	0x72, 0x49, 0x44, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x53, 0x65, 0x72, 0x76, 0x65,
	0x72, 0x49, 0x44, 0x22, 0x82, 0x01, 0x0a, 0x08, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x31,
	0x12, 0x16, 0x0a, 0x06, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x31, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03,
	0x52, 0x06, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x31, 0x12, 0x16, 0x0a, 0x06, 0x50, 0x61, 0x72, 0x61,
	0x6d, 0x32, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x32,
	0x12, 0x21, 0x0a, 0x02, 0x4d, 0x32, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x11, 0x2e, 0x73,
	0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x32, 0x52,
	0x02, 0x4d, 0x32, 0x12, 0x23, 0x0a, 0x03, 0x4d, 0x32, 0x53, 0x18, 0x04, 0x20, 0x03, 0x28, 0x0b,
	0x32, 0x11, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61,
	0x67, 0x65, 0x32, 0x52, 0x03, 0x4d, 0x32, 0x53, 0x22, 0x65, 0x0a, 0x08, 0x4d, 0x65, 0x73, 0x73,
	0x61, 0x67, 0x65, 0x32, 0x12, 0x16, 0x0a, 0x06, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x31, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x07, 0x52, 0x06, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x31, 0x12, 0x0e, 0x0a, 0x02,
	0x50, 0x32, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x02, 0x50, 0x32, 0x12, 0x0e, 0x0a, 0x02,
	0x50, 0x33, 0x18, 0x03, 0x20, 0x03, 0x28, 0x05, 0x52, 0x02, 0x50, 0x33, 0x12, 0x21, 0x0a, 0x02,
	0x4d, 0x31, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x11, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x69,
	0x63, 0x65, 0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x31, 0x52, 0x02, 0x4d, 0x31, 0x32,
	0xda, 0x03, 0x0a, 0x06, 0x53, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x12, 0x47, 0x0a, 0x04, 0x45, 0x63,
	0x68, 0x6f, 0x12, 0x14, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x45, 0x63, 0x68,
	0x6f, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x15, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x69,
	0x63, 0x65, 0x2e, 0x45, 0x63, 0x68, 0x6f, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22,
	0x12, 0x9a, 0xb5, 0x18, 0x0e, 0x0a, 0x03, 0x67, 0x65, 0x74, 0x12, 0x05, 0x2f, 0x65, 0x63, 0x68,
	0x6f, 0x20, 0x01, 0x12, 0x42, 0x0a, 0x03, 0x53, 0x65, 0x74, 0x12, 0x13, 0x2e, 0x73, 0x65, 0x72,
	0x76, 0x69, 0x63, 0x65, 0x2e, 0x53, 0x65, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a,
	0x14, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x53, 0x65, 0x74, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x10, 0x9a, 0xb5, 0x18, 0x0c, 0x0a, 0x04, 0x70, 0x6f, 0x73,
	0x74, 0x12, 0x04, 0x2f, 0x73, 0x65, 0x74, 0x12, 0x50, 0x0a, 0x03, 0x47, 0x65, 0x74, 0x12, 0x13,
	0x2e, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x47, 0x65, 0x74, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x1a, 0x14, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x47, 0x65,
	0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x1e, 0x9a, 0xb5, 0x18, 0x1a, 0x0a,
	0x03, 0x67, 0x65, 0x74, 0x12, 0x13, 0x2f, 0x72, 0x65, 0x71, 0x2f, 0x3a, 0x4b, 0x65, 0x79, 0x2f,
	0x73, 0x6f, 0x6d, 0x65, 0x74, 0x68, 0x69, 0x6e, 0x67, 0x12, 0x6e, 0x0a, 0x0a, 0x45, 0x63, 0x68,
	0x6f, 0x54, 0x75, 0x6e, 0x6e, 0x65, 0x6c, 0x12, 0x14, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63,
	0x65, 0x2e, 0x45, 0x63, 0x68, 0x6f, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x15, 0x2e,
	0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x45, 0x63, 0x68, 0x6f, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x22, 0x33, 0x9a, 0xb5, 0x18, 0x2f, 0x0a, 0x03, 0x67, 0x65, 0x74, 0x12,
	0x13, 0x2f, 0x65, 0x63, 0x68, 0x6f, 0x5f, 0x74, 0x75, 0x6e, 0x6e, 0x65, 0x6c, 0x2f, 0x3a, 0x58,
	0x2f, 0x3a, 0x59, 0x59, 0x1a, 0x05, 0x58, 0x3d, 0x49, 0x6e, 0x74, 0x1a, 0x0c, 0x59, 0x59, 0x3d,
	0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x12, 0x41, 0x0a, 0x0c, 0x45, 0x63, 0x68,
	0x6f, 0x49, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x12, 0x14, 0x2e, 0x73, 0x65, 0x72, 0x76,
	0x69, 0x63, 0x65, 0x2e, 0x45, 0x63, 0x68, 0x6f, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a,
	0x15, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x45, 0x63, 0x68, 0x6f, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x04, 0x90, 0xb5, 0x18, 0x01, 0x12, 0x38, 0x0a, 0x09,
	0x45, 0x63, 0x68, 0x6f, 0x44, 0x65, 0x6c, 0x61, 0x79, 0x12, 0x14, 0x2e, 0x73, 0x65, 0x72, 0x76,
	0x69, 0x63, 0x65, 0x2e, 0x45, 0x63, 0x68, 0x6f, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a,
	0x15, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x45, 0x63, 0x68, 0x6f, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x1a, 0x04, 0x88, 0xb5, 0x18, 0x01, 0x42, 0x0c, 0x5a, 0x0a,
	0x2e, 0x2f, 0x3b, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x33,
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

var file_service_proto_msgTypes = make([]protoimpl.MessageInfo, 8)
var file_service_proto_goTypes = []interface{}{
	(*GetRequest)(nil),   // 0: service.GetRequest
	(*GetResponse)(nil),  // 1: service.GetResponse
	(*SetRequest)(nil),   // 2: service.SetRequest
	(*SetResponse)(nil),  // 3: service.SetResponse
	(*EchoRequest)(nil),  // 4: service.EchoRequest
	(*EchoResponse)(nil), // 5: service.EchoResponse
	(*Message1)(nil),     // 6: service.Message1
	(*Message2)(nil),     // 7: service.Message2
}
var file_service_proto_depIdxs = []int32{
	7, // 0: service.Message1.M2:type_name -> service.Message2
	7, // 1: service.Message1.M2S:type_name -> service.Message2
	6, // 2: service.Message2.M1:type_name -> service.Message1
	4, // 3: service.Sample.Echo:input_type -> service.EchoRequest
	2, // 4: service.Sample.Set:input_type -> service.SetRequest
	0, // 5: service.Sample.Get:input_type -> service.GetRequest
	4, // 6: service.Sample.EchoTunnel:input_type -> service.EchoRequest
	4, // 7: service.Sample.EchoInternal:input_type -> service.EchoRequest
	4, // 8: service.Sample.EchoDelay:input_type -> service.EchoRequest
	5, // 9: service.Sample.Echo:output_type -> service.EchoResponse
	3, // 10: service.Sample.Set:output_type -> service.SetResponse
	1, // 11: service.Sample.Get:output_type -> service.GetResponse
	5, // 12: service.Sample.EchoTunnel:output_type -> service.EchoResponse
	5, // 13: service.Sample.EchoInternal:output_type -> service.EchoResponse
	5, // 14: service.Sample.EchoDelay:output_type -> service.EchoResponse
	9, // [9:15] is the sub-list for method output_type
	3, // [3:9] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_service_proto_init() }
func file_service_proto_init() {
	if File_service_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_service_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
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
		file_service_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetResponse); i {
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
		file_service_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SetRequest); i {
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
		file_service_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SetResponse); i {
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
		file_service_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
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
		file_service_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
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
		file_service_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
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
		file_service_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
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
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_service_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   8,
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
