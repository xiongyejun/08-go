package main

import (
	"archive/zip"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/lxn/walk"
	"github.com/lxn/walk/declarative"
)

func fileMenu() declarative.Menu {
	fmt.Println("fileMenu")
	return declarative.Menu{
		Text: "文件(&F)",
		Items: []declarative.MenuItem{
			declarative.Action{
				AssignTo:    &ct.miSelectFile,
				Text:        "选择文件(&S)",
				OnTriggered: selectFile,
			},

			declarative.Action{
				AssignTo: &ct.miSaveXml,
				Text:     "保存xml(&S)",
				OnTriggered: func() {
					if ct.fileName == "" {
						walk.MsgBox(ct.form, "", "还没选择文件。", walk.MsgBoxIconInformation)
						return
					}
					if newFile, err := saveXmlToFile(); err != nil {
						walk.MsgBox(ct.form, "保存出错。", err.Error(), walk.MsgBoxIconError)
					} else {
						if ct.miBak.Checked() {
							if err := os.Rename(ct.fileName, ct.fileName+".bak"); err != nil {
								walk.MsgBox(ct.form, "备份文件出错。", err.Error(), walk.MsgBoxIconError)
								return
							}
						}
						if err := os.Rename(newFile, ct.fileName); err != nil {
							walk.MsgBox(ct.form, "覆盖文件出错。", err.Error(), walk.MsgBoxIconError)
						}
					}
				},
			},
			declarative.Action{
				AssignTo: &ct.miOpenFile,
				Text:     "打开文件(&O)",
				OnTriggered: func() {
					if ct.fileName != "" {
						openFolderFile(ct.fileName)
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

func saveXmlToFile() (newFile string, err error) {
	var zipReader *zip.ReadCloser
	var wr io.Writer

	// 读取zip文件
	if zipReader, err = zip.OpenReader(ct.fileName); err != nil {
		return
	}
	defer zipReader.Close()

	// 设置新文件名
	strExt := filepath.Ext(ct.fileName)
	newFile = ct.fileName[:len(ct.fileName)-len(strExt)] + "(new)" + strExt
	// 创建新文件
	var fw *os.File
	fw, err = os.OpenFile(newFile, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return
	}
	defer fw.Close()
	// 创建zip writer
	zipWriter := zip.NewWriter(fw)
	defer zipWriter.Close()
	// 循环zip文件中的文件
	for _, f := range zipReader.File {
		// 打开子文件
		var fr io.ReadCloser
		if fr, err = f.Open(); err != nil {
			return
		}
		defer fr.Close()

		// 在zipwriter中创建新文件
		if wr, err = zipWriter.Create(f.Name); err != nil {
			return
		}

		// 如果是custom，就用改写了的流
		if f.Name == ct.customUIName {
			if _, err = io.Copy(wr, strings.NewReader(ct.tbXml.Text())); err != nil {
				return
			}
		} else {
			if f.Name == "_rels/.rels" && ct.customUIName == "" {
				// 在最后添加关联
				var b []byte
				if b, err = ioutil.ReadAll(fr); err != nil {
					return
				} else {
					tmp := string(b[:len(b)-len("</Relationships>")]) + `<Relationship Id="sd23s" Type="http://schemas.microsoft.com/office/2006/relationships/ui/extensibility" Target="customUI/customUI.xml"/></Relationships>`
					if _, err = io.Copy(wr, strings.NewReader(tmp)); err != nil {
						return
					}
				}

			} else {
				if _, err = io.Copy(wr, fr); err != nil {
					return
				}
			}
		}

	}
	// 如果没有的话，需要创建1个
	if ct.customUIName == "" {
		// 在zipwriter中创建新文件
		if wr, err = zipWriter.Create("customUI/customUI.xml"); err != nil {
			return
		}
		if _, err = io.Copy(wr, strings.NewReader(ct.tbXml.Text())); err != nil {
			return
		}
		// 然后记录下来，后面再修改就不会再添加了
		ct.customUIName = "customUI/customUI.xml"
	}
	return
}

func selectFile() {
	fd := walk.FileDialog{}
	if b, _ := fd.ShowOpen(ct.form); b {
		ct.fileName = fd.FilePath
		// 读取zip文件中的customUI
		if zipReader, err := zip.OpenReader(ct.fileName); err != nil {
			ct.tbXml.SetText(err.Error())
		} else {
			defer zipReader.Close()
			for _, f := range zipReader.File {
				if f.Name == "customUI/customUI.xml" || f.Name == "customUI/customUI14.xml" {
					ct.customUIName = f.Name

					if frc, err := f.Open(); err != nil {
						ct.tbXml.SetText(err.Error())
					} else {
						if bstr, err := ioutil.ReadAll(frc); err != nil {
							ct.tbXml.SetText(err.Error())
						} else {
							// 解析到ribbon_Xml结构里
							customUIToStruct(bstr)
							ct.tbXml.SetText(string(bstr))
						}
					}

					return
				}
			}
			ct.customUIName = ""
			ct.tbXml.SetText("没有找到customUI.xml")
		}
	} else {
		ct.fileName = ""
	}
}

func customUIToStruct(b []byte) (err error) {
	fmt.Println(123)
	if err = xml.Unmarshal(b, custom_UI); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("%#v\r\n", *custom_UI)
	fmt.Println(custom_UI.Xmlns)
	fmt.Println(custom_UI.XmlnsQ)
	fmt.Printf("%#v\r\n", custom_UI.Ribbon.Tabs.TabSlice[0].GroupSlice[0])

	return
}
