// 常用密码数据库

package keysInSQLite

import (
	"strconv"
)

const (
	NUM          = 0x1               // 0000 0000 0000 0001 纯数字
	LCASE_LETTER = NUM << 1          // 0000 0000 0000 0010 小写英文字母
	UCASE_LETTER = LCASE_LETTER << 1 // 0000 0000 0000 0100 大写英文字母
	PUNCTUATION  = UCASE_LETTER << 1 // 0000 0000 0000 1000 标点符号

	type_Counts = PUNCTUATION
)

// 有多少基础种类
var typeCounts int = 0
var allType int = NUM | LCASE_LETTER | PUNCTUATION | UCASE_LETTER

// 纯数字……
var strSQL_NUM string = "type=" + strconv.Itoa(NUM)
var strSQL_LCASE_LETTER string = "type=" + strconv.Itoa(LCASE_LETTER)
var strSQL_UCASE_LETTER string = "type=" + strconv.Itoa(UCASE_LETTER)
var strSQL_PUNCTUATION string = "type=" + strconv.Itoa(PUNCTUATION)

// 包含数字……
var strSQL_include_NUM string = "type&" + strconv.Itoa(NUM) + ">0"
var strSQL_include_LCASE_LETTER string = "type&" + strconv.Itoa(LCASE_LETTER) + ">0"
var strSQL_include_UCASE_LETTER string = "type&" + strconv.Itoa(UCASE_LETTER) + ">0"
var strSQL_include_PUNCTUATION string = "type&" + strconv.Itoa(PUNCTUATION) + ">0"

// 包含2种
var strSQL_include_NUM_and_LCASE_LETTER string = strSQL_include_NUM + " and " + strSQL_include_LCASE_LETTER

//var strSQL_include_LCASE_LETTER string = "type&" + strconv.Itoa(LCASE_LETTER) + ">0"
//var strSQL_include_UCASE_LETTER string = "type&" + strconv.Itoa(UCASE_LETTER) + ">0"
//var strSQL_include_PUNCTUATION string = "type&" + strconv.Itoa(PUNCTUATION) + ">0"

func init() {
	var tmp int = type_Counts
	for tmp > 0 {
		typeCounts++
		tmp = tmp >> 1
	}

	d = new(DataStruct)
	d.DBPath = `E:\keys.db`
	d.tableName = "password"
}

func getType(str string) (iType int) {
	for i := range str {
		if str[i] >= '0' && str[i] <= '9' {
			iType |= NUM
		} else if str[i] >= 'a' && str[i] <= 'z' {
			iType |= LCASE_LETTER
		} else if str[i] >= 'A' && str[i] <= 'Z' {
			iType |= UCASE_LETTER
		} else {
			iType |= PUNCTUATION
		}

		if iType == allType {
			return
		}
	}

	return
}
