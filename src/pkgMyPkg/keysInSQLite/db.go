package keysInSQLite

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

type dbField struct {
	password     string
	length       int
	typeI        int
	successCount int
	chineseHabit int
}

type DataStruct struct {
	DBPath    string
	tableName string

	db *sql.DB
}

var d *DataStruct

func CloseDB() {
	d.db.Close()
}
func GetDB() (err error) {
	return d.getDB()
}

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
			// 创建表									密码									   密码长度				密码类型					密码曾经成功几次					中国人习惯用的密码用 1表示
			sqlStmt := `create table ` + d.tableName + `(password TEXT not null primary key, len integer not null, type integer not null, successCount integer not null, chineseHabit integer not null);`
			if _, err = d.db.Exec(sqlStmt); err != nil {
				return
			} else {
				fmt.Println("成功创建数据库。")
				// 创建索引
				_, err = d.db.Exec("CREATE INDEX inx_password ON " + d.tableName + "(password);")
				if err != nil {
					fmt.Println("create index error->%q: %s\n", err, sqlStmt)
					return
				}
			}
		}
	}
	return nil
}

func Insert(password []string) (err error) {
	return d.insert(password)
}

// 插入数据
func (me *DataStruct) insert(password []string) (err error) {
	var tx *sql.Tx
	if tx, err = me.db.Begin(); err != nil {
		return err
	}
	defer tx.Commit()

	var stmt *sql.Stmt
	if stmt, err = tx.Prepare("insert into " + me.tableName + "(password,len,type,successCount,chineseHabit) values(?,?,?,?,?)"); err != nil {
		return err
	}
	defer stmt.Close()

	for i := range password {
		if _, err = stmt.Exec(password[i], len(password[i]), getType(password[i]), 0, 0); err != nil {
			// 这里只打印错误，不退出，让后续的正常insert
			fmt.Println(password[i], err)
		}
	}

	return nil
}

// 标记为中国人习惯
func (me *DataStruct) signChineseHabit(password []string) (err error) {
	var tx *sql.Tx
	if tx, err = me.db.Begin(); err != nil {
		return err
	}
	defer tx.Commit()

	var stmt *sql.Stmt
	if stmt, err = tx.Prepare("update " + me.tableName + " set chineseHabit=(?) where password=(?)"); err != nil {
		return err
	}
	defer stmt.Close()

	for i := range password {
		if _, err = stmt.Exec(1, password[i]); err != nil {
			// 这里只打印错误，不退出，让后续的正常
			fmt.Println(password[i], err)
		}
	}

	return nil
}

// 删除
func (me *DataStruct) del(password string) (err error) {
	sqlStmt := `delete from ` + me.tableName + ` where password = ` + password
	if _, err = me.db.Exec(sqlStmt); err != nil {
		return
	}
	return nil
}

func SuccessAdd(password string) (err error) {
	return d.update(password)
}

// 修改	successCount++
func (me *DataStruct) update(password string) (err error) {
	sqlStmt := `update ` + me.tableName + ` set successCount=successCount+1 where password = ` + password
	if _, err = me.db.Exec(sqlStmt); err != nil {
		return
	}
	return nil
}

func SelectValue(strWhere string, ch chan []byte, count *uint64) (err error) {
	return d.list(strWhere, ch, count)
}

// 列出所有文件
func (me *DataStruct) list(strWhere string, ch chan []byte, count *uint64) (err error) {
	var rows *sql.Rows
	strSql := "select password from " + me.tableName
	if strWhere != "" {
		strSql += " where "
		strSql += strWhere
	}
	if rows, err = d.db.Query(strSql); err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var password string
		if err = rows.Scan(&password); err != nil {
			return
		}
		ch <- []byte(password)
		*count++
	}

	if err = rows.Err(); err != nil {
		return
	}

	return
}
