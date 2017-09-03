// 将md文件转换为html文件
package main

import (
	"bufio"
	//	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

var f_html *os.File
var w_html *bufio.Writer

func main() {
	if len(os.Args) == 0 {
		return
	}
	str_dir := os.Args[1]
	save_dir := str_dir + "_html"
	os.Mkdir(save_dir, 0777)
	fmt.Println("ok")

	entrys, _ := ioutil.ReadDir(str_dir)
	for _, en := range entrys {
		if !en.IsDir() {
			fmt.Println(en.Name())
			f_html, err := os.Create(save_dir + "\\" + en.Name() + ".html") // OpenFile	, os.O_RDWR|os.O_CREATE, 0644)
			w_html = bufio.NewWriter(f_html)
			if err != nil {
				fmt.Println("create file err")
			}

			fmt.Fprintln(w_html, "<meta charset=\"UTF-8\">")
			fmt.Fprintln(w_html, "<style type=\"text/css\"> body {line-height:200%} </style>")
			readLine(str_dir+"\\"+en.Name(), md2Html)
			w_html.Flush()
			f_html.Close()
			//			break
		}
	}
}

// 将md的str转化为html的str
func md2Html(strMD string) string {
	runestr := []rune(strMD)

	const str_pre string = "![](../images"
	if strings.HasPrefix(strMD, str_pre) {
		return fmt.Sprintf("<img src=\"%s\" />", string(runestr[len("![]("):]))
	}

	if strings.Index(strMD, "](") > 0 && strings.Index(strMD, ".md") > 0 {
		//		index := strings.Index(strMD, "(")
		linkName := strings.Split(strings.Split(strMD, "(")[1], ")")[0]
		return fmt.Sprintf("<li><a href=\"%s.html\">%s</a></li>", linkName, strMD)
	}

	if string(runestr[0:1]) == "#" {
		var i int
		for i = 2; i < len(runestr); i++ {
			if string(runestr[i:(i+1)]) != "#" {
				break
			}
		}
		return fmt.Sprintf("<h%d>%s</h%d>\n", i-1, string(runestr[i:]), i-1)
	}

	if string(runestr[0:1]) == "-" {
		return fmt.Sprintf("<li>%s</li>\n", string(runestr[2:]))
	}

	return fmt.Sprintf("<p>%s</p>", strMD)
}

func readLine(fileName string, handler func(string) string) error {
	f, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer f.Close()

	buf := bufio.NewReader(f)
	for {
		var line string
		var err error
		line, err = buf.ReadString('\n')
		//		line = strings.TrimSpace(line)

		var s string
		if strings.HasPrefix(line, "```") {
			s = strings.Replace(line, "\t", "&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;", -1)
			s = fmt.Sprintf("<span style=\"color:#ee1b2e;\"><ol><li>%s<br/></li>\n", s)
			line = ""
			for {
				//				fmt.Println(line)
				line = strings.Replace(line, "\t", "&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;", -1)
				s = fmt.Sprintf("%s<li>%s<br/></li>\n", s, line)
				line, err = buf.ReadString('\n')
				//				line = strings.TrimSpace(line)
				if strings.Index(line, "```") > -1 || err == io.EOF {
					break
				}
			}
			s = fmt.Sprintf("%s</ol></span>\n", s)
		} else {
			s = handler(line)
		}
		fmt.Fprintln(w_html, s)

		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
	}
	return nil
}
