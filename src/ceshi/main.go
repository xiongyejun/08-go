package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

type insertData struct {
	name  string
	t     string
	value float64
}

type dataStruct struct {
	DBPath string
	db     *sql.DB
}

var d *dataStruct

func init() {
	d = new(dataStruct)
	d.DBPath, _ = os.Getwd()
	d.DBPath += `\test.db`
}

func main() {
	if err := d.getDB(); err != nil {
		fmt.Println(err)
		return
	}
}

// 如何高效的实现
func (me *dataStruct) insertData(data []inserData) {

}

func (me *dataStruct) getDB() (err error) {
	if _, err = os.Stat(d.DBPath); err == nil {
		if me.db, err = sql.Open("sqlite3", d.DBPath); err != nil {
			return
		} else {
			fmt.Println("成功打开数据库。")
			return nil
		}
	} else {
		if me.db, err = sql.Open("sqlite3", d.DBPath); err != nil {
			return
		} else {
			sqlStmt := `create table ta (id integer not null primary key autoincrement, name text not null);`
			if _, err = d.db.Exec(sqlStmt); err != nil {
				return
			} else {
				fmt.Println("成功创建表a。")

				sqlStmt = `create table tb (aid integer not null, t DATE, value DOUBLE not null, primary key(aid,t));`
				if _, err = d.db.Exec(sqlStmt); err != nil {
					return
				} else {
					fmt.Println("成功创建表b。")
					return nil
				}
			}

		}
	}
}
