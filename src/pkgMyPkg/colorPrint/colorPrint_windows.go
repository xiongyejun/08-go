// 打印彩色
package colorPrint

import (
	"syscall"
)

const (
	Black uintptr = iota
	DarkBlue
	DarkGreen
	DarkCyan
	DarkRed
	DarkMagenta
	DarkYellow
	Gray
	DarkGray
	Blue
	Green
	Cyan
	Red
	Magenta
	Yellow
	White
	// 背景色 background 是按16循环的
	// 0 是Black背景，Black字体
	// 16就是DarkBlue背景，Black字体
	// 32就是DarkGreen背景，Black字体
	// 94就是DarkMagenta背景，Yellow字体
)

var kernel32 *syscall.LazyDLL
var proc *syscall.LazyProc

func SetColor(fontColor uintptr, backGroundColor uintptr) {
	proc.Call(uintptr(syscall.Stdout), backGroundColor*16+fontColor)
}

func UnSetColor() {
	handle, _, _ := proc.Call(uintptr(syscall.Stdout), White) // 恢复白色
	closeHandle := kernel32.NewProc("CloseHandle")
	closeHandle.Call(handle)
}

func init() {
	kernel32 = syscall.NewLazyDLL("kernel32.dll")
	proc = kernel32.NewProc("SetConsoleTextAttribute")
}
