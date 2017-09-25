package compdocFile

import (
	"encoding/binary"
	"io"
	"os"
)

type xlsFile struct {
	fileName string
	fileSize uint64
	cfs      *cfStruct

	rc io.ReadCloser
}

func (xf *xlsFile) GetCFStruct() *cfStruct {
	return xf.cfs
}

func (xf *xlsFile) GetFileName() string {
	return xf.fileName
}

func (xf *xlsFile) GetFileByte() *[]byte {
	return &xf.cfs.fileByte
}

func (xf *xlsFile) GetFileSize() uint64 {
	return xf.fileSize
}

func (this *xlsFile) readFileByte() (err error) {
	f, err := os.Open(this.fileName)
	if err != nil {
		return
	}
	fi, err := f.Stat()

	if err != nil {
		return
	}
	defer f.Close()

	this.fileSize = uint64(fi.Size())
	this.cfs.fileByte = make([]byte, this.fileSize)
	f.Read(this.cfs.fileByte)

	iSizeHeader := binary.Size(this.cfs.header)
	byte2struct(this.cfs.fileByte[:iSizeHeader], &this.cfs.header)

	return
}

func (xf *xlsFile) reWriteFile() {

}
func (xf *xlsFile) GetModuleString(strModuleName string) string {
	if streamIndex, ok := xf.cfs.dic[strModuleName]; ok {
		if dirInfoIndex, ok2 := xf.cfs.dicModule[strModuleName]; ok2 {
			b := xf.cfs.arrStream[streamIndex].stream.Bytes()[xf.cfs.arrDirInfo[dirInfoIndex].textOffset+1:]
			b = unCompressStream(b)
			b, _ = gbkToUtf8(b)
			return string(b)
		} else {
			return "不存在的模块名称。"
		}

	}
	return "不存在的目录名称。"
}

func NewXlsFile(fileName string) *xlsFile {
	xf := new(xlsFile)
	xf.fileName = fileName

	cfs := new(cfStruct)
	xf.cfs = cfs
	return xf
}
