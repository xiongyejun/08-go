// 将用十六进制表示的文件转成二进制文件
package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"pkgAPI/comdlg32"
	"strconv"
)

func main() {
	var hexFile string

	if len(os.Args) == 1 {
		fd := comdlg32.NewFileDialog()
		b, _ := fd.GetOpenFileName()
		if !b {
			return
		}
		hexFile = fd.FilePath
	} else {
		hexFile = os.Args[1]
	}

	f, err := os.Open(hexFile)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()
	bHex, err := ioutil.ReadAll(f)

	if bHex[0] == 0xef && bHex[1] == 0xbb && bHex[2] == 0xbf {
		bHex = bHex[3:] // 跳过UTF-8的头
	}

	if len(bHex)%2 > 0 {
		fmt.Println("字节个数是奇数。")
		return
	}
	resultB := make([]byte, len(bHex)/2)
	for i := 0; i < len(bHex); i += 2 {
		n, err := strconv.ParseInt(string(bHex[i:i+2]), 16, 32)
		if err != nil {
			fmt.Println(err)
			return
		}
		resultB[i/2] = byte(n)
	}

	// fs 保存文件

	fs, err := os.OpenFile(hexFile+getExt(resultB), os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer fs.Close()
	fs.Write(resultB)
	fmt.Println("ok")
}

// 获取扩展名
func getExt(b []byte) string {
	if b[0] == '7' && b[1] == 'z' {
		return ".7z"
	} else if b[0] == 'P' && b[1] == 'K' {
		return ".zip"
	} else {
		return ".bin"
	}
}
