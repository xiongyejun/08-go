package main

import (
	"github.com/lxn/walk"
)

func initNotifyIcon() error {
	var err1 error
	ct.ni, err1 = walk.NewNotifyIcon()
	if err1 != nil {
		return err1
	}

	ic, err := walk.Resources.Icon("\\image\\go.ico")
	if err != nil {
		return err
	}
	if err := ct.ni.SetIcon(ic); err != nil {
		return err
	}
	if err := ct.ni.SetToolTip("go\r\nhere"); err != nil {
		return err
	}

	ct.ni.MouseDown().Attach(func(x, y int, button walk.MouseButton) {
		if button == walk.LeftButton {
			setShowToRight()
		}
	})

	if err := ct.ni.SetVisible(true); err != nil {
		return err
	}

	if err := ct.ni.ShowInfo("NotifyIcon", "Click"); err != nil {
		return err
	}

	return nil
}
