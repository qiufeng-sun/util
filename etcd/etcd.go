package etcd

import (
	"time"

	"util/logs"
)

var _ = logs.Debug

//
type SrvCfg struct {
	EtcdAddrs []string // etcd地址数组

	SrvAddr      string // 服务监听addr
	SrvRegPath   string // etcd上服务位置
	SrvRegUpTick int64  // etcd上服务更新ttl间隔(ms)

	WatchPaths []string // 需监听的服务器位置
}

//
func RegAndWatchs(logMsg string, cfg *SrvCfg, fUpdate func(svc string, svcAddrs []string)) {
	logs.Info("etcd init start: %v", logMsg)
	defer logs.Info("etcd init end: %v", logMsg)

	// register
	Register(cfg)

	// services watched
	if len(cfg.WatchPaths) > 0 {
		AddWatchs(cfg.EtcdAddrs, cfg.WatchPaths, fUpdate)
	}
}

//
func Register(cfg *SrvCfg) {
	//
	etcdRegUpTick := time.Duration(cfg.SrvRegUpTick) * time.Millisecond
	etcdRegTTL := etcdRegUpTick * 3

	//
	r := NewEtcdRegister(cfg.EtcdAddrs, etcdRegUpTick, etcdRegTTL)
	r.Register(cfg.SrvRegPath, cfg.SrvAddr, "")
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
