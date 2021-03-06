// Code generated by protoc-gen-go. DO NOT EDIT.
// source: prof.proto

package prof

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type ProfMsg struct {
	RequestID            []byte   `protobuf:"bytes,1,opt,name=requestID,proto3" json:"requestID,omitempty"`
	Content              []byte   `protobuf:"bytes,2,opt,name=content,proto3" json:"content,omitempty"`
	GroupID              []byte   `protobuf:"bytes,3,opt,name=groupID,proto3" json:"groupID,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ProfMsg) Reset()         { *m = ProfMsg{} }
func (m *ProfMsg) String() string { return proto.CompactTextString(m) }
func (*ProfMsg) ProtoMessage()    {}
func (*ProfMsg) Descriptor() ([]byte, []int) {
	return fileDescriptor_prof_043c76f9d6cc050b, []int{0}
}
func (m *ProfMsg) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ProfMsg.Unmarshal(m, b)
}
func (m *ProfMsg) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ProfMsg.Marshal(b, m, deterministic)
}
func (dst *ProfMsg) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ProfMsg.Merge(dst, src)
}
func (m *ProfMsg) XXX_Size() int {
	return xxx_messageInfo_ProfMsg.Size(m)
}
func (m *ProfMsg) XXX_DiscardUnknown() {
	xxx_messageInfo_ProfMsg.DiscardUnknown(m)
}

var xxx_messageInfo_ProfMsg proto.InternalMessageInfo

func (m *ProfMsg) GetRequestID() []byte {
	if m != nil {
		return m.RequestID
	}
	return nil
}

func (m *ProfMsg) GetContent() []byte {
	if m != nil {
		return m.Content
	}
	return nil
}

func (m *ProfMsg) GetGroupID() []byte {
	if m != nil {
		return m.GroupID
	}
	return nil
}

func init() {
	proto.RegisterType((*ProfMsg)(nil), "prof.ProfMsg")
}

func init() { proto.RegisterFile("prof.proto", fileDescriptor_prof_043c76f9d6cc050b) }

var fileDescriptor_prof_043c76f9d6cc050b = []byte{
	// 109 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x2a, 0x28, 0xca, 0x4f,
	0xd3, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x01, 0xb1, 0x95, 0xa2, 0xb9, 0xd8, 0x03, 0x8a,
	0xf2, 0xd3, 0x7c, 0x8b, 0xd3, 0x85, 0x64, 0xb8, 0x38, 0x8b, 0x52, 0x0b, 0x4b, 0x53, 0x8b, 0x4b,
	0x3c, 0x5d, 0x24, 0x18, 0x15, 0x18, 0x35, 0x78, 0x82, 0x10, 0x02, 0x42, 0x12, 0x5c, 0xec, 0xc9,
	0xf9, 0x79, 0x25, 0xa9, 0x79, 0x25, 0x12, 0x4c, 0x60, 0x39, 0x18, 0x17, 0x24, 0x93, 0x5e, 0x94,
	0x5f, 0x5a, 0xe0, 0xe9, 0x22, 0xc1, 0x0c, 0x91, 0x81, 0x72, 0x93, 0xd8, 0xc0, 0x36, 0x19, 0x03,
	0x02, 0x00, 0x00, 0xff, 0xff, 0x17, 0xd5, 0x4b, 0x12, 0x77, 0x00, 0x00, 0x00,
}
