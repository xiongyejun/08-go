// 获取天气
package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"pkgMySelf/colorPrint"
	"regexp"
)

const WEATHER_URL string = "http://www.weather.com.cn/weather/101240101.shtml"

type info struct {
	date            string
	week            string
	wea             string
	highTemperature string // 高温
	lowTemperature  string
}

func main() {
	cd := colorPrint.NewColorDll()
	rsp, err := http.Get(WEATHER_URL)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer rsp.Body.Close()
	b, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	result, err := getInfo(string(b))
	if err != nil {
		fmt.Println(err)
		return
	}
	cd.SetColor(colorPrint.White, colorPrint.DarkMagenta)
	for _, v := range result {
		fmt.Printf("%3s%s%8s\t%2s-%2s℃<\n", v.date, v.week, v.wea, v.lowTemperature, v.highTemperature)
	}
	cd.UnSetColor()
}

// 用正则获取需要的天气数据
func getInfo(strHtml string) (result []info, err error) {
	reg, err := regexp.Compile(`<h1>(\d{1,2}日)（(.{2})）</h1>(?s:.*?)class="wea">(.*?)</p>(?s:.*?)<span>(.*?)℃</span>/<i>(.*?)℃</i>`) // 不用\d是怕有负号
	if err != nil {
		fmt.Println(err)
		return
	}
	regResult := reg.FindAllStringSubmatch(strHtml, -1)
	result = make([]info, len(regResult))
	for i := 0; i < len(regResult); i++ {
		j := 1
		result[i].date = regResult[i][j]
		j++
		result[i].week = regResult[i][j]
		j++
		result[i].wea = regResult[i][j]
		j++
		result[i].highTemperature = regResult[i][j]
		j++
		result[i].lowTemperature = regResult[i][j]
	}
	return
}
