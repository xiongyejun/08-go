// 解压dir流和模块流
package compdocFile

import (
	"bytes"
	"encoding/binary"
	"io/ioutil"
	"syscall"
	"unsafe"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"

	"github.com/axgle/mahonia"
)

func byte2struct(b []byte, pStruct interface{}) {
	buf := bytes.NewBuffer(b)
	binary.Read(buf, binary.LittleEndian, pStruct)
}

// 解压steam
//func unCompressStream(compressByre []byte) (unCompressByte []byte, err error) {
//	ntdll, err := syscall.LoadLibrary("NTDLL.dll")
//	if err != nil {
//		return []byte{}, err
//	}
//	tdb, err := syscall.GetProcAddress(ntdll, "RtlDecompressBuffer")
//	if err != nil {
//		return []byte{}, err
//	}

//	var outSize int32 = 0
//	unCompressByte = make([]byte, 2*len(compressByre))
//	_, _, err = syscall.Syscall6(tdb,
//		4,
//		2,
//		uintptr(unsafe.Pointer(&unCompressByte[0])),
//		uintptr(len(unCompressByte)),
//		uintptr(unsafe.Pointer(&compressByre[0])),
//		uintptr(len(unCompressByte)),
//		uintptr(unsafe.Pointer(&outSize)))

//	fmt.Println("outSize", outSize)
//	return
//}

func unCompressStream(compressByre []byte) (unCompressByte []byte) {
	ntdll := syscall.NewLazyDLL("NTDLL.dll")
	tdb := ntdll.NewProc("RtlDecompressBuffer")

	var k int32 = 1
	iLen := int32(len(compressByre))
	var outSize int32 = iLen + 1

	for k*iLen <= outSize || outSize == 4096 {
		k++
		unCompressByte = make([]byte, k*iLen)
		tdb.Call(2,
			uintptr(unsafe.Pointer(&unCompressByte[0])),
			uintptr(len(unCompressByte)),
			uintptr(unsafe.Pointer(&compressByre[0])),
			uintptr(len(unCompressByte)),
			uintptr(unsafe.Pointer(&outSize)))
	}

	return unCompressByte[:outSize]
}

// RtlDecompressBuffer(CShort(2), p1, Result.Length, p2, Origin2.Length, ResultSize)

// Private Declare Function RtlDecompressBuffer Lib "NTDLL" (ByVal flags As Short,
//                    ByVal BuffUnCompressed As IntPtr, ByVal UnCompSize As Integer,
//                    ByVal BuffCompressed As IntPtr, ByVal CompBuffSize As Integer,
//                    ByRef OutputSize As Integer) As Integer

func gbkToUtf8(b []byte) ([]byte, error) {
	reader := transform.NewReader(bytes.NewReader(b), simplifiedchinese.GBK.NewDecoder())
	d, e := ioutil.ReadAll(reader)
	if e != nil {
		return nil, e
	}
	return d, nil

	//			simplifiedchinese.HZGB2312.NewDecoder()
}

func utf8ToGbk(src string) string {
	srcCoder := mahonia.NewEncoder("gbk")
	return srcCoder.ConvertString(src)
}
