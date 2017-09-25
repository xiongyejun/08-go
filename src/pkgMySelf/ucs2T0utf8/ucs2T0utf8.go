// http://liangguanhui.iteye.com/blog/558849
// 1 所有以0开始的字节，都与原来的ASCII码兼容，也就是说，0xxxxxxx不需要额外转换，就是我们平时用的ASCII码
// 2 所有以10开始的字节，都不是每个UNICODE的第一个字节，都是紧跟着前一位。
//   例如：10110101，这个字节不可以单独解析，必须通过前一个字节来解析
//   如果前一个也是10开头，就继续前嗍
// 3 所有以11开始的字节，都表示是UNICODE的第一个字节，而且后面紧跟着若干个以10开头的字节。
//   如果是110xxxxx（就是最左边的0的左边有2个1），代表后面还有1个10xxxxxx；
//   如果是1110xxxx（就是最左边的0的左边有3个1），代表后面还有2个10xxxxxx；以此类推，一直到1111110x

//1字节 0xxxxxxx
//2字节 110xxxxx 10xxxxxx
//3字节 1110xxxx 10xxxxxx 10xxxxxx
//4字节 11110xxx 10xxxxxx 10xxxxxx 10xxxxxx
//5字节 111110xx 10xxxxxx 10xxxxxx 10xxxxxx 10xxxxxx
//6字节 1111110x 10xxxxxx 10xxxxxx 10xxxxxx 10xxxxxx 10xxxxxx

// 对于UCS-2，仅仅有2 ^ 16个字符，只需要三个字节就可以

package ucs2T0utf8

import (
	"bytes"
	"errors"
)

const (
	b_1000_0000 = 128
	b_1100_0000 = 192
	b_1110_0000 = 224
	b_1111_0000 = 240
	b_0001_1100 = 28
	b_0000_0011 = 3
	b_0011_1111 = 63
	b_0000_1111 = 15
	b_0011_1100 = 60
)

// UCS-2转UTF-8
// 1 对于不大于0x007F（即00000000 01111111）的，直接把它转成一个字节，变成ASCII
// 2 对于不大于0x07FF（即00000111 11111111）的，转换成两个字节
//   转换的时候把右边的11位分别放到110xxxxx 10yyyyyy里边
//   即00000aaa bbbbbbbb ==> 110aa abb   10bb bbbb
// 3 剩下的回转换成三个字节，转换的时候也是把16个位分别填写到那三个字节里面
//   即aaaaaaaa bbbbbbbb ==> 1110 aaaa 	 10aa aabb   10bb bbbb

func UCS2toUTF8(bUCS2 []byte) ([]byte, error) {
	if len(bUCS2)%2 > 0 {
		return []byte{}, errors.New("err:输入的UCS2字节数组不是偶数！")
	}

	buf := bytes.NewBuffer([]byte{})

	for i := 0; i < len(bUCS2); i += 2 {
		ub1 := uint16(bUCS2[i+1])
		ub2 := uint16(bUCS2[i])

		var b uint16 = (ub1 << 8) | ub2 // 读2个

		if b <= 0x007f {
			buf.WriteByte(byte(b))
		} else if b <= 0x07ff {
			b1 := b_1100_0000 + (ub1 << 2) + (ub2 >> 6)
			b2 := b_1000_0000 + (ub2 & b_0011_1111)
			buf.WriteByte(byte(b1))
			buf.WriteByte(byte(b2))
		} else {
			b1 := b_1110_0000 + (ub1 >> 4)
			b2 := b_1000_0000 + ((ub1 & b_0000_1111) << 2) + (ub2 >> 6)
			b3 := b_1000_0000 + (ub2 & b_0011_1111)
			buf.WriteByte(byte(b1))
			buf.WriteByte(byte(b2))
			buf.WriteByte(byte(b3))
		}
	}
	return buf.Bytes(), nil
}
