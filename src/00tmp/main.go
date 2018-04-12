package main

import (
	"fmt"
)

func main() {
	str := "abc,hao,de"
	var key []byte = []byte("19850522")
	fmt.Printf("src = %s\r\n", str)

	c, err := desEncryptString(str, key)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("encrypt = %s\r\n", c)

	d, err := desDecryptString(c, key)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf(" desDecrypt = %s\r\n", d)
}
