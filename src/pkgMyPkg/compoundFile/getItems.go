package compoundFile

import (
	"errors"
	"fmt"
	"pkgMyPkg/colorPrint"
	"strconv"
	"strings"
	"unicode"
)

//func (me *CompoundFile) GetStorage(StorageName string) (s *Storage, err error) {
//	return nil, nil
//}
// [6]DataSpaces\DataSpaceInfo\streamName

var cp *colorPrint.ColorDll

func (me *CompoundFile) GetStream(streamPath string) (b []byte, err error) {
	var s *Storage = me.Root
	var arr []string = strings.Split(streamPath, `\`)

	for i := 0; i < len(arr)-1; i++ {
		if index, ok := s.storageDic[arr[i]]; !ok {
			return nil, errors.New("不存在的路径，" + arr[i])
		} else {
			s = s.Storages[index]
		}
	}

	if index, ok := s.streamDic[arr[len(arr)-1]]; !ok {
		return nil, errors.New("不存在的Stream名称。" + arr[len(arr)-1])
	} else {
		// 读取stream是按照512的大小读取的，但最后1个可能没有512
		return s.Streams[index].stream.Bytes()[:s.streamSize[index]], nil
	}
}

func (me *CompoundFile) PrintOut() {
	cp = colorPrint.NewColorDll()
	me.printOut(me.Root, "")
}

func (me *CompoundFile) printOut(s *Storage, strPre string) {
	cp.SetColor(colorPrint.White, colorPrint.DarkYellow)
	fmt.Printf("%s%s [Storage]\r\n", strPre, getPrintString(s.dir.name))
	cp.UnSetColor()

	for i := range s.Storages {
		//		fmt.Printf("%s%s\r\n", strPre+"--", s.Storages[i].dir.name)
		me.printOut(s.Storages[i], strPre+"----")
	}

	for i := range s.Streams {
		fmt.Printf("%s%s [Stream]\r\n", strPre+"----", getPrintString(s.Streams[i].name))
	}
}

func getPrintString(str string) string {
	var strRet string

	for _, r := range []rune(str) {
		if !unicode.IsPrint(r) {
			strRet = strRet + "[" + strconv.Itoa(int(r)) + "]"
		} else {
			strRet += string(r)
		}
	}

	return strRet
}

type retStream struct {
	B    []byte
	Name string
}

// 保存所有的stream
func (me *CompoundFile) GetStreams() (ret []retStream) {
	ret = make([]retStream, len(me.cfs.arrStream))
	var k int = 0
	for i := range me.cfs.arrStream {
		if me.cfs.arrStream[i].stream != nil {
			ret[k].Name = me.cfs.arrStream[i].name
			ret[k].B = me.cfs.arrStream[i].stream.Bytes()
			k++
		}

	}
	return ret[:k]
}
