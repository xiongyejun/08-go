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
			if fileName == "Thumbs.db" {
				tokens[1] = tokens[1] + `.jpg`
				b = b[24:]
			}
			// 每一个缩略图IStream的前12个字节（3个整形）不是缩略图的内容，不能用的，因此在读取的时候跳过那三个字节好了
			if err := ioutil.WriteFile(tokens[1], b, 0666); err != nil {
				fmt.Println(err)
			}
		}
	case "saveall":
		arr := cf.GetStreams()
		if fileName == "Thumbs.db" {
			os.Mkdir("Thumbs", 0666)
		}

		for i := range arr {
			if fileName == "Thumbs.db" {
				arr[i].Name = `Thumbs\` + arr[i].Name + `.jpg`
				arr[i].B = arr[i].B[24:]
			}

			if err := ioutil.WriteFile(arr[i].Name, arr[i].B, 0666); err != nil {
				fmt.Println(err)
			}
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
 saveall -- 保存所有流
 e或者q -- 退出 
 `)

	cl.UnSetColor()
}
