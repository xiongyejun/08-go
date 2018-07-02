// officeCrypto office文件找回打开密码

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
	"strconv"
	"time"
)

var COUNT_CHECK int = runtime.NumCPU() // 启动多少个进程测试
var arrCount []uint64                  // 统计已经测试的个数

var ies []offCrypto.IEncryptedType
var chanPsw chan []byte
var flag bool = false

func main() {
	fmt.Println("officeCrypto <FileName> <Password txt>")
	if len(os.Args) != 3 {
		return
	}

	// 读取检查密码所需数据
	var bEncryptionInfo []byte
	var err error
	var bECMA376 bool = true
	if bEncryptionInfo, bECMA376, err = readEncryptionInfo(os.Args[1]); err != nil {
		fmt.Println(err)
		return
	}

	// 初始化offCrypto.IEncryptedType
	ies = make([]offCrypto.IEncryptedType, COUNT_CHECK)
	for i := 0; i < COUNT_CHECK; i++ {
		if ies[i], err = offCrypto.NewIEncrypted(bEncryptionInfo, bECMA376); err != nil {
			fmt.Println(err)
			return
		}
	}
	// 读取记录密码的txt文件
	chanPsw = make(chan []byte, COUNT_CHECK*10)
	if f, err := os.Open(os.Args[2]); err != nil {
		fmt.Println(err)
		return
	} else {
		defer f.Close()
		go readPassword(chanPsw, f)
	}

	// 开始测试密码
	arrCount = make([]uint64, COUNT_CHECK)
	for i := 0; i < COUNT_CHECK; i++ {
		go checkPassword(ies[i], chanPsw, &arrCount[i])
	}

	timeStart := time.Now()
	for !flag {
		time.Sleep(0.01 * 1e9)
		printSumCount(timeStart)
	}

	printSumCount(timeStart)
	fmt.Println("\r\n未能找到密码！")
}

// 输出测试了多少个
func printSumCount(timeStart time.Time) {
	// 统计测试的总个数
	var sum uint64 = 0
	for i := 0; i < COUNT_CHECK; i++ {
		sum += arrCount[i]
	}

	fmt.Printf("\r正在测试第%10s个密码，用时：%s……", formatUint64(sum), formatTime(time.Now().Unix()-timeStart.Unix()))
}

// 读取解密所需的数据
func readEncryptionInfo(fileName string) (bEncryptionInfo []byte, bECMA376 bool, err error) {
	if bEncryptionInfo, err = ioutil.ReadFile(fileName); err != nil {
		return
	} else {
		// 可以直接读取加密的文件EncryptionInfo和Workbook
		if fileName == "EncryptionInfo" {
			bECMA376 = true
		} else if fileName == "Workbook" {
			bECMA376 = false
		} else {
			// 否则解析复合文档
			if cf, err1 := compoundFile.NewCompoundFile(bEncryptionInfo); err1 != nil {
				err = err1
				return
			} else {
				cf.Parse()
				// 先尝试读取ECMA376的
				if bEncryptionInfo, err = cf.GetStream(`EncryptionInfo`); err != nil {
					// 再尝试offBinary
					if bEncryptionInfo, err = cf.GetStream("Workbook"); err != nil {
						err = errors.New("没有找到EncryptionInfo或者Workbook流。")
						return
					} else {
						bECMA376 = false
					}
				} else {
					bECMA376 = true
				}
				cf = nil
			}
		}
	}
	return
}

func formatUint64(i uint64) (str string) {
	str = strconv.Itoa(int(i % 10000))
	i /= 10000
	if i > 0 {
		str = "万" + str
	} else {
		return
	}

	str = strconv.Itoa(int(i%10000)) + str
	i /= 10000
	if i > 0 {
		str = "亿" + str
	} else {
		return
	}
	str = strconv.Itoa(int(i)) + str

	return
}

func formatTime(second int64) string {
	m := second / 60
	second = second % 60

	h := m / 60
	m = m % 60

	return fmt.Sprintf("%2d时%2d分%2d秒", h, m, second)
}

// ie		IEncryptedType密码测试的接口
// chanPsw	记录密码的channel
// count	统计测试的密码的个数
func checkPassword(ie offCrypto.IEncryptedType, chanPsw chan []byte, count *uint64) {
	var pswUnicode []byte

	for {
		*count++
		pswUnicode = <-chanPsw
		if err := ie.CheckPassword(pswUnicode); err == nil {
			if b, err := unicode2Asc(pswUnicode); err != nil {
				fmt.Println(err)
			} else {
				fmt.Printf("\r\n找到密码：%s\r\n", b)
				flag = true
			}
			os.Exit(0)
		}

	}
}

// 从f中读取密码，放入chanPsw
func readPassword(chanPsw chan []byte, f *os.File) {
	bf := bufio.NewReader(f)
	for {
		b, _, err := bf.ReadLine()
		if err == io.EOF {
			break
		}
		b = asc2Unicode(b)
		chanPsw <- b
	}
	// 读完了密码字典txt
	flag = true
}

// asc的byte转换为unicode的byte
func asc2Unicode(b []byte) []byte {
	var bb []byte = make([]byte, len(b)*2)
	for i := range b {
		bb[i*2] = b[i]
		bb[i*2+1] = 0
	}
	return bb
}

func unicode2Asc(b []byte) (bb []byte, err error) {
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
