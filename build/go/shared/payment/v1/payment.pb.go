// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        v3.18.1
// source: shared/payment/v1/payment.proto

package v1

import (
	v1 "github.com/jacktantram/payments-api/build/go/shared/amount/v1"
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

// Represents the current status of the payment.
type PaymentStatus int32

const (
	// If the payment status is not provided.
	PaymentStatus_PAYMENT_STATUS_UNSPECIFIED PaymentStatus = 0
	// The payment is in a pending status.
	PaymentStatus_PAYMENT_STATUS_PENDING PaymentStatus = 1
	// The payment is currently partially captured.
	PaymentStatus_PAYMENT_STATUS_PARTIALLY_CAPTURED PaymentStatus = 3
	// The payment has been completely captured.
	PaymentStatus_PAYMENT_STATUS_CAPTURED PaymentStatus = 4
	// The payment has been partially refunded.
	PaymentStatus_PAYMENT_STATUS_PARTIALLY_REFUNDED PaymentStatus = 5
	// The payment has been fully refunded.
	PaymentStatus_PAYMENT_STATUS_REFUNDED PaymentStatus = 6
	// The payment has been voided.
	PaymentStatus_PAYMENT_STATUS_VOIDED PaymentStatus = 7
)

// Enum value maps for PaymentStatus.
var (
	PaymentStatus_name = map[int32]string{
		0: "PAYMENT_STATUS_UNSPECIFIED",
		1: "PAYMENT_STATUS_PENDING",
		3: "PAYMENT_STATUS_PARTIALLY_CAPTURED",
		4: "PAYMENT_STATUS_CAPTURED",
		5: "PAYMENT_STATUS_PARTIALLY_REFUNDED",
		6: "PAYMENT_STATUS_REFUNDED",
		7: "PAYMENT_STATUS_VOIDED",
	}
	PaymentStatus_value = map[string]int32{
		"PAYMENT_STATUS_UNSPECIFIED":        0,
		"PAYMENT_STATUS_PENDING":            1,
		"PAYMENT_STATUS_PARTIALLY_CAPTURED": 3,
		"PAYMENT_STATUS_CAPTURED":           4,
		"PAYMENT_STATUS_PARTIALLY_REFUNDED": 5,
		"PAYMENT_STATUS_REFUNDED":           6,
		"PAYMENT_STATUS_VOIDED":             7,
	}
)

func (x PaymentStatus) Enum() *PaymentStatus {
	p := new(PaymentStatus)
	*p = x
	return p
}

func (x PaymentStatus) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (PaymentStatus) Descriptor() protoreflect.EnumDescriptor {
	return file_shared_payment_v1_payment_proto_enumTypes[0].Descriptor()
}

func (PaymentStatus) Type() protoreflect.EnumType {
	return &file_shared_payment_v1_payment_proto_enumTypes[0]
}

func (x PaymentStatus) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use PaymentStatus.Descriptor instead.
func (PaymentStatus) EnumDescriptor() ([]byte, []int) {
	return file_shared_payment_v1_payment_proto_rawDescGZIP(), []int{0}
}

// Defines a payment entity.
type Payment struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The unique payment identifier.
	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	// The payment amount.
	Amount *v1.Money `protobuf:"bytes,2,opt,name=amount,proto3" json:"amount,omitempty"`
	// The status of the payment/
	PaymentStatus PaymentStatus `protobuf:"varint,3,opt,name=payment_status,json=paymentStatus,proto3,enum=shared.payment.v1.PaymentStatus" json:"payment_status,omitempty"`
	// The linking action id to the payment. See PaymentAction.
	ActionId string `protobuf:"bytes,4,opt,name=action_id,json=actionId,proto3" json:"action_id,omitempty"`
	// The date the payment was created.
	CreatedAt *timestamppb.Timestamp `protobuf:"bytes,5,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
	// The date the payment was updated.
	UpdatedAt *timestamppb.Timestamp `protobuf:"bytes,6,opt,name=updated_at,json=updatedAt,proto3" json:"updated_at,omitempty"`
}

func (x *Payment) Reset() {
	*x = Payment{}
	if protoimpl.UnsafeEnabled {
		mi := &file_shared_payment_v1_payment_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Payment) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Payment) ProtoMessage() {}

func (x *Payment) ProtoReflect() protoreflect.Message {
	mi := &file_shared_payment_v1_payment_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Payment.ProtoReflect.Descriptor instead.
func (*Payment) Descriptor() ([]byte, []int) {
	return file_shared_payment_v1_payment_proto_rawDescGZIP(), []int{0}
}

func (x *Payment) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Payment) GetAmount() *v1.Money {
	if x != nil {
		return x.Amount
	}
	return nil
}

func (x *Payment) GetPaymentStatus() PaymentStatus {
	if x != nil {
		return x.PaymentStatus
	}
	return PaymentStatus_PAYMENT_STATUS_UNSPECIFIED
}

func (x *Payment) GetActionId() string {
	if x != nil {
		return x.ActionId
	}
	return ""
}

func (x *Payment) GetCreatedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.CreatedAt
	}
	return nil
}

func (x *Payment) GetUpdatedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.UpdatedAt
	}
	return nil
}

var File_shared_payment_v1_payment_proto protoreflect.FileDescriptor

var file_shared_payment_v1_payment_proto_rawDesc = []byte{
	0x0a, 0x1f, 0x73, 0x68, 0x61, 0x72, 0x65, 0x64, 0x2f, 0x70, 0x61, 0x79, 0x6d, 0x65, 0x6e, 0x74,
	0x2f, 0x76, 0x31, 0x2f, 0x70, 0x61, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x12, 0x11, 0x73, 0x68, 0x61, 0x72, 0x65, 0x64, 0x2e, 0x70, 0x61, 0x79, 0x6d, 0x65, 0x6e,
	0x74, 0x2e, 0x76, 0x31, 0x1a, 0x1c, 0x73, 0x68, 0x61, 0x72, 0x65, 0x64, 0x2f, 0x61, 0x6d, 0x6f,
	0x75, 0x6e, 0x74, 0x2f, 0x76, 0x31, 0x2f, 0x6d, 0x6f, 0x6e, 0x65, 0x79, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x1a, 0x26, 0x73, 0x68, 0x61, 0x72, 0x65, 0x64, 0x2f, 0x70, 0x61, 0x79, 0x6d,
	0x65, 0x6e, 0x74, 0x2f, 0x76, 0x31, 0x2f, 0x70, 0x61, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x5f, 0x6d,
	0x65, 0x74, 0x68, 0x6f, 0x64, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xa6, 0x02, 0x0a, 0x07,
	0x50, 0x61, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x2f, 0x0a, 0x06, 0x61, 0x6d, 0x6f, 0x75, 0x6e,
	0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x17, 0x2e, 0x73, 0x68, 0x61, 0x72, 0x65, 0x64,
	0x2e, 0x61, 0x6d, 0x6f, 0x75, 0x6e, 0x74, 0x2e, 0x76, 0x31, 0x2e, 0x4d, 0x6f, 0x6e, 0x65, 0x79,
	0x52, 0x06, 0x61, 0x6d, 0x6f, 0x75, 0x6e, 0x74, 0x12, 0x47, 0x0a, 0x0e, 0x70, 0x61, 0x79, 0x6d,
	0x65, 0x6e, 0x74, 0x5f, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0e,
	0x32, 0x20, 0x2e, 0x73, 0x68, 0x61, 0x72, 0x65, 0x64, 0x2e, 0x70, 0x61, 0x79, 0x6d, 0x65, 0x6e,
	0x74, 0x2e, 0x76, 0x31, 0x2e, 0x50, 0x61, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x53, 0x74, 0x61, 0x74,
	0x75, 0x73, 0x52, 0x0d, 0x70, 0x61, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x53, 0x74, 0x61, 0x74, 0x75,
	0x73, 0x12, 0x1b, 0x0a, 0x09, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x69, 0x64, 0x18, 0x04,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x49, 0x64, 0x12, 0x39,
	0x0a, 0x0a, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x18, 0x05, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x09,
	0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x12, 0x39, 0x0a, 0x0a, 0x75, 0x70, 0x64,
	0x61, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e,
	0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e,
	0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x09, 0x75, 0x70, 0x64, 0x61, 0x74,
	0x65, 0x64, 0x41, 0x74, 0x2a, 0xee, 0x01, 0x0a, 0x0d, 0x50, 0x61, 0x79, 0x6d, 0x65, 0x6e, 0x74,
	0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x1e, 0x0a, 0x1a, 0x50, 0x41, 0x59, 0x4d, 0x45, 0x4e,
	0x54, 0x5f, 0x53, 0x54, 0x41, 0x54, 0x55, 0x53, 0x5f, 0x55, 0x4e, 0x53, 0x50, 0x45, 0x43, 0x49,
	0x46, 0x49, 0x45, 0x44, 0x10, 0x00, 0x12, 0x1a, 0x0a, 0x16, 0x50, 0x41, 0x59, 0x4d, 0x45, 0x4e,
	0x54, 0x5f, 0x53, 0x54, 0x41, 0x54, 0x55, 0x53, 0x5f, 0x50, 0x45, 0x4e, 0x44, 0x49, 0x4e, 0x47,
	0x10, 0x01, 0x12, 0x25, 0x0a, 0x21, 0x50, 0x41, 0x59, 0x4d, 0x45, 0x4e, 0x54, 0x5f, 0x53, 0x54,
	0x41, 0x54, 0x55, 0x53, 0x5f, 0x50, 0x41, 0x52, 0x54, 0x49, 0x41, 0x4c, 0x4c, 0x59, 0x5f, 0x43,
	0x41, 0x50, 0x54, 0x55, 0x52, 0x45, 0x44, 0x10, 0x03, 0x12, 0x1b, 0x0a, 0x17, 0x50, 0x41, 0x59,
	0x4d, 0x45, 0x4e, 0x54, 0x5f, 0x53, 0x54, 0x41, 0x54, 0x55, 0x53, 0x5f, 0x43, 0x41, 0x50, 0x54,
	0x55, 0x52, 0x45, 0x44, 0x10, 0x04, 0x12, 0x25, 0x0a, 0x21, 0x50, 0x41, 0x59, 0x4d, 0x45, 0x4e,
	0x54, 0x5f, 0x53, 0x54, 0x41, 0x54, 0x55, 0x53, 0x5f, 0x50, 0x41, 0x52, 0x54, 0x49, 0x41, 0x4c,
	0x4c, 0x59, 0x5f, 0x52, 0x45, 0x46, 0x55, 0x4e, 0x44, 0x45, 0x44, 0x10, 0x05, 0x12, 0x1b, 0x0a,
	0x17, 0x50, 0x41, 0x59, 0x4d, 0x45, 0x4e, 0x54, 0x5f, 0x53, 0x54, 0x41, 0x54, 0x55, 0x53, 0x5f,
	0x52, 0x45, 0x46, 0x55, 0x4e, 0x44, 0x45, 0x44, 0x10, 0x06, 0x12, 0x19, 0x0a, 0x15, 0x50, 0x41,
	0x59, 0x4d, 0x45, 0x4e, 0x54, 0x5f, 0x53, 0x54, 0x41, 0x54, 0x55, 0x53, 0x5f, 0x56, 0x4f, 0x49,
	0x44, 0x45, 0x44, 0x10, 0x07, 0x42, 0x40, 0x5a, 0x3e, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e,
	0x63, 0x6f, 0x6d, 0x2f, 0x6a, 0x61, 0x63, 0x6b, 0x74, 0x61, 0x6e, 0x74, 0x72, 0x61, 0x6d, 0x2f,
	0x70, 0x61, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x2d, 0x61, 0x70, 0x69, 0x2f, 0x62, 0x75, 0x69,
	0x6c, 0x64, 0x2f, 0x67, 0x6f, 0x2f, 0x73, 0x68, 0x61, 0x72, 0x65, 0x64, 0x2f, 0x70, 0x61, 0x79,
	0x6d, 0x65, 0x6e, 0x74, 0x2f, 0x76, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_shared_payment_v1_payment_proto_rawDescOnce sync.Once
	file_shared_payment_v1_payment_proto_rawDescData = file_shared_payment_v1_payment_proto_rawDesc
)

func file_shared_payment_v1_payment_proto_rawDescGZIP() []byte {
	file_shared_payment_v1_payment_proto_rawDescOnce.Do(func() {
		file_shared_payment_v1_payment_proto_rawDescData = protoimpl.X.CompressGZIP(file_shared_payment_v1_payment_proto_rawDescData)
	})
	return file_shared_payment_v1_payment_proto_rawDescData
}

var file_shared_payment_v1_payment_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_shared_payment_v1_payment_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_shared_payment_v1_payment_proto_goTypes = []interface{}{
	(PaymentStatus)(0),            // 0: shared.payment.v1.PaymentStatus
	(*Payment)(nil),               // 1: shared.payment.v1.Payment
	(*v1.Money)(nil),              // 2: shared.amount.v1.Money
	(*timestamppb.Timestamp)(nil), // 3: google.protobuf.Timestamp
}
var file_shared_payment_v1_payment_proto_depIdxs = []int32{
	2, // 0: shared.payment.v1.Payment.amount:type_name -> shared.amount.v1.Money
	0, // 1: shared.payment.v1.Payment.payment_status:type_name -> shared.payment.v1.PaymentStatus
	3, // 2: shared.payment.v1.Payment.created_at:type_name -> google.protobuf.Timestamp
	3, // 3: shared.payment.v1.Payment.updated_at:type_name -> google.protobuf.Timestamp
	4, // [4:4] is the sub-list for method output_type
	4, // [4:4] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_shared_payment_v1_payment_proto_init() }
func file_shared_payment_v1_payment_proto_init() {
	if File_shared_payment_v1_payment_proto != nil {
		return
	}
	file_shared_payment_v1_payment_method_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_shared_payment_v1_payment_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Payment); i {
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
			RawDescriptor: file_shared_payment_v1_payment_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_shared_payment_v1_payment_proto_goTypes,
		DependencyIndexes: file_shared_payment_v1_payment_proto_depIdxs,
		EnumInfos:         file_shared_payment_v1_payment_proto_enumTypes,
		MessageInfos:      file_shared_payment_v1_payment_proto_msgTypes,
	}.Build()
	File_shared_payment_v1_payment_proto = out.File
	file_shared_payment_v1_payment_proto_rawDesc = nil
	file_shared_payment_v1_payment_proto_goTypes = nil
	file_shared_payment_v1_payment_proto_depIdxs = nil
}