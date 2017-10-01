package compdocFile

import (
	"archive/zip"
	"encoding/binary"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
)

// zipFile格式需要先解压缩读取xl\vbaProject.bin
type zipFile struct {
	xlsFile
	vbaProjectIndex int
}

// 取消隐藏模块
func (me *zipFile) UnHideModule(moduleName string) (err error) {
	err = me.unHideModule(moduleName)
	if err != nil {
		return
	}
	return me.reWriteFile()
}
func (me *zipFile) UnProtectProject() (err error) {
	err = me.unProtectProject()
	if err != nil {
		return
	}
	return me.reWriteFile()
}

// 隐藏模块
func (me *zipFile) HideModule(moduleName string) (err error) {
	err = me.hideModule(moduleName)
	if err != nil {
		return
	}
	return me.reWriteFile()
}

func (me *zipFile) readFileByte() (err error) {
	reader, err := zip.OpenReader(me.fileName)
	if err != nil {
		return
	}
	defer reader.Close()

	for i, f := range reader.File {
		if f.Name == "xl/vbaProject.bin" {
			rc, err := f.Open() // readCloser	rc
			if err != nil {
				return err
			}

			me.vbaProjectIndex = i

			me.fileSize = f.UncompressedSize64
			me.cfs.fileByte = make([]byte, me.fileSize)

			var pFileByte uint64 = 0
			for pFileByte < me.fileSize {
				n, _ := rc.Read(me.cfs.fileByte[pFileByte:]) // 一次只能读取32768个byte，不知道为什么
				pFileByte += uint64(n)
				//				fmt.Println("pFileByte=", pFileByte, "f.UncompressedSize64=", f.UncompressedSize64)
			}

			iSizeHeader := binary.Size(me.cfs.header)
			byte2struct(me.cfs.fileByte[:iSizeHeader], &me.cfs.header)

			return nil
		}

	}
	return errors.New("err: 没有找到 vbaProject.bin")
}

func (me *zipFile) reWriteFile() (err error) {
	zipReader, err := zip.OpenReader(me.fileName)
	if err != nil {
		return err
	}
	defer zipReader.Close()

	strExt := filepath.Ext(me.fileName)
	strFileSave := me.fileName[:len(me.fileName)-len(strExt)] + "(new)" + strExt

	fw, err := os.OpenFile(strFileSave, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer fw.Close()

	zipWriter := zip.NewWriter(fw)
	defer zipWriter.Close()

	for _, f := range zipReader.File {
		fr, err := f.Open()
		if err != nil {
			return err
		}
		b, err := ioutil.ReadAll(fr)
		if err != nil {
			return err
		}
		defer fr.Close()

		if f.Name == "xl/vbaProject.bin" {
			b = me.cfs.fileByte
		}

		wr, err := zipWriter.Create(f.Name)
		if err != nil {
			return err
		}
		wr.Write(b)
		zipWriter.Flush()
	}

	return nil
}

func NewZipFile(fileName string) *zipFile {
	zf := new(zipFile)
	zf.fileName = fileName

	cfs := new(cfStruct)
	zf.cfs = cfs

	return zf
}
