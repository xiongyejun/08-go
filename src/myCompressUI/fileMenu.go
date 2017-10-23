package main

import (
	"io/ioutil"

	"github.com/lxn/walk"
	"github.com/lxn/walk/declarative"
)

func fileMenu() declarative.Menu {
	return declarative.Menu{
		Text: "文件(&F)",
		Items: []declarative.MenuItem{
			declarative.Action{
				AssignTo:    &ct.miSelectFile,
				Text:        "&选择文件",
				OnTriggered: selectFile, // 触发，相当于click
			},
			declarative.Action{
				AssignTo: &ct.miSelectFile,
				Text:     "&退出",
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
	mt.fileName = fd.FilePath
	ct.lbFileName.SetText(mt.fileName)

	if mt.fileName == "" {
		return
	}
	var err error
	if mt.fileByte, err = ioutil.ReadFile(mt.fileName); err != nil {
		walk.MsgBox(ct.form, "err", err.Error(), walk.MsgBoxIconWarning)
	}
	if mt.fileByte[0] == 0xef && mt.fileByte[1] == 0xbb && mt.fileByte[2] == 0xbf {
		mt.fileByte = mt.fileByte[3:] // 跳过UTF-8的头
	}
	return
}
