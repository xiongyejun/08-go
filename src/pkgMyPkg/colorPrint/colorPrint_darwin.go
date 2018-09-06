// 打印彩色
package colorPrint

import (
	"fmt"
)

const (
	Black uintptr = iota
	Red
	Green
	Yellow
	Blue
	DarkRed
	DarkBlue
	White

	DarkGreen   = Green
	DarkCyan    = Blue
	DarkMagenta = Blue
	DarkYellow  = Yellow
	Gray        = Blue
	DarkGray    = Blue
	Cyan        = Blue
	Magenta     = Blue

	// 前景 背景 颜色
	// ---------------------------------------
	// 30  40  黑色
	// 31  41  红色
	// 32  42  绿色
	// 33  43  黄色
	// 34  44  蓝色
	// 35  45  紫红色
	// 36  46  青蓝色
	// 37  47  白色

	// 代码 意义
	// -------------------------
	//  0  终端默认设置
	//  1  高亮显示
	//  4  使用下划线
	//  5  闪烁
	//  7  反白显示
	//  8  不可见
)

//fmt.Printf("\n %c[1;40;32m%s%c[0m\n\n", 0x1B, "testPrintColor", 0x1B)
//其中0x1B是标记，[开始定义颜色，1代表高亮，40代表黑色背景，32代表绿色前景，0代表恢复默认颜色
func SetColor(fontColor uintptr, backGroundColor uintptr) {
	fmt.Printf("%c[1;%d;%dm%c", 0x1b, fontColor+30, backGroundColor+40, 0x1b)
}

func UnSetColor() {
	fmt.Printf("%c[0m%c", 0x1b, 0x1b)
}
