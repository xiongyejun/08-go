package main

import (
	"fmt"
	"pkgMyPkg/colorPrint"
)

func main() {
	colorPrint.SetColor(colorPrint.White, colorPrint.DarkCyan)

	fmt.Println("test")

	colorPrint.UnSetColor()
}
