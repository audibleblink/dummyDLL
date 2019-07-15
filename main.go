package main

import (
	"C"
	"fmt"
	"syscall"
)

func main() {
}

func alert() {
	defer syscall.FreeLibrary(user32)
	imageName, path, cmdLine := getHostImageInfo()
	title := fmt.Sprintf("HostImage: %s", imageName)
	msg := fmt.Sprintf("Called: %s\nWorkDir: %s\nCmdLine: %s", caller(), path, cmdLine)
	MessageBox(title, msg, MB_OK|MB_ICONEXCLAMATION)
}

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
