// 操作的菜单项

package main

import (
	"pkgMyPkg/inputBox"
	"strconv"

	"github.com/lxn/walk"
	"github.com/lxn/walk/declarative"
)

func menuActions() declarative.Menu {
	return declarative.Menu{
		Text: "操作(&C)",
		Items: []declarative.MenuItem{
			declarative.Action{
				AssignTo:    &ct.miUnProtectProject,
				Text:        "破解工程密码(&P)",
				OnTriggered: unProtectProject,
			},

			declarative.Action{
				AssignTo:    &ct.miHideModule,
				Text:        "隐藏模块(&H)",
				OnTriggered: hideModule,
			},
			declarative.Action{
				AssignTo:    &ct.miUnHideModule,
				Text:        "取消隐藏模块(&U)",
				OnTriggered: unHideModule,
			},
			declarative.Action{
				AssignTo:    &ct.miModifyFile,
				Text:        "改写文件(&M)",
				OnTriggered: modifyFile,
			},
		},
	} // "操作(&C)"
}

func modifyFile() {
	str := []string{"地址", "内容"}
	strRet := inputBox.Show(str, ct.form)
	if strRet[0] == "" || strRet[1] == "" {
		walk.MsgBox(ct.form, "出错了", "没有输入内容。", walk.MsgBoxIconWarning)
		return
	} else {
		var address int64
		var err error
		if address, err = strconv.ParseInt(strRet[0], 16, 32); err != nil {
			walk.MsgBox(ct.form, "出错了", err.Error(), walk.MsgBoxIconWarning)
			return
		}
		modifyB := []byte(strRet[1])
		if newFile, err := cf.ReWriteFile(int(address), modifyB); err != nil {
			walk.MsgBox(ct.form, "出错了", err.Error(), walk.MsgBoxIconWarning)
			return
		} else {
			walk.MsgBox(ct.form, "success", "修改成功，新文件:\r\n"+newFile, walk.MsgBoxIconInformation)
		}
	}
}
func unHideModule() {
	if cfflag {
		if index := ct.tableview.CurrentIndex(); index > -1 {
			if ct.tableModle.items[index].Type == "模块流" {
				newFile, err := cf.UnHideModule(ct.tableModle.items[index].Name)
				if err != nil {
					walk.MsgBox(ct.form, "出错了", err.Error(), walk.MsgBoxIconWarning)
				} else {
					walk.MsgBox(ct.form, "成功了", "取消隐藏成功，新文件:\r\n"+newFile, walk.MsgBoxIconInformation)
				}
			}
		}
	}
}

func hideModule() {
	if cfflag {
		if index := ct.tableview.CurrentIndex(); index > -1 {
			if ct.tableModle.items[index].Type == "模块流" {
				newFile, err := cf.HideModule(ct.tableModle.items[index].Name)
				if err != nil {
					walk.MsgBox(ct.form, "出错了", err.Error(), walk.MsgBoxIconWarning)
				} else {
					walk.MsgBox(ct.form, "成功了", "隐藏成功，新文件:\r\n"+newFile, walk.MsgBoxIconInformation)
				}
			}
		}
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
