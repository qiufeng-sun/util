package redis

import (
	"math/rand"
	"time"
)

//
var g_samplePools *RedisPools

//
func init() {
	InitByFile("redis.conf")

	g_samplePools = GetRedisPools("sample")

	rand.Seed(time.Now().Unix())
}

//
func getSampleConn(uid string) *RedisConn {
	return g_samplePools.GetConn()
	// return g_samplePools.GetConnByRand()
	// return g_samplePools.GetConnByIndex()
}
