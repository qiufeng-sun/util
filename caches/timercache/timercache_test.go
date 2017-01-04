package timercache

import (
	"testing"

	"fmt"
	"time"
)

//
type cacheString struct {
}

func (this *cacheString) LogName() string {
	return "test cache string"
}

var dataNum = 0
var dataStr = "cache data "
var refreshTime = time.Millisecond

func getCacheData() string {
	return fmt.Sprintf("%v%v", dataStr, dataNum)
}

func (this *cacheString) LoadData() (interface{}, error) {
	dataNum++
	d := getCacheData()
	fmt.Println(d)
	return d, nil
}

func (this *cacheString) GetRefreshTime() time.Duration {
	return refreshTime
}

//
func TestGetCache(t *testing.T) {
	caches := NewAtomicCache()

	key := "string"
	caches.Register(key, &cacheString{})
	caches.Start()

	time.Sleep(refreshTime / 2)
	for i := 0; i < 5; i++ {
		d := getCacheData()
		d1 := caches.GetCache(key)
		fmt.Printf("cur: %v, cache:%v\n", d, d1)

		if d != d1 {
			t.Errorf("cur: %v, cache:%v\n", d, d1)
		}

		time.Sleep(refreshTime)
	}
}

// benchmark
// cpu: go test -v -run=^$ -bench=^BenchmarkGetCache$ -benchtime=2s -cpuprofile=prof.cpu
// 		go tool pprof memcache.test prof.cpu
// mem: go test -v -run=^$ -bench=^BenchmarkGetCache$ -benchtime=3s -memprofile=prof.mem
//      go tool pprof -alloc_space memcache.test prof.mem
func BenchmarkGetCache(b *testing.B) {
	b.ReportAllocs()

	caches := NewAtomicCache()

	key := "string"
	caches.Register(key, &cacheString{})
	caches.Start()

	for i := 0; i < b.N; i++ {
		caches.GetCache(key)
	}
}

// parallel: go test -bench=BenchmarkGetCacheParallel -blockprofile=proc.block
//           go tool pprof memcache.test proc.block
func BenchmarkGetCacheParallel(b *testing.B) {
	caches := NewAtomicCache()

	key := "string"
	caches.Register(key, &cacheString{})
	caches.Start()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			caches.GetCache(key)
		}
	})
}
