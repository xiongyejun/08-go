package keysInSQLite

import (
	//	"bufio"
	"fmt"
	//	"io"
	//	"os"
	"testing"
)

func TestPrintnum(t *testing.T) {
	//	//	fmt.Printf("NUM %d = %4b\r\n", NUM, NUM)
	//	//	fmt.Printf("LCASE_LETTER %d = %4b\r\n", LCASE_LETTER, LCASE_LETTER)
	//	//	fmt.Printf("UCASE_LETTER %d = %4b\r\n", UCASE_LETTER, UCASE_LETTER)
	//	//	fmt.Printf("PUNCTUATION %d = %4b\r\n", PUNCTUATION, PUNCTUATION)
	fmt.Printf("strSQL_include_NUM_and_LCASE_LETTER = %s\r\n", strSQL_include_NUM_and_LCASE_LETTER)

	//	//	fmt.Printf("%d, NUM | LCASE_LETTER = %d \r\n", getType("str1"), NUM|LCASE_LETTER)
	//	//	fmt.Printf("%d, NUM | LCASE_LETTER | PUNCTUATION = %d\r\n", getType("str1."), NUM|LCASE_LETTER|PUNCTUATION)
	//	//	fmt.Printf("%d, NUM | LCASE_LETTER | PUNCTUATION | UCASE_LETTER = %d\r\n", getType("str1.A"), NUM|LCASE_LETTER|PUNCTUATION|UCASE_LETTER)

	//	//	ch := make(chan []byte, 2)
	//	//	var count uint64
	//	//	var err error

	//	//	go func() {
	//	//		if err = SelectValue("type&"+strconv.Itoa(NUM)+">0", ch, &count); err != nil {
	//	//			t.Error(err.Error())
	//	//			return
	//	//		}
	//	//	}()

	//	//	for i := range ch {
	//	//		fmt.Println(i)
	//	//	}

	//	if err := GetDB(); err != nil {
	//		t.Error(err.Error())
	//		return
	//	}
	//	defer CloseDB()
	//	//	d.checkLen()
	//	//	Insert([]string{"0"})

	//	addTxt(`C:\Users\Administrator\Desktop\OfficeEncryption\6000常用密码字典.txt`)
	//	signChineseHabitTxt(`C:\Users\Administrator\Desktop\OfficeEncryption\6000常用密码字典.txt`)

	//	//		addTxt(`E:\keys.txt`)
}

//// 将某个文件里的所有数据添加到数据库
//func addTxt(filename string) {
//	if f, err := os.Open(filename); err != nil {
//		fmt.Println(err)
//		return
//	} else {
//		bf := bufio.NewReader(f)
//		src := make([]string, 10000)
//		var k int = 0
//		for {
//			b, _, err := bf.ReadLine()
//			if err == io.EOF {
//				break
//			}
//			src[k] = string(b)
//			k++
//			if k == 10000 {
//				Insert(src)
//				k = 0
//			}
//		}
//		if k > 0 {
//			Insert(src)
//		}
//	}
//}

//// 将某个文件里的所有数据标记为中国人习惯的
//func signChineseHabitTxt(filename string) {
//	if f, err := os.Open(filename); err != nil {
//		fmt.Println(err)
//		return
//	} else {
//		bf := bufio.NewReader(f)
//		src := make([]string, 10000)
//		var k int = 0
//		for {
//			b, _, err := bf.ReadLine()
//			if err == io.EOF {
//				break
//			}
//			src[k] = string(b)
//			k++
//			if k == 10000 {
//				d.signChineseHabit(src)
//				k = 0
//			}
//		}
//		if k > 0 {
//			d.signChineseHabit(src)
//		}
//	}
//}
