// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        v3.17.3
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
	// bind_path_variables is a list of pair of bindings. For example if we have a path '/part1/:var1/:var2'
	// and we want to bind 'var1' to field 'fieldVar1': "var1=fieldVar1"
	BindPathParam []string `protobuf:"bytes,3,rep,name=bind_path_param,json=bindPathParam" json:"bind_path_param,omitempty"`
	// bind_query_param is a list of pair of bindings. For example if we have a query param var1
	// and we want to bind 'var1' to field 'fieldVar1': "var1=fieldVar1"
	BindQueryParam []string `protobuf:"bytes,4,rep,name=bind_query_param,json=bindQueryParam" json:"bind_query_param,omitempty"`
	// json_encode if is set then input and output data are json encoded version of the proto messages
	JsonEncode *bool `protobuf:"varint,5,opt,name=json_encode,json=jsonEncode" json:"json_encode,omitempty"`
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

func (x *RestOpt) GetBindPathParam() []string {
	if x != nil {
		return x.BindPathParam
	}
	return nil
}

func (x *RestOpt) GetBindQueryParam() []string {
	if x != nil {
		return x.BindQueryParam
	}
	return nil
}

func (x *RestOpt) GetJsonEncode() bool {
	if x != nil && x.JsonEncode != nil {
		return *x.JsonEncode
	}
	return false
}

// ModelOpt holds all the information if you want to create a model from your proto message.
type ModelOpt struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// singleton marks this model as a singleton.
	Singleton *bool `protobuf:"varint,1,opt,name=singleton" json:"singleton,omitempty"`
	// global_datasource generates the code for remote repository
	// POSSIBLE VALUES: cql, sql
	GlobalDatasource *string `protobuf:"bytes,2,opt,name=global_datasource,json=globalDatasource" json:"global_datasource,omitempty"`
	// local_datasource generates the code for local repository
	// POSSIBLE VALUES: store, cql, sql
	LocalDatasource *string `protobuf:"bytes,3,opt,name=local_datasource,json=localDatasource" json:"local_datasource,omitempty"`
	// table a virtual table presentation to hold instances of this message, like rows in a table
	// PRIMARY KEY FORMAT: ( (partitionKey1, partitionKey2, ...), clusteringKey1, clusteringKey2, ...)
	Table *PrimaryKeyOpt `protobuf:"bytes,4,opt,name=table" json:"table,omitempty"`
	// view creates a materialized view of the aggregate based on the primary key.
	// PRIMARY KEY FORMAT: ( (partitionKey1, partitionKey2, ...), clusteringKey1, clusteringKey2, ...)
	// NOTE: The primary key of the view must contains all the primary key items of the table. They don't need to
	//           follow the same order as table. for example the following is correct:
	//                  table = ((a, b), c)
	//                  view = ((c, a),  d, b)
	View []*PrimaryKeyOpt `protobuf:"bytes,5,rep,name=view" json:"view,omitempty"`
}

func (x *ModelOpt) Reset() {
	*x = ModelOpt{}
	if protoimpl.UnsafeEnabled {
		mi := &file_options_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ModelOpt) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ModelOpt) ProtoMessage() {}

func (x *ModelOpt) ProtoReflect() protoreflect.Message {
	mi := &file_options_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ModelOpt.ProtoReflect.Descriptor instead.
func (*ModelOpt) Descriptor() ([]byte, []int) {
	return file_options_proto_rawDescGZIP(), []int{1}
}

func (x *ModelOpt) GetSingleton() bool {
	if x != nil && x.Singleton != nil {
		return *x.Singleton
	}
	return false
}

func (x *ModelOpt) GetGlobalDatasource() string {
	if x != nil && x.GlobalDatasource != nil {
		return *x.GlobalDatasource
	}
	return ""
}

func (x *ModelOpt) GetLocalDatasource() string {
	if x != nil && x.LocalDatasource != nil {
		return *x.LocalDatasource
	}
	return ""
}

func (x *ModelOpt) GetTable() *PrimaryKeyOpt {
	if x != nil {
		return x.Table
	}
	return nil
}

func (x *ModelOpt) GetView() []*PrimaryKeyOpt {
	if x != nil {
		return x.View
	}
	return nil
}

// PrimaryKeyOpt is the container for aggregate primary key settings. It is going to replace old
// rony_aggregate_table & rony_aggregate_view options.
type PrimaryKeyOpt struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	PartKey []string `protobuf:"bytes,1,rep,name=part_key,json=partKey" json:"part_key,omitempty"`
	SortKey []string `protobuf:"bytes,2,rep,name=sort_key,json=sortKey" json:"sort_key,omitempty"`
	// Alias could be used to override the name of materialized view. If this is not set, Rony automatically
	// generates for you.
	Alias *string `protobuf:"bytes,4,opt,name=alias" json:"alias,omitempty"`
}

func (x *PrimaryKeyOpt) Reset() {
	*x = PrimaryKeyOpt{}
	if protoimpl.UnsafeEnabled {
		mi := &file_options_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PrimaryKeyOpt) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PrimaryKeyOpt) ProtoMessage() {}

func (x *PrimaryKeyOpt) ProtoReflect() protoreflect.Message {
	mi := &file_options_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PrimaryKeyOpt.ProtoReflect.Descriptor instead.
func (*PrimaryKeyOpt) Descriptor() ([]byte, []int) {
	return file_options_proto_rawDescGZIP(), []int{2}
}

func (x *PrimaryKeyOpt) GetPartKey() []string {
	if x != nil {
		return x.PartKey
	}
	return nil
}

func (x *PrimaryKeyOpt) GetSortKey() []string {
	if x != nil {
		return x.SortKey
	}
	return nil
}

func (x *PrimaryKeyOpt) GetAlias() string {
	if x != nil && x.Alias != nil {
		return *x.Alias
	}
	return ""
}

type CliOpt struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// protocol defines what protocol should client use to communicate with server.
	// POSSIBLE VALUES: "ws", "http"
	Protocol *string `protobuf:"bytes,1,opt,name=protocol" json:"protocol,omitempty"`
}

func (x *CliOpt) Reset() {
	*x = CliOpt{}
	if protoimpl.UnsafeEnabled {
		mi := &file_options_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CliOpt) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CliOpt) ProtoMessage() {}

func (x *CliOpt) ProtoReflect() protoreflect.Message {
	mi := &file_options_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CliOpt.ProtoReflect.Descriptor instead.
func (*CliOpt) Descriptor() ([]byte, []int) {
	return file_options_proto_rawDescGZIP(), []int{3}
}

func (x *CliOpt) GetProtocol() string {
	if x != nil && x.Protocol != nil {
		return *x.Protocol
	}
	return ""
}

var file_options_proto_extTypes = []protoimpl.ExtensionInfo{
	{
		ExtendedType:  (*descriptorpb.ServiceOptions)(nil),
		ExtensionType: (*CliOpt)(nil),
		Field:         50001,
		Name:          "rony_cli",
		Tag:           "bytes,50001,opt,name=rony_cli",
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
		ExtensionType: (*ModelOpt)(nil),
		Field:         50001,
		Name:          "rony_model",
		Tag:           "bytes,50001,opt,name=rony_model",
		Filename:      "options.proto",
	},
	{
		ExtendedType:  (*descriptorpb.MessageOptions)(nil),
		ExtensionType: (*bool)(nil),
		Field:         50031,
		Name:          "rony_envelope",
		Tag:           "varint,50031,opt,name=rony_envelope",
		Filename:      "options.proto",
	},
	{
		ExtendedType:  (*descriptorpb.MessageOptions)(nil),
		ExtensionType: (*bool)(nil),
		Field:         50032,
		Name:          "rony_skip_json",
		Tag:           "varint,50032,opt,name=rony_skip_json",
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
	{
		ExtendedType:  (*descriptorpb.FieldOptions)(nil),
		ExtensionType: (*string)(nil),
		Field:         50002,
		Name:          "rony_help",
		Tag:           "bytes,50002,opt,name=rony_help",
		Filename:      "options.proto",
	},
	{
		ExtendedType:  (*descriptorpb.FieldOptions)(nil),
		ExtensionType: (*string)(nil),
		Field:         50003,
		Name:          "rony_default",
		Tag:           "bytes,50003,opt,name=rony_default",
		Filename:      "options.proto",
	},
}

// Extension fields to descriptorpb.ServiceOptions.
var (
	// rony_cli generates the boiler plate code for client stub of rpc methods.
	//
	// optional CliOpt rony_cli = 50001;
	E_RonyCli = &file_options_proto_extTypes[0]
	// rony_no_client if is set then no client code will be generated. This flag is for internal usage.
	// DO NOT USE IT.
	//
	// optional bool rony_no_client = 50003;
	E_RonyNoClient = &file_options_proto_extTypes[1]
)

// Extension fields to descriptorpb.MethodOptions.
var (
	// rony_internal marks this method internal, hence only edges could execute this rpc through tunnel messages. In other words,
	// this command is not exposed to external clients connected through th gateway.
	//
	// optional bool rony_internal = 50002;
	E_RonyInternal = &file_options_proto_extTypes[2]
	// optional RestOpt rony_rest = 50003;
	E_RonyRest = &file_options_proto_extTypes[3]
)

// Extension fields to descriptorpb.MessageOptions.
var (
	// optional ModelOpt rony_model = 50001;
	E_RonyModel = &file_options_proto_extTypes[4]
	// optional bool rony_envelope = 50031;
	E_RonyEnvelope = &file_options_proto_extTypes[5]
	// rony_skip_json will not generate de/serializer for JSON. If you use this option, then you need to
	// write them manually otherwise you will get compiler error.
	//
	// optional bool rony_skip_json = 50032;
	E_RonySkipJson = &file_options_proto_extTypes[6]
)

// Extension fields to descriptorpb.FieldOptions.
var (
	// rony_index marks this field as an indexed field. Some queries will be generated for this indexed field.
	//
	// optional bool rony_index = 50001;
	E_RonyIndex = &file_options_proto_extTypes[7]
	// rony_default set the help text in generated cli client. This flag only works if rony_cobra_cmd is set
	// in the service.
	//
	// optional string rony_help = 50002;
	E_RonyHelp = &file_options_proto_extTypes[8]
	// rony_default set the default value in generated cli client. This flag only works if rony_cobra_cmd is set
	// in the service.
	//
	// optional string rony_default = 50003;
	E_RonyDefault = &file_options_proto_extTypes[9]
)

var File_options_proto protoreflect.FileDescriptor

var file_options_proto_rawDesc = []byte{
	0x0a, 0x0d, 0x6f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a,
	0x20, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2f, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x6f, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x22, 0xa8, 0x01, 0x0a, 0x07, 0x52, 0x65, 0x73, 0x74, 0x4f, 0x70, 0x74, 0x12, 0x16, 0x0a,
	0x06, 0x6d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x18, 0x01, 0x20, 0x02, 0x28, 0x09, 0x52, 0x06, 0x6d,
	0x65, 0x74, 0x68, 0x6f, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x70, 0x61, 0x74, 0x68, 0x18, 0x02, 0x20,
	0x02, 0x28, 0x09, 0x52, 0x04, 0x70, 0x61, 0x74, 0x68, 0x12, 0x26, 0x0a, 0x0f, 0x62, 0x69, 0x6e,
	0x64, 0x5f, 0x70, 0x61, 0x74, 0x68, 0x5f, 0x70, 0x61, 0x72, 0x61, 0x6d, 0x18, 0x03, 0x20, 0x03,
	0x28, 0x09, 0x52, 0x0d, 0x62, 0x69, 0x6e, 0x64, 0x50, 0x61, 0x74, 0x68, 0x50, 0x61, 0x72, 0x61,
	0x6d, 0x12, 0x28, 0x0a, 0x10, 0x62, 0x69, 0x6e, 0x64, 0x5f, 0x71, 0x75, 0x65, 0x72, 0x79, 0x5f,
	0x70, 0x61, 0x72, 0x61, 0x6d, 0x18, 0x04, 0x20, 0x03, 0x28, 0x09, 0x52, 0x0e, 0x62, 0x69, 0x6e,
	0x64, 0x51, 0x75, 0x65, 0x72, 0x79, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x12, 0x1f, 0x0a, 0x0b, 0x6a,
	0x73, 0x6f, 0x6e, 0x5f, 0x65, 0x6e, 0x63, 0x6f, 0x64, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x08,
	0x52, 0x0a, 0x6a, 0x73, 0x6f, 0x6e, 0x45, 0x6e, 0x63, 0x6f, 0x64, 0x65, 0x22, 0xca, 0x01, 0x0a,
	0x08, 0x4d, 0x6f, 0x64, 0x65, 0x6c, 0x4f, 0x70, 0x74, 0x12, 0x1c, 0x0a, 0x09, 0x73, 0x69, 0x6e,
	0x67, 0x6c, 0x65, 0x74, 0x6f, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x09, 0x73, 0x69,
	0x6e, 0x67, 0x6c, 0x65, 0x74, 0x6f, 0x6e, 0x12, 0x2b, 0x0a, 0x11, 0x67, 0x6c, 0x6f, 0x62, 0x61,
	0x6c, 0x5f, 0x64, 0x61, 0x74, 0x61, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x10, 0x67, 0x6c, 0x6f, 0x62, 0x61, 0x6c, 0x44, 0x61, 0x74, 0x61, 0x73, 0x6f,
	0x75, 0x72, 0x63, 0x65, 0x12, 0x29, 0x0a, 0x10, 0x6c, 0x6f, 0x63, 0x61, 0x6c, 0x5f, 0x64, 0x61,
	0x74, 0x61, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0f,
	0x6c, 0x6f, 0x63, 0x61, 0x6c, 0x44, 0x61, 0x74, 0x61, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x12,
	0x24, 0x0a, 0x05, 0x74, 0x61, 0x62, 0x6c, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0e,
	0x2e, 0x50, 0x72, 0x69, 0x6d, 0x61, 0x72, 0x79, 0x4b, 0x65, 0x79, 0x4f, 0x70, 0x74, 0x52, 0x05,
	0x74, 0x61, 0x62, 0x6c, 0x65, 0x12, 0x22, 0x0a, 0x04, 0x76, 0x69, 0x65, 0x77, 0x18, 0x05, 0x20,
	0x03, 0x28, 0x0b, 0x32, 0x0e, 0x2e, 0x50, 0x72, 0x69, 0x6d, 0x61, 0x72, 0x79, 0x4b, 0x65, 0x79,
	0x4f, 0x70, 0x74, 0x52, 0x04, 0x76, 0x69, 0x65, 0x77, 0x22, 0x5b, 0x0a, 0x0d, 0x50, 0x72, 0x69,
	0x6d, 0x61, 0x72, 0x79, 0x4b, 0x65, 0x79, 0x4f, 0x70, 0x74, 0x12, 0x19, 0x0a, 0x08, 0x70, 0x61,
	0x72, 0x74, 0x5f, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x03, 0x28, 0x09, 0x52, 0x07, 0x70, 0x61,
	0x72, 0x74, 0x4b, 0x65, 0x79, 0x12, 0x19, 0x0a, 0x08, 0x73, 0x6f, 0x72, 0x74, 0x5f, 0x6b, 0x65,
	0x79, 0x18, 0x02, 0x20, 0x03, 0x28, 0x09, 0x52, 0x07, 0x73, 0x6f, 0x72, 0x74, 0x4b, 0x65, 0x79,
	0x12, 0x14, 0x0a, 0x05, 0x61, 0x6c, 0x69, 0x61, 0x73, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x05, 0x61, 0x6c, 0x69, 0x61, 0x73, 0x22, 0x24, 0x0a, 0x06, 0x43, 0x6c, 0x69, 0x4f, 0x70, 0x74,
	0x12, 0x1a, 0x0a, 0x08, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6c, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x08, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6c, 0x3a, 0x45, 0x0a, 0x08,
	0x72, 0x6f, 0x6e, 0x79, 0x5f, 0x63, 0x6c, 0x69, 0x12, 0x1f, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c,
	0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x53, 0x65, 0x72, 0x76, 0x69,
	0x63, 0x65, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0xd1, 0x86, 0x03, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x07, 0x2e, 0x43, 0x6c, 0x69, 0x4f, 0x70, 0x74, 0x52, 0x07, 0x72, 0x6f, 0x6e, 0x79,
	0x43, 0x6c, 0x69, 0x3a, 0x47, 0x0a, 0x0e, 0x72, 0x6f, 0x6e, 0x79, 0x5f, 0x6e, 0x6f, 0x5f, 0x63,
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
	0x70, 0x74, 0x52, 0x08, 0x72, 0x6f, 0x6e, 0x79, 0x52, 0x65, 0x73, 0x74, 0x3a, 0x4b, 0x0a, 0x0a,
	0x72, 0x6f, 0x6e, 0x79, 0x5f, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x12, 0x1f, 0x2e, 0x67, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x4d, 0x65, 0x73,
	0x73, 0x61, 0x67, 0x65, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0xd1, 0x86, 0x03, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x09, 0x2e, 0x4d, 0x6f, 0x64, 0x65, 0x6c, 0x4f, 0x70, 0x74, 0x52, 0x09,
	0x72, 0x6f, 0x6e, 0x79, 0x4d, 0x6f, 0x64, 0x65, 0x6c, 0x3a, 0x46, 0x0a, 0x0d, 0x72, 0x6f, 0x6e,
	0x79, 0x5f, 0x65, 0x6e, 0x76, 0x65, 0x6c, 0x6f, 0x70, 0x65, 0x12, 0x1f, 0x2e, 0x67, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x4d, 0x65, 0x73,
	0x73, 0x61, 0x67, 0x65, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0xef, 0x86, 0x03, 0x20,
	0x01, 0x28, 0x08, 0x52, 0x0c, 0x72, 0x6f, 0x6e, 0x79, 0x45, 0x6e, 0x76, 0x65, 0x6c, 0x6f, 0x70,
	0x65, 0x3a, 0x47, 0x0a, 0x0e, 0x72, 0x6f, 0x6e, 0x79, 0x5f, 0x73, 0x6b, 0x69, 0x70, 0x5f, 0x6a,
	0x73, 0x6f, 0x6e, 0x12, 0x1f, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x4f, 0x70, 0x74,
	0x69, 0x6f, 0x6e, 0x73, 0x18, 0xf0, 0x86, 0x03, 0x20, 0x01, 0x28, 0x08, 0x52, 0x0c, 0x72, 0x6f,
	0x6e, 0x79, 0x53, 0x6b, 0x69, 0x70, 0x4a, 0x73, 0x6f, 0x6e, 0x3a, 0x3e, 0x0a, 0x0a, 0x72, 0x6f,
	0x6e, 0x79, 0x5f, 0x69, 0x6e, 0x64, 0x65, 0x78, 0x12, 0x1d, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c,
	0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x46, 0x69, 0x65, 0x6c, 0x64,
	0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0xd1, 0x86, 0x03, 0x20, 0x01, 0x28, 0x08, 0x52,
	0x09, 0x72, 0x6f, 0x6e, 0x79, 0x49, 0x6e, 0x64, 0x65, 0x78, 0x3a, 0x3c, 0x0a, 0x09, 0x72, 0x6f,
	0x6e, 0x79, 0x5f, 0x68, 0x65, 0x6c, 0x70, 0x12, 0x1d, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x4f,
	0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0xd2, 0x86, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08,
	0x72, 0x6f, 0x6e, 0x79, 0x48, 0x65, 0x6c, 0x70, 0x3a, 0x42, 0x0a, 0x0c, 0x72, 0x6f, 0x6e, 0x79,
	0x5f, 0x64, 0x65, 0x66, 0x61, 0x75, 0x6c, 0x74, 0x12, 0x1d, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c,
	0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x46, 0x69, 0x65, 0x6c, 0x64,
	0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0xd3, 0x86, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x0b, 0x72, 0x6f, 0x6e, 0x79, 0x44, 0x65, 0x66, 0x61, 0x75, 0x6c, 0x74, 0x42, 0x1b, 0x5a, 0x19,
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

var file_options_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_options_proto_goTypes = []interface{}{
	(*RestOpt)(nil),                     // 0: RestOpt
	(*ModelOpt)(nil),                    // 1: ModelOpt
	(*PrimaryKeyOpt)(nil),               // 2: PrimaryKeyOpt
	(*CliOpt)(nil),                      // 3: CliOpt
	(*descriptorpb.ServiceOptions)(nil), // 4: google.protobuf.ServiceOptions
	(*descriptorpb.MethodOptions)(nil),  // 5: google.protobuf.MethodOptions
	(*descriptorpb.MessageOptions)(nil), // 6: google.protobuf.MessageOptions
	(*descriptorpb.FieldOptions)(nil),   // 7: google.protobuf.FieldOptions
}
var file_options_proto_depIdxs = []int32{
	2,  // 0: ModelOpt.table:type_name -> PrimaryKeyOpt
	2,  // 1: ModelOpt.view:type_name -> PrimaryKeyOpt
	4,  // 2: rony_cli:extendee -> google.protobuf.ServiceOptions
	4,  // 3: rony_no_client:extendee -> google.protobuf.ServiceOptions
	5,  // 4: rony_internal:extendee -> google.protobuf.MethodOptions
	5,  // 5: rony_rest:extendee -> google.protobuf.MethodOptions
	6,  // 6: rony_model:extendee -> google.protobuf.MessageOptions
	6,  // 7: rony_envelope:extendee -> google.protobuf.MessageOptions
	6,  // 8: rony_skip_json:extendee -> google.protobuf.MessageOptions
	7,  // 9: rony_index:extendee -> google.protobuf.FieldOptions
	7,  // 10: rony_help:extendee -> google.protobuf.FieldOptions
	7,  // 11: rony_default:extendee -> google.protobuf.FieldOptions
	3,  // 12: rony_cli:type_name -> CliOpt
	0,  // 13: rony_rest:type_name -> RestOpt
	1,  // 14: rony_model:type_name -> ModelOpt
	15, // [15:15] is the sub-list for method output_type
	15, // [15:15] is the sub-list for method input_type
	12, // [12:15] is the sub-list for extension type_name
	2,  // [2:12] is the sub-list for extension extendee
	0,  // [0:2] is the sub-list for field type_name
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
		file_options_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ModelOpt); i {
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
		file_options_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PrimaryKeyOpt); i {
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
		file_options_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CliOpt); i {
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
			NumMessages:   4,
			NumExtensions: 10,
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
