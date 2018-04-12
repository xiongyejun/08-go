// 删除文件，根据文件大小

package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
)

var fileSize int64 = 0
var rmFiles []string

func main() {
	if len(os.Args) == 1 {
		fmt.Println("没有指定文件大小。")
		return
	}

	if tmp, err := strconv.Atoi(os.Args[1]); err != nil {
		fmt.Println(err)
		return
	} else {

		fileSize = int64(tmp)
	}

	str, _ := os.Getwd() // 获得cmd命令行cd的路径
	scanDir(str + string(os.PathSeparator))

	if len(rmFiles) > 0 {
		fmt.Printf("将要删除[%s]等%d个文件\r\n是否继续？\t[Y]继续\t[N]退出\r\n", rmFiles[0], len(rmFiles))
	}
	var tmpGoOn string
	fmt.Scanln(&tmpGoOn)
	if tmpGoOn == "Y" {
		fmt.Println("delete……")
		for i := range rmFiles {
			os.Remove(rmFiles[i])
		}
	}
}

func scanDir(strDir string) {
	entrys, _ := ioutil.ReadDir(strDir)

	for _, v := range entrys {
		if !v.IsDir() {
			if v.Size() <= fileSize {
				rmFiles = append(rmFiles, strDir+v.Name())
			}
		}
	}
}
