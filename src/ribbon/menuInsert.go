package main

import (
	"github.com/lxn/walk/declarative"
)

const (
	START_customUI string = "<customUI xmlns=\"http://schemas.microsoft.com/office/2006/01/customui\">\r\n<ribbon>\r\n<tabs>\r\n"
	END_customUI   string = "</tabs>\r\n</ribbon>\r\n</customUI>"

	START_TAB string = "  <tab id=\"TabID\" label=\"tabName\" insertAfterMso=\"TabDeveloper\">\r\n"
	END_TAB   string = "  </tab>\r\n"

	START_GROUP string = "    <group id=\"GroupID\" label=\"GroupName\">\r\n\r\n"
	END_GROUP   string = "    </group>\r\n"
)

func insertMenu() declarative.Menu {
	return declarative.Menu{
		Text: "插入(&I)",
		Items: []declarative.MenuItem{
			declarative.Action{
				AssignTo: &ct.miCustomUI,
				Text:     "CustomUI(&C)",
				OnTriggered: func() {
					insertXml(START_customUI + START_TAB + START_GROUP + END_GROUP + END_TAB + END_customUI)
				},
			},

			declarative.Action{
				AssignTo: &ct.miButton,
				Text:     "Button(&B)",
				OnTriggered: func() {
					insertXml("      <button id=\"Button1\" label=\"buttonname&#13;\" size=\"large\" onAction=\"Macro\" imageMso=\"HappyFace\" />\r\n")
				},
			},
			declarative.Action{
				AssignTo: &ct.miCallBack,
				Text:     "CallBack(&A)",
				OnTriggered: func() {
					insertXml("Sub CallBack(control As IRibbonControl)\r\n\r\nEnd Sub")
				},
			},
		},
	}
}

func insertXml(strXml string) {
	ct.tbXml.ReplaceSelectedText(strXml, true)
}
