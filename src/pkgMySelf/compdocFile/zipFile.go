package compdocFile

type zipFile struct {
	fileName string
	fileByte *[]byte
}

func (xf zipFile) GetCFStruct() *cfStruct {
	return nil
}

func (zf zipFile) readFileByte() error {
	return nil
}

func (zf zipFile) GetFileByte() *[]byte {
	return nil
}
func (zf zipFile) GetFileName() string {
	return zf.fileName
}
func (zf zipFile) GetFileSize() int32 {
	return 1
}
func (zf zipFile) reWriteFile() {

}

func NewZipFile(fileName string) *zipFile {
	zf := new(zipFile)
	zf.fileName = fileName
	return zf
}
