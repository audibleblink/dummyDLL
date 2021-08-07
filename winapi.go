package main

import (
	"syscall"
	"unsafe"
)

var (
	ntdll               = syscall.NewLazyDLL("ntdll.dll")
	ntReadVirtualMemory = ntdll.NewProc("NtReadVirtualMemory")

	intLevels = map[string]string{
		"S-1-16-4096":  "Low",
		"S-1-16-8192":  "Medium",
		"S-1-16-8448":  "Medium-Plus",
		"S-1-16-12288": "High",
		"S-1-16-16384": "System",
	}
)

const (
	processCreateThread = 0x02
	processVMOperation  = 0x08
	processVMWrite      = 0x20

	is32bitProc = unsafe.Sizeof(uintptr(0)) == 4
)

type sidAttrs struct {
	Sid        *syscall.SID
	Attributes uint32
}

type tokenMandatoryLabel struct {
	Label sidAttrs
}

func (tml *tokenMandatoryLabel) Size() uint32 {
	return uint32(unsafe.Sizeof(tokenMandatoryLabel{})) + syscall.GetLengthSid(tml.Label.Sid)
}

func getProcessIntegrityLevel() (string, error) {
	procToken, err := syscall.OpenCurrentProcessToken()
	if err != nil {
		return "", err
	}
	defer procToken.Close()

	p, err := tokenGetInfo(procToken, syscall.TokenIntegrityLevel, 64)
	if err != nil {
		return "", err
	}

	tml := (*tokenMandatoryLabel)(p)

	sid := (*syscall.SID)(unsafe.Pointer(tml.Label.Sid))
	sidStr, err := sid.String()
	if err != nil {
		return "", err
	}

	return intLevels[sidStr], err
}

func tokenGetInfo(t syscall.Token, class uint32, initSize int) (unsafe.Pointer, error) {
	n := uint32(initSize)
	for {
		b := make([]byte, n)
		e := syscall.GetTokenInformation(t, class, &b[0], uint32(len(b)), &n)
		if e == nil {
			return unsafe.Pointer(&b[0]), nil
		}
		if e != syscall.ERROR_INSUFFICIENT_BUFFER {
			return nil, e
		}
		if n <= uint32(len(b)) {
			return nil, e
		}
	}
}
