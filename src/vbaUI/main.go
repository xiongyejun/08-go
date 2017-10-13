// 用UI来做VBA模块的查看，破解工程密码等
package main

import (
	//	"fmt"
	//	"bytes"
	//	"fmt"
	"pkgMyPkg/compdocFile"
	//	"strings"

	"github.com/lxn/walk"
	"github.com/lxn/walk/declarative"
)

type TableItem struct {
	ModuleName string
	ModuleType string
}

type tableItemModle struct {
	walk.SortedReflectTableModelBase
	items []*TableItem
}

// 这一句做什么，不懂
var _ walk.ReflectTableModel = new(tableItemModle)

type control struct {
	form *walk.MainWindow

	// MenuItem
	miSelectFile       *walk.Action
	miUnProtectProject *walk.Action // 破解工程密码

	lbFileName *walk.Label

	hsplitter  *walk.Splitter
	tableview  *walk.TableView // tableview显示模块名称、模块type
	tableModle *tableItemModle // tableview添加items
	tb         *walk.TextEdit
}

var ct *control
var mw *declarative.MainWindow
var cf compdocFile.CF
var cfflag bool

func init() {
	ct = new(control)
	ct.tableModle = new(tableItemModle)

	mw = &declarative.MainWindow{
		AssignTo: &ct.form,
		Title:    "测试",
		Size:     declarative.Size{600, 600},
		Font:     declarative.Font{PointSize: 10},
		// 菜单
		MenuItems: []declarative.MenuItem{
			declarative.Menu{
				Text: "&Action",
				Items: []declarative.MenuItem{
					declarative.Action{
						AssignTo:    &ct.miSelectFile,
						Text:        "&选择文件",
						OnTriggered: selectFile, // 触发，相当于click
					},

					declarative.Action{
						AssignTo:    &ct.miUnProtectProject,
						Text:        "&破解工程密码",
						OnTriggered: unProtectProject,
					},
				}, // Items
			},
		}, // MenuItems

		// 布局
		Layout: declarative.VBox{},
		// 控件
		Children: []declarative.Widget{ // widget小部件
			declarative.Label{
				AssignTo: &ct.lbFileName,
				Text:     "测试",
			},

			declarative.HSplitter{
				AssignTo: &ct.hsplitter,
				Children: []declarative.Widget{
					declarative.TableView{
						AssignTo:      &ct.tableview,
						StretchFactor: 2,

						Columns: []declarative.TableViewColumn{
							declarative.TableViewColumn{
								DataMember: "ModuleName",
								Width:      200,
							},

							declarative.TableViewColumn{
								DataMember: "ModuleType",
								Width:      100,
							},
						}, // TableView Columns
						Model:   ct.tableModle,
						MinSize: declarative.Size{300, 500},
						MaxSize: declarative.Size{300, 500},
						OnCurrentIndexChanged: func() {
							if index := ct.tableview.CurrentIndex(); index > -1 {
								showCode(ct.tableModle.items[index].ModuleName)
							}

						},
					}, // TableView

					declarative.TextEdit{
						AssignTo: &ct.tb,
						ReadOnly: true,
						HScroll:  true,
						VScroll:  true,
					},
				},
			}, // HSplitter Children
		}, // MainWindow Children
	} // MainWindow

	// init end
}

func main() {
	mw.Run()
}

func selectFile() {
	fd := new(walk.FileDialog)
	fd.ShowOpen(ct.form)
	ct.lbFileName.SetText(fd.FilePath)
	// 清空textbox
	ct.tb.SetText("")
	// 清空table
	ct.tableModle.items = nil
	ct.tableModle.addItem([][2]string{})

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
	err := compdocFile.CFInit(cf)
	if err != nil {
		walk.MsgBox(ct.form, "title", err.Error(), walk.MsgBoxIconInformation)
		cfflag = false
	} else {
		showModule()
	}
	return
}

func showModule() {
	if cfflag {
		modules := cf.GetModuleName()
		//		ct.tb.SetText(strings.Join(modules, "\r\n"))
		ct.tableModle.addItem(modules)
	}
}
func showCode(moduleName string) {
	if cfflag {
		str := cf.GetModuleCode(moduleName)

		//		// 不能有NUL字符，会出错——可能原因：C的char数组是以\0来结尾的
		//		// 正常不会有的，除非是解压缩模块代码时候出问题了
		//		b := []byte(str)
		//		b = bytes.Replace(b, []byte{0}, []byte{}, -1)
		//		str = string(b)

		ct.tb.SetText(str)
	}
}
func unProtectProject() {
	if cfflag {
		newFile, err := cf.UnProtectProject()

		var str string
		if err != nil {
			str = err.Error()
		} else {
			str = "破解成功，新文件名：\r\n" + newFile
		}
		ct.tb.SetText(str)
	}
}

func (me *tableItemModle) addItem(modules [][2]string) {
	me.items = nil

	for _, v := range modules {
		item := &TableItem{}
		item.ModuleName = v[0]
		item.ModuleType = v[1]
		me.items = append(me.items, item)
	}
	me.PublishRowsReset()
}

func (me *tableItemModle) Items() interface{} {
	return me.items
}
