package main

import (
	"io/ioutil"
	"pkgMyPkg/compoundFile"
)

type MyData struct {
	cf *compoundFile.CompoundFile
}

func (me *MyData) Parse(fileName string) (err error) {
	var b []byte
	if b, err = ioutil.ReadFile(fileName); err != nil {
		return
	}

	if me.cf, err = compoundFile.NewCompoundFile(b); err != nil {
		return
	}

	// 读取DataSpaceMap
	return nil
}

func (me *MyData) getDataSpaceMap() (err error) {

	return nil
}
