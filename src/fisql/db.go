package main

import (
	"database/sql"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

type DataStruct struct {
	DBPath       string
	tableName    string
	table2Name   string
	fileSavePath string // show的时候，文件读取保存的位置

	db *sql.DB

	key     []byte         // 密码
	dicShow map[int]string // key:id	item:saveName
	pID     []int          // 记录id的slice，方便用0、1、2……的序号
}

const MAX_SIZE int = 1 * 1024 * 1024 // 文件切分为1M大小来存储
var d *DataStruct

// 打开数据库
func (me *DataStruct) getDB() (err error) {
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
			// 2018-09-02
			// 用单独1个表存储：如果数据bytes很大时，内存占用太大，运行很慢
			// 修改为：
			// files表仅存储文件名称等信息
			// filedata表来存储具体的数据，每一个MAX_SIZE大小作为1条记录
			sqlStmt := `create table files (id integer not null primary key autoincrement, name text not null, star integer not null, ext text not null);`
			if _, err = d.db.Exec(sqlStmt); err != nil {
				return
			} else {
				fmt.Println("成功创建files数据库。")
			}

			sqlStmt = `create table filedata (file_id integer references files(id) on delete cascade deferrable initially deferred, file_index integer not null, bytes blob not null, primary key (file_id,file_index));`
			if _, err = d.db.Exec(sqlStmt); err != nil {
				return
			} else {
				fmt.Println("成功创建filedata数据库。")
				return nil
			}
		}
	}
}

// 插入数据
// filesPath 文件的路径
func (me *DataStruct) insert(filesPath []string) (err error) {
	// 获取当前files表中id的最大值
	var file_id int = -1
	var stmt *sql.Stmt
	if stmt, err = d.db.Prepare("select max(id) from " + me.tableName); err != nil {
		return
	}
	defer stmt.Close()

	if err = stmt.QueryRow().Scan(&file_id); err != nil {
		if err.Error() == `sql: Scan error on column index 0: converting driver.Value type <nil> ("<nil>") to a int: invalid syntax` {

		} else if err.Error() != "sql: no rows in result set" {
			return
		}
	}
	file_id++
	fmt.Println(file_id)

	for i := range filesPath {
		if err = me.insertItem(filesPath[i], file_id); err != nil {
			fmt.Printf("%s 添加出错：%s", filesPath[i], err.Error())
		} else {
			file_id++
		}
	}

	return nil
}

func (me *DataStruct) insertItem(filePath string, file_id int) (err error) {
	var tx *sql.Tx
	if tx, err = me.db.Begin(); err != nil {
		return err
	}
	defer tx.Commit()

	var stmt *sql.Stmt
	if stmt, err = tx.Prepare("insert into " + me.tableName + "(id,name,star,ext) values(?,?,?,?)"); err != nil {
		return err
	}
	defer stmt.Close()

	var stmt2 *sql.Stmt
	if stmt2, err = tx.Prepare("insert into " + me.table2Name + "(file_id,file_index,bytes) values(?,?,?)"); err != nil {
		return err
	}
	defer stmt2.Close()

	var fi *os.File
	if fi, err = os.Open(filePath); err != nil {
		return
	}
	defer fi.Close()

	var file_index int = 0
	var n int = 0
	// 按MAX_SIZE来读取数据
	for {
		buf := make([]byte, MAX_SIZE)
		n, err = fi.Read(buf)
		if err != nil && err != io.EOF {
			return
		}
		if 0 == n {
			break
		}
		// 加密文件byte---只加密第1个，也就是file_index=0的
		if file_index == 0 {
			if buf, err = desEncrypt(buf, d.key); err != nil {
				return
			}
		}

		if _, err = stmt2.Exec(file_id, file_index, buf); err != nil {
			return
		}
		file_index++
	}

	strExt := filepath.Ext(filePath)
	// 去除文件名的后缀
	name := strings.TrimSuffix(filepath.Base(filePath), strExt)
	// 加密文件名称
	if name, err = desEncryptString(name, d.key); err != nil {
		return
	}
	if _, err = stmt.Exec(file_id, name, 0, strExt); err != nil {
		return
	}

	return nil
}

// 清理数据库
func (me *DataStruct) cls() (err error) {
	if _, err = me.db.Exec(`VACUUM`); err != nil {
		return
	}
	return nil
}

// 删除文件
func (me *DataStruct) del(pID int) (err error) {
	id := me.pID[pID]
	sqlStmt := `delete from ` + me.table2Name + ` where file_id = ` + strconv.Itoa(id)
	if _, err = me.db.Exec(sqlStmt); err != nil {
		return
	}

	sqlStmt = `delete from ` + me.tableName + ` where id = ` + strconv.Itoa(id)
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
		if stmt, err = d.db.Prepare("select name,ext from " + me.tableName + " where id = ?"); err != nil {
			return
		}
		defer stmt.Close()

		if err = stmt.QueryRow(strconv.Itoa(id)).Scan(&name, &ext); err != nil {
			return
		}
		// 文件保存路径
		name = me.fileSavePath + strconv.Itoa(id) + ext

		// 读取文件的byte
		var rows *sql.Rows
		if rows, err = d.db.Query("select file_index,bytes from " + me.table2Name + " where file_id = " + strconv.Itoa(id) + " order by file_index"); err != nil {
			return
		}
		defer rows.Close()

		// 记录文件index，以防止存在不连续的
		var file_index_tmp int
		// 保存文件
		var file_append *os.File
		if file_append, err = os.OpenFile(name, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666); err != nil {
			return
		}
		defer file_append.Close()

		for rows.Next() {
			var file_index int
			var bi interface{}
			if err = rows.Scan(&file_index, &bi); err != nil {
				return
			}

			if file_index_tmp != file_index {
				return errors.New("不连续的file_index")
			}

			if b, ok := bi.([]byte); ok {
				// 第1段是被加密的
				if file_index == 0 {
					// 解密byte
					if b, err = desDecrypt(b, d.key); err != nil {
						return
					}
				}

				// 保存文件
				if _, err = file_append.Write(b); err != nil {
					return
				}

			} else {
				return errors.New("bi.([]byte)出错")
			}
			file_index_tmp++
		}
		// 记录打开过的，退出时删除
		me.dicShow[id] = name
		if err = rows.Err(); err != nil {
			return
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
		cmd = exec.Command("open", "-a", "/Applications/暴风影音.app", path)
	}
	if err := cmd.Start(); err != nil {
		return err
	}
	return nil
}
