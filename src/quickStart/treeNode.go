package main

import (
	"github.com/lxn/walk"
)

type node struct {
	Path string // 路径
	Name string // 显示的名称

	parent   *node
	Children []*node
	index    int // 记录在parent的Children里的下标、或是roots里的下标
}

var _ walk.TreeItem = new(node)

// 实现TreeItem接口
func (me *node) Text() string { return me.Name }
func (me *node) Parent() walk.TreeItem {
	// We can't simply return d.parent in this case, because the interface
	// value then would not be nil.
	if me.parent == nil {
		return nil
	}
	return me.parent
}

func (me *node) ChildCount() int {
	if me.Children == nil {
		return 0
	}
	return len(me.Children)
}

func (me *node) ChildAt(index int) walk.TreeItem { return me.Children[index] }

// 实现Imager接口，为了在treeview上显示图标
func (me *node) Image() interface{} { return me.Path }

//func (me *node) ResetChildren() error {
//	return nil
//}

func newNode(path, text string, parent *node) *node {
	n := &node{
		Path: path,
		Name: text,
		//		Level:  level,
		parent: parent,
		//		left:   left,
		//		right:  right,
	}
	return n
}
