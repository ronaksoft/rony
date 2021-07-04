// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.15.8
// source: options.proto

package rony

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	descriptorpb "google.golang.org/protobuf/types/descriptorpb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// RestOpt is the container message for all the required and optional parameters for REST setup.
// By adding 'rony_rest' option to the method definition a REST wrapper will be generated for that
// method and that RPC handler will become accessible using REST endpoint.
type RestOpt struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// method identifies the HTTP command: i.e. 'get', 'post', 'put', ...
	Method *string `protobuf:"bytes,1,req,name=method" json:"method,omitempty"`
	// path identifies the path that this RPC method will be available as REST endpoint.
	Path *string `protobuf:"bytes,2,req,name=path" json:"path,omitempty"`
	// bind_variables is a comma separated pair of bindings. For example if we have a path '/part1/:var1/:var2'
	// and we want to bind 'var1' to field 'fieldVar1' and bind 'var2' to field 'fieldVar2' then it will be:
	// `var1=fieldVar1,var2=fieldVar2`
	BindVariables *string `protobuf:"bytes,3,opt,name=bind_variables,json=bindVariables" json:"bind_variables,omitempty"`
	// json_encode if is set then input and output data are json encoded version of the proto messages
	JsonEncode *bool `protobuf:"varint,4,opt,name=json_encode,json=jsonEncode" json:"json_encode,omitempty"`
}

func (x *RestOpt) Reset() {
	*x = RestOpt{}
	if protoimpl.UnsafeEnabled {
		mi := &file_options_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RestOpt) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RestOpt) ProtoMessage() {}

func (x *RestOpt) ProtoReflect() protoreflect.Message {
	mi := &file_options_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RestOpt.ProtoReflect.Descriptor instead.
func (*RestOpt) Descriptor() ([]byte, []int) {
	return file_options_proto_rawDescGZIP(), []int{0}
}

func (x *RestOpt) GetMethod() string {
	if x != nil && x.Method != nil {
		return *x.Method
	}
	return ""
}

func (x *RestOpt) GetPath() string {
	if x != nil && x.Path != nil {
		return *x.Path
	}
	return ""
}

func (x *RestOpt) GetBindVariables() string {
	if x != nil && x.BindVariables != nil {
		return *x.BindVariables
	}
	return ""
}

func (x *RestOpt) GetJsonEncode() bool {
	if x != nil && x.JsonEncode != nil {
		return *x.JsonEncode
	}
	return false
}

var file_options_proto_extTypes = []protoimpl.ExtensionInfo{
	{
		ExtendedType:  (*descriptorpb.ServiceOptions)(nil),
		ExtensionType: (*bool)(nil),
		Field:         50001,
		Name:          "rony_cobra_cmd",
		Tag:           "varint,50001,opt,name=rony_cobra_cmd",
		Filename:      "options.proto",
	},
	{
		ExtendedType:  (*descriptorpb.ServiceOptions)(nil),
		ExtensionType: (*string)(nil),
		Field:         50002,
		Name:          "rony_cobra_cmd_protocol",
		Tag:           "bytes,50002,opt,name=rony_cobra_cmd_protocol",
		Filename:      "options.proto",
	},
	{
		ExtendedType:  (*descriptorpb.ServiceOptions)(nil),
		ExtensionType: (*bool)(nil),
		Field:         50003,
		Name:          "rony_no_client",
		Tag:           "varint,50003,opt,name=rony_no_client",
		Filename:      "options.proto",
	},
	{
		ExtendedType:  (*descriptorpb.MethodOptions)(nil),
		ExtensionType: (*bool)(nil),
		Field:         50002,
		Name:          "rony_internal",
		Tag:           "varint,50002,opt,name=rony_internal",
		Filename:      "options.proto",
	},
	{
		ExtendedType:  (*descriptorpb.MethodOptions)(nil),
		ExtensionType: (*RestOpt)(nil),
		Field:         50003,
		Name:          "rony_rest",
		Tag:           "bytes,50003,opt,name=rony_rest",
		Filename:      "options.proto",
	},
	{
		ExtendedType:  (*descriptorpb.MessageOptions)(nil),
		ExtensionType: (*bool)(nil),
		Field:         50001,
		Name:          "rony_aggregate",
		Tag:           "varint,50001,opt,name=rony_aggregate",
		Filename:      "options.proto",
	},
	{
		ExtendedType:  (*descriptorpb.MessageOptions)(nil),
		ExtensionType: (*bool)(nil),
		Field:         50002,
		Name:          "rony_singleton",
		Tag:           "varint,50002,opt,name=rony_singleton",
		Filename:      "options.proto",
	},
	{
		ExtendedType:  (*descriptorpb.MessageOptions)(nil),
		ExtensionType: (*string)(nil),
		Field:         50010,
		Name:          "rony_aggregate_type",
		Tag:           "bytes,50010,opt,name=rony_aggregate_type",
		Filename:      "options.proto",
	},
	{
		ExtendedType:  (*descriptorpb.MessageOptions)(nil),
		ExtensionType: (*bool)(nil),
		Field:         50011,
		Name:          "rony_aggregate_command",
		Tag:           "varint,50011,opt,name=rony_aggregate_command",
		Filename:      "options.proto",
	},
	{
		ExtendedType:  (*descriptorpb.MessageOptions)(nil),
		ExtensionType: (*bool)(nil),
		Field:         50012,
		Name:          "rony_aggregate_event",
		Tag:           "varint,50012,opt,name=rony_aggregate_event",
		Filename:      "options.proto",
	},
	{
		ExtendedType:  (*descriptorpb.MessageOptions)(nil),
		ExtensionType: (*string)(nil),
		Field:         50013,
		Name:          "rony_aggregate_link",
		Tag:           "bytes,50013,opt,name=rony_aggregate_link",
		Filename:      "options.proto",
	},
	{
		ExtendedType:  (*descriptorpb.MessageOptions)(nil),
		ExtensionType: (*string)(nil),
		Field:         50014,
		Name:          "rony_aggregate_table",
		Tag:           "bytes,50014,opt,name=rony_aggregate_table",
		Filename:      "options.proto",
	},
	{
		ExtendedType:  (*descriptorpb.MessageOptions)(nil),
		ExtensionType: ([]string)(nil),
		Field:         50015,
		Name:          "rony_aggregate_view",
		Tag:           "bytes,50015,rep,name=rony_aggregate_view",
		Filename:      "options.proto",
	},
	{
		ExtendedType:  (*descriptorpb.FieldOptions)(nil),
		ExtensionType: (*bool)(nil),
		Field:         50001,
		Name:          "rony_index",
		Tag:           "varint,50001,opt,name=rony_index",
		Filename:      "options.proto",
	},
}

// Extension fields to descriptorpb.ServiceOptions.
var (
	// rony_cobra_cmd generates the boiler plate code for client stub of rpc methods, using cobra package.
	//
	// optional bool rony_cobra_cmd = 50001;
	E_RonyCobraCmd = &file_options_proto_extTypes[0]
	// rony_cobra_cmd_protocol defines what protocol should client use to communicate with server.
	// POSSIBLE VALUES: "ws", "http"
	//
	// optional string rony_cobra_cmd_protocol = 50002;
	E_RonyCobraCmdProtocol = &file_options_proto_extTypes[1]
	// rony_no_client if is set then no client code will be generated. This flag is for internal usage.
	// DO NOT USE IT.
	//
	// optional bool rony_no_client = 50003;
	E_RonyNoClient = &file_options_proto_extTypes[2]
)

// Extension fields to descriptorpb.MethodOptions.
var (
	// rony_internal marks this method internal, hence only edges could execute this rpc through tunnel messages. In other words,
	// this command is not exposed to external clients connected through th gateway.
	//
	// optional bool rony_internal = 50002;
	E_RonyInternal = &file_options_proto_extTypes[3]
	// optional RestOpt rony_rest = 50003;
	E_RonyRest = &file_options_proto_extTypes[4]
)

// Extension fields to descriptorpb.MessageOptions.
var (
	// rony_aggregate marks this message as an aggregate, then 'rony_table' and 'rony_view' options become available for this
	// message.
	//
	// optional bool rony_aggregate = 50001;
	E_RonyAggregate = &file_options_proto_extTypes[5]
	// rony_singleton marks this message as a singleton.
	// NOTE: a message could either have 'rony_aggregate' ro 'rony_singleton' options at a same time. Setting both
	// cause unpredictable results.
	//
	// optional bool rony_singleton = 50002;
	E_RonySingleton = &file_options_proto_extTypes[6]
	// rony_aggregate_type makes the code generator to generate appropriate functions based on the way you are going
	// to handle actions on the aggregate.
	// POSSIBLE_VALUES: "crud", "eventsource"
	//
	// optional string rony_aggregate_type = 50010;
	E_RonyAggregateType = &file_options_proto_extTypes[7]
	// rony_aggregate_command makes this message as a command which is ONLY used if rony_aggregate_type is set to 'eventsource'.
	// If this option is set then you MUST also define rony_aggregate_link to identify which aggregate this command will work on.
	//
	// optional bool rony_aggregate_command = 50011;
	E_RonyAggregateCommand = &file_options_proto_extTypes[8]
	// rony_aggregate_event makes this message as an event which is ONLY used if rony_aggregate_type is set to 'eventsource'
	// If this option is set then you MUST also define rony_aggregate_link to identify which aggregate this event will be read from.
	//
	// optional bool rony_aggregate_event = 50012;
	E_RonyAggregateEvent = &file_options_proto_extTypes[9]
	// rony_aggregate_link is the name of the aggregate message which we link this message to.
	//
	// optional string rony_aggregate_link = 50013;
	E_RonyAggregateLink = &file_options_proto_extTypes[10]
	// rony_aggregate_table creates a virtual table presentation to hold instances of this message, like rows in a table
	// PRIMARY KEY FORMAT: ( (partitionKey1, partitionKey2, ...), clusteringKey1, clusteringKey2, ...)
	// NOTE: If there is only one partition key then you could safely drop the parenthesis, i.e. (pk1, ck1, ck2)
	//
	// optional string rony_aggregate_table = 50014;
	E_RonyAggregateTable = &file_options_proto_extTypes[11]
	// rony_aggregate_view creates a materialized view of the aggregate based on the primary key.
	// PRIMARY KEY FORMAT: ( (partitionKey1, partitionKey2, ...), clusteringKey1, clusteringKey2, ...)
	// NOTE (1): If there is only one partition key then you could safely drop the parenthesis, i.e. (pk1, ck1, ck2)
	// NOTE (2): The primary key of the model must contains all the primary key items of the table. They don't need to
	//           follow the same order as table. for example the following is correct:
	//                  rony_aggregate_table = ((a, b), c)
	//                  rony_aggregate_view = ((c, a), d, b)
	//
	// repeated string rony_aggregate_view = 50015;
	E_RonyAggregateView = &file_options_proto_extTypes[12]
)

// Extension fields to descriptorpb.FieldOptions.
var (
	// optional bool rony_index = 50001;
	E_RonyIndex = &file_options_proto_extTypes[13]
)

var File_options_proto protoreflect.FileDescriptor

var file_options_proto_rawDesc = []byte{
	0x0a, 0x0d, 0x6f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a,
	0x20, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2f, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x6f, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x22, 0x7d, 0x0a, 0x07, 0x52, 0x65, 0x73, 0x74, 0x4f, 0x70, 0x74, 0x12, 0x16, 0x0a, 0x06,
	0x6d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x18, 0x01, 0x20, 0x02, 0x28, 0x09, 0x52, 0x06, 0x6d, 0x65,
	0x74, 0x68, 0x6f, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x70, 0x61, 0x74, 0x68, 0x18, 0x02, 0x20, 0x02,
	0x28, 0x09, 0x52, 0x04, 0x70, 0x61, 0x74, 0x68, 0x12, 0x25, 0x0a, 0x0e, 0x62, 0x69, 0x6e, 0x64,
	0x5f, 0x76, 0x61, 0x72, 0x69, 0x61, 0x62, 0x6c, 0x65, 0x73, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x0d, 0x62, 0x69, 0x6e, 0x64, 0x56, 0x61, 0x72, 0x69, 0x61, 0x62, 0x6c, 0x65, 0x73, 0x12,
	0x1f, 0x0a, 0x0b, 0x6a, 0x73, 0x6f, 0x6e, 0x5f, 0x65, 0x6e, 0x63, 0x6f, 0x64, 0x65, 0x18, 0x04,
	0x20, 0x01, 0x28, 0x08, 0x52, 0x0a, 0x6a, 0x73, 0x6f, 0x6e, 0x45, 0x6e, 0x63, 0x6f, 0x64, 0x65,
	0x3a, 0x47, 0x0a, 0x0e, 0x72, 0x6f, 0x6e, 0x79, 0x5f, 0x63, 0x6f, 0x62, 0x72, 0x61, 0x5f, 0x63,
	0x6d, 0x64, 0x12, 0x1f, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2e, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x4f, 0x70, 0x74, 0x69,
	0x6f, 0x6e, 0x73, 0x18, 0xd1, 0x86, 0x03, 0x20, 0x01, 0x28, 0x08, 0x52, 0x0c, 0x72, 0x6f, 0x6e,
	0x79, 0x43, 0x6f, 0x62, 0x72, 0x61, 0x43, 0x6d, 0x64, 0x3a, 0x58, 0x0a, 0x17, 0x72, 0x6f, 0x6e,
	0x79, 0x5f, 0x63, 0x6f, 0x62, 0x72, 0x61, 0x5f, 0x63, 0x6d, 0x64, 0x5f, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x63, 0x6f, 0x6c, 0x12, 0x1f, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x4f, 0x70,
	0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0xd2, 0x86, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x14, 0x72,
	0x6f, 0x6e, 0x79, 0x43, 0x6f, 0x62, 0x72, 0x61, 0x43, 0x6d, 0x64, 0x50, 0x72, 0x6f, 0x74, 0x6f,
	0x63, 0x6f, 0x6c, 0x3a, 0x47, 0x0a, 0x0e, 0x72, 0x6f, 0x6e, 0x79, 0x5f, 0x6e, 0x6f, 0x5f, 0x63,
	0x6c, 0x69, 0x65, 0x6e, 0x74, 0x12, 0x1f, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x4f,
	0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0xd3, 0x86, 0x03, 0x20, 0x01, 0x28, 0x08, 0x52, 0x0c,
	0x72, 0x6f, 0x6e, 0x79, 0x4e, 0x6f, 0x43, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x3a, 0x45, 0x0a, 0x0d,
	0x72, 0x6f, 0x6e, 0x79, 0x5f, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x12, 0x1e, 0x2e,
	0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e,
	0x4d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0xd2, 0x86,
	0x03, 0x20, 0x01, 0x28, 0x08, 0x52, 0x0c, 0x72, 0x6f, 0x6e, 0x79, 0x49, 0x6e, 0x74, 0x65, 0x72,
	0x6e, 0x61, 0x6c, 0x3a, 0x47, 0x0a, 0x09, 0x72, 0x6f, 0x6e, 0x79, 0x5f, 0x72, 0x65, 0x73, 0x74,
	0x12, 0x1e, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62,
	0x75, 0x66, 0x2e, 0x4d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73,
	0x18, 0xd3, 0x86, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x08, 0x2e, 0x52, 0x65, 0x73, 0x74, 0x4f,
	0x70, 0x74, 0x52, 0x08, 0x72, 0x6f, 0x6e, 0x79, 0x52, 0x65, 0x73, 0x74, 0x3a, 0x48, 0x0a, 0x0e,
	0x72, 0x6f, 0x6e, 0x79, 0x5f, 0x61, 0x67, 0x67, 0x72, 0x65, 0x67, 0x61, 0x74, 0x65, 0x12, 0x1f,
	0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18,
	0xd1, 0x86, 0x03, 0x20, 0x01, 0x28, 0x08, 0x52, 0x0d, 0x72, 0x6f, 0x6e, 0x79, 0x41, 0x67, 0x67,
	0x72, 0x65, 0x67, 0x61, 0x74, 0x65, 0x3a, 0x48, 0x0a, 0x0e, 0x72, 0x6f, 0x6e, 0x79, 0x5f, 0x73,
	0x69, 0x6e, 0x67, 0x6c, 0x65, 0x74, 0x6f, 0x6e, 0x12, 0x1f, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c,
	0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61,
	0x67, 0x65, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0xd2, 0x86, 0x03, 0x20, 0x01, 0x28,
	0x08, 0x52, 0x0d, 0x72, 0x6f, 0x6e, 0x79, 0x53, 0x69, 0x6e, 0x67, 0x6c, 0x65, 0x74, 0x6f, 0x6e,
	0x3a, 0x51, 0x0a, 0x13, 0x72, 0x6f, 0x6e, 0x79, 0x5f, 0x61, 0x67, 0x67, 0x72, 0x65, 0x67, 0x61,
	0x74, 0x65, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x12, 0x1f, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67,
	0x65, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0xda, 0x86, 0x03, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x11, 0x72, 0x6f, 0x6e, 0x79, 0x41, 0x67, 0x67, 0x72, 0x65, 0x67, 0x61, 0x74, 0x65, 0x54,
	0x79, 0x70, 0x65, 0x3a, 0x57, 0x0a, 0x16, 0x72, 0x6f, 0x6e, 0x79, 0x5f, 0x61, 0x67, 0x67, 0x72,
	0x65, 0x67, 0x61, 0x74, 0x65, 0x5f, 0x63, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x12, 0x1f, 0x2e,
	0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e,
	0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0xdb,
	0x86, 0x03, 0x20, 0x01, 0x28, 0x08, 0x52, 0x14, 0x72, 0x6f, 0x6e, 0x79, 0x41, 0x67, 0x67, 0x72,
	0x65, 0x67, 0x61, 0x74, 0x65, 0x43, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x3a, 0x53, 0x0a, 0x14,
	0x72, 0x6f, 0x6e, 0x79, 0x5f, 0x61, 0x67, 0x67, 0x72, 0x65, 0x67, 0x61, 0x74, 0x65, 0x5f, 0x65,
	0x76, 0x65, 0x6e, 0x74, 0x12, 0x1f, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x4f, 0x70,
	0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0xdc, 0x86, 0x03, 0x20, 0x01, 0x28, 0x08, 0x52, 0x12, 0x72,
	0x6f, 0x6e, 0x79, 0x41, 0x67, 0x67, 0x72, 0x65, 0x67, 0x61, 0x74, 0x65, 0x45, 0x76, 0x65, 0x6e,
	0x74, 0x3a, 0x51, 0x0a, 0x13, 0x72, 0x6f, 0x6e, 0x79, 0x5f, 0x61, 0x67, 0x67, 0x72, 0x65, 0x67,
	0x61, 0x74, 0x65, 0x5f, 0x6c, 0x69, 0x6e, 0x6b, 0x12, 0x1f, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c,
	0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61,
	0x67, 0x65, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0xdd, 0x86, 0x03, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x11, 0x72, 0x6f, 0x6e, 0x79, 0x41, 0x67, 0x67, 0x72, 0x65, 0x67, 0x61, 0x74, 0x65,
	0x4c, 0x69, 0x6e, 0x6b, 0x3a, 0x53, 0x0a, 0x14, 0x72, 0x6f, 0x6e, 0x79, 0x5f, 0x61, 0x67, 0x67,
	0x72, 0x65, 0x67, 0x61, 0x74, 0x65, 0x5f, 0x74, 0x61, 0x62, 0x6c, 0x65, 0x12, 0x1f, 0x2e, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x4d,
	0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0xde, 0x86,
	0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x12, 0x72, 0x6f, 0x6e, 0x79, 0x41, 0x67, 0x67, 0x72, 0x65,
	0x67, 0x61, 0x74, 0x65, 0x54, 0x61, 0x62, 0x6c, 0x65, 0x3a, 0x51, 0x0a, 0x13, 0x72, 0x6f, 0x6e,
	0x79, 0x5f, 0x61, 0x67, 0x67, 0x72, 0x65, 0x67, 0x61, 0x74, 0x65, 0x5f, 0x76, 0x69, 0x65, 0x77,
	0x12, 0x1f, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62,
	0x75, 0x66, 0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e,
	0x73, 0x18, 0xdf, 0x86, 0x03, 0x20, 0x03, 0x28, 0x09, 0x52, 0x11, 0x72, 0x6f, 0x6e, 0x79, 0x41,
	0x67, 0x67, 0x72, 0x65, 0x67, 0x61, 0x74, 0x65, 0x56, 0x69, 0x65, 0x77, 0x3a, 0x3e, 0x0a, 0x0a,
	0x72, 0x6f, 0x6e, 0x79, 0x5f, 0x69, 0x6e, 0x64, 0x65, 0x78, 0x12, 0x1d, 0x2e, 0x67, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x46, 0x69, 0x65,
	0x6c, 0x64, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0xd1, 0x86, 0x03, 0x20, 0x01, 0x28,
	0x08, 0x52, 0x09, 0x72, 0x6f, 0x6e, 0x79, 0x49, 0x6e, 0x64, 0x65, 0x78, 0x42, 0x1b, 0x5a, 0x19,
	0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x72, 0x6f, 0x6e, 0x61, 0x6b,
	0x73, 0x6f, 0x66, 0x74, 0x2f, 0x72, 0x6f, 0x6e, 0x79,
}

var (
	file_options_proto_rawDescOnce sync.Once
	file_options_proto_rawDescData = file_options_proto_rawDesc
)

func file_options_proto_rawDescGZIP() []byte {
	file_options_proto_rawDescOnce.Do(func() {
		file_options_proto_rawDescData = protoimpl.X.CompressGZIP(file_options_proto_rawDescData)
	})
	return file_options_proto_rawDescData
}

var file_options_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_options_proto_goTypes = []interface{}{
	(*RestOpt)(nil),                     // 0: RestOpt
	(*descriptorpb.ServiceOptions)(nil), // 1: google.protobuf.ServiceOptions
	(*descriptorpb.MethodOptions)(nil),  // 2: google.protobuf.MethodOptions
	(*descriptorpb.MessageOptions)(nil), // 3: google.protobuf.MessageOptions
	(*descriptorpb.FieldOptions)(nil),   // 4: google.protobuf.FieldOptions
}
var file_options_proto_depIdxs = []int32{
	1,  // 0: rony_cobra_cmd:extendee -> google.protobuf.ServiceOptions
	1,  // 1: rony_cobra_cmd_protocol:extendee -> google.protobuf.ServiceOptions
	1,  // 2: rony_no_client:extendee -> google.protobuf.ServiceOptions
	2,  // 3: rony_internal:extendee -> google.protobuf.MethodOptions
	2,  // 4: rony_rest:extendee -> google.protobuf.MethodOptions
	3,  // 5: rony_aggregate:extendee -> google.protobuf.MessageOptions
	3,  // 6: rony_singleton:extendee -> google.protobuf.MessageOptions
	3,  // 7: rony_aggregate_type:extendee -> google.protobuf.MessageOptions
	3,  // 8: rony_aggregate_command:extendee -> google.protobuf.MessageOptions
	3,  // 9: rony_aggregate_event:extendee -> google.protobuf.MessageOptions
	3,  // 10: rony_aggregate_link:extendee -> google.protobuf.MessageOptions
	3,  // 11: rony_aggregate_table:extendee -> google.protobuf.MessageOptions
	3,  // 12: rony_aggregate_view:extendee -> google.protobuf.MessageOptions
	4,  // 13: rony_index:extendee -> google.protobuf.FieldOptions
	0,  // 14: rony_rest:type_name -> RestOpt
	15, // [15:15] is the sub-list for method output_type
	15, // [15:15] is the sub-list for method input_type
	14, // [14:15] is the sub-list for extension type_name
	0,  // [0:14] is the sub-list for extension extendee
	0,  // [0:0] is the sub-list for field type_name
}

func init() { file_options_proto_init() }
func file_options_proto_init() {
	if File_options_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_options_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RestOpt); i {
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
			RawDescriptor: file_options_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 14,
			NumServices:   0,
		},
		GoTypes:           file_options_proto_goTypes,
		DependencyIndexes: file_options_proto_depIdxs,
		MessageInfos:      file_options_proto_msgTypes,
		ExtensionInfos:    file_options_proto_extTypes,
	}.Build()
	File_options_proto = out.File
	file_options_proto_rawDesc = nil
	file_options_proto_goTypes = nil
	file_options_proto_depIdxs = nil
}
