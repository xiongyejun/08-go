// 调用微信截图dll

package main

import (
	"fmt"
	//	"os"
	"syscall"
)

func main() {
	//	dllPath := os.Getenv("GOPATH") + `\src\Prscrn\PrScrn.dll`
	//	fmt.Println(dllPath)
	if dll, err := syscall.LoadLibrary("PrScrn.dll"); err != nil {
		fmt.Println(err.Error() + " LoadLibrary")
	} else {
		defer syscall.FreeLibrary(dll)
		if ps, err := syscall.GetProcAddress(dll, "PrScrn"); err != nil {
			fmt.Println(err.Error() + " GetProcAddress")
		} else {

			syscall.Syscall(ps, 0, 0, 0, 0)
		}
	}

	//	prscrn := syscall.NewLazyDLL("PrScrn.dll")
	//	ps := prscrn.NewProc("PrScrn")
	//	ps.Call(uintptr(0))
	//	fmt.Println(1)
}
