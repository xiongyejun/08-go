// cd := colorPrint.NewColorDll()
// cd.SetColor(colorPrint.White, colorPrint.DarkRed)
// fmt.Println("nihao")
// …………
// cd.UnSetColor()
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

type ColorDll struct {
	kernel32 *syscall.LazyDLL
	proc     *syscall.LazyProc
}

func (this *ColorDll) SetColor(fontColor uintptr, backGroundColor uintptr) {
	this.proc.Call(uintptr(syscall.Stdout), backGroundColor*16+fontColor)
}

func (this *ColorDll) UnSetColor() {
	handle, _, _ := this.proc.Call(uintptr(syscall.Stdout), White) // 恢复白色
	closeHandle := this.kernel32.NewProc("CloseHandle")
	closeHandle.Call(handle)
}

func NewColorDll() *ColorDll {
	cd := new(ColorDll)
	cd.kernel32 = syscall.NewLazyDLL("kernel32.dll")
	cd.proc = cd.kernel32.NewProc("SetConsoleTextAttribute")
	return cd
}
