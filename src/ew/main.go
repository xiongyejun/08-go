// enumWindows枚举窗口

package main

import (
	"fmt"
	"os"

	"pkgAPI/user32"

	//	"github.com/lxn/win"
)

func main() {
	var hwnd int = 0
	if len(os.Args) > 1 {
		hwnd = user32.FindWindow("", os.Args[1])
		fmt.Println("hwnd", hwnd)
	}
	if hwnd == -1 {
		fmt.Printf("没有找到窗口[%s]\r\n", os.Args[1])
		return
	}

	if wins, err := user32.EnumWindows(hwnd); err != nil {
		fmt.Println(err)
	} else {
		for k, v := range wins {
			if v != "" {
				fmt.Printf("hwnd=%d,\ttext=%s\r\n", k, v)
			}
		}
	}
}
