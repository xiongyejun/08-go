package main

import (
	"fmt"
)

func main() {
	if in, err := Parse(`C:\Users\Administrator\Desktop\加密\agile.xlsm`); err != nil {
		fmt.Println(err)
		return
	} else {
		if err := in.CheckPassword("12"); err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("密码正确。")
		}
	}

}
