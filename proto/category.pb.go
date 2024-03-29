// Code generated by protoc-gen-go. DO NOT EDIT.
// source: proto/category.proto

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

// Category defines the supported categories.
// Full list supported by Douban are:
// - 广播 broadcast (所有的状态更新)
// - 书籍 book
// - 电影 movie
// - 音乐 music
// - 游戏 game
// - 移动应用 app
// - 评论 review
// - 舞台剧 drama
// - 小组 group (not supported)
// - 日记 note (not supported)
// - 图片 album (not supported)
// - 小站 site (not supported)
// - 同城活动 activity (not supported)
// - 豆品 thing (not supported)
type Category int32

const (
	Category_broadcast Category = 0
	Category_book      Category = 1
	Category_movie     Category = 2
	Category_game      Category = 3
	Category_music     Category = 4
	Category_app       Category = 5
	Category_review    Category = 6
	Category_drama     Category = 7
)

var Category_name = map[int32]string{
	0: "broadcast",
	1: "book",
	2: "movie",
	3: "game",
	4: "music",
	5: "app",
	6: "review",
	7: "drama",
}

var Category_value = map[string]int32{
	"broadcast": 0,
	"book":      1,
	"movie":     2,
	"game":      3,
	"music":     4,
	"app":       5,
	"review":    6,
	"drama":     7,
}

func (x Category) String() string {
	return proto.EnumName(Category_name, int32(x))
}

func (Category) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_1bfb247fa9b1cc73, []int{0}
}

func init() {
	proto.RegisterEnum("proto.Category", Category_name, Category_value)
}

func init() { proto.RegisterFile("proto/category.proto", fileDescriptor_1bfb247fa9b1cc73) }

var fileDescriptor_1bfb247fa9b1cc73 = []byte{
	// 164 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x2c, 0x8e, 0x4d, 0x0a, 0xc2, 0x30,
	0x10, 0x46, 0xd5, 0xfe, 0x0f, 0x08, 0x43, 0xf0, 0x06, 0xe2, 0x46, 0x68, 0xb3, 0xf0, 0x06, 0x7a,
	0x0b, 0x77, 0x93, 0x1f, 0x6a, 0x28, 0x61, 0x4a, 0x9a, 0x54, 0x7a, 0x7b, 0x69, 0x70, 0xf5, 0xf8,
	0xde, 0xdb, 0x7c, 0x70, 0x99, 0x03, 0x47, 0x96, 0x9a, 0xa2, 0x1d, 0x39, 0x6c, 0x43, 0x9e, 0xa2,
	0xca, 0xb8, 0x6b, 0x68, 0x5f, 0xff, 0x20, 0xce, 0xd0, 0xa9, 0xc0, 0x64, 0x34, 0x2d, 0x11, 0x0f,
	0xa2, 0x85, 0x52, 0x31, 0x4f, 0x78, 0x14, 0x1d, 0x54, 0x9e, 0x57, 0x67, 0xf1, 0xb4, 0xcb, 0x91,
	0xbc, 0xc5, 0x22, 0xcb, 0xb4, 0x38, 0x8d, 0xa5, 0x68, 0xa0, 0xa0, 0x79, 0xc6, 0x4a, 0x00, 0xd4,
	0xc1, 0xae, 0xce, 0x7e, 0xb1, 0xde, 0xbb, 0x09, 0xe4, 0x09, 0x9b, 0xe7, 0xed, 0x7d, 0x1d, 0x5d,
	0xfc, 0x24, 0x35, 0x68, 0xf6, 0xd2, 0xc5, 0xa5, 0xf7, 0x5b, 0x6f, 0x28, 0x92, 0x34, 0x9c, 0x14,
	0x4d, 0x32, 0x7f, 0x51, 0x75, 0xc6, 0xe3, 0x17, 0x00, 0x00, 0xff, 0xff, 0x60, 0xa6, 0xb8, 0x02,
	0xb1, 0x00, 0x00, 0x00,
}
