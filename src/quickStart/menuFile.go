package main

import (
	"fmt"
	"time"

	"github.com/lxn/walk"
	"github.com/lxn/walk/declarative"
)

func fileMenu() declarative.Menu {
	fmt.Println("fileMenu")
	return declarative.Menu{
		Text: "文件(&F)",
		Items: []declarative.MenuItem{
			declarative.Action{
				AssignTo: &ct.miBackUp,
				Text:     "备份(&B)",
				OnTriggered: func() {
					fd := walk.FileDialog{Title: "输入文件名称.", FilePath: "bkjson" + time.Now().Format("2006年1月2日15时4分5秒") + ".txt", InitialDirPath: exePath, Filter: "所有文件|*.*"}
					if b, _ := fd.ShowSave(ct.form); b {
						ct.treeModle.saveNodeToFile(fd.FilePath)
					}
				},
			},

			declarative.Action{
				AssignTo: &ct.miQuit,
				Text:     "退出(&Q)",
				OnTriggered: func() {
					ct.form.Close()
				},
			},
		}, // Items
	}
}
