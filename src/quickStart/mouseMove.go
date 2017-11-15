package main

import (
	//	"fmt"
	"syscall"
	"time"

	"github.com/lxn/walk"
	"github.com/lxn/win"
)

func mouseMove() {
	lib, err := syscall.LoadLibrary("user32.dll")
	if err != nil {
		walk.MsgBox(ct.form, "user32.dll 读取出错", err.Error(), walk.MsgBoxIconInformation)
	}
	addr, err := syscall.GetProcAddress(lib, "mouse_event")
	if err != nil {
		walk.MsgBox(ct.form, "mouse_event 读取出错", err.Error(), walk.MsgBoxIconInformation)
	}

	const MOUSEEVENTF_MOVE = 0x1        //移动鼠标
	const MOUSEEVENTF_ABSOLUTE = 0x8000 //指定鼠标使用绝对坐标系，此时，屏幕在水平和垂直方向上均匀分割成65535×65535个单元
	srcWidth := win.GetSystemMetrics(win.SM_CXSCREEN)
	srcHeight := win.GetSystemMetrics(win.SM_CYSCREEN)

	for {
		p := new(win.POINT)
		win.GetCursorPos(p)
		syscall.Syscall6(addr, 5, MOUSEEVENTF_MOVE|MOUSEEVENTF_ABSOLUTE, uintptr((p.X+100)/srcWidth*65535), uintptr((p.Y+100)/srcHeight*65535), 0, 0, 0)
		//		syscall.Syscall6(addr, 5, MOUSEEVENTF_MOVE|MOUSEEVENTF_ABSOLUTE, uintptr((p.X)/srcWidth*65535), uintptr((p.Y)/srcHeight*65535), 0, 0, 0)
		win.SetCursorPos(p.X, p.Y)
		time.Sleep(1e9 * 110)
		// mouse_event(MOUSEEVENTF_ABSOLUTE Or MOUSEEVENTF_MOVE, (x + 100) / w * 65535, (y + 100) / h * 65535, 0, 0)
	}
}

func getMousePos() (x, y int) {
	p := new(win.POINT)
	win.GetCursorPos(p)

	return int(p.X), int(p.Y)
}
