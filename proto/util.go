package proto

import (
	"strings"
)

// ConcatProtoEnum concats enum proto values to a string.
func ConcatProtoEnum(nameMap map[int32]string, separator string) string {
	list := []string{}
	for _, v := range nameMap {
		list = append(list, v)
	}
	return strings.Join(list, separator)
}
