//+build windows

package main

import (
	"syscall"
	"unsafe"
)

func setConsoleTitle() {
	mod := syscall.NewLazyDLL("kernel32.dll")
	proc := mod.NewProc("SetConsoleTitleW")
	i,_:=syscall.UTF16PtrFromString(title)
	proc.Call(uintptr(unsafe.Pointer(i)))
}