package main

import (
	"regexp"

	"fmt"
)

func main() {
	b := []byte("abcdefg")
	m, err := regexp.Match(`cd`, b)
	fmt.Println(m, err)
}
