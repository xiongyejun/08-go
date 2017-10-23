// 选项菜单

package main

import (
	//	"github.com/lxn/walk"
	"github.com/lxn/walk/declarative"
)

func menuCheck() declarative.Menu {
	return declarative.Menu{
		Text: "选项(&B)",
		Items: []declarative.MenuItem{
			declarative.Action{
				AssignTo:    &ct.miShowCode,
				Text:        "显示模块代码(&M)",
				Checkable:   true,
				OnTriggered: showCodeCheckChange,
			},
		},
	}
}

func showCodeCheckChange() {
	ct.miShowCode.SetChecked(!ct.miShowCode.Checked())
	tableviewShowInfoIndexChanged()
}
