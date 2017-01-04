package gamerec

import (
	//	"errors"
	"fmt"
	"time"

	"git.apache.org/thrift.git/lib/go/thrift"

	. "util/thrift"
	. "util/thrift/thrift_example/gen-go/gamerec"

	"util/thrift/hyzpool"
)

/**
* 需更新的代码
*   ServicePath  = "/services/com.xiaomi.misearch.apprec.thrift.GameRec"
*   ServiceClientProxy.*GameRecClient
*   NewServiceClientProxy.client := NewGameRecClientFactory(s.Transport, protocolFactory)
**/

var (
	env map[string]string
)

func init() {
	env = make(map[string]string)
	// staging环境
	env["staging"] = "staging"
	// 线上环境，此服务当前只在c3部署
	env["lugu"] = "c3"
	env["c3"] = "c3"
	// 默认
	env["default"] = "staging"
}

func myEnv(envString string) string {
	r, ok := env[envString]
	if ok {
		return r
	}
	return env["default"]
}

const (
	ServicePath  = "/services/com.xiaomi.misearch.apprec.thrift.GameRec"
	ServiceLevel = 10
)

var (
	ServiceTimeout = time.Millisecond * 500 // 500ms
)

// 1. thrift service client
type ServiceClientProxy struct {
	*GameRecClient
	*ThriftService
}

func NewServiceClientProxy(server string) (*ServiceClientProxy, error) {
	startTime := CurrentTimeMillis()
	s, err := NewThriftServiceWithTimeout(server, ServiceTimeout)
	if err != nil {
		return nil, err
	}

	protocolFactory := thrift.NewTBinaryProtocolFactoryDefault()

	client := NewGameRecClientFactory(s.Transport, protocolFactory)

	service := &ServiceClientProxy{client, s}

	endTime := CurrentTimeMillis()
	if debug {
		fmt.Println("new GameRecClient wast time:", (endTime - startTime), "ms")
	}
	return service, nil
}

// 2. connection pool for package
var (
	pool  hyzpool.Pool
	zk    *hyzpool.ZooData
	debug = true

//	ErrZkPool   = errors.New("no pool in zk")
//	ErrInitPool = errors.New("pool is not initiallized")
)

// 包第一次加载是初始化一个全局的pool
func InitPool(envString string, initialCap, maxCap int) {
	environment := myEnv(envString)
	zk = hyzpool.NewZkData(environment, ServicePath, ServiceLevel)
	if debug {
		fmt.Println("InitPool", environment, zk)
	}

	factory := func() (hyzpool.Client, error) {
		node := zk.GetThriftNode()
		if node != "" {
			return NewServiceClientProxy(node)
		}
		return nil, ErrZkPool
	}
	pool, _ = hyzpool.NewPool(initialCap, maxCap, factory)
}

func GetPoolClient() (*wrappedClientProxy, error) {
	if pool == nil {
		fmt.Println("pool is not initiallized")
		return nil, ErrInitPool
	}
	client, err := pool.GetByTime(zk.GetUpdateTime())
	if err != nil {
		return nil, err
	}
	c, _ := client.(*ServiceClientProxy)
	return &wrappedClientProxy{c, pool}, nil
}

//wrapped struct, add method: Return (return the client back to pool)
type wrappedClientProxy struct {
	*ServiceClientProxy
	p hyzpool.Pool
}

func (c *wrappedClientProxy) Return() error {
	if c.ServiceClientProxy.IsClose() {
		return nil
	}
	return c.p.Put(c.ServiceClientProxy)
}
