// MsgBox

package user32

import (
	"syscall"
	"unsafe"
)

const (
	MB_OK                = 0x00000000
	MB_OKCANCEL          = 0x00000001
	MB_ABORTRETRYIGNORE  = 0x00000002
	MB_YESNOCANCEL       = 0x00000003
	MB_YESNO             = 0x00000004
	MB_RETRYCANCEL       = 0x00000005
	MB_CANCELTRYCONTINUE = 0x00000006
	MB_ICONHAND          = 0x00000010
	MB_ICONQUESTION      = 0x00000020
	MB_ICONEXCLAMATION   = 0x00000030
	MB_ICONASTERISK      = 0x00000040
	MB_USERICON          = 0x00000080
	MB_ICONWARNING       = MB_ICONEXCLAMATION
	MB_ICONERROR         = MB_ICONHAND
	MB_ICONINFORMATION   = MB_ICONASTERISK
	MB_ICONSTOP          = MB_ICONHAND
	MB_DEFBUTTON1        = 0x00000000
	MB_DEFBUTTON2        = 0x00000100
	MB_DEFBUTTON3        = 0x00000200
	MB_DEFBUTTON4        = 0x00000300

	MB_RETURN_YES = 0x00000006
	MB_RETURN_NO  = 0x00000007
)

var user32 syscall.Handle

func init() {
	user32, _ = syscall.LoadLibrary("user32.dll")
}
func MsgBox(title, msg string, style uintptr) int {
	messageBox, _ := syscall.GetProcAddress(user32, "MessageBoxW")

	defer syscall.FreeLibrary(user32)
	ret, _, err := syscall.Syscall6(uintptr(messageBox),
		4,
		0,
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(msg))),
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(title))),
		style,
		0,
		0)

	if err != 0 {
		abort("Call MessageBox", int(err))

	}
	return int(ret)
}

func abort(funcName string, err int) {
	panic(funcName + " failed:" + syscall.Errno(err).Error())
}
