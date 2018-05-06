// 导出文件的vba模块代码

package main

import (
	"archive/zip"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"pkgMyPkg/compoundFile"
	"pkgMyPkg/rleVBA"
	"pkgMyPkg/vbaDir"
)

type dataStruct struct {
	fileName string
	cf       *compoundFile.CompoundFile
	saveDir  string
}

var d *dataStruct = new(dataStruct)

func main() {
	if len(os.Args) == 1 {
		fmt.Println("请输入文件。")
		return
	}
	d.fileName = os.Args[1]
	if _, err := os.Stat(d.fileName); err != nil {
		fmt.Println(err)
		return
	}
	d.saveDir = d.fileName + "Codes" + string(os.PathSeparator)

	if err := d.getVBACode(); err != nil {
		fmt.Println(err)
		return
	}
}

func (me *dataStruct) getVBACode() (err error) {
	// 判断一下是否是zip文件
	var b []byte
	if IsZip(me.fileName) {
		if b, err = readVbaProjectBin(me.fileName); err != nil {
			return
		}
	} else {
		if b, err = ioutil.ReadFile(me.fileName); err != nil {
			return
		}
	}
	if me.cf, err = compoundFile.NewCompoundFile(b); err != nil {
		return
	}

	if err = me.cf.Parse(); err != nil {
		return
	}

	// 读取dir目录
	var bDir []byte
	if bDir, err = me.cf.GetStream(`VBA\dir`); err != nil {
		return
	}
	// 读取的dir byte需压进行解压
	rle := rleVBA.NewRLE(bDir)
	bDir = rle.UnCompress()
	// 并分析有哪些模块
	if mi, err1 := vbaDir.GetModuleInfo(bDir); err1 != nil {
		return err1
	} else {
		if err = os.Mkdir(me.saveDir, 0666); err != nil {
			return
		}
		// 保存文件
		for i := range mi {
			fmt.Println(mi[i].Name)
			if b, err = me.cf.GetStream(`VBA\` + mi[i].Name); err != nil {
				return
			} else {
				rle = rleVBA.NewRLE(b[mi[i].TextOffset:])
				b = rle.UnCompress()

				var strExt string
				if mi[i].ModuleType == 33 {
					// 标准模块
					strExt = ".bas"
				} else {
					strExt = ".cls"
				}

				if err = ioutil.WriteFile(me.saveDir+mi[i].Name+strExt, b, 0666); err != nil {
					return
				}
			}
		}
	}
	return nil
}

func IsZip(fileName string) bool {
	b := make([]byte, 2)
	f, _ := os.Open(fileName)
	defer f.Close()
	f.Read(b)
	if b[0] == 'P' && b[1] == 'K' {
		return true
	}
	return false
}

func readVbaProjectBin(fileName string) (b []byte, err error) {
	reader, err := zip.OpenReader(fileName)
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	for _, f := range reader.File {
		if f.Name == "xl/vbaProject.bin" {
			rc, err := f.Open() // readCloser	rc
			if err != nil {
				return nil, err
			}

			if b, err = ioutil.ReadAll(rc); err != nil {
				return nil, err
			}

			return b, nil
		}
	}
	return nil, errors.New("err: 没有找到 vbaProject.bin")
}
