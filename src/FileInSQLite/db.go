package main

import (
	"database/sql"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
)

type DataStruct struct {
	DBPath       string
	tableName    string
	fileSavePath string

	db *sql.DB

	dicShow map[int]string // key:id	item:saveName
}

var d *DataStruct

// 打开数据库
func (me *DataStruct) getDB() (err error) {
	if me.db, err = sql.Open("sqlite3", d.DBPath); err != nil {
		return err
	} else {
		return nil
	}
}

// 插入数据
// filesPath 文件的路径
func (me *DataStruct) insert(filesPath []string) (err error) {
	var tx *sql.Tx
	if tx, err = me.db.Begin(); err != nil {
		return err
	}
	defer tx.Commit()

	var stmt *sql.Stmt
	if stmt, err = tx.Prepare("insert into " + me.tableName + "(id,name,bytes) values(?,?,?)"); err != nil {
		return err
	}
	defer stmt.Close()

	for i := range filesPath {
		if _, err = os.Stat(filesPath[i]); err != nil {
			fmt.Println(err)
		} else {
			// 读取文件字节
			var b []byte
			if b, err = ioutil.ReadFile(filesPath[i]); err != nil {
				fmt.Println(err)
			} else {
				if _, err = stmt.Exec(nil, filepath.Base(filesPath[i]), b); err != nil {
					fmt.Println(err)
				}
			}
		}
	}

	return nil
}

// 列出所有文件
func (me *DataStruct) list() (err error) {
	var rows *sql.Rows
	if rows, err = d.db.Query("select id,name from " + me.tableName); err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var name string
		if err = rows.Scan(&id, &name); err != nil {
			return
		}
		fmt.Printf("%3d %s\r\n", id, name)
	}

	if err = rows.Err(); err != nil {
		return
	}

	return nil
}

// 读取文件bytes，保存在exe的路径下，并打开
func (me *DataStruct) show(id int) (err error) {
	var name string
	var ok bool
	// 先判断是否已经存在了
	if name, ok = me.dicShow[id]; !ok {
		var stmt *sql.Stmt
		if stmt, err = d.db.Prepare("select name,bytes from " + me.tableName + " where id = ?"); err != nil {
			return
		}
		defer stmt.Close()

		var b interface{}

		if err = stmt.QueryRow(strconv.Itoa(id)).Scan(&name, &b); err != nil {
			return errors.New("stmt.QueryRow\r\n") //+ err.Error())
		}

		name = me.fileSavePath + strconv.Itoa(id) + filepath.Ext(name)
		fmt.Println(name)
		if err = ioutil.WriteFile(name, b.([]byte), 0666); err != nil {
			return
		}
		me.dicShow[id] = name
	} /* else {
		fmt.Println("已经有了")
	}*/

	//	fmt.Println(name)
	return openFolderFile(name)
}

// 删除已经释放的文件
func (me *DataStruct) deleteShow() {
	for _, item := range me.dicShow {
		os.Remove(item)
	}
}

// 使用cmd打开文件和文件夹
func openFolderFile(path string) error {
	fmt.Println("open:", path)
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
