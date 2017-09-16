package compdocFile

import (
	"bytes"
)

// 复合文档头部512个byte
type cfHeader struct {
	id                   [8]byte
	file_id              [16]byte
	file_format_revision int16
	file_format_version  int16
	memory_endian        int16
	sector_size          int16 // '扇区的大小 2的幂 通常为2^9=512
	short_sector_size    int16
	not_used_1           [10]byte
	sat_count            int32 //'分区表扇区的总数
	dir_first_sid        int32
	not_used_2           [4]byte
	min_stream_size      int32
	ssat_first_sid       int32
	ssat_count           int32
	msat_first_sid       int32
	msat_count           int32
	arr_sid              [109]int32
}

type cfDir struct {
	dir_name    [64]byte
	len_name    int16
	cfType      byte  // 1仓storage 2流 5根
	color       byte  // 0红色 1黑色
	left_child  int32 // -1表示叶子
	right_child int32
	sub_dir     int32
	arr_keep    [20]byte
	time_create [8]byte
	time_modify [8]byte
	first_SID   int32 // 目录入口所表示的第1个扇区编码
	stream_size int32 // 目录入口流尺寸，可判断是否是短扇区
	not_used    int32
}

type cfStruct struct {
	fileByte []byte   // 文件的byte
	header   cfHeader // 文件头部512个字节

	arrMSAT   []int32     // 主分区表
	arrSAT    []int32     // 分区表
	arrDir    []cfDir     // 目录
	arrStream []*cfStream // 目录对应的流
	arrSSAT   []int32     // 短分区表
}

type cfStream struct {
	stream  *bytes.Buffer // 流的信息
	step    int32         // 短流是64，正常的是512，如果是0就不是流
	address []int32       // 记录每个地址的开始	，也就是记录arrSAT或者arrSSAT
}
