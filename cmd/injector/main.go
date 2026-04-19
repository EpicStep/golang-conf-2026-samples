package main

import (
	"flag"
	"log"
	"syscall"

	"golang.org/x/sys/windows"

	"github.com/EpicStep/golang-conf-2026-samples/internal/winutils"
)

var (
	pathToDLL   = flag.String("dll", "", "Path to DLL file")
	windowName  = flag.String("window", "", "Name of window to inject")
	dllProcName = flag.String("dll-proc", "GoMain", "Name of DLL procedure to inject")
)

var (
	kernel32 = syscall.NewLazyDLL("kernel32.dll")
)

func main() {
	flag.Parse()

	hwnd, err := winutils.FindWindow(*windowName)
	if err != nil {
		log.Fatal(err)
	}

	var pid uint32
	if _, err = windows.GetWindowThreadProcessId(hwnd, &pid); err != nil {
		log.Fatal(err)
	}

	handle, err := windows.OpenProcess(
		windows.PROCESS_CREATE_THREAD|windows.PROCESS_VM_OPERATION|windows.PROCESS_VM_WRITE|
			windows.PROCESS_VM_READ|windows.PROCESS_QUERY_INFORMATION,
		false, pid)
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		if err = windows.CloseHandle(handle); err != nil {
			log.Println(err)
		}
	}()

	vAddr, err := winutils.VirtualAllocEx(handle, 0, uintptr(len(*pathToDLL)), windows.MEM_RESERVE|windows.MEM_COMMIT, windows.PAGE_EXECUTE_READWRITE)
	if err != nil {
		log.Fatal(err)
	}

	pathToDLLBytePtr, err := windows.BytePtrFromString(*pathToDLL)
	if err != nil {
		log.Fatal(err)
	}

	err = windows.WriteProcessMemory(handle, vAddr, pathToDLLBytePtr, uintptr(len(*pathToDLL)), nil)
	if err != nil {
		log.Fatal(err)
	}

	loadLibAddr, err := syscall.GetProcAddress(syscall.Handle(kernel32.Handle()), "LoadLibraryA")
	if err != nil {
		log.Fatal(err)
	}

	tHandle, err := winutils.CreateRemoteThread(handle, loadLibAddr, vAddr)
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		if err = windows.CloseHandle(tHandle); err != nil {
			log.Println(err)
		}
	}()

	localDLL, err := windows.LoadLibrary(*pathToDLL)
	if err != nil {
		log.Fatal(err)
	}

	address, err := windows.GetProcAddress(localDLL, *dllProcName)
	if err != nil {
		log.Fatal(err)
	}

	mainHandle, err := winutils.CreateRemoteThread(handle, address, 0)
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		if err = windows.CloseHandle(mainHandle); err != nil {
			log.Println(err)
		}
	}()
}
