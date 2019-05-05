package tomysql

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"util/dbs/redis/redigo"
)

//
var g_samplePools *redigo.RedisPools

//
func init() {
	Init("test.conf")
	redigo.InitByFile("../redis.conf")

	g_samplePools = redigo.GetRedisPools("sample")
}

//
func getSampleConn() *redigo.RedisConn {
	return g_samplePools.GetConn()
}

//
func TestUpdateHashSet(t *testing.T) {
	c := getSampleConn()
	defer c.Close()

	num := 10
	for i := 0; i < num; i++ {
		key := fmt.Sprintf("h:test:%v", i)

		c.Do("hset", key, "name", i)
		c.Do("hset", key, "sex", rand.Intn(1))
		c.Do("hset", key, "score", rand.Intn(100))
		c.Do("hset", key, "timestamp", time.Now().Unix())

		UpdateHashSet(c, key)
	}
	time.Sleep(time.Second)
}

//
func TestUpdateList(t *testing.T) {
	c := getSampleConn()
	defer c.Close()

	num := 5
	for i := 0; i < num; i++ {
		key := fmt.Sprintf("l:test:%v", i)
		c.Do("rpush", key, i)
		c.Do("rpush", key, rand.Intn(10))
		c.Do("rpush", key, 3)
		c.Do("rpush", key, time.Now().Unix())

		UpdateList(c, key)
	}
	time.Sleep(time.Second)
}
