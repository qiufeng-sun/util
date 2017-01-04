// 读取并监听zk节点子节点中数据
package zk

import (
	"time"

	"github.com/samuel/go-zookeeper/zk"

	"util/logs"
)

////////////////////////////////////////////////////////////////////////////////
//
type NodeParser interface {
	Parser      // 解析并保存数据
	ResetData() // 清理数据
}

//
type nodeParserMgr map[string]NodeParser

//
var (
	g_zkNodeParsers     nodeParserMgr = make(nodeParserMgr)
	g_zkChildrenParsers nodeParserMgr = make(nodeParserMgr)
)

//
func AddZkNodeParser(path string, parser NodeParser) {
	if _, ok := g_zkNodeParsers[path]; ok {
		panic("duplicate zk node parser! path=" + path)
	}

	g_zkNodeParsers[path] = parser
}

////////////////////////////////////////////////////////////////////////////////
//
var g_zkNodeCon *zk.Conn
var g_zkNodeEvt <-chan zk.Event

//
func LoadZkNodesW(servers []string, recvTimeout time.Duration) {
	zkCon, zkEvt, err := zk.Connect(servers, recvTimeout)
	if err != nil {
		panic(err)
	}

	g_zkNodeCon = zkCon
	g_zkNodeEvt = zkEvt

	loadZkNodes(true)

	go onZkNodeEvent()
}

//
func loadZkNodes(first bool) {
	g_zkChildrenParsers = make(nodeParserMgr)

	for path, parser := range g_zkNodeParsers {
		loadZkNode(path, parser, first)
	}
}

func loadZkNode(path string, parser NodeParser, first bool) {
	nodes, stat, _, e := g_zkNodeCon.ChildrenW(path)
	if !checkError(path, e, first) {
		return
	}
	logs.Info("zk node: %+v\nstat: %+v\n", nodes, stat)

	parser.ResetData()
	for _, node := range nodes {
		pn := path + "/" + node

		g_zkChildrenParsers[pn] = parser
		loadZkNodeChildData(pn, parser, first)
	}
}

func onZkNodeEvent() {
	for evt := range g_zkNodeEvt {
		if evt.Err != nil {
			continue
		}

		path := evt.Path

		logs.Info("zk node event: %+v\n", evt)
		switch evt.Type {
		case zk.EventSession:
			if zk.StateHasSession == evt.State {
				loadZkNodes(false)
			}

		case zk.EventNodeDataChanged:
			if parser, ok := g_zkChildrenParsers[path]; ok {
				loadZkNodeChildData(path, parser, false)
			}

		case zk.EventNodeChildrenChanged:
			if parser, ok := g_zkNodeParsers[path]; ok {
				loadZkNode(path, parser, false)
			}
		}
	}
}

func loadZkNodeChildData(path string, parse Parser, first bool) {
	loadZkDataEx(g_zkNodeCon, path, parse, first)
}
