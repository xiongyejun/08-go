package main

import (
	"bufio"
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
	d.DBPath = `E:\Files.db`
	d.tableName = "files"
	d.fileSavePath, _ = os.Getwd()
	d.fileSavePath += string(os.PathSeparator)

	d.dicShow = make(map[int]string)
}
func main() {
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

	//	sqlStmt := `create table files (id integer primary key autoincrement, name text, bytes blob);`
	//	if _, err := d.db.Exec(sqlStmt); err != nil {
	//		fmt.Println(err)
	//	}

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

	case "list":
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
	default:
		fmt.Println("Unrecognized lib command:", tokens)
	}
}

func printCmd() {
	fmt.Println(`
 Enter following commands to control:
 add -- 添加文件
 list -- 查看文件列表
 show <id> -- 打开文件
 rn <id> <newName> -- 重命名
 e或者q -- 退出 
 `)
}
