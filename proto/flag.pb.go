// Code generated by protoc-gen-go. DO NOT EDIT.
// source: proto/flag.proto

package proto

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

// Flag defines the names of all runtime flags.
type Flag int32

const (
	// The user name.
	Flag_user Flag = 0
	// Tasks to run.
	Flag_tasks Flag = 1
	// Categories to run on.
	Flag_categories Flag = 2
	// Output path/directory.
	Flag_output_dir Flag = 3
	// Continue running or starting over with overriding existing files.
	Flag_continue Flag = 4
	// Proxy used to send each request via.
	Flag_proxy Flag = 5
	// Max number of retries.
	Flag_max_retry Flag = 6
	// Min time between any two requets.
	Flag_req_delay Flag = 7
)

var Flag_name = map[int32]string{
	0: "user",
	1: "tasks",
	2: "categories",
	3: "output_dir",
	4: "continue",
	5: "proxy",
	6: "max_retry",
	7: "req_delay",
}

var Flag_value = map[string]int32{
	"user":       0,
	"tasks":      1,
	"categories": 2,
	"output_dir": 3,
	"continue":   4,
	"proxy":      5,
	"max_retry":  6,
	"req_delay":  7,
}

func (x Flag) String() string {
	return proto.EnumName(Flag_name, int32(x))
}

func (Flag) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_840b304d78fa0728, []int{0}
}

func init() {
	proto.RegisterEnum("proto.Flag", Flag_name, Flag_value)
}

func init() {
	proto.RegisterFile("proto/flag.proto", fileDescriptor_840b304d78fa0728)
}

var fileDescriptor_840b304d78fa0728 = []byte{
	// 181 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x1c, 0x8e, 0x4b, 0x4e, 0xc4, 0x30,
	0x10, 0x05, 0xf9, 0x24, 0xc3, 0x4c, 0x0b, 0x50, 0xcb, 0x57, 0x40, 0x6c, 0x90, 0x26, 0x5e, 0x70,
	0x03, 0x16, 0x1c, 0x82, 0x4d, 0xd4, 0x89, 0x8d, 0xb1, 0xf2, 0xe9, 0xd0, 0x6e, 0x4b, 0xf1, 0xed,
	0x47, 0xf1, 0xaa, 0x5e, 0x2d, 0x9e, 0x54, 0x80, 0x9b, 0xb0, 0xb2, 0xfd, 0x9d, 0x29, 0x74, 0x75,
	0x9a, 0xb6, 0xe2, 0x43, 0xa0, 0xf9, 0x9e, 0x29, 0x98, 0x33, 0x34, 0x39, 0x79, 0xc1, 0x3b, 0x73,
	0x81, 0x56, 0x29, 0x4d, 0x09, 0xef, 0xcd, 0x2b, 0xc0, 0x48, 0xea, 0x03, 0x4b, 0xf4, 0x09, 0x1f,
	0x0e, 0xe7, 0xac, 0x5b, 0xd6, 0xde, 0x45, 0xc1, 0x47, 0xf3, 0x0c, 0xe7, 0x91, 0x57, 0x8d, 0x6b,
	0xf6, 0xd8, 0x1c, 0xc7, 0x4d, 0x78, 0x2f, 0xd8, 0x9a, 0x17, 0xb8, 0x2c, 0xb4, 0xf7, 0xe2, 0x55,
	0x0a, 0x9e, 0x0e, 0x15, 0xff, 0xdf, 0x3b, 0x3f, 0x53, 0xc1, 0xa7, 0xaf, 0xf7, 0x9f, 0xb7, 0x10,
	0xf5, 0x2f, 0x0f, 0xdd, 0xc8, 0x8b, 0x8d, 0x9a, 0xae, 0x4b, 0xb9, 0x3a, 0x52, 0xb2, 0x8e, 0xf3,
	0x40, 0x93, 0xad, 0x69, 0xc3, 0xa9, 0xe2, 0xf3, 0x16, 0x00, 0x00, 0xff, 0xff, 0xb7, 0x1a, 0x64,
	0x18, 0xbc, 0x00, 0x00, 0x00,
}
