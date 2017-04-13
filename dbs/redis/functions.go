package redis

import (
	"util/logs"
)

//
var g_val = 1

// locker
func Lock(key string, lockSec int, conn *RedisConn) bool {
	// SET key value [EX seconds] [PX milliseconds] [NX|XX]
	r, e := conn.Do("SET", key, g_val, "EX", lockSec, "NX")
	if e != nil || nil == r {
		logs.Info("SET NX expire failed! reply=%v, error=%v\n", r, e)
		return false
	}

	return true
}

func Unlock(key string, conn *RedisConn) {
	conn.Do("DEL", key)
}
