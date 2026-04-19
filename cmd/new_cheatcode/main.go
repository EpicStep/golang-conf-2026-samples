package main

import "C"
import (
	"hash/crc32"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

//export GoMain
func GoMain() {
	addr := uintptr(0x8A5CC8 + 0xC)
	value := crc32.ChecksumIEEE([]byte("GNALOG")) ^ 0xFFFFFFFF

	ptr := (*uint32)(unsafe.Pointer(addr))
	*ptr = value

	self, err := getSelfModule()
	if err != nil {
		return
	}

	cheatProc, err := windows.GetProcAddress(self, "Cheat")
	if err != nil {
		return
	}

	addr = uintptr(0x8A5B58 + 0xC)
	*(*uintptr)(unsafe.Pointer(addr)) = cheatProc

	windows.MessageBox(windows.HWND(0), syscall.StringToUTF16Ptr("Injected from Go"), syscall.StringToUTF16Ptr("Injection works"), windows.MB_OK)
}

//export Cheat
func Cheat() {
	textPtr, err := syscall.BytePtrFromString("Hello, GolangConf!")
	if err != nil {
		panic(err)
	}

	syscall.SyscallN(
		uintptr(0x69F2B0),
		uintptr(unsafe.Pointer(textPtr)),
		uintptr(5000),
		uintptr(0),
	)

	syscall.SyscallN(
		uintptr(0x43A0B0),
		uintptr(425),
	)
}

var moduleByte byte

func getSelfModule() (windows.Handle, error) {
	var h windows.Handle

	err := windows.GetModuleHandleEx(
		windows.GET_MODULE_HANDLE_EX_FLAG_FROM_ADDRESS|
			windows.GET_MODULE_HANDLE_EX_FLAG_UNCHANGED_REFCOUNT,
		(*uint16)(unsafe.Pointer(&moduleByte)),
		&h,
	)
	if err != nil {
		return 0, err
	}

	return h, nil
}

func main() {}
