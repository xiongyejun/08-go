package main

import (
	"bufio"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"os"
	"strconv"
	"strings"

	//	"pkgAPI/comdlg32"
	"pkgMyPkg/colorPrint"

	_ "github.com/mattn/go-sqlite3"
)

func init() {
	d = new(DataStruct)
	d.DBPath = `/Users/xiongyejun/Documents/Visual Studio 2010/Backup Files/ConsoleApplication1/asldjfldsajfwo2`
	d.tableName = "files"
	d.table2Name = "filedata"
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
		if len(tokens) != 2 {
			fmt.Println("add <filePath>")
			return
		}
		files := make([]string, 0)
		files = append(files, tokens[1])
		if err := d.insert(files); err != nil {
			fmt.Println(err)
		}

		// 清理数据库
	case "clear":
		if err := d.cls(); err != nil {
			fmt.Println(err)
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
		colorPrint.SetColor(colorPrint.White, colorPrint.DarkMagenta)

		if err := d.list(); err != nil {
			fmt.Println(err)
		}
		colorPrint.UnSetColor()

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
	colorPrint.SetColor(colorPrint.Green, colorPrint.Black)

	fmt.Println(" Enter following commands to control:\n" +
		" add <filePath>-- 添加文件          del <id> -- 删除文件\n" +
		" ls -- 查看文件列表                        show <id> -- 打开文件\n" +
		" rn <id> <newName> -- 重命名   star <id> <int> -- 标星\n" +
		" clear -- 清理数据库                       e或者q -- 退出")

	colorPrint.UnSetColor()
}
