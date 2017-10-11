// 仿unix里的ls命令
package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"pkgMySelf/colorPrint"
	"time"
)

type ls struct {
	dir          string
	subDir       bool
	fScanDir     func(d string, entryName string)
	fullName     bool
	fGetFileName func(d string, entryName string) string
	sep          string
	numDir       int32
	numFile      int32

	chanDir  chan string  // 控制搜索
	chanFile chan outType // 控制输出

	cd          *colorPrint.ColorDll
	dicExtColor map[string]uintptr
}

type outType struct {
	isDir bool
	name  string
}

var l ls

//var chanEnd chan int

func main() {
	// 判断dir是否存在
	finfo, err := os.Stat(l.dir)
	if err != nil {
		fmt.Println("不存在的文件夹。")
		return
	}
	if !finfo.IsDir() {
		fmt.Println(l.dir, " 不是文件夹。")
		return
	}

	l.chanDir = make(chan string, 100)
	l.chanFile = make(chan outType, 1000)

	go l.scanDir(l.dir)

	l.cd = colorPrint.NewColorDll()
	l.initExtColor()

	go l.printOut()
	time.Sleep(1e8)
	for len(l.chanDir) != 0 || len(l.chanFile) != 0 {
		time.Sleep(1e8)
	}
	fmt.Printf("dir Count = %d\r\nfile Count = %d\r\n", l.numDir, l.numFile)
}

func init() {
	str, _ := os.Getwd() // 获得cmd命令行cd的路径
	var strDir = flag.String("d", str, "scan dir path")
	var subDir = flag.Bool("s", false, "scan sub dir")
	var fullName = flag.Bool("b", false, "full name")

	flag.PrintDefaults()
	flag.Parse()

	l = ls{dir: *strDir, subDir: *subDir, fullName: *fullName, sep: string(os.PathSeparator)}
	if string(l.dir[len(l.dir)-1]) != l.sep {
		l.dir = l.dir + l.sep
	}
	// 在这里判断是否要遍历子文件夹
	if l.subDir {
		l.fScanDir = scanSubDir
	} else {
		l.fScanDir = scanNoSubDir
	}
	// 在这里判断是否要带路径的文件名
	if l.fullName {
		l.fGetFileName = getFullName
	} else {
		l.fGetFileName = getName
	}

	//	fmt.Printf("%#v\r\n", l)
}

func (this *ls) scanDir(dirName string) {
	this.chanDir <- dirName
	defer func() {
		<-this.chanDir
	}()

	entrys, err := ioutil.ReadDir(dirName)
	if err != nil {
		return
	}
	outtype := outType{}
	for _, entry := range entrys {
		outtype.isDir = false
		if entry.IsDir() {
			outtype.isDir = true
			this.fScanDir(dirName, entry.Name())
		}
		outtype.name = this.fGetFileName(dirName, entry.Name())
		this.chanFile <- outtype
	}
}

func scanSubDir(d string, entryName string) {
	go l.scanDir(d + entryName + l.sep)
}
func scanNoSubDir(d string, entryName string) {

}

func getFullName(d string, entryName string) string {
	return d + entryName
}
func getName(d string, entryName string) string {
	return entryName
}

func (this *ls) printOut() {
	for f := range this.chanFile {
		if f.isDir {
			this.numDir++
			this.cd.SetColor(colorPrint.White, colorPrint.DarkYellow)
		} else {
			this.numFile++
			strExtension := path.Ext(f.name)
			if v, ok := this.dicExtColor[strExtension]; ok {
				this.cd.SetColor(colorPrint.White, v)
			}
		}

		fmt.Printf("%s", f.name)
		this.cd.UnSetColor()
		fmt.Printf("\r\n") // 回车要这里输，在前面输了下一行的空白也有颜色，不知道为什么
	}
}

func (this *ls) initExtColor() {
	this.dicExtColor = make(map[string]uintptr)
	this.dicExtColor[".xls"] = colorPrint.DarkMagenta
	this.dicExtColor[".xlsm"] = colorPrint.DarkMagenta
	this.dicExtColor[".xlsx"] = colorPrint.DarkMagenta

	this.dicExtColor[".doc"] = colorPrint.Blue
	this.dicExtColor[".docx"] = colorPrint.Blue
	this.dicExtColor[".docm"] = colorPrint.Blue

	this.dicExtColor[".txt"] = colorPrint.DarkGreen
	this.dicExtColor[".exe"] = colorPrint.DarkRed
	this.dicExtColor[".go"] = colorPrint.Red
}
