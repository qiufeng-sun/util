package etcd

import (
	"context"
	"net/url"
	"strings"
	"time"

	"github.com/coreos/etcd/client"

	"util/logs"
)

//
const x_updateInterval = time.Second
const x_serverNodeTTL = x_updateInterval + time.Second
const x_registerTimeout = time.Second * 2

//
type EtcdRegister struct {
	EtcdAddrs      []string      // etcd server address
	ServerPath     string        // like "/game"
	ServerAddr     string        // like "xx.xx.xxx.x:1234"
	UpdateInterval time.Duration // update register interval
	ServerNodeTTL  time.Duration // server node ttl
	ServerValue    string

	KeysAPI      client.KeysAPI
	ticker       *time.Ticker // update ticker
	chUpdateStop chan bool    // update stop sync

	debug bool // test flag
}

//
func NewEtcdRegister(etcdAddrs []string,
	updateInterval, serverNodeTTL time.Duration) *EtcdRegister {

	r := &EtcdRegister{
		EtcdAddrs:      etcdAddrs,
		UpdateInterval: x_updateInterval,
		ServerNodeTTL:  x_serverNodeTTL,
	}

	if updateInterval > 0 {
		r.UpdateInterval = updateInterval
	}

	if serverNodeTTL >= time.Second {
		r.ServerNodeTTL = serverNodeTTL
	}

	if r.ServerNodeTTL < r.UpdateInterval {
		panic("etcd register: invalid update interval or server node ttl")
	}

	//
	r.start()

	return r
}

//
func (this *EtcdRegister) start() {
	//
	config := client.Config{
		Endpoints:               this.EtcdAddrs,
		Transport:               client.DefaultTransport,
		HeaderTimeoutPerRequest: x_registerTimeout,
	}

	//
	c, e := client.New(config)
	if e != nil {
		// log and exit
		panic("etcd register:" + e.Error())
	}

	this.KeysAPI = client.NewKeysAPI(c)
}

//
func (this *EtcdRegister) Register(serverPath, serverAddr, value string) {
	// check server path to start with "/"
	if !strings.HasPrefix(serverPath, "/") {
		serverPath = "/" + serverPath
	}

	serverAddr = url.QueryEscape(serverAddr)

	//
	this.ServerPath = serverPath
	this.ServerAddr = serverAddr
	this.ServerValue = value

	//
	this.register()

	//
	this.chUpdateStop = make(chan bool, 1)
	this.ticker = time.NewTicker(this.UpdateInterval)

	// update ttl
	go this.update()
}

func (this *EtcdRegister) getRegisterPath() string {
	return this.ServerPath + "/" + this.ServerAddr
}

func (this *EtcdRegister) register() {
	//
	path := this.getRegisterPath()
	option := &client.SetOptions{
		PrevExist: client.PrevIgnore,
		TTL:       this.ServerNodeTTL,
	}

	// reg
	_, e := this.KeysAPI.Set(context.Background(), path, this.ServerValue, option)
	if e != nil {
		logs.Warn("etcd register: register server failed! path:%v, error:%v", path, e)
	} else {
		logs.Info("etcd register: ok! path: %v", path)
	}
}

//
func (this *EtcdRegister) Unregister() {
	//
	this.stopUpdate()

	//
	path := this.ServerPath + "/" + this.ServerAddr

	//
	this.KeysAPI.Delete(context.Background(), path, nil)
}

//
func (this *EtcdRegister) update() {
	// Refresh set to true, to update ttl and without firing a watch
	option := &client.SetOptions{TTL: this.ServerNodeTTL, Refresh: true}

	// register path
	path := this.getRegisterPath()

	//
	for {
		select {
		case <-this.ticker.C:
			_, e := this.KeysAPI.Set(context.Background(), path, "", option)
			if e != nil {
				logs.Warn("etcd update: set ttl failed! path: %v, error:%v!", path, e)

				this.register()
			}

			if this.debug {
				logs.Debug("update<%v>: error=%v", path, e)
			}
		case <-this.chUpdateStop:
			break
		}
	}

	this.chUpdateStop <- true

	logs.Info("etcd update stop! path: %v", path)
}

//
func (this *EtcdRegister) stopUpdate() {
	// stop update ticker
	this.ticker.Stop()

	// stop update goroutine
	this.chUpdateStop <- true

	// wait update goroutine stop
	<-this.chUpdateStop
}
