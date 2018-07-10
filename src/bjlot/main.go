// http://www.bjlot.com/
// 北京体彩网

package main

import (
	"fmt"
)

type dataStructItem struct {
	url      string   // 网抓的地址
	attr     string   // 属性
	attrItem string   // 属性
	rets     []string // 返回的内容
}

type dataStruct struct {
	year        *dataStructItem
	month       []*dataStructItem
	drawno      []*dataStructItem
	drawnoCount int

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

	d.data = new(datas)
}

func main() {
	//	getDrawno()

	if err := d.data.getData(`http://www.bjlot.com/data/230/draw/2009/90511.js`); err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("%#v\r\n", d.data)
	}

	fmt.Printf("%#v\r\n", d.data.Drawresult[0])

}

func getDrawno() (err error) {
	// 获取年
	if err = d.year.getValue(""); err != nil {
		return
	}
	fmt.Println(d.year.rets)
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
		fmt.Println(d.year.rets[i], d.month[i].rets)

		d.drawnoCount += len(d.month[i].rets)
	}

	// 获取drawno
	d.drawno = make([]*dataStructItem, d.drawnoCount)
	d.drawnoCount = 0
	for i := range d.month {
		for j := range d.month[i].rets {
			d.drawno[d.drawnoCount] = new(dataStructItem)
			d.drawno[d.drawnoCount].url = `http://www.bjlot.com/data/230/control/drawnolist_` + d.month[i].rets[j] + `.js`
			d.drawno[d.drawnoCount].attr = "drawnolist"
			d.drawno[d.drawnoCount].attrItem = "drawno"

			if err = d.drawno[d.drawnoCount].getValue(d.month[i].rets[j][:4] + `/`); err != nil {
				return
			}
			d.drawnoCount++
		}
	}

	return nil
}
