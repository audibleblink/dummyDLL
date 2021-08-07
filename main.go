package main

//#include "dllmain.h"
import (
	"C"
	"fmt"
	"os/user"
	"runtime"
	"unsafe"

	"golang.org/x/sys/windows"
)

const template = `
FnCall: %s

WorkDir: %s

CmdLine: %s

Arch: %s

User: %s

Integrity: %s
`

//export OnProcessAttach
func OnProcessAttach(
	hinstDLL unsafe.Pointer, // handle to DLL module
	fdwReason uint32, // reason for calling function
	lpReserved unsafe.Pointer, // reserved
) {
	alert()
}

func alert() {
	imageName, path, cmdLine := hostingImageInfo()
	title := fmt.Sprintf("Host Image: %s", imageName)
	arch := runtime.GOARCH

	usr, err := user.Current()
	if err != nil {
		usr.Username = "Unknown Error"
	}

	integrity, err := getProcessIntegrityLevel()
	if err != nil {
		integrity = "Unknown Error"
	}

	msg := fmt.Sprintf(template, caller(), path, cmdLine, arch, usr.Username, integrity)
	MessageBox(title, msg, MB_OK|MB_ICONEXCLAMATION|MB_TOPMOST)
}

func hostingImageInfo() (imageName, path, cmdLine string) {
	peb := windows.RtlGetCurrentPeb()
	userProcParams := peb.ProcessParameters
	imageName = userProcParams.ImagePathName.String()
	path = userProcParams.CurrentDirectory.DosPath.String()
	cmdLine = userProcParams.CommandLine.String()
	return
}

//export CallNtPowerInformation
func CallNtPowerInformation() { alert() }

//export ClrCreateManagedInstance
func ClrCreateManagedInstance() { alert() }

//export ConstructPartialMsgVW
func ConstructPartialMsgVW() { alert() }

//export CorBindToRuntimeEx
func CorBindToRuntimeEx() { alert() }

//export CreateUri
func CreateUri() { alert() }

//export CurrentIP
func CurrentIP() { alert() }

//export DevObjCreateDeviceInfoList
func DevObjCreateDeviceInfoList() { alert() }

//export DevObjDestroyDeviceInfoList
func DevObjDestroyDeviceInfoList() { alert() }

//export DevObjEnumDeviceInterfaces
func DevObjEnumDeviceInterfaces() { alert() }

//export DevObjGetClassDevs
func DevObjGetClassDevs() { alert() }

//export DllCanUnloadNow
func DllCanUnloadNow() { alert() }

//export DllGetClassObject
func DllGetClassObject() { alert() }

//export DllProcessAttach
func DllProcessAttach() { alert() }

//export DevObjOpenDeviceInfo
func DevObjOpenDeviceInfo() { alert() }

//export DllRegisterServer
func DllRegisterServer() { alert() }

//export DllUnregisterServer
func DllUnregisterServer() { alert() }

//export DpxNewJob
func DpxNewJob() { alert() }

//export ExtractMachineName
func ExtractMachineName() { alert() }

//export FveCloseVolume
func FveCloseVolume() { alert() }

//export FveCommitChanges
func FveCommitChanges() { alert() }

//export FveConversionDecrypt
func FveConversionDecrypt() { alert() }

//export FveDeleteAuthMethod
func FveDeleteAuthMethod() { alert() }

//export FveDeleteDeviceEncryptionOptOutForVolumeW
func FveDeleteDeviceEncryptionOptOutForVolumeW() { alert() }

//export FveGetAuthMethodInformation
func FveGetAuthMethodInformation() { alert() }

//export FveGetStatus
func FveGetStatus() { alert() }

//export FveOpenVolume
func FveOpenVolume() { alert() }

//export FveRevertVolume
func FveRevertVolume() { alert() }

//export GenerateActionQueue
func GenerateActionQueue() { alert() }

//export GetFQDN_Ipv4
func GetFQDN_Ipv4() { alert() }

//export GetMemLogObject
func GetMemLogObject() { alert() }

//export GetFQDN_Ipv6
func GetFQDN_Ipv6() { alert() }

//export InitCommonControlsEx
func InitCommonControlsEx() { alert() }

//export IsLocalConnection
func IsLocalConnection() { alert() }

//export LoadLibraryShim
func LoadLibraryShim() { alert() }

//export NetApiBufferAllocate
func NetApiBufferAllocate() { alert() }

//export NetApiBufferFree
func NetApiBufferFree() { alert() }

//export NetApiBufferReallocate
func NetApiBufferReallocate() { alert() }

//export NetApiBufferSize
func NetApiBufferSize() { alert() }

//export NetRemoteComputerSupports
func NetRemoteComputerSupports() { alert() }

//export NetapipBufferAllocate
func NetapipBufferAllocate() { alert() }

//export NetpIsComputerNameValid
func NetpIsComputerNameValid() { alert() }

//export NetpIsDomainNameValid
func NetpIsDomainNameValid() { alert() }

//export NetpIsGroupNameValid
func NetpIsGroupNameValid() { alert() }

//export NetpIsRemote
func NetpIsRemote() { alert() }

//export NetpIsRemoteNameValid
func NetpIsRemoteNameValid() { alert() }

//export NetpIsShareNameValid
func NetpIsShareNameValid() { alert() }

//export NetpIsUncComputerNameValid
func NetpIsUncComputerNameValid() { alert() }

//export NetpIsUserNameValid
func NetpIsUserNameValid() { alert() }

//export NetpwListCanonicalize
func NetpwListCanonicalize() { alert() }

//export NetpwListTraverse
func NetpwListTraverse() { alert() }

//export NetpwNameCanonicalize
func NetpwNameCanonicalize() { alert() }

//export NetpwNameCompare
func NetpwNameCompare() { alert() }

//export NetpwNameValidate
func NetpwNameValidate() { alert() }

//export NetpwPathCanonicalize
func NetpwPathCanonicalize() { alert() }

//export NetpwPathCompare
func NetpwPathCompare() { alert() }

//export NetpwPathType
func NetpwPathType() { alert() }

//export PowerGetActiveScheme
func PowerGetActiveScheme() { alert() }

//export PrivateCoInternetCombineUri
func PrivateCoInternetCombineUri() { alert() }

//export ProcessActionQueue
func ProcessActionQueue() { alert() }

//export RegisterDLL
func RegisterDLL() { alert() }

//export Run
func Run() { alert() }

//export SLGetWindowsInformation
func SLGetWindowsInformation() { alert() }

//export UnRegisterDLL
func UnRegisterDLL() { alert() }

//export WdsAbortBlackboa
func WdsAbortBlackboa() { alert() }

func main() {}
