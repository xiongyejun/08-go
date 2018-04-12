package main

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
)

func main() {
	str := "19850522"

	sh1 := sha1.New()

	if _, err := sh1.Write([]byte(str)); err != nil {
		fmt.Println(err)
		return
	}
	b := sh1.Sum(nil)
	fmt.Println(hex.EncodeToString(b))
}
