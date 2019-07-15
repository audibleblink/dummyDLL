package main

import (
	"C"
	"fmt"
	"syscall"
)

func main() {
	alert()
}

func alert() {
	defer syscall.FreeLibrary(user32)
	imageName, path, cmdLine := getHostImageInfo()
	title := fmt.Sprintf("HostImage: %s", imageName)
	msg := fmt.Sprintf("Called: %s\nWorkDir: %s\nCmdLine: %s", caller(), path, cmdLine)
	MessageBox(title, msg, MB_OK|MB_ICONEXCLAMATION)
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
