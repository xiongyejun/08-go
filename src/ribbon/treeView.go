package main

import (
	"os/exec"

	"github.com/lxn/walk"
	"github.com/lxn/walk/declarative"
)

type TreeModle struct {
	walk.TreeModelBase
	roots []*node
}

var _ walk.TreeModel = new(TreeModle)

func initTreeView() declarative.TreeView {
	return declarative.TreeView{
		MaxSize:          declarative.Size{Width: 200},
		AssignTo:         &ct.treeView,
		Model:            ct.treeModle,
		Font:             declarative.Font{PointSize: 11},
		ContextMenuItems: contextMenu(),
		OnItemActivated: func() {
			//			if !ct.miEdit.Checked() {
			//				openFolderFile(ct.treeView.CurrentItem().(*node).Path)
			//			}
		},
	}
}

func newTreeModle() {
	ct.treeModle = new(TreeModle)
	//	if err := ct.treeModle.readNodeFromFile(); err != nil {
	//		fmt.Println(err)
	//	}

	return
}

// 连接到parent
func setParent(nd *node) {
	for i, _ := range nd.Children {
		setParent(nd.Children[i])
		nd.Children[i].parent = nd
		nd.Children[i].index = i
	}
}

// 使用cmd打开文件和文件夹
func openFolderFile(path string) error {
	// 第4个参数，是作为start的title，不加的话有空格的path是打不开的
	cmd := exec.Command("cmd.exe", "/c", "start", "", path)
	if err := cmd.Start(); err != nil {
		return err
	}
	return nil
}

// 实现walk TreeModle接口
func (*TreeModle) LazyPopulation() bool              { return false }
func (me *TreeModle) RootCount() int                 { return len(me.roots) }
func (me *TreeModle) RootAt(index int) walk.TreeItem { return me.roots[index] }
