package gamerec

import (
	"reflect"

	. "util/thrift"
	. "util/thrift/new_example/gen-go/gamerec"
)

var (
	pool *Pool
	cfg  *ClientConfig
)

func init() {
	cfg = NewClientConfig()
	cfg.ServicePath = "/services/com.xiaomi.misearch.apprec.thrift.GameRec"
	cfg.NewThriftClient = reflect.ValueOf(NewGameRecClientFactory)

	cfg.DeployEnv = NewDeployEnv()
	cfg.DeployEnv.Set("default", "staging")
	cfg.DeployEnv.Set("staging", "staging")
	cfg.DeployEnv.Set("c3", "c3")
}

//
func (this *Wrap) Service() GameRec {
	return this.ThriftClient.(GameRec)
}

//no need change below
type Wrap struct {
	*ClientProxy
}

// 包第一次加载是初始化一个全局的pool
func InitPool(envString string, initialCap, maxCap int) {
	pool = NewPool(envString, initialCap, maxCap, cfg)
}

// GameRecClient
func GetPoolClient() (*Wrap, error) {
	c, e := pool.GetClient()
	if e != nil {
		return nil, e
	}

	return &Wrap{c}, nil
}
