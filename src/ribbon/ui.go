package main

import (
	"fmt"

	"github.com/lxn/walk"
	"github.com/lxn/walk/declarative"
	"github.com/lxn/win"
)

type uiControl struct {
	form *walk.MainWindow

	//Menu -- 文件
	miSelectFile *walk.Action
	miSaveXml    *walk.Action
	miOpenFile   *walk.Action
	miQuit       *walk.Action

	// 插入
	miCustomUI *walk.Action
	miButton   *walk.Action
	miCallBack *walk.Action

	//Menu -- 选项
	miBak     *walk.Action // 备份原文件
	miTopMost *walk.Action

	// ContextMenuItems
	cmiAdd       *walk.Action
	cmiAddRoot   *walk.Action
	cmiDel       *walk.Action
	cmiEdit      *walk.Action
	cmiExpandAll *walk.Action

	// treeView
	treeView  *walk.TreeView
	treeModle *TreeModle

	// tb
	tbXml        *walk.TextEdit
	fileName     string
	customUIName string // 记录读取到的customUI，有3种情况，1没有，2custom，3custom14
}

var ct *uiControl = new(uiControl)

func uiInit() {
	newTreeModle()

	mw := declarative.MainWindow{
		AssignTo: &ct.form,
		Title:    "go",
		Size:     declarative.Size{800, 500},
		Icon: func() *walk.Icon {
			i, _ := walk.Resources.Icon("image\\favicon.ico")
			return i
		}(),

		MenuItems: []declarative.MenuItem{
			fileMenu(),
			optionMenu(),
			insertMenu(),
		},

		Layout: declarative.VBox{Margins: declarative.Margins{1, 1, 1, 1}},
		Children: []declarative.Widget{
			declarative.HSplitter{
				Children: []declarative.Widget{
					initTreeView(),

					declarative.TextEdit{
						AssignTo: &ct.tbXml,
						Font:     declarative.Font{PointSize: 10},
						HScroll:  true,
						VScroll:  true,
					},
				},
			},
		},
	}
	if err := mw.Create(); err != nil {
		walk.MsgBox(nil, "", err.Error(), walk.MsgBoxIconInformation)
		fmt.Println(err)
		return
	}

	setMiddlePos(ct.form.Handle(), nil)
	// 最后一个参数若取1，表示窗口大小保持不变，取2表示保持位置不变，因此，取3（=1＋2）表示大小和位置均保持不变，取0表示将窗口的大小和位置改变为指定值
	win.SetWindowPos(ct.form.Handle(), win.HWND_NOTOPMOST, 0, 0, 0, 0, 1|2)
	ct.form.Run()
}

// 设置窗口居中
func setMiddlePos(hwnd win.HWND, owner walk.Form) {
	var srcWidth, srcHeight int32

	if owner == nil {
		srcWidth = win.GetSystemMetrics(win.SM_CXSCREEN)
		srcHeight = win.GetSystemMetrics(win.SM_CYSCREEN)
	} else {
		srcWidth = int32(owner.Width()) + 2*int32(owner.X())
		srcHeight = int32(owner.Height()) + 2*int32(owner.Y())
	}

	rect := new(win.RECT)
	win.GetWindowRect(hwnd, rect)
	win.SetWindowPos(hwnd, win.HWND_TOPMOST,
		(srcWidth-(rect.Right-rect.Left))/2,
		(srcHeight-(rect.Bottom-rect.Top))/2,
		rect.Right-rect.Left,
		rect.Bottom-rect.Top,
		win.SWP_SHOWWINDOW)
}
