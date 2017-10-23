package main

import (
	//	"fmt"

	"github.com/lxn/walk"
	"github.com/lxn/walk/declarative"
)

func actionMenu() declarative.Menu {
	return declarative.Menu{
		Text: "操作(&C)",
		Items: []declarative.MenuItem{

			declarative.Action{
				AssignTo:  &ct.miCompress,
				Text:      "压缩(&C)",
				Checkable: true,
				OnTriggered: func() {
					if str, err := mt.compress(); err != nil {
						walk.MsgBox(ct.form, "err", err.Error(), walk.MsgBoxIconWarning)
						return
					} else {
						ct.imageName = savePic(str)
						ct.pImage = 0
						showImage()
					}
				},
			},
			declarative.Action{
				AssignTo: &ct.miUnCompress,
				Text:     "解压缩(&U)",
				OnTriggered: func() {
					if err := mt.unCompress(); err != nil {
						walk.MsgBox(ct.form, "err", err.Error(), walk.MsgBoxIconWarning)
					}
				},
			},
		},
	} // "操作(&C)"
}
