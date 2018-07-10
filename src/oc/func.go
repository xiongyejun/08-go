package main

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"
)

func getPPTEncryptionInfo(bEncryptionInfo []byte) (b []byte, err error) {
	// ppt加密信息放在文件的后面
	bFind := asc2Unicode([]byte("Microsoft Enhanced Cryptographic Provider v1"))
	index := bytes.LastIndex(bEncryptionInfo, bFind)
	if index == -1 {
		return nil, errors.New("PPT byte 未能找到 Microsoft Enhanced Cryptographic Provider v1")
	}
	index = index - 2*16 - 12
	return bEncryptionInfo[index:], err
}
func formatUint64(i uint64) (str string) {
	str = strconv.Itoa(int(i % 10000))
	i /= 10000
	if i > 0 {
		str = "万" + str
	} else {
		return
	}

	str = strconv.Itoa(int(i%10000)) + str
	i /= 10000
	if i > 0 {
		str = "亿" + str
	} else {
		return
	}
	str = strconv.Itoa(int(i)) + str

	return
}

func formatTime(second int64) string {
	m := second / 60
	second = second % 60

	h := m / 60
	m = m % 60

	return fmt.Sprintf("%2d时%2d分%2d秒", h, m, second)
}

// asc的byte转换为unicode的byte
func asc2Unicode(b []byte) []byte {
	var bb []byte = make([]byte, len(b)*2)
	for i := range b {
		bb[i*2] = b[i]
		bb[i*2+1] = 0
	}
	return bb
}

func unicode2Asc(b []byte) (bb []byte, err error) {
	var iLen int = len(b)
	if iLen%2 != 0 {
		return nil, errors.New("输入的unicode byte长度必须是整数。")
	}

	bb = make([]byte, iLen/2)
	for i := 0; i < iLen; i += 2 {
		bb[i/2] = b[i]
	}
	return bb, nil
}

func getCount(m, n int) (icount int) {
	icount = 1
	for i := 0; i < n; i++ {
		icount *= m
	}
	return
}
