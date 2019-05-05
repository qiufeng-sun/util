// redis常用接口封装
package redix

import (
	"util/logs"
)

// locker
func Lock(pool Pool, key string, lockSec int) bool {
	// key
	key = genLockKey(key)

	// SET key value [EX seconds] [PX milliseconds] [NX|XX]
	_, e := pool.Cmd("SET", key, 1, "EX", lockSec, "NX").Str()
	if e != nil {
		logs.Info("SET NX expire failed! key=%v, error=%v", key, e)
		return false
	}

	return true
}

func Unlock(pool Pool, key string) {
	// key
	key = genLockKey(key)

	pool.Cmd("DEL", key)
}

func genLockKey(key string) string {
	return "lock@" + key
}
