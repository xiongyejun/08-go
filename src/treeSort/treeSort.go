package main

import (
	"fmt"
)

type tree struct {
	value       int
	left, right *tree
}

func main() {
	values := []int{1, 6, 4, 8, 2, 4, 7, 2}
	fmt.Println(values)
	Sort(values)
	fmt.Println(values)

}

func Sort(values []int) {
	var root *tree
	for _, v := range values {
		root = add(root, v) // add 构建了1个2叉树
	}
	appendValues(values[:0], root) // 前序遍历？，将tree结果输出到了初始为0的slice
}

func appendValues(values []int, t *tree) []int {
	if t != nil {
		values = appendValues(values, t.left)
		values = append(values, t.value)
		values = appendValues(values, t.right)
	}
	return values
}

func add(t *tree, value int) *tree {
	if t == nil {
		t = new(tree)
		t.value = value
		return t
	}

	if value < t.value {
		t.left = add(t.left, value)
	} else {
		t.right = add(t.right, value)
	}
	return t
}
