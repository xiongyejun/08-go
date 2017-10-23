package main

import (
	"fmt"
	"strconv"

	qrcode "github.com/skip2/go-qrcode"
)

func savePic(str string) {
	n := len(str)

	for i := 0; i < n; i += 2048 {
		end := i + 2048
		if end > n {
			end = n
		}
		if err := qrcode.WriteFile(str[i:end], qrcode.Low, 512, strconv.Itoa(i/2048)+"qrcode.png"); err != nil {
			fmt.Println(err)
		}
	}
}
