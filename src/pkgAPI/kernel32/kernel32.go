package kernel32

import (
	"syscall"
	"unsafe"
)

var kernel32 syscall.Handle

func init() {
	kernel32, _ = syscall.LoadLibrary("kernel32.dll")
}

func MoveMemory(des, src unsafe.Pointer, length uintptr) {
	moveMemory, _ := syscall.GetProcAddress(kernel32, "RtlMoveMemory")
	defer syscall.FreeLibrary(kernel32)

	syscall.Syscall(moveMemory,
		3,
		uintptr(unsafe.Pointer(des)),
		uintptr(src),
		uintptr(length))
}
