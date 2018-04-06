package main

import (
	"fmt"
	"io/ioutil"

	"pkgMyPkg/compoundFile"
)

func main() {

	var cf *compoundFile.CompoundFile
	var err error
	var b []byte
	if b, err = ioutil.ReadFile(`C:\Users\Administrator\Desktop\加密\密码1.xlsm`); err != nil {
		fmt.Println(err)
		return
	}

	if cf, err = compoundFile.NewCompoundFile(b); err != nil {
		fmt.Println(err)
		return
	}

	if err = cf.Parse(); err != nil {
		fmt.Println(err)
		return
	}

	b, err = cf.GetStream(string([]byte{6}) + `DataSpaces\DataSpaceMap`)
	fmt.Println(len(b), err)

	cf.PrintOut()
}
