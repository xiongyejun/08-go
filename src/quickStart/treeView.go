package main

import (
	//	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
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
		AssignTo:         &ct.treeView,
		Model:            ct.treeModle,
		ContextMenuItems: contextMenu(),
		OnItemActivated: func() {
			nd := ct.treeView.CurrentItem().(*node)
			openFolderFile(nd.Path)
			// 打开文件后隐藏
			ct.form.Hide()
			// 双击会切换expend
			if len(nd.Children) > 0 {
				ct.treeView.SetExpanded(nd, !ct.treeView.Expanded(nd))
			}
		},

		OnKeyPress: func(key walk.Key) {
			if key == walk.KeyH {
				ct.form.Hide()
			}
		},
	}
}

func newTreeModle() {
	ct.treeModle = new(TreeModle)
	if err := ct.treeModle.readNodeFromFile(); err != nil {
		fmt.Println(err)
	}

	return
}

// 利用json读取保存的node信息
func (me *TreeModle) readNodeFromFile() (err error) {
	if b, err := ioutil.ReadFile(exePath + SAVE_FILE); err != nil {
		return err
	} else {
		if err := json.Unmarshal(b, &me.roots); err != nil {
			return err
		}
	}
	// 读取的node是没有连接到parent的
	for i, _ := range me.roots {
		me.roots[i].index = i
		setParent(me.roots[i])
	}

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

// 利用json保存node记录
func (me *TreeModle) saveNodeToFile(fileName string) (err error) {
	if js, err := json.Marshal(me.roots); err != nil {
		return err
	} else {
		if err = ioutil.WriteFile(fileName, js, 0666); err != nil {
			return err
		}
	}

	return nil
}

// 使用cmd打开文件和文件夹
func openFolderFile(path string) error {
	// 第4个参数，是作为start的title，不加的话有空格的path是打不开的
	cmd := exec.Command("cmd.exe", "/c", "start", "", path, "0")
	if err := cmd.Start(); err != nil {
		return err
	}
	return nil
}

// 实现walk TreeModle接口
func (*TreeModle) LazyPopulation() bool              { return false }
func (me *TreeModle) RootCount() int                 { return len(me.roots) }
func (me *TreeModle) RootAt(index int) walk.TreeItem { return me.roots[index] }
