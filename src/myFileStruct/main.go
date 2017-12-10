package main

import (
	"bufio"
	"fmt"
	"os"
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
	fmt.Println(`
 Enter following commands to control:
 list -- View the files
 unpack <index> -- 释放某个文件
 unpackinit <packfile> -- 读取打包文件的信息
 pack <dir><saveName> -- 打包文件
 `)

	defer unPack_File.deleteUnPackedFIle()

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
		for i := range unPack_File.files {
			fmt.Println(i, ":", unPack_File.files[i].FileName)
		}
	case "unpack":
		if len(tokens) != 2 {
			fmt.Println("输入的命令不正确\r\nunpack <index> -- 释放某个文件")
			return
		}
		if i, err := strconv.Atoi(tokens[1]); err != nil {
			fmt.Println(err)
		} else {
			// 释放文件
			unPack_File.unPackFile(i)
			// 打开文件

		}
	case "pack":
		if len(tokens) != 3 {
			fmt.Println("输入的命令不正确\r\npack <dir><saveName> -- 打包文件")
			return
		}
		PackFile(tokens[1], tokens[2])
	case "unpackinit":
		if len(tokens) != 2 {
			fmt.Println("输入的命令不正确\r\n unpackinit <packfile> -- 读取打包文件的信息")
			return
		}
		if err := unPack_File.unPackInit(tokens[1]); err != nil {
			fmt.Println(err)
		}
	default:
		fmt.Println("Unrecognized lib command:", tokens)
	}
}
