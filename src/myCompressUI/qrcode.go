package main

import (
	"fmt"
	"strconv"

	qrcode "github.com/skip2/go-qrcode"
)

func savePic(str string) []string {
	n := len(str)
	var picname []string

	for i := 0; i < n; i += 2048 {
		end := i + 2048
		if end > n {
			end = n
		}
		savename := strconv.Itoa(i/2048) + "qrcode.png"
		picname = append(picname, savename)
		if err := qrcode.WriteFile(str[i:end], qrcode.Low, 512, savename); err != nil {
			fmt.Println(err)
		}
		fmt.Println(i / 2048)
	}
	return picname
}
