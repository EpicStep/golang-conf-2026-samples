package main

import (
	"encoding/binary"
	"log"

	"golang.org/x/sys/windows"

	"github.com/EpicStep/golang-conf-2026-samples/internal/winutils"
)

func main() {
	hwnd, err := winutils.FindWindow("GTA: San Andreas")
	if err != nil {
		log.Fatal(err)
	}

	var pid uint32
	if _, err = windows.GetWindowThreadProcessId(hwnd, &pid); err != nil {
		log.Fatal(err)
	}

	handle, err := windows.OpenProcess(
		windows.PROCESS_VM_OPERATION|windows.PROCESS_VM_WRITE,
		false, pid,
	)
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		if err = windows.CloseHandle(handle); err != nil {
			log.Println(err)
		}
	}()

	addr := uintptr(0x00438E5A) + 2 // addr from IDA + skip opcode (add) and ModRM byte (ECX)

	writeBuf := make([]byte, 4)
	newIncValue := int32(-1000000)
	binary.LittleEndian.PutUint32(writeBuf, uint32(newIncValue))

	if err = windows.WriteProcessMemory(handle, addr, &writeBuf[0], uintptr(len(writeBuf)), nil); err != nil {
		log.Fatal(err)
	}
}
