package compdocFile

import (
	"math"
	"os"
)

type xlsFile struct {
	fileName string
	fileByte []byte
	fileSize int32
}

func (xf *xlsFile) GetFileName() string {
	return xf.fileName
}

func (xf *xlsFile) GetFileByte() *[]byte {
	return &xf.fileByte
}

func (xf *xlsFile) GetFileSize() int32 {
	return xf.fileSize
}

func (xf *xlsFile) ReadFileByte() {
	getCfHeader(xf.fileName)
	xf.fileSize = int32(math.Pow(2, float64(cfs.header.sector_size))) * cfs.header.sat_count * 127 // 1个分区表最多记录127个
	xf.fileByte = make([]byte, xf.fileSize)
	f, _ := os.Open(xf.fileName)
	n, _ := f.Read(xf.fileByte)
	xf.fileSize = int32(n)
	xf.fileByte = xf.fileByte[:n]

}

func (xf *xlsFile) reWriteFile() {

}

func NewXlsFile(fileName string) *xlsFile {
	xf := new(xlsFile)
	xf.fileName = fileName
	return xf
}
