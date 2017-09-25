package compdocFile

import (
	"bytes"
)

// 复合文档头部512个byte
type cfHeader struct {
	Id                   [8]byte
	File_id              [16]byte
	File_format_revision int16
	File_format_version  int16
	Memory_endian        int16
	Sector_size          int16 // '扇区的大小 2的幂 通常为2^9=512
	Short_sector_size    int16
	Not_used_1           [10]byte
	Sat_count            int32 //'分区表扇区的总数
	Dir_first_sid        int32
	Not_used_2           [4]byte
	Min_stream_size      int32
	Ssat_first_sid       int32
	Ssat_count           int32
	Msat_first_sid       int32
	Msat_count           int32
	Arr_sid              [109]int32
}

type cfDir struct {
	Dir_name    [64]byte
	Len_name    int16
	CfType      byte  // 1仓storage 2流 5根
	Color       byte  // 0红色 1黑色
	Left_child  int32 // -1表示叶子
	Right_child int32
	Sub_dir     int32
	Arr_keep    [20]byte
	Time_create [8]byte
	Time_modify [8]byte
	First_SID   int32 // 目录入口所表示的第1个扇区编码
	Stream_size int32 // 目录入口流尺寸，可判断是否是短扇区
	Not_used    int32
}

type cfStruct struct {
	fileByte   []byte   // 文件的byte
	header     cfHeader // 文件头部512个字节
	byteHeader []byte   // rc io.ReadCloser读取头结构获取文件大小时，向前移动了位置，要保留下来

	arrMSAT   []int32     // 主分区表
	arrSAT    []int32     // 分区表
	arrDir    []cfDir     // 目录
	arrStream []*cfStream // 目录对应的流
	arrSSAT   []int32     // 短分区表

	arrDirInfo []dirInfo
	dic        map[string]int32
	dicModule  map[string]int32
}

type cfStream struct {
	stream  *bytes.Buffer // 流的信息
	step    int32         // 短流是64，正常的是512，如果是0就不是流
	address []int32       // 记录每个地址的开始	，也就是记录arrSAT或者arrSSAT
	name    string
}
