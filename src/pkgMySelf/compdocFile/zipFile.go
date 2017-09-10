package compdocFile

type zipFile struct {
	fileName string
	fileByte *[]byte
}

func (zf *zipFile) ReadFileByte() {

}

func (zf *zipFile) reWriteFile() {

}

func NewZipFile(fileName string) *zipFile {
	zf := new(zipFile)
	zf.fileName = fileName
	return zf
}
