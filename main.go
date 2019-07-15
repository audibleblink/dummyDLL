package main

import (
	"C"
	"fmt"
	"syscall"
)

//export DllMain
func DllMain() {
	alert()
}

//export DllCanUnloadNow
func DllCanUnloadNow() {
	alert()
}

//export DllGetClassObject
func DllGetClassObject() {
	alert()
}

//export DllRegisterServer
func DllRegisterServer() {
	alert()
}

//export DllUnregisterServer
func DllUnregisterServer() {
	alert()
}

//export FveCloseVolume
func FveCloseVolume() {
	alert()
}

//export FveOpenVolume
func FveOpenVolume() {
	alert()
}

//export FveGetStatus
func FveGetStatus() {
	alert()
}

//export FveDeleteDeviceEncryptionOptOutForVolumeW
func FveDeleteDeviceEncryptionOptOutForVolumeW() {
	alert()
}

//export FveCommitChanges
func FveCommitChanges() {
	alert()
}

//export FveConversionDecrypt
func FveConversionDecrypt() {
	alert()
}

//export FveDeleteAuthMethod
func FveDeleteAuthMethod() {
	alert()
}

//export FveGetAuthMethodInformation
func FveGetAuthMethodInformation() {
	alert()
}

//export FveRevertVolume
func FveRevertVolume() {
	alert()
}

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

func main() {}
