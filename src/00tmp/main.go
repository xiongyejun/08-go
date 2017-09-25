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

	var cf compdocFile.CF
	if compdocFile.IsCompdocFile(fileName) {
		cf = compdocFile.NewXlsFile(fileName)
	} else if compdocFile.IsZip(fileName) {
		cf = compdocFile.NewZipFile(fileName)
	}

	err := compdocFile.CFInit(cf)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("size", cf.GetFileSize())

	fmt.Println(cf.GetModuleString("ThisWorkbook"))
	fmt.Println(cf.GetModuleString("Sheet1"))
	fmt.Println(cf.GetModuleString("CCompdocFile"))

}
