package main

import (
	"fmt"
	"runtime"
)

func main() {
	fmt.Println(formatTime(234))
	fmt.Println(formatTime(2341))
	fmt.Println(formatTime(23411))
	fmt.Println(formatTime(234111))
	fmt.Println(formatTime(234111))

	fmt.Println(runtime.NumCPU())
}

func formatTime(second int64) string {
	m := second / 60
	second = second % 60

	h := m / 60
	m = m % 60

	return fmt.Sprintf("%2d时%2d分%2d秒", h, m, second)
}
