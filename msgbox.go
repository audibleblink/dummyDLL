package main

import (
	"os"
	"syscall"
	"unsafe"
)

var (
	user32, _     = syscall.LoadLibrary("user32.dll")
	messageBox, _ = syscall.GetProcAddress(user32, "MessageBoxW")
)

const (
	MB_OK              = 0x00000000
	MB_ICONEXCLAMATION = 0x00000030
	MB_SETFOREGROUND   = 0x00010000
)

func MessageBox(caption, text string, style uintptr) (result int) {
	var nargs uintptr = 4
	ret, _, callErr := syscall.Syscall9(uintptr(messageBox),
		nargs,
		0,
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(text))),
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(caption))),
		style,
		0,
		0,
		0,
		0,
		0)
	if callErr != 0 {
		os.Exit(1)
	}
	result = int(ret)
	return
}
