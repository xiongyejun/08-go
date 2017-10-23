package main

import (
	"pkgMyPkg/compdocFile"

	"github.com/lxn/walk"
	"github.com/lxn/walk/declarative"
)

func menuFile() declarative.Menu {
	return declarative.Menu{
		Text: "文件(&F)",
		Items: []declarative.MenuItem{
			declarative.Action{
				AssignTo:    &ct.miSelectFile,
				Text:        "选择文件(&S)",
				OnTriggered: selectFile, // 触发，相当于click
			},

			declarative.Action{
				AssignTo: &ct.miSelectFile,
				Text:     "退出(&Q)",
				OnTriggered: func() {
					ct.form.Close()
				},
			},
		}, // Items
	}
}

func selectFile() {
	fd := new(walk.FileDialog)
	fd.ShowOpen(ct.form)
	ct.lbFileName.SetText(fd.FilePath)
	// 清空textbox
	ct.tb.SetText("")
	// 清空table
	ct.tableModle.items = nil
	ct.tableModle.addItem([]*compdocFile.OutStruct{})

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
		showVBAInfo()

	}

	return
}
