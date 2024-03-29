// Code generated by mojo. DO NOT EDIT.

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.0
// 	protoc        v3.21.12
// source: ncraft/messaging/push_endpoint.proto

package messaging

import (
	core "github.com/mojo-lang/core/go/pkg/mojo/core"
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

type PushEndpoint struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Service string    `protobuf:"bytes,1,opt,name=service,proto3" json:"service,omitempty"`
	Method  string    `protobuf:"bytes,2,opt,name=method,proto3" json:"method,omitempty"`
	Url     *core.Url `protobuf:"bytes,3,opt,name=url,proto3" json:"url,omitempty"`
}

func (x *PushEndpoint) Reset() {
	*x = PushEndpoint{}
	if protoimpl.UnsafeEnabled {
		mi := &file_ncraft_messaging_push_endpoint_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PushEndpoint) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PushEndpoint) ProtoMessage() {}

func (x *PushEndpoint) ProtoReflect() protoreflect.Message {
	mi := &file_ncraft_messaging_push_endpoint_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PushEndpoint.ProtoReflect.Descriptor instead.
func (*PushEndpoint) Descriptor() ([]byte, []int) {
	return file_ncraft_messaging_push_endpoint_proto_rawDescGZIP(), []int{0}
}

func (x *PushEndpoint) GetService() string {
	if x != nil {
		return x.Service
	}
	return ""
}

func (x *PushEndpoint) GetMethod() string {
	if x != nil {
		return x.Method
	}
	return ""
}

func (x *PushEndpoint) GetUrl() *core.Url {
	if x != nil {
		return x.Url
	}
	return nil
}

var File_ncraft_messaging_push_endpoint_proto protoreflect.FileDescriptor

var file_ncraft_messaging_push_endpoint_proto_rawDesc = []byte{
	0x0a, 0x24, 0x6e, 0x63, 0x72, 0x61, 0x66, 0x74, 0x2f, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x69,
	0x6e, 0x67, 0x2f, 0x70, 0x75, 0x73, 0x68, 0x5f, 0x65, 0x6e, 0x64, 0x70, 0x6f, 0x69, 0x6e, 0x74,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x10, 0x6e, 0x63, 0x72, 0x61, 0x66, 0x74, 0x2e, 0x6d,
	0x65, 0x73, 0x73, 0x61, 0x67, 0x69, 0x6e, 0x67, 0x1a, 0x13, 0x6d, 0x6f, 0x6a, 0x6f, 0x2f, 0x63,
	0x6f, 0x72, 0x65, 0x2f, 0x75, 0x72, 0x6c, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x62, 0x0a,
	0x0c, 0x50, 0x75, 0x73, 0x68, 0x45, 0x6e, 0x64, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x12, 0x18, 0x0a,
	0x07, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07,
	0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x6d, 0x65, 0x74, 0x68, 0x6f,
	0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x6d, 0x65, 0x74, 0x68, 0x6f, 0x64, 0x12,
	0x20, 0x0a, 0x03, 0x75, 0x72, 0x6c, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0e, 0x2e, 0x6d,
	0x6f, 0x6a, 0x6f, 0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x55, 0x72, 0x6c, 0x52, 0x03, 0x75, 0x72,
	0x6c, 0x42, 0x70, 0x0a, 0x1a, 0x69, 0x6f, 0x2e, 0x6e, 0x63, 0x72, 0x61, 0x66, 0x74, 0x2e, 0x6e,
	0x63, 0x72, 0x61, 0x66, 0x74, 0x2e, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x69, 0x6e, 0x67, 0x42,
	0x11, 0x50, 0x75, 0x73, 0x68, 0x45, 0x6e, 0x64, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x50, 0x72, 0x6f,
	0x74, 0x6f, 0x50, 0x01, 0x5a, 0x3d, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d,
	0x2f, 0x6e, 0x63, 0x72, 0x61, 0x66, 0x74, 0x2d, 0x69, 0x6f, 0x2f, 0x6e, 0x63, 0x72, 0x61, 0x66,
	0x74, 0x2f, 0x67, 0x6f, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x6e, 0x63, 0x72, 0x61, 0x66, 0x74, 0x2f,
	0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x69, 0x6e, 0x67, 0x3b, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67,
	0x69, 0x6e, 0x67, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_ncraft_messaging_push_endpoint_proto_rawDescOnce sync.Once
	file_ncraft_messaging_push_endpoint_proto_rawDescData = file_ncraft_messaging_push_endpoint_proto_rawDesc
)

func file_ncraft_messaging_push_endpoint_proto_rawDescGZIP() []byte {
	file_ncraft_messaging_push_endpoint_proto_rawDescOnce.Do(func() {
		file_ncraft_messaging_push_endpoint_proto_rawDescData = protoimpl.X.CompressGZIP(file_ncraft_messaging_push_endpoint_proto_rawDescData)
	})
	return file_ncraft_messaging_push_endpoint_proto_rawDescData
}

var file_ncraft_messaging_push_endpoint_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_ncraft_messaging_push_endpoint_proto_goTypes = []interface{}{
	(*PushEndpoint)(nil), // 0: ncraft.messaging.PushEndpoint
	(*core.Url)(nil),     // 1: mojo.core.Url
}
var file_ncraft_messaging_push_endpoint_proto_depIdxs = []int32{
	1, // 0: ncraft.messaging.PushEndpoint.url:type_name -> mojo.core.Url
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_ncraft_messaging_push_endpoint_proto_init() }
func file_ncraft_messaging_push_endpoint_proto_init() {
	if File_ncraft_messaging_push_endpoint_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_ncraft_messaging_push_endpoint_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PushEndpoint); i {
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
			RawDescriptor: file_ncraft_messaging_push_endpoint_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_ncraft_messaging_push_endpoint_proto_goTypes,
		DependencyIndexes: file_ncraft_messaging_push_endpoint_proto_depIdxs,
		MessageInfos:      file_ncraft_messaging_push_endpoint_proto_msgTypes,
	}.Build()
	File_ncraft_messaging_push_endpoint_proto = out.File
	file_ncraft_messaging_push_endpoint_proto_rawDesc = nil
	file_ncraft_messaging_push_endpoint_proto_goTypes = nil
	file_ncraft_messaging_push_endpoint_proto_depIdxs = nil
}
