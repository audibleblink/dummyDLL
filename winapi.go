package main

import (
	"C"
	"fmt"
	"os"
	"syscall"
	"unsafe"

	"github.com/elastic/go-windows"
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
	processCreateThread = 2
	processVMOperation  = 8
	processVMWrite      = 32

	is32bitProc        = unsafe.Sizeof(uintptr(0)) == 4
	userProcParamsSize = unsafe.Sizeof(rtlUserProcessParameters{})
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

type memoryBasicInformation struct {
	BaseAddress       uintptr
	AllocationBase    uintptr
	AllocationProtect uint32
	RegionSize        uintptr
	State             uint32
	Protect           uint32
	Type              uint32
}

type rtlUserProcessParameters struct {
	Reserved1 [16]byte
	Reserved2 [5]uintptr

	// <undocumented>
	CurrentDirectoryPath   unicodeString
	CurrentDirectoryHandle uintptr
	DllPath                unicodeString
	// </undocumented>

	ImagePathName unicodeString
	CommandLine   unicodeString
	Environment   uintptr
}

// https://github.com/shenwei356/rush/blob/3699d8775d5f4d429351700fea4231de0ec1e281/process/process_windows.go#L259
func getUserProcessParams(
	processHandle syscall.Handle,
	pbi windows.ProcessBasicInformationStruct) (rtlUserProcessParameters, error) {

	// Offset of params field within PEB structure.
	// This structure is different in 32 and 64 bit.
	paramsOffset := 0x20
	if is32bitProc {
		paramsOffset = 0x10
	}

	// Read the PEB from the target process memory
	pebSize := paramsOffset + 8
	peb := make([]byte, pebSize)

	nRead, err := readMem(processHandle, pbi.PebBaseAddress, peb)
	var userProcParams rtlUserProcessParameters
	if err != nil {
		return userProcParams, err
	}

	if nRead != uintptr(pebSize) {
		return userProcParams, fmt.Errorf("PEB: short read (%d/%d)", nRead, pebSize)
	}

	// Get the RTL_USER_PROCESS_PARAMETERS struct pointer from the PEB
	paramsAddr := *(*uintptr)(unsafe.Pointer(&peb[paramsOffset]))

	// Read the RTL_USER_PROCESS_PARAMETERS from the target process memory
	paramsBuf := make([]byte, userProcParamsSize)
	nRead, err = readMem(processHandle, paramsAddr, paramsBuf)
	if err != nil {
		return userProcParams, err
	}

	if nRead != uintptr(userProcParamsSize) {
		return userProcParams, fmt.Errorf("RTL_USER_PROCESS_PARAMETERS: short read (%d/%d)",
			nRead,
			userProcParamsSize)
	}

	userProcParams = *(*rtlUserProcessParameters)(unsafe.Pointer(&paramsBuf[0]))
	return userProcParams, err
}

func handleToSelf() (handle syscall.Handle, err error) {
	var attrs uint32 = processCreateThread | syscall.PROCESS_QUERY_INFORMATION |
		processVMOperation | processVMWrite | windows.PROCESS_VM_READ
	pid := uint32(os.Getpid())
	handle, err = syscall.OpenProcess(attrs, false, pid)
	return
}

// https://github.com/shenwei356/rush/blob/3699d8775d5f4d429351700fea4231de0ec1e281/process/process_windows.go#L251
func getProcessBasicInformation(processHandle syscall.Handle) (pbi windows.ProcessBasicInformationStruct, err error) {

	actualSize, err := windows.NtQueryInformationProcess(
		processHandle,
		windows.ProcessBasicInformation,
		unsafe.Pointer(&pbi),
		uint32(windows.SizeOfProcessBasicInformationStruct))

	if actualSize < uint32(windows.SizeOfProcessBasicInformationStruct) {
		return pbi, fmt.Errorf("Bad size for PROCESS_BASIC_INFORMATION")
	}
	return pbi, err
}

// https://github.com/shenwei356/rush/blob/3699d8775d5f4d429351700fea4231de0ec1e281/process/process_windows.go#L232
func readMem(handle syscall.Handle, baseAddress uintptr, dest []byte) (numRead uintptr, err error) {
	n := len(dest)
	if n == 0 {
		return 0, nil
	}

	ntReadVirtualMemory.Call(
		uintptr(handle),
		uintptr(baseAddress),
		uintptr(unsafe.Pointer(&dest[0])),
		uintptr(n),
		uintptr(unsafe.Pointer(&numRead)))
	return
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
