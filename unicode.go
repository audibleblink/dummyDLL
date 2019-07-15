package main

import (
	"C"
	"unsafe"
)

type unicodeString struct {
	Size          uint16
	MaximumLength uint16
	Buffer        uintptr
}

func (u *unicodeString) String() string {
	//https://stackoverflow.com/questions/43766471/accessing-a-memory-address-from-a-string-in-go
	//https://forum.golangbridge.org/t/nice-way-to-convert-int-to-c-int/5351/3
	dataPtr := unsafe.Pointer(u.Buffer)
	dataSize := C.int(u.Size)
	data := C.GoBytes(dataPtr, dataSize)
	return string(_iconv(data))
}

// naively convert utf16
func _iconv(utf16Bytes []byte) (utf8Bytes []byte) {
	for _, bite := range utf16Bytes {
		if bite > 0 {
			utf8Bytes = append(utf8Bytes, bite)
		}
	}
	return
}
