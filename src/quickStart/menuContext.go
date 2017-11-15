package main

import (
	"fmt"
	"path/filepath"

	"github.com/lxn/walk"
	"github.com/lxn/walk/declarative"
)

func contextMenu() []declarative.MenuItem {
	return []declarative.MenuItem{
		declarative.Action{
			AssignTo: &ct.cmiAdd,
			Text:     "&Add",
			OnTriggered: func() {
				nd := new(node)
				showDialog(ct.form, nd)
				if nd.Path == "" || nd.Name == "" || nd == nil {
					return
				}
				if currentNd, ok := ct.treeView.CurrentItem().(*node); !ok {
					nd.index = len(ct.treeModle.roots)
					ct.treeModle.roots = append(ct.treeModle.roots, nd)
					ct.treeView.SetModel(ct.treeModle)
				} else {
					nd.parent = currentNd
					nd.index = len(currentNd.Children)
					currentNd.Children = append(currentNd.Children, nd)
					ct.treeModle.PublishItemsReset(currentNd)
				}
			},
		},

		declarative.Action{
			AssignTo: &ct.cmiMoveUp,
			Text:     "Move&Up",
			OnTriggered: func() {
				if currentNd, ok := ct.treeView.CurrentItem().(*node); ok {
					parentNd := currentNd.parent
					index := currentNd.index
					swapNode(parentNd, index, index-1)
				}
			},
		},

		declarative.Action{
			AssignTo: &ct.cmiMoveDown,
			Text:     "&MoveDown",
			OnTriggered: func() {
				if currentNd, ok := ct.treeView.CurrentItem().(*node); ok {
					parentNd := currentNd.parent
					index := currentNd.index
					swapNode(parentNd, index, index+1)
				}
			},
		},

		declarative.Action{
			AssignTo: &ct.cmiAddRoot,
			Text:     "Add&Root",
			OnTriggered: func() {
				nd := new(node)
				showDialog(ct.form, nd)
				if nd.Path == "" || nd.Name == "" || nd == nil {
					return
				}
				nd.index = len(ct.treeModle.roots)
				ct.treeModle.roots = append(ct.treeModle.roots, nd)
				ct.treeView.SetModel(ct.treeModle)
			},
		},

		declarative.Action{
			AssignTo: &ct.cmiDel,
			Text:     "&Delete",
			OnTriggered: func() {
				if currentNd, ok := ct.treeView.CurrentItem().(*node); ok {
					parentNd := currentNd.parent
					defer ct.treeModle.PublishItemsReset(parentNd)
					if parentNd == nil {
						ct.treeModle.roots = removeNode(ct.treeModle.roots, currentNd)
						ct.treeView.SetModel(ct.treeModle)
						return
					}
					parentNd.Children = removeNode(parentNd.Children, currentNd)
				}
			},
		},

		declarative.Action{
			AssignTo: &ct.cmiEdit,
			Text:     "&Edit",
			OnTriggered: func() {
				if currentNd, ok := ct.treeView.CurrentItem().(*node); ok {
					showDialog(ct.form, currentNd)
					defer ct.treeModle.PublishItemsReset(currentNd)
				}
			},
		},

		declarative.Action{
			AssignTo: &ct.cmiExpandAll,
			Text:     "展开所有(&E)",
			OnTriggered: func() {
				for i, _ := range ct.treeModle.roots {
					expandChildren(ct.treeModle.roots[i], true)
				}
			},
		},
	}
}

// 交换节点的位置
func swapNode(parentNd *node, i, j int) {
	if parentNd == nil {
		return
	}
	if i < 0 || j < 0 {
		walk.MsgBox(ct.form, "越界", "已到达最上面了。", walk.MsgBoxIconInformation)
		return
	}
	iLen := len(parentNd.Children)
	if i >= iLen || j >= iLen {
		walk.MsgBox(ct.form, "越界", "已到达最底下了。", walk.MsgBoxIconInformation)
		return
	}
	parentNd.Children[i], parentNd.Children[j] = parentNd.Children[j], parentNd.Children[i]
	parentNd.Children[i].index = i
	parentNd.Children[j].index = j

	ct.treeModle.PublishItemsReset(parentNd)
	expandChildren(parentNd, true)
	ct.treeView.SetCurrentItem(parentNd.Children[j])
}

// 展开节点
func expandChildren(nd *node, bDG bool) {
	ct.treeView.SetExpanded(nd, true)
	if bDG {
		for i, _ := range nd.Children {
			expandChildren(nd.Children[i], true)
		}
	}
}

// 在nodes中删除n，返回删除后的
func removeNode(nodes []*node, n *node) []*node {
	i := n.index
	for j := i; j < len(nodes)-1; j++ {
		nodes[j] = nodes[j+1]
		nodes[j].index = j
	}

	return nodes[:len(nodes)-1]
}

// 使用dialog来输入nd
func showDialog(owner walk.Form, nd *node) {
	dlg := new(walk.Dialog)
	btnDefault := new(walk.PushButton)
	btnCancel := new(walk.PushButton)
	db := new(walk.DataBinder)
	tbPath := new(walk.TextEdit)
	btnSelectFile := new(walk.PushButton)
	btnSelectFolder := new(walk.PushButton)
	initDir := ""
	if currentNd, ok := ct.treeView.CurrentItem().(*node); ok {
		initDir = filepath.Dir(currentNd.Path)
	}

	declarative.Dialog{
		AssignTo:      &dlg,
		Title:         "选择文件或文件夹",
		DefaultButton: &btnDefault,
		CancelButton:  &btnCancel,
		DataBinder: declarative.DataBinder{
			AssignTo:   &db,
			DataSource: nd,
		},
		MinSize: declarative.Size{300, 100},

		Layout: declarative.VBox{},
		Children: []declarative.Widget{
			declarative.Composite{
				Layout: declarative.Grid{Columns: 2},
				Children: []declarative.Widget{
					declarative.PushButton{
						AssignTo: &btnSelectFile,
						Text:     "选择文件",
						OnClicked: func() {
							fd := walk.FileDialog{InitialDirPath: initDir}
							if b, _ := fd.ShowOpen(dlg); b {
								tbPath.SetText(fd.FilePath)
							}
						},
					},

					declarative.PushButton{
						AssignTo: &btnSelectFolder,
						Text:     "选择文件夹",
						OnClicked: func() {
							fd := walk.FileDialog{InitialDirPath: initDir}
							if b, _ := fd.ShowBrowseFolder(dlg); b {
								tbPath.SetText(fd.FilePath)
							}
						},
					},
					declarative.Label{
						Text: "Path:",
					},
					declarative.TextEdit{
						AssignTo: &tbPath,
						Text:     declarative.Bind("Path"),
						ReadOnly: true,
					},

					declarative.Label{
						Text: "Name:",
					},
					declarative.TextEdit{
						Text: declarative.Bind("Name"),
					},
				},
			},
			declarative.HSpacer{},

			declarative.Composite{
				Layout: declarative.HBox{},
				Children: []declarative.Widget{
					declarative.PushButton{
						AssignTo: &btnDefault,
						Text:     "OK",
						OnClicked: func() {
							if err := db.Submit(); err != nil {
								fmt.Println(err)
								return
							}
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

	return
}
