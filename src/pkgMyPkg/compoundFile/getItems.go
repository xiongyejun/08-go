package compoundFile

import (
	"errors"
	"fmt"
)

//func (me *CompoundFile) GetStorage(StorageName string) (s *Storage, err error) {
//	return nil, nil
//}

func (me *CompoundFile) GetStream(StreamName string) (b []byte, err error) {
	if index, ok := me.cfs.dic[StreamName]; !ok {
		return nil, errors.New("不存在的Stream名称。")
	} else {
		if me.cfs.arrDir[index].CfType != 2 {
			return nil, errors.New("不是Stream。")
		} else {
			return me.cfs.arrStream[index].stream.Bytes(), nil
		}
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
