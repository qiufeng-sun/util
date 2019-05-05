// redis客户端封装(github.com/mediocregopher/radix.v2)
package redix

import (
	"github.com/mediocregopher/radix.v2/cluster"
	"github.com/mediocregopher/radix.v2/pool"
	"github.com/mediocregopher/radix.v2/redis"
)

//
type Pool interface {
	Cmd(cmd string, args ...interface{}) *redis.Resp
}

//
func dial(pwd string) func(network, addr string) (*redis.Client, error) {
	if "" == pwd {
		return redis.Dial
	}

	return func(network, addr string) (*redis.Client, error) {
		c, e := redis.Dial(network, addr)
		if e != nil {
			return nil, e
		}

		e = c.Cmd("AUTH", pwd).Err
		if e != nil {
			return nil, e
		}

		return c, nil
	}
}

// redis pool
func NewPool(addr, pwd string, maxIdle int) (*pool.Pool, error) {
	return pool.NewCustom("tcp", addr, maxIdle, dial(pwd))
}

// cluster
func NewCluster(addr, pwd string, maxIdle int) (*cluster.Cluster, error) {
	opts := cluster.Opts{
		Addr:     addr,
		PoolSize: maxIdle,
		Dialer:   dial(pwd),
	}
	return cluster.NewWithOpts(opts)
}

//
type Config struct {
	Addr    string // 10.38.164.197:6379
	Pwd     string // 密码
	MaxIdle int    // 最大空闲连接
	Cluster bool   // 是否为cluster. true=cluster, false=普通
}

//
func New(cfg *Config) (Pool, error) {
	if cfg.Cluster {
		return NewCluster(cfg.Addr, cfg.Pwd, cfg.MaxIdle)
	}
	return NewPool(cfg.Addr, cfg.Pwd, cfg.MaxIdle)
}
