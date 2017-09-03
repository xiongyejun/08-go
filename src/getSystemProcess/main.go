package main

import (
	"fmt"
	"strings"
	//	"strconv"
	"syscall"
	"unsafe"
)

type ulong int32
type ulong_ptr uintptr

type PROCESSENTRY32 struct {
	dwSize              ulong
	cntUsage            ulong
	th32ProcessID       ulong
	th32DefaultHeapID   ulong_ptr
	th32ModuleID        ulong
	cntThreads          ulong
	th32ParentProcessID ulong
	pcPriClassBase      ulong
	dwFlags             ulong
	szExeFile           [260]byte
}

func main() {
	processName, processID := getProcess()
	for i, v := range processName {
		fmt.Printf("%d\t%30s\t%d\n", i, v, processID[i])
	}
	//	fmt.Println(processID)
	//	fmt.Println(processName[1])
	//	for i, v := range processName[1] {
	//		fmt.Printf("i=%d v=%c v=%d\n", i, v, v)
	//	}
}

func getProcess() ([]string, []int) {
	const MAX_NUM int = 500
	var processName [MAX_NUM]string
	var processID [MAX_NUM]int
	var pNum int = 0

	kernel32 := syscall.NewLazyDLL("kernel32.dll")
	CreateToolhelp32Snapshot := kernel32.NewProc("CreateToolhelp32Snapshot")
	pHandle, _, _ := CreateToolhelp32Snapshot.Call(uintptr(0x2), uintptr(0x0))
	if int(pHandle) == -1 {
		return nil, nil
	}
	Process32Next := kernel32.NewProc("Process32Next")
	for {
		var proc PROCESSENTRY32
		proc.dwSize = ulong(unsafe.Sizeof(proc))
		if rt, _, _ := Process32Next.Call(uintptr(pHandle), uintptr(unsafe.Pointer(&proc))); int(rt) == 1 {
			processName[pNum] = strings.TrimRight(string(proc.szExeFile[0:]), string(byte(0)))
			processID[pNum] = int(proc.th32ProcessID)
			pNum++
		} else {
			break
		}
	}
	CloseHandle := kernel32.NewProc("CloseHandle")
	_, _, _ = CloseHandle.Call(pHandle)

	return processName[0:pNum], processID[0:pNum]
}
