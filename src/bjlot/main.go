// http://www.bjlot.com/
// 北京体彩网

package main

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"

	termbox "github.com/nsf/termbox-go"
)

type dataStructItem struct {
	url      string   // 网抓的地址
	attr     string   // 属性
	attrItem string   // 属性
	rets     []string // 返回的内容
}

type dataStruct struct {
	year      *dataStructItem
	month     []*dataStructItem
	drawno    []*dataStructItem
	yearmonth []string

	recordYear      string      // 数据库记录的最大的年
	recordYearmonth string      // 上次下载的数据的年月
	recordDrawno    map[int]int // recordYearmonth下，数据库已经抓取的drawno

	data *datas

	DBPath string
	db     *sql.DB
}

var d *dataStruct

func init() {
	d = new(dataStruct)
	d.year = new(dataStructItem)
	d.year.url = `http://www.bjlot.com/data/230/control/drawyearlist.js`
	d.year.attr = "drawyears"
	d.year.attrItem = "year"

	d.DBPath, _ = filepath.Abs(filepath.Dir(os.Args[0]))
	d.DBPath += `\bjlot.db`
	fmt.Println("数据库地址：", d.DBPath)

	if err := termbox.Init(); err != nil {
		panic(err)
	}
	//	termbox.SetCursor(0, 0)
	//	termbox.HideCursor()

}

func main() {
	defer pause()
	var err error

	if err = d.getDB(); err != nil {
		fmt.Println(err)
		return
	}
	defer d.db.Close()

	// 数据库为空的时候初始化
	if d.recordYearmonth, d.recordDrawno, err = d.getNewData(); err != nil {
		d.recordYear = "0000"
		d.recordYearmonth = "000000"
		d.recordDrawno = make(map[int]int)
	}
	d.recordYear = d.recordYearmonth[:4]

	// 获取要下载的drawno列表
	getDrawno()

	for i := range d.drawno {
		var iyearmonth int = 0
		var idrawno int = 0
		var err error
		if iyearmonth, err = strconv.Atoi(d.yearmonth[i]); err != nil {
			fmt.Println(err)
			return
		}

		for j := range d.drawno[i].rets {
			//												####/###### 年/no
			if idrawno, err = strconv.Atoi(d.drawno[i].rets[j][5:]); err != nil {
				fmt.Println(err)
				return
			}

			// 数据库不存在的进行抓取
			if _, ok := d.recordDrawno[idrawno]; !ok {
				fmt.Println("\r\n正在下载数据：", d.yearmonth[i], idrawno)

				d.data = new(datas)
				// 下载网页数据
				if err := d.data.getData(`http://www.bjlot.com/data/230/draw/` + d.drawno[i].rets[j] + `.js`); err != nil {
					fmt.Println("\r\n从网页读取数据出错：", err)
					return
				}

				// data转换为InsertData结构 *datas存放到数据库
				var data []*insertData
				if data, err = d.data.toInsertData(iyearmonth, idrawno); err != nil {
					fmt.Println(err)
					return
				}

				// 插入数据到数据库
				if err = d.insertData(data); err != nil {
					fmt.Println("\r\n数据存入数据库出错：", err)
				} else {
					fmt.Println("成功下载数据：", d.yearmonth[i], idrawno)
				}
			}
		}

	}

	fmt.Println("数据下载完成。")
}

func (me *datas) toInsertData(yearmonth int, idrawno int) (data []*insertData, err error) {
	data = make([]*insertData, len(me.Drawresult))

	// result有可能出现*，替换为_,sqlite的单字通配符
	for i := range me.Drawresult {
		if "*" == me.Drawresult[i].Result {
			me.Drawresult[i].Result = "_"
		}
	}

	// result有可能出现7+这种情况
	for i := range me.Drawresult {
		data[i] = new(insertData)
		me.Drawresult[i].Result = me.Drawresult[i].Result[:1]
	}
	// 获取Spvalue---要计算和、积
	for i := range me.Drawresult {
		if data[i].Spvalue, err = strconv.ParseFloat(me.Drawresult[i].Spvalue, 64); err != nil {
			return nil, errors.New(me.Drawresult[i].Matchno + " Spvalue:" + err.Error())
		}
		data[i].twoPRODUCTsp = 1.0
		data[i].threePRODUCTsp = 1.0
		data[i].fourPRODUCTsp = 1.0
		data[i].fivePRODUCTsp = 1.0
		data[i].sixPRODUCTsp = 1.0
		data[i].sevenPRODUCTsp = 1.0
		data[i].eightPRODUCTsp = 1.0
		data[i].ninePRODUCTsp = 1.0
		data[i].tenPRODUCTsp = 1.0

	}

	// 做好2-10的连接，后面查询的时候方便判断
	var n int = 2
	for i := 0; i < len(me.Drawresult)-n+1; i++ {
		for j := i; j < i+n; j++ {
			data[i].two += me.Drawresult[j].Result
			data[i].twoPRODUCTsp *= data[j].Spvalue
			data[i].twoSUMsp += data[j].Spvalue
		}
	}
	n++
	for i := 0; i < len(me.Drawresult)-n+1; i++ {
		for j := i; j < i+n; j++ {
			data[i].three += me.Drawresult[j].Result
			data[i].threePRODUCTsp *= data[j].Spvalue
			data[i].threeSUMsp += data[j].Spvalue
		}
	}
	n++
	for i := 0; i < len(me.Drawresult)-n+1; i++ {
		for j := i; j < i+n; j++ {
			data[i].four += me.Drawresult[j].Result
			data[i].fourPRODUCTsp *= data[j].Spvalue
			data[i].fourSUMsp += data[j].Spvalue
		}
	}
	n++
	for i := 0; i < len(me.Drawresult)-n+1; i++ {
		for j := i; j < i+n; j++ {
			data[i].five += me.Drawresult[j].Result
			data[i].fivePRODUCTsp *= data[j].Spvalue
			data[i].fiveSUMsp += data[j].Spvalue
		}
	}
	n++
	for i := 0; i < len(me.Drawresult)-n+1; i++ {
		for j := i; j < i+n; j++ {
			data[i].six += me.Drawresult[j].Result
			data[i].sixPRODUCTsp *= data[j].Spvalue
			data[i].sixSUMsp += data[j].Spvalue
		}
	}
	n++
	for i := 0; i < len(me.Drawresult)-n+1; i++ {
		for j := i; j < i+n; j++ {
			data[i].seven += me.Drawresult[j].Result
			data[i].sevenPRODUCTsp *= data[j].Spvalue
			data[i].sevenSUMsp += data[j].Spvalue
		}
	}
	n++
	for i := 0; i < len(me.Drawresult)-n+1; i++ {
		for j := i; j < i+n; j++ {
			data[i].eight += me.Drawresult[j].Result
			data[i].eightPRODUCTsp *= data[j].Spvalue
			data[i].eightSUMsp += data[j].Spvalue
		}
	}
	n++
	for i := 0; i < len(me.Drawresult)-n+1; i++ {
		for j := i; j < i+n; j++ {
			data[i].nine += me.Drawresult[j].Result
			data[i].ninePRODUCTsp *= data[j].Spvalue
			data[i].nineSUMsp += data[j].Spvalue
		}
	}
	n++
	for i := 0; i < len(me.Drawresult)-n+1; i++ {
		for j := i; j < i+n; j++ {
			data[i].ten += me.Drawresult[j].Result
			data[i].tenPRODUCTsp *= data[j].Spvalue
			data[i].tenSUMsp += data[j].Spvalue
		}
	}
	n++

	// 再替换为数字
	for i := range me.Drawresult {
		if "_" == me.Drawresult[i].Result {
			me.Drawresult[i].Result = "9"
		}
	}

	for i := range me.Drawresult {
		data[i].yearmonth = yearmonth
		data[i].drawno = idrawno

		data[i].Competitions = me.Drawresult[i].Competitions
		data[i].Datetime = me.Drawresult[i].Datetime
		if data[i].Matchno, err = strconv.Atoi(me.Drawresult[i].Matchno); err != nil {
			return nil, errors.New(data[i].Datetime + ":" + err.Error())
		}
		if data[i].Dispatchamt, err = strconv.ParseFloat(me.Drawresult[i].Dispatchamt, 64); err != nil {
			return nil, errors.New(me.Drawresult[i].Matchno + " Dispatchamt:" + err.Error())
		}
		data[i].Guestteam = me.Drawresult[i].Guestteam
		if data[i].Handicap, err = strconv.Atoi(me.Drawresult[i].Handicap); err != nil {
			return nil, errors.New(me.Drawresult[i].Matchno + " Handicap:" + err.Error())
		}
		if data[i].Hitcount, err = strconv.ParseFloat(me.Drawresult[i].Hitcount, 64); err != nil {
			return nil, errors.New(me.Drawresult[i].Matchno + " Hitcount:" + err.Error())
		}

		data[i].Hostteam = me.Drawresult[i].Hostteam

		if data[i].Result, err = strconv.Atoi(me.Drawresult[i].Result[:1]); err != nil {
			return nil, errors.New(me.Drawresult[i].Matchno + " Result:" + err.Error())
		}
		data[i].Score = me.Drawresult[i].Score
		if data[i].Spvalue, err = strconv.ParseFloat(me.Drawresult[i].Spvalue, 64); err != nil {
			return nil, errors.New(me.Drawresult[i].Matchno + " Spvalue:" + err.Error())
		}
		if data[i].Stake, err = strconv.ParseFloat(me.Drawresult[i].Stake, 64); err != nil {
			return nil, errors.New(me.Drawresult[i].Matchno + " Stake:" + err.Error())
		}
	}

	return
}

func getDrawno() (err error) {
	var i int = 0

	// 获取年
	if err = d.year.getValue(""); err != nil {
		return
	}

	// 与数据库中最晚的年份对比，只需要下载大于等于最晚的年份
	sort.Strings(d.year.rets)
	fmt.Println("网站年份：\r\n", d.year.rets)
	fmt.Println("数据库中最新年份：", d.recordYear)

	i = len(d.year.rets) - 1
	for d.year.rets[i] != d.recordYear && i > 0 {
		i--
	}
	if i > 0 {
		d.year.rets = d.year.rets[i:]
	}
	fmt.Println("需要抓取的年份：\r\n", d.year.rets)

	// 获取月
	d.month = make([]*dataStructItem, len(d.year.rets))
	for i := range d.month {
		d.month[i] = new(dataStructItem)
		d.month[i].url = `http://www.bjlot.com/data/230/control/` + d.year.rets[i] + `.js`
		d.month[i].attr = "monthlist"
		d.month[i].attrItem = "month"

		if err = d.month[i].getValue(d.year.rets[i]); err != nil {
			return
		}
	}
	for i := range d.month {
		d.yearmonth = append(d.yearmonth, d.month[i].rets...)
	}
	// 排序----升序
	sort.Strings(d.yearmonth)
	// 与记录的对比，只要找后面的就可以
	i = len(d.yearmonth) - 1
	for d.yearmonth[i] != d.recordYearmonth && i > 0 {
		i--
	}
	if i > 0 {
		d.yearmonth = d.yearmonth[i:]
	}
	fmt.Println("需要抓取的月份数据：\r\n", d.yearmonth)

	// 获取drawno
	d.drawno = make([]*dataStructItem, len(d.yearmonth))
	for i := range d.yearmonth {
		d.drawno[i] = new(dataStructItem)
		d.drawno[i].url = `http://www.bjlot.com/data/230/control/drawnolist_` + d.yearmonth[i] + `.js`
		d.drawno[i].attr = "drawnolist"
		d.drawno[i].attrItem = "drawno"

		if err = d.drawno[i].getValue(d.yearmonth[i][:4] + `/`); err != nil {
			return
		}
	}

	return nil
}

func pause() {
	fmt.Println("\r\n请按任意键继续...")
Loop:
	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			break Loop
		}
	}
}
