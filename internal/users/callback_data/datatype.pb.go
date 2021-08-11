// Code generated by protoc-gen-go. DO NOT EDIT.
// source: internal/users/callback_data/datatype.proto

package callback_data

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type KeyboardType int32

const (
	KeyboardType_KT_NULL           KeyboardType = 0
	KeyboardType_KT_ACCOUNT        KeyboardType = 1
	KeyboardType_KT_TRANSACTIONS   KeyboardType = 2
	KeyboardType_KT_MULTI_OUT      KeyboardType = 3
	KeyboardType_KT_MULTI_OUT_SAME KeyboardType = 4
	KeyboardType_KT_BLOCKS         KeyboardType = 5
)

var KeyboardType_name = map[int32]string{
	0: "KT_NULL",
	1: "KT_ACCOUNT",
	2: "KT_TRANSACTIONS",
	3: "KT_MULTI_OUT",
	4: "KT_MULTI_OUT_SAME",
	5: "KT_BLOCKS",
}

var KeyboardType_value = map[string]int32{
	"KT_NULL":           0,
	"KT_ACCOUNT":        1,
	"KT_TRANSACTIONS":   2,
	"KT_MULTI_OUT":      3,
	"KT_MULTI_OUT_SAME": 4,
	"KT_BLOCKS":         5,
}

func (x KeyboardType) String() string {
	return proto.EnumName(KeyboardType_name, int32(x))
}

func (KeyboardType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_2c12396cafb085bf, []int{0}
}

type ActionType int32

const (
	ActionType_AT_NULL                     ActionType = 0
	ActionType_AT_REFRESH                  ActionType = 1
	ActionType_AT_TRANSACTIONS             ActionType = 2
	ActionType_AT_MULTI_OUT                ActionType = 3
	ActionType_AT_MULTI_OUT_SAME           ActionType = 4
	ActionType_AT_BLOCKS                   ActionType = 5
	ActionType_AT_BACK                     ActionType = 6
	ActionType_AT_NEXT                     ActionType = 7
	ActionType_AT_PREV                     ActionType = 8
	ActionType_AT_ENABLE_INCOME_TX_NOTIFY  ActionType = 9
	ActionType_AT_DISABLE_INCOME_TX_NOTIFY ActionType = 10
	ActionType_AT_ENABLE_BLOCK_NOTIFY      ActionType = 11
	ActionType_AT_DISABLE_BLOCK_NOTIFY     ActionType = 12
	ActionType_AT_ENABLE_OUTGO_TX_NOTIFY   ActionType = 13
	ActionType_AT_DISABLE_OUTGO_TX_NOTIFY  ActionType = 14
)

var ActionType_name = map[int32]string{
	0:  "AT_NULL",
	1:  "AT_REFRESH",
	2:  "AT_TRANSACTIONS",
	3:  "AT_MULTI_OUT",
	4:  "AT_MULTI_OUT_SAME",
	5:  "AT_BLOCKS",
	6:  "AT_BACK",
	7:  "AT_NEXT",
	8:  "AT_PREV",
	9:  "AT_ENABLE_INCOME_TX_NOTIFY",
	10: "AT_DISABLE_INCOME_TX_NOTIFY",
	11: "AT_ENABLE_BLOCK_NOTIFY",
	12: "AT_DISABLE_BLOCK_NOTIFY",
	13: "AT_ENABLE_OUTGO_TX_NOTIFY",
	14: "AT_DISABLE_OUTGO_TX_NOTIFY",
}

var ActionType_value = map[string]int32{
	"AT_NULL":                     0,
	"AT_REFRESH":                  1,
	"AT_TRANSACTIONS":             2,
	"AT_MULTI_OUT":                3,
	"AT_MULTI_OUT_SAME":           4,
	"AT_BLOCKS":                   5,
	"AT_BACK":                     6,
	"AT_NEXT":                     7,
	"AT_PREV":                     8,
	"AT_ENABLE_INCOME_TX_NOTIFY":  9,
	"AT_DISABLE_INCOME_TX_NOTIFY": 10,
	"AT_ENABLE_BLOCK_NOTIFY":      11,
	"AT_DISABLE_BLOCK_NOTIFY":     12,
	"AT_ENABLE_OUTGO_TX_NOTIFY":   13,
	"AT_DISABLE_OUTGO_TX_NOTIFY":  14,
}

func (x ActionType) String() string {
	return proto.EnumName(ActionType_name, int32(x))
}

func (ActionType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_2c12396cafb085bf, []int{1}
}

type QueryDataType struct {
	MessageId            int64        `protobuf:"varint,1,opt,name=message_id,json=messageId,proto3" json:"message_id,omitempty"`
	Account              string       `protobuf:"bytes,2,opt,name=account,proto3" json:"account,omitempty"`
	Keyboard             KeyboardType `protobuf:"varint,3,opt,name=keyboard,proto3,enum=callback_data.KeyboardType" json:"keyboard,omitempty"`
	Action               ActionType   `protobuf:"varint,4,opt,name=action,proto3,enum=callback_data.ActionType" json:"action,omitempty"`
	XXX_NoUnkeyedLiteral struct{}     `json:"-"`
	XXX_unrecognized     []byte       `json:"-"`
	XXX_sizecache        int32        `json:"-"`
}

func (m *QueryDataType) Reset()         { *m = QueryDataType{} }
func (m *QueryDataType) String() string { return proto.CompactTextString(m) }
func (*QueryDataType) ProtoMessage()    {}
func (*QueryDataType) Descriptor() ([]byte, []int) {
	return fileDescriptor_2c12396cafb085bf, []int{0}
}

func (m *QueryDataType) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_QueryDataType.Unmarshal(m, b)
}
func (m *QueryDataType) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_QueryDataType.Marshal(b, m, deterministic)
}
func (m *QueryDataType) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryDataType.Merge(m, src)
}
func (m *QueryDataType) XXX_Size() int {
	return xxx_messageInfo_QueryDataType.Size(m)
}
func (m *QueryDataType) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryDataType.DiscardUnknown(m)
}

var xxx_messageInfo_QueryDataType proto.InternalMessageInfo

func (m *QueryDataType) GetMessageId() int64 {
	if m != nil {
		return m.MessageId
	}
	return 0
}

func (m *QueryDataType) GetAccount() string {
	if m != nil {
		return m.Account
	}
	return ""
}

func (m *QueryDataType) GetKeyboard() KeyboardType {
	if m != nil {
		return m.Keyboard
	}
	return KeyboardType_KT_NULL
}

func (m *QueryDataType) GetAction() ActionType {
	if m != nil {
		return m.Action
	}
	return ActionType_AT_NULL
}

func init() {
	proto.RegisterEnum("callback_data.KeyboardType", KeyboardType_name, KeyboardType_value)
	proto.RegisterEnum("callback_data.ActionType", ActionType_name, ActionType_value)
	proto.RegisterType((*QueryDataType)(nil), "callback_data.QueryDataType")
}

func init() {
	proto.RegisterFile("internal/users/callback_data/datatype.proto", fileDescriptor_2c12396cafb085bf)
}

var fileDescriptor_2c12396cafb085bf = []byte{
	// 416 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x6c, 0x92, 0x51, 0x8f, 0xd2, 0x40,
	0x10, 0xc7, 0x2d, 0x9c, 0x70, 0xcc, 0x01, 0xae, 0x63, 0xd4, 0xde, 0x91, 0x53, 0xe2, 0x13, 0x39,
	0x13, 0x88, 0xfa, 0xe0, 0xf3, 0xd2, 0xeb, 0x69, 0xd3, 0xd2, 0xd5, 0x76, 0x6a, 0xce, 0xa7, 0xcd,
	0x52, 0x36, 0x86, 0x1c, 0x52, 0x52, 0x4a, 0x22, 0xdf, 0xcb, 0xcf, 0xe2, 0xe7, 0x31, 0xd4, 0x62,
	0xcb, 0x1d, 0x2f, 0x9b, 0xcc, 0xfc, 0xe7, 0x37, 0xbf, 0x79, 0x58, 0x78, 0x3b, 0x5f, 0x66, 0x3a,
	0x5d, 0xaa, 0xc5, 0x68, 0xb3, 0xd6, 0xe9, 0x7a, 0x14, 0xab, 0xc5, 0x62, 0xaa, 0xe2, 0x3b, 0x39,
	0x53, 0x99, 0x1a, 0xed, 0x9e, 0x6c, 0xbb, 0xd2, 0xc3, 0x55, 0x9a, 0x64, 0x09, 0x76, 0x0e, 0xd2,
	0x37, 0xbf, 0x0d, 0xe8, 0x7c, 0xdd, 0xe8, 0x74, 0x7b, 0xad, 0x32, 0x45, 0xdb, 0x95, 0xc6, 0x4b,
	0x80, 0x9f, 0x7a, 0xbd, 0x56, 0x3f, 0xb4, 0x9c, 0xcf, 0x4c, 0xa3, 0x6f, 0x0c, 0xea, 0x41, 0xab,
	0xe8, 0x38, 0x33, 0x34, 0xa1, 0xa9, 0xe2, 0x38, 0xd9, 0x2c, 0x33, 0xb3, 0xd6, 0x37, 0x06, 0xad,
	0x60, 0x5f, 0xe2, 0x47, 0x38, 0xbd, 0xd3, 0xdb, 0x69, 0xa2, 0xd2, 0x99, 0x59, 0xef, 0x1b, 0x83,
	0xee, 0xfb, 0xde, 0xf0, 0x40, 0x36, 0x74, 0x8b, 0x78, 0xe7, 0x09, 0xfe, 0x0f, 0xe3, 0x3b, 0x68,
	0xa8, 0x38, 0x9b, 0x27, 0x4b, 0xf3, 0x24, 0xc7, 0xce, 0xef, 0x61, 0x3c, 0x0f, 0x73, 0xa8, 0x18,
	0xbc, 0xfa, 0x05, 0xed, 0xea, 0x32, 0x3c, 0x83, 0xa6, 0x4b, 0xd2, 0x8f, 0x3c, 0x8f, 0x3d, 0xc2,
	0x2e, 0x80, 0x4b, 0x92, 0x5b, 0x96, 0x88, 0x7c, 0x62, 0x06, 0x3e, 0x83, 0x27, 0x2e, 0x49, 0x0a,
	0xb8, 0x1f, 0x72, 0x8b, 0x1c, 0xe1, 0x87, 0xac, 0x86, 0x0c, 0xda, 0x2e, 0xc9, 0x49, 0xe4, 0x91,
	0x23, 0x45, 0x44, 0xac, 0x8e, 0xcf, 0xe1, 0x69, 0xb5, 0x23, 0x43, 0x3e, 0xb1, 0xd9, 0x09, 0x76,
	0xa0, 0xe5, 0x92, 0x1c, 0x7b, 0xc2, 0x72, 0x43, 0xf6, 0xf8, 0xea, 0x4f, 0x0d, 0xa0, 0x3c, 0x68,
	0x27, 0xe6, 0x55, 0x31, 0x27, 0x19, 0xd8, 0x37, 0x81, 0x1d, 0x7e, 0xfe, 0x27, 0xe6, 0xc7, 0xc4,
	0xfc, 0x81, 0x98, 0x1f, 0x17, 0xf3, 0x52, 0x5c, 0x98, 0xc6, 0xdc, 0x72, 0x59, 0x63, 0xaf, 0xb5,
	0x6f, 0x89, 0x35, 0x8b, 0xe2, 0x4b, 0x60, 0x7f, 0x63, 0xa7, 0xf8, 0x0a, 0x2e, 0x38, 0x49, 0xdb,
	0xe7, 0x63, 0xcf, 0x96, 0x8e, 0x6f, 0x89, 0x89, 0x2d, 0xe9, 0x56, 0xfa, 0x82, 0x9c, 0x9b, 0xef,
	0xac, 0x85, 0xaf, 0xa1, 0xc7, 0x49, 0x5e, 0x3b, 0xe1, 0xf1, 0x01, 0xc0, 0x0b, 0x78, 0x51, 0x2e,
	0xc8, 0xed, 0xfb, 0xec, 0x0c, 0x7b, 0xf0, 0xb2, 0x02, 0x1f, 0x84, 0x6d, 0xbc, 0x84, 0xf3, 0x12,
	0x14, 0x11, 0x7d, 0x12, 0x95, 0xbd, 0x9d, 0xe2, 0xb0, 0x3d, 0x7b, 0x3f, 0xef, 0x4e, 0x1b, 0xf9,
	0xff, 0xfc, 0xf0, 0x37, 0x00, 0x00, 0xff, 0xff, 0x70, 0x42, 0x85, 0x94, 0xce, 0x02, 0x00, 0x00,
}
