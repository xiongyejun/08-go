// 用tableview来显示信息

package main

import (
	"fmt"
	"pkgMyPkg/compdocFile"
	"strings"

	"github.com/lxn/walk"
	"github.com/lxn/walk/declarative"
	"github.com/lxn/win"
)

func tableviewShowInfo() declarative.TableView {
	return declarative.TableView{
		AssignTo: &ct.tableview,
		//		StretchFactor: 2,

		Columns: []declarative.TableViewColumn{
			declarative.TableViewColumn{
				DataMember: "Name",
				Width:      200,
			},

			declarative.TableViewColumn{
				DataMember: "Type",
				Width:      100,
			},
		}, // TableView Columns
		Model:                 ct.tableModle,
		MinSize:               declarative.Size{300, 500},
		MaxSize:               declarative.Size{300, 500},
		OnCurrentIndexChanged: tableviewShowInfoIndexChanged,
	}
}

func tableviewShowInfoIndexChanged() {
	if index := ct.tableview.CurrentIndex(); index > -1 {
		if !ct.miShowCode.Checked() || ct.tableModle.items[index].Type == "流" {
			strASCII := showByte(ct.tableModle.items[index].Name)
			selectAttribut(index, strASCII)
		} else if strings.HasSuffix(ct.tableModle.items[index].Type, "模块流") {
			showCode(ct.tableModle.items[index].Name)
		} else {
			ct.tb.SetText("")
		}
	}
}

// 如果是模块流的时候，就找到Attribut模块开始的位置
func selectAttribut(tableviewCurrentindex int, strASCII string) {
	sc := new(win.SCROLLINFO)
	sc.FMask = win.SIF_ALL
	win.GetScrollInfo(ct.tb.Handle(), win.SB_VERT, sc)

	if ct.tableModle.items[tableviewCurrentindex].Type == "模块流" {
		start := strings.Index(strASCII, "Attribut")
		// start是Attribut在strASCII里的位置，在tb里，一行的字数量是不同的
		lines := start / 16 // 占的行数
		tbstart := (2*9 + 3*16 + 16 + 2) * lines
		iMod := start % 16 // 剩下的个数
		if iMod > 0 {
			tbstart = tbstart + iMod + (2*9 + 3*16 + 16 + 2) - 16
		}

		fmt.Println(start, tbstart, len(ct.tb.Text()))
		ct.tb.SetFocus()
		ct.tb.SetTextSelection(tbstart, tbstart+200)
		sc.NPos = sc.NMax * int32(start) / int32(len(strASCII))
		win.SetScrollInfo(ct.tb.Handle(), win.SB_VERT, sc, true)
		ct.tb.SetFocus()
		fmt.Println("ok")
	}
}

type tableItemModle struct {
	walk.SortedReflectTableModelBase
	items []*compdocFile.OutStruct
}

// 这一句做什么，不懂
var _ walk.ReflectTableModel = new(tableItemModle)

func (me *tableItemModle) addItem(value []*compdocFile.OutStruct) {
	me.items = value

	me.PublishRowsReset()
}

func (me *tableItemModle) Items() interface{} {
	return me.items
}
