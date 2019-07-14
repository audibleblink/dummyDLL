package main

import (
	"github.com/gen2brain/beeep"
)

func main() {
	alert()
}

func alert() {
	name := getHostImagePath()
	beeep.Notify(name, caller(), "")
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
