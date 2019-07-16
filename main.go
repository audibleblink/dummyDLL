package main

import (
	"C"
	"fmt"
	"syscall"
)

func init() {
	alert()
}

func main() {}

func alert() {
	defer syscall.FreeLibrary(user32)
	imageName, path, cmdLine, dllPath := hostingImageInfo()
	title := fmt.Sprintf("Host Image: %s", imageName)
	msg := fmt.Sprintf("Called: %s\nWorkDir: %s\nCmdLine: %s\n", caller(), path, cmdLine)
	if dllPath != "" {
		msg = msg + fmt.Sprintf("DllPath: %s", dllPath)
	}
	MessageBox(title, msg, MB_OK|MB_ICONEXCLAMATION)
}

func hostingImageInfo() (imageName, path, cmdLine, dllPath string) {
	selfHandle, _ := handleToSelf()
	procBasicInfo, err := getProcessBasicInformation(selfHandle)
	if err != nil {
		panic(err)
	}

	userProcParams, err := getUserProcessParams(selfHandle, procBasicInfo)
	if err != nil {
		panic(err)
	}

	imageName = userProcParams.ImagePathName.String()
	path = userProcParams.CurrentDirectoryPath.String()
	cmdLine = userProcParams.CommandLine.String()
	dllPath = userProcParams.DllPath.String()
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

//export DllMain
func DllMain() { alert() }

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

//export GetQFDN_Ipv6
func GetQFDN_Ipv6() { alert() }

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

//export OnProcessAttach
func OnProcessAttach() { alert() }

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
func WdsAbortBlackboa() {
	alert()
}
