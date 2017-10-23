// 将文件的二进制数据用16进制表示后
// 为了节省字符，利用0x21(!)——0x7E(~)中
// 除0x30(0)——0x39(9)、0x41(A)——0x46(F)
// 之外的可打印其他字符，不用16进制的2位数表示，直接用其本身
// 不知道什么原因，好像有些字符用了就扫不出来

package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

// 如果是0x30(0)——0x39(9)、0x41(A)——0x46(F)，记录为true
type myType struct {
	arrPrint   [0x7E + 1]bool
	fileName   string
	fileByte   []byte
	doCompress *bool
}

var mt = new(myType)

func init() {
	for i := 0; i < 0x7E; i++ {
		mt.arrPrint[i] = true
	}

	for i := 'G'; i < ('~' + 1); i++ {
		mt.arrPrint[i] = false
	}
	for i := ':'; i < ('@' + 1); i++ {
		mt.arrPrint[i] = false
	}
	for i := '!'; i < ('/' + 1); i++ {
		mt.arrPrint[i] = false
	}
	// 下面这几个有什么特殊意思？？
	for i := '%'; i < ('&' + 1); i++ {
		mt.arrPrint[i] = true
	}
}

func (me *myType) IsPrintHex(b byte) bool {
	return me.arrPrint[b]
}

func main() {
	initForm()

	fmt.Println("ok")
}

func (me *myType) unCompress() (err error) {

	var buf = bytes.NewBuffer([]byte{})
	var p int = 0
	n := len(me.fileByte)

	for p < n {
		if me.IsPrintHex(me.fileByte[p]) {
			strhex := string(me.fileByte[p : p+2])
			if tmp, err := strconv.ParseUint(strhex, 16, 8); err != nil {
				return err
			} else {
				buf.WriteByte(byte(tmp))
				p += 2
			}
		} else {
			buf.WriteByte(me.fileByte[p])
			p++
		}
	}
	me.byteToFile(buf.Bytes())

	return nil
}

func (me *myType) compress() (str string, err error) {
	str = me.byteToMyHex(me.fileByte)
	if err = ioutil.WriteFile(me.fileName+".txt", []byte(str), 0666); err != nil {
		return
	}

	return
}

func (me *myType) byteToMyHex(b []byte) string {
	str := make([]string, len(b))

	for i, v := range b {
		if v > 0x7E || me.IsPrintHex(v) {
			str[i] = fmt.Sprintf("%02X", v)
		} else {
			str[i] = fmt.Sprintf("%c", v)
		}
	}

	return strings.Join(str, "")
}

// byte保存文件
func (me *myType) byteToFile(b []byte) {
	fs, err := os.OpenFile(me.fileName+getExt(b), os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer fs.Close()
	fs.Write(b)
}

// 获取扩展名
func getExt(b []byte) string {
	if b[0] == '7' && b[1] == 'z' {
		return ".7z"
	} else if b[0] == 'P' && b[1] == 'K' {
		return ".zip"
	} else {
		return ".bin"
	}
}
