// 读取并监听zk节点中数据
package zk

import (
	"fmt"
	"time"

	"github.com/samuel/go-zookeeper/zk"

	"util/logs"
)

////////////////////////////////////////////////////////////////////////////////
//
type Parser interface {
	Unmarshal(data []byte) (interface{}, error) // 解析并保存数据
}

//
type IZkStat interface {
	Stat(stat *zk.Stat)
}

var g_zkParsers map[string]Parser = make(map[string]Parser)

func AddZkParser(path string, parser Parser) {
	if _, ok := g_zkParsers[path]; ok {
		panic("duplicate zk parser! path=" + path)
	}

	g_zkParsers[path] = parser
}

////////////////////////////////////////////////////////////////////////////////
//
var g_zkCon *zk.Conn
var g_zkEvt <-chan zk.Event

func LoadZkDatasW(servers []string, recvTimeout time.Duration) {
	zkCon, zkEvt, err := zk.Connect(servers, recvTimeout)
	if err != nil {
		panic(err)
	}

	g_zkCon = zkCon
	g_zkEvt = zkEvt

	loadZkDatas(true)

	go onZkEvent()
}

func loadZkDatas(first bool) {
	for path, parser := range g_zkParsers {
		loadZkData(path, parser, first)
	}
}

func loadZkData(path string, parse Parser, first bool) {
	loadZkDataEx(g_zkCon, path, parse, first)
}

func loadZkDataEx(zkCon *zk.Conn, path string, parse Parser, first bool) {
	data, stat, _, err := zkCon.GetW(path)
	if !checkError(path, err, first) {
		return
	}
	logs.Info("zk data: %+v\n %+v\n", string(data), stat)

	val, err := parse.Unmarshal(data)
	if !checkError(path, err, first) {
		return
	}

	if pI, ok := parse.(IZkStat); ok {
		pI.Stat(stat)
	}

	logs.Info("go data: %+v\n", val)
}

func onZkEvent() {
	for evt := range g_zkEvt {
		if evt.Err != nil {
			continue
		}

		logs.Info("zk event: path=%v, type=%v, state=%v\n", evt.Path, evt.Type, evt.State)
		switch evt.Type {
		case zk.EventSession:
			if zk.StateHasSession == evt.State {
				loadZkDatas(false)
			}

		case zk.EventNodeDataChanged:
			if parser, ok := g_zkParsers[evt.Path]; ok {
				loadZkData(evt.Path, parser, false)
			}
		}
	}
}

func checkError(path string, err error, first bool) bool {
	if err != nil {
		msg := fmt.Sprintf("get zk data failed! path=%v, err=%v\n", path, err.Error())
		if first {
			panic(msg)
		} else {
			logs.Warn(msg)
		}

		return false
	}

	return true
}
