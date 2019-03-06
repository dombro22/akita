// Code generated by protoc-gen-go. DO NOT EDIT.
// source: syncoffset.proto

package db

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

type SyncOffset struct {
	Offset               int64    `protobuf:"varint,2,opt,name=offset,proto3" json:"offset,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SyncOffset) Reset()         { *m = SyncOffset{} }
func (m *SyncOffset) String() string { return proto.CompactTextString(m) }
func (*SyncOffset) ProtoMessage()    {}
func (*SyncOffset) Descriptor() ([]byte, []int) {
	return fileDescriptor_37ffb859e0585169, []int{0}
}

func (m *SyncOffset) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SyncOffset.Unmarshal(m, b)
}
func (m *SyncOffset) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SyncOffset.Marshal(b, m, deterministic)
}
func (m *SyncOffset) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SyncOffset.Merge(m, src)
}
func (m *SyncOffset) XXX_Size() int {
	return xxx_messageInfo_SyncOffset.Size(m)
}
func (m *SyncOffset) XXX_DiscardUnknown() {
	xxx_messageInfo_SyncOffset.DiscardUnknown(m)
}

var xxx_messageInfo_SyncOffset proto.InternalMessageInfo

func (m *SyncOffset) GetOffset() int64 {
	if m != nil {
		return m.Offset
	}
	return 0
}

func init() {
	proto.RegisterType((*SyncOffset)(nil), "db.SyncOffset")
}

func init() { proto.RegisterFile("syncoffset.proto", fileDescriptor_37ffb859e0585169) }

var fileDescriptor_37ffb859e0585169 = []byte{
	// 77 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x12, 0x28, 0xae, 0xcc, 0x4b,
	0xce, 0x4f, 0x4b, 0x2b, 0x4e, 0x2d, 0xd1, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x4a, 0x49,
	0x52, 0x52, 0xe1, 0xe2, 0x0a, 0xae, 0xcc, 0x4b, 0xf6, 0x07, 0x8b, 0x0b, 0x89, 0x71, 0xb1, 0x41,
	0x54, 0x48, 0x30, 0x29, 0x30, 0x6a, 0x30, 0x07, 0x41, 0x79, 0x49, 0x6c, 0x60, 0x0d, 0xc6, 0x80,
	0x00, 0x00, 0x00, 0xff, 0xff, 0x04, 0xed, 0xb1, 0x02, 0x44, 0x00, 0x00, 0x00,
}
