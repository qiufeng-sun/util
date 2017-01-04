package etcd

import (
	"context"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/coreos/etcd/client"

	"util/logs"
)

//
const x_requestTimeout = time.Second * 2

//
type ServerInfo struct {
	Addr  string
	Value string
}

//
type EtcdClient struct {
	// config
	EtcdAddrs  []string // etcd server address
	ServerPath string   // like "/game"
	Watch      chan bool

	// attr
	KeysAPI     client.KeysAPI
	ServerInfos []*ServerInfo // server infos
	LastIndex   uint64        // last modify index
	sync.Mutex                // mutex for this.ServerAddrs op
}

//
func NewEtcdClient(etcdAddrs []string, serverPath string) *EtcdClient {
	// check server path to start with "/"
	if !strings.HasPrefix(serverPath, "/") {
		serverPath = "/" + serverPath
	}

	//
	c := &EtcdClient{
		EtcdAddrs:  etcdAddrs,
		ServerPath: serverPath,
		Watch:      make(chan bool, 1),
	}

	//
	c.start()

	return c
}

//
func (this *EtcdClient) start() {
	//
	config := client.Config{
		Endpoints:               this.EtcdAddrs,
		Transport:               client.DefaultTransport,
		HeaderTimeoutPerRequest: x_requestTimeout,
	}

	//
	c, e := client.New(config)
	if e != nil {
		panic("etcd:" + e.Error())
	}

	this.KeysAPI = client.NewKeysAPI(c)

	// get server address
	if e := this.pullServers(); e != nil {
		// log
		logs.Warn("etcd start:get server address failed! error:%v", e)
	}

	//
	this.Watch <- true
	go this.watch()
}

//
func (this *EtcdClient) GetServerInfos() []*ServerInfo {
	this.Lock()
	ret := this.ServerInfos
	this.Unlock()

	return ret
}

//
func (this *EtcdClient) setServerInfos(infos []*ServerInfo) {
	this.Lock()
	this.ServerInfos = infos
	this.Unlock()
}

//
func (this *EtcdClient) pullServers() error {
	//
	option := &client.GetOptions{Recursive: true}

	//
	resp, e := this.KeysAPI.Get(context.Background(), this.ServerPath, option)
	if e != nil {
		return e
	}

	// log
	logs.Info("pullServers<etcd>: path=%v, resp=%v", this.ServerPath, resp)

	//
	this.LastIndex = resp.Index

	// none server
	if nil == resp.Node || len(resp.Node.Nodes) <= 0 {
		this.ServerInfos = nil
		return nil
	}

	// log nodes
	logs.Info("pullServers<etcd>: path=%v, nodes=%v", this.ServerPath, resp.Node.Nodes)

	// proc server address
	num := len(resp.Node.Nodes)
	infos := make([]*ServerInfo, num)
	for i, node := range resp.Node.Nodes {
		key, _ := url.QueryUnescape(node.Key)
		info := &ServerInfo{
			Addr:  strings.TrimPrefix(key, this.ServerPath+"/"),
			Value: node.Value,
		}
		infos[i] = info
	}
	this.setServerInfos(infos)

	return nil
}

//
func (this *EtcdClient) watch() {
	for {
		//
		opts := &client.WatcherOptions{
			AfterIndex: this.LastIndex,
			Recursive:  true,
		}

		// rebuild watcher to use the lastest modified index which last got
		w := this.KeysAPI.Watcher(this.ServerPath, opts)

		// watch
		if _, e := w.Next(context.Background()); e != nil {
			logs.Warn("watch<etcd>: watch failed! error:%v", e)
			continue
		}

		// pull server address
		if e := this.pullServers(); e != nil {
			logs.Warn("watch<etcd>: pull server address failed! error:%v", e)
			continue
		}

		// active watch
		select {
		case this.Watch <- true: // do nothing
		default: // do not get yet or not matter about it
		}
	}
}

//
func AddWatchs(etcdAddrs, services []string, fUpdate func(svc string, svcAddrs []string)) {
	logs.Info("etcd add watchs: etcdAdrs=%v, services=%v", etcdAddrs, services)

	for _, v := range services {
		go func(svc string) {
			c := NewEtcdClient(etcdAddrs, svc)
			for range c.Watch {
				infos := c.GetServerInfos()
				addrs := make([]string, len(infos))
				for i, info := range infos {
					addrs[i] = info.Addr
				}
				fUpdate(svc, addrs)
			}
		}(v)
	}
}
