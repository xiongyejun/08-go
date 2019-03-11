package main

import (
	"crypto/sha256"
	"fmt"
	"os"
)

//func main() {
//	sha := sha1.New()

//	if _, err := sha.Write([]byte("a1ee74d59d5c2eaf9e54df0dd8f2b5815f8dd62bbc1e240ddc57e13c33235a5b")); err != nil {
//		fmt.Println(err)
//		return
//	}

//	b := sha.Sum(nil)
//	fmt.Println(hex.EncodeToString(b))

//}

func main() {
	if len(os.Args) == 1 {
		fmt.Println("sha256 <value>")
		return
	}
	sha := sha256.New()
	if _, err := sha.Write([]byte(os.Args[1])); err != nil {
		fmt.Println(err)
		return
	} else {
		b := sha.Sum(nil)
		fmt.Printf("%x\r\n", b)
	}
}
