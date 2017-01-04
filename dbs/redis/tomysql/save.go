package tomysql

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	ghRedis "github.com/garyburd/redigo/redis"
	_ "github.com/go-sql-driver/mysql"

	"util/dbs/database"
	"util/dbs/redis"
	"util/loader"
	"util/logs"
)

var _ = fmt.Print
var _ = logs.Debug

//
type Config struct {
	Redis RedisCfg `json:"redis"`
	Mysql MysqlCfg `json:"mysql"`
}

type RedisCfg struct {
	HashSet  string           `json:"hashset"`
	List     string           `json:"list"`
	MaxCache int              `json:"maxcache"`
	DBCfgs   []redis.RedisCfg `json:"db"`
}

type MysqlCfg struct {
	HashSet string                             `json:"hashset"`
	List    string                             `json:"list"`
	DBCfgs  map[string]database.DatabaseConfig `json:"db"`
}

//
type Saver struct {
	Cfg *Config
	redis.RedisPools
}

func (this *Saver) Init(confName string) {
	// load config
	cfg := &Config{}
	if e := loader.ParseJsonFile(confName, &cfg); e != nil {
		logs.Panicln("error", e)
		return
	}

	logs.Info("tomysql config: %v", cfg)

	this.Cfg = cfg

	// init redis
	this.RedisPools = redis.CreatePools(cfg.Redis.DBCfgs, len(cfg.Redis.DBCfgs))

	// init mysql
	database.InitByCfgs(cfg.Mysql.DBCfgs)

	// loop to mysql
	for i := range this.RedisPools {
		go this.save(i, this.Cfg.Redis.HashSet, this.Cfg.Mysql.HashSet, this.getHashSetVal)
		go this.save(i, this.Cfg.Redis.List, this.Cfg.Mysql.List, this.getListVal)
	}
}

func (this *Saver) save(connIndex int, redisKey, table string,
	fget func(c *redis.RedisConn, key string) (string, error)) {

	var c *redis.RedisConn
	for {
		if c != nil {
			c.Close()
			time.Sleep(time.Second)
		}
		c = this.RedisPools.GetConnByIndex(connIndex)

		r, e := ghRedis.String(c.Do("LPOP", redisKey))
		if ghRedis.ErrNil == e {
			continue
		}

		if e != nil {
			logs.Warn("redis save: get save list failed! key=%v, error=%v", redisKey, e)
			continue
		}

		// get save string
		val, e := fget(c, r)
		if e != nil {
			logs.Warn("redis save: get save value failed! key=%v, error=%v", r, e)
			continue
		}

		// save 2 mysql
		if e := this.tomysql(table, r, val); e != nil {
			logs.Warn("to msyql: failed! table:%v, key:%v, val:%v, error:%v",
				table, r, val, e)
			continue
		}

		c.Close()
		c = nil
	}
}

func (this *Saver) tomysql(table, key, val string) error {
	logs.Info("to mysql: table=%v, key=%v, val=%v", table, key, val)

	//
	sql := fmt.Sprintf("insert into %v set rkey=?,rval=? on duplicate key update rkey=?,rval=?", table)

	o := database.GetDefOrm()
	_, e := o.Raw(sql, key, val, key, val).Exec()

	return e
}

func (this *Saver) getHashSetVal(c *redis.RedisConn, key string) (string, error) {
	r, e := ghRedis.Strings(c.Do("HGETAll", key))
	if ghRedis.ErrNil == e {
		return "", nil
	}
	if e != nil || len(r) <= 0 {
		return "", e
	}

	vals := make([]string, len(r))
	for i, s := range r {
		vals[i] = url.QueryEscape(s)
	}

	return strings.Join(vals, " "), nil
}

func (this *Saver) getListVal(c *redis.RedisConn, key string) (string, error) {
	r, e := ghRedis.Strings(c.Do("LRANGE", key, 0, -1))
	if ghRedis.ErrNil == e {
		return "", nil
	}
	if e != nil || len(r) <= 0 {
		return "", e
	}

	vals := make([]string, len(r))
	for i, s := range r {
		vals[i] = url.QueryEscape(s)
	}

	return strings.Join(vals, " "), nil
}

//
func (this *Saver) UdpateHashSet(c *redis.RedisConn, hsKey string) bool {
	return this.toCache(c, this.Cfg.Redis.HashSet, hsKey)
}

//
func (this *Saver) UpdateList(c *redis.RedisConn, lstKey string) bool {
	return this.toCache(c, this.Cfg.Redis.List, lstKey)
}

func (this *Saver) toCache(c *redis.RedisConn, key, val string) bool {
	// check len
	num, e := ghRedis.Int(c.Do("LLEN", key))
	if e != nil {
		logs.Warn("list cache: llen failed! key:%v, val:%v, error:%v", key, val, e)
		return false
	}
	if num > this.Cfg.Redis.MaxCache {
		logs.Warn("list cache: too many wait! key:%v, val:%v, num:%v", key, val, num)
		return false
	}

	_, e = c.Do("RPUSH", key, val)
	if e != nil {
		logs.Warn("list cache: rpush failed! key:%v, val:%v, error:%v", key, val, e)
		return false
	}
	return true
}

//
var g_saver *Saver

func Init(confName string) {
	g_saver = &Saver{}
	g_saver.Init(confName)
}

func UpdateHashSet(conn *redis.RedisConn, hsKey string) bool {
	return g_saver.UdpateHashSet(conn, hsKey)
}

//
func UpdateList(conn *redis.RedisConn, lstKey string) bool {
	return g_saver.UpdateList(conn, lstKey)
}
