package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// 记录某种后缀的文件个数
var dic map[string]int = make(map[string]int)

func main() {
	dir, _ := os.Getwd()
	if err := scanDir(dir); err != nil {
		fmt.Println(err)
	}
}

func scanDir(strDir string) (err error) {
	return filepath.Walk(strDir, func(path string, f os.FileInfo, err error) error {
		if f == nil {
			return err
		}
		if f.IsDir() {
			return nil
		}

		if f.Name() == ".DS_Store" || f.Name() == "Thums.db" {
			return nil
		}
		ext := strings.ToLower(filepath.Ext(f.Name()))
		dic[ext] += 1
		newpath := strconv.Itoa(dic[ext]) + ext
		if err := os.Rename(f.Name(), newpath); err != nil {
			return err
		} else {
			fmt.Println(f.Name(), "—>", newpath)
		}

		return nil
	})
}
