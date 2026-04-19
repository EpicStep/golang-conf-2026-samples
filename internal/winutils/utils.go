package winutils

import (
	"fmt"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

var (
	kernel32 = windows.NewLazySystemDLL("kernel32.dll")
	user32   = windows.NewLazySystemDLL("user32.dll")

	procVirtualAllocEx     = kernel32.NewProc("VirtualAllocEx")
	procCreateRemoteThread = kernel32.NewProc("CreateRemoteThread")
	procFindWindow         = user32.NewProc("FindWindowW")
)

func FindWindow(name string) (windows.HWND, error) {
	nameUTF16, err := syscall.UTF16PtrFromString(name)
	if err != nil {
		return 0, fmt.Errorf("syscall.UTF16PtrFromString: %w", err)
	}

	r1, _, e1 := syscall.SyscallN(procFindWindow.Addr(), uintptr(unsafe.Pointer(nil)), uintptr(unsafe.Pointer(nameUTF16)))
	if r1 == 0 {
		return 0, fmt.Errorf("syscall.SyscallN: %w", error(e1))
	}

	return windows.HWND(r1), nil
}

func VirtualAllocEx(process windows.Handle, address uintptr, size uintptr, allocationType uint32, protect uint32) (uintptr, error) {
	r1, _, e1 := syscall.SyscallN(procVirtualAllocEx.Addr(), uintptr(process), address, size, uintptr(allocationType), uintptr(protect))
	if r1 == 0 {
		return 0, fmt.Errorf("syscall.SyscallN: %w", error(e1))
	}

	return r1, nil
}

func CreateRemoteThread(process windows.Handle, startAddress uintptr, param uintptr) (windows.Handle, error) {
	r1, _, e1 := syscall.SyscallN(procCreateRemoteThread.Addr(), uintptr(process), 0, 0, startAddress, param, 0, 0)
	if r1 == 0 {
		return 0, fmt.Errorf("syscall.SyscallN: %w", error(e1))
	}

	return windows.Handle(r1), nil
}
