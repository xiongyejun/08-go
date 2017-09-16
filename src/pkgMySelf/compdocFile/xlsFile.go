package compdocFile

import (
	"math"
	"os"
)

type xlsFile struct {
	fileName string
	fileSize int32
	cfs      *cfStruct
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

func (xf *xlsFile) GetFileSize() int32 {
	return xf.fileSize
}

func (xf *xlsFile) readFileByte() (err error) {
	getCfHeader(xf.cfs, xf.fileName)
	xf.fileSize = int32(math.Pow(2, float64(xf.cfs.header.sector_size))) * xf.cfs.header.sat_count * 128 // 1个分区表最多记录128个id
	xf.cfs.fileByte = make([]byte, xf.fileSize)
	f, err := os.Open(xf.fileName)
	if err != nil {
		return err
	}
	defer f.Close()

	n, err := f.Read(xf.cfs.fileByte)
	if err != nil {
		return err
	}

	xf.fileSize = int32(n)
	xf.cfs.fileByte = xf.cfs.fileByte[:n]

	return nil
}

func (xf *xlsFile) reWriteFile() {

}

func NewXlsFile(fileName string) *xlsFile {
	xf := new(xlsFile)
	xf.fileName = fileName

	cfs := new(cfStruct)
	xf.cfs = cfs
	return xf
}
