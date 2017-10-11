package compdocFile

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"path/filepath"
	"pkgMySelf/colorPrint"
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

func (me *xlsFile) UnProtectProject() (err error) {
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
func (me *xlsFile) UnHideModule(moduleName string) (err error) {
	err = me.unHideModule(moduleName)
	if err != nil {
		return
	}
	return me.reWriteFile()
}
func (me *xlsFile) unHideModule(moduleName string) (err error) {
	// HelpFile="" 在这个前面添加 Module=moduleNameODOA
	if streamIndex, ok := me.cfs.dic["PROJECT"]; ok {
		// 读取PROJECT的byte
		b := me.cfs.arrStream[streamIndex].stream.Bytes()
		bModule := []byte(utf8ToGbk(`Module=` + moduleName))
		bModule = append(bModule, '\r')
		bModule = append(bModule, '\n')

		bOld := []byte(`HelpFile=""`)
		bNew := make([]byte, len(bModule)+len(bOld))
		copy(bNew[0:], bModule)
		copy(bNew[len(bModule):], bOld)

		b1 := bytes.Replace(b, bOld, bNew, 1)
		fmt.Println("b=", len(b), "b1=", len(b1))
		err = me.modifyProject(b, b1, streamIndex)
		return err
	} else {
		return errors.New("未找到PROJECT。")
	}
	return nil
}

// 隐藏模块
func (me *xlsFile) HideModule(moduleName string) (err error) {
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

// 在清除工程密码、隐藏模块后等操作后，将filebyte重新保存文件
func (me *xlsFile) reWriteFile() (err error) {
	strExt := filepath.Ext(me.fileName)
	strFileSave := me.fileName[:len(me.fileName)-len(strExt)] + "(new)" + strExt

	return ioutil.WriteFile(strFileSave, me.cfs.fileByte, 0666)
	//	fs, err := os.OpenFile(strFileSave, os.O_CREATE|os.O_WRONLY, 0666)
	//	if err != nil {
	//		return
	//	}
	//	fs.Write(me.cfs.fileByte)
	//	return
}

func (me *xlsFile) GetModuleName() (modules []string) {
	modules = make([]string, len(me.cfs.arrDirInfo))
	var strType string
	for i := 0; i < len(me.cfs.arrDirInfo); i++ {
		if me.cfs.arrDirInfo[i].moduleType == 0x21 {
			strType = "标准模块"
		} else {
			strType = "类模块"
		}
		modules[i] = me.cfs.arrDirInfo[i].name + "\t" + strType
	}
	return
}

func (xf *xlsFile) GetModuleString(strModuleName string) string {
	if streamIndex, ok := xf.cfs.dic[strModuleName]; ok {
		if dirInfoIndex, ok2 := xf.cfs.dicModule[strModuleName]; ok2 {
			b := xf.cfs.arrStream[streamIndex].stream.Bytes()[xf.cfs.arrDirInfo[dirInfoIndex].textOffset:]
			b = unCompressStream(b)
			b, _ = gbkToUtf8(b)
			return string(b)
		} else {
			return "不存在的模块名称。"
		}

	}
	return "不存在的目录名称。"
}

func (me *xlsFile) PrintAllCode() {
	cd := colorPrint.NewColorDll()
	for i, v := range me.cfs.arrDirInfo {
		cd.SetColor(colorPrint.White, colorPrint.DarkMagenta)
		fmt.Print("\r\n")
		fmt.Printf("%2d--------%s.moduleType(33是标准模块，34是其他)=%d--------", i, v.name, v.moduleType)
		cd.SetColor(colorPrint.White, colorPrint.DarkGreen)
		fmt.Print("\r\n")
		if streamIndex, ok := me.cfs.dic[v.name]; ok {
			b := me.cfs.arrStream[streamIndex].stream.Bytes()[v.textOffset:]
			b = unCompressStream(b)
			b, _ = gbkToUtf8(b)
			fmt.Print(string(b))
			cd.UnSetColor()
			fmt.Print("\r\n")
		}
	}
}

func (me *xlsFile) GetAllCode() string {
	str := make([]string, 0)

	for i, v := range me.cfs.arrDirInfo {

		str = append(str, fmt.Sprintf("%2d--------%s.moduleType(33是标准模块，34是其他)=%d--------\r\n", i, v.name, v.moduleType))
		if streamIndex, ok := me.cfs.dic[v.name]; ok {
			b := me.cfs.arrStream[streamIndex].stream.Bytes()[v.textOffset:]
			b = unCompressStream(b)
			b, _ = gbkToUtf8(b)
			str = append(str, string(b))
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
