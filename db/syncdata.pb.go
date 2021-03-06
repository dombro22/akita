// Code generated by protoc-gen-go. DO NOT EDIT.
// source: syncdata.proto

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

type SyncData struct {
	Code                 int32    `protobuf:"varint,8,opt,name=Code,proto3" json:"Code,omitempty"`
	Data                 []byte   `protobuf:"bytes,7,opt,name=Data,proto3" json:"Data,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SyncData) Reset()         { *m = SyncData{} }
func (m *SyncData) String() string { return proto.CompactTextString(m) }
func (*SyncData) ProtoMessage()    {}
func (*SyncData) Descriptor() ([]byte, []int) {
	return fileDescriptor_0825ce093e9e0b31, []int{0}
}

func (m *SyncData) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SyncData.Unmarshal(m, b)
}
func (m *SyncData) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SyncData.Marshal(b, m, deterministic)
}
func (m *SyncData) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SyncData.Merge(m, src)
}
func (m *SyncData) XXX_Size() int {
	return xxx_messageInfo_SyncData.Size(m)
}
func (m *SyncData) XXX_DiscardUnknown() {
	xxx_messageInfo_SyncData.DiscardUnknown(m)
}

var xxx_messageInfo_SyncData proto.InternalMessageInfo

func (m *SyncData) GetCode() int32 {
	if m != nil {
		return m.Code
	}
	return 0
}

func (m *SyncData) GetData() []byte {
	if m != nil {
		return m.Data
	}
	return nil
}

func init() {
	proto.RegisterType((*SyncData)(nil), "db.SyncData")
}

func init() { proto.RegisterFile("syncdata.proto", fileDescriptor_0825ce093e9e0b31) }

var fileDescriptor_0825ce093e9e0b31 = []byte{
	// 90 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x2b, 0xae, 0xcc, 0x4b,
	0x4e, 0x49, 0x2c, 0x49, 0xd4, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x4a, 0x49, 0x52, 0x32,
	0xe2, 0xe2, 0x08, 0xae, 0xcc, 0x4b, 0x76, 0x49, 0x2c, 0x49, 0x14, 0x12, 0xe2, 0x62, 0x71, 0xce,
	0x4f, 0x49, 0x95, 0xe0, 0x50, 0x60, 0xd4, 0x60, 0x0d, 0x02, 0xb3, 0x41, 0x62, 0x20, 0x39, 0x09,
	0x76, 0x05, 0x46, 0x0d, 0x9e, 0x20, 0x30, 0x3b, 0x89, 0x0d, 0xac, 0xdd, 0x18, 0x10, 0x00, 0x00,
	0xff, 0xff, 0x78, 0x94, 0xff, 0x9a, 0x50, 0x00, 0x00, 0x00,
}
