package main

import (
	"os"

	"fmt"
)

func main() {
	for k, str := range os.Environ() {
		fmt.Println(k, str)
	}

	fmt.Println(os.Getenv("USERPROFILE") + "\\Desktop")
}
