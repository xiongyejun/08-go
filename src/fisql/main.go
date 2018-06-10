package main

import (
	"bufio"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"os"
	"strconv"
	"strings"

	"pkgAPI/comdlg32"
	"pkgMyPkg/colorPrint"

	_ "github.com/mattn/go-sqlite3"
)

func init() {
	d = new(DataStruct)
	d.DBPath = `E:\VisualStudio2015\Templates\ItemTemplates\扩展性\asldjfldsajfwo`
	d.tableName = "files"
	d.fileSavePath, _ = os.Getwd()
	d.fileSavePath += string(os.PathSeparator)

	d.dicShow = make(map[int]string)
}
func main() {
	if len(os.Args) > 1 {
		d.key = []byte(os.Args[1])
	} else {
		fmt.Println("未输入密码。")
		return
	}

	if !checkKey(d.key) {
		fmt.Println("密码错误。")
		return
	}

	if err := d.getDB(); err != nil {
		fmt.Println(err)
	}
	defer d.db.Close()
	defer d.deleteShow()

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

func checkKey(key []byte) bool {
	sha := sha1.New()

	if _, err := sha.Write(key); err != nil {
		fmt.Println(err)
		return false
	}

	b := sha.Sum(nil)
	if hex.EncodeToString(b) == "583af6c52d1ce8900d767c0f73c2db88f60ce8d1" {
		return true
	} else {
		return false
	}
}

func handleCommands(tokens []string) {
	switch tokens[0] {
	case "add":
		fd := comdlg32.NewFileDialog()
		if b, err := fd.GetOpenFileNames(); !b || err != nil {
			fmt.Println(b, err)
		} else {
			if err := d.insert(fd.FilePaths); err != nil {
				fmt.Println(err)
			}

		}
	case "del":
		if len(tokens) != 2 {
			fmt.Println(`输入的命令不正确del <id> -- 删除文件`)
			return
		}
		if n, err := strconv.Atoi(tokens[1]); err != nil {
			fmt.Println(err)
		} else {
			if err := d.del(n); err != nil {
				fmt.Println(err)
			}
		}

	case "ls":
		cl := colorPrint.NewColorDll()
		cl.SetColor(colorPrint.White, colorPrint.DarkMagenta)

		if err := d.list(); err != nil {
			fmt.Println(err)
		}
		cl.UnSetColor()

	case "show":
		if len(tokens) != 2 {
			fmt.Println(`输入的命令不正确show <id> -- 打开文件`)
			return
		}
		if n, err := strconv.Atoi(tokens[1]); err != nil {
			fmt.Println(err)
		} else {
			if err := d.show(n); err != nil {
				fmt.Println(err)
			}
		}

	case "rn":
		if len(tokens) != 3 {
			fmt.Println(`输入的命令不正确 rn <id> <newName> -- 重命名`)
			return
		}
		if n, err := strconv.Atoi(tokens[1]); err != nil {
			fmt.Println(err)
		} else {
			if err := d.rn(n, tokens[2]); err != nil {
				fmt.Println(err)
			}
		}
	case "star":
		if len(tokens) != 3 {
			fmt.Println(`输入的命令不正确 star <id> <int> -- 标星`)
			return
		}
		if id, err := strconv.Atoi(tokens[1]); err != nil {
			fmt.Println(err)
		} else {
			if n, err := strconv.Atoi(tokens[2]); err != nil {
				fmt.Println(err)
			} else {
				if err := d.star(id, n); err != nil {
					fmt.Println(err)
				}
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
 add -- 添加文件
 del <id> -- 删除文件
 ls -- 查看文件列表
 show <id> -- 打开文件
 rn <id> <newName> -- 重命名
 star <id> <int> -- 标星
 e或者q -- 退出 
 `)

	cl.UnSetColor()
}
