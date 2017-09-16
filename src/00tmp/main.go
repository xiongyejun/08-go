package main

import (
	"fmt"
	"pkgAPI/comdlg32"
	"pkgMySelf/compdocFile"
)

func main() {
	fd := comdlg32.NewFileDialog()
	b, _ := fd.GetOpenFileName()
	if !b {
		return
	}

	fileName := fd.FilePath
	fmt.Println(fileName)
	if compdocFile.CheckCompdocFile(fileName) {
		cf := compdocFile.NewXlsFile(fileName)
		compdocFile.CFInit(cf)

		fmt.Println("size", cf.GetFileSize())
	}

}
