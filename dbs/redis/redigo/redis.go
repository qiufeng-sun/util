package redigo

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/garyburd/redigo/redis"

	"util"
	"util/loader"
	"util/logs"
)

//
type RedisConn struct {
	redis.Conn
}

//
type RedisPool struct {
	*redis.Pool
}

func (this *RedisPool) Get() *RedisConn {
	return &RedisConn{Conn: this.Pool.Get()}
}

//
func CreateRedisPool(server, password string, maxActive, maxIdle int) *RedisPool {
	ret := &RedisPool{}

	ret.Pool = &redis.Pool{
		MaxActive:   maxActive,
		MaxIdle:     maxIdle,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", server)
			if err != nil {
				return nil, err
			}

			if password != "" {
				if _, err := c.Do("AUTH", password); err != nil {
					c.Close()
					return nil, err
				}
			}

			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}

	return ret
}

// redis config
type RedisCfg struct {
	Addr     string `json:"address"`
	Port     int    `json:"port"`
	Pwd      string `json:"pwd"`
	MaxIdle  int    `json:"maxIdle"`
	MaxTotal int    `json::maxTotal"`
}

func (this RedisCfg) GetHost() string {
	return fmt.Sprintf("%v:%v", this.Addr, this.Port)
}

//
func CreateRedisPoolByCfg(cfg RedisCfg) *RedisPool {
	return CreateRedisPool(cfg.GetHost(), cfg.Pwd, cfg.MaxTotal, cfg.MaxIdle)
}

//
func CreateRedisPoolByFile(cfgFile string) *RedisPool {
	var cfg RedisCfg
	if err := loader.ParseJsonFile(cfgFile, &cfg); err != nil {
		return nil
	}

	logs.Info("file %v parse ok!\n go data:%v\n", cfgFile, util.ToJsonString(cfg))

	return CreateRedisPoolByCfg(cfg)
}

//
type RedisCfgArr struct {
	CfgArr []RedisCfg `json:"cfg"`
	UseNum int        `json:"useNum"`
}

//
type RedisPools []*RedisPool

func (this RedisPools) GetConn() *RedisConn {
	return ([]*RedisPool(this))[0].Get()
}

func (this RedisPools) GetConnByIndex(index int) *RedisConn {
	return ([]*RedisPool(this))[index%len(this)].Get()
}

func (this RedisPools) GetConnByRand() *RedisConn {
	return ([]*RedisPool(this))[rand.Intn(len(this))].Get()
}

func (this *RedisPools) test() (int, error) {
	checkNum := len(*this)

	for i := 0; i < checkNum; i++ {
		conn := this.GetConnByIndex(i)
		defer conn.Close()

		_, e := conn.Do("PING")
		if e != nil {
			return i, e
		}
	}

	return -1, nil
}

func CreatePools(cfgs []RedisCfg, num int) []*RedisPool {
	size := num
	if size > len(cfgs) {
		size = len(cfgs)
	}

	pool := make([]*RedisPool, size)
	for i := 0; i < size; i++ {
		cfg := cfgs[i]
		pool[i] = CreateRedisPool(cfg.GetHost(), cfg.Pwd, cfg.MaxTotal, cfg.MaxIdle)
	}

	return pool
}

//
type RedisMgr struct {
	Cfgs     map[string]RedisCfgArr
	MapPools map[string]RedisPools
}

//
func (this *RedisMgr) InitByFile(fileName string) {
	// load config
	cfg := make(map[string]RedisCfgArr)
	if err := loader.ParseJsonFile(fileName, &cfg); err != nil {
		panic(fmt.Sprintf("Error:%v\n", err.Error()))
	}

	logs.Info("file %v parse ok!\n go data:%v\n", fileName, util.ToJsonString(cfg))

	// init pool
	this.InitByCfg(cfg)
}

func (this *RedisMgr) InitByCfg(cfg map[string]RedisCfgArr) {
	// cfg
	this.Cfgs = cfg
	this.MapPools = make(map[string]RedisPools, len(cfg))

	// init pool
	for k, v := range cfg {
		if v.UseNum <= 0 || v.UseNum > len(v.CfgArr) {
			v.UseNum = len(v.CfgArr)
		}

		if v.UseNum <= 0 {
			panic(fmt.Sprintf("not found redis config! name=%v\n", k))
		}

		if _, ok := this.MapPools[k]; ok {
			panic("redis key already exist! please check config!\n")
		}

		this.MapPools[k] = CreatePools(v.CfgArr, v.UseNum)
	}
}

func (this *RedisMgr) InitByData(data []byte) {
	cfg := make(map[string]RedisCfgArr)
	if err := json.Unmarshal(data, &cfg); err != nil {
		panic(fmt.Sprintf("Error:%v\n, data:%v\n", err.Error(), string(data)))
	}

	logs.Info("redis config parse ok!\n go data:%v\n", util.ToJsonString(cfg))

	this.InitByCfg(cfg)
}

func (this *RedisMgr) GetRedisPools(name string) *RedisPools {
	pools := this.MapPools[name]
	return &pools
}

//
func (this *RedisMgr) HealthCheck() error {
	for k, v := range this.MapPools {
		if index, err := v.test(); err != nil {
			return errors.New(fmt.Sprintf("%v redis ping failed! index=%v, error=%v", k, index, err.Error()))
		}
	}

	return nil
}

//
func NewRedisMgr() *RedisMgr {
	return &RedisMgr{}
}

//
var g_redisMgr = NewRedisMgr()

//
const x_redisCfgFile = "conf/redis.conf"

func Init() {
	g_redisMgr.InitByFile(x_redisCfgFile)
}

func InitByFile(fileName string) {
	g_redisMgr.InitByFile(fileName)
}

func InitByCfg(cfg map[string]RedisCfgArr) {
	g_redisMgr.InitByCfg(cfg)
}

func InitByData(data []byte) {
	g_redisMgr.InitByData(data)
}

func GetRedisPools(name string) *RedisPools {
	return g_redisMgr.GetRedisPools(name)
}

//
func HealthCheck() error {
	return g_redisMgr.HealthCheck()
}
