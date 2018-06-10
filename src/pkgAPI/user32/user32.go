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
var t2 string
var wins map[int]string

func init() {
	user32, _ = syscall.LoadLibrary("user32.dll")
}

func EnumWindows(parentHwnd int) (map[int]string, error) {
	wins = make(map[int]string)

	if parentHwnd > 0 {
		wins[parentHwnd] = GetWindowText(parentHwnd)
	}

	if ew, err := syscall.GetProcAddress(user32, "EnumChildWindows"); err != nil {
		return nil, err
	} else {
		defer syscall.FreeLibrary(user32)
		if _, _, err := syscall.Syscall(uintptr(ew), 3, uintptr(parentHwnd), syscall.NewCallback(enumProc), 0); err != 0 {
			abort("Call EnumChildWindows", int(err))
		}
		return wins, nil
	}
	return wins, nil
}

func enumProc(hwnd, lParam int) int {
	if _, ok := wins[hwnd]; ok {
		return 0
	}
	wins[hwnd] = GetWindowText(hwnd)

	return 1
}

func GetWindowText(hwnd int) string {
	if gw, err := syscall.GetProcAddress(user32, "GetWindowTextW"); err != nil {
		return ""
	} else {
		defer syscall.FreeLibrary(user32)
		iLen := GetWindowTextLength(hwnd) + 1
		buf := make([]uint16, iLen)

		if _, _, err := syscall.Syscall(uintptr(gw), 3, uintptr(hwnd), uintptr(unsafe.Pointer(&buf[0])), uintptr(iLen)); err != 0 {
			abort("Call GetWindowText", int(err))
		} else {
			return syscall.UTF16ToString(buf)
		}
	}
	return ""
}

func GetWindowTextLength(hwnd int) int {
	if gw, err := syscall.GetProcAddress(user32, "GetWindowTextLengthW"); err != nil {
		return 0
	} else {
		defer syscall.FreeLibrary(user32)
		if ret, _, err := syscall.Syscall(uintptr(gw), 1, uintptr(hwnd), 0, 0); err != 0 {
			abort("Call GetWindowTextLengthW", int(err))
		} else {
			return int(ret)
		}
	}
	return 0
}

func FindWindow(lpClassName, lpWindowName string) int {
	var pClassName uintptr

	if lpClassName == "" {
		pClassName = 0
	} else {
		pClassName = uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(lpClassName)))
	}

	if fw, err := syscall.GetProcAddress(user32, "FindWindowW"); err != nil {
		return -1
	} else {
		defer syscall.FreeLibrary(user32)
		if ret, _, err := syscall.Syscall(uintptr(fw), 2,
			pClassName,
			uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(lpWindowName))),
			0); err != 0 {
			return -1
			abort("Call FindWindowW", int(err))
		} else {
			if ret == 0 {
				return -1
			}
			return int(ret)
		}
	}

	return -1
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
