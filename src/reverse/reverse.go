package main

import (
	"fmt"
)

func main() {
	a := []int{1, 2, 3, 4, 5, 5, 4, 7, 34, 74, 2}
	fmt.Println(a)
	reverse(a[:])
	//	reverse(a[:2])
	//	reverse(a[2:])
	fmt.Println(a)
}

func reverse(s []int) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}
