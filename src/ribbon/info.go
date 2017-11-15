// <commands>元素用来重复利用内置控件
// <contextualTabs>元素创建上下文选项卡
// box元素用来组织控件的排列
// separator元素用来放置分隔条
package main

import (
	"encoding/xml"
)

//	"strconv"

//const (
//	officeMenu = "officeMenu" // 用来定制Office菜单
//	qat        = "qat"        // 用来定制快速访问工具栏
//	tab        = "tab"        // 自定义功能区
//)

type ribbonType struct {
	Type   string
	TypeID int
}

var ribbonTypes []*ribbonType = []*ribbonType{
	{"tab", 1},
	{"group", 2},
	{"button", 3},
}

type customUI struct {
	XMLName   xml.Name `xml:"customUI"`
	Xmlns     string   `xml:"xmlns,attr"`
	XmlnsQ    string   `xml:"xmlns:Q,attr"`
	OnLoad    string   `xml:"onLoad,attr"`    // Sub onLoad(Ribbon as IRibbonUI)
	LoadImage string   `xml:"loadImage,attr"` // Sub loadImage(imageID as string,ByRef returnedVal)

	Ribbon *ribbon `xml:"ribbon"`
}

type ribbon struct {
	StartFromScratch string `xml:"startFromScratch,attr"` // startFromScratch属性能够隐藏整个内置的功能区。默认值为false

	Tabs *tabs `xml:"tabs"`
}

type tabs struct {
	TabSlice []*tab `xml:"tab"`
}

type tab struct {
	ID    string `xml:"id,attr"`
	Lable string `xml:"lable,attr"`

	GroupSlice []*group `xml:"group"`
}

type group struct {
	ID    string `xml:"id,attr"`
	Lable string `xml:"lable,attr"`

	ButtonSlice []*button `xml:"button"`
}

type button struct {
	// id	当创建自已的选项卡时
	// idMso	当使用现有的Microsoft选项卡时
	// idQ	当创建在命名空间之间共享的选项卡时
	ID    string `xml:"id,attr"`
	IdMso string `xml:"idMso,attr"`
	IdQ   string `xml:"idQ,attr"`
	// Sub OnAction(control As IRibbonControl)
	// 重新使用（或重利用）	Sub OnAction(control As IRibbonControl,byRef CancelDefaultcancelDefault)
	OnAction        string `xml:"onAction,attr"`
	InsertAfterMso  string `xml:"insertAfterMso,attr"`
	InsertBeforeMso string `xml:"insertBeforeMso,attr"`
	InsertAfterQ    string `xml:"insertAfterQ,attr"`
	InsertBeforeQ   string `xml:"insertBeforeQ,attr"`

	Description    string `xml:"description,attr"`
	Enabled        string `xml:"enabled,attr"`
	Image          string `xml:"image,attr"`
	ImageMso       string `xml:"imageMso,attr"`
	Keytip         string `xml:"keytip,attr"`
	Label          string `xml:"label,attr"`
	Screentip      string `xml:"screentip,attr"`
	ShowImage      string `xml:"showImage,attr"`
	ShowLabel      string `xml:"showLabel,attr"`
	Size           string `xml:"size,attr"`
	Supertip       string `xml:"supertip,attr"`
	Tag            string `xml:"tag,attr"`
	Visible        string `xml:"visible,attr"`
	GetDescription string `xml:"getDescription,attr"`
	GetEnabled     string `xml:"getEnabled,attr"`
	GetImage       string `xml:"getImage,attr"`
	GetKeytip      string `xml:"getKeytip,attr"`
	GetLabel       string `xml:"getLabel,attr"`
	GetScreentip   string `xml:"getScreentip,attr"`
	GetShowImage   string `xml:"getShowImage,attr"`
	GetShowLabel   string `xml:"getShowLabel,attr"`
	GetSize        string `xml:"getSize,attr"`
	GetSupertip    string `xml:"getSupertip,attr"`
	GetVisible     string `xml:"getVisible,attr"`
}
