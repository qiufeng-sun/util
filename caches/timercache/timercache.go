// 定时更新cache: 每类cache有一个更新时间, 到期后, 整体更新该类数据
package timercache

import (
	"fmt"
	"sync"
	"time"

	"util/logs"
)

////////////////////////////////////////////////////////////////////////////////
type Cacher interface {
	LoadData() (interface{}, error)
	GetRefreshTime() time.Duration
	LogName() string
}

////////////////////////////////////////////////////////////////////////////////
//
type CacheData struct {
	Cacher
	//	Key string
	Data interface{}
}

//
type AtomicCache struct {
	Cache map[string]*CacheData
	*sync.RWMutex
}

//
func NewAtomicCache() *AtomicCache {
	return &AtomicCache{Cache: make(map[string]*CacheData), RWMutex: &sync.RWMutex{}}
}

//
func (this *AtomicCache) Register(key string, c Cacher) {
	if _, ok := this.Cache[key]; ok {
		panic(fmt.Sprintf("duplicate register cache! key=%v\n", key))
	}

	this.Cache[key] = &CacheData{Cacher: c, Data: nil}
}

//
func (this *AtomicCache) loadAll() {
	for k, d := range this.Cache {
		c := d.Cacher
		dd, e := c.LoadData()
		if e != nil {
			panic(fmt.Sprintf("%v load failed! error:%v", c.LogName(), e))
		}

		this.Cache[k].Data = dd
	}
}

//
func (this *AtomicCache) GetCache(key string) interface{} {
	this.RLock()
	defer this.RUnlock()

	val, ok := this.Cache[key]
	if !ok || nil == val {
		return nil
	}

	return val.Data
}

//
func (this *AtomicCache) Start() {
	this.loadAll()

	for k, d := range this.Cache {
		go this.updateData(k, d)
	}
}

//
func (this *AtomicCache) updateData(k string, d *CacheData) {
	for {
		time.Sleep(d.GetRefreshTime())

		c := d.Cacher
		dd, e := c.LoadData()
		if e != nil {
			logs.Warn("%v load failed! error:%v", c.LogName(), e)
			continue
		}

		this.Lock()
		this.Cache[k].Data = dd
		this.Unlock()
	}
}

///////////////////////////////////////////////////////////////////////////////
// default cache
var g_defCache = NewAtomicCache()

//
func Register(key string, c Cacher) {
	g_defCache.Register(key, c)
}

//
func GetCache(key string) interface{} {
	return g_defCache.GetCache(key)
}

//
func Start() {
	g_defCache.Start()
}

////////////////////////////////////////////////////////////////////////////////
