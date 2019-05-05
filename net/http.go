package net

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"util/logs"
)

////////////////////////////////////////////////////////////////////////////////
//
func HttpGet(url string) ([]byte, error) {
	//
	//logs.Debug("url:%v\n", url)

	//
	resp, e := http.Get(url)
	if e != nil {
		logs.Error("http get failed! url=%v, error=%v\n", url, e)
		return nil, e
	}
	defer resp.Body.Close()

	data, e := ioutil.ReadAll(resp.Body)
	if e != nil {
		logs.Error("http get body read failed! url=%v, error=%v\n", url, e)
		return nil, e
	}

	return data, nil
}

//
func HttpGetJson(url string, out interface{}) error {
	data, e := HttpGet(url)
	if e != nil {
		return e
	}

	if e := json.Unmarshal(data, out); e != nil {
		logs.Error("HttpGetJson|url:%v, error:%v, data:%v", url, e, string(data))
		return e
	}

	return nil
}
