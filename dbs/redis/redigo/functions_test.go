package redigo

import (
	"testing"
)

//
func TestLock(t *testing.T) {
	c := getSampleConn("")
	defer c.Close()

	key := "test:lock"
	sec := 5

	if ok := Lock(key, sec, c); !ok {
		t.Fatal("lock failed!")
	}

	if ok := Lock(key, sec, c); ok {
		t.Fatal("lock failed!")
	}

	Unlock(key, c)
}
