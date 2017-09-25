package main

import (
	"fmt"

	"github.com/skip2/go-qrcode"
)

func main() {
	err := qrcode.WriteFile("ceshi,hahah,123", qrcode.Medium, 256, "qr.png")
	fmt.Println(err)
}
