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
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

// 打开数据库
func (me *dataStruct) getDB() (err error) {
	if _, err = os.Stat(d.DBPath); err == nil {
		if me.db, err = sql.Open("sqlite3", d.DBPath); err != nil {
			return
		} else {
			fmt.Println("成功打开数据库。")
			return nil
		}
	} else {
		// 不存在数据库的情况下进行创建
		if me.db, err = sql.Open("sqlite3", d.DBPath); err != nil {
			return
		} else {
			// 														{Competitions:"J联赛", Datetime:"09-05-16 11:55:00", Dispatchamt:"70.29", Guestteam:"大宫", Handicap:"0", Hitcount:"9", Hostteam:"名古屋", Matchno:"1", Result:"2", Score:"1:1", Spvalue:"6.011832", Stake:"7.8153"}
			sqlStmt := `create table datas (drawno integer not null, Competitions text, matchno integer not null, primary key (drawno,matchno));`
			if _, err = d.db.Exec(sqlStmt); err != nil {
				return
			} else {
				fmt.Println("成功创建数据库。")
				return nil
			}
		}
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
	if stmt, err = tx.Prepare("insert into " + me.tableName + "(id,name,star,ext,bytes) values(?,?,?,?,?)"); err != nil {
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
				return err
			} else {
				strExt := filepath.Ext(filesPath[i])
				// 去除文件名的后缀
				name := strings.TrimSuffix(filepath.Base(filesPath[i]), strExt)
				// 加密文件byte
				if b, err = desEncrypt(b, d.key); err != nil {
					return
				}
				// 加密文件名称
				if name, err = desEncryptString(name, d.key); err != nil {
					return
				}
				if _, err = stmt.Exec(nil, name, 0, strExt, b); err != nil {
					fmt.Println(err)
				}
			}
		}
	}

	return nil
}

// 删除文件
func (me *DataStruct) del(pID int) (err error) {
	id := me.pID[pID]
	sqlStmt := `delete from ` + me.tableName + ` where id = ` + strconv.Itoa(id)
	if _, err = me.db.Exec(sqlStmt); err != nil {
		return
	}
	return nil
}

// 重命名
func (me *DataStruct) rn(pID int, newName string) (err error) {
	id := me.pID[pID]
	if newName, err = desEncryptString(newName, d.key); err != nil {
		return
	}
	sqlStmt := `update ` + me.tableName + ` set name="` + newName + `" where id = ` + strconv.Itoa(id)
	if _, err = me.db.Exec(sqlStmt); err != nil {
		return
	}
	return nil
}

// 标星
func (me *DataStruct) star(pID int, iStar int) (err error) {
	id := me.pID[pID]
	sqlStmt := `update ` + me.tableName + ` set star=` + strconv.Itoa(iStar) + ` where id = ` + strconv.Itoa(id)
	if _, err = me.db.Exec(sqlStmt); err != nil {
		return
	}
	return nil
}

// 列出所有文件
func (me *DataStruct) list() (err error) {
	var rows *sql.Rows
	if rows, err = d.db.Query("select id,star,name,ext from " + me.tableName); err != nil {
		return
	}
	defer rows.Close()

	me.pID = make([]int, 0)
	var pIDCount int = 0
	for rows.Next() {
		var id int
		var star int
		var name string
		var ext string
		if err = rows.Scan(&id, &star, &name, &ext); err != nil {
			return
		}
		// 解密文件名
		if name, err = desDecryptString(name, d.key); err != nil {
			return
		}

		me.pID = append(me.pID, id)
		fmt.Printf("%3d\t%3d\t%s\r\n", pIDCount, star, name+ext)
		pIDCount++
	}

	if err = rows.Err(); err != nil {
		return
	}

	return nil
}

// 读取文件bytes，保存在当前程序的路径下，并打开
func (me *DataStruct) show(pID int) (err error) {
	var name string
	var ext string
	var ok bool

	if pID >= len(me.pID) {
		return errors.New("不存在的索引。")
	}
	id := me.pID[pID]
	// 先判断是否已经存在了
	if name, ok = me.dicShow[id]; !ok {
		var stmt *sql.Stmt
		if stmt, err = d.db.Prepare("select name,ext,bytes from " + me.tableName + " where id = ?"); err != nil {
			return
		}
		defer stmt.Close()

		var bi interface{}
		if err = stmt.QueryRow(strconv.Itoa(id)).Scan(&name, &ext, &bi); err != nil {
			return
		}
		// 文件保存路径
		name = me.fileSavePath + strconv.Itoa(id) + ext
		// 读取文件的byte
		if b, ok := bi.([]byte); ok {
			// 解密byte
			if b, err = desDecrypt(b, d.key); err != nil {
				return
			}
			// 保存文件
			if err = ioutil.WriteFile(name, b, 0666); err != nil {
				return
			}
			// 记录打开过的，退出时删除
			me.dicShow[id] = name

		}

	} /* else {
		fmt.Println("已经有了")
	}*/

	//	fmt.Println(name)
	return openFolderFile(name)
}

// 删除已经释放的文件
func (me *DataStruct) deleteShow() {
	for _, item := range me.dicShow {
		if err := os.Remove(item); err != nil {
			fmt.Println(err)
		}
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