package main

import (
	"fmt"
	"path/filepath"
)

func main() {
	fmt.Println(filepath.Dir(`C:\Users\Administrator\Desktop\vbaProject.bin`))
	fmt.Println(filepath.Clean(`C:\Users\Administrator\Desktop\vbaProject.bin`))
}
