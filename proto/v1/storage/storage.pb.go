// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.0
// 	protoc        v3.19.4
// source: proto/v1/storage/storage.proto

package storage_v1

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

type PingRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *PingRequest) Reset() {
	*x = PingRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_v1_storage_storage_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PingRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PingRequest) ProtoMessage() {}

func (x *PingRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_v1_storage_storage_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PingRequest.ProtoReflect.Descriptor instead.
func (*PingRequest) Descriptor() ([]byte, []int) {
	return file_proto_v1_storage_storage_proto_rawDescGZIP(), []int{0}
}

type PingResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Msg string `protobuf:"bytes,1,opt,name=msg,proto3" json:"msg,omitempty"`
}

func (x *PingResponse) Reset() {
	*x = PingResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_v1_storage_storage_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PingResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PingResponse) ProtoMessage() {}

func (x *PingResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_v1_storage_storage_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PingResponse.ProtoReflect.Descriptor instead.
func (*PingResponse) Descriptor() ([]byte, []int) {
	return file_proto_v1_storage_storage_proto_rawDescGZIP(), []int{1}
}

func (x *PingResponse) GetMsg() string {
	if x != nil {
		return x.Msg
	}
	return ""
}

type CreateIngestionJobRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	JobID string `protobuf:"bytes,1,opt,name=jobID,proto3" json:"jobID,omitempty"`
}

func (x *CreateIngestionJobRequest) Reset() {
	*x = CreateIngestionJobRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_v1_storage_storage_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateIngestionJobRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateIngestionJobRequest) ProtoMessage() {}

func (x *CreateIngestionJobRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_v1_storage_storage_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateIngestionJobRequest.ProtoReflect.Descriptor instead.
func (*CreateIngestionJobRequest) Descriptor() ([]byte, []int) {
	return file_proto_v1_storage_storage_proto_rawDescGZIP(), []int{2}
}

func (x *CreateIngestionJobRequest) GetJobID() string {
	if x != nil {
		return x.JobID
	}
	return ""
}

type CreateIngestionJobResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	StorageID string `protobuf:"bytes,1,opt,name=storageID,proto3" json:"storageID,omitempty"`
	Message   string `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
}

func (x *CreateIngestionJobResponse) Reset() {
	*x = CreateIngestionJobResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_v1_storage_storage_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateIngestionJobResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateIngestionJobResponse) ProtoMessage() {}

func (x *CreateIngestionJobResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_v1_storage_storage_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateIngestionJobResponse.ProtoReflect.Descriptor instead.
func (*CreateIngestionJobResponse) Descriptor() ([]byte, []int) {
	return file_proto_v1_storage_storage_proto_rawDescGZIP(), []int{3}
}

func (x *CreateIngestionJobResponse) GetStorageID() string {
	if x != nil {
		return x.StorageID
	}
	return ""
}

func (x *CreateIngestionJobResponse) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

type QueryTableRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	BrokerID string `protobuf:"bytes,1,opt,name=brokerID,proto3" json:"brokerID,omitempty"`
	Table    string `protobuf:"bytes,2,opt,name=table,proto3" json:"table,omitempty"`
	Query    string `protobuf:"bytes,3,opt,name=query,proto3" json:"query,omitempty"`
}

func (x *QueryTableRequest) Reset() {
	*x = QueryTableRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_v1_storage_storage_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *QueryTableRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*QueryTableRequest) ProtoMessage() {}

func (x *QueryTableRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_v1_storage_storage_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use QueryTableRequest.ProtoReflect.Descriptor instead.
func (*QueryTableRequest) Descriptor() ([]byte, []int) {
	return file_proto_v1_storage_storage_proto_rawDescGZIP(), []int{4}
}

func (x *QueryTableRequest) GetBrokerID() string {
	if x != nil {
		return x.BrokerID
	}
	return ""
}

func (x *QueryTableRequest) GetTable() string {
	if x != nil {
		return x.Table
	}
	return ""
}

func (x *QueryTableRequest) GetQuery() string {
	if x != nil {
		return x.Query
	}
	return ""
}

type QueryTableResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	BrokerID string `protobuf:"bytes,1,opt,name=brokerID,proto3" json:"brokerID,omitempty"`
	Result   []byte `protobuf:"bytes,2,opt,name=result,proto3" json:"result,omitempty"`
}

func (x *QueryTableResponse) Reset() {
	*x = QueryTableResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_v1_storage_storage_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *QueryTableResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*QueryTableResponse) ProtoMessage() {}

func (x *QueryTableResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_v1_storage_storage_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use QueryTableResponse.ProtoReflect.Descriptor instead.
func (*QueryTableResponse) Descriptor() ([]byte, []int) {
	return file_proto_v1_storage_storage_proto_rawDescGZIP(), []int{5}
}

func (x *QueryTableResponse) GetBrokerID() string {
	if x != nil {
		return x.BrokerID
	}
	return ""
}

func (x *QueryTableResponse) GetResult() []byte {
	if x != nil {
		return x.Result
	}
	return nil
}

var File_proto_v1_storage_storage_proto protoreflect.FileDescriptor

var file_proto_v1_storage_storage_proto_rawDesc = []byte{
	0x0a, 0x1e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x76, 0x31, 0x2f, 0x73, 0x74, 0x6f, 0x72, 0x61,
	0x67, 0x65, 0x2f, 0x73, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x12, 0x0a, 0x73, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x2e, 0x76, 0x31, 0x22, 0x0d, 0x0a, 0x0b,
	0x50, 0x69, 0x6e, 0x67, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x22, 0x20, 0x0a, 0x0c, 0x50,
	0x69, 0x6e, 0x67, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x10, 0x0a, 0x03, 0x6d,
	0x73, 0x67, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6d, 0x73, 0x67, 0x22, 0x31, 0x0a,
	0x19, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x49, 0x6e, 0x67, 0x65, 0x73, 0x74, 0x69, 0x6f, 0x6e,
	0x4a, 0x6f, 0x62, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x14, 0x0a, 0x05, 0x6a, 0x6f,
	0x62, 0x49, 0x44, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x6a, 0x6f, 0x62, 0x49, 0x44,
	0x22, 0x54, 0x0a, 0x1a, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x49, 0x6e, 0x67, 0x65, 0x73, 0x74,
	0x69, 0x6f, 0x6e, 0x4a, 0x6f, 0x62, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x1c,
	0x0a, 0x09, 0x73, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x49, 0x44, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x09, 0x73, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x49, 0x44, 0x12, 0x18, 0x0a, 0x07,
	0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6d,
	0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x22, 0x5b, 0x0a, 0x11, 0x51, 0x75, 0x65, 0x72, 0x79, 0x54,
	0x61, 0x62, 0x6c, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1a, 0x0a, 0x08, 0x62,
	0x72, 0x6f, 0x6b, 0x65, 0x72, 0x49, 0x44, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x62,
	0x72, 0x6f, 0x6b, 0x65, 0x72, 0x49, 0x44, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x61, 0x62, 0x6c, 0x65,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x74, 0x61, 0x62, 0x6c, 0x65, 0x12, 0x14, 0x0a,
	0x05, 0x71, 0x75, 0x65, 0x72, 0x79, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x71, 0x75,
	0x65, 0x72, 0x79, 0x22, 0x48, 0x0a, 0x12, 0x51, 0x75, 0x65, 0x72, 0x79, 0x54, 0x61, 0x62, 0x6c,
	0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x62, 0x72, 0x6f,
	0x6b, 0x65, 0x72, 0x49, 0x44, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x62, 0x72, 0x6f,
	0x6b, 0x65, 0x72, 0x49, 0x44, 0x12, 0x16, 0x0a, 0x06, 0x72, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x06, 0x72, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x32, 0xfa, 0x01,
	0x0a, 0x07, 0x53, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x12, 0x39, 0x0a, 0x04, 0x50, 0x69, 0x6e,
	0x67, 0x12, 0x17, 0x2e, 0x73, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x50,
	0x69, 0x6e, 0x67, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x18, 0x2e, 0x73, 0x74, 0x6f,
	0x72, 0x61, 0x67, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x50, 0x69, 0x6e, 0x67, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x12, 0x65, 0x0a, 0x12, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x49, 0x6e,
	0x67, 0x65, 0x73, 0x74, 0x69, 0x6f, 0x6e, 0x4a, 0x6f, 0x62, 0x12, 0x25, 0x2e, 0x73, 0x74, 0x6f,
	0x72, 0x61, 0x67, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x49, 0x6e,
	0x67, 0x65, 0x73, 0x74, 0x69, 0x6f, 0x6e, 0x4a, 0x6f, 0x62, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x1a, 0x26, 0x2e, 0x73, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x43,
	0x72, 0x65, 0x61, 0x74, 0x65, 0x49, 0x6e, 0x67, 0x65, 0x73, 0x74, 0x69, 0x6f, 0x6e, 0x4a, 0x6f,
	0x62, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x4d, 0x0a, 0x0a, 0x51,
	0x75, 0x65, 0x72, 0x79, 0x54, 0x61, 0x62, 0x6c, 0x65, 0x12, 0x1d, 0x2e, 0x73, 0x74, 0x6f, 0x72,
	0x61, 0x67, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x51, 0x75, 0x65, 0x72, 0x79, 0x54, 0x61, 0x62, 0x6c,
	0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1e, 0x2e, 0x73, 0x74, 0x6f, 0x72, 0x61,
	0x67, 0x65, 0x2e, 0x76, 0x31, 0x2e, 0x51, 0x75, 0x65, 0x72, 0x79, 0x54, 0x61, 0x62, 0x6c, 0x65,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x42, 0x3a, 0x5a, 0x38, 0x67, 0x69,
	0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x73, 0x70, 0x61, 0x67, 0x68, 0x65, 0x74,
	0x74, 0x69, 0x66, 0x75, 0x6e, 0x6b, 0x2f, 0x6e, 0x6f, 0x72, 0x6d, 0x61, 0x6e, 0x2f, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x2f, 0x73, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x2f, 0x73, 0x74, 0x6f, 0x72,
	0x61, 0x67, 0x65, 0x5f, 0x76, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_proto_v1_storage_storage_proto_rawDescOnce sync.Once
	file_proto_v1_storage_storage_proto_rawDescData = file_proto_v1_storage_storage_proto_rawDesc
)

func file_proto_v1_storage_storage_proto_rawDescGZIP() []byte {
	file_proto_v1_storage_storage_proto_rawDescOnce.Do(func() {
		file_proto_v1_storage_storage_proto_rawDescData = protoimpl.X.CompressGZIP(file_proto_v1_storage_storage_proto_rawDescData)
	})
	return file_proto_v1_storage_storage_proto_rawDescData
}

var file_proto_v1_storage_storage_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_proto_v1_storage_storage_proto_goTypes = []interface{}{
	(*PingRequest)(nil),                // 0: storage.v1.PingRequest
	(*PingResponse)(nil),               // 1: storage.v1.PingResponse
	(*CreateIngestionJobRequest)(nil),  // 2: storage.v1.CreateIngestionJobRequest
	(*CreateIngestionJobResponse)(nil), // 3: storage.v1.CreateIngestionJobResponse
	(*QueryTableRequest)(nil),          // 4: storage.v1.QueryTableRequest
	(*QueryTableResponse)(nil),         // 5: storage.v1.QueryTableResponse
}
var file_proto_v1_storage_storage_proto_depIdxs = []int32{
	0, // 0: storage.v1.Storage.Ping:input_type -> storage.v1.PingRequest
	2, // 1: storage.v1.Storage.CreateIngestionJob:input_type -> storage.v1.CreateIngestionJobRequest
	4, // 2: storage.v1.Storage.QueryTable:input_type -> storage.v1.QueryTableRequest
	1, // 3: storage.v1.Storage.Ping:output_type -> storage.v1.PingResponse
	3, // 4: storage.v1.Storage.CreateIngestionJob:output_type -> storage.v1.CreateIngestionJobResponse
	5, // 5: storage.v1.Storage.QueryTable:output_type -> storage.v1.QueryTableResponse
	3, // [3:6] is the sub-list for method output_type
	0, // [0:3] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_proto_v1_storage_storage_proto_init() }
func file_proto_v1_storage_storage_proto_init() {
	if File_proto_v1_storage_storage_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_proto_v1_storage_storage_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PingRequest); i {
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
		file_proto_v1_storage_storage_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PingResponse); i {
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
		file_proto_v1_storage_storage_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CreateIngestionJobRequest); i {
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
		file_proto_v1_storage_storage_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CreateIngestionJobResponse); i {
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
		file_proto_v1_storage_storage_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*QueryTableRequest); i {
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
		file_proto_v1_storage_storage_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*QueryTableResponse); i {
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
			RawDescriptor: file_proto_v1_storage_storage_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_proto_v1_storage_storage_proto_goTypes,
		DependencyIndexes: file_proto_v1_storage_storage_proto_depIdxs,
		MessageInfos:      file_proto_v1_storage_storage_proto_msgTypes,
	}.Build()
	File_proto_v1_storage_storage_proto = out.File
	file_proto_v1_storage_storage_proto_rawDesc = nil
	file_proto_v1_storage_storage_proto_goTypes = nil
	file_proto_v1_storage_storage_proto_depIdxs = nil
}