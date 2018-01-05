package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"pkgMyPkg/colorPrint"
	"runtime"
	"sort"
	"strconv"
	"strings"
)

var currentPath string
var pack_File = new(packFile)
var unPack_File = new(unPackFile)

func init() {
	unPack_File.unPacked = make(map[int]bool, 0)
	unPack_File.unPackedFiles = make([]string, 0)
	unPack_File.files = make([]*dirInfo, 0)

	currentPath, _ = os.Getwd()
	currentPath = currentPath + Path_Separator
}

func main() {
	printCmd()

	defer unPack_File.deleteUnPackedFIle()
	defer unPack_File.saveDir()

	r := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("Enter Cmd-> ")
		rawLine, _, _ := r.ReadLine()
		line := string(rawLine)
		if line == "q" || line == "e" {
			break
		}
		tokens := strings.Split(line, " ")
		handleCommands(tokens)
	}

}

func handleCommands(tokens []string) {
	switch tokens[0] {
	case "list":
		var f func(i int, fileName string)
		if len(tokens) > 1 {
			f = func(i int, fileName string) {
				if strings.Contains(fileName, tokens[1]) {
					fmt.Print(i, ":", fileName)
					if i%2 == 1 {
						fmt.Println()
					} else {
						fmt.Print("\t")
					}
				}
			}
		} else {
			f = func(i int, fileName string) {
				fmt.Print(i, ":", fileName)
				if i%2 == 1 {
					fmt.Println()
				} else {
					fmt.Print("\t")
				}
			}
		}

		cl := colorPrint.NewColorDll()
		cl.SetColor(colorPrint.White, colorPrint.DarkMagenta)
		for i := range unPack_File.files {
			f(i, unPack_File.files[i].FileName)
		}
		cl.UnSetColor()
		fmt.Println()
		printCmd()
	case "lists":
		if len(tokens) != 2 {
			fmt.Println("输入的命令不正确\r\nlists <num> -- star大于num数来输出")
			return
		}
		if n, err := strconv.Atoi(tokens[1]); err != nil {
			fmt.Println(err)
		} else {
			cl := colorPrint.NewColorDll()
			cl.SetColor(colorPrint.White, colorPrint.DarkMagenta)

			var k int = 0
			for i := range unPack_File.files {
				if unPack_File.files[i].Star >= n {
					fmt.Print(i, unPack_File.files[i].FileName)
					k++
					if k%2 == 1 {
						fmt.Println()
					} else {
						fmt.Print("\t")
					}
				}
			}
			cl.UnSetColor()
			fmt.Println()
			printCmd()
		}
	case "star":
		if len(tokens) != 3 {
			fmt.Println("输入的命令不正确\r\n star <index><num> -- 标记star")
			return
		}

		if index, err := strconv.Atoi(tokens[1]); err != nil {
			fmt.Println(err)
		} else if num, err := strconv.Atoi(tokens[2]); err != nil {
			fmt.Println(err)
		} else {
			unPack_File.files[index].Star = num
			unPack_File.bSave = true
		}
	case "u":
		if len(tokens) != 2 {
			fmt.Println("输入的命令不正确\r\nu<index> -- 释放某个文件")
			return
		}
		if i, err := strconv.Atoi(tokens[1]); err != nil {
			fmt.Println(err)
		} else {
			// 释放文件
			if i >= len(unPack_File.files) {
				fmt.Println("没有这么多文件。")
				return
			}
			if sf, err := unPack_File.unPackFile(i); err != nil {
				fmt.Println(err)
			} else {
				// 打开文件
				openFolderFile(sf)
			}

		}
	case "pack":
		if len(tokens) != 3 {
			fmt.Println("输入的命令不正确\r\npack <dir><saveName> -- 打包文件")
			return
		}
		PackFile(tokens[1], tokens[2])
	case "ui":
		if len(tokens) != 2 {
			fmt.Println("输入的命令不正确\r\n unpackinit <packfile> -- 读取打包文件的信息")
			return
		}
		if err := unPack_File.unPackInit(tokens[1]); err != nil {
			fmt.Println(err)
		}
	case "sort":
		// 排序前先删除原来释放的文件，因为排序后下标已经变了
		unPack_File.deleteUnPackedFIle()
		unPack_File.unPackedFiles = make([]string, 0)
		unPack_File.unPacked = make(map[int]bool, 0)

		sort.Sort(unPack_File)

	default:
		fmt.Println("Unrecognized lib command:", tokens)
	}
}

// 使用cmd打开文件和文件夹
func openFolderFile(path string) error {
	// 第4个参数，是作为start的title，不加的话有空格的path是打不开的
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd.exe", "/c", "start", "", path)
	} else {
		cmd = exec.Command("open", path)
	}
	if err := cmd.Start(); err != nil {
		return err
	}
	return nil
}

func printCmd() {
	fmt.Println(`
 Enter following commands to control:
 list <filter> -- View the files
 lists <num> -- star大于num数来输出
 star <index><num> -- 标记star
 u <index> -- 释放某个文件(unpack)
 ui <packfile> -- 读取打包文件的信息(unpackinit)
 pack <dir><saveName> -- 打包文件
 sort  -- 按文件后缀名排序
 `)
}
