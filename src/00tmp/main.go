package main

import (
	"fmt"
	"os"
)

func main() {
	f, _ := os.OpenFile("msg.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	f.WriteString(fmt.Sprintln("a<'", "msg"))
	f.Close()

}
