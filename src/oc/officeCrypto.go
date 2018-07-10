// officeCrypto office文件找回打开密码

package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"pkgMyPkg/compoundFile"
	"pkgMyPkg/keysInSQLite"
	"pkgMyPkg/offCrypto"
	"pkgMyPkg/permuCombin"
	"runtime"
	"strconv"
	"time"
)

var COUNT_CHECK int = runtime.NumCPU() // 启动多少个进程测试
var arrCount []uint64                  // 统计已经测试的个数

var ies []offCrypto.IEncryptedType // 密码测试的接口
var chanPsw chan []byte            // 用通道里接受读取的密码
var bflag bool = false
var totalsum uint64 = 0

// 命令参数，选择用什么来获取密码
var bKey *bool    // 通过txt文件
var bPermu *bool  // 排列组合
var bSQLite *bool // 读取数据库

// 密码信息存放的文件名
const EI_OFFICE_07 string = "EncryptionInfo"
const EI_XLS_03 string = "Workbook"
const EI_DOC_03 string = "1Table"
const EI_PPT_03 string = "PowerPoint Document"

func main() {
	fmt.Println("1 oc <-s> <FileName> <sqlWhere txt>") // sqlWhere 只需要在第1行设置条件
	fmt.Println("2 oc <-k> <FileName> <key txt>")
	fmt.Println("3 oc <-p> <FileName> <srcPermu txt> <selectCount>")

	if *bKey && *bPermu && *bSQLite {
		fmt.Println("err: -k、-p、-s只能设置一个。")
		return
	} else if !*bKey && !*bPermu && !*bSQLite {
		fmt.Println("err: -k、-p、-s至少设置一个。")
		return
	}

	if *bKey {
		if len(os.Args) != 4 {
			return
		}
	} else if *bPermu {
		if len(os.Args) != 5 {
			return
		}
	} else if *bSQLite {
		if len(os.Args) != 4 {
			return
		}
	}

	// 读取检查密码所需数据
	var bEncryptionInfo []byte
	var err error
	var bECMA376 bool = true
	if bEncryptionInfo, bECMA376, err = readEncryptionInfo(os.Args[2]); err != nil {
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

	// 获取测试的密码
	chanPsw = make(chan []byte, COUNT_CHECK*100)
	if *bKey {
		// 读取记录密码的txt文件
		if f, err := os.Open(os.Args[3]); err != nil {
			fmt.Println(err)
			return
		} else {
			defer f.Close()
			go readPassword(chanPsw, f)
		}
	} else if *bPermu {
		//  使用排列组合方式
		var src []string
		if f, err := os.Open(os.Args[3]); err != nil {
			fmt.Println(err)
			return
		} else {
			defer f.Close()
			src = readPermuSrc(f)
		}
		if selectCount, err := strconv.Atoi(os.Args[4]); err != nil {
			fmt.Println(err)
			return
		} else {
			fmt.Println("src =", src)
			go permuCombin.PermuString(src, uint(selectCount), chanPsw)
			totalsum = uint64(getCount(len(src), int(selectCount)))
			fmt.Printf("共有排列组合%d个。\r\n", totalsum)
		}
	} else if *bSQLite {
		// 读取SQLite数据库
		if err := keysInSQLite.GetDB(); err != nil {
			fmt.Println(err)
			return
		}
		defer keysInSQLite.CloseDB()

		go keysInSQLite.SelectValue(os.Args[3], chanPsw, &totalsum)
	}

	// 开始测试密码
	arrCount = make([]uint64, COUNT_CHECK)
	for i := 0; i < COUNT_CHECK; i++ {
		go checkPassword(ies[i], chanPsw, &arrCount[i])
	}

	var iTime int64 = 0
	for !bflag {
		time.Sleep(1 * 1e9)
		iTime++
		var sum uint64 = getSumCount()
		fmt.Printf("\r正在测试第%10s个密码，用时：%s……", formatUint64(sum), formatTime(iTime))
		if sum == totalsum {
			fmt.Println("\r\n未能找到密码！")
			return
		}
	}

	fmt.Println("\r\n未能找到密码！")
}

func init() {
	bSQLite = flag.Bool("s", false, "从SQLite读取")
	bKey = flag.Bool("k", false, "读取密码字典")
	bPermu = flag.Bool("p", false, "排列组合")

	flag.PrintDefaults()
	flag.Parse()
}

// ie		IEncryptedType密码测试的接口
// chanPsw	记录密码的channel
// count	统计测试的密码的个数
func checkPassword(ie offCrypto.IEncryptedType, chanPsw chan []byte, count *uint64) {
	var pswUnicode []byte

	for {
		pswUnicode = <-chanPsw
		//		tmp := []byte("1")
		//		pswUnicode = append(tmp, pswUnicode...)
		pswUnicode = asc2Unicode(pswUnicode)
		if err := ie.CheckPassword(pswUnicode); err == nil {
			if b, err := unicode2Asc(pswUnicode); err != nil {
				fmt.Println(err)
			} else {
				fmt.Printf("\r\n找到密码%s密码长度=%d\r\n", b, len(b))
				bflag = true
			}
			os.Exit(0)
		}
		*count++
	}
}

// 测试了多少个，将多个goroutine统计的数据加起来
func getSumCount() uint64 {
	// 统计测试的总个数
	var sum uint64 = 0
	for i := 0; i < COUNT_CHECK; i++ {
		sum += arrCount[i]
	}
	return sum
}

// 读取解密所需的数据
func readEncryptionInfo(fileName string) (bEncryptionInfo []byte, bECMA376 bool, err error) {
	if bEncryptionInfo, err = ioutil.ReadFile(fileName); err != nil {
		return
	} else {
		// 可以直接读取加密的文件EncryptionInfo和Workbook……
		if fileName == EI_OFFICE_07 {
			bECMA376 = true
		} else if fileName == EI_XLS_03 {
			bECMA376 = false
			bEncryptionInfo = bEncryptionInfo[0x1A:]
		} else if fileName == EI_DOC_03 {
			bECMA376 = false
		} else if fileName == EI_PPT_03 {
			bECMA376 = false
			// ppt加密信息放在文件的后面
			if bEncryptionInfo, err = getPPTEncryptionInfo(bEncryptionInfo); err != nil {
				return
			}
		} else {
			// 否则解析复合文档
			if cf, err1 := compoundFile.NewCompoundFile(bEncryptionInfo); err1 != nil {
				err = err1
				return
			} else {
				if err = cf.Parse(); err != nil {
					fmt.Println(err)
					return
				}
				// 先尝试读取ECMA376的
				if bEncryptionInfo, err = cf.GetStream(EI_OFFICE_07); err != nil {
					// 再尝试offBinary
					if bEncryptionInfo, err = cf.GetStream(EI_XLS_03); err != nil {
						// word的加密信息就是0开始
						if bEncryptionInfo, err = cf.GetStream(EI_DOC_03); err != nil {
							// ppt的加密信息放在最后面
							if bEncryptionInfo, err = cf.GetStream(EI_PPT_03); err != nil {
								err = errors.New("没有找到EncryptionInfo或者Workbook或1Table或PowerPoint Document流。")
								return
							} else {
								if bEncryptionInfo, err = getPPTEncryptionInfo(bEncryptionInfo); err != nil {
									return
								}
								bECMA376 = false
							}
						}
					} else {
						// Workbook Stream的加密信息是从0x1A开始的，这个是通过查看字节信息猜的！
						bEncryptionInfo = bEncryptionInfo[0x1A:]
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

// 从f中逐行读取密码，放入chanPsw
func readPassword(chanPsw chan []byte, f *os.File) {
	bf := bufio.NewReader(f)
	for {
		b, _, err := bf.ReadLine()
		if err == io.EOF {
			break
		}
		chanPsw <- b
		// 统计读取到了多少个密码
		totalsum++
	}
}

// 读取排列组合用的数据源，逐行读取
func readPermuSrc(f *os.File) (src []string) {
	src = make([]string, 0)
	bf := bufio.NewReader(f)
	for {
		b, _, err := bf.ReadLine()
		if err == io.EOF {
			break
		}
		src = append(src, string(b))
	}
	return
}
