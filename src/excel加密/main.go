package main

import (
	"fmt"
	"io/ioutil"

	"pkgMyPkg/compoundFile"
)

func main() {
	b, _ := ioutil.ReadFile(`C:\Users\Administrator\Desktop\加密\nn.xlsm`)
	cf, err := compoundFile.NewCompoundFile(b)

	err = cf.Parse()
	fmt.Println(err)

	b, err = cf.GetStream("EncryptionInfo")
	fmt.Println(len(b), err)
}
