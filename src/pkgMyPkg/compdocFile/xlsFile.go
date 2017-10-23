package compdocFile

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"strings"
)

type xlsFile struct {
	fileName string
	fileSize uint64
	cfs      *cfStruct

	rc io.ReadCloser
}

func (xf *xlsFile) GetCFStruct() *cfStruct {
	return xf.cfs
}

func (xf *xlsFile) GetFileName() string {
	return xf.fileName
}

func (xf *xlsFile) GetFileByte() *[]byte {
	return &xf.cfs.fileByte
}

func (xf *xlsFile) GetFileSize() uint64 {
	return xf.fileSize
}

func (me *xlsFile) readFileByte() (err error) {
	me.cfs.fileByte, err = ioutil.ReadFile(me.fileName)
	if err != nil {
		return
	}
	me.fileSize = uint64(len(me.cfs.fileByte))

	iSizeHeader := binary.Size(me.cfs.header)
	byte2struct(me.cfs.fileByte[:iSizeHeader], &me.cfs.header)

	return
}

func (me *xlsFile) UnProtectProject() (newFile string, err error) {
	err = me.unProtectProject()
	if err != nil {
		return
	}
	return me.reWriteFile()
}

// 清除vba工程密码
func (me *xlsFile) unProtectProject() (err error) {
	// "CMG", "DPB", "GC"
	if streamIndex, ok := me.cfs.dic["PROJECT"]; ok {
		// 读取PROJECT的byte
		b := me.cfs.arrStream[streamIndex].stream.Bytes()
		var b1 []byte
		pattern := `CMG="[A-Z\d]+"\r\n|DPB="[A-Z\d]+"\r\n|GC="[A-Z\d]+"\r\n`

		if bMatch, _ := regexp.Match(pattern, b); !bMatch {
			err = errors.New("err:没有找到查找的内容。")
			return
		}
		reg, err := regexp.Compile(pattern)
		if err != nil {
			return err
		}
		// 替换后的byte
		b1 = reg.ReplaceAll(b, []byte{})
		err = me.modifyProject(b, b1, streamIndex)
		if err != nil {
			err = errors.New(err.Error() + "破解工程密码出错。")
			return err
		}
	} else {
		return errors.New("未找到PROJECT。")
	}
	return nil
}

// 取消隐藏模块
func (me *xlsFile) UnHideModule(moduleName string) (newFile string, err error) {
	err = me.unHideModule(moduleName)
	if err != nil {
		return
	}
	return me.reWriteFile()
}
func (me *xlsFile) unHideModule(moduleName string) (err error) {
	// HelpFile="" 在这个前面添加 Module=moduleNameODOA
	if streamIndex, ok := me.cfs.dic["PROJECT"]; ok {
		if _, ok := me.cfs.dicModule[moduleName]; ok {
			// 读取PROJECT的byte
			b := me.cfs.arrStream[streamIndex].stream.Bytes()
			bModule := []byte(utf8ToGbk(`Module=` + moduleName))
			bModule = append(bModule, '\r')
			bModule = append(bModule, '\n')
			// 判断是否是被隐藏了
			if bytes.Contains(b, bModule) {
				return errors.New("模块并没有被隐藏。")
			}

			bOld := []byte(`HelpFile=""`)
			bNew := make([]byte, len(bModule)+len(bOld))
			copy(bNew[0:], bModule)
			copy(bNew[len(bModule):], bOld)

			b1 := bytes.Replace(b, bOld, bNew, 1)
			err = me.modifyProject(b, b1, streamIndex)
			return err
		} else {
			return errors.New("不存在的模块名称。")
		}

	} else {
		return errors.New("未找到PROJECT。")
	}
	return nil
}

// 隐藏模块
func (me *xlsFile) HideModule(moduleName string) (newFile string, err error) {
	err = me.hideModule(moduleName)
	if err != nil {
		return
	}
	return me.reWriteFile()
}

func (me *xlsFile) hideModule(moduleName string) (err error) {
	// arrDirInfo 记录了dir中解压缩出来的模块名称、模块类型、模块偏移信息
	if dirInfoIndex, ok := me.cfs.dicModule[moduleName]; ok {
		if me.cfs.arrDirInfo[dirInfoIndex].moduleType == 0x22 {
			err = errors.New("不能隐藏类模块。")
			return
		}

		pattern := `Module=` + moduleName // + `\r\n` //|` + moduleName + `.*?\r\n`
		if streamIndex, ok := me.cfs.dic["PROJECT"]; ok {
			// 读取PROJECT的byte
			b := me.cfs.arrStream[streamIndex].stream.Bytes()
			var b1 []byte

			pattern = utf8ToGbk(pattern)
			bReplace := []byte(pattern)
			bReplace = append(bReplace, '\r')
			bReplace = append(bReplace, '\n')

			b1 = bytes.Replace(b, bReplace, []byte{}, -1)
			// project流中没有模块的信息
			if len(b) == len(b1) {
				err = errors.New(err.Error() + "已经是被隐藏的模块。")
				return
			}

			err = me.modifyProject(b, b1, streamIndex)

			if err != nil {
				err = errors.New(err.Error() + "隐藏模块出错。")
				return
			}
			return nil
		} else {
			return errors.New("未找到PROJECT。")
		}

	} else {
		return errors.New("不存在的模块名称。")
	}
	return nil
}

// 修改PROJECT目录流，主要是清除工程密码、隐藏模块等需要
func (me *xlsFile) modifyProject(oldB, newB []byte, streamIndex int32) (err error) {
	// b2保持大小不变，方便复制到filebyte
	b2 := make([]byte, len(oldB))
	copy(b2, newB)
	// 修改替换后的文件byte
	for i, v := range me.cfs.arrStream[streamIndex].address {
		bStart := int32(i) * me.cfs.arrStream[streamIndex].step
		bEnd := bStart + me.cfs.arrStream[streamIndex].step
		copy(me.cfs.fileByte[v:], b2[bStart:bEnd])
	}
	// 修改dir中的Stream_size
	// b中实际仅有me.cfs.arrDir[streamIndex].Stream_size的大小，但是为了上面循环方便按照step复制，在这里来扣除多余的
	iSub := int32(len(oldB)) - me.cfs.arrDir[streamIndex].Stream_size
	var iLen int32 = int32(len(newB)) - iSub
	// int32转byte
	buf := bytes.NewBuffer([]byte{})
	binary.Write(buf, binary.LittleEndian, &iLen)
	// 内存中的数据也修改下
	me.cfs.arrDir[streamIndex].Stream_size = iLen
	// fileByte的下标
	indexStreamSize := me.cfs.arrDirAddr[streamIndex] + DIR_SIZE - 8 // -8是因为在倒数第2个，减2个int32
	copy(me.cfs.fileByte[indexStreamSize:], buf.Bytes())

	return err
}

func (me *xlsFile) ReWriteFile(startAddress int, modifyByte []byte) (newFile string, err error) {
	copy(me.cfs.fileByte[startAddress:], modifyByte)
	return me.reWriteFile()
}

// 在清除工程密码、隐藏模块后等操作后，将filebyte重新保存文件
func (me *xlsFile) reWriteFile() (newFile string, err error) {
	strExt := filepath.Ext(me.fileName)
	newFile = me.fileName[:len(me.fileName)-len(strExt)] + "(new)" + strExt

	return newFile, ioutil.WriteFile(newFile, me.cfs.fileByte, 0666)
	//	fs, err := os.OpenFile(strFileSave, os.O_CREATE|os.O_WRONLY, 0666)
	//	if err != nil {
	//		return
	//	}
	//	fs.Write(me.cfs.fileByte)
	//	return
}

func (me *xlsFile) GetVBAInfo() (out []*OutStruct) {
	out = make([]*OutStruct, len(me.cfs.arrStream))

	for i := 0; i < len(me.cfs.arrStream); i++ {
		out[i] = NewOutStruct()
		out[i].Name = me.cfs.arrStream[i].name // stream已经记录了name
		if me.cfs.arrDir[i].CfType == 1 {      // 1仓 2流 5根
			out[i].Type = "仓"
		} else if me.cfs.arrDir[i].CfType == 5 { // 1仓 2流 5根
			out[i].Type = "根"
		} else if me.cfs.arrDir[i].CfType == 2 { // 1仓 2流 5根
			if dirInfoIndex, ok := me.cfs.dicModule[out[i].Name]; ok {
				if me.cfs.arrDirInfo[dirInfoIndex].moduleType == 0x22 {
					out[i].Type = "类模块流"
				} else {
					out[i].Type = "模块流"
				}
			} else {
				out[i].Type = "流"
			}
		}
	}
	return
}

func (me *xlsFile) GetModuleName() (modules [][2]string) {
	modules = make([][2]string, len(me.cfs.arrDirInfo))
	var strType string

	for i := 0; i < len(me.cfs.arrDirInfo); i++ {
		if me.cfs.arrDirInfo[i].moduleType == 0x21 {
			strType = "标准模块"
		} else {
			strType = "类模块"
		}
		modules[i][0] = me.cfs.arrDirInfo[i].name
		modules[i][1] = strType
	}
	return
}

// 返回数据的steam和数据的地址
func (me *xlsFile) GetStream(name string) (bStream []byte, bAddress []int32, step int32) {
	if streamIndex, ok := me.cfs.dic[name]; ok {
		if me.cfs.arrStream[streamIndex].stream == nil {
			return []byte{}, []int32{}, 0
		} else {
			bStream = me.cfs.arrStream[streamIndex].stream.Bytes()[:me.cfs.arrDir[streamIndex].Stream_size]
			bAddress = me.cfs.arrStream[streamIndex].address
			step = me.cfs.arrStream[streamIndex].step
			return
		}
	}

	return []byte("不存在的流。"), []int32{}, 0
}

func (me *xlsFile) GetModuleCode(strModuleName string) string {
	if streamIndex, ok := me.cfs.dic[strModuleName]; ok {
		if dirInfoIndex, ok2 := me.cfs.dicModule[strModuleName]; ok2 {
			// 看modifyProject里的说明，streamSize为什么没有调整
			b := me.cfs.arrStream[streamIndex].stream.Bytes()[me.cfs.arrDirInfo[dirInfoIndex].textOffset:me.cfs.arrDir[streamIndex].Stream_size]
			b = unCompressStream(b)
			return gbkToUtf8(b)
		} else {
			return "不存在的模块名称。"
		}

	}
	return "不存在的目录名称。"
}

func (me *xlsFile) GetAllCode() string {
	str := make([]string, 0)

	for i, v := range me.cfs.arrDirInfo {

		str = append(str, fmt.Sprintf("%2d--------%s.moduleType(33是标准模块，34是其他)=%d--------\r\n", i, v.name, v.moduleType))
		if streamIndex, ok := me.cfs.dic[v.name]; ok {
			// 看modifyProject里的说明，streamSize为什么没有调整
			b := me.cfs.arrStream[streamIndex].stream.Bytes()[v.textOffset:me.cfs.arrDir[streamIndex].Stream_size]
			b = unCompressStream(b)
			str = append(str, gbkToUtf8(b))
		}
	}
	return strings.Join(str, "\r\n")
}
func NewXlsFile(fileName string) *xlsFile {
	xf := new(xlsFile)
	xf.fileName = fileName

	cfs := new(cfStruct)
	xf.cfs = cfs
	return xf
}
