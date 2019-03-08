package main

import (
	"fmt"
	"os"
	"regexp"
	"runtime"
	"strconv"
	"strings"

	"github.com/opesun/goquery"
)

var strSep string = string(os.PathSeparator)
var saveTxt string
var iCount int = 0
var f *os.File // 保存的文件，只需要写入网页代码

func init() {
	if runtime.GOOS == "darwin" {
		saveTxt = os.Getenv("HOME") + strSep + `Desktop` + strSep + `pic.html`
	} else if runtime.GOOS == "windows" {
		saveTxt = os.Getenv("USERPROFILE") + strSep + `Desktop` + strSep + `pic.html`
	}
}

func main() {
	fmt.Println("downPic <点赞num> <from> <to>")
	if len(os.Args) != 4 {
		return
	}
	nimNum, err := strconv.Atoi(os.Args[1])
	if err != nil {
		fmt.Println(err)
		return
	}
	ifrom, err := strconv.Atoi(os.Args[2])
	if err != nil {
		fmt.Println(err)
		return
	}
	ito, err := strconv.Atoi(os.Args[3])
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(saveTxt)
	os.Remove(saveTxt)
	var err1 error
	if f, err1 = os.OpenFile(saveTxt, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644); err1 != nil {
		fmt.Println(err1)
		return
	}
	defer f.Close()

	if _, err1 = f.WriteString(`<p>` + os.Args[2] + `-` + os.Args[3] + `</p>`); err1 != nil {
		fmt.Println(err1)
		return
	}
	// 记录下上次的页面

	var url string = `https://www.haha.mx/topic/1/new/`

	for j := ifrom; j <= ito; j++ {
		p, err := goquery.ParseUrl(url + strconv.Itoa(j))
		if err != nil {
			fmt.Println(err)
			return
		}
		t := p.Find(".joke-main-img-suspend")
		t2 := p.Find(".joke-list-item-footer")

		for i := 0; i < t.Length(); i++ {
			d := t.Eq(i).Attr("data-original")
			d2 := t2.Eq(i).Html()
			if num, err := getNum(d2); err != nil {
				fmt.Println(err)
				return
			} else {
				if num >= nimNum {
					if _, err1 = f.WriteString(fmt.Sprintf("\r\n\r\n<p><a href=\"https://www.haha.mx/topic/1/new/%d\">%2d-%d</a></p><img src=\"https:%s\">", j, iCount, num, strings.Replace(d, "normal", "middle", -1))); err1 != nil {
						fmt.Println(err1)
						return
					}
					fmt.Printf("\rOK%4d", iCount)
					iCount++
				}

			}
		}
	}

	fmt.Println()
}

func getNum(strHtml string) (num int, err error) {
	var reg *regexp.Regexp
	reg, err = regexp.Compile(`<a href="javascript:;" title="称赞" class="btn-icon-good" data="g">(\d{1,999})</a>`)
	if err != nil {
		return
	}
	regResult := reg.FindAllStringSubmatch(strHtml, -1)
	if len(regResult) > 0 {
		return strconv.Atoi(regResult[0][1])
	}

	return
}
