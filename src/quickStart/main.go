// 快速启动窗体
// 用树形控件记录

package main

import (
	"fmt"
	"os"
)

const SAVE_FILE string = "json.txt"

var exePath string

func main() {
	exePath, _ = os.Getwd()
	exePath += string(os.PathSeparator)
	//	fmt.Println(exePath)

	//	return
	uiInit()
	// 关闭时记录节点
	if err := ct.treeModle.saveNodeToFile(exePath + SAVE_FILE); err != nil {
		fmt.Println(err)
	}

	fmt.Println("ok")
}
