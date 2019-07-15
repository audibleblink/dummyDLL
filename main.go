package main

import (
	"C"
	"syscall"
)

func main() {
	alert()
}

func alert() {
	defer syscall.FreeLibrary(user32)
	name := getHostImagePath()
	MessageBox(name, caller(), MB_OK | MB_ICONEXCLAMATION | MB_TOPMOST)
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
