// 根据输入的参数来做具体的事
package main

import (
	"flag"
	"fmt"
	"pkgAPI/comdlg32"
	"pkgMySelf/colorPrint"
	"pkgMySelf/compdocFile"
)

type todo struct {
	file       *string
	moduleName *bool   // -m	打印模块名称
	code       *bool   // -c	打印模块代码
	project    *bool   // -p	破解工程密码
	hideModule *string // -h	隐藏某个模块
}

var td *todo

func init() {
	td = new(todo)

	td.file = flag.String("f", "", "文件名称")
	td.moduleName = flag.Bool("m", false, "打印模块名称")
	td.code = flag.Bool("c", false, "打印模块代码")
	td.project = flag.Bool("p", false, "破解工程密码")
	td.hideModule = flag.String("h", "", "隐藏某个模块")

	flag.PrintDefaults()
	flag.Parse()
}

func main() {
	cd := colorPrint.NewColorDll()

	if !(*td.moduleName || *td.code || *td.project) && *td.hideModule == "" {
		return
	}

	if *td.file == "" {
		fd := comdlg32.NewFileDialog()
		b, _ := fd.GetOpenFileName()
		if !b {
			return
		}

		*td.file = fd.FilePath
		fmt.Println(*td.file)
	}

	var cf compdocFile.CF
	if compdocFile.IsCompdocFile(*td.file) {
		cf = compdocFile.NewXlsFile(*td.file)
	} else if compdocFile.IsZip(*td.file) {
		cf = compdocFile.NewZipFile(*td.file)
	} else {
		cd.SetColor(colorPrint.White, colorPrint.DarkMagenta)
		fmt.Println("未知文件：", *td.file)
		cd.UnSetColor()
		return
	}

	err := compdocFile.CFInit(cf)
	if err != nil {
		fmt.Println(err)
		return
	}

	if *td.moduleName {
		cd.SetColor(colorPrint.White, colorPrint.DarkMagenta)
		modules := cf.GetModuleName()
		for _, v := range modules {
			fmt.Println(v)
		}
	}

	if *td.code {
		cf.PrintAllCode()
	}

	if *td.project {
		err := cf.UnProtectProject()
		cd.SetColor(colorPrint.White, colorPrint.DarkMagenta)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("破解工程密码成功。")
		}
	}

	if *td.hideModule != "" {
		fmt.Println("hide")
		err := cf.HideModule(*td.hideModule)
		cd.SetColor(colorPrint.White, colorPrint.DarkMagenta)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("隐藏模块成功。")
		}
	}
	cd.UnSetColor()
}
