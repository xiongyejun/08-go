# go

# 编译压缩体积 #

	go build -ldflags "-s -w"
	-s 去掉符号信息， -w 去掉DWARF调试信息，得到的程序就不能用gdb调试了

# 运行隐藏黑窗口 #
	go build -ldflags "-H windowsgui" project.go

## 编译exe添加icon ##

	rsrc.exe -manifest ico.manifest -o myapp.syso -ico myapp.ico

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