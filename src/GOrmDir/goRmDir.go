// 删除名称为xx的文件夹

package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

var rmDir string = ".git"

func main() {
	if len(os.Args) == 1 {
		fmt.Println("没有指定文件夹名称。")
		return
	}

	str, _ := os.Getwd() // 获得cmd命令行cd的路径
	scanDir(str)
}

func scanDir(strDir string) {
	entrys, _ := ioutil.ReadDir(strDir)

	for _, v := range entrys {
		if v.IsDir() {
			if v.Name() == rmDir {
				err := os.RemoveAll(strDir + string(os.PathSeparator) + v.Name())
				if err != nil {
					fmt.Println(err)
				} else {
					fmt.Println("delete dir:", strDir+string(os.PathSeparator)+v.Name())
				}

			} else {
				scanDir(strDir + string(os.PathSeparator) + v.Name())
			}
		}
	}
}
