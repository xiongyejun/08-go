package main

import (
	"fmt"

	"github.com/lxn/walk"
	"github.com/lxn/walk/declarative"
	"github.com/lxn/win"
)

type control struct {
	form *walk.MainWindow

	// MenuItem
	miSelectFile *walk.Action
	miExit       *walk.Action

	miCompress   *walk.Action
	miUnCompress *walk.Action

	lbFileName *walk.Label
	lbTiShi    *walk.Label
	imageView  *walk.ImageView
	imageName  []string
	pImage     int
}

var ct *control
var mw *declarative.MainWindow

func initForm() {
	ct = new(control)

	mw = &declarative.MainWindow{
		AssignTo: &ct.form,
		Title:    "测试",
		Size:     declarative.Size{600, 600},
		Font:     declarative.Font{PointSize: 10},
		// 菜单
		MenuItems: []declarative.MenuItem{
			fileMenu(), // 文件下拉菜单
			actionMenu(),
		}, // MenuItems

		// 布局
		Layout: declarative.VBox{},
		// 控件
		Children: []declarative.Widget{ // widget小部件
			declarative.Label{
				AssignTo: &ct.lbFileName,
				Text:     "测试",
			},
			declarative.Label{
				AssignTo: &ct.lbTiShi,
				Text:     "提示",
			},

			declarative.ImageView{
				AssignTo: &ct.imageView,
				OnMouseUp: func(x, y int, button walk.MouseButton) {
					showImage()
				},
			},
		}, // MainWindow Children

	} // MainWindow

	// init end

	mw.Create()
	setMiddlePos(ct.form.Handle(), nil)
	ct.form.Run()
}

func showImage() {
	ct.lbTiShi.SetText(fmt.Sprintf("共%d张 当前第%d张", len(ct.imageName), ct.pImage+1))
	im, _ := walk.NewBitmapFromFile(ct.imageName[ct.pImage])
	ct.pImage++
	ct.pImage %= len(ct.imageName)
	ct.imageView.SetImage(im)

}

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
		(srcWidth-rect.Right)/2,
		(srcHeight-rect.Bottom)/2,
		rect.Right-rect.Left,
		rect.Bottom-rect.Top,
		win.SWP_SHOWWINDOW)
}
