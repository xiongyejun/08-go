package compdocFile

import (
	"archive/zip"
	"encoding/binary"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// zipFile格式需要先解压缩读取xl\vbaProject.bin
type zipFile struct {
	xlsFile
	vbaProjectIndex int
}

// 取消隐藏模块
func (me *zipFile) UnHideModule(moduleName string) (newFile string, err error) {
	err = me.unHideModule(moduleName)
	if err != nil {
		return
	}
	return me.reWriteFile()
}
func (me *zipFile) UnProtectProject() (newFile string, err error) {
	err = me.unProtectProject()
	if err != nil {
		return
	}
	return me.reWriteFile()
}
func (me *zipFile) UnProtectSheetProtection() (newFile string, err error) {
	// 在xl\worksheets\ 下，找每个sheet 的：
	// <sheetProtection algorithmName="SHA-512" hashValue="wX1JS/iCwuxonczqbHLNhh/z0pPa+PBgEf3lErY+va1dcRSoIGoLDtDs7fF6J3HtvUGeIovMVENm6cea6xwqkg==" saltValue="Bl2sMnDmaODE073NETbEuA==" spinCount="100000" sheet="1" objects="1" scenarios="1" formatCells="0" formatColumns="0" formatRows="0"/>
	// 替换为空
	// 读取zip文件
	zipReader, err := zip.OpenReader(me.fileName)
	if err != nil {
		return newFile, err
	}
	defer zipReader.Close()
	// 设置新文件名
	strExt := filepath.Ext(me.fileName)
	newFile = me.fileName[:len(me.fileName)-len(strExt)] + "(new)" + strExt
	// 创建新文件
	fw, err := os.OpenFile(newFile, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return newFile, err
	}
	defer fw.Close()
	// 创建zip writer
	zipWriter := zip.NewWriter(fw)
	defer zipWriter.Close()
	// 循环zip文件中的文件
	for _, f := range zipReader.File {
		// 打开子文件
		fr, err := f.Open()
		if err != nil {
			return newFile, err
		}
		// 读取子文件流
		b, err := ioutil.ReadAll(fr)
		if err != nil {
			return newFile, err
		}
		defer fr.Close()
		// 如果是sheet，就改写流
		if strings.HasPrefix(f.Name, "xl/worksheets/") {
			reg, _ := regexp.Compile("<sheetProtection .*?/>")
			b = reg.ReplaceAll(b, []byte{})
		}
		// 在zipwriter中创建新文件
		wr, err := zipWriter.Create(f.Name)
		if err != nil {
			return newFile, err
		}
		// 写入新文件的数据
		n := 0
		n, err = wr.Write(b)
		if err != nil {
			return newFile, err
		}
		if n < len(b) {
			return newFile, errors.New("写入不完整")
		}
	}
	err = zipWriter.Flush()

	return newFile, err
}

// 隐藏模块
func (me *zipFile) HideModule(moduleName string) (newFile string, err error) {
	err = me.hideModule(moduleName)
	if err != nil {
		return
	}
	return me.reWriteFile()
}

func (me *zipFile) ReWriteFile(startAddress int, modifyByte []byte) (newFile string, err error) {
	copy(me.cfs.fileByte[startAddress:], modifyByte)
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

			//			var pFileByte uint64 = 0
			//			for pFileByte < me.fileSize {
			//				n, _ := rc.Read(me.cfs.fileByte[pFileByte:]) // 一次只能读取32768个byte，不知道为什么
			//				pFileByte += uint64(n)
			//				//				fmt.Println("pFileByte=", pFileByte, "f.UncompressedSize64=", f.UncompressedSize64)
			//			}
			if me.cfs.fileByte, err = ioutil.ReadAll(rc); err != nil {
				return err
			}

			iSizeHeader := binary.Size(me.cfs.header)
			byte2struct(me.cfs.fileByte[:iSizeHeader], &me.cfs.header)

			return nil
		}

	}
	return errors.New("err: 没有找到 vbaProject.bin")
}

func (me *zipFile) reWriteFile() (newFile string, err error) {
	// 读取zip文件
	zipReader, err := zip.OpenReader(me.fileName)
	if err != nil {
		return newFile, err
	}
	defer zipReader.Close()
	// 设置新文件名
	strExt := filepath.Ext(me.fileName)
	newFile = me.fileName[:len(me.fileName)-len(strExt)] + "(new)" + strExt
	// 创建新文件
	fw, err := os.OpenFile(newFile, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return newFile, err
	}
	defer fw.Close()
	// 创建zip writer
	zipWriter := zip.NewWriter(fw)
	defer zipWriter.Close()
	// 循环zip文件中的文件
	for _, f := range zipReader.File {
		// 打开子文件
		fr, err := f.Open()
		if err != nil {
			return newFile, err
		}
		// 读取子文件流
		b, err := ioutil.ReadAll(fr)
		if err != nil {
			return newFile, err
		}
		defer fr.Close()
		// 如果是vba，就用改写了的流
		if f.Name == "xl/vbaProject.bin" {
			b = me.cfs.fileByte
		}
		// 在zipwriter中创建新文件
		wr, err := zipWriter.Create(f.Name)
		if err != nil {
			return newFile, err
		}
		// 写入新文件的数据
		n := 0
		n, err = wr.Write(b)
		if err != nil {
			return newFile, err
		}
		if n < len(b) {
			return newFile, errors.New("写入不完整")
		}
	}
	err = zipWriter.Flush()

	return newFile, err
}

func NewZipFile(fileName string) *zipFile {
	zf := new(zipFile)
	zf.fileName = fileName

	cfs := new(cfStruct)
	zf.cfs = cfs

	return zf
}
