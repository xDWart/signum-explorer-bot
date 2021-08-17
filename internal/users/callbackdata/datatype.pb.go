// Code generated by protoc-gen-go. DO NOT EDIT.
// source: internal/users/callbackdata/datatype.proto

package callbackdata

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
	KeyboardType_KT_PRICE_CHART    KeyboardType = 6
	KeyboardType_KT_NETWORK_CHART  KeyboardType = 7
)

var KeyboardType_name = map[int32]string{
	0: "KT_NULL",
	1: "KT_ACCOUNT",
	2: "KT_TRANSACTIONS",
	3: "KT_MULTI_OUT",
	4: "KT_MULTI_OUT_SAME",
	5: "KT_BLOCKS",
	6: "KT_PRICE_CHART",
	7: "KT_NETWORK_CHART",
}

var KeyboardType_value = map[string]int32{
	"KT_NULL":           0,
	"KT_ACCOUNT":        1,
	"KT_TRANSACTIONS":   2,
	"KT_MULTI_OUT":      3,
	"KT_MULTI_OUT_SAME": 4,
	"KT_BLOCKS":         5,
	"KT_PRICE_CHART":    6,
	"KT_NETWORK_CHART":  7,
}

func (x KeyboardType) String() string {
	return proto.EnumName(KeyboardType_name, int32(x))
}

func (KeyboardType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_18009145167ae7f1, []int{0}
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
	ActionType_AT_PRICE_CHART_1_DAY        ActionType = 15
	ActionType_AT_PRICE_CHART_1_WEEK       ActionType = 16
	ActionType_AT_PRICE_CHART_1_MONTH      ActionType = 17
	ActionType_AT_PRICE_CHART_ALL          ActionType = 18
	ActionType_AT_NETWORK_CHART_1_DAY      ActionType = 19
	ActionType_AT_NETWORK_CHART_1_WEEK     ActionType = 20
	ActionType_AT_NETWORK_CHART_1_MONTH    ActionType = 21
	ActionType_AT_NETWORK_CHART_ALL        ActionType = 22
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
	15: "AT_PRICE_CHART_1_DAY",
	16: "AT_PRICE_CHART_1_WEEK",
	17: "AT_PRICE_CHART_1_MONTH",
	18: "AT_PRICE_CHART_ALL",
	19: "AT_NETWORK_CHART_1_DAY",
	20: "AT_NETWORK_CHART_1_WEEK",
	21: "AT_NETWORK_CHART_1_MONTH",
	22: "AT_NETWORK_CHART_ALL",
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
	"AT_PRICE_CHART_1_DAY":        15,
	"AT_PRICE_CHART_1_WEEK":       16,
	"AT_PRICE_CHART_1_MONTH":      17,
	"AT_PRICE_CHART_ALL":          18,
	"AT_NETWORK_CHART_1_DAY":      19,
	"AT_NETWORK_CHART_1_WEEK":     20,
	"AT_NETWORK_CHART_1_MONTH":    21,
	"AT_NETWORK_CHART_ALL":        22,
}

func (x ActionType) String() string {
	return proto.EnumName(ActionType_name, int32(x))
}

func (ActionType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_18009145167ae7f1, []int{1}
}

type QueryDataType struct {
	MessageId            int64        `protobuf:"varint,1,opt,name=message_id,json=messageId,proto3" json:"message_id,omitempty"`
	Account              string       `protobuf:"bytes,2,opt,name=account,proto3" json:"account,omitempty"`
	Keyboard             KeyboardType `protobuf:"varint,3,opt,name=keyboard,proto3,enum=callbackdata.KeyboardType" json:"keyboard,omitempty"`
	Action               ActionType   `protobuf:"varint,4,opt,name=action,proto3,enum=callbackdata.ActionType" json:"action,omitempty"`
	XXX_NoUnkeyedLiteral struct{}     `json:"-"`
	XXX_unrecognized     []byte       `json:"-"`
	XXX_sizecache        int32        `json:"-"`
}

func (m *QueryDataType) Reset()         { *m = QueryDataType{} }
func (m *QueryDataType) String() string { return proto.CompactTextString(m) }
func (*QueryDataType) ProtoMessage()    {}
func (*QueryDataType) Descriptor() ([]byte, []int) {
	return fileDescriptor_18009145167ae7f1, []int{0}
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
	proto.RegisterEnum("callbackdata.KeyboardType", KeyboardType_name, KeyboardType_value)
	proto.RegisterEnum("callbackdata.ActionType", ActionType_name, ActionType_value)
	proto.RegisterType((*QueryDataType)(nil), "callbackdata.QueryDataType")
}

func init() {
	proto.RegisterFile("internal/users/callbackdata/datatype.proto", fileDescriptor_18009145167ae7f1)
}

var fileDescriptor_18009145167ae7f1 = []byte{
	// 509 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x6c, 0x93, 0xcf, 0x6e, 0xda, 0x40,
	0x10, 0xc6, 0xeb, 0x24, 0x85, 0x30, 0x21, 0x64, 0x33, 0x01, 0xea, 0x90, 0xa6, 0x45, 0x3d, 0x21,
	0x0e, 0xd0, 0x3f, 0x52, 0xef, 0x8b, 0xd9, 0x14, 0xcb, 0xc6, 0x4e, 0xed, 0xa1, 0x49, 0x4e, 0x2b,
	0x03, 0x56, 0x85, 0x42, 0x01, 0x81, 0x39, 0xf0, 0x3a, 0x7d, 0x82, 0xaa, 0x4f, 0x58, 0x61, 0x96,
	0x62, 0x13, 0x2e, 0x48, 0xbb, 0xbf, 0xf9, 0xe6, 0xfb, 0xa1, 0x95, 0xa1, 0x3e, 0x9a, 0x44, 0xe1,
	0x7c, 0x12, 0x8c, 0x9b, 0xcb, 0x45, 0x38, 0x5f, 0x34, 0x07, 0xc1, 0x78, 0xdc, 0x0f, 0x06, 0xcf,
	0xc3, 0x20, 0x0a, 0x9a, 0xeb, 0x9f, 0x68, 0x35, 0x0b, 0x1b, 0xb3, 0xf9, 0x34, 0x9a, 0x62, 0x3e,
	0x09, 0x3f, 0xfc, 0xd5, 0xe0, 0xfc, 0xfb, 0x32, 0x9c, 0xaf, 0xda, 0x41, 0x14, 0xd0, 0x6a, 0x16,
	0xe2, 0x2d, 0xc0, 0xaf, 0x70, 0xb1, 0x08, 0x7e, 0x86, 0x72, 0x34, 0xd4, 0xb5, 0xaa, 0x56, 0x3b,
	0xf6, 0x72, 0xea, 0xc6, 0x1c, 0xa2, 0x0e, 0xd9, 0x60, 0x30, 0x98, 0x2e, 0x27, 0x91, 0x7e, 0x54,
	0xd5, 0x6a, 0x39, 0x6f, 0x7b, 0xc4, 0xaf, 0x70, 0xfa, 0x1c, 0xae, 0xfa, 0xd3, 0x60, 0x3e, 0xd4,
	0x8f, 0xab, 0x5a, 0xad, 0xf0, 0xb9, 0xd2, 0x48, 0x76, 0x35, 0x2c, 0x45, 0xd7, 0x35, 0xde, 0xff,
	0x59, 0xfc, 0x08, 0x99, 0x60, 0x10, 0x8d, 0xa6, 0x13, 0xfd, 0x24, 0x4e, 0xe9, 0xe9, 0x14, 0x8f,
	0x59, 0x9c, 0x51, 0x73, 0xf5, 0xdf, 0x1a, 0xe4, 0x93, 0xcb, 0xf0, 0x0c, 0xb2, 0x16, 0x49, 0xa7,
	0x67, 0xdb, 0xec, 0x15, 0x16, 0x00, 0x2c, 0x92, 0xdc, 0x30, 0xdc, 0x9e, 0x43, 0x4c, 0xc3, 0x2b,
	0xb8, 0xb0, 0x48, 0x92, 0xc7, 0x1d, 0x9f, 0x1b, 0x64, 0xba, 0x8e, 0xcf, 0x8e, 0x90, 0x41, 0xde,
	0x22, 0xd9, 0xed, 0xd9, 0x64, 0x4a, 0xb7, 0x47, 0xec, 0x18, 0x4b, 0x70, 0x99, 0xbc, 0x91, 0x3e,
	0xef, 0x0a, 0x76, 0x82, 0xe7, 0x90, 0xb3, 0x48, 0xb6, 0x6c, 0xd7, 0xb0, 0x7c, 0xf6, 0x1a, 0x11,
	0x0a, 0x16, 0xc9, 0x7b, 0xcf, 0x34, 0x84, 0x34, 0x3a, 0xdc, 0x23, 0x96, 0xc1, 0x22, 0xb0, 0x75,
	0xbb, 0xa0, 0x07, 0xd7, 0xb3, 0xd4, 0x6d, 0xb6, 0xfe, 0xe7, 0x04, 0x60, 0xe7, 0xbe, 0x56, 0xe4,
	0x49, 0x45, 0x4e, 0xd2, 0x13, 0x77, 0x9e, 0xf0, 0x3b, 0x1b, 0x45, 0x7e, 0x48, 0x91, 0xbf, 0x50,
	0xe4, 0x87, 0x15, 0x79, 0x42, 0x71, 0xd3, 0xd4, 0xe2, 0x86, 0xc5, 0x32, 0xdb, 0x5a, 0xf1, 0x48,
	0x2c, 0xab, 0x0e, 0xf7, 0x9e, 0xf8, 0xc1, 0x4e, 0xf1, 0x1d, 0x54, 0x38, 0x49, 0xe1, 0xf0, 0x96,
	0x2d, 0xa4, 0xe9, 0x18, 0x6e, 0x57, 0x48, 0x7a, 0x94, 0x8e, 0x4b, 0xe6, 0xdd, 0x13, 0xcb, 0xe1,
	0x7b, 0xb8, 0xe1, 0x24, 0xdb, 0xa6, 0x7f, 0x78, 0x00, 0xb0, 0x02, 0xe5, 0xdd, 0x82, 0xb8, 0x7d,
	0xcb, 0xce, 0xf0, 0x06, 0xde, 0x24, 0xc2, 0x29, 0x98, 0xc7, 0x5b, 0xb8, 0xde, 0x05, 0xdd, 0x1e,
	0x7d, 0x73, 0x13, 0x7b, 0xcf, 0x95, 0xd8, 0x36, 0xbb, 0xcf, 0x0b, 0xa8, 0x43, 0x91, 0xa7, 0x9e,
	0x40, 0x7e, 0x92, 0x6d, 0xfe, 0xc4, 0x2e, 0xf0, 0x1a, 0x4a, 0x2f, 0xc8, 0x83, 0x10, 0x16, 0x63,
	0x4a, 0x36, 0x8d, 0xba, 0xae, 0x43, 0x1d, 0x76, 0x89, 0x65, 0xc0, 0x3d, 0xc6, 0x6d, 0x9b, 0xa1,
	0xca, 0xa4, 0xde, 0x55, 0x55, 0x5d, 0xa9, 0x3f, 0xb8, 0xcf, 0xe2, 0xb2, 0x22, 0xbe, 0x05, 0xfd,
	0x00, 0xdc, 0xd4, 0x95, 0x94, 0x7f, 0x9a, 0xae, 0x0b, 0xcb, 0xfd, 0x4c, 0xfc, 0x85, 0x7e, 0xf9,
	0x17, 0x00, 0x00, 0xff, 0xff, 0x74, 0x88, 0x7e, 0x4a, 0xcf, 0x03, 0x00, 0x00,
}