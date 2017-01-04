package net

import (
	"net"
	"strings"

	"util/logs"
)

var g_localIp string = "0.0.0.0"

func LocalIp() string {
	return g_localIp
}

func initLocalIp() {
	conn, e := net.Dial("udp", "www.mi.com:80")
	if e != nil {
		logs.Warn("error:%v", e)
		return
	}
	defer conn.Close()

	g_localIp = strings.Split(conn.LocalAddr().String(), ":")[0]

	logs.Info("local ip:%v", LocalIp())
}

func init() {
	initLocalIp()
}
