package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

func (me *dataStructItem) getValue(strPre string) (err error) {
	resp := new(http.Response)
	if resp, err = http.Get(me.url); err != nil {
		return
	}
	defer resp.Body.Close()

	var b []byte
	if b, err = ioutil.ReadAll(resp.Body); err != nil {
		return
	}
	// jsonString={"drawyears":[{"year":"2018"},{"year":"2017"},{"year":"2016"},{"year":"2015"},{"year":"2014"},{"year":"2013"},{"year":"2012"},{"year":"2011"},{"year":"2010"},{"year":"2009"}]}
	// 获取=后面的内容
	index := bytes.Index(b, []byte("="))
	b = b[index+1:]

	// 解析
	var d map[string]interface{}
	if err = json.Unmarshal(b, &d); err != nil {
		return
	}
	if v, ok := d[me.attr]; ok {
		arr := v.([]interface{})
		for _, item := range arr {
			itemd := item.(map[string]interface{})
			if y, ok := itemd[me.attrItem]; ok {
				me.rets = append(me.rets, strPre+y.(string))
			}
		}
	}

	return nil
}

type drawResult struct {
	Competitions string `json:"Competitions"`
	Datetime     string `json:"datetime"`
	Dispatchamt  string `json:"dispatchamt"`
	Guestteam    string `json:"guestteam"`
	Handicap     string `json:"handicap"`
	Hitcount     string `json:"hitcount"`
	Hostteam     string `json:"hostteam"`
	Matchno      string `json:"matchno"`
	Result       string `json:"result"`
	Score        string `json:"score"`
	Spvalue      string `json:"spvalue"`
	Stake        string `json:"stake"`
}

type datas struct {
	Datetime   string        `json:"datetime"`
	Drawresult []*drawResult `json:"drawresult"`
}

func (me *datas) getData(url string) (err error) {
	//	println(url)

	resp := new(http.Response)
	if resp, err = http.Get(url); err != nil {
		return
	}
	defer resp.Body.Close()

	var b []byte
	if b, err = ioutil.ReadAll(resp.Body); err != nil {
		return
	}
	// 获取=后面的内容
	index := bytes.Index(b, []byte("="))
	b = b[index+1:]

	if err = json.Unmarshal(b, me); err != nil {
		return
	}

	return nil
}
