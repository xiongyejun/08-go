package main

import (
	"fmt"
	"io/ioutil"
	"pkgMySelf/rleVBA"
)

func main() {
	b, _ := ioutil.ReadFile("m")
	rle := rleVBA.NewRLE(b[0x398:]) // 16*57+8) 0x398
	b = rle.UnCompress()
	ioutil.WriteFile("mm", b, 0666)
	fmt.Println("ok")
}
