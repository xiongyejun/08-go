// 选择icons

package main

import (
	"github.com/lxn/walk"
	"github.com/lxn/walk/declarative"
)

func showIcon(owner walk.Form, imageMso *string) {
	dlg := new(walk.Dialog)
	btnDefault := new(walk.PushButton)
	btnCancel := new(walk.PushButton)
	iv := new(walk.ImageView)
	im, _ := walk.NewImageFromFile("C:\\Users\\Administrator\\Desktop\\16.bmp")

	declarative.Dialog{
		AssignTo:      &dlg,
		Title:         "选择icon",
		DefaultButton: &btnDefault,
		CancelButton:  &btnCancel,
		MinSize:       declarative.Size{300, 300},

		Layout: declarative.VBox{},
		Children: []declarative.Widget{

			declarative.ImageView{
				AssignTo: &iv,
				Image:    im,
			},

			declarative.HSpacer{},

			declarative.Composite{
				Layout: declarative.HBox{},
				Children: []declarative.Widget{
					declarative.PushButton{
						AssignTo: &btnDefault,
						Text:     "OK",
						OnClicked: func() {
							*imageMso = "1"
							dlg.Accept()
						},
					},
					declarative.PushButton{
						AssignTo:  &btnCancel,
						Text:      "Cancel",
						OnClicked: func() { dlg.Cancel() },
					},
				},
			},
		},
	}.Run(owner)

	iv.MouseUp().Attach(ImageViewMouseUp)

	return
}

func ImageViewMouseUp(x, y int, button walk.MouseButton) {

}
