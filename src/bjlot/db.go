package main

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
)

const TABLE_DATA string = "datas"

type insertData struct {
	yearmonth int
	drawno    int

	Competitions string  // Competitions:"J联赛"
	Datetime     string  // Datetime:"09-05-16 11:55:00"
	Dispatchamt  float64 // Dispatchamt:"70.29"
	Guestteam    string  // Guestteam:"大宫"
	Handicap     int     // Handicap:"0"
	Hitcount     float64 // Hitcount:"9"
	Hostteam     string  // Hostteam:"名古屋"
	Matchno      int     // Matchno:"1"
	Result       int     // Result:"2"
	Score        string  // Score:"1:1"
	Spvalue      float64 // Spvalue:"6.011832"
	Stake        float64 // Stake:"7.8153"

	two   string
	three string
	four  string
	five  string
	six   string
	seven string
	eight string
	nine  string
	ten   string

	// sp之和
	twoSUMsp   float64
	threeSUMsp float64
	fourSUMsp  float64
	fiveSUMsp  float64
	sixSUMsp   float64
	sevenSUMsp float64
	eightSUMsp float64
	nineSUMsp  float64
	tenSUMsp   float64

	// sp之积
	twoPRODUCTsp   float64
	threePRODUCTsp float64
	fourPRODUCTsp  float64
	fivePRODUCTsp  float64
	sixPRODUCTsp   float64
	sevenPRODUCTsp float64
	eightPRODUCTsp float64
	ninePRODUCTsp  float64
	tenPRODUCTsp   float64
}

// 打开数据库
func (me *dataStruct) getDB() (err error) {
	if _, err = os.Stat(d.DBPath); err == nil {
		if me.db, err = sql.Open("sqlite3", d.DBPath); err != nil {
			return
		} else {
			fmt.Println("成功连接数据库。")
			return nil
		}
	} else {
		// 不存在数据库的情况下进行创建
		if me.db, err = sql.Open("sqlite3", d.DBPath); err != nil {
			return
		} else {
			// 								year+month=201101								    {Competitions:"J联赛", Datetime:"09-05-16 11:55:00", Dispatchamt:"70.29", Guestteam:"大宫", Handicap:"0", Hitcount:"9", Hostteam:"名古屋", Matchno:"1", Result:"2", Score:"1:1", Spvalue:"6.011832", Stake:"7.8153"}
			sqlStmt := `create table ` + TABLE_DATA + ` (yearmonth integer not null,drawno integer not null,Competitions text,Datetime text,Dispatchamt DOUBLE,Guestteam text,Handicap integer,Hitcount float64,Hostteam text,Matchno integer not null,Result integer,Score text,Spvalue DOUBLE,Stake DOUBLE,two varchar(2)  default '',three varchar(3)  default '',four varchar(4)  default '',five varchar(5)  default '',six varchar(6)  default '',seven varchar(7)  default '',eight varchar(8)  default '',nine varchar(9)  default '',ten varchar(10) default '',twoSUMsp DOUBLE,threeSUMsp DOUBLE,fourSUMsp DOUBLE,fiveSUMsp DOUBLE,sixSUMsp DOUBLE,sevenSUMsp DOUBLE,eightSUMsp DOUBLE,nineSUMsp DOUBLE,tenSUMsp DOUBLE,twoPRODUCTsp DOUBLE,threePRODUCTsp DOUBLE,fourPRODUCTsp DOUBLE,fivePRODUCTsp DOUBLE,sixPRODUCTsp DOUBLE,sevenPRODUCTsp DOUBLE,eightPRODUCTsp DOUBLE,ninePRODUCTsp DOUBLE,tenPRODUCTsp DOUBLE, primary key (yearmonth,drawno,Matchno));`
			if _, err = d.db.Exec(sqlStmt); err != nil {
				return
			} else {
				fmt.Println("成功创建表datas。")
				return nil
			}
		}
	}
}

// 插入数据
func (me *dataStruct) insertData(data []*insertData) (err error) {
	var tx *sql.Tx
	if tx, err = me.db.Begin(); err != nil {
		return err
	}
	defer tx.Commit()

	var stmt *sql.Stmt
	if stmt, err = tx.Prepare("insert into " + TABLE_DATA + "(yearmonth,drawno,Competitions,Datetime,Dispatchamt,Guestteam,Handicap,Hitcount,Hostteam,Matchno,Result,Score,Spvalue,Stake,two,three,four,five,six,seven,eight,nine,ten,twoSUMsp,threeSUMsp,fourSUMsp,fiveSUMsp,sixSUMsp,sevenSUMsp,eightSUMsp,nineSUMsp,tenSUMsp,twoPRODUCTsp,threePRODUCTsp,fourPRODUCTsp,fivePRODUCTsp,sixPRODUCTsp,sevenPRODUCTsp,eightPRODUCTsp,ninePRODUCTsp,tenPRODUCTsp) values(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)"); err != nil {
		return err
	}
	defer stmt.Close()

	for i := range data {
		if _, err = stmt.Exec(data[i].yearmonth, data[i].drawno, data[i].Competitions, data[i].Datetime, data[i].Dispatchamt, data[i].Guestteam, data[i].Handicap, data[i].Hitcount, data[i].Hostteam, data[i].Matchno, data[i].Result, data[i].Score, data[i].Spvalue, data[i].Stake, data[i].two, data[i].three, data[i].four, data[i].five, data[i].six, data[i].seven, data[i].eight, data[i].nine, data[i].ten, data[i].twoSUMsp, data[i].threeSUMsp, data[i].fourSUMsp, data[i].fiveSUMsp, data[i].sixSUMsp, data[i].sevenSUMsp, data[i].eightSUMsp, data[i].nineSUMsp, data[i].tenSUMsp, data[i].twoPRODUCTsp, data[i].threePRODUCTsp, data[i].fourPRODUCTsp, data[i].fivePRODUCTsp, data[i].sixPRODUCTsp, data[i].sevenPRODUCTsp, data[i].eightPRODUCTsp, data[i].ninePRODUCTsp, data[i].tenPRODUCTsp); err != nil {
			return
		}
	}

	return nil
}

// 获取数据库中最新的数据
func (me *dataStruct) getNewData() (strYearmonth string, recordDrawno map[int]int, err error) {
	var rows *sql.Rows
	if rows, err = d.db.Query("select MAX(yearmonth) from " + TABLE_DATA); err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var yearmonth int
		if err = rows.Scan(&yearmonth); err != nil {
			return
		}
		strYearmonth = strconv.Itoa(yearmonth)
	}

	if err = rows.Err(); err != nil {
		return
	}

	if rows, err = d.db.Query("select drawno from " + TABLE_DATA + " WHERE yearmonth=" + strYearmonth); err != nil {
		return
	}
	recordDrawno = make(map[int]int)
	for rows.Next() {
		var drawno int
		if err = rows.Scan(&drawno); err != nil {
			return
		}
		recordDrawno[drawno] = 0
	}

	if err = rows.Err(); err != nil {
		return
	}

	return
}

//// 读取文件bytes，保存在当前程序的路径下，并打开
//func (me *dataStruct) show(pID int) (err error) {
//	var name string
//	var ext string
//	var ok bool

//	if pID >= len(me.pID) {
//		return errors.New("不存在的索引。")
//	}
//	id := me.pID[pID]
//	// 先判断是否已经存在了
//	if name, ok = me.dicShow[id]; !ok {
//		var stmt *sql.Stmt
//		if stmt, err = d.db.Prepare("select name,ext,bytes from " + me.tableName + " where id = ?"); err != nil {
//			return
//		}
//		defer stmt.Close()

//		var bi interface{}
//		if err = stmt.QueryRow(strconv.Itoa(id)).Scan(&name, &ext, &bi); err != nil {
//			return
//		}
//		// 文件保存路径
//		name = me.fileSavePath + strconv.Itoa(id) + ext
//		// 读取文件的byte
//		if b, ok := bi.([]byte); ok {
//			// 解密byte
//			if b, err = desDecrypt(b, d.key); err != nil {
//				return
//			}
//			// 保存文件
//			if err = ioutil.WriteFile(name, b, 0666); err != nil {
//				return
//			}
//			// 记录打开过的，退出时删除
//			me.dicShow[id] = name

//		}

//	} /* else {
//		fmt.Println("已经有了")
//	}*/

//	//	fmt.Println(name)
//	return openFolderFile(name)
//}
