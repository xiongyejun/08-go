package compdocFile

import (
	"archive/zip"
	"encoding/binary"
	"errors"
)

// zipFile格式需要先解压缩读取xl\vbaProject.bin
type zipFile struct {
	xlsFile
	vbaProjectIndex int
}

func (this *zipFile) readFileByte() (err error) {

	reader, err := zip.OpenReader(this.fileName)
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

			this.vbaProjectIndex = i

			this.fileSize = f.UncompressedSize64
			this.cfs.fileByte = make([]byte, this.fileSize)

			var pFileByte uint64 = 0
			for pFileByte < this.fileSize {
				n, _ := rc.Read(this.cfs.fileByte[pFileByte:]) // 一次只能读取32768个byte，不知道为什么
				pFileByte += uint64(n)
				//				fmt.Println("pFileByte=", pFileByte, "f.UncompressedSize64=", f.UncompressedSize64)
			}

			iSizeHeader := binary.Size(this.cfs.header)
			byte2struct(this.cfs.fileByte[:iSizeHeader], &this.cfs.header)

			return nil
		}

	}
	return errors.New("err: 没有找到 vbaProject.bin")
}

func (this *zipFile) reWriteFile() {

}

func NewZipFile(fileName string) *zipFile {
	zf := new(zipFile)
	zf.fileName = fileName

	cfs := new(cfStruct)
	zf.cfs = cfs

	return zf
}
