package redix

import (
	"fmt"
	"testing"
)

var t_pool Pool

func init() {
	p, e := NewPool("127.0.0.1:6379", "", 10)
	if e != nil {
		fmt.Printf("create pool failed! error=%v", e)
	}
	t_pool = p
}

//
func TestInc(t *testing.T) {
	key := "test:incr"

	r, e := t_pool.Cmd("INCR", key).Int()
	if e != nil || r < 0 {
		t.Fatalf("incr failed! r=%v, error=%v", r, e)
	}
}

//
func TestLock(t *testing.T) {
	key := "test"
	sec := 50

	if ok := Lock(t_pool, key, sec); !ok {
		t.Fatal("lock failed!")
	}

	if ok := Lock(t_pool, key, sec); ok {
		t.Fatal("lock failed!")
	}

	Unlock(t_pool, key)
}
