package main

import (
	"bytes"
	"encoding/binary"
	"errors"
)

func byteToUint32(src []byte) (x uint32, err error) {
	if len(src) != 4 {
		return 0, errors.New("转uint32必须是4个字节。")
	}
	return uint32(src[0]) | uint32(src[1])<<8 | uint32(src[2])<<16 | uint32(src[3])<<24, nil
}
func byteToUint16(src []byte) (x uint16, err error) {
	if len(src) != 2 {
		return 0, errors.New("转uint32必须是2个字节。")
	}
	return uint16(src[0]) | uint16(src[1])<<8, nil
}

func uintToByte(i uint) []byte {
	tmp := int32(i)
	buf := bytes.NewBuffer([]byte{})
	binary.Write(buf, binary.LittleEndian, tmp)
	return buf.Bytes()
}

func appendByte(src []byte, iSize int, bb byte) []byte {
	if len(src) < iSize {
		// padded by appending bytes with a value of 0x36
		var b []byte = make([]byte, iSize-len(src))
		for j := range b {
			b[j] = bb
		}
		src = append(src, b...)
	} else {
		src = src[:iSize]
	}
	return src
}

func string2Unicode(str string) []byte {
	b := []byte(str)
	var bb []byte = make([]byte, len(b)*2)
	for i := range b {
		bb[i*2] = b[i]
		bb[i*2+1] = 0
	}
	return bb
}
