package inputBox

import (
	"fmt"

	"github.com/lxn/walk"
	"github.com/lxn/walk/declarative"
	"github.com/lxn/win"
)

var returnValues []string

//func main() {
//	fmt.Println(Show([]string{"开始地址", "修改文本", "输出编码", "ceshi"}, nil))
//}

func Show(values []string, owner walk.Form) []string {
	n := len(values)
	if n == 0 {
		return []string{}
	}
	form := new(walk.Dialog) // 用dialog会自动居中，还有设置窗口样式，还能绑定数据，暂时不会用

	// lable放在左边1个垂直框
	lbSplitter := new(walk.Splitter)
	lbs := make([]declarative.Widget, n+1)
	// textEdit放在右边1个垂直框
	tbSplitter := new(walk.Splitter)
	walkTb := make([]*walk.TextEdit, n)
	tbs := make([]declarative.Widget, n+1)
	// 初始化控件
	for i, v := range values {
		lbs[i] = declarative.Label{Text: v}
		tbs[i] = declarative.TextEdit{AssignTo: &walkTb[i]}
	}

	// 2个框里最后面都加1个pushbutton
	lbs[n] = declarative.PushButton{Text: "取消", OnClicked: func() {
		returnValues = make([]string, n)
		form.Close(1)
	}}
	tbs[n] = declarative.PushButton{Text: "确定", OnClicked: func() {
		returnValues = make([]string, n)
		for i := 0; i < n; i++ {
			returnValues[i] = walkTb[i].Text()
		}
		form.Close(1)
	}}

	mw := &declarative.Dialog{
		AssignTo: &form,
		Title:    "ceshi",
		Size:     declarative.Size{100, n * 20},

		Layout: declarative.HBox{},

		Children: []declarative.Widget{
			declarative.VSplitter{
				AssignTo: &lbSplitter,
				MaxSize:  declarative.Size{Width: 30},
				Children: lbs,
			},

			declarative.VSplitter{
				AssignTo: &tbSplitter,
				Children: tbs,
			},
		},
	}

	err := mw.Create(owner)
	if err != nil {
		fmt.Println(err)
		return []string{"0"}
	}
	//	form.SetOwner(owner)

	//	// 设置窗口样式
	//	style := win.GetWindowLong(form.Handle(), win.GWL_STYLE)
	//	style = style&(^win.WS_CAPTION)&(^win.WS_THICKFRAME) | win.WS_BORDER
	//	win.SetWindowLong(form.Handle(), win.GWL_STYLE, style)
	//	win.DrawMenuBar(form.Handle())

	//	// 设置窗口位置
	//	setMiddlePos(form.Handle(), owner)
	form.Run()

	return returnValues
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
