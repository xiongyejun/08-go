// 复合文档
package compdocFile

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"pkgMySelf/ucs2T0utf8"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

// 复合文档接口
type CF interface {
	readFileByte() error
	reWriteFile()

	GetFileName() string
	GetFileByte() *[]byte
	GetFileSize() uint64

	GetCFStruct() *cfStruct
	GetModuleString(strModuleName string) string
}

const (
	CFHEADER_SIZE int32 = 512
	DIR_SIZE      int32 = 128
)

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

func IsZip(fileName string) bool {
	b := make([]byte, 2)
	f, _ := os.Open(fileName)
	defer f.Close()
	f.Read(b)
	if b[0] == 'P' && b[1] == 'K' {
		return true
	}
	return false
}

func CFInit(c CF) (err error) {
	fmt.Println("cfinit")

	err = c.readFileByte()
	if err != nil {
		return err
	}

	err = getMSAT(c.GetCFStruct())
	if err != nil {
		return err
	}

	err = getSAT(c.GetCFStruct())
	if err != nil {
		return err
	}

	err = getDir(c.GetCFStruct())
	if err != nil {
		return err
	}

	err = getSSAT(c.GetCFStruct())
	if err != nil {
		return err
	}

	err = getStream(c.GetCFStruct())
	if err != nil {
		return err
	}
	err = getDirInfo(c.GetCFStruct())
	if err != nil {
		return err
	}

	return nil
}

func getDirInfo(cfs *cfStruct) (err error) {
	dirIndex := cfs.dic["dir"]
	b := cfs.arrStream[dirIndex].stream.Bytes()[:cfs.arrDir[dirIndex].Stream_size]
	b = unCompressStream(b[1:]) // 解压的时候要跳过第1个标志位

	cfs.arrDirInfo = getModuleInfo(b)
	cfs.dicModule = make(map[string]int32, 10)
	for i := 0; i < len(cfs.arrDirInfo); i++ {
		cfs.dicModule[cfs.arrDirInfo[i].name] = int32(i)
	}
	return nil
}

func printTest(cfs *cfStruct) {
	for i := 0; i < len(cfs.arrDir); i++ {
		b := cfs.arrDir[i].Dir_name[:cfs.arrDir[i].Len_name-2]
		b, err := ucs2T0utf8.UCS2toUTF8(b)
		if err != nil {
			fmt.Println(err)
			return
		}
		name := string(b)
		//		fmt.Println(name)

		//				if name == "PROJECT" {
		//					b, _ := gbkToUtf8(cfs.arrStream[i].stream.Bytes())
		//					fmt.Println(string(b))
		//				}

		if name == "dir" {
			b := cfs.arrStream[i].stream.Bytes()[:cfs.arrDir[i].Stream_size]
			fmt.Println("dirstream=", len(b))
			b = unCompressStream(b[1:]) // 解压的时候要跳过第1个标志位

			arrDirInfo := getModuleInfo(b)
			fmt.Println(len(arrDirInfo))
			for j := 0; j < len(arrDirInfo); j++ {
				fmt.Println(arrDirInfo[j].name)
			}
		}
	}
}

// 获取主分区表
func getMSAT(cfs *cfStruct) (err error) {
	cfs.arrMSAT = make([]int32, cfs.header.Sat_count)

	for i := 0; i < 109; i++ {
		if cfs.header.Arr_sid[i] == -1 {
			return nil
		}
		cfs.arrMSAT[i] = cfs.header.Arr_sid[i]
	}

	// 获取109个另外的
	p_MSAT := 109
	nextSID := cfs.header.Msat_first_sid
	for {
		arr := [128]int32{}
		byte2struct(cfs.fileByte[CFHEADER_SIZE+CFHEADER_SIZE*nextSID:], &arr)
		//		kernel32.MoveMemory(unsafe.Pointer(&arr[0]), unsafe.Pointer(&cfs.fileByte[CFHEADER_SIZE+CFHEADER_SIZE*nextSID]), uintptr(CFHEADER_SIZE))

		for i := 0; i < 127; i++ {
			if arr[i] == -1 {
				return
			}

			cfs.arrMSAT[p_MSAT] = arr[i]
			p_MSAT++
		}
		nextSID = arr[127]
		if nextSID == -2 {
			break
		}
	}

	return nil
}

// 获取分区表
func getSAT(cfs *cfStruct) (err error) {
	cfs.arrSAT = make([]int32, cfs.header.Sat_count*128)
	tmpArrSat := [128]int32{}
	pSAT := 0
	var i int32 = 0
	for ; i < cfs.header.Sat_count; i++ {
		byte2struct(cfs.fileByte[CFHEADER_SIZE+CFHEADER_SIZE*cfs.arrMSAT[i]:], &tmpArrSat)
		//		kernel32.MoveMemory(unsafe.Pointer(&cfs.arrSAT[pSAT]), unsafe.Pointer(&cfs.fileByte[CFHEADER_SIZE+CFHEADER_SIZE*cfs.arrMSAT[i]]), uintptr(CFHEADER_SIZE))
		copy(cfs.arrSAT[pSAT:], tmpArrSat[:])
		pSAT += 128
	}
	return nil
}

// 获取目录
func getDir(cfs *cfStruct) (err error) {
	pSID := cfs.header.Dir_first_sid
	cfs.arrDir = make([]cfDir, 0, 10)
	var pDir int32 = 0

	for pSID != -2 {
		tmpDir := cfDir{}
		byte2struct(cfs.fileByte[CFHEADER_SIZE+CFHEADER_SIZE*pSID+DIR_SIZE*(pDir%4):], &tmpDir)
		//		kernel32.MoveMemory(unsafe.Pointer(&tmpDir.dir_name[0]), unsafe.Pointer(&cfs.fileByte[CFHEADER_SIZE+CFHEADER_SIZE*pSID+DIR_SIZE*(pDir%4)]), uintptr(DIR_SIZE))
		cfs.arrDir = append(cfs.arrDir, tmpDir)
		pDir++
		if pDir%4 == 0 {
			pSID = cfs.arrSAT[pSID]
		}
	}

	return nil
}

// 获取短扇区分区表
func getSSAT(cfs *cfStruct) (err error) {
	var nSSAT int32 = 0
	if cfs.header.Ssat_count == 0 {
		return
	}
	// 根目录的 stream_size 表示短流存放流的大小，每64个为一个short sector
	nSSAT = cfs.arrDir[0].Stream_size / 64
	cfs.arrSSAT = make([]int32, nSSAT)
	// 短流起始SID
	pSID := cfs.arrDir[0].First_SID
	var i int32 = 0
	for ; i < nSSAT; i++ {
		// 指向偏移地址，实际地址要加上 &file_byte[0]
		cfs.arrSSAT[i] = pSID*CFHEADER_SIZE + CFHEADER_SIZE + (i%8)*64
		// 到下一个SID
		if (i+1)%8 == 0 {
			pSID = cfs.arrSAT[pSID]
		}
	}

	return nil
}

// 把目录里的每个流信息读取出来，存放在结构cfStream里
func getStream(cfs *cfStruct) (err error) {
	var i int32 = 0
	var n int32 = int32(len(cfs.arrDir))
	cfs.arrStream = make([]*cfStream, n)
	cfs.dic = make(map[string]int32, 10)

	for ; i < n; i++ {
		if 0 == cfs.arrDir[i].Len_name { // dir读取的时候可能出现空的dir
			continue
		}
		b := cfs.arrDir[i].Dir_name[:cfs.arrDir[i].Len_name-2]
		b, err := ucs2T0utf8.UCS2toUTF8(b)

		if err != nil {
			return err
		}
		name := string(b)
		cfs.dic[name] = i //记录每个dir name 所在的下标
		cfs.arrStream[i] = new(cfStream)
		cfs.arrStream[i].name = name

		if cfs.arrDir[i].CfType == 2 && cfs.arrDir[i].First_SID != -1 {
			// 1仓 2流 5根

			cfs.arrStream[i].stream = bytes.NewBuffer([]byte{})
			if cfs.arrDir[i].Stream_size < cfs.header.Min_stream_size {
				// short_sector，是短流
				cfs.arrStream[i].step = 64
				var shortSID int32 = cfs.arrDir[i].First_SID
				for int32(len(cfs.arrStream[i].stream.Bytes())) < cfs.arrDir[i].Stream_size {
					cfs.arrStream[i].address = append(cfs.arrStream[i].address, cfs.arrSSAT[shortSID])
					cfs.arrStream[i].stream.Write(cfs.fileByte[cfs.arrSSAT[shortSID] : cfs.arrSSAT[shortSID]+64])
					shortSID++
				}

			} else {
				cfs.arrStream[i].step = 512
				var pSID int32 = cfs.arrDir[i].First_SID
				for int32(len(cfs.arrStream[i].stream.Bytes())) < cfs.arrDir[i].Stream_size {
					cfs.arrStream[i].address = append(cfs.arrStream[i].address, CFHEADER_SIZE+CFHEADER_SIZE*pSID)
					cfs.arrStream[i].stream.Write(cfs.fileByte[CFHEADER_SIZE+CFHEADER_SIZE*pSID : CFHEADER_SIZE+CFHEADER_SIZE*pSID+512])
					pSID = cfs.arrSAT[pSID]
				}
			}

		}

	}

	return nil
}

func gbkToUtf8(b []byte) ([]byte, error) {
	reader := transform.NewReader(bytes.NewReader(b), simplifiedchinese.GBK.NewDecoder())
	d, e := ioutil.ReadAll(reader)
	if e != nil {
		return nil, e
	}
	return d, nil

	//			simplifiedchinese.HZGB2312.NewDecoder()
}
