package net

import (
	"io/ioutil"
	"net/http"

	"util/logs"
)

////////////////////////////////////////////////////////////////////////////////
//
func HttpGet(url string) ([]byte, error) {
	//
	logs.Debug("url:%v\n", url)

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
