// Code generated by protoc-gen-go. DO NOT EDIT.
// source: proto/exchangeRate.proto

package auth

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

type GetRateReq struct {
	BaseCurrency         string   `protobuf:"bytes,1,opt,name=BaseCurrency" json:"BaseCurrency,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetRateReq) Reset()         { *m = GetRateReq{} }
func (m *GetRateReq) String() string { return proto.CompactTextString(m) }
func (*GetRateReq) ProtoMessage()    {}
func (*GetRateReq) Descriptor() ([]byte, []int) {
	return fileDescriptor_exchangeRate_2c1ac57600542583, []int{0}
}
func (m *GetRateReq) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetRateReq.Unmarshal(m, b)
}
func (m *GetRateReq) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetRateReq.Marshal(b, m, deterministic)
}
func (dst *GetRateReq) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetRateReq.Merge(dst, src)
}
func (m *GetRateReq) XXX_Size() int {
	return xxx_messageInfo_GetRateReq.Size(m)
}
func (m *GetRateReq) XXX_DiscardUnknown() {
	xxx_messageInfo_GetRateReq.DiscardUnknown(m)
}

var xxx_messageInfo_GetRateReq proto.InternalMessageInfo

func (m *GetRateReq) GetBaseCurrency() string {
	if m != nil {
		return m.BaseCurrency
	}
	return ""
}

type GetRateResp struct {
	CurrencyRatio        string   `protobuf:"bytes,1,opt,name=CurrencyRatio" json:"CurrencyRatio,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetRateResp) Reset()         { *m = GetRateResp{} }
func (m *GetRateResp) String() string { return proto.CompactTextString(m) }
func (*GetRateResp) ProtoMessage()    {}
func (*GetRateResp) Descriptor() ([]byte, []int) {
	return fileDescriptor_exchangeRate_2c1ac57600542583, []int{1}
}
func (m *GetRateResp) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetRateResp.Unmarshal(m, b)
}
func (m *GetRateResp) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetRateResp.Marshal(b, m, deterministic)
}
func (dst *GetRateResp) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetRateResp.Merge(dst, src)
}
func (m *GetRateResp) XXX_Size() int {
	return xxx_messageInfo_GetRateResp.Size(m)
}
func (m *GetRateResp) XXX_DiscardUnknown() {
	xxx_messageInfo_GetRateResp.DiscardUnknown(m)
}

var xxx_messageInfo_GetRateResp proto.InternalMessageInfo

func (m *GetRateResp) GetCurrencyRatio() string {
	if m != nil {
		return m.CurrencyRatio
	}
	return ""
}

func init() {
	proto.RegisterType((*GetRateReq)(nil), "auth.GetRateReq")
	proto.RegisterType((*GetRateResp)(nil), "auth.GetRateResp")
}

func init() {
	proto.RegisterFile("proto/exchangeRate.proto", fileDescriptor_exchangeRate_2c1ac57600542583)
}

var fileDescriptor_exchangeRate_2c1ac57600542583 = []byte{
	// 146 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x92, 0x28, 0x28, 0xca, 0x2f,
	0xc9, 0xd7, 0x4f, 0xad, 0x48, 0xce, 0x48, 0xcc, 0x4b, 0x4f, 0x0d, 0x4a, 0x2c, 0x49, 0xd5, 0x03,
	0x0b, 0x09, 0xb1, 0x24, 0x96, 0x96, 0x64, 0x28, 0x19, 0x70, 0x71, 0xb9, 0xa7, 0x96, 0x80, 0x84,
	0x83, 0x52, 0x0b, 0x85, 0x94, 0xb8, 0x78, 0x9c, 0x12, 0x8b, 0x53, 0x9d, 0x4b, 0x8b, 0x8a, 0x52,
	0xf3, 0x92, 0x2b, 0x25, 0x18, 0x15, 0x18, 0x35, 0x38, 0x83, 0x50, 0xc4, 0x94, 0x8c, 0xb9, 0xb8,
	0xe1, 0x3a, 0x8a, 0x0b, 0x84, 0x54, 0xb8, 0x78, 0x61, 0x52, 0x41, 0x89, 0x25, 0x99, 0xf9, 0x50,
	0x3d, 0xa8, 0x82, 0x46, 0x36, 0x5c, 0x1c, 0x30, 0x01, 0x21, 0x03, 0x2e, 0x76, 0xa8, 0x01, 0x42,
	0x02, 0x7a, 0x20, 0x47, 0xe8, 0x21, 0x5c, 0x20, 0x25, 0x88, 0x26, 0x52, 0x5c, 0xa0, 0xc4, 0x90,
	0xc4, 0x06, 0x76, 0xb1, 0x31, 0x20, 0x00, 0x00, 0xff, 0xff, 0x69, 0x4f, 0xe6, 0xe9, 0xcd, 0x00,
	0x00, 0x00,
}
