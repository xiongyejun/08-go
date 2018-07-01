package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"pkgMyPkg/compoundFile"
	"pkgMyPkg/offCrypto"
	"runtime"
	"time"
)

var COUNT_CHECK int = runtime.NumCPU() // 启动多少个测试
var arrCount []int                     // 统计已经测试的个数

var ies []offCrypto.IEncryptedType
var chanPsw chan []byte
var flag bool = false

func main() {
	if len(os.Args) != 3 {
		fmt.Println("CompoundFileEncryption <FileName> <Password Txt>")
		return
	}

	var bEncryptionInfo []byte
	var err1 error
	var bECMA376 bool = true

	if b, err := ioutil.ReadFile(os.Args[1]); err != nil {
		fmt.Println(err)
		return
	} else {
		// 可以直接读取加密的文件EncryptionInfo和Workbook
		if os.Args[1] == "EncryptionInfo" {
			bEncryptionInfo = b
		} else if os.Args[1] == "Workbook" {
			bEncryptionInfo = b
			bECMA376 = false
		} else {
			// 否则解析复合文档
			if cf, err := compoundFile.NewCompoundFile(b); err != nil {
				fmt.Println(err)
				return
			} else {
				cf.Parse()
				// 先尝试读取ECMA376的
				if bEncryptionInfo, err1 = cf.GetStream(`EncryptionInfo`); err1 != nil {
					// 再尝试offBinary
					if bEncryptionInfo, err1 = cf.GetStream("Workbook"); err1 != nil {
						fmt.Println("没有找到EncryptionInfo或者Workbook流。")
						return
					} else {
						bECMA376 = false
					}
				}
				cf = nil
			}
		}

	}

	var err error
	ies = make([]offCrypto.IEncryptedType, COUNT_CHECK)
	// 初始化offCrypto.IEncryptedType
	for i := 0; i < COUNT_CHECK; i++ {
		if ies[i], err = offCrypto.NewIEncrypted(bEncryptionInfo, bECMA376); err != nil {
			fmt.Println(err)
			return
		}
	}
	// 读取记录密码的txt文件
	chanPsw = make(chan []byte, 10)
	if f, err := os.Open(os.Args[2]); err != nil {
		fmt.Println(err)
		return
	} else {
		defer f.Close()
		go readPassword(chanPsw, f)
	}

	// 开始测试密码
	arrCount = make([]int, COUNT_CHECK)
	for i := 0; i < COUNT_CHECK; i++ {
		go checkPassword(ies[i], chanPsw, &arrCount[i])
	}

	timeStart := time.Now()
	for !flag {
		time.Sleep(3 * 1e9)

		var sum int = 0
		for i := 0; i < COUNT_CHECK; i++ {
			sum += arrCount[i]
		}
		fmt.Printf("正在测试第%7d个密码，用时：%s……\r\n", sum, formatTime(time.Now().Unix()-timeStart.Unix()))
	}
	fmt.Println("未能找到密码！")
}

func formatTime(second int64) string {
	m := second / 60
	second = second % 60

	h := m / 60
	m = m % 60

	return fmt.Sprintf("%2d时%2d分%2d秒", h, m, second)
}

func checkPassword(ie offCrypto.IEncryptedType, chanPsw chan []byte, count *int) {
	var pswUnicode []byte

	for {
		*count++
		pswUnicode = <-chanPsw
		if err := ie.CheckPassword(pswUnicode); err == nil {
			if b, err := unicode2Byte(pswUnicode); err != nil {
				fmt.Println(err)
			} else {
				fmt.Printf("找到密码：%s\r\n", b)
				flag = true
			}
			os.Exit(0)
		}
	}

}

func readPassword(chanPsw chan []byte, f *os.File) {
	bf := bufio.NewReader(f)
	for {
		b, _, err := bf.ReadLine()
		if err == io.EOF {
			break
		}
		b = byte2Unicode(b)
		chanPsw <- b
	}
	// 读完了密码字典txt
	flag = true
}

func byte2Unicode(b []byte) []byte {
	var bb []byte = make([]byte, len(b)*2)
	for i := range b {
		bb[i*2] = b[i]
		bb[i*2+1] = 0
	}
	return bb
}
func unicode2Byte(b []byte) (bb []byte, err error) {
	var iLen int = len(b)
	if iLen%2 != 0 {
		return nil, errors.New("输入的unicode byte长度必须是整数。")
	}

	bb = make([]byte, iLen)
	for i := 0; i < iLen; i += 2 {
		bb[i/2] = b[i]
	}
	return bb, nil
}

func string2Unicode(str string) []byte {
	b := []byte(str)
	var bb []byte = make([]byte, len(b)*2)
	for i := range b {
		bb[i*2] = b[i]
		bb[i*2+1] = 0
	}
	return bb
}
