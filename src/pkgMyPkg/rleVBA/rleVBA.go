// run length encoding
// To preserve space, VBA uses data compression on a contiguous
// sequence of records on various streams.
// The data compression technique is run length encoding

// VBA代码的压缩：
// 	  每4096个byte作为一个chunk进行压缩，然后在前面添加1个2byte的head
//    head里记录了压缩后的byte大小及压缩标志（有可能是没有进行压缩的数据）
//        如果是压缩的，在压缩数据块里，第1个byte是个标识，指示后面8个copytoken是否是经过了压缩的
//        没有压缩的话就直接读取1个byte
//        如果是压缩了的，用2个字节copytoken来记录当前压缩的数据怎么来读取
//            从copytoken里读取出offset和len，这里是相对解压缩后的数据的offset和len
//            根据offset和len来读取解压缩后的数据
package rleVBA

import (
	"bytes"
	"math"
)

const UnCompressChunkSize = 4096 // 未压缩块的大小、解压出来的块大小，每4096个byte进行压缩

type rle struct {
	compressByte         []byte // 压缩的字节
	compressed_Current   int    // 指向压缩字节正在解压的位置
	compressed_RecordEnd int    // 最后位置，也就是compressByte的长度-1

	deCompressBuffer        *bytes.Buffer // 解压后的buffer
	deCompressed_Current    int           // 将被写入Decompressed_Buffer的位置
	deCompressedChunk_Start int           // 解压缩数据块开始写入deCompressBuffer的位置

	compressed_Flag       int16 // 数据块压缩标识，1是压缩，0是没有压缩
	compressedChunk_Start int   // 压缩数据块开始的地方
	compressedChunk_Size  int   // 压缩数据块的大小
	compressedChunk_End   int   // 压缩数据块结束的地方
}

func (me *rle) UnCompress() []byte {
	// SignatureByte 压缩标识为0x1才是压缩过的
	if me.compressByte[0] != 0x1 {
		return me.compressByte
	}

	me.compressed_Current++
	for me.compressed_Current < me.compressed_RecordEnd {
		me.compressedChunk_Start = me.compressed_Current
		me.deCompressingCompressedChunk()
	}
	return me.deCompressBuffer.Bytes()
}

// 解压缩 压缩块
func (me *rle) deCompressingCompressedChunk() {
	var Compressed_Header uint16
	// 每个输出块前面都有一个两个字节的头，表示块中的字节数和块的格式。
	// 每个压缩块被解码成4096字节的未压缩数据，被写入输出。
	// 对于每个块，从块标题中提取大小和格式样式。然后根据标题中指定的格式读取和解码该块
	Compressed_Header = (uint16(me.compressByte[me.compressed_Current+1]) << 8) | uint16(me.compressByte[me.compressed_Current])
	// 获得压缩数据块的大小
	me.compressedChunk_Size = me.getCompressedChunkSize(Compressed_Header)
	// 获取数据块压缩标识，1是压缩，0是没有压缩
	me.compressed_Flag = me.getCompressedFlag(Compressed_Header)
	// 压缩数据块的最后位置
	if me.compressed_RecordEnd > (me.compressedChunk_Start + me.compressedChunk_Size) {
		me.compressedChunk_End = me.compressedChunk_Start + me.compressedChunk_Size
	} else {
		me.compressedChunk_End = me.compressed_RecordEnd
	}
	me.compressed_Current = me.compressedChunk_Start + 2

	if me.compressed_Flag == 1 {
		for me.compressed_Current < me.compressedChunk_End {
			me.deCompressingTokenSequence()
		}
	} else {
		// 未压缩的块，直接读取
		me.deCompressingRawChunk()
	}

	me.compressedChunk_Start = me.compressed_Current
	me.deCompressedChunk_Start = me.deCompressed_Current
}

// 获取压缩数据块的大小
func (me *rle) getCompressedChunkSize(HeaderByte uint16) int {
	return 3 + int((HeaderByte & 0xfff))
}

//获取数据块压缩标识，最后1个bit，1是压缩，0是没有压缩
func (me *rle) getCompressedFlag(HeaderByte uint16) int16 {
	flag := HeaderByte & 0x8000 // 1000 0000 0000 0000
	if flag == 0x8000 {
		return 1
	} else {
		return 0
	}
}

// 直接读取未压缩的4096个数据
// contain 4096 bytes of uncompressed data
func (me *rle) deCompressingRawChunk() {
	i_end := me.compressed_Current + UnCompressChunkSize
	if i_end > me.compressed_RecordEnd {
		i_end = me.compressed_RecordEnd
	}

	me.deCompressBuffer.Write(me.compressByte[me.compressed_Current:i_end])
	me.compressed_Current += (i_end - me.compressed_Current)
	me.deCompressed_Current += (i_end - me.compressed_Current)
}

// 解压缩TokenSequence
func (me *rle) deCompressingTokenSequence() {
	// flagByte的8位对应了8个Tokens
	// 0表示没有压缩，1表示是1个copyToken
	flagByte := me.compressByte[me.compressed_Current]
	me.compressed_Current++

	var i uint = 0
	for ; i < 8; i++ {
		if me.compressed_Current < me.compressedChunk_End { // 有可能没有8个token
			// CALL Decompressing a Token (section 2.4.1.3.5) with index and Byte
			me.deCompressingToken(i, flagByte)
		}
	}
}

// 解压缩Token，正常1个TokenSequence包含8个token
func (me *rle) deCompressingToken(TokenSequenceIndex uint, FlagByte byte) {
	// 读取FlagByte的TokenSequenceIndex的bit
	flag := me.getFlagBit(TokenSequenceIndex, FlagByte)

	if flag == 0 {
		// COPY the byte at CompressedCurrent TO DecompressedCurrent
		me.deCompressBuffer.Write(me.compressByte[me.compressed_Current:(me.compressed_Current + 1)])
		me.compressed_Current++
		me.deCompressed_Current++
	} else {
		// CopyToken 2byte，包含的是Offset and Length
		// CALL Unpack CopyToken (section 2.4.1.3.19.2) with Token returning Offset and Length
		var token uint16
		token = (uint16(me.compressByte[me.compressed_Current+1]) << 8) | uint16(me.compressByte[me.compressed_Current])
		Offset, Length := me.unpackCopyToken(token)
		// SET CopySource TO DecompressedCurrent - Offset
		// CALL Byte Copy (section 2.4.1.3.11) with CopySource, DecompressedCurrent, and Length
		var i_start int = me.deCompressed_Current - int(Offset)
		var i_end int = me.deCompressed_Current - int(Offset) + int(Length)

		for i := i_start; i < i_end; i++ {
			me.deCompressBuffer.Write(me.deCompressBuffer.Bytes()[i : i+1])
			me.deCompressed_Current++
		}
		//		me.deCompressBuffer.Write(me.deCompressBuffer.Bytes()[i_start:i_end])
		//		me.deCompressed_Current += (i_end - i_start)

		me.compressed_Current += 2
	}
}

// 获取TokenSequence的FlagByte 中第BitIndex的flag
func (me *rle) getFlagBit(bitIndex uint, FlagByte byte) uint {
	// 2.4.1.3.17	Extract FlagBit
	// Index: An unsigned integer specifying which FlagBit to extract. MUST be greater than or equal to zero and less than eight
	// Byte (1 byte): An instance of a FlagByte
	// return-- Flag: An integer. The value of the bit in Byte at location Index. The value returned MUST be zero or one.
	// SET Flag TO (Byte RIGHT SHIFT BY Index) BITWISE AND 1
	return (uint(FlagByte) >> bitIndex) & 1
}

// 从CopyToken中获取Offset和Length
func (me *rle) unpackCopyToken(token uint16) (Offset, Length uint16) {
	// 2.4.1.3.19.2	Unpack CopyToken
	// Offset (2 bytes): An unsigned 16-bit integer that specifies the beginning of a CopySequence (section 2.4.1.3.19).
	// Length (2 bytes): An unsigned 16-bit integer that specifies the length of a CopySequence

	//1.	CALL CopyToken Help (section 2.4.1.3.19.1) returning LengthMask, OffsetMask, and BitCount.
	LengthMask, OffsetMask, BitCount, _ := me.copyTokenHelp()
	//2.	SET Length TO (Token BITWISE AND LengthMask) PLUS 3.
	Length = (token & LengthMask) + 3
	//3.	SET temp1 TO Token BITWISE AND OffsetMask.
	temp1 := token & OffsetMask
	//4.	SET temp2 TO 16 MINUS BitCount.
	temp2 := 16 - BitCount
	//5.	SET Offset TO (temp1 RIGHT SHIFT BY temp2) PLUS 1.
	Offset = (temp1 >> temp2) + 1

	return
}

// CopyToken Help
func (me *rle) copyTokenHelp() (LengthMask, OffsetMask, BitCount, MaximumLength uint16) {
	// LengthMask (2 bytes): An unsigned 16-bit integer. A bitmask used to access CopyToken.Length.
	// OffsetMask (2 bytes): An unsigned 16-bit integer. A bitmask used to access CopyToken.Offset.
	// BitCount (2 bytes): An unsigned 16-bit integer. The number of bits set to 0b1 in OffsetMask.
	// MaximumLength (2 bytes): An unsigned 16-bit integer. The largest possible integral积分value that can fit into CopyToken.Length

	//§	SET difference TO DecompressedCurrent MINUS DecompressedChunkStart
	difference := me.deCompressed_Current - me.deCompressedChunk_Start
	//§	SET BitCount TO the smallest integer that is GREATER THAN OR EQUAL TO LOGARITHM base 2 of difference
	// 大于或者等于log2(different)的最小整数，要向上取整
	BitCount = uint16(math.Ceil(math.Log2(float64(difference))))

	//§	SET BitCount TO the maximum of BitCount and 4
	if BitCount < 4 {
		BitCount = 4
	}

	//§	SET LengthMask TO 0xFFFF RIGHT SHIFT BY BitCount
	LengthMask = 0xffff >> BitCount
	//§	SET OffsetMask TO BITWISE NOT LengthMask
	OffsetMask = ^LengthMask
	//§	SET MaximumLength TO (0xFFFF RIGHT SHIFT BY BitCount) PLUS 3
	MaximumLength = (0xffff >> BitCount) + 3

	return
}

func NewRLE(compressByte []byte) *rle {
	r := new(rle)
	r.compressByte = compressByte[:]

	r.deCompressBuffer = bytes.NewBuffer([]byte{})
	r.compressed_RecordEnd = len(compressByte) - 1

	return r
}
