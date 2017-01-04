package etcd

import (
	"time"

	"github.com/astaxie/beego/config"

	"util/logs"
)

var _ = logs.Debug

//
func RegAndWatchs(logMsg string, confd config.Configer, fUpdate func(svc string, svcAddrs []string)) {
	logs.Info("etcd init start: %v", logMsg)
	defer logs.Info("etcd init end: %v", logMsg)

	//
	etcdAddrs := confd.Strings("etcd_addrs")

	// register
	etcdRegAddr := confd.String("server_addr")
	etcdRegPath := confd.String("etcd_reg_path")
	etcdRegUpTick := time.Duration(confd.DefaultInt64("etcd_reg_uptick", 1000)) * time.Millisecond
	etcdRegTTL := etcdRegUpTick * 3

	r := NewEtcdRegister(etcdAddrs, etcdRegUpTick, etcdRegTTL)
	r.Register(etcdRegPath, etcdRegAddr, "")

	// services watched
	etcdWatchPaths := confd.Strings("etcd_watch_path")

	if len(etcdWatchPaths) > 0 {
		AddWatchs(etcdAddrs, etcdWatchPaths, fUpdate)
	}
}
