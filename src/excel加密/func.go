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
	return uint32(src[0]) + uint32(src[1]<<8) + uint32(src[2]<<16) + uint32(src[3]<<24), nil
}
func byteToUint16(src []byte) (x uint16, err error) {
	if len(src) != 2 {
		return 0, errors.New("转uint32必须是2个字节。")
	}
	return uint16(src[0]) + uint16(src[1]<<8), nil
}

func uintToByte(i uint) []byte {
	tmp := int32(i)
	buf := bytes.NewBuffer([]byte{})
	binary.Write(buf, binary.LittleEndian, tmp)
	return buf.Bytes()
}

func append36(src []byte, iSize int) []byte {
	if len(src) < iSize {
		// padded by appending bytes with a value of 0x36
		var b []byte = make([]byte, iSize-len(src))
		for j := range b {
			b[j] = 0x36
		}
		src = append(src, b...)
	} else {
		src = src[:iSize]
	}
	return src
}
