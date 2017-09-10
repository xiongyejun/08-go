// 复合文档
package compdocFile

import (
	"os"
	"unsafe"

	"pkgAPI/kernel32"
	//	"unsafe"
)

type cf interface {
	ReadFileByte()
	reWriteFile()

	GetFileName() string
	GetFileByte() *[]byte
	GetFileSize() int32
}

func NewCF() cf {
	var c cf
	return c
}

type cfHeader struct {
	id                   [8]byte
	file_id              [16]byte
	file_format_revision int16
	file_format_version  int16
	memory_endian        int16
	sector_size          int16 // '扇区的大小 2的幂 通常为2^9=512
	short_sector_size    int16
	not_used_1           [10]byte
	sat_count            int32 //'分区表扇区的总数
	dir_first_sid        int32
	not_used_2           [4]byte
	min_stream_size      int32
	ssat_first_sid       int32
	ssat_count           int32
	msat_first_sid       int32
	msat_count           int32
	arr_sid              [109]int32
}

type cfStruct struct {
	fileByte []byte
	header   cfHeader
}

var cfs cfStruct

// 判断是否是复合文档
func IsCompdocFile(fileName string) bool {
	var id []byte = make([]byte, 8)
	f, _ := os.Open(fileName)
	defer f.Close()
	f.Read(id)

	for i, v := range []byte{208, 207, 17, 224, 161, 177, 26, 225} {
		if id[i] != v {
			return false
		}
	}
	return true
}

func getCfHeader(fileName string) {
	f, _ := os.Open(fileName)
	var iSizeHeader uintptr = unsafe.Sizeof(cfs.header)
	var b = make([]byte, int(iSizeHeader))
	f.Read(b)
	kernel32.MoveMemory(unsafe.Pointer(&cfs.header), unsafe.Pointer(&b[0]), iSizeHeader)
}
