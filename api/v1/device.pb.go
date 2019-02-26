// Code generated by protoc-gen-go. DO NOT EDIT.
// source: github.com/antha-lang/antha/api/v1/device.proto

package org_antha_lang_antha_v1

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

type DeviceMetadata struct {
	Tags []string `protobuf:"bytes,1,rep,name=tags" json:"tags,omitempty"`
}

func (m *DeviceMetadata) Reset()                    { *m = DeviceMetadata{} }
func (m *DeviceMetadata) String() string            { return proto.CompactTextString(m) }
func (*DeviceMetadata) ProtoMessage()               {}
func (*DeviceMetadata) Descriptor() ([]byte, []int) { return fileDescriptor11, []int{0} }

func (m *DeviceMetadata) GetTags() []string {
	if m != nil {
		return m.Tags
	}
	return nil
}

func init() {
	proto.RegisterType((*DeviceMetadata)(nil), "org.antha_lang.antha.v1.DeviceMetadata")
}

func init() { proto.RegisterFile("github.com/antha-lang/antha/api/v1/device.proto", fileDescriptor11) }

var fileDescriptor11 = []byte{
	// 121 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xd2, 0x4f, 0xcf, 0x2c, 0xc9,
	0x28, 0x4d, 0xd2, 0x4b, 0xce, 0xcf, 0xd5, 0x4f, 0xcc, 0x2b, 0xc9, 0x48, 0xd4, 0xcd, 0x49, 0xcc,
	0x4b, 0x87, 0x30, 0xf5, 0x13, 0x0b, 0x32, 0xf5, 0xcb, 0x0c, 0xf5, 0x53, 0x52, 0xcb, 0x32, 0x93,
	0x53, 0xf5, 0x0a, 0x8a, 0xf2, 0x4b, 0xf2, 0x85, 0xc4, 0xf3, 0x8b, 0xd2, 0xf5, 0xc0, 0xd2, 0xf1,
	0x20, 0x95, 0x10, 0xa6, 0x5e, 0x99, 0xa1, 0x92, 0x0a, 0x17, 0x9f, 0x0b, 0x58, 0xa1, 0x6f, 0x6a,
	0x49, 0x62, 0x4a, 0x62, 0x49, 0xa2, 0x90, 0x10, 0x17, 0x4b, 0x49, 0x62, 0x7a, 0xb1, 0x04, 0xa3,
	0x02, 0xb3, 0x06, 0x67, 0x10, 0x98, 0x9d, 0xc4, 0x06, 0x36, 0xc5, 0x18, 0x10, 0x00, 0x00, 0xff,
	0xff, 0x3a, 0x64, 0x72, 0x0a, 0x78, 0x00, 0x00, 0x00,
}