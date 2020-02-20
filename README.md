# go

# 编译压缩体积 #

	go build -ldflags "-s -w"
	-s 去掉符号信息， -w 去掉DWARF调试信息，得到的程序就不能用gdb调试了

# 运行隐藏黑窗口 #
	go build -ldflags "-H windowsgui" project.go

## 编译exe添加icon ##

	rsrc.exe -manifest ico.manifest -o myapp.syso -ico myapp.ico

# 查看汇编代码 #

	go tool compile -S main.go >> main.S


# 学习资源 #

- [A golang ebook intro how to build a web with golang](https://github.com/astaxie/build-web-application-with-golang)
- [Go语言圣经中文版](https://github.com/gopl-zh/gopl-zh.github.com)
- [《The Way to Go》中文译本，中文正式名《Go入门指南》](https://github.com/Unknwon/the-way-to-go_ZH_CN)
- [微信](https://github.com/liushuchun/wechatcmd)
- [网页版微信API，包含终端版微信及微信机器人](https://github.com/Urinx/WeixinBot)
- [win api](https://github.com/lxn/win "win API")
- [ui3 win](https://github.com/lxn/walk "https://github.com/lxn/walk")
- [ui1](https://github.com/visualfc/goqt "UI")
- [ui2](https://github.com/google/gxui "https://github.com/google/gxui")
- [二维码](https://github.com/skip2/go-qrcode "https://github.com/skip2/go-qrcode")
- [编码转换](github.com/axgle/mahonia)
- [带附件mail](https://github.com/scorredoira/email)
- [excel](https://github.com/aswjh/excel)
- [go-sqlite3](http://godoc.org/github.com/mattn/go-sqlite3 "go-sqlite3")
https://github.com/andlabs/ui
- [鼠标键盘截图等](https://github.com/go-vgo/robotgo "鼠标键盘截图等")
- [用Go开发可以内网活跃主机嗅探器](https://studygolang.com/articles/11517 "用Go开发可以内网活跃主机嗅探器")
- [Go标准库所有方法使用例子](https://github.com/zc2638/go-standard "Go标准库所有方法使用例子")
- [免费的编程中文书籍索引](https://github.com/justjavac/free-programming-books-zh_CN "免费的编程中文书籍索引")
- [go知识图谱](https://www.processon.com/view/link/5a9ba4c8e4b0a9d22eb3bdf0#map "go知识图谱")
- [go操作excel](https://github.com/360EntSecGroup-Skylar/excelize "go操作excel")
- [go操作office](https://github.com/unidoc/unioffice "go操作office")
- [算法学习等](https://github.com/studygolang/leetcode "算法学习等")
- [生成pdf](https://github.com/tiechui1994/gopdf "生成pdf")
