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
	xf.fileSize = int32(math.Pow(2, float64(xf.cfs.header.Sector_size))) * xf.cfs.header.Sat_count * 128 // 1个分区表最多记录128个id
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
