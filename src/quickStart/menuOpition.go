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
				AssignTo:  &ct.miTopmost,
				Checkable: true,
				Text:      "Topmost(&T)",
				OnTriggered: func() {
					ct.miTopmost.SetChecked(!ct.miTopmost.Checked())
					if ct.miTopmost.Checked() {
						win.SetWindowPos(ct.form.Handle(), win.HWND_TOPMOST, 0, 0, 0, 0, 1|2)
					} else {
						win.SetWindowPos(ct.form.Handle(), win.HWND_NOTOPMOST, 0, 0, 0, 0, 1|2)
					}
				},
			},
		}, // Items
	}
}
