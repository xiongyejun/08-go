package main

import (
	"fmt"
	"io/ioutil"
	"regexp"

	"github.com/axgle/mahonia"
)

func main() {
	moduleName := charSetChange("utf-8", "gbk", "模块1")
	fmt.Println(moduleName)

	b, _ := ioutil.ReadFile("PROJECT")

	pattern := `Module=` + moduleName + `\r\n|` + moduleName + `.*?\r\n`
	reg, _ := regexp.Compile(pattern)

	b1 := reg.ReplaceAll(b, []byte{})

	ioutil.WriteFile("PROJECTre", b1, 0666)
	fmt.Println("ok")
}
func charSetChange(srcCharset string, desCharset string, src string) string {
	srcCoder := mahonia.NewEncoder(desCharset)
	srcResult := srcCoder.ConvertString(src)
	return srcResult

}
