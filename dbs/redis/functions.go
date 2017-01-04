package redis

import (
	"util/logs"

	ghRedis "github.com/garyburd/redigo/redis"
)

// locker
func Lock(key string, lockSec int, conn *RedisConn) bool {
	logs.DebugFunc()

	//
	val := 1

	reply, e := ghRedis.Int(conn.Do("SETNX", key, val))
	if e != nil {
		logs.Error("SETNX failed! error=%v\n", e)
		return false
	}

	// lock success
	if 1 == reply {
		conn.Do("EXPIRE", key, lockSec)
		return true
	}

	// lock failed
	reply, e = ghRedis.Int(conn.Do("TTL", key))
	if e != nil {
		logs.Error("TTL failed! error=%v\n", e)
		return false
	}

	if -1 == reply {
		conn.Do("EXPIRE", key, lockSec)
	}

	return false
}

func UnLock(key string, conn *RedisConn) {
	conn.Do("DEL", key)
}
