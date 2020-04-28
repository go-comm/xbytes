package xbytes

import (
	"unsafe"
)

func BytesToString(data []byte) string {
	return *(*string)(unsafe.Pointer(&data))
}

func StringToBytes(s *string) []byte {
	bs := (*[2]uintptr)(unsafe.Pointer(s))
	b := [3]uintptr{bs[0], bs[1], bs[1]}
	return *(*[]byte)(unsafe.Pointer(&b))
}
