[https://studygolang.com/articles/11453](https://studygolang.com/articles/11453)

1. 在golang代码开始部分(package xxx之后)，添加注释，注释中编写需要使用的C语言代码
1. 紧挨着注释结束，另起一行增加import "C"，注意跟注释中的C代码紧挨，不要有空行，且不要跟其他golang的import放在一起
1. 这样在golang语言的正文中就可以用C.xxx的方式调用注释中的C代码了

    	package main

		// #include <stdio.h>
		// #include <stdlib.h>
		/*
		void print(char *s) {
		    printf("print used by C language:%s\n", s);
		}
		*/
		import "C" //和上一行"*/"直接不能有空行或其他注释
		
		import "unsafe"
		
		func main() {
		    s := "hello"
		    cs := C.CString(s)
		    defer C.free(unsafe.Pointer(cs))
		    C.print(cs)
		}


**原理**

其实cgo就是先由编译器识别出import "C"的位置，然后在其上的注释中提取C代码，最后调用C编译器进行分开编译

**使用cgo要点我觉得有两个**

golang和C直接的类型转换
静态库和动态库的链接
补充一点当时业务中遇到的问题，要链接的动态库文件，不知道相对路径怎么取，后来找到了解决办法：${SRCDIR}
	
**When the cgo directives are parsed, any occurrence of the string ${SRCDIR} will be replaced by the absolute path to the directory containing the source file.**

		package aes
		
		/*
		#cgo LDFLAGS: -L${SRCDIR} -lyourfile -ldl
		#include <stdio.h>
		#include <stdlib.h>
		#include "yourcode.h"
		*/
		import "C"
		
		......


[http://cholerae.com/2015/05/17/%E4%BD%BF%E7%94%A8Cgo%E7%9A%84%E4%B8%80%E7%82%B9%E6%80%BB%E7%BB%93/](http://cholerae.com/2015/05/17/%E4%BD%BF%E7%94%A8Cgo%E7%9A%84%E4%B8%80%E7%82%B9%E6%80%BB%E7%BB%93/ "使用Cgo的一点总结")