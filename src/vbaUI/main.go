// 用UI来做VBA模块的查看，破解工程密码等
package main

import (
	"fmt"
	"pkgMyPkg/compdocFile"
	"strings"
	"unicode"

	"github.com/lxn/walk"
	"github.com/lxn/walk/declarative"
	"github.com/lxn/win"
)

type control struct {
	form *walk.MainWindow

	// MenuItem
	miSelectFile *walk.Action
	miExit       *walk.Action

	miUnProtectProject         *walk.Action // 破解工程密码
	miUnProtectSheetProtection *walk.Action // 破解工作表保护密码
	miHideModule               *walk.Action // 隐藏模块
	miUnHideModule             *walk.Action // 取消隐藏模块
	miModifyFile               *walk.Action // 根据输入的地址和内容改写文件
	miShowCode                 *walk.Action // 是否需要解压缩模块流，显示模块的代码

	lbFileName *walk.Label

	hsplitter  *walk.Splitter
	tableModle *tableItemModle // tableview添加items
	tableview  *walk.TableView // tableview显示信息
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
		Title:    "vba查看",
		//		Size:     declarative.Size{800, 800},
		Font: declarative.Font{PointSize: 10},
		Icon: getExcelIcon(),
		// 菜单
		MenuItems: []declarative.MenuItem{
			menuFile(), // "文件(&F)
			menuActions(),
			menuCheck(),
		},

		// 布局
		Layout: declarative.VBox{},
		// 控件
		Children: []declarative.Widget{ // widget小部件
			declarative.Label{
				AssignTo: &ct.lbFileName,
				Text:     "文件名称",
			},

			declarative.HSplitter{
				AssignTo: &ct.hsplitter,
				Children: []declarative.Widget{
					tableviewShowInfo(), // TableView

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
	mw.Create()
	ct.miShowCode.SetChecked(true)
	win.ShowWindow(ct.form.Handle(), win.SW_MAXIMIZE)
	ct.form.Run()

}

func showVBAInfo() {
	if cfflag {
		info := cf.GetVBAInfo()
		// 添加到tableview
		ct.tableModle.addItem(info)
	}
}

// 如果是其他流就显示字节
// 返回ASCII的string，方便查找attribute
// 如果不单独返回，就可能被换行了
func showByte(name string) string {
	b, address, addStep32 := cf.GetStream(name)

	n := len(b)
	var str []string = make([]string, 0, n/16+2)
	var strASCII []string = make([]string, 0, n/16+2)
	str = append(str, fmt.Sprintf("   index Address  % X ------ASCII-----", []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}))

	addStep := int(addStep32)
	// bAddress是addStep个记录1个开始地址
	for i := 0; i < n/16*16; i += 16 {
		pAddress := int(i / addStep)
		str_tmp := fmt.Sprintf("%08X %08X % X ", i/16, address[pAddress]+int32(i%addStep), b[i:i+16])
		strASCII_tmp := ""
		for _, v := range b[i : i+16] {
			if v < 127 && unicode.IsPrint(rune(v)) {
				str_tmp += fmt.Sprintf("%c", v)
				strASCII_tmp += fmt.Sprintf("%c", v)
			} else {
				str_tmp += "^"
				strASCII_tmp += "^"
			}
		}
		str = append(str, str_tmp)
		strASCII = append(strASCII, strASCII_tmp)
	}
	// 最后可能还剩下0-15个
	nleft := n % 16

	if nleft > 0 {
		var strChar string
		var str_tmp string = fmt.Sprintf("%08X %08X ", n/16+1, address[len(address)-1]+int32(n/16*16%addStep))

		strASCII_tmp := ""
		for i := n / 16 * 16; i < n; i++ {
			str_tmp += fmt.Sprintf("%02X ", b[i])
			if unicode.IsPrint(rune(b[i])) {
				strChar += fmt.Sprintf("%c", b[i])
				strASCII_tmp += fmt.Sprintf("%c", b[i])
			} else {
				strChar += "^"
				strASCII_tmp += "^"
			}
		}

		str_tmp = str_tmp + strings.Repeat(" ", 3*(16-nleft))
		str_tmp = str_tmp + strChar
		str = append(str, str_tmp)
		strASCII = append(strASCII, strASCII_tmp)
	}

	ct.tb.SetText(strings.Join(str, "\r\n"))
	return strings.Join(strASCII, "")
}

// 如果是模块流就显示模块代码
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

func getExcelIcon() *walk.Icon {
	ic, err := walk.NewIconFromFile("E:\\04-github\\98-pic\\ico\\Excel32px.ico")
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return ic
}
