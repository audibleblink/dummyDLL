package main

import (
	"C"
	"fmt"
	"os"
	"runtime"
	"syscall"
	"unsafe"

	"github.com/elastic/go-windows"
)

var (
	ntdll                   = syscall.NewLazyDLL("ntdll.dll")
	procNtReadVirtualMemory = ntdll.NewProc("NtReadVirtualMemory")
)

const (
	processCreateThread = 2
	processVMOperation  = 8
	processVMWrite      = 32

	sizeOfRtlUserProcessParameters = unsafe.Sizeof(rtlUserProcessParameters{})
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

func getHostImageInfo() (imageName, path, cmdLine, dllPath string) {
	winProcessHandle, _ := openProcess()
	pbi, err := getProcessBasicInformation(winProcessHandle)
	if err != nil {
		panic(err)
	}

	params, err := getUserProcessParams(winProcessHandle, pbi)
	if err != nil {
		panic(err)
	}

	imageName = params.ImagePathName.String()
	path = params.CurrentDirectoryPath.String()
	cmdLine = params.CommandLine.String()
	dllPath = params.DllPath.String()
	return imageName, path, cmdLine, dllPath
}

// returns the name of the function 2 levels up in the call stack
func caller() string {
	// Skip `caller` and the function to get the caller of
	return getFrame(2).Function
}

// https://stackoverflow.com/questions/35212985/is-it-possible-get-information-about-caller-function-in-golang
func getFrame(skipFrames int) runtime.Frame {
	// We need the frame at index skipFrames+2, since we never want runtime.Callers and getFrame
	targetFrameIndex := skipFrames + 2

	// Set size to targetFrameIndex+2 to ensure we have room for one more caller than we need
	programCounters := make([]uintptr, targetFrameIndex+2)
	n := runtime.Callers(0, programCounters)

	frame := runtime.Frame{Function: "unknown"}
	if n > 0 {
		frames := runtime.CallersFrames(programCounters[:n])
		for more, frameIndex := true, 0; more && frameIndex <= targetFrameIndex; frameIndex++ {
			var frameCandidate runtime.Frame
			frameCandidate, more = frames.Next()
			if frameIndex == targetFrameIndex {
				frame = frameCandidate
			}
		}
	}

	return frame
}

// https://github.com/shenwei356/rush/blob/3699d8775d5f4d429351700fea4231de0ec1e281/process/process_windows.go#L259
func getUserProcessParams(
	processHandle syscall.Handle,
	pbi windows.ProcessBasicInformationStruct) (params rtlUserProcessParameters, err error) {

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

	nRead, err := ntReadVirtualMemory(processHandle, pbi.PebBaseAddress, peb)
	if err != nil {
		return params, err
	}

	if nRead != uintptr(pebSize) {
		return params, fmt.Errorf("PEB: short read (%d/%d)", nRead, pebSize)
	}

	// Get the RTL_USER_PROCESS_PARAMETERS struct pointer from the PEB
	paramsAddr := *(*uintptr)(unsafe.Pointer(&peb[paramsOffset]))

	// Read the RTL_USER_PROCESS_PARAMETERS from the target process memory
	paramsBuf := make([]byte, sizeOfRtlUserProcessParameters)
	nRead, err = ntReadVirtualMemory(processHandle, paramsAddr, paramsBuf)
	if err != nil {
		return params, err
	}
	if nRead != uintptr(sizeOfRtlUserProcessParameters) {
		return params, fmt.Errorf("RTL_USER_PROCESS_PARAMETERS: short read (%d/%d)", nRead, sizeOfRtlUserProcessParameters)
	}

	params = *(*rtlUserProcessParameters)(unsafe.Pointer(&paramsBuf[0]))
	return params, nil
}

func openProcess() (handle syscall.Handle, err error) {
	var da uint32 = processCreateThread | syscall.PROCESS_QUERY_INFORMATION |
		processVMOperation | processVMWrite | windows.PROCESS_VM_READ
	pid := uint32(os.Getpid())
	handle, err = syscall.OpenProcess(da, false, pid)
	return
}

// https://github.com/shenwei356/rush/blob/3699d8775d5f4d429351700fea4231de0ec1e281/process/process_windows.go#L251
func getProcessBasicInformation(processHandle syscall.Handle) (pbi windows.ProcessBasicInformationStruct, err error) {
	actualSize, err := windows.NtQueryInformationProcess(processHandle, windows.ProcessBasicInformation, unsafe.Pointer(&pbi), uint32(windows.SizeOfProcessBasicInformationStruct))
	if actualSize < uint32(windows.SizeOfProcessBasicInformationStruct) {
		return pbi, fmt.Errorf("Bad size for PROCESS_BASIC_INFORMATION")
	}
	return pbi, err
}

func _ntReadVirtualMemory(handle syscall.Handle, baseAddress uintptr, buffer uintptr, size uintptr, numRead *uintptr) (ntStatus uint32) {
	r0, _, _ := procNtReadVirtualMemory.Call(
		uintptr(handle), uintptr(baseAddress), uintptr(buffer), uintptr(size), uintptr(unsafe.Pointer(numRead)))
	ntStatus = uint32(r0)
	return
}

// https://github.com/shenwei356/rush/blob/3699d8775d5f4d429351700fea4231de0ec1e281/process/process_windows.go#L232
func ntReadVirtualMemory(handle syscall.Handle, baseAddress uintptr, dest []byte) (numRead uintptr, err error) {
	n := len(dest)
	if n == 0 {
		return 0, nil
	}
	status := _ntReadVirtualMemory(handle, baseAddress, uintptr(unsafe.Pointer(&dest[0])), uintptr(n), &numRead)
	if status != 0 {
		return numRead, fmt.Errorf("*kaboom noises*")
	}
	return numRead, nil
}

// naively convert utf16 to utf8
func _iconv(utf16Bytes []byte) (utf8Bytes []byte) {
	for _, bite := range utf16Bytes {
		if bite > 0 {
			utf8Bytes = append(utf8Bytes, bite)
		}
	}
	return
}
