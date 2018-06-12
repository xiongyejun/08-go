package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Println("第1个参数filename，第2个参数Password")
		return
	}

	if in, err := Parse(os.Args[1]); err != nil {
		fmt.Println(err)
		return
	} else {
		if err := in.CheckPassword(os.Args[2]); err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("密码正确。")
		}
	}

}
