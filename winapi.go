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
)

const (
	processCreateThread = 2
	processVMOperation  = 8
	processVMWrite      = 32

	userProcParamsSize = unsafe.Sizeof(rtlUserProcessParameters{})
)

// MemoryBasicInformation is Go's equivalent for the
// _MEMORY_BASIC_INFORMATION struct.
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

	var userProcParams rtlUserProcessParameters
	const is32bitProc = unsafe.Sizeof(uintptr(0)) == 4

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
