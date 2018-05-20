package main

import (
	"fmt"
)

func main() {
	d := new(MyData)

	if err := d.Parse(`C:\Users\Administrator\Desktop\加密\密码1.xlsm`); err != nil {
		fmt.Println(err)
		return
	}

}
