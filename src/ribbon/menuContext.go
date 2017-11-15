package main

import (
	"fmt"

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
				if nd.Name == "" || nd.TypeID == 0 || nd == nil {
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
					expandChildren(currentNd, false)
				}
			},
		},

		//		declarative.Action{
		//			AssignTo: &ct.cmiAddRoot,
		//			Text:     "Add&Root",
		//			OnTriggered: func() {
		//				nd := new(node)
		//				showDialog(ct.form, nd)
		//				if nd.CustomType == "" || nd.Lable == "" || nd == nil {
		//					return
		//				}
		//				nd.index = len(ct.treeModle.roots)
		//				ct.treeModle.roots = append(ct.treeModle.roots, nd)
		//				ct.treeView.SetModel(ct.treeModle)
		//			},
		//		},

		//		declarative.Action{
		//			AssignTo: &ct.cmiDel,
		//			Text:     "&Delete",
		//			OnTriggered: func() {
		//				if currentNd, ok := ct.treeView.CurrentItem().(*node); ok {
		//					parentNd := currentNd.parent
		//					defer ct.treeModle.PublishItemsReset(parentNd)
		//					if parentNd == nil {
		//						ct.treeModle.roots = removeNode(ct.treeModle.roots, currentNd)
		//						ct.treeView.SetModel(ct.treeModle)
		//						return
		//					}
		//					parentNd.Children = removeNode(parentNd.Children, currentNd)
		//				}
		//			},
		//		},

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
	cb := new(walk.ComboBox)

	declarative.Dialog{
		AssignTo:      &dlg,
		Title:         "设置节点属性",
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
					declarative.Label{
						Text: "Type:",
					},
					declarative.ComboBox{
						AssignTo:      &cb,
						Value:         declarative.Bind("TypeID"),
						BindingMember: "Id",
						DisplayMember: "Type",
						Model:         ribbonTypes,
						OnCurrentIndexChanged: func() {
							fmt.Println(cb.Text())
							pb, _ := walk.NewPushButton(dlg)
							pb.SetText("ce")
						},
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
