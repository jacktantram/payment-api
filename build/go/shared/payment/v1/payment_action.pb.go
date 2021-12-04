// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        v3.18.1
// source: shared/payment/v1/payment_action.proto

package v1

import (
	v1 "github.com/jacktantram/build/go/shared/amount/v1"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// The type of the payment.
type PaymentType int32

const (
	// The payment type is unspecified. This should not happen.
	PaymentType_PAYMENT_TYPE_UNSPECIFIED PaymentType = 0
	// The payment type is a capture type
	PaymentType_PAYMENT_TYPE_CAPTURE PaymentType = 1
	// The payment type is a refund type
	PaymentType_PAYMENT_TYPE_REFUND PaymentType = 2
	// The payment type is a void type
	PaymentType_PAYMENT_TYPE_VOID PaymentType = 3
)

// Enum value maps for PaymentType.
var (
	PaymentType_name = map[int32]string{
		0: "PAYMENT_TYPE_UNSPECIFIED",
		1: "PAYMENT_TYPE_CAPTURE",
		2: "PAYMENT_TYPE_REFUND",
		3: "PAYMENT_TYPE_VOID",
	}
	PaymentType_value = map[string]int32{
		"PAYMENT_TYPE_UNSPECIFIED": 0,
		"PAYMENT_TYPE_CAPTURE":     1,
		"PAYMENT_TYPE_REFUND":      2,
		"PAYMENT_TYPE_VOID":        3,
	}
)

func (x PaymentType) Enum() *PaymentType {
	p := new(PaymentType)
	*p = x
	return p
}

func (x PaymentType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (PaymentType) Descriptor() protoreflect.EnumDescriptor {
	return file_shared_payment_v1_payment_action_proto_enumTypes[0].Descriptor()
}

func (PaymentType) Type() protoreflect.EnumType {
	return &file_shared_payment_v1_payment_action_proto_enumTypes[0]
}

func (x PaymentType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use PaymentType.Descriptor instead.
func (PaymentType) EnumDescriptor() ([]byte, []int) {
	return file_shared_payment_v1_payment_action_proto_rawDescGZIP(), []int{0}
}

// The action made towards a payment
type PaymentAction struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// payment action id
	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	// The payment action amount.
	Amount *v1.Money `protobuf:"bytes,2,opt,name=amount,proto3" json:"amount,omitempty"`
	// The payment type associated to the action.
	PaymentType PaymentType `protobuf:"varint,3,opt,name=payment_type,json=paymentType,proto3,enum=shared.payment.v1.PaymentType" json:"payment_type,omitempty"`
	// ISO-1987 response code
	ResponseCode string `protobuf:"bytes,4,opt,name=response_code,json=responseCode,proto3" json:"response_code,omitempty"`
	// The time in which the action was created.
	CreatedAt *timestamppb.Timestamp `protobuf:"bytes,5,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
	// The time in which the action was successfully processed.
	ProcessedAt *timestamppb.Timestamp `protobuf:"bytes,6,opt,name=processed_at,json=processedAt,proto3" json:"processed_at,omitempty"`
}

func (x *PaymentAction) Reset() {
	*x = PaymentAction{}
	if protoimpl.UnsafeEnabled {
		mi := &file_shared_payment_v1_payment_action_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PaymentAction) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PaymentAction) ProtoMessage() {}

func (x *PaymentAction) ProtoReflect() protoreflect.Message {
	mi := &file_shared_payment_v1_payment_action_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PaymentAction.ProtoReflect.Descriptor instead.
func (*PaymentAction) Descriptor() ([]byte, []int) {
	return file_shared_payment_v1_payment_action_proto_rawDescGZIP(), []int{0}
}

func (x *PaymentAction) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *PaymentAction) GetAmount() *v1.Money {
	if x != nil {
		return x.Amount
	}
	return nil
}

func (x *PaymentAction) GetPaymentType() PaymentType {
	if x != nil {
		return x.PaymentType
	}
	return PaymentType_PAYMENT_TYPE_UNSPECIFIED
}

func (x *PaymentAction) GetResponseCode() string {
	if x != nil {
		return x.ResponseCode
	}
	return ""
}

func (x *PaymentAction) GetCreatedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.CreatedAt
	}
	return nil
}

func (x *PaymentAction) GetProcessedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.ProcessedAt
	}
	return nil
}

var File_shared_payment_v1_payment_action_proto protoreflect.FileDescriptor

var file_shared_payment_v1_payment_action_proto_rawDesc = []byte{
	0x0a, 0x26, 0x73, 0x68, 0x61, 0x72, 0x65, 0x64, 0x2f, 0x70, 0x61, 0x79, 0x6d, 0x65, 0x6e, 0x74,
	0x2f, 0x76, 0x31, 0x2f, 0x70, 0x61, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x5f, 0x61, 0x63, 0x74, 0x69,
	0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x11, 0x73, 0x68, 0x61, 0x72, 0x65, 0x64,
	0x2e, 0x70, 0x61, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x2e, 0x76, 0x31, 0x1a, 0x1c, 0x73, 0x68, 0x61,
	0x72, 0x65, 0x64, 0x2f, 0x61, 0x6d, 0x6f, 0x75, 0x6e, 0x74, 0x2f, 0x76, 0x31, 0x2f, 0x6d, 0x6f,
	0x6e, 0x65, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c,
	0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73,
	0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xb2, 0x02, 0x0a, 0x0d, 0x50,
	0x61, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x41, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x0e, 0x0a, 0x02,
	0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x2f, 0x0a, 0x06,
	0x61, 0x6d, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x17, 0x2e, 0x73,
	0x68, 0x61, 0x72, 0x65, 0x64, 0x2e, 0x61, 0x6d, 0x6f, 0x75, 0x6e, 0x74, 0x2e, 0x76, 0x31, 0x2e,
	0x4d, 0x6f, 0x6e, 0x65, 0x79, 0x52, 0x06, 0x61, 0x6d, 0x6f, 0x75, 0x6e, 0x74, 0x12, 0x41, 0x0a,
	0x0c, 0x70, 0x61, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x0e, 0x32, 0x1e, 0x2e, 0x73, 0x68, 0x61, 0x72, 0x65, 0x64, 0x2e, 0x70, 0x61, 0x79,
	0x6d, 0x65, 0x6e, 0x74, 0x2e, 0x76, 0x31, 0x2e, 0x50, 0x61, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x54,
	0x79, 0x70, 0x65, 0x52, 0x0b, 0x70, 0x61, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x54, 0x79, 0x70, 0x65,
	0x12, 0x23, 0x0a, 0x0d, 0x72, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x5f, 0x63, 0x6f, 0x64,
	0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x72, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x43, 0x6f, 0x64, 0x65, 0x12, 0x39, 0x0a, 0x0a, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64,
	0x5f, 0x61, 0x74, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65,
	0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x09, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74,
	0x12, 0x3d, 0x0a, 0x0c, 0x70, 0x72, 0x6f, 0x63, 0x65, 0x73, 0x73, 0x65, 0x64, 0x5f, 0x61, 0x74,
	0x18, 0x06, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61,
	0x6d, 0x70, 0x52, 0x0b, 0x70, 0x72, 0x6f, 0x63, 0x65, 0x73, 0x73, 0x65, 0x64, 0x41, 0x74, 0x2a,
	0x75, 0x0a, 0x0b, 0x50, 0x61, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x54, 0x79, 0x70, 0x65, 0x12, 0x1c,
	0x0a, 0x18, 0x50, 0x41, 0x59, 0x4d, 0x45, 0x4e, 0x54, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x55,
	0x4e, 0x53, 0x50, 0x45, 0x43, 0x49, 0x46, 0x49, 0x45, 0x44, 0x10, 0x00, 0x12, 0x18, 0x0a, 0x14,
	0x50, 0x41, 0x59, 0x4d, 0x45, 0x4e, 0x54, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x43, 0x41, 0x50,
	0x54, 0x55, 0x52, 0x45, 0x10, 0x01, 0x12, 0x17, 0x0a, 0x13, 0x50, 0x41, 0x59, 0x4d, 0x45, 0x4e,
	0x54, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x52, 0x45, 0x46, 0x55, 0x4e, 0x44, 0x10, 0x02, 0x12,
	0x15, 0x0a, 0x11, 0x50, 0x41, 0x59, 0x4d, 0x45, 0x4e, 0x54, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f,
	0x56, 0x4f, 0x49, 0x44, 0x10, 0x03, 0x42, 0x33, 0x5a, 0x31, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62,
	0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x6a, 0x61, 0x63, 0x6b, 0x74, 0x61, 0x6e, 0x74, 0x72, 0x61, 0x6d,
	0x2f, 0x62, 0x75, 0x69, 0x6c, 0x64, 0x2f, 0x67, 0x6f, 0x2f, 0x73, 0x68, 0x61, 0x72, 0x65, 0x64,
	0x2f, 0x70, 0x61, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x2f, 0x76, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x33,
}

var (
	file_shared_payment_v1_payment_action_proto_rawDescOnce sync.Once
	file_shared_payment_v1_payment_action_proto_rawDescData = file_shared_payment_v1_payment_action_proto_rawDesc
)

func file_shared_payment_v1_payment_action_proto_rawDescGZIP() []byte {
	file_shared_payment_v1_payment_action_proto_rawDescOnce.Do(func() {
		file_shared_payment_v1_payment_action_proto_rawDescData = protoimpl.X.CompressGZIP(file_shared_payment_v1_payment_action_proto_rawDescData)
	})
	return file_shared_payment_v1_payment_action_proto_rawDescData
}

var file_shared_payment_v1_payment_action_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_shared_payment_v1_payment_action_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_shared_payment_v1_payment_action_proto_goTypes = []interface{}{
	(PaymentType)(0),              // 0: shared.payment.v1.PaymentType
	(*PaymentAction)(nil),         // 1: shared.payment.v1.PaymentAction
	(*v1.Money)(nil),              // 2: shared.amount.v1.Money
	(*timestamppb.Timestamp)(nil), // 3: google.protobuf.Timestamp
}
var file_shared_payment_v1_payment_action_proto_depIdxs = []int32{
	2, // 0: shared.payment.v1.PaymentAction.amount:type_name -> shared.amount.v1.Money
	0, // 1: shared.payment.v1.PaymentAction.payment_type:type_name -> shared.payment.v1.PaymentType
	3, // 2: shared.payment.v1.PaymentAction.created_at:type_name -> google.protobuf.Timestamp
	3, // 3: shared.payment.v1.PaymentAction.processed_at:type_name -> google.protobuf.Timestamp
	4, // [4:4] is the sub-list for method output_type
	4, // [4:4] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_shared_payment_v1_payment_action_proto_init() }
func file_shared_payment_v1_payment_action_proto_init() {
	if File_shared_payment_v1_payment_action_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_shared_payment_v1_payment_action_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PaymentAction); i {
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
			RawDescriptor: file_shared_payment_v1_payment_action_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_shared_payment_v1_payment_action_proto_goTypes,
		DependencyIndexes: file_shared_payment_v1_payment_action_proto_depIdxs,
		EnumInfos:         file_shared_payment_v1_payment_action_proto_enumTypes,
		MessageInfos:      file_shared_payment_v1_payment_action_proto_msgTypes,
	}.Build()
	File_shared_payment_v1_payment_action_proto = out.File
	file_shared_payment_v1_payment_action_proto_rawDesc = nil
	file_shared_payment_v1_payment_action_proto_goTypes = nil
	file_shared_payment_v1_payment_action_proto_depIdxs = nil
}
