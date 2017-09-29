// 读取文件的byte
package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"pkgMySelf/colorPrint"
	"strings"
)

type filebyte struct {
	file  *string
	pause *bool

	f func(b []byte, p_iPre *int)
}

var fb *filebyte
var cd *colorPrint.ColorDll

func main() {
	if *fb.file == "" {
		return
	}

	f, err := os.Open(*fb.file)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()

	N := 512
	iPre := 0
	var n int = N

	for n == N {
		fmt.Print("\r\n")
		b := make([]byte, N)
		n, _ = f.Read(b)
		fb.f(b[:n], &iPre)
	}

	cd.UnSetColor()
}

func printOutPause(b []byte, p_iPre *int) {
	printOut(b[:], p_iPre)

	fmt.Print("pause ")
	var c string
	fmt.Scan(&c)
}

func printOut(b []byte, p_iPre *int) {
	cd.SetColor(colorPrint.White, colorPrint.DarkMagenta)
	fmt.Printf("   index % X\r\n", []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15})
	fmt.Print(strings.Repeat("-", 8+16*3))

	cd.SetColor(colorPrint.White, colorPrint.DarkCyan)
	fmt.Print("\r\n")

	for i := 0; i < len(b); i += 16 {
		fmt.Printf("%08X % X ", *p_iPre, b[i:i+16])
		bb := bytes.Replace(b[i:i+16], []byte{'\n'}, []byte{'^'}, -1)
		bb = bytes.Replace(bb, []byte{'\r'}, []byte{'^'}, -1)
		for _, v := range bb {
			fmt.Printf("%c", v)
		}
		//		str := string(b[i : i+16])
		//		str = strings.Replace(str, "\n", "^", -1)
		//		str = strings.Replace(str, "\r", "^", -1)
		//		fmt.Print(str)

		fmt.Print("\r\n")
		*p_iPre += 16
	}
}

func init() {
	fb = new(filebyte)

	fb.file = flag.String("f", "", "需要查看的文件名称")
	fb.pause = flag.Bool("p", false, "打印完一段就pause")

	flag.PrintDefaults()
	flag.Parse()

	if *fb.pause {
		fb.f = printOutPause
	} else {
		fb.f = printOut
	}
	cd = colorPrint.NewColorDll()
}
