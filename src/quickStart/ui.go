package main

import (
	"fmt"
	"time"

	"github.com/lxn/walk"
	"github.com/lxn/walk/declarative"
	"github.com/lxn/win"
)

type control struct {
	form *walk.MainWindow

	//Menu -- 文件
	miBackUp *walk.Action
	miQuit   *walk.Action

	//Menu -- 选项
	miTopmost *walk.Action

	// ContextMenuItems
	cmiAdd       *walk.Action
	cmiMoveUp    *walk.Action
	cmiMoveDown  *walk.Action
	cmiAddRoot   *walk.Action
	cmiDel       *walk.Action
	cmiEdit      *walk.Action
	cmiExpandAll *walk.Action

	// treeView
	treeView                *walk.TreeView
	treeModle               *TreeModle
	treeViewMouseMoveHandle int

	// notifyicon
	ni *walk.NotifyIcon
}

var ct *control = new(control)

func uiInit() {
	newTreeModle()

	mw := declarative.MainWindow{
		AssignTo: &ct.form,
		Title:    "go",
		Size:     declarative.Size{300, int(win.GetSystemMetrics(win.SM_CYSCREEN) - 300)},
		Icon: func() *walk.Icon {
			i, _ := walk.Resources.Icon("\\image\\go.ico")
			return i
		}(),
		Font: declarative.Font{PointSize: 9, Family: "宋体"},

		MenuItems: []declarative.MenuItem{
			fileMenu(),
			optionMenu(),
		},

		Layout: declarative.VBox{Margins: declarative.Margins{1, 1, 1, 1}},
		Children: []declarative.Widget{
			initTreeView(),
		},
	}

	if err := initNotifyIcon(); err != nil {
		walk.MsgBox(nil, "", err.Error(), walk.MsgBoxIconInformation)
		fmt.Println(err)
		return
	}
	defer ct.ni.Dispose()

	if err := mw.Create(); err != nil {
		walk.MsgBox(nil, "", err.Error(), walk.MsgBoxIconInformation)
		fmt.Println(err)
		return
	}

	setFormStyle(ct.form)
	setHideToRight(ct.form.Handle(), nil)
	//	ct.miTopmost.SetChecked(true)
	ct.treeViewMouseMoveHandle = ct.treeView.MouseMove().Attach(treeMouseIn)
	go mouseMove() // 移动鼠标，避免锁屏
	ct.form.Run()
}

// 鼠标进入，就设置窗体在右边显示出来
func treeMouseIn(x, y int, button walk.MouseButton) {
	setShowToRight()
	ct.treeView.MouseMove().Detach(ct.treeViewMouseMoveHandle)
	// 启动1个进程，获取鼠标的坐标，移出窗体的话，就设置窗体到右边
	go func() { // 不使用上面的formX是因为有可能会手动移动窗体
		for {
			x, y := getMousePos()
			if x < ct.form.X() || x > ct.form.X()+ct.form.Width() || y < ct.form.Y() || y > ct.form.Y()+ct.form.Height() {
				setHideToRight(ct.form.Handle(), nil)
				// 隐藏了就再把MouseMove事件加进去
				ct.treeViewMouseMoveHandle = ct.treeView.MouseMove().Attach(treeMouseIn)
				return
			}
			time.Sleep(1e8)
		}
	}()
}

func setFormStyle(form *walk.MainWindow) {
	// 设置窗口样式
	style := win.GetWindowLong(form.Handle(), win.GWL_STYLE)
	style = style&(^win.WS_THICKFRAME)&(^win.WS_MAXIMIZEBOX)&(^win.WS_MINIMIZEBOX) | win.WS_EX_TOOLWINDOW // 设置不能拉伸 // WS_EX_TOOLWINDOW没有效果！
	win.SetWindowLong(form.Handle(), win.GWL_STYLE, style)
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

func setShowToRight() {
	formX := int(win.GetSystemMetrics(win.SM_CXSCREEN)) - ct.form.Width()
	ct.form.SetX(formX)
}

// 设置窗口在右边隐藏
func setHideToRight(hwnd win.HWND, owner walk.Form) {
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
		srcWidth-8,
		0,
		200,
		srcHeight-25,
		win.SWP_SHOWWINDOW)
}
