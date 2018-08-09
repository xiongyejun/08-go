package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"

	"github.com/opesun/goquery"
)

var saveTxt string = os.Getenv("USERPROFILE") + `\Desktop\pic.html`
var iCount int = 0
var f *os.File // 保存的文件，只需要写入网页代码

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
	os.RemoveAll(saveTxt)
	var err1 error
	if f, err1 = os.OpenFile(saveTxt, os.O_APPEND|os.O_CREATE, 0666); err1 != nil {
		fmt.Println(err1)
		return
	}
	f.WriteString(`<p>` + os.Args[2] + `-` + os.Args[3] + `</p>`)
	// 记录下上次的页面

	var url string = `https://www.haha.mx/topic/1/new/`

	for j := ifrom; j <= ito; j++ {
		p, err := goquery.ParseUrl(url + strconv.Itoa(j))
		if err != nil {
			fmt.Println(err)
			return
		}
		t := p.Find(".joke-main-img")
		t2 := p.Find(".joke-list-item-footer")

		for i := 0; i < t.Length(); i++ {
			d := t.Eq(i).Attr("src")
			d2 := t2.Eq(i).Html()
			if num, err := getNum(d2); err != nil {
				fmt.Println(err)
				return
			} else {
				if num >= nimNum {
					f.WriteString(fmt.Sprintf("<p><a href=\"https://www.haha.mx/topic/1/new/%d\">%2d-%d</a></p><img src=\"https:%s\">", j, iCount, num, d))
					fmt.Printf("\rOK%4d", iCount)
					iCount++
				}

			}
		}
	}

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
