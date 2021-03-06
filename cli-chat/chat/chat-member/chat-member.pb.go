// Code generated by protoc-gen-go. DO NOT EDIT.
// source: chat-member.proto

package chatmember

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import timestamp "github.com/golang/protobuf/ptypes/timestamp"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type Member_MessageType int32

const (
	// messages to handle subscriptions
	Member_JOIN  Member_MessageType = 0
	Member_LEAVE Member_MessageType = 1
	Member_PING  Member_MessageType = 2
)

var Member_MessageType_name = map[int32]string{
	0: "JOIN",
	1: "LEAVE",
	2: "PING",
}
var Member_MessageType_value = map[string]int32{
	"JOIN":  0,
	"LEAVE": 1,
	"PING":  2,
}

func (x Member_MessageType) String() string {
	return proto.EnumName(Member_MessageType_name, int32(x))
}
func (Member_MessageType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_chat_member_a304ca18d77ccdb1, []int{0, 0}
}

// Chat member
type Member struct {
	MsgType              Member_MessageType   `protobuf:"varint,1,opt,name=msgType,enum=chatmember.Member_MessageType" json:"msgType,omitempty"`
	Name                 string               `protobuf:"bytes,2,opt,name=name" json:"name,omitempty"`
	Sender               string               `protobuf:"bytes,3,opt,name=sender" json:"sender,omitempty"`
	Timestamp            *timestamp.Timestamp `protobuf:"bytes,4,opt,name=timestamp" json:"timestamp,omitempty"`
	XXX_NoUnkeyedLiteral struct{}             `json:"-"`
	XXX_unrecognized     []byte               `json:"-"`
	XXX_sizecache        int32                `json:"-"`
}

func (m *Member) Reset()         { *m = Member{} }
func (m *Member) String() string { return proto.CompactTextString(m) }
func (*Member) ProtoMessage()    {}
func (*Member) Descriptor() ([]byte, []int) {
	return fileDescriptor_chat_member_a304ca18d77ccdb1, []int{0}
}
func (m *Member) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Member.Unmarshal(m, b)
}
func (m *Member) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Member.Marshal(b, m, deterministic)
}
func (dst *Member) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Member.Merge(dst, src)
}
func (m *Member) XXX_Size() int {
	return xxx_messageInfo_Member.Size(m)
}
func (m *Member) XXX_DiscardUnknown() {
	xxx_messageInfo_Member.DiscardUnknown(m)
}

var xxx_messageInfo_Member proto.InternalMessageInfo

func (m *Member) GetMsgType() Member_MessageType {
	if m != nil {
		return m.MsgType
	}
	return Member_JOIN
}

func (m *Member) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *Member) GetSender() string {
	if m != nil {
		return m.Sender
	}
	return ""
}

func (m *Member) GetTimestamp() *timestamp.Timestamp {
	if m != nil {
		return m.Timestamp
	}
	return nil
}

func init() {
	proto.RegisterType((*Member)(nil), "chatmember.Member")
	proto.RegisterEnum("chatmember.Member_MessageType", Member_MessageType_name, Member_MessageType_value)
}

func init() { proto.RegisterFile("chat-member.proto", fileDescriptor_chat_member_a304ca18d77ccdb1) }

var fileDescriptor_chat_member_a304ca18d77ccdb1 = []byte{
	// 219 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x12, 0x4c, 0xce, 0x48, 0x2c,
	0xd1, 0xcd, 0x4d, 0xcd, 0x4d, 0x4a, 0x2d, 0xd2, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0xe2, 0x02,
	0x09, 0x41, 0x44, 0xa4, 0xe4, 0xd3, 0xf3, 0xf3, 0xd3, 0x73, 0x52, 0xf5, 0xc1, 0x32, 0x49, 0xa5,
	0x69, 0xfa, 0x25, 0x99, 0xb9, 0xa9, 0xc5, 0x25, 0x89, 0xb9, 0x05, 0x10, 0xc5, 0x4a, 0xd7, 0x18,
	0xb9, 0xd8, 0x7c, 0xc1, 0x6a, 0x85, 0x2c, 0xb8, 0xd8, 0x73, 0x8b, 0xd3, 0x43, 0x2a, 0x0b, 0x52,
	0x25, 0x18, 0x15, 0x18, 0x35, 0xf8, 0x8c, 0xe4, 0xf4, 0x10, 0x26, 0xe9, 0xf9, 0xc2, 0xa8, 0xe2,
	0xe2, 0xc4, 0xf4, 0x54, 0x90, 0xaa, 0x20, 0x98, 0x72, 0x21, 0x21, 0x2e, 0x96, 0xbc, 0xc4, 0xdc,
	0x54, 0x09, 0x26, 0x05, 0x46, 0x0d, 0xce, 0x20, 0x30, 0x5b, 0x48, 0x8c, 0x8b, 0xad, 0x38, 0x35,
	0x2f, 0x25, 0xb5, 0x48, 0x82, 0x19, 0x2c, 0x0a, 0xe5, 0x09, 0x59, 0x70, 0x71, 0xc2, 0xdd, 0x20,
	0xc1, 0xa2, 0xc0, 0xa8, 0xc1, 0x6d, 0x24, 0xa5, 0x07, 0x71, 0xa5, 0x1e, 0xcc, 0x95, 0x7a, 0x21,
	0x30, 0x15, 0x41, 0x08, 0xc5, 0x4a, 0x3a, 0x5c, 0xdc, 0x48, 0xb6, 0x0b, 0x71, 0x70, 0xb1, 0x78,
	0xf9, 0x7b, 0xfa, 0x09, 0x30, 0x08, 0x71, 0x72, 0xb1, 0xfa, 0xb8, 0x3a, 0x86, 0xb9, 0x0a, 0x30,
	0x82, 0x04, 0x03, 0x3c, 0xfd, 0xdc, 0x05, 0x98, 0x92, 0xd8, 0xc0, 0x86, 0x19, 0x03, 0x02, 0x00,
	0x00, 0xff, 0xff, 0x24, 0xed, 0x8d, 0x52, 0x21, 0x01, 0x00, 0x00,
}
