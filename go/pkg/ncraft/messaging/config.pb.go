// Code generated by mojo. DO NOT EDIT.

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.0
// 	protoc        v3.21.12
// source: ncraft/messaging/config.proto

package messaging

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

type Config struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Provider      string          `protobuf:"bytes,1,opt,name=provider,proto3" json:"provider,omitempty"`
	Broker        string          `protobuf:"bytes,2,opt,name=broker,proto3" json:"broker,omitempty"`
	ServiceName   string          `protobuf:"bytes,3,opt,name=service_name,json=serviceName,proto3" json:"serviceName,omitempty"`
	Subscriptions []*Subscription `protobuf:"bytes,10,rep,name=subscriptions,proto3" json:"subscriptions,omitempty"`
	Nats          *Nats           `protobuf:"bytes,15,opt,name=nats,proto3" json:"nats,omitempty"`
}

func (x *Config) Reset() {
	*x = Config{}
	if protoimpl.UnsafeEnabled {
		mi := &file_ncraft_messaging_config_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Config) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Config) ProtoMessage() {}

func (x *Config) ProtoReflect() protoreflect.Message {
	mi := &file_ncraft_messaging_config_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Config.ProtoReflect.Descriptor instead.
func (*Config) Descriptor() ([]byte, []int) {
	return file_ncraft_messaging_config_proto_rawDescGZIP(), []int{0}
}

func (x *Config) GetProvider() string {
	if x != nil {
		return x.Provider
	}
	return ""
}

func (x *Config) GetBroker() string {
	if x != nil {
		return x.Broker
	}
	return ""
}

func (x *Config) GetServiceName() string {
	if x != nil {
		return x.ServiceName
	}
	return ""
}

func (x *Config) GetSubscriptions() []*Subscription {
	if x != nil {
		return x.Subscriptions
	}
	return nil
}

func (x *Config) GetNats() *Nats {
	if x != nil {
		return x.Nats
	}
	return nil
}

type Nats struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	JetStream  string   `protobuf:"bytes,1,opt,name=jet_stream,json=jetStream,proto3" json:"jetStream,omitempty"`
	TopicNames []string `protobuf:"bytes,2,rep,name=topic_names,json=topicNames,proto3" json:"topicNames,omitempty"`
	MaxMsgs    int64    `protobuf:"varint,3,opt,name=max_msgs,json=maxMsgs,proto3" json:"maxMsgs,omitempty"`
	MaxAge     int64    `protobuf:"varint,4,opt,name=max_age,json=maxAge,proto3" json:"maxAge,omitempty"`
}

func (x *Nats) Reset() {
	*x = Nats{}
	if protoimpl.UnsafeEnabled {
		mi := &file_ncraft_messaging_config_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Nats) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Nats) ProtoMessage() {}

func (x *Nats) ProtoReflect() protoreflect.Message {
	mi := &file_ncraft_messaging_config_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Nats.ProtoReflect.Descriptor instead.
func (*Nats) Descriptor() ([]byte, []int) {
	return file_ncraft_messaging_config_proto_rawDescGZIP(), []int{1}
}

func (x *Nats) GetJetStream() string {
	if x != nil {
		return x.JetStream
	}
	return ""
}

func (x *Nats) GetTopicNames() []string {
	if x != nil {
		return x.TopicNames
	}
	return nil
}

func (x *Nats) GetMaxMsgs() int64 {
	if x != nil {
		return x.MaxMsgs
	}
	return 0
}

func (x *Nats) GetMaxAge() int64 {
	if x != nil {
		return x.MaxAge
	}
	return 0
}

var File_ncraft_messaging_config_proto protoreflect.FileDescriptor

var file_ncraft_messaging_config_proto_rawDesc = []byte{
	0x0a, 0x1d, 0x6e, 0x63, 0x72, 0x61, 0x66, 0x74, 0x2f, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x69,
	0x6e, 0x67, 0x2f, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12,
	0x10, 0x6e, 0x63, 0x72, 0x61, 0x66, 0x74, 0x2e, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x69, 0x6e,
	0x67, 0x1a, 0x23, 0x6e, 0x63, 0x72, 0x61, 0x66, 0x74, 0x2f, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67,
	0x69, 0x6e, 0x67, 0x2f, 0x73, 0x75, 0x62, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xd1, 0x01, 0x0a, 0x06, 0x43, 0x6f, 0x6e, 0x66, 0x69,
	0x67, 0x12, 0x1a, 0x0a, 0x08, 0x70, 0x72, 0x6f, 0x76, 0x69, 0x64, 0x65, 0x72, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x08, 0x70, 0x72, 0x6f, 0x76, 0x69, 0x64, 0x65, 0x72, 0x12, 0x16, 0x0a,
	0x06, 0x62, 0x72, 0x6f, 0x6b, 0x65, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x62,
	0x72, 0x6f, 0x6b, 0x65, 0x72, 0x12, 0x21, 0x0a, 0x0c, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65,
	0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x73, 0x65, 0x72,
	0x76, 0x69, 0x63, 0x65, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x44, 0x0a, 0x0d, 0x73, 0x75, 0x62, 0x73,
	0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0x0a, 0x20, 0x03, 0x28, 0x0b, 0x32,
	0x1e, 0x2e, 0x6e, 0x63, 0x72, 0x61, 0x66, 0x74, 0x2e, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x69,
	0x6e, 0x67, 0x2e, 0x53, 0x75, 0x62, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x52,
	0x0d, 0x73, 0x75, 0x62, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x12, 0x2a,
	0x0a, 0x04, 0x6e, 0x61, 0x74, 0x73, 0x18, 0x0f, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x16, 0x2e, 0x6e,
	0x63, 0x72, 0x61, 0x66, 0x74, 0x2e, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x69, 0x6e, 0x67, 0x2e,
	0x4e, 0x61, 0x74, 0x73, 0x52, 0x04, 0x6e, 0x61, 0x74, 0x73, 0x22, 0x7a, 0x0a, 0x04, 0x4e, 0x61,
	0x74, 0x73, 0x12, 0x1d, 0x0a, 0x0a, 0x6a, 0x65, 0x74, 0x5f, 0x73, 0x74, 0x72, 0x65, 0x61, 0x6d,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x6a, 0x65, 0x74, 0x53, 0x74, 0x72, 0x65, 0x61,
	0x6d, 0x12, 0x1f, 0x0a, 0x0b, 0x74, 0x6f, 0x70, 0x69, 0x63, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x73,
	0x18, 0x02, 0x20, 0x03, 0x28, 0x09, 0x52, 0x0a, 0x74, 0x6f, 0x70, 0x69, 0x63, 0x4e, 0x61, 0x6d,
	0x65, 0x73, 0x12, 0x19, 0x0a, 0x08, 0x6d, 0x61, 0x78, 0x5f, 0x6d, 0x73, 0x67, 0x73, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x03, 0x52, 0x07, 0x6d, 0x61, 0x78, 0x4d, 0x73, 0x67, 0x73, 0x12, 0x17, 0x0a,
	0x07, 0x6d, 0x61, 0x78, 0x5f, 0x61, 0x67, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x03, 0x52, 0x06,
	0x6d, 0x61, 0x78, 0x41, 0x67, 0x65, 0x42, 0x6a, 0x0a, 0x1a, 0x69, 0x6f, 0x2e, 0x6e, 0x63, 0x72,
	0x61, 0x66, 0x74, 0x2e, 0x6e, 0x63, 0x72, 0x61, 0x66, 0x74, 0x2e, 0x6d, 0x65, 0x73, 0x73, 0x61,
	0x67, 0x69, 0x6e, 0x67, 0x42, 0x0b, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x50, 0x72, 0x6f, 0x74,
	0x6f, 0x50, 0x01, 0x5a, 0x3d, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f,
	0x6e, 0x63, 0x72, 0x61, 0x66, 0x74, 0x2d, 0x69, 0x6f, 0x2f, 0x6e, 0x63, 0x72, 0x61, 0x66, 0x74,
	0x2f, 0x67, 0x6f, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x6e, 0x63, 0x72, 0x61, 0x66, 0x74, 0x2f, 0x6d,
	0x65, 0x73, 0x73, 0x61, 0x67, 0x69, 0x6e, 0x67, 0x3b, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x69,
	0x6e, 0x67, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_ncraft_messaging_config_proto_rawDescOnce sync.Once
	file_ncraft_messaging_config_proto_rawDescData = file_ncraft_messaging_config_proto_rawDesc
)

func file_ncraft_messaging_config_proto_rawDescGZIP() []byte {
	file_ncraft_messaging_config_proto_rawDescOnce.Do(func() {
		file_ncraft_messaging_config_proto_rawDescData = protoimpl.X.CompressGZIP(file_ncraft_messaging_config_proto_rawDescData)
	})
	return file_ncraft_messaging_config_proto_rawDescData
}

var file_ncraft_messaging_config_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_ncraft_messaging_config_proto_goTypes = []interface{}{
	(*Config)(nil),       // 0: ncraft.messaging.Config
	(*Nats)(nil),         // 1: ncraft.messaging.Nats
	(*Subscription)(nil), // 2: ncraft.messaging.Subscription
}
var file_ncraft_messaging_config_proto_depIdxs = []int32{
	2, // 0: ncraft.messaging.Config.subscriptions:type_name -> ncraft.messaging.Subscription
	1, // 1: ncraft.messaging.Config.nats:type_name -> ncraft.messaging.Nats
	2, // [2:2] is the sub-list for method output_type
	2, // [2:2] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_ncraft_messaging_config_proto_init() }
func file_ncraft_messaging_config_proto_init() {
	if File_ncraft_messaging_config_proto != nil {
		return
	}
	file_ncraft_messaging_subscription_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_ncraft_messaging_config_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Config); i {
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
		file_ncraft_messaging_config_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Nats); i {
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
			RawDescriptor: file_ncraft_messaging_config_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_ncraft_messaging_config_proto_goTypes,
		DependencyIndexes: file_ncraft_messaging_config_proto_depIdxs,
		MessageInfos:      file_ncraft_messaging_config_proto_msgTypes,
	}.Build()
	File_ncraft_messaging_config_proto = out.File
	file_ncraft_messaging_config_proto_rawDesc = nil
	file_ncraft_messaging_config_proto_goTypes = nil
	file_ncraft_messaging_config_proto_depIdxs = nil
}
