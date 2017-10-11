// 用UI来做VBA模块的查看，破解工程密码等
package main

import (
	"fmt"
	"pkgMySelf/compdocFile"
	_ "runtime/cgo"
	"strings"

	"github.com/lxn/walk"
	"github.com/lxn/walk/declarative"
)

type control struct {
	form *walk.MainWindow

	// MenuItem
	miSelectFile *walk.Action
	miShowModule *walk.Action
	miShowCode   *walk.Action

	vbox       *declarative.VBox
	lbFileName *walk.Label
	tb         *walk.TextEdit
}

var ct *control
var mw *declarative.MainWindow
var cf compdocFile.CF
var cfflag bool

func init() {
	ct = new(control)
	ct.lbFileName = &walk.Label{}
	ct.tb = &walk.TextEdit{}

	ct.vbox = &declarative.VBox{}
	//	ct.vbox.Margins = declarative.Margins{Right: 300, Bottom: 200}

	mw = &declarative.MainWindow{
		AssignTo: &ct.form,
		Title:    "测试",
		Size:     declarative.Size{300, 300},
		MaxSize:  declarative.Size{300, 300},
		// 菜单
		MenuItems: []declarative.MenuItem{
			declarative.Menu{
				Text: "&File",
				Items: []declarative.MenuItem{
					declarative.Action{
						AssignTo:    &ct.miSelectFile,
						Text:        "&Select",
						OnTriggered: selectFile, // 触发，相当于click
					},

					declarative.Action{
						AssignTo:    &ct.miShowModule,
						Text:        "&ShowModule",
						OnTriggered: showModule,
					},

					declarative.Action{
						AssignTo:    &ct.miShowCode,
						Text:        "&ShowCode",
						OnTriggered: showCode,
					},
				}, // Items
			},
		}, // MenuItems

		// 布局
		Layout: ct.vbox,
		// 控件
		Children: []declarative.Widget{ // widget小部件
			declarative.Label{
				AssignTo: &ct.lbFileName,
				Text:     "测试",
			},

			declarative.TextEdit{
				AssignTo: &ct.tb,
				Enabled:  false,
			},
		}, // Children
	} // MainWindow

	// init end
}

func main() {
	//	mw.Create()
	mw.Run()
}

func selectFile() {
	fd := new(walk.FileDialog)
	fd.ShowOpen(ct.form)
	ct.lbFileName.SetText(fd.FilePath)

	cfflag = true
	if compdocFile.IsCompdocFile(fd.FilePath) {
		cf = compdocFile.NewXlsFile(fd.FilePath)
	} else if compdocFile.IsZip(fd.FilePath) {
		cf = compdocFile.NewZipFile(fd.FilePath)
	} else {
		walk.MsgBox(ct.form, "title", "未知文件："+fd.FilePath, walk.MsgBoxIconInformation)
		cfflag = false
		return
	}
	compdocFile.CFInit(cf)
	return
}

func showModule() {
	if cfflag {
		modules := cf.GetModuleName()
		ct.tb.SetText(strings.Join(modules, "\r\n"))
	}
}
func showCode() {
	if cfflag {
		str := cf.GetAllCode()
		fmt.Println(str)
		ct.tb.SetText(str)
	}
}
