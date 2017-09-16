package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"pkgMySelf/ucs2T0utf8"
)

func main() {
	f, _ := os.Open("ucs2.txt")
	defer f.Close()
	b, _ := ioutil.ReadAll(f)
	b = b[2:]
	//	fmt.Printf("%s\n", b)
	b, err := ucs2T0utf8.UCS2toUTF8(b)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("%s\n", b)

}
