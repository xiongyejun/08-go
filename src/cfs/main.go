// compoundFileStruct
package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"pkgMyPkg/colorPrint"
	"pkgMyPkg/compoundFile"
	"strings"
)

var cf *compoundFile.CompoundFile
var fileName string

func main() {
	if len(os.Args) == 1 {
		fmt.Println("请输入文件名。")
	} else {
		fileName = os.Args[1]
		if b, err := ioutil.ReadFile(fileName); err != nil {
			fmt.Println(err)
		} else {
			var err1 error
			if cf, err1 = compoundFile.NewCompoundFile(b); err1 != nil {
				fmt.Println(err1)
			} else {
				if err := cf.Parse(); err != nil {
					fmt.Println(err)
				} else {
					r := bufio.NewReader(os.Stdin)
					for {
						printCmd()
						fmt.Print("Enter Cmd->")
						rawLine, _, _ := r.ReadLine()
						line := string(rawLine)
						if line == "q" || line == "e" {
							break
						}
						tokens := strings.Split(line, " ")
						handleCommands(tokens)
					}
				}
			}
		}
	}
}
func handleCommands(tokens []string) {
	switch tokens[0] {
	case "ls":
		cl := colorPrint.NewColorDll()
		cf.PrintOut()
		cl.UnSetColor()

	case "show":
		if len(tokens) != 2 {
			fmt.Println(`输入的命令不正确show <name> -- 输出文件数据`)
			fmt.Printf("%d, %#v\r\n", len(tokens), tokens)
			return
		}
		if b, err := cf.GetStream(tokens[1]); err != nil {
			fmt.Println(err)
		} else {
			var saveFile string
			arr := strings.Split(tokens[1], `\`)
			saveFile = arr[len(arr)-1]

			if fileName == "Thumbs.db" {
				saveFile += `.jpg`
				b = b[24:]
			}
			// 每一个缩略图IStream的前12个字节（3个整形）不是缩略图的内容
			// 不能用的，因此在读取的时候跳过那三个字节好了
			// 64位系统是24位？ 不确定
			if err := ioutil.WriteFile(saveFile, b, 0666); err != nil {
				fmt.Println(err)
			}
		}
	case "rel":
		var bOffset int = 0
		var strExt string = ""

		if fileName == "Thumbs.db" {
			bOffset = 24
			strExt = ".jpg"
		}
		if err := cf.Release(bOffset, strExt); err != nil {
			fmt.Println(err)
			return
		}

	default:
		fmt.Println("Unrecognized lib command:", tokens)
	}
}
func printCmd() {
	cl := colorPrint.NewColorDll()
	cl.SetColor(colorPrint.Green, colorPrint.Black)

	fmt.Println(`
 Enter following commands to control:
 ls -- 查看文件列表
 show <name> -- 输出文件数据
 rel -- release释放所有流
 e或者q -- 退出 
 `)

	cl.UnSetColor()
}
