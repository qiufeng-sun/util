package redis

import (
	"testing"

	"math/rand"
	"time"

	ghRedis "github.com/garyburd/redigo/redis"
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

//
func TestHealthCheck(t *testing.T) {
	e := HealthCheck()
	t.Logf("health check: error=%v\n", e)

	if e != nil {
		t.Error(e)
	}
}

//
func TestRedis(t *testing.T) {
	c := getSampleConn("")
	defer c.Close()

	k := "testing:unit"
	v := rand.Int()

	t.Logf("input: k=%v, v=%v\n", k, v)

	r, e := c.Do("SET", k, v)
	t.Logf("set res: k=%v, v=%v, reply=%v, error=%v\n", k, v, r, e)

	if e != nil {
		t.Error(e)
		return
	}

	r1, e := ghRedis.Int(c.Do("GET", k))
	t.Logf("get res: k=%v, v=%v, error=%v\n", k, r1, e)

	if e != nil {
		t.Error(e)
		return
	}

	if r1 != v {
		t.Errorf("test failed! k=%v, res=%v, expected=%v\n", k, r1, v)
		return
	}
}
