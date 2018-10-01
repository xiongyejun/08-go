// 读取文件的byte
package main

import (
	//	"bytes"
	"flag"
	"fmt"
	"os"
	"pkgMyPkg/colorPrint"
	"strings"
	"unicode"
)

type filebyte struct {
	file  string
	pause *bool

	f func(b []byte, p_iPre *int)
}

var fb *filebyte

const N_READ = 512 // 每次读取的byte个数

func main() {
	if len(os.Args) == 1 {
		return
	}
	// 文件是否存在
	fb.file = os.Args[len(os.Args)-1] // 先输参数，最后数文件名
	_, err := os.Stat(fb.file)
	if err != nil {
		fmt.Println(err)
		return
	}

	f, err := os.Open(fb.file)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()

	iNo := 0 // 输出的序号
	var n int = N_READ

	for n == N_READ {
		fmt.Print("\r\n")
		b := make([]byte, N_READ)
		n, _ = f.Read(b)
		fb.f(b[:n], &iNo)
	}

	colorPrint.UnSetColor()
}

func printOutPause(b []byte, p_iNo *int) {
	printOut(b[:], p_iNo)

	if len(b) < N_READ {
		return
	}
	var c string
	colorPrint.SetColor(colorPrint.White, colorPrint.Red)
	fmt.Print("\r\npause ")
	fmt.Scan(&c)
	colorPrint.UnSetColor()
	fmt.Print("\r\n")

	if c == "e" || c == "q" {
		os.Exit(1)
	}
}

func printOut(b []byte, p_iNo *int) {
	colorPrint.SetColor(colorPrint.White, colorPrint.DarkMagenta)
	fmt.Printf("   index % X ------ASCII-----", []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15})
	//	fmt.Print(strings.Repeat("-", 8+16*3+16+1))

	colorPrint.SetColor(colorPrint.White, colorPrint.DarkCyan)
	fmt.Print("\r\n")

	var str []string = make([]string, 0, N_READ/16)
	n := len(b)
	for i := 0; i < n; i += 16 {
		str_tmp := fmt.Sprintf("%08X % X ", *p_iNo, b[i:i+16])
		for _, v := range b[i : i+16] {
			if unicode.IsPrint(rune(v)) {
				str_tmp += fmt.Sprintf("%c", v)
			} else {
				str_tmp += "^"
			}
		}
		*p_iNo += 16
		str = append(str, str_tmp)
	}
	//	nstr := len(str)
	strPrint := strings.Join(str, "\r\n")
	// 将最后一些多余的清空
	if n < N_READ {
		n = n % 16
		if n > 0 {
			n = 16 - n
			nbyte := 3 * n // 空byte所占的个数
			strPrint = strPrint[:len(strPrint)-16-nbyte] +
				strings.Repeat(" ", nbyte) +
				strPrint[len(strPrint)-16:(len(strPrint)-16+(16-n))] // 最后一行的字符
		}
	}
	fmt.Print(strPrint)
}

func init() {
	fb = new(filebyte)

	fb.pause = flag.Bool("p", false, "直接打印完(pause的时候输入e或q直接退出)。")

	flag.PrintDefaults()
	flag.Parse()

	if *fb.pause {
		fb.f = printOut
	} else {
		fb.f = printOutPause
	}
}
