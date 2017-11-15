package main

import (
	//	"encoding/xml"
	"fmt"
	//	"io/ioutil"
)

var custom_UI = new(customUI)

func main() {
	uiInit()
	//	saveNodeToXml()

	fmt.Println("ok")
}

//func saveNodeToXml() {
//	if b, err := xml.Marshal(ct.treeModle.roots[0]); err != nil {
//		fmt.Println(err)
//	} else {
//		if err := ioutil.WriteFile("xml.txt", b, 0666); err != nil {
//			fmt.Println(err)
//		}
//	}
//}
