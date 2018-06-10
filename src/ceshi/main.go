package main

import (
	"encoding/base64"
	"fmt"
)

func main() {
	str := "f7u+rzRljbybSV+jwnCDzg=="
	if b, err := base64.StdEncoding.DecodeString(str); err != nil {
		fmt.Println(b)
	} else {
		fmt.Println(b)
		fmt.Printf("%s\r\n", b)
	}

}
