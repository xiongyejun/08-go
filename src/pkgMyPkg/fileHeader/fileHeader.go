// 常用文件头判断
package fileHeader

import (
	"bytes"
	"errors"
)

type fileHeader struct {
	ext     string // 文件后缀
	bHeader []byte
}

const (
	JPG int = iota
	PNG
	GIF
	BMP
	ZIP
	RAR
	PDF
)

var fh []*fileHeader

func IsZIP(bFile []byte) bool {
	return isWhat(bFile, ZIP)
}

func isWhat(bFile []byte, index int) bool {
	return bytes.HasPrefix(bFile, fh[index].bHeader)
}

// 根据文件的字节，判断文件后缀
func GetExt(bFile []byte) (ext string, err error) {
	if len(bFile) < 4 {
		return "", errors.New("文件至少需要4个字节。")
	}

	for i := range fh {
		if bytes.HasPrefix(bFile, fh[i].bHeader) {
			return fh[i].ext, nil
		}
	}

	return "", errors.New("未知文件类型。")
}

func init() {
	fh = []*fileHeader{
		//JPEG (jpg)，文件头：FFD8FFE1
		&fileHeader{"jpg", []byte{0xff, 0xd8, 0xff, 0xe1}},
		//PNG (png)，文件头：89504E47
		&fileHeader{"png", []byte{0x89, 0x50, 0x4e, 0x47}},
		//GIF (gif)，文件头：47494638
		&fileHeader{"gif", []byte{0x47, 0x49, 0x46, 0x38}},
		//Windows Bitmap (bmp)，文件头：424DC001
		&fileHeader{"bmp", []byte{0x42, 0x4D, 0xC0, 0x01}},
		//ZIP Archive (zip)，文件头：504B0304
		&fileHeader{"zip", []byte{0x50, 0x4B, 0x03, 0x04}},
		//RAR Archive (rar)，文件头：52617221
		&fileHeader{"rar", []byte{0x52, 0x61, 0x72, 0x21}},
		// 25 50 44 46
		&fileHeader{"pdf", []byte{0x25, 0x50, 0x44, 0x46}},
	}

	//TIFF (tif)，文件头：49492A00

	//CAD (dwg)，文件头：41433130

	//Adobe Photoshop (psd)，文件头：38425053

	//Rich Text Format (rtf)，文件头：7B5C727466

	//XML (xml)，文件头：3C3F786D6C

	//HTML (html)，文件头：68746D6C3E

	//Email [thorough only] (eml)，文件头：44656C69766572792D646174653A

	//Outlook Express (dbx)，文件头：CFAD12FEC5FD746F

	//Outlook (pst)，文件头：2142444E

	//MS Word/Excel (xls.or.doc)，文件头：D0CF11E0

	//MS Access (mdb)，文件头：5374616E64617264204A

	//WordPerfect (wpd)，文件头：FF575043

	//Adobe Acrobat (pdf)，文件头：255044462D312E

	//Quicken (qdf)，文件头：AC9EBD8F

	//Windows Password (pwl)，文件头：E3828596

	//Wave (wav)，文件头：57415645

	//AVI (avi)，文件头：41564920

	//Real Audio (ram)，文件头：2E7261FD

	//Real Media (rm)，文件头：2E524D46

	//MPEG (mpg)，文件头：000001BA

	//MPEG (mpg)，文件头：000001B3

	//Quicktime (mov)，文件头：6D6F6F76

	//Windows Media (asf)，文件头：3026B2758E66CF11

	//MIDI (mid)，文件头：4D546864

}
