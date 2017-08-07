package main

import (
	"fmt"
	"os"
)

func main() {
	s := os.Args[1]
	fmt.Println(s, "  ", comma(s))
}

func comma(s string) string {
	n := len(s)

	if n <= 3 {
		return s
	}

	return comma(s[:n-3]) + "," + s[n-3:]
}
