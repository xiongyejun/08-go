package main

import (
	"fmt"
	"path/filepath"
)

func main() {
	str := `C:\Users\Administrator\Desktop\111客户档案0403.xlsx`

	fmt.Println(str)
	fmt.Println(filepath.Base(str))
}
