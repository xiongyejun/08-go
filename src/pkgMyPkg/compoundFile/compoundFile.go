// 复合文档

package compoundFile

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"pkgMyPkg/ucs2T0utf8"
)

const (
	CFHEADER_SIZE int32 = 512
	DIR_SIZE      int32 = 128
)

type Storage struct {
	dir *dirInfo // Storage在文件byte里的目录信息

	Streams  []*cfStream
	Storages []*Storage
	Parent   *Storage
}

type CompoundFile struct {
	SrcByte []byte
	Root    *Storage

	cfs *cfStruct
}

// new复合文档对象
func NewCompoundFile(b []byte) (cf *CompoundFile, err error) {
	id := []byte{208, 207, 17, 224, 161, 177, 26, 225}

	if len(b) < int(CFHEADER_SIZE) {
		return nil, errors.New("字节数还不够Header的字节数。")
	}
	for i, v := range id {
		if b[i] != v {
			return nil, errors.New("不符合复合文档结构。")
		}
	}

	cf = new(CompoundFile)
	cf.SrcByte = b

	return cf, nil
}

// 解析复合文档
func (me *CompoundFile) Parse() (err error) {
	if err = me.getHeader(); err != nil {
		return
	}

	if err = me.getMSAT(); err != nil {
		return
	}

	if err = me.getSAT(); err != nil {
		return
	}

	if err = me.getDir(); err != nil {
		return
	}

	if err = me.getSSAT(); err != nil {
		return
	}

	if err = me.getStream(); err != nil {
		return
	}

	me.initDirTree()

	return nil
}

// 获取文件头
func (me *CompoundFile) getHeader() (err error) {
	me.cfs = new(cfStruct)
	me.cfs.header = new(cfHeader)
	iSizeHeader := binary.Size(me.cfs.header)
	return byte2struct(me.SrcByte[:iSizeHeader], me.cfs.header)
}

// 获取主分区表
func (me *CompoundFile) getMSAT() (err error) {
	me.cfs.arrMSAT = make([]int32, me.cfs.header.Sat_count)

	for i := 0; i < 109; i++ {
		if me.cfs.header.Arr_sid[i] == -1 {
			return nil
		}
		me.cfs.arrMSAT[i] = me.cfs.header.Arr_sid[i]
	}

	// 获取109个另外的
	p_MSAT := 109
	nextSID := me.cfs.header.Msat_first_sid
	for {
		arr := [128]int32{}
		byte2struct(me.SrcByte[CFHEADER_SIZE+CFHEADER_SIZE*nextSID:], &arr)

		for i := 0; i < 127; i++ {
			if arr[i] == -1 {
				return
			}

			me.cfs.arrMSAT[p_MSAT] = arr[i]
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
func (me *CompoundFile) getSAT() (err error) {
	me.cfs.arrSAT = make([]int32, me.cfs.header.Sat_count*128)
	tmpArrSat := [128]int32{}
	pSAT := 0
	var i int32 = 0
	for ; i < me.cfs.header.Sat_count; i++ {
		byte2struct(me.SrcByte[CFHEADER_SIZE+CFHEADER_SIZE*me.cfs.arrMSAT[i]:], &tmpArrSat)
		copy(me.cfs.arrSAT[pSAT:], tmpArrSat[:])
		pSAT += 128
	}
	return nil
}

// 获取目录
func (me *CompoundFile) getDir() (err error) {
	pSID := me.cfs.header.Dir_first_sid
	me.cfs.arrDir = make([]cfDir, 0)
	me.cfs.arrDirAddr = make([]int32, 0)
	var pDir int32 = 0

	for pSID != -2 {
		tmpDir := cfDir{}
		tmp := CFHEADER_SIZE + CFHEADER_SIZE*pSID + DIR_SIZE*(pDir%4)
		byte2struct(me.SrcByte[tmp:], &tmpDir)
		me.cfs.arrDir = append(me.cfs.arrDir, tmpDir)
		me.cfs.arrDirAddr = append(me.cfs.arrDirAddr, tmp)
		pDir++
		if pDir%4 == 0 {
			pSID = me.cfs.arrSAT[pSID]
		}
	}
	// 因为是4个一次性读取，所以最后可能有0-3个空白
	var k int = len(me.cfs.arrDir) - 1
	for me.cfs.arrDir[k].Len_name == 0 {
		k--
	}
	k++
	me.cfs.arrDir = me.cfs.arrDir[:k]

	return nil
}

// 获取短扇区分区表
func (me *CompoundFile) getSSAT() (err error) {
	var nSSAT int32 = 0
	if me.cfs.header.Ssat_count == 0 {
		return
	}
	// 根目录的 stream_size 表示短流存放流的大小，每64个为一个short sector
	nSSAT = me.cfs.arrDir[0].Stream_size / 64
	me.cfs.arrSSAT = make([]int32, nSSAT)
	// 短流起始SID
	pSID := me.cfs.arrDir[0].First_SID
	var i int32 = 0
	for ; i < nSSAT; i++ {
		// 指向偏移地址，实际地址要加上 &file_byte[0]
		me.cfs.arrSSAT[i] = pSID*CFHEADER_SIZE + CFHEADER_SIZE + (i%8)*64
		// 到下一个SID
		if (i+1)%8 == 0 {
			pSID = me.cfs.arrSAT[pSID]
		}
	}

	return nil
}

// 把目录里的每个流信息读取出来，存放在结构cfStream里
func (me *CompoundFile) getStream() (err error) {
	var i int32 = 0
	var n int32 = int32(len(me.cfs.arrDir))
	me.cfs.arrStream = make([]*cfStream, n)
	me.cfs.dic = make(map[string]int32, 10)

	for ; i < n; i++ {
		me.cfs.arrStream[i] = new(cfStream)
		// 读取name的byte
		b := me.cfs.arrDir[i].Dir_name[:me.cfs.arrDir[i].Len_name-2]
		b, err := ucs2T0utf8.UCS2toUTF8(b)
		if err != nil {
			return err
		}
		name := string(b)
		//		fmt.Println(name, b, name == "DataSpaces")
		me.cfs.arrStream[i].name = name

		// 窗体的时候，会出现2个在dir中，但是type1个是2,1个是1
		if _, ok := me.cfs.dic[name]; ok {
			// 存在了就不执行后面语句
			continue
		}
		me.cfs.dic[name] = i //记录每个dir name 所在的下标

		if me.cfs.arrDir[i].CfType == 2 && me.cfs.arrDir[i].First_SID != -1 {
			// 1仓 2流 5根

			me.cfs.arrStream[i].stream = bytes.NewBuffer([]byte{})
			if me.cfs.arrDir[i].Stream_size < me.cfs.header.Min_stream_size {
				// short_sector，是短流
				me.cfs.arrStream[i].step = 64
				var shortSID int32 = me.cfs.arrDir[i].First_SID
				for int32(me.cfs.arrStream[i].stream.Len()) < me.cfs.arrDir[i].Stream_size {
					me.cfs.arrStream[i].address = append(me.cfs.arrStream[i].address, me.cfs.arrSSAT[shortSID])
					me.cfs.arrStream[i].stream.Write(me.SrcByte[me.cfs.arrSSAT[shortSID] : me.cfs.arrSSAT[shortSID]+64])
					shortSID++
				}

			} else {
				me.cfs.arrStream[i].step = 512
				var pSID int32 = me.cfs.arrDir[i].First_SID
				for int32(me.cfs.arrStream[i].stream.Len()) < me.cfs.arrDir[i].Stream_size {
					me.cfs.arrStream[i].address = append(me.cfs.arrStream[i].address, CFHEADER_SIZE+CFHEADER_SIZE*pSID)
					me.cfs.arrStream[i].stream.Write(me.SrcByte[CFHEADER_SIZE+CFHEADER_SIZE*pSID : CFHEADER_SIZE+CFHEADER_SIZE*pSID+512])
					pSID = me.cfs.arrSAT[pSID]
				}
			}
		}

	}

	return nil
}

func (me *CompoundFile) initDirTree() (err error) {
	// 弄成树形结构
	// 1仓storage 2流 5根
	// root的left和right也肯定要是-1
	if me.cfs.arrDir[0].CfType != 5 || me.cfs.arrDir[0].Left_child != -1 || me.cfs.arrDir[0].Right_child != -1 {
		return errors.New("第1个dir不是根。")
	}
	me.Root = me.newStorage(0)

	me.initStorage(me.Root)

	return nil
}

// 对于1个Storage s，它的left和right是和s挂在一个级别的，他们有同1个parent
func (me *CompoundFile) initStorage(s *Storage) {
	me.appendLeftRightSub(s, s.dir.dir.Sub_dir)
}

// 左右是添加到s的parent
func (me *CompoundFile) appendLeftRightSub(s *Storage, index int32) {
	if index != -1 {
		me.appendItemToStorage(s, index)

		me.appendLeftRightSub(s, me.cfs.arrDir[index].Left_child)
		me.appendLeftRightSub(s, me.cfs.arrDir[index].Right_child)
	}
}

// 添加元素到Storage
func (me *CompoundFile) appendItemToStorage(s *Storage, index int32) int {
	if me.cfs.arrDir[index].CfType == 1 {
		tmp := me.newStorage(index)
		tmp.Parent = s
		s.Storages = append(s.Storages, tmp)
		// 是个storage，继续初始化
		me.initStorage(tmp)
		return 1
	} else if me.cfs.arrDir[index].CfType == 2 {
		s.Streams = append(s.Streams, me.cfs.arrStream[index])
		return 2
	} else {
		fmt.Println("err")
		return 0
	}
}

// 新建1个Storage
func (me *CompoundFile) newStorage(index int32) *Storage {
	s := new(Storage)
	s.dir = new(dirInfo)
	s.dir.dir = me.cfs.arrDir[index]
	s.dir.dirAddr = me.cfs.arrDirAddr[index]
	s.dir.name = me.cfs.arrStream[index].name
	return s
}

func byte2struct(b []byte, pStruct interface{}) error {
	buf := bytes.NewBuffer(b)
	return binary.Read(buf, binary.LittleEndian, pStruct)
}
