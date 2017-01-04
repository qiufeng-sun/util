// 有有效期的cache:每条记录有单独到有效期, 获取记录时判断是否有效, 失效则重新获取
package tempcache

import (
	"fmt"
	"sync"
	"time"

	"util/logs"
)

////////////////////////////////////////////////////////////////////////////////
//
var g_debug = false

//
type Cacher interface {
	Load(keys []string) ([]interface{}, error)
	GetExpiredSec() int64
	LogName() string
}

////////////////////////////////////////////////////////////////////////////////
//
type TempValue struct {
	ExpiredTime int64
	Value       interface{}
}

//
type Cache struct {
	*sync.RWMutex
	Cacher

	Data map[string]*TempValue
}

//
func NewCache(cacher Cacher) *Cache {
	return &Cache{
		RWMutex: &sync.RWMutex{},
		Cacher:  cacher,
		Data:    make(map[string]*TempValue),
	}
}

//
func (this *Cache) GetValues(keys []string) []interface{} {
	//
	num := len(keys)
	if num <= 0 {
		return nil
	}

	//
	ret := make([]interface{}, num)
	uKeys := make([]string, num)
	uKeyIndexs := make([]int, num)

	//
	ds := this.Data
	uIndex := 0
	now := time.Now().Unix()

	// from cache 1st
	this.RLock()
	for i := 0; i < num; i++ {
		k := keys[i]
		if "" == k {
			continue
		}
		v, ok := ds[k]
		if ok && v.ExpiredTime > now {
			ret[i] = v.Value

			// for unit test
			if g_debug {
				logs.Info("found one")
			}
		} else {
			uKeys[uIndex] = k
			uKeyIndexs[uIndex] = i
			uIndex++
		}
	}
	this.RUnlock()

	// found all
	if uIndex <= 0 {
		return ret
	}

	// load unfound 2nd
	uKeys = uKeys[:uIndex]
	res, e := this.Load(uKeys)
	if e != nil {
		logs.Warn("load %v failed! error=%v", this.LogName(), e)
	} else {
		expiredSec := this.GetExpiredSec() + now
		this.Lock()
		for i, v := range res {
			if nil == v {
				continue
			}

			k := uKeys[i]
			ds[k] = &TempValue{ExpiredTime: expiredSec, Value: v}

			index := uKeyIndexs[i]
			ret[index] = v
		}
		this.Unlock()
	}

	return ret
}

////////////////////////////////////////////////////////////////////////////////
// cache mgr
type CacheMgr map[string]*Cache

//
func NewCacheMgr() *CacheMgr {
	mgr := CacheMgr(make(map[string]*Cache))
	return &mgr
}

//
func (this *CacheMgr) Register(key string, cacher Cacher) {
	if _, ok := (*this)[key]; ok {
		panic(fmt.Sprintf("duplicate register cache! key=%v\n", key))
	}

	(*this)[key] = NewCache(cacher)
}

//
func (this *CacheMgr) GetCache(key string) *Cache {
	return (*this)[key]
}

////////////////////////////////////////////////////////////////////////////////
// default cache mgr
var g_defMgr = NewCacheMgr()

//
func Register(key string, c Cacher) {
	g_defMgr.Register(key, c)
}

//
func GetCache(key string) *Cache {
	return g_defMgr.GetCache(key)
}

//
func Start() {
	// do nothing
}

////////////////////////////////////////////////////////////////////////////////
