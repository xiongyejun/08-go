package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

var dic = map[string]string{}
var pid int

func main() {
	pid = os.Getpid()
	getSet()
	openURL("http://localhost:9090/?strURL=" + dic["str_URL"] + "^&dir=1")

	http.HandleFunc("/", handleFunc)
	http.ListenAndServe(":9090", nil)
}

// 读取一些设置
func getSet() {
	// 不需要删除，直接更新就可以了
	//	for k, _ := range dic {
	//		delete(dic, k)
	//	}
	file, _ := exec.LookPath(os.Args[0])
	path := filepath.Dir(file)

	//	fmt.Println(path)

	f, err := os.Open(path + "\\set.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()
	br := bufio.NewReader(f)
	for {
		a, _, c := br.ReadLine()
		if c == io.EOF {
			break
		}
		str := string(a)
		arr := strings.Split(str, ":=")
		if len(arr) == 2 {
			dic[arr[0]] = arr[1]
		}
	}

}

func openURL(str_URL string) {
	cmd := exec.Command("cmd", "/c", "start", str_URL)
	err := cmd.Start()
	if err != nil {
		fmt.Println("err")
	}
}

//http://localhost:9090/?strURL=http://www.biquge3.com/xs/1181206
func handleFunc(w http.ResponseWriter, r *http.Request) {
	str_URL := dic["str_URL"]
	getSet()
	// 如果网址更新了就重新打开
	if str_URL != dic["str_URL"] {
		openURL("http://localhost:9090/?strURL=" + dic["str_URL"])
	}

	r.ParseForm() //解析参数
	ex := r.Form["exit"]
	if len(ex) > 0 {
		os.Exit(0)
	}
	//	fmt.Fprintf(w, "\npath", r.URL.Path, "\nscheme", r.URL.Scheme, r.Form["url_long"])
	strURL := r.Form["strURL"]
	if len(strURL) == 0 {
		fmt.Fprintf(w, "没有url")
		return
	}

	fmt.Fprintf(w, "<!DOCTYPE html>\n<html>\n")
	fmt.Fprintf(w, "<style type=\"text/css\"> body {font-size:24px;line-height:1.5;}</style>\n <body onunload=\"closeGo()\" bgcolor=\"#C7EDCC\">\n ")
	// 插入退出按钮
	fmt.Fprintf(w, exitHtml())
	// 如果请求中含有strURL，就去获取strURL地址的网页源码
	str_html := getHtml(strURL[0])
	isDir := r.Form["dir"]
	if len(isDir) == 0 {
		str_html = getHtmlId(str_html, dic["str_id"]) // "<div id=\"zjneirong\">")
		fmt.Fprintf(w, str_html)
	} else {
		// 如果是目录，就按照正则获取所有需要的章节
		str_patten := dic["str_patten"] //"<a href=\"http:\\/\\/www\\.biquge3\\.com\\/xs\\/1181206/\\d{8,8}\\.htm.*?第.{1,6}.*?<\\/a>"
		reg, _ := regexp.Compile(str_patten)
		submatchall := reg.FindAllString(str_html, -1)
		str_html = strings.Join(submatchall, "<br>")

		// 有的是完整的地址，有的不是
		if strings.Contains(str_html, "http://www.") {
			var str_old string = "href=\"http://www."
			str_html = strings.Replace(str_html, str_old, "href=\"http://localhost:9090/?strURL=http://www.", 99999)
		} else {
			str_html = strings.Replace(str_html, "href=\"", "href=\"http://localhost:9090/?strURL="+strURL[0], 99999)
		}
		fmt.Fprintf(w, str_html, '\n')
	}
	fmt.Fprintf(w, "</body>\n")
	fmt.Fprintf(w, jsExecCmd(fmt.Sprintf("TaskKill /PID %d", pid)))
	fmt.Fprintf(w, "</html>")
}

func getHtmlId(strHtml string, id string) string {
	if id == "" {
		return strHtml
	}

	arr := strings.Split(strHtml, id)

	if len(arr) > 1 {
		str := arr[1]
		arr = strings.Split(str, "<div")
		if len(arr) > 1 {
			str = arr[0]
		}
		return str
	}
	return strHtml
}

func getHtml(url string) string {
	resp, err := http.Get(url)
	if err != nil {
		return "err"
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "geting err"
	}
	b, _ := ioutil.ReadAll(resp.Body)
	if dic["charset"] == "UTF-8" {
		return string(b)
	} else {
		c, _ := gbkToUtf8(b)
		return string(c)
	}
}

func gbkToUtf8(b []byte) ([]byte, error) {
	reader := transform.NewReader(bytes.NewReader(b), simplifiedchinese.GBK.NewDecoder())
	d, e := ioutil.ReadAll(reader)
	if e != nil {
		return nil, e
	}
	return d, nil
}

func exitHtml() string {
	return "<button type=\"button\" onclick=\"closeGo()\">Exit</button><br>\n"
	//	return "<a href=\"javascript:void(0);\" onclick=\"closeGo()\">Exit</a><br><br>"
}

// cmd.exe /c TaskKill /PID 123
func jsExecCmd(strCmd string) string {
	var str string

	str = fmt.Sprintf("<script type=\"text/javascript\">\n")
	str = fmt.Sprintf("%sfunction closeGo() {\n", str)
	str = fmt.Sprintf("%s    alert(\"test1\");\n", str)
	str = fmt.Sprintf("%s    var cmd = new ActiveXObject(\"WScript.Shell\");\n", str)
	str = fmt.Sprintf("%s    cmd.run(\"%s\");\n", str, strCmd)
	//	str = fmt.Sprintf("%s    cmd = null;\n", str)
	str = fmt.Sprintf("%s    alert(\"test2\");\n", str)
	str = fmt.Sprintf("%s}\n", str)
	str = fmt.Sprintf("%s</script>\n", str)
	return str
}
