// 读取customUI.xml的内容到ribbonTypes里
// 然后放到tree控件中

package main

import (
	"encoding/xml"
	"fmt"
)

func customUIToStruct(bXML []byte) (err error) {
	Custom_UI := new(rbcustomUI)

	if err = xml.Unmarshal(bXML, Custom_UI); err != nil {
		return
	}

	fmt.Println(Custom_UI.OnLoad)

	return
}
