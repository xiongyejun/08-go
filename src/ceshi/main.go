package main

import (
	"fmt"
)

func main() {
	fmt.Println(isNumber("145"))
}

func isNumber(str string) bool {
	for i := range str {
		fmt.Println(str[i])
	}
	return true
}
