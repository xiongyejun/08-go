package main

import (
	"fmt"
)

var count int = 0

func main() {
	fmt.Println(getCount(10, 4))
}

func getCount(m, n int) (icount int) {
	icount = 1
	for i := 0; i < n; i++ {
		icount *= m
	}
	return
}
