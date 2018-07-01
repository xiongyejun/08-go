// Compound File Encryption

package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Println("CompoundFileEncryption <FileName> <Password>")
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
