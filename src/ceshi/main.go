package main

import (
	"encoding/binary"
	"fmt"
)

func main() {
	var i int64 = 23421323

	var b = make([]byte, 8)
	binary.PutVarint(b, i)
	fmt.Println(b)

	fmt.Println(binary.Varint(b))
}
