package main

import (
	"fmt"
	"pkgAPI/comdlg32"
)

func main() {
	opf := comdlg32.NewFileDialog()

	//	opf.GetOpenFileNames()

	//	fmt.Println(opf)
	//	fmt.Println(opf.FilePath)
	//	for _, r := range opf.FilePaths {
	//		fmt.Println(r)
	//	}

	opf.GetSaveFileName()
	fmt.Println(opf.FilePath)

}
