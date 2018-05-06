package compoundFile

import (
	"errors"
	"fmt"
	"strings"
)

//func (me *CompoundFile) GetStorage(StorageName string) (s *Storage, err error) {
//	return nil, nil
//}
// [6]DataSpaces\DataSpaceInfo\streamName
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
		return nil, errors.New("不存在的Stream名称。")
	} else {
		// 读取stream是按照512的大小读取的，但最后1个可能没有512
		return s.Streams[index].stream.Bytes()[:s.streamSize[index]], nil
	}
}

func (me *CompoundFile) PrintOut() {
	me.printOut(me.Root, "")
}

func (me *CompoundFile) printOut(s *Storage, strPre string) {
	fmt.Printf("%s%s[Storage]\r\n", strPre, firstIs6(s.dir.name))

	for i := range s.Storages {
		//		fmt.Printf("%s%s\r\n", strPre+"--", s.Storages[i].dir.name)
		me.printOut(s.Storages[i], strPre+"----")
	}

	for i := range s.Streams {
		fmt.Printf("%s%s[Stream]\r\n", strPre+"----", firstIs6(s.Streams[i].name))
	}
}

func firstIs6(str string) string {
	if str[0] == 6 {
		return "[6]" + str[1:]
	}
	return str
}
