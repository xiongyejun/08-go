package main

import (
	"github.com/lxn/walk/declarative"
	"github.com/lxn/win"
)

func optionMenu() declarative.Menu {
	return declarative.Menu{
		Text: "选项(&X)",
		Items: []declarative.MenuItem{
			declarative.Action{
				AssignTo:  &ct.miBak,
				Checkable: true,
				Text:      "备份源文件(&B)",
				OnTriggered: func() {
					ct.miBak.SetChecked(!ct.miBak.Checked())
				},
			},
			declarative.Action{
				AssignTo:  &ct.miTopMost,
				Checkable: true,
				Text:      "窗口置顶(&T)",
				OnTriggered: func() {
					ct.miTopMost.SetChecked(!ct.miTopMost.Checked())

					if ct.miTopMost.Checked() {
						win.SetWindowPos(ct.form.Handle(), win.HWND_TOPMOST, 0, 0, 0, 0, 1|2)
					} else {
						win.SetWindowPos(ct.form.Handle(), win.HWND_NOTOPMOST, 0, 0, 0, 0, 1|2)
					}
				},
			},
		}, // Items
	}
}
