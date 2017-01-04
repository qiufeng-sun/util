package thrift

import (
	"errors"
	"reflect"
	"time"

	"git.apache.org/thrift.git/lib/go/thrift"

	"util"
	"util/logs"
	"util/thrift/pools"
)

//
type ClientConfig struct {
	ServicePath     string
	NewThriftClient reflect.Value // thrift client function

	DeployEnv

	ServiceLevel   int
	ServiceTimeout time.Duration
	Debug          bool
}

func NewClientConfig() *ClientConfig {
	return &ClientConfig{
		ServiceLevel:   10,
		ServiceTimeout: time.Millisecond * 500,
		Debug:          false,
	}
}

//
type Client struct {
	ThriftClient interface{}
	*ThriftService
}

func NewClient(addr string, cfg *ClientConfig) (*Client, error) {
	defer logs.CheckTime(200, 0, cfg.Debug)()

	s, e := NewThriftServiceWithTimeout(addr, cfg.ServiceTimeout)
	if e != nil {
		return nil, e
	}

	protocolFactory := thrift.NewTBinaryProtocolFactoryDefault()
	client := cfg.NewThriftClient.Call(util.ToValues(s.Transport, protocolFactory))

	return &Client{client[0].Interface(), s}, nil
}

//
var (
	ErrZkPool   = errors.New("no pool in zk")
	ErrInitPool = errors.New("pool is not initiallized")
)

// 2. connection pool for package
type Pool struct {
	pools.Pool
	ZkData *pools.ZooData
}

// 包第一次加载是初始化一个全局的pool
func NewPool(envString string, initialCap, maxCap int, cfg *ClientConfig) *Pool {
	env := cfg.DeployEnv.RealEnv(envString)
	zkData := pools.NewZkData(env, cfg.ServicePath, cfg.ServiceLevel)
	logs.Debugln("InitPool", env, zkData)

	factory := func() (pools.Client, error) {
		node := zkData.GetThriftNode()
		if node != "" {
			return NewClient(node, cfg)
		}
		return nil, ErrZkPool
	}
	pool, _ := pools.NewPool(initialCap, maxCap, factory)

	return &Pool{Pool: pool, ZkData: zkData}
}

func (this *Pool) GetClient() (*ClientProxy, error) {
	v, e := this.Pool.GetByTime(this.ZkData.GetUpdateTime())
	if e != nil {
		return nil, e
	}

	c, _ := v.(*Client)

	return &ClientProxy{Client: c, Pool: this}, nil
}

//
type ClientProxy struct {
	*Client
	*Pool
}

func (this *ClientProxy) Return() error {
	if this.IsClose() {
		return nil
	}
	return this.Put(this.Client)
}
